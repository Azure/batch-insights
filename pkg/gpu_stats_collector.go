package batchinsights

import (
	"fmt"
	"github.com/Azure/batch-insights/nvml"
)

type GPUStatsCollector struct {
	nvml        nvml.NvmlClient
	deviceCount uint
}

type GPUUsage struct {
	GPU    float64
	Memory float64
}

func NewGPUStatsCollector() GPUStatsCollector {
	nvmlClient, err := nvml.New()

	if err != nil {
		fmt.Println("No GPU detected. Nvidia driver might be missing")
	} else {
		err = nvmlClient.Init()

		if err != nil {
			fmt.Println("No GPU detected. Nvidia driver might be missing. Error while initializing NVML", err)
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

func (gpu GPUStatsCollector) GetStats() []GPUUsage {
	if gpu.nvml == nil {
		return nil
	}

	var uses []GPUUsage

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

		usage := GPUUsage{
			GPU:    float64(use.GPU),
			Memory: float64(memory.Used) / float64(memory.Total) * 100,
		}
		uses = append(uses, usage)
	}
	return uses
}

func (gpu GPUStatsCollector) Shutdown() {
	if gpu.nvml == nil {
		return
	}
	gpu.nvml.Shutdown()
}
