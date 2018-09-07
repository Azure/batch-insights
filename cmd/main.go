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
)

type NodeStats struct {
	memoryTotal uint64
	memoryFree  uint64
	cpuPercents []float64
	diskUsage   []*disk.UsageStat
}

var IS_PLATFORM_WINDOWS = runtime.GOOS == "windows"

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

	var stats = getStats()
	printStats(stats)
}

func getStats() NodeStats {
	v, _ := mem.VirtualMemory()
	var cpus, _ = cpu.Percent(time.Duration(5), true)
	return NodeStats{
		memoryTotal: v.Total,
		memoryFree:  v.Free,
		cpuPercents: cpus,
		diskUsage:   getDiskUsage(),
	}
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
	fmt.Printf("Cpu percent:           %f%% %v, cpu(s)\n", avg(stats.cpuPercents), len(stats.cpuPercents))
	fmt.Printf("Memory used:           %s/%s\n", humanize.Bytes(stats.memoryTotal-stats.memoryFree), humanize.Bytes(stats.memoryTotal))
	fmt.Printf("Disk usage:\n")
	for _, usage := range stats.diskUsage {
		fmt.Printf("  - %s: %s/%s (%v%%)\n", usage.Path, humanize.Bytes(usage.Used), humanize.Bytes(usage.Total), usage.UsedPercent)
	}

}

func avg(array []float64) float64 {
	var total float64 = 0
	for _, value := range array {
		total += value
	}
	return total / float64(len(array))
}
