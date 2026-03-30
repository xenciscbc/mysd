---
spec-version: "1.0"
capability: Per-Spec Execution Mode
delta: ADDED
status: draft
---

## ADDED Requirements

### Requirement: Per-spec execution mode

The `mysd:apply` skill SHALL support `execution_mode: "spec"` as a third execution mode alongside `single` and `wave`.

In spec mode, the orchestrator SHALL:
1. Group pending tasks by their `spec` field
2. For each spec group, spawn one executor agent with all tasks in that group
3. Execute spec groups sequentially (one agent at a time)

The executor agent SHALL receive an `assigned_tasks` array (instead of a single `assigned_task`) containing all tasks for that spec, and SHALL execute them in order within a single agent session.

Spec mode SHALL NOT use worktree isolation — the agent SHALL execute directly in the repo root (same as single mode).

#### Scenario: Spec mode groups tasks by spec field

- **WHEN** execution_mode is "spec"
- **AND** pending tasks have specs: auth (3 tasks), billing (2 tasks)
- **THEN** the orchestrator SHALL spawn 2 executor agents sequentially
- **AND** the first agent SHALL receive 3 tasks for auth
- **AND** the second agent SHALL receive 2 tasks for billing

#### Scenario: Spec mode with --spec flag

- **WHEN** execution_mode is "spec"
- **AND** `--spec auth` is specified
- **THEN** the orchestrator SHALL spawn only 1 executor agent for the auth spec

#### Scenario: Change-level tasks in spec mode

- **WHEN** execution_mode is "spec"
- **AND** tasks exist with no spec field (change-level tasks)
- **THEN** change-level tasks SHALL be grouped together and executed by a separate executor agent after all spec groups

### Requirement: spec-executor model role

The `DefaultModelMap` SHALL include a `spec-executor` role with the following profile mappings:

| Profile | Model |
|---------|-------|
| quality | opus |
| balanced | opus |
| budget | sonnet |

The `mysd:apply` skill SHALL use `ResolveModel("spec-executor", ...)` when `execution_mode` is "spec", and `ResolveModel("executor", ...)` for single and wave modes.

#### Scenario: Spec mode uses spec-executor role

- **WHEN** execution_mode is "spec" and profile is "balanced"
- **THEN** the orchestrator SHALL resolve the model using the "spec-executor" role
- **AND** the resolved model SHALL be "opus"

#### Scenario: Single mode still uses executor role

- **WHEN** execution_mode is "single" and profile is "balanced"
- **THEN** the orchestrator SHALL resolve the model using the "executor" role
- **AND** the resolved model SHALL be "sonnet"

### Requirement: Executor handles assigned_tasks array

The `mysd-executor` agent SHALL accept either `assigned_task` (single object) or `assigned_tasks` (array of task objects) in its context.

When `assigned_tasks` is present, the agent SHALL execute each task in array order, performing the full task execution cycle (pre-task checks, mark in-progress, implement, mark done) for each task within the same agent session.

The agent SHALL re-read the relevant spec and design sections before each task to maintain accuracy despite context compression.

#### Scenario: Executor with single task (backward compatible)

- **WHEN** context contains `assigned_task` (single object)
- **THEN** the executor SHALL execute that single task (existing behavior)

#### Scenario: Executor with multiple tasks

- **WHEN** context contains `assigned_tasks` (array of 3 tasks)
- **THEN** the executor SHALL execute all 3 tasks sequentially within the same session
- **AND** SHALL re-read spec/design before each task
