## 1. State files stored in .mysd directory（實作 spec: State files stored in .mysd directory）

- [ ] 1.1 State files stored in .mysd directory：修改 `internal/state/state.go` 的 `SaveState` 和 `LoadState`，將 STATE.json 路徑從 `specsDir/STATE.json` 改為專案根目錄的 `.mysd/STATE.json`（透過從 specsDir 回推專案根目錄，或新增參數接收專案根目錄）
- [ ] 1.2 加入 backward compatibility：`LoadState` 若 `.mysd/STATE.json` 不存在但 `openspec/STATE.json`（或 `.specs/STATE.json`）存在，則從 legacy 位置讀取
- [ ] 1.3 修改 `internal/roadmap/roadmap.go` 的 `UpdateTracking`，將 `roadmap-tracking.json` 和 `roadmap-timeline.md` 路徑從 specsDir 改為 `.mysd/`
- [ ] 1.4 更新 `internal/state/state_test.go` 測試，驗證新路徑和 backward compatibility
- [ ] 1.5 更新所有呼叫 `SaveState`/`LoadState` 的 cmd 檔案，確保傳入正確的路徑參數
- [ ] 1.6 確認 `go test ./internal/state/ ./internal/roadmap/ ./cmd/` 全部通過

## 2. STATE.json cleanup after archive（實作 spec: STATE.json cleanup after archive）

- [ ] 2.1 STATE.json cleanup after archive：修改 `cmd/archive.go` 的 `runArchive`，在 archive 成功後刪除 `.mysd/STATE.json`（best-effort，失敗時印 warning 到 stderr）
- [ ] 2.2 更新 `cmd/archive_test.go` 和 `cmd/integration_test.go`，驗證 archive 後 STATE.json 被刪除
- [ ] 2.3 確認 `go test ./cmd/` 全部通過

## 3. Conversation context option in plan spec selector（實作 spec: Conversation context option in plan spec selector）

- [ ] 3.1 Conversation context option in plan spec selector：修改 `mysd/skills/plan/SKILL.md` 的 Step 2b，在 spec 選擇器中新增「From conversation context」選項（在 [All] 之後）
- [ ] 3.2 在 Step 2b 中加入 conversation context 處理邏輯：orchestrator 從對話中提取相關需求描述 → 寫入 `changes/<name>/conversation-context.md` → 用 `--from` 傳給 binary
- [ ] 3.3 在 plan pipeline 完成後（Step 7 或結束時）加入 conversation-context.md 暫存檔清理邏輯

## 4. 文件更新

- [ ] 4.1 更新 `README.md` 的 Planning with Context section，加入 conversation context 說明
- [ ] 4.2 更新 `README.zh-TW.md` 同步 conversation context 說明
