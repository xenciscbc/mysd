package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTaskUpdateFixture creates a temp dir with .specs structure, a STATE.json,
// and a tasks.md with one task. Returns the temp dir path.
func setupTaskUpdateFixture(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	specsDir := filepath.Join(tmpDir, ".specs")
	changeDir := filepath.Join(specsDir, "changes", "my-change")
	require.NoError(t, os.MkdirAll(changeDir, 0755))

	// Write tasks.md with one task
	tasksContent := `---
spec-version: "1"
total: 1
completed: 0
tasks:
  - id: 1
    name: "implement feature"
    status: pending
---

## Tasks
`
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644))

	// Write STATE.json
	ws := state.WorkflowState{
		ChangeName: "my-change",
		Phase:      state.PhasePlanned,
	}
	require.NoError(t, state.SaveState(specsDir, ws))

	return tmpDir
}

func TestTaskUpdate_ValidUpdatesDone(t *testing.T) {
	tmpDir := setupTaskUpdateFixture(t)
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"task-update", "1", "done"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Verify tasks.md was updated
	tasksPath := filepath.Join(tmpDir, ".specs", "changes", "my-change", "tasks.md")
	fm, _, parseErr := spec.ParseTasksV2(tasksPath)
	require.NoError(t, parseErr)
	require.Len(t, fm.Tasks, 1)
	assert.Equal(t, spec.StatusDone, fm.Tasks[0].Status)

	// Check success message
	assert.True(t, strings.Contains(buf.String(), "1"), "output should reference task ID")
}

func TestTaskUpdate_InvalidIDReturnsError(t *testing.T) {
	tmpDir := setupTaskUpdateFixture(t)
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"task-update", "abc", "done"})

	err = rootCmd.Execute()
	assert.Error(t, err, "non-integer task ID should return error")
	assert.True(t, strings.Contains(err.Error(), "invalid task ID"), "error should mention invalid task ID")
}

// setupTaskUpdateFixtureMulti creates a temp dir with multiple tasks for auto-transition tests.
func setupTaskUpdateFixtureMulti(t *testing.T, phase state.Phase) string {
	t.Helper()
	tmpDir := t.TempDir()

	specsDir := filepath.Join(tmpDir, ".specs")
	changeDir := filepath.Join(specsDir, "changes", "my-change")
	require.NoError(t, os.MkdirAll(changeDir, 0755))

	tasksContent := `---
spec-version: "1"
total: 2
completed: 1
tasks:
  - id: 1
    name: "first task"
    status: done
  - id: 2
    name: "second task"
    status: pending
---

## Tasks
`
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644))

	ws := state.WorkflowState{
		ChangeName: "my-change",
		Phase:      phase,
	}
	require.NoError(t, state.SaveState(specsDir, ws))

	return tmpDir
}

func TestTaskUpdate_AutoTransitionOnLastTaskDone(t *testing.T) {
	tmpDir := setupTaskUpdateFixtureMulti(t, state.PhasePlanned)
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"task-update", "2", "done"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Phase should have advanced to executed
	ws, loadErr := state.LoadState(filepath.Join(tmpDir, ".specs"))
	require.NoError(t, loadErr)
	assert.Equal(t, state.PhaseExecuted, ws.Phase)

	// Output should include auto-transition message
	assert.Contains(t, buf.String(), "All tasks complete — phase advanced to executed")
}

func TestTaskUpdate_NoTransitionWhenTasksRemain(t *testing.T) {
	tmpDir := setupTaskUpdateFixtureMulti(t, state.PhasePlanned)
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Mark task 1 as in_progress (task 2 is still pending) — so not all terminal
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"task-update", "1", "in_progress"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Phase should remain planned
	ws, loadErr := state.LoadState(filepath.Join(tmpDir, ".specs"))
	require.NoError(t, loadErr)
	assert.Equal(t, state.PhasePlanned, ws.Phase)

	// Output should NOT include auto-transition message
	assert.NotContains(t, buf.String(), "phase advanced to executed")
}

func TestTaskUpdate_NoTransitionWhenAlreadyExecuted(t *testing.T) {
	tmpDir := setupTaskUpdateFixtureMulti(t, state.PhaseExecuted)
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"task-update", "2", "done"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Phase should remain executed (not double-transition)
	ws, loadErr := state.LoadState(filepath.Join(tmpDir, ".specs"))
	require.NoError(t, loadErr)
	assert.Equal(t, state.PhaseExecuted, ws.Phase)

	// Output should NOT include auto-transition message
	assert.NotContains(t, buf.String(), "phase advanced to executed")
}

func TestTaskUpdate_InvalidStatusReturnsError(t *testing.T) {
	tmpDir := setupTaskUpdateFixture(t)
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"task-update", "1", "invalid_status"})

	err = rootCmd.Execute()
	assert.Error(t, err, "invalid status should return error")
	assert.True(t, strings.Contains(err.Error(), "invalid status"), "error should mention invalid status")
}
