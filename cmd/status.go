package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/xenciscbc/mysd/internal/executor"
	"github.com/xenciscbc/mysd/internal/output"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current spec state and progress",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	ws, err := state.LoadState(specDir)
	if err != nil {
		return err
	}

	if ws.ChangeName == "" {
		p.Info("No active change. Run: mysd propose <name>")
		return nil
	}

	changeDir := filepath.Join(specDir, "changes", ws.ChangeName)

	// Load tasks via ParseTasks (returns []spec.Task, compatible with BuildStatusSummary)
	tasks, _, err := spec.ParseTasks(filepath.Join(changeDir, "tasks.md"))
	if err != nil {
		// Non-fatal — render with empty tasks
		tasks = nil
	}

	// Parse change for requirements
	change, err := spec.ParseChange(changeDir)
	if err != nil {
		// Non-fatal — render with empty reqs
		change = spec.Change{}
	}

	summary := executor.BuildStatusSummary(ws, tasks, change.Specs)
	executor.RenderStatus(cmd.OutOrStdout(), summary)

	// Deferred notes count (D-09)
	count, _ := spec.CountDeferredNotes(specDir)
	if count > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "\nDeferred notes: %d — run /mysd:note to browse\n", count)
	}
	return nil
}
