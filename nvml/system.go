package nvml

import (
	"C"
	"unsafe"
)

// SystemGetCudaDriverVersion retrieves the version of the CUDA driver.
// The returned CUDA driver version is the same as the CUDA API cuDriverGetVersion() would return on the system.
func (a API) SystemGetCudaDriverVersion() (cudaDriverVersion int32, err error) {
	err = a.call(a.nvmlSystemGetCudaDriverVersion, uintptr(unsafe.Pointer(&cudaDriverVersion)))
	return
}

// SystemGetDriverVersion retrieves the version of the system's graphics driver.
func (a API) SystemGetDriverVersion() (string, error) {
	buffer := [systemDriverVersionBufferSize]C.char{}
	if err := a.call(a.nvmlSystemGetDriverVersion, uintptr(unsafe.Pointer(&buffer[0])), systemDriverVersionBufferSize); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}

// SystemGetNVMLVersion retrieves the version of the NVML library.
func (a API) SystemGetNVMLVersion() (string, error) {
	buffer := [systemDriverVersionBufferSize]C.char{}
	if err := a.call(a.nvmlSystemGetNVMLVersion, uintptr(unsafe.Pointer(&buffer[0])), systemDriverVersionBufferSize); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}

// SystemGetProcessName gets name of the process with provided process id
func (a API) SystemGetProcessName(pid uint) (string, error) {
	const maxLength = 256

	buffer := [maxLength]C.char{}
	if err := a.call(a.nvmlSystemGetProcessName, uintptr(pid), uintptr(unsafe.Pointer(&buffer[0])), maxLength); err != nil {
		return "", err
	}

	return C.GoString(&buffer[0]), nil
}
