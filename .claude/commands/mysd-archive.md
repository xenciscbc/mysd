---
model: claude-sonnet-4-5
description: Archive a verified spec change to .specs/archive/.
argument-hint: ""
allowed-tools:
  - Bash
  - Read
---

# /mysd:archive — Archive Verified Change

You are the mysd archive orchestrator. Your job is to archive a fully verified spec change.

## Step 1: Archive the Change

Run:
```
mysd archive
```

**If the command succeeds:**

Inform the user:
```
Change archived successfully.

Archive location: .specs/archive/{change_name}/
All spec files, verification report, and history preserved.

Next step: Start your next change with `/mysd:propose`.
```

**If the command fails with "must be verified" or "not in verified state":**

```
Archive blocked: This change has not been fully verified.

All MUST requirements must pass verification before archiving.

Next step: Run `/mysd:verify` to verify this change first.
```

**If the command fails with "not done" or "MUST items not complete":**

```
Archive blocked: Not all MUST requirements are marked as DONE.

Some requirements may have failed verification or were not executed.

Next steps:
1. Run `/mysd:verify` to see which items failed
2. Run `/mysd:execute` to fix the failing items
3. Run `/mysd:verify` again to re-verify
4. Then run `/mysd:archive` once all MUST items pass
```

**If the command fails with any other error:**

Display the error message and suggest:
```
If the problem persists, check:
- Run `mysd status` to see the current change state
- Ensure you are in a project directory initialized with `mysd init`
```

## Step 2: Post-Archive Guidance

After a successful archive, remind the user of the workflow:
```
Your change is now archived and the spec lifecycle is complete.

To start a new feature or change:
  /mysd:propose — Create a new spec proposal
```
