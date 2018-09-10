package main

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type NodeStats struct {
	memory      *mem.VirtualMemoryStat
	cpuPercents []float64
	diskUsage   []*disk.UsageStat
	diskIO      *IOStats
	netIO       *IOStats
}
