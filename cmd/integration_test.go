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
	err := runArchive(specsDir, ws, true) // --yes skips interactive UAT prompt
	require.NoError(t, err)

	// 1. Archive directory exists
	archiveDir := filepath.Join(specsDir, "archive", changeName)
	assert.DirExists(t, archiveDir)

	// 2. Change directory was removed
	_, statErr := os.Stat(changeDir)
	assert.True(t, os.IsNotExist(statErr), "changes/ directory must be removed after archive")

	// 3. ARCHIVED-STATE.json exists in archive dir
	archivedStatePath := filepath.Join(archiveDir, "ARCHIVED-STATE.json")
	assert.FileExists(t, archivedStatePath)

	// 4. STATE.json has phase == archived
	loadedWS, err := state.LoadState(specsDir)
	require.NoError(t, err)
	assert.Equal(t, state.PhaseArchived, loadedWS.Phase)
}

// TestArchiveIntegration_GateRejectsExecuted tests archive fails when state is executed.
func TestArchiveIntegration_GateRejectsExecuted(t *testing.T) {
	specsDir, changeName, _ := setupVerifiedChangeDir(t, state.PhaseExecuted, false)

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseExecuted}
	err := runArchive(specsDir, ws, true)
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
	err := runArchive(specsDir, ws, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not done", "error should mention MUST item is not done")
}

// TestArchiveIntegration_NoUATCheck tests that archive succeeds without any UAT files (UAT-02).
func TestArchiveIntegration_NoUATCheck(t *testing.T) {
	specsDir, changeName, _ := setupVerifiedChangeDir(t, state.PhaseVerified, false)

	// Ensure no .mysd/uat/ directory exists at all
	uatDir := filepath.Join(specsDir, ".mysd", "uat")
	_, statErr := os.Stat(uatDir)
	assert.True(t, os.IsNotExist(statErr), ".mysd/uat/ should not exist before test")

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseVerified}
	// Archive should succeed even without any UAT files
	err := runArchive(specsDir, ws, true) // --yes skips interactive UAT prompt
	assert.NoError(t, err, "archive should succeed regardless of UAT file absence")

	// Verify archive completed
	archiveDir := filepath.Join(specsDir, "archive", changeName)
	assert.DirExists(t, archiveDir)
}
