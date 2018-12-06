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
		// client.TrackMetric("Disk usage", disk_usage.used, properties={"Disk": name})
		// client.TrackMetric("Disk free", disk_usage.free, properties={"Disk": name})
		usedMetric := appinsights.NewMetricTelemetry("Disk usage", float64(usage.Used))
		usedMetric.Properties["Disk"] = usage.Path
		client.Track(usedMetric)
		freeMetric := appinsights.NewMetricTelemetry("Disk free", float64(usage.Free))
		freeMetric.Properties["Disk"] = usage.Path
		client.Track(freeMetric)
	}

	client.TrackMetric("Memory used", float64(stats.memory.Used))
	client.TrackMetric("Memory available", float64(stats.memory.Free))
	client.TrackMetric("Disk read", float64(stats.diskIO.readBps))
	client.TrackMetric("Disk write", float64(stats.diskIO.writeBps))
	client.TrackMetric("Network read", float64(stats.netIO.readBps))
	client.TrackMetric("Network write", float64(stats.netIO.writeBps))
	client.Channel().Flush()
}
