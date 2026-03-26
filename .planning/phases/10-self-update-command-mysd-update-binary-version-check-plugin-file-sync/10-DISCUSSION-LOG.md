# Phase 10: Self-Update Command — Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-26
**Phase:** 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync
**Areas discussed:** 版本檢查策略, Binary 更新機制, Plugin 檔案同步範圍, 指令形式與 UX

---

## 版本檢查策略

### 偵測最新版本方式

| Option | Description | Selected |
|--------|-------------|----------|
| GitHub Releases API | 透過 GitHub API 查詢 latest release tag。不需 Go 環境，只需網路。 | ✓ |
| go list -m | 透過 Go module proxy 查詢最新 tag。需要 Go 環境安裝。 | |
| 雙路徑偵測 | 優先用 GitHub API，失敗時 fallback 到 go list。 | |

**User's choice:** GitHub Releases API
**Notes:** 無

### Dev build 處理

| Option | Description | Selected |
|--------|-------------|----------|
| 當作過期，提示更新 | version=="dev" 視為過期，建議更新到最新 release。 | ✓ |
| 跳過檢查，顯示警告 | dev build 無法比較版本，顯示警告並跳過 binary 更新。 | |
| 由 Claude 決定 | Claude 根據實作情況決定。 | |

**User's choice:** 當作過期，提示更新

### API 認證

| Option | Description | Selected |
|--------|-------------|----------|
| 無認證即可 | 更新檢查頻率低，60 req/hr 綽綽有餘。 | ✓ |
| 支援可選 GITHUB_TOKEN | 預設無認證，有 token 則自動使用。 | |

**User's choice:** 無認證即可

### 自動檢查時機

| Option | Description | Selected |
|--------|-------------|----------|
| 只在 /mysd:update 時 | 不干擾正常工作流程，使用者主動執行時才檢查。 | ✓ |
| 每次執行 mysd 時檢查 | 類似 gh CLI，每 24 小時檢查一次。 | |
| SessionStart hook | 透過 Claude Code hook 在每次 session 開始時檢查。 | |

**User's choice:** 只在 /mysd:update 時

### Cache 機制

| Option | Description | Selected |
|--------|-------------|----------|
| 不需要 cache | 手動觸發，每次即時查詢即可。 | ✓ |
| 簡單 cache | cache 結果到本地檔案，設 TTL。 | |

**User's choice:** 不需要 cache

### 錯誤處理

| Option | Description | Selected |
|--------|-------------|----------|
| 顯示錯誤並繼續 plugin sync | binary 版本檢查失敗不阻止 plugin 同步。 | ✓ |
| 顯示錯誤並中止 | 任何網路問題都完全中止。 | |
| 由 Claude 決定 | | |

**User's choice:** 顯示錯誤並繼續 plugin sync

### Semver 比較

| Option | Description | Selected |
|--------|-------------|----------|
| Semver 解析 | 使用 semver 庫解析 v1.0.0 格式，正確比較 major.minor.patch。 | ✓ |
| 字串比較 | 直接比較 tag 字串是否相同。 | |
| 由 Claude 決定 | | |

**User's choice:** Semver 解析

---

## Binary 更新機制

### 更新方式

| Option | Description | Selected |
|--------|-------------|----------|
| 就地替換 | 用 os.Executable() 找到目前 binary 路徑，直接下載新版替換。 | ✓ |
| go install | 執行 go install。需要 Go 環境。 | |
| 只顯示指引 | 不自動下載，印出更新指令讓使用者手動執行。 | |
| 智慧偵測 | 先檢查 Go 環境是否可用，有則 go install，無則下載。 | |

**User's choice:** 就地替換
**Notes:** 使用者提出有的人 binary 路徑不在 go bin 中，os.Executable() 解決此問題

### Windows 處理

| Option | Description | Selected |
|--------|-------------|----------|
| Rename-then-replace | 舊 binary → .old → 寫入新 binary → 刪除 .old。 | ✓ |
| 下載到 temp 再提示 | 下載到 temp 目錄，顯示手動替換指引。 | |
| 由 Claude 決定 | | |

**User's choice:** Rename-then-replace

### Checksum 驗證

| Option | Description | Selected |
|--------|-------------|----------|
| 驗證 SHA256 | 下載 checksums.txt，比對 archive 的 SHA256。 | ✓ |
| 不驗證 | 直接信任 GitHub Releases 的下載。 | |
| 由 Claude 決定 | | |

**User's choice:** 驗證 SHA256

### 平台偵測

| Option | Description | Selected |
|--------|-------------|----------|
| runtime 常數即可 | runtime.GOOS + runtime.GOARCH 直接對應 GoReleaser naming。 | ✓ |
| 由 Claude 決定 | | |

**User's choice:** 用 runtime 常數即可

### 權限處理

| Option | Description | Selected |
|--------|-------------|----------|
| 自動 chmod +x | 下載解壓後自動設定 0755。Windows 上不需要。 | ✓ |
| 由 Claude 決定 | | |

**User's choice:** 自動 chmod +x

### Rollback 機制

| Option | Description | Selected |
|--------|-------------|----------|
| 自動 rollback | rename-then-replace 保留 .old，失敗時自動還原。 | ✓ |
| 不做 rollback | 失敗時顯示錯誤，使用者自行處理。 | |
| 由 Claude 決定 | | |

**User's choice:** 自動 rollback

---

## Plugin 檔案同步範圍

### 同步來源

| Option | Description | Selected |
|--------|-------------|----------|
| GitHub Release artifact | GoReleaser 將 plugin/ 包在 release 中。 | ✓ |
| 從 repo clone/fetch | git 操作獲取 plugin/ 目錄。 | |
| Binary 內嵌 | go:embed 將 plugin 嵌入 binary。 | |

**User's choice:** GitHub Release artifact

### 同步範圍

| Option | Description | Selected |
|--------|-------------|----------|
| commands + agents | 同步兩個目錄。完整 plugin 更新。 | ✓ |
| 只同步 commands | 只同步 SKILL.md orchestrators。 | |
| 由 Claude 決定 | | |

**User's choice:** commands + agents

### 自訂處理

| Option | Description | Selected |
|--------|-------------|----------|
| 直接覆寫 | plugin 檔案是發佈內容，不預期使用者自訂。 | ✓ |
| 備份後覆寫 | 覆寫前將現有檔案備份到 .old/。 | |
| 差異比對後確認 | 顯示差異讓使用者逐檔確認。 | |

**User's choice:** 直接覆寫

### 刪除策略

| Option | Description | Selected |
|--------|-------------|----------|
| 用 manifest 追蹤 | plugin-manifest.json 追蹤官方檔案，比對新舊 manifest 決定刪除。 | ✓ |
| 只新增/更新，不刪除 | 安全優先，不刪除任何檔案。 | |

**User's choice:** 用 manifest，但舊版無 manifest 時 fallback 為只新增/更新不刪除
**Notes:** 使用者提問如何區分舊版遺留和使用者自訂檔案，確認 manifest 是正確方案

---

## 指令形式與 UX

### 實作形式

| Option | Description | Selected |
|--------|-------------|----------|
| SKILL.md + binary 子命令 | binary 新增 mysd update 子命令，SKILL.md 薄 wrapper。 | ✓ |
| 純 SKILL.md orchestrator | 全部在 SKILL.md 中用 Bash 完成。 | |
| 純 binary 子命令 | 只在 binary 中實作，不建 SKILL.md。 | |

**User's choice:** SKILL.md + binary 子命令

### 輸出格式

| Option | Description | Selected |
|--------|-------------|----------|
| JSON 輸出 | Binary 輸出結構化 JSON，SKILL.md 解析後格式化。 | ✓ |
| 人類可讀輸出 | Binary 直接輸出格式化文字。 | |
| 雙模式 | 預設人類可讀，--json 切換 JSON。 | |

**User's choice:** JSON 輸出

### --check 模式

| Option | Description | Selected |
|--------|-------------|----------|
| 支援 --check | mysd update --check 只檢查不更新。 | ✓ |
| 不需要 | 每次都直接更新。 | |
| 由 Claude 決定 | | |

**User's choice:** 支援 --check

### 確認流程

| Option | Description | Selected |
|--------|-------------|----------|
| 預設確認，--force 跳過 | 顯示更新資訊後確認，--force 跳過。 | ✓ |
| 直接更新 | 不問直接更新。 | |
| 由 Claude 決定 | | |

**User's choice:** 預設確認，--force 跳過

### argument-hint

| Option | Description | Selected |
|--------|-------------|----------|
| 加 argument-hint | SKILL.md 使用 argument-hint 欄位。 | ✓ |
| 全部現有指令也補上 | 順便給所有 mysd SKILL.md 加 argument-hint。 | |

**User's choice:** 全部現有 mysd 指令也補上
**Notes:** 此需求已記為 deferred idea（scope creep），Phase 10 只對 mysd-update 加 argument-hint

---

## Claude's Discretion

- Semver 解析的具體實作方式
- GoReleaser plugin 打包方式
- plugin-manifest.json 結構
- JSON 輸出欄位結構
- HTTP client timeout 和重試策略

## Deferred Ideas

- 更新 README.md 到 v1.1 — 文件更新任務
- 全部現有 mysd 指令補上 argument-hint — UI 增強任務
