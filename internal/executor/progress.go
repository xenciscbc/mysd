package executor

import "github.com/xenciscbc/mysd/internal/spec"

// PendingTasks returns tasks that have not been completed or blocked.
// Done and blocked tasks are excluded to support execution resumption (EXEC-05).
func PendingTasks(tasks []spec.TaskEntry) []spec.TaskEntry {
	var pending []spec.TaskEntry
	for _, t := range tasks {
		if t.Status != spec.StatusDone && t.Status != spec.StatusBlocked {
			pending = append(pending, t)
		}
	}
	return pending
}

// CalcProgress returns the count of done tasks and total tasks.
func CalcProgress(tasks []spec.TaskEntry) (done int, total int) {
	total = len(tasks)
	for _, t := range tasks {
		if t.Status == spec.StatusDone {
			done++
		}
	}
	return done, total
}
