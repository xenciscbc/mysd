package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/xenciscbc/mysd/internal/roadmap"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/xenciscbc/mysd/internal/verifier"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive completed spec to history",
	RunE:  runArchiveCmd,
}

func init() {
	archiveCmd.Flags().Bool("yes", false, "Skip UAT prompt (for non-interactive/CI usage)")
	rootCmd.AddCommand(archiveCmd)
}

func runArchiveCmd(cmd *cobra.Command, args []string) error {
	skipPrompt, _ := cmd.Flags().GetBool("yes")

	specsDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	ws, err := state.LoadState(specsDir)
	if err != nil {
		return err
	}

	// Handle interactive UAT prompt before the gate checks
	if !skipPrompt && isInteractive() {
		fmt.Fprint(cmd.OutOrStdout(), "Run UAT first? [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		resp, _ := reader.ReadString('\n')
		resp = strings.TrimSpace(strings.ToLower(resp))
		if resp == "y" || resp == "yes" {
			fmt.Fprintln(cmd.OutOrStdout(), "Hint: Run /mysd:uat to start UAT. Continuing with archive...")
		}
	}

	return runArchive(specsDir, ws, skipPrompt)
}

// runArchive is the testable core of the archive command.
// skipPrompt=true suppresses the interactive UAT prompt.
func runArchive(specsDir string, ws state.WorkflowState, skipPrompt bool) error {
	// Gate 1: state must be verified
	if ws.Phase != state.PhaseVerified {
		return fmt.Errorf("cannot archive: phase is %s, must be verified", ws.Phase)
	}

	// Gate 2: all MUST items must be DONE
	changeDir := filepath.Join(specsDir, "changes", ws.ChangeName)
	if err := checkMustItemsDone(changeDir, ws.ChangeName); err != nil {
		return err
	}

	// Archive: move directory
	archiveDir := filepath.Join(specsDir, "archive", ws.ChangeName)
	if err := os.MkdirAll(filepath.Dir(archiveDir), 0755); err != nil {
		return fmt.Errorf("create archive parent dir: %w", err)
	}

	// Save STATE.json snapshot to change dir as ARCHIVED-STATE.json before moving (Pitfall 5)
	if err := saveArchivedState(changeDir, ws); err != nil {
		// Log warning but don't fail — snapshot is best-effort
		fmt.Fprintf(os.Stderr, "warning: could not save ARCHIVED-STATE.json: %v\n", err)
	}

	// Delete discuss research cache before moving (D-18)
	deleteResearchCache(changeDir)

	// Move change directory to archive
	if err := moveDir(changeDir, archiveDir); err != nil {
		return fmt.Errorf("move change directory: %w", err)
	}

	// Transition state to archived
	if err := state.Transition(&ws, state.PhaseArchived); err != nil {
		return fmt.Errorf("state transition: %w", err)
	}
	if err := state.SaveState(specsDir, ws); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	if trackErr := roadmap.UpdateTracking(specsDir, ws); trackErr != nil {
		fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
	}

	fmt.Printf("Archived %s to %s\n", ws.ChangeName, archiveDir)
	return nil
}

// deleteResearchCache removes the discuss research cache file from changeDir (best-effort, D-18).
func deleteResearchCache(changeDir string) {
	_ = os.Remove(filepath.Join(changeDir, "discuss-research-cache.json"))
}

// checkMustItemsDone verifies that all MUST items in the change are done in verification-status.json.
// Returns an error naming the first undone MUST item found.
func checkMustItemsDone(changeDir, changeName string) error {
	// Parse the change to get MUST requirements
	change, err := spec.ParseChange(changeDir)
	if err != nil {
		return fmt.Errorf("parse change for gate check: %w", err)
	}

	// Read verification status
	vs, err := spec.ReadVerificationStatus(changeDir)
	if err != nil {
		return fmt.Errorf("read verification status: %w", err)
	}

	// Check each MUST requirement
	for _, r := range change.Specs {
		if r.Keyword != spec.Must {
			continue
		}
		id := verifier.StableID(r)
		status, ok := vs.Requirements[id]
		if !ok || status != spec.StatusDone {
			return fmt.Errorf("cannot archive: MUST item %q is not done (status: %s)", truncate(r.Text, 60), status)
		}
	}

	return nil
}

// saveArchivedState writes a ARCHIVED-STATE.json snapshot to changeDir.
func saveArchivedState(changeDir string, ws state.WorkflowState) error {
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal archived state: %w", err)
	}
	archivedStatePath := filepath.Join(changeDir, "ARCHIVED-STATE.json")
	return os.WriteFile(archivedStatePath, data, 0644)
}

// moveDir moves a directory from src to dst.
// It tries os.Rename first (atomic on same volume).
// On failure (cross-volume on Windows), falls back to recursive copy + os.RemoveAll.
func moveDir(src, dst string) error {
	// Try atomic rename first
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// Fallback: recursive copy then delete source
	if err := copyDir(src, dst); err != nil {
		return fmt.Errorf("copy dir (rename fallback): %w", err)
	}
	if err := os.RemoveAll(src); err != nil {
		return fmt.Errorf("remove source after copy (rename fallback): %w", err)
	}
	return nil
}

// copyDir recursively copies src directory to dst.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// isInteractive returns true if stdin is a terminal.
func isInteractive() bool {
	return term.IsTerminal(os.Stdin.Fd())
}

// truncate returns s truncated to max characters, with "..." appended if truncated.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
