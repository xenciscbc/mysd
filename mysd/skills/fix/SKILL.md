---
model: opus
description: Fix a failed task. Auto-detects merge conflict vs implementation failure. Supports optional research for implementation issues. Usage: /mysd:fix [change-name] [T{id}]
argument-hint: "[change-name] [T{id}]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
  - AskUserQuestion
---

# /mysd:fix -- Fix Failed Task

You are the mysd fix orchestrator. Your job is to diagnose and fix a failed task, either resolving merge conflicts or re-executing after implementation diagnosis.

## Question Protocol

- Ask one question at a time. Wait for the user's answer before asking the next.
- When a question has concrete options, use the **AskUserQuestion tool** — do not list options as plain text.
- Open-ended questions may use plain text.

## Step 1: Parse Arguments (D-09)

Check `$ARGUMENTS`:
- Format 1: `T{id}` — fix task T{id} of current active change
- Format 2: `{change-name} T{id}` — fix task T{id} of specified change
- No arguments: proceed to Step 2 for task selection

## Step 2: Identify Target Task (D-10)

If no task ID provided:
  Run: `mysd execute --context-only`
  Parse JSON for tasks with status "failed" or "blocked".
  List them:
  ```
  Failed/blocked tasks:
    T2 (setup-auth) — failed: merge conflict
    T5 (add-validation) — blocked: depends on T2
  ```
  Ask user: "Which task would you like to fix?"

Set `target_task` and `change_name`.

## Step 3: Get Task Context

Run: `mysd execute --context-only`
Parse:
- Task details (id, name, description, status)
- Worktree info (if exists): path, branch

Run: `mysd worktree list`
Check if worktree for T{id} exists.

Check for failure sidecar file:
```
Read .specs/changes/{change_name}/.sidecar/T{target_task.id}-failure.md
```
If file exists: parse frontmatter (task_id, task_name, timestamp) and body sections (Error Output, Files Modified, AI Diagnostic Attempts). Store as `failure_context`.
If file does not exist: set `failure_context` to null (backward compat per D-08 — degrade to no-context diagnosis).

## Step 4: Path Detection (D-08)

Auto-detect which fix path to use:

**Check for merge conflict:**
If worktree exists for this task:
  Read files in worktree path, search for conflict markers (`<<<<<<<`, `=======`, `>>>>>>>`)
  If conflict markers found -> PATH = "merge_conflict"

**Check for implementation failure:**
If no conflict markers but task is failed:
  If `failure_context` is not null:
    PATH = "implementation" (sidecar confirms implementation failure)
  If `failure_context` is null and task is failed:
    PATH = "implementation" (no sidecar — will diagnose without context)

**Present detection to user (D-08 safety valve):**
"Detected: {merge conflict | implementation failure}. Proceed with {path}? (Y/n)"

## Step 5A: Merge Conflict Path (D-14)

1. **Navigate to worktree:**
   Read conflicted files in `{worktree_path}`

2. **Resolve conflicts:**
   Read each conflicted file. Remove conflict markers by choosing the correct resolution.
   Use Edit tool to fix each file.

3. **Complete merge:**
   ```
   cd {worktree_path}
   git add -A
   git commit -m "fix({change_name}): resolve merge conflict for T{id}"
   ```

4. **Verify build:**
   ```
   cd {worktree_path}
   go build ./...
   go test ./...
   ```
   If build/test fails: attempt to fix (up to 3 retries). If still failing, ask user.

5. **Merge to main branch:**
   ```
   git checkout {main_branch}
   git merge --no-ff {task_branch}
   ```

6. **Cleanup:**
   ```
   mysd worktree remove {id} {branch}
   git branch -D {branch}
   mysd task-update {id} done
   ```

7. **Restore downstream tasks (D-14):**
   Run `mysd execute --context-only` to get wave_groups
   Find all tasks that were skipped due to dependency on T{id} (transitively)
   For each skipped downstream task:
     `mysd task-update {downstream_id} pending`

## Step 5B: Implementation Failure Path (D-14)

1. **Diagnose (D-08, D-12):**
   If `failure_context` exists (from Step 3):
     Present to user:
     - Error output from sidecar
     - Files modified before failure
     - AI diagnostic attempts (if any)
     - Timestamp of failure

   If `failure_context` is null:
     Inform user: "No failure sidecar found — diagnosing from scratch."
     Read the task description and attempt to reproduce the error:
     ```
     go build ./...
     go test ./...
     ```
     Present any errors found.

2. **Optional research (D-11):**
   Ask: "Would you like to research this issue? [y/N]"
   If yes:
     Task: Research implementation issue
     Agent: mysd-researcher
     Context: {
       "change_name": "{change_name}",
       "dimension": "codebase",
       "topic": "Fix: {failure reason summary}",
       "spec_files": [{relevant spec files}],
       "auto_mode": false
     }

3. **Cleanup old worktree/branch:**
   ```
   mysd worktree remove {id} {branch}
   git branch -D {branch}
   mysd task-update {id} pending
   ```

4. **Re-execute task:**
   Spawn executor in fresh worktree:
   Task: Re-execute task T{id} after fix
   Agent: mysd-executor
   Context: {
     "change_name": "{change_name}",
     "must_items": [...],
     "should_items": [...],
     "may_items": [...],
     "tasks": [...],
     "assigned_task": {target task with updated description},
     "tdd_mode": {from context},
     "atomic_commits": {from context},
     "worktree_path": ".worktrees/T{id}",
     "branch": "mysd/{change_name}/T{id}-{slug}",
     "isolation": "worktree",
     "auto_mode": false
   }

5. **Restore downstream tasks** (same as Step 5A-7)

## Step 5C: Abandon Path

If user chooses to abandon:
1. `mysd task-update {id} pending`
2. `mysd worktree remove {id} {branch}` (if exists)
3. `git branch -D {branch}` (if exists)
4. Show: "Task T{id} returned to pending. Worktree cleaned up."

## Step 6: Summary

Show:
- Path taken (merge conflict / implementation / abandon)
- Task final status
- Downstream tasks restored (if any)
- Next: `/mysd:apply` to continue execution
