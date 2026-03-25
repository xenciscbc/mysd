---
phase: 06-executor-wave-grouping-worktree-engine
plan: "02"
subsystem: worktree
tags: [git-worktree, worktree-manager, disk-space, windows-longpaths, cobra, cli]

requires:
  - phase: 01-foundation
    provides: "config.ProjectConfig with WorktreeDir field, cobra rootCmd pattern"
  - phase: 02-execution-engine
    provides: "state.LoadState for ChangeName, spec.DetectSpecDir pattern"

provides:
  - "WorktreeManager.Create: creates git worktree at .worktrees/T{id}/ with branch mysd/{change}/T{id}-{slug}"
  - "WorktreeManager.Remove: removes worktree with --force + branch delete + prune"
  - "WorktreeManager.CheckDiskSpace: 500MB guard before creating worktree"
  - "Windows core.longpaths auto-set via setLongPaths() on Create"
  - "Platform-specific disk space: syscall.Statfs (Unix) / GetDiskFreeSpaceExW (Windows)"
  - "mysd worktree create/remove subcommands outputting JSON for SKILL.md consumption"

affects:
  - 06-03-waveexecutor
  - 06-04-skill-md-integration
  - plugin/commands

tech-stack:
  added: []
  patterns:
    - "Build-tag separation: diskspace_unix.go (!windows) / diskspace_windows.go (windows)"
    - "platform-specific syscall via build tags, no golang.org/x/sys needed for Windows (kernel32.dll via syscall.NewLazyDLL)"
    - "ToSlug exported for independent testability"
    - "JSON output to stdout for CLI-to-SKILL.md communication pattern"
    - "Non-fatal branch delete warning to stderr (worktree already removed = success)"

key-files:
  created:
    - internal/worktree/worktree.go
    - internal/worktree/worktree_test.go
    - internal/worktree/diskspace_unix.go
    - internal/worktree/diskspace_windows.go
    - cmd/worktree.go
  modified: []

key-decisions:
  - "Used kernel32.dll GetDiskFreeSpaceExW via syscall.NewLazyDLL instead of golang.org/x/sys/windows — avoids new dependency, kernel32.dll always available on Windows"
  - "ToSlug exported (not unexported toSlug) to enable independent unit testing from worktree_test package"
  - "Branch delete failure non-fatal: worktree already removed is the success state; branch cleanup is best-effort"

patterns-established:
  - "Pattern 1: Build-tag platform split — diskspace_unix.go vs diskspace_windows.go, package private function getAvailableBytes"
  - "Pattern 2: worktree create outputs JSON {path, branch} to stdout — all CLI-to-SKILL.md data via JSON stdout"

requirements-completed: [FEXEC-04, FEXEC-05, FEXEC-08, FEXEC-10, FEXEC-11]

duration: 3min
completed: "2026-03-25"
---

# Phase 06 Plan 02: Worktree Lifecycle Management Summary

**Git worktree lifecycle package with WorktreeManager (Create/Remove/CheckDiskSpace), cross-platform disk space via build tags, Windows core.longpaths auto-config, and mysd worktree create/remove CLI subcommands outputting JSON**

## Performance

- **Duration:** ~3 min
- **Started:** 2026-03-25T08:09:57Z
- **Completed:** 2026-03-25T08:12:56Z
- **Tasks:** 2 (TDD: 3 commits for Task 1)
- **Files modified:** 5 created

## Accomplishments

- WorktreeManager with Create (path + branch naming), Remove (--force + prune), CheckDiskSpace (500MB guard)
- Cross-platform disk space check: `syscall.Statfs` on Unix, `GetDiskFreeSpaceExW` via `kernel32.dll` on Windows
- Windows core.longpaths automatically set on Create (FEXEC-11)
- `mysd worktree create <id> <name>` and `mysd worktree remove <id> <branch>` CLI subcommands with JSON output
- All 8 worktree tests pass (including Windows-specific longpaths test running natively on Windows CI)

## Task Commits

Each task was committed atomically:

1. **Task 1 RED: Add failing tests for WorktreeManager** - `5aab498` (test)
2. **Task 1 GREEN: Implement WorktreeManager with Create, Remove, CheckDiskSpace** - `61ce6ea` (feat)
3. **Task 2: Add mysd worktree create/remove CLI subcommands** - `a24a820` (feat)

## Files Created/Modified

- `internal/worktree/worktree.go` - WorktreeManager struct, Create/Remove/CheckDiskSpace methods, ToSlug helper
- `internal/worktree/worktree_test.go` - 8 integration tests using t.TempDir + git init
- `internal/worktree/diskspace_unix.go` - `//go:build !windows` disk space via syscall.Statfs
- `internal/worktree/diskspace_windows.go` - `//go:build windows` disk space via GetDiskFreeSpaceExW
- `cmd/worktree.go` - cobra worktreeCmd + worktreeCreateCmd + worktreeRemoveCmd, JSON output

## Decisions Made

- Used `kernel32.dll` via `syscall.NewLazyDLL` instead of `golang.org/x/sys/windows.GetDiskFreeSpaceEx` to avoid adding a new direct dependency. `kernel32.dll` is always present on Windows.
- `ToSlug` is exported (capitalized) to allow `worktree_test` package (black-box test) to call it directly for independent slug logic testing.
- Branch delete failure in `Remove()` is non-fatal (logs to stderr). The critical operation is `git worktree remove --force`; branch cleanup is best-effort.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None. The tests ran cleanly on Windows (the development environment), including `TestCreate_WindowsLongPaths` which verified `git config core.longpaths` returns "true" after Create.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- `internal/worktree` package ready for import by Phase 06-03 WaveExecutor
- `mysd worktree create/remove` CLI available for SKILL.md calling convention
- Concern: Windows MAX_PATH for deeply nested Go packages inside worktrees still needs empirical CI validation (existing blocker in STATE.md)

---
*Phase: 06-executor-wave-grouping-worktree-engine*
*Completed: 2026-03-25*
