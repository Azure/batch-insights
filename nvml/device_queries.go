package nvml

// #define NVML_DEVICE_PCI_BUS_ID_BUFFER_SIZE 16
// typedef struct nvmlPciInfo_st {
//     char busId[NVML_DEVICE_PCI_BUS_ID_BUFFER_SIZE]; //!< The tuple domain:bus:device.function PCI identifier
//     unsigned int domain; //!< The PCI domain on which the device's bus resides, 0 to 0xffff
//     unsigned int bus; //!< The bus on which the device resides, 0 to 0xff
//     unsigned int device; //!< The device's id on the bus, 0 to 31
//     unsigned int pciDeviceId; //!< The combined 16-bit device id and 16-bit vendor id
//
//     // Added in NVML 2.285 API
//     unsigned int pciSubSystemId;
//
//     // NVIDIA reserved for internal use only
//     unsigned int reserved0;
//     unsigned int reserved1;
//     unsigned int reserved2;
//     unsigned int reserved3;
// } nvmlPciInfo_t;
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// DeviceGetAPIRestriction retrieves the root/admin permissions on the target API.
// See nvmlRestrictedAPI_t for the list of supported APIs.
// If an API is restricted only root users can call that API.
// See nvmlDeviceSetAPIRestriction to change current permissions.
func (a API) DeviceGetAPIRestriction(device Device, apiType RestrictedAPI) (bool, error) {
	var state int32
	if err := a.call(a.nvmlDeviceGetAPIRestriction, uintptr(device), uintptr(apiType), uintptr(unsafe.Pointer(&state))); err != nil {
		return false, err
	}

	if state > 0 {
		return true, nil
	}

	return false, nil
}

// DeviceGetApplicationsClock retrieves the current setting of a clock that applications will use unless an overspec
// situation occurs. Can be changed using DeviceSetApplicationsClocks.
func (a API) DeviceGetApplicationsClock(device Device, clockType ClockType) (clockMHz uint32, err error) {
	err = a.call(a.nvmlDeviceGetApplicationsClock, uintptr(device), uintptr(clockType), uintptr(unsafe.Pointer(&clockMHz)))
	return
}

// DeviceGetAutoBoostedClocksEnabled retrieve the current state of Auto Boosted clocks on a device and store it in isEnabled.
// Auto Boosted clocks are enabled by default on some hardware, allowing the GPU to run at higher clock rates to
// maximize performance as thermal limits allow.
// On Pascal and newer hardware, Auto Aoosted clocks are controlled through application clocks.
func (a API) DeviceGetAutoBoostedClocksEnabled(device Device) (isEnabled, defaultIsEnabled bool, err error) {
	var isEnabledInt int32
	var defaultIsEnabledInt int32

	err = a.call(a.nvmlDeviceGetAutoBoostedClocksEnabled, uintptr(device), uintptr(unsafe.Pointer(&isEnabledInt)), uintptr(unsafe.Pointer(&defaultIsEnabledInt)))
	if err != nil {
		return
	}

	if isEnabledInt > 0 {
		isEnabled = true
	} else {
		isEnabled = false
	}

	if defaultIsEnabledInt > 0 {
		defaultIsEnabled = true
	} else {
		defaultIsEnabled = false
	}

	return
}

// DeviceGetBAR1MemoryInfo gets Total, Available and Used size of BAR1 memory.
// BAR1 is used to map the FB (device memory) so that it can be directly accessed by the CPU or
// by 3rd party devices (peer-to-peer on the PCIE bus).
func (a API) DeviceGetBAR1MemoryInfo(device Device) (mem BAR1Memory, err error) {
	err = a.call(a.nvmlDeviceGetBAR1MemoryInfo, uintptr(device), uintptr(unsafe.Pointer(&mem)))
	return
}

// DeviceGetBoardID retrieves the device boardId from 0-N. Devices with the same boardId indicate GPUs connected to
// the same PLX. Use in conjunction with DeviceGetMultiGpuBoard() to decide if they are on the same board as well.
// The boardId returned is a unique ID for the current configuration.
// Uniqueness and ordering across reboots and system configurations is not guaranteed (i.e. if a Tesla K40c returns
// 0x100 and the two GPUs on a Tesla K10 in the same system returns 0x200 it is not guaranteed they will always return
// those values but they will always be different from each other).
func (a API) DeviceGetBoardID(device Device) (boardID uint32, err error) {
	err = a.call(a.nvmlDeviceGetBoardId, uintptr(device), uintptr(unsafe.Pointer(&boardID)))
	return
}

// DeviceGetBoardPartNumber retrieves the the device board part number which is programmed into the board's InfoROM
func (a API) DeviceGetBoardPartNumber(device Device) (string, error) {
	const bufferSize = 128

	buffer := [bufferSize]C.char{}
	if err := a.call(a.nvmlDeviceGetBoardPartNumber, uintptr(device), uintptr(unsafe.Pointer(&buffer[0])), bufferSize); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}

// DeviceGetBrand retrieves the brand of this device.
func (a API) DeviceGetBrand(device Device) (brand BrandType, err error) {
	err = a.call(a.nvmlDeviceGetBrand, uintptr(device), uintptr(unsafe.Pointer(&brand)))
	return
}

func (a API) DeviceGetBridgeChipInfo() {

}

// DeviceGetClock retrieves the clock speed for the clock specified by the clock type and clock ID.
func (a API) DeviceGetClock(device Device, clockType ClockType, clockID ClockID) (clockMHz uint32, err error) {
	err = a.call(a.nvmlDeviceGetClock, uintptr(device), uintptr(clockType), uintptr(clockID), uintptr(unsafe.Pointer(&clockMHz)))
	return
}

// DeviceGetClockInfo retrieves the current clock speeds for the device.
func (a API) DeviceGetClockInfo(device Device, clockType ClockType) (clock uint32, err error) {
	err = a.call(a.nvmlDeviceGetClockInfo, uintptr(device), uintptr(clockType), uintptr(unsafe.Pointer(&clock)))
	return
}

// DeviceGetComputeMode retrieves the current compute mode for the device.
func (a API) DeviceGetComputeMode(device Device) (mode ComputeMode, err error) {
	err = a.call(a.nvmlDeviceGetComputeMode, uintptr(device), uintptr(unsafe.Pointer(&mode)))
	return
}

// DeviceGetComputeRunningProcesses gets information about processes with a compute context on a device.
// This function returns information only about compute running processes (e.g. CUDA application which have
// active context). Any graphics applications (e.g. using OpenGL, DirectX) won't be listed by this function.
// Keep in mind that information returned by this call is dynamic and the number of elements might change in time.
// Allocate more space for infos table in case new compute processes are spawned.
func (a API) DeviceGetComputeRunningProcesses(device Device) ([]ProcessInfo, error) {
	var infoCount uint32

	// Query the current number of running compute processes
	err := a.call(a.nvmlDeviceGetComputeRunningProcesses, uintptr(device), uintptr(unsafe.Pointer(&infoCount)), 0)

	// None are running
	if err == nil || infoCount == 0 {
		return []ProcessInfo{}, nil
	}

	if err != ErrInsufficientSize {
		return nil, err
	}

	list := make([]ProcessInfo, infoCount)
	err = a.call(a.nvmlDeviceGetComputeRunningProcesses, uintptr(device), uintptr(unsafe.Pointer(&infoCount)), uintptr(unsafe.Pointer(&list[0])))
	if err != nil {
		return nil, err
	}

	return list[:infoCount], nil
}

// DeviceGetCount retrieves the number of compute devices in the system. A compute device is a single GPU.
func (a API) DeviceGetCount() (count uint32, err error) {
	err = a.call(a.nvmlDeviceGetCount, uintptr(unsafe.Pointer(&count)))
	return
}

// DeviceGetCudaComputeCapability retrieves the CUDA compute capability of the device.
// Returns the major and minor compute capability version numbers of the device.
// The major and minor versions are equivalent to the CU_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MINOR and
// CU_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MAJOR attributes that would be returned by CUDA's cuDeviceGetAttribute().
func (a API) DeviceGetCudaComputeCapability(device Device) (major, minor int32, err error) {
	err = a.call(a.nvmlDeviceGetCudaComputeCapability, uintptr(device), uintptr(unsafe.Pointer(&major)), uintptr(unsafe.Pointer(&minor)))
	return
}

// DeviceGetCurrPcieLinkGeneration retrieves the current PCIe link generation.
func (a API) DeviceGetCurrPcieLinkGeneration(device Device) (currLinkGen uint32, err error) {
	err = a.call(a.nvmlDeviceGetCurrPcieLinkGeneration, uintptr(device), uintptr(unsafe.Pointer(&currLinkGen)))
	return
}

// DeviceGetCurrPcieLinkWidth retrieves the current PCIe link width.
func (a API) DeviceGetCurrPcieLinkWidth(device Device) (currLinkWidth uint32, err error) {
	err = a.call(a.nvmlDeviceGetCurrPcieLinkWidth, uintptr(device), uintptr(unsafe.Pointer(&currLinkWidth)))
	return
}

// DeviceGetCurrentClocksThrottleReasons retrieves current clocks throttling reasons.
// More than one bit can be enabled at the same time. Multiple reasons can be affecting clocks at once.
func (a API) DeviceGetCurrentClocksThrottleReasons(device Device) (clocksThrottleReasons ClocksThrottleReason, err error) {
	err = a.call(a.nvmlDeviceGetCurrentClocksThrottleReasons, uintptr(device), uintptr(unsafe.Pointer(&clocksThrottleReasons)))
	return
}

// DeviceGetDecoderUtilization retrieves the current utilization and sampling size in microseconds for the Decoder.
func (a API) DeviceGetDecoderUtilization(device Device) (utilization, samplingPeriodUs uint32, err error) {
	err = a.call(a.nvmlDeviceGetDecoderUtilization, uintptr(device), uintptr(unsafe.Pointer(&utilization)), uintptr(unsafe.Pointer(&samplingPeriodUs)))
	return
}

// DeviceGetDefaultApplicationsClock retrieves the default applications clock that GPU boots with or
// defaults to after DeviceResetApplicationsClocks call.
func (a API) DeviceGetDefaultApplicationsClock(device Device, clockType ClockType) (clockMHz uint32, err error) {
	err = a.call(a.nvmlDeviceGetDefaultApplicationsClock, uintptr(device), uintptr(clockType), uintptr(unsafe.Pointer(&clockMHz)))
	return
}

// DeviceGetDetailedECCErrors retrieves the detailed ECC error counts for the device.
// Only applicable to devices with ECC. Requires NVML_INFOROM_ECC version 2.0 or higher to report aggregate
// location-based ECC counts. Requires NVML_INFOROM_ECC version 1.0 or higher to report all other ECC counts.
// Requires ECC Mode to be enabled.
// Detailed errors provide separate ECC counts for specific parts of the memory system.
// Reports zero for unsupported ECC error counters when a subset of ECC error counters are supported.
// Deprecated: This API supports only a fixed set of ECC error locations.
// On different GPU architectures different locations are supported, see DeviceGetMemoryErrorCounter
func (a API) DeviceGetDetailedECCErrors(device Device, errorType MemoryErrorType, counterType ECCCounterType) (*ECCErrorCounts, error) {
	counts := &ECCErrorCounts{}
	if err := a.call(a.nvmlDeviceGetDetailedEccErrors, uintptr(device), uintptr(errorType), uintptr(counterType), uintptr(unsafe.Pointer(counts))); err != nil {
		return nil, err
	}

	return counts, nil
}

// DeviceGetDisplayActive retrieves the display active state for the device.
// This method indicates whether a display is initialized on the device.
// For example whether X Server is attached to this device and has allocated memory for the screen.
// Display can be active even when no monitor is physically attached.
func (a API) DeviceGetDisplayActive(device Device) (bool, error) {
	var state int32
	if err := a.call(a.nvmlDeviceGetDisplayActive, uintptr(device), uintptr(unsafe.Pointer(&state))); err != nil {
		return false, err
	}

	if state > 0 {
		return true, nil
	}

	return false, nil
}

// DeviceGetDisplayMode retrieves the display mode for the device. This method indicates whether a physical display
// (e.g. monitor) is currently connected to any of the device's connectors.
func (a API) DeviceGetDisplayMode(device Device) (bool, error) {
	var state int32
	if err := a.call(a.nvmlDeviceGetDisplayMode, uintptr(device), uintptr(unsafe.Pointer(&state))); err != nil {
		return false, err
	}

	if state > 0 {
		return true, nil
	}

	return false, nil
}

// DeviceGetDriverModel retrieves the current and pending driver model for the device.
// On Windows platforms the device driver can run in either WDDM or WDM (TCC) mode.
// If a display is attached to the device it must run in WDDM mode. TCC mode is preferred if a display is not attached.
func (a API) DeviceGetDriverModel(device Device) (current, pending DriverModel, err error) {
	err = a.call(a.nvmlDeviceGetDriverModel, uintptr(device), uintptr(unsafe.Pointer(&current)), uintptr(unsafe.Pointer(&pending)))
	return
}

// DeviceGetECCMode retrieves the current and pending ECC modes for the device.
// Only applicable to devices with ECC. Requires NVML_INFOROM_ECC version 1.0 or higher.
// Changing ECC modes requires a reboot. The "pending" ECC mode refers to the target mode following the next reboot.
func (a API) DeviceGetECCMode(device Device) (current, pending bool, err error) {
	var currentInt int32
	var pendingInt int32

	err = a.call(a.nvmlDeviceGetEccMode, uintptr(device), uintptr(unsafe.Pointer(&currentInt)), uintptr(unsafe.Pointer(&pendingInt)))
	if err != nil {
		return
	}

	if currentInt > 0 {
		current = true
	} else {
		current = false
	}

	if pendingInt > 0 {
		pending = true
	} else {
		pending = false
	}

	return
}

// DeviceGetEncoderCapacity retrieves the current capacity of the device's encoder, in macroblocks per second.
func (a API) DeviceGetEncoderCapacity(device Device, encoderQueryType EncoderType) (encoderCapacity uint32, err error) {
	err = a.call(a.nvmlDeviceGetEncoderCapacity, uintptr(device), uintptr(encoderQueryType), uintptr(unsafe.Pointer(&encoderCapacity)))
	return
}

func (a API) DeviceGetEncoderSessions() error {
	return ErrNotImplemented
}

// DeviceGetEncoderStats retrieves the current encoder statistics for a given device.
func (a API) DeviceGetEncoderStats(device Device) (sessionCount, averageFPS, averageLatency uint32, err error) {
	err = a.call(
		a.nvmlDeviceGetEncoderStats,
		uintptr(device),
		uintptr(unsafe.Pointer(&sessionCount)),
		uintptr(unsafe.Pointer(&averageFPS)),
		uintptr(unsafe.Pointer(&averageLatency)))
	return
}

// DeviceGetEncoderUtilization retrieves the current utilization and sampling size in microseconds for the Encoder
func (a API) DeviceGetEncoderUtilization(device Device) (utilization, samplingPeriodUs uint32, err error) {
	err = a.call(a.nvmlDeviceGetEncoderUtilization, uintptr(device), uintptr(unsafe.Pointer(&utilization)), uintptr(unsafe.Pointer(&samplingPeriodUs)))
	return
}

// DeviceGetEnforcedPowerLimit gets the effective power limit that the driver enforces after taking into account all limiters.
// Note: This can be different from the DeviceGetPowerManagementLimit if other limits are set elsewhere.
// This includes the out of band power limit interface
func (a API) DeviceGetEnforcedPowerLimit(device Device) (limit uint32, err error) {
	err = a.call(a.nvmlDeviceGetEnforcedPowerLimit, uintptr(device), uintptr(unsafe.Pointer(&limit)))
	return
}

// DeviceGetFanSpeed retrieves the intended operating speed of the device's fan.
// Note: The reported speed is the intended fan speed. If the fan is physically blocked and unable to spin,
// the output will not match the actual fan speed.
// The fan speed is expressed as a percent of the maximum, i.e. full speed is 100%.
func (a API) DeviceGetFanSpeed(device Device) (speed uint32, err error) {
	err = a.call(a.nvmlDeviceGetFanSpeed, uintptr(device), uintptr(unsafe.Pointer(&speed)))
	return
}

// DeviceGetGPUOperationMode retrieves the current GOM and pending GOM (the one that GPU will switch to after reboot).
// For GK110 M-class and X-class Tesla products from the Kepler family.
// Modes NVML_GOM_LOW_DP and NVML_GOM_ALL_ON are supported on fully supported GeForce products.
// Not supported on Quadro and Tesla C-class products.
func (a API) DeviceGetGPUOperationMode(device Device) (current, pending GPUOperationMode, err error) {
	err = a.call(a.nvmlDeviceGetGpuOperationMode, uintptr(device), uintptr(unsafe.Pointer(&current)), uintptr(unsafe.Pointer(&pending)))
	return
}

// DeviceGetGraphicsRunningProcesses get information about processes with a graphics context on a device.
// This function returns information only about graphics based processes (eg. applications using OpenGL, DirectX).
// Keep in mind that information returned by this call is dynamic and the number of elements might change in time.
// Allocate more space for infos table in case new graphics processes are spawned.
func (a API) DeviceGetGraphicsRunningProcesses(device Device) ([]ProcessInfo, error) {
	var infoCount uint32

	// Query the current number of running compute processes
	err := a.call(a.nvmlDeviceGetGraphicsRunningProcesses, uintptr(device), uintptr(unsafe.Pointer(&infoCount)), 0)

	// None are running
	if err == nil || infoCount == 0 {
		return []ProcessInfo{}, nil
	}

	if err != ErrInsufficientSize {
		return nil, err
	}

	list := make([]ProcessInfo, infoCount)
	err = a.call(a.nvmlDeviceGetGraphicsRunningProcesses, uintptr(device), uintptr(unsafe.Pointer(&infoCount)), uintptr(unsafe.Pointer(&list[0])))
	if err != nil {
		return nil, err
	}

	return list[:infoCount], nil
}

// DeviceGetHandleByIndex acquires the handle for a particular device, based on its index.
func (a API) DeviceGetHandleByIndex(index uint32) (device Device, err error) {
	err = a.call(a.nvmlDeviceGetHandleByIndex, uintptr(index), uintptr(unsafe.Pointer(&device)))
	return
}

// DeviceGetHandleByPciBusId acquires the handle for a particular device, based on its PCI bus id.
// This value corresponds to the nvmlPciInfo_t::busId returned by DeviceGetPciInfo().
// Starting from NVML 5, this API causes NVML to initialize the target GPU NVML may initialize additional GPUs if:
//   - The target GPU is an SLI slave
// Note: NVML 4.304 and older version of nvmlDeviceGetHandleByPciBusId"_v1" returns NVML_ERROR_NOT_FOUND instead of NVML_ERROR_NO_PERMISSION.
func (a API) DeviceGetHandleByPCIBusID(pciBusID string) (device Device, err error) {
	cstr := C.CString(pciBusID)
	defer C.free(unsafe.Pointer(cstr))

	err = a.call(a.nvmlDeviceGetHandleByPciBusId, uintptr(unsafe.Pointer(cstr)), uintptr(unsafe.Pointer(&device)))
	return
}

// DeviceGetHandleBySerial acquires the handle for a particular device, based on its board serial number.
// Starting from NVML 5, this API causes NVML to initialize the target GPU, NVML may initialize additional
// GPUs as it searches for the target GPU
func (a API) DeviceGetHandleBySerial(serial string) (device Device, err error) {
	cstr := C.CString(serial)
	defer C.free(unsafe.Pointer(cstr))

	err = a.call(a.nvmlDeviceGetHandleBySerial, uintptr(unsafe.Pointer(cstr)), uintptr(unsafe.Pointer(&device)))
	return
}

// DeviceGetHandleByUUID acquires the handle for a particular device,
// based on its globally unique immutable UUID associated with each device.
func (a API) DeviceGetHandleByUUID(uuid string) (device Device, err error) {
	cstr := C.CString(uuid)
	defer C.free(unsafe.Pointer(cstr))

	err = a.call(a.nvmlDeviceGetHandleByUUID, uintptr(unsafe.Pointer(cstr)), uintptr(unsafe.Pointer(&device)))
	return
}

// DeviceGetIndex retrieves the NVML index of this device.
func (a API) DeviceGetIndex(device Device) (index uint32, err error) {
	err = a.call(a.nvmlDeviceGetIndex, uintptr(device), uintptr(unsafe.Pointer(&index)))
	return
}

// DeviceGetInforomConfigurationChecksum retrieves the checksum of the configuration stored in the device's infoROM.
// Can be used to make sure that two GPUs have the exact same configuration.
// Current checksum takes into account configuration stored in PWR and ECC infoROM objects.
// Checksum can change between driver releases or when user changes configuration (e.g. disable/enable ECC)
func (a API) DeviceGetInforomConfigurationChecksum(device Device) (checksum uint32, err error) {
	err = a.call(a.nvmlDeviceGetInforomConfigurationChecksum, uintptr(device), uintptr(unsafe.Pointer(&checksum)))
	return
}

// DeviceGetInforomImageVersion retrieves the global infoROM image version. Image version just like VBIOS version
// uniquely describes the exact version of the infoROM flashed on the board in contrast to infoROM object version
// which is only an indicator of supported features.
func (a API) DeviceGetInfoROMImageVersion(device Device) (string, error) {
	buffer := [deviceInfoROMVersionBufferSize]C.char{}
	if err := a.call(a.nvmlDeviceGetInforomImageVersion, uintptr(device), uintptr(unsafe.Pointer(&buffer[0])), deviceInfoROMVersionBufferSize); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}

// DeviceGetInfoROMVersion retrieves the version information for the device's infoROM object.
func (a API) DeviceGetInfoROMVersion(device Device, object InfoROMObject) (string, error) {
	buffer := [deviceInfoROMVersionBufferSize]C.char{}
	if err := a.call(a.nvmlDeviceGetInforomVersion, uintptr(device), uintptr(object), uintptr(unsafe.Pointer(&buffer[0])), deviceInfoROMVersionBufferSize); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}

// DeviceGetMaxClockInfo retrieves the maximum clock speeds for the device.
func (a API) DeviceGetMaxClockInfo(device Device, clockType ClockType) (clock uint32, err error) {
	err = a.call(a.nvmlDeviceGetMaxClockInfo, uintptr(device), uintptr(clockType), uintptr(unsafe.Pointer(&clock)))
	return
}

// DeviceGetMaxCustomerBoostClock retrieves the customer defined maximum boost clock speed specified by the given clock type.
func (a API) DeviceGetMaxCustomerBoostClock(device Device, clockType ClockType) (clockMHz uint32, err error) {
	err = a.call(a.nvmlDeviceGetMaxCustomerBoostClock, uintptr(device), uintptr(clockType), uintptr(unsafe.Pointer(&clockMHz)))
	return
}

// DeviceGetMaxPcieLinkGeneration retrieves the maximum PCIe link generation possible with this device and system.
// I.E. for a generation 2 PCIe device attached to a generation 1 PCIe bus the max link generation this function will
// report is generation 1.
func (a API) DeviceGetMaxPcieLinkGeneration(device Device) (maxLinkGen uint32, err error) {
	err = a.call(a.nvmlDeviceGetMaxPcieLinkGeneration, uintptr(device), uintptr(unsafe.Pointer(&maxLinkGen)))
	return
}

// DeviceGetMaxPcieLinkWidth retrieves the maximum PCIe link width possible with this device and system
// I.E. for a device with a 16x PCIe bus width attached to a 8x PCIe system bus this function will report a max link width of 8.
func (a API) DeviceGetMaxPcieLinkWidth(device Device) (maxLinkWidth uint32, err error) {
	err = a.call(a.nvmlDeviceGetMaxPcieLinkWidth, uintptr(device), uintptr(unsafe.Pointer(&maxLinkWidth)))
	return
}

// DeviceGetMemoryErrorCounter retrieves the requested memory error counter for the device.
// Requires NVML_INFOROM_ECC version 2.0 or higher to report aggregate location-based memory error counts.
// Requires NVML_INFOROM_ECC version 1.0 or higher to report all other memory error counts.
// Only applicable to devices with ECC. Requires ECC Mode to be enabled.
func (a API) DeviceGetMemoryErrorCounter(device Device, errorType MemoryErrorType, counterType ECCCounterType, locationType MemoryLocation) (count uint64, err error) {
	err = a.call(a.nvmlDeviceGetMemoryErrorCounter,
		uintptr(device),
		uintptr(errorType),
		uintptr(counterType),
		uintptr(locationType),
		uintptr(unsafe.Pointer(&count)))
	return
}

// DeviceGetMemoryInfo retrieves the amount of used, free and total memory available on the device, in bytes.
func (a API) DeviceGetMemoryInfo(device Device) (mem Memory, err error) {
	err = a.call(a.nvmlDeviceGetMemoryInfo, uintptr(device), uintptr(unsafe.Pointer(&mem)))
	return
}

// DeviceGetMinorNumber retrieves minor number for the device. The minor number for the device is such that
// the Nvidia device node file for each GPU will have the form /dev/nvidia[minor number].
func (a API) DeviceGetMinorNumber(device Device) (minorNumber uint32, err error) {
	err = a.call(a.nvmlDeviceGetMinorNumber, uintptr(device), uintptr(unsafe.Pointer(&minorNumber)))
	return
}

// DeviceGetMultiGpuBoard retrieves whether the device is on a Multi-GPU Board.
func (a API) DeviceGetMultiGpuBoard(device Device) (multiGpu bool, err error) {
	var multiGpuBool uint
	err = a.call(a.nvmlDeviceGetMultiGpuBoard, uintptr(device), uintptr(unsafe.Pointer(&multiGpuBool)))
	if err != nil {
		return
	}

	// Non-zero value indicates whether the device is on a multi GPU board
	multiGpu = false
	if multiGpuBool != 0 {
		multiGpu = true
	}

	return
}

// DeviceGetName retrieves the name of this device.
func (a API) DeviceGetName(device Device) (string, error) {
	buffer := [deviceNameBufferSize]C.char{}
	if err := a.call(a.nvmlDeviceGetName, uintptr(device), uintptr(unsafe.Pointer(&buffer[0])), deviceNameBufferSize); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}

func (a API) DeviceGetP2PStatus() error {
	return ErrNotImplemented
}

func (a API) DeviceGetPCIInfo(device Device) (*PCIInfo, error) {
	var pci C.nvmlPciInfo_t
	if err := a.call(a.nvmlDeviceGetPciInfo, uintptr(device), uintptr(unsafe.Pointer(&pci))); err != nil {
		return nil, err
	}

	return &PCIInfo{
		BusID:          C.GoString(&pci.busId[0]),
		Domain:         uint32(pci.domain),
		Bus:            uint32(pci.bus),
		Device:         uint32(pci.device),
		PCIDeviceID:    uint32(pci.pciDeviceId),
		PCISubsystemID: uint32(pci.pciSubSystemId),
	}, nil
}

// DeviceGetPcieReplayCounter retrieve the PCIe replay counter.
func (a API) DeviceGetPcieReplayCounter(device Device) (value uint32, err error) {
	err = a.call(a.nvmlDeviceGetPcieReplayCounter, uintptr(device), uintptr(unsafe.Pointer(&value)))
	return
}

// DeviceGetPCIeThroughput eetrieve PCIe utilization information.
// This function is querying a byte counter over a 20ms interval and thus is the PCIe throughput over that interval.
// This method is not supported in virtual machines running virtual GPU (vGPU).
func (a API) DeviceGetPCIeThroughput(device Device, counter PCIeUtilCounter) (value uint32, err error) {
	err = a.call(a.nvmlDeviceGetPcieThroughput, uintptr(device), uintptr(counter), uintptr(unsafe.Pointer(&value)))
	return
}

// DeviceGetPerformanceState retrieves the current performance state for the device.
func (a API) DeviceGetPerformanceState(device Device) (state PState, err error) {
	err = a.call(a.nvmlDeviceGetPerformanceState, uintptr(device), uintptr(unsafe.Pointer(&state)))
	return
}

// DeviceGetPowerManagementDefaultLimit retrieves default power management limit on this device, in milliwatts.
// Default power management limit is a power management limit that the device boots with.
func (a API) DeviceGetPowerManagementDefaultLimit(device Device) (defaultLimit uint32, err error) {
	err = a.call(a.nvmlDeviceGetPowerManagementDefaultLimit, uintptr(device), uintptr(unsafe.Pointer(&defaultLimit)))
	return
}

// DeviceGetPowerManagementLimit retrieves the power management limit associated with this device.
// The power limit defines the upper boundary for the card's power draw.
// If the card's total power draw reaches this limit the power management algorithm kicks in.
// This reading is only available if power management mode is supported, see DeviceGetPowerManagementMode.
func (a API) DeviceGetPowerManagementLimit(device Device) (limit uint32, err error) {
	err = a.call(a.nvmlDeviceGetPowerManagementLimit, uintptr(device), uintptr(unsafe.Pointer(&limit)))
	return
}

// DeviceGetPowerManagementLimitConstraints retrieves information about possible values of power management limits on this device.
func (a API) DeviceGetPowerManagementLimitConstraints(device Device) (minLimit, maxLimit uint32, err error) {
	err = a.call(a.nvmlDeviceGetPowerManagementLimitConstraints, uintptr(device), uintptr(unsafe.Pointer(&minLimit)), uintptr(unsafe.Pointer(&maxLimit)))
	return
}

// DeviceGetPowerManagementMode retrieves the power management mode associated with this device.
// This API has been deprecated.
// This flag indicates whether any power management algorithm is currently active on the device.
// An enabled state does not necessarily mean the device is being actively throttled -- only that that the driver will
// do so if the appropriate conditions are met.
func (a API) DeviceGetPowerManagementMode(device Device) (bool, error) {
	var state int32
	if err := a.call(a.nvmlDeviceGetPowerManagementMode, uintptr(device), uintptr(unsafe.Pointer(&state))); err != nil {
		return false, nil
	}

	if state > 0 {
		return true, nil
	}

	return false, nil
}

// DeviceGetPowerState retrieve the current performance state for the device.
// Deprecated: Use DeviceGetPerformanceState.
// This function exposes an incorrect generalization.
func (a API) DeviceGetPowerState(device Device) (state PState, err error) {
	err = a.call(a.nvmlDeviceGetPowerState, uintptr(device), uintptr(unsafe.Pointer(&state)))
	return
}

// DeviceGetPowerUsage retrieves power usage for this GPU in milliwatts and its associated circuitry (e.g. memory)
func (a API) DeviceGetPowerUsage(device Device) (power uint32, err error) {
	err = a.call(a.nvmlDeviceGetPowerUsage, uintptr(device), uintptr(unsafe.Pointer(&power)))
	return
}

// DeviceGetRetiredPages returns the list of retired pages by source, including pages that are pending retirement.
// The address information provided from this API is the hardware address of the page that was retired.
// Note that this does not match the virtual address used in CUDA, but will match the address information in XID 63
func (a API) DeviceGetRetiredPages(device Device, cause PageRetirementCause) ([]uint64, error) {
	// Get array size
	var count uint32
	err := a.call(a.nvmlDeviceGetRetiredPages, uintptr(device), uintptr(cause), uintptr(unsafe.Pointer(&count)), 0)
	if err == nil {
		return []uint64{}, nil
	}

	if err != ErrInsufficientSize {
		return nil, err
	}

	// Query data
	list := make([]uint64, count)
	err = a.call(a.nvmlDeviceGetRetiredPages,
		uintptr(device),
		uintptr(cause),
		uintptr(unsafe.Pointer(&count)),
		uintptr(unsafe.Pointer(&list[0])))

	if err != nil {
		return nil, err
	}

	return list, nil
}

// DeviceGetRetiredPagesPendingStatus checks if any pages are pending retirement and need a reboot to fully retire.
func (a API) DeviceGetRetiredPagesPendingStatus(device Device) (isPending bool, err error) {
	var state int32 = 0
	err = a.call(a.nvmlDeviceGetRetiredPagesPendingStatus, uintptr(device), uintptr(unsafe.Pointer(&state)))
	if err != nil {
		return
	}

	if state > 0 {
		isPending = true
	} else {
		isPending = false
	}

	return
}

func (a API) DeviceGetSamples() error {
	return ErrNotImplemented
}

// DeviceGetSerial retrieves the globally unique board serial number associated with this device's board.
func (a API) DeviceGetSerial(device Device) (serial string, err error) {
	buffer := [deviceSerialBufferSize]C.char{}
	err = a.call(a.nvmlDeviceGetSerial, uintptr(device), uintptr(unsafe.Pointer(&buffer[0])), deviceSerialBufferSize)
	return
}

// DeviceGetSupportedClocksThrottleReasons retrieves bitmask of supported clocks throttle reasons that can be
// returned by DeviceGetCurrentClocksThrottleReasons. This method is not supported in virtual machines
// running virtual GPU (vGPU).
func (a API) DeviceGetSupportedClocksThrottleReasons(device Device) (supportedClocksThrottleReasons ClocksThrottleReason, err error) {
	err = a.call(a.nvmlDeviceGetSupportedClocksThrottleReasons, uintptr(device), uintptr(unsafe.Pointer(&supportedClocksThrottleReasons)))
	return
}

// DeviceGetSupportedGraphicsClocks retrieves the list of possible graphics clocks that can be used
// as an argument for DeviceSetApplicationsClocks.
func (a API) DeviceGetSupportedGraphicsClocks(device Device, memoryClockMHz uint32) ([]uint32, error) {
	// Get array size
	var count uint32
	err := a.call(a.nvmlDeviceGetSupportedGraphicsClocks, uintptr(device), uintptr(memoryClockMHz), uintptr(unsafe.Pointer(&count)), 0)
	if err == nil {
		return []uint32{}, nil
	}

	if err != ErrInsufficientSize {
		return nil, err
	}

	// Query data
	list := make([]uint32, count)
	if err := a.call(a.nvmlDeviceGetSupportedGraphicsClocks, uintptr(device), uintptr(memoryClockMHz), uintptr(unsafe.Pointer(&count)), uintptr(unsafe.Pointer(&list[0]))); err != nil {
		return nil, err
	}

	return list, nil
}

// DeviceGetSupportedMemoryClocks retrieves the list of possible memory clocks that can be used
// as an argument for DeviceSetApplicationsClocks.
func (a API) DeviceGetSupportedMemoryClocks(device Device) ([]uint32, error) {
	// Get array size
	var count uint32

	err := a.call(a.nvmlDeviceGetSupportedMemoryClocks, uintptr(device), uintptr(unsafe.Pointer(&count)), 0)
	if err == nil {
		return []uint32{}, nil
	}

	if err != ErrInsufficientSize {
		return nil, err
	}

	// Query data
	list := make([]uint32, count)
	if err := a.call(a.nvmlDeviceGetSupportedMemoryClocks, uintptr(device), uintptr(unsafe.Pointer(&count)), uintptr(unsafe.Pointer(&list[0]))); err != nil {
		return nil, err
	}

	return list, nil
}

// DeviceGetTemperature retrieves the current temperature readings for the device, in degrees C.
func (a API) DeviceGetTemperature(device Device, sensorType TemperatureSensor) (temp uint32, err error) {
	err = a.call(a.nvmlDeviceGetTemperature, uintptr(device), uintptr(sensorType), uintptr(unsafe.Pointer(&temp)))
	return
}

// DeviceGetTemperatureThreshold retrieves the temperature threshold for the GPU with the specified threshold type in degrees C.
func (a API) DeviceGetTemperatureThreshold(device Device, thresholdType TemperatureThreshold) (temp uint32, err error) {
	err = a.call(a.nvmlDeviceGetTemperatureThreshold, uintptr(device), uintptr(thresholdType), uintptr(unsafe.Pointer(&temp)))
	return
}

// DeviceGetTopologyCommonAncestor retrieves the common ancestor for two devices. Supported on Linux only.
func (a API) DeviceGetTopologyCommonAncestor(device1 Device, device2 Device) (pathInfo GPUTopologyLevel, err error) {
	err = a.call(a.nvmlDeviceGetTopologyCommonAncestor, uintptr(device1), uintptr(device2), uintptr(unsafe.Pointer(&pathInfo)))
	return
}

func (a API) DeviceGetTopologyNearestGpus() error {
	return ErrNotImplemented
}

// DeviceGetTotalECCErrors retrieves the total ECC error counts for the device.
// Only applicable to devices with ECC. Requires NVML_INFOROM_ECC version 1.0 or higher. Requires ECC Mode to be enabled.
// The total error count is the sum of errors across each of the separate memory systems, i.e. the total set of errors across the entire device.
func (a API) DeviceGetTotalECCErrors(device Device, errorType MemoryErrorType, counterType ECCCounterType) (eccCount uint64, err error) {
	err = a.call(a.nvmlDeviceGetTotalEccErrors, uintptr(device), uintptr(errorType), uintptr(counterType), uintptr(unsafe.Pointer(&eccCount)))
	return
}

// DeviceGetTotalEnergyConsumption retrieves total energy consumption for this GPU in millijoules (mJ)
// since the driver was last reloaded.
func (a API) DeviceGetTotalEnergyConsumption(device Device) (energy uint64, err error) {
	err = a.call(a.nvmlDeviceGetTotalEnergyConsumption, uintptr(device), uintptr(unsafe.Pointer(&energy)))
	return
}

// DeviceGetUUID retrieves the globally unique immutable UUID associated with this device,
// as a 5 part hexadecimal string, that augments the immutable, board serial identifier.
func (a API) DeviceGetUUID(device Device) (string, error) {
	buffer := [deviceUUIDBufferSize]C.char{}
	if err := a.call(a.nvmlDeviceGetUUID, uintptr(device), uintptr(unsafe.Pointer(&buffer[0])), deviceUUIDBufferSize); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}

// DeviceGetUtilizationRates retrieves the current utilization rates for the device's major subsystems.
func (a API) DeviceGetUtilizationRates(device Device) (u Utilization, err error) {
	u.GPU = 0
	u.Memory = 0
	err = a.call(a.nvmlDeviceGetUtilizationRates, uintptr(device), uintptr(unsafe.Pointer(&u)))
	return
}

// DeviceGetVbiosVersion gets VBIOS version of the device. The VBIOS version may change from time to time.
func (a API) DeviceGetVbiosVersion(device Device) (string, error) {
	buffer := [deviceVBIOSVersionBufferSize]C.char{}
	if err := a.call(a.nvmlDeviceGetVbiosVersion, uintptr(device), uintptr(unsafe.Pointer(&buffer[0])), deviceVBIOSVersionBufferSize); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}

// DeviceGetViolationStatus gets the duration of time during which the device was throttled (lower than requested
// clocks) due to power or thermal constraints.
// The method is important to users who are tying to understand if their GPUs throttle at any point during their
// applications. The difference in violation times at two different reference times gives the indication of
// GPU throttling event.
func (a API) DeviceGetViolationStatus(device Device, policyType PerfPolicyType) (violTime ViolationTime, err error) {
	err = a.call(a.nvmlDeviceGetViolationStatus, uintptr(device), uintptr(policyType), uintptr(unsafe.Pointer(&violTime)))
	return
}

// DeviceOnSameBoard checks if the GPU devices are on the same physical board.
func (a API) DeviceOnSameBoard(device1 Device, device2 Device) (bool, error) {
	var onSameBoard int32 = 0

	if err := a.call(a.nvmlDeviceOnSameBoard, uintptr(device1), uintptr(device2), uintptr(unsafe.Pointer(&onSameBoard))); err != nil {
		return false, err
	}

	if onSameBoard == 0 {
		return false, nil
	}

	return true, nil
}

// DeviceResetApplicationsClocks resets the application clock to the default value.
func (a API) DeviceResetApplicationsClocks(device Device) error {
	return a.call(a.nvmlDeviceResetApplicationsClocks, uintptr(device))
}

// DeviceSetAutoBoostedClocksEnabled tries to set the current state of Auto Boosted clocks on a device.
// Auto Boosted clocks are enabled by default on some hardware, allowing the GPU to run at higher clock rates to
// maximize performance as thermal limits allow. Auto Boosted clocks should be disabled if fixed clock rates are desired.
// Non-root users may use this API by default but can be restricted by root from using this API by calling
// DeviceSetAPIRestriction with apiType=NVML_RESTRICTED_API_SET_AUTO_BOOSTED_CLOCKS.
// Note: Persistence Mode is required to modify current Auto Boost settings, therefore, it must be enabled.
func (a API) DeviceSetAutoBoostedClocksEnabled(device Device, enabled bool) error {
	var state int32 = 0
	if enabled {
		state = 1
	}

	return a.call(a.nvmlDeviceSetAutoBoostedClocksEnabled, uintptr(device), uintptr(state))
}

// DeviceSetDefaultAutoBoostedClocksEnabled tries to set the default state of Auto Boosted clocks on a device.
// This is the default state that Auto Boosted clocks will return to when no compute running processes (e.g. CUDA
// application which have an active context) are running.
func (a API) DeviceSetDefaultAutoBoostedClocksEnabled(device Device, enabled bool) error {
	var state int32 = 0
	if enabled {
		state = 1
	}

	return a.call(a.nvmlDeviceSetDefaultAutoBoostedClocksEnabled, uintptr(device), uintptr(state), 0)
}

// DeviceValidateInforom reads the infoROM from the flash and verifies the checksums.
func (a API) DeviceValidateInforom(device Device) (err error) {
	err = a.call(a.nvmlDeviceValidateInforom, uintptr(device))
	return
}
