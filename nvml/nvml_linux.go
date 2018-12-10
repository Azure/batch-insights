// +build linux

package nvml

import (
	nvml_linux "github.com/mindprince/gonvml"
)

type LinuxNvmlClient struct {
}

func New() (*LinuxNvmlClient, error) {
	client := LinuxNvmlClient{}

	return &client, nil
}

func (client *LinuxNvmlClient) Init() error {
	return nvml_linux.Init()
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
	value, err := nvml_linux.UtilizationRates(nvml_linux.Device({dev: device}))
	if err != nil {
		return GPUUtilization{GPU: 0, Memory: 0}, err
	}

	use := GPUUtilization{
		GPU:    uint(value.GPU),
		Memory: uint(value.Memory),
	}
	return use, nil
}

func (client *LinuxNvmlClient) DeviceGetMemoryInfo(device Device) (Memory, error) {
	use, err := nvml_linux.MemoryInfo(nvml_linux.Device({dev: device}))
	if err != nil {
		return Memory(use), err
	}
	return Memory(use), nil
}

func (client *LinuxNvmlClient) DeviceGetHandleByIndex(index uint) (Device, error) {
	device, err := nvml_linux.DeviceHandleByIndex(uint32(index))
	if err != nil {
		return Device(device.dev), err
	}
	return Device(device.dev), nil
}
