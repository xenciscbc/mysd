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

The `--context-only` flag MUST output the JSON without triggering execution.

### Scenario: Single Agent Execution

WHEN execute is called with agent-count=1
THEN tasks are executed sequentially in dependency order

### Scenario: Wave Parallel Execution

WHEN execute is called with agent-count > 1
THEN independent tasks are grouped into waves
AND each wave's tasks run in parallel across agents

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

The system MUST support two execution modes:
- **Single mode**: Tasks executed sequentially by one agent
- **Wave mode**: Tasks grouped into parallel waves, each wave executed by multiple agents

The `--agent-count` flag MUST control the number of parallel agents in wave mode.

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

## Covered Packages

- `cmd/execute.go`, `cmd/task_update.go`
- `internal/executor/` — ExecutionContext building, wave coordination
- `internal/worktree/` — git worktree lifecycle management
- `plugin/commands/mysd-apply.md` — apply command skill with mandatory verification

## Requirements