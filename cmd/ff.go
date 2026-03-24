package cmd

import (
	"fmt"

	"github.com/mysd/internal/output"
	"github.com/mysd/internal/spec"
	"github.com/mysd/internal/state"
	"github.com/spf13/cobra"
)

var ffCmd = &cobra.Command{
	Use:   "ff [name]",
	Short: "Fast-forward: propose through plan (skip interactive confirmations)",
	Args:  cobra.ExactArgs(1),
	RunE:  runFF,
}

func init() {
	rootCmd.AddCommand(ffCmd)
}

func runFF(cmd *cobra.Command, args []string) error {
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

	// Load state and transition through the pipeline: None -> Proposed -> Specced -> Designed -> Planned
	ws, _ := state.LoadState(specDir)
	ws.ChangeName = name

	transitions := []state.Phase{
		state.PhaseProposed,
		state.PhaseSpecced,
		state.PhaseDesigned,
		state.PhasePlanned,
	}

	for _, nextPhase := range transitions {
		if transErr := state.Transition(&ws, nextPhase); transErr != nil {
			return fmt.Errorf("transition to %s: %w", nextPhase, transErr)
		}
		if saveErr := state.SaveState(specDir, ws); saveErr != nil {
			return fmt.Errorf("save state after %s: %w", nextPhase, saveErr)
		}
		p.Info(fmt.Sprintf("-> %s", nextPhase))
	}

	p.Success("Fast-forward complete. Ready for: mysd execute")
	return nil
}
