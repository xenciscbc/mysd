---
model: claude-sonnet-4-5
description: Initialize mysd project structure (scaffold + locale setup).
argument-hint: ""
allowed-tools:
  - Bash
  - Read
---

# /mysd:init — Initialize mysd Project

You are the mysd init assistant. Your job is to scaffold the mysd project structure and configure the locale.

## Step 1: Run Init Command

Run:
```
mysd init
```

This creates the scaffold:
- `.claude/mysd.yaml` — project configuration file
- `openspec/` — OpenSpec root directory
- `openspec/specs/` — spec storage directory

If there are any errors, report them to the user and stop.

## Step 2: Set Language

Ask the user:
```
What language should mysd use? (e.g., zh-TW, en-US, ja-JP)
Press Enter to skip and configure later with /mysd:lang.
```

If the user provides a locale, run:
```
mysd lang set {user_choice}
```

This atomically sets `response_language` in `.claude/mysd.yaml` and `locale` in `openspec/config.yaml`.

If the user presses Enter (skips), continue without setting language.

## Step 3: Confirm

Tell the user:
```
Project initialized. Run /mysd:scan to discover your codebase and generate specs.
```

If language was set, also mention:
```
Language set to: {locale}
```
