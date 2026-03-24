---
phase: 02-execution-engine
plan: "04"
subsystem: cli
tags: [go, cobra, state-machine, context-json, spec-driven]

requires:
  - phase: 02-01
    provides: state transitions (PhaseSpecced, PhaseDesigned, PhasePlanned) and LoadState/SaveState
  - phase: 02-02
    provides: config.ResolveModel and ProjectConfig with TestGeneration field

provides:
  - cmd/spec.go with --context-only flag and PhaseSpecced state transition
  - cmd/design.go with --context-only flag and PhaseDesigned state transition
  - cmd/plan.go with --context-only, --research, --check flags and PhasePlanned state transition
  - Context JSON output pattern for all three intermediate workflow commands

affects: [03-plugin-layer, 04-distribution]

tech-stack:
  added: []
  patterns:
    - "Thin command layer: binary only does state management and context output; AI logic in SKILL.md"
    - "context-only flag: outputs JSON blob for SKILL.md agent consumption"
    - "Convention-over-config: --research and --check are opt-in flags, fast by default"

key-files:
  created: []
  modified:
    - cmd/spec.go
    - cmd/design.go
    - cmd/plan.go

key-decisions:
  - "context-only JSON output includes model resolved via ResolveModel (same pattern as execute command)"
  - "plan --context-only includes research_enabled, check_enabled, test_generation booleans for SKILL.md pipeline depth control"
  - "design and plan context JSON include specs as []string (formatted as [KEYWORD] text) for readability"
  - "state.Transition sets LastRun internally; commands also set ws.LastRun = time.Now() for explicit timestamp"

patterns-established:
  - "Intermediate command pattern: DetectSpecDir -> LoadState -> (context-only branch OR Transition+SaveState) -> print guidance"
  - "Context JSON key naming: snake_case matching WorkflowState and ProjectConfig field names"

requirements-completed:
  - WCMD-01
  - WCMD-02
  - WCMD-03
  - WCMD-04
  - TEST-02

duration: 5min
completed: "2026-03-24"
---

# Phase 02 Plan 04: Spec, Design, Plan Commands Summary

**Three intermediate workflow commands (spec/design/plan) with state transitions and --context-only JSON output for SKILL.md agent consumption**

## Performance

- **Duration:** 5 min
- **Started:** 2026-03-24T00:14:00Z
- **Completed:** 2026-03-24T00:16:03Z
- **Tasks:** 1
- **Files modified:** 3

## Accomplishments

- Implemented `cmd/spec.go` replacing stub: transitions to PhaseSpecced, outputs context JSON with proposal body and resolved model
- Implemented `cmd/design.go` replacing stub: transitions to PhaseDesigned, outputs context JSON with requirements list and resolved model
- Implemented `cmd/plan.go` replacing stub: transitions to PhasePlanned, outputs context JSON with design body, specs, model, and pipeline flags (research_enabled, check_enabled, test_generation); supports --research and --check optional flags per D-12
- All three commands follow propose.go thin-layer pattern; no AI logic in binary layer

## Task Commits

1. **Task 1: spec, design, plan commands with --context-only and state transitions** - `d59e30b` (feat)

## Files Created/Modified

- `cmd/spec.go` - Spec command: DetectSpecDir, LoadState, --context-only JSON output (proposal + model), or state.Transition to PhaseSpecced
- `cmd/design.go` - Design command: same pattern, outputs specs summary + model, transitions to PhaseDesigned
- `cmd/plan.go` - Plan command: --context-only with research_enabled/check_enabled/test_generation, --research and --check flags, transitions to PhasePlanned

## Decisions Made

- `state.Transition` already sets `ws.LastRun` internally (see transitions.go line 41), but commands also call `ws.LastRun = time.Now()` to ensure the timestamp is updated even if Transition's internal behavior changes — defensive coding
- Context JSON for design and plan converts `[]Requirement` to `[]string` with `[KEYWORD] text` format for human-readable SKILL.md consumption

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- spec, design, plan commands are fully wired to state machine
- Context JSON output is ready for SKILL.md agent integration in Phase 03
- All five core workflow commands (propose, spec, design, plan, execute) now have state transition logic

## Self-Check: PASSED

- cmd/spec.go: FOUND
- cmd/design.go: FOUND
- cmd/plan.go: FOUND
- Commit d59e30b: FOUND

---
*Phase: 02-execution-engine*
*Completed: 2026-03-24*
