---
model: claude-sonnet-4-5
description: Archive a verified spec change to .specs/archive/.
argument-hint: "[--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:archive — Archive Verified Change

You are the mysd archive orchestrator. Your job is to archive a fully verified spec change and update documentation afterward.

## Step 0: Parse Arguments and Read Config (D-18)

Check `$ARGUMENTS` for `--auto`. Set `auto_mode` = true if `--auto` is present, false otherwise.

Read docs_to_update config:
```
mysd execute --context-only
```
Parse JSON output. Extract `docs_to_update` array and `change_name`.
If `docs_to_update` is null or empty: set `has_docs_to_update` = false.
Otherwise: set `has_docs_to_update` = true, store the file paths list.

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
```

Proceed to Step 2.

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

## Step 2: Doc Maintenance (D-11, D-11b, D-13, D-14)

If `has_docs_to_update` is false:
  Skip this step. Go to Step 3.

### Step 2a: Confirm with User (D-13)

If `auto_mode` is true: skip confirmation, proceed directly to Step 2b.

If `auto_mode` is false:
```
The following files are configured for post-archive update:
{list each file path, one per line}

Press Enter to proceed, or type 'n' to skip doc updates.
```
If user types 'n': skip to Step 3.

### Step 2b: Read Context for Updates (D-11b)

For the archived change, read:
- `.specs/archive/{change_name}/proposal.md` — what and why
- `.specs/archive/{change_name}/tasks.md` — what was done
- `.specs/archive/{change_name}/specs/` — all spec files (MUST/SHOULD/MAY requirements)

Combine into `update_context`.

### Step 2c: Update Each File (D-11, D-11b)

For each file path in `docs_to_update`:

1. Read the current file content (if file exists; if not, treat as new file)

2. Determine update strategy based on filename:
   - If filename is `CHANGELOG.md` (case-insensitive):
     **Prepend strategy** — Generate ONLY the new changelog entry for this change.
     The entry should include: date, change name, summary of what was done (from proposal + tasks).
     Use Edit tool to prepend the new entry at the top of the file, below any existing header.
     Do NOT rewrite or modify existing changelog entries.

   - If filename is `README.md` (case-insensitive):
     **Full rewrite strategy** — Read the entire current README.md content.
     Rewrite the full file incorporating information from the archived change.
     Preserve the overall structure and tone. Update sections affected by the change.

   - For any other file:
     **Auto-detect strategy** — Read the file content. Based on the filename and content structure,
     determine the best approach (usually full rewrite). Apply the update.

3. After updating, show: "Updated: {file_path}"

### Step 2d: Summary

Show: "Updated {N} documentation file(s) based on archived change `{change_name}`."

## Step 3: Post-Archive Guidance

After a successful archive, remind the user of the workflow:
```
Your change is now archived and the spec lifecycle is complete.

To start a new feature or change:
  /mysd:propose — Create a new spec proposal
```
