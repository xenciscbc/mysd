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

### Step 1: Classify and Confirm

Determine the question type, then use AskUserQuestion to confirm with the user:

**Use AskUserQuestion:**
> 我對你的問題的理解是 [一句話重述]。我建議用以下模式處理：
>
> A) **4 維度深度分析** — 派 4 個 subagent 分別從 Coverage/Consistency/Ambiguity/Gaps 維度分析問題，蒐集證據，產出 Decision Doc
> B) **直接分析** — 用 LLM 本身的能力直接分析回答，不跑 subagent

若使用者的提問明確指定要檢查 spec 品質，則額外顯示：
> C) **Spec 健康檢查** — 檢查 openspec/ 下的 spec 檔案品質（Coverage/Consistency/Ambiguity/Gaps）

推薦選項和理由。讓使用者選擇。

- 使用者選 **A** → 進入 Step 2（Context Gathering），Step 3 用 4 個 subagent 做維度分析
- 使用者選 **B** → 直接分析回答，不經過 subagent 流程
- 使用者選 **C**（僅 spec 相關問題時出現）→ 跳到 [Spec Health Check Mode](#spec-health-check-mode)

### Step 2: Context Gathering

Gather evidence in this order — stop when you have enough to frame 2+ options:

1. Codebase — Grep/Glob/Read for existing patterns, prior decisions, constraints
2. Spec health (only if user chose C, i.e. spec-related question) — run the 4-dimension health check (read `formats/health-check.md`) against the relevant change or spec directory.
3. Git history — `git log --oneline` or `git diff` for recent relevant changes
4. Project docs — CLAUDE.md, README, any spec files in `openspec/`
5. WebSearch — only if the above leave critical gaps and WebSearch is available

### Step 3: Option Framing

**若使用者選 A（4 維度深度分析）：**

派 4 個 subagent 平行分析，每個 subagent 負責一個維度：
1. **Coverage subagent** — 這個問題涵蓋了哪些面向？有沒有遺漏的選項或情境？
2. **Consistency subagent** — 各選項的論點是否自洽？有沒有矛盾的證據？
3. **Ambiguity subagent** — 問題本身或各選項中有哪些模糊地帶？需要哪些澄清？
4. **Gaps subagent** — 還缺少什麼關鍵資訊？哪些假設尚未驗證？

每個 subagent 回傳結構化分析結果，彙整後進入 Option Framing。

**若使用者選 B（直接分析）：** 跳過 subagent，直接由 LLM 分析。

---

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
