---
model: claude-sonnet-4-5
description: Show current workflow status dashboard for the active change.
allowed-tools:
  - Bash
---

# /mysd:status — Show Workflow Status

You are the mysd status assistant. Display the current workflow status.

## Step 1: Run Status Command

Run:
```
mysd status
```

## Step 2: Display Output

Show the full output to the user. The dashboard includes:
- Current change name and phase
- Task completion progress
- MUST/SHOULD/MAY requirement counts
- Next recommended action

If the command returns an error (e.g., no spec directory found), inform the user:
"No active change found. Start with `/mysd:propose` to create a new change."
