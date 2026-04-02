## 1. Archive 路徑與日期前綴

- [x] 1.1 修改 archive path with date prefix：`cmd/archive.go` 的 `archiveDir` 改為 `filepath.Join(specsDir, "changes", "archive", time.Now().Format("2006-01-02")+"-"+ws.ChangeName)`，引入 `"time"` import
- [x] 1.2 更新 `cmd/archive_test.go` 和 `cmd/integration_test.go` 中所有 archive 路徑斷言，改為 `changes/archive/YYYY-MM-DD-<name>/` 格式

## 2. RENAMED delta 操作支援

- [x] 2.1 在 `internal/spec/schema.go` 新增 `DeltaRenamed DeltaOp = "RENAMED"` 和 `RenamedRequirement` struct（含 From/To 欄位），實作 RENAMED delta operation support
- [x] 2.2 擴展 `internal/spec/delta.go` 的 `reDeltaHeading` 和 `DetectDeltaOp` 支援 RENAMED，修改 `ParseDelta` 簽名回傳 renamed 切片，解析 `### FROM:` / `### TO:` 格式
- [x] 2.3 為 RENAMED 解析撰寫單元測試（parse RENAMED section、多對 FROM/TO、混合 ADDED+RENAMED）

## 3. Delta spec merge 引擎

- [x] 3.1 新增 `internal/spec/merge.go`，實作 delta spec merge on archive 的 `MergeSpecs(mainSpecPath, deltaBody string) (string, error)` 函數，按 delta merge operation order（RENAMED → REMOVED → MODIFIED → ADDED）執行合併
- [x] 3.2 實作 merge failure handling：heading 不匹配時 emit warning 並跳過，不阻斷
- [x] 3.3 為 `MergeSpecs` 撰寫單元測試：merge ADDED/MODIFIED/REMOVED/RENAMED requirements、no matching main spec exists、heading mismatch during MODIFIED
- [x] 3.4 在 `cmd/archive.go` 的 `runArchive` 中，搬目錄前遍歷 `changes/<name>/specs/` 下的 delta specs，呼叫 `MergeSpecs` 合併回 `openspec/specs/`

## 4. Spec frontmatter 擴展

- [x] 4.1 擴展 `internal/spec/schema.go` 的 `SpecFrontmatter`，新增 `Name`、`Description`、`Version`、`GeneratedBy` 欄位，實作 spec frontmatter on merge
- [x] 4.2 在 `MergeSpecs` 中處理 frontmatter：MODIFIED 時 version increment on MODIFIED（minor +1），新建時 new spec gets initial frontmatter（`version: 1.0.0`、`generatedBy: mysd v<version>`）
- [x] 4.3 為 frontmatter 處理撰寫單元測試

## 5. Scenario GIVEN 驗證

- [x] 5.1 新增 scenario GIVEN validation 函數於 `internal/verifier/`：`ValidateScenarioFormat(specBody string) []string`，檢查每個 `#### Scenario:` 區塊是否包含 GIVEN/WHEN/THEN，scenario missing GIVEN 時 emit warning，complete scenario passes validation
- [x] 5.2 為 `ValidateScenarioFormat` 撰寫單元測試

## 6. Skipped task 標記

- [x] 6.1 實作 skipped task 標記與 spec 影響分析的解析端：擴展 `internal/spec/parser.go` 的 task 解析邏輯，識別 skipped task marker `[~]`，`Task` struct 新增 `Skipped bool` 和 `SkipReason string` 欄位，parse skipped task with reason、skipped task without reason emit warning
- [x] 6.2 修改 `cmd/archive.go` 的 gate 檢查，實作 archive gate accepts skipped tasks：`[x]` 和 `[~]` 視為已處理，只有 `[ ]` 阻斷
- [x] 6.3 新增 `--analyze-skipped` flag，實作 skipped task spec impact analysis output：輸出 JSON 格式的 skipped tasks 與 requirement 對應關係，analyze skipped tasks 和 no skipped tasks 兩種場景
- [x] 6.4 為 skipped task 解析、archive gate、`--analyze-skipped` 撰寫單元測試

## 7. 整合測試

- [x] 7.1 撰寫端到端整合測試：建立 change 含 delta specs → archive → 驗證 main specs 已合併、archive 路徑正確、skipped tasks 處理正確
