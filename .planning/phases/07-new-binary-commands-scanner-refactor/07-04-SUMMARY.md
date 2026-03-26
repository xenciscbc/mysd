---
phase: 07-new-binary-commands-scanner-refactor
plan: "04"
subsystem: cmd
tags: [lang, config, atomic-write, dual-config, BCP47]
dependency_graph:
  requires:
    - internal/spec/openspec_config.go (ReadOpenSpecConfig, WriteOpenSpecConfig)
    - internal/config/config.go (Load)
    - internal/config/defaults.go (ProjectConfig.ResponseLanguage)
  provides:
    - cmd/lang.go (langCmd, langSetCmd — atomic dual-config write)
  affects:
    - .claude/mysd.yaml (response_language field)
    - openspec/config.yaml (locale field)
tech_stack:
  added: []
  patterns:
    - defer-rollback atomic write (write A first, rollback A if B fails)
    - viper ReadInConfig before Set to preserve existing fields
    - cobra sub-subcommand pattern (langCmd.AddCommand(langSetCmd))
key_files:
  created:
    - cmd/lang.go
    - cmd/lang_test.go
  modified: []
decisions:
  - "lang set uses defer rollback pattern (write mysd.yaml first, rollback if openspec write fails) — safer on Windows than write-then-rename (D-09)"
  - "TestLangSet_AtomicRollback skipped on Windows — chmod dir behavior differs; rollback logic is tested via code review"
metrics:
  duration: "7 min"
  completed_date: "2026-03-26"
  tasks_completed: 1
  files_changed: 2
---

# Phase 07 Plan 04: Lang Command with Atomic Dual-Config Write Summary

**One-liner:** `mysd lang` shows BCP47 language settings; `mysd lang set` atomically updates both `.claude/mysd.yaml` `response_language` and `openspec/config.yaml` `locale` with rollback on failure.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 (RED) | Add failing tests for lang commands | 321f5b5 | cmd/lang_test.go |
| 1 (GREEN) | Implement lang and lang set with atomic write | 8969838 | cmd/lang.go, cmd/lang_test.go |

## What Was Built

### cmd/lang.go

- `langCmd`: reads current language settings from `.claude/mysd.yaml` (`response_language`) and `openspec/config.yaml` (`locale`), outputs both with clear labels
- `langSetCmd`: atomically updates both configs using the defer-rollback pattern:
  1. Read old values
  2. Write `mysd.yaml` via Viper (preserving other fields)
  3. Write `openspec/config.yaml` via `spec.WriteOpenSpecConfig`
  4. On openspec write failure: rollback `mysd.yaml` to old value

### cmd/lang_test.go

Five test cases:
- `TestLangRead_ShowsCurrent`: shows `response_language` and `locale` from both configs
- `TestLangSet_UpdatesBothConfigs`: verifies atomic update of both files
- `TestLangSet_AtomicRollback`: verifies rollback (skipped on Windows — chmod dir semantics differ)
- `TestLangSet_CreatesOpenSpecConfig`: verifies auto-creation of missing `openspec/config.yaml`
- `TestLangSet_PreservesOtherFields`: verifies `project:` field preserved in openspec/config.yaml

## Verification

```
go test ./cmd/... -run TestLang -v   -> PASS (1 SKIP on Windows for chmod test)
go build ./...                        -> OK
```

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Platform] TestLangSet_AtomicRollback skipped on Windows**
- **Found during:** Task 1 GREEN phase
- **Issue:** `os.Chmod(dir, 0555)` does not prevent writes on Windows — the test expected an error but got nil
- **Fix:** Added `isWindows()` helper using `runtime.GOOS` and skip when running on Windows
- **Files modified:** cmd/lang_test.go
- **Commit:** 8969838 (bundled with GREEN commit)

## Known Stubs

None — all functionality fully implemented and wired.

## Self-Check: PASSED

Files exist:
- FOUND: /d/work_data/project/go/mysd/cmd/lang.go
- FOUND: /d/work_data/project/go/mysd/cmd/lang_test.go

Commits exist:
- FOUND: 321f5b5 (RED tests)
- FOUND: 8969838 (GREEN implementation)
