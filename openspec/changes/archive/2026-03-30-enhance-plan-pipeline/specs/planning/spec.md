---
spec-version: "1.0"
capability: Task Planning & Coverage Validation
delta: MODIFIED
status: draft
---

## MODIFIED Requirements

### Requirement: Task Planning

The `mysd plan` command MUST break a design into executable tasks stored in `tasks.md`.

Tasks frontmatter MUST include: `spec-version`, `change`, `status`.

Each task MUST have: ID (T{n}), description, status, `satisfies` field mapping to requirement IDs, and an optional `spec` field identifying which spec the task belongs to.

The `--research` flag MUST spawn a researcher agent before planning.

The `--context-only` flag MUST output planning context as JSON without writing files.

The `--spec <name>` flag MUST restrict planning to the specified spec only. When used with `--context-only`, the output SHALL include only the requirements and design sections relevant to that spec.

The `--from <path>` flag MUST read the specified file and include its content as `external_input` in the plan context JSON. The planner agent SHALL use this content as reference context (equivalent to research findings) when generating tasks.

#### Scenario: Full Coverage

WHEN all MUST requirements have at least one task with a matching `satisfies` entry
THEN CheckCoverage() returns Passed=true

#### Scenario: Missing Coverage

WHEN a MUST requirement has no matching task
THEN CheckCoverage() returns Passed=false with the uncovered ID in UncoveredIDs

#### Scenario: Per-spec planning with --spec flag

- **WHEN** `mysd plan --spec material-selection` is executed
- **THEN** the planner SHALL generate tasks only for the `material-selection` spec
- **AND** each generated task SHALL have `spec: "material-selection"`
- **AND** new task IDs SHALL start from the current maximum ID + 1 in tasks.md

#### Scenario: Per-spec planning merges into existing tasks

- **WHEN** tasks.md already contains tasks for spec `planning` (IDs 1-5)
- **AND** `mysd plan --spec material-selection` is executed
- **THEN** new tasks SHALL be appended with IDs starting from 6
- **AND** existing tasks for spec `planning` SHALL remain unchanged

#### Scenario: External input via --from flag

- **WHEN** `mysd plan --from ~/gstack-plan.md` is executed
- **THEN** the plan context JSON SHALL include `external_input` with the file content
- **AND** the planner agent SHALL reference this content when generating tasks

#### Scenario: Interactive spec selection without --spec flag

- **WHEN** `mysd plan` is executed without `--spec` and `auto_mode` is false
- **AND** the change has 3 specs: `material-selection`, `planning`, `execution`
- **AND** `material-selection` has no tasks yet
- **THEN** the orchestrator SHALL present a selection list:
  1. material-selection (no tasks)
  2. planning (5 tasks)
  3. execution (3 tasks)
  4. [All]
- **AND** wait for user selection

#### Scenario: Auto mode plans all specs

- **WHEN** `mysd plan --auto` is executed without `--spec`
- **THEN** the planner SHALL generate tasks for all specs in the change

## ADDED Requirements

### Requirement: Plan pipeline uses mysd instructions for agent guidance

The `mysd:plan` skill SHALL call `mysd instructions <artifact-id> --change <name> --json` before spawning each agent (designer, planner).

The orchestrator SHALL pass the instructions output (template, rules, instruction, selfReviewChecklist) as part of the agent's context.

The agent SHALL use `template` as the output structure, `rules` as constraints, and `selfReviewChecklist` as a verification guide before completing.

#### Scenario: Planner receives instructions

- **WHEN** the plan orchestrator is about to spawn mysd-planner
- **THEN** it SHALL first call `mysd instructions tasks --change <name> --json`
- **AND** pass the result in the agent's context

#### Scenario: Designer receives instructions

- **WHEN** the plan orchestrator is about to spawn mysd-designer
- **THEN** it SHALL first call `mysd instructions design --change <name> --json`
- **AND** pass the result in the agent's context
