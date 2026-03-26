# Phase 7: New Binary Commands & Scanner Refactor — Research

**Researched:** 2026-03-26
**Domain:** Go CLI subcommand design, language-agnostic file scanner, atomic config write, SKILL.md agent UX layer
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**D-01:** 語言偵測使用 file-based markers：`go.mod` → Go、`package.json` → Node.js、`requirements.txt`/`pyproject.toml` → Python。未匹配到已知 marker 時 primary_language 標為 `unknown`

**D-02:** ScanContext struct 完全替換為語言無關通用 struct：`primary_language`（string）、`files`（副檔名統計）、`modules`（偵測到的 module/package 列表）、`existing_specs`。移除 Go-specific `PackageInfo` 陣列

**D-03:** Binary 只收集 metadata，LLM agent 負責理解任意語言並生成 spec。Go binary 不做語言特定的 spec 決策

**D-04:** 保持 `--context-only` 執行模式。scan 輸出 JSON metadata，spec 寫入由 SKILL.md agent 負責。Binary 不新增 `--write-specs` flag

**D-05:** `mysd init` 內部直接展開為 `scan --scaffold-only` 執行（不顯示 deprecation warning，完全回展相容）

**D-06:** 首次建立 `openspec/config.yaml` 時的互動式 locale 詢問（FSCAN-04）在 SKILL.md agent 層發生：agent 詢問使用者後呼叫 `mysd lang set {locale}` 寫入。Go binary 的 scaffold-only 只建立空結構

**D-07:** SKILL-01 的 skills 推薦邏輯在 mysd-planner agent 層（LLM 根據 task 內容推斷），Go binary 不實作規則式 skills 對映

**D-08:** SKILL-02/03 的表格顯示與使用者確認流程在 SKILL.md 層：plan 完成後 SKILL.md 讀取 `--context-only` JSON，Claude 直接呈現 task↔skills 對應表並互動確認

**D-09:** 批次同意 UX（SKILL-03）：呈現完整對應表後詢問 `Accept all recommended? Y/n`，預設 accept（Enter 即同意）。使用者選 n 才進入逐一調整流程

**D-10:** `ffe` 模式（SKILL-04）跳過互動，直接使用 planner 推薦值

**D-11:** `mysd model`（讀）輸出 lipgloss table 格式：Profile 標題行 + Role｜Model 兩欄表格

**D-12:** `mysd model set <profile>` 在 Go binary 層直接寫入 `.claude/mysd.yaml` 的 `model_profile` 欄位

**D-13:** 新增 `FilterBlockedTasks(tasks []TaskItem, failedIDs []int) []TaskItem`（`internal/executor/waves.go`）：給定已失敗的 task ID set，返回需略過的下游 tasks（遞移性）

**D-14:** SKILL.md 執行每個 wave 前，將累積的 `failed_ids` 傳給 binary，binary 回傳可執行清單；被略過的 tasks 標記為 `skipped`

**Phase 5 D-08/09（繼承）：** locale 使用 BCP47 格式（zh-TW, en-US, ja-JP）；`openspec/config.yaml` 的 locale 為 source of truth；`/mysd:lang` 修改 locale 時兩個 config 原子同步更新

### Claude's Discretion

- ScanContext 的 `modules` 欄位具體結構（per-language module metadata 細節）
- `unknown` language 時 scan 回傳的 fallback metadata 格式
- `mysd model` table 的 lipgloss 樣式細節（顏色、對齊）
- `scan --scaffold-only` 建立的空結構具體目錄和文件列表
- `mysd model set` 的 profile 驗證邏輯（無效 profile 時的錯誤訊息）

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope

</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| FCMD-03 | `/mysd:model` 顯示目前 profile 及所有 agent role 的 resolved model | `config.ResolveModel` 已存在，`DefaultModelMap` 有 10 個 roles；新增 `cmd/model.go` subcommand，使用 lipgloss table 格式輸出 |
| FCMD-04 | `/mysd:lang` 互動式設定 response_language 和 document_language，同步 mysd.yaml 和 openspec/config.yaml | `spec.WriteOpenSpecConfig` + Viper `WriteConfig` 已有基礎；需實作 atomic write（write-then-rename 或 defer rollback 模式） |
| FCMD-05 | `/mysd:lang` 使用者可選擇或輸入語言，自動轉換為合法 BCP47 locale 值 | locale 轉換邏輯屬 SKILL.md agent 層；binary 只接受已標準化的 BCP47 值 |
| FSCAN-01 | `/mysd:scan` 升級為語言無關通用掃描器 | `internal/scanner/scanner.go` 整體重構；移除 `PackageInfo`，新增 `primary_language`、`files`（副檔名統計）、`modules` |
| FSCAN-02 | Scan 偵測語言、模組結構，產生 `openspec/config.yaml` + `openspec/specs/` 下的 spec 文件 | binary `--context-only` 輸出 JSON；spec 寫入由 SKILL.md agent 層負責（不在 binary 層）|
| FSCAN-03 | 已存在 `openspec/config.yaml` 時只增量更新 specs，不覆蓋 config | `spec.ReadOpenSpecConfig` 已有 zero-value 回傳模式；scan 邏輯讀取後判斷 `existing: true` |
| FSCAN-04 | 首次建立 config.yaml 時互動式詢問 locale | 在 SKILL.md agent 層處理互動；binary 的 scaffold-only 只建空結構 |
| FSCAN-05 | `/mysd:init` 改為 `scan --scaffold-only`，只建空結構 + 互動式設定 locale | `cmd/init_cmd.go` 的 `runInit` 改為呼叫 `runScan` with `--scaffold-only` flag；SKILL.md 接管 locale 互動 |
| SKILL-01 | Planner 自動依 task 內容推薦 `skills` 欄位 | mysd-planner agent 定義需擴展；`tasks.md` 的 `skills` 欄位已在 `TaskEntry` 中（Phase 5 FSCHEMA-04） |
| SKILL-02 | Plan 完成後列出所有 task 與推薦 skills 對應表，互動式讓使用者確認 | SKILL.md 消費 `plan --context-only` JSON（已有 `wave_groups`、`tasks` 欄位）；表格顯示在 SKILL.md 層 |
| SKILL-03 | 使用者可逐一調整或批次同意推薦的 skills | SKILL.md 呈現表格後詢問 `Accept all recommended? Y/n`；預設 accept |
| SKILL-04 | ffe 模式跳過互動，直接使用推薦值 | `auto_mode: true` 已在 `plan --context-only` JSON 的 `auto_mode` 欄位；SKILL.md 讀取此值決定是否顯示確認提示 |

</phase_requirements>

---

## Summary

Phase 7 有四個並行工作流：(1) Go binary 新增 `model` 和 `lang` subcommands，(2) `scanner` package 完整重構為語言無關通用掃描器，(3) `init` 命令橋接到 `scan --scaffold-only`，(4) SKILL.md 和 agent 層新增 skills 推薦 UX。

所有 Go binary 的工作都有清楚的前例可循：`cmd/worktree.go` 示範了多 subcommand 設計和 JSON stdout 輸出；`internal/output.Printer` 提供 lipgloss TTY 偵測；`internal/config.ResolveModel` 直接可呼叫；`spec.WriteOpenSpecConfig` + `spec.ReadOpenSpecConfig` 已有 Phase 5 基礎。

最大的設計挑戰是 **atomic config write**（`/mysd:lang` 需要同時更新兩個 YAML 文件）以及 **scanner struct 破壞性替換**（現有 `ScanContext` 有外部消費者 `cmd/scan.go` 和 `mysd-scan.md` SKILL.md，兩者都需要同步更新）。

**Primary recommendation:** 按工作流分組成 5 個 wave-able tasks：(1) FilterBlockedTasks，(2) Scanner 重構，(3) `mysd model` subcommand，(4) `mysd lang` subcommand + atomic write，(5) `mysd init` 橋接 + SKILL.md updates。SKILL.md 的 skills 推薦 UX 可與 binary tasks 同步進行。

---

## Standard Stack

### Core（已在 go.mod，版本已驗證）

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| github.com/spf13/cobra | v1.10.2 | CLI subcommand（新增 `model` 子命令） | 專案已用，`worktree.go` 是精確範本 |
| github.com/spf13/viper | v1.21.0 | Config 讀寫（`model set` 寫 `model_profile`） | 專案已用；`viper.WriteConfig()` 是標準寫入路徑 |
| github.com/charmbracelet/lipgloss | v1.1.0 | `mysd model` table 格式輸出 | 專案已用；`internal/output/colors.go` 已定義 style constants |
| github.com/charmbracelet/x/term | v0.2.1 | TTY 偵測（`output.NewPrinter` 已用） | 專案已有，無需再 import |
| gopkg.in/yaml.v3 | v3.0.1 | YAML marshal/unmarshal（lang atomic write） | 專案已用；`spec.WriteOpenSpecConfig` 範本 |

### 無需新增依賴

所有 Phase 7 實作均可使用現有 go.mod 依賴，不需要 `go get` 任何新 package。

---

## Architecture Patterns

### Recommended Project Structure（Phase 7 異動）

```
cmd/
├── model.go          # 新增 — mysd model + model set subcommands
├── lang.go           # 新增 — mysd lang + lang set subcommands
├── scan.go           # 修改 — 新增 --scaffold-only flag，更新 ScanContext 參照
├── init_cmd.go       # 修改 — runInit 改為呼叫 runScan with --scaffold-only
└── root.go           # 修改 — 新增 modelCmd, langCmd 到 rootCmd

internal/
├── scanner/
│   ├── scanner.go        # 完整重構 — 新 ScanContext struct
│   └── scanner_test.go   # 同步更新 — 所有舊 PackageInfo 測試替換
└── executor/
    ├── waves.go          # 新增 FilterBlockedTasks function
    └── waves_test.go     # 新增依賴失敗傳播測試

.claude/
├── commands/
│   ├── mysd-model.md     # 新增 SKILL.md
│   ├── mysd-lang.md      # 新增 SKILL.md
│   └── mysd-scan.md      # 修改 — 更新 JSON fields 參照新 ScanContext
└── agents/
    └── mysd-planner.md   # 修改 — 新增 skills 推薦邏輯
```

### Pattern 1: Cobra Subcommand with Sub-subcommands（`mysd model`）

`cmd/worktree.go` 是精確範本：一個 parent command（`worktreeCmd`）+ 兩個 subcommands（`worktreeCreateCmd`, `worktreeRemoveCmd`），全部 wire 到 `rootCmd`。

`mysd model` 沿用相同模式：

```go
// Source: cmd/worktree.go pattern
var modelCmd = &cobra.Command{
    Use:   "model",
    Short: "Display or set model profile",
}

var modelSetCmd = &cobra.Command{
    Use:  "set <profile>",
    Args: cobra.ExactArgs(1),
    RunE: runModelSet,
}

func init() {
    modelCmd.AddCommand(modelSetCmd)
    rootCmd.AddCommand(modelCmd)
}
```

`mysd model`（讀）不需要 subcommand——`modelCmd.RunE = runModelRead` 即可，呼叫 `config.ResolveModel` 取出所有 roles。

### Pattern 2: Lipgloss Table（`mysd model` 讀取輸出）

lipgloss v1.1.0 有 `lipgloss/table` 子套件。專案已有 `output.Printer` 的 TTY 偵測，table 渲染需在 TTY 模式下才啟用，non-TTY 輸出純文字（`|` 分隔格式）。

```go
// TTY 模式下使用 lipgloss table
// Non-TTY fallback: fmt.Fprintf(w, "Profile: %s\n", profile) + 逐行 role|model 輸出
```

`DefaultModelMap` 在 `internal/config/config.go` 已包含所有 10 個 roles。`runModelRead` 直接迭代 `DefaultModelMap[profile]` 即可，role 順序需固定（slice of known roles 而非 map iteration）。

### Pattern 3: Viper WriteConfig（`mysd model set`）

`config.Load` 使用 Viper 讀取 `.claude/mysd.yaml`。`model set` 寫入需要：

```go
// 1. Load current config
cfg, err := config.Load(".")

// 2. Validate profile is one of: quality, balanced, budget
// 3. Set value and write back via viper
v := viper.New()
v.SetConfigFile(filepath.Join(".", ".claude", "mysd.yaml"))
v.ReadInConfig()         // 讀取現有 config（保留其他欄位）
v.Set("model_profile", profile)
v.WriteConfig()          // 原地覆寫
```

注意：`viper.WriteConfig()` 需要先有 config file 存在，`viper.SafeWriteConfig()` 用於初次建立。Phase 7 的 `model set` 假設 `.claude/mysd.yaml` 已存在（由 `mysd init` 或 `mysd scan --scaffold-only` 建立）。

### Pattern 4: Atomic Config Write（`mysd lang set`）

需要同時更新 `.claude/mysd.yaml`（`response_language`）和 `openspec/config.yaml`（`locale`）。任一失敗都必須 rollback，確保兩者不會不一致。

**Write-then-rename 模式（推薦）：**

```go
// Phase 1: 寫入兩個 temp file
// Phase 2: os.Rename(tmpA, fileA) + os.Rename(tmpB, fileB)
// 如果任何步驟失敗，defer 清理 temp files
```

或使用 **defer rollback 模式**（讀取舊值 → 寫入新值 → 失敗時 defer 還原）。

`spec.WriteOpenSpecConfig` 直接 `os.WriteFile`（非原子），需要在 `cmd/lang.go` 層實作 atomic 包裝。

### Pattern 5: Language-Agnostic Scanner（新 ScanContext）

現有 `BuildScanContext` 只收集 `.go` 文件。新設計：

```go
// Source: CONTEXT.md D-01, D-02
type ScanContext struct {
    RootDir         string              `json:"root_dir"`
    PrimaryLanguage string              `json:"primary_language"` // "go" | "nodejs" | "python" | "unknown"
    Files           map[string]int      `json:"files"`            // 副檔名 → 計數，e.g. {".go": 42, ".md": 10}
    Modules         []ModuleInfo        `json:"modules"`          // 語言無關的 module/package 列表
    ExistingSpecs   []string            `json:"existing_specs"`
    ExcludedDirs    []string            `json:"excluded_dirs"`
    TotalFiles      int                 `json:"total_files"`
    ConfigExists    bool                `json:"config_exists"`    // openspec/config.yaml 是否已存在
}

type ModuleInfo struct {
    Name string `json:"name"`  // 模組/套件名稱（語言無關）
    Dir  string `json:"dir"`   // 相對路徑
    // Claude's Discretion: per-language details 可在此擴展
}
```

語言偵測（file-based markers，D-01）：

```go
func detectPrimaryLanguage(root string) string {
    markers := []struct{ file, lang string }{
        {"go.mod", "go"},
        {"package.json", "nodejs"},
        {"requirements.txt", "python"},
        {"pyproject.toml", "python"},
    }
    for _, m := range markers {
        if _, err := os.Stat(filepath.Join(root, m.file)); err == nil {
            return m.lang
        }
    }
    return "unknown"
}
```

### Pattern 6: FilterBlockedTasks（`internal/executor/waves.go`）

給定一組已失敗的 task IDs，使用現有的 dependency graph 做遞移性傳播：

```go
// D-13: 遞移性依賴失敗傳播
func FilterBlockedTasks(tasks []TaskItem, failedIDs []int) []TaskItem {
    // 1. 建立 failedSet
    // 2. BFS/DFS 從 failedSet 出發，順著 adj（依賴反向）找出所有下游
    // 3. 返回 tasks 中 NOT in blockedSet 的項目
}
```

這與 `BuildWaveGroups` 使用的 `adj` 圖結構相同（`id -> []id` that depend on it），可以複用建圖邏輯。

### Anti-Patterns to Avoid

- **直接 viper.Set + viper.WriteConfig 不讀取現有值：** 會清空 config 中的其他欄位（tdd、execution_mode 等）。必須先 `v.ReadInConfig()` 再 `v.Set(key, val)` 再 `v.WriteConfig()`
- **Scanner 保留 PackageInfo 做向後相容：** CONTEXT.md D-02 明確指出完全替換，不保留舊 struct。消費者（`cmd/scan.go`、`mysd-scan.md`）必須同步更新
- **在 binary 層做語言特定決策：** D-03 原則——binary 只是 metadata collector，LLM 理解語言並生成 spec
- **Lipgloss table 在 non-TTY 環境（CI）輸出 ANSI：** 必須繼承 `output.NewPrinter` 的 TTY 偵測模式；non-TTY 用 plain text

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Config 讀寫 | 自己 YAML marshal/unmarshal | `viper.ReadInConfig()` + `viper.WriteConfig()` | 保留現有欄位，不破壞未 touch 的 config keys |
| TTY 偵測 | `os.Stdin.Fd()` 直接呼叫 | `output.NewPrinter` 已封裝 `charmbracelet/x/term` | 專案已建立 non-TTY fallback 模式 |
| Model resolve | 自建 map lookup | `config.ResolveModel(role, profile, overrides)` | 已有完整的 overrides 邏輯 |
| OpenSpec config 讀寫 | 自己 YAML 處理 | `spec.WriteOpenSpecConfig` / `spec.ReadOpenSpecConfig` | Phase 5 已有完整實作，含 zero-value 慣例 |
| 依賴圖遍歷 | 自建 BFS | 複用 `waves.go` 內的 adj graph 建構邏輯 | `FilterBlockedTasks` 需要的遍歷和 `BuildWaveGroups` 共享相同 graph |

**Key insight:** Phase 7 的所有 Go binary 工作都是在現有基礎設施上疊加，沒有需要從頭設計的新型態問題。

---

## Common Pitfalls

### Pitfall 1: Viper WriteConfig 清空 Config
**What goes wrong:** `viper.Set("model_profile", val)` 後直接 `viper.WriteConfig()` ——如果 Viper 的 config instance 是剛建立的且未 `ReadInConfig()`，會只寫出 `model_profile` 一個欄位，清除其他欄位
**Why it happens:** Viper 的 in-memory state 初始為空；`WriteConfig()` 序列化當前 in-memory state
**How to avoid:** 永遠先 `ReadInConfig()`（或使用 `viper.MergeConfig`），確認 SetConfigFile 指向正確路徑再 WriteConfig
**Warning signs:** 寫入後 `mysd plan --context-only` 回傳的 execution_mode 變成空字串

### Pitfall 2: ScanContext 破壞性替換的消費者同步
**What goes wrong:** `scanner.ScanContext` 的 `Packages []PackageInfo` 被移除後，`cmd/scan.go`（直接引用）和 `mysd-scan.md` SKILL.md（引用 JSON 欄位名 `packages`、`go_files` 等）都會失效
**Why it happens:** 兩個消費者位於不同層（Go code vs Markdown instruction）
**How to avoid:** 必須在同一個 task 或同一個 plan 中同步更新：`scanner.go`、`scanner_test.go`、`cmd/scan.go`、`mysd-scan.md` 四個文件
**Warning signs:** `mysd scan --context-only` 成功但 SKILL.md agent 找不到 `packages` 欄位

### Pitfall 3: Atomic Write 的 Windows os.Rename 限制
**What goes wrong:** 在 Windows 上，`os.Rename(tmp, dst)` 若 `dst` 已存在且被另一個進程開啟，會回傳 error（不像 POSIX 原子替換）
**Why it happens:** Windows 檔案系統語意與 POSIX 不同
**How to avoid:** 使用 defer rollback 模式（讀取舊值 → 寫入新值 → 失敗時 defer 還原）而非 write-then-rename。若兩個文件寫入之間 process 被中斷，最多只有一個文件被更新，下次執行 `mysd lang` 可修復
**Warning signs:** Windows CI 的 lang set test 失敗

### Pitfall 4: Lipgloss Table 在 Non-TTY 輸出 ANSI
**What goes wrong:** `lipgloss` 預設會自動偵測 TTY，但 `cmd.OutOrStdout()` 在 testing 時不是真實 TTY，導致 table render 在測試中輸出 ANSI escape codes
**Why it happens:** lipgloss 的 color profile 偵測依賴 os.Stdout，不是 cmd.OutOrStdout()
**How to avoid:** 使用 `output.NewPrinter` 的 TTY 偵測；table 渲染包裝在 `if p.isTTY` 分支；non-TTY 輸出 plain text
**Warning signs:** cmd 層的 unit test snapshot 包含 `\x1b[` escape codes

### Pitfall 5: mysd model set 無效 profile 的 UX
**What goes wrong:** 使用者輸入 `mysd model set fast`（不存在的 profile），binary 靜默寫入，`ResolveModel` fallback 到 `claude-sonnet-4-5`，使用者不知道 profile 名稱無效
**Why it happens:** `DefaultModelMap` 只有 quality/balanced/budget 三個 key
**How to avoid:** 在寫入前驗證 profile 存在於 `DefaultModelMap`，回傳清楚的錯誤訊息："unknown profile %q; valid profiles: quality, balanced, budget"
**Warning signs:** `mysd model` 顯示 Profile: fast，但 ResolveModel 全部回傳 sonnet-4-5 fallback

### Pitfall 6: FilterBlockedTasks 遞移性傳播方向
**What goes wrong:** `FilterBlockedTasks` 只過濾直接依賴失敗 task 的 tasks，沒有做遞移性傳播（T1 失敗 → T2 depends T1 被過濾，但 T3 depends T2 未被過濾）
**Why it happens:** BFS 只走一層，忘記繼續傳播 blocked set
**How to avoid:** BFS 使用 queue：把初始 failedIDs 加入 blocked set 後，從每個 blocked task 出發繼續遍歷 adj（downstream），全部加入 blocked set 直到 queue 空
**Warning signs:** `waves_test.go` 的 transitive propagation test 失敗

---

## Code Examples

### mysd model 讀取輸出（TTY mode）

```
Profile: balanced

Role              Model
──────────────────────────────────────
spec-writer       claude-sonnet-4-5
designer          claude-sonnet-4-5
planner           claude-sonnet-4-5
executor          claude-sonnet-4-5
verifier          claude-sonnet-4-5
fast-forward      claude-sonnet-4-5
researcher        claude-sonnet-4-5
advisor           claude-sonnet-4-5
proposal-writer   claude-sonnet-4-5
plan-checker      claude-sonnet-4-5
```

Non-TTY（pipes/CI）：`fmt.Fprintf` 純文字，每行 `{role}\t{model}`。

### mysd lang set 互動設計（SKILL.md 層）

```
Current language settings:
  mysd.yaml response_language: zh-TW
  openspec/config.yaml locale: zh-TW

Select language:
  1. zh-TW (Traditional Chinese)
  2. en-US (English)
  3. ja-JP (Japanese)
  4. Enter custom BCP47 code

> 2

Updating language settings... done
  response_language: en-US
  locale: en-US
```

### plan --context-only JSON（skills 欄位已存在）

tasks 的 `skills` 欄位在 `executor.TaskItem` 已有（`Skills []string`）。planner agent 推薦後，plan `--context-only` JSON 的 `tasks` 陣列中每個 task 的 `skills` 欄位將被填充。SKILL.md 讀取此 JSON 並展示 task↔skills 對應表。

### scan --context-only 新 JSON 格式

```json
{
  "root_dir": "/path/to/project",
  "primary_language": "nodejs",
  "files": {
    ".ts": 42,
    ".json": 8,
    ".md": 12
  },
  "modules": [
    {"name": "src/auth", "dir": "src/auth"},
    {"name": "src/api", "dir": "src/api"}
  ],
  "existing_specs": ["src/auth"],
  "excluded_dirs": ["node_modules"],
  "total_files": 62,
  "config_exists": false
}
```

---

## Runtime State Inventory

> 此 phase 為 Go binary 新增 subcommands 和修改現有 scanner，屬於程式碼層修改，無 rename/refactor 語意。

| Category | Items Found | Action Required |
|----------|-------------|------------------|
| Stored data | 無 — 此 phase 不修改任何 DB 或 key-value store | None |
| Live service config | 無 — mysd 不依賴外部 service | None |
| OS-registered state | 無 | None |
| Secrets/env vars | 無 — model_profile 不是 secret | None |
| Build artifacts | `mysd.exe`（根目錄）為舊 binary，每次 `go build` 重新生成 | 執行 `go build` 後自動更新 |

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | Binary build | 確認（go.mod: 1.25.5） | 1.25.5 | — |
| git | worktree operations | 確認（Phase 6 前提） | 系統 git | — |
| All go.mod deps | Build | 確認（go test ./... 全通過） | 見 go.mod | — |

**Missing dependencies with no fallback:** None

**Step 2.6 assessment:** Phase 7 為純 Go code 修改，無需額外安裝任何外部工具。現有測試環境（`go test ./... OK`）確認 build toolchain 完整。

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | testify v1.11.1（`github.com/stretchr/testify`）|
| Config file | none — 使用 `go test ./...` |
| Quick run command | `go test ./internal/executor/... ./internal/scanner/... ./cmd/...` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FCMD-03 | `mysd model` 輸出 Profile 標題 + Role/Model 表格 | unit | `go test ./cmd/... -run TestModelCmd` | ❌ Wave 0 |
| FCMD-04 | `mysd lang set zh-TW` 同步更新兩個 config 文件 | unit | `go test ./cmd/... -run TestLangSet` | ❌ Wave 0 |
| FCMD-04 | 一個寫入失敗時兩個 config 保持原狀（atomic rollback） | unit | `go test ./cmd/... -run TestLangSetAtomic` | ❌ Wave 0 |
| FSCAN-01 | Go 專案偵測 `primary_language: "go"` | unit | `go test ./internal/scanner/... -run TestBuildScanContext_GoProject` | ❌ Wave 0（替換） |
| FSCAN-01 | Node.js 專案偵測 `primary_language: "nodejs"` | unit | `go test ./internal/scanner/... -run TestBuildScanContext_NodeProject` | ❌ Wave 0（新增） |
| FSCAN-01 | Python 專案偵測 `primary_language: "python"` | unit | `go test ./internal/scanner/... -run TestBuildScanContext_PythonProject` | ❌ Wave 0（新增） |
| FSCAN-01 | 未知語言回傳 `primary_language: "unknown"` | unit | `go test ./internal/scanner/... -run TestBuildScanContext_Unknown` | ❌ Wave 0（新增） |
| FSCAN-02 | `--context-only` 輸出新 ScanContext JSON 欄位 | unit | `go test ./cmd/... -run TestScanContextOnly` | ❌ Wave 0（更新） |
| FSCAN-03 | `config_exists: true` 時 scan 不重建 config | unit | `go test ./internal/scanner/... -run TestScanConfigExists` | ❌ Wave 0 |
| FSCAN-05 | `mysd init` 執行後產生 scaffold 結構 | unit | `go test ./cmd/... -run TestInitScaffold` | ❌ Wave 0（更新） |
| D-13 | `FilterBlockedTasks` 直接依賴失敗的 task 被過濾 | unit | `go test ./internal/executor/... -run TestFilterBlockedTasks_Direct` | ❌ Wave 0 |
| D-13 | `FilterBlockedTasks` 遞移性依賴失敗傳播 | unit | `go test ./internal/executor/... -run TestFilterBlockedTasks_Transitive` | ❌ Wave 0 |
| D-13 | `FilterBlockedTasks` 空 failedIDs 回傳全部 tasks | unit | `go test ./internal/executor/... -run TestFilterBlockedTasks_Empty` | ❌ Wave 0 |

**Note:** SKILL-01~04 和 FSCAN-04 屬於 SKILL.md/agent 層，為 manual-only 驗證（Claude Code session 中端對端測試）。

### Sampling Rate

- **Per task commit:** `go test ./internal/executor/... ./internal/scanner/... ./cmd/...`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/scanner/scanner_test.go` — 現有 PackageInfo 測試需完整替換為新 ScanContext 測試結構
- [ ] `internal/executor/waves_test.go` — 新增 `FilterBlockedTasks` 測試（3 個 test cases）
- [ ] `cmd/model_test.go` — 新建：`TestModelRead`、`TestModelSet`、`TestModelSet_InvalidProfile`
- [ ] `cmd/lang_test.go` — 新建：`TestLangSet`、`TestLangSetAtomic`、`TestLangSet_InvalidLocale`
- [ ] `cmd/scan_test.go` — 更新：`TestScanContextOnly` 驗證新 JSON schema

---

## Open Questions

1. **mysd model read 的 role 順序**
   - What we know: `DefaultModelMap` 是 `map[string]map[string]string`，map iteration 順序不確定
   - What's unclear: 輸出時應使用哪個固定順序？
   - Recommendation: 在 `cmd/model.go` 定義 `knownRoles = []string{"spec-writer", "designer", "planner", "executor", "verifier", "fast-forward", "researcher", "advisor", "proposal-writer", "plan-checker"}`，以此 slice 順序輸出

2. **scan --scaffold-only 的空結構具體內容（Claude's Discretion）**
   - What we know: 需建立 `openspec/` 目錄，但 CONTEXT.md 將細節列為 Claude's Discretion
   - What's unclear: 是否要建立 `openspec/specs/` 子目錄？是否建立空白 `openspec/config.yaml`？
   - Recommendation: 建立 `openspec/` + `openspec/specs/`；不建立 `openspec/config.yaml`（由 `mysd lang set` 建立，確保 locale 有效值）

3. **openspec/config.yaml 的 `schema` 欄位**
   - What we know: 現有 `openspec/config.yaml` 有 `schema: spec-driven` 和 `locale: tw`（非 BCP47）
   - What's unclear: `lang set` 是否需要修正現有 `locale: tw` → `locale: zh-TW`？
   - Recommendation: `lang set` 只更新 `locale` 欄位；如果讀取到非 BCP47 格式，輸出 warning 但仍寫入新值

4. **`mysd init` SKILL.md 更新**
   - What we know: `mysd-init.md` 目前說明 `mysd init` 建立 `.mysd.yaml` 並顯示 config 欄位
   - What's unclear: 改為 scaffold-only 後，`mysd-init.md` 的 SKILL.md 應如何修改？
   - Recommendation: `mysd-init.md` 改為：執行 `mysd init`（等同 scan --scaffold-only），建立後呼叫 `mysd lang set` 互動設定 locale，完成後提示 `Run /mysd:scan to discover existing codebase`

---

## Sources

### Primary (HIGH confidence)

- `internal/scanner/scanner.go` — 現有 ScanContext 結構（直接讀取）
- `internal/config/config.go` — DefaultModelMap、ResolveModel（直接讀取）
- `internal/config/defaults.go` — ProjectConfig struct（直接讀取）
- `internal/spec/openspec_config.go` — WriteOpenSpecConfig、ReadOpenSpecConfig（直接讀取）
- `internal/executor/waves.go` — BuildWaveGroups、adj graph 設計（直接讀取）
- `internal/executor/context.go` — TaskItem struct with Skills field（直接讀取）
- `cmd/worktree.go` — sub-subcommand 設計範本（直接讀取）
- `cmd/scan.go` — --context-only pattern（直接讀取）
- `internal/output/printer.go` — TTY 偵測 + Printer 設計（直接讀取）
- `.planning/phases/07-new-binary-commands-scanner-refactor/07-CONTEXT.md` — 所有鎖定決策（直接讀取）
- `go test ./...` — 確認現有測試全通過（執行驗證）

### Secondary (MEDIUM confidence)

- charmbracelet/lipgloss v1.1.0 table package — TTY-aware table rendering 已在 go.mod 中（未直接查 Context7，但 go.sum 確認版本）

### Tertiary (LOW confidence)

- viper.WriteConfig 保留現有欄位的行為：基於 Viper 的已知行為，未在此 session 中針對 v1.21.0 驗證此具體場景

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — 所有依賴已在 go.mod，直接讀取源碼確認
- Architecture: HIGH — 基於現有程式碼 pattern（worktree.go, scan.go, printer.go）直接推導
- Pitfalls: HIGH (P1-P2, P5-P6) / MEDIUM (P3-P4) — P3 Windows atomic 和 P4 lipgloss TTY 基於已知問題，未在此環境實測
- FilterBlockedTasks: HIGH — waves.go 結構清楚，算法直觀
- Atomic config write: MEDIUM — 需要 Windows 環境實測 os.Rename 行為

**Research date:** 2026-03-26
**Valid until:** 2026-04-26（依賴均為穩定版本，無 fast-moving library）
