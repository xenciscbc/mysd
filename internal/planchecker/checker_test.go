package planchecker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xenciscbc/mysd/internal/planchecker"
	"github.com/xenciscbc/mysd/internal/spec"
)

func TestCheckCoverage_AllCovered(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Satisfies: []string{"REQ-01"}},
		{ID: 2, Name: "Task B", Satisfies: []string{"REQ-02"}},
		{ID: 3, Name: "Task C", Satisfies: []string{"REQ-03"}},
	}
	mustIDs := []string{"REQ-01", "REQ-02", "REQ-03"}

	result := planchecker.CheckCoverage(tasks, mustIDs)

	assert.True(t, result.Passed)
	assert.Equal(t, 3, result.TotalMust)
	assert.Equal(t, 3, result.CoveredCount)
	assert.Empty(t, result.UncoveredIDs)
	assert.InDelta(t, 1.0, result.CoverageRatio, 0.001)
}

func TestCheckCoverage_PartialCoverage(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Satisfies: []string{"REQ-01"}},
		{ID: 2, Name: "Task B", Satisfies: []string{"REQ-02"}},
	}
	mustIDs := []string{"REQ-01", "REQ-02", "REQ-03"}

	result := planchecker.CheckCoverage(tasks, mustIDs)

	assert.False(t, result.Passed)
	assert.Equal(t, 3, result.TotalMust)
	assert.Equal(t, 2, result.CoveredCount)
	assert.Equal(t, []string{"REQ-03"}, result.UncoveredIDs)
	assert.InDelta(t, 2.0/3.0, result.CoverageRatio, 0.001)
}

func TestCheckCoverage_NoneCovered(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Satisfies: []string{"UNRELATED-01"}},
	}
	mustIDs := []string{"REQ-01", "REQ-02", "REQ-03"}

	result := planchecker.CheckCoverage(tasks, mustIDs)

	assert.False(t, result.Passed)
	assert.Equal(t, 3, result.TotalMust)
	assert.Equal(t, 0, result.CoveredCount)
	assert.Len(t, result.UncoveredIDs, 3)
	assert.Contains(t, result.UncoveredIDs, "REQ-01")
	assert.Contains(t, result.UncoveredIDs, "REQ-02")
	assert.Contains(t, result.UncoveredIDs, "REQ-03")
	assert.InDelta(t, 0.0, result.CoverageRatio, 0.001)
}

func TestCheckCoverage_EmptyMustIDs(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Satisfies: []string{"REQ-01"}},
	}
	mustIDs := []string{}

	result := planchecker.CheckCoverage(tasks, mustIDs)

	assert.True(t, result.Passed)
	assert.Equal(t, 0, result.TotalMust)
	assert.Equal(t, 0, result.CoveredCount)
}

func TestCheckCoverage_EmptyTasks(t *testing.T) {
	tasks := []spec.TaskEntry{}
	mustIDs := []string{"REQ-01", "REQ-02", "REQ-03"}

	result := planchecker.CheckCoverage(tasks, mustIDs)

	assert.False(t, result.Passed)
	assert.Equal(t, 3, result.TotalMust)
	assert.Equal(t, 0, result.CoveredCount)
	assert.Len(t, result.UncoveredIDs, 3)
	assert.InDelta(t, 0.0, result.CoverageRatio, 0.001)
}

func TestCheckCoverage_DuplicateCoverage(t *testing.T) {
	// Two tasks both satisfy the same ID — should count as 1, not 2
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Satisfies: []string{"REQ-01"}},
		{ID: 2, Name: "Task B", Satisfies: []string{"REQ-01"}},
	}
	mustIDs := []string{"REQ-01"}

	result := planchecker.CheckCoverage(tasks, mustIDs)

	assert.True(t, result.Passed)
	assert.Equal(t, 1, result.TotalMust)
	assert.Equal(t, 1, result.CoveredCount)
	assert.Empty(t, result.UncoveredIDs)
	assert.InDelta(t, 1.0, result.CoverageRatio, 0.001)
}

func TestCheckCoverage_MultiSatisfies(t *testing.T) {
	// One task satisfies multiple MUST IDs
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Satisfies: []string{"REQ-01", "REQ-02"}},
	}
	mustIDs := []string{"REQ-01", "REQ-02"}

	result := planchecker.CheckCoverage(tasks, mustIDs)

	assert.True(t, result.Passed)
	assert.Equal(t, 2, result.TotalMust)
	assert.Equal(t, 2, result.CoveredCount)
	assert.Empty(t, result.UncoveredIDs)
	assert.InDelta(t, 1.0, result.CoverageRatio, 0.001)
}
