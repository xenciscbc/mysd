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

func TestInit_ExistingFile_NoForce_DoesNotOverwrite(t *testing.T) {
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

	// File should not be overwritten
	data, err := os.ReadFile(filepath.Join(claudeDir, "mysd.yaml"))
	require.NoError(t, err)
	assert.Equal(t, existingContent, string(data), "existing config should not be overwritten without --force")

	// Should print a warning
	output := buf.String()
	assert.True(t, strings.Contains(strings.ToLower(output), "exist") ||
		strings.Contains(strings.ToLower(output), "warn") ||
		strings.Contains(strings.ToLower(output), "already"),
		"output should contain warning about existing file, got: %s", output)
}

func TestInit_ExistingFile_WithForce_Overwrites(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create existing config
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	existingContent := "existing: content\n"
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(existingContent), 0644))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"init", "--force"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// File should now have default config content
	data, err := os.ReadFile(filepath.Join(claudeDir, "mysd.yaml"))
	require.NoError(t, err)
	content := string(data)
	assert.NotEqual(t, existingContent, content, "config should be overwritten with --force")
	assert.True(t, strings.Contains(content, "execution_mode"), "overwritten config should contain execution_mode")
}
