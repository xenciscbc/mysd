---
model: claude-sonnet-4-5
description: Create a new spec change with proposal scaffolding. Supports source auto-detection. Usage: /mysd:propose [change-name|file-path|dir-path] [--auto]
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:propose — Create a New Change Proposal

You are the mysd propose orchestrator. Your job is to scaffold a new change, detect the input source, and invoke the proposal writer agent.

## Step 1: Parse Arguments

Check `$ARGUMENTS` for `--auto`. Remove it from the arguments list.
Set `auto_mode` = true if `--auto` is present, false otherwise.

The remaining arguments (after removing `--auto`) are the `source_arg`.

## Step 2: Source Detection

Apply the following priority order to determine the input source and change name:

**Priority 1:** If `source_arg` matches a directory `.specs/changes/{source_arg}/`
→ Use `source_arg` as the change name (mysd change mode)
→ Read existing proposal.md if present as initial content

**Priority 2:** If `source_arg` is a file path (ends with `.md` or file exists on disk)
→ Single file mode: read the file as initial content
→ Derive change name from filename (strip extension, kebab-case)

**Priority 3:** If `source_arg` is a directory path (directory exists on disk)
→ Selection mode: list all `.md` files in the directory
→ If `auto_mode` is true: use all files as initial content
→ If `auto_mode` is false: present list and let user multi-select

**Priority 4:** If no `source_arg` and there is an active change (check `mysd status` output)
→ Use the current active change

**Priority 5:** If no `source_arg` and no active change → auto-detect from known sources:
→ Check `~/.gstack/projects/{project}/` for `.md` files (design docs, test plans, etc.)
→ Check conversation context for mentioned plan documents or design files
→ Do NOT check `.claude/plans/` (hash filenames have no project info)
→ If `auto_mode` is true: use first detected source
→ If `auto_mode` is false: present detected sources and let user choose

**Priority 6:** If nothing found
→ If `auto_mode` is true: auto-generate change name from conversation context
→ If `auto_mode` is false: ask user for change name and brief description

## Step 3: Scaffold the Change

Run:
```
mysd propose {change-name}
```

This creates `.specs/changes/{change-name}/` with a template `proposal.md`.

If source content was detected in Step 2 (file/directory mode), read that content now.

## Step 4: Invoke Proposal Writer

Use the Task tool to invoke `mysd-proposal-writer`:

```
Task: Write proposal for {change_name}
Agent: mysd-proposal-writer
Context: {
  "change_name": "{change_name}",
  "conclusions": "{source content or user description}",
  "existing_proposal": null,
  "auto_mode": {auto_mode}
}
```

The proposal writer will fill in the proposal.md with structured content based on the source material.

## Step 5: Confirm

Show the user:
1. The proposal file path: `.specs/changes/{change_name}/proposal.md`
2. A brief summary of what was written
3. Next step: "Run `/mysd:spec` to define detailed requirements, or `/mysd:plan` if specs are ready"
