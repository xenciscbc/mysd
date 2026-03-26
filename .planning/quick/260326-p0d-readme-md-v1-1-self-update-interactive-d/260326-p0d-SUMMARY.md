---
phase: quick
plan: 260326-p0d
subsystem: documentation
tags: [readme, v1.1, documentation, self-update, wave-execution, interactive-discovery]
dependency_graph:
  requires: []
  provides: [updated-readme-v1.1]
  affects: [README.md]
tech_stack:
  added: []
  patterns: []
key_files:
  created: []
  modified:
    - README.md
decisions:
  - "README.md 保留英文為主語言（符合原有風格），技術術語不翻譯"
  - "Commands 表格分拆為三個子表格（Spec Workflow / Utility / Fast-Forward）提升可讀性"
  - "/mysd:execute 更名為 /mysd:apply — 符合 Phase 8 SKILL.md 層改名決策"
metrics:
  duration: "72s"
  completed_date: "2026-03-26"
  tasks: 1
  files_changed: 1
---

# Quick Task 260326-p0d: README.md v1.1 Update Summary

Updated README.md to document all v1.1 features including self-update, wave parallel execution, interactive discovery, and 6 new commands (discuss/fix/update/note/model/lang).

## What Was Done

### Task 1: Update README.md with v1.1 features

- **Why mysd? 區段** — 新增 4 個 v1.1 亮點：wave parallel execution、interactive discovery、self-update、deferred notes
- **Quick Start 區段** — 加入 `/mysd:discuss` 步驟；將 `/mysd:execute` 更名為 `/mysd:apply`；加入 wave 平行模式說明
- **Commands 表格** — 分拆為 Spec Workflow / Utility / Fast-Forward 三個子表格；補齊所有 v1.1 指令（discuss、fix、update、note、model、lang）；更新 ff/ffe 描述
- **How It Works 生命週期** — 加入 Discuss 為可選步驟；Plan 步驟加入 plan-checker 說明；Apply 步驟加入 wave 執行說明
- **Wave Parallel Execution 子章節** — 新增說明依賴分析、worktree 平行執行、AI 衝突解決、模式選擇
- **Self-Update 章節** — 新增包含 CLI 範例、SHA256 驗證、rollback、plugin file sync 的完整說明
- **Configuration 區段** — 新增 `response_language` 欄位；補充 `execution_mode: wave` 和 `model_profile` 說明
- **Tech Stack** — 修正 Go 版本：1.25+ → 1.23+（符合 CLAUDE.md 規範）

**Commit:** `626f275`
**Lines:** 155 → 220 (+65 lines)
**Keyword matches (wave/worktree/parallel/self-update/discuss/fix/note/model/lang):** 20 occurrences

## Deviations from Plan

None — plan executed exactly as written. One minor structural decision: Commands table split into three sub-tables (Spec Workflow, Utility, Fast-Forward) for better readability — not in plan but improves UX without contradicting requirements.

## Known Stubs

None.

## Self-Check: PASSED

- [x] `README.md` exists and has 220 lines (min 180)
- [x] Commit `626f275` exists
- [x] All 6 new v1.1 commands present in Commands table
- [x] Self-Update section with CLI examples present
- [x] Wave Parallel Execution section present
- [x] Interactive discovery (`/mysd:discuss`) in Quick Start and lifecycle
- [x] Go version shows 1.23+
- [x] No broken markdown formatting
