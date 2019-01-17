// +build windows

package disk

import (
	"context"
	"fmt"

	"github.com/Azure/batch-insights/pkg/utils"
	"github.com/Azure/batch-insights/pkg/wmi"
)

type win32_PerfRawData_PerfDisk_PhysicalDisk struct {
	Name                      string
	AvgDiskBytesPerRead       uint64
	AvgDiskBytesPerRead_Base  uint64
	AvgDiskBytesPerWrite      uint64
	AvgDiskBytesPerWrite_Base uint64
	AvgDiskReadQueueLength    uint64
	AvgDiskWriteQueueLength   uint64
	AvgDisksecPerRead         uint64
	AvgDisksecPerWrite        uint64
}

var diskIO = utils.IOAggregator{}

func DiskIO() *utils.IOStats {
	return DiskIOWithContext(context.Background())
}

func DiskIOWithContext(ctx context.Context, names ...string) *utils.IOStats {
	var ret []win32_PerfRawData_PerfDisk_PhysicalDisk

	q := wmi.CreateQuery(&ret, "WHERE NOT Name LIKE '%_Total'")
	err := wmi.QueryWithContext(ctx, q, &ret)
	if err != nil {
		fmt.Println("Error while retrieving DISK IO", err)
		return nil
	}

	var readBytes uint64 = 0
	var writeBytes uint64 = 0
	for _, v := range ret {
		readBytes += v.AvgDiskBytesPerRead
		writeBytes += v.AvgDiskBytesPerWrite
	}
	stats := diskIO.UpdateAggregates(readBytes, writeBytes)

	return &stats
}
