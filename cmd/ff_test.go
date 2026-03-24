package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mysd/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFFStateTransitions(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"ff", "test-feature"})

	err = rootCmd.Execute()
	require.NoError(t, err, "ff should complete without error")

	// Verify STATE.json exists and has phase == planned
	specsDir := filepath.Join(tmpDir, ".specs")
	ws, loadErr := state.LoadState(specsDir)
	require.NoError(t, loadErr)
	assert.Equal(t, state.PhasePlanned, ws.Phase, "ff should end at planned phase")
	assert.Equal(t, "test-feature", ws.ChangeName)

	// Verify proposal.md exists (scaffold ran)
	proposalPath := filepath.Join(specsDir, "changes", "test-feature", "proposal.md")
	assert.FileExists(t, proposalPath, "ff should create proposal.md via Scaffold")
}

func TestFFEStateTransitions(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"ffe", "test-feature"})

	err = rootCmd.Execute()
	require.NoError(t, err, "ffe should complete without error")

	// Verify STATE.json exists and has phase == executed
	specsDir := filepath.Join(tmpDir, ".specs")
	ws, loadErr := state.LoadState(specsDir)
	require.NoError(t, loadErr)
	assert.Equal(t, state.PhaseExecuted, ws.Phase, "ffe should end at executed phase")
	assert.Equal(t, "test-feature", ws.ChangeName)
}

func TestFFAlreadyProposed(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Pre-create state with phase == proposed (change already went through propose step)
	specsDir := filepath.Join(tmpDir, ".specs")
	require.NoError(t, os.MkdirAll(filepath.Join(specsDir, "changes", "test-feature"), 0755))
	ws := state.WorkflowState{
		ChangeName: "test-feature",
		Phase:      state.PhaseProposed,
	}
	require.NoError(t, state.SaveState(specsDir, ws))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"ff", "test-feature"})

	err = rootCmd.Execute()
	// ff should return an error: state.Transition rejects proposed -> proposed
	assert.Error(t, err, "ff on already-proposed change should return an error")

	errStr := strings.ToLower(err.Error())
	hasRelevantMsg := strings.Contains(errStr, "already") ||
		strings.Contains(errStr, "cannot") ||
		strings.Contains(errStr, "invalid") ||
		strings.Contains(errStr, "transition")
	assert.True(t, hasRelevantMsg, "error should indicate invalid transition, got: %s", err.Error())
}
