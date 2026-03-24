package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/output"
	"github.com/xenciscbc/mysd/internal/roadmap"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/spf13/cobra"
)

var specContextOnly bool

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Define detailed requirements with RFC 2119 keywords",
	RunE:  runSpec,
}

func init() {
	specCmd.Flags().BoolVar(&specContextOnly, "context-only", false, "output spec context as JSON for SKILL.md consumption")
	rootCmd.AddCommand(specCmd)
}

func runSpec(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return fmt.Errorf("no spec directory: %w", err)
	}

	ws, err := state.LoadState(specDir)
	if err != nil {
		return err
	}

	cfg, err := config.Load(".")
	if err != nil {
		return err
	}

	if specContextOnly {
		changeDir := filepath.Join(specDir, "changes", ws.ChangeName)
		change, _ := spec.ParseChange(changeDir)

		ctx := map[string]interface{}{
			"change_name": ws.ChangeName,
			"phase":       ws.Phase,
			"proposal":    change.Proposal.Body,
			"model":       config.ResolveModel("spec-writer", cfg.ModelProfile, cfg.ModelOverrides),
		}
		data, err := json.MarshalIndent(ctx, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal context: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	if err := state.Transition(&ws, state.PhaseSpecced); err != nil {
		return fmt.Errorf("cannot transition to specced: %w", err)
	}
	ws.LastRun = time.Now()
	if err := state.SaveState(specDir, ws); err != nil {
		return err
	}
	if trackErr := roadmap.UpdateTracking(specDir, ws); trackErr != nil {
		fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
	}

	p.Success("State transitioned to: specced")
	p.Info("Use /mysd:spec in Claude Code for AI-assisted spec writing")
	return nil
}
