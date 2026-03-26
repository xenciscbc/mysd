package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// isWindows returns true when running on Windows.
func isWindows() bool {
	return runtime.GOOS == "windows"
}

// TestLangRead_ShowsCurrent verifies that mysd lang with both config files present
// outputs the current response_language and locale values.
func TestLangRead_ShowsCurrent(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	openspecDir := filepath.Join(tmpDir, "openspec")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	require.NoError(t, os.MkdirAll(openspecDir, 0755))

	require.NoError(t, os.WriteFile(
		filepath.Join(claudeDir, "mysd.yaml"),
		[]byte("response_language: zh-TW\n"),
		0644,
	))
	require.NoError(t, os.WriteFile(
		filepath.Join(openspecDir, "config.yaml"),
		[]byte("locale: zh-TW\n"),
		0644,
	))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"lang"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "response_language: zh-TW")
	assert.Contains(t, out, "locale: zh-TW")
}

// TestLangSet_UpdatesBothConfigs verifies that lang set en-US writes response_language
// to .claude/mysd.yaml AND locale to openspec/config.yaml.
func TestLangSet_UpdatesBothConfigs(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	openspecDir := filepath.Join(tmpDir, "openspec")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	require.NoError(t, os.MkdirAll(openspecDir, 0755))

	require.NoError(t, os.WriteFile(
		filepath.Join(claudeDir, "mysd.yaml"),
		[]byte("response_language: zh-TW\n"),
		0644,
	))
	require.NoError(t, os.WriteFile(
		filepath.Join(openspecDir, "config.yaml"),
		[]byte("locale: zh-TW\n"),
		0644,
	))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"lang", "set", "en-US"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Verify .claude/mysd.yaml updated
	mysdData, err := os.ReadFile(filepath.Join(claudeDir, "mysd.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(mysdData), "response_language: en-US")

	// Verify openspec/config.yaml updated
	osData, err := os.ReadFile(filepath.Join(openspecDir, "config.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(osData), "locale: en-US")
}

// TestLangSet_AtomicRollback verifies that if openspec/config.yaml write fails,
// .claude/mysd.yaml is rolled back to its original value.
// Uses a read-only directory to cause the write failure.
// Skipped on Windows where directory permissions work differently.
func TestLangSet_AtomicRollback(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping read-only test: running as root")
	}
	if isWindows() {
		t.Skip("skipping read-only test: Windows directory chmod does not block writes the same way")
	}

	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	openspecDir := filepath.Join(tmpDir, "openspec")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	require.NoError(t, os.MkdirAll(openspecDir, 0755))

	require.NoError(t, os.WriteFile(
		filepath.Join(claudeDir, "mysd.yaml"),
		[]byte("response_language: zh-TW\n"),
		0644,
	))
	require.NoError(t, os.WriteFile(
		filepath.Join(openspecDir, "config.yaml"),
		[]byte("locale: zh-TW\n"),
		0644,
	))

	// Make openspec dir read-only to cause write failure
	require.NoError(t, os.Chmod(openspecDir, 0555))
	defer func() { _ = os.Chmod(openspecDir, 0755) }()

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"lang", "set", "en-US"})

	err = rootCmd.Execute()
	require.Error(t, err, "should return error when openspec write fails")

	// Verify .claude/mysd.yaml was rolled back
	mysdData, err := os.ReadFile(filepath.Join(claudeDir, "mysd.yaml"))
	require.NoError(t, err)
	content := string(mysdData)
	assert.True(t,
		strings.Contains(content, "response_language: zh-TW") || !strings.Contains(content, "response_language: en-US"),
		"mysd.yaml should be rolled back to zh-TW (not en-US), got: %s", content,
	)
}

// TestLangSet_CreatesOpenSpecConfig verifies that if openspec/config.yaml does not exist,
// lang set creates it with the locale set.
func TestLangSet_CreatesOpenSpecConfig(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	// No openspec dir — lang set should create it

	require.NoError(t, os.WriteFile(
		filepath.Join(claudeDir, "mysd.yaml"),
		[]byte("response_language: en-US\n"),
		0644,
	))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"lang", "set", "zh-TW"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Verify openspec/config.yaml was created
	osData, err := os.ReadFile(filepath.Join(tmpDir, "openspec", "config.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(osData), "locale: zh-TW")
}

// TestLangSet_PreservesOtherFields verifies that lang set preserves existing fields
// in openspec/config.yaml (e.g. project name).
func TestLangSet_PreservesOtherFields(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	openspecDir := filepath.Join(tmpDir, "openspec")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))
	require.NoError(t, os.MkdirAll(openspecDir, 0755))

	require.NoError(t, os.WriteFile(
		filepath.Join(claudeDir, "mysd.yaml"),
		[]byte("response_language: en-US\n"),
		0644,
	))
	require.NoError(t, os.WriteFile(
		filepath.Join(openspecDir, "config.yaml"),
		[]byte("project: myproject\nlocale: en-US\n"),
		0644,
	))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"lang", "set", "zh-TW"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	osData, err := os.ReadFile(filepath.Join(openspecDir, "config.yaml"))
	require.NoError(t, err)

	content := string(osData)
	assert.Contains(t, content, "locale: zh-TW", "locale should be updated")
	assert.Contains(t, content, "project: myproject", "project field should be preserved")
}
