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

You are the mysd execute orchestrator. Your job is to gather context, decide execution mode, and invoke executor agent(s) with mandatory alignment.

## Step 1: Get Execution Context

Run:
```
mysd execute --context-only
```

Parse the JSON output. It contains:
- `change_name`: The current change
- `pending_tasks`: Tasks not yet done (id, name, description, status, depends, files, skills)
- `wave_groups`: Pre-computed wave groups — array of arrays, each inner array is one wave of tasks that can execute in parallel. Computed by the binary based on `depends` and `files` fields.
- `worktree_dir`: Base directory for worktrees (e.g., `.worktrees`)
- `auto_mode`: If true, skip asking and use wave mode when parallel opportunity exists
- `has_parallel_opportunity`: True if wave_groups contains at least one wave with 2+ tasks
- `must_items`: Array of MUST requirements (id, text)
- `should_items`: Array of SHOULD requirements (id, text)
- `may_items`: Array of MAY requirements (id, text)
- `tdd_mode`: Whether to use TDD (write tests first)
- `atomic_commits`: Whether to commit after each task

If this returns an error (e.g., not in planned phase), guide the user to complete `/mysd:plan` first.

If `pending_tasks` is empty, inform the user: "All tasks are already complete. Nothing to execute."

---

## Step 2: Decide Execution Mode

Apply these rules in order:

1. **If `has_parallel_opportunity` is false** (all tasks are sequential with no parallel opportunity):
   - Use sequential mode. Do NOT ask the user. Skip to Step 3A.

2. **If `auto_mode` is true** (user ran with `--auto` or `ffe`):
   - If `wave_groups` has any wave with 2+ tasks: use wave mode. Skip to Step 3B.
   - Otherwise: use sequential mode. Skip to Step 3A.

3. **Otherwise** (parallel opportunity exists, user did not set auto mode):
   - Ask the user: "Choose execution mode:"
     - Option A: "Sequential (safe, one task at a time)"
     - Option B: "Wave parallel ({total_tasks} tasks in {total_waves} waves)" — use actual counts from `wave_groups`
   - Proceed based on user choice.

---

## Step 3A: Sequential Mode

Use the Task tool to invoke ONE mysd-executor agent with the full context:

```
Task: Invoke mysd-executor agent for single-agent execution
Agent: mysd-executor
Context: {
  change_name: {change_name},
  must_items: {must_items},
  should_items: {should_items},
  may_items: {may_items},
  tdd_mode: {tdd_mode},
  atomic_commits: {atomic_commits},
  execution_mode: "sequential",
  pending_tasks: {pending_tasks}
}
```

The executor agent will handle all pending tasks sequentially with the mandatory alignment gate.

Skip to Step 4.

---

## Step 3B: Wave Mode

Process each wave in `wave_groups` in order. **CRITICAL: Do NOT start the next wave until the current wave's merge loop is fully complete.**

For each `wave_index`, `wave` in `wave_groups`:

### 3B-1: Announce Wave

Print:
```
Wave {wave_index+1}/{total_waves}: T{task_ids_in_wave} executing in parallel...
```
(e.g., "Wave 1/3: T1, T2, T3 executing in parallel...")

### 3B-2: Create Worktrees

For each task in the current wave, run:
```
mysd worktree create {task.id} "{task.name}"
```

Parse the JSON output:
```json
{"path": ".worktrees/T1", "branch": "mysd/my-feature/T1-setup-auth"}
```

Record `path` and `branch` for each task. If worktree creation fails for a task, skip that task in this wave (mark as failed, do not block other tasks).

### 3B-3: Spawn Parallel Executors

Use the Task tool to spawn ONE executor agent PER task in the current wave (all in parallel):

For each task in the wave:
```
Task: Invoke mysd-executor agent for wave task T{task.id}
Agent: mysd-executor
Context: {
  change_name: {change_name},
  must_items: {must_items},
  should_items: {should_items},
  may_items: {may_items},
  tdd_mode: {tdd_mode},
  atomic_commits: {atomic_commits},
  execution_mode: "wave",
  assigned_task: {task},
  worktree_path: {path},
  branch: {branch},
  isolation: "worktree"
}
```

**Wait for ALL executor agents to complete before proceeding.** Do NOT abort the wave if one task fails — continue-on-failure policy. Other tasks in the same wave must complete regardless.

### 3B-4: Merge Loop (ascending task ID order)

After ALL executors in this wave have completed (success or failure), run the merge loop.

Sort completed tasks by task ID in **ascending order** (deterministic merge order, FEXEC-06).

For each task in ascending ID order:

**If task succeeded:**

Run:
```
git merge --no-ff {branch}
```
(from the repo root, NOT the worktree directory)

If merge is clean:
- Run: `mysd worktree remove {task.id} "{branch}"`
- Print: `T{id} ✓ merged and cleaned up`

If merge conflict detected:
- Retry up to **3 attempts** (FEXEC-07). For each attempt:
  1. Resolve the conflict markers in the conflicting files (edit files to remove `<<<<<<<`, `=======`, `>>>>>>>` markers, keeping the correct content from both branches)
  2. Run: `git add {resolved_files}`
  3. Run: `git commit` (complete the merge)
  4. Run: `go build ./...` — must succeed
  5. Run: `go test ./...` — must pass
  6. If build or tests fail: run `git merge --abort`, undo the merge, start next attempt from step 1

If all 3 attempts fail:
- Print error:
  ```
  Merge failed for T{id} after 3 attempts: {reason}
  Worktree preserved at: {worktree_path}
  Branch: {branch}
  To resolve manually: cd {worktree_path} && git status
  ```
- Do NOT run `mysd worktree remove` (preserve the worktree for manual resolution)
- Continue to the next task in the merge loop (do not abort the wave)

**If task failed (executor reported failure):**
- Print: `T{id} ✗ failed — worktree preserved at {worktree_path} on branch {branch}`
- Do NOT attempt merge
- Do NOT run `mysd worktree remove`

### 3B-5: Wave Summary

After the merge loop completes for this wave:
```
Wave {wave_index+1} complete: {succeeded} succeeded, {failed} failed
```

**CRITICAL: Do NOT create worktrees for the next wave until this wave's merge loop is fully complete.** Creating worktrees before merges complete means those branches fork from an outdated HEAD, creating unnecessary merge conflicts.

---

## Step 4: Post-Execution Summary

After all waves (or sequential execution) complete, print a summary:

```
## Execution Summary

Tasks completed: {count}
Tasks failed: {count}

{If any failed tasks}
Preserved worktrees (require manual resolution):
- T{id}: {worktree_path} (branch: {branch})
  Reason: {failure reason}

Next steps:
- Run `mysd status` to check overall progress
- For preserved worktrees: resolve conflicts manually, then run `mysd execute` again
```

If all tasks completed successfully with no preserved worktrees:
```
All {count} tasks completed successfully.
Run `mysd status` to check overall progress.
```
