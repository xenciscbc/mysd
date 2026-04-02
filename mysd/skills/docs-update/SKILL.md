---
description: Trigger documentation updates independently. Usage: /mysd:docs-update [--change <name> | --last N | --full | "description"]
argument-hint: "[--change <name> | --last N | --full | \"description\"]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Glob
  - Grep
  - AskUserQuestion
---

# /mysd:docs-update -- Standalone Documentation Update

You are the mysd docs-update assistant. Your job is to update project documentation based on archived changes or codebase state, independently from the archive/ff/ffe workflows.

## Step 1: Parse Arguments

Check `$ARGUMENTS` to determine the update scope:

- **No arguments** → default scope (latest archived change)
- `--change <name>` → specified change scope
- `--last N` → last N archived changes scope
- `--full` → full codebase scan scope
- Any other text → free-text description scope

## Step 2: Gather Update Context

Based on the scope determined in Step 1, gather the update context.

### Default scope (no arguments)

Read the most recent archived change from `openspec/changes/archive/`. The archive directory uses date-prefixed names (`YYYY-MM-DD-<changeName>/`). Sort by directory name descending to find the most recent.

If no archived changes exist, inform the user:
```
No archived changes found in openspec/changes/archive/.
Nothing to update from. Use --full to scan the codebase instead.
```

For the most recent archived change, read:
- `openspec/changes/archive/YYYY-MM-DD-{change_name}/proposal.md` -- what and why
- `openspec/changes/archive/YYYY-MM-DD-{change_name}/tasks.md` -- what was done
- `openspec/changes/archive/YYYY-MM-DD-{change_name}/specs/` -- all spec files

Combine into `update_context`.

### `--change <name>` scope

Locate the archived change matching `<name>` in `openspec/changes/archive/`. The directory name includes a date prefix, so search for directories ending with `-<name>`.

If no match found:
```
Archived change "{name}" not found in openspec/changes/archive/.
Available archived changes:
{list available archived change directories}
```

If found, read the same files as default scope (proposal, tasks, specs) and combine into `update_context`.

### `--last N` scope

Read the N most recent archived changes from `openspec/changes/archive/`, sorted by directory name descending (date prefix ensures chronological order).

If fewer than N archived changes exist, use all available and inform the user.

For each archived change, read proposal, tasks, and specs. Combine all into `update_context`.

### `--full` scope

Scan the current codebase to build update context:
- Read existing commands: `mysd --help` output
- Read configuration options: `openspec/config.yaml`
- Read project structure: list key directories and files
- Read existing skills: `mysd/skills/*/SKILL.md` descriptions
- Read any existing docs in `docs_to_update` for current state

Combine into `update_context` that reflects the actual project state.

### Free-text description scope

Use the provided description text directly as `update_context`. No file reading needed.

## Step 3: Check docs_to_update Configuration

Run:
```
mysd docs
```

Parse the output to get the list of files in `docs_to_update`.

If the list is empty or no files are configured:
```
No files configured for post-archive doc updates.
Use `mysd docs add <path>` to add files (e.g., README.md, CHANGELOG.md).
```
Stop here.

If files are configured, proceed to Step 4.

## Step 4: Update Each File

For each file path in `docs_to_update`:

1. Read the current file content (if file exists; if not, treat as new file)

2. Determine update strategy based on filename:
   - If filename is `CHANGELOG.md` (case-insensitive):
     **Prepend strategy** -- Generate ONLY the new changelog entry based on `update_context`.
     The entry should include: date, change name/summary, what was done.
     Use Edit tool to prepend the new entry at the top of the file, below any existing header.
     Do NOT rewrite or modify existing changelog entries.

   - If filename is `README.md` (case-insensitive):
     **Full rewrite strategy** -- Read the entire current README.md content.
     Rewrite the full file incorporating information from `update_context`.
     Preserve the overall structure and tone. Update sections affected by the change.

   - For any other file:
     **Auto-detect strategy** -- Read the file content. Based on the filename and content structure,
     determine the best approach (usually full rewrite). Apply the update.

3. After updating, show: "Updated: {file_path}"

## Step 5: Summary

Show:
```
Updated {N} documentation file(s).
Scope: {scope description}
```
