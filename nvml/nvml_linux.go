// +build linux

package nvml

import (
	nvml_linux "github.com/mindprince/gonvml"
)

type LinuxDevice = nvml_linux.Device

type LinuxNvmlClient struct {
}

func New() (*LinuxNvmlClient, error) {
	client := LinuxNvmlClient{}

	return &client, nil
}

func (client *LinuxNvmlClient) Init() error {
	return nvml_linux.Initialize()
}

func (client *LinuxNvmlClient) Shutdown() error {
	return nvml_linux.Shutdown()
}

func (client *LinuxNvmlClient) GetDeviceCount() (uint, error) {
	value, err := nvml_linux.DeviceCount()
	if err != nil {
		return 0, err
	}

	return uint(value), nil
}

func (client *LinuxNvmlClient) DeviceGetUtilizationRates(device Device) (GPUUtilization, error) {
	linuxDevice := device.(LinuxDevice)
	gpu, memory, err := linuxDevice.UtilizationRates()
	if err != nil {
		return GPUUtilization{GPU: 0, Memory: 0}, err
	}

	use := GPUUtilization{
		GPU:    gpu,
		Memory: memory,
	}
	return use, nil
}

func (client *LinuxNvmlClient) DeviceGetMemoryInfo(device Device) (Memory, error) {
	linuxDevice := device.(LinuxDevice)
	total, used, err := linuxDevice.MemoryInfo()
	if err != nil {
		return Memory{Used: used, Total: total}, err
	}
	return Memory{Used: used, Total: total}, nil
}

func (client *LinuxNvmlClient) DeviceGetHandleByIndex(index uint) (Device, error) {
	device, err := nvml_linux.DeviceHandleByIndex(uint(index))
	if err != nil {
		return Device(device), err
	}
	return Device(device), nil
}
