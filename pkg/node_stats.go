package batchinsights

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/Azure/batch-insights/pkg/utils"
)

type NodeStats struct {
	memory      *mem.VirtualMemoryStat
	cpuPercents []float64
	diskUsage   []*disk.UsageStat
	diskIO      *utils.IOStats
	netIO       *utils.IOStats
	gpus        []GPUUsage
}
