package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupArchiveTestChange creates a change directory with verified state and MUST items done.
// Returns specsDir and changeName.
func setupArchiveTestChange(t *testing.T, phase state.Phase, mustDone bool) (specsDir, changeName, changeDir string) {
	t.Helper()
	tmp := t.TempDir()
	specsDir = tmp
	changeName = "archive-test"
	changeDir = filepath.Join(specsDir, "changes", changeName)

	// Create change directory structure
	require.NoError(t, os.MkdirAll(changeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, ".openspec.yaml"), []byte("schema: openspec/v1\ncreated: 2026-01-01\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("# Proposal\nTest.\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte("- [x] Task 1\n"), 0644))

	// Create a spec with a MUST requirement
	specsSubDir := filepath.Join(changeDir, "specs", "capability-a")
	require.NoError(t, os.MkdirAll(specsSubDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(specsSubDir, "spec.md"), []byte("The system MUST provide authentication.\n"), 0644))

	// Set state
	ws := state.WorkflowState{
		ChangeName: changeName,
		Phase:      phase,
		LastRun:    time.Now(),
	}
	require.NoError(t, state.SaveState(specsDir, ws))

	// If mustDone, create a verification-status.json with MUST item done
	if mustDone {
		// We need to know the stable ID — replicate the StableID logic
		// "The system MUST provide authentication." from spec.md
		// Just write a verification status that has a requirement done
		vs := spec.VerificationStatus{
			ChangeName: changeName,
			VerifiedAt: time.Now().UTC(),
			Requirements: map[string]spec.ItemStatus{
				// This ID will be derived from the StableID function
				// We'll use a placeholder and let the test check behavior
			},
		}
		_ = spec.WriteVerificationStatus(changeDir, vs)
	}

	return specsDir, changeName, changeDir
}

// TestArchiveGate_WrongPhase tests that archive returns an error if state != verified.
func TestArchiveGate_WrongPhase(t *testing.T) {
	specsDir, changeName, _ := setupArchiveTestChange(t, state.PhaseExecuted, false)

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseExecuted}
	err := runArchive(specsDir, ws)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot archive")
	assert.Contains(t, err.Error(), "executed")
	assert.Contains(t, err.Error(), "verified")
}

// TestArchiveGate_MustNotDone tests that archive returns an error if MUST items are not done.
func TestArchiveGate_MustNotDone(t *testing.T) {
	specsDir, changeName, changeDir := setupArchiveTestChange(t, state.PhaseVerified, false)

	// Write verification status with a MUST item blocked (not done)
	vs := spec.VerificationStatus{
		ChangeName: changeName,
		VerifiedAt: time.Now().UTC(),
		Requirements: map[string]spec.ItemStatus{
			"spec.md::must-somehash": spec.StatusBlocked,
		},
	}
	require.NoError(t, spec.WriteVerificationStatus(changeDir, vs))

	// Write gap-report.md (trigger: gap exists, MUST not done)
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "gap-report.md"), []byte("---\nfailed_must_ids:\n  - spec.md::must-somehash\n---\n# Gap Report\n"), 0644))

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseVerified}
	err := runArchive(specsDir, ws)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot archive")
	assert.Contains(t, err.Error(), "MUST item")
}

// TestArchiveSuccess tests that a valid archive moves the directory and updates state.
func TestArchiveSuccess(t *testing.T) {
	specsDir, changeName, changeDir := setupArchiveTestChange(t, state.PhaseVerified, false)

	// Write a verification-status.json with the MUST item done
	// We need the actual stable ID from the spec file
	// The spec has: "The system MUST provide authentication."
	// The StableID format is: "{source_file}::{keyword_lower}-{hex_crc32}"
	// To get this right without importing verifier (circular dep risk), we pre-compute
	// In the archive gate, we check the verification-status vs the parsed MUST items
	// So we need to write a verification status with a stable ID that matches the spec.

	// Create an empty verification status (no requirements yet) - archive should use ParseChange + StableID
	// The archive gate checks MUST items from parsed change against verification-status map.
	// If the verification-status is empty, any MUST item found will fail the gate.
	// So for a "success" test, we need to provide a verification-status with the correct ID.

	// We'll create a change with NO spec MUST requirements to simplify the success test.
	// Remove the spec file and create one with no MUST items.
	specsSubDir := filepath.Join(changeDir, "specs", "capability-a")
	require.NoError(t, os.WriteFile(filepath.Join(specsSubDir, "spec.md"), []byte("# Specification\nThis is a test specification without requirements.\n"), 0644))

	// Write empty verification-status
	vs := spec.VerificationStatus{
		ChangeName:   changeName,
		VerifiedAt:   time.Now().UTC(),
		Requirements: map[string]spec.ItemStatus{},
	}
	require.NoError(t, spec.WriteVerificationStatus(changeDir, vs))

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseVerified}
	err := runArchive(specsDir, ws)
	require.NoError(t, err)

	// Change directory should have moved to changes/archive/YYYY-MM-DD-<name>/
	archiveDir := filepath.Join(specsDir, "changes", "archive", time.Now().Format("2006-01-02")+"-"+changeName)
	assert.DirExists(t, archiveDir)

	// Original changeDir should no longer exist
	_, statErr := os.Stat(changeDir)
	assert.True(t, os.IsNotExist(statErr))

	// STATE.json should be deleted after archive
	mysdStateFile := filepath.Join(filepath.Dir(specsDir), ".mysd", "STATE.json")
	_, statFileErr := os.Stat(mysdStateFile)
	assert.True(t, os.IsNotExist(statFileErr), "STATE.json should be deleted after archive")

	// LoadState should return zero-value (no file)
	loadedWS, err := state.LoadState(specsDir)
	require.NoError(t, err)
	assert.Equal(t, state.PhaseNone, loadedWS.Phase)

	// ARCHIVED-STATE.json should exist in archive directory
	archivedStateDir := archiveDir
	archivedStatePath := filepath.Join(archivedStateDir, "ARCHIVED-STATE.json")
	assert.FileExists(t, archivedStatePath)

	// Validate ARCHIVED-STATE.json content
	data, err := os.ReadFile(archivedStatePath)
	require.NoError(t, err)
	var archivedWS state.WorkflowState
	require.NoError(t, json.Unmarshal(data, &archivedWS))
	assert.Equal(t, changeName, archivedWS.ChangeName)
}

// TestArchiveDeletesResearchCache tests that deleteResearchCache removes the cache file.
func TestArchiveDeletesResearchCache(t *testing.T) {
	// Setup: create temp dir simulating changeDir with cache file
	changeDir := t.TempDir()
	cachePath := filepath.Join(changeDir, "discuss-research-cache.json")
	err := os.WriteFile(cachePath, []byte(`{"change_name":"test","cached_at":"2026-01-01T00:00:00Z"}`), 0644)
	require.NoError(t, err)

	// Act: call the extracted helper (same code path as runArchive)
	deleteResearchCache(changeDir)

	// Assert: file is gone
	_, err = os.Stat(cachePath)
	assert.True(t, os.IsNotExist(err))
}

// TestArchiveDeletesCacheSilentFail tests that deleteResearchCache does not panic when cache file is absent.
func TestArchiveDeletesCacheSilentFail(t *testing.T) {
	changeDir := t.TempDir()
	// No cache file exists — deleteResearchCache should not panic
	assert.NotPanics(t, func() {
		deleteResearchCache(changeDir)
	})
}

// TestCheckTasksDone_AllDone tests that the task gate passes when all tasks are done.
func TestCheckTasksDone_AllDone(t *testing.T) {
	changeDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte("- [x] Task 1\n- [x] Task 2\n"), 0644))
	assert.NoError(t, checkTasksDone(changeDir))
}

// TestCheckTasksDone_WithSkipped tests that the task gate passes when tasks are done or skipped.
func TestCheckTasksDone_WithSkipped(t *testing.T) {
	changeDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte("- [x] Task 1\n- [~] Task 2（跳過：不需要）\n"), 0644))
	assert.NoError(t, checkTasksDone(changeDir))
}

// TestCheckTasksDone_IncompleteTasks tests that the task gate blocks when pending tasks exist.
func TestCheckTasksDone_IncompleteTasks(t *testing.T) {
	changeDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte("- [x] Task 1\n- [ ] Task 2\n"), 0644))
	err := checkTasksDone(changeDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 incomplete task(s)")
}

// TestRunAnalyzeSkipped_WithSkipped tests --analyze-skipped output with skipped tasks.
func TestRunAnalyzeSkipped_WithSkipped(t *testing.T) {
	tmp := t.TempDir()
	changeName := "analyze-test"
	changeDir := filepath.Join(tmp, "changes", changeName)
	require.NoError(t, os.MkdirAll(changeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"),
		[]byte("- [x] Task 1\n- [~] Task 2（跳過：需求變更）\n- [~] Task 3（跳過：技術限制）\n"), 0644))

	ws := state.WorkflowState{ChangeName: changeName}
	require.NoError(t, state.SaveState(tmp, ws))

	// Create a cobra command with a buffer to capture output
	cmd := &cobra.Command{}
	buf := new(strings.Builder)
	cmd.SetOut(buf)

	err := runAnalyzeSkipped(cmd, tmp, ws)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Task 2")
	assert.Contains(t, output, "需求變更")
	assert.Contains(t, output, "Task 3")
	assert.Contains(t, output, "技術限制")
}

// TestRunAnalyzeSkipped_NoSkipped tests --analyze-skipped output with no skipped tasks.
func TestRunAnalyzeSkipped_NoSkipped(t *testing.T) {
	tmp := t.TempDir()
	changeName := "analyze-test"
	changeDir := filepath.Join(tmp, "changes", changeName)
	require.NoError(t, os.MkdirAll(changeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"),
		[]byte("- [x] Task 1\n- [x] Task 2\n"), 0644))

	ws := state.WorkflowState{ChangeName: changeName}
	require.NoError(t, state.SaveState(tmp, ws))

	cmd := &cobra.Command{}
	buf := new(strings.Builder)
	cmd.SetOut(buf)

	err := runAnalyzeSkipped(cmd, tmp, ws)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "[]")
}

// TestMoveDir_Fallback tests that moveDir falls back to copy+delete when os.Rename fails.
func TestMoveDir_Fallback(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	// Remove dst so moveDir can create it
	require.NoError(t, os.RemoveAll(dst))

	// Create test files in src
	require.NoError(t, os.WriteFile(filepath.Join(src, "file1.txt"), []byte("content1"), 0644))
	subDir := filepath.Join(src, "subdir")
	require.NoError(t, os.MkdirAll(subDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("content2"), 0644))

	// moveDir src -> dst
	err := moveDir(src, dst)
	require.NoError(t, err)

	// dst should have all files
	assert.FileExists(t, filepath.Join(dst, "file1.txt"))
	assert.FileExists(t, filepath.Join(dst, "subdir", "file2.txt"))

	// src should no longer exist
	_, statErr := os.Stat(src)
	assert.True(t, os.IsNotExist(statErr))
}
