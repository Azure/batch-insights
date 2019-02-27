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
	aggregation              time.Duration
	aggregateCollectionStart *time.Time
	aggregates               map[string]*appinsights.AggregateMetricTelemetry
}

func NewAppInsightsService(instrumentationKey string, poolId string, nodeId string, aggregation time.Duration) AppInsightsService {
	client := appinsights.NewTelemetryClient(instrumentationKey)
	client.Context().Tags.Cloud().SetRole(poolId)
	client.Context().Tags.Cloud().SetRoleInstance(nodeId)

	return AppInsightsService{
		client:      client,
		aggregation: aggregation,
		aggregates:  make(map[string]*appinsights.AggregateMetricTelemetry),
	}
}

func (service *AppInsightsService) track(metric *appinsights.MetricTelemetry) {
	t := time.Now()

	if service.aggregateCollectionStart != nil {
		elapsed := t.Sub(*service.aggregateCollectionStart)

		if elapsed > AGGREGATE_TIME {
			for k, aggregate := range service.aggregates {
				service.client.Track(aggregate)
			}
			service.aggregates = make(map[string]*appinsights.AggregateMetricTelemetry)
			service.aggregateCollectionStart = &t
		}
	} else {
		service.aggregateCollectionStart = &t
	}

	id := GetMetricID(metric)

	aggregate, ok := service.aggregates[id]
	if !ok {
		aggregate = appinsights.NewAggregateMetricTelemetry(metric.Name)
		aggregate.Properties = metric.Properties
		service.aggregates[id] = aggregate
	}
	aggregate.AddData([]float64{metric.Value})
}

// UploadStats will register the given stats for upload. They will be first aggregated during the given aggregation interval
func (service *AppInsightsService) UploadStats(stats NodeStats) {
	client := service.client

	for cpuN, percent := range stats.CpuPercents {
		metric := appinsights.NewMetricTelemetry("Cpu usage", percent)
		metric.Properties["CPU #"] = strconv.Itoa(cpuN)
		metric.Properties["Core count"] = strconv.Itoa(len(stats.CpuPercents))
		service.track(metric)
	}

	for _, usage := range stats.DiskUsage {
		usedMetric := appinsights.NewMetricTelemetry("Disk usage", float64(usage.Used))
		usedMetric.Properties["Disk"] = usage.Path
		service.track(usedMetric)
		freeMetric := appinsights.NewMetricTelemetry("Disk free", float64(usage.Free))
		freeMetric.Properties["Disk"] = usage.Path
		service.track(freeMetric)
	}

	if stats.Memory != nil {
		service.track(appinsights.NewMetricTelemetry("Memory used", float64(stats.Memory.Used)))
		service.track(appinsights.NewMetricTelemetry("Memory available", float64(stats.Memory.Total-stats.Memory.Used)))
	}
	if stats.DiskIO != nil {
		service.track(appinsights.NewMetricTelemetry("Disk read", float64(stats.DiskIO.ReadBps)))
		service.track(appinsights.NewMetricTelemetry("Disk write", float64(stats.DiskIO.WriteBps)))
	}

	if stats.NetIO != nil {
		service.track(appinsights.NewMetricTelemetry("Network read", float64(stats.NetIO.ReadBps)))
		service.track(appinsights.NewMetricTelemetry("Network write", float64(stats.NetIO.WriteBps)))
	}

	if len(stats.Gpus) > 0 {
		for cpuN, usage := range stats.Gpus {
			gpuMetric := appinsights.NewMetricTelemetry("Gpu usage", usage.GPU)
			gpuMetric.Properties["GPU #"] = strconv.Itoa(cpuN)
			service.track(gpuMetric)

			gpuMemoryMetric := appinsights.NewMetricTelemetry("Gpu memory usage", usage.Memory)
			gpuMemoryMetric.Properties["GPU #"] = strconv.Itoa(cpuN)
			service.track(gpuMemoryMetric)
		}
	}

	if len(stats.Processes) > 0 {
		for _, processStats := range stats.Processes {

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

// GetMetricId compute an group id for this metric so it can be aggregated
func GetMetricID(metric *appinsights.MetricTelemetry) string {
	groupBy := createKeyValuePairs(metric.Properties)
	return fmt.Sprintf("%s/%s", metric.Name, groupBy)
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	first := true
	for key, value := range m {
		if first {
			first = false
		} else {
			fmt.Fprintf(b, ",")
		}
		fmt.Fprintf(b, "%s=%s", key, value)
	}
	return b.String()
}
