---
phase: 07-new-binary-commands-scanner-refactor
plan: "02"
subsystem: scanner
tags: [scanner, language-detection, init, scaffold, refactor]
dependency_graph:
  requires: []
  provides: [scanner.ScanContext, scanner.ModuleInfo, scanner.BuildScanContext, cmd.scaffoldOpenSpecDir]
  affects: [internal/scanner, cmd/scan.go, cmd/init_cmd.go]
tech_stack:
  added: []
  patterns: [language-marker-detection, idempotent-scaffold, file-extension-counting]
key_files:
  created: []
  modified:
    - internal/scanner/scanner.go
    - internal/scanner/scanner_test.go
    - cmd/scan.go
    - cmd/init_cmd.go
    - cmd/scan_test.go
    - cmd/init_cmd_test.go
decisions:
  - "ScanContext replaces PackageInfo entirely — no backward compat (D-02)"
  - "detectPrimaryLanguage uses marker file priority: go.mod > package.json > requirements.txt > pyproject.toml"
  - "Files map counts ALL file extensions (not just primary language)"
  - "scaffoldOpenSpecDir is idempotent via os.MkdirAll — no --force flag needed"
  - "init does not create openspec/config.yaml per D-06 (locale set by SKILL.md agent)"
  - "init preserves existing .claude/mysd.yaml (idempotent, no overwrite)"
metrics:
  duration: "6 min"
  completed: "2026-03-26T01:43:25Z"
  tasks_completed: 2
  files_modified: 6
---

# Phase 07 Plan 02: Language-Agnostic Scanner Refactor Summary

Language-agnostic ScanContext replacing Go-specific PackageInfo, with detectPrimaryLanguage() via file markers (go.mod/package.json/pyproject.toml) and scaffoldOpenSpecDir() shared by scan --scaffold-only and init.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Replace ScanContext with language-agnostic struct | 81207c1 | internal/scanner/scanner.go, internal/scanner/scanner_test.go |
| 2 | Add --scaffold-only flag to scan and rewire init | fc6dab1 | cmd/scan.go, cmd/init_cmd.go, cmd/scan_test.go, cmd/init_cmd_test.go |

## What Was Built

### Task 1: Language-Agnostic Scanner (TDD)

Completely replaced the Go-specific scanner with a universal implementation:

- **Removed:** `PackageInfo` struct, `Packages []PackageInfo` field — no backward compat per D-02
- **New `ScanContext`:** `PrimaryLanguage string`, `Files map[string]int`, `Modules []ModuleInfo`, `ConfigExists bool`
- **`detectPrimaryLanguage(root)`:** Checks markers in priority order — go.mod → "go", package.json → "nodejs", requirements.txt/pyproject.toml → "python", fallback → "unknown"
- **`BuildScanContext`:** Counts ALL files by extension (not just .go), builds module list per language, detects openspec/changes/ for ExistingSpecs, checks openspec/config.yaml for ConfigExists
- **Module detection per language:** Go = dirs with .go files; Node.js = dirs with index.js/ts or package.json; Python = dirs with __init__.py; unknown = any dir with files

10 tests cover all language scenarios, edge cases, and config detection.

### Task 2: Scaffold Mode and Init Rewire

- **`scaffoldOpenSpecDir(root string) error`:** Creates openspec/ and openspec/specs/ via os.MkdirAll (idempotent). Does NOT create openspec/config.yaml per D-06.
- **`scan --scaffold-only`:** New flag routes to runScanScaffoldOnly → scaffoldOpenSpecDir
- **Updated `runScan`:** Returns error if neither --context-only nor --scaffold-only is provided
- **Rewired `init`:** Delegates to scaffoldOpenSpecDir, creates .claude/mysd.yaml only if absent (idempotent), removed --force flag
- **Updated tests:** scan_test.go updated for new ScanContext fields; init_cmd_test.go rewritten for new idempotent behavior

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Updated cmd/scan_test.go for removed Packages field**
- **Found during:** Task 2
- **Issue:** scan_test.go referenced `ctx.Packages` and `PackageInfo` fields removed in Task 1
- **Fix:** Rewrote tests to use new `ctx.PrimaryLanguage` and `ctx.Files` fields; added TestScaffoldOpenSpecDir tests
- **Files modified:** cmd/scan_test.go
- **Commit:** fc6dab1

**2. [Rule 1 - Bug] Updated cmd/init_cmd_test.go for removed --force flag**
- **Found during:** Task 2
- **Issue:** init_cmd_test.go tested `--force` flag behavior and "existing file warning" — both removed by new idempotent design
- **Fix:** Rewrote TestInit_ExistingFile_NoForce_DoesNotOverwrite to test idempotent behavior; removed TestInit_ExistingFile_WithForce_Overwrites; added TestInit_CreatesOpenspecStructure and TestInit_Idempotent
- **Files modified:** cmd/init_cmd_test.go
- **Commit:** fc6dab1

## Verification Results

```
ok  github.com/xenciscbc/mysd/internal/scanner  0.619s
ok  github.com/xenciscbc/mysd/cmd               1.522s
go build ./...  (exit 0)
```

## Known Stubs

None — all data sources are wired. ScanContext fields are populated from actual filesystem walk, detectPrimaryLanguage uses real os.Stat calls, and ConfigExists checks the actual path.

## Self-Check: PASSED
