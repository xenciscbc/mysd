package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
)

// setupMinimalChangeForPlan creates a minimal spec dir for plan command tests.
// Unlike setupTestChange (execute_test.go), this one is simpler — plan only needs
// proposal.md, design.md, and specs/ to exist.
func setupMinimalChangeForPlan(t *testing.T, dir string) {
	t.Helper()

	specsDir := filepath.Join(dir, ".specs")
	changeDir := filepath.Join(specsDir, "changes", "test-change")

	require.NoError(t, os.MkdirAll(changeDir, 0755))

	// Write .openspec.yaml
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, ".openspec.yaml"),
		[]byte("schema: spec-driven\ncreated: 2026-01-01\n"),
		0644,
	))

	// Write proposal.md
	proposalContent := "---\nspec-version: \"1\"\nchange: test-change\nstatus: proposed\ncreated: 2026-01-01\nupdated: 2026-01-01\n---\n\n## Summary\n\nTest change.\n"
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644))

	// Write design.md
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "design.md"), []byte("## Architecture\n\nTest design.\n"), 0644))

	// Write STATE.json
	ws := state.WorkflowState{
		ChangeName: "test-change",
		Phase:      state.PhaseDesigned,
	}
	require.NoError(t, state.SaveState(specsDir, ws))
}

func TestPlanContextOnly_NewFields(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	setupMinimalChangeForPlan(t, tmpDir)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"plan", "--context-only"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output, "output should not be empty")

	// Validate JSON contains the new fields
	var ctx map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(output), &ctx), "output must be valid JSON")

	// FSCHEMA-05: wave_groups field present
	_, hasWaveGroups := ctx["wave_groups"]
	assert.True(t, hasWaveGroups, "JSON output must contain 'wave_groups' field")

	// FSCHEMA-06: worktree_dir field present with default value ".worktrees"
	worktreeDir, hasWorktreeDir := ctx["worktree_dir"]
	assert.True(t, hasWorktreeDir, "JSON output must contain 'worktree_dir' field")
	assert.Equal(t, ".worktrees", worktreeDir, "worktree_dir should default to '.worktrees'")

	// auto_mode field present and defaults to false
	autoMode, hasAutoMode := ctx["auto_mode"]
	assert.True(t, hasAutoMode, "JSON output must contain 'auto_mode' field")
	assert.Equal(t, false, autoMode, "auto_mode should default to false")

	// Existing fields still present (regression guard)
	assert.Contains(t, ctx, "change_name")
	assert.Contains(t, ctx, "phase")
	assert.Contains(t, ctx, "model")

	// --check not passed: coverage field should not be present
	_, hasCoverage := ctx["coverage"]
	assert.False(t, hasCoverage, "coverage field should not be present without --check flag")
}

func TestPlanContextOnly_ExistingFieldsPresent(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	setupMinimalChangeForPlan(t, tmpDir)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"plan", "--context-only"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	var ctx map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(buf.String()), &ctx))

	// Verify all pre-existing fields from original implementation still exist
	requiredFields := []string{
		"change_name", "phase", "specs", "design",
		"model", "research_enabled", "check_enabled", "test_generation",
	}
	for _, field := range requiredFields {
		assert.Contains(t, ctx, field, "existing field %q must still be present", field)
	}
}

func TestPlanContextOnly_CheckFlag_NoCoverage_WhenNoTasksFile(t *testing.T) {
	// When --check is passed but tasks.md doesn't exist, coverage field should be absent
	// (ParseTasksV2 returns error, planCheck block skips ctx["coverage"] assignment)
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	setupMinimalChangeForPlan(t, tmpDir)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"plan", "--context-only", "--check"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	var ctx map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(buf.String()), &ctx))

	// coverage absent when tasks.md doesn't exist
	_, hasCoverage := ctx["coverage"]
	assert.False(t, hasCoverage, "coverage should be absent when tasks.md does not exist")
}

func TestPlanContextOnly_CheckFlag_WithCoverage(t *testing.T) {
	// When --check is passed and tasks.md exists with satisfies fields, coverage is present
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	setupMinimalChangeForPlan(t, tmpDir)

	// Write tasks.md with a task that has satisfies
	specsDir := filepath.Join(tmpDir, ".specs")
	changeDir := filepath.Join(specsDir, "changes", "test-change")
	fm := spec.TasksFrontmatterV2{
		SpecVersion: "1",
		Total:       1,
		Completed:   0,
		Tasks: []spec.TaskEntry{
			{ID: 1, Name: "Implement auth", Status: spec.StatusPending, Satisfies: []string{"REQ-01"}},
		},
	}
	require.NoError(t, spec.WriteTasks(filepath.Join(changeDir, "tasks.md"), fm, "\n## Tasks\n\n- [ ] Implement auth\n"))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"plan", "--context-only", "--check"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	var ctx map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(buf.String()), &ctx))

	// coverage present when tasks.md exists (even with no MUST IDs, Passed=true)
	coverage, hasCoverage := ctx["coverage"]
	assert.True(t, hasCoverage, "coverage should be present when tasks.md exists and --check passed")
	assert.NotNil(t, coverage)

	coverageMap, ok := coverage.(map[string]interface{})
	require.True(t, ok, "coverage should be a JSON object")
	assert.Contains(t, coverageMap, "passed")
	assert.Contains(t, coverageMap, "total_must")
	assert.Contains(t, coverageMap, "coverage_ratio")
}
