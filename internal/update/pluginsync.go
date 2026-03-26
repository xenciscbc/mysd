package update

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SyncResult reports what sync operations were performed.
type SyncResult struct {
	Added   int
	Updated int
	Deleted int
	Errors  []string // non-fatal errors (e.g., delete failed for non-existent file)
}

// SyncPlugins applies ManifestDiff by copying files from sourceDir to targetDir.
//
// sourceDir: extracted release plugin/ directory (contains commands/ and agents/ subdirs)
// targetDir: project .claude/ directory (contains commands/ and agents/ subdirs)
//
// Per D-16: existing files are overwritten without backup.
// Delete errors are non-fatal — appended to Errors slice, execution continues.
func SyncPlugins(sourceDir, targetDir string, diff ManifestDiff) SyncResult {
	var result SyncResult

	// Process commands
	if len(diff.AddCommands) > 0 || len(diff.UpdateCommands) > 0 {
		if err := os.MkdirAll(filepath.Join(targetDir, "commands"), 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("mkdir commands: %v", err))
		}
	}

	for _, file := range diff.AddCommands {
		if err := copyFile(
			filepath.Join(sourceDir, "commands", file),
			filepath.Join(targetDir, "commands", file),
		); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("add command %s: %v", file, err))
		} else {
			result.Added++
		}
	}

	for _, file := range diff.UpdateCommands {
		if err := copyFile(
			filepath.Join(sourceDir, "commands", file),
			filepath.Join(targetDir, "commands", file),
		); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("update command %s: %v", file, err))
		} else {
			result.Updated++
		}
	}

	for _, file := range diff.DeleteCommands {
		if err := os.Remove(filepath.Join(targetDir, "commands", file)); err != nil {
			// Non-fatal: append error and continue
			result.Errors = append(result.Errors, fmt.Sprintf("delete command %s: %v", file, err))
		} else {
			result.Deleted++
		}
	}

	// Process agents
	if len(diff.AddAgents) > 0 || len(diff.UpdateAgents) > 0 {
		if err := os.MkdirAll(filepath.Join(targetDir, "agents"), 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("mkdir agents: %v", err))
		}
	}

	for _, file := range diff.AddAgents {
		if err := copyFile(
			filepath.Join(sourceDir, "agents", file),
			filepath.Join(targetDir, "agents", file),
		); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("add agent %s: %v", file, err))
		} else {
			result.Added++
		}
	}

	for _, file := range diff.UpdateAgents {
		if err := copyFile(
			filepath.Join(sourceDir, "agents", file),
			filepath.Join(targetDir, "agents", file),
		); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("update agent %s: %v", file, err))
		} else {
			result.Updated++
		}
	}

	for _, file := range diff.DeleteAgents {
		if err := os.Remove(filepath.Join(targetDir, "agents", file)); err != nil {
			// Non-fatal: append error and continue
			result.Errors = append(result.Errors, fmt.Sprintf("delete agent %s: %v", file, err))
		} else {
			result.Deleted++
		}
	}

	return result
}

// copyFile copies a file from src to dst, creating parent directories as needed.
// If dst already exists it is overwritten (per D-16).
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer in.Close()

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return nil
}
