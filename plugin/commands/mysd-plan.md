---
model: claude-sonnet-4-5
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
- `research_enabled`, `check_enabled`, `test_generation`

If error (not in designed/specced phase), guide user to complete prerequisites.

## Step 3: Research Phase (if research_enabled)

If `research_enabled` is true (from context JSON) or `--research` flag present:

  If `auto_mode` is false:
    Ask: "Would you like to run focused research on implementation details? [y/N]"
    If user declines: skip to Step 4.

  If `auto_mode` is true:
    Skip research entirely. Go to Step 4.

  Spawn ONE `mysd-researcher` agent (single, NOT parallel):

  Task: Research implementation details for {change_name}
  Agent: mysd-researcher
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

Use the Task tool to invoke `mysd-designer`:

Task: Create design document for {change_name}
Agent: mysd-designer
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

Use the Task tool to invoke `mysd-planner`:

Task: Create task list for {change_name}
Agent: mysd-planner
Context: {full context JSON from Step 2, plus research_findings and design content, plus auto_mode}

After planner completes, run state transition:
```
mysd plan
```

## Step 6: Plan Check (if check_enabled)

If `check_enabled` is true:

Run:
```
mysd plan --check --context-only
```

Use the Task tool to invoke `mysd-plan-checker`:

Task: Validate plan coverage for {change_name}
Agent: mysd-plan-checker
Context: {check output JSON}

## Step 7: Confirm

Show: "Planning complete. Pipeline: {research if enabled} -> design -> plan {-> check if enabled}. Next: `/mysd:apply`"
