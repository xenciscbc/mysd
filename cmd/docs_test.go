package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDocsListEmpty verifies `mysd docs` with no config prints the empty message.
func TestDocsListEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "No docs_to_update configured")
	assert.Contains(t, out, "mysd docs add")
}

// TestDocsAddAndList verifies `mysd docs add <path>` appends to config and list shows it.
func TestDocsAddAndList(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Add a path
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs", "add", "README.md"})

	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Added: README.md")

	// List — should show README.md
	buf.Reset()
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs"})

	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "README.md")
}

// TestDocsAddMultiple verifies multiple paths can be added and all appear in list.
func TestDocsAddMultiple(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	for _, p := range []string{"README.md", "CHANGELOG.md"} {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&bytes.Buffer{})
		rootCmd.SetArgs([]string{"docs", "add", p})
		require.NoError(t, rootCmd.Execute())
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs"})
	require.NoError(t, rootCmd.Execute())

	out := buf.String()
	assert.Contains(t, out, "README.md")
	assert.Contains(t, out, "CHANGELOG.md")
}

// TestDocsRemove verifies `mysd docs remove <path>` removes path from config.
func TestDocsRemove(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Add then remove
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs", "add", "README.md"})
	require.NoError(t, rootCmd.Execute())

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs", "remove", "README.md"})

	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Removed: README.md")

	// List should be empty
	buf.Reset()
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs"})
	require.NoError(t, rootCmd.Execute())
	assert.Contains(t, buf.String(), "No docs_to_update configured")
}

// TestDocsRemoveNotFound verifies `mysd docs remove <path>` returns error for missing path.
func TestDocsRemoveNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs", "remove", "nonexistent.md"})

	err = rootCmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path not in docs_to_update")
	assert.Contains(t, err.Error(), "nonexistent.md")
}

// TestDocsAddDuplicate verifies adding the same path twice prints "already configured".
func TestDocsAddDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	require.NoError(t, os.MkdirAll(claudeDir, 0755))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// First add
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs", "add", "README.md"})
	require.NoError(t, rootCmd.Execute())

	// Second add (duplicate)
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"docs", "add", "README.md"})

	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "already configured")
}
