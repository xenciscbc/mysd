---
phase: 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync
plan: "03"
subsystem: cli
tags: [cobra, update, self-update, plugin-sync, json-output, skill-md]

requires:
  - phase: 10-01
    provides: version.go CheckLatestVersion, IsUpdateAvailable, ReleaseInfo, ApplyUpdate
  - phase: 10-02
    provides: manifest.go LoadManifest, DiffManifests, SaveManifest, GenerateManifest; pluginsync.go SyncPlugins

provides:
  - cmd/update.go Cobra update command with --check and --force flags, JSON output
  - cmd/update_test.go tests for update command flags and JSON struct
  - .claude/commands/mysd-update.md SKILL.md thin wrapper for /mysd:update
  - plugin/commands/mysd-update.md distribution copy (identical to .claude/ copy)

affects:
  - phase-11 (any future phases that invoke mysd update)
  - plugin-distribution (plugin/commands now includes mysd-update.md)

tech-stack:
  added: []
  patterns:
    - "update command wires internal/update/ library functions into Cobra RunE"
    - "plugin sync walks up cwd to locate .claude/ directory (findClaudeDir)"
    - "network failure non-fatal: error recorded in JSON output, plugin sync continues regardless"
    - "SKILL.md thin wrapper pattern: argument parsing + binary invocation + JSON parsing only"

key-files:
  created:
    - cmd/update.go
    - cmd/update_test.go
    - .claude/commands/mysd-update.md
    - plugin/commands/mysd-update.md
  modified: []

key-decisions:
  - "findClaudeDir walks up from cwd to filesystem root to locate .claude/ — supports running from any subdirectory"
  - "Plugin sync source is plugin/ (sibling of .claude/) not a release archive — uses current dev copy for local sync"
  - "Binary update only runs when --force AND update_available — SKILL.md layer handles user confirmation flow"

patterns-established:
  - "Update JSON output: current_version, latest_version, update_available, check_only, force, binary_updated, plugin_sync, error"
  - "SyncOutput mirrors SyncResult from internal/update/pluginsync.go for JSON serialization"

requirements-completed: [UPD-05, UPD-06]

duration: 3min
completed: 2026-03-26
---

# Phase 10 Plan 03: Update Command & SKILL.md Wrapper Summary

**`mysd update [--check] [--force]` Cobra command with JSON output wired to internal/update/ library, plus `/mysd:update` SKILL.md thin wrapper with interactive confirmation flow**

## Performance

- **Duration:** 3 min
- **Started:** 2026-03-26T09:01:00Z
- **Completed:** 2026-03-26T09:03:33Z
- **Tasks:** 2
- **Files modified:** 4

## Accomplishments

- Cobra `update` command registered on rootCmd with `--check` (version check only) and `--force` (skip confirmation) flags
- JSON output structure covers version info, update status, binary update result, and plugin sync summary
- Plugin sync integrated: scans plugin/ directory, diffs against old manifest, syncs .claude/ commands and agents
- SKILL.md wrapper provides interactive 5-step flow (parse args → check → confirm → execute → verify)

## Task Commits

1. **Task 1: Cobra update command with JSON output** - `cc7a10a` (feat)
2. **Task 2: SKILL.md thin wrapper and plugin distribution copy** - `5a3bbf5` (feat)

**Plan metadata:** (added in final commit)

## Files Created/Modified

- `cmd/update.go` - Cobra update command, UpdateOutput/SyncOutput structs, runUpdate, runPluginSync, findClaudeDir, printJSON
- `cmd/update_test.go` - 7 tests covering flag registration, JSON struct marshaling, omitempty behavior
- `.claude/commands/mysd-update.md` - SKILL.md thin wrapper with argument-hint, 5-step interactive flow
- `plugin/commands/mysd-update.md` - Distribution copy, identical content to .claude/ copy

## Decisions Made

- `findClaudeDir` walks upward from cwd to filesystem root — supports running the command from any project subdirectory
- Plugin sync source uses `plugin/` directory (sibling of `.claude/`) rather than extracting from a release archive — for local/dev sync this is correct; in production the binary update would extract plugin files from the downloaded archive
- Binary update only applies when both `--force` is set AND an update is available — SKILL.md handles the confirmation prompt before re-running with `--force`

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Phase 10 complete: all 3 plans (version check, manifest/plugin sync, update command + SKILL.md) are done
- `mysd update --check` outputs JSON with current/latest version comparison
- `mysd update --force` applies binary update and syncs plugin files
- `/mysd:update` SKILL.md provides interactive wrapper for end users

## Self-Check: PASSED

- [x] cmd/update.go exists
- [x] cmd/update_test.go exists
- [x] .claude/commands/mysd-update.md exists
- [x] plugin/commands/mysd-update.md exists
- [x] go build ./... exits 0
- [x] diff .claude/commands/mysd-update.md plugin/commands/mysd-update.md produces no output
- [x] commits cc7a10a and 5a3bbf5 exist

---
*Phase: 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync*
*Completed: 2026-03-26*
