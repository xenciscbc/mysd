package uat_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mysd/internal/uat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewUATChecklist verifies NewUATChecklist creates a checklist with correct defaults.
func TestNewUATChecklist(t *testing.T) {
	items := []uat.UATItem{
		{Description: "User can see the login button on homepage"},
		{Description: "User can submit the login form"},
	}

	cl := uat.NewUATChecklist("my-feature", items)

	assert.Equal(t, "1", cl.SpecVersion)
	assert.Equal(t, "my-feature", cl.Change)
	assert.NotEmpty(t, cl.Generated)
	assert.Equal(t, 2, cl.Summary.Total)
	assert.Equal(t, 0, cl.Summary.Pass)
	assert.Equal(t, 0, cl.Summary.Fail)
	assert.Equal(t, 0, cl.Summary.Skip)
	assert.Len(t, cl.Results, 2)

	// All items should default to "pending" status
	for _, item := range cl.Results {
		assert.Equal(t, "pending", item.Status)
	}

	// IDs should be uat-1, uat-2, ...
	assert.Equal(t, "uat-1", cl.Results[0].ID)
	assert.Equal(t, "uat-2", cl.Results[1].ID)
}

// TestWriteUAT_CreatesDirectory verifies WriteUAT creates the .mysd/uat/ directory.
func TestWriteUAT_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := uat.UATFilePath(tmpDir, "my-feature")

	items := []uat.UATItem{
		{Description: "User can see the login button"},
	}
	cl := uat.NewUATChecklist("my-feature", items)

	err := uat.WriteUAT(filePath, cl)
	require.NoError(t, err)

	// Directory must be created
	assert.DirExists(t, filepath.Dir(filePath))
	assert.FileExists(t, filePath)
}

// TestWriteUAT_FrontmatterFormat verifies UAT file has correct YAML frontmatter.
func TestWriteUAT_FrontmatterFormat(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := uat.UATFilePath(tmpDir, "my-feature")

	items := []uat.UATItem{
		{Description: "User can see the login button on homepage"},
	}
	cl := uat.NewUATChecklist("my-feature", items)

	err := uat.WriteUAT(filePath, cl)
	require.NoError(t, err)

	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	body := string(content)
	assert.Contains(t, body, `spec-version: "1"`)
	assert.Contains(t, body, "change: my-feature")
	assert.Contains(t, body, "generated:")
	assert.Contains(t, body, "summary:")
	assert.Contains(t, body, "total: 1")
}

// TestWriteUAT_MarkdownBody verifies UAT file has correct markdown body.
func TestWriteUAT_MarkdownBody(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := uat.UATFilePath(tmpDir, "my-feature")

	items := []uat.UATItem{
		{Description: "User can see the login button on homepage"},
		{Description: "User can submit the login form"},
	}
	cl := uat.NewUATChecklist("my-feature", items)

	err := uat.WriteUAT(filePath, cl)
	require.NoError(t, err)

	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	body := string(content)
	assert.Contains(t, body, "## UAT Checklist: my-feature")
	assert.Contains(t, body, "- [ ] User can see the login button on homepage")
	assert.Contains(t, body, "- [ ] User can submit the login form")
}

// TestWriteUAT_PreservesHistory verifies that WriteUAT appends to run_history
// instead of overwriting it (Pitfall 4).
func TestWriteUAT_PreservesHistory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := uat.UATFilePath(tmpDir, "my-feature")

	// First write
	items := []uat.UATItem{
		{ID: "uat-1", Description: "User can see the login button", Status: "pass"},
	}
	cl1 := uat.NewUATChecklist("my-feature", items)
	cl1.Results = items
	cl1.Summary = uat.UATSummary{Total: 1, Pass: 1}
	err := uat.WriteUAT(filePath, cl1)
	require.NoError(t, err)

	// Second write with different results
	items2 := []uat.UATItem{
		{ID: "uat-1", Description: "User can see the login button", Status: "fail"},
	}
	cl2 := uat.NewUATChecklist("my-feature", items2)
	cl2.Results = items2
	cl2.Summary = uat.UATSummary{Total: 1, Fail: 1}
	err = uat.WriteUAT(filePath, cl2)
	require.NoError(t, err)

	// Read back and verify history is preserved
	result, err := uat.ReadUAT(filePath)
	require.NoError(t, err)

	// run_history should have the first run's results
	assert.Len(t, result.RunHistory, 1, "first run should be in run_history")
	assert.Equal(t, 1, result.RunHistory[0].Summary.Pass)
}

// TestReadUAT verifies ReadUAT can read a previously written UAT file.
func TestReadUAT(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := uat.UATFilePath(tmpDir, "my-feature")

	items := []uat.UATItem{
		{Description: "User can see the login button"},
		{Description: "User can submit the login form"},
	}
	cl := uat.NewUATChecklist("my-feature", items)

	err := uat.WriteUAT(filePath, cl)
	require.NoError(t, err)

	result, err := uat.ReadUAT(filePath)
	require.NoError(t, err)

	assert.Equal(t, "1", result.SpecVersion)
	assert.Equal(t, "my-feature", result.Change)
	assert.Equal(t, 2, result.Summary.Total)
	assert.Len(t, result.Results, 2)
}

// TestReadUAT_FileNotExist verifies ReadUAT returns zero-value (no error) for missing files.
func TestReadUAT_FileNotExist(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "nonexistent-uat.md")

	result, err := uat.ReadUAT(filePath)
	assert.NoError(t, err, "missing file should not be an error")
	assert.Equal(t, uat.UATChecklist{}, result, "missing file should return zero-value")
}

// TestUATFilePath verifies the file path convention.
func TestUATFilePath(t *testing.T) {
	path := uat.UATFilePath("/project", "my-feature")
	// Normalize separators for cross-platform comparison
	normalized := strings.ReplaceAll(path, "\\", "/")
	assert.Equal(t, "/project/.mysd/uat/my-feature-uat.md", normalized)
}
