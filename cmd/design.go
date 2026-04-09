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

var designContextOnly bool

var designCmd = &cobra.Command{
	Use:   "design",
	Short: "Capture technical decisions and architecture",
	RunE:  runDesign,
}

func init() {
	designCmd.Flags().BoolVar(&designContextOnly, "context-only", false, "output design context as JSON for SKILL.md consumption")
	rootCmd.AddCommand(designCmd)
}

func runDesign(cmd *cobra.Command, args []string) error {
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

	if designContextOnly {
		changeDir := filepath.Join(specDir, "changes", ws.ChangeName)
		change, _ := spec.ParseChange(changeDir)

		// Build a summary of requirements for context
		var reqTexts []string
		for _, r := range change.Specs {
			reqTexts = append(reqTexts, fmt.Sprintf("[%s] %s", r.Keyword, r.Text))
		}

		ctx := map[string]interface{}{
			"spec_dir":         specDir,
			"change_name":      ws.ChangeName,
			"phase":            ws.Phase,
			"proposal_summary": change.Proposal.Body,
			"specs":            reqTexts,
			"model":            config.ResolveModel("designer", cfg.ModelProfile, cfg.ModelOverrides, cfg.CustomProfiles),
		}
		data, err := json.MarshalIndent(ctx, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal context: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	if err := state.Transition(&ws, state.PhaseDesigned); err != nil {
		return fmt.Errorf("cannot transition to designed: %w", err)
	}
	ws.LastRun = time.Now()
	if err := state.SaveState(specDir, ws); err != nil {
		return err
	}
	if trackErr := roadmap.UpdateTracking(specDir, ws); trackErr != nil {
		fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
	}

	p.Success("State transitioned to: designed")
	p.Info("Design is integrated into /mysd:plan")
	return nil
}
