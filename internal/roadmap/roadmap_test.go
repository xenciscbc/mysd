package roadmap_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mysd/internal/roadmap"
	"github.com/mysd/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpdateTracking_NewFile verifies that tracking.yaml is created when it doesn't exist.
func TestUpdateTracking_NewFile(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, ".specs")
	require.NoError(t, os.MkdirAll(specsDir, 0755))

	ws := state.WorkflowState{
		ChangeName: "auth",
		Phase:      state.PhaseProposed,
		LastRun:    time.Now(),
	}

	err := roadmap.UpdateTracking(specsDir, ws)
	require.NoError(t, err)

	trackingPath := filepath.Join(tmpDir, ".mysd", "roadmap", "tracking.yaml")
	assert.FileExists(t, trackingPath)

	tf, err := roadmap.ReadTracking(filepath.Join(tmpDir, ".mysd", "roadmap"))
	require.NoError(t, err)
	assert.Equal(t, "1", tf.SchemaVersion)
	require.Len(t, tf.Changes, 1)
	assert.Equal(t, "auth", tf.Changes[0].Name)
	assert.Equal(t, "proposed", tf.Changes[0].Status)
	assert.NotNil(t, tf.Changes[0].StartedAt)
}

// TestUpdateTracking_UpsertExisting verifies that an existing record is updated (not duplicated).
func TestUpdateTracking_UpsertExisting(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, ".specs")
	require.NoError(t, os.MkdirAll(specsDir, 0755))

	ws := state.WorkflowState{
		ChangeName: "auth",
		Phase:      state.PhaseProposed,
		LastRun:    time.Now(),
	}
	require.NoError(t, roadmap.UpdateTracking(specsDir, ws))

	// Now update to specced
	ws.Phase = state.PhaseSpecced
	err := roadmap.UpdateTracking(specsDir, ws)
	require.NoError(t, err)

	tf, err := roadmap.ReadTracking(filepath.Join(tmpDir, ".mysd", "roadmap"))
	require.NoError(t, err)
	// Should not be duplicated
	assert.Len(t, tf.Changes, 1)
	assert.Equal(t, "specced", tf.Changes[0].Status)
}

// TestUpdateTracking_MultipleChanges verifies multiple distinct changes are tracked.
func TestUpdateTracking_MultipleChanges(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, ".specs")
	require.NoError(t, os.MkdirAll(specsDir, 0755))

	ws1 := state.WorkflowState{
		ChangeName: "auth",
		Phase:      state.PhaseProposed,
		LastRun:    time.Now(),
	}
	require.NoError(t, roadmap.UpdateTracking(specsDir, ws1))

	ws2 := state.WorkflowState{
		ChangeName: "payments",
		Phase:      state.PhaseProposed,
		LastRun:    time.Now(),
	}
	err := roadmap.UpdateTracking(specsDir, ws2)
	require.NoError(t, err)

	tf, err := roadmap.ReadTracking(filepath.Join(tmpDir, ".mysd", "roadmap"))
	require.NoError(t, err)
	assert.Len(t, tf.Changes, 2)
}

// TestUpdateTracking_CompletedAt verifies that CompletedAt is set when phase is archived.
func TestUpdateTracking_CompletedAt(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, ".specs")
	require.NoError(t, os.MkdirAll(specsDir, 0755))

	ws := state.WorkflowState{
		ChangeName: "auth",
		Phase:      state.PhaseArchived,
		LastRun:    time.Now(),
	}

	err := roadmap.UpdateTracking(specsDir, ws)
	require.NoError(t, err)

	tf, err := roadmap.ReadTracking(filepath.Join(tmpDir, ".mysd", "roadmap"))
	require.NoError(t, err)
	require.Len(t, tf.Changes, 1)
	assert.NotNil(t, tf.Changes[0].CompletedAt, "CompletedAt should be set when phase is archived")
}

// TestUpdateTracking_TimelineMdGenerated verifies that timeline.md is created alongside tracking.yaml.
func TestUpdateTracking_TimelineMdGenerated(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, ".specs")
	require.NoError(t, os.MkdirAll(specsDir, 0755))

	ws := state.WorkflowState{
		ChangeName: "auth",
		Phase:      state.PhaseProposed,
		LastRun:    time.Now(),
	}

	err := roadmap.UpdateTracking(specsDir, ws)
	require.NoError(t, err)

	timelinePath := filepath.Join(tmpDir, ".mysd", "roadmap", "timeline.md")
	assert.FileExists(t, timelinePath)

	content, err := os.ReadFile(timelinePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "gantt")
}

// TestUpdateTracking_ProjectRootDerivation verifies tracking.yaml lands in {root}/.mysd/roadmap/,
// NOT inside .specs/.
func TestUpdateTracking_ProjectRootDerivation(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, ".specs")
	require.NoError(t, os.MkdirAll(specsDir, 0755))

	ws := state.WorkflowState{
		ChangeName: "auth",
		Phase:      state.PhaseProposed,
		LastRun:    time.Now(),
	}

	err := roadmap.UpdateTracking(specsDir, ws)
	require.NoError(t, err)

	// tracking.yaml must be at {root}/.mysd/roadmap/tracking.yaml
	expectedPath := filepath.Join(tmpDir, ".mysd", "roadmap", "tracking.yaml")
	assert.FileExists(t, expectedPath)

	// tracking.yaml must NOT be inside .specs/
	wrongPath := filepath.Join(specsDir, ".mysd", "roadmap", "tracking.yaml")
	_, err = os.Stat(wrongPath)
	assert.True(t, os.IsNotExist(err), "tracking.yaml should NOT be inside .specs/")
}
