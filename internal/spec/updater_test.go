package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tasksV2Content is a tasks.md with YAML frontmatter containing a tasks list.
const tasksV2Content = `---
spec-version: "1"
total: 3
completed: 1
tasks:
  - id: 1
    name: "Design API schema"
    description: "Define OpenAPI spec"
    status: done
  - id: 2
    name: "Implement handler"
    status: pending
  - id: 3
    name: "Write integration tests"
    status: pending
---

## Tasks

Progress: 1/3 tasks complete.
`

// tasksV2EmptyContent is a tasks.md with zero tasks.
const tasksV2EmptyContent = `---
spec-version: "1"
total: 0
completed: 0
---

## Tasks

No tasks yet.
`

// writeTempTasksFile creates a temporary tasks.md file and returns its path.
func writeTempTasksFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

// TestParseTasksV2_WithTasksList verifies ParseTasks reads tasks list and body string correctly.
func TestParseTasksV2_WithTasksList(t *testing.T) {
	path := writeTempTasksFile(t, tasksV2Content)

	fm, body, err := ParseTasksV2(path)
	require.NoError(t, err)

	assert.Equal(t, "1", fm.SpecVersion)
	assert.Equal(t, 3, fm.Total)
	assert.Equal(t, 1, fm.Completed)
	assert.Len(t, fm.Tasks, 3)

	assert.Equal(t, 1, fm.Tasks[0].ID)
	assert.Equal(t, "Design API schema", fm.Tasks[0].Name)
	assert.Equal(t, StatusDone, fm.Tasks[0].Status)

	assert.Equal(t, 2, fm.Tasks[1].ID)
	assert.Equal(t, "Implement handler", fm.Tasks[1].Name)
	assert.Equal(t, StatusPending, fm.Tasks[1].Status)

	assert.Contains(t, body, "## Tasks")
}

// TestUpdateTaskStatus_PendingToInProgress verifies status update and Completed recomputation.
func TestUpdateTaskStatus_PendingToInProgress(t *testing.T) {
	path := writeTempTasksFile(t, tasksV2Content)

	err := UpdateTaskStatus(path, 2, StatusInProgress)
	require.NoError(t, err)

	fm, _, err := ParseTasksV2(path)
	require.NoError(t, err)

	assert.Equal(t, StatusInProgress, fm.Tasks[1].Status)
	// Completed count: only task 1 is done, task 2 is in_progress
	assert.Equal(t, 1, fm.Completed)
}

// TestUpdateTaskStatus_MarkDone verifies Completed increments when a task is marked done.
func TestUpdateTaskStatus_MarkDone(t *testing.T) {
	path := writeTempTasksFile(t, tasksV2Content)

	err := UpdateTaskStatus(path, 2, StatusDone)
	require.NoError(t, err)

	fm, _, err := ParseTasksV2(path)
	require.NoError(t, err)

	assert.Equal(t, StatusDone, fm.Tasks[1].Status)
	// Now tasks 1 and 2 are done → Completed = 2
	assert.Equal(t, 2, fm.Completed)
}

// TestUpdateTaskStatus_NotFound verifies error message for nonexistent task ID.
func TestUpdateTaskStatus_NotFound(t *testing.T) {
	path := writeTempTasksFile(t, tasksV2Content)

	err := UpdateTaskStatus(path, 99, StatusDone)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task 99 not found")
}

// TestWriteTasks_PreservesBody verifies round-trip preserves markdown body.
func TestWriteTasks_PreservesBody(t *testing.T) {
	path := writeTempTasksFile(t, tasksV2Content)

	fm, body, err := ParseTasksV2(path)
	require.NoError(t, err)

	// Modify a task and write back
	fm.Tasks[0].Status = StatusDone
	err = WriteTasks(path, fm, body)
	require.NoError(t, err)

	// Re-read and verify body is preserved
	_, newBody, err := ParseTasksV2(path)
	require.NoError(t, err)
	assert.Equal(t, body, newBody, "body content must be preserved after round-trip")
}

// TestParseTasksV2_EmptyTasks verifies zero tasks returns empty slice without error.
func TestParseTasksV2_EmptyTasks(t *testing.T) {
	path := writeTempTasksFile(t, tasksV2EmptyContent)

	fm, body, err := ParseTasksV2(path)
	require.NoError(t, err)

	assert.Empty(t, fm.Tasks, "empty tasks list should return empty slice")
	assert.Equal(t, 0, fm.Total)
	assert.Contains(t, body, "No tasks yet.")
}
