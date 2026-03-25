package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/xenciscbc/mysd/internal/executor"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestChange creates a minimal .specs/ change structure for execute integration tests.
// tasks is a TasksFrontmatterV2 with pre-defined entries and statuses.
// phase is the WorkflowState phase to save.
//
// The tasks.md is written in V2 YAML frontmatter format (for executor.BuildContext),
// with a markdown body listing tasks as checkboxes (for executor.BuildStatusSummary via ParseTasks).
func setupTestChange(t *testing.T, dir string, tasks spec.TasksFrontmatterV2, phase state.Phase) {
	t.Helper()

	specsDir := filepath.Join(dir, ".specs")
	changeDir := filepath.Join(specsDir, "changes", "test-change")

	require.NoError(t, os.MkdirAll(changeDir, 0755))

	// Create specs/auth-spec/spec.md with a MUST requirement
	specCapDir := filepath.Join(changeDir, "specs", "auth-spec")
	require.NoError(t, os.MkdirAll(specCapDir, 0755))
	specContent := "## Requirement: Auth\n\nThe system MUST authenticate users.\n"
	require.NoError(t, os.WriteFile(filepath.Join(specCapDir, "spec.md"), []byte(specContent), 0644))

	// Write .openspec.yaml
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, ".openspec.yaml"), []byte("schema: spec-driven\ncreated: 2026-01-01\n"), 0644))

	// Write proposal.md
	proposalContent := "---\nspec-version: \"1\"\nchange: test-change\nstatus: proposed\ncreated: 2026-01-01\nupdated: 2026-01-01\n---\n\n## Summary\n\nTest change.\n"
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644))

	// Write design.md
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "design.md"), []byte("## Architecture\n\nTest design.\n"), 0644))

	// Build markdown task body: ParseTasks reads `- [x] name` / `- [ ] name` from body.
	// Also write V2 YAML frontmatter so executor.BuildContext (ParseTasksV2) reads status correctly.
	var bodyLines []string
	bodyLines = append(bodyLines, "\n## Tasks\n")
	for _, te := range tasks.Tasks {
		if te.Status == spec.StatusDone {
			bodyLines = append(bodyLines, "- [x] "+te.Name)
		} else {
			bodyLines = append(bodyLines, "- [ ] "+te.Name)
		}
	}
	body := ""
	for _, l := range bodyLines {
		body += l + "\n"
	}

	require.NoError(t, spec.WriteTasks(filepath.Join(changeDir, "tasks.md"), tasks, body))

	// Write STATE.json
	ws := state.WorkflowState{
		ChangeName: "test-change",
		Phase:      phase,
	}
	require.NoError(t, state.SaveState(specsDir, ws))
}

func TestExecuteContextOnly(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// 3 tasks: 1 done, 2 pending
	tasks := spec.TasksFrontmatterV2{
		SpecVersion: "1",
		Total:       3,
		Completed:   1,
		Tasks: []spec.TaskEntry{
			{ID: 1, Name: "Task One", Status: spec.StatusDone},
			{ID: 2, Name: "Task Two", Status: spec.StatusPending},
			{ID: 3, Name: "Task Three", Status: spec.StatusPending},
		},
	}
	setupTestChange(t, tmpDir, tasks, state.PhasePlanned)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"execute", "--context-only"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output, "output should not be empty")

	// Validate JSON output
	var ctx executor.ExecutionContext
	require.NoError(t, json.Unmarshal([]byte(output), &ctx), "output must be valid JSON ExecutionContext")

	assert.Equal(t, "test-change", ctx.ChangeName)
	assert.Len(t, ctx.Tasks, 3, "should have 3 total tasks")
	assert.Len(t, ctx.PendingTasks, 2, "should have 2 pending tasks (excluding done)")
	assert.Len(t, ctx.MustItems, 1, "should have 1 MUST requirement")
}

func TestExecuteResumeFromInterruption(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// 3 tasks: 2 done, 1 pending — simulates resume after interruption
	tasks := spec.TasksFrontmatterV2{
		SpecVersion: "1",
		Total:       3,
		Completed:   2,
		Tasks: []spec.TaskEntry{
			{ID: 1, Name: "Task One", Status: spec.StatusDone},
			{ID: 2, Name: "Task Two", Status: spec.StatusDone},
			{ID: 3, Name: "Task Three", Status: spec.StatusPending},
		},
	}
	setupTestChange(t, tmpDir, tasks, state.PhasePlanned)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"execute", "--context-only"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	var ctx executor.ExecutionContext
	require.NoError(t, json.Unmarshal([]byte(output), &ctx), "output must be valid JSON")

	assert.Len(t, ctx.Tasks, 3, "should have 3 total tasks")
	assert.Len(t, ctx.PendingTasks, 1, "should have 1 pending task (resume from interruption)")
	assert.Equal(t, 3, ctx.PendingTasks[0].ID, "pending task should be task 3")
}

func TestExecuteTDDFlag(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Config file with tdd: false
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte("tdd: false\n"), 0644))

	tasks := spec.TasksFrontmatterV2{
		SpecVersion: "1",
		Total:       1,
		Completed:   0,
		Tasks: []spec.TaskEntry{
			{ID: 1, Name: "Task One", Status: spec.StatusPending},
		},
	}
	setupTestChange(t, tmpDir, tasks, state.PhasePlanned)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	// --tdd flag overrides config tdd: false (per CONF-04)
	rootCmd.SetArgs([]string{"execute", "--context-only", "--tdd"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	var ctx executor.ExecutionContext
	require.NoError(t, json.Unmarshal([]byte(output), &ctx), "output must be valid JSON")

	assert.True(t, ctx.TDDMode, "--tdd flag should override config tdd: false")
}

func TestExecuteContextOnly_WaveGroups(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Task 2 depends on Task 1 — this creates a wave structure:
	// Wave 0: [Task 1], Wave 1: [Task 2]
	tasks := spec.TasksFrontmatterV2{
		SpecVersion: "1",
		Total:       2,
		Completed:   0,
		Tasks: []spec.TaskEntry{
			{ID: 1, Name: "Task One", Status: spec.StatusPending, Files: []string{"internal/auth.go"}},
			{ID: 2, Name: "Task Two", Status: spec.StatusPending, Depends: []int{1}, Files: []string{"cmd/auth.go"}},
		},
	}
	setupTestChange(t, tmpDir, tasks, state.PhasePlanned)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"execute", "--context-only"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	var ctx executor.ExecutionContext
	require.NoError(t, json.Unmarshal([]byte(output), &ctx), "output must be valid JSON ExecutionContext")

	// wave_groups must be present and non-nil
	assert.NotNil(t, ctx.WaveGroups, "wave_groups should not be nil")
	assert.Len(t, ctx.WaveGroups, 2, "should have 2 waves: [Task1], [Task2]")
	if len(ctx.WaveGroups) == 2 {
		assert.Len(t, ctx.WaveGroups[0], 1, "wave 0 should have 1 task")
		assert.Len(t, ctx.WaveGroups[1], 1, "wave 1 should have 1 task")
		assert.Equal(t, 1, ctx.WaveGroups[0][0].ID, "wave 0 task should be Task 1")
		assert.Equal(t, 2, ctx.WaveGroups[1][0].ID, "wave 1 task should be Task 2")
	}

	// has_parallel_opportunity should be true (tasks have Depends/Files)
	assert.True(t, ctx.HasParallelOpp, "has_parallel_opportunity should be true when tasks have depends/files")
}

func TestExecuteWaveModeFlag(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	tasks := spec.TasksFrontmatterV2{
		SpecVersion: "1",
		Total:       1,
		Completed:   0,
		Tasks: []spec.TaskEntry{
			{ID: 1, Name: "Task One", Status: spec.StatusPending},
		},
	}
	setupTestChange(t, tmpDir, tasks, state.PhasePlanned)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"execute", "--context-only", "--execution-mode=wave", "--agent-count=3"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	var ctx executor.ExecutionContext
	require.NoError(t, json.Unmarshal([]byte(output), &ctx), "output must be valid JSON")

	assert.Equal(t, "wave", ctx.ExecutionMode, "execution_mode should be 'wave'")
	assert.Equal(t, 3, ctx.AgentCount, "agent_count should be 3")
}
