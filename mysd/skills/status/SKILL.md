---
model: sonnet
description: Show current workflow status dashboard with task progress and next step recommendation.
allowed-tools:
  - Bash
  - Read
---

# /mysd:status — Show Workflow Status

You are the mysd status assistant. Display the current workflow status dashboard with workflow stage indicator, task list, and next step recommendation.

## Step 1: Get Status

Run:
```
mysd status
```

Capture the output. If this returns an error (e.g., no spec directory found), inform the user:
"No active change found. Start with `/mysd:propose` to create a new change."

## Step 2: Get Task Details

Run:
```
mysd execute --context-only
```

Parse the JSON output to get the full task list with statuses. If this fails (not yet in planned phase), use only the basic status from Step 1.

## Step 3: Display Enhanced Dashboard

### Workflow Stage Indicator

Show the workflow stages with a position indicator marking the current stage:

```
Workflow: propose > plan > apply > archive
                    ^^^^ (current)
```

Determine the current stage from the `phase` field in the status output:
- `proposed` or `specced` → current stage is `propose`
- `designed` → current stage is `plan` (design is now part of plan pipeline)
- `planned` → current stage is `apply`
- `executed` → current stage is `archive`
- `verified` → current stage is `archive`
- `archived` → workflow complete

### Task List with Status Symbols

If task details are available, show each task with its status symbol:

```
Tasks:
  T1 (setup-models) .............. done
  T2 (create-api) ................ in_progress
  T3 (add-validation) ............ pending
  T4 (write-tests) ............... pending
```

Status symbol mapping:
- `done` — task completed successfully
- `failed` — task failed
- `skipped` — task skipped (dependency not met)
- `pending` — task not yet started
- `in_progress` — task currently running

### Progress Summary

Show task completion counts:
```
Progress: {completed}/{total} tasks done
```

### Next Step Recommendation

Show the recommended next command based on current stage:

```
Next: /mysd:apply
```

Mapping:
- Stage `propose` (in proposed phase) → `Next: /mysd:plan`
- Stage `plan` (in designed phase) → `Next: /mysd:plan`
- Stage `apply` → `Next: /mysd:apply`
- Stage `archive` (executed/verified) → `Next: /mysd:archive`
- Workflow complete → "Change is archived. Start a new change with `/mysd:propose`"

### Deferred Notes Count

Run: `mysd note`
Parse the output to count notes (each line starting with `[` is a note).

If notes exist, show at the bottom:
```
Deferred notes: {N} — run /mysd:note to browse
```

If no notes exist, do not show this line.
