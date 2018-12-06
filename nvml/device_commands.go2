package nvml

// DeviceClearECCErrorCounts clears the ECC error and other memory error counts for the device.
// Only applicable to devices with ECC. Requires NVML_INFOROM_ECC version 2.0 or higher to clear aggregate
// location-based ECC counts. Requires NVML_INFOROM_ECC version 1.0 or higher to clear all other ECC counts.
// Requires root/admin permissions. Requires ECC Mode to be enabled.
// Sets all of the specified ECC counters to 0, including both detailed and total counts.
// This operation takes effect immediately.
func (a API) DeviceClearECCErrorCounts(device Device, counterType ECCCounterType) error {
	return a.call(a.nvmlDeviceClearEccErrorCounts, uintptr(device), uintptr(counterType))
}

// DeviceSetAPIRestriction changes the root/admin restructions on certain APIs.
// See nvmlRestrictedAPI_t for the list of supported APIs.
// This method can be used by a root/admin user to give non-root/admin access to certain otherwise-restricted APIs.
// The new setting lasts for the lifetime of the NVIDIA driver; it is not persistent.
// See DeviceGetAPIRestriction to query the current restriction settings.
func (a API) DeviceSetAPIRestriction(device Device, apiType RestrictedAPI, isRestricted bool) error {
	var isRestrictedInt int32 = 0
	if isRestricted {
		isRestrictedInt = 1
	}

	return a.call(a.nvmlDeviceSetAPIRestriction, uintptr(device), uintptr(apiType), uintptr(isRestrictedInt))
}

// DeviceSetApplicationsClocks set clocks that applications will lock to.
// Sets the clocks that compute and graphics applications will be running at. e.g. CUDA driver requests these clocks
// during context creation which means this property defines clocks at which CUDA applications will be running unless
// some overspec event occurs (e.g. over power, over thermal or external HW brake).
// Can be used as a setting to request constant performance.
// On Pascal and newer hardware, this will automatically disable automatic boosting of clocks.
// On K80 and newer Kepler and Maxwell GPUs, users desiring fixed performance should also call
// DeviceSetAutoBoostedClocksEnabled to prevent clocks from automatically boosting above the clock value being set.
// After system reboot or driver reload applications clocks go back to their default value.
func (a API) DeviceSetApplicationsClocks(device Device, memClockMHz, graphicsClockMHz uint32) error {
	return a.call(a.nvmlDeviceSetApplicationsClocks, uintptr(device), uintptr(memClockMHz), uintptr(graphicsClockMHz))
}

// DeviceSetComputeMode sets the compute mode for the device.
// Requires root/admin permissions.
// The compute mode determines whether a GPU can be used for compute operations and whether it can be shared across contexts.
// This operation takes effect immediately.
// Under Linux it is not persistent across reboots and always resets to "Default". Under windows it is persistent.
// Under windows compute mode may only be set to DEFAULT when running in WDDM.
func (a API) DeviceSetComputeMode(device Device, mode ComputeMode) error {
	return a.call(a.nvmlDeviceSetComputeMode, uintptr(device), uintptr(mode))
}

// DeviceSetDriverModel sets the driver model for the device.
// For windows only. Requires root/admin permissions.
// On Windows platforms the device driver can run in either WDDM or WDM (TCC) mode.
// If a display is attached to the device it must run in WDDM mode.
// It is possible to force the change to WDM (TCC) while the display is still attached with a force flag (nvmlFlagForce).
// This should only be done if the host is subsequently powered down and the display is detached from the device before
// the next reboot.
// This operation takes effect after the next reboot.
// Windows driver model may only be set to WDDM when running in DEFAULT compute mode. Change driver model to WDDM is not
// supported when GPU doesn't support graphics acceleration or will not support it after reboot.
func (a API) DeviceSetDriverModel(device Device, model DriverModel, flags uint32) error {
	return a.call(a.nvmlDeviceSetDriverModel, uintptr(device), uintptr(model), uintptr(flags))
}

// DeviceSetECCMode sets the ECC mode for the device.
// Only applicable to devices with ECC. Requires NVML_INFOROM_ECC version 1.0 or higher.
// Requires root/admin permissions.
// The ECC mode determines whether the GPU enables its ECC support.
// This operation takes effect after the next reboot.
func (a API) DeviceSetECCMode(device Device, ecc bool) error {
	var eccInt int32 = 0
	if ecc {
		eccInt = 1
	}

	return a.call(a.nvmlDeviceSetEccMode, uintptr(device), uintptr(eccInt))
}

// DeviceSetGPUOperationMode sets new GOM. See nvmlGpuOperationMode_t for details.
// For GK110 M-class and X-class Tesla products from the Kepler family.
// Modes NVML_GOM_LOW_DP and NVML_GOM_ALL_ON are supported on fully supported GeForce products.
// Not supported on Quadro and Tesla C-class products.
// Requires root/admin permissions.
// Changing GOMs requires a reboot. The reboot requirement might be removed in the future.
// Compute only GOMs don't support graphics acceleration.
// Under windows switching to these GOMs when pending driver model is WDDM is not supported.
func (a API) DeviceSetGPUOperationMode(device Device, mode GPUOperationMode) error {
	return a.call(a.nvmlDeviceSetGpuOperationMode, uintptr(device), uintptr(mode))
}

// DeviceSetPowerManagementLimit set new power limit of this device.
// Requires root/admin permissions.
// Note: Limit is not persistent across reboots or driver unloads.
// Enable persistent mode to prevent driver from unloading when no application is using the device.
func (a API) DeviceSetPowerManagementLimit(device Device, limit uint32) error {
	return a.call(a.nvmlDeviceSetPowerManagementLimit, uintptr(device), uintptr(limit))
}
