package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/output"
	"github.com/xenciscbc/mysd/internal/planchecker"
	"github.com/xenciscbc/mysd/internal/roadmap"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
)

var (
	planContextOnly     bool
	planResearch        bool
	planCheck           bool
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Break design into executable task list",
	RunE:  runPlan,
}

func init() {
	planCmd.Flags().BoolVar(&planContextOnly, "context-only", false, "output plan context as JSON for SKILL.md consumption")
	planCmd.Flags().BoolVar(&planResearch, "research", false, "enable research phase before planning (deeper analysis)")
	planCmd.Flags().BoolVar(&planCheck, "check", false, "enable plan check phase after planning (validation pass)")
	rootCmd.AddCommand(planCmd)
}

func runPlan(cmd *cobra.Command, args []string) error {
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

	if planContextOnly {
		changeDir := filepath.Join(specDir, "changes", ws.ChangeName)
		change, _ := spec.ParseChange(changeDir)

		// Build a summary of requirements for context
		var reqTexts []string
		for _, r := range change.Specs {
			reqTexts = append(reqTexts, fmt.Sprintf("[%s] %s", r.Keyword, r.Text))
		}

		ctx := map[string]interface{}{
			"change_name":      ws.ChangeName,
			"phase":            ws.Phase,
			"specs":            reqTexts,
			"design":           change.Design.Body,
			"model":            config.ResolveModel("planner", cfg.ModelProfile, cfg.ModelOverrides),
			"research_enabled": planResearch,
			"check_enabled":    planCheck,
			"test_generation":  cfg.TestGeneration,
			"wave_groups":      [][]int{},       // empty for Phase 5, populated in Phase 6
			"worktree_dir":     cfg.WorktreeDir, // from ProjectConfig (default ".worktrees")
			"auto_mode":        cfg.AutoMode,    // from ProjectConfig (default false)
		}

		if planCheck {
			tasksPath := filepath.Join(changeDir, "tasks.md")
			fm, _, parseErr := spec.ParseTasksV2(tasksPath)
			if parseErr == nil {
				var mustIDs []string
				for _, r := range change.Specs {
					if r.Keyword == spec.Must && r.ID != "" {
						mustIDs = append(mustIDs, r.ID)
					}
				}
				ctx["coverage"] = planchecker.CheckCoverage(fm.Tasks, mustIDs)
			}
		}

		data, err := json.MarshalIndent(ctx, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal context: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	if err := state.Transition(&ws, state.PhasePlanned); err != nil {
		return fmt.Errorf("cannot transition to planned: %w", err)
	}
	ws.LastRun = time.Now()
	if err := state.SaveState(specDir, ws); err != nil {
		return err
	}
	if trackErr := roadmap.UpdateTracking(specDir, ws); trackErr != nil {
		fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
	}

	p.Success("State transitioned to: planned")
	p.Info("Use /mysd:plan in Claude Code for AI-assisted planning")
	return nil
}
