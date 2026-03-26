---
phase: 08-skill-md-orchestrators-agent-definitions
plan: "05"
subsystem: plugin
tags: [skill-md, orchestrator, fix, fast-forward, pipeline, auto-mode]

requires:
  - phase: 08-03
    provides: apply orchestrator and mysd-executor per-task spawn pattern

provides:
  - /mysd:fix orchestrator with dual-path detection (merge conflict vs implementation failure)
  - /mysd:ff rewritten as plan+apply+archive pipeline without research
  - /mysd:ffe rewritten as research+plan+apply+archive pipeline
  - Both ff/ffe hardcode auto_mode=true without using mysd-fast-forward agent

affects:
  - Phase 09 interactive discovery (fix, ff, ffe are core pipeline commands)
  - Any human audit of FAGENT-05 (Task tool in agent definitions)

tech-stack:
  added: []
  patterns:
    - "fix dual-path: conflict marker detection -> merge path OR sidecar failure -> implementation path"
    - "ff/ffe direct pipeline: orchestrator spawns designer/planner/executor directly (no intermediate agent)"
    - "auto_mode hardcoded true in ff/ffe — not parsed from $ARGUMENTS"

key-files:
  created:
    - .claude/commands/mysd-fix.md
    - plugin/commands/mysd-fix.md
  modified:
    - .claude/commands/mysd-ff.md
    - .claude/commands/mysd-ffe.md
    - plugin/commands/mysd-ff.md
    - plugin/commands/mysd-ffe.md

key-decisions:
  - "fix uses safety valve: auto-detects path but confirms with user before proceeding (D-08)"
  - "ff has no research step per FAUTO-04; ffe has 4-dim parallel research per D-25"
  - "ff/ffe do not use mysd-fast-forward agent — they directly orchestrate the pipeline (D-24/D-25)"
  - "fix transitive downstream task recovery: all skipped tasks restored to pending after successful merge"
  - "abandon path returns task to pending and cleans up worktree+branch without re-execution"

patterns-established:
  - "Pipeline orchestrator pattern: SKILL.md spawns designer -> planner -> executor directly"
  - "auto_mode=true hardcoded in ff/ffe — not user-configurable, always autonomous"
  - "Dual-path fix: conflict markers check first, implementation failure fallback"

requirements-completed:
  - FCMD-02
  - FAUTO-03
  - FAUTO-04

duration: 10min
completed: 2026-03-26
---

# Phase 08 Plan 05: Fix Orchestrator & FF/FFE Pipeline Rewrite Summary

**/mysd:fix with dual-path conflict/implementation detection, /mysd:ff and /mysd:ffe rewritten as direct pipeline orchestrators with hardcoded auto_mode=true**

## Performance

- **Duration:** ~10 min
- **Started:** 2026-03-26T04:54:00Z
- **Completed:** 2026-03-26T05:04:41Z
- **Tasks:** 2 of 3 (Task 3 is checkpoint:human-verify — awaiting human audit)
- **Files modified:** 6

## Accomplishments

- Created `/mysd:fix` orchestrator with full dual-path detection: merge conflict path (detect markers, resolve, build/test, merge, cleanup) and implementation failure path (diagnose, optional research via mysd-researcher, re-execute via mysd-executor)
- Added abandon path to `/mysd:fix` (task-update to pending, worktree+branch cleanup)
- Added transitive downstream task recovery: after successful fix, all skipped dependent tasks restored to pending
- Rewrote `/mysd:ff` as direct plan+apply+archive pipeline (no mysd-fast-forward agent, no research, auto_mode=true always)
- Rewrote `/mysd:ffe` as research+plan+apply+archive pipeline (4-dim parallel mysd-researcher, no mysd-fast-forward agent, auto_mode=true always)

## Task Commits

1. **Task 1: Create /mysd:fix SKILL.md orchestrator** - `730dcf4` (feat)
2. **Task 2: Rewrite /mysd:ff and /mysd:ffe for direct pipeline model** - `325a103` (feat)

## Files Created/Modified

- `.claude/commands/mysd-fix.md` - Fix orchestrator with dual-path detection (merge conflict, implementation, abandon)
- `plugin/commands/mysd-fix.md` - Identical copy for plugin distribution
- `.claude/commands/mysd-ff.md` - Rewritten: plan+apply+archive (no research, auto_mode=true)
- `plugin/commands/mysd-ff.md` - Identical copy for plugin distribution
- `.claude/commands/mysd-ffe.md` - Rewritten: research+plan+apply+archive (auto_mode=true)
- `plugin/commands/mysd-ffe.md` - Identical copy for plugin distribution

## Decisions Made

- fix uses D-08 safety valve: auto-detects path (conflict markers vs sidecar failure) but confirms with user before acting — detection is not silent
- ff has zero research steps per FAUTO-04 — research_findings passed as empty array to designer
- ffe spawns 4 mysd-researcher agents in parallel per D-25 for 4 dimensions: codebase/domain/architecture/pitfalls
- Both ff and ffe remove dependency on mysd-fast-forward agent; they directly orchestrate designer → planner → executor pipeline
- fix abandon path: `mysd task-update {id} pending` + worktree/branch cleanup — no re-execution attempted
- Transitive downstream recovery: after fix merge, runs execute --context-only, finds all skipped tasks depending on fixed task, restores each to pending

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Task 3 checkpoint: human audit required for FAGENT-05 (zero Task tool references in all 9 agent definitions)
- After audit approval, Phase 8 Plan 05 is complete
- Phase 8 is the final phase of v1.1 milestone — all SKILL.md orchestrators and agent definitions are ready

## Self-Check: PASSED

Files verified:
- `.claude/commands/mysd-fix.md` - FOUND
- `plugin/commands/mysd-fix.md` - FOUND (identical to .claude copy)
- `.claude/commands/mysd-ff.md` - FOUND
- `plugin/commands/mysd-ff.md` - FOUND (identical to .claude copy)
- `.claude/commands/mysd-ffe.md` - FOUND
- `plugin/commands/mysd-ffe.md` - FOUND (identical to .claude copy)

Commits verified:
- `730dcf4` - Task 1 commit FOUND
- `325a103` - Task 2 commit FOUND

---
*Phase: 08-skill-md-orchestrators-agent-definitions*
*Completed: 2026-03-26*
