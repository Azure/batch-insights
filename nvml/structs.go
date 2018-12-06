package nvml

import "math"

// Device represents native NVML device handle.
type Device uintptr

const (
	systemDriverVersionBufferSize  = 80
	deviceNameBufferSize           = 64
	deviceSerialBufferSize         = 30
	deviceUUIDBufferSize           = 80
	deviceVBIOSVersionBufferSize   = 32
	deviceInfoROMVersionBufferSize = 16
)

// Memory holds allocation information for a device.
type Memory struct {
	Total uint64 // Total installed FB memory (in bytes).
	Free  uint64 // Unallocated FB memory (in bytes).
	Used  uint64 // Allocated FB memory (in bytes).
}

// BAR1Memory holds BAR1 memory allocation information for a device.
type BAR1Memory struct {
	Total uint64 // Total BAR1 Memory (in bytes)
	Free  uint64 // Unallocated BAR1 Memory (in bytes)
	Used  uint64 // Allocated Used Memory (in bytes)
}

// Temperature sensors.
type TemperatureSensor int32

//noinspection GoUnusedConst
const (
	TemperatureGPU = TemperatureSensor(0) // Temperature sensor for the GPU die.
)

// Temperature thresholds.
type TemperatureThreshold int32

//noinspection GoUnusedConst
const (
	// Temperature at which the GPU will shut down for HW protection
	TemperatureThresholdShutdown = TemperatureThreshold(0)
	// Temperature at which the GPU will begin HW slowdown
	TemperatureThresholdSlowdown = TemperatureThreshold(1)
	// Memory Temperature at which the GPU will begin SW slowdown
	TemperatureThresholdMemMax = TemperatureThreshold(2)
	// GPU Temperature at which the GPU can be throttled below base clock
	TemperatureThresholdGPUMax = TemperatureThreshold(3)
)

// Clock types. All speeds are in Mhz.
type ClockType int32

//noinspection GoUnusedConst
const (
	ClockGraphics = ClockType(0) // Graphics clock domain
	ClockSM       = ClockType(1) // SM clock domain
	ClockMem      = ClockType(2) // Memory clock domain
	ClockVideo    = ClockType(3) // Video encoder/decoder clock domain
)

// ProcessInfo holds information about running compute processes on the GPU.
type ProcessInfo struct {
	// Process ID
	PID uint32
	// Amount of used GPU memory in bytes. Under WDDM, NVML_VALUE_NOT_AVAILABLE is always reported because Windows KMD
	// manages all the memory and not the NVIDIA driver.
	UsedGPUMemory uint64
}

func (i ProcessInfo) MemoryInfoAvailable() bool {
	return i.UsedGPUMemory != math.MaxUint64
}

// Utilization information for a device.
// Each sample period may be between 1 second and 1/6 second, depending on the product being queried.
type Utilization struct {
	GPU    uint32 // Percent of time over the past sample period during which one or more kernels was executing on the GPU.
	Memory uint32 // Percent of time over the past sample period during which global (device) memory was being read or written.
}

// The Brand of the GPU.
type BrandType int32

//noinspection GoUnusedConst
const (
	BrandUnknown = BrandType(0)
	BrandQuadro  = BrandType(1)
	BrandTesla   = BrandType(2)
	BrandNVS     = BrandType(3)
	BrandGrid    = BrandType(4)
	BrandGeforce = BrandType(5)
)

func (b BrandType) String() string {
	switch b {
	case BrandQuadro:
		return "Quadro"
	case BrandTesla:
		return "Tesla"
	case BrandNVS:
		return "NVS"
	case BrandGrid:
		return "Grid"
	case BrandGeforce:
		return "Geforce"
	default:
		return "Unknown"
	}
}

// Clock Ids. These are used in combination with ClockType to specify a single clock value.
type ClockID int32

//noinspection GoUnusedConst
const (
	ClockIDCurrent          = ClockID(0) // Current actual clock value.
	ClockIDAppClockTarget   = ClockID(1) // Target application clock.
	ClockIDAppClockDefault  = ClockID(2) // Default application clock target.
	ClockIDCustomerBoostMax = ClockID(3) // OEM-defined maximum clock rate.
)

type ClocksThrottleReason uint64

//noinspection GoUnusedConst
const (
	// Bit mask representing no clocks throttling. Clocks are as high as possible.
	ClocksThrottleReasonNone = ClocksThrottleReason(0)
	// Nothing is running on the GPU and the clocks are dropping to Idle state.
	ClocksThrottleReasonGPUIdle = ClocksThrottleReason(0x0000000000000001)
	// GPU clocks are limited by current setting of applications clocks.
	ClocksThrottleReasonApplicationsClocksSetting = ClocksThrottleReason(0x0000000000000002)
	// Renamed to ClocksThrottleReasonApplicationsClocksSetting as the name describes the situation more accurately.
	ClocksThrottleReasonUserDefinedClocks = ClocksThrottleReason(0x0000000000000002)
	// SW Power Scaling algorithm is reducing the clocks below requested clocks.
	ClocksThrottleReasonSWPowerCap = ClocksThrottleReason(0x0000000000000004)
	// HW Slowdown (reducing the core clocks by a factor of 2 or more) is engaged.
	// This is an indicator of:
	// 	 - Temperature being too high
	//   - External Power Brake Assertion is triggered (e.g. by the system power supply)
	//   - Power draw is too high and Fast Trigger protection is reducing the clocks
	//   - May be also reported during PState or clock change
	//     - This behavior may be removed in a later release.
	ClocksThrottleReasonHWSlowdown = ClocksThrottleReason(0x0000000000000008)
	// This GPU has been added to a Sync boost group with nvidia-smi or DCGM in
	// order to maximize performance per watt. All GPUs in the sync boost group
	// will boost to the minimum possible clocks across the entire group. Look at
	// the throttle reasons for other GPUs in the system to see why those GPUs are
	// holding this one at lower clocks.
	ClocksThrottleReasonSyncBoost = ClocksThrottleReason(0x0000000000000010)
	// SW Thermal Slowdown
	// This is an indicator of one or more of the following:
	//   - Current GPU temperature above the GPU Max Operating Temperature
	//   - Current memory temperature above the Memory Max Operating Temperature
	ClocksThrottleReasonSWThermalSlowdown = ClocksThrottleReason(0x0000000000000020)
	// HW Thermal Slowdown (reducing the core clocks by a factor of 2 or more) is engaged.
	// This is an indicator of:
	//   - Temperature being too high
	ClocksThrottleReasonHwThermalSlowdown = ClocksThrottleReason(0x0000000000000040)
	// HW Power Brake Slowdown (reducing the core clocks by a factor of 2 or more) is engaged.
	// This is an indicator of:
	//   - External Power Brake Assertion being triggered (e.g. by the system power supply)
	ClocksThrottleReasonHwPowerBrakeSlowdown = ClocksThrottleReason(0x0000000000000080)
)

// GPUOperationMode represents GPU Operation Mode.
// GOM allows to reduce power usage and optimize GPU throughput by disabling GPU features.
// Each GOM is designed to meet specific user needs.
type GPUOperationMode int32

//noinspection GoUnusedConst
const (
	// Everything is enabled and running at full speed.
	GPUOperationModeAllOn = GPUOperationMode(0)
	// Designed for running only compute tasks. Graphics operations are not allowed
	GPUOperationModeCompute = GPUOperationMode(1)
	// Designed for running graphics applications that don't require high bandwidth double precision
	GPUOperationModeLowDoublePrecision = GPUOperationMode(2)
)

// PCIInfo represents PCI information about a GPU device.
type PCIInfo struct {
	BusID string
	// The legacy tuple domain:bus:device.function PCI identifier
	BusIDLegacy string
	// The PCI domain on which the device's bus resides, 0 to 0xffffffff
	Domain uint32
	// The bus on which the device resides, 0 to 0xff
	Bus uint32
	// The device's id on the bus, 0 to 31
	Device uint32
	// The combined 16-bit device id and 16-bit vendor id
	PCIDeviceID uint32
	// The 32-bit Sub System Device ID. Added in NVML 2.285 API
	PCISubsystemID uint32
}

// Driver models. Windows only.
type DriverModel int32

//noinspection GoUnusedConst
const (
	// WDDM driver model -- GPU treated as a display device.
	DriverModelWDDM = DriverModel(0)
	// WDM (TCC) model (recommended) -- GPU treated as a generic device.
	DriverModelWDM = DriverModel(1)
)

// Compute mode.
type ComputeMode int32

//noinspection GoUnusedConst
const (
	// Default compute mode - multiple contexts per device.
	ComputeModeDefault = ComputeMode(0)
	// Support Removed.
	ComputeModeExclusiveThread = ComputeMode(1)
	// No contexts per device.
	ComputeModeProhibited = ComputeMode(2)
	// Only one context per device, usable from multiple threads at a time.
	ComputeModeExclusiveProcess = ComputeMode(3)
)

// API types that allow changes to default permission restrictions.
type RestrictedAPI int32

//noinspection GoUnusedConst
const (
	// APIs that change application clocks
	RestrictedAPISetApplicationClocks = RestrictedAPI(0)
	// APIs that enable/disable Auto Boosted clocks
	RestrictedAPISetAutoBoostedClocks = RestrictedAPI(1)
)

// Available infoROM objects.
type InfoROMObject int32

//noinspection GoUnusedConst
const (
	// An object defined by OEM.
	InfoROMObjectOEM = InfoROMObject(0)
	// The ECC object determining the level of ECC support.
	InfoROMObjectECC = InfoROMObject(1)
	// The power management object.
	InfoROMObjectPower = InfoROMObject(2)
)

// Represents type of encoder for capacity can be queried.
type EncoderType int32

//noinspection GoUnusedConst
const (
	EncoderTypeQueryH264 = EncoderType(0)
	EncoderTypeQueryHEVC = EncoderType(1)
)

// Memory error types
type MemoryErrorType int32

//noinspection GoUnusedConst
const (
	// A memory error that was corrected for ECC errors, these are single bit errors For Texture memory, these are errors fixed by resend.
	MemoryErrorTypeCorrected = MemoryErrorType(0)
	// A memory error that was not corrected for ECC errors, these are double bit errors For Texture memory, these are errors where the resend fails.
	MemoryErrorTypeUncorrected = MemoryErrorType(1)
)

// ECC counter types.
// Note: Volatile counts are reset each time the driver loads.
// On Windows this is once per boot. On Linux this can be more frequent.
// On Linux the driver unloads when no active clients exist.
// If persistence mode is enabled or there is always a driver client active (e.g. X11), then Linux also sees per-boot
// behavior. If not, volatile counts are reset each time a compute app is run.
type ECCCounterType int32

//noinspection GoUnusedConst
const (
	// Volatile counts are reset each time the driver loads.
	VolatileECC = ECCCounterType(0)
	// Aggregate counts persist across reboots (i.e. for the lifetime of the device).
	AggregateECC = ECCCounterType(1)
)

// Memory locations.
type MemoryLocation int32

//noinspection GoUnusedConst
const (
	// GPU L1 Cache.
	MemoryLocationL1Cache = MemoryLocation(0)
	// GPU L2 Cache.
	MemoryLocationL2Cache = MemoryLocation(1)
	// GPU Device Memory.
	MemoryLocationDeviceMemory = MemoryLocation(2)
	// GPU Register File.
	MemoryLocationRegisterFile = MemoryLocation(3)
	// GPU Texture Memory.
	MemoryLocationTextureMemory = MemoryLocation(4)
	// Shared memory.
	MemoryLocationTextureSHM = MemoryLocation(5)
	// CBU.
	MemoryLocationCBU = MemoryLocation(6)
)

// Represents the queryable PCIe utilization counters.
type PCIeUtilCounter int32

//noinspection GoUnusedConst
const (
	PCIeUtilTXBytes = PCIeUtilCounter(0)
	PCIeUtilRXBytes = PCIeUtilCounter(1)
)

// PState represents allowed PStates.
type PState int32

//noinspection GoUnusedConst
const (
	PState0       = PState(0) // Performance state 0 -- Maximum Performance.
	PState1       = PState(1)
	PState2       = PState(2)
	PState3       = PState(3)
	PState4       = PState(4)
	PState5       = PState(5)
	PState6       = PState(6)
	PState7       = PState(7)
	PState8       = PState(8)
	PState9       = PState(9)
	PState10      = PState(10)
	PState11      = PState(11)
	PState12      = PState(12)
	PState13      = PState(13)
	PState14      = PState(14)
	PState15      = PState(15) // Performance state 15 -- Minimum Performance.
	PStateUnknown = PState(32)
)

// Causes for page retirement.
type PageRetirementCause int32

//noinspection GoUnusedConst
const (
	// Page was retired due to multiple single bit ECC error.
	PageRetirementCauseMultipleSingleBitECCErrors = PageRetirementCause(0)
	// Page was retired due to double bit ECC error.
	PageRetirementCauseDoubleBitECCError = PageRetirementCause(1)
)

// Represents level relationships within a system between two GPUs.
// The enums are spaced to allow for future relationships.
type GPUTopologyLevel int32

//noinspection GoUnusedConst
const (
	TopologyInternal   = GPUTopologyLevel(0)
	TopologySingle     = GPUTopologyLevel(10)
	TopologyMultiple   = GPUTopologyLevel(20)
	TopologyHostbridge = GPUTopologyLevel(30)
	TopologyNode       = GPUTopologyLevel(40)
	TopologySystem     = GPUTopologyLevel(50)
)

// Detailed ECC error counts for a device.
// Different GPU families can have different memory error counters.
type ECCErrorCounts struct {
	L1Cache      uint64 // L1 cache errors.
	L2Cache      uint64 // L2 cache errors.
	DeviceMemory uint64 // Device memory errors.
	RegisterFile uint64 // Register file errors.
}

// Represents type of perf policy for which violation times can be queried.
type PerfPolicyType int32

//noinspection GoUnusedConst
const (
	// How long did power violations cause the GPU to be below application clocks.
	PerfPolicyPower = PerfPolicyType(0)
	// How long did thermal violations cause the GPU to be below application clocks.
	PerfPolicyThermal = PerfPolicyType(1)
	// How long did sync boost cause the GPU to be below application clocks.
	PerfPolicySyncBoost = PerfPolicyType(2)
	// How long did the board limit cause the GPU to be below application clocks.
	PerfPolicyBoardLimit = PerfPolicyType(3)
	// How long did low utilization cause the GPU to be below application clocks.
	PerfPolicyLowUtilization = PerfPolicyType(4)
	// How long did the board reliability limit cause the GPU to be below application clocks.
	PerfPolicyReliability = PerfPolicyType(5)
	// Total time the GPU was held below application clocks by any limiter (0 - 5 above).
	PerfPolicyTotalAppClocks = PerfPolicyType(10)
	// Total time the GPU was held below base clocks.
	PerfPolicyTotalBaseClocks = PerfPolicyType(11)
)

// ViolationTime holds perf policy violation status data.
type ViolationTime struct {
	ReferenceTime uint64 // ReferenceTime represents CPU timestamp in microseconds
	ViolationTime uint64 // ViolationTime in Nanoseconds
}
