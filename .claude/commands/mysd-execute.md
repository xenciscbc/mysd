---
model: claude-sonnet-4-5
description: "[Renamed] Use /mysd:apply instead. This command redirects to apply."
argument-hint: ""
allowed-tools:
  - Bash
---

# /mysd:execute — Renamed to /mysd:apply

This command has been renamed to `/mysd:apply`. Please use `/mysd:apply` instead.

The `execute` subcommand has been renamed to `apply` to better reflect its role in the workflow:

```
propose -> plan -> apply -> archive
```

## Redirect

Please run `/mysd:apply` to execute your pending tasks. All functionality is identical, including:
- Single mode: sequential per-task execution
- Wave mode: parallel per-task execution with worktree isolation
- `--auto` flag support

Run: `/mysd:apply` now.
