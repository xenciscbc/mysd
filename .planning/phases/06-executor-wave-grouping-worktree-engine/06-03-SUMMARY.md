---
phase: 06-executor-wave-grouping-worktree-engine
plan: "03"
subsystem: cmd
tags: [wave-grouping, cmd-layer, execute-command, plan-command, context-only, json-output]

requires:
  - phase: 06-executor-wave-grouping-worktree-engine
    plan: "01"
    provides: BuildWaveGroups, HasParallelOpportunity, ExecutionContext.WaveGroups
  - phase: 06-executor-wave-grouping-worktree-engine
    plan: "02"
    provides: WorktreeManager (cfg.WorktreeDir, cfg.AutoMode in config)

provides:
  - execute --context-only JSON with wave_groups (TaskItem arrays, not placeholder)
  - plan --context-only JSON with wave_groups (real computed groups, not [][]int{})
  - plan --context-only JSON with has_parallel_opportunity boolean

affects:
  - 06-04 (SKILL.md wave mode reads wave_groups from both execute and plan commands)

tech-stack:
  added: []
  patterns:
    - "cmd layer is thin: no wave logic in cmd — delegates entirely to executor.BuildWaveGroups"
    - "plan.go lazy tasks load: wave groups computed only when tasks.md exists, nil on parse failure"
    - "execute.go unchanged: BuildContext already wired wave_groups via Plan 01 BuildContextFromParts"

key-files:
  created: []
  modified:
    - cmd/execute_test.go
    - cmd/plan.go

key-decisions:
  - "execute --context-only required zero cmd changes — Plan 01 already wired WaveGroups via BuildContextFromParts"
  - "plan.go wave groups computed only in --context-only path (tasks.md may not exist at planning time — graceful nil)"
  - "plan.go reuses same tasksPath in --check branch (no double declaration, inner scope fm variable)"

requirements-completed:
  - FEXEC-03

duration: 2min
completed: "2026-03-25"
---

# Phase 06 Plan 03: Wire Wave Groups into cmd Layer Summary

**cmd/plan.go placeholder `[][]int{}` replaced with real executor.BuildWaveGroups call; cmd/execute.go confirmed already emitting wave_groups via Plan 01 wire-up; TestExecuteContextOnly_WaveGroups added to verify 2-wave structure**

## Performance

- **Duration:** ~2 min
- **Started:** 2026-03-25T08:16:29Z
- **Completed:** 2026-03-25T08:19:00Z
- **Tasks:** 2
- **Files modified:** 2

## Accomplishments

- Verified `execute --context-only` already outputs `wave_groups` with no code changes needed — `BuildContextFromParts` (Plan 01) already calls `BuildWaveGroups` and populates `ExecutionContext.WaveGroups`, `WorktreeDir`, `AutoMode`, `HasParallelOpp`
- Added `TestExecuteContextOnly_WaveGroups` test: sets up 2 tasks with `Depends`/`Files`, asserts 2-wave structure (`wave_groups[0]=[Task1]`, `wave_groups[1]=[Task2]`) and `has_parallel_opportunity=true`
- Replaced `cmd/plan.go` line 75 `[][]int{}` placeholder with real computation: convert `spec.TaskEntry` -> `executor.TaskItem`, call `executor.BuildWaveGroups`, add `has_parallel_opportunity` to context map
- Added `executor` import to `cmd/plan.go`
- All 43 cmd tests pass, `go build ./...` clean, `go vet ./...` clean

## Task Commits

1. **Task 1: TestExecuteContextOnly_WaveGroups** - `f3d4eb8` (test)
2. **Task 2: Replace plan wave_groups placeholder** - `88eb2f2` (feat)

## Files Created/Modified

- `cmd/execute_test.go` — added `TestExecuteContextOnly_WaveGroups` (46 lines): 2-task setup with depends/files, asserts wave structure and `has_parallel_opportunity`
- `cmd/plan.go` — replaced `[][]int{}` with real `executor.BuildWaveGroups` call, added `has_parallel_opportunity`, added `executor` import

## Decisions Made

- **No changes to execute.go:** Plan 01's `BuildContextFromParts` already wires `WaveGroups` — the cmd layer was already correct. Only a test was needed to confirm.
- **Lazy wave computation in plan.go:** Wave groups are computed inside the `if fm, _, parseErr := spec.ParseTasksV2(tasksPath); parseErr == nil` block. If tasks.md does not exist yet (valid state during planning), `waveGroups` remains nil and `has_parallel_opportunity` stays false — graceful degradation.
- **Inner scope for planCheck branch:** The `--check` branch now uses an inner-scope `fm` variable to avoid redeclaring the outer `fm` used for wave computation.

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None.

## Next Phase Readiness

- Both `mysd execute --context-only` and `mysd plan --context-only` now emit complete `wave_groups` JSON
- SKILL.md consumers (Phase 06-04) can read `wave_groups` and `has_parallel_opportunity` from both commands
- `worktree_dir` and `auto_mode` are also present in both outputs

## Self-Check: PASSED

Files confirmed:
- `cmd/execute_test.go` — FOUND (TestExecuteContextOnly_WaveGroups present)
- `cmd/plan.go` — FOUND (no `[][]int{}`, executor.BuildWaveGroups present)

Commits confirmed:
- `f3d4eb8` test(06-03): add TestExecuteContextOnly_WaveGroups — FOUND
- `88eb2f2` feat(06-03): replace plan --context-only wave_groups placeholder — FOUND

---
*Phase: 06-executor-wave-grouping-worktree-engine*
*Completed: 2026-03-25*
