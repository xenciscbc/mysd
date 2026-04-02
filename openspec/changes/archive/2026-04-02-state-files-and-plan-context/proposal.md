## Why

mysd 的工作狀態檔（STATE.json、roadmap-tracking.json、roadmap-timeline.md）目前散落在 `openspec/` 根目錄，與 OpenSpec 的 spec 內容混在一起。這造成 `openspec/` 目錄不乾淨，且 STATE.json 在 archive 後殘留為過期狀態。此外，`/mysd:plan` 目前只能從 spec 檔案或 `--from` 外部檔案取得 context，無法從對話中已有的需求描述直接轉換為 tasks。

## What Changes

- 將 `STATE.json`、`roadmap-tracking.json`、`roadmap-timeline.md` 的儲存路徑從 `openspec/` 搬遷到 `.mysd/`
- `mysd archive` 完成後自動刪除 `.mysd/STATE.json`
- 在 plan SKILL.md 的 Step 2b spec 選擇器中新增「From conversation context」選項，讓使用者能從對話 context 提取需求並轉換為 tasks（透過 orchestrator 寫入暫存檔 + 複用 `--from` 機制）
- 更新 README.md 和 README.zh-TW.md 反映上述變更

## Non-Goals

- 不搬遷 `discuss-research-cache.json` — 維持在 `changes/<name>/` 中，隨 change 生死
- 不搬遷 `openspec/config.yaml` — 這是 OpenSpec 配置，屬於 spec 內容
- 不搬遷 `deferred.json` — 延遲筆記跨 change 存在，目前位置合理
- 不修改 `--from` 的 binary 實作 — conversation context 純粹在 SKILL.md 層處理

## Capabilities

### New Capabilities

- `plan-conversation-context`: plan SKILL.md 的 spec 選擇器新增從對話 context 提取需求的選項

### Modified Capabilities

(none)

## Impact

- Affected code:
  - `internal/state/state.go` — STATE.json 讀寫路徑從 `specsDir/STATE.json` 改為 `.mysd/STATE.json`
  - `internal/roadmap/roadmap.go` — tracking/timeline 路徑改為 `.mysd/`
  - `cmd/archive.go` — archive 完成後刪除 `.mysd/STATE.json`
  - `internal/state/state_test.go` — 更新測試
  - `internal/roadmap/roadmap_test.go` — 更新測試（如存在）
  - `cmd/archive_test.go` — 更新測試
  - `cmd/integration_test.go` — 更新測試
- Affected skills:
  - `mysd/skills/plan/SKILL.md` — Step 2b 新增 conversation context 選項
- Affected docs:
  - `README.md` — 更新 plan 功能描述
  - `README.zh-TW.md` — 同步更新
