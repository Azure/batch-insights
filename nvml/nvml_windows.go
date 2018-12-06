// +build !linux

package nvml

import (
	nvml_win "github.com/mxpv/nvml-go"
)

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
