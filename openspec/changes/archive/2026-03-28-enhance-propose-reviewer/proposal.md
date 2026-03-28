## Why

目前 propose skill 的 Step 12 使用 inline self-review（由 propose orchestrator 自己執行 4 個 check），但 mysd-reviewer agent 已經實作了相同的品質檢查邏輯，且支援 `phase: "propose"` 模式。這造成兩個問題：

1. **邏輯重複** — propose 的 inline check 和 mysd-reviewer 的 Check 1-4 是相同內容，維護兩份
2. **品質落差** — propose 階段沒有 `mysd validate` CLI 驗證，reviewer agent 有整合 validate output 的能力但在 propose 階段未被使用

## What Changes

- propose SKILL.md Step 3b：擴充 model resolution，新增 `reviewer_model` 欄位（根據 profile 映射）
- propose SKILL.md Step 12：移除 inline self-review 邏輯，改為先跑 `mysd validate`，再 spawn `mysd-reviewer` with `phase: "propose"` 和 `reviewer_model`
- mysd-reviewer agent：在品質檢查前加入 Rationalization Table（防止 AI 偷懶跳過檢查的 anti-pattern 對照表）

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `reviewer-agent`: 新增 Rationalization Table section，強化 reviewer 的 self-check 能力
- `model-passthrough`: propose skill 需要解析 reviewer role 的 model mapping，擴充 Step 3b 的 model resolution 邏輯

## Impact

- Affected specs: `reviewer-agent`、`model-passthrough`
- Affected code:
  - `mysd/skills/propose/SKILL.md` — Step 3b model resolution + Step 12 reviewer invocation
  - `mysd/agents/mysd-reviewer.md` — Rationalization Table
