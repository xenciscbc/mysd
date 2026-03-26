---
phase: 08-skill-md-orchestrators-agent-definitions
plan: 02
subsystem: plugin
tags: [agent-definitions, executor, spec-writer, task-tool-audit, per-task-spawn]

requires:
  - phase: 07-new-binary-commands-scanner-refactor
    provides: mysd-fast-forward agent and all existing agent definitions

provides:
  - mysd-executor rewritten for per-task spawn model (assigned_task only, no pending_tasks loop)
  - mysd-spec-writer rewritten for per-spec-file model (capability_area, one file per invocation)
  - All 9 agent definitions pass Task tool audit (zero violations)
  - .claude/agents/ and plugin/agents/ in sync for executor and spec-writer

affects:
  - 08-03 onwards (SKILL.md orchestrators spawn these agents per-task/per-capability)

tech-stack:
  added: []
  patterns:
    - "Per-task spawn: SKILL.md orchestrator loops over tasks and spawns executor once per task"
    - "Per-spec spawn: SKILL.md orchestrator loops over capability areas and spawns spec-writer once per area"
    - "No Task tool in agent definitions: leaf agents cannot spawn sub-agents"

key-files:
  created: []
  modified:
    - .claude/agents/mysd-executor.md
    - plugin/agents/mysd-executor.md
    - .claude/agents/mysd-spec-writer.md
    - plugin/agents/mysd-spec-writer.md
    - .claude/agents/mysd-fast-forward.md

key-decisions:
  - "mysd-executor: assigned_task is now the ONLY task input — no pending_tasks list, no execution_mode field"
  - "mysd-spec-writer: capability_area + auto_mode added; Discuss step and state transition (mysd spec) removed — these are SKILL.md orchestrator responsibilities"
  - "Task tool audit: all 8 .claude/agents/ files have zero '  - Task' matches in allowed-tools; mysd-fast-forward label fix was cosmetic (Tasks: -> tasks.md:) to satisfy grep acceptance criteria"

patterns-established:
  - "Agent leaf contract: agents receive exactly one unit of work, no iteration, no spawn"
  - "State transitions belong to SKILL.md orchestrators, not agent definitions"

requirements-completed:
  - FAGENT-05
  - FAGENT-06
  - FAGENT-07

duration: 4min
completed: 2026-03-26
---

# Phase 08 Plan 02: Agent Definitions Rewrite — Per-Task Spawn Model Summary

**mysd-executor rewritten to single-task spawn model and mysd-spec-writer to single-capability model; all 9 agent definitions pass Task tool audit with zero violations**

## Performance

- **Duration:** 4 min
- **Started:** 2026-03-26T04:30:16Z
- **Completed:** 2026-03-26T04:33:41Z
- **Tasks:** 2
- **Files modified:** 5

## Accomplishments

- Rewrote mysd-executor: removed `pending_tasks` loop and `execution_mode` field; `assigned_task` is now the sole task input; added `auto_mode`, `worktree_path`, `branch`, `isolation` fields; completion summary changed to singular "Task completed: {name}"
- Rewrote mysd-spec-writer: added `capability_area`, `auto_mode`, `existing_spec_body` inputs; removed "Discuss Capability Priorities" step; removed "mysd spec" state transition; agent now writes exactly one spec file per invocation
- Audited all 9 agent definitions for Task tool violations: zero violations found in any allowed-tools frontmatter

## Task Commits

1. **Task 1: Rewrite mysd-executor for per-task spawn model** - `e571839` (feat)
2. **Task 2: Rewrite mysd-spec-writer per-spec-file model + audit all 9 agents** - `db6c559` (feat)

## Files Created/Modified

- `.claude/agents/mysd-executor.md` - Per-task executor: single assigned_task, worktree isolation, auto_mode
- `plugin/agents/mysd-executor.md` - Identical copy (synced)
- `.claude/agents/mysd-spec-writer.md` - Per-spec writer: single capability_area, no Discuss step, no state transition
- `plugin/agents/mysd-spec-writer.md` - Identical copy (synced)
- `.claude/agents/mysd-fast-forward.md` - Minor label fix: "Tasks:" -> "tasks.md:" to satisfy Task tool grep audit

## Decisions Made

- `assigned_task` is the only task input to mysd-executor — SKILL.md orchestrator is responsible for the loop over tasks
- State transitions (mysd spec, mysd plan, etc.) belong to SKILL.md orchestrators, not agent definitions — agents are stateless from the workflow perspective
- mysd-fast-forward Phase 4 "For each task in tasks.md" loop is self-execution (not Task tool spawn), so no violation; only label cosmetically adjusted

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] mysd-fast-forward.md grep false-positive**
- **Found during:** Task 2 (agent audit)
- **Issue:** Line `  - Tasks: .specs/changes/...` matched `grep -c '  - Task'` acceptance criteria returning 1 instead of 0
- **Fix:** Renamed label from `Tasks:` to `tasks.md:` to distinguish from allowed-tools entries
- **Files modified:** `.claude/agents/mysd-fast-forward.md`
- **Verification:** `grep -c '  - Task' mysd-fast-forward.md` returns 0
- **Committed in:** db6c559 (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (Rule 1 - cosmetic label fix in fast-forward agent)
**Impact on plan:** Minimal — label change only, no behavioral change to the agent.

## Issues Encountered

None — all agents already had zero Task tool entries in allowed-tools. Only the grep false-positive in fast-forward required a fix.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- executor and spec-writer agents are ready for SKILL.md orchestrators that will spawn them per-task/per-capability
- All 9 agents comply with FAGENT-05 (no Task tool)
- Phase 08 Plans 03+ can safely reference these rewritten agents

---
*Phase: 08-skill-md-orchestrators-agent-definitions*
*Completed: 2026-03-26*
