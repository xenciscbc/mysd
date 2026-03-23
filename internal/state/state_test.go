package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadState_MissingFile_ReturnsZeroState(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := LoadState(tmpDir)
	require.NoError(t, err, "missing STATE.json should not return error (convention over config)")
	assert.Equal(t, WorkflowState{}, ws, "missing file should return zero-value WorkflowState")
}

func TestSaveState_LoadState_Roundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	verifyPass := true
	original := WorkflowState{
		ChangeName: "add-auth",
		Phase:      PhaseProposed,
		LastRun:    time.Date(2026, 3, 23, 10, 0, 0, 0, time.UTC),
		VerifyPass: &verifyPass,
	}

	err := SaveState(tmpDir, original)
	require.NoError(t, err, "SaveState should not error")

	loaded, err := LoadState(tmpDir)
	require.NoError(t, err, "LoadState should not error")
	assert.Equal(t, original.ChangeName, loaded.ChangeName)
	assert.Equal(t, original.Phase, loaded.Phase)
	assert.Equal(t, original.LastRun.UTC(), loaded.LastRun.UTC())
	require.NotNil(t, loaded.VerifyPass)
	assert.Equal(t, *original.VerifyPass, *loaded.VerifyPass)
}

func TestSaveState_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "nested", ".specs")
	ws := WorkflowState{ChangeName: "test", Phase: PhaseNone}
	err := SaveState(nestedDir, ws)
	require.NoError(t, err, "SaveState should create nested directories")
	_, err = os.Stat(filepath.Join(nestedDir, "STATE.json"))
	require.NoError(t, err, "STATE.json should exist after SaveState")
}

func TestWorkflowState_JSONMarshaling(t *testing.T) {
	ws := WorkflowState{
		ChangeName: "my-change",
		Phase:      PhaseExecuted,
		LastRun:    time.Date(2026, 1, 15, 9, 30, 0, 0, time.UTC),
	}
	data, err := json.Marshal(ws)
	require.NoError(t, err)

	var decoded WorkflowState
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, ws.ChangeName, decoded.ChangeName)
	assert.Equal(t, ws.Phase, decoded.Phase)
	assert.Nil(t, decoded.VerifyPass, "omitempty nil pointer should not appear in JSON")
}

func TestWorkflowState_JSONFields(t *testing.T) {
	ws := WorkflowState{
		ChangeName: "test-change",
		Phase:      PhaseVerified,
		LastRun:    time.Now().UTC(),
	}
	data, err := json.Marshal(ws)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	assert.Contains(t, raw, "change_name", "should use JSON tag 'change_name'")
	assert.Contains(t, raw, "phase", "should use JSON tag 'phase'")
	assert.Contains(t, raw, "last_run", "should use JSON tag 'last_run'")
	assert.NotContains(t, raw, "verify_pass", "omitempty nil should not appear")
}
