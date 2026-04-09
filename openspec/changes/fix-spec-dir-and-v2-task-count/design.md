## Context

mysd 使用 `spec.DetectSpecDir()` 偵測專案的 spec 目錄（`.specs/` 或 `openspec/`），Go binary 的所有命令都能正確處理。但 agent 定義檔和部分 orchestrator 都硬寫 `.specs/`，在 `openspec/` 專案中依賴 AI 自行猜測路徑。

另外，validator 的 task count 檢查使用 V1 parser，無法正確處理 TasksFrontmatterV2 格式。

## Goals / Non-Goals

**Goals:**

- 所有 agent 使用動態 `{spec_dir}` 取代硬寫的 `.specs/`
- Go binary `--context-only` JSON 輸出包含 `spec_dir` 欄位
- Orchestrator 提取 `spec_dir` 並傳入 agent context
- Validator 正確處理 V1 和 V2 task 格式

**Non-Goals:**

- 不改變 `DetectSpecDir()` 的偵測邏輯
- 不新增第三種 spec 目錄慣例
- 不修改 V1 task 格式的相容性

## Decisions

### D1: spec_dir 傳遞機制

在每個 `--context-only` JSON 輸出中加入 `spec_dir` 欄位。值來自 `spec.DetectSpecDir()` 的回傳值（`.specs` 或 `openspec`）。

Orchestrator 從 JSON 提取後，加入 agent 的 Task context JSON：
```json
{
  "spec_dir": ".specs",
  "change_name": "my-feature",
  ...
}
```

Agent 使用 `{spec_dir}/changes/{change_name}/` 路徑。

**為何不用環境變數或 config：** spec_dir 是 per-project 的，由目錄結構決定。透過 context JSON 傳遞最直接，不需要額外的配置層。

### D2: Agent 路徑替換策略

所有 agent 中 `.specs/` 替換為 `{spec_dir}/`。替換是機械式的：

- `.specs/changes/{change_name}/` → `{spec_dir}/changes/{change_name}/`
- `.specs/changes/{name}/` → `{spec_dir}/changes/{name}/`

Agent 開頭加入 `{spec_dir}` 的說明：
```
- `spec_dir`: The detected spec directory for this project (`.specs` or `openspec`)
```

### D3: Validator V2 task count

修改 `validateTasks` 的計數邏輯：

1. 先嘗試 `ParseTasksV2`（讀 frontmatter `tasks` 陣列）
2. 如果 `fm.Tasks` 非空 → 用 `len(fm.Tasks)` 作為實際 task 數
3. 如果 `fm.Tasks` 為空 → fallback 到 V1（數 `- [ ]` 行）

這保持了 V1 相容性，同時正確處理 V2。

### D4: 受影響的 Go cmd 檔案

以下 6 個命令的 `--context-only` JSON 需要加 `spec_dir`：

| 檔案 | 命令 | JSON 建構位置 |
|------|------|--------------|
| cmd/spec.go | `mysd spec --context-only` | context map |
| cmd/plan.go | `mysd plan --context-only` | context map |
| cmd/design.go | `mysd design --context-only` | context map |
| cmd/execute.go | `mysd execute --context-only` | ExecutionContext struct |
| cmd/scan.go | `mysd scan --context-only` | ScanContext struct |
| cmd/verify.go | `mysd verify --context-only` | VerificationContext struct |

前三個直接加 `"spec_dir": specDir` 到 map。後三個需要在 struct 加欄位。

## Risks / Trade-offs

- **改動面大（47 處 agent 替換）**：但全是機械式替換，邏輯不變。[Risk] → 用 grep 確認替換完整性。
- **Orchestrator 必須正確傳遞 spec_dir**：如果 orchestrator 漏傳，agent 會找不到 `{spec_dir}` 變數。[Risk] → 每個 orchestrator 的 context JSON 都加 spec_dir，並在 agent 的 Input 段落明確列出。
