package executor

import (
	"testing"

	"github.com/mysd/internal/spec"
	"github.com/stretchr/testify/assert"
)

// TestPendingTasks_FiltersDoneAndBlocked verifies done and blocked tasks are excluded.
func TestPendingTasks_FiltersDoneAndBlocked(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Status: spec.StatusPending},
		{ID: 2, Name: "Task B", Status: spec.StatusDone},
		{ID: 3, Name: "Task C", Status: spec.StatusBlocked},
		{ID: 4, Name: "Task D", Status: spec.StatusInProgress},
	}

	pending := PendingTasks(tasks)

	assert.Len(t, pending, 2)
	assert.Equal(t, 1, pending[0].ID)
	assert.Equal(t, 4, pending[1].ID)
}

// TestCalcProgress_CorrectCounts verifies done/total calculation.
func TestCalcProgress_CorrectCounts(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Status: spec.StatusDone},
		{ID: 2, Status: spec.StatusPending},
		{ID: 3, Status: spec.StatusDone},
		{ID: 4, Status: spec.StatusBlocked},
	}

	done, total := CalcProgress(tasks)

	assert.Equal(t, 2, done)
	assert.Equal(t, 4, total)
}
