---
spec-version: "1.0"
capability: Task Execution Engine
delta: MODIFIED
status: draft
---

## ADDED Requirements

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

## MODIFIED Requirements

### Requirement: Execution Modes

The system MUST support three execution modes:
- **Single mode**: Tasks executed sequentially by one agent per task
- **Wave mode**: Tasks grouped into parallel waves, each wave executed by multiple agents
- **Spec mode**: Tasks grouped by spec field, each spec executed by one agent handling all tasks in that spec

The `--agent-count` flag MUST control the number of parallel agents in wave mode.

The `execution_mode` configuration value MUST accept "single", "wave", or "spec".

#### Scenario: Single Agent Execution

WHEN execute is called with execution_mode=single
THEN tasks are executed sequentially, one agent per task

#### Scenario: Wave Parallel Execution

WHEN execute is called with execution_mode=wave and agent-count > 1
THEN independent tasks are grouped into waves
AND each wave's tasks run in parallel across agents

#### Scenario: Spec Execution

- **WHEN** execute is called with execution_mode=spec
- **THEN** tasks are grouped by spec field
- **AND** each spec group is handled by one agent executing all tasks in that group sequentially
