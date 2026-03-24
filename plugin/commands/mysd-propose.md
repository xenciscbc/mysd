---
description: Create a new spec change with proposal scaffolding. Usage: /mysd:propose [change-name]
allowed-tools:
  - Bash
  - Read
  - Write
---

# /mysd:propose — Create a New Change Proposal

You are the mysd propose assistant. Your job is to scaffold a new change and fill in the proposal.

## Step 1: Get Change Name

If `$ARGUMENTS` is provided, use it as the change name. Otherwise ask:
"What is the name for this change? (use kebab-case, e.g. `add-user-auth`)"

Also ask for a brief description if not already provided:
"Briefly describe what this change does and why it's needed."

## Step 2: Scaffold the Change

Run:
```
mysd propose {change-name}
```

This creates `.specs/changes/{change-name}/` with a template `proposal.md`.

## Step 3: Fill in the Proposal

Read the scaffolded proposal file:
```
.specs/changes/{change-name}/proposal.md
```

Fill in the following sections based on the user's description:
- **Summary**: 1-2 sentence description of the change
- **Motivation**: Why is this change needed? What problem does it solve?
- **Scope**: What is in scope? What is explicitly out of scope?
- **Success Criteria**: How will we know the change is complete?

Write the filled proposal back to the file.

## Step 4: Confirm and Guide Next Steps

Show the user:
1. The proposal file path
2. A brief summary of what was written
3. Next step: "Run `/mysd:spec` to define detailed requirements"
