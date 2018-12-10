// +build linux

package cpu

import (
	psutils_cpu "github.com/shirou/gopsutil/cpu"
)

func PerCpuPercent() ([]float64, error) {
	return psutils_cpu.Percent(0, true)
}
