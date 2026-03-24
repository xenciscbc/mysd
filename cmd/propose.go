package cmd

import (
	"fmt"
	"os"

	"github.com/mysd/internal/output"
	"github.com/mysd/internal/roadmap"
	"github.com/mysd/internal/spec"
	"github.com/mysd/internal/state"
	"github.com/spf13/cobra"
)

var proposeCmd = &cobra.Command{
	Use:   "propose [name]",
	Short: "Create a new spec from description",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runPropose,
}

func init() {
	rootCmd.AddCommand(proposeCmd)
}

func runPropose(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	// Detect spec dir; if not found, default to .specs (new project)
	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		specDir = ".specs"
	}

	// Scaffold the change directory
	change, err := spec.Scaffold(args[0], specDir)
	if err != nil {
		p.Error(fmt.Sprintf("Failed to create spec: %s", err))
		return err
	}

	// Load state, transition to PhaseProposed, save
	ws, _ := state.LoadState(specDir)
	if transErr := state.Transition(&ws, state.PhaseProposed); transErr != nil {
		// If already proposed (e.g. re-running), that's fine — just set the change name
		ws.Phase = state.PhaseProposed
	}
	ws.ChangeName = args[0]
	if saveErr := state.SaveState(specDir, ws); saveErr != nil {
		p.Error(fmt.Sprintf("Failed to save state: %s", saveErr))
		return saveErr
	}
	if trackErr := roadmap.UpdateTracking(specDir, ws); trackErr != nil {
		fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
	}

	p.Success(fmt.Sprintf("Created spec: %s", change.Dir))
	p.Info("Next: mysd spec — define detailed requirements")
	return nil
}
