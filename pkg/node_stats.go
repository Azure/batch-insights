package batchinsights

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/Azure/batch-insights/pkg/utils"
)

// ProcessPerfInfo Process specific information
type ProcessPerfInfo struct {
	pid    int32
	name   string
	cpu    float64
	memory uint64
}

// NodeStats Combined model for all metrics being collected at the given interal
type NodeStats struct {
	Memory      *mem.VirtualMemoryStat
	CPUPercents []float64
	DiskUsage   []*disk.UsageStat
	DiskIO      *utils.IOStats
	NetIO       *utils.IOStats
	Gpus        []GPUUsage
	Processes   []*ProcessPerfInfo
}
