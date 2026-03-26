# Phase 10: Self-Update Command — /mysd:update - Context

**Gathered:** 2026-03-26
**Status:** Ready for planning

<domain>
## Phase Boundary

提供 `/mysd:update` 指令，實現統一的自我更新機制：binary 版本檢查 + 下載替換 + plugin 檔案同步。使用者透過一個指令即可將 mysd binary 和 plugin 檔案（commands + agents）更新到最新版本。

</domain>

<decisions>
## Implementation Decisions

### 版本檢查策略
- **D-01:** 使用 GitHub Releases API 偵測最新版本（`/repos/{owner}/{repo}/releases/latest`）
- **D-02:** 本地版本為 `dev`（ldflags 未注入）時視為過期，提示更新到最新 release
- **D-03:** 無需認證 — GitHub API 未認證 rate limit（60 req/hr）對手動觸發已足夠
- **D-04:** 只在使用者主動執行 `/mysd:update` 時檢查版本，不做自動背景檢查或 SessionStart hook
- **D-05:** 不 cache API 查詢結果 — 手動觸發頻率低，每次即時查詢即可
- **D-06:** 網路失敗時顯示錯誤但繼續 plugin sync — binary 版本檢查失敗不阻止 plugin 同步
- **D-07:** 使用 semver 解析比較版本號（需 semver library 或自行解析 major.minor.patch）

### Binary 更新機制
- **D-08:** 用 `os.Executable()` 找到目前 binary 路徑，從 GitHub Release 下載對應平台的 archive 後就地替換。不依賴 Go 環境
- **D-09:** Windows 使用 rename-then-replace 模式：舊 binary → `.old` → 寫入新 binary → 刪除 `.old`（Windows 不允許覆寫執行中的 binary）
- **D-10:** 下載後驗證 SHA256 checksum — 從 release 的 `checksums.txt` 取得預期 hash 比對
- **D-11:** 用 `runtime.GOOS` + `runtime.GOARCH` 對應 GoReleaser 命名樣板 `mysd_{version}_{os}_{arch}`
- **D-12:** Linux/macOS 下載解壓後自動設定 0755 執行權限
- **D-13:** 更新失敗時自動 rollback — 將 `.old` rename 回原位，恢復到更新前狀態

### Plugin 檔案同步
- **D-14:** Plugin 檔案來源為 GitHub Release artifact（需修改 .goreleaser.yaml 將 plugin/ 包含在 release 中）
- **D-15:** 同步 commands + agents 兩個目錄（plugin/commands/*.md → .claude/commands/、plugin/agents/*.md → .claude/agents/）
- **D-16:** 本地已存在的同名檔案直接覆寫，不保留使用者自訂內容
- **D-17:** 使用 `plugin-manifest.json` 追蹤官方檔案清單。同步時比對新舊 manifest：新版有但舊版沒有 → 新增；兩者都有 → 更新；舊版有但新版沒有 → 刪除（被官方移除的檔案）。舊版無 manifest（v1.1 之前）→ fallback 為只新增/更新，不刪除

### 指令形式與 UX
- **D-18:** SKILL.md 薄 wrapper + binary `mysd update` 子命令。與現有 model/lang/note 指令模式一致
- **D-19:** Binary 輸出結構化 JSON（版本資訊、更新狀態、同步結果）。SKILL.md 解析後格式化顯示
- **D-20:** 支援 `mysd update --check` — 只檢查版本和可用更新，不執行更新
- **D-21:** 預設顯示更新資訊後要求使用者確認，`--force` 跳過確認直接更新
- **D-22:** SKILL.md 使用 `argument-hint: "[--check] [--force]"` 欄位顯示參數提示

### Claude's Discretion
- Semver 解析的具體實作方式（標準庫自行解析 or 第三方 library）
- GoReleaser 如何包含 plugin/ 目錄（extra_files or 獨立 archive）
- plugin-manifest.json 的具體結構
- JSON 輸出的欄位結構
- HTTP client 的 timeout 和重試策略

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Build & Distribution
- `.goreleaser.yaml` — GoReleaser v2 設定，定義跨平台 build 和 archive 命名樣板
- `main.go` — version/commit/date ldflags 變數，SetVersion 呼叫入口

### Existing Command Patterns
- `cmd/root.go` — Cobra root command 設定，SetVersion()，initConfig()
- `cmd/model.go` — 現有 binary 子命令模式（model set/get）
- `cmd/lang.go` — 現有 binary 子命令模式（lang set）
- `cmd/note.go` — 現有 binary 子命令模式（note add/delete/list）

### Plugin Structure
- `.claude/commands/mysd-note.md` — 現有薄 wrapper SKILL.md 範例（Bash + Read only）
- `.claude/commands/mysd-model.md` — 現有薄 wrapper SKILL.md 範例
- `plugin/commands/` — Distribution copy 目錄（19 個 SKILL.md commands）
- `plugin/agents/` — Distribution copy 目錄（12 個 agent definitions）

### Prior Decisions
- `.planning/phases/08-skill-md-orchestrators-agent-definitions/08-CONTEXT.md` — Plugin sync pattern: .claude/ is authoritative, plugin/ is distribution

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `cmd.SetVersion(v string)` — 已有版本設定機制，version 透過 ldflags 注入
- `rootCmd.Version` — Cobra 內建版本輸出（`mysd --version`）
- `.goreleaser.yaml` — 跨平台 build 設定已就位（linux/darwin/windows × amd64/arm64）
- `internal/config/` — Viper 配置管理已建立
- `internal/spec/deferred.go` — JSON 檔案讀寫模式（可參考 manifest 存取）

### Established Patterns
- **JSON 輸出模式:** `--context-only` flag 輸出 JSON 供 SKILL.md 消費（cmd/plan.go, cmd/execute.go）
- **子命令模式:** Cobra subcommand 加 RunE 函式（cmd/model.go, cmd/lang.go）
- **薄 wrapper SKILL.md:** Bash 呼叫 binary → 讀取 JSON 輸出 → 格式化顯示（mysd-note.md, mysd-model.md）
- **Convention over config:** 預設即好用，不需配置（cmd/root.go initConfig）

### Integration Points
- `cmd/root.go init()` — 新增 `updateCmd` 子命令
- `.claude/commands/` — 新增 `mysd-update.md` SKILL.md
- `plugin/commands/` — 新增 `mysd-update.md` distribution copy
- `.goreleaser.yaml` — 新增 `extra_files` 包含 plugin/ 和 manifest

</code_context>

<specifics>
## Specific Ideas

- GoReleaser 命名樣板 `mysd_{version}_{os}_{arch}` 必須與 runtime 常數精確對應
- Windows rename-then-replace 是標準 Go CLI self-update 模式（參考 gh CLI、gum 等）
- plugin-manifest.json 需在 GoReleaser build 流程中自動產生，不應手動維護
- `--check` 模式的 JSON 應包含 `update_available: bool`、`current_version`、`latest_version` 欄位

</specifics>

<deferred>
## Deferred Ideas

- **更新 README.md 到 v1.1** — 文件更新任務，不在 self-update command 範圍內
- **全部現有 mysd 指令補上 argument-hint** — 增強所有 SKILL.md 的 UI 提示，可作為獨立 quick task

</deferred>

---

*Phase: 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync*
*Context gathered: 2026-03-26*
