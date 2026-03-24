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

func TestStatusOutput(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// 3 tasks: 2 done, 1 pending
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
	rootCmd.SetArgs([]string{"status"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output, "status output should not be empty")
	assert.True(t, strings.Contains(output, "test-change"), "output should contain the change name")

	// Should show task counts — either "2/3" or "2" and "3" separately
	hasTaskCount := strings.Contains(output, "2/3") ||
		(strings.Contains(output, "2") && strings.Contains(output, "3"))
	assert.True(t, hasTaskCount, "output should contain task count indicators, got: %s", output)
}

func TestStatusNoChange(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create a .specs/ dir with empty STATE.json (no active change)
	specsDir := filepath.Join(tmpDir, ".specs")
	require.NoError(t, os.MkdirAll(specsDir, 0755))
	ws := state.WorkflowState{} // empty: no change_name
	require.NoError(t, state.SaveState(specsDir, ws))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"status"})

	// Must not panic — error or graceful message both acceptable
	err = rootCmd.Execute()
	// Even if an error is returned, the test should not panic
	output := buf.String()
	_ = err
	_ = output
	// No assertion on specific content — just verify no panic
}
