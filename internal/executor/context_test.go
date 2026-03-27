package executor

import (
	"encoding/json"
	"testing"

	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildContext_PendingTasksFiltered verifies that BuildContext returns
// only pending tasks (done tasks excluded per EXEC-02, EXEC-05).
func TestBuildContext_PendingTasksFiltered(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Status: spec.StatusDone},
		{ID: 2, Name: "Task B", Status: spec.StatusPending},
		{ID: 3, Name: "Task C", Status: spec.StatusPending},
	}
	reqs := []spec.Requirement{}
	cfg := config.Defaults()

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)

	assert.Len(t, ctx.PendingTasks, 2)
	assert.Equal(t, 1, ctx.Tasks[0].ID)  // all tasks included
	assert.Len(t, ctx.Tasks, 3)
}

// TestBuildContext_MustItemsFiltered verifies MUST requirements are extracted.
func TestBuildContext_MustItemsFiltered(t *testing.T) {
	reqs := []spec.Requirement{
		{ID: "REQ-01", Text: "System MUST validate input", Keyword: spec.Must, Status: spec.StatusPending},
		{ID: "REQ-02", Text: "System SHOULD log errors", Keyword: spec.Should, Status: spec.StatusPending},
		{ID: "REQ-03", Text: "System MAY cache results", Keyword: spec.May, Status: spec.StatusPending},
	}
	tasks := []spec.TaskEntry{}
	cfg := config.Defaults()

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)

	assert.Len(t, ctx.MustItems, 1)
	assert.Equal(t, "REQ-01", ctx.MustItems[0].ID)
	assert.Len(t, ctx.ShouldItems, 1)
	assert.Equal(t, "REQ-02", ctx.ShouldItems[0].ID)
	assert.Len(t, ctx.MayItems, 1)
	assert.Equal(t, "REQ-03", ctx.MayItems[0].ID)
}

// TestBuildContext_ConfigFieldsPopulated verifies TDD and ExecutionMode come from config.
func TestBuildContext_ConfigFieldsPopulated(t *testing.T) {
	cfg := config.ProjectConfig{
		ExecutionMode: "wave",
		AgentCount:    3,
		TDD:           true,
		AtomicCommits: true,
	}
	tasks := []spec.TaskEntry{}
	reqs := []spec.Requirement{}

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)

	assert.True(t, ctx.TDDMode)
	assert.Equal(t, "wave", ctx.ExecutionMode)
	assert.Equal(t, 3, ctx.AgentCount)
	assert.True(t, ctx.AtomicCommits)
}

// TestBuildContext_JSONMarshal verifies ExecutionContext marshals to valid JSON
// with required keys per EXEC-01.
func TestBuildContext_JSONMarshal(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Status: spec.StatusPending},
	}
	reqs := []spec.Requirement{
		{ID: "REQ-01", Text: "System MUST validate", Keyword: spec.Must},
	}
	cfg := config.Defaults()
	cfg.TDD = true

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)
	data, err := json.Marshal(ctx)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &m))

	assert.Contains(t, m, "change_name")
	assert.Contains(t, m, "must_items")
	assert.Contains(t, m, "tasks")
	assert.Contains(t, m, "pending_tasks")
	assert.Contains(t, m, "tdd_mode")
	assert.Contains(t, m, "execution_mode")
	assert.Contains(t, m, "agent_count")
	assert.Contains(t, m, "has_parallel_opportunity")
}

// TestBuildContextFromParts_NewFields verifies new fields are copied from TaskEntry to TaskItem.
func TestBuildContextFromParts_NewFields(t *testing.T) {
	tasks := []spec.TaskEntry{
		{
			ID:        1,
			Name:      "Build auth",
			Status:    spec.StatusPending,
			Depends:   []int{2, 3},
			Files:     []string{"auth.go"},
			Satisfies: []string{"REQ-01"},
			Skills:    []string{"/mysd:apply"},
		},
	}
	cfg := config.Defaults()

	ctx := BuildContextFromParts("my-change", tasks, nil, cfg)

	require.Len(t, ctx.Tasks, 1)
	ti := ctx.Tasks[0]
	assert.Equal(t, []int{2, 3}, ti.Depends)
	assert.Equal(t, []string{"auth.go"}, ti.Files)
	assert.Equal(t, []string{"REQ-01"}, ti.Satisfies)
	assert.Equal(t, []string{"/mysd:apply"}, ti.Skills)

	// Also check PendingTasks
	require.Len(t, ctx.PendingTasks, 1)
	pt := ctx.PendingTasks[0]
	assert.Equal(t, []int{2, 3}, pt.Depends)
	assert.Equal(t, []string{"REQ-01"}, pt.Satisfies)
}

// TestBuildContextFromParts_WaveGroups verifies WaveGroups is populated from pending tasks with depends.
func TestBuildContextFromParts_WaveGroups(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Status: spec.StatusPending},
		{ID: 2, Name: "Task B", Status: spec.StatusPending, Depends: []int{1}},
	}
	cfg := config.Defaults()

	ctx := BuildContextFromParts("my-change", tasks, nil, cfg)

	// WaveGroups should be computed: wave 1 = [Task A], wave 2 = [Task B]
	require.Len(t, ctx.WaveGroups, 2)
	assert.Equal(t, 1, ctx.WaveGroups[0][0].ID)
	assert.Equal(t, 2, ctx.WaveGroups[1][0].ID)
	assert.True(t, ctx.HasParallelOpp)
	assert.Equal(t, ".worktrees", ctx.WorktreeDir)
	assert.False(t, ctx.AutoMode)
}

// TestBuildContextFromParts_WaveGroups_NoParallel verifies WaveGroups and HasParallelOpp
// when tasks have no deps or files.
func TestBuildContextFromParts_WaveGroups_NoParallel(t *testing.T) {
	tasks := []spec.TaskEntry{
		{ID: 1, Name: "Task A", Status: spec.StatusPending},
		{ID: 2, Name: "Task B", Status: spec.StatusPending},
	}
	cfg := config.Defaults()

	ctx := BuildContextFromParts("my-change", tasks, nil, cfg)

	// No depends/files — still computes wave groups (single wave)
	require.Len(t, ctx.WaveGroups, 1)
	assert.Len(t, ctx.WaveGroups[0], 2)
	assert.False(t, ctx.HasParallelOpp)
}

// TestBuildContextFromParts_AutoMode verifies WorktreeDir and AutoMode come from config.
func TestBuildContextFromParts_AutoMode(t *testing.T) {
	cfg := config.ProjectConfig{
		ExecutionMode: "wave",
		AgentCount:    2,
		WorktreeDir:   ".worktrees",
		AutoMode:      true,
	}
	tasks := []spec.TaskEntry{}
	reqs := []spec.Requirement{}

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)

	assert.Equal(t, ".worktrees", ctx.WorktreeDir)
	assert.True(t, ctx.AutoMode)
}

// TestBuildContextFromParts_DocsToUpdate verifies DocsToUpdate is passed from config to ExecutionContext.
func TestBuildContextFromParts_DocsToUpdate(t *testing.T) {
	cfg := config.ProjectConfig{
		ExecutionMode: "single",
		AgentCount:    1,
		DocsToUpdate:  []string{"README.md", "CHANGELOG.md"},
	}
	tasks := []spec.TaskEntry{}
	reqs := []spec.Requirement{}

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)

	assert.Equal(t, []string{"README.md", "CHANGELOG.md"}, ctx.DocsToUpdate)
}

// TestBuildContextFromParts_DocsToUpdateNil verifies DocsToUpdate nil (default) produces nil in ExecutionContext.
func TestBuildContextFromParts_DocsToUpdateNil(t *testing.T) {
	cfg := config.Defaults()
	tasks := []spec.TaskEntry{}
	reqs := []spec.Requirement{}

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)

	assert.Nil(t, ctx.DocsToUpdate)
}

// TestBuildContextFromParts_DocsToUpdateJSON verifies JSON contains docs_to_update when set.
func TestBuildContextFromParts_DocsToUpdateJSON(t *testing.T) {
	cfg := config.ProjectConfig{
		ExecutionMode: "single",
		AgentCount:    1,
		DocsToUpdate:  []string{"README.md"},
	}
	tasks := []spec.TaskEntry{}
	reqs := []spec.Requirement{}

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)
	data, err := json.Marshal(ctx)
	require.NoError(t, err)

	output := string(data)
	assert.Contains(t, output, `"docs_to_update"`)
	assert.Contains(t, output, `"README.md"`)
}

// TestBuildContextFromParts_DocsToUpdateJSONOmitEmpty verifies JSON does NOT contain docs_to_update when nil.
func TestBuildContextFromParts_DocsToUpdateJSONOmitEmpty(t *testing.T) {
	cfg := config.Defaults()
	tasks := []spec.TaskEntry{}
	reqs := []spec.Requirement{}

	ctx := BuildContextFromParts("my-change", tasks, reqs, cfg)
	data, err := json.Marshal(ctx)
	require.NoError(t, err)

	output := string(data)
	assert.NotContains(t, output, `"docs_to_update"`)
}

// TestTaskItemJSON_OmitEmpty verifies TaskItem with nil new fields does NOT emit those keys in JSON.
func TestTaskItemJSON_OmitEmpty(t *testing.T) {
	ti := TaskItem{
		ID:     1,
		Name:   "A",
		Status: "pending",
	}

	data, err := json.Marshal(ti)
	require.NoError(t, err)

	output := string(data)
	assert.NotContains(t, output, `"depends"`)
	assert.NotContains(t, output, `"files"`)
	assert.NotContains(t, output, `"satisfies"`)
	assert.NotContains(t, output, `"skills"`)
}
