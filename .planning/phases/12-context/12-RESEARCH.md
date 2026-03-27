# Phase 12: 加入 context 的 % 數及色條 - Research

**Researched:** 2026-03-27
**Domain:** Node.js hook scripting (statusline), Go CLI extension (statusline subcommand + archive extension), Claude Code settings.json merge
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Statusline 內容與格式**
- D-01: 顯示格式 `{model} │ {change} │ {dir} │ {bar} {pct}%`（無 change 時省略 change 欄）
- D-02: Model 簡稱提取：含 `opus`/`sonnet`/`haiku` 關鍵字對應簡稱；否則取 display_name 第一個空格前；fallback `claude`
- D-03: Context 計算公式：`AUTO_COMPACT_BUFFER_PCT = 16.5`；`usableRemaining = max(0,(remaining-16.5)/83.5*100)`；`used = max(0,min(100,round(100-usableRemaining)))`
- D-04: Bridge file 只在 GSD 並存時寫入（偵測 `gsd-context-monitor.js` 任一路徑）；不並存則完全不寫
- D-05: 顏色閾值：`<50%` 綠、`<65%` 黃、`<80%` 橙 `\x1b[38;5;208m`、`>=80%` 紅閃爍 + 🥵（不用 💀）
- D-06: `mysd init` 安裝：複製 `plugin/hooks/mysd-statusline.js` → `.claude/hooks/mysd-statusline.js`；寫入 settings.json `statusLine` 欄位
- D-07: settings.json 寫入為 merge 操作（保留既有 keys）；不存在則新建
- D-08: hook 原始檔位於 `plugin/hooks/mysd-statusline.js`（distribution source of truth）
- D-09: Phase 12 不擴展 PluginManifest 去 track hooks
- D-10: `.specs/state.yaml` 讀取 `change_name`；使用 regex/line-scan（不引入 YAML library）；silent fail
- D-11: `ProjectConfig` 新增 `statusline_enabled: true`（預設開啟），存於 `.claude/mysd.yaml`
- D-12: hook 啟動時讀 `statusline_enabled`：false → 不輸出（空字串），但仍執行 D-04 bridge file 邏輯
- D-13: 新增 `/mysd:statusline` SKILL.md（對應 `mysd statusline` binary 子指令）：無參數 toggle；`on`/`off` 設定；顯示當前狀態

**Discuss Research Cache**
- D-14: Cache 位置：`.specs/changes/{change_name}/discuss-research-cache.json`
- D-15: Cache 內容：`{ change_name, cached_at, research: { architecture, codebase, ux, security } }`
- D-16: 寫入時機：research step 完成後立即寫入（proactive）；已存在則覆蓋
- D-17: Discuss 啟動偵測：存在 → 顯示時間 + 三選一（重用/重新/都不要）；不存在 → 正常流程後寫入
- D-18: `mysd archive` 執行時刪除 cache（silent fail if not found）
- D-19: `discuss-research-cache.json` 加入 `.gitignore`

### Claude's Discretion
- Stdin timeout guard 邏輯（沿用 GSD 的 3s timeout 模式）
- 錯誤處理細節（所有錯誤 silent fail，不破壞 statusline）
- YAML 解析方式（regex line scan 即可）

### Deferred Ideas (OUT OF SCOPE)
- 擴展 PluginManifest 去 track hooks
- 加入 task 名稱顯示
- `mysd update` 自動更新 hook 檔案
</user_constraints>

---

## Summary

Phase 12 包含兩個獨立子系統：**mysd-statusline hook**（Node.js）和 **discuss research cache**（Go + SKILL.md）。

**Statusline 子系統**的核心工作是將 GSD 的 `gsd-statusline.js` port 為 mysd 版本，差異點在於：加入 `change_name` 欄位（從 `.specs/state.yaml` 讀取）、改用 mysd 的 model 簡稱邏輯、bridge file 只在 GSD 並存時才寫入、顏色閾值使用 🥵 而非 💀。安裝面需要擴展 `mysd init`（Go code）寫入 `.claude/settings.json` 的 `statusLine` 欄位，以及新增 `mysd statusline` 子指令控制開關。

**Discuss research cache 子系統**完全在 SKILL.md 層（`mysd-discuss.md`）和 Go binary（`cmd/archive.go`）實作，不需要新的 Go package——只需修改既有的 `runArchive` 函式加入 cache 刪除邏輯，並修改 `mysd-discuss.md` 加入 Step 4.5（cache 偵測）與 Step 6 後的 cache 寫入。

**Primary recommendation:** 以 gsd-statusline.js 為模板直接 port，保留 stdin timeout guard 和 silent fail 模式；settings.json merge 使用 `encoding/json` 的 `map[string]interface{}` read-modify-write 模式（與 cmd/update.go 的既有 JSON 處理一致）。Hook 安裝使用 Go `embed` package 將 JS 檔案嵌入 binary（專案首次使用，但 stdlib 直接支援）。

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Node.js stdlib (`fs`, `path`, `os`) | Node.js (任意版本) | Statusline hook 的所有 I/O | 零依賴；hook 執行環境由 Claude Code 保證有 Node.js |
| Go stdlib (`encoding/json`, `os`, `path/filepath`, `embed`) | Go 1.23+ | settings.json merge、hook 安裝、embed JS 檔案 | 專案已建立的 stdlib-first 模式；`embed` 是 Go 1.16+ stdlib |
| `github.com/spf13/viper` v1 | 已在 go.mod | `statusline_enabled` 在 mysd.yaml 的讀寫 | 專案中所有 mysd.yaml 讀寫都透過 viper（參考 `cmd/model.go` 的 read-modify-write 模式） |
| `github.com/spf13/cobra` v1.10.2 | 已在 go.mod | 新增 `statusline` 子指令 | 專案 CLI 框架 |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `gopkg.in/yaml.v3` | 已在 go.mod | mysd.yaml 讀寫（viper 底層） | 不需直接引入；viper 已處理 |
| ANSI escape codes | N/A（hardcoded） | terminal 顏色輸出 | 直接在 statusline hook 中 hardcode，與 gsd-statusline.js 相同模式 |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| `encoding/json` map merge for settings.json | viper for settings.json | settings.json 不是 mysd.yaml 格式，是 Claude Code 原生 JSON；用 encoding/json 更直接，避免 viper 的 YAML 假設 |
| regex line scan for state.yaml | `gopkg.in/yaml.v3` | D-10 明確指定 regex/line scan；hook 是純 Node.js stdlib，不能引入 Go YAML |
| Go embed package for hook install | 從 plugin/ 目錄 runtime copy | plugin/ 是 distribution 目錄，不保證在用戶的 PATH 旁邊；embed 是唯一可靠的 binary-carries-assets 方案 |

---

## Architecture Patterns

### Files to Create / Modify

```
plugin/
└── hooks/
    └── mysd-statusline.js      (NEW — statusline hook，distribution source)

plugin/commands/
└── mysd-statusline.md          (NEW — /mysd:statusline SKILL.md)

cmd/
├── init_cmd.go                 (MODIFY — 加入 hook embed 寫出 + settings.json statusLine 寫入)
├── statusline.go               (NEW — mysd statusline on/off/toggle 子指令)
└── archive.go                  (MODIFY — 加入 discuss-research-cache.json 刪除)

internal/config/
└── defaults.go                 (MODIFY — ProjectConfig 新增 StatuslineEnabled *bool)

internal/hooks/
└── embed.go                    (NEW — //go:embed directive 持有 mysd-statusline.js bytes)

plugin/commands/
└── mysd-discuss.md             (MODIFY — 加入 D-14~D-17 cache 邏輯)

.gitignore                      (MODIFY — 加入 discuss-research-cache.json)

.claude/hooks/
└── mysd-statusline.js          (由 mysd init 安裝時從 embed 寫出，非直接修改)
```

### Pattern 1: Node.js Statusline Hook (gsd-statusline.js port)

**What:** 讀取 stdin JSON → 解析 context_window → 建立色條 → 輸出 statusline 字串

**完整流程：**
1. `setTimeout(() => process.exit(0), 3000)` — stdin timeout guard（Claude's Discretion）
2. 讀取並 JSON.parse stdin
3. 提取 `data.model?.display_name`、`data.model?.id`、`data.session_id`、`data.context_window?.remaining_percentage`、`data.workspace?.current_dir`
4. 讀取 `{workspace}/.claude/mysd.yaml` 的 `statusline_enabled`（regex line scan）
5. 讀取 `{workspace}/.specs/state.yaml` 的 `change_name`（regex line scan，D-10）
6. 偵測 GSD 並存（D-04）
7. 計算 `used` percentage（D-03 公式）
8. 有 GSD 並存 + 有 `session` → 寫入 bridge file（D-04）
9. 建立 10 段色條 + 顏色（D-05）
10. 組合 statusline 字串輸出（D-01 格式）

**Key differences vs gsd-statusline.js:**
- 無 `todos` task 讀取邏輯（不需要）
- 無 `gsd-update-check` cache 邏輯
- 增加 `change_name` 讀取（step 5）
- 增加 `statusline_enabled` 讀取（step 4）；disabled → 空字串輸出，但仍執行 step 6~8
- Bridge file 條件寫入（只在 GSD 並存時，D-04）
- 使用 🥵 而非 💀（D-05）

### Pattern 2: settings.json Merge (encoding/json map)

**What:** 讀取現有 `.claude/settings.json` → 解析為 `map[string]interface{}` → 設定 `statusLine` key → 序列化回寫

**Why encoding/json over viper:** settings.json 是 Claude Code 的原生格式，不是 mysd.yaml 格式。直接用 `encoding/json` 避免 viper 的 YAML/config-file 假設。

```go
// Source: 基於 cmd/update.go 的 encoding/json 使用模式
func writeSettingsStatusLine(claudeDir string) error {
    settingsPath := filepath.Join(claudeDir, "settings.json")

    // Read existing (or start fresh)
    raw := map[string]interface{}{}
    if data, err := os.ReadFile(settingsPath); err == nil {
        _ = json.Unmarshal(data, &raw) // silent fail on parse error
    }

    // Merge: only set statusLine key, preserve all other keys
    raw["statusLine"] = map[string]interface{}{
        "type":    "command",
        "command": "node .claude/hooks/mysd-statusline.js",
    }

    out, err := json.MarshalIndent(raw, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(settingsPath, out, 0644)
}
```

### Pattern 3: Go embed for hook installation

**What:** 使用 Go `//go:embed` directive 將 `plugin/hooks/mysd-statusline.js` 嵌入 binary，`mysd init` 時寫出到 `.claude/hooks/`

**Why embed is the only viable option:** `plugin/` 是 distribution 目錄（在 GitHub Release 中），不保證在用戶機器的 binary 旁邊。embed 讓 binary 自帶 hook 內容，無需外部依賴。

```go
// Source: Go stdlib embed package (Go 1.16+)
// 建議新建 internal/hooks/embed.go 或 cmd/hooks_embed.go

package hooks // or cmd

import _ "embed"

//go:embed ../../plugin/hooks/mysd-statusline.js
var StatuslineHookContent []byte
```

然後在 init_cmd.go 中：
```go
import "github.com/xenciscbc/mysd/internal/hooks"

// 在 runInit 中
hookDest := filepath.Join(".", ".claude", "hooks", "mysd-statusline.js")
if err := os.MkdirAll(filepath.Dir(hookDest), 0755); err != nil {
    return err
}
if err := os.WriteFile(hookDest, hooks.StatuslineHookContent, 0644); err != nil {
    return err
}
```

**embed 路徑注意事項：**
- `//go:embed` 路徑相對於包含該 directive 的 `.go` 檔案
- 若 embed.go 在 `internal/hooks/`，路徑為 `../../plugin/hooks/mysd-statusline.js`
- 若 embed.go 在 `cmd/`，路徑為 `../plugin/hooks/mysd-statusline.js`
- embed 路徑不能使用 `..` 跨 module root，但可以在同一 module 內跨目錄

### Pattern 4: viper read-modify-write for mysd.yaml (statusline_enabled)

**What:** 用 viper 讀取 `.claude/mysd.yaml`，修改 `statusline_enabled`，寫回

**Source:** `cmd/model.go` 的 `runModelSet` 函式——現有 established pattern：

```go
// Source: cmd/model.go runModelSet（直接沿用）
v := viper.New()
v.SetConfigFile(configPath)
if err := v.ReadInConfig(); err != nil {
    // File may not exist yet — SafeWriteConfig will create it
}
v.Set("statusline_enabled", newValue) // true or false
if err := v.WriteConfig(); err != nil {
    if err2 := v.SafeWriteConfig(); err2 != nil {
        return fmt.Errorf("write config: %w", err2)
    }
}
```

### Pattern 5: ProjectConfig Extension (additive-only)

**What:** 在 `internal/config/defaults.go` 的 `ProjectConfig` struct 末尾新增 `StatuslineEnabled` 欄位

**Rule:** Phase 05-01 確立的原則：新欄位必須加在 struct 末尾（D-11/D-12）以保持 YAML field order 穩定。

```go
// 在現有 DocsToUpdate 欄位之後新增（struct 末尾）
StatuslineEnabled *bool `yaml:"statusline_enabled,omitempty" mapstructure:"statusline_enabled"`
```

使用 `*bool`（指標）而非 `bool`，可以區分「欄位不存在」（nil）與「明確設為 false」。D-12 要求：欄位不存在時視為 `true`（正常顯示）。

在 hook 端（Node.js）用 string 比對來處理 nil 情況：regex 找不到 `statusline_enabled` 欄位 → 視為 true。

### Pattern 6: regex line scan for YAML (Node.js hook)

**What:** 不引入 YAML library，用 regex 掃描 state.yaml 找 `change_name` 值

```javascript
// Source: 基於 D-10 設計決策
function readChangeName(workspaceDir) {
  try {
    const stateYamlPath = path.join(workspaceDir, '.specs', 'state.yaml');
    const content = fs.readFileSync(stateYamlPath, 'utf8');
    const match = content.match(/^change_name:\s*["']?([^"'\n\r]+)["']?\s*$/m);
    return match ? match[1].trim() : null;
  } catch (e) {
    return null; // silent fail
  }
}
```

### Pattern 7: discuss-research-cache.json operations (SKILL.md layer)

**What:** 所有 cache 的讀寫邏輯在 SKILL.md 層（`mysd-discuss.md`），不在 Go binary

**Why:** cache 格式是 free-form JSON，只有 SKILL.md orchestrator 知道 research 輸出的完整內容

**mysd-discuss.md 需要新增的兩個步驟：**

Step 4.5（在現有 Step 4 Deferred Notes 之後，Step 5 Research 之前）：
```
## Step 4.5: Research Cache Detection (D-17)

Check if cache exists:
Run: bash cat .specs/changes/{change_name}/discuss-research-cache.json 2>/dev/null

If file exists and is valid JSON:
- Extract cached_at field
- Ask user:
  "Found cached research from {cached_at}. Would you like to:
   1. Reuse cached research (skip Step 6)
   2. Run fresh research (overwrite cache after Step 6)
   3. Skip research entirely (cache unchanged)"
- If option 1: load cache, skip Step 6, go to Step 7 with cached research
- If option 2: proceed normally; cache will be overwritten after Step 6
- If option 3: go directly to Step 9 (no research)

If auto_mode: treat as option 2 (always run fresh research).
```

After Step 6 (parallel research complete), add cache write:
```
## Step 6.5: Write Research Cache (D-16)

After collecting all 4 research outputs, write cache:
Use Write tool to create .specs/changes/{change_name}/discuss-research-cache.json
Content: JSON with change_name, cached_at (ISO8601 timestamp), research object
```

**archive 時刪除：**
```go
// cmd/archive.go runArchive 函式中，在 moveDir 之前
cachePath := filepath.Join(changeDir, "discuss-research-cache.json")
_ = os.Remove(cachePath) // D-18: best-effort, ignore error
```

### Anti-Patterns to Avoid

- **Hook 中引入 npm 依賴：** statusline hook 必須 zero-dependency（只用 Node.js stdlib），否則在沒有 `node_modules` 的用戶機器上失敗
- **Stdin blocking without timeout：** 沒有 3s timeout guard 的 hook 在 Windows/Git Bash pipe 問題時會永久 hang（GSD #775）
- **settings.json 用 JSON.stringify 覆蓋：** 必須先讀取再 merge，保留既有 hooks 設定
- **statusline 輸出 stderr：** 任何錯誤都必須 silent fail；輸出 stderr 會破壞 Claude Code statusline 顯示
- **bool 而非 *bool：** `ProjectConfig.StatuslineEnabled` 用 `bool` 無法區分「未設定」（nil/omitempty）vs 「明確 false」，違反 D-12 要求
- **embed 路徑用 `..` 跨 module root：** Go embed 不允許跨越 module root 的 `..` 路徑；embed.go 檔案必須放在 module root 內

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| settings.json merge | 自訂 JSON merge utility | `encoding/json` + `map[string]interface{}` | 標準 Go stdlib 模式，專案已有 `cmd/update.go` 使用範例 |
| mysd.yaml 讀寫 | 直接用 `os.WriteFile` | viper read-modify-write | 保持與 `cmd/model.go` 一致，避免覆蓋其他欄位 |
| YAML 解析（Node.js hook） | 引入 js-yaml 或 yaml npm | regex line scan | Hook 必須 zero-dependency；`change_name` 是簡單字串欄位，regex 足夠 |
| context 計算邏輯 | 重新設計計算方式 | 直接 port 自 `gsd-statusline.js` 的公式 | D-03 已明確指定公式，與 GSD 保持相容 |
| hook 檔案安裝機制 | runtime 從 plugin/ 目錄 copy | Go `embed` package | plugin/ 不保證在用戶機器上可找到；embed 是唯一可靠的 binary-carries-assets 方案 |

**Key insight:** 這個 phase 是一個 port + extension，不是原創設計。大部分邏輯直接來自 `gsd-statusline.js`，Go 層只需小幅擴展既有 patterns。

---

## Common Pitfalls

### Pitfall 1: settings.json 的 statusLine 覆蓋 GSD 設定

**What goes wrong:** 用戶的全域 `~/.claude/settings.json` 可能已有 GSD 的 `statusLine`。`mysd init` 寫入專案層 `.claude/settings.json` 的 `statusLine` 會在 project scope 覆蓋 global scope，導致 GSD statusline 在 mysd 專案中消失。

**Why it happens:** Claude Code settings.json 有 global/project scope 層疊，project-level 覆蓋 global-level。

**How to avoid:** 這是 D-06 的設計意圖——mysd statusline **取代** GSD statusline 在 mysd 專案中的顯示。bridge file 機制（D-04）確保 GSD context-monitor 仍然能讀到 context 資料。這是預期行為，**不是 bug**。在 `mysd init` 的輸出訊息中說明清楚即可。

**Warning signs:** 用戶在 mysd 專案中看不到 GSD statusline 時的困惑。

### Pitfall 2: bridge file 寫入的 session_id 為空

**What goes wrong:** `data.session_id` 為 undefined 或空字串，寫入的 bridge file 路徑為 `/tmp/claude-ctx-.json`，與 gsd-context-monitor 讀取的路徑不匹配。

**Why it happens:** 某些 Claude Code 環境或版本可能不提供 session_id。

**How to avoid:** 寫入 bridge file 前先 guard check：`if (session && gsdCoexists)` 才寫入（gsd-statusline.js 已有此模式）。

### Pitfall 3: state.yaml 路徑 vs workspace.current_dir

**What goes wrong:** hook 用 `process.cwd()` 而非 `data.workspace?.current_dir` 來組合 state.yaml 路徑，在某些 Claude Code 設定下 cwd 可能不是 workspace root。

**Why it happens:** Claude Code hook 的 cwd 可能是 Claude Code binary 所在目錄，而不是用戶開啟的 workspace。

**How to avoid:** 始終優先用 `data.workspace?.current_dir`，fallback 才用 `process.cwd()`（gsd-statusline.js 的相同模式）。

### Pitfall 4: ProjectConfig StatuslineEnabled 的零值問題

**What goes wrong:** 如果用 `bool` 型別，Go 的零值是 `false`，舊 mysd.yaml 沒有此欄位的使用者一旦 viper 讀取，`statusline_enabled` 會被 resolved 為 `false`，導致 statusline 不顯示。

**Why it happens:** viper 對沒有在 config file 的 key 返回型別零值。

**How to avoid:** 使用 `*bool`（pointer to bool）搭配 `omitempty`。nil 表示「欄位不存在，視為 true」；`false` 表示「明確停用」。在 hook 端（Node.js），讀不到欄位時 default 為 enabled。

### Pitfall 5: discuss-research-cache.json 格式中的研究內容轉義

**What goes wrong:** 如果 research 輸出包含 Markdown 程式碼區塊、backticks、換行符，SKILL.md 層用字串插值組裝 JSON 時可能產生非法 JSON。

**Why it happens:** SKILL.md 層沒有 JSON serializer；AI agent 手動組裝 JSON 字串。

**How to avoid:** SKILL.md 中指示 AI 用 Write tool 寫入 JSON 時，說明 research 內容應該是 escaped 字串（AI 在 Write 工具中自然處理 JSON 轉義）。讀取時用 try/catch，parse 失敗 = silent fail，視同 cache 不存在（重新 research）。

### Pitfall 6: mysd init 的冪等性 — settings.json 修改

**What goes wrong:** 重複執行 `mysd init` 時，第二次呼叫把 settings.json 中其他 hooks（如 GSD 的 PostToolUse、PreToolUse）清掉。

**Why it happens:** 沒有正確 merge，直接覆蓋整個 settings.json。

**How to avoid:** merge 操作：先讀取現有 settings.json，只修改 `statusLine` key，保留其他所有 keys（Pattern 2 已說明正確做法）。

### Pitfall 7: Go embed 的路徑限制

**What goes wrong:** `//go:embed` directive 的路徑必須相對於包含該 directive 的 `.go` 檔案所在目錄，且不能包含跨越 module root 的 `..`。

**Why it happens:** embed 有意設計成禁止跨越 module 邊界。

**How to avoid:** embed.go 檔案必須放在相對路徑能正確指向 `plugin/hooks/` 的目錄。建議放在 `cmd/hooks_embed.go` 路徑，embed path 為 `../plugin/hooks/mysd-statusline.js`（從 `cmd/` 到 `plugin/hooks/` 是向上一層再向下，合法）。

---

## Code Examples

Verified patterns from existing codebase:

### Model 簡稱提取邏輯（D-02）
```javascript
// Source: 基於 D-02 決策，新撰寫
function extractModelShortname(data) {
  const name = (data.model?.display_name || data.model?.id || '').toLowerCase();
  if (name.includes('opus'))   return 'opus';
  if (name.includes('sonnet')) return 'sonnet';
  if (name.includes('haiku'))  return 'haiku';
  // fallback: first word of display_name
  const displayName = data.model?.display_name || '';
  const firstWord = displayName.split(' ')[0];
  return firstWord || 'claude';
}
```

### GSD 並存偵測（D-04）
```javascript
// Source: 基於 D-04 決策
function detectGsdCoexistence(workspaceDir) {
  const homedir = require('os').homedir();
  const claudeConfigDir = process.env.CLAUDE_CONFIG_DIR || path.join(homedir, '.claude');
  const checkPaths = [
    path.join(workspaceDir, '.claude', 'hooks', 'gsd-context-monitor.js'),
    path.join(claudeConfigDir, 'hooks', 'gsd-context-monitor.js')
  ];
  return checkPaths.some(p => {
    try { return fs.existsSync(p); } catch(e) { return false; }
  });
}
```

### statusline_enabled 讀取（hook 端，D-12）
```javascript
// Source: 基於 D-10/D-12 決策，regex line scan
function readStatuslineEnabled(workspaceDir) {
  try {
    const yamlPath = path.join(workspaceDir, '.claude', 'mysd.yaml');
    const content = fs.readFileSync(yamlPath, 'utf8');
    const match = content.match(/^statusline_enabled:\s*(.+)$/m);
    if (match) {
      const val = match[1].trim().toLowerCase();
      return val !== 'false';
    }
    return true; // not found = enabled (D-12)
  } catch(e) {
    return true; // file not found = enabled
  }
}
```

### archive.go 刪除 cache（D-18）
```go
// Source: 基於 D-18 決策，加入 runArchive 函式中，在 moveDir 之前
cachePath := filepath.Join(changeDir, "discuss-research-cache.json")
_ = os.Remove(cachePath) // best-effort, D-18: ignore error
```

### mysd statusline subcommand pattern（D-13）
```go
// Source: 基於 cmd/model.go runModelSet 模式
var statuslineCmd = &cobra.Command{
    Use:   "statusline [on|off]",
    Short: "Toggle or set statusline display",
    RunE:  runStatusline,
}
// 無參數 = toggle（讀當前值，寫入相反值）
// "on" = set true，"off" = set false
// 輸出純文字：Statusline: on / Statusline: off
```

### Go embed for hook（new pattern for this project）
```go
// Source: Go stdlib embed package
// File: cmd/hooks_embed.go (or internal/hooks/embed.go)

package cmd

import _ "embed"

//go:embed ../plugin/hooks/mysd-statusline.js
var statuslineHookBytes []byte
```

---

## Environment Availability

Step 2.6: SKIPPED (no external dependencies beyond project's existing packages and Node.js runtime which Claude Code guarantees)

---

## Validation Architecture

nyquist_validation 為 true（`.planning/config.json`），需要包含此 section。

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go testing (stdlib) + testify v1 |
| Config file | 無獨立 config，使用 `go test ./...` |
| Quick run command | `go test ./cmd/... -run TestStatusline -v` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Behavior | Test Type | Automated Command | File Exists? |
|----------|-----------|-------------------|-------------|
| `mysd statusline` toggle/on/off 正確寫入 mysd.yaml | unit | `go test ./cmd/... -run TestRunStatusline` | No — Wave 0 |
| `mysd init` 安裝 hook 檔案到 `.claude/hooks/` | unit | `go test ./cmd/... -run TestInitStatuslineInstall` | No — Wave 0 |
| `mysd init` 寫入 settings.json `statusLine` key 且保留既有 keys | unit | `go test ./cmd/... -run TestWriteSettingsStatusLine` | No — Wave 0 |
| `mysd archive` 刪除 discuss-research-cache.json | unit | `go test ./cmd/... -run TestArchiveDeletesResearchCache` | No — Wave 0 |
| `ProjectConfig` 新欄位 `statusline_enabled` 零值相容 | unit | `go test ./internal/config/...` | No — Wave 0 |
| mysd-statusline.js：正確輸出 statusline 格式 | manual smoke | N/A | N/A |
| mysd-statusline.js：disabled 時不輸出 | manual smoke | N/A | N/A |

### Sampling Rate
- **Per task commit:** `go test ./cmd/... -run TestStatusline -v`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `cmd/statusline_test.go` — 新檔案，covers `TestRunStatusline` (on/off/toggle)
- [ ] `cmd/init_cmd_test.go` 擴展 — 新增 `TestInitStatuslineInstall`、`TestWriteSettingsStatusLine`
- [ ] `cmd/archive_test.go` 擴展 — 新增 `TestArchiveDeletesResearchCache`
- [ ] `internal/config/config_test.go` 或同層測試 — 新增 `StatuslineEnabled` 零值相容測試

---

## Integration Points: 現有系統的影響範圍

### 1. mysd init 的變更

現有 `cmd/init_cmd.go` 的 `runInit` 函式需要新增：
1. 從 embed 寫出 hook 檔案到 `.claude/hooks/mysd-statusline.js`
2. 呼叫 `writeSettingsStatusLine(".claude")` 寫入 settings.json

**確認（已研究）：** 專案目前**沒有任何 `//go:embed` 用法**（codebase 全域 grep 確認）。`plugin/` 是純 distribution 目錄，在 `mysd update` 時由 GitHub Release tarball 解壓縮後 copy 到 `.claude/`。

**解法（唯一 viable option）：**
使用 Go `//go:embed` 將 hook 內容嵌入 binary。建議新建 `cmd/hooks_embed.go` 持有 embed directive，路徑為 `../plugin/hooks/mysd-statusline.js`（從 cmd/ 到 plugin/ 的相對路徑，合法）。這是專案第一個 embed 用法，但 Go stdlib `embed` package（1.16+）完全支援，無需新依賴。

### 2. .claude/settings.json 的 scope

現有的 `.claude/settings.json`（project level）已有 GSD 的 hooks 設定（PostToolUse context-monitor 等）。`mysd init` 的 merge 操作只修改 `statusLine` key，不影響其他 hooks。

執行後此 repo 的 `.claude/settings.json` 將新增 `"statusLine": { "type": "command", "command": "node .claude/hooks/mysd-statusline.js" }`，覆蓋掉原本 GSD 的 statusLine（若有）——這是預期行為。

### 3. archive 的非破壞性擴展

`runArchive` 函式目前的執行順序：
1. Gate 1: phase check
2. Gate 2: MUST items check
3. `saveArchivedState` (snapshot)
4. `moveDir(changeDir, archiveDir)` — 此後 changeDir 不再存在
5. State transition + save

cache 刪除**必須在 step 4 (moveDir) 之前**插入：
```go
// 在 saveArchivedState 之後、moveDir 之前
cachePath := filepath.Join(changeDir, "discuss-research-cache.json")
_ = os.Remove(cachePath) // D-18: best-effort
```

---

## State of the Art

| Old Approach | Current Approach | Notes |
|--------------|------------------|-------|
| GSD statusline（全域，顯示 model+task+dir+context） | mysd statusline（project-level，顯示 model+change+dir+context） | mysd 版本更接近 spec workflow，用 change 替換 GSD 的 todo task |
| Bridge file 無條件寫入（GSD statusline） | Bridge file 條件寫入（只在 GSD 並存時） | 避免在純 mysd 環境產生無人消費的 /tmp 檔案 |
| 無 research cache（discuss 每次重新 research） | discuss-research-cache.json proactive 寫入 | 重啟 session 後可重用 research 結果，節省時間 |

---

## Open Questions

1. **[RESOLVED] Plugin/hooks embed 機制**
   - What we found: 專案目前沒有任何 `//go:embed` 用法（grep 確認）
   - Resolution: 需要引入 Go `embed` package 來嵌入 `plugin/hooks/mysd-statusline.js`；這是專案第一個 embed 用法，但技術上直接（stdlib，無新依賴）
   - Planner action: Plan 中需包含建立 `cmd/hooks_embed.go` 持有 `//go:embed` directive 的步驟

2. **mysd-discuss.md 的 cache 寫入如何處理 research 內容的特殊字元？**
   - What we know: SKILL.md 層使用 Write tool，AI agent 需要組裝 JSON 字串
   - What's unclear: 如果 research 輸出包含 backticks 或 `"` 引號，JSON 可能非法
   - Recommendation: 在 SKILL.md 中指示 AI 用 Write tool 寫入時，research 內容的值以 escaped JSON string 形式表達；讀取時 try/catch，parse 失敗 = silent fail，視同 cache 不存在

3. **`/mysd:statusline` SKILL.md 需要確認的 binary 輸出格式**
   - What we know: D-13 說「對應 `mysd statusline` binary 子指令」，需要新增 Go cobra subcommand
   - What's unclear: SKILL.md 是否需要解析 binary 輸出做額外格式化
   - Recommendation: SKILL.md 為薄 wrapper（thin wrapper pattern，與 mysd-docs.md、mysd-note.md 一致），只呼叫 `mysd statusline [on|off]` 並直接顯示輸出；binary 輸出純文字 "Statusline: on" 或 "Statusline: off"

---

## Sources

### Primary (HIGH confidence)

- Direct code read: `.claude/hooks/gsd-statusline.js` — 完整 GSD statusline 實作，mysd 版本直接 port；bridge file schema 和 context 計算公式確認
- Direct code read: `.claude/hooks/gsd-context-monitor.js` — bridge file 消費者，schema 確認（`{ session_id, remaining_percentage, used_pct, timestamp }`）
- Direct code read: `cmd/init_cmd.go` — init 命令現有框架，擴展點確認
- Direct code read: `internal/config/defaults.go` — ProjectConfig struct，新欄位插入位置確認（DocsToUpdate 之後）
- Direct code read: `cmd/model.go` — viper read-modify-write 模式（runModelSet）
- Direct code read: `cmd/archive.go` — runArchive 函式，cache 刪除插入點確認（moveDir 之前）
- Direct code read: `.claude/settings.json` — 現有 statusLine 格式和 hooks 結構確認
- Direct code read: `plugin/hooks/hooks.json` — plugin hooks 目錄結構
- Direct code read: `plugin/commands/mysd-discuss.md` — discuss SKILL.md 現有步驟，Step 4.5 和 Step 6.5 插入點確認
- Codebase grep: `//go:embed` — 確認無現有 embed 用法，embed 是首次引入
- Direct code read: `internal/update/pluginsync.go` — 確認 plugin/ 是 distribution 目錄（非 embed source）
- 12-CONTEXT.md decisions D-01 ~ D-19 — 所有設計決策

### Secondary (MEDIUM confidence)

- 專案慣例（Phase 05-01 established pattern）：新欄位加在 struct 末尾；additive-only extension
- Phase 07 決策：`mysd model` 使用 plain text fmt.Fprintf（非 lipgloss），statusline subcommand 應沿用相同模式
- Phase 11 決策：thin wrapper pattern 用於 mysd-docs.md、mysd-note.md

---

## Metadata

**Confidence breakdown:**
- Statusline hook implementation: HIGH — gsd-statusline.js 可直接 port，差異點 D-01~D-05 明確
- Go init extension: HIGH — encoding/json merge pattern 清楚，init_cmd.go 框架已存在；embed 機制已確認（首次引入 Go embed package，但 stdlib 直接支援，無新依賴）
- ProjectConfig extension: HIGH — defaults.go 結構清楚，`*bool` 處理方式明確
- archive.go extension: HIGH — 插入點確認，os.Remove silent fail 模式簡單
- SKILL.md (mysd-discuss cache): HIGH — discuss SKILL.md 結構已讀，插入 Step 4.5 和 Step 6.5 位置清楚
- Pitfalls: HIGH — 基於直接閱讀 gsd-statusline.js 和現有 Go code

**Research date:** 2026-03-27
**Valid until:** 2026-04-27（Go stdlib 和 Node.js stdlib 穩定）
