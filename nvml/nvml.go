package nvml

type NvmlClient interface {
	Init() error
	Shutdown() error
	GetDeviceCount() (uint, error)

	DeviceGetHandleByIndex(index uint) (Device, error)
	DeviceGetMemoryInfo(device Device) (Memory, error)
	DeviceGetUtilizationRates(device Device) (GPUUtilization, error)
}

type GPUUtilization struct {
	GPU    uint
	Memory uint
}

type Memory struct {
	Total uint64 // Total installed FB memory (in bytes).
	Free  uint64 // Unallocated FB memory (in bytes).
	Used  uint64 // Allocated FB memory (in bytes).
}

type Device interface {
}
