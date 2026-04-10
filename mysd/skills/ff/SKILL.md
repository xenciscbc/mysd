---
description: Fast-forward through plan + apply + archive. Assumes spec is ready. No research. Implies --auto. Usage: /mysd:ff [change-name]
argument-hint: "[change-name]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
  - AskUserQuestion
---

# /mysd:ff -- Fast-Forward (Plan + Apply + Archive)

You are the mysd fast-forward orchestrator. Run the pipeline: plan (no research) -> apply -> archive. Auto mode is always on.

## Question Protocol

- Ask one question at a time. Wait for the user's answer before asking the next.
- When a question has concrete options, use the **AskUserQuestion tool** — do not list options as plain text.
- Open-ended questions may use plain text.

## Step 1: Get Change Name

Get change name from `$ARGUMENTS`. If not provided, check `mysd status` for active change. If none, ask user.

Set `auto_mode = true` (always, per D-19/FAUTO-03).

## Step 2: Plan Phase (no research, per FAUTO-04)

Run: `mysd plan --context-only`
Parse JSON. Extract `spec_dir` field. Pass `spec_dir` to all agents.

Resolve models:
```
mysd model resolve designer
mysd model resolve planner
mysd model resolve executor
mysd model resolve verifier
```
Capture as `designer_model`, `planner_model`, `executor_model`, `verifier_model`.

Show: "Spawning mysd-designer ({designer_model})..."
Spawn designer with `model` parameter set to `{designer_model}`:
  Task: Create design for {change_name} (ff mode)
  Agent: mysd-designer
  Model: {designer_model}
  Context: { "spec_dir": "{spec_dir}", "change_name": "...", "specs": [...], "research_findings": [], "auto_mode": true }

Run: `mysd design`

Show: "Spawning mysd-planner ({planner_model})..."
Spawn planner with `model` parameter set to `{planner_model}`:
  Task: Create task list for {change_name} (ff mode)
  Agent: mysd-planner
  Model: {planner_model}
  Context: { full context JSON including spec_dir, "auto_mode": true }

Run: `mysd plan`

## Step 3: Apply Phase

Run: `mysd execute --context-only`
Parse JSON.

Execute tasks using the same logic as /mysd:apply Step 3, passing `executor_model` to each executor:
- Single mode: sequential per-task spawn of mysd-executor with auto_mode: true, model: {executor_model}
- Wave mode: parallel per-task spawn with worktree isolation, auto_mode: true, model: {executor_model}
- Show "Spawning mysd-executor ({executor_model})..." before each spawn

Run: `mysd execute` (state transition)

## Step 4: Inline Auto-Verify (D-17a)

Run build and test (auto_mode=true, no confirmation):
```
go build ./...
```
If build fails: display error, show "Build failed. Archive skipped. Run `/mysd:fix` to fix." STOP.

```
go test ./...
```
If tests fail: display error, show "Tests failed. Archive skipped. Run `/mysd:fix` to fix." STOP.

If both pass, invoke verifier:
```
mysd execute --context-only
```
Parse JSON for must_items, should_items, may_items.

Show: "Spawning mysd-verifier ({verifier_model})..."
Use Task tool to invoke `mysd-verifier` with `model` parameter set to `{verifier_model}`:
```
Task: Verify spec coverage for {change_name} (ff auto-verify)
Agent: mysd-verifier
Model: {verifier_model}
Context: {
  "spec_dir": "{spec_dir}",
  "change_name": "{change_name}",
  "must_items": [...],
  "should_items": [...],
  "may_items": [...]
}
```

If MUST items fail: display results, show "Verification failed. Archive skipped." STOP.
If all pass: proceed to archive.

## Step 5: Archive

Run: `mysd archive`

## Step 6: Inline Docs Update (D-17b)

Read docs_to_update config:
```
mysd execute --context-only
```
Parse JSON, extract `docs_to_update`.

If `docs_to_update` is null or empty: skip to Step 7.

For each file in `docs_to_update`:
1. Read `{spec_dir}/archive/{change_name}/proposal.md`, `tasks.md`, and `specs/` as context
2. Read current file content
3. Apply update strategy:
   - `CHANGELOG.md`: prepend new entry only
   - `README.md`: full rewrite
   - Others: auto-detect from content
4. Show: "Updated: {file_path}"

No user confirmation (auto_mode=true always in ff).

## Step 7: Confirm

Show: "Fast-forward complete. Change `{change_name}` has been planned, executed, verified, and archived."

If any documentation files were updated, show: "Documentation updated: {list of updated files}"
