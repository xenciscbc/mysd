---
phase: 09-interactive-discovery-integration
plan: 01
subsystem: cli
tags: [go, cobra, json, deferred-notes, tdd]

requires:
  - phase: 07-new-binary-commands-scanner-refactor
    provides: DetectSpecDir, cobra subcommand patterns, cmd test patterns with t.TempDir+os.Chdir

provides:
  - DeferredNote struct with JSON tags in internal/spec/deferred.go
  - DeferredStore CRUD (Load/Save/Add/Delete/Count) in internal/spec/deferred.go
  - mysd note list/add/delete cobra subcommands in cmd/note.go
  - mysd status deferred notes count line (shown only when count > 0)

affects:
  - 09-02 (SKILL.md orchestrators that use note command via binary invocation)
  - 09-03 (agent definitions that reference deferred notes)
  - 09-04 (interactive discovery integration that uses DeferredStore)

tech-stack:
  added: []
  patterns:
    - "DeferredStore zero-value on missing file (convention-over-config, same as ReadOpenSpecConfig)"
    - "Auto-increment ID uses max(existing)+1, never reuses deleted IDs"
    - "SaveDeferredStore uses json.MarshalIndent for human-readable deferred.json"
    - "noteCmd default RunE = runNoteList (list as default action)"
    - "CountDeferredNotes delegates to LoadDeferredStore (single responsibility)"

key-files:
  created:
    - internal/spec/deferred.go
    - internal/spec/deferred_test.go
    - cmd/note.go
    - cmd/note_test.go
  modified:
    - cmd/status.go

key-decisions:
  - "DeferredStore stored as deferred.json in specDir root (not changes/<name>/) — scope-free, not tied to active change"
  - "ID auto-increment: max(existing IDs)+1 ensures no reuse after delete, O(n) scan acceptable for small note lists"
  - "status shows deferred count only when count > 0 — no visual noise for clean projects"
  - "noteCmd.RunE = runNoteList so `mysd note` without subcommand lists notes"

patterns-established:
  - "LoadDeferredStore: os.IsNotExist returns zero-value (not error) — convention-over-config from Phase 1"
  - "Error wrapping: fmt.Errorf(deferred: %w, err) prefix pattern"
  - "cmd test pattern: t.TempDir() + os.Chdir + rootCmd.SetOut/SetErr/SetArgs + Execute()"

requirements-completed:
  - DISC-08

duration: 4min
completed: 2026-03-26
---

# Phase 09 Plan 01: Deferred Notes CRUD + mysd note CLI Summary

**DeferredNote CRUD package (LoadDeferredStore/Add/Delete/CountDeferredNotes) with auto-increment IDs, `mysd note list/add/delete` cobra subcommands, and `mysd status` deferred count integration**

## Performance

- **Duration:** 4 min
- **Started:** 2026-03-26T06:31:16Z
- **Completed:** 2026-03-26T06:34:59Z
- **Tasks:** 2
- **Files modified:** 5

## Accomplishments

- DeferredNote CRUD package: LoadDeferredStore (zero-value on missing), SaveDeferredStore (indented JSON), Add (auto-increment ID, RFC3339 CreatedAt), Delete (bool return), CountDeferredNotes
- mysd note cobra subcommands: list (default action), add (joins multi-word args), delete (error on not-found)
- mysd status shows "Deferred notes: N — run /mysd:note to browse" when count > 0, omits line when empty
- 17 tests total: 9 unit tests in internal/spec, 8 integration tests in cmd — all green with no regressions

## Task Commits

Each task was committed atomically:

1. **Task 1: DeferredNote CRUD package + unit tests** - `b0774bb` (feat)
2. **Task 2: mysd note cobra subcommand + status deferred count** - `3faa2cf` (feat)

## Files Created/Modified

- `internal/spec/deferred.go` - DeferredNote/DeferredStore structs, Load/Save/Add/Delete/Count functions
- `internal/spec/deferred_test.go` - 9 unit tests covering all behaviors including ID non-reuse
- `cmd/note.go` - noteCmd (list default), noteAddCmd, noteDeleteCmd cobra subcommands
- `cmd/note_test.go` - 8 integration tests: note list/add/delete + status deferred count
- `cmd/status.go` - Added fmt import + CountDeferredNotes call after RenderStatus

## Decisions Made

- DeferredStore stored as `deferred.json` in specDir root (not tied to active change) — scope-free design, notes persist across change lifecycle
- ID auto-increment via `max(existing)+1` — simple O(n) scan, ensures deleted IDs never reused, consistent with DISC-08 requirement
- `noteCmd.RunE = runNoteList` so `mysd note` without subcommand lists notes (user-friendly default)
- status shows deferred count only when `count > 0` — no visual noise for clean projects (D-09 scope guardrail)

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- DeferredStore binary infrastructure complete — SKILL.md orchestrators in Plan 02 can invoke `mysd note add/list/delete` directly
- CountDeferredNotes exposed for scope guardrail checks (DISC-08)
- All existing tests pass — no regressions from Plan 01 changes

---
*Phase: 09-interactive-discovery-integration*
*Completed: 2026-03-26*
