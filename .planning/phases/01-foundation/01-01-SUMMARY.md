---
phase: 01-foundation
plan: 01
subsystem: spec
tags: [go, cobra, yaml, frontmatter, rfc2119, delta-specs, openspec]

requires: []

provides:
  - "Go module github.com/mysd with all core dependencies"
  - "internal/spec/schema.go: complete typed domain model (RFC2119Keyword, DeltaOp, ItemStatus, SpecDirFlavor, Change, Requirement, ProposalDoc, DesignDoc, Task, ChangeMeta)"
  - "internal/spec/detector.go: DetectSpecDir and ListChanges"
  - "internal/spec/parser.go: ParseProposal, ParseSpec, ParseTasks, ParseChangeMeta, ParseChange"
  - "internal/spec/delta.go: DetectDeltaOp, ParseDelta"
  - "internal/spec/writer.go: Scaffold"
  - "Test fixtures for both OpenSpec brownfield and my-ssd native formats"

affects:
  - "02-state-machine (depends on Change, Requirement, ItemStatus types)"
  - "03-cli-skeleton (depends on all internal/spec functions)"
  - "04-release (depends on compilable Go module)"

tech-stack:
  added:
    - "github.com/spf13/cobra v1.10.2 (CLI framework)"
    - "gopkg.in/yaml.v3 (YAML parsing)"
    - "github.com/adrg/frontmatter v0.2.0 (Markdown frontmatter extraction)"
    - "github.com/spf13/viper v1.21.0 (configuration management)"
    - "github.com/charmbracelet/lipgloss v1.1.0 (terminal styling)"
    - "github.com/stretchr/testify v1.11.1 (test assertions)"
  patterns:
    - "Package-level sentinel errors (ErrNoSpecDir, ErrInvalidTransition) for typed error handling"
    - "Graceful brownfield degradation: frontmatter absent = zero-value struct + raw body"
    - "Case-sensitive RFC 2119 regex: \\bMUST\\b matches only uppercase, not 'must'"
    - "TDD: RED (failing tests) -> GREEN (minimal implementation) -> REFACTOR"

key-files:
  created:
    - "go.mod - Go module definition with all dependencies"
    - "main.go - placeholder main package"
    - "internal/spec/schema.go - all typed structs and constants"
    - "internal/spec/detector.go - DetectSpecDir and ListChanges"
    - "internal/spec/parser.go - all parse functions with frontmatter support"
    - "internal/spec/delta.go - delta operation identification and parsing"
    - "internal/spec/writer.go - Scaffold for creating change directories"
    - "internal/spec/schema_test.go - schema constant and struct tests"
    - "internal/spec/detector_test.go - spec directory detection tests"
    - "internal/spec/parser_test.go - parser tests including RFC 2119 case-sensitivity"
    - "internal/spec/delta_test.go - delta operation tests"
    - "internal/spec/writer_test.go - scaffold creation tests"
    - "testdata/fixtures/openspec-project/openspec/... - brownfield fixture"
    - "testdata/fixtures/mysd-project/.specs/... - native my-ssd fixture"
  modified: []

key-decisions:
  - "OpenSpec brownfield fixtures placed under openspec/ subdirectory to match real OpenSpec project structure (not under changes/ directly)"
  - "extractKeyword uses uppercase-only regex; lowercase 'must'/'should'/'may' intentionally not matched"
  - "ParseProposal returns zero-value frontmatter (not error) when no frontmatter present - enables brownfield support"
  - "Scaffold uses text/template stdlib for file generation - no external template dependency needed"

patterns-established:
  - "Pattern: Graceful frontmatter degradation — open file, attempt frontmatter.Parse, on error fall back to reading raw file content"
  - "Pattern: RFC 2119 parsing uses \\bWORD\\b word boundaries to avoid false matches in compound words"
  - "Pattern: Test fixtures split by format: openspec/ (brownfield) vs .specs/ (native) for clear intent"

requirements-completed:
  - SPEC-01
  - SPEC-02
  - SPEC-03
  - SPEC-04
  - SPEC-07
  - OPSX-01
  - OPSX-02
  - OPSX-03
  - OPSX-04
  - DIST-01
  - DIST-02

duration: 10min
completed: 2026-03-23
---

# Phase 01 Plan 01: Go Module Init + Spec Domain Model Summary

**Go module with typed spec domain model (RFC2119Keyword, DeltaOp, Change, Requirement), brownfield-compatible parser using adrg/frontmatter, and Scaffold writer — all tested via TDD with 36 passing tests**

## Performance

- **Duration:** ~10 min
- **Started:** 2026-03-23T08:34:23Z
- **Completed:** 2026-03-23T08:43:54Z
- **Tasks:** 2 (both TDD)
- **Files modified:** 18 created, 0 modified

## Accomplishments

- Go module `github.com/mysd` initialized with Cobra, yaml.v3, adrg/frontmatter, viper, lipgloss, testify
- Complete typed spec domain model in `internal/spec/schema.go` with all constants (RFC2119Keyword, DeltaOp, ItemStatus, SpecDirFlavor) and structs (Change, Requirement, ProposalDoc, DesignDoc, Task, ChangeMeta)
- Full parser stack: DetectSpecDir, ParseProposal, ParseSpec, ParseTasks, ParseChangeMeta, ParseChange — all handle brownfield (no frontmatter) and native (with frontmatter) gracefully
- Delta spec parsing: DetectDeltaOp and ParseDelta correctly categorize requirements into ADDED/MODIFIED/REMOVED
- Scaffold writer creates complete change directory with correctly templated frontmatter using text/template
- 36 tests passing, go vet clean, go build ./... exits 0

## Task Commits

Each task was committed atomically:

1. **Task 1: Go module init + spec schema types + test fixtures** - `cbe00ca` (feat)
2. **Task 2: Spec parser + detector + delta + writer implementation** - `b9b8fc1` (feat)

**Plan metadata:** (pending docs commit)

_Note: Both tasks used TDD (RED → GREEN cycle)_

## Files Created/Modified

- `go.mod` - module github.com/mysd with all dependencies
- `main.go` - placeholder main package
- `internal/spec/schema.go` - complete typed domain model
- `internal/spec/detector.go` - DetectSpecDir (.specs/ vs openspec/) and ListChanges
- `internal/spec/parser.go` - ParseChange + all sub-parsers with brownfield fallback
- `internal/spec/delta.go` - DetectDeltaOp and ParseDelta
- `internal/spec/writer.go` - Scaffold creates change directory with frontmatter templates
- `internal/spec/schema_test.go` - 7 tests for constants and struct fields
- `internal/spec/detector_test.go` - 4 tests for spec dir detection
- `internal/spec/parser_test.go` - 15 tests including lowercase RFC 2119 negative case
- `internal/spec/delta_test.go` - 7 tests for delta operations
- `internal/spec/writer_test.go` - 2 tests for scaffold file creation
- `testdata/fixtures/openspec-project/openspec/...` - brownfield OpenSpec fixture (no frontmatter)
- `testdata/fixtures/mysd-project/.specs/...` - native my-ssd fixture (with frontmatter)

## Decisions Made

- OpenSpec fixture placed under `openspec/` subdirectory — matches real OpenSpec project layout where `openspec/` is the spec root directory
- RFC 2119 parsing is strictly case-sensitive via uppercase-only regex; this is correct per RFC 2119 which specifies keywords "appear in uppercase" to be semantically meaningful
- `ParseProposal` returns zero-value frontmatter struct (not an error) when no frontmatter present — this is the brownfield compatibility design (OPSX-04)
- Used `text/template` stdlib for Scaffold templates — no external dependency needed for simple `{{ .Field }}` substitutions

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Corrected OpenSpec brownfield fixture directory structure**
- **Found during:** Task 2 (TestDetectSpecDir_OpenSpec failure)
- **Issue:** Plan specified `testdata/fixtures/openspec-project/changes/sample-change/` but DetectSpecDir checks for an `openspec/` subdirectory within the project root. The original fixture had no `openspec/` directory, so DetectSpecDir returned ErrNoSpecDir.
- **Fix:** Created `testdata/fixtures/openspec-project/openspec/changes/sample-change/` and moved fixture files there. Updated parser_test.go paths accordingly.
- **Files modified:** Added openspec/ fixture directory; updated 4 parser_test.go path references
- **Verification:** TestDetectSpecDir_OpenSpec passes; all 36 tests pass
- **Committed in:** b9b8fc1 (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (Rule 1 - Bug)
**Impact on plan:** Fix was necessary for correctness — the detector logic was correct but the fixture structure did not match the expected OpenSpec layout. No scope creep.

## Issues Encountered

- testify module needed explicit `go get github.com/stretchr/testify/assert@v1.11.1` to resolve transitive dependency (go.sum entry missing after initial `go get github.com/stretchr/testify@v1`)

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- All typed structs available for Phase 2 state machine (STAT-01, STAT-02, STAT-03)
- Parser handles both OpenSpec brownfield and my-ssd native formats — Phase 2 can build state machine on top of these types
- Scaffold foundation ready — Phase 3 CLI will wire `mysd propose` to `Scaffold()`
- No blockers for Phase 2

## Self-Check: PASSED

- FOUND: go.mod
- FOUND: internal/spec/schema.go
- FOUND: internal/spec/parser.go
- FOUND: internal/spec/detector.go
- FOUND: internal/spec/delta.go
- FOUND: internal/spec/writer.go
- FOUND: .planning/phases/01-foundation/01-01-SUMMARY.md
- FOUND commit: cbe00ca (feat(01-01): Go module init + spec schema types + test fixtures)
- FOUND commit: b9b8fc1 (feat(01-01): spec parser + detector + delta + writer implementation)
- FOUND commit: 8025f58 (docs(01-01): complete plan)
- 36 tests passing, go build ./... exits 0, go vet ./... clean

---
*Phase: 01-foundation*
*Completed: 2026-03-23*
