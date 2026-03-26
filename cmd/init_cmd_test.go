package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit_CreatesConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"init"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	configPath := filepath.Join(tmpDir, ".claude", "mysd.yaml")
	require.FileExists(t, configPath)

	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	content := string(data)
	assert.True(t, strings.Contains(content, "execution_mode"), "config should contain execution_mode")
}

func TestInit_CreatesOpenspecStructure(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"init"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Check openspec/ and openspec/specs/ were created
	info, statErr := os.Stat(filepath.Join(tmpDir, "openspec"))
	require.NoError(t, statErr)
	assert.True(t, info.IsDir())

	info, statErr = os.Stat(filepath.Join(tmpDir, "openspec", "specs"))
	require.NoError(t, statErr)
	assert.True(t, info.IsDir())

	// openspec/config.yaml should NOT be created by init (per D-06)
	_, statErr = os.Stat(filepath.Join(tmpDir, "openspec", "config.yaml"))
	assert.True(t, os.IsNotExist(statErr), "init should not create openspec/config.yaml")
}

func TestInit_ExistingFile_DoesNotOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create .claude/ and an existing config
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	existingContent := "existing: content\n"
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(existingContent), 0644))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"init"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// File should not be overwritten (idempotent)
	data, err := os.ReadFile(filepath.Join(claudeDir, "mysd.yaml"))
	require.NoError(t, err)
	assert.Equal(t, existingContent, string(data), "existing config should not be overwritten by init (idempotent)")
}

func TestInit_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Run init twice — should not error either time
	rootCmd.SetArgs([]string{"init"})
	require.NoError(t, rootCmd.Execute())

	buf.Reset()
	rootCmd.SetArgs([]string{"init"})
	require.NoError(t, rootCmd.Execute())

	// openspec/ should still exist
	info, statErr := os.Stat(filepath.Join(tmpDir, "openspec"))
	require.NoError(t, statErr)
	assert.True(t, info.IsDir())
}
