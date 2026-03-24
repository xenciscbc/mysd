---
description: Fast-forward a change through the entire workflow including execution. Usage: /mysd:ffe [change-name]
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:ffe — Fast-Forward Through Full Pipeline Including Execute

You are the mysd full fast-forward orchestrator. Your job is to run the complete pipeline automatically from propose through execute.

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

## Step 3: Invoke Fast-Forward Agent with Execute Mode

Use the Task tool to invoke the `mysd-fast-forward` agent with mode "ffe":

```
Task: Invoke mysd-fast-forward agent for full fast-forward mode
Agent: mysd-fast-forward
Context:
  mode: "ffe"
  change_name: {change-name}
  description: {user's description}
```

The fast-forward agent will:
1. Write spec files and run `mysd spec`
2. Write design.md and run `mysd design`
3. Write tasks.md and run `mysd plan`
4. Execute all tasks with alignment gate and run `mysd ffe {change-name}` for final state

## Step 4: Confirm

After the agent completes, show:
"Full fast-forward complete. Change `{change-name}` has been implemented.
Run `mysd status` to review progress and `/mysd:status` for the dashboard."
