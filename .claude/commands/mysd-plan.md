---
model: claude-sonnet-4-5
description: Break design into an executable task list. Invokes mysd-planner agent. Usage: /mysd:plan [--research] [--check]
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:plan — Create Executable Task List

You are the mysd plan orchestrator. Your job is to gather context and invoke the planner agent.

## Step 1: Parse Options

Check `$ARGUMENTS` for optional flags:
- `--research`: Enable research phase before planning (deeper analysis)
- `--check`: Enable plan check/validation phase after planning

Build the command accordingly.

## Step 2: Get Execution Context

Run with appropriate flags:
```
mysd plan --context-only [--research] [--check]
```

Parse the JSON output. It contains:
- `change_name`: The current change being worked on
- `phase`: Current workflow phase
- `specs`: Array of requirements
- `design`: Design document body
- `model`: The model to use for planning
- `research_enabled`: Whether research mode is active
- `check_enabled`: Whether plan validation is active
- `test_generation`: Whether to generate tests post-execution

If this returns an error (e.g., not in designed phase), guide the user to complete `/mysd:design` first.

## Step 3: Invoke Planner Agent

Use the Task tool to invoke the `mysd-planner` agent with the full context JSON:

```
Task: Invoke mysd-planner agent
Agent: mysd-planner
Context: {context JSON from Step 2}
```

Pass the entire context so the agent has specs, design, and all configuration flags.

## Step 4: State Transition

After the agent completes task planning, run:
```
mysd plan
```

This transitions the workflow state to `planned`.

## Step 5: Confirm

Show the user: "Planning complete. State transitioned to planned. Next: `/mysd:execute`"
