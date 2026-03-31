package executor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateInstruction_FirstRun(t *testing.T) {
	ctx := ExecutionContext{
		Tasks: []TaskItem{
			{ID: 1, Status: "pending"},
			{ID: 2, Status: "pending"},
			{ID: 3, Status: "pending"},
		},
		PendingTasks: []TaskItem{
			{ID: 1, Status: "pending"},
			{ID: 2, Status: "pending"},
			{ID: 3, Status: "pending"},
		},
	}
	result := GenerateInstruction(ctx, nil)
	assert.Contains(t, result, "3 tasks pending")
	assert.Contains(t, result, "T1")
}

func TestGenerateInstruction_Resume(t *testing.T) {
	ctx := ExecutionContext{
		Tasks: []TaskItem{
			{ID: 1, Status: "done"},
			{ID: 2, Status: "done"},
			{ID: 3, Status: "done"},
			{ID: 4, Status: "pending"},
			{ID: 5, Status: "pending"},
		},
		PendingTasks: []TaskItem{
			{ID: 4, Status: "pending"},
			{ID: 5, Status: "pending"},
		},
	}
	result := GenerateInstruction(ctx, nil)
	assert.Contains(t, result, "3/5 complete")
	assert.Contains(t, result, "T4")
}

func TestGenerateInstruction_LastTask(t *testing.T) {
	ctx := ExecutionContext{
		Tasks: []TaskItem{
			{ID: 1, Status: "done"},
			{ID: 2, Status: "done"},
			{ID: 3, Status: "done"},
			{ID: 4, Status: "done"},
			{ID: 5, Status: "pending"},
		},
		PendingTasks: []TaskItem{
			{ID: 5, Status: "pending"},
		},
	}
	result := GenerateInstruction(ctx, nil)
	assert.Contains(t, result, "T5")
	assert.Contains(t, result, "erify")
}

func TestGenerateInstruction_AllDone(t *testing.T) {
	ctx := ExecutionContext{
		Tasks: []TaskItem{
			{ID: 1, Status: "done"},
			{ID: 2, Status: "done"},
			{ID: 3, Status: "done"},
			{ID: 4, Status: "done"},
			{ID: 5, Status: "done"},
		},
		PendingTasks: []TaskItem{},
	}
	result := GenerateInstruction(ctx, nil)
	assert.Contains(t, result, "All 5 tasks complete")
	assert.Contains(t, result, "erify")
}

func TestGenerateInstruction_HasFailed(t *testing.T) {
	ctx := ExecutionContext{
		Tasks: []TaskItem{
			{ID: 1, Status: "done"},
			{ID: 2, Status: "done"},
			{ID: 3, Status: "blocked"},
			{ID: 4, Status: "pending"},
		},
		PendingTasks: []TaskItem{
			{ID: 4, Status: "pending"},
		},
	}
	result := GenerateInstruction(ctx, nil)
	assert.Contains(t, result, "T3")
	assert.Contains(t, result, "etry")
}

func TestGenerateInstruction_Stale(t *testing.T) {
	ctx := ExecutionContext{
		Tasks:        []TaskItem{{ID: 1, Status: "pending"}},
		PendingTasks: []TaskItem{{ID: 1, Status: "pending"}},
	}
	preflight := &PreflightReport{
		Status: "warning",
		Checks: PreflightChecks{
			MissingFiles: []string{},
			Staleness: StalenessCheck{
				DaysSinceLastPlan: 15,
				IsStale:           true,
			},
		},
	}
	result := GenerateInstruction(ctx, preflight)
	assert.Contains(t, result, "15 days since last plan")
}

func TestGenerateInstruction_MissingFiles(t *testing.T) {
	ctx := ExecutionContext{
		Tasks:        []TaskItem{{ID: 1, Status: "pending"}},
		PendingTasks: []TaskItem{{ID: 1, Status: "pending"}},
	}
	preflight := &PreflightReport{
		Status: "warning",
		Checks: PreflightChecks{
			MissingFiles: []string{"foo.go", "bar.go"},
			Staleness:    StalenessCheck{DaysSinceLastPlan: 1, IsStale: false},
		},
	}
	result := GenerateInstruction(ctx, preflight)
	assert.Contains(t, result, "2 missing files")
}

func TestGenerateInstruction_Combined_ResumeAndStale(t *testing.T) {
	ctx := ExecutionContext{
		Tasks: []TaskItem{
			{ID: 1, Status: "done"},
			{ID: 2, Status: "done"},
			{ID: 3, Status: "done"},
			{ID: 4, Status: "pending"},
			{ID: 5, Status: "pending"},
		},
		PendingTasks: []TaskItem{
			{ID: 4, Status: "pending"},
			{ID: 5, Status: "pending"},
		},
	}
	preflight := &PreflightReport{
		Status: "warning",
		Checks: PreflightChecks{
			MissingFiles: []string{},
			Staleness: StalenessCheck{
				DaysSinceLastPlan: 10,
				IsStale:           true,
			},
		},
	}
	result := GenerateInstruction(ctx, preflight)
	lines := strings.Split(result, "\n")
	assert.GreaterOrEqual(t, len(lines), 2)
	assert.Contains(t, result, "3/5 complete")
	assert.Contains(t, result, "10 days since last plan")
}

func TestGenerateInstruction_NilPreflight(t *testing.T) {
	ctx := ExecutionContext{
		Tasks:        []TaskItem{{ID: 1, Status: "pending"}},
		PendingTasks: []TaskItem{{ID: 1, Status: "pending"}},
	}
	result := GenerateInstruction(ctx, nil)
	assert.NotContains(t, result, "missing files")
	assert.NotContains(t, result, "days since")
}
