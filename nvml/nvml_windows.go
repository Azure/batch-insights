// +build windows

package nvml

import (
	nvml_win "github.com/mxpv/nvml-go"
)

type WinDevice struct {
	handle nvml_win.Device
}

func New() (*WinNvmlClient, error) {
	api, err := nvml_win.New("")

	if err != nil {
		return nil, err
	}

	client := WinNvmlClient{
		api: api,
	}

	return &client, nil
}

type WinNvmlClient struct {
	api *nvml_win.API
}

func (client *WinNvmlClient) Init() error {
	return client.api.Init()
}

func (client *WinNvmlClient) Shutdown() error {
	return client.api.Shutdown()
}

func (client *WinNvmlClient) GetDeviceCount() (uint, error) {
	value, err := client.api.DeviceGetCount()
	if err != nil {
		return 0, err
	}

	return uint(value), nil
}

func (client *WinNvmlClient) DeviceGetUtilizationRates(device Device) (GPUUtilization, error) {
	winDevice := device.(WinDevice)
	value, err := client.api.DeviceGetUtilizationRates(winDevice.handle)
	if err != nil {
		return GPUUtilization{GPU: 0, Memory: 0}, err
	}

	use := GPUUtilization{
		GPU:    uint(value.GPU),
		Memory: uint(value.Memory),
	}
	return use, nil
}

func (client *WinNvmlClient) DeviceGetMemoryInfo(device Device) (Memory, error) {
	winDevice := device.(WinDevice)
	use, err := client.api.DeviceGetMemoryInfo(winDevice.handle)
	if err != nil {
		return Memory(use), err
	}
	return Memory(use), nil
}

func (client *WinNvmlClient) DeviceGetHandleByIndex(index uint) (Device, error) {
	handle, err := client.api.DeviceGetHandleByIndex(uint32(index))
	if err != nil {
		return Device(WinDevice{handle: handle}), err
	}
	return Device(WinDevice{handle: handle}), nil
}
