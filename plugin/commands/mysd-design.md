---
model: claude-sonnet-4-5
description: Capture technical decisions and architecture for the current change. Invokes mysd-designer agent.
argument-hint: ""
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:design — Capture Technical Design

You are the mysd design orchestrator. Your job is to gather context and invoke the designer agent.

## Step 1: Get Execution Context

Run:
```
mysd design --context-only
```

Parse the JSON output. It contains:
- `change_name`: The current change being worked on
- `phase`: Current workflow phase
- `proposal_summary`: Summary of the proposal
- `specs`: Array of requirements with RFC 2119 keywords
- `model`: The model to use for design

If this returns an error (e.g., not in specced phase), guide the user to complete `/mysd:spec` first.

## Step 2: Invoke Designer Agent

Use the Task tool to invoke the `mysd-designer` agent with the full context JSON:

```
Task: Invoke mysd-designer agent
Agent: mysd-designer
Context: {context JSON from Step 1}
```

Pass the entire context so the agent has change name, specs, and model preference.

## Step 3: State Transition

After the agent completes design documentation, run:
```
mysd design
```

This transitions the workflow state to `designed`.

## Step 4: Confirm

Show the user: "Design complete. State transitioned to designed. Next: `/mysd:plan`"
