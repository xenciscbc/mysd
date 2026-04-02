package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupVerifiedChangeDir creates a temp specsDir with a change in the given phase.
// By default the spec has no MUST requirements, so archive gate passes without
// needing to compute stable IDs. Pass withMust=true to include a MUST requirement.
func setupVerifiedChangeDir(t *testing.T, phase state.Phase, withMust bool) (specsDir, changeName, changeDir string) {
	t.Helper()
	tmp := t.TempDir()
	specsDir = tmp
	changeName = "archive-integration-change"
	changeDir = filepath.Join(specsDir, "changes", changeName)

	specsSubDir := filepath.Join(changeDir, "specs", "capability-a")
	require.NoError(t, os.MkdirAll(specsSubDir, 0755))

	// .openspec.yaml
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, ".openspec.yaml"),
		[]byte("schema: openspec/v1\ncreated: 2026-01-01\n"),
		0644,
	))
	// proposal.md
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte("# Proposal\nIntegration test.\n"),
		0644,
	))
	// tasks.md
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "tasks.md"),
		[]byte("- [x] Task One\n"),
		0644,
	))

	if withMust {
		// spec.md with one MUST requirement
		require.NoError(t, os.WriteFile(
			filepath.Join(specsSubDir, "spec.md"),
			[]byte("The system MUST provide authentication.\n"),
			0644,
		))
	} else {
		// spec.md with no MUST requirements
		require.NoError(t, os.WriteFile(
			filepath.Join(specsSubDir, "spec.md"),
			[]byte("# Specification\nThis is a test specification without requirements.\n"),
			0644,
		))
	}

	// Save state
	ws := state.WorkflowState{
		ChangeName: changeName,
		Phase:      phase,
		LastRun:    time.Now(),
	}
	require.NoError(t, state.SaveState(specsDir, ws))

	// Write empty verification-status (archive gate reads this)
	vs := spec.VerificationStatus{
		ChangeName:   changeName,
		VerifiedAt:   time.Now().UTC(),
		Requirements: map[string]spec.ItemStatus{},
	}
	require.NoError(t, spec.WriteVerificationStatus(changeDir, vs))

	return specsDir, changeName, changeDir
}

// TestArchiveIntegration_Success tests the full archive pipeline succeeds with verified state.
func TestArchiveIntegration_Success(t *testing.T) {
	specsDir, changeName, changeDir := setupVerifiedChangeDir(t, state.PhaseVerified, false)

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseVerified}
	err := runArchive(specsDir, ws)
	require.NoError(t, err)

	// 1. Archive directory exists (with date prefix)
	archiveDir := filepath.Join(specsDir, "changes", "archive", time.Now().Format("2006-01-02")+"-"+changeName)
	assert.DirExists(t, archiveDir)

	// 2. Change directory was removed
	_, statErr := os.Stat(changeDir)
	assert.True(t, os.IsNotExist(statErr), "changes/ directory must be removed after archive")

	// 3. ARCHIVED-STATE.json exists in archive dir
	archivedStatePath := filepath.Join(archiveDir, "ARCHIVED-STATE.json")
	assert.FileExists(t, archivedStatePath)

	// 4. STATE.json should be deleted after archive
	loadedWS, err := state.LoadState(specsDir)
	require.NoError(t, err)
	assert.Equal(t, state.PhaseNone, loadedWS.Phase, "STATE.json should be cleaned up after archive")
}

// TestArchiveIntegration_GateRejectsExecuted tests archive fails when state is executed.
func TestArchiveIntegration_GateRejectsExecuted(t *testing.T) {
	specsDir, changeName, _ := setupVerifiedChangeDir(t, state.PhaseExecuted, false)

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseExecuted}
	err := runArchive(specsDir, ws)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be verified", "error should mention 'must be verified'")
}

// TestArchiveIntegration_GateRejectsMustNotDone tests archive fails when MUST item is blocked.
func TestArchiveIntegration_GateRejectsMustNotDone(t *testing.T) {
	specsDir, changeName, changeDir := setupVerifiedChangeDir(t, state.PhaseVerified, false)

	// Overwrite the spec to have a MUST requirement
	specsSubDir := filepath.Join(changeDir, "specs", "capability-a")
	require.NoError(t, os.WriteFile(
		filepath.Join(specsSubDir, "spec.md"),
		[]byte("The system MUST provide authentication.\n"),
		0644,
	))

	// The verification-status already has an empty requirements map (from setupVerifiedChangeDir).
	// checkMustItemsDone will find the MUST item in the parsed spec but NOT find it in the
	// empty verification map, triggering "not done" error.

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseVerified}
	err := runArchive(specsDir, ws)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not done", "error should mention MUST item is not done")
}

// TestArchiveIntegration_DeltaSpecMerge tests the full pipeline:
// create change with delta specs → archive → verify main specs merged, archive path correct, skipped tasks handled.
func TestArchiveIntegration_DeltaSpecMerge(t *testing.T) {
	tmp := t.TempDir()
	specsDir := tmp
	changeName := "e2e-merge-test"
	changeDir := filepath.Join(specsDir, "changes", changeName)

	// Create change directory structure
	deltaSpecDir := filepath.Join(changeDir, "specs", "auth")
	require.NoError(t, os.MkdirAll(deltaSpecDir, 0755))

	// .openspec.yaml
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, ".openspec.yaml"),
		[]byte("schema: openspec/v1\ncreated: 2026-04-02\n"),
		0644,
	))

	// proposal.md
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte("# Proposal\nE2E merge test.\n"),
		0644,
	))

	// tasks.md with completed and skipped tasks
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "tasks.md"),
		[]byte("- [x] Task 1\n- [x] Task 2\n- [~] Task 3（跳過：需求變更）\n"),
		0644,
	))

	// Delta spec with ADDED requirement (use SHOULD to avoid MUST gate complexity)
	require.NoError(t, os.WriteFile(
		filepath.Join(deltaSpecDir, "spec.md"),
		[]byte("## ADDED Requirements\n\nThe system SHOULD support OAuth2 authentication.\n"),
		0644,
	))

	// Create a main specs directory (empty — simulates no existing main spec)
	mainSpecsDir := filepath.Join(specsDir, "specs", "auth")
	require.NoError(t, os.MkdirAll(mainSpecsDir, 0755))

	// Save state as verified
	ws := state.WorkflowState{
		ChangeName: changeName,
		Phase:      state.PhaseVerified,
		LastRun:    time.Now(),
	}
	require.NoError(t, state.SaveState(specsDir, ws))

	// Write empty verification-status
	vs := spec.VerificationStatus{
		ChangeName:   changeName,
		VerifiedAt:   time.Now().UTC(),
		Requirements: map[string]spec.ItemStatus{},
	}
	require.NoError(t, spec.WriteVerificationStatus(changeDir, vs))

	// Run archive
	err := runArchive(specsDir, ws)
	require.NoError(t, err)

	// 1. Archive directory exists with date prefix
	archiveDir := filepath.Join(specsDir, "changes", "archive", time.Now().Format("2006-01-02")+"-"+changeName)
	assert.DirExists(t, archiveDir)

	// 2. Original change directory is gone
	_, statErr := os.Stat(changeDir)
	assert.True(t, os.IsNotExist(statErr))

	// 3. Main spec was created with merged content
	mainSpecPath := filepath.Join(mainSpecsDir, "spec.md")
	assert.FileExists(t, mainSpecPath)
	content, err := os.ReadFile(mainSpecPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "SHOULD support OAuth2 authentication")
	assert.Contains(t, string(content), "version: 1.0.0") // new spec gets initial frontmatter

	// 4. STATE.json should be deleted after archive
	loadedWS, err := state.LoadState(specsDir)
	require.NoError(t, err)
	assert.Equal(t, state.PhaseNone, loadedWS.Phase, "STATE.json should be cleaned up after archive")
}

// TestArchiveIntegration_TaskGateBlocksIncomplete tests that archive fails when tasks are incomplete.
func TestArchiveIntegration_TaskGateBlocksIncomplete(t *testing.T) {
	specsDir, changeName, changeDir := setupVerifiedChangeDir(t, state.PhaseVerified, false)

	// Overwrite tasks.md with an incomplete task
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "tasks.md"),
		[]byte("- [x] Task 1\n- [ ] Task 2\n"),
		0644,
	))

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseVerified}
	err := runArchive(specsDir, ws)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete task")
}

// TestArchiveIntegration_SkippedTasksPassGate tests that archive succeeds with skipped tasks.
func TestArchiveIntegration_SkippedTasksPassGate(t *testing.T) {
	specsDir, changeName, changeDir := setupVerifiedChangeDir(t, state.PhaseVerified, false)

	// Overwrite tasks.md with completed and skipped tasks
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "tasks.md"),
		[]byte("- [x] Task 1\n- [~] Task 2（跳過：不需要）\n"),
		0644,
	))

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseVerified}
	err := runArchive(specsDir, ws)
	require.NoError(t, err)

	archiveDir := filepath.Join(specsDir, "changes", "archive", time.Now().Format("2006-01-02")+"-"+changeName)
	assert.DirExists(t, archiveDir)
}
