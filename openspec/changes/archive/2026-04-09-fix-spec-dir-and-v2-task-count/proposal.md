## Problem

Agent 定義檔（mysd/agents/*.md）中所有檔案路徑都硬寫 `.specs/`，但 mysd 支援兩種 spec 目錄：`.specs/`（FlavorMySD）和 `openspec/`（FlavorOpenSpec）。當專案使用 `openspec/` 時，agent 讀寫的路徑是錯的，靠 AI 自行猜測才能找到正確路徑，不可靠。

此外，`internal/validator/validator.go` 的 task count 驗證使用 V1 parser（`ParseTasks`，數 markdown `- [ ]` 行），但 mysd 已改用 TasksFrontmatterV2 格式，task 定義放在 YAML frontmatter 的 `tasks` 陣列中，body 沒有 checkbox。導致 `mysd validate` 永遠報告 "total does not match actual task count" 的誤報。

## Root Cause

1. **Agent 路徑硬編碼**：11 個 agent 共 47 處 `.specs/` 引用，Go binary 的 `--context-only` JSON 輸出沒有包含 `spec_dir` 欄位，orchestrator 無法傳遞正確路徑給 agent。

2. **Validator V1/V2 不相容**：`validateTasks` 呼叫 `spec.ParseTasks()`（V1），只會數 body 裡的 `- [ ]` 行。TasksFrontmatterV2 的 task 在 YAML frontmatter 裡，V1 parser 數不到。

## Proposed Solution

### Fix 1：spec_dir 傳遞鏈

1. Go binary 的所有 `--context-only` JSON 輸出加入 `spec_dir` 欄位（6 個 cmd 檔）
2. Orchestrator SKILL.md 從 JSON 提取 `spec_dir`，傳入 agent context
3. Agent 定義檔中 `.specs/` 全部改成 `{spec_dir}/`

### Fix 2：Validator V2 相容

修改 `validateTasks` 使用 `ParseTasksV2`，當 `fm.Tasks` 非空時用 `len(fm.Tasks)` 比對 `total`。

## Success Criteria

- `mysd validate` 在 `openspec/` 專案上不報錯
- `mysd validate` 在 TasksFrontmatterV2 格式的 tasks.md 上不報誤報
- Agent 在 `openspec/` 專案中使用正確路徑讀寫 artifact
- Agent 在 `.specs/` 專案中行為不變

## Impact

- Affected code:
  - `cmd/spec.go`, `cmd/plan.go`, `cmd/design.go`, `cmd/execute.go`, `cmd/scan.go`, `cmd/verify.go` — 加 `spec_dir` 欄位
  - `internal/validator/validator.go` — V2 task count 修正
  - `mysd/skills/*/SKILL.md` — 提取並傳遞 `spec_dir`（~10 個 orchestrator）
  - `mysd/agents/*.md` — `.specs/` → `{spec_dir}/`（11 個 agent、47 處）
