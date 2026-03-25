---
phase: 05-schema-foundation-plan-checker
plan: "01"
subsystem: spec
tags: [go, yaml, taskentry, taskitem, model-profile, openspec, config]

requires:
  - phase: 04-plugin-layer-distribution
    provides: v1.0 codebase with spec/schema.go, executor/context.go, config package

provides:
  - TaskEntry struct extended with Depends, Files, Satisfies, Skills (omitempty YAML)
  - TaskItem struct extended with matching JSON fields (omitempty)
  - BuildContextFromParts copies new fields from TaskEntry to TaskItem
  - DefaultModelMap extended with 4 new agent roles (researcher, advisor, proposal-writer, plan-checker)
  - ProjectConfig extended with WorktreeDir (".worktrees") and AutoMode (false) fields
  - OpenSpecConfig struct with WriteOpenSpecConfig/ReadOpenSpecConfig functions
  - Backward compat: old tasks.md round-trips without new field artifacts

affects:
  - 05-02: plan-checker uses TaskEntry.Satisfies for requirement coverage check
  - phase-06: worktree execution uses ProjectConfig.WorktreeDir
  - phase-07: interactive commands use ProjectConfig.AutoMode
  - phase-08: new agent definitions use 4 new roles from DefaultModelMap
  - phase-09: researcher/advisor agents resolve via updated model profile

tech-stack:
  added: []
  patterns:
    - "Additive struct extension: new fields at END of struct, always omitempty to preserve backward compat"
    - "TDD convention: test file before implementation, RED confirms compile failure, GREEN confirms all pass"
    - "Viper binding: every new ProjectConfig field needs both yaml+mapstructure tags AND v.SetDefault() call"
    - "Convention-over-config read: os.IsNotExist returns zero-value, not error"

key-files:
  created:
    - internal/spec/openspec_config.go
    - internal/spec/openspec_config_test.go
  modified:
    - internal/spec/schema.go
    - internal/spec/schema_test.go
    - internal/spec/updater_test.go
    - internal/executor/context.go
    - internal/executor/context_test.go
    - internal/config/config.go
    - internal/config/defaults.go
    - internal/config/config_test.go

key-decisions:
  - "New fields appended at END of structs to preserve field order for stable YAML output (D-11, D-12)"
  - "Budget profile new roles (researcher/advisor/proposal-writer/plan-checker) all use sonnet-4-5 not haiku (D-06)"
  - "ReadOpenSpecConfig returns zero-value OpenSpecConfig (not error) when file absent — convention-over-config"
  - "nil slices passed directly to TaskItem fields — no make([]int,0) to preserve JSON omitempty behavior"

patterns-established:
  - "Additive-only struct extension: append fields at end, use omitempty — zero migration required"
  - "Convention-over-config file reads: os.IsNotExist check returns zero-value and nil error"

requirements-completed:
  - FSCHEMA-01
  - FSCHEMA-02
  - FSCHEMA-03
  - FSCHEMA-04
  - FSCHEMA-07
  - FMODEL-01
  - FMODEL-02
  - FMODEL-03

duration: 18min
completed: 2026-03-25
---

# Phase 05 Plan 01: Schema Foundation Summary

**TaskEntry/TaskItem extended with 4 dependency fields, DefaultModelMap gains 4 new agent roles, and OpenSpecConfig reader/writer created — zero migration required for existing tasks.md files**

## Performance

- **Duration:** ~18 min
- **Started:** 2026-03-25T07:10:00Z
- **Completed:** 2026-03-25T07:28:00Z
- **Tasks:** 3 (all TDD)
- **Files modified:** 9 (7 modified, 2 created)

## Accomplishments

- Extended TaskEntry with Depends/Files/Satisfies/Skills fields — old tasks.md files round-trip without adding empty field keys
- Extended DefaultModelMap with 4 new agent roles across all 3 profiles; budget profile new roles use sonnet-4-5 per D-06
- Added WorktreeDir (".worktrees") and AutoMode (false) to ProjectConfig with proper viper binding
- Created WriteOpenSpecConfig/ReadOpenSpecConfig with convention-over-config absent-file handling

## Task Commits

Each task was committed atomically:

1. **Task 1: Extend TaskEntry and TaskItem structs with 4 new fields** - `3eb8942` (feat)
2. **Task 2: Extend model profile and ProjectConfig for new agent roles** - `7168c85` (feat)
3. **Task 3: Create openspec/config.yaml reader and writer** - `570fe7c` (feat)

## Files Created/Modified

- `internal/spec/schema.go` - TaskEntry: added Depends, Files, Satisfies, Skills fields with omitempty YAML tags
- `internal/spec/schema_test.go` - Added TestTaskEntryNewFields_YAMLRoundTrip, TestTaskEntryNewFields_OmitEmpty
- `internal/spec/updater_test.go` - Added TestParseTasksV2_BackwardCompat_NoNewFields
- `internal/spec/openspec_config.go` - New: OpenSpecConfig struct, WriteOpenSpecConfig, ReadOpenSpecConfig
- `internal/spec/openspec_config_test.go` - New: 6 tests for write/read/round-trip/not-exist/malformed/BCP47
- `internal/executor/context.go` - TaskItem: added 4 matching JSON fields; BuildContextFromParts copies them
- `internal/executor/context_test.go` - Added TestBuildContextFromParts_NewFields, TestTaskItemJSON_OmitEmpty
- `internal/config/config.go` - DefaultModelMap: 4 new roles in all 3 tiers; Load(): 2 new SetDefault calls
- `internal/config/defaults.go` - ProjectConfig: WorktreeDir, AutoMode fields; Defaults() returns correct values
- `internal/config/config_test.go` - Added TestResolveModel_NewRoles (12 combos), TestDefaults_NewFields, TestLoad_NewFields

## Decisions Made

- New fields appended at END of structs (D-11, D-12) for stable YAML output field order
- Budget profile new roles use sonnet-4-5 not haiku — consistent with D-06 rationale that new subagent roles require quality model
- nil slices for empty new fields (not `make([]int,0)`) — preserves omitempty behavior in both YAML and JSON

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None — all tests passed on first GREEN implementation.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

- Plan 02 (plan-checker) can now use TaskEntry.Satisfies for requirement coverage checks
- Phase 06 worktree execution can use ProjectConfig.WorktreeDir
- Phase 07 interactive commands can use ProjectConfig.AutoMode
- Phase 08 new agent definitions can resolve model via 4 new roles in DefaultModelMap

---
*Phase: 05-schema-foundation-plan-checker*
*Completed: 2026-03-25*
