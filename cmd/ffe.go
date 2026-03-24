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

var ffeCmd = &cobra.Command{
	Use:   "ffe [name]",
	Short: "Fast-forward execute: propose through execute completion",
	Args:  cobra.ExactArgs(1),
	RunE:  runFFE,
}

func init() {
	rootCmd.AddCommand(ffeCmd)
}

func runFFE(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())
	name := args[0]

	// Detect or default spec dir
	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		specDir = ".specs"
	}

	// Scaffold the change directory (propose step)
	_, err = spec.Scaffold(name, specDir)
	if err != nil {
		return fmt.Errorf("scaffold change: %w", err)
	}

	// Load state and transition through full pipeline: None -> Proposed -> Specced -> Designed -> Planned -> Executed
	ws, _ := state.LoadState(specDir)
	ws.ChangeName = name

	transitions := []state.Phase{
		state.PhaseProposed,
		state.PhaseSpecced,
		state.PhaseDesigned,
		state.PhasePlanned,
		state.PhaseExecuted,
	}

	for _, nextPhase := range transitions {
		if transErr := state.Transition(&ws, nextPhase); transErr != nil {
			return fmt.Errorf("transition to %s: %w", nextPhase, transErr)
		}
		if saveErr := state.SaveState(specDir, ws); saveErr != nil {
			return fmt.Errorf("save state after %s: %w", nextPhase, saveErr)
		}
		if trackErr := roadmap.UpdateTracking(specDir, ws); trackErr != nil {
			fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
		}
		p.Info(fmt.Sprintf("-> %s", nextPhase))
	}

	p.Success("Fast-forward execute complete. Ready for: mysd verify")
	return nil
}
