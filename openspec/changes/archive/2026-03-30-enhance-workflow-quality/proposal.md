## Why

與 spectra:propose 做功能對比後，發現 mysd 的 propose + plan 流程有三個品質缺口：

1. **propose 的 proposal template 是固定的** — proposal-writer agent 不管 change type（Feature/Bug Fix/Refactor），一律使用 Summary/Motivation/Scope/Success Criteria 結構，導致 Bug Fix 缺少 Root Cause/Solution 欄位，Feature 缺少 Capabilities 欄位
2. **plan 固定跑 design + tasks** — 即使小 change 不需要 design.md，plan 仍然生成它，浪費 token 和時間
3. **plan 沒有 artifact 結構分析** — reviewer agent 做 AI self-check，但缺少 CLI 層面的 cross-artifact 分析（如 coverage、consistency、ambiguity、gaps），需要新增 `mysd analyze` command

## What Changes

- propose SKILL.md Step 9：傳入 `change_type` 給 proposal-writer agent，讓 agent 根據 type 使用對應 template
- proposal-writer agent：新增三種 template（Feature / Bug Fix / Refactor），取代現有的固定 template
- plan SKILL.md：在 design phase 前加入判斷邏輯，允許跳過 design（小 change、無 cross-cutting concern）
- 新增 `cmd/analyze.go`：實作 `mysd analyze [change-name]` CLI command，做 cross-artifact 結構分析
- plan SKILL.md：在 reviewer 之後加入 analyze-fix loop（最多 2 輪）

## Non-Goals

- 不改 spec-writer agent 的 workflow
- 不改 apply 的 parallel execution 機制
- 不加 Park/Unpark 機制到 propose 結束流程

## Capabilities

### New Capabilities

- `artifact-analysis`: `mysd analyze` CLI command，提供 cross-artifact 結構分析（coverage、consistency、ambiguity、gaps），輸出 JSON 格式的 findings

### Modified Capabilities

- `planning`: plan pipeline 新增 optional design skip 邏輯和 analyze-fix loop step
- `reviewer-agent`: propose pipeline 傳入 `change_type` 作為 context，reviewer 可驗證 proposal template 是否符合 type

## Impact

- Affected specs: `artifact-analysis`（新）、`planning`（改）、`reviewer-agent`（改）
- Affected code:
  - `cmd/analyze.go`（新）+ `cmd/analyze_test.go`（新）
  - `internal/analyzer/` 目錄（新）— analyze 邏輯
  - `mysd/skills/propose/SKILL.md` — Step 9 傳入 change_type
  - `mysd/agents/mysd-proposal-writer.md` — 三種 template
  - `mysd/skills/plan/SKILL.md` — optional design skip + analyze-fix loop
