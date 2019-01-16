// +build linux

package disk

import (
	"fmt"

	"github.com/Azure/batch-insights/pkg/utils"
	psutils_disk "github.com/shirou/gopsutil/disk"
)

var diskIO = utils.IOAggregator{}

func DiskIO() *utils.IOStats {
	var counters, err = psutils_disk.IOCounters()

	if err != nil {
		fmt.Println("Error while retrieving Disk IO", err)
		return nil
	}
	var readBytes uint64 = 0
	var writeBytes uint64 = 0

	for _, v := range counters {
		fmt.Println("stats", v.WriteBytes, v.WriteTime, v.WriteCount)
		readBytes += v.ReadBytes
		writeBytes += v.WriteBytes
	}
	var stats = diskIO.UpdateAggregates(readBytes, writeBytes)
	return &stats
}
