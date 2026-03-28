---
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

Resolve model: run `mysd model`, parse profile, determine model short name (quality/balanced → sonnet, budget → haiku/sonnet per role).

Show: "Spawning 4 mysd-researcher agents ({model})..."
Spawn 4 `mysd-researcher` agents in parallel, each with `model` parameter set to `{model}`:

For each dimension in ["codebase", "domain", "architecture", "pitfalls"]:
  Task: Research {dimension} for {change_name} (ffe mode)
  Agent: mysd-researcher
  Model: {model}
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
Parse JSON. Extract `model` field.

Show: "Spawning mysd-designer ({model})..."
Spawn designer with `model` parameter set to `{model}`:
  Task: Create design for {change_name} (ffe mode)
  Agent: mysd-designer
  Model: {model}
  Context: { "change_name": "...", "specs": [...], "research_findings": [{from Step 2}], "auto_mode": true }

Run: `mysd design`

Show: "Spawning mysd-planner ({model})..."
Spawn planner with `model` parameter set to `{model}`:
  Task: Create task list for {change_name} (ffe mode)
  Agent: mysd-planner
  Model: {model}
  Context: { full context JSON, "auto_mode": true }

Run: `mysd plan`

## Step 4: Apply Phase

Same as /mysd:ff Step 3 — execute all tasks with auto_mode: true.

Run: `mysd execute --context-only`
Parse JSON. Extract `model` field.

Execute tasks using the same logic as /mysd:apply Step 3, passing `model` to each executor:
- Single mode: sequential per-task spawn of mysd-executor with auto_mode: true, model: {model}
- Wave mode: parallel per-task spawn with worktree isolation, auto_mode: true, model: {model}
- Show "Spawning mysd-executor ({model})..." before each spawn

Run: `mysd execute` (state transition)

## Step 5: Inline Auto-Verify (D-17a)

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

Show: "Spawning mysd-verifier ({model})..."
Use Task tool to invoke `mysd-verifier` with `model` parameter set to `{model}`:
```
Task: Verify spec coverage for {change_name} (ffe auto-verify)
Agent: mysd-verifier
Model: {model}
Context: {
  "change_name": "{change_name}",
  "must_items": [...],
  "should_items": [...],
  "may_items": [...]
}
```

If MUST items fail: display results, show "Verification failed. Archive skipped." STOP.
If all pass: proceed to archive.

## Step 6: Archive

Run: `mysd archive`

## Step 7: Inline Docs Update (D-17b)

Read docs_to_update config:
```
mysd execute --context-only
```
Parse JSON, extract `docs_to_update`.

If `docs_to_update` is null or empty: skip to Step 8.

For each file in `docs_to_update`:
1. Read `.specs/archive/{change_name}/proposal.md`, `tasks.md`, and `specs/` as context
2. Read current file content
3. Apply update strategy:
   - `CHANGELOG.md`: prepend new entry only
   - `README.md`: full rewrite
   - Others: auto-detect from content
4. Show: "Updated: {file_path}"

No user confirmation (auto_mode=true always in ffe).

## Step 8: Confirm

Show: "Full fast-forward complete. Change `{change_name}` has been researched, planned, executed, verified, and archived."

If any documentation files were updated, show: "Documentation updated: {list of updated files}"
