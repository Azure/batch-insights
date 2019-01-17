package disk

import (
	"fmt"
	"os"
	"runtime"

	psutils_disk "github.com/shirou/gopsutil/disk"
)

var IS_PLATFORM_WINDOWS = runtime.GOOS == "windows"

func GetDiskUsage() []*psutils_disk.UsageStat {
	var disks = getDiskToWatch()
	var stats []*psutils_disk.UsageStat

	for _, diskPath := range disks {
		usage, err := psutils_disk.Usage(diskPath)
		if err == nil {
			stats = append(stats, usage)
		} else {
			fmt.Println(err)
		}
	}

	return stats
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
