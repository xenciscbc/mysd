package executor

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: RenderStatus with change "my-feature" in PhasePlanned with 2/5 done
// outputs string containing "my-feature" and "2/5"
func TestRenderStatus_OutputsChangeNameAndProgress(t *testing.T) {
	summary := StatusSummary{
		ChangeName: "my-feature",
		Phase:      string(state.PhasePlanned),
		TasksDone:  2,
		TasksTotal: 5,
	}
	var buf bytes.Buffer
	RenderStatus(&buf, summary)
	output := buf.String()
	assert.Contains(t, output, "my-feature", "output should contain change name")
	assert.Contains(t, output, "2", "output should contain tasks done count")
	assert.Contains(t, output, "5", "output should contain tasks total count")
}

// Test 2: RenderStatus with 3 MUST items (2 done, 1 pending) outputs string containing "2" and "3"
func TestRenderStatus_OutputsMUSTCounts(t *testing.T) {
	summary := StatusSummary{
		ChangeName: "test-change",
		Phase:      string(state.PhaseExecuted),
		TasksDone:  0,
		TasksTotal: 0,
		MustDone:   2,
		MustTotal:  3,
	}
	var buf bytes.Buffer
	RenderStatus(&buf, summary)
	output := buf.String()
	assert.Contains(t, output, "2", "output should contain MUST done count")
	assert.Contains(t, output, "3", "output should contain MUST total count")
}

// Test 3: RenderStatus with zero tasks outputs "0/0" or similar zero progress indication
func TestRenderStatus_ZeroTasks(t *testing.T) {
	summary := StatusSummary{
		ChangeName: "empty-change",
		Phase:      string(state.PhaseProposed),
		TasksDone:  0,
		TasksTotal: 0,
	}
	var buf bytes.Buffer
	RenderStatus(&buf, summary)
	output := buf.String()
	assert.NotEmpty(t, output, "output should not be empty even with zero tasks")
	// Should show 0 done and 0 total
	assert.Contains(t, output, "0", "output should contain zero value")
}

// Test 4: RenderStatus with LastRun set outputs formatted date string
func TestRenderStatus_OutputsLastRun(t *testing.T) {
	summary := StatusSummary{
		ChangeName: "dated-change",
		Phase:      string(state.PhasePlanned),
		TasksDone:  1,
		TasksTotal: 2,
		LastRun:    "2026-03-24 12:00",
	}
	var buf bytes.Buffer
	RenderStatus(&buf, summary)
	output := buf.String()
	assert.Contains(t, output, "2026-03-24", "output should contain formatted last run date")
}

// Test 5: BuildStatusSummary correctly computes StatusSummary from WorkflowState + tasks + requirements
func TestBuildStatusSummary_ComputesCorrectly(t *testing.T) {
	lastRun := time.Date(2026, 3, 24, 10, 30, 0, 0, time.UTC)
	ws := state.WorkflowState{
		ChangeName: "my-change",
		Phase:      state.PhasePlanned,
		LastRun:    lastRun,
	}

	tasks := []spec.Task{
		{ID: 1, Status: spec.StatusDone},
		{ID: 2, Status: spec.StatusDone},
		{ID: 3, Status: spec.StatusPending},
		{ID: 4, Status: spec.StatusInProgress},
		{ID: 5, Status: spec.StatusDone},
	}

	reqs := []spec.Requirement{
		{ID: "R1", Keyword: spec.Must, Status: spec.StatusDone},
		{ID: "R2", Keyword: spec.Must, Status: spec.StatusDone},
		{ID: "R3", Keyword: spec.Must, Status: spec.StatusPending},
		{ID: "R4", Keyword: spec.Should, Status: spec.StatusDone},
		{ID: "R5", Keyword: spec.Should, Status: spec.StatusPending},
		{ID: "R6", Keyword: spec.May, Status: spec.StatusPending},
	}

	summary := BuildStatusSummary(ws, tasks, reqs)

	require.Equal(t, "my-change", summary.ChangeName)
	require.Equal(t, string(state.PhasePlanned), summary.Phase)
	require.Equal(t, 3, summary.TasksDone, "should count 3 done tasks")
	require.Equal(t, 5, summary.TasksTotal, "should count 5 total tasks")
	require.Equal(t, 2, summary.MustDone, "should count 2 done MUST requirements")
	require.Equal(t, 3, summary.MustTotal, "should count 3 total MUST requirements")
	require.Equal(t, 1, summary.ShouldDone, "should count 1 done SHOULD requirements")
	require.Equal(t, 2, summary.ShouldTotal, "should count 2 total SHOULD requirements")
	require.Equal(t, 1, summary.MayTotal, "should count 1 total MAY requirements")
	assert.Equal(t, "2026-03-24 10:30", summary.LastRun, "LastRun should be formatted as '2006-01-02 15:04'")
}

// Test 6: BuildStatusSummary with zero LastRun returns "never"
func TestBuildStatusSummary_ZeroLastRun_ReturnsNever(t *testing.T) {
	ws := state.WorkflowState{
		ChangeName: "fresh-change",
		Phase:      state.PhaseProposed,
	}
	summary := BuildStatusSummary(ws, nil, nil)
	assert.Equal(t, "never", summary.LastRun, "zero LastRun should display as 'never'")
}

// Helper to check output contains a string (case-insensitive for label matching)
func assertContainsCI(t *testing.T, output, substr string) {
	t.Helper()
	assert.True(t, strings.Contains(strings.ToLower(output), strings.ToLower(substr)),
		"expected output to contain %q (case-insensitive), got:\n%s", substr, output)
}
