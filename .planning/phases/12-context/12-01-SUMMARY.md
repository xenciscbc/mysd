---
phase: 12-context
plan: 01
subsystem: binary
tags: [statusline, config, embed, init, settings]
dependency_graph:
  requires: []
  provides: [statusline-subcommand, hook-embed, init-hook-install, settings-merge]
  affects: [cmd/init_cmd.go, internal/config/defaults.go]
tech_stack:
  added: [go:embed]
  patterns: [viper-read-modify-write, json-merge, tdd-red-green]
key_files:
  created:
    - plugin/hooks/mysd-statusline.js
    - cmd/hooks/mysd-statusline.js
    - cmd/hooks_embed.go
    - cmd/statusline.go
    - cmd/statusline_test.go
  modified:
    - internal/config/defaults.go
    - cmd/init_cmd.go
    - cmd/init_cmd_test.go
decisions:
  - "Go embed cannot use ../ path prefix; workaround: copy JS to cmd/hooks/ subdirectory for embed (deviation Rule 3)"
  - "runStatuslineInDir extracted for testability — tests pass baseDir instead of using cwd"
  - "MkdirAll before v.SafeWriteConfig required to avoid 'missing configuration for configPath' on fresh TempDir"
metrics:
  duration: 258s
  completed: "2026-03-27"
  tasks_completed: 2
  files_changed: 8
---

# Phase 12 Plan 01: Go Binary Statusline Infrastructure Summary

Go binary extensions for statusline: ProjectConfig field, embed hook content, statusline subcommand, init hook install, settings.json merge.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | ProjectConfig extension + statusline subcommand + embed + tests | bf9a3e1 | internal/config/defaults.go, cmd/hooks_embed.go, cmd/statusline.go, cmd/statusline_test.go, cmd/hooks/mysd-statusline.js, plugin/hooks/mysd-statusline.js |
| 2 | Init hook install + settings.json merge + tests | 59ad258 | cmd/init_cmd.go, cmd/init_cmd_test.go |

## What Was Built

**Task 1:**
- `ProjectConfig.StatuslineEnabled *bool` field added at struct end (last field, omitempty, nil = enabled per D-12)
- `plugin/hooks/mysd-statusline.js` — the real implementation was already in place (from Plan 02 work done concurrently)
- `cmd/hooks/mysd-statusline.js` — copy in cmd/ for Go embed (see Deviations)
- `cmd/hooks_embed.go` — `//go:embed hooks/mysd-statusline.js` directive, exposes `statuslineHookBytes []byte`
- `cmd/statusline.go` — `mysd statusline [on|off]` subcommand with toggle support via viper read-modify-write
- 5 unit tests: TestRunStatuslineOn, TestRunStatuslineOff, TestRunStatuslineToggle, TestRunStatuslineToggleDefault, TestStatuslineOutputFormat

**Task 2:**
- `cmd/init_cmd.go` extended with two new steps after `.claude/` creation:
  - Installs `statuslineHookBytes` to `.claude/hooks/mysd-statusline.js`
  - Calls `writeSettingsStatusLine` to merge `statusLine` key into `.claude/settings.json`
- `writeSettingsStatusLine(claudeDir string) error` helper: reads existing JSON, sets `statusLine` key, writes back (D-07 merge pattern)
- 5 unit tests: TestInitStatuslineInstall, TestInitStatuslineInstallIdempotent, TestWriteSettingsStatusLine, TestWriteSettingsStatusLineMerge, TestWriteSettingsStatusLineNew

## Verification Results

```
go build ./...          PASS
go test ./cmd/... -run "TestRunStatusline|TestInitStatusline|TestWriteSettings"  10/10 PASS
go test ./internal/config/... -v -count=1  21/21 PASS
```

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Go embed cannot use ../ path prefix**
- **Found during:** Task 1, creating cmd/hooks_embed.go
- **Issue:** Plan specified `//go:embed ../plugin/hooks/mysd-statusline.js` but Go embed spec explicitly disallows `..` in embed paths. Build error: `invalid pattern syntax`
- **Fix:** Created `cmd/hooks/` subdirectory, copied JS to `cmd/hooks/mysd-statusline.js`, changed embed directive to `//go:embed hooks/mysd-statusline.js`
- **Files modified:** cmd/hooks_embed.go, cmd/hooks/mysd-statusline.js (new)
- **Note:** `plugin/hooks/mysd-statusline.js` remains the authoritative source; `cmd/hooks/` copy is for Go embed only

**2. [Rule 1 - Bug] MkdirAll needed before viper SafeWriteConfig**
- **Found during:** Task 1, TestRunStatuslineOn failing
- **Issue:** `v.SafeWriteConfig()` returns "missing configuration for configPath" when `.claude/` directory does not exist
- **Fix:** Added `os.MkdirAll(claudeDir, 0755)` in `runStatuslineInDir` before calling WriteConfig
- **Files modified:** cmd/statusline.go
- **Commit:** bf9a3e1 (same task commit)

## Known Stubs

None — `plugin/hooks/mysd-statusline.js` contains the full implementation (not a placeholder).

## Self-Check: PASSED

Files exist:
- internal/config/defaults.go: FOUND
- cmd/hooks_embed.go: FOUND
- cmd/statusline.go: FOUND
- cmd/statusline_test.go: FOUND
- cmd/init_cmd.go: FOUND
- cmd/init_cmd_test.go: FOUND
- plugin/hooks/mysd-statusline.js: FOUND
- cmd/hooks/mysd-statusline.js: FOUND

Commits exist:
- bf9a3e1: FOUND
- 59ad258: FOUND
