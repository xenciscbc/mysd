package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/executor"
	"github.com/xenciscbc/mysd/internal/output"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/spf13/cobra"
)

var (
	contextOnly bool
	executeSpec string
	preflight   bool
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Run tasks with pre-execution alignment",
	RunE:  runExecute,
}

func init() {
	executeCmd.Flags().BoolVar(&contextOnly, "context-only", false, "output execution context as JSON (for SKILL.md consumption)")
	executeCmd.Flags().StringVar(&executeSpec, "spec", "", "filter execution to tasks matching the specified spec")
	executeCmd.Flags().BoolVar(&preflight, "preflight", false, "run pre-execution validation (file existence, artifact staleness)")
	rootCmd.AddCommand(executeCmd)
}

func runExecute(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	ws, err := state.LoadState(specDir)
	if err != nil {
		return err
	}

	cfg, err := config.Load(".")
	if err != nil {
		return err
	}

	// Override config with flags
	if cmd.Flags().Changed("tdd") {
		cfg.TDD, _ = cmd.Flags().GetBool("tdd")
	}
	if cmd.Flags().Changed("atomic-commits") {
		cfg.AtomicCommits, _ = cmd.Flags().GetBool("atomic-commits")
	}
	if cmd.Flags().Changed("execution-mode") {
		cfg.ExecutionMode, _ = cmd.Flags().GetString("execution-mode")
	}
	if cmd.Flags().Changed("agent-count") {
		cfg.AgentCount, _ = cmd.Flags().GetInt("agent-count")
	}

	if preflight {
		report, err := runPreflight(specDir, ws)
		if err != nil {
			return err
		}
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	if contextOnly {
		ctx, err := executor.BuildContext(specDir, ws.ChangeName, cfg)
		if err != nil {
			return err
		}

		// Filter by --spec flag
		if executeSpec != "" {
			ctx.PendingTasks = filterTasksBySpec(ctx.PendingTasks, executeSpec)
			// Recompute wave groups from filtered tasks
			wg, _ := executor.BuildWaveGroups(ctx.PendingTasks)
			ctx.WaveGroups = wg
			ctx.HasParallelOpp = executor.HasParallelOpportunity(ctx.PendingTasks)
		}

		// Generate dynamic instruction from current state + preflight
		report, _ := runPreflight(specDir, ws)
		ctx.Instruction = executor.GenerateInstruction(ctx, &report)

		data, _ := json.MarshalIndent(ctx, "", "  ")
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	// Non-context-only: print guidance directing to SKILL.md
	p.Info("Use /mysd:apply in Claude Code to run with AI alignment gate")
	p.Info("Or use: mysd execute --context-only | jq")
	return nil
}

// runPreflight performs pre-execution validation: file existence and artifact staleness.
func runPreflight(specDir string, ws state.WorkflowState) (executor.PreflightReport, error) {
	changeDir := filepath.Join(specDir, "changes", ws.ChangeName)
	tasksPath := filepath.Join(changeDir, "tasks.md")

	fm, _, err := spec.ParseTasksV2(tasksPath)
	if err != nil {
		return executor.PreflightReport{}, fmt.Errorf("parse tasks: %w", err)
	}

	// Check file existence for pending tasks
	var missingFiles []string
	for _, task := range fm.Tasks {
		if task.Status == spec.StatusDone {
			continue
		}
		nameLower := strings.ToLower(task.Name)
		descLower := strings.ToLower(task.Description)
		isCreateTask := strings.Contains(nameLower, "create") || strings.Contains(nameLower, "add") ||
			strings.Contains(descLower, "create") || strings.Contains(descLower, "add")

		for _, f := range task.Files {
			if isCreateTask {
				continue
			}
			if _, err := os.Stat(f); os.IsNotExist(err) {
				missingFiles = append(missingFiles, f)
			}
		}
	}

	// Check staleness from STATE.json last_run
	staleness := executor.StalenessCheck{}
	if ws.LastRun.IsZero() {
		staleness.DaysSinceLastPlan = -1
		staleness.IsStale = true
	} else {
		days := int(math.Floor(time.Since(ws.LastRun).Hours() / 24))
		staleness.DaysSinceLastPlan = days
		staleness.IsStale = days > 7
	}

	// Determine overall status
	status := "ok"
	if len(missingFiles) > 0 || staleness.IsStale {
		status = "warning"
	}
	if staleness.DaysSinceLastPlan > 30 || staleness.DaysSinceLastPlan == -1 {
		status = "critical"
	}

	if missingFiles == nil {
		missingFiles = []string{}
	}

	return executor.PreflightReport{
		Status: status,
		Checks: executor.PreflightChecks{
			MissingFiles: missingFiles,
			Staleness:    staleness,
		},
	}, nil
}

// filterTasksBySpec returns only tasks whose Spec field matches the given spec name.
// Tasks with an empty Spec field (change-level tasks) are excluded from per-spec filtering.
func filterTasksBySpec(tasks []executor.TaskItem, specName string) []executor.TaskItem {
	var filtered []executor.TaskItem
	for _, t := range tasks {
		if t.Spec == specName {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
