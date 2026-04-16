---
name: mysd:sync
description: >
  Sync all content artifacts after a change: research → doc → spec.
  Chains the three mysd skills using subagents, passing context between them.
  Each skill runs in its own context window. Use when you want to sync everything
  after a significant change. For individual tasks, use mysd:research, mysd:doc, or mysd:spec directly.
---

## When to Use

USE when: "run everything", "do everything", "full pipeline", or a significant change that
needs all three artifacts updated (decision doc + docs + spec).

DO NOT USE when only one artifact needs updating — call `mysd:research`, `mysd:doc`, or
`mysd:spec` directly instead.

---

## Flow

### Step 1: Understand Scope

Ask (or infer): What changed? Which skills are needed — all three or a subset?
- **Research**: gray area decision or spec health check
- **Doc**: README, CHANGELOG, or other prose docs
- **Spec**: OpenSpec spec files

Default to all three if the user said "run everything". Confirm scope before dispatching.

### Step 2: Run Research (if needed)

Dispatch via the Agent tool (general-purpose subagent). Prompt:
```
Read mysd-skills/research/SKILL.md and follow its instructions.
User question: {user_question}
Repo root: {repo_root}
```
Capture: Decision Doc title and confidence score from the output.
On failure: record "research failed", suggest `mysd:research`, continue to Step 3.

### Step 3: Run Doc Writer (if needed)

Dispatch via the Agent tool. Prompt:
```
Read mysd-skills/doc/SKILL.md and follow its instructions.
Research context: {research_summary}   ← 2–3 sentence digest of Step 2, or "N/A"
Diff range: {diff_range}               ← defaults to HEAD~1
Repo root: {repo_root}
```
Capture: list of files updated.
On failure: record "doc failed", suggest `mysd:doc`, continue to Step 4.

### Step 4: Run Spec Writer (if needed)

Dispatch via the Agent tool. Prompt:
```
Read mysd-skills/spec/SKILL.md and follow its instructions.
Research context: {research_summary}
Change name: {change_name}             ← OpenSpec change directory name, or "unknown"
Repo root: {repo_root}
```
Capture: list of spec files written or updated.
On failure: record "spec failed", suggest `mysd:spec`.

### Step 5: Summary

```
## mysd:sync — Pipeline Summary

Research:  {Decision Doc title} (confidence {N}/10)  |  skipped / failed: {reason}
Docs:      {N} file(s) updated: {list}               |  skipped / failed: {reason}
Specs:     {N} file(s) written/updated: {list}       |  skipped / failed: {reason}
```

If any step failed, list the direct commands to run to complete the pipeline.
