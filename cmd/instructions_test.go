package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupInstructionsTestChange(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	changeDir := filepath.Join(tmp, "changes", "test-inst")

	// proposal.md
	require.NoError(t, os.MkdirAll(changeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("## Why\n\nTest\n"), 0644))

	// specs dir with one capability
	specDir := filepath.Join(changeDir, "specs", "auth")
	require.NoError(t, os.MkdirAll(specDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(specDir, "spec.md"), []byte("### Requirement: auth\n"), 0644))

	// design.md
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "design.md"), []byte("## Decisions\n"), 0644))

	return tmp
}

func TestBuildDesignInstructions(t *testing.T) {
	specsDir := setupInstructionsTestChange(t)
	changeDir := filepath.Join(specsDir, "changes", "test-inst")

	out := buildDesignInstructions("test-inst", specsDir, changeDir)

	assert.Equal(t, "design", out.ArtifactID)
	assert.Equal(t, "test-inst", out.ChangeName)
	assert.Contains(t, out.OutputPath, "design.md")
	assert.NotEmpty(t, out.Template)
	assert.NotEmpty(t, out.Rules)
	assert.NotEmpty(t, out.Instruction)
	assert.NotEmpty(t, out.SelfReviewChecklist)

	// Dependencies
	require.Len(t, out.Dependencies, 2)
	assert.Equal(t, "proposal", out.Dependencies[0].ID)
	assert.True(t, out.Dependencies[0].Done, "proposal should be done")
	assert.Equal(t, "specs", out.Dependencies[1].ID)
	assert.True(t, out.Dependencies[1].Done, "specs should be done")
}

func TestBuildTasksInstructions(t *testing.T) {
	specsDir := setupInstructionsTestChange(t)
	changeDir := filepath.Join(specsDir, "changes", "test-inst")

	out := buildTasksInstructions("test-inst", specsDir, changeDir)

	assert.Equal(t, "tasks", out.ArtifactID)
	assert.Equal(t, "test-inst", out.ChangeName)
	assert.Contains(t, out.OutputPath, "tasks.md")
	assert.NotEmpty(t, out.Template)
	assert.Contains(t, out.Template, "spec:")
	assert.NotEmpty(t, out.Rules)
	assert.NotEmpty(t, out.SelfReviewChecklist)

	// Dependencies
	require.Len(t, out.Dependencies, 3)
	assert.Equal(t, "proposal", out.Dependencies[0].ID)
	assert.True(t, out.Dependencies[0].Done)
	assert.Equal(t, "specs", out.Dependencies[1].ID)
	assert.True(t, out.Dependencies[1].Done)
	assert.Equal(t, "design", out.Dependencies[2].ID)
	assert.True(t, out.Dependencies[2].Done)
}

func TestBuildTasksInstructions_MissingDesign(t *testing.T) {
	specsDir := setupInstructionsTestChange(t)
	changeDir := filepath.Join(specsDir, "changes", "test-inst")

	// Remove design.md
	os.Remove(filepath.Join(changeDir, "design.md"))

	out := buildTasksInstructions("test-inst", specsDir, changeDir)

	designDep := out.Dependencies[2]
	assert.Equal(t, "design", designDep.ID)
	assert.False(t, designDep.Done, "design should not be done when file is missing")
}

func TestRunInstructions_UnknownArtifact(t *testing.T) {
	// Test the validation logic directly since cobra command state is global
	buf := new(bytes.Buffer)
	cmd := instructionsCmd
	cmd.SetOut(buf)

	err := runInstructions(cmd, []string{"unknown"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown artifact ID")
}

func TestRunInstructions_JSONOutput(t *testing.T) {
	specsDir := setupInstructionsTestChange(t)

	// Override the working directory detection by running buildInstructions directly
	changeDir := filepath.Join(specsDir, "changes", "test-inst")
	out := buildInstructions("design", "test-inst", specsDir, changeDir)

	data, err := json.MarshalIndent(out, "", "  ")
	require.NoError(t, err)

	var parsed InstructionsOutput
	require.NoError(t, json.Unmarshal(data, &parsed))

	assert.Equal(t, "design", parsed.ArtifactID)
	assert.Equal(t, "test-inst", parsed.ChangeName)
	assert.NotEmpty(t, parsed.Template)
	assert.NotEmpty(t, parsed.Rules)
	assert.NotEmpty(t, parsed.SelfReviewChecklist)
}
