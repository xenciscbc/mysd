---
spec-version: "1.0"
capability: Task Planning & Coverage Validation
delta: ADDED
status: done
---

## Requirement: Task Planning

The `mysd plan` command MUST break a design into executable tasks stored in `tasks.md`.

Tasks frontmatter MUST include: `spec-version`, `change`, `status`.

Each task MUST have: ID (T{n}), description, status, and `satisfies` field mapping to requirement IDs.

The `--research` flag MUST spawn a researcher agent before planning.

The `--context-only` flag MUST output planning context as JSON without writing files.

## Requirement: Plan Coverage Checking

The `planchecker` package MUST validate that all MUST-level requirements are covered by at least one task.

`CheckCoverage()` MUST return a `CoverageResult` with: TotalMust, CoveredCount, UncoveredIDs, CoverageRatio, Passed.

The `--check` flag on `mysd plan` MUST invoke the plan checker after task generation.

Coverage MUST be computed by mapping each task's `satisfies[]` field to MUST requirement IDs.

## Requirement: Wave Grouping

The `executor` package MUST compute parallel execution waves via `ComputeWaves()`.

Tasks with no inter-dependencies MUST be grouped into the same wave for parallel execution.

`BuildAlignmentReport()` MUST validate task dependencies and ordering before execution.

### Scenario: Full Coverage

WHEN all MUST requirements have at least one task with a matching `satisfies` entry
THEN CheckCoverage() returns Passed=true

### Scenario: Missing Coverage

WHEN a MUST requirement has no matching task
THEN CheckCoverage() returns Passed=false with the uncovered ID in UncoveredIDs

## Covered Packages

- `cmd/plan.go`
- `internal/planchecker/` — MUST coverage validation
- `internal/executor/` — wave computation and alignment (shared with execution spec)
