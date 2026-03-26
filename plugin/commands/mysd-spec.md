---
model: claude-sonnet-4-5
description: Write detailed requirements with optional focused research. Invokes mysd-spec-writer agent. Usage: /mysd:spec [--auto]
argument-hint: "[--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:spec — Write Detailed Requirements

You are the mysd spec orchestrator. Your job is to gather context, optionally run focused research, and invoke the spec-writer agent.

## Step 1: Get Execution Context

Check `$ARGUMENTS` for `--auto`. Set `auto_mode` = true if present.

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

## Step 2: Optional Research (DISC-02, DISC-04)

If `auto_mode` is true: skip research. Go to Step 3.

If `auto_mode` is false:
  Ask: "Would you like to run focused research on how to implement this spec? [y/N]"

If user chooses research:
  Spawn ONE `mysd-researcher` agent:

  Task: Research implementation approach for {change_name} spec
  Agent: mysd-researcher
  Context: {
    "change_name": "{change_name}",
    "dimension": "codebase",
    "topic": "how to implement the requirements in {change_name} — focus on existing codebase patterns and integration points",
    "spec_files": ["{proposal.md path}"],
    "auto_mode": false
  }

  Present research findings to user.
  These findings become additional context for the spec writer in Step 3.

## Step 3: Invoke Spec Writer Agent

Use the Task tool to invoke the `mysd-spec-writer` agent with the full context JSON:

```
Task: Invoke mysd-spec-writer agent
Agent: mysd-spec-writer
Context: {
  ...context JSON from Step 1,
  "research_findings": "{from Step 2, or empty if no research}"
}
```

Pass the entire context to the agent so it has change name, proposal body, model preference, and any research findings.

## Step 4: State Transition

After the agent completes spec writing, run:
```
mysd spec
```

This transitions the workflow state to `specced`.

## Step 5: Confirm

Show the user: "Spec writing complete. State transitioned to specced. Next: `/mysd:design`"
