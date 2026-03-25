//go:build windows

package worktree

import (
	"syscall"
	"unsafe"
)

// getDiskFreeSpaceEx is the Windows API function to get disk space information.
var getDiskFreeSpaceEx = syscall.NewLazyDLL("kernel32.dll").NewProc("GetDiskFreeSpaceExW")

// getAvailableBytes returns the available disk space in bytes at the given path.
// Uses GetDiskFreeSpaceExW from kernel32.dll.
func getAvailableBytes(path string) (uint64, error) {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	var freeBytesAvailable uint64
	var totalBytes uint64
	var totalFreeBytes uint64

	ret, _, err := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)
	if ret == 0 {
		return 0, err
	}
	return freeBytesAvailable, nil
}
