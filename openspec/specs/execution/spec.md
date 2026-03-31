---
spec-version: "1.0"
capability: Task Execution Engine
delta: MODIFIED
status: done
---

## Requirement: Execution Context

The `mysd execute` command MUST build an `ExecutionContext` JSON for executor agents containing:
- ChangeName, MustItems, ShouldItems, MayItems
- Tasks, PendingTasks, WaveGroups
- AtomicCommits flag, ExecutionMode, AgentCount
- Instruction (dynamically generated guidance string)

The `--context-only` flag MUST output the JSON without triggering execution.

The `--context-only` path SHALL call `runPreflight` internally to obtain a `PreflightReport`, then pass both the `ExecutionContext` and the report to `GenerateInstruction` to populate the `Instruction` field before JSON serialization. The preflight data SHALL NOT appear in the `--context-only` JSON output — it is consumed only by the instruction generator.

### Scenario: Single Agent Execution

- **WHEN** execute is called with agent-count=1
- **THEN** tasks are executed sequentially in dependency order

### Scenario: Wave Parallel Execution

- **WHEN** execute is called with agent-count > 1
- **THEN** independent tasks are grouped into waves
- **AND** each wave's tasks run in parallel across agents

### Scenario: Context-only includes instruction field

- **WHEN** `mysd execute --context-only` is run
- **THEN** the JSON output SHALL contain an `instruction` field
- **AND** the instruction SHALL reflect the current task state and any preflight issues

## Requirement: Apply command verification is mandatory

The `/mysd:apply` command SHALL always run spec verification after task execution completes successfully. The verification step SHALL NOT be skippable by user interaction.

In auto mode, verification SHALL proceed without confirmation. In interactive mode, verification SHALL also proceed without confirmation — the user prompt asking whether to run verification SHALL be removed.

### Scenario: Apply runs verification automatically

WHEN `/mysd:apply` completes all tasks and build+tests pass
THEN the verifier agent SHALL be invoked automatically without asking the user

### Scenario: Apply skips verification only on build failure

WHEN `/mysd:apply` completes tasks but `go build` or `go test` fails
THEN verification SHALL be skipped
AND the user SHALL be informed to run `/mysd:fix`

## Requirement: Execution Modes

The system MUST support three execution modes:
- **Single mode**: Tasks executed sequentially by one agent per task
- **Wave mode**: Tasks grouped into parallel waves, each wave executed by multiple agents
- **Spec mode**: Tasks grouped by spec field, each spec executed by one agent handling all tasks in that spec

The `--agent-count` flag MUST control the number of parallel agents in wave mode.

The `execution_mode` configuration value MUST accept "single", "wave", or "spec".

### Scenario: Spec Execution

- **WHEN** execute is called with execution_mode=spec
- **THEN** tasks are grouped by spec field
- **AND** each spec group is handled by one agent executing all tasks in that group sequentially

## Requirement: Task Status Updates

The `task_update` helper MUST update task status in `tasks.md` (pending → in_progress → done/blocked).

Status updates MUST preserve the YAML frontmatter and markdown structure.

## Requirement: Atomic Commits

The `--atomic-commits` flag MUST instruct executor agents to commit after each task completion.

## Requirement: TDD Mode

The `--tdd` flag MUST instruct executor agents to write tests before implementation code.

## Requirement: Git Worktree Isolation

The `worktree` package MUST create isolated git worktrees at `.worktrees/T{id}/` for parallel task execution.

`Create()` MUST create a branch named `mysd/{change}/T{id}-{slug}`.

`Clean()` MUST prune completed worktrees.

`CheckDiskSpace()` MUST verify at least 500MB available before creating a worktree.

## Removed Requirements

### ~~Requirement: Execute command skill~~

**Removed**: The `/mysd:execute` command skill has been renamed to `/mysd:apply`. The redirect stub is no longer needed as all references have been updated.

**Migration**: Use `/mysd:apply` for task execution.

### ~~Requirement: Spec command skill~~

**Removed**: Spec writing is now embedded within `/mysd:propose` (Step 11) and `/mysd:discuss` (Step 10). A standalone `/mysd:spec` command is redundant.

**Migration**: Use `/mysd:propose` for initial spec generation or `/mysd:discuss` for spec refinement.

### ~~Requirement: Design command skill~~

**Removed**: Design document creation is now embedded within `/mysd:plan` (Step 4). A standalone `/mysd:design` command is redundant.

**Migration**: Use `/mysd:plan` which includes the design phase automatically.

### ~~Requirement: Capture command skill~~

**Removed**: Conversation capture is fully covered by `/mysd:discuss` which provides the same functionality plus research, gray area exploration, and spec updates.

**Migration**: Use `/mysd:discuss` to capture conversation context into structured proposals and specs.

## Requirement: Executor pre-task checks

The `mysd-executor` agent SHALL perform 4 pre-task checks before writing any implementation code for each task. These checks SHALL be executed after marking the task in-progress and before TDD or implementation steps.

The checks SHALL be:
1. **Reuse**: Search adjacent modules and shared utilities for existing implementations that could be reused instead of writing new code
2. **Quality**: Verify that existing types and constants are used instead of redefining; derive values from existing state instead of duplicating
3. **Efficiency**: Verify that independent async operations are parallelized; match operation scope to actual need
4. **No Placeholders**: Read the spec and design sections relevant to this task and verify no TBD/TODO/FIXME/vague language exists. If placeholders are found, the executor SHALL pause and report instead of implementing against vague requirements.

### Scenario: Pre-task checks find reuse opportunity

- **WHEN** the executor is about to implement a string formatting utility
- **AND** a similar utility exists in an adjacent module
- **THEN** the executor SHALL use the existing utility instead of creating a new one

### Scenario: Pre-task checks find placeholder in spec

- **WHEN** the executor reads the spec for the current task
- **AND** the spec contains "TBD" in the relevant requirement
- **THEN** the executor SHALL pause and report the placeholder instead of implementing

## Requirement: Executor pause conditions

The `mysd-executor` agent SHALL pause and report (instead of guessing or continuing) when any of the following conditions are met:

1. **Task unclear**: Task description is ambiguous or contradictory
2. **Design issue discovered**: Implementation reveals a missing or contradictory design decision
3. **Error/blocker**: Build failure, dependency issue, or technical obstacle
4. **User interrupt**: User explicitly requests a pause

When pausing, the executor SHALL output:
- A description of the issue
- Suggested resolution options (2-3 concrete choices)
- A request for guidance before continuing

### Scenario: Task description is ambiguous

- **WHEN** the executor cannot determine the expected behavior from the task description
- **THEN** the executor SHALL pause, describe the ambiguity, and present resolution options

### Scenario: Design gap discovered during implementation

- **WHEN** the executor discovers that the design does not cover a necessary decision
- **THEN** the executor SHALL pause, describe the missing design decision, and suggest updating the design artifact

## Covered Packages

- `cmd/execute.go`, `cmd/task_update.go`
- `internal/executor/` — ExecutionContext building, wave coordination
- `internal/worktree/` — git worktree lifecycle management
- `plugin/commands/mysd-apply.md` — apply command skill with mandatory verification

## Requirements

### Requirement: TaskItem includes spec field

The `TaskItem` struct SHALL include a `Spec` field (`json:"spec,omitempty"`).

The `spec` field value SHALL correspond to the spec directory name (e.g., `material-selection` for `specs/material-selection/spec.md`).

Tasks without a `spec` field SHALL be treated as change-level tasks.

The YAML frontmatter parser SHALL read and write the `spec` field in TasksFrontmatterV2 format.

#### Scenario: TaskItem with spec field

- **WHEN** a task has `spec: "material-selection"` in the YAML frontmatter
- **THEN** the parsed `TaskItem` SHALL have `Spec: "material-selection"`

#### Scenario: TaskItem without spec field

- **WHEN** a task has no `spec` field in the YAML frontmatter
- **THEN** the parsed `TaskItem` SHALL have `Spec: ""` (empty string)

<!-- @trace
source: enhance-plan-pipeline, enhance-apply-pipeline
updated: 2026-03-30
code:
  - mysd/skills/apply/SKILL.md
  - mysd/agents/mysd-planner.md
  - mysd/agents/mysd-designer.md
  - internal/executor/context.go
  - mysd/skills/plan/SKILL.md
  - cmd/execute.go
  - cmd/instructions.go
  - cmd/plan.go
tests:
  - cmd/execute_test.go
  - internal/spec/schema.go
  - cmd/plan_test.go
  - internal/spec/updater.go
  - cmd/instructions_test.go
  - internal/spec/schema_test.go
  - internal/executor/context_test.go
-->

---
### Requirement: Executor pre-task checks

The `mysd-executor` agent SHALL perform 4 pre-task checks before writing any implementation code for each task. These checks SHALL be executed after marking the task in-progress and before TDD or implementation steps.

The checks SHALL be:
1. **Reuse**: Search adjacent modules and shared utilities for existing implementations that could be reused instead of writing new code
2. **Quality**: Verify that existing types and constants are used instead of redefining; derive values from existing state instead of duplicating
3. **Efficiency**: Verify that independent async operations are parallelized; match operation scope to actual need
4. **No Placeholders**: Read the spec and design sections relevant to this task and verify no TBD/TODO/FIXME/vague language exists. If placeholders are found, the executor SHALL pause and report instead of implementing against vague requirements.

#### Scenario: Pre-task checks find reuse opportunity

- **WHEN** the executor is about to implement a string formatting utility
- **AND** a similar utility exists in an adjacent module
- **THEN** the executor SHALL use the existing utility instead of creating a new one

#### Scenario: Pre-task checks find placeholder in spec

- **WHEN** the executor reads the spec for the current task
- **AND** the spec contains "TBD" in the relevant requirement
- **THEN** the executor SHALL pause and report the placeholder instead of implementing


<!-- @trace
source: enhance-apply-pipeline
updated: 2026-03-30
code:
  - mysd/agents/mysd-executor.md
  - internal/config/config.go
  - cmd/execute.go
  - mysd/skills/apply/SKILL.md
  - internal/config/defaults.go
tests:
  - cmd/execute_test.go
  - internal/executor/context_test.go
  - internal/config/config_test.go
-->

---
### Requirement: Executor pause conditions

The `mysd-executor` agent SHALL pause and report (instead of guessing or continuing) when any of the following conditions are met:

1. **Task unclear**: Task description is ambiguous or contradictory
2. **Design issue discovered**: Implementation reveals a missing or contradictory design decision
3. **Error/blocker**: Build failure, dependency issue, or technical obstacle
4. **User interrupt**: User explicitly requests a pause

When pausing, the executor SHALL output:
- A description of the issue
- Suggested resolution options (2-3 concrete choices)
- A request for guidance before continuing

#### Scenario: Task description is ambiguous

- **WHEN** the executor cannot determine the expected behavior from the task description
- **THEN** the executor SHALL pause, describe the ambiguity, and present resolution options

#### Scenario: Design gap discovered during implementation

- **WHEN** the executor discovers that the design does not cover a necessary decision
- **THEN** the executor SHALL pause, describe the missing design decision, and suggest updating the design artifact

<!-- @trace
source: enhance-apply-pipeline
updated: 2026-03-30
code:
  - mysd/agents/mysd-executor.md
  - internal/config/config.go
  - cmd/execute.go
  - mysd/skills/apply/SKILL.md
  - internal/config/defaults.go
tests:
  - cmd/execute_test.go
  - internal/executor/context_test.go
  - internal/config/config_test.go
-->