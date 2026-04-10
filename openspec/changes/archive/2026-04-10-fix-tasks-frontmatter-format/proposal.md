## Why

Planner agent 產出的 `tasks.md` 是純 markdown checkbox 格式，缺少 `TasksFrontmatterV2` YAML frontmatter。Go 代碼（`ParseTasksV2`、`UpdateTaskStatus`、`CheckCoverage`）全部依賴此 frontmatter schema，但因為 silent fallback（回傳 zero-value struct），問題被掩蓋。結果是 task 狀態追蹤、覆蓋率檢查、依賴排序全部失效。

## What Changes

- 修正 planner agent prompt（`mysd/agents/mysd-planner.md`），確保產出的 `tasks.md` 包含完整的 `TasksFrontmatterV2` YAML frontmatter
- 修正 planning spec，使 frontmatter 欄位定義與 `TasksFrontmatterV2` Go struct 一致（`spec-version`、`total`、`completed`、`tasks` 陣列）
- 確保 `ParseTasksV2` → `UpdateTaskStatus` → `CheckCoverage` 完整鏈路在新格式下正常運作

## Non-Goals

- 不回溯修正既有的無 frontmatter tasks.md 檔案 — 現有的 fallback 機制保留
- 不變更 `TasksFrontmatterV2` Go struct 本身 — schema 設計正確，問題在 agent prompt

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `planning`: Task planning spec 的 frontmatter 欄位定義需與 `TasksFrontmatterV2` struct 對齊（目前 spec 寫 `spec-version, change, status`，但 struct 是 `spec-version, total, completed, tasks`）

## Impact

- Affected specs: `openspec/specs/planning/spec.md`
- Affected code: `mysd/agents/mysd-planner.md`（agent prompt）、`internal/spec/updater.go`（已正確，驗證即可）、`internal/planchecker/checker.go`（已正確，驗證即可）
