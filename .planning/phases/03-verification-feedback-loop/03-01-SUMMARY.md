---
phase: 03-verification-feedback-loop
plan: "01"
subsystem: verifier
tags: [go, crc32, verification, spec-parsing, json, markdown]

requires:
  - phase: 02-execution-engine
    provides: "spec.Requirement, spec.Task, spec.ParseChange, BuildContextFromParts pattern"

provides:
  - "internal/verifier package: VerificationContext, VerifyItem, TaskItem, StableID, BuildVerificationContext, BuildVerificationContextFromParts"
  - "internal/verifier/report.go: VerifierReport, ParseVerifierReport, WriteGapReport, WriteVerificationReport"
  - "internal/spec/updater.go: VerificationStatus sidecar CRUD (Read/Write/UpdateItemStatus)"
  - "spec.Requirement.SourceFile field populated by ParseSpec"

affects:
  - 03-02 (verify CLI command uses VerificationContext and report writers)
  - 03-03 (verifier agent SKILL.md consumes VerificationContext JSON)

tech-stack:
  added: ["hash/crc32 (stdlib)", "encoding/json (stdlib)", "time (stdlib)"]
  patterns:
    - "StableID: CRC32 hash of requirement text produces deterministic IDs"
    - "Pure-function constructors (BuildVerificationContextFromParts) for test isolation"
    - "Sidecar JSON pattern: verification-status.json tracks state without modifying spec.md"
    - "TDD: RED test first, then GREEN implementation"

key-files:
  created:
    - internal/verifier/context.go
    - internal/verifier/context_test.go
    - internal/verifier/report.go
    - internal/verifier/report_test.go
  modified:
    - internal/spec/schema.go
    - internal/spec/parser.go
    - internal/spec/parser_test.go
    - internal/spec/updater.go
    - internal/spec/updater_test.go

key-decisions:
  - "StableID uses CRC32 hash of requirement text (not sequential counter) — stable across re-parses as long as text is unchanged"
  - "BuildVerificationContextFromParts is a pure function (no I/O) — enables test isolation without filesystem fixtures"
  - "WriteGapReport skips file creation when no MUST failures — avoids empty gap-report.md clutter"
  - "VerificationStatus sidecar (verification-status.json) does not modify spec.md — per D-04 research recommendation"
  - "failed_task_ids in gap-report.md frontmatter is left empty — task ID mapping delegated to CLI layer"

patterns-established:
  - "Verifier context: pure function + disk-loading wrapper (same pattern as executor BuildContextFromParts/BuildContext)"
  - "Report writers use strings.Builder for in-memory construction before single WriteFile call"
  - "Sidecar JSON: ReadXxx returns zero-value with empty map on missing file (not error)"

requirements-completed: [VRFY-01, VRFY-03, VRFY-04, VRFY-05, SPEC-05]

duration: 4min
completed: 2026-03-24
---

# Phase 3 Plan 01: Verification Engine Core Summary

**CRC32-stable requirement IDs, VerificationContext builder, VerifierReport parser, gap/verification markdown report writers, and verification-status.json sidecar CRUD — the complete verification engine foundation**

## Performance

- **Duration:** 4 min
- **Started:** 2026-03-24T01:45:30Z
- **Completed:** 2026-03-24T01:49:18Z
- **Tasks:** 2
- **Files modified:** 9 (4 created, 5 modified)

## Accomplishments

- Added `SourceFile string` field to `spec.Requirement` and populated it via `filepath.Base(path)` in `ParseSpec`
- Built `internal/verifier` package with `VerificationContext` builder using CRC32-based stable IDs in format `{source_file}::{keyword}-{hex_hash}`
- Implemented `VerifierReport` parser (`ParseVerifierReport`), gap report writer (`WriteGapReport`), and verification report writer (`WriteVerificationReport`)
- Added `VerificationStatus` sidecar (read/write/update) to `internal/spec/updater.go` — tracks verification state without modifying spec.md

## Task Commits

Each task was committed atomically:

1. **Task 1: Parser SourceFile enhancement + Stable Requirement IDs + VerificationContext builder** - `8cc9982` (feat)
2. **Task 2: VerifierReport parser + GapReport writer + VerificationReport writer + VerificationStatus sidecar** - `ba926ca` (feat)

## Files Created/Modified

- `internal/spec/schema.go` - Added `SourceFile string` field to `Requirement` struct
- `internal/spec/parser.go` - `ParseSpec` now fills `SourceFile` via `filepath.Base(path)` on each returned requirement
- `internal/spec/parser_test.go` - Added `TestParseSpec_FillsSourceFile` and `TestParseRequirementsFromBody_RegressionKeywords`
- `internal/spec/updater.go` - Added `VerificationStatus`, `ReadVerificationStatus`, `WriteVerificationStatus`, `UpdateItemStatus`
- `internal/spec/updater_test.go` - Added `TestReadVerificationStatus_NoFile`, `TestWriteVerificationStatus`, `TestUpdateItemStatus`, `TestUpdateItemStatus_CreatesFileIfNotExist`
- `internal/verifier/context.go` - New file: `VerificationContext`, `VerifyItem`, `TaskItem`, `StableID`, `BuildVerificationContextFromParts`, `BuildVerificationContext`
- `internal/verifier/context_test.go` - New file: `TestStableID`, `TestBuildVerificationContextFromParts`, `TestVerificationContext_ItemClassification`, `TestBuildVerificationContext_EmptySpecs`, `TestVerifyItem_IDFormat`, `TestBuildVerificationContext_TasksSummary`
- `internal/verifier/report.go` - New file: `VerifierReport`, `VerifierResultItem`, `UIItem`, `ParseVerifierReport`, `WriteGapReport`, `WriteVerificationReport`
- `internal/verifier/report_test.go` - New file: `TestParseVerifierReport`, `TestParseVerifierReport_InvalidJSON`, `TestWriteGapReport`, `TestWriteGapReport_NoFailures`, `TestWriteVerificationReport`

## Decisions Made

- **StableID uses CRC32 (not sequential)**: Re-parsing the same spec.md produces identical IDs as long as the text is unchanged. Sequential counters would break when requirements are reordered.
- **Pure-function constructor**: `BuildVerificationContextFromParts` takes pre-loaded data — no filesystem I/O — enabling unit tests without filesystem fixtures (same pattern as `executor.BuildContextFromParts`).
- **WriteGapReport skips file on zero failures**: Avoids creating empty/misleading gap-report.md when verification passes entirely.
- **Sidecar does not modify spec.md**: Per research D-04 — direct spec.md modification creates merge conflicts and breaks brownfield compatibility; sidecar JSON is the safe alternative.
- **failed_task_ids empty in gap-report.md**: Task ID mapping is done at the CLI layer when it has full context of which task produced which requirement.

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

- `VerificationContext` JSON structure ready for `mysd verify --context-only` (Plan 03-02)
- `VerifierReport` parsing ready for `mysd verify --write-results` (Plan 03-02)
- Gap report and verification report writers ready for CLI layer integration (Plan 03-02)
- Verification-status.json sidecar replaces direct spec.md modification (per research recommendation)
- All `spec.Requirement` structs now have `SourceFile` populated after `ParseSpec`

---
*Phase: 03-verification-feedback-loop*
*Completed: 2026-03-24*
