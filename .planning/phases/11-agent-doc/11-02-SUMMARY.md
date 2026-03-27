---
phase: 11-agent-doc
plan: "02"
subsystem: plugin
tags: [skill-md, auto-chain, spec-writer, verifier, workflow-automation]

requires:
  - phase: 11-agent-doc plan 01
    provides: mysd docs command and DocsToUpdate config field

provides:
  - Auto-spec chain in mysd-propose.md (Step 11 invokes mysd-spec-writer after proposal)
  - Auto-verify chain in mysd-apply.md (Step 5 invokes go build + go test + mysd-verifier after execution)
  - --skip-spec flag in propose to bypass auto-spec generation

affects:
  - mysd-propose workflow (now auto-generates specs after proposal)
  - mysd-apply workflow (now auto-verifies after task execution)
  - Users who previously ran /mysd:spec and /mysd:verify manually

tech-stack:
  added: []
  patterns:
    - "Auto-chain pattern: SKILL.md orchestrator invokes downstream agent via Task tool after completing its primary flow"
    - "--skip-spec flag pattern: optional bypass for auto-chain steps"
    - "auto_mode propagation: --auto flag skips confirmation prompts in downstream verifier step"

key-files:
  created: []
  modified:
    - .claude/commands/mysd-propose.md
    - .claude/commands/mysd-apply.md

key-decisions:
  - "propose Step 11 auto-chains to mysd-spec-writer immediately after proposal completion (D-01)"
  - "apply Step 5 auto-chains to verifier only after build+test pass — build/test failures are early-exit gates (D-02)"
  - "--skip-spec flag in propose bypasses Step 11 auto-spec chain; documented in argument-hint and Step 1 parsing"
  - "auto_mode=true in apply Step 5b skips confirmation prompt before verifier (D-05)"
  - "Step 4 in apply no longer says 'Run /mysd:verify' — replaced by auto-chain marker 'Proceeding to auto-verify...'"

patterns-established:
  - "Auto-chain pattern: after completing primary task, SKILL.md automatically invokes next-stage agent via Task tool"
  - "Early-exit gate pattern: build and test checks in Step 5a prevent wasted verifier invocation on broken code"

requirements-completed:
  - D-01
  - D-02
  - D-04
  - D-05

duration: 11min
completed: 2026-03-27
---

# Phase 11 Plan 02: Workflow auto-chain — propose auto-spec + apply auto-verify Summary

**propose auto-chains to mysd-spec-writer after proposal (--skip-spec bypass), apply auto-chains to verifier after go build+test pass (auto_mode skips confirmation)**

## Performance

- **Duration:** 11 min
- **Started:** 2026-03-27T02:11:28Z
- **Completed:** 2026-03-27T02:22:28Z
- **Tasks:** 2
- **Files modified:** 2

## Accomplishments
- mysd-propose.md now auto-invokes mysd-spec-writer via Task tool after Step 10 (proposal complete)
- --skip-spec flag added to propose to bypass auto-spec generation when user wants manual control
- mysd-apply.md now auto-invokes build check, test check, and mysd-verifier via Task tool after Step 4 (execution complete)
- Build or test failures in Step 5a produce clear messages and stop before wasting verifier invocation
- auto_mode=true in apply skips the verification confirmation prompt (D-05 compliance)

## Task Commits

Each task was committed atomically:

1. **Task 1: Add Step 11 auto-spec to mysd-propose.md** - `a9e4225` (feat)
2. **Task 2: Add Step 5 auto-verify to mysd-apply.md** - `c5d0bf9` (feat)

## Files Created/Modified
- `.claude/commands/mysd-propose.md` - Added --skip-spec flag, updated argument-hint, Step 1 parsing, Step 10 point 6, new Step 11 auto-spec chain with mysd-spec-writer Task tool invocation and next steps display
- `.claude/commands/mysd-apply.md` - Updated Step 4 next-step hint, added Step 5 with Step 5a build+test gates and Step 5b verifier invocation with auto_mode conditional

## Decisions Made
- propose Step 11 wraps the entire auto-spec flow: read proposal, list existing specs, invoke mysd-spec-writer per capability area, display MUST/SHOULD/MAY summary and next steps
- apply Step 5a uses two sequential gates (build, then test) before allowing verifier to run — prevents false negatives from broken code
- apply Step 5b uses `mysd execute --context-only` (not `mysd verify --context-only`) to get fresh must/should/may items — reuses the execution context endpoint already wired in Step 2
- Step 10 in propose removes the "Run /mysd:spec" manual step hint since Step 11 handles it automatically

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- Workflow auto-chain is complete for propose → spec → apply → verify
- Plans 03-05 in Phase 11 can proceed independently
- The auto-chain pattern established here (Task tool invocation at end of SKILL.md flow) is reusable for future orchestrators

---
*Phase: 11-agent-doc*
*Completed: 2026-03-27*
