---
phase: 06-executor-wave-grouping-worktree-engine
plan: "01"
subsystem: executor
tags: [wave-grouping, topological-sort, kahn-algorithm, parallel-execution, execution-context]

requires:
  - phase: 05-schema-foundation-plan-checker
    provides: TaskItem struct with Depends and Files fields in executor package

provides:
  - BuildWaveGroups pure function (Kahn's BFS topological sort + file overlap split)
  - HasParallelOpportunity pure function for SKILL.md decision gate
  - ErrCyclicDependency sentinel error
  - ExecutionContext extended with wave_groups, worktree_dir, auto_mode, has_parallel_opportunity

affects:
  - 06-02 (WorktreeManager uses worktree_dir from ExecutionContext)
  - 06-03 (SKILL.md wave mode reads wave_groups and has_parallel_opportunity)
  - 06-04 (conflict resolution reads wave_groups)

tech-stack:
  added: []
  patterns:
    - "Pure function wave computation: BuildWaveGroups is I/O-free, testable in isolation"
    - "Additive-only struct extension (D-11): new fields appended at END of ExecutionContext"
    - "Cycle error silent in BuildContextFromParts: nil WaveGroups triggers sequential fallback in SKILL.md"
    - "TDD RED-GREEN pattern: failing tests committed before implementation"

key-files:
  created:
    - internal/executor/waves.go
    - internal/executor/waves_test.go
  modified:
    - internal/executor/context.go
    - internal/executor/context_test.go

key-decisions:
  - "ErrCyclicDependency returned (not silent skip) — RESEARCH.md pattern updated to add cycle detection"
  - "BuildContextFromParts ignores cycle error: WaveGroups nil, SKILL.md falls back to sequential execution"
  - "splitLayer uses greedy first-fit (not optimal) — sufficient for task counts in practice"
  - "hasFileConflict uses exact string matching (case-sensitive) — consistent with file path semantics"

patterns-established:
  - "Pattern: wave computation is a pure function layer — no I/O, takes []TaskItem, returns [][]TaskItem"
  - "Pattern: ExecutionContext additive extension — new fields at end of struct, all omitempty where optional"

requirements-completed:
  - FEXEC-01
  - FEXEC-02

duration: 12min
completed: "2026-03-25"
---

# Phase 06 Plan 01: Wave Grouping Algorithm Summary

**Kahn's topological sort + file overlap split in pure Go, with ExecutionContext extended to emit wave_groups, worktree_dir, auto_mode, has_parallel_opportunity for SKILL.md wave execution**

## Performance

- **Duration:** ~12 min
- **Started:** 2026-03-25T08:10:00Z
- **Completed:** 2026-03-25T08:22:00Z
- **Tasks:** 2 (TDD: 3 commits for Task 1)
- **Files modified:** 4

## Accomplishments

- `BuildWaveGroups` implements Kahn's BFS algorithm with cycle detection — returns `ErrCyclicDependency` instead of silently dropping tasks
- `HasParallelOpportunity` provides the SKILL.md D-03 decision gate (returns true if any task has Depends or Files)
- `splitByFileOverlap` / `splitLayer` / `hasFileConflict` enforce file isolation within waves using greedy first-fit
- `ExecutionContext` extended with 4 new fields per D-11 additive-only pattern — all existing tests still pass
- 11 wave grouping tests + 3 context integration tests cover all specified behaviors

## Task Commits

Each task was committed atomically:

1. **Task 1 RED: Failing tests for BuildWaveGroups + HasParallelOpportunity** - `bca26c5` (test)
2. **Task 1 GREEN: Implement waves.go** - `90bff63` (feat)
3. **Task 2: Extend ExecutionContext** - `c9688e9` (feat)

## Files Created/Modified

- `internal/executor/waves.go` — BuildWaveGroups, HasParallelOpportunity, ErrCyclicDependency, splitByFileOverlap, splitLayer, hasFileConflict (132 lines, pure functions)
- `internal/executor/waves_test.go` — 11 test cases covering empty/no-deps/linear-chain/diamond/cycle/file-overlap/no-overlap/deterministic/parallel-opportunity (151 lines)
- `internal/executor/context.go` — WaveGroups, WorktreeDir, AutoMode, HasParallelOpp fields added; BuildContextFromParts calls BuildWaveGroups
- `internal/executor/context_test.go` — 3 new tests: TestBuildContextFromParts_WaveGroups, _NoParallel, _AutoMode; JSONMarshal updated to assert has_parallel_opportunity

## Decisions Made

- **ErrCyclicDependency on cycle:** RESEARCH.md code silently skipped cycles — plan specifies explicit error. Added cycle detection (processed count check) to make contract clear.
- **Silent ignore in BuildContextFromParts:** Cycle error deliberately ignored so context building never fails. SKILL.md falls back to sequential when WaveGroups is nil.
- **Greedy first-fit for splitLayer:** Simpler than optimal bin-packing; sufficient for practical task counts (rarely >10 tasks per wave).

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- `BuildWaveGroups` and `HasParallelOpportunity` ready for use in Phase 06-02 (WorktreeManager)
- `ExecutionContext.WaveGroups` populated automatically in all `BuildContextFromParts` calls
- `ExecutionContext.WorktreeDir` and `AutoMode` available for worktree lifecycle management
- All 31 executor tests pass; `go build ./...` and `go vet ./internal/executor/...` clean

## Self-Check: PASSED

All files confirmed present:
- `internal/executor/waves.go` - FOUND
- `internal/executor/waves_test.go` - FOUND
- `.planning/phases/06-executor-wave-grouping-worktree-engine/06-01-SUMMARY.md` - FOUND

All commits confirmed:
- `bca26c5` test(06-01): add failing tests - FOUND
- `90bff63` feat(06-01): implement waves.go - FOUND
- `c9688e9` feat(06-01): extend ExecutionContext - FOUND

---
*Phase: 06-executor-wave-grouping-worktree-engine*
*Completed: 2026-03-25*
