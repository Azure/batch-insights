package nvml

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrUninitialized        = errors.New("NVML was not first initialized with Init")
	ErrInvalidArgument      = errors.New("A supplied argument is invalid")
	ErrNotSupported         = errors.New("The requested operation is not available on target device")
	ErrNoPermission         = errors.New("The current user does not have permission for operation")
	ErrAlreadyInititlized   = errors.New("Multiple initializations are now allowed through ref counting")
	ErrNotFound             = errors.New("A query to find an object was unsuccessful")
	ErrInsufficientSize     = errors.New("An input argument is not large enough")
	ErrInsufficientPower    = errors.New("A device's external power cables are not properly attached")
	ErrDriverNotLoaded      = errors.New("NVIDIA driver is not loaded")
	ErrTimeout              = errors.New("User provided timeout passed")
	ErrIRQIssue             = errors.New("NVIDIA Kernel detected an interrupt issue with a GPU")
	ErrLibraryNotFound      = errors.New("NVML Shared Library couldn't be found or loaded")
	ErrFunctionNotFound     = errors.New("Local version of NVML doesn't implement this function")
	ErrCorruptedInfoROM     = errors.New("infoROM is corrupted")
	ErrGPULost              = errors.New("The GPU has fallen off the bus or has otherwise become inaccessible")
	ErrResetRequired        = errors.New("The GPU requires a reset before it can be used again")
	ErrOperatingSystem      = errors.New("The GPU control device has been blocked by the operating system/cgroups")
	ErrLibRMVersionMismatch = errors.New("RM detects a driver/library version mismatch")
	ErrInUse                = errors.New("An operation cannot be performed because the GPU is currently in use")
	ErrMemory               = errors.New("Insufficient memory")
	ErrNoData               = errors.New("No data")
	ErrVGPUECCNotSupported  = errors.New("The requested vgpu operation is not available on target device, because ECC is enabled")
	ErrUnknown              = errors.New("An internal driver error occurred")
)

var errorCodeMappings = map[int]error{
	0:   nil,
	1:   ErrUninitialized,
	2:   ErrInvalidArgument,
	3:   ErrNotSupported,
	4:   ErrNoPermission,
	5:   ErrAlreadyInititlized,
	6:   ErrNotFound,
	7:   ErrInsufficientSize,
	8:   ErrInsufficientPower,
	9:   ErrDriverNotLoaded,
	10:  ErrTimeout,
	11:  ErrIRQIssue,
	12:  ErrLibraryNotFound,
	13:  ErrFunctionNotFound,
	14:  ErrCorruptedInfoROM,
	15:  ErrGPULost,
	16:  ErrResetRequired,
	17:  ErrOperatingSystem,
	18:  ErrLibRMVersionMismatch,
	19:  ErrInUse,
	20:  ErrMemory,
	21:  ErrNoData,
	22:  ErrVGPUECCNotSupported,
	999: ErrUnknown,
}

func returnValueToError(code int) error {
	if code == 0 {
		return nil
	}

	err, ok := errorCodeMappings[code]
	if ok {
		return err
	}

	return fmt.Errorf("NVML call failed with error code %d", code)
}
