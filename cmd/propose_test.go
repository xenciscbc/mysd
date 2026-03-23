package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetRootCmd() {
	rootCmd.ResetFlags()
	rootCmd.ResetCommands()
	// Re-register all commands and flags by re-running init functions is not
	// easily possible. Instead we reuse the package-level rootCmd and just
	// reset its args/output between tests.
}

func TestPropose_CreatesDirectoryStructure(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"propose", "test-feature"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Check directory structure
	changeDir := filepath.Join(tmpDir, ".specs", "changes", "test-feature")
	assert.DirExists(t, changeDir)
	assert.FileExists(t, filepath.Join(changeDir, "proposal.md"))
	assert.FileExists(t, filepath.Join(changeDir, "design.md"))
	assert.FileExists(t, filepath.Join(changeDir, "tasks.md"))
	assert.DirExists(t, filepath.Join(changeDir, "specs"))
}

func TestPropose_PrintsSuccessMessage(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"propose", "my-feature"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, strings.Contains(output, "my-feature"), "output should contain the feature name")
}

func TestPropose_NoArgs_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"propose"})

	err := rootCmd.Execute()
	assert.Error(t, err, "propose with no args should return an error")
}

func TestPropose_CreatesStateJSON(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"propose", "state-feature"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// STATE.json should be created in .specs/
	stateFile := filepath.Join(tmpDir, ".specs", "STATE.json")
	require.FileExists(t, stateFile)

	data, err := os.ReadFile(stateFile)
	require.NoError(t, err)

	var ws map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &ws))
	assert.Equal(t, "proposed", ws["phase"])
	assert.Equal(t, "state-feature", ws["change_name"])
}
