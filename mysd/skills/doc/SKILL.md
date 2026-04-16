---
name: mysd:doc
description: >
  Update documentation files based on code changes. Detects which docs need updating
  from git diff, generates content matching existing style, and applies changes with
  user confirmation. DO NOT use for spec files (use mysd:spec) or research/decisions
  (use mysd:research).
---

## When to Use

USE this skill when:
- The user says "update docs", "sync README", "add to CHANGELOG"
- Post-feature doc sync after code changes are committed or staged
- Keeping multi-language README variants in sync (e.g. README.zh-TW.md)

DO NOT USE when:
- The file to update is an OpenSpec spec file → use `mysd:spec`
- The task is research, analysis, or a technical decision → use `mysd:research`
- Writing new docs from scratch with no code change context → just write directly

---

## Flow

### Step 1: Detect Changes

Run `git diff --name-only HEAD~1` by default.

Accept overrides:
- `HEAD~N` — go back N commits
- `<sha1>..<sha2>` — explicit range
- Explicit file list — skip git diff entirely

If the diff is empty, report "no changes detected" and stop.

### Step 2: Impact Analysis

Map each changed file to the docs that need updating:

| Change Type | Detection Pattern | Docs to Update |
|---|---|---|
| New/removed command | `cmd/*.go` added/removed | README.md, README.zh-TW.md, CLAUDE.md |
| API change | exported function signature changed | API docs, CHANGELOG.md |
| Config change | config struct or yaml schema changed | README (config section), example configs |
| Bug fix | `.go` modified + "fix" in commit message | CHANGELOG.md |
| Architecture change | new package directory or major refactor | ARCHITECTURE.md, README |
| Dependency update | `go.mod`, `package.json` changed | README (installation section) |

Fallback: grep changed keywords across `.md` files to find additional affected docs.

### Step 2.5: Confirm Scope

Use AskUserQuestion to present the impact analysis results and let the user choose:

> 根據變更分析，以下文件需要更新：
>
> A) 全部更新 — [列出所有受影響的文件]
> B) 只更新部分 — [讓使用者勾選要更新哪些]
> C) 跳過 — 這次不需要更新文件

如果使用者選 C，結束。選 A 或 B，帶著選定的文件清單繼續到 Step 3。

### Step 3: Style Matching

Before writing, read the first 50 lines of each target doc. Match:
- Heading levels and capitalization style
- List style (dashes vs asterisks, ordered vs unordered)
- Language and tone (formal, casual, zh-TW vs en)
- Code block conventions (language tags, indentation)

### Step 4: Multi-Language Sync

When updating `README.md`, check for `README.{locale}.md` or `README-{locale}.md` variants in the repo root. For each variant found, generate the equivalent update translated into that locale's language, matching the same section structure.

### Step 5: Apply and Confirm

For each doc with proposed changes:
1. Show the proposed edit using the Edit tool format (old block → new block)
2. Wait for user confirmation before proceeding
3. Apply one doc at a time — do not batch without confirmation

If the user rejects a change, note it and move to the next doc.

### Step 6: Summary

After all changes are handled, list:
- What was updated (file path + one-line description of change)
- What was skipped and why (rejected, already current, etc.)
