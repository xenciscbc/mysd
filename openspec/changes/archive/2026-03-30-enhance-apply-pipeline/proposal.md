---
spec-version: "1"
change: enhance-apply-pipeline
type: feature
status: proposed
---

## Why

mysd apply 目前缺乏執行前品質保障和靈活的執行粒度。Executor agent 收到 task 後直接寫 code，沒有結構化的 pre-task 檢查（reuse、quality、efficiency、no placeholders）。遇到問題時的暫停條件不明確。此外，只有 per-task（single）和 per-wave（wave）兩種模式，缺少 per-spec 模式讓 agent 在同一 spec 的 tasks 間保持 context 連續性。最後，執行前沒有 preflight 檢查來偵測 missing files 或 stale artifacts。

## What Changes

- 新增 `execution_mode: "spec"` — per-spec 執行模式，每個 spec spawn 一個 executor agent 處理該 spec 所有 tasks
- 新增 `spec-executor` role 到 model profile 表 — quality/balanced 用 opus，budget 用 sonnet
- 強化 `mysd-executor.md` agent — 加入 4 項 pre-task checks（Reuse、Quality、Efficiency、No Placeholders）
- 強化 `mysd-executor.md` agent — 加入 4 種明確暫停條件（task unclear、design issue、error/blocker、user interrupt）
- 新增 `mysd execute --preflight` CLI flag — 檢查 file existence 和 artifact staleness
- 更新 `apply SKILL.md` — Step 2 後呼叫 preflight，有問題時 warn + confirm

## Capabilities

### New Capabilities

- `per-spec-execution`: 新增 `execution_mode: "spec"` 模式，orchestrator 對每個 spec spawn 一個 spec-executor agent 持續處理該 spec 的所有 tasks
- `execution-preflight`: CLI 端 `--preflight` flag 在執行前檢查 file existence 和 artifact staleness，回傳結構化 JSON

### Modified Capabilities

- `execution`: 支援 `execution_mode: "spec"` 模式，新增 `spec-executor` role 到 model profile，executor agent 加入 pre-task checks 和 pause conditions

## Impact

- 新增 Go code: `internal/executor/spec_mode.go`（per-spec execution context 建構）
- 修改 Go code: `internal/config/config.go`（DefaultModelMap 加 `spec-executor` role）
- 修改 Go code: `cmd/execute.go`（`--preflight` flag、spec mode context）
- 修改 skill: `mysd/skills/apply/SKILL.md`（preflight step、spec mode 分派）
- 修改 agent: `mysd/agents/mysd-executor.md`（pre-task checks、pause conditions）
- 影響 spec: `openspec/specs/execution/spec.md`
