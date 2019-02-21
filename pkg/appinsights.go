package batchinsights

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
)

type AppInsightsService struct {
	client                   appinsights.TelemetryClient
	aggregateCollectionStart *time.Time
	aggregates               map[string]*appinsights.AggregateMetricTelemetry
}

const AGGREGATE_TIME = time.Duration(1) * time.Minute

func NewAppInsightsService(instrumentationKey string, poolId string, nodeId string) AppInsightsService {
	client := appinsights.NewTelemetryClient(instrumentationKey)
	client.Context().Tags.Cloud().SetRole(poolId)
	client.Context().Tags.Cloud().SetRoleInstance(nodeId)

	return AppInsightsService{
		client:     client,
		aggregates: make(map[string]*appinsights.AggregateMetricTelemetry),
	}
}

func (service AppInsightsService) track(metric *appinsights.MetricTelemetry) {
	t := time.Now()
	fmt.Printf("Last time %v\n", service.aggregateCollectionStart)

	if service.aggregateCollectionStart != nil {
		elapsed := t.Sub(*service.aggregateCollectionStart)
		fmt.Printf("Time elapsed %f > %f\n", elapsed, AGGREGATE_TIME)
		if elapsed > AGGREGATE_TIME {
			fmt.Println("Sending aggregated data")
			for _, aggregate := range service.aggregates {
				service.client.Track(aggregate)
			}
			service.aggregates = make(map[string]*appinsights.AggregateMetricTelemetry)
			service.aggregateCollectionStart = &t
		}
	} else {
		service.aggregateCollectionStart = &t
	}

	id := getMetricId(metric)

	aggregate, ok := service.aggregates[id]
	if !ok {
		aggregate = appinsights.NewAggregateMetricTelemetry(metric.Name)
		aggregate.Properties = metric.Properties
		service.aggregates[id] = aggregate
	}
	aggregate.AddData([]float64{metric.Value})
}

func (service AppInsightsService) UploadStats(stats NodeStats) {
	client := service.client

	for cpuN, percent := range stats.cpuPercents {
		metric := appinsights.NewMetricTelemetry("Cpu usage", percent)
		metric.Properties["CPU #"] = strconv.Itoa(cpuN)
		// metric.Properties["Cores count"] = strconv.Itoa(len(stats.cpuPercents))
		service.track(metric)
	}

	for _, usage := range stats.diskUsage {
		usedMetric := appinsights.NewMetricTelemetry("Disk usage", float64(usage.Used))
		usedMetric.Properties["Disk"] = usage.Path
		service.track(usedMetric)
		freeMetric := appinsights.NewMetricTelemetry("Disk free", float64(usage.Free))
		freeMetric.Properties["Disk"] = usage.Path
		service.track(freeMetric)
	}

	if stats.memory != nil {
		service.track(appinsights.NewMetricTelemetry("Memory used", float64(stats.memory.Used)))
		service.track(appinsights.NewMetricTelemetry("Memory available", float64(stats.memory.Total-stats.memory.Used)))
	}
	if stats.diskIO != nil {
		service.track(appinsights.NewMetricTelemetry("Disk read", float64(stats.diskIO.ReadBps)))
		service.track(appinsights.NewMetricTelemetry("Disk write", float64(stats.diskIO.WriteBps)))
	}

	if stats.netIO != nil {
		service.track(appinsights.NewMetricTelemetry("Network read", float64(stats.netIO.ReadBps)))
		service.track(appinsights.NewMetricTelemetry("Network write", float64(stats.netIO.WriteBps)))
	}

	if len(stats.gpus) > 0 {
		for cpuN, usage := range stats.gpus {
			gpuMetric := appinsights.NewMetricTelemetry("Gpu usage", usage.GPU)
			gpuMetric.Properties["GPU #"] = strconv.Itoa(cpuN)
			service.track(gpuMetric)

			gpuMemoryMetric := appinsights.NewMetricTelemetry("Gpu memory usage", usage.Memory)
			gpuMemoryMetric.Properties["GPU #"] = strconv.Itoa(cpuN)
			service.track(gpuMemoryMetric)
		}
	}

	if len(stats.processes) > 0 {
		for _, processStats := range stats.processes {

			pidStr := strconv.FormatInt(int64(processStats.pid), 10)

			{
				cpuMetric := appinsights.NewMetricTelemetry("Process CPU", processStats.cpu)
				cpuMetric.Properties["Process Name"] = processStats.name
				cpuMetric.Properties["PID"] = pidStr
				service.track(cpuMetric)
			}

			{
				memMetric := appinsights.NewMetricTelemetry("Process Memory", float64(processStats.memory))
				memMetric.Properties["Process Name"] = processStats.name
				memMetric.Properties["PID"] = pidStr
				service.track(memMetric)
			}

		}
	}

	client.Channel().Flush()
}

func getMetricId(metric *appinsights.MetricTelemetry) string {
	groupBy := createKeyValuePairs(metric.Properties)
	return fmt.Sprintf("%s=%f/%s", metric.Name, metric.Value, groupBy)
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
