package batchinsights

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Azure/batch-insights/pkg/cpu"
	"github.com/Azure/batch-insights/pkg/disk"
	"github.com/Azure/batch-insights/pkg/utils"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

const STATS_POLL_RATE = time.Duration(5) * time.Second

func ListenForStats(poolId string, nodeId string, appInsightsKey string) {
	var netIO = utils.IOAggregator{}
	var gpuStatsCollector = NewGPUStatsCollector()
	defer gpuStatsCollector.Shutdown()

	var appInsightsService = createAppInsightsService(poolId, nodeId, appInsightsKey)

	for _ = range time.Tick(STATS_POLL_RATE) {
		gpuStatsCollector.GetStats()

		v, _ := mem.VirtualMemory()
		var cpus, err = cpu.PerCpuPercent()
		if err != nil {
			fmt.Println(err)
		}
		var stats = NodeStats{
			memory:      v,
			cpuPercents: cpus,
			diskUsage:   disk.GetDiskUsage(),
			diskIO:      disk.DiskIO(),
			netIO:       getNetIO(&netIO),
			gpus:        gpuStatsCollector.GetStats(),
		}

		if appInsightsService != nil {
			appInsightsService.UploadStats(stats)
		} else {
			printStats(stats)
		}
	}
}

func getNetIO(diskIO *utils.IOAggregator) *utils.IOStats {
	var counters, err = net.IOCounters(false)

	if err != nil {
		fmt.Println(err)
	} else if len(counters) >= 1 {
		var stats = diskIO.UpdateAggregates(counters[0].BytesRecv, counters[0].BytesSent)
		return &stats
	}
	return nil
}

func PrintSystemInfo() {
	fmt.Printf("System information:\n")
	fmt.Printf("   OS: %s\n", runtime.GOOS)
}

func getConfiguration() {

}

func printStats(stats NodeStats) {
	fmt.Printf("========================= Stats =========================\n")
	fmt.Printf("Cpu percent:           %f%%, %v cpu(s)\n", avg(stats.cpuPercents), len(stats.cpuPercents))
	fmt.Printf("Memory used:           %s/%s\n", humanize.Bytes(stats.memory.Used), humanize.Bytes(stats.memory.Total))
	fmt.Printf("Disk usage:\n")
	for _, usage := range stats.diskUsage {
		fmt.Printf("  - %s: %s/%s (%v%%)\n", usage.Path, humanize.Bytes(usage.Used), humanize.Bytes(usage.Total), usage.UsedPercent)
	}

	if stats.diskIO != nil {
		fmt.Printf("Disk IO: R:%sps, W:%sps\n", humanize.Bytes(stats.diskIO.ReadBps), humanize.Bytes(stats.diskIO.WriteBps))
	}

	if stats.netIO != nil {
		fmt.Printf("NET IO: R:%sps, S:%sps\n", humanize.Bytes(stats.netIO.ReadBps), humanize.Bytes(stats.netIO.WriteBps))
	}

	if len(stats.gpus) > 0 {
		fmt.Printf("GPU(s) usage:\n")
		for _, usage := range stats.gpus {
			fmt.Printf("  - GPU: %f%%, Memory: %f%%\n", usage.GPU, usage.Memory)
		}
	}
	fmt.Println()
	fmt.Println()
}

func avg(array []float64) float64 {
	var total float64 = 0
	for _, value := range array {
		total += value
	}
	return total / float64(len(array))
}

func createAppInsightsService(poolId string, nodeId string, appInsightsKey string) *AppInsightsService {
	if appInsightsKey != "" {
		service := NewAppInsightsService(appInsightsKey, poolId, nodeId)
		return &service
	} else {
		fmt.Println("APP_INSIGHTS_INSTRUMENTATION_KEY is not set will to upload to app insights")
		return nil
	}
}
