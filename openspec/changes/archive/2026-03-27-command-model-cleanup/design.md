## Context

mysd 的 plugin 系統由兩層組成：
- **Go binary**（`cmd/`、`internal/`）— 提供 `--context-only` JSON 輸出，包含 config 解析後的 model 欄位
- **Plugin skills**（`plugin/commands/*.md`）— 讀取 JSON context，orchestrate agent spawning

目前 profile 系統（`internal/config/config.go` 的 `DefaultModelMap` + `ResolveModel`）已實作但未被 skill 層使用。所有 command 和 agent 的 frontmatter 都 hardcode `model: claude-sonnet-4-5`，繞過了 profile。

## Goals / Non-Goals

**Goals:**

- Profile 設定能實際影響工作流中 agent 的 model 選擇
- 獨立工具命令有明確的固定 model（不受 profile 影響）
- 清理已整合進其他命令的廢棄 skill 檔案
- 所有過時引用統一更新

**Non-Goals:**

- 不修改 profile 系統本身的架構（DefaultModelMap 結構不變，只改值）
- 不新增 agent role（uat-guide、scanner 等未納入 profile 管理）
- 不修改 Go binary 的子命令邏輯（`cmd/execute.go` 等保留，只改提示訊息）

## Decisions

### Model 短名格式

`DefaultModelMap` 的值從全名 `claude-sonnet-4-5` 改為短名 `sonnet`。

**理由**：Claude Code 的 Agent/Task tool 的 `model` 參數只接受 `sonnet`、`opus`、`haiku`。使用短名讓 binary context JSON 的 model 值可直接傳給 Task tool，無需轉換。

**替代方案**：在 skill 層做 full-name → short-name mapping。但這增加了每個 command skill 的複雜度，且 mapping 需要隨 model 更新而維護。

### Command model 分類策略

三類處理：

1. **工作流 command**（propose、discuss、plan、apply、archive、ff、ffe、uat）— 移除 model frontmatter，繼承呼叫者 session 的 model。它們 spawn 的 agent model 由 profile 控制。
2. **獨立輕量工具**（status、lang、model、note、docs、statusline、update）— 固定 `claude-sonnet-4-5`。這些是使用者直接呼叫的簡單操作，不需要高階推理。
3. **獨立重度工具**（init、scan、fix）— 固定 `claude-opus-4-6`。init 做專案初始化、scan 分析 codebase、fix 做故障診斷，需要深度推理能力。

**理由**：工作流 command 是 orchestrator，它們的 model 可以跟隨 session；但 agent 的 model 需要精確控制（不同 role 不同能力需求），所以由 profile 管理。獨立工具不在 profile 管理範圍（不呼叫 profile-managed agents），需要固定 model。

### Model passthrough 機制

Command skill 讀取 `--context-only` JSON 中的 `model` 欄位，在 Task tool 呼叫時傳入 `model` 參數：

```
1. Command 執行 `mysd plan --context-only`
2. 解析 JSON，取得 model: "sonnet"
3. 顯示 "Spawning mysd-designer (sonnet)..."
4. Task tool 呼叫 mysd-designer，model 參數設為 "sonnet"
```

**問題**：不同 agent role 可能有不同 model（例如 planner 用 sonnet，但 executor 用 haiku）。目前 `--context-only` 只輸出一個 model 欄位。

**解決**：對於 plan command，binary 輸出的是 planner 的 model。但 plan 也需要 spawn designer — designer 的 model 需要另外取得。方案：
- `plan --context-only` 輸出包含 `model`（planner role）欄位，command skill 直接用
- 對於需要多個 role 的 command，binary 的 context JSON 需擴展為包含所有相關 role 的 model mapping，或 command 額外呼叫 `mysd model` 解析

目前先採用最簡方案：每個 `--context-only` 輸出對應其主要 agent role 的 model。若 command 需要多個 role，在 skill 中多次讀取或由 binary 擴展 context JSON 輸出多個 model 欄位。

### 廢棄命令處理

直接刪除檔案而非保留 redirect。

**理由**：redirect 會增加維護負擔，且所有引用處都會同步更新。使用者如果呼叫已刪除的命令，Claude Code 會提示命令不存在，這比 redirect 更清楚。

## Risks / Trade-offs

- **[Risk] Binary context JSON 只有單一 model 欄位** → command 需要 spawn 多個不同 role 的 agent 時，可能需要額外的 binary 呼叫或 JSON 擴展。短期先用主要 role 的 model，後續視需要擴展。
- **[Risk] 使用者可能習慣舊命令名稱** → 刪除 execute/spec/design/capture 後，舊習慣會報錯。透過引用清理確保所有提示都指向正確命令。
- **[Trade-off] 獨立工具固定 model 而非繼承** → 失去靈活性但保證穩定體驗。使用者用 opus session 跑 `/mysd:status` 時不會浪費高價 model 在簡單狀態查詢上。
