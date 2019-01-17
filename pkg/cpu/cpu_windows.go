// +build windows

package cpu

import (
	"context"
	"github.com/Azure/batch-insights/pkg/wmi"
)

type CPUStat struct {
	value     uint64
	timestamp uint64
}

var lastCpus map[string]CPUStat = make(map[string]CPUStat)

type win32_PerfRawData_Counters_ProcessorInformation struct {
	Name                  string
	PercentDPCTime        uint64
	PercentIdleTime       uint64
	PercentUserTime       uint64
	PercentProcessorTime  uint64
	PercentInterruptTime  uint64
	PercentPriorityTime   uint64
	PercentPrivilegedTime uint64
	TimeStamp_Sys100NS    uint64
	InterruptsPerSec      uint32
	ProcessorFrequency    uint32
	DPCRate               uint32
}

func PerCpuPercent() ([]float64, error) {
	return perCPUPercentWithContext(context.Background())
}

func perCPUPercentWithContext(ctx context.Context) ([]float64, error) {
	var ret []float64
	stats, err := perfInfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range stats {
		last := lastCpus[v.Name]

		lastCpus[v.Name] = CPUStat{
			value:     v.PercentProcessorTime,
			timestamp: v.TimeStamp_Sys100NS,
		}

		if last.timestamp != 0 {
			value := (1 - (float64(v.PercentProcessorTime-last.value) / float64(v.TimeStamp_Sys100NS-last.timestamp))) * 100
			ret = append(ret, value)
		}
	}
	return ret, nil
}

// PerfInfo returns the performance counter's instance value for ProcessorInformation.
// Name property is the key by which overall, per cpu and per core metric is known.
func perfInfoWithContext(ctx context.Context) ([]win32_PerfRawData_Counters_ProcessorInformation, error) {
	var ret []win32_PerfRawData_Counters_ProcessorInformation

	q := wmi.CreateQuery(&ret, "WHERE NOT Name LIKE '%_Total'")
	err := wmi.QueryWithContext(ctx, q, &ret)
	if err != nil {
		return []win32_PerfRawData_Counters_ProcessorInformation{}, err
	}

	return ret, err
}
