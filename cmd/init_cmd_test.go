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

func TestInitStatuslineInstall(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"init"})
	require.NoError(t, rootCmd.Execute())

	// .claude/hooks/mysd-statusline.js should exist
	hookPath := filepath.Join(tmpDir, ".claude", "hooks", "mysd-statusline.js")
	require.FileExists(t, hookPath, "mysd-statusline.js should be installed by init")

	data, err := os.ReadFile(hookPath)
	require.NoError(t, err)
	assert.Equal(t, statuslineHookBytes, data, "hook file content should match embedded bytes")
}

func TestInitStatuslineInstallIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Run init twice — should not error, file should be overwritten
	rootCmd.SetArgs([]string{"init"})
	require.NoError(t, rootCmd.Execute())

	buf.Reset()
	rootCmd.SetArgs([]string{"init"})
	require.NoError(t, rootCmd.Execute())

	hookPath := filepath.Join(tmpDir, ".claude", "hooks", "mysd-statusline.js")
	require.FileExists(t, hookPath)
}

func TestWriteSettingsStatusLine(t *testing.T) {
	claudeDir := t.TempDir()

	require.NoError(t, writeSettingsStatusLine(claudeDir))

	settingsPath := filepath.Join(claudeDir, "settings.json")
	data, err := os.ReadFile(settingsPath)
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))

	statusLine, ok := raw["statusLine"].(map[string]interface{})
	require.True(t, ok, "statusLine key should be present")
	assert.Equal(t, "command", statusLine["type"])
	assert.Equal(t, "node .claude/hooks/mysd-statusline.js", statusLine["command"])
}

func TestWriteSettingsStatusLineMerge(t *testing.T) {
	claudeDir := t.TempDir()

	// Pre-write settings.json with existing hooks key
	existing := map[string]interface{}{
		"hooks": map[string]interface{}{
			"SessionStart": []string{"echo hello"},
		},
	}
	existingData, err := json.MarshalIndent(existing, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "settings.json"), existingData, 0644))

	require.NoError(t, writeSettingsStatusLine(claudeDir))

	data, err := os.ReadFile(filepath.Join(claudeDir, "settings.json"))
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))

	// Both hooks and statusLine keys should be present
	_, hasHooks := raw["hooks"]
	assert.True(t, hasHooks, "existing hooks key should be preserved after merge")

	statusLine, ok := raw["statusLine"].(map[string]interface{})
	require.True(t, ok, "statusLine key should be added")
	assert.Equal(t, "node .claude/hooks/mysd-statusline.js", statusLine["command"])
}

func TestWriteSettingsStatusLineNew(t *testing.T) {
	claudeDir := t.TempDir()
	// No pre-existing settings.json

	require.NoError(t, writeSettingsStatusLine(claudeDir))

	settingsPath := filepath.Join(claudeDir, "settings.json")
	require.FileExists(t, settingsPath, "settings.json should be created when it doesn't exist")

	data, err := os.ReadFile(settingsPath)
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))

	_, ok := raw["statusLine"]
	assert.True(t, ok, "statusLine key should be in new settings.json")
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
