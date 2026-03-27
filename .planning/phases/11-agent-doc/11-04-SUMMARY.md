---
phase: 11-agent-doc
plan: "04"
subsystem: plugin
tags: [skill-md, archive, ff, ffe, docs-update, auto-verify]

requires:
  - phase: 11-01
    provides: docs_to_update field in ExecutionContext (mysd execute --context-only JSON)
  - phase: 11-02
    provides: mysd-apply.md Step 5 auto-verify pattern to reference

provides:
  - mysd-archive.md with doc maintenance flow (Step 0 config read, Step 2 doc update)
  - mysd-ff.md with inline auto-verify (Step 4) and docs update (Step 6); 7 steps total
  - mysd-ffe.md with inline auto-verify (Step 5) and docs update (Step 7); 8 steps total

affects:
  - plugin-sync
  - ff-pipeline
  - ffe-pipeline
  - archive-workflow

tech-stack:
  added: []
  patterns:
    - "docs_to_update read via mysd execute --context-only JSON (no direct YAML parsing in SKILL.md)"
    - "CHANGELOG.md prepend strategy: generate new entry only, insert at top, preserve existing entries"
    - "README.md full rewrite strategy: read current content, rewrite incorporating change context"
    - "ff/ffe inline auto-verify: build+test+verifier before archive, early exit with Archive skipped message"
    - "ff/ffe inline docs update: auto_mode=true always, no user confirmation"

key-files:
  created: []
  modified:
    - .claude/commands/mysd-archive.md
    - .claude/commands/mysd-ff.md
    - .claude/commands/mysd-ffe.md

key-decisions:
  - "archive Step 0 reads docs_to_update before archive runs — enables confirmation flow before irreversible action"
  - "ff/ffe use auto_mode=true for docs update (no confirmation) — consistent with ff/ffe always-on auto mode"
  - "docs update context reads from archived .specs/archive/{change_name}/ not active change — archive must complete first"

patterns-established:
  - "Inline pipeline steps follow same pattern as SKILL.md orchestrator steps: distinct numbered steps, explicit STOP conditions"
  - "Build-fail early exit pattern: display error + message + STOP, do not proceed to next pipeline stage"
  - "Convention over config: when docs_to_update is empty, skip silently (no prompt)"

requirements-completed:
  - D-03
  - D-11
  - D-11b
  - D-13
  - D-14
  - D-17
  - D-18

duration: 8min
completed: 2026-03-27
---

# Phase 11 Plan 04: Archive doc maintenance flow + ff/ffe inline additions Summary

**Archive SKILL.md rewritten with 4-step doc maintenance flow (config read, confirm, context read, per-file update); ff/ffe pipelines extended with inline auto-verify and docs_to_update steps**

## Performance

- **Duration:** ~8 min
- **Started:** 2026-03-27T05:59:37Z
- **Completed:** 2026-03-27T06:07:00Z
- **Tasks:** 2
- **Files modified:** 3

## Accomplishments

- Rewrote mysd-archive.md from 2 steps to 4 steps with full doc maintenance flow including --auto flag support, per-file update strategies (CHANGELOG prepend, README rewrite, auto-detect for others), and user confirmation
- Extended mysd-ff.md from 5 steps to 7 steps with Step 4 Inline Auto-Verify (go build + go test + mysd-verifier with early exit) and Step 6 Inline Docs Update (reads docs_to_update, applies update strategies)
- Extended mysd-ffe.md from 6 steps to 8 steps with Step 5 Inline Auto-Verify and Step 7 Inline Docs Update, matching ff.md's inline patterns

## Task Commits

Each task was committed atomically:

1. **Task 1: Rewrite mysd-archive.md with doc maintenance flow** - `00e88bd` (feat)
2. **Task 2: Add inline auto-verify + docs_to_update to ff.md and ffe.md** - `bdea85e` (feat)

**Plan metadata:** (docs commit — see below)

## Files Created/Modified

- `.claude/commands/mysd-archive.md` - Rewritten with Step 0 (parse --auto + read docs_to_update), Step 1 (archive with all original error handling), Step 2 (doc maintenance with confirm/context/update/summary sub-steps), Step 3 (post-archive guidance)
- `.claude/commands/mysd-ff.md` - Step 4 Inline Auto-Verify + Step 6 Inline Docs Update inserted; Archive renumbered to Step 5, Confirm renumbered to Step 7 (7 steps total)
- `.claude/commands/mysd-ffe.md` - Step 5 Inline Auto-Verify + Step 7 Inline Docs Update inserted; Archive renumbered to Step 6, Confirm renumbered to Step 8 (8 steps total)

## Decisions Made

- archive reads docs_to_update before running `mysd archive` (Step 0) so the confirmation flow can happen before the irreversible archive action
- ff/ffe inline docs update always uses auto_mode=true (no user confirmation) — consistent with ff/ffe being fully automatic pipelines
- docs update context reads from `.specs/archive/{change_name}/` (archived copy) not the active change directory — archive must succeed first before doc update can reference archived content
- Build-fail in ff/ffe auto-verify shows "Archive skipped" message, not just an error — clear signal that the pipeline halted before the irreversible step

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- All three files ready for plugin sync (plan 11-05 copies to plugin/commands/)
- archive.md, ff.md, ffe.md now have identical docs_to_update patterns — consistent UX across all archive paths
- Phase 11 complete after plugin sync

## Self-Check: PASSED

- `.claude/commands/mysd-archive.md` — FOUND
- `.claude/commands/mysd-ff.md` — FOUND
- `.claude/commands/mysd-ffe.md` — FOUND
- Commit `00e88bd` — FOUND
- Commit `bdea85e` — FOUND

---
*Phase: 11-agent-doc*
*Completed: 2026-03-27*
