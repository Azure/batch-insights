package nvml

import (
	"C"
	"os"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

var ErrNotImplemented = errors.New("Not implemented")

type API struct {
	dll *syscall.DLL
	// Initialization and cleanup
	nvmlInit,
	nvmlShutdown,
	// Error reporting
	nvmlErrorString,
	// System Queries
	nvmlSystemGetCudaDriverVersion,
	nvmlSystemGetDriverVersion,
	nvmlSystemGetNVMLVersion,
	nvmlSystemGetProcessName,
	// Device Queries
	nvmlDeviceClearCpuAffinity,
	nvmlDeviceGetAPIRestriction,
	nvmlDeviceGetApplicationsClock,
	nvmlDeviceGetAutoBoostedClocksEnabled,
	nvmlDeviceGetBAR1MemoryInfo,
	nvmlDeviceGetBoardId,
	nvmlDeviceGetBoardPartNumber,
	nvmlDeviceGetBrand,
	nvmlDeviceGetBridgeChipInfo,
	nvmlDeviceGetClock,
	nvmlDeviceGetClockInfo,
	nvmlDeviceGetComputeMode,
	nvmlDeviceGetComputeRunningProcesses,
	nvmlDeviceGetCount,
	nvmlDeviceGetCpuAffinity,
	nvmlDeviceGetCudaComputeCapability,
	nvmlDeviceGetCurrPcieLinkGeneration,
	nvmlDeviceGetCurrPcieLinkWidth,
	nvmlDeviceGetCurrentClocksThrottleReasons,
	nvmlDeviceGetDecoderUtilization,
	nvmlDeviceGetDefaultApplicationsClock,
	nvmlDeviceGetDetailedEccErrors,
	nvmlDeviceGetDisplayActive,
	nvmlDeviceGetDisplayMode,
	nvmlDeviceGetDriverModel,
	nvmlDeviceGetEccMode,
	nvmlDeviceGetEncoderCapacity,
	nvmlDeviceGetEncoderSessions,
	nvmlDeviceGetEncoderStats,
	nvmlDeviceGetEncoderUtilization,
	nvmlDeviceGetEnforcedPowerLimit,
	nvmlDeviceGetFanSpeed,
	nvmlDeviceGetGpuOperationMode,
	nvmlDeviceGetGraphicsRunningProcesses,
	nvmlDeviceGetHandleByIndex,
	nvmlDeviceGetHandleByPciBusId,
	nvmlDeviceGetHandleBySerial,
	nvmlDeviceGetHandleByUUID,
	nvmlDeviceGetIndex,
	nvmlDeviceGetInforomConfigurationChecksum,
	nvmlDeviceGetInforomImageVersion,
	nvmlDeviceGetInforomVersion,
	nvmlDeviceGetMaxClockInfo,
	nvmlDeviceGetMaxCustomerBoostClock,
	nvmlDeviceGetMaxPcieLinkGeneration,
	nvmlDeviceGetMaxPcieLinkWidth,
	nvmlDeviceGetMemoryErrorCounter,
	nvmlDeviceGetMemoryInfo,
	nvmlDeviceGetMinorNumber,
	nvmlDeviceGetMultiGpuBoard,
	nvmlDeviceGetName,
	nvmlDeviceGetP2PStatus,
	nvmlDeviceGetPciInfo,
	nvmlDeviceGetPcieReplayCounter,
	nvmlDeviceGetPcieThroughput,
	nvmlDeviceGetPerformanceState,
	nvmlDeviceGetPersistenceMode,
	nvmlDeviceGetPowerManagementDefaultLimit,
	nvmlDeviceGetPowerManagementLimit,
	nvmlDeviceGetPowerManagementLimitConstraints,
	nvmlDeviceGetPowerManagementMode,
	nvmlDeviceGetPowerState,
	nvmlDeviceGetPowerUsage,
	nvmlDeviceGetRetiredPages,
	nvmlDeviceGetRetiredPagesPendingStatus,
	nvmlDeviceGetSamples,
	nvmlDeviceGetSerial,
	nvmlDeviceGetSupportedClocksThrottleReasons,
	nvmlDeviceGetSupportedGraphicsClocks,
	nvmlDeviceGetSupportedMemoryClocks,
	nvmlDeviceGetTemperature,
	nvmlDeviceGetTemperatureThreshold,
	nvmlDeviceGetTopologyCommonAncestor,
	nvmlDeviceGetTopologyNearestGpus,
	nvmlDeviceGetTotalEccErrors,
	nvmlDeviceGetTotalEnergyConsumption,
	nvmlDeviceGetUUID,
	nvmlDeviceGetUtilizationRates,
	nvmlDeviceGetVbiosVersion,
	nvmlDeviceGetViolationStatus,
	nvmlDeviceOnSameBoard,
	nvmlDeviceResetApplicationsClocks,
	nvmlDeviceSetAutoBoostedClocksEnabled,
	nvmlDeviceSetCpuAffinity,
	nvmlDeviceSetDefaultAutoBoostedClocksEnabled,
	nvmlDeviceValidateInforom,
	nvmlSystemGetTopologyGpuSet,
	// Device commands
	nvmlDeviceClearEccErrorCounts,
	nvmlDeviceSetAPIRestriction,
	nvmlDeviceSetApplicationsClocks,
	nvmlDeviceSetComputeMode,
	nvmlDeviceSetDriverModel,
	nvmlDeviceSetEccMode,
	nvmlDeviceSetGpuOperationMode,
	nvmlDeviceSetPersistenceMode,
	nvmlDeviceSetPowerManagementLimit *syscall.Proc
}

func (a API) call(p *syscall.Proc, args ...uintptr) error {
	ret, _, _ := p.Call(args...)
	if ret != 0 {
		return returnValueToError(int(ret))
	}

	return nil
}

// Init initializes NVML, but don't initialize any GPUs yet.
func (a API) Init() error {
	return a.call(a.nvmlInit)
}

// Shutdown shut downs NVML by releasing all GPU resources previously allocated with Init() and
// unloads nvml.dll via UnloadLibrary call.
func (a API) Shutdown() error {
	err := a.call(a.nvmlShutdown)
	a.ReleaseDLL()
	return err
}

func (a API) ReleaseDLL() error {
	return a.dll.Release()
}

// ErrorString returns a string representation of the error.
func (a API) ErrorString(result uintptr) string {
	ret, _, _ := a.nvmlErrorString.Call(uintptr(result))
	buf := (*C.char)(unsafe.Pointer(ret))
	return C.GoString(buf)
}

// New creates nvml.dll wrapper
func NewNVML(path string) (*API, error) {
	if path == "" {
		path = os.ExpandEnv("$ProgramW6432\\NVIDIA Corporation\\NVSMI\\nvml.dll")
	}

	dll, err := syscall.LoadDLL(path)
	if err != nil {
		return nil, err
	}

	bindings := &API{
		dll:                                          dll,
		nvmlInit:                                     dll.MustFindProc("nvmlInit"),
		nvmlShutdown:                                 dll.MustFindProc("nvmlShutdown"),
		nvmlErrorString:                              dll.MustFindProc("nvmlErrorString"),
		nvmlSystemGetCudaDriverVersion:               dll.MustFindProc("nvmlSystemGetCudaDriverVersion"),
		nvmlSystemGetDriverVersion:                   dll.MustFindProc("nvmlSystemGetDriverVersion"),
		nvmlSystemGetNVMLVersion:                     dll.MustFindProc("nvmlSystemGetNVMLVersion"),
		nvmlSystemGetProcessName:                     dll.MustFindProc("nvmlSystemGetProcessName"),
		nvmlDeviceClearCpuAffinity:                   dll.MustFindProc("nvmlDeviceClearCpuAffinity"),
		nvmlDeviceGetAPIRestriction:                  dll.MustFindProc("nvmlDeviceGetAPIRestriction"),
		nvmlDeviceGetApplicationsClock:               dll.MustFindProc("nvmlDeviceGetApplicationsClock"),
		nvmlDeviceGetAutoBoostedClocksEnabled:        dll.MustFindProc("nvmlDeviceGetAutoBoostedClocksEnabled"),
		nvmlDeviceGetBAR1MemoryInfo:                  dll.MustFindProc("nvmlDeviceGetBAR1MemoryInfo"),
		nvmlDeviceGetBoardId:                         dll.MustFindProc("nvmlDeviceGetBoardId"),
		nvmlDeviceGetBoardPartNumber:                 dll.MustFindProc("nvmlDeviceGetBoardPartNumber"),
		nvmlDeviceGetBrand:                           dll.MustFindProc("nvmlDeviceGetBrand"),
		nvmlDeviceGetBridgeChipInfo:                  dll.MustFindProc("nvmlDeviceGetBridgeChipInfo"),
		nvmlDeviceGetClock:                           dll.MustFindProc("nvmlDeviceGetClock"),
		nvmlDeviceGetClockInfo:                       dll.MustFindProc("nvmlDeviceGetClockInfo"),
		nvmlDeviceGetComputeMode:                     dll.MustFindProc("nvmlDeviceGetComputeMode"),
		nvmlDeviceGetComputeRunningProcesses:         dll.MustFindProc("nvmlDeviceGetComputeRunningProcesses"),
		nvmlDeviceGetCount:                           dll.MustFindProc("nvmlDeviceGetCount"),
		nvmlDeviceGetCpuAffinity:                     dll.MustFindProc("nvmlDeviceGetCpuAffinity"),
		nvmlDeviceGetCudaComputeCapability:           dll.MustFindProc("nvmlDeviceGetCudaComputeCapability"),
		nvmlDeviceGetCurrPcieLinkGeneration:          dll.MustFindProc("nvmlDeviceGetCurrPcieLinkGeneration"),
		nvmlDeviceGetCurrPcieLinkWidth:               dll.MustFindProc("nvmlDeviceGetCurrPcieLinkWidth"),
		nvmlDeviceGetCurrentClocksThrottleReasons:    dll.MustFindProc("nvmlDeviceGetCurrentClocksThrottleReasons"),
		nvmlDeviceGetDecoderUtilization:              dll.MustFindProc("nvmlDeviceGetDecoderUtilization"),
		nvmlDeviceGetDefaultApplicationsClock:        dll.MustFindProc("nvmlDeviceGetDefaultApplicationsClock"),
		nvmlDeviceGetDetailedEccErrors:               dll.MustFindProc("nvmlDeviceGetDetailedEccErrors"),
		nvmlDeviceGetDisplayActive:                   dll.MustFindProc("nvmlDeviceGetDisplayActive"),
		nvmlDeviceGetDisplayMode:                     dll.MustFindProc("nvmlDeviceGetDisplayMode"),
		nvmlDeviceGetDriverModel:                     dll.MustFindProc("nvmlDeviceGetDriverModel"),
		nvmlDeviceGetEccMode:                         dll.MustFindProc("nvmlDeviceGetEccMode"),
		nvmlDeviceGetEncoderCapacity:                 dll.MustFindProc("nvmlDeviceGetEncoderCapacity"),
		nvmlDeviceGetEncoderSessions:                 dll.MustFindProc("nvmlDeviceGetEncoderSessions"),
		nvmlDeviceGetEncoderStats:                    dll.MustFindProc("nvmlDeviceGetEncoderStats"),
		nvmlDeviceGetEncoderUtilization:              dll.MustFindProc("nvmlDeviceGetEncoderUtilization"),
		nvmlDeviceGetEnforcedPowerLimit:              dll.MustFindProc("nvmlDeviceGetEnforcedPowerLimit"),
		nvmlDeviceGetFanSpeed:                        dll.MustFindProc("nvmlDeviceGetFanSpeed"),
		nvmlDeviceGetGpuOperationMode:                dll.MustFindProc("nvmlDeviceGetGpuOperationMode"),
		nvmlDeviceGetGraphicsRunningProcesses:        dll.MustFindProc("nvmlDeviceGetGraphicsRunningProcesses"),
		nvmlDeviceGetHandleByIndex:                   dll.MustFindProc("nvmlDeviceGetHandleByIndex"),
		nvmlDeviceGetHandleByPciBusId:                dll.MustFindProc("nvmlDeviceGetHandleByPciBusId"),
		nvmlDeviceGetHandleBySerial:                  dll.MustFindProc("nvmlDeviceGetHandleBySerial"),
		nvmlDeviceGetHandleByUUID:                    dll.MustFindProc("nvmlDeviceGetHandleByUUID"),
		nvmlDeviceGetIndex:                           dll.MustFindProc("nvmlDeviceGetIndex"),
		nvmlDeviceGetInforomConfigurationChecksum:    dll.MustFindProc("nvmlDeviceGetInforomConfigurationChecksum"),
		nvmlDeviceGetInforomImageVersion:             dll.MustFindProc("nvmlDeviceGetInforomImageVersion"),
		nvmlDeviceGetInforomVersion:                  dll.MustFindProc("nvmlDeviceGetInforomVersion"),
		nvmlDeviceGetMaxClockInfo:                    dll.MustFindProc("nvmlDeviceGetMaxClockInfo"),
		nvmlDeviceGetMaxCustomerBoostClock:           dll.MustFindProc("nvmlDeviceGetMaxCustomerBoostClock"),
		nvmlDeviceGetMaxPcieLinkGeneration:           dll.MustFindProc("nvmlDeviceGetMaxPcieLinkGeneration"),
		nvmlDeviceGetMaxPcieLinkWidth:                dll.MustFindProc("nvmlDeviceGetMaxPcieLinkWidth"),
		nvmlDeviceGetMemoryErrorCounter:              dll.MustFindProc("nvmlDeviceGetMemoryErrorCounter"),
		nvmlDeviceGetMemoryInfo:                      dll.MustFindProc("nvmlDeviceGetMemoryInfo"),
		nvmlDeviceGetMinorNumber:                     dll.MustFindProc("nvmlDeviceGetMinorNumber"),
		nvmlDeviceGetMultiGpuBoard:                   dll.MustFindProc("nvmlDeviceGetMultiGpuBoard"),
		nvmlDeviceGetName:                            dll.MustFindProc("nvmlDeviceGetName"),
		nvmlDeviceGetP2PStatus:                       dll.MustFindProc("nvmlDeviceGetP2PStatus"),
		nvmlDeviceGetPciInfo:                         dll.MustFindProc("nvmlDeviceGetPciInfo"),
		nvmlDeviceGetPcieReplayCounter:               dll.MustFindProc("nvmlDeviceGetPcieReplayCounter"),
		nvmlDeviceGetPcieThroughput:                  dll.MustFindProc("nvmlDeviceGetPcieThroughput"),
		nvmlDeviceGetPerformanceState:                dll.MustFindProc("nvmlDeviceGetPerformanceState"),
		nvmlDeviceGetPersistenceMode:                 dll.MustFindProc("nvmlDeviceGetPersistenceMode"),
		nvmlDeviceGetPowerManagementDefaultLimit:     dll.MustFindProc("nvmlDeviceGetPowerManagementDefaultLimit"),
		nvmlDeviceGetPowerManagementLimit:            dll.MustFindProc("nvmlDeviceGetPowerManagementLimit"),
		nvmlDeviceGetPowerManagementLimitConstraints: dll.MustFindProc("nvmlDeviceGetPowerManagementLimitConstraints"),
		nvmlDeviceGetPowerManagementMode:             dll.MustFindProc("nvmlDeviceGetPowerManagementMode"),
		nvmlDeviceGetPowerState:                      dll.MustFindProc("nvmlDeviceGetPowerState"),
		nvmlDeviceGetPowerUsage:                      dll.MustFindProc("nvmlDeviceGetPowerUsage"),
		nvmlDeviceGetRetiredPages:                    dll.MustFindProc("nvmlDeviceGetRetiredPages"),
		nvmlDeviceGetRetiredPagesPendingStatus:       dll.MustFindProc("nvmlDeviceGetRetiredPagesPendingStatus"),
		nvmlDeviceGetSamples:                         dll.MustFindProc("nvmlDeviceGetSamples"),
		nvmlDeviceGetSerial:                          dll.MustFindProc("nvmlDeviceGetSerial"),
		nvmlDeviceGetSupportedClocksThrottleReasons:  dll.MustFindProc("nvmlDeviceGetSupportedClocksThrottleReasons"),
		nvmlDeviceGetSupportedGraphicsClocks:         dll.MustFindProc("nvmlDeviceGetSupportedGraphicsClocks"),
		nvmlDeviceGetSupportedMemoryClocks:           dll.MustFindProc("nvmlDeviceGetSupportedMemoryClocks"),
		nvmlDeviceGetTemperature:                     dll.MustFindProc("nvmlDeviceGetTemperature"),
		nvmlDeviceGetTemperatureThreshold:            dll.MustFindProc("nvmlDeviceGetTemperatureThreshold"),
		nvmlDeviceGetTopologyCommonAncestor:          dll.MustFindProc("nvmlDeviceGetTopologyCommonAncestor"),
		nvmlDeviceGetTopologyNearestGpus:             dll.MustFindProc("nvmlDeviceGetTopologyNearestGpus"),
		nvmlDeviceGetTotalEccErrors:                  dll.MustFindProc("nvmlDeviceGetTotalEccErrors"),
		nvmlDeviceGetTotalEnergyConsumption:          dll.MustFindProc("nvmlDeviceGetTotalEnergyConsumption"),
		nvmlDeviceGetUUID:                            dll.MustFindProc("nvmlDeviceGetUUID"),
		nvmlDeviceGetUtilizationRates:                dll.MustFindProc("nvmlDeviceGetUtilizationRates"),
		nvmlDeviceGetVbiosVersion:                    dll.MustFindProc("nvmlDeviceGetVbiosVersion"),
		nvmlDeviceGetViolationStatus:                 dll.MustFindProc("nvmlDeviceGetViolationStatus"),
		nvmlDeviceOnSameBoard:                        dll.MustFindProc("nvmlDeviceOnSameBoard"),
		nvmlDeviceResetApplicationsClocks:            dll.MustFindProc("nvmlDeviceResetApplicationsClocks"),
		nvmlDeviceSetAutoBoostedClocksEnabled:        dll.MustFindProc("nvmlDeviceSetAutoBoostedClocksEnabled"),
		nvmlDeviceSetCpuAffinity:                     dll.MustFindProc("nvmlDeviceSetCpuAffinity"),
		nvmlDeviceSetDefaultAutoBoostedClocksEnabled: dll.MustFindProc("nvmlDeviceSetDefaultAutoBoostedClocksEnabled"),
		nvmlDeviceValidateInforom:                    dll.MustFindProc("nvmlDeviceValidateInforom"),
		nvmlSystemGetTopologyGpuSet:                  dll.MustFindProc("nvmlSystemGetTopologyGpuSet"),
		nvmlDeviceClearEccErrorCounts:                dll.MustFindProc("nvmlDeviceClearEccErrorCounts"),
		nvmlDeviceSetAPIRestriction:                  dll.MustFindProc("nvmlDeviceSetAPIRestriction"),
		nvmlDeviceSetApplicationsClocks:              dll.MustFindProc("nvmlDeviceSetApplicationsClocks"),
		nvmlDeviceSetComputeMode:                     dll.MustFindProc("nvmlDeviceSetComputeMode"),
		nvmlDeviceSetDriverModel:                     dll.MustFindProc("nvmlDeviceSetDriverModel"),
		nvmlDeviceSetEccMode:                         dll.MustFindProc("nvmlDeviceSetEccMode"),
		nvmlDeviceSetGpuOperationMode:                dll.MustFindProc("nvmlDeviceSetGpuOperationMode"),
		nvmlDeviceSetPersistenceMode:                 dll.MustFindProc("nvmlDeviceSetPersistenceMode"),
		nvmlDeviceSetPowerManagementLimit:            dll.MustFindProc("nvmlDeviceSetPowerManagementLimit"),
	}

	return bindings, nil
}
