---
name: mysd:research
description: >
  Research ambiguous problems and make gray-area decisions with evidence.
  Use when facing technical choices with 2+ viable options and no clear consensus,
  or when analyzing spec quality. DO NOT use for documentation updates (use mysd:doc)
  or spec writing (use mysd:spec). DO NOT use for questions with clear best practices
  or official documentation answers.
---

## When to Use

USE this skill when:
- There are 2+ viable approaches with no community consensus
- A best practice exists but specific constraints make it a poor fit
- A decision must be made with incomplete information
- The user asks about spec quality, health, or completeness

DO NOT USE when:
- Official docs or a clear best practice answers the question directly → answer directly
- The question is about syntax, API usage, or error messages → answer directly
- The task is updating documentation → use `mysd:doc`
- The task is writing or editing a spec → use `mysd:spec`

---

## Flow

### Step 1: Classify

Determine if this is a gray area. Gray area means one of:
- **(a)** Multiple viable approaches, no consensus
- **(b)** Best practice exists but doesn't apply given the constraints
- **(c)** Must decide with incomplete information

If **not** gray area: answer directly and stop here.

If **spec health check** request: skip to [Spec Health Check Mode](#spec-health-check-mode).

### Step 2: Context Gathering

Gather evidence in this order — stop when you have enough to frame 2+ options:

1. Codebase — Grep/Glob/Read for existing patterns, prior decisions, constraints
2. Spec health — if the question involves an area with OpenSpec specs, run the 4-dimension health check (read `formats/health-check.md`) against the relevant change or spec directory. Coverage gaps, ambiguous language, inconsistencies, and missing scenarios are decision-relevant context.
3. Git history — `git log --oneline` or `git diff` for recent relevant changes
4. Project docs — CLAUDE.md, README, any spec files in `openspec/`
5. WebSearch — only if the above leave critical gaps and WebSearch is available

### Step 3: Option Framing

Frame **2–4 options** (fewer = no real decision; more = consolidate first).

For each option, capture:
- Evidence (concrete: docs, benchmarks, observed behavior — not speculation)
- Pros
- Cons
- Effort: S (hours) / M (days) / L (weeks or structural change)

### Step 4: Recommendation

Pick one option. State:
- **Confidence:** 1–10 (most gray area decisions land 6–8; avoid 5 or 10)
- **Reasoning:** why this option wins, what trade-off is accepted
- **What would change my mind:** specific conditions that would reverse the call

### Step 5: Output

Read `formats/decision-doc.md` for the exact template. Produce a complete Decision Doc using that template.

---

## Spec Health Check Mode

Triggered when the user asks about spec quality, health, or completeness for a change.

1. Read `formats/health-check.md` for the full procedure and finding formats.
2. Identify the change directory (`openspec/changes/{change-name}/`).
3. Run all 4 dimensions in order: **Coverage → Consistency → Ambiguity → Gaps**
4. Skip dimensions whose required artifacts are missing (skip rules are in the format file).
5. Present findings using the Output Summary Format defined in `formats/health-check.md`.
