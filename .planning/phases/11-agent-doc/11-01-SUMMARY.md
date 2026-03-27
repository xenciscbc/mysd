---
phase: 11-agent-doc
plan: "01"
subsystem: config
tags: [go, cobra, viper, yaml, config, executor]

# Dependency graph
requires:
  - phase: 07-new-binary-commands-scanner-refactor
    provides: ProjectConfig struct and config.Load/Save patterns established
  - phase: 06-executor-wave-grouping-worktree-engine
    provides: ExecutionContext struct and BuildContextFromParts function
provides:
  - DocsToUpdate []string field in ProjectConfig (internal/config/defaults.go)
  - DocsToUpdate []string field in ExecutionContext with json:omitempty
  - cfg.DocsToUpdate wired through BuildContextFromParts
  - mysd docs list/add/remove subcommands (cmd/docs.go)
affects:
  - 11-02 (SKILL.md archive layer reads docs_to_update from context-only JSON)
  - 11-03 (ff/ffe pipeline uses docs_to_update)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Additive-only struct extension: new fields appended at END of ProjectConfig and ExecutionContext per D-11/D-12"
    - "viper-based config read-modify-write pattern preserving existing fields (same as cmd/model.go)"
    - "docs command follows note.go cobra subcommand pattern (list as default RunE)"

key-files:
  created:
    - cmd/docs.go
    - cmd/docs_test.go
  modified:
    - internal/config/defaults.go
    - internal/config/config.go
    - internal/executor/context.go
    - internal/executor/context_test.go

key-decisions:
  - "DocsToUpdate defaults to nil (not empty slice) — convention over config per D-14; omitempty ensures JSON omits the field when unconfigured"
  - "viper preserves other config fields when writing docs_to_update — same pattern as runModelSet in cmd/model.go"
  - "Duplicate detection in docs add is string equality — exact path match required"

patterns-established:
  - "Config extension pattern: append field to ProjectConfig struct, add nil default in Defaults(), add v.SetDefault() in Load()"
  - "ExecutionContext extension pattern: append field with omitempty, wire in BuildContextFromParts after HasParallelOpp line"

requirements-completed:
  - D-10
  - D-12
  - D-19

# Metrics
duration: 3min
completed: 2026-03-27
---

# Phase 11 Plan 01: Binary Go code — DocsToUpdate config + context + mysd docs command Summary

**DocsToUpdate []string field added to ProjectConfig and ExecutionContext with full viper wiring, plus mysd docs list/add/remove subcommands for managing docs_to_update in mysd.yaml**

## Performance

- **Duration:** ~3 min
- **Started:** 2026-03-27T02:06:04Z
- **Completed:** 2026-03-27T02:09:19Z
- **Tasks:** 2
- **Files modified:** 6

## Accomplishments
- DocsToUpdate field flows from ProjectConfig through BuildContextFromParts to ExecutionContext JSON (omitempty when nil)
- viper default registered in config.Load() for docs_to_update key
- mysd docs command with list/add/remove subcommands, following cmd/model.go viper write pattern
- 10 new tests covering DocsToUpdate field behavior and docs command operations

## Task Commits

Each task was committed atomically:

1. **TDD RED: DocsToUpdate test stubs** - `c6778af` (test)
2. **Task 1: Add DocsToUpdate to ProjectConfig and ExecutionContext** - `b85481e` (feat)
3. **Task 2: Create mysd docs command (list/add/remove)** - `5f0864e` (feat)

**Plan metadata:** (docs commit follows)

_Note: TDD task has two commits (test RED → feat GREEN)_

## Files Created/Modified
- `internal/config/defaults.go` - Added DocsToUpdate []string field at END of ProjectConfig, nil default in Defaults()
- `internal/config/config.go` - Added v.SetDefault("docs_to_update") after auto_mode line
- `internal/executor/context.go` - Added DocsToUpdate []string field to ExecutionContext, wire in BuildContextFromParts
- `internal/executor/context_test.go` - Added 4 DocsToUpdate tests (field passthrough, nil default, JSON presence, JSON omitempty)
- `cmd/docs.go` - New file: docsCmd, docsAddCmd, docsRemoveCmd with viper-based config write
- `cmd/docs_test.go` - New file: 6 tests covering list/add/remove/duplicate/notfound cases

## Decisions Made
- DocsToUpdate defaults to nil (not empty slice) so omitempty suppresses the JSON key when unconfigured — same convention as existing optional fields
- viper read-modify-write in writeDocsToUpdate preserves all other .claude/mysd.yaml fields, matching runModelSet pattern
- docs add duplicate detection uses exact string equality (path must match exactly as provided)

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## Known Stubs

None — DocsToUpdate data is wired end-to-end from config load through JSON output. The SKILL.md layer (plans 11-02/11-03) will consume the field.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- DocsToUpdate field available in `mysd execute --context-only` JSON output when configured
- `mysd docs add/remove` enables users to manage the list in mysd.yaml
- Ready for Phase 11-02 (SKILL.md archive layer integration)

---
*Phase: 11-agent-doc*
*Completed: 2026-03-27*
