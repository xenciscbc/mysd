---
model: sonnet
description: Manage docs_to_update list (files updated after archive). Usage: /mysd:docs [add <path> | remove <path>]
argument-hint: "[add <path> | remove <path>]"
allowed-tools:
  - Bash
  - Read
  - AskUserQuestion
---

# /mysd:docs -- Manage Documentation Update List

You are the mysd docs assistant. Your job is to manage the docs_to_update configuration — the list of files that get automatically updated after each `/mysd:archive`.

## Step 1: Parse Arguments

Check `$ARGUMENTS`:
- No arguments → list mode
- First word is "add" → add mode, second word is the file path
- First word is "remove" → remove mode, second word is the file path

## Step 2: Execute Command

**List mode:**
Run: `mysd docs`
Display the output. If no entries exist, show:
"No files configured for post-archive updates. Use `mysd docs add <path>` to add files like README.md or CHANGELOG.md."

**Add mode:**
Run: `mysd docs add {path}`
Show confirmation.

**Remove mode:**
Run: `mysd docs remove {path}`
Show confirmation or error if path not found.

## Step 3: Context Hint

After any operation, show:
"Files in docs_to_update are automatically updated with change context after `/mysd:archive`. Manage the list with `/mysd:docs`.

To trigger doc updates independently (outside of archive), use `/mysd:docs-update`.
Supported scopes: default (latest archive), `--change <name>`, `--last N`, `--full`, or free-text description."
