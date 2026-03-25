package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildWaveGroups_Empty verifies that empty input returns nil.
func TestBuildWaveGroups_Empty(t *testing.T) {
	groups, err := BuildWaveGroups(nil)
	require.NoError(t, err)
	assert.Nil(t, groups)
}

// TestBuildWaveGroups_NoDeps verifies that tasks with no dependencies land in one wave.
func TestBuildWaveGroups_NoDeps(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A"},
		{ID: 2, Name: "Task B"},
		{ID: 3, Name: "Task C"},
	}
	groups, err := BuildWaveGroups(tasks)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	assert.Len(t, groups[0], 3)
}

// TestBuildWaveGroups_LinearChain verifies A->B->C produces 3 waves of 1 task each.
func TestBuildWaveGroups_LinearChain(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A"},
		{ID: 2, Name: "Task B", Depends: []int{1}},
		{ID: 3, Name: "Task C", Depends: []int{2}},
	}
	groups, err := BuildWaveGroups(tasks)
	require.NoError(t, err)
	require.Len(t, groups, 3)
	assert.Len(t, groups[0], 1)
	assert.Equal(t, 1, groups[0][0].ID)
	assert.Len(t, groups[1], 1)
	assert.Equal(t, 2, groups[1][0].ID)
	assert.Len(t, groups[2], 1)
	assert.Equal(t, 3, groups[2][0].ID)
}

// TestBuildWaveGroups_Diamond verifies diamond pattern: A,B (no deps) -> C (depends A,B)
// produces 2 waves: [A,B], [C].
func TestBuildWaveGroups_Diamond(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A"},
		{ID: 2, Name: "Task B"},
		{ID: 3, Name: "Task C", Depends: []int{1, 2}},
	}
	groups, err := BuildWaveGroups(tasks)
	require.NoError(t, err)
	require.Len(t, groups, 2)
	assert.Len(t, groups[0], 2)
	assert.Equal(t, 1, groups[0][0].ID)
	assert.Equal(t, 2, groups[0][1].ID)
	assert.Len(t, groups[1], 1)
	assert.Equal(t, 3, groups[1][0].ID)
}

// TestBuildWaveGroups_Cycle verifies that A depends B, B depends A returns ErrCyclicDependency.
func TestBuildWaveGroups_Cycle(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A", Depends: []int{2}},
		{ID: 2, Name: "Task B", Depends: []int{1}},
	}
	groups, err := BuildWaveGroups(tasks)
	assert.Nil(t, groups)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrCyclicDependency))
}

// TestBuildWaveGroups_FileOverlap verifies same-layer tasks with overlapping Files
// get split into separate waves.
func TestBuildWaveGroups_FileOverlap(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A", Files: []string{"shared.go", "a.go"}},
		{ID: 2, Name: "Task B", Files: []string{"shared.go", "b.go"}},
	}
	groups, err := BuildWaveGroups(tasks)
	require.NoError(t, err)
	// Both tasks are in the same dependency layer, but share "shared.go"
	// so they must be split into 2 waves.
	require.Len(t, groups, 2)
	assert.Len(t, groups[0], 1)
	assert.Len(t, groups[1], 1)
}

// TestBuildWaveGroups_NoOverlap verifies same-layer tasks with distinct Files stay in one wave.
func TestBuildWaveGroups_NoOverlap(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A", Files: []string{"a.go"}},
		{ID: 2, Name: "Task B", Files: []string{"b.go"}},
	}
	groups, err := BuildWaveGroups(tasks)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	assert.Len(t, groups[0], 2)
}

// TestBuildWaveGroups_DeterministicOrder verifies tasks within a wave are sorted by ID ascending.
func TestBuildWaveGroups_DeterministicOrder(t *testing.T) {
	tasks := []TaskItem{
		{ID: 5, Name: "Task E"},
		{ID: 3, Name: "Task C"},
		{ID: 1, Name: "Task A"},
		{ID: 4, Name: "Task D"},
		{ID: 2, Name: "Task B"},
	}
	groups, err := BuildWaveGroups(tasks)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Len(t, groups[0], 5)
	assert.Equal(t, 1, groups[0][0].ID)
	assert.Equal(t, 2, groups[0][1].ID)
	assert.Equal(t, 3, groups[0][2].ID)
	assert.Equal(t, 4, groups[0][3].ID)
	assert.Equal(t, 5, groups[0][4].ID)
}

// TestHasParallelOpportunity_NoDepsNoFiles verifies returns false when no deps or files.
func TestHasParallelOpportunity_NoDepsNoFiles(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A"},
		{ID: 2, Name: "Task B"},
	}
	assert.False(t, HasParallelOpportunity(tasks))
}

// TestHasParallelOpportunity_HasDepends verifies returns true when task has Depends.
func TestHasParallelOpportunity_HasDepends(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A"},
		{ID: 2, Name: "Task B", Depends: []int{1}},
	}
	assert.True(t, HasParallelOpportunity(tasks))
}

// TestHasParallelOpportunity_HasFiles verifies returns true when task has Files.
func TestHasParallelOpportunity_HasFiles(t *testing.T) {
	tasks := []TaskItem{
		{ID: 1, Name: "Task A", Files: []string{"main.go"}},
	}
	assert.True(t, HasParallelOpportunity(tasks))
}
