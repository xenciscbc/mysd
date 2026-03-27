---
spec-version: "1.0"
capability: Self-Update & Plugin Sync
delta: ADDED
status: done
---

## Requirement: Version Check

The `mysd update` command MUST check for new releases on GitHub.

The `--check` flag MUST only report available updates without applying them.

## Requirement: Binary Self-Update

The update mechanism MUST download and replace the running binary with the new version.

Platform-specific logic MUST handle:
- Unix: Direct binary replacement
- Windows: Rename-then-replace strategy (cannot overwrite running executable)

## Requirement: Plugin File Sync

After binary update, the system MUST sync plugin files (commands and agents) to `.claude/`.

Sync MUST use `plugin-manifest.json` to compute the diff between old and new versions:
- Files in new but not in old → add
- Files in both → overwrite
- Files in old but not in new → delete

If no previous manifest exists (pre-v1.1 installation), all new files MUST be added and no files MUST be deleted.

## Requirement: Plugin Manifest

`GenerateManifest()` MUST scan `plugin/commands/*.md` and `plugin/agents/*.md` to build the manifest.

`CLAUDE.md` files MUST be excluded from the manifest.

`SaveManifest()` MUST write the manifest as indented JSON to `.claude/plugin-manifest.json`.

The `--force` flag MUST re-sync all plugin files regardless of manifest state.

### Scenario: First Update from Pre-v1.1

WHEN no plugin-manifest.json exists
THEN DiffManifests treats all new files as "add"
AND no existing files are deleted

### Scenario: Normal Update

WHEN plugin-manifest.json exists with version 1.0.2
AND new release is 1.1.0
THEN DiffManifests computes add/update/delete operations
AND plugin files are synced accordingly
AND manifest is updated to 1.1.0

## Covered Packages

- `cmd/update.go`
- `internal/update/manifest.go` — manifest generation, diff, save/load
- `internal/update/selfupdate.go` — binary download and replacement
- `internal/update/pluginsync.go` — plugin file synchronization
- `internal/update/selfupdate_unix.go`, `selfupdate_windows.go` — platform-specific logic
