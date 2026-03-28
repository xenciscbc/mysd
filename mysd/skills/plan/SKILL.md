---
description: Plan orchestrator. Optional single-agent research, then design, then task planning. Usage: /mysd:plan [--research] [--check] [--auto]
argument-hint: "[--research] [--check] [--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:plan -- Create Executable Task List

You are the mysd plan orchestrator. Your job is to run the planning pipeline: optional research, then design, then task planning.

## Step 1: Parse Arguments

Check `$ARGUMENTS` for flags:
- `--research`: Enable focused research before design
- `--check`: Enable plan-checker validation after planning
- `--auto`: Auto mode — skip interactive prompts, use AI recommendations

Set `auto_mode` = true if `--auto` is present, false otherwise.

## Step 2: Get Planning Context

Run:
```
mysd plan --context-only [--research] [--check]
```

Parse the JSON output. It contains:
- `change_name`, `phase`, `specs`, `design`, `model`
- `reviewer_model`, `plan_checker_model`
- `research_enabled`, `check_enabled`, `test_generation`

Extract the following fields:
- `model`: profile-resolved model for designer and planner agents
- `reviewer_model`: profile-resolved model for mysd-reviewer
- `plan_checker_model`: profile-resolved model for mysd-plan-checker

If error (not in designed/specced phase), guide user to complete prerequisites.

## Step 3: Research Phase (if research_enabled)

If `research_enabled` is true (from context JSON) or `--research` flag present:

  If `auto_mode` is false:
    Ask: "Would you like to run focused research on implementation details? [y/N]"
    If user declines: skip to Step 4.

  If `auto_mode` is true:
    Skip research entirely. Go to Step 4.

  Show: "Spawning mysd-researcher ({model})..."
  Spawn ONE `mysd-researcher` agent (single, NOT parallel), with `model` parameter set to `{model}`:

  Task: Research implementation details for {change_name}
  Agent: mysd-researcher
  Model: {model}
  Context: {
    "change_name": "{change_name}",
    "dimension": "architecture",
    "topic": "implementation of {change_name} — validate technical feasibility and supplement implementation details",
    "spec_files": [{all spec file paths from context + design.md path}],
    "auto_mode": {auto_mode}
  }

  Collect research output. This becomes additional input for the designer in Step 4.

  If `auto_mode` is false:
    Present research summary and ask: "Research complete. Proceed to design? (Y/n)"

## Step 4: Design Phase

Show: "Spawning mysd-designer ({model})..."
Use the Task tool to invoke `mysd-designer` with `model` parameter set to `{model}`:

Task: Create design document for {change_name}
Agent: mysd-designer
Model: {model}
Context: {
  "change_name": "{change_name}",
  "specs": [{spec content}],
  "research_findings": [{from Step 3, or empty if no research}],
  "auto_mode": {auto_mode}
}

The designer produces `design.md`.

After designer completes, run state transition:
```
mysd design
```

## Step 5: Planning Phase

Show: "Spawning mysd-planner ({model})..."
Use the Task tool to invoke `mysd-planner` with `model` parameter set to `{model}`:

Task: Create task list for {change_name}
Agent: mysd-planner
Model: {model}
Context: {full context JSON from Step 2, plus research_findings and design content, plus auto_mode}

After planner completes, run state transition:
```
mysd plan
```

## Step 5b: Invoke Reviewer

Run:
```
mysd validate {change_name}
```
Capture output as `validate_output` (empty string if command not found or fails).

Show: "Spawning mysd-reviewer ({reviewer_model})..."
Use the Task tool to invoke `mysd-reviewer` with `model` parameter set to `{reviewer_model}`:

Task: Review artifacts for {change_name}
Agent: mysd-reviewer
Model: {reviewer_model}
Context: {
  "change_name": "{change_name}",
  "phase": "plan",
  "validate_output": "{validate_output}",
  "auto_mode": {auto_mode}
}

Collect the reviewer summary. Include it in Step 7 output.

## Step 6: Plan Check (if check_enabled)

If `check_enabled` is true:

Run:
```
mysd plan --check --context-only
```

Show: "Spawning mysd-plan-checker ({plan_checker_model})..."
Use the Task tool to invoke `mysd-plan-checker` with `model` parameter set to `{plan_checker_model}`:

Task: Validate plan coverage for {change_name}
Agent: mysd-plan-checker
Model: {model}
Context: {check output JSON}

## Step 7: Confirm

Show:
1. "Planning complete. Pipeline: {research if enabled} -> design -> plan -> reviewer {-> check if enabled}."
2. Reviewer summary from Step 5b (issues fixed, cannot-auto-fix items if any)
3. Next: `/mysd:apply`
