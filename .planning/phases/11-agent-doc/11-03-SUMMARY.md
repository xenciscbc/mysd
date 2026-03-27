---
phase: 11-agent-doc
plan: "03"
subsystem: agent-doc
tags: [executor, fix, sidecar, failure-context, gitignore]

# Dependency graph
requires:
  - phase: 08-skill-md-orchestrators-agent-definitions
    provides: mysd-executor agent and mysd-fix SKILL.md baseline
provides:
  - On-failure sidecar writing in mysd-executor (Steps F1-F3)
  - Aligned sidecar reading in mysd-fix with D-06 path format
  - .sidecar/ excluded from git via .gitignore
affects:
  - 11-agent-doc (subsequent plans referencing executor/fix flow)
  - any plan that uses /mysd:fix to diagnose implementation failures

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Failure sidecar pattern: executor writes T{id}-failure.md on task failure, fix reads it for diagnosis context"
    - "Backward compat via null check: fix agent degrades gracefully when sidecar missing (D-08)"

key-files:
  created:
    - .gitignore
  modified:
    - .claude/agents/mysd-executor.md
    - .claude/commands/mysd-fix.md

key-decisions:
  - "On Failure is alternative exit path — executor must NOT proceed to Mark Task Done or Atomic Commit when failure occurs (D-06)"
  - "failure_context null when sidecar missing — backward compat for pre-D-06 task failures without sidecars (D-08)"
  - "Sidecar path follows D-06 format: .specs/changes/{change_name}/.sidecar/T{id}-failure.md"

patterns-established:
  - "Sidecar write before mark-failed: F2 writes file, F3 marks status, output guides user to /mysd:fix"
  - "Conditional diagnosis in fix agent: present sidecar context when available, reproduce error from scratch when not"

requirements-completed:
  - D-06
  - D-07
  - D-08
  - D-09

# Metrics
duration: 2min
completed: 2026-03-27
---

# Phase 11 Plan 03: Executor failure sidecar + fix alignment + .gitignore Summary

**Executor agent writes structured T{id}-failure.md sidecar on task failure; fix agent reads from D-06 canonical path with backward-compat null fallback; .sidecar/ excluded from git**

## Performance

- **Duration:** 2 min
- **Started:** 2026-03-27T03:44:48Z
- **Completed:** 2026-03-27T03:46:xx Z
- **Tasks:** 2
- **Files modified:** 3

## Accomplishments
- Added On Failure section (Steps F1-F3) to mysd-executor.md between Task Execution and Post-Execution Test Generation — executor now writes structured sidecar with frontmatter (task_id, task_name, timestamp, change_name) + body sections (Error Output, Task Description, Files Modified Before Failure, AI Diagnostic Attempts)
- Aligned mysd-fix.md Step 3 sidecar reading to exact D-06 path format `.specs/changes/{change_name}/.sidecar/T{target_task.id}-failure.md` with failure_context variable and null fallback
- Updated Step 4 path detection and Step 5B diagnose to use failure_context conditionally, with "No failure sidecar found" fallback that runs `go build ./... && go test ./...`
- Created .gitignore with .sidecar/ exclusion pattern (file did not exist previously)

## Task Commits

Each task was committed atomically:

1. **Task 1: Add on-failure sidecar writing to mysd-executor.md + update .gitignore** - `6068d4a` (feat)
2. **Task 2: Align mysd-fix.md sidecar reading path to D-06 format** - `2cca45a` (feat)

## Files Created/Modified
- `.claude/agents/mysd-executor.md` - Added On Failure section (Steps F1-F3) after Step 5 Atomic Commit, before Post-Execution Test Generation
- `.claude/commands/mysd-fix.md` - Step 3 explicit sidecar read with failure_context; Step 4 path detection using failure_context; Step 5B conditional diagnosis
- `.gitignore` - Created new file with .spectra/, openspec/.vector-search.db*, .sidecar/ entries

## Decisions Made
- On Failure path is alternative exit: executor MUST NOT proceed to Mark Task Done or Atomic Commit after writing sidecar and marking failed
- Sidecar frontmatter fields: task_id, task_name, timestamp, change_name — minimal set for fix agent identification
- Body sections: Error Output (truncated to 200 lines), Task Description, Files Modified Before Failure, AI Diagnostic Attempts — matches D-12 diagnosis requirements
- Null failure_context triggers fresh reproduction via go build/test rather than blocking — resilient to pre-D-06 task states

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- Executor now writes failure context on any task failure
- Fix agent can diagnose implementation failures with full context from sidecar
- .sidecar/ directories are gitignored project-wide
- Plan 04 and 05 of phase 11 can proceed

## Self-Check: PASSED

- FOUND: .planning/phases/11-agent-doc/11-03-SUMMARY.md
- FOUND: .claude/agents/mysd-executor.md (On Failure section at line 182)
- FOUND: .claude/commands/mysd-fix.md (.sidecar/T path at line 51)
- FOUND: .gitignore (.sidecar/ at line 3)
- FOUND: commit 6068d4a (Task 1)
- FOUND: commit 2cca45a (Task 2)
- FOUND: commit a900b4f (metadata)

---
*Phase: 11-agent-doc*
*Completed: 2026-03-27*
