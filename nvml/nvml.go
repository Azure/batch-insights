package nvml

type NvmlClient interface {
	Init() error
	Shutdown() error
	GetDeviceCount() (uint, error)
}
