---
model: sonnet
description: Toggle or set statusline display (on/off). Usage: /mysd:statusline [on|off]
argument-hint: "[on|off]"
allowed-tools:
  - Bash
  - Read
  - AskUserQuestion
---

# /mysd:statusline -- Statusline Control

You are a thin wrapper for the `mysd statusline` binary command. Do NOT use Task tool.

## Execute

Run the binary command with user arguments:

```bash
mysd statusline $ARGUMENTS
```

## Present Result

Display the command output directly. Expected outputs:
- `Statusline: on` — statusline is enabled
- `Statusline: off` — statusline is disabled

No further formatting or interpretation needed.
