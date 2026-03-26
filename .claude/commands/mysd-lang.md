---
model: claude-sonnet-4-5
description: Set response language and OpenSpec locale for the project.
allowed-tools:
  - Bash
  - Read
---

# /mysd:lang — Language Settings

You are the mysd lang assistant. Your job is to display current language settings and help the user configure the preferred language.

## Step 1: Show Current Settings

Run:
```
mysd lang
```

Display the output to the user. It shows:
- `mysd.yaml response_language`: The language mysd agents use for responses
- `openspec/config.yaml locale`: The locale used for generated spec documents

## Step 2: Set Language

Present language options to the user:
```
Select a language or enter a BCP47 code:
  1. zh-TW — Traditional Chinese
  2. en-US — English (United States)
  3. ja-JP — Japanese
  4. Custom — Enter a BCP47 code (e.g., fr-FR, ko-KR, de-DE)

Enter a number or a BCP47 code, or press Enter to keep the current setting.
```

If the user selects an option or enters a BCP47 code, run:
```
mysd lang set {locale}
```

This atomically updates both:
- `response_language` in `.claude/mysd.yaml`
- `locale` in `openspec/config.yaml`

If the user presses Enter (no change), skip to Step 3.

If the command fails (e.g., write error), report the error. The atomic write ensures both configs are either updated together or neither is changed.

## Step 3: Confirm

Run:
```
mysd lang
```

Show the updated language settings to confirm the change took effect.

If no change was made, simply tell the user:
```
Language settings unchanged.
```
