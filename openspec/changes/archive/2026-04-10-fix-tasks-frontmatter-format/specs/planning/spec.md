---
spec-version: "1.0"
capability: Task Planning & Coverage Validation
delta: MODIFIED
status: draft
---

## MODIFIED Requirements

### Requirement: Task Planning

The `mysd plan` command MUST break a design into executable tasks stored in `tasks.md`.

Tasks MUST use YAML frontmatter with the following fields:
- `spec-version`: string (e.g., `"1.0"`)
- `total`: integer — total number of tasks
- `completed`: integer — number of tasks with status `done`
- `tasks`: array of task entries

Each task entry in the `tasks` array MUST have:
- `id`: integer — sequential task identifier starting from 1
- `name`: string — short task name
- `status`: one of `pending`, `in-progress`, `done`, `skipped`

Each task entry MUST include a `satisfies` field — an array of requirement IDs (strings) that the task fulfills.

Each task entry MAY include:
- `description`: string — detailed description of what to implement
- `spec`: string — spec directory name the task belongs to
- `depends`: array of integer task IDs — tasks that must complete first
- `files`: array of strings — file paths touched by this task
- `skills`: array of strings — slash commands used for this task

The markdown body after the frontmatter MUST be empty or contain only supplementary notes. The structured task data MUST live exclusively in the YAML frontmatter `tasks` array.

The `--research` flag MUST spawn a researcher agent before planning.

The `--context-only` flag MUST output planning context as JSON without writing files.

#### Scenario: Valid tasks.md with frontmatter

- **WHEN** the planner agent writes tasks.md
- **THEN** the file SHALL contain YAML frontmatter with `spec-version`, `total`, `completed`, and `tasks` array
- **AND** each task in the array SHALL have `id`, `name`, `status`, and `satisfies` fields

#### Scenario: Task status round-trip via UpdateTaskStatus

- **WHEN** `UpdateTaskStatus()` is called with a task ID and new status
- **THEN** it SHALL parse the YAML frontmatter, update the matching task's status, recompute `completed`, and write back the file preserving body content

#### Scenario: Tasks without frontmatter fallback

- **WHEN** `ParseTasksV2()` reads a tasks.md file without YAML frontmatter
- **THEN** it SHALL return a zero-value `TasksFrontmatterV2` and the entire file content as body
- **AND** `UpdateTaskStatus()` SHALL fail gracefully (no matching task ID found)

### Requirement: Plan Coverage Checking

The `planchecker` package MUST validate that all MUST-level requirements are covered by at least one task.

`CheckCoverage()` MUST return a `CoverageResult` with: TotalMust, CoveredCount, UncoveredIDs, CoverageRatio, Passed.

The `--check` flag on `mysd plan` MUST invoke the plan checker after task generation.

Coverage MUST be computed by mapping each task's `satisfies[]` field to MUST requirement IDs.

#### Scenario: Full Coverage

- **WHEN** all MUST requirements have at least one task with a matching `satisfies` entry
- **THEN** CheckCoverage() returns Passed=true

#### Scenario: Missing Coverage

- **WHEN** a MUST requirement has no matching task
- **THEN** CheckCoverage() returns Passed=false with the uncovered ID in UncoveredIDs

### Requirement: Wave Grouping

The `executor` package MUST compute parallel execution waves via `ComputeWaves()`.

Tasks with no inter-dependencies MUST be grouped into the same wave for parallel execution.

`BuildAlignmentReport()` MUST validate task dependencies and ordering before execution.

#### Scenario: Independent tasks in same wave

- **WHEN** tasks T1 and T2 have no `depends` entries referencing each other
- **THEN** `ComputeWaves()` SHALL place them in the same wave

#### Scenario: Dependent tasks in sequential waves

- **WHEN** task T2 has `depends: [1]`
- **THEN** `ComputeWaves()` SHALL place T2 in a later wave than T1
