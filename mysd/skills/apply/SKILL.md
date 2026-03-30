---
description: Execute pending tasks with mandatory alignment gate. Supports single (sequential per-task) and wave (parallel per-task) modes. Usage: /mysd:apply [--auto]
argument-hint: "[--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
  - AskUserQuestion
---

# /mysd:apply — Execute Change Tasks

You are the mysd apply orchestrator. Your job is to execute pending tasks by spawning executor agents, with mandatory alignment gate.

## Step 1: Parse Arguments

Check `$ARGUMENTS` for `--auto`. Set `auto_mode` = true if present, false otherwise.

## Step 2: Get Execution Context

Run:
```
mysd execute --context-only
```

Parse the JSON output. It contains:
- `change_name`: The current change
- `model`: Profile-resolved short model name (e.g., "sonnet", "opus", "haiku") for agent spawning
- `must_items`: Array of MUST requirements (id, text)
- `should_items`: Array of SHOULD requirements (id, text)
- `may_items`: Array of MAY requirements (id, text)
- `tasks`: All tasks (id, name, description, status)
- `pending_tasks`: Tasks not yet done (id, name, description, status)
- `tdd_mode`: Whether to use TDD (write tests first)
- `atomic_commits`: Whether to commit after each task
- `execution_mode`: "single" or "wave"
- `agent_count`: Number of parallel agents for wave mode
- `wave_groups`: Task groups for wave execution
- `worktree_dir`: Base directory for worktrees
- `auto_mode`: From binary config (may be overridden by --auto flag)

If `--auto` was parsed in Step 1, override `auto_mode` to true.

If this returns an error (e.g., not in planned phase), guide the user to complete `/mysd:plan` first.

If `pending_tasks` is empty, inform the user: "All tasks are already complete. Nothing to execute."

## Step 3: Execute Based on Mode

### Single Mode (execution_mode == "single")

For each task in `pending_tasks` (sequential, one at a time):

Show: "Spawning mysd-executor ({model})..."
Use the Task tool to invoke `mysd-executor` with `model` parameter set to `{model}`:
```
Task: Execute task T{task.id}: {task.name}
Agent: mysd-executor
Model: {model}
Context: {
  "change_name": "{change_name}",
  "must_items": [...],
  "should_items": [...],
  "may_items": [...],
  "tasks": [...],
  "assigned_task": {current task object},
  "tdd_mode": {tdd_mode},
  "atomic_commits": {atomic_commits},
  "auto_mode": {auto_mode}
}
```

Wait for completion before spawning next executor. This ensures sequential execution with proper context from each completed task.

### Wave Mode (execution_mode == "wave")

Process `wave_groups` sequentially. Within each wave, spawn executors in parallel:

For each wave in `wave_groups`:
  For each task in wave (spawn in parallel):
    Show: "Spawning mysd-executor ({model})..."
    Use the Task tool to invoke `mysd-executor` with `model` parameter set to `{model}`:
    ```
    Task: Execute wave task T{task.id}: {task.name}
    Agent: mysd-executor
    Model: {model}
    Context: {
      "change_name": "{change_name}",
      "must_items": [...],
      "should_items": [...],
      "may_items": [...],
      "tasks": [...],
      "assigned_task": {current task object},
      "tdd_mode": {tdd_mode},
      "atomic_commits": {atomic_commits},
      "auto_mode": {auto_mode},
      "worktree_path": "{worktree_dir}/T{task.id}",
      "branch": "mysd/{change_name}/T{task.id}-{task.slug}",
      "isolation": "worktree"
    }
    ```

  After all tasks in the wave complete, run merge step for this wave:
  - For each completed task in the wave (ascending task ID order):
    ```
    git checkout main
    git merge --no-ff mysd/{change_name}/T{task.id}-{task.slug}
    ```
  - If merge conflicts occur: attempt AI resolution (up to 3 retries with go build+test verification)
  - On conflict resolution failure: preserve worktree, continue with next task (continue-on-failure policy)

  Then proceed to next wave.

## Step 4: Post-Execution

Run state transition:
```
mysd execute
```

Show summary:
- Tasks completed (count and names)
- Tasks failed or blocked (count and names, if any)
- Proceeding to auto-verify...

If any tasks failed: "Run `/mysd:fix T{id}` to fix failed tasks"

## Step 5: Auto-Verify (D-02, D-05)

After all tasks complete, automatically run verification.

### Step 5a: Build and Test

Run build check:
```
go build ./...
```

If build fails:
- Display the build error output
- Show: "Build failed. Skipping spec verification. Run `/mysd:fix` to address build errors."
- STOP here — do not proceed to verifier agent.

Run test check:
```
go test ./...
```

If tests fail:
- Display the test failure output
- Show: "Tests failed. Skipping spec verification. Run `/mysd:fix` to address test failures."
- STOP here — do not proceed to verifier agent.

### Step 5b: Spec Verification

If build and tests both pass, invoke the verifier agent. Verification is mandatory and cannot be skipped.

First, get fresh context for verification:
```
mysd execute --context-only
```
Parse JSON to get `must_items`, `should_items`, `may_items`.

Show: "Spawning mysd-verifier ({model})..."
Use the Task tool to invoke `mysd-verifier` with `model` parameter set to `{model}`:
```
Task: Verify spec coverage for {change_name}
Agent: mysd-verifier
Model: {model}
Context: {
  "change_name": "{change_name}",
  "must_items": [...],
  "should_items": [...],
  "may_items": [...]
}
```

After verifier completes, show:
- Verification result summary
- If all MUST items pass: "Ready to archive. Run `/mysd:archive`."
- If MUST items fail: "Some MUST requirements not met. Run `/mysd:fix` or re-execute."
