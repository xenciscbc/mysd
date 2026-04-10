## Why

目前 skill 取得 model profile 有三種方式：直接呼叫 `mysd model` 解析人類表格、從 `--context-only` JSON 取 `model` 欄位、Go 命令內部直接 resolve。這導致每個 SKILL.md 的 model 取得邏輯不一致，解析表格的 skill（propose、discuss、scan、uat）尤其脆弱且重複。統一為單一取得方式可以提升一致性與可維護性。

## What Changes

- 新增 `mysd model resolve <role>` 子命令，接受一個 role 名稱，回傳該 role 對應的 model short name（如 `sonnet`、`opus`）
- **BREAKING**: 從 `plan --context-only`、`execute --context-only`、`verify --context-only` 的 JSON 輸出中移除 `model`、`verifier_model`、`reviewer_model`、`plan_checker_model` 等 model 相關欄位
- 更新所有需要 model 的 SKILL.md，統一改用 `mysd model resolve <role>` 取得 model
- `design --context-only` 和 `spec --context-only` 同樣移除 model 欄位（如有）

## Non-Goals

- 不改變 `mysd model`（無參數）顯示完整表格的行為 — 那是給人看的 UI
- 不改變 `mysd model set` 子命令
- 不改變 model profile 的 resolve 邏輯本身（DefaultModelMap、custom profiles 等）
- 不處理不需要 model 的 skill（docs、note、status、lang、init、update、statusline）

## Capabilities

### New Capabilities

- `model-resolve-command`: 新增 `mysd model resolve <role>` CLI 子命令，為 SKILL.md 提供統一的 model 查詢介面

### Modified Capabilities

- `model-passthrough`: 移除 `--context-only` JSON 中的 model 欄位，model 查詢職責完全轉移至 `model resolve` 子命令

## Impact

- Affected specs: `model-passthrough`（行為變更）、`model-resolve-command`（新增）
- Affected code:
  - `cmd/model.go` — 新增 `resolve` 子命令
  - `cmd/plan.go` — 移除 context JSON 中的 model 欄位
  - `cmd/execute.go` — 移除 context JSON 中的 model/verifier_model 欄位
  - `cmd/verify.go` — 移除 context JSON 中的 model 欄位
  - `cmd/design.go` — 移除 context JSON 中的 model 欄位
  - `cmd/spec.go` — 移除 context JSON 中的 model 欄位
  - `mysd/skills/propose/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/discuss/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/scan/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/uat/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/ff/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/ffe/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/fix/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/apply/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/verify/SKILL.md` — 改用 `mysd model resolve`
  - `mysd/skills/plan/SKILL.md` — 改用 `mysd model resolve`
  - 相關測試檔案
