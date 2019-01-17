package batchinsights

import (
	"strconv"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
)

type AppInsightsService struct {
	client appinsights.TelemetryClient
}

func NewAppInsightsService(instrumentationKey string, poolId string, nodeId string) AppInsightsService {
	client := appinsights.NewTelemetryClient(instrumentationKey)
	client.Context().Tags.Cloud().SetRole(poolId)
	client.Context().Tags.Cloud().SetRoleInstance(nodeId)

	return AppInsightsService{
		client: client,
	}
}

func (service AppInsightsService) UploadStats(stats NodeStats) {
	client := service.client

	for cpuN, percent := range stats.cpuPercents {
		metric := appinsights.NewMetricTelemetry("Cpu usage", percent)
		metric.Properties["CPU #"] = strconv.Itoa(cpuN)
		client.Track(metric)
	}

	for _, usage := range stats.diskUsage {
		usedMetric := appinsights.NewMetricTelemetry("Disk usage", float64(usage.Used))
		usedMetric.Properties["Disk"] = usage.Path
		client.Track(usedMetric)
		freeMetric := appinsights.NewMetricTelemetry("Disk free", float64(usage.Free))
		freeMetric.Properties["Disk"] = usage.Path
		client.Track(freeMetric)
	}

	if stats.memory != nil {
		client.TrackMetric("Memory used", float64(stats.memory.Used))
		client.TrackMetric("Memory available", float64(stats.memory.Total-stats.memory.Used))
	}
	if stats.diskIO != nil {
		client.TrackMetric("Disk read", float64(stats.diskIO.ReadBps))
		client.TrackMetric("Disk write", float64(stats.diskIO.WriteBps))
	}

	if stats.netIO != nil {
		client.TrackMetric("Network read", float64(stats.netIO.ReadBps))
		client.TrackMetric("Network write", float64(stats.netIO.WriteBps))
	}

	if len(stats.gpus) > 0 {
		for cpuN, usage := range stats.gpus {
			gpuMetric := appinsights.NewMetricTelemetry("Gpu usage", usage.GPU)
			gpuMetric.Properties["GPU #"] = strconv.Itoa(cpuN)
			client.Track(gpuMetric)

			gpuMemoryMetric := appinsights.NewMetricTelemetry("Gpu memory usage", usage.Memory)
			gpuMemoryMetric.Properties["GPU #"] = strconv.Itoa(cpuN)
			client.Track(gpuMemoryMetric)
		}
	}

	client.Channel().Flush()
}
