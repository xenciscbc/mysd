---
phase: 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync
plan: 02
subsystem: update
tags: [plugin-sync, manifest, goreleaser, json, tdd]

requires:
  - phase: 10-01
    provides: version.go ReleaseInfo/Asset types, selfupdate.go infrastructure already committed

provides:
  - PluginManifest struct with Version/Commands/Agents fields (internal/update/manifest.go)
  - LoadManifest/SaveManifest JSON CRUD with nil-safe missing-file handling
  - DiffManifests three-way comparison with backward-compat nil old manifest
  - GenerateManifest scans plugin directories for .md files (excludes CLAUDE.md)
  - SyncPlugins copies/deletes files based on ManifestDiff with non-fatal delete errors
  - SyncResult struct with Added/Updated/Deleted/Errors fields
  - .goreleaser.yaml updated to include plugin/commands/*.md, plugin/agents/*.md, plugin-manifest.json in release archives

affects:
  - 10-03 (cmd layer will call SyncPlugins after downloading release archive)
  - future GoReleaser releases (plugin files now bundled in every platform archive)

tech-stack:
  added: []
  patterns:
    - "Nil-safe manifest loading: LoadManifest returns (nil, nil) for missing file — same convention-over-config pattern as deferred.go"
    - "DiffManifests nil-old backward compat: pre-v1.1 installations without manifest never get delete operations (D-17)"
    - "Non-fatal delete errors: SyncPlugins appends to Errors slice and continues — sync is best-effort for removals"

key-files:
  created:
    - internal/update/manifest.go
    - internal/update/manifest_test.go
    - internal/update/pluginsync.go
    - internal/update/pluginsync_test.go
  modified:
    - .goreleaser.yaml

key-decisions:
  - "LoadManifest returns (nil, nil) for missing file — matches deferred.go pattern, represents pre-v1.1 installation"
  - "DiffManifests with nil old: all new files are add, zero deletes — backward compat per D-17"
  - "SyncPlugins delete errors are non-fatal: append to Errors slice, continue — sync does not fail on missing-file delete"
  - "GoReleaser files: includes plugin/commands/*.md, plugin/agents/*.md, plugin-manifest.json in every platform archive"

patterns-established:
  - "Manifest JSON I/O: os.ReadFile + json.Unmarshal + os.IsNotExist nil return (deferred.go pattern)"
  - "Three-way diff via set intersection: build oldSet and newSet from string slices, classify into add/update/delete"

requirements-completed: [UPD-03, UPD-07]

duration: 6min
completed: 2026-03-26
---

# Phase 10 Plan 02: Plugin Manifest and Sync Summary

**Plugin manifest CRUD with three-way diff and file sync executor, GoReleaser updated to bundle plugin files in every release archive**

## Performance

- **Duration:** 6 min
- **Started:** 2026-03-26T08:51:33Z
- **Completed:** 2026-03-26T08:57:18Z
- **Tasks:** 2
- **Files modified:** 5

## Accomplishments

- PluginManifest struct and JSON Load/Save with nil-safe missing-file convention (D-17 backward compat)
- DiffManifests three-way comparison: add/update/delete with nil old manifest producing zero deletes
- GenerateManifest scans plugin/commands/ and plugin/agents/, excludes CLAUDE.md
- SyncPlugins copies and deletes files based on ManifestDiff, non-fatal delete errors
- .goreleaser.yaml `archives.files` now includes plugin .md files and plugin-manifest.json

## Task Commits

Each task was committed atomically:

1. **Task 1: Plugin manifest and diff logic** - `70554b3` (feat)
2. **Task 2: Plugin sync executor and GoReleaser config** — committed as part of `e0d1ba5` by 10-01 parallel agent (see Deviations)

**Plan metadata:** (see final commit)

_Note: TDD tasks have RED (failing test) → GREEN (implementation) cycle verified._

## Files Created/Modified

- `internal/update/manifest.go` — PluginManifest struct, LoadManifest, SaveManifest, DiffManifests, GenerateManifest
- `internal/update/manifest_test.go` — Table-driven tests for all DiffManifests edge cases + nil backward compat
- `internal/update/pluginsync.go` — SyncResult struct, SyncPlugins function with non-fatal delete errors
- `internal/update/pluginsync_test.go` — Tests for add/update/delete commands and agents, combined scenario
- `.goreleaser.yaml` — Added `files:` section to `archives` with plugin/commands/*.md, plugin/agents/*.md, plugin-manifest.json

## Decisions Made

- LoadManifest returns (nil, nil) for missing file — same convention-over-config pattern as deferred.go
- DiffManifests with nil old manifest: all new files are "add", zero "delete" operations — D-17 backward compat for pre-v1.1 installations
- SyncPlugins delete errors are non-fatal: appended to Errors slice, sync continues — prevents hard failure when previously-tracked file was manually removed
- GoReleaser archives.files: plugin/* and manifest bundled in binary archive so update command only needs one download

## Deviations from Plan

### Parallel Agent Overlap

**1. [Rule 3 - Blocking] Parallel agent 10-01 pre-committed pluginsync.go, pluginsync_test.go, .goreleaser.yaml**
- **Found during:** Task 2 post-implementation
- **Issue:** When attempting to commit Task 2 files, git status showed they had already been committed in `e0d1ba5` by the 10-01 parallel agent. The 10-01 agent committed files beyond its own scope (selfupdate + version) and included 10-02 files.
- **Resolution:** Verified all 10-02 acceptance criteria were met by the existing committed code. Confirmed pluginsync.go, pluginsync_test.go, and .goreleaser.yaml match 10-02 specifications exactly. Tests pass. No action required — plan outcome delivered.
- **Impact:** No functional impact. All success criteria met. Task 2 implementation correct.

---

**Total deviations:** 1 (parallel agent overlap — no corrective action needed)
**Impact on plan:** All plan requirements delivered. Parallel execution produced identical implementation to what this agent would have written.

## Issues Encountered

- Parallel 10-01 agent created selfupdate_test.go with references to undefined functions (VerifyChecksum, ExtractBinary, etc.), blocking package compilation during Task 2 verification. Resolved when selfupdate.go and selfupdate_unix/windows.go (also pre-created by 10-01 agent) were found to already define these functions, allowing `go build ./...` to succeed.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- plugin manifest and sync infrastructure complete
- SyncPlugins ready for integration in cmd/update.go (10-03)
- .goreleaser.yaml configured — next release will bundle plugin files automatically
- GoReleaser will need actual plugin-manifest.json generated at release time (future CI step)

## Self-Check: PASSED

- manifest.go: FOUND
- manifest_test.go: FOUND
- pluginsync.go: FOUND
- pluginsync_test.go: FOUND
- .goreleaser.yaml: FOUND
- 10-02-SUMMARY.md: FOUND
- Commit 70554b3 (Task 1): FOUND
- Commit e0d1ba5 (Task 2, via parallel agent): FOUND

---
*Phase: 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync*
*Completed: 2026-03-26*
