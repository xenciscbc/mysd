//go:build !windows

package update

import (
	"fmt"
	"io"
	"os"
)

// replaceExecutable implements Unix direct-replace strategy with 0755 permissions.
//  1. Rename current exe to exe + ".old"
//  2. Copy new binary to exe path with 0755 permissions
//  3. Remove .old file
func replaceExecutable(currentExe, newBinaryPath string) error {
	oldPath := currentExe + ".old"

	// Step 1: Rename current exe to .old
	if err := os.Rename(currentExe, oldPath); err != nil {
		return fmt.Errorf("update: failed to rename current executable to .old: %w", err)
	}

	// Step 2: Copy new binary to exe path with executable permissions
	if err := copyFile(newBinaryPath, currentExe, 0755); err != nil {
		// Try to restore original on failure
		_ = os.Rename(oldPath, currentExe)
		return fmt.Errorf("update: failed to copy new binary into place: %w", err)
	}

	// Step 3: Remove .old file
	_ = os.Remove(oldPath)

	return nil
}

// copyFile copies src to dst with the given permissions.
func copyFile(src, dst string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("update: cannot open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("update: cannot create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("update: copy failed: %w", err)
	}

	return nil
}
