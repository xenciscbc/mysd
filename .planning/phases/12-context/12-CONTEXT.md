# Phase 12: 加入 context 的 % 數及色條 - Context

**Gathered:** 2026-03-27 (discuss mode)
**Status:** Ready for planning

<domain>
## Phase Boundary

Context 管理功能，包含兩個子系統：

1. **mysd statusline**：建立 mysd 自己的 statusline hook（`mysd-statusline.js`），顯示 model 名稱、當前 change 名稱、目錄、及 Claude context 窗口使用百分比色條。透過 `mysd init` 自動安裝，可用 `/mysd:statusline` 切換開關。

2. **discuss research cache**：`/mysd:discuss` research 完成後自動暫存 4 維度研究結果，下次同一 change 重啟 discuss 時可選擇重用、重新 research 或略過。archive 時自動清除。

不包含：修改 GSD statusline、修改 context-monitor 邏輯、其他 hook 類型。
</domain>

<decisions>
## Implementation Decisions

### Statusline 內容與格式

- **D-01:** 顯示格式為 `{model} │ {change} │ {dir} │ {bar} {pct}%`
  - model：從 `data.model?.display_name` 或 `data.model?.id` 提取**簡稱**（見 D-02）
  - change：從 `.specs/state.yaml` 讀取 `change_name`，若無 active change 則省略此欄
  - dir：`path.basename(data.workspace?.current_dir || process.cwd())`
  - bar：10 段 `█░` 色條 + 百分比數字

- **D-02:** Model 簡稱提取邏輯（優先順序）：
  1. 從 display_name 或 model id 中比對關鍵字 → 對應簡稱：
     - 含 `opus` → `opus`
     - 含 `sonnet` → `sonnet`
     - 含 `haiku` → `haiku`
  2. 若無匹配：取 display_name 第一個空格前的部分，或 fallback `claude`

  範例輸出：
  - 有 change：`sonnet │ my-change │ mysd │ █████░░░░░ 50%`
  - 無 change：`sonnet │ mysd │ █████░░░░░ 50%`

### Context 資料計算

- **D-03:** 直接從 stdin JSON 讀取 `data.context_window?.remaining_percentage`，套用 GSD 相同的 normalization 邏輯：
  - `AUTO_COMPACT_BUFFER_PCT = 16.5`
  - `usableRemaining = max(0, (remaining - 16.5) / 83.5 * 100)`
  - `used = max(0, min(100, round(100 - usableRemaining)))`

- **D-04:** Bridge file 寫入邏輯：**只在 GSD 並存時寫入，不並存時完全不寫**（與 statusline_enabled 無關）。
  - GSD 並存偵測：檢查以下任一路徑是否存在 `gsd-context-monitor.js`：
    1. `{workspace}/.claude/hooks/gsd-context-monitor.js`
    2. `~/.claude/hooks/gsd-context-monitor.js`（`CLAUDE_CONFIG_DIR` 或 homedir）
  - 並存 → 寫入 bridge file（`/tmp/claude-ctx-{session}.json`），格式與 GSD 相同：`{ session_id, remaining_percentage, used_pct, timestamp }`
  - 不並存 → 跳過，不寫入
  - 寫入為 best-effort（silent fail）

- **D-05:** 顏色閾值（改自 GSD，移除骷髏頭換成 🥵）：
  - `used < 50%`：綠色 `\x1b[32m`
  - `used < 65%`：黃色 `\x1b[33m`
  - `used < 80%`：橙色 `\x1b[38;5;208m`
  - `used >= 80%`：紅色閃爍 + 🥵 `\x1b[5;31m`（不用 💀，改用 🥵）

### 安裝機制

- **D-06:** `mysd init` 負責安裝，執行兩個動作：
  1. 將 `plugin/hooks/mysd-statusline.js` 複製到 `.claude/hooks/mysd-statusline.js`
  2. 寫入/更新 `.claude/settings.json`，加入 `"statusLine": { "type": "command", "command": "node .claude/hooks/mysd-statusline.js" }`

- **D-07:** `.claude/settings.json` 的寫入為 merge 操作（保留既有 keys，只新增/覆蓋 `statusLine` 欄位）。若 settings.json 不存在則新建。

- **D-08:** hook 原始檔位於 `plugin/hooks/mysd-statusline.js`，這是 plugin distribution 的 source of truth。

### Statusline On/Off 控制

- **D-11:** `ProjectConfig` 新增 `statusline_enabled: true`（預設開啟），存於 `.claude/mysd.yaml`。

- **D-12:** `mysd-statusline.js` 啟動時讀取 `{workspace}/.claude/mysd.yaml` 的 `statusline_enabled`：
  - `true`（或欄位不存在）→ 正常顯示 statusline
  - `false` → 不輸出（空字串），但仍執行 D-04 的 GSD 並存偵測與 bridge file 寫入

- **D-13:** 新增 `/mysd:statusline` SKILL.md 指令（對應 `mysd statusline` binary 子指令）：
  - 無參數 → toggle（讀當前值，寫入相反值，並顯示新狀態）
  - `on` → 設為 true
  - `off` → 設為 false
  - 顯示當前狀態：`Statusline: on` 或 `Statusline: off`

### Plugin Sync

- **D-09:** `mysd update`（Phase 10 的 PluginManifest）目前只 track commands + agents。Phase 12 **不** 擴展 manifest 去 track hooks（hooks 是 one-time init，不需要 self-update sync）。

### State.yaml 讀取

- **D-10:** 讀取 change 名稱時，路徑為 `.specs/state.yaml`（相對於 `data.workspace?.current_dir`）。YAML key 為 `change_name`。使用簡單 regex 或 line parsing 讀取（不引入 YAML library，statusline hook 只用 Node stdlib）。若讀取失敗 silent fail，省略 change 欄位。

### Discuss Research Cache

- **D-14:** Cache 檔案位置：`.specs/changes/{change_name}/discuss-research-cache.json`（與 proposal.md、tasks.md 同層）

- **D-15:** Cache 內容：4 維度研究完整輸出（architecture / codebase / ux / security），加上 metadata：
  ```json
  {
    "change_name": "...",
    "cached_at": "2026-03-27T03:00:00Z",
    "research": {
      "architecture": "...",
      "codebase": "...",
      "ux": "...",
      "security": "..."
    }
  }
  ```

- **D-16:** 寫入時機：`/mysd:discuss` 中 research step 完成後**立即寫入**（proactive，非等中止才存）。若 cache 已存在則覆蓋。

- **D-17:** Discuss 啟動時的 cache 偵測邏輯：
  1. 檢查 `.specs/changes/{change}/discuss-research-cache.json` 是否存在
  2. 若存在 → 顯示 cache 時間，詢問使用者：
     - **重用** — 跳過 research step，直接載入 cache 結果
     - **重新 research** — 執行 research，完成後覆蓋 cache
     - **都不要** — 跳過 research step，cache 保持不動
  3. 若不存在 → 正常流程（research 後自動寫入）

- **D-18:** Archive 清除：`mysd archive` 執行時，若 `.specs/changes/{change}/discuss-research-cache.json` 存在則刪除（silent fail if not found）

- **D-19:** Cache 檔案加入 `.gitignore`（`discuss-research-cache.json`），不納入版控。

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

### Discuss Research Cache
- `plugin/commands/mysd-discuss.md` — discuss SKILL.md，需加入 cache 偵測、寫入、重用邏輯
- `cmd/archive.go` — archive 指令，需加入 cache 清除邏輯
- `.specs/changes/` — cache 檔案存放路徑範例
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
