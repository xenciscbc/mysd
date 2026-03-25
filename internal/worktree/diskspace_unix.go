//go:build !windows

package worktree

import "syscall"

// getAvailableBytes returns the available disk space in bytes at the given path.
// Uses syscall.Statfs which is available on Linux, macOS, and other Unix systems.
func getAvailableBytes(path string) (uint64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}
	// Bavail is blocks available to unprivileged users; Bsize is block size
	return stat.Bavail * uint64(stat.Bsize), nil
}
