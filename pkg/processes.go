package batchinsights

import (
	"strings"

	"github.com/shirou/gopsutil/process"
)

// love to program in go
func containsCaseInsensitive(xs []string, str string) bool {

	for _, x := range xs {
		if strings.EqualFold(str, x) {
			return true
		}
	}

	return false

}

// ListProcesses Retrieve process cpu, memory, etc usage for the given list of process names
func ListProcesses(processNames []string) ([]*ProcessPerfInfo, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	ps := []*ProcessPerfInfo{}
	for _, pid := range pids {

		// if err != nil, process has probably disappeared, continue on
		if p, err := process.NewProcess(pid); err == nil {

			name, err := p.Name()
			if err != nil {
				// process might have disappeared
				continue
			}

			// check if we should include it
			if !containsCaseInsensitive(processNames, name) {
				continue
			}

			cpuPercent, err := p.CPUPercent()
			if err != nil {
				// process might have disappeared
				continue
			}

			memoryInfoStat, err := p.MemoryInfo()
			if err != nil {
				// process might have disappeared
				continue
			}

			ps = append(ps, &ProcessPerfInfo{
				pid:    pid,
				name:   name,
				cpu:    cpuPercent,
				memory: memoryInfoStat.VMS,
			})
		}

	}

	return ps, err
}
