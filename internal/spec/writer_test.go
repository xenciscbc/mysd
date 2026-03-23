package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScaffold_CreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()

	c, err := Scaffold("my-new-feature", tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "my-new-feature", c.Name)

	changeDir := filepath.Join(tmpDir, "changes", "my-new-feature")

	// Check .openspec.yaml
	metaPath := filepath.Join(changeDir, ".openspec.yaml")
	assert.FileExists(t, metaPath)
	metaBytes, _ := os.ReadFile(metaPath)
	assert.Contains(t, string(metaBytes), "spec-driven")

	// Check proposal.md
	proposalPath := filepath.Join(changeDir, "proposal.md")
	assert.FileExists(t, proposalPath)
	proposalBytes, _ := os.ReadFile(proposalPath)
	proposalContent := string(proposalBytes)
	assert.Contains(t, proposalContent, `spec-version: "1"`)
	assert.Contains(t, proposalContent, "my-new-feature")
	assert.Contains(t, proposalContent, "proposed")
	assert.Contains(t, proposalContent, "## Summary")

	// Check specs/ directory exists
	specsDir := filepath.Join(changeDir, "specs")
	info, err := os.Stat(specsDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Check design.md
	designPath := filepath.Join(changeDir, "design.md")
	assert.FileExists(t, designPath)
	designBytes, _ := os.ReadFile(designPath)
	assert.Contains(t, string(designBytes), "## Architecture")

	// Check tasks.md
	tasksPath := filepath.Join(changeDir, "tasks.md")
	assert.FileExists(t, tasksPath)
	tasksBytes, _ := os.ReadFile(tasksPath)
	tasksContent := string(tasksBytes)
	assert.Contains(t, tasksContent, `spec-version: "1"`)
}

func TestScaffold_ReturnedChangeHasCorrectDir(t *testing.T) {
	tmpDir := t.TempDir()
	c, err := Scaffold("feature-xyz", tmpDir)
	require.NoError(t, err)

	expectedDir := filepath.Join(tmpDir, "changes", "feature-xyz")
	assert.Equal(t, expectedDir, c.Dir)
}
