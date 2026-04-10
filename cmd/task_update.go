package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xenciscbc/mysd/internal/output"
	"github.com/xenciscbc/mysd/internal/roadmap"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/spf13/cobra"
)

var taskUpdateCmd = &cobra.Command{
	Use:   "task-update <id> <status>",
	Short: "Update task status in tasks.md",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskUpdate,
}

func init() {
	rootCmd.AddCommand(taskUpdateCmd)
}

func runTaskUpdate(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	// Parse task ID
	taskID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID %q: must be an integer", args[0])
	}

	// Validate status
	newStatus := spec.ItemStatus(args[1])
	switch newStatus {
	case spec.StatusPending, spec.StatusInProgress, spec.StatusDone, spec.StatusBlocked:
		// valid
	default:
		return fmt.Errorf("invalid status %q: must be one of pending, in_progress, done, blocked", args[1])
	}

	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	ws, err := state.LoadState(specDir)
	if err != nil {
		return err
	}

	if ws.ChangeName == "" {
		return fmt.Errorf("no active change: run mysd propose <name> first")
	}

	tasksPath := filepath.Join(specDir, "changes", ws.ChangeName, "tasks.md")

	if err := spec.UpdateTaskStatus(tasksPath, taskID, newStatus); err != nil {
		return err
	}

	// Auto-transition: if all tasks are terminal (done/skipped) and phase is planned, advance to executed
	if ws.Phase == state.PhasePlanned {
		fm, _, fmErr := spec.ParseTasksV2(tasksPath)
		if fmErr == nil && allTasksTerminal(fm) {
			if tErr := state.Transition(&ws, state.PhaseExecuted); tErr == nil {
				p.Success("All tasks complete — phase advanced to executed")
			}
		}
	}

	// Update STATE.json LastRun
	ws.LastRun = time.Now()
	if saveErr := state.SaveState(specDir, ws); saveErr != nil {
		p.Warning("Updated tasks.md but failed to update STATE.json: " + saveErr.Error())
	}
	if trackErr := roadmap.UpdateTracking(specDir, ws); trackErr != nil {
		fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
	}

	p.Success(fmt.Sprintf("Task %d status updated to %s", taskID, newStatus))
	return nil
}

// allTasksTerminal returns true if every task in the frontmatter has a terminal status (done or skipped).
// Returns false if there are no tasks.
func allTasksTerminal(fm spec.TasksFrontmatterV2) bool {
	if len(fm.Tasks) == 0 {
		return false
	}
	for _, t := range fm.Tasks {
		if t.Status != spec.StatusDone && t.Status != "skipped" {
			return false
		}
	}
	return true
}
