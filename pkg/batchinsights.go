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

func getSamplingRate(rate time.Duration) time.Duration {
	if rate <= time.Duration(0) {
		return DefaultSamplingRate
	}
	return rate
}

// ListenForStats Start the sanpling of node metrics
func ListenForStats(config Config) {
	var netIO = utils.IOAggregator{}

	var gpuStatsCollector = NewGPUStatsCollector()
	defer gpuStatsCollector.Shutdown()

	var appInsightsService = createAppInsightsService(config)

	for range time.Tick(getSamplingRate(config.SamplingRate)) {
		gpuStatsCollector.GetStats()

		var stats = NodeStats{}

		if !config.Disable.Memory {
			v, err := mem.VirtualMemory()
			if err == nil {
				stats.Memory = v
			} else {
				fmt.Println(err)
			}
		}
		if !config.Disable.CPU {
			cpus, err := cpu.PerCpuPercent()
			if err == nil {
				stats.CpuPercents = cpus
			} else {
				fmt.Println(err)
			}
		}
		if !config.Disable.DiskUsage {
			stats.DiskUsage = disk.GetDiskUsage()
		}
		if !config.Disable.DiskIO {
			stats.DiskIO = disk.DiskIO()
		}
		if !config.Disable.NetworkIO {
			stats.NetIO = getNetIO(&netIO)
		}
		if !config.Disable.GPU {
			stats.Gpus = gpuStatsCollector.GetStats()
		}

		processes, err := ListProcesses(config.Processes)
		if err == nil {
			stats.Processes = processes
		} else {
			fmt.Println(err)
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

// PrintSystemInfo print system info needed
func PrintSystemInfo() {
	fmt.Printf("System information:\n")
	fmt.Printf("   OS: %s\n", runtime.GOOS)
}

func getConfiguration() {

}

func printStats(stats NodeStats) {
	fmt.Printf("========================= Stats =========================\n")
	fmt.Printf("Cpu percent:           %f%%, %v cpu(s)\n", avg(stats.CpuPercents), len(stats.CpuPercents))
	fmt.Printf("Memory used:           %s/%s\n", humanize.Bytes(stats.Memory.Used), humanize.Bytes(stats.Memory.Total))

	if len(stats.DiskUsage) > 0 {
		fmt.Printf("Disk usage:\n")
		for _, usage := range stats.DiskUsage {
			fmt.Printf("  - %s: %s/%s (%v%%)\n", usage.Path, humanize.Bytes(usage.Used), humanize.Bytes(usage.Total), usage.UsedPercent)
		}
	}

	if stats.DiskIO != nil {
		fmt.Printf("Disk IO: R:%sps, W:%sps\n", humanize.Bytes(stats.DiskIO.ReadBps), humanize.Bytes(stats.DiskIO.WriteBps))
	}

	if stats.NetIO != nil {
		fmt.Printf("NET IO: R:%sps, S:%sps\n", humanize.Bytes(stats.NetIO.ReadBps), humanize.Bytes(stats.NetIO.WriteBps))
	}

	if len(stats.Gpus) > 0 {
		fmt.Printf("GPU(s) usage:\n")
		for _, usage := range stats.Gpus {
			fmt.Printf("  - GPU: %f%%, Memory: %f%%\n", usage.GPU, usage.Memory)
		}
	}

	if len(stats.Processes) > 0 {
		fmt.Printf("Tracked processes:\n")
		for _, process := range stats.Processes {
			fmt.Printf("  - %s (%d), CPU: %f%%, Memory: %s\n", process.name, process.pid, process.cpu, humanize.Bytes(process.memory))
		}
	}

	fmt.Println()
	fmt.Println()
}

func avg(array []float64) float64 {
	var total float64
	for _, value := range array {
		total += value
	}
	return total / float64(len(array))
}

func createAppInsightsService(config Config) *AppInsightsService {
	if config.InstrumentationKey != "" {
		service := NewAppInsightsService(config.InstrumentationKey, config.PoolID, config.NodeID)
		return &service
	} else {
		fmt.Println("APP_INSIGHTS_INSTRUMENTATION_KEY is not set; will not upload to Application Insights")
		return nil
	}
}
