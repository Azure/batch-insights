package batchinsights

import (
	"fmt"
	"github.com/Azure/batch-insights/nvml"
)

type GPUStatsCollector struct {
	nvml        nvml.NvmlClient
	deviceCount uint
}

type GPUStats struct {
}

func NewGPUStatsCollector() GPUStatsCollector {
	nvmlClient, err := nvml.New()

	if err != nil {
		fmt.Println("No GPU detected. Nvidia driver might be missing")
	} else {
		err = nvmlClient.Init()

		if err != nil {
			fmt.Println("Error while loading the GPU", err)
			nvmlClient = nil
		} else {
			deviceCount, err := nvmlClient.GetDeviceCount()

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("NVML is loaded found %d gpus\n", deviceCount)
			}

			return GPUStatsCollector{
				nvml:        nvmlClient,
				deviceCount: deviceCount,
			}
		}
	}
	return GPUStatsCollector{}
}

func (gpu GPUStatsCollector) GetStats() {
	if gpu.nvml == nil {
		return
	}

	for i := uint(0); i < gpu.deviceCount; i++ {
		device, err := gpu.nvml.DeviceGetHandleByIndex(i)
		if err != nil {
			fmt.Println(err)
			continue
		}

		memory, err := gpu.nvml.DeviceGetMemoryInfo(device)

		if err != nil {
			fmt.Println(err)
		}

		use, err := gpu.nvml.DeviceGetUtilizationRates(device)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Utilization (#%d): GPU: %d%%, Memory: %d%%\n", i, use.GPU, use.Memory)
		fmt.Printf("Memory usage (#%d): Free: %d, Total: %d\n", i, memory.Free, memory.Total)
	}
}

func (gpu GPUStatsCollector) Shutdown() {
	if gpu.nvml == nil {
		return
	}
	gpu.nvml.Shutdown()
}
