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

// TestModelRead_DefaultProfile verifies that mysd model with no config file
// outputs "Profile: balanced" header followed by 10 role-model rows.
func TestModelRead_DefaultProfile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"model"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Profile: balanced")
}

// TestModelRead_ContainsAllRoles verifies that all 10 agent roles appear in the output.
func TestModelRead_ContainsAllRoles(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"model"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	roles := []string{
		"spec-writer", "designer", "planner", "executor", "verifier",
		"fast-forward", "researcher", "advisor", "proposal-writer", "plan-checker",
	}
	for _, role := range roles {
		assert.Contains(t, out, role, "output should contain role %q", role)
	}
}

// TestModelRead_NonTTY verifies that in non-TTY mode, output uses plain text format
// with space separation and no ANSI escape codes.
func TestModelRead_NonTTY(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"model"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	// bytes.Buffer is not a TTY fd, so no ANSI codes should appear
	assert.NotContains(t, out, "\x1b[", "output should not contain ANSI escape codes in non-TTY mode")
	// Should have columnar output
	assert.True(t, strings.Contains(out, "Role") || strings.Contains(out, "spec-writer"),
		"output should contain role column")
}

// TestModelSet_ValidProfile verifies that model set quality writes model_profile to
// .claude/mysd.yaml and outputs a success message.
func TestModelSet_ValidProfile(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	// Write initial config with tdd: true
	initialConfig := "tdd: true\n"
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(initialConfig), 0644))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"model", "set", "quality"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Read back the config file
	data, err := os.ReadFile(filepath.Join(claudeDir, "mysd.yaml"))
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "model_profile: quality", "should write model_profile")
	assert.Contains(t, content, "tdd: true", "should preserve existing tdd field")
}

// TestModelSet_InvalidProfile verifies that model set with unknown profile name
// returns an error containing "unknown profile" and lists valid profiles.
func TestModelSet_InvalidProfile(t *testing.T) {
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(errBuf)
	rootCmd.SetArgs([]string{"model", "set", "fast"})

	err := rootCmd.Execute()
	require.Error(t, err)

	errMsg := err.Error()
	assert.Contains(t, errMsg, "unknown profile", "error should contain 'unknown profile'")
	assert.Contains(t, errMsg, "quality", "error should list valid profiles")
	assert.Contains(t, errMsg, "balanced", "error should list valid profiles")
	assert.Contains(t, errMsg, "budget", "error should list valid profiles")
}

// TestModelSet_PreservesOtherConfig verifies that model set preserves other config fields.
func TestModelSet_PreservesOtherConfig(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	// Write initial config with multiple fields
	initialConfig := "tdd: true\nexecution_mode: wave\n"
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(initialConfig), 0644))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"model", "set", "quality"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(claudeDir, "mysd.yaml"))
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "model_profile: quality")
	assert.Contains(t, content, "tdd: true")
	assert.Contains(t, content, "execution_mode: wave")
}

// TestModelSet_CustomProfile verifies that model set accepts a custom profile name.
func TestModelSet_CustomProfile(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	initialConfig := `custom_profiles:
  my-team:
    base: balanced
    models:
      executor: opus
`
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(initialConfig), 0644))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"model", "set", "my-team"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(claudeDir, "mysd.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "model_profile: my-team")
}

// TestModelSet_UnknownProfile_ListsCustomProfiles verifies error lists both built-in and custom profiles.
func TestModelSet_UnknownProfile_WithCustom(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	initialConfig := `custom_profiles:
  my-team:
    base: balanced
    models:
      executor: opus
`
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(initialConfig), 0644))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(errBuf)
	rootCmd.SetArgs([]string{"model", "set", "nonexistent"})

	err = rootCmd.Execute()
	require.Error(t, err)

	errMsg := err.Error()
	assert.Contains(t, errMsg, "unknown profile")
	assert.Contains(t, errMsg, "my-team")
}

// TestModelRead_CustomProfile verifies model display works with a custom profile.
func TestModelRead_CustomProfile(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	configContent := `model_profile: my-team
custom_profiles:
  my-team:
    base: balanced
    models:
      executor: opus
`
	require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(configContent), 0644))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"model"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Profile: my-team")
	assert.Contains(t, out, "executor")
}
