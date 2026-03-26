---
model: claude-sonnet-4-5
description: Fast-forward through plan + apply + archive. Assumes spec is ready. No research. Implies --auto. Usage: /mysd:ff [change-name]
argument-hint: "[change-name]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:ff -- Fast-Forward (Plan + Apply + Archive)

You are the mysd fast-forward orchestrator. Run the pipeline: plan (no research) -> apply -> archive. Auto mode is always on.

## Step 1: Get Change Name

Get change name from `$ARGUMENTS`. If not provided, check `mysd status` for active change. If none, ask user.

Set `auto_mode = true` (always, per D-19/FAUTO-03).

## Step 2: Plan Phase (no research, per FAUTO-04)

Run: `mysd plan --context-only`
Parse JSON.

Spawn designer:
  Task: Create design for {change_name} (ff mode)
  Agent: mysd-designer
  Context: { "change_name": "...", "specs": [...], "research_findings": [], "auto_mode": true }

Run: `mysd design`

Spawn planner:
  Task: Create task list for {change_name} (ff mode)
  Agent: mysd-planner
  Context: { full context JSON, "auto_mode": true }

Run: `mysd plan`

## Step 3: Apply Phase

Run: `mysd execute --context-only`
Parse JSON (tasks, pending_tasks, wave_groups, etc.).

Execute tasks using the same logic as /mysd:apply Step 3:
- Single mode: sequential per-task spawn of mysd-executor with auto_mode: true
- Wave mode: parallel per-task spawn with worktree isolation, auto_mode: true

Run: `mysd execute` (state transition)

## Step 4: Archive

Run: `mysd archive`

## Step 5: Confirm

Show: "Fast-forward complete. Change `{change_name}` has been planned, executed, and archived."
