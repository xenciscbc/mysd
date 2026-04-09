package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/xenciscbc/mysd/internal/verifier"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupVerifyTestChange creates a minimal change directory in a temp specsDir
// with a valid state in "executed" phase, and one MUST requirement.
func setupVerifyTestChange(t *testing.T) (specsDir, changeName, changeDir string) {
	t.Helper()
	tmp := t.TempDir()
	specsDir = tmp
	changeName = "test-change"
	changeDir = filepath.Join(specsDir, "changes", changeName)

	// Create .openspec.yaml
	require.NoError(t, os.MkdirAll(changeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, ".openspec.yaml"), []byte("schema: openspec/v1\ncreated: 2026-01-01\n"), 0644))

	// Create proposal.md
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("# Proposal\nTest proposal.\n"), 0644))

	// Create tasks.md
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte("- [x] Task 1\n"), 0644))

	// Create a spec with a MUST requirement
	specsSubDir := filepath.Join(changeDir, "specs", "capability-a")
	require.NoError(t, os.MkdirAll(specsSubDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(specsSubDir, "spec.md"), []byte("The system MUST provide authentication.\n"), 0644))

	// Write STATE.json with executed phase
	ws := state.WorkflowState{
		ChangeName: changeName,
		Phase:      state.PhaseExecuted,
	}
	require.NoError(t, state.SaveState(specsDir, ws))

	return specsDir, changeName, changeDir
}

// TestVerifyContextOnly tests that --context-only outputs valid JSON with required keys.
func TestVerifyContextOnly(t *testing.T) {
	specsDir, changeName, _ := setupVerifyTestChange(t)

	// Save working dir and change to specsDir parent so DetectSpecDir finds it
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(filepath.Dir(specsDir)))

	// Directly call the function with the specsDir and state
	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseExecuted}
	var buf bytes.Buffer
	err = runVerifyContextOnly(&buf, specsDir, ws, config.Defaults())
	require.NoError(t, err)

	// Parse output as JSON
	var ctx verifier.VerificationContext
	require.NoError(t, json.Unmarshal(buf.Bytes(), &ctx))

	assert.Equal(t, changeName, ctx.ChangeName)
	assert.NotNil(t, ctx.MustItems)
	assert.NotNil(t, ctx.ShouldItems)
	assert.NotNil(t, ctx.MayItems)
	// Should have at least 1 MUST item from spec.md
	assert.Len(t, ctx.MustItems, 1)
	assert.Equal(t, "MUST", ctx.MustItems[0].Keyword)
}

// TestVerifyContextOnly_NoChange tests that context-only with empty ChangeName returns error.
func TestVerifyContextOnly_NoChange(t *testing.T) {
	tmp := t.TempDir()
	ws := state.WorkflowState{ChangeName: "", Phase: state.PhaseNone}
	var buf bytes.Buffer
	err := runVerifyContextOnly(&buf, tmp, ws, config.Defaults())
	assert.Error(t, err)
}

// TestVerifyWriteResults_MustPass tests that write-results transitions state to verified when must_pass is true.
func TestVerifyWriteResults_MustPass(t *testing.T) {
	specsDir, changeName, changeDir := setupVerifyTestChange(t)

	// Create a verifier report with must_pass=true
	report := verifier.VerifierReport{
		ChangeName:  changeName,
		OverallPass: true,
		MustPass:    true,
		Results: []verifier.VerifierResultItem{
			{
				ID:      "spec.md::must-abc123",
				Text:    "The system MUST provide authentication.",
				Keyword: "MUST",
				Pass:    true,
				Evidence: "Auth module found in code",
			},
		},
	}
	reportData, err := json.Marshal(report)
	require.NoError(t, err)

	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, os.WriteFile(reportPath, reportData, 0644))

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseExecuted}
	var buf bytes.Buffer
	err = runVerifyWriteResults(&buf, specsDir, &ws, reportPath)
	require.NoError(t, err)

	// State should now be verified
	assert.Equal(t, state.PhaseVerified, ws.Phase)
	require.NotNil(t, ws.VerifyPass)
	assert.True(t, *ws.VerifyPass)

	// verification.md should exist
	assert.FileExists(t, filepath.Join(changeDir, "verification.md"))

	// verification-status.json should exist
	assert.FileExists(t, filepath.Join(changeDir, "verification-status.json"))

	// gap-report.md should NOT exist (all passed)
	_, statErr := os.Stat(filepath.Join(changeDir, "gap-report.md"))
	assert.True(t, os.IsNotExist(statErr))
}

// TestVerifyWriteResults_MustFail tests that state stays executed when must_pass is false.
func TestVerifyWriteResults_MustFail(t *testing.T) {
	specsDir, changeName, changeDir := setupVerifyTestChange(t)

	// Create a verifier report with must_pass=false
	report := verifier.VerifierReport{
		ChangeName:  changeName,
		OverallPass: false,
		MustPass:    false,
		Results: []verifier.VerifierResultItem{
			{
				ID:         "spec.md::must-abc123",
				Text:       "The system MUST provide authentication.",
				Keyword:    "MUST",
				Pass:       false,
				Evidence:   "No auth module found",
				Suggestion: "Implement auth middleware",
			},
		},
	}
	reportData, err := json.Marshal(report)
	require.NoError(t, err)

	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, os.WriteFile(reportPath, reportData, 0644))

	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseExecuted}
	var buf bytes.Buffer
	err = runVerifyWriteResults(&buf, specsDir, &ws, reportPath)
	require.NoError(t, err)

	// State should still be executed (not verified)
	assert.Equal(t, state.PhaseExecuted, ws.Phase)
	require.NotNil(t, ws.VerifyPass)
	assert.False(t, *ws.VerifyPass)

	// verification.md should still be written
	assert.FileExists(t, filepath.Join(changeDir, "verification.md"))

	// gap-report.md should be written (there are failures)
	assert.FileExists(t, filepath.Join(changeDir, "gap-report.md"))

	// verification-status.json should exist
	assert.FileExists(t, filepath.Join(changeDir, "verification-status.json"))

	// Check status file has "blocked" for the MUST item
	vs, readErr := spec.ReadVerificationStatus(changeDir)
	require.NoError(t, readErr)
	assert.Equal(t, spec.StatusBlocked, vs.Requirements["spec.md::must-abc123"])
}

// TestVerifyNoFlags tests that verify with no flags returns usage hint error.
func TestVerifyNoFlags(t *testing.T) {
	cmd := &cobra.Command{}
	err := runVerifyNoFlags(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage:")
}
