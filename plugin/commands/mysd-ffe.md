---
model: claude-sonnet-4-5
description: Full fast-forward with research + plan + apply + archive. Implies --auto. Usage: /mysd:ffe [change-name]
argument-hint: "[change-name]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:ffe -- Full Fast-Forward (Research + Plan + Apply + Archive)

You are the mysd full fast-forward orchestrator. Run the pipeline: 4-dimension research -> plan -> apply -> archive. Auto mode is always on.

## Step 1: Get Change Name

Get change name from `$ARGUMENTS`. If not provided, check `mysd status` for active change. If none, ask user.

Set `auto_mode = true` (always, per D-19/FAUTO-03).

## Step 2: Research Phase

Spawn 4 `mysd-researcher` agents in parallel:

For each dimension in ["codebase", "domain", "architecture", "pitfalls"]:
  Task: Research {dimension} for {change_name} (ffe mode)
  Agent: mysd-researcher
  Context: {
    "change_name": "{change_name}",
    "dimension": "{dimension}",
    "topic": "{change description from spec}",
    "spec_files": [{spec file paths}],
    "auto_mode": true
  }

Collect all 4 research outputs.

## Step 3: Plan Phase (with research findings)

Run: `mysd plan --context-only`
Parse JSON.

Spawn designer with research:
  Task: Create design for {change_name} (ffe mode)
  Agent: mysd-designer
  Context: { "change_name": "...", "specs": [...], "research_findings": [{from Step 2}], "auto_mode": true }

Run: `mysd design`

Spawn planner:
  Task: Create task list for {change_name} (ffe mode)
  Agent: mysd-planner
  Context: { full context JSON, "auto_mode": true }

Run: `mysd plan`

## Step 4: Apply Phase

Same as /mysd:ff Step 3 — execute all tasks with auto_mode: true.

Execute tasks using the same logic as /mysd:apply Step 3:
- Single mode: sequential per-task spawn of mysd-executor with auto_mode: true
- Wave mode: parallel per-task spawn with worktree isolation, auto_mode: true

Run: `mysd execute` (state transition)

## Step 5: Archive

Run: `mysd archive`

## Step 6: Confirm

Show: "Full fast-forward complete. Change `{change_name}` has been researched, planned, executed, and archived."
