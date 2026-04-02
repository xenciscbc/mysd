## 1. UAT prompt removal（實作 spec: UAT prompt removal）

- [x] 1.1 UAT prompt removal：移除 `cmd/archive.go` 的 `runArchiveCmd` 中互動式 UAT prompt 區塊（`isInteractive()` 檢查、`bufio.NewReader`、`Run UAT first?` 提示），移除不再使用的 `isInteractive` 函數和 `bufio` import
- [x] 1.2 移除 `cmd/archive_test.go` 中 UAT 相關測試（`TestArchiveGateNoUAT`），更新 `setupArchiveTestChange` 中與 UAT 相關的註解
- [x] 1.3 移除 `cmd/integration_test.go` 中 `TestArchiveIntegration_NoUATCheck` 測試
- [x] 1.4 確認 `go test ./cmd/` 全部通過

## 2. Archive SKILL.md path accuracy（實作 spec: Archive SKILL.md path accuracy）

- [x] 2.1 Archive SKILL.md path accuracy：修正 `mysd/skills/archive/SKILL.md` Step 1 成功訊息的 archive 路徑，從 `.specs/archive/{change_name}/` 改為 `openspec/changes/archive/YYYY-MM-DD-{change_name}/`
- [x] 2.2 修正 `mysd/skills/archive/SKILL.md` Step 2b 的 context 讀取路徑，從 `.specs/archive/{change_name}/` 改為 `openspec/changes/archive/YYYY-MM-DD-{change_name}/`

## 3. Standalone docs update skill（實作 spec: Standalone docs update skill）

- [x] 3.1 Standalone docs update skill：建立 `mysd/skills/docs-update/SKILL.md`，實作 `/mysd:docs update` skill 的基礎框架：argument 解析（無參數、`--change`、`--last N`、`--full`、自由文字）
- [x] 3.2 實作預設 scope — 讀取 `openspec/changes/archive/` 最近一次 archived change，組合 proposal + tasks + specs 作為 update context
- [x] 3.3 實作 `--change <name>` scope — 在 archive 目錄中定位指定 change
- [x] 3.4 實作 `--last N` scope — 按日期前綴排序取最近 N 個 archived changes，合併 context
- [x] 3.5 實作 `--full` scope — 掃描 codebase 現有程式碼、commands、config 作為 update context
- [x] 3.6 實作自由文字 scope — 使用者提供的描述文字直接作為 update context
- [x] 3.7 實作 docs_to_update 檢查 — 無設定時提示使用者，有設定時對每個檔案套用 archive SKILL.md Step 2c 同樣的更新策略（CHANGELOG prepend、README full rewrite、其他 auto-detect）
- [x] 3.8 更新 `mysd/skills/docs/SKILL.md` 的說明，加入 `/mysd:docs update` 的使用提示
