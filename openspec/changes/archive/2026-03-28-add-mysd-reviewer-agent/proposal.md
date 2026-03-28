## Why

`mysd:plan` 目前沒有品質閘道：planner 完成後直接進入確認，中間沒有任何機制驗證 artifacts 的一致性、完整性與範圍合理性。inline self-review 邏輯分散在 propose/SKILL.md 且無法被 plan 複用，且無法透過 model profile 控制品質層級。

## What Changes

- 新增 `mysd-reviewer` agent（`mysd/agents/mysd-reviewer.md`），負責掃描並修復 plan phase 所有 4 個 artifacts（proposal + specs + design + tasks）
- `mysd:plan` pipeline 在 planner 完成後、plan-checker 之前插入 Step 5b，自動呼叫 mysd-reviewer
- `internal/config/config.go` 的 `DefaultModelMap` 新增 `reviewer` role（quality: opus, balanced: sonnet, budget: sonnet）
- `mysd plan --context-only` 回傳的 JSON 新增 `reviewer_model` 與 `plan_checker_model` 欄位，讓 SKILL.md 能為各 agent 套用正確的 profile-resolved model
- `mysd:plan` Step 6 改用 `{plan_checker_model}` 呼叫 mysd-plan-checker（修正現有的 model 共用問題）
- `mysd:discuss` 加入討論品質規範（一次一問、具體選項、禁止空話）與強制收斂機制（結論 summary + 不允許無結論結束）
- `mysd:propose` 的 Step 12 inline self-review 保持不變（不引入 reviewer）

## Non-Goals

- propose phase 不引入 mysd-reviewer（propose 的 inline review 已足夠）
- 不修改 mysd-plan-checker 的行為或邏輯
- 不改變 mysd:propose 的 model 解析方式

## Capabilities

### New Capabilities

- `reviewer-agent`: 獨立的 artifact quality reviewer agent，在 plan phase 對 4 個 artifacts 執行 4 項品質檢查（no placeholders、internal consistency、scope check、ambiguity check），inline 修復問題並回傳修復摘要

### Modified Capabilities

- `planning`: plan pipeline 新增 reviewer step（Step 5b），plan context JSON 新增 `reviewer_model` 與 `plan_checker_model` 欄位，discuss skill 加入品質規範與收斂機制
- `model-passthrough`: `DefaultModelMap` 新增 `reviewer` role；plan context JSON 擴充為每個 plan-phase agent 回傳獨立的 model 欄位

## Impact

- Affected specs: `reviewer-agent`（新增）、`planning`（修改）、`model-passthrough`（修改）
- Affected code:
  - `mysd/agents/mysd-reviewer.md`（新增）
  - `mysd/skills/plan/SKILL.md`（修改）
  - `mysd/skills/discuss/SKILL.md`（修改）
  - `internal/config/config.go`（修改 DefaultModelMap）
  - `cmd/plan.go` 或 plan context 相關程式碼（修改 --context-only JSON 輸出）
