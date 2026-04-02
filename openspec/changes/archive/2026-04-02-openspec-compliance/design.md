## Context

mysd 的 archive 流程和 spec 格式與 OpenSpec 規範存在偏差。現有相關程式碼：

- `cmd/archive.go` — archive 到 `openspec/archive/<name>/`，只搬目錄不做 spec merge
- `internal/spec/delta.go` — 已有 `ParseDelta` 函數，支援 ADDED/MODIFIED/REMOVED，但缺少 RENAMED
- `internal/spec/schema.go` — 已有 `SpecFrontmatter` struct（`spec-version`, `capability`, `delta`, `status`），但缺少 OpenSpec 要求的 `name`, `description`, `version`, `generatedBy`
- `internal/spec/parser.go` — `ParseSpec` 用 frontmatter library 解析，已支援 frontmatter 讀取
- `internal/verifier/` — 驗證 MUST items done，不驗證 scenario 格式

## Goals / Non-Goals

**Goals:**

- Archive 路徑改為 `openspec/changes/archive/YYYY-MM-DD-<name>/`
- 實作 delta spec merge 邏輯，歸檔時自動合併回 main specs
- 支援 RENAMED delta 操作
- Scenario 格式加入 GIVEN 支援（產生端和驗證端）
- Spec frontmatter 加入 OpenSpec 要求的欄位
- 引入 `[~]` skipped task 標記和 AI 驅動的 spec 影響分析

**Non-Goals:**

- 不遷移現有 archive 目錄（已歸檔的資料不搬）
- 不改變 mysd 的 phase gate workflow
- 不修改 spectra 的行為

## Decisions

### Archive 路徑與日期前綴

`cmd/archive.go` 中的 `archiveDir` 從：
```go
archiveDir := filepath.Join(specsDir, "archive", ws.ChangeName)
```
改為：
```go
archiveDir := filepath.Join(specsDir, "changes", "archive",
    time.Now().Format("2006-01-02") + "-" + ws.ChangeName)
```

使用 Go 的 `time.Now().Format("2006-01-02")` 產生 YYYY-MM-DD 前綴。

### Delta spec merge 引擎

新增 `internal/spec/merge.go`，提供 `MergeSpecs` 函數：

```go
func MergeSpecs(mainSpecPath string, deltaBody string) (string, error)
```

流程：
1. 讀取 main spec 的完整內容
2. 用 `ParseDelta` 解析 delta body（擴展後支援 RENAMED）
3. 按 OpenSpec 規範順序執行合併：RENAMED → REMOVED → MODIFIED → ADDED
4. 回傳合併後的完整 spec 內容

合併策略：
- **RENAMED**：在 main spec 中找到 `### Requirement: <old>` heading，替換為 `### Requirement: <new>`
- **REMOVED**：找到對應 heading，刪除整個 requirement block（到下一個 `### Requirement:` 或 EOF）
- **MODIFIED**：找到對應 heading，用 delta 內容替換整個 requirement block
- **ADDED**：追加到 main spec 的 `## Requirements` 區段末尾

在 `cmd/archive.go` 的 `runArchive` 中，在搬目錄前遍歷 delta specs 並呼叫 `MergeSpecs`。

### RENAMED delta 操作支援

擴展 `internal/spec/schema.go`：
```go
DeltaRenamed DeltaOp = "RENAMED"
```

擴展 `internal/spec/delta.go` 的 `reDeltaHeading` 和 `DetectDeltaOp` 支援 RENAMED。

`ParseDelta` 簽名擴展回傳 renamed：
```go
func ParseDelta(body string) (added, modified, removed []Requirement, renamed []RenamedRequirement)
```

新增 `RenamedRequirement` struct：
```go
type RenamedRequirement struct {
    From string
    To   string
}
```

解析格式：`### FROM: <old>` 後接 `### TO: <new>`。

### Spec frontmatter 擴展

擴展 `SpecFrontmatter` struct：
```go
type SpecFrontmatter struct {
    Name        string     `yaml:"name"`
    Description string     `yaml:"description"`
    Version     string     `yaml:"version"`
    GeneratedBy string     `yaml:"generatedBy"`
    SpecVersion string     `yaml:"spec-version"`
    Capability  string     `yaml:"capability"`
    Delta       DeltaOp    `yaml:"delta"`
    Status      ItemStatus `yaml:"status"`
}
```

`version` 語意：spec 內容版本，新建為 `1.0.0`，每次 MODIFIED 時 minor +1（例如 `1.0.0` → `1.1.0`）。merge 時自動遞增。

`generatedBy` 格式：`mysd v<version>`，從 build 時注入的版本號取得。

### Scenario GIVEN 驗證

在 `internal/verifier/` 新增 scenario 格式驗證函數：
```go
func ValidateScenarioFormat(specBody string) []string
```

解析 `#### Scenario:` 區塊，檢查是否包含 `**GIVEN**`、`**WHEN**`、`**THEN**` 三個關鍵字。缺少任一則回傳 warning。

在 `cmd/spec.go` 的 spec 產生 context 中加入 GIVEN/WHEN/THEN 格式要求。

### Skipped task 標記與 spec 影響分析

Task 標記擴展：
- `- [ ]` — 未完成
- `- [x]` — 已完成
- `- [~]` — 刻意跳過（必須附理由）

`internal/spec/parser.go` 的 task 解析邏輯擴展，識別 `[~]` 標記。`Task` struct 新增 `Skipped bool` 和 `SkipReason string` 欄位。

`cmd/archive.go` 的 gate 檢查改為：所有 task 必須是 `[x]` 或 `[~]`，不允許 `[ ]`。

Spec 影響分析不在 mysd binary 中實作——這部分由 archive skill（SKILL.md）驅動 AI agent 執行：
1. 讀取 `[~]` tasks 和對應的 spec requirements
2. AI 產出 spec 修改建議
3. 呈現 diff 給使用者確認
4. 使用者核可後修改 delta spec
5. 繼續 archive

mysd binary 只負責提供 `mysd archive --analyze-skipped` flag 輸出 skipped tasks 與 requirement 的對應關係（JSON 格式），供 skill 使用。

## Risks / Trade-offs

- [Risk] Delta spec merge 的 heading 匹配可能因格式差異失敗 → 使用 whitespace-insensitive 比較，失敗時 warning 不阻斷
- [Risk] `ParseDelta` 簽名變更影響現有呼叫端 → 呼叫端數量有限，影響可控
- [Risk] GIVEN 驗證可能對現有 specs 產生大量 warning → 不阻斷，只 warning，允許漸進式修正
- [Trade-off] Spec 影響分析放在 skill 層而非 binary 層 → 保持 binary 簡潔，AI 驅動的分析更適合在 skill 中處理
