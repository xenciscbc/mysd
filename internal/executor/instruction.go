package executor

import (
	"fmt"
	"strings"
)

// GenerateInstruction produces a dynamic natural-language instruction string
// based on the current execution state and optional preflight results.
// It is a pure function with no I/O side effects.
func GenerateInstruction(ctx ExecutionContext, preflight *PreflightReport) string {
	var segments []string

	// Task state segment (mutually exclusive, priority order)
	segments = append(segments, taskStateSegment(ctx))

	// Preflight segments (additive)
	if preflight != nil {
		if s := staleSegment(preflight); s != "" {
			segments = append(segments, s)
		}
		if s := missingFilesSegment(preflight); s != "" {
			segments = append(segments, s)
		}
	}

	return strings.Join(segments, "\n")
}

// taskStateSegment returns exactly one segment based on task status.
// Priority: all_done > has_failed > last_task > resume > first_run.
func taskStateSegment(ctx ExecutionContext) string {
	total := len(ctx.Tasks)
	if total == 0 {
		return "No tasks found."
	}

	var doneCount, failedCount int
	var failedIDs []string
	for _, t := range ctx.Tasks {
		switch t.Status {
		case "done":
			doneCount++
		case "blocked", "failed":
			failedCount++
			failedIDs = append(failedIDs, fmt.Sprintf("T%d", t.ID))
		}
	}

	pendingCount := len(ctx.PendingTasks)

	// 1. all_done
	if doneCount == total {
		return fmt.Sprintf("All %d tasks complete. Proceed to verify or archive.", total)
	}

	// 2. has_failed
	if failedCount > 0 {
		return fmt.Sprintf("%s failed. Retry or skip before continuing.", strings.Join(failedIDs, ", "))
	}

	// 3. last_task (exactly 1 pending, rest done)
	if pendingCount == 1 && doneCount == total-1 {
		return fmt.Sprintf("Last task: T%d. Verify follows after completion.", ctx.PendingTasks[0].ID)
	}

	// 4. resume (some done, some pending)
	if doneCount > 0 && pendingCount > 0 {
		return fmt.Sprintf("Resuming: %d/%d complete. Continue from T%d.", doneCount, total, ctx.PendingTasks[0].ID)
	}

	// 5. first_run (all pending)
	if pendingCount > 0 {
		return fmt.Sprintf("%d tasks pending. Start from T%d.", total, ctx.PendingTasks[0].ID)
	}

	return "No actionable tasks."
}

func staleSegment(preflight *PreflightReport) string {
	if preflight.Checks.Staleness.IsStale && preflight.Checks.Staleness.DaysSinceLastPlan > 0 {
		return fmt.Sprintf("Warning: %d days since last plan. Consider re-planning.", preflight.Checks.Staleness.DaysSinceLastPlan)
	}
	return ""
}

func missingFilesSegment(preflight *PreflightReport) string {
	count := len(preflight.Checks.MissingFiles)
	if count > 0 {
		return fmt.Sprintf("Warning: %d missing files detected. Review before starting.", count)
	}
	return ""
}
