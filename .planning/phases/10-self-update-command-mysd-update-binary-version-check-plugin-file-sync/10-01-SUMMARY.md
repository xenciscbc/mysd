---
phase: 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync
plan: "01"
subsystem: internal/update
tags: [go, self-update, semver, checksum, platform-specific, tdd]
dependency_graph:
  requires: []
  provides: [internal/update package — version check, semver comparison, binary self-update]
  affects: [cmd/update — Plan 03 will use these functions to implement the CLI command]
tech_stack:
  added: []
  patterns:
    - TDD (RED-GREEN) for all public functions
    - httptest.NewServer for HTTP mocking without external dependencies
    - Build tags for platform-specific implementations (windows / !windows)
    - rename-then-replace strategy for Windows running-binary safety
key_files:
  created:
    - internal/update/version.go
    - internal/update/version_test.go
    - internal/update/selfupdate.go
    - internal/update/selfupdate_windows.go
    - internal/update/selfupdate_unix.go
    - internal/update/selfupdate_test.go
  modified: []
decisions:
  - "CheckLatestVersionWithBase added as testable variant of CheckLatestVersion to enable httptest mocking without changing public API"
  - "AssetNameForPlatform uses runtime.GOOS/GOARCH so tests must also use runtime values to match platform"
  - "ApplyUpdate calls replaceExecutable internally; rollback is a public function for cmd layer to call on post-update failure"
metrics:
  duration: "297s"
  completed: "2026-03-26"
  tasks_completed: 2
  files_created: 6
---

# Phase 10 Plan 01: Core Update Library — Version Check and Binary Self-Update Summary

**One-liner:** GitHub API version check with semver comparison and platform-aware binary self-update using SHA256 checksum verification and automatic rollback.

## What Was Built

The `internal/update` package provides the foundational update mechanics for `mysd update` command (Plan 03):

**version.go** — GitHub API client and semver utilities:
- `ParseSemver`: parses vX.Y.Z or X.Y.Z, rejects "dev" and malformed strings
- `IsUpdateAvailable`: treats "dev" as always outdated (D-02), compares semver otherwise
- `CheckLatestVersion`: queries GitHub Releases API with proper Accept header and 15s timeout
- `CheckLatestVersionWithBase`: testable variant accepting custom API base URL
- `AssetNameForPlatform`: generates GoReleaser archive names (`.tar.gz` / `.zip` for Windows)
- `FindAssetURL` / `FindChecksumURL`: locate specific assets in ReleaseInfo

**selfupdate.go** — Binary download and replacement:
- `DownloadFile`: downloads URL to temp file with context support
- `VerifyChecksum`: SHA256 comparison, returns error with "checksum mismatch"
- `ParseChecksumFile`: parses GoReleaser double-space separated checksums.txt
- `ExtractBinary`: extracts mysd binary from tar.gz or zip using stdlib only
- `ApplyUpdate`: full pipeline — download, verify, extract, replace (with rollback on failure)
- `Rollback`: restores .old to original path

**selfupdate_windows.go** (build tag `windows`):
- `replaceExecutable`: rename-then-replace pattern (safe for running binaries on Windows)

**selfupdate_unix.go** (build tag `!windows`):
- `replaceExecutable`: rename-then-copy with 0755 permissions, removes .old

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Version check — GitHub API client and semver comparison | 2c7d795 | version.go, version_test.go |
| 2 | Binary self-update — download, checksum verify, platform replace, rollback | e0d1ba5 | selfupdate.go, selfupdate_windows.go, selfupdate_unix.go, selfupdate_test.go |

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed TestApplyUpdateChecksumMismatch platform mismatch**
- **Found during:** Task 2 test execution
- **Issue:** Test hardcoded `linux_amd64.tar.gz` asset name, but `ApplyUpdate` uses `runtime.GOOS/GOARCH` to determine asset name. On Windows, the function looked for `windows_amd64.zip` but the test's mock release only had the Linux asset.
- **Fix:** Changed test to call `AssetNameForPlatform(runtime.GOOS, runtime.GOARCH, "1.0.0")` so both ApplyUpdate and the mock release use the same platform-appropriate asset name.
- **Files modified:** internal/update/selfupdate_test.go
- **Commit:** e0d1ba5

## Known Stubs

None — all functions are fully implemented with no placeholder values or hardcoded empty returns.

## Self-Check: PASSED

- internal/update/version.go: FOUND
- internal/update/version_test.go: FOUND
- internal/update/selfupdate.go: FOUND
- internal/update/selfupdate_windows.go: FOUND
- internal/update/selfupdate_unix.go: FOUND
- internal/update/selfupdate_test.go: FOUND
- Commit 2c7d795: FOUND
- Commit e0d1ba5: FOUND
- go test ./internal/update/: PASSED (all tests green)
- go build ./...: PASSED
- go vet ./internal/update/: PASSED
