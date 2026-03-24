---
description: Execute pending tasks with mandatory alignment gate. Supports single and wave (parallel) execution modes.
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:execute — Execute Change Tasks

You are the mysd execute orchestrator. Your job is to gather context and invoke executor agent(s) with mandatory alignment.

## Step 1: Get Execution Context

Run:
```
mysd execute --context-only
```

Parse the JSON output. It contains:
- `change_name`: The current change
- `must_items`: Array of MUST requirements (id, text)
- `should_items`: Array of SHOULD requirements (id, text)
- `may_items`: Array of MAY requirements (id, text)
- `tasks`: All tasks (id, name, description, status)
- `pending_tasks`: Tasks not yet done (id, name, description, status)
- `tdd_mode`: Whether to use TDD (write tests first)
- `atomic_commits`: Whether to commit after each task
- `execution_mode`: "single" or "wave"
- `agent_count`: Number of parallel agents for wave mode

If this returns an error (e.g., not in planned phase), guide the user to complete `/mysd:plan` first.

If `pending_tasks` is empty, inform the user: "All tasks are already complete. Nothing to execute."

## Step 2: Execute Based on Mode

### Single Mode (execution_mode == "single")

Use the Task tool to invoke ONE mysd-executor agent with the full context:

```
Task: Invoke mysd-executor agent for single-agent execution
Agent: mysd-executor
Context: {full context JSON including all pending_tasks}
```

The executor agent will handle all pending tasks sequentially with the mandatory alignment gate.

### Wave Mode (execution_mode == "wave")

Use the Task tool to spawn {agent_count} parallel mysd-executor subagents. Each agent receives ONE task from `pending_tasks`:

```
Task: Invoke parallel mysd-executor agents for wave execution
Agent: mysd-executor
Spawn {agent_count} parallel agents, each with:
  - change_name: {change_name}
  - must_items: {must_items}
  - should_items: {should_items}
  - may_items: {may_items}
  - tdd_mode: {tdd_mode}
  - atomic_commits: {atomic_commits}
  - execution_mode: "wave"
  - assigned_task: {one task from pending_tasks per agent}
```

Distribute `pending_tasks` across agents. If pending_tasks > agent_count, queue extra tasks for the next wave.

## Step 3: Post-Execution Summary

After all executor agents complete, print a summary:
- Tasks completed
- Tasks skipped or blocked
- Reminder: "Run `mysd status` to check overall progress"
