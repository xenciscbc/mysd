---
phase: 07-new-binary-commands-scanner-refactor
plan: "01"
subsystem: executor
tags: [tdd, wave-groups, dependency-propagation, executor]
dependency_graph:
  requires: []
  provides: [FilterBlockedTasks]
  affects: [internal/executor/waves.go, internal/executor/waves_test.go]
tech_stack:
  added: []
  patterns: [BFS adjacency graph, TDD red-green]
key_files:
  created: []
  modified:
    - internal/executor/waves.go
    - internal/executor/waves_test.go
decisions:
  - "FilterBlockedTasks uses BFS over adjacency map (same pattern as BuildWaveGroups) — no new data structures needed"
  - "Failed tasks excluded from result alongside blocked tasks — single blocked set covers both"
metrics:
  duration_min: 8
  completed_date: "2026-03-26"
  tasks_completed: 1
  files_modified: 2
---

# Phase 07 Plan 01: FilterBlockedTasks — Dependency Failure Propagation Summary

**One-liner:** BFS-based FilterBlockedTasks in waves.go that transitively excludes failed and blocked tasks from execution queue.

## What Was Built

Added `FilterBlockedTasks(tasks []TaskItem, failedIDs []int) []TaskItem` to `internal/executor/waves.go`. When a task fails during wave execution, all downstream dependent tasks (transitively) are excluded from the returned slice. This enables SKILL.md to query the binary for executable tasks per wave without re-running blocked work.

## Tasks Completed

| Task | Description | Commit | Files |
|------|-------------|--------|-------|
| RED | Write 6 failing test cases for FilterBlockedTasks | ef33882 | waves_test.go |
| GREEN | Implement FilterBlockedTasks BFS algorithm | 54ed33a | waves.go |

## Implementation Details

Algorithm (per D-13):
1. Build adjacency map: `dep -> []dependents` (identical to BuildWaveGroups adj pattern)
2. Seed `blocked` set from `failedIDs`
3. BFS from each failed ID: add all downstream dependents to blocked set
4. Return filtered slice: tasks where `task.ID` not in blocked set

The BFS is fully transitive — a queue ensures T1 fail → T2 (depends T1) → T3 (depends T2) are ALL blocked, not just direct dependents.

## Test Coverage

6 test cases in `TestFilterBlockedTasks_*`:
- `EmptyFailedIDs`: all tasks returned when no failures
- `DirectDependency`: direct dependent blocked, independent task unaffected
- `TransitivePropagation`: T1→T2→T3 chain all blocked when T1 fails
- `MultipleFailures`: both failure subtrees blocked, result empty
- `FailedExcluded`: failed task itself not in returned slice
- `NoDependencies`: independent tasks unaffected by peer failure

## Verification

```
go test ./internal/executor/... -v -run TestFilterBlockedTasks
# All 6 PASS

go test ./internal/executor/...
# ok — all existing BuildWaveGroups and HasParallelOpportunity tests still pass
```

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None.

## Self-Check: PASSED

- `internal/executor/waves.go` exists with `FilterBlockedTasks` function
- `internal/executor/waves_test.go` exists with `TestFilterBlockedTasks_*` tests
- Commit `ef33882` exists (RED phase)
- Commit `54ed33a` exists (GREEN phase)
- `go test ./internal/executor/...` exits 0
