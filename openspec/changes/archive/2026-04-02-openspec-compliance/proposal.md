## Summary

mysd 的 archive 流程和 spec 格式與 OpenSpec（Fission-AI/OpenSpec）規範存在多處偏差，需要修正以確保相容性。

## Motivation

mysd 參考 OpenSpec 規範管理 `openspec/` 目錄結構和 spec 格式。經過系統性比對，發現以下偏差：

1. **Archive 路徑錯誤** — mysd 歸檔到 `openspec/archive/<name>/`，規範要求 `openspec/changes/archive/YYYY-MM-DD-<name>/`
2. **缺少 delta spec merge** — 歸檔時不合併 delta specs 回 main specs，導致 `openspec/specs/` 逐漸失真
3. **Scenario 格式不完整** — 缺少 GIVEN 前置條件，規範要求 GIVEN/WHEN/THEN
4. **Spec 缺少 frontmatter** — 規範要求 name, description, version, generatedBy 等 YAML frontmatter
5. **無法處理不可完成的 task** — 沒有 skipped task 機制，遇到不可行的 task 時 archive 會被卡住

## Proposed Solution

- 修正 archive 路徑和日期前綴
- 實作 delta spec merge 邏輯（RENAMED → REMOVED → MODIFIED → ADDED）
- spec 產生和驗證加入 GIVEN 支援
- spec 產生時加入 YAML frontmatter
- 引入 `[~]` skipped task 標記，archive 時 AI 驅動分析 skipped tasks 對 spec 的影響，人確認後調整

## Alternatives Considered

- 依賴 spectra CLI 做 delta spec merge → 拒絕，mysd 是獨立工具，不應依賴外部 CLI
- `--mark-tasks-complete` 強制標完 → 拒絕，不誠實，`[~]` 標記附理由更嚴謹
- 放鬆 phase gate 改為 fluid → 拒絕，嚴格 gate 是 mysd 的設計選擇，不算偏差

## Impact

- Affected specs: 所有現有 specs 需要補 frontmatter（可漸進式處理）
- Affected code:
  - `cmd/archive.go` — archive 路徑、日期前綴、delta spec merge、skipped task 處理
  - `cmd/archive_test.go` — 對應測試更新
  - `internal/spec/` — delta spec parser、merge 邏輯、frontmatter 處理
  - `internal/verifier/` — GIVEN 驗證
  - `cmd/spec.go` — spec 產生時加入 GIVEN 和 frontmatter
  - `cmd/task-update.go` 或相關 — `[~]` skipped 標記支援
  - `cmd/integration_test.go` — 整合測試路徑更新

## Capabilities

### New Capabilities

- `delta-spec-merge`: 在 archive 時自動合併 delta specs 回 main specs，支援 RENAMED/REMOVED/MODIFIED/ADDED 四種操作，按規範順序執行
- `skipped-task-handling`: 引入 `[~]` skipped task 標記，archive gate 視為已處理，AI 驅動分析 skipped tasks 對 spec requirement 的影響並由人確認

### Modified Capabilities

(none)
