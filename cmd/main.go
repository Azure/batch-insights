package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var IS_PLATFORM_WINDOWS = runtime.GOOS == "windows"

const STATS_POLL_RATE = time.Duration(5) * time.Second

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getDiskToWatch() []string {
	if IS_PLATFORM_WINDOWS == true {
		return []string{"C:/", "D:/"}
	} else {
		var osDisk = "/"
		var userDisk = "/mnt/resources"
		var exists, _ = pathExists(userDisk)

		if !exists {
			userDisk = "/mnt"
		}
		return []string{osDisk, userDisk}
	}
}

func main() {
	printSystemInfo()
	listenForStats()
}

func listenForStats() {
	var diskIO = IOAggregator{}
	var netIO = IOAggregator{}
	var appInsightsService = createAppInsightsService()

	for _ = range time.Tick(STATS_POLL_RATE) {

		v, _ := mem.VirtualMemory()
		var cpus, _ = cpu.Percent(time.Duration(5), true)

		var stats = NodeStats{
			memory:      v,
			cpuPercents: cpus,
			diskUsage:   getDiskUsage(),
			diskIO:      getDiskIO(&diskIO),
			netIO:       getNetIO(&netIO),
		}

		if appInsightsService != nil {
			appInsightsService.UploadStats(stats)
		} else {
			printStats(stats)
		}
	}
}

func getDiskIO(diskIO *IOAggregator) *IOStats {
	var counters, err = disk.IOCounters()

	if err != nil {
		fmt.Println(err)
		return nil
	}
	var readBytes uint64 = 0
	var writeBytes uint64 = 0

	for _, v := range counters {
		readBytes += v.ReadBytes
		writeBytes += v.WriteBytes
	}
	var stats = diskIO.UpdateAggregates(readBytes, writeBytes)
	return &stats
}

func getNetIO(diskIO *IOAggregator) *IOStats {
	var counters, err = net.IOCounters(false)

	if err != nil {
		fmt.Println(err)
	} else if len(counters) >= 1 {
		var stats = diskIO.UpdateAggregates(counters[0].BytesRecv, counters[0].BytesSent)
		return &stats
	}
	return nil
}

func getDiskUsage() []*disk.UsageStat {
	var disks = getDiskToWatch()
	var stats []*disk.UsageStat

	for _, diskPath := range disks {
		usage, err := disk.Usage(diskPath)
		if err == nil {
			stats = append(stats, usage)
		} else {
			fmt.Println(err)
		}
	}

	return stats
}

func printSystemInfo() {
	fmt.Printf("System information:\n")
	fmt.Printf("   OS: %s\n", runtime.GOOS)
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
		fmt.Printf("Disk IO: R:%sps, W:%sps\n", humanize.Bytes(stats.diskIO.readBps), humanize.Bytes(stats.diskIO.writeBps))
	}

	if stats.netIO != nil {
		fmt.Printf("NET IO: R:%sps, S:%sps\n", humanize.Bytes(stats.netIO.readBps), humanize.Bytes(stats.netIO.writeBps))
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

func createAppInsightsService() *AppInsightsService {
	var key, keySet = os.LookupEnv("APP_INSIGHTS_INSTRUMENTATION_KEY")
	var poolId = os.Getenv("AZ_BATCH_POOL_ID")
	var nodeId = os.Getenv("AZ_BATCH_NODE_ID")

	if keySet {
		service := NewAppInsightsService(key, poolId, nodeId)
		return &service
	} else {
		fmt.Println("APP_INSIGHTS_INSTRUMENTATION_KEY is not set will to upload to app insights")
		return nil
	}
}
