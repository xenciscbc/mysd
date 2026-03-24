---
model: claude-sonnet-4-5
description: Write detailed requirements for the current change using RFC 2119 keywords. Invokes mysd-spec-writer agent.
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:spec — Write Detailed Requirements

You are the mysd spec orchestrator. Your job is to gather context and invoke the spec-writer agent.

## Step 1: Get Execution Context

Run:
```
mysd spec --context-only
```

Parse the JSON output. It contains:
- `change_name`: The current change being worked on
- `phase`: Current workflow phase
- `proposal`: The proposal body content
- `model`: The model to use for spec writing

If this returns an error (e.g., no spec directory), guide the user to run `/mysd:propose` first.

## Step 2: Invoke Spec Writer Agent

Use the Task tool to invoke the `mysd-spec-writer` agent with the full context JSON:

```
Task: Invoke mysd-spec-writer agent
Agent: mysd-spec-writer
Context: {context JSON from Step 1}
```

Pass the entire context to the agent so it has change name, proposal body, and model preference.

## Step 3: State Transition

After the agent completes spec writing, run:
```
mysd spec
```

This transitions the workflow state to `specced`.

## Step 4: Confirm

Show the user: "Spec writing complete. State transitioned to specced. Next: `/mysd:design`"
