---
description: Fast-forward a change through propose, spec, design, and plan in one command. Usage: /mysd:ff [change-name]
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:ff — Fast-Forward Through Planning

You are the mysd fast-forward orchestrator. Your job is to run the full planning pipeline automatically.

## Step 1: Get Change Name

Get the change name from `$ARGUMENTS`. If not provided, ask:
"What is the name for this change? (use kebab-case, e.g. `add-user-auth`)"

Also ask for a brief description of the change.

## Step 2: Scaffold the Change

Run:
```
mysd propose {change-name}
```

This creates the change directory and sets state to `proposed`.

## Step 3: Invoke Fast-Forward Agent

Use the Task tool to invoke the `mysd-fast-forward` agent with mode "ff":

```
Task: Invoke mysd-fast-forward agent for fast-forward mode
Agent: mysd-fast-forward
Context:
  mode: "ff"
  change_name: {change-name}
  description: {user's description}
```

The fast-forward agent will:
1. Write spec files and run `mysd spec`
2. Write design.md and run `mysd design`
3. Write tasks.md and run `mysd plan`

It stops at the `planned` state — ready for execution review.

## Step 4: Confirm

After the agent completes, show:
"Fast-forward complete. Change `{change-name}` is now at planned state.
Review the generated artifacts, then run `/mysd:execute` to implement."
