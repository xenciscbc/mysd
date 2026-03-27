# Phase 12: 加入 context 的 % 數及色條 - Context

**Gathered:** 2026-03-27 (discuss mode)
**Status:** Ready for planning

<domain>
## Phase Boundary

建立 mysd 自己的 statusline hook（`mysd-statusline.js`），顯示 model 名稱、當前 change 名稱、目錄、及 Claude context 窗口使用百分比色條。透過 `mysd init` 自動安裝（複製 hook 檔案 + 寫入 `.claude/settings.json`）。

不包含：修改 GSD statusline、新增其他 hook 類型、修改 context-monitor 邏輯。
</domain>

<decisions>
## Implementation Decisions

### Statusline 內容與格式

- **D-01:** 顯示格式為 `{model} │ {change} │ {dir} │ {bar} {pct}%`
  - model：`data.model?.display_name || 'Claude'`
  - change：從 `.specs/state.yaml` 讀取 `change_name`，若無 active change 則省略此欄
  - dir：`path.basename(data.workspace?.current_dir || process.cwd())`
  - bar：10 段 `█░` 色條 + 百分比數字

- **D-02:** 若有 active change：`claude-sonnet-4-5 │ my-change │ mysd │ █████░░░░░ 50%`
  若無 active change：`claude-sonnet-4-5 │ mysd │ █████░░░░░ 50%`

### Context 資料計算

- **D-03:** 直接從 stdin JSON 讀取 `data.context_window?.remaining_percentage`，套用 GSD 相同的 normalization 邏輯：
  - `AUTO_COMPACT_BUFFER_PCT = 16.5`
  - `usableRemaining = max(0, (remaining - 16.5) / 83.5 * 100)`
  - `used = max(0, min(100, round(100 - usableRemaining)))`

- **D-04:** 計算完後，**寫入** GSD bridge file（`/tmp/claude-ctx-{session}.json`），使 GSD context-monitor hook 繼續正常運作（即使 GSD statusline 被 mysd statusline 取代）。Bridge file 寫入為 best-effort（silent fail）。

- **D-05:** 顏色閾值與 GSD 相同：
  - `used < 50%`：綠色 `\x1b[32m`
  - `used < 65%`：黃色 `\x1b[33m`
  - `used < 80%`：橙色 `\x1b[38;5;208m`
  - `used >= 80%`：紅色閃爍 + 💀 `\x1b[5;31m`

### 安裝機制

- **D-06:** `mysd init` 負責安裝，執行兩個動作：
  1. 將 `plugin/hooks/mysd-statusline.js` 複製到 `.claude/hooks/mysd-statusline.js`
  2. 寫入/更新 `.claude/settings.json`，加入 `"statusLine": { "type": "command", "command": "node .claude/hooks/mysd-statusline.js" }`

- **D-07:** `.claude/settings.json` 的寫入為 merge 操作（保留既有 keys，只新增/覆蓋 `statusLine` 欄位）。若 settings.json 不存在則新建。

- **D-08:** hook 原始檔位於 `plugin/hooks/mysd-statusline.js`，這是 plugin distribution 的 source of truth。

### Plugin Sync

- **D-09:** `mysd update`（Phase 10 的 PluginManifest）目前只 track commands + agents。Phase 12 **不** 擴展 manifest 去 track hooks（hooks 是 one-time init，不需要 self-update sync）。

### State.yaml 讀取

- **D-10:** 讀取 change 名稱時，路徑為 `.specs/state.yaml`（相對於 `data.workspace?.current_dir`）。YAML key 為 `change_name`。使用簡單 regex 或 line parsing 讀取（不引入 YAML library，statusline hook 只用 Node stdlib）。若讀取失敗 silent fail，省略 change 欄位。

### Claude's Discretion

- Stdin timeout guard 邏輯（沿用 GSD 的 3s timeout 模式）
- 錯誤處理細節（所有錯誤 silent fail，不破壞 statusline）
- YAML 解析方式（regex line scan 即可）
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### GSD Statusline（參考實作）
- `.claude/hooks/gsd-statusline.js` — 完整的 GSD statusline 實作，mysd-statusline 的設計基礎

### GSD Context Monitor（相容性目標）
- `.claude/hooks/gsd-context-monitor.js` — bridge file 的消費者，需確保 mysd statusline 寫入格式相容

### Init Command（安裝目標）
- `cmd/init_cmd.go` — `mysd init` 的實作，Phase 12 需擴展此命令
- `internal/config/defaults.go` — config 結構參考

### Plugin Distribution
- `plugin/hooks/hooks.json` — hooks 目錄結構
- `internal/update/manifest.go` — PluginManifest 結構（確認不需要修改）
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `gsd-statusline.js`：完整參考實作，直接 port 大部分邏輯（stdin 讀取、normalization、bar 建立、顏色、bridge file 寫入）
- `gsd-context-monitor.js`：bridge file 的 schema（`{ session_id, remaining_percentage, used_pct, timestamp }`）
- `cmd/init_cmd.go`：init 命令框架，已有 `.claude/` 目錄建立邏輯，擴展即可
- `internal/output/`：Printer 模式（用於 init 的 success/error 輸出）

### Established Patterns
- Settings.json 寫入：參考 `cmd/model.go` 的 viper 讀寫模式，或直接用 `encoding/json` merge
- Plugin hook 檔案：`plugin/hooks/` 目錄已存在（有 `hooks.json`），新增 JS 檔案即可
- Silent fail pattern：GSD hook 的所有 FS/parse 錯誤都 silent fail，mysd 保持相同風格

### Integration Points
- `mysd init` 輸出新增一行：`"Statusline configured. Context bar will appear in Claude Code."`
- `.claude/settings.json` 的 `statusLine` 欄位（project-level 覆蓋 global GSD statusline）
- Bridge file：`path.join(os.tmpdir(), \`claude-ctx-\${session}.json\`)` — GSD context-monitor 讀取此檔
</code_context>

<specifics>
## Specific Ideas

- 格式完全對齊 GSD style（`│` 分隔符、`█░` 字元、ANSI 顏色碼），使用者在 mysd 專案中看到的 statusline 與 GSD 一致
- 如果 `.specs/state.yaml` 不存在或讀取失敗，靜默省略 change 欄位（不顯示 undefined 或 error）
- `mysd init` 是 idempotent：重複執行不會造成問題（覆蓋 statusLine 設定是可接受的）
</specifics>

<deferred>
## Deferred Ideas

- 擴展 PluginManifest 去 track hooks（D-09 決定不做，hooks 是 one-time init）
- 加入 task 名稱顯示（類似 GSD 從 todos JSON 讀取）— 可在未來 phase 加入
- `mysd update` 自動更新 hook 檔案 — 未來 milestone

None — discussion stayed within phase scope
</deferred>

---

*Phase: 12-context*
*Context gathered: 2026-03-27*
