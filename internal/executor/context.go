package executor

import (
	"fmt"
	"path/filepath"

	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/spec"
)

// ExecutionContext is the JSON-serializable context passed to SKILL.md consumers (per EXEC-01).
// It describes all pending work, requirements, and configuration for the current execution run.
type ExecutionContext struct {
	SpecDir       string            `json:"spec_dir"`
	ChangeName    string            `json:"change_name"`
	MustItems     []RequirementItem `json:"must_items"`
	ShouldItems   []RequirementItem `json:"should_items"`
	MayItems      []RequirementItem `json:"may_items"`
	Tasks         []TaskItem        `json:"tasks"`
	PendingTasks  []TaskItem        `json:"pending_tasks"`
	TDDMode       bool              `json:"tdd_mode"`
	AtomicCommits bool              `json:"atomic_commits"`
	ExecutionMode string            `json:"execution_mode"`
	AgentCount    int               `json:"agent_count"`
	// Wave grouping fields (Phase 06 extension — additive only per D-11)
	WaveGroups     [][]TaskItem `json:"wave_groups,omitempty"`
	WorktreeDir    string       `json:"worktree_dir,omitempty"`
	AutoMode       bool         `json:"auto_mode,omitempty"`
	HasParallelOpp bool         `json:"has_parallel_opportunity"`
	// Doc update fields (Phase 11 extension — additive only per D-11/D-12)
	DocsToUpdate []string `json:"docs_to_update,omitempty"`
	// Dynamic instruction for orchestrator (dynamic-apply-instructions change)
	Instruction string `json:"instruction"`
	// Profile-resolved model names for agent spawning
	Model         string `json:"model"`
	VerifierModel string `json:"verifier_model"`
}

// RequirementItem is a flattened requirement for JSON output.
type RequirementItem struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// TaskItem is a flattened task for JSON output.
type TaskItem struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status"`
	Spec        string   `json:"spec,omitempty"`
	Depends     []int    `json:"depends,omitempty"`
	Files       []string `json:"files,omitempty"`
	Satisfies   []string `json:"satisfies,omitempty"`
	Skills      []string `json:"skills,omitempty"`
}

// BuildContextFromParts constructs an ExecutionContext from pre-loaded parts.
// This is the primary constructor for testing and internal composition.
func BuildContextFromParts(
	changeName string,
	tasks []spec.TaskEntry,
	reqs []spec.Requirement,
	cfg config.ProjectConfig,
) ExecutionContext {
	ctx := ExecutionContext{
		ChangeName:    changeName,
		TDDMode:       cfg.TDD,
		AtomicCommits: cfg.AtomicCommits,
		ExecutionMode: cfg.ExecutionMode,
		AgentCount:    cfg.AgentCount,
	}

	// Populate all tasks
	for _, t := range tasks {
		ctx.Tasks = append(ctx.Tasks, TaskItem{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Status:      string(t.Status),
			Spec:        t.Spec,
			Depends:     t.Depends,
			Files:       t.Files,
			Satisfies:   t.Satisfies,
			Skills:      t.Skills,
		})
	}

	// Populate pending tasks (excludes done and blocked)
	for _, t := range PendingTasks(tasks) {
		ctx.PendingTasks = append(ctx.PendingTasks, TaskItem{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Status:      string(t.Status),
			Spec:        t.Spec,
			Depends:     t.Depends,
			Files:       t.Files,
			Satisfies:   t.Satisfies,
			Skills:      t.Skills,
		})
	}

	// Compute wave groups and parallel opportunity from pending tasks
	wg, _ := BuildWaveGroups(ctx.PendingTasks)
	ctx.WaveGroups = wg
	ctx.WorktreeDir = cfg.WorktreeDir
	ctx.AutoMode = cfg.AutoMode
	ctx.HasParallelOpp = HasParallelOpportunity(ctx.PendingTasks)
	ctx.DocsToUpdate = cfg.DocsToUpdate

	// Classify requirements by RFC 2119 keyword
	for _, r := range reqs {
		item := RequirementItem{ID: r.ID, Text: r.Text}
		switch r.Keyword {
		case spec.Must:
			ctx.MustItems = append(ctx.MustItems, item)
		case spec.Should:
			ctx.ShouldItems = append(ctx.ShouldItems, item)
		case spec.May:
			ctx.MayItems = append(ctx.MayItems, item)
		}
	}

	return ctx
}

// BuildContext loads change data from disk and constructs an ExecutionContext.
// specsDir is the root specs directory (e.g., ".specs") and changeName identifies the change.
// This is the primary entrypoint for the `mysd execute --context-only` command (EXEC-01).
func BuildContext(specsDir string, changeName string, cfg config.ProjectConfig) (ExecutionContext, error) {
	changeDir := filepath.Join(specsDir, "changes", changeName)

	// Load tasks via updater (supports per-task status tracking)
	tasksPath := filepath.Join(changeDir, "tasks.md")
	fm, _, err := spec.ParseTasksV2(tasksPath)
	if err != nil {
		return ExecutionContext{}, fmt.Errorf("load tasks: %w", err)
	}

	// Parse change for requirements
	change, err := spec.ParseChange(changeDir)
	if err != nil {
		return ExecutionContext{}, fmt.Errorf("parse change: %w", err)
	}

	return BuildContextFromParts(changeName, fm.Tasks, change.Specs, cfg), nil
}
