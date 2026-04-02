package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	archiveCmd.Flags().Bool("analyze-skipped", false, "Output skipped tasks and their spec requirement relationships as JSON")
	rootCmd.AddCommand(archiveCmd)
}

func runArchiveCmd(cmd *cobra.Command, args []string) error {
	skipPrompt, _ := cmd.Flags().GetBool("yes")
	analyzeSkipped, _ := cmd.Flags().GetBool("analyze-skipped")

	specsDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	ws, err := state.LoadState(specsDir)
	if err != nil {
		return err
	}

	// --analyze-skipped: output skipped tasks as JSON and exit
	if analyzeSkipped {
		return runAnalyzeSkipped(cmd, specsDir, ws)
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

	// Gate 2: all tasks must be completed or skipped
	changeDir := filepath.Join(specsDir, "changes", ws.ChangeName)
	if err := checkTasksDone(changeDir); err != nil {
		return err
	}

	// Gate 3: all MUST items must be DONE
	if err := checkMustItemsDone(changeDir, ws.ChangeName); err != nil {
		return err
	}

	// Archive: move directory
	archiveDir := filepath.Join(specsDir, "changes", "archive", time.Now().Format("2006-01-02")+"-"+ws.ChangeName)
	if err := os.MkdirAll(filepath.Dir(archiveDir), 0755); err != nil {
		return fmt.Errorf("create archive parent dir: %w", err)
	}

	// Save STATE.json snapshot to change dir as ARCHIVED-STATE.json before moving (Pitfall 5)
	if err := saveArchivedState(changeDir, ws); err != nil {
		// Log warning but don't fail — snapshot is best-effort
		fmt.Fprintf(os.Stderr, "warning: could not save ARCHIVED-STATE.json: %v\n", err)
	}

	// Delta spec merge: merge delta specs back into main specs before moving
	if mergeWarnings := mergeDeltaSpecs(specsDir, changeDir); len(mergeWarnings) > 0 {
		for _, w := range mergeWarnings {
			fmt.Fprintf(os.Stderr, "warning: %s\n", w)
		}
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

// SkippedTaskInfo represents a skipped task and its relationship to spec requirements.
type SkippedTaskInfo struct {
	TaskID     int    `json:"task_id"`
	TaskName   string `json:"task_name"`
	SkipReason string `json:"skip_reason"`
}

// runAnalyzeSkipped outputs skipped tasks as JSON without performing the archive.
func runAnalyzeSkipped(cmd *cobra.Command, specsDir string, ws state.WorkflowState) error {
	changeDir := filepath.Join(specsDir, "changes", ws.ChangeName)
	tasksPath := filepath.Join(changeDir, "tasks.md")

	tasks, _, err := spec.ParseTasks(tasksPath)
	if err != nil {
		return fmt.Errorf("parse tasks: %w", err)
	}

	var skipped []SkippedTaskInfo
	for _, t := range tasks {
		if t.Skipped {
			skipped = append(skipped, SkippedTaskInfo{
				TaskID:     t.ID,
				TaskName:   t.Name,
				SkipReason: t.SkipReason,
			})
		}
	}

	if skipped == nil {
		skipped = []SkippedTaskInfo{}
	}

	data, err := json.MarshalIndent(skipped, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal skipped tasks: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

// checkTasksDone verifies that all tasks in tasks.md are completed ([x]) or skipped ([~]).
// Returns an error if any task is still pending ([ ]).
func checkTasksDone(changeDir string) error {
	tasksPath := filepath.Join(changeDir, "tasks.md")
	tasks, _, err := spec.ParseTasks(tasksPath)
	if err != nil {
		// If tasks.md doesn't exist, no task gate to enforce
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("parse tasks for gate check: %w", err)
	}

	var incomplete int
	for _, t := range tasks {
		if t.Status == spec.StatusPending {
			incomplete++
		}
	}
	if incomplete > 0 {
		return fmt.Errorf("cannot archive: %d incomplete task(s) remain", incomplete)
	}
	return nil
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

// mergeDeltaSpecs iterates over delta spec directories in changes/<name>/specs/
// and merges each delta spec into the corresponding main spec at openspec/specs/<capability>/spec.md.
// Returns accumulated warnings from all merge operations.
func mergeDeltaSpecs(specsDir, changeDir string) []string {
	var allWarnings []string

	deltaSpecsDir := filepath.Join(changeDir, "specs")
	entries, err := os.ReadDir(deltaSpecsDir)
	if err != nil {
		// No specs directory — nothing to merge
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		capName := entry.Name()
		deltaSpecPath := filepath.Join(deltaSpecsDir, capName, "spec.md")

		deltaContent, err := os.ReadFile(deltaSpecPath)
		if err != nil {
			continue // no spec.md in this capability dir
		}

		mainSpecPath := filepath.Join(specsDir, "specs", capName, "spec.md")

		merged, warnings, err := spec.MergeSpecs(mainSpecPath, string(deltaContent))
		if err != nil {
			allWarnings = append(allWarnings, fmt.Sprintf("merge %s: %v", capName, err))
			continue
		}
		allWarnings = append(allWarnings, warnings...)

		// Ensure the main spec directory exists
		if err := os.MkdirAll(filepath.Dir(mainSpecPath), 0755); err != nil {
			allWarnings = append(allWarnings, fmt.Sprintf("create dir for %s: %v", capName, err))
			continue
		}

		if err := os.WriteFile(mainSpecPath, []byte(merged), 0644); err != nil {
			allWarnings = append(allWarnings, fmt.Sprintf("write merged spec %s: %v", capName, err))
		}
	}

	return allWarnings
}
