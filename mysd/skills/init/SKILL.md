---
model: opus
description: Initialize mysd project structure (scaffold + locale setup).
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

## Step 2: Set Project Name

Read `openspec/config.yaml` and check if `project_name` is set. If not set or empty:

1. Detect the current folder name (use `basename $(pwd)`)
2. Suggest it as default:
   ```
   Project name? (default: {folder_name})
   Press Enter to use the folder name, or type a custom name.
   ```
3. Run:
   ```
   mysd config set project_name "{chosen_name}"
   ```

If `project_name` is already set, skip this step.

## Step 3: Set Language

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

## Step 4: Confirm

Tell the user:

```
Project initialized. Next steps:

- /mysd:scan    — discover your codebase and generate specs automatically
- /mysd:propose — create a spec for a specific change or feature
- /mysd:discuss — explore ideas and build a spec through discussion
```

If language was set, also mention:
```
Language set to: {locale}
```
