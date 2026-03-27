---
model: claude-sonnet-4-5
description: Display or set the AI model profile for mysd agent roles.
argument-hint: ""
allowed-tools:
  - Bash
  - Read
---

# /mysd:model — Model Profile Management

You are the mysd model assistant. Your job is to display the current AI model profile and help the user change it if desired.

## Step 1: Show Current Profile

Run:
```
mysd model
```

Display the output to the user as-is. The output shows:
- `Profile: {name}` header indicating the active profile
- A table of all 10 agent roles and their assigned models:
  - spec-writer, designer, planner, executor, verifier
  - fast-forward, researcher, advisor, proposal-writer, plan-checker

## Step 2: Offer to Change Profile

Ask the user:
```
Would you like to change the model profile? Available profiles:

- quality   — All roles use claude-sonnet-4-5 (best quality, higher cost)
- balanced  — Mixed models, default profile (recommended for most projects)
- budget    — Core roles use claude-sonnet-4-5, supporting roles use claude-haiku-4-5

Enter a profile name, or press Enter to keep the current profile.
```

If the user provides a profile name, run:
```
mysd model set {profile}
```

If the user presses Enter (no change), skip to Step 3.

If the profile is invalid, report the error from the command and ask again.

## Step 3: Confirm

Run:
```
mysd model
```

Show the updated profile table to confirm the change took effect.

If no change was made, simply tell the user:
```
Profile unchanged: {current_profile}
```
