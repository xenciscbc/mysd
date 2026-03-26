//go:build windows

package update

import (
	"fmt"
	"os"
)

// replaceExecutable implements Windows rename-then-replace strategy.
// Windows cannot delete or overwrite a running executable, so we:
//  1. Rename current exe to exe + ".old"
//  2. Rename/move new binary to exe path
//  3. Attempt to delete .old (may fail if still locked — that's OK)
func replaceExecutable(currentExe, newBinaryPath string) error {
	oldPath := currentExe + ".old"

	// Step 1: Rename current exe to .old
	if err := os.Rename(currentExe, oldPath); err != nil {
		return fmt.Errorf("update: failed to rename current executable to .old: %w", err)
	}

	// Step 2: Move new binary to exe path
	if err := os.Rename(newBinaryPath, currentExe); err != nil {
		// Try to restore original on failure
		_ = os.Rename(oldPath, currentExe)
		return fmt.Errorf("update: failed to move new binary into place: %w", err)
	}

	// Step 3: Attempt to delete .old (may fail if binary is still locked — that's OK on Windows)
	_ = os.Remove(oldPath)

	return nil
}
