---
phase: 04-plugin-layer-distribution
plan: 01
subsystem: cli
tags: [go, cobra, scanner, filepath-walkdir, version-ldflags, json]

requires:
  - phase: 03-verification-feedback-loop
    provides: verifier package patterns (BuildVerificationContextFromParts, context-only JSON output pattern)
  - phase: 01-foundation
    provides: spec.DetectSpecDir for specs directory detection

provides:
  - internal/scanner package with BuildScanContext function
  - cmd/scan subcommand with --context-only and --exclude flags
  - SetVersion function in cmd/root.go for ldflags injection
  - version/commit/date variables in main.go for GoReleaser

affects:
  - 04-02 (GoReleaser config — consumes version variables and ldflags pattern)
  - 04-03 (scan agent SKILL.md — invokes mysd scan --context-only)

tech-stack:
  added: []
  patterns:
    - "WalkDir root-skip guard: path != root check prevents hidden-dir logic from aborting entire walk when root is '.'"
    - "context-only command pattern: runScanContextOnly(out io.Writer, root string, exclude []string) separates cobra from logic for test isolation"
    - "Version ldflags pattern: var version='dev' in main.go, cmd.SetVersion(version) before cmd.Execute()"

key-files:
  created:
    - internal/scanner/scanner.go
    - internal/scanner/scanner_test.go
    - cmd/scan.go
    - cmd/scan_test.go
  modified:
    - main.go
    - cmd/root.go

key-decisions:
  - "WalkDir root skip guard: path != root required because WalkDir calls root with name '.' which HasPrefix('.') — without guard entire walk aborts immediately"
  - "ExcludedDirs returns passed-in slice (empty slice not nil) for clean JSON output"
  - "PackageInfo.Name uses filepath.ToSlash for cross-platform forward-slash path convention"
  - "SetVersion added to cmd/root.go (not main.go) to keep rootCmd mutation inside cmd package"

requirements-completed: [WCMD-09, DIST-03]

duration: 22min
completed: 2026-03-24
---

# Phase 4 Plan 01: Scan Command and Version Wiring Summary

**codebase scanner (BuildScanContext) with WalkDir + exclusion + HasSpec detection, cobra scan subcommand with --context-only JSON output, and version ldflags wiring for GoReleaser**

## Performance

- **Duration:** 22 min
- **Started:** 2026-03-24T03:25:38Z
- **Completed:** 2026-03-24T03:47:00Z
- **Tasks:** 2
- **Files modified:** 6

## Accomplishments

- `internal/scanner.BuildScanContext` walks Go codebase with hidden-dir skip, exclude list, _test.go separation, and HasSpec detection via spec.DetectSpecDir
- `mysd scan --context-only` outputs structured JSON of all Go packages for AI agent spec generation
- `mysd --version` prints version string; main.go wired with version/commit/date variables for GoReleaser ldflags

## Task Commits

Each task was committed atomically:

1. **Task 1 RED: TestBuildScanContext (6 tests)** - `dbaed7e` (test)
2. **Task 1 GREEN: scanner.BuildScanContext implementation** - `7e91363` (feat)
3. **Task 2: scan command + version wiring + WalkDir bug fix** - `7978ae6` (feat)

## Files Created/Modified

- `internal/scanner/scanner.go` - BuildScanContext, ScanContext, PackageInfo types
- `internal/scanner/scanner_test.go` - 6 test cases covering all behavior scenarios
- `cmd/scan.go` - cobra scan subcommand with --context-only and --exclude flags
- `cmd/scan_test.go` - tests for runScanContextOnly and SetVersion
- `cmd/root.go` - added SetVersion() function
- `main.go` - added version/commit/date variables + cmd.SetVersion(version) call

## Decisions Made

- WalkDir root-skip guard (`path != root`) required: WalkDir calls root itself first with name `"."` which starts with `"."`, causing entire walk to abort without the guard. Found and fixed during Task 2 integration testing.
- SetVersion kept in cmd package (not main.go) to keep rootCmd mutation encapsulated within cmd package.
- PackageInfo.Name uses `filepath.ToSlash` for cross-platform forward-slash paths (convention over OS-native paths).

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed WalkDir root directory skip bug**
- **Found during:** Task 2 (integration testing of mysd scan --context-only)
- **Issue:** When root is `"."`, `filepath.WalkDir` calls the root entry with `d.Name()` returning `"."` — which `strings.HasPrefix(".", ".")` matches as true, causing `filepath.SkipDir` to abort the entire walk immediately. Output was always `"packages": []`.
- **Fix:** Added `path != root` guard to both hidden-dir and exclude-dir checks in WalkDir callback
- **Files modified:** `internal/scanner/scanner.go`
- **Verification:** All 6 scanner unit tests still pass; `mysd scan --context-only` now returns 10 packages from the worktree
- **Committed in:** `7978ae6` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (Rule 1 - Bug)
**Impact on plan:** Critical correctness fix — without it, scan command returned empty results for any real project scan from `"."` root.

## Issues Encountered

- WalkDir hidden-dir guard logic interacted with `"."` root name — resolved with `path != root` guard. All unit tests passed because tests use `t.TempDir()` absolute paths, not `"."`. Integration testing revealed the bug.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Scanner package ready for Plan 03 (scan agent SKILL.md) to invoke `mysd scan --context-only`
- Version ldflags pattern ready for Plan 02 (GoReleaser config) to inject real version values
- All existing tests still pass (`go test ./internal/scanner/... ./cmd/...`)

---
*Phase: 04-plugin-layer-distribution*
*Completed: 2026-03-24*
