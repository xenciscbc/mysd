---
spec-version: "1.0"
capability: Task Execution Engine
delta: MODIFIED
status: draft
---

## MODIFIED Requirements

### Requirement: Execution Context

The `mysd execute` command MUST build an `ExecutionContext` JSON for executor agents containing:
- ChangeName, MustItems, ShouldItems, MayItems
- Tasks, PendingTasks, WaveGroups
- AtomicCommits flag, ExecutionMode, AgentCount

The `--context-only` flag MUST output the JSON without triggering execution.

The `--spec <name>` flag MUST filter `PendingTasks` to only tasks with the matching `spec` field. `WaveGroups` SHALL be recomputed from the filtered task set. Tasks without a `spec` field (change-level tasks) SHALL NOT be included in per-spec filtering unless "All" is selected.

#### Scenario: Single Agent Execution

WHEN execute is called with agent-count=1
THEN tasks are executed sequentially in dependency order

#### Scenario: Wave Parallel Execution

WHEN execute is called with agent-count > 1
THEN independent tasks are grouped into waves
AND each wave's tasks run in parallel across agents

#### Scenario: Per-spec execution with --spec flag

- **WHEN** `mysd execute --spec material-selection --context-only` is executed
- **AND** tasks.md contains tasks for specs `material-selection`, `planning`, and `execution`
- **THEN** `pending_tasks` SHALL contain only tasks with `spec: "material-selection"`
- **AND** `wave_groups` SHALL be computed from the filtered tasks only

#### Scenario: Interactive spec selection without --spec flag

- **WHEN** `mysd execute` is executed without `--spec` and `auto_mode` is false
- **AND** tasks.md has pending tasks in 2 specs: `material-selection` (3 pending), `planning` (2 pending)
- **THEN** the orchestrator SHALL present a selection list:
  1. material-selection (3 pending tasks)
  2. planning (2 pending tasks)
  3. [All] (5 pending tasks)
- **AND** wait for user selection

#### Scenario: Auto mode executes all pending tasks

- **WHEN** `mysd execute --auto` is executed without `--spec`
- **THEN** all pending tasks SHALL be included regardless of spec

#### Scenario: Change-level tasks included in All

- **WHEN** tasks.md contains tasks without a `spec` field (change-level tasks)
- **AND** user selects "All" or `--auto` mode is active
- **THEN** change-level tasks SHALL be included in execution

## ADDED Requirements

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
