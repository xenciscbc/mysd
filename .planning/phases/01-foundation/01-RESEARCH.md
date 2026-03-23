# Phase 1: Foundation - Research

**Researched:** 2026-03-23
**Domain:** Go CLI 骨架、OpenSpec 格式 parser、spec 狀態機、設定檔管理
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Spec 檔案存放在專案根目錄的 `.specs/` 目錄下（相容 OpenSpec 的 `openspec/` 目錄）
- **D-02:** 全域追蹤資訊（roadmap 歷史、UAT 結果等）存放在 `.mysd/` 目錄
- **D-03:** `.mysd/roadmap/` 目錄記錄每個 change 的名稱、狀態、完成日期時間，格式需可被第三方工具解析（用於 roadmap 視覺化）
- **D-04:** 設定檔存放在 `.claude/mysd.yaml`
- **D-05:** 指令風格為 `mysd <verb>` — 簡潔直觀（如 `mysd propose`, `mysd verify`, `mysd ff`）
- **D-06:** 使用 Cobra CLI 框架
- **D-07:** 輸出使用彩色終端輸出（lipgloss），TTY 偵測自動降級為純文字
- **D-08:** Spec 狀態存放在每個 spec 檔案的 YAML frontmatter 中（跟著檔案走）
- **D-09:** Frontmatter 包含 `spec-version` 欄位用於 schema 版本控制（forward compatibility）
- **D-10:** `mysd propose` 先進行 GSD 式的互動提問（了解使用者想做什麼）
- **D-11:** 提問結束後，可選擇是否使用 agent 進行領域研究（像 GSD 的研究模式）
- **D-12:** 最後產出完整的 spec artifacts（proposal.md / specs/ / design.md / tasks.md）
- **D-13:** 實作完成後，自動產生或更新 `.mysd/roadmap/` 下的文件
- **D-14:** 新增 `/mysd:capture` 指令 — 從當前對話中分析並提取要做的變更，然後自動進入 propose 的討論模式
- **D-15:** Parser 自動偵測 `openspec/` 或 `.specs/` 目錄
- **D-16:** 支援讀寫 OpenSpec 的完整 artifact 結構（proposal.md, specs/, design.md, tasks.md）
- **D-17:** Delta Specs 語義（ADDED / MODIFIED / REMOVED）在 parser 層就被識別

### Claude's Discretion

- 具體的 frontmatter schema 欄位設計
- lipgloss 的配色方案
- 錯誤訊息的措辭風格
- `.mysd/roadmap/` 的具體檔案格式（JSON / YAML / Markdown，只要能被第三方工具讀取即可）

### Deferred Ideas (OUT OF SCOPE)

- Debug session 功能 — 可考慮在 v1.x 加入
- Todo / Notes 管理 — 可考慮在 v1.x 加入
- Session pause/resume — 可考慮在 v1.x 加入
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| SPEC-01 | User can create structured spec artifacts (proposal.md, specs/, design.md, tasks.md) via `mysd propose` | `internal/spec/writer.go` scaffold 設計；OpenSpec 目錄結構確認 |
| SPEC-02 | Spec files support RFC 2119 semantic keywords (MUST / SHOULD / MAY) with machine-parseable priority levels | RFC 2119 case-sensitive parser 設計；正則表達式模式確認 |
| SPEC-03 | User can use Delta Specs semantics (ADDED / MODIFIED / REMOVED) to describe changes to existing specs | OpenSpec delta spec 格式研究；`internal/spec/delta.go` 設計 |
| SPEC-04 | Spec status is tracked per-item (PENDING / IN_PROGRESS / DONE / BLOCKED) in spec metadata | YAML frontmatter 設計；`adrg/frontmatter` 讀寫模式確認 |
| SPEC-07 | Spec format uses schema-versioned frontmatter (`spec-version` field) for forward compatibility | frontmatter schema 版本設計；migration 路徑規劃 |
| OPSX-01 | Parser can read existing OpenSpec `openspec/` directory structure | 目錄自動偵測邏輯；`openspec/` vs `.specs/` 相容性確認 |
| OPSX-02 | Parser can read and write OpenSpec's proposal.md / specs/ / design.md / tasks.md format | OpenSpec 實際格式研究（BDD 格式、無 frontmatter 主體）確認 |
| OPSX-03 | Delta Specs support matches OpenSpec's ADDED / MODIFIED / REMOVED semantics | OpenSpec delta 語義研究；`## ADDED` / `## MODIFIED` / `## REMOVED` heading 識別 |
| OPSX-04 | User can point my-ssd at an existing OpenSpec project and run execute/verify without migration | brownfield 相容性設計；`.openspec.yaml` 解析 |
| STAT-01 | Project state tracked in `.specs/STATE.md` for cross-session continuity | WorkflowState struct 設計；JSON state file 路徑確認 |
| STAT-02 | State machine enforces valid transitions (proposed → specced → designed → planned → executed → verified → archived) | Phase enum 設計；transition validation 邏輯 |
| STAT-03 | User can resume interrupted workflow from last valid state | state file read-on-startup；`--resume` flag 設計 |
| CONF-01 | 專案設定檔存放於 `.claude/mysd.yaml`，記憶使用者的偏好預設值 | Viper config 讀取；`~/.claude/mysd.yaml` global + `.claude/mysd.yaml` project 設計 |
| CONF-02 | 設定檔支援：執行模式（single/wave）、agent 數量、atomic commits、TDD 模式、測試產出等可選項目的預設值 | Config struct 欄位設計；yaml.v3 unmarshal 模式確認 |
| CONF-03 | 設定檔支援預設回應語言（response_language）和文件產出語言（document_language） | Config struct 擴充設計 |
| CONF-04 | 所有可選項目在指令執行時可被 flag 覆蓋（flag 優先於設定檔） | Cobra PersistentFlags + Viper BindPFlag 模式確認 |
| DIST-01 | Single Go binary with zero runtime dependencies | Go 模組設計；stdlib-only 執行路徑確認 |
| DIST-02 | Cross-platform support (macOS / Linux / Windows) | GOOS/GOARCH 交叉編譯；filepath.Join 路徑處理確認 |
</phase_requirements>

---

## Summary

Phase 1 建立 my-ssd 所有後續功能的基礎。這個 phase 的核心挑戰有三個：(1) 設計一個能容忍真實世界 OpenSpec 格式變異的 spec parser，(2) 建立清晰的狀態機架構使後續 phase 可以安全地讀取和轉換狀態，(3) 讓 CLI skeleton 從第一天就能被完整測試（thin commands, fat internal）。

研究發現 OpenSpec 的規格文件**沒有** YAML frontmatter — 它們純粹是 Markdown，使用 BDD 格式（Requirement/Scenario，WHEN/THEN/AND）。只有 `.openspec.yaml` 這個 change-level 設定文件含有 metadata（`schema:` 和 `created:` 欄位）。這意味著 my-ssd 的 frontmatter 設計是**自定義的延伸**，不是 OpenSpec 的原生格式 — 這讓 OPSX-04（brownfield 讀取現有 OpenSpec 專案）有一個清楚的處理策略：讀取時優雅降級（frontmatter 缺失 = 使用預設值），寫入時加入 frontmatter 版本資訊。

技術選型已由先前研究確認：Go 1.25.5（環境已有）、Cobra v1.10.2、adrg/frontmatter v0.2.0、yaml.v3 v3.0.1、lipgloss v1.1.0、viper v1.21.0、testify v1.11.1。所有版本已在本機環境驗證。

**Primary recommendation:** 按 ARCHITECTURE.md 定義的 build order 實作：Storage schema → `internal/spec/` → `internal/state/` → `cmd/` skeleton → `internal/config/`。CLI binary 本身在 Phase 1 只需要 `mysd propose` 和 `mysd init` 兩個可執行指令，其他 command stub 只需要 cobra 架構。

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.25.5 (env) | Primary language | 已在環境安裝確認；單一 binary 無 runtime 依賴 |
| github.com/spf13/cobra | v1.10.2 | CLI framework | 已確認版本；kubectl/helm/gh 使用；CONTEXT.md D-06 locked |
| gopkg.in/yaml.v3 | v3.0.1 | YAML parsing | 已確認版本；frontmatter 序列化/反序列化 |
| github.com/adrg/frontmatter | v0.2.0 | Markdown frontmatter extraction | 已確認版本；spec 文件的 frontmatter/body split |
| github.com/spf13/viper | v1.21.0 | Configuration management | 已確認版本；.claude/mysd.yaml 讀寫 + env override |
| github.com/charmbracelet/lipgloss | v1.1.0 | Terminal output styling | 已確認版本；TTY 偵測 + 彩色輸出（CONTEXT.md D-07 locked）|

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| github.com/stretchr/testify | v1.11.1 | Test assertions | 所有 unit test；`assert` (non-fatal) / `require` (fatal) |
| encoding/json (stdlib) | Go stdlib | State file serialization | `.specs/STATE.json` 讀寫；無額外依賴 |
| filepath (stdlib) | Go stdlib | Cross-platform path handling | 所有檔案路徑操作；避免字串拼接 |
| regexp (stdlib) | Go stdlib | RFC 2119 keyword parsing | MUST/SHOULD/MAY 大小寫敏感識別 |
| text/template (stdlib) | Go stdlib | Spec artifact scaffolding | 產生 proposal.md / design.md / tasks.md 模板 |
| os/signal (stdlib) | Go stdlib | TTY detection | `lipgloss.DefaultRenderer().HasDarkBackground()` 前的 TTY 偵測 |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| adrg/frontmatter | goldmark-frontmatter | goldmark-frontmatter 強迫引入 goldmark 依賴；Phase 1 不需要 Markdown 渲染 |
| gopkg.in/yaml.v3 | goccy/go-yaml | goccy 較快但複雜；OpenSpec frontmatter 是簡單 key-value；yaml.v3 已在 cobra 依賴樹中 |
| spf13/viper v1 | koanf | koanf API 更乾淨但社群資源較少；viper 在 cobra 生態系中是標準配對 |
| encoding/json (stdlib) | bbolt / sqlite | state 文件是單一小型 JSON；不需要資料庫 |

**Installation:**

```bash
go get github.com/spf13/cobra@v1.10.2
go get gopkg.in/yaml.v3
go get github.com/adrg/frontmatter@v0.2.0
go get github.com/spf13/viper@v1
go get github.com/charmbracelet/lipgloss@v1
go get github.com/stretchr/testify@v1
```

**Version verification (已驗證 2026-03-23):**

```
github.com/spf13/cobra v1.10.2      ✓
gopkg.in/yaml.v3 v3.0.1             ✓
github.com/adrg/frontmatter v0.2.0  ✓
github.com/spf13/viper v1.21.0      ✓
github.com/charmbracelet/lipgloss v1.1.0  ✓
github.com/stretchr/testify v1.11.1 ✓
Go runtime: 1.25.5 (windows/amd64)  ✓
```

---

## Architecture Patterns

### Recommended Project Structure

```
mysd/
├── main.go                    # cobra root cmd entry point
├── cmd/
│   ├── root.go                # persistent flags, version, TTY detection
│   ├── propose.go             # mysd propose (Phase 1: interactive scaffold)
│   ├── init.go                # mysd init (Phase 1: .claude/mysd.yaml setup)
│   ├── spec.go                # stub — Phase 2
│   ├── design.go              # stub — Phase 2
│   ├── plan.go                # stub — Phase 2
│   ├── execute.go             # stub — Phase 2
│   ├── verify.go              # stub — Phase 3
│   └── archive.go             # stub — Phase 3
├── internal/
│   ├── spec/
│   │   ├── parser.go          # parse .openspec.yaml + artifact files
│   │   ├── schema.go          # Go structs: Change, Requirement, Task, ProposalDoc
│   │   ├── writer.go          # scaffold .specs/changes/[name]/ 目錄結構
│   │   ├── delta.go           # ADDED/MODIFIED/REMOVED delta 識別
│   │   └── detector.go        # openspec/ vs .specs/ 目錄自動偵測
│   ├── state/
│   │   ├── state.go           # WorkflowState struct + read/write STATE.json
│   │   └── transitions.go     # valid state transition 驗證
│   ├── config/
│   │   ├── config.go          # ProjectConfig struct + viper load
│   │   └── defaults.go        # convention-over-config 預設值
│   └── output/
│       ├── printer.go         # lipgloss styled output + TTY fallback
│       └── colors.go          # 配色方案常數
├── testdata/
│   └── fixtures/
│       ├── openspec-project/  # 真實 OpenSpec 目錄結構 brownfield fixture
│       └── mysd-project/      # my-ssd 格式 fixture（含 frontmatter）
├── go.mod
├── go.sum
└── Makefile
```

### Pattern 1: Thin Commands, Fat Internal

**What:** `cmd/` 文件不含業務邏輯。每個 command 呼叫 `internal/` 函數，然後格式化輸出。

**When to use:** 從第一天開始，沒有例外。

**Example:**

```go
// cmd/propose.go
func runPropose(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load(".")
    if err != nil {
        return err
    }
    change, err := spec.Scaffold(cfg, args)
    if err != nil {
        return err
    }
    output.PrintScaffoldResult(cmd.OutOrStdout(), change)
    return nil
}
```

### Pattern 2: Spec as Struct, Not String

**What:** 在邊界層（`internal/spec/parser.go`）解析 Markdown 為 typed Go structs。所有下游代碼操作 `spec.Change`、`spec.Requirement`、`spec.Task`，永不操作原始字串。

**When to use:** 從第一個 parser 實作開始。

**Example:**

```go
// internal/spec/schema.go
type RFC2119Keyword string
const (
    Must   RFC2119Keyword = "MUST"
    Should RFC2119Keyword = "SHOULD"
    May    RFC2119Keyword = "MAY"
)

type DeltaOp string
const (
    DeltaAdded    DeltaOp = "ADDED"
    DeltaModified DeltaOp = "MODIFIED"
    DeltaRemoved  DeltaOp = "REMOVED"
    DeltaNone     DeltaOp = ""
)

type Requirement struct {
    ID      string
    Text    string
    Keyword RFC2119Keyword
    DeltaOp DeltaOp
    Status  ItemStatus  // PENDING | IN_PROGRESS | DONE | BLOCKED
}

type Change struct {
    Name     string
    Dir      string
    Proposal ProposalDoc
    Specs    []Requirement
    Design   DesignDoc
    Tasks    []Task
    Meta     ChangeMeta  // .openspec.yaml 的內容
}
```

### Pattern 3: Explicit Workflow State Machine

**What:** 狀態存在 `.specs/STATE.json`。Commands 在執行前驗證狀態轉換是否合法。

**When to use:** 每個 command 的第一行就讀取並驗證狀態。

**Example:**

```go
// internal/state/state.go
type Phase string
const (
    PhaseNone     Phase = ""
    PhaseProposed Phase = "proposed"
    PhaseSpecced  Phase = "specced"
    PhaseDesigned Phase = "designed"
    PhasePlanned  Phase = "planned"
    PhaseExecuted Phase = "executed"
    PhaseVerified Phase = "verified"
    PhaseArchived Phase = "archived"
)

// ValidTransitions 定義合法的狀態轉換
var ValidTransitions = map[Phase][]Phase{
    PhaseNone:     {PhaseProposed},
    PhaseProposed: {PhaseSpecced},
    PhaseSpecced:  {PhaseDesigned},
    PhaseDesigned: {PhasePlanned},
    PhasePlanned:  {PhaseExecuted},
    PhaseExecuted: {PhaseVerified},
    PhaseVerified: {PhaseArchived, PhaseExecuted}, // FAIL 可以重新 execute
    PhaseArchived: {PhaseProposed},                // 新的 change
}

type WorkflowState struct {
    ChangeName string    `json:"change_name"`
    Phase      Phase     `json:"phase"`
    LastRun    time.Time `json:"last_run"`
    VerifyPass *bool     `json:"verify_pass,omitempty"`
}
```

### Pattern 4: OpenSpec Brownfield Detection

**What:** Parser 自動偵測 `openspec/` 或 `.specs/` 目錄（D-15）。對現有 OpenSpec 項目（無 frontmatter 的規格文件），使用優雅降級策略。

**When to use:** `spec.Parser` 初始化時。

**Example:**

```go
// internal/spec/detector.go
func DetectSpecDir(root string) (string, SpecDirFlavor, error) {
    // 優先 .specs/ (my-ssd native)
    if _, err := os.Stat(filepath.Join(root, ".specs")); err == nil {
        return filepath.Join(root, ".specs"), FlavorMySD, nil
    }
    // fallback openspec/ (brownfield)
    if _, err := os.Stat(filepath.Join(root, "openspec")); err == nil {
        return filepath.Join(root, "openspec"), FlavorOpenSpec, nil
    }
    return "", FlavorNone, ErrNoSpecDir
}
```

### Pattern 5: Viper + Cobra Flag Override

**What:** Viper 讀取 `.claude/mysd.yaml`，Cobra flag 透過 `BindPFlag` 覆蓋設定值（D-04, CONF-04）。

**Example:**

```go
// cmd/root.go
func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().String("config", "", "config file (default .claude/mysd.yaml)")
    rootCmd.PersistentFlags().String("lang", "", "response language override")
    viper.BindPFlag("response_language", rootCmd.PersistentFlags().Lookup("lang"))
}

func initConfig() {
    viper.SetConfigName("mysd")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".claude")            // project-level
    viper.AddConfigPath("$HOME/.claude")     // user-level
    viper.AutomaticEnv()
    viper.ReadInConfig() // 忽略 not found 錯誤（convention over config）
}
```

### Anti-Patterns to Avoid

- **直接在 `cmd/` 中解析 Markdown：** 破壞測試性，格式變更時需改動 command 文件
- **以檔案系統狀態推斷 workflow phase：** 不可靠，使用明確的 `STATE.json`
- **硬編碼 spec heading 名稱：** 如 `strings.Contains(line, "## Requirements")` — 使用 schema-driven 解析
- **在 SKILL.md 中放業務邏輯：** 超出 context budget，不可測試
- **RFC 2119 case-insensitive 解析：** `must` (lowercase) 不是需求；只有 `MUST` 是

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| CLI argument parsing | 自製 arg parser | cobra v1.10.2 | 子命令、flag 繼承、help 生成、shell completion — 手寫有大量 edge case |
| YAML frontmatter extraction | 手寫 `---` 邊界解析 | adrg/frontmatter v0.2.0 | 處理 TOML/JSON/YAML 混合；邊界條件（空 frontmatter、多行值）複雜 |
| YAML serialization | 手寫 YAML | gopkg.in/yaml.v3 | anchor/alias、多行字串、類型推斷複雜 |
| Config file management | 手寫 YAML config reader | spf13/viper v1.21.0 | env override、多層設定文件、default values 的優先級邏輯 |
| Terminal color output | ANSI escape codes 手寫 | lipgloss v1.1.0 | TTY 偵測、Windows 相容性、顏色深度偵測 — 手寫在 Windows 上會失敗 |
| Test assertions | 自製 assert helpers | testify v1.11.1 | diff output、深層比較、mock — stdlib testing 缺少這些 |

**Key insight:** Go stdlib 處理了最低層的 I/O，但 CLI 慣用模式（subcommands、config layering、TTY detection）都有已知的 edge case — 使用 battle-tested 函式庫可以避免 2-3 個月的 bug hunting。

---

## OpenSpec Format Research Findings

### 實際 OpenSpec 文件格式（HIGH confidence — 直接驗證）

來源：`github.com/Fission-AI/OpenSpec` 真實 changes 目錄驗證

**Change 目錄結構：**
```
openspec/changes/[change-name]/
├── .openspec.yaml          # change-level 設定（schema, created）
├── proposal.md             # 提案文件（純 Markdown，無 frontmatter）
├── specs/
│   └── [capability]/
│       └── spec.md         # BDD 格式規格（無 frontmatter）
├── design.md               # 技術決策（如存在）
└── tasks.md                # 實作清單（如存在）
```

**`.openspec.yaml` 格式：**
```yaml
schema: spec-driven
created: 2026-01-20
```

**`spec.md` 文件格式（BDD 風格）：**
```markdown
## Requirement: [Name]

[Description]

### Scenario: [Name]
WHEN [condition]
THEN [expected outcome]
AND [additional condition]
```

**關鍵發現：OpenSpec spec 文件沒有 YAML frontmatter**

這意味著：
1. my-ssd 的 frontmatter（status、spec-version 等）是**額外的 metadata 延伸**
2. 讀取現有 OpenSpec 項目（OPSX-04）時，frontmatter 缺失是正常狀態，不是錯誤
3. `adrg/frontmatter` 的 `Parse()` 返回空 struct + body 完全有效

**Delta Spec 格式：**
```markdown
## ADDED Requirements

### Requirement: [新需求名稱]
[需求描述]

## MODIFIED Requirements

### Requirement: [現有需求名稱]
[修改後的完整需求文字]

## REMOVED Requirements

### Requirement: [要刪除的需求名稱]
```

### my-ssd Frontmatter Schema 設計（Claude's Discretion）

my-ssd 為其自有格式的 spec 文件加入 frontmatter（D-08, D-09, SPEC-07）：

**proposal.md frontmatter:**
```yaml
---
spec-version: "1"
change: add-dark-mode
status: proposed
created: 2026-03-23
updated: 2026-03-23
---
```

**specs/[capability]/spec.md frontmatter:**
```yaml
---
spec-version: "1"
capability: user-authentication
delta: ADDED           # ADDED | MODIFIED | REMOVED | "" (無 delta)
status: pending        # pending | in_progress | done | blocked
---
```

**tasks.md frontmatter:**
```yaml
---
spec-version: "1"
total: 5
completed: 0
---
```

**設計原則：**
- `spec-version: "1"` 用引號包圍（避免 YAML 將版本解析為整數）
- brownfield OpenSpec 文件（無 frontmatter）：`spec-version` 缺失 = 視為 `"0"` (legacy)
- 寫入時若 frontmatter 已存在則更新，不存在則插入

### RFC 2119 Parsing Pattern

```go
// internal/spec/parser.go
// Source: RFC 2119 https://datatracker.ietf.org/doc/html/rfc2119
// CRITICAL: case-sensitive — only UPPERCASE keywords are normative
var (
    reMust    = regexp.MustCompile(`\bMUST\b|\bMUST NOT\b|\bREQUIRED\b|\bSHALL\b|\bSHALL NOT\b`)
    reShould  = regexp.MustCompile(`\bSHOULD\b|\bSHOULD NOT\b|\bRECOMMENDED\b`)
    reMay     = regexp.MustCompile(`\bMAY\b|\bOPTIONAL\b`)
)

func extractKeyword(line string) RFC2119Keyword {
    if reMust.MatchString(line) {
        return Must
    }
    if reShould.MatchString(line) {
        return Should
    }
    if reMay.MatchString(line) {
        return May
    }
    return ""
}
```

### Roadmap 追蹤格式（Claude's Discretion — 建議 YAML）

`.mysd/roadmap/` 目錄，每個 change 一個文件：

```yaml
# .mysd/roadmap/add-dark-mode.yaml
name: add-dark-mode
status: verified          # proposed | specced | designed | planned | executed | verified | archived
started: 2026-03-23
completed: 2026-03-24
verify_pass: true
phases_completed:
  - proposed: 2026-03-23T10:00:00Z
  - specced: 2026-03-23T11:30:00Z
  - designed: 2026-03-23T13:00:00Z
  - planned: 2026-03-23T14:00:00Z
  - executed: 2026-03-24T09:00:00Z
  - verified: 2026-03-24T10:00:00Z
```

YAML 格式原因：可被 Python/Node.js 工具直接讀取；比 JSON 更易用 yq 查詢；比 Markdown 更易機器解析 (D-03 要求)。Mermaid gantt chart 可以直接從此格式生成。

---

## Common Pitfalls

### Pitfall 1: Spec Format Lock-In（CRITICAL）

**What goes wrong:** Parser 硬編碼 `"## Requirements"` 等 heading 名稱。真實 OpenSpec 項目使用 `"## Requirement: [Name]"` 格式（帶冒號和名稱），導致 brownfield 讀取失敗。

**Why it happens:** 寫 `strings.Contains(line, "## Requirements")` 比 schema-driven 解析快 — 但在 Phase 1 就固化這個選擇。

**How to avoid:** 使用 regex 而不是 literal string matching：
```go
var reRequirementHeading = regexp.MustCompile(`^#{1,3}\s+Requirement:?\s*(.*)$`)
```
接受 `## Requirements`、`## Requirement: User Auth`、`### Requirement` 等變體。

**Warning signs:** Parser code 出現 `strings.HasPrefix(line, "## Requirements")`

### Pitfall 2: RFC 2119 Case-Insensitive Matching

**What goes wrong:** `strings.ToLower(line)` 後再找 `must` — 將 prose 中的 "you must understand" 解析為 MUST requirement。

**Why it happens:** 大小寫不敏感 matching 看起來更「容錯」。

**How to avoid:** RFC 2119 明確規定只有 UPPERCASE 關鍵字是規範性的。Parser 必須 case-sensitive。測試用例必須包含 `"you must not"` (lowercase) 返回 0 requirements。

**Warning signs:** 測試只用大寫 `MUST` 測試，沒有小寫反例

### Pitfall 3: Inferring Phase from Filesystem

**What goes wrong:** 檢查 "if design.md exists → phase is designed"，沒有明確的 state file。

**Why it happens:** 看起來省去了維護 state file 的工作。

**How to avoid:** 維護 `.specs/STATE.json` 作為 workflow cursor 的唯一事實來源。檔案系統內容是 data；state file 是 workflow cursor。interrupted command 後的部分寫入會造成 filesystem inference 誤判。

**Warning signs:** `cmd/` 文件中出現 `os.Stat(filepath.Join(dir, "design.md"))`

### Pitfall 4: Plugin SKILL.md 超出 Context Budget

**What goes wrong:** SKILL.md 寫入太多指令，超過 500 行，Claude 在 context budget 下默默排除部分 skill。

**Why it happens:** Phase 1 建立 plugin skeleton 時容易過度寫入。

**How to avoid:** Phase 1 只寫 stub SKILL.md（50 行以內）。詳細指令在 Phase 2 實際功能完成後再補充。使用 `disable-model-invocation: true` 只在顯式呼叫時載入。

**Warning signs:** SKILL.md 文件超過 400 行

### Pitfall 5: Windows 路徑處理（filepath.Join vs 字串拼接）

**What goes wrong:** `root + "/.specs/" + changeName` 在 Windows 上生成錯誤路徑（混合 `/` 和 `\`）。

**Why it happens:** 開發者在 macOS/Linux 上開發，字串拼接在 Unix 上正確。

**How to avoid:** 所有路徑操作使用 `filepath.Join()`。跨平台測試 fixture 使用 `filepath.ToSlash()` 做比較。

**Warning signs:** 任何 `dir + "/" + file` 型的路徑拼接

---

## Code Examples

Verified patterns from official sources and Go stdlib:

### adrg/frontmatter: Parse Spec File

```go
// Source: github.com/adrg/frontmatter v0.2.0 README
// Handles: YAML frontmatter present, absent, or empty — all cases valid
import "github.com/adrg/frontmatter"

type SpecFrontmatter struct {
    SpecVersion string `yaml:"spec-version"`
    Capability  string `yaml:"capability"`
    Delta       string `yaml:"delta"`
    Status      string `yaml:"status"`
}

func ParseSpecFile(path string) (*SpecFrontmatter, string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, "", err
    }
    defer f.Close()

    var fm SpecFrontmatter
    body, err := frontmatter.Parse(f, &fm)
    if err != nil {
        return nil, "", err
    }
    // 若無 frontmatter，fm 為零值（SpecVersion == ""）— 這是合法的 brownfield 狀態
    return &fm, string(body), nil
}
```

### yaml.v3: Write Frontmatter Back to File

```go
// Source: gopkg.in/yaml.v3 documentation
import "gopkg.in/yaml.v3"

func WriteFrontmatter(path string, fm interface{}, body string) error {
    fmBytes, err := yaml.Marshal(fm)
    if err != nil {
        return err
    }
    content := fmt.Sprintf("---\n%s---\n\n%s", string(fmBytes), body)
    return os.WriteFile(path, []byte(content), 0644)
}
```

### Cobra: Command with Viper Config Override

```go
// Source: STACK.md + CONTEXT.md D-06, CONF-04
// cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "mysd",
    Short: "Spec-Driven Development for Claude Code",
}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().String("lang", "", "response language (overrides config)")
    viper.BindPFlag("response_language", rootCmd.PersistentFlags().Lookup("lang"))
}

func initConfig() {
    viper.SetConfigName("mysd")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".claude")
    home, _ := os.UserHomeDir()
    viper.AddConfigPath(filepath.Join(home, ".claude"))
    viper.AutomaticEnv()
    viper.ReadInConfig() // silent on not-found
}
```

### lipgloss: TTY Detection + Styled Output

```go
// Source: charmbracelet/lipgloss v1.1.0 + CONTEXT.md D-07
// internal/output/printer.go
import (
    "github.com/charmbracelet/lipgloss"
    "os"
)

var (
    isTTY    = isTerminal()
    StyleOK  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))  // green
    StyleErr = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))   // red
    StyleInfo = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // blue
)

func isTerminal() bool {
    fi, err := os.Stdout.Stat()
    if err != nil {
        return false
    }
    return (fi.Mode() & os.ModeCharDevice) != 0
}

func PrintSuccess(w io.Writer, msg string) {
    if isTTY {
        fmt.Fprintln(w, StyleOK.Render("✓ "+msg))
    } else {
        fmt.Fprintln(w, "OK: "+msg)  // plain text fallback
    }
}
```

### State Machine: Transition Validation

```go
// Source: ARCHITECTURE.md Pattern 3
// internal/state/transitions.go
func (s *WorkflowState) CanTransitionTo(next Phase) bool {
    allowed, ok := ValidTransitions[s.Phase]
    if !ok {
        return false
    }
    for _, p := range allowed {
        if p == next {
            return true
        }
    }
    return false
}

func (s *WorkflowState) Transition(next Phase) error {
    if !s.CanTransitionTo(next) {
        return fmt.Errorf("cannot transition from %q to %q", s.Phase, next)
    }
    s.Phase = next
    s.LastRun = time.Now()
    return s.Save()
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| OpenSpec `commands/` directory | OpenSpec 繼續使用其自有命令系統 | N/A | my-ssd 使用 Claude Code SKILL.md 格式（`skills/<name>/SKILL.md`）|
| cobra v1.8.x | cobra v1.10.2 | Dec 2024 | pflag v1.0.9 required；go.yaml.in/yaml/v3 替代 gopkg.in/yaml.v3 |
| golangci-lint v1 | golangci-lint v2 | March 2025 | 新設定格式：`linters.default: standard`；Go 1.22+ 最低需求 |
| goreleaser Homebrew formulae | goreleaser Homebrew casks | June 2025 | formulae 已廢棄；用 `brews` config 的 cask 設定 |
| lipgloss v0.x | lipgloss v1.1.0 | 2025 | API 穩定化；v1.1.0 已確認為最新版本 |

**Deprecated/outdated:**
- GoReleaser `brews` formulae 設定：2025 年 6 月起廢棄；Phase 4 分發時使用 cask 設定
- `cobra-cli` 加入 go.mod：只應作為 global binary 安裝，不是 module 依賴
- viper v2：尚未發布，不要依賴 pre-release

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go runtime | All compilation | ✓ | 1.25.5 (windows/amd64) | — |
| git | `.specs/` version control | ✓ | (assumed, standard dev env) | — |
| go test | Unit testing | ✓ | Go 1.25.5 stdlib | — |
| golangci-lint | Linting | ✗ | — | 跳過 lint gate，Phase 1 後補 |
| cobra-cli | Command scaffolding | ✗ | — | 手寫 command boilerplate（可接受）|

**Missing dependencies with no fallback:**
- 無 — Go runtime 已確認，Phase 1 的所有核心功能可以執行

**Missing dependencies with fallback:**
- golangci-lint：Phase 1 可用 `go vet` 替代；Phase 2 前正式安裝
- cobra-cli：Phase 1 手寫 command 文件即可，不需要 scaffolding tool

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` + `github.com/stretchr/testify` v1.11.1 |
| Config file | 無獨立設定文件 — Go test 慣例（`*_test.go` 放在同一 package） |
| Quick run command | `go test ./internal/spec/... -v -run TestParse` |
| Full suite command | `go test ./... -v` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| SPEC-01 | Scaffold 建立正確的目錄結構 | unit | `go test ./internal/spec/... -run TestScaffold` | ❌ Wave 0 |
| SPEC-02 | RFC 2119 MUST/SHOULD/MAY 大寫解析正確 | unit | `go test ./internal/spec/... -run TestRFC2119` | ❌ Wave 0 |
| SPEC-02 | RFC 2119 lowercase `must` 返回 0 requirements（反例）| unit | `go test ./internal/spec/... -run TestRFC2119CaseSensitive` | ❌ Wave 0 |
| SPEC-03 | Delta ADDED/MODIFIED/REMOVED 在 heading 中正確識別 | unit | `go test ./internal/spec/... -run TestDelta` | ❌ Wave 0 |
| SPEC-04 | Frontmatter status 欄位讀寫正確 | unit | `go test ./internal/spec/... -run TestFrontmatter` | ❌ Wave 0 |
| SPEC-07 | spec-version 欄位出現在所有新建文件中 | unit | `go test ./internal/spec/... -run TestSpecVersion` | ❌ Wave 0 |
| OPSX-01 | 自動偵測 `openspec/` 目錄 | unit | `go test ./internal/spec/... -run TestDetector` | ❌ Wave 0 |
| OPSX-02 | 讀取無 frontmatter 的 OpenSpec 文件不報錯 | unit | `go test ./internal/spec/... -run TestBrownfield` | ❌ Wave 0 |
| OPSX-03 | Delta spec round-trip（parse → write → parse）正確 | integration | `go test ./internal/spec/... -run TestDeltaRoundTrip` | ❌ Wave 0 |
| OPSX-04 | 指向 OpenSpec `openspec/` 目錄的命令可執行 | integration | `go test ./internal/spec/... -run TestOpenSpecCompat` | ❌ Wave 0 |
| STAT-01 | STATE.json 正確讀寫 | unit | `go test ./internal/state/... -run TestStateReadWrite` | ❌ Wave 0 |
| STAT-02 | 非法狀態轉換返回 error | unit | `go test ./internal/state/... -run TestTransitions` | ❌ Wave 0 |
| STAT-03 | State resume：讀取現有 STATE.json 恢復正確狀態 | unit | `go test ./internal/state/... -run TestResume` | ❌ Wave 0 |
| CONF-01 | `.claude/mysd.yaml` 讀取正確 | unit | `go test ./internal/config/... -run TestConfigLoad` | ❌ Wave 0 |
| CONF-04 | CLI flag 優先於設定檔 | unit | `go test ./internal/config/... -run TestFlagOverride` | ❌ Wave 0 |
| DIST-02 | `filepath.Join` 用於所有路徑（無字串拼接）| unit | `go test ./internal/... -run TestPaths` | ❌ Wave 0 |

### Sampling Rate

- **Per task commit:** `go test ./internal/spec/... -v` 或 `go test ./internal/state/...`（依修改的 package 而定）
- **Per wave merge:** `go test ./... -v`
- **Phase gate:** `go test ./... -v` 全部 GREEN 後才進入 `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/spec/parser_test.go` — covers SPEC-01, SPEC-02, SPEC-03, SPEC-04, SPEC-07, OPSX-01, OPSX-02, OPSX-03
- [ ] `internal/spec/detector_test.go` — covers OPSX-01, OPSX-04
- [ ] `internal/spec/delta_test.go` — covers SPEC-03, OPSX-03
- [ ] `internal/state/state_test.go` — covers STAT-01, STAT-02, STAT-03
- [ ] `internal/config/config_test.go` — covers CONF-01, CONF-02, CONF-04
- [ ] `testdata/fixtures/openspec-project/` — brownfield fixture（無 frontmatter 的 OpenSpec 結構）
- [ ] `testdata/fixtures/mysd-project/` — my-ssd 格式 fixture（含 frontmatter）
- [ ] Framework install: `go get github.com/stretchr/testify@v1` — if not already in go.mod

---

## Project Constraints (from CLAUDE.md)

| Directive | Type | Constraint |
|-----------|------|-----------|
| Tech stack: Go | Hard | 單一 binary，跨平台編譯 — 不使用 Node.js / Python |
| 相容性: OpenSpec 格式 | Hard | 必須能讀寫 OpenSpec 的 `openspec/` 目錄結構 |
| Plugin 形式: Claude Code slash commands + agent definitions | Hard | 使用 SKILL.md 格式，不是 MCP server |
| Convention over configuration | Design principle | 預設即好用；所有 defaults 開箱即用 |
| GSD Workflow Enforcement | Process | 所有代碼變更通過 GSD command（`/gsd:execute-phase`）進行 |
| Cobra CLI | Locked (D-06) | 使用 cobra v1.10.2，不用 urfave/cli |
| `.claude/mysd.yaml` config path | Locked (D-04) | 不用 `.mysdrc`、`mysd.config.yaml` 等其他路徑 |
| lipgloss for terminal output | Locked (D-07) | 使用 lipgloss，不用手寫 ANSI codes 或 bubbletea |
| Frontmatter 含 spec-version | Locked (D-09) | 所有新建文件必須有 `spec-version` 欄位 |

---

## Open Questions

1. **`mysd propose` 的互動提問機制（D-10）**
   - What we know: 應該是 GSD 式深度提問（不是表單填寫），但 Phase 1 的 Go binary 只能印出 prompt 讓 Claude 在 SKILL.md 層做提問
   - What's unclear: `propose` 指令的 Go binary 部分做什麼？只是 scaffold artifact 結構，還是也做 interactive prompting？
   - Recommendation: Phase 1 的 Go binary `propose` 只做 scaffold（建立目錄結構 + 空文件）；interactive questioning 在 SKILL.md 層由 Claude 執行。Binary 接受 `--name <change-name>` flag 作為 scaffold 目標。

2. **`.mysd/roadmap/` 格式 selection（Claude's Discretion）**
   - What we know: 需要可被第三方工具讀取；需支援 roadmap 視覺化
   - What's unclear: YAML vs JSON vs Markdown 哪個最易被常見工具（yq, jq, Python）讀取？
   - Recommendation: YAML（每個 change 一個 `.yaml` 文件）— 可用 `yq` 直接查詢，也可被 Python `yaml.safe_load()` 讀取；比 JSON 更易閱讀；比 Markdown 更易機器解析

3. **`mysd init` 指令的互動機制**
   - What we know: 初始化 `.claude/mysd.yaml`（WCMD-11 在 Phase 2）；CONF-01-04 要求設定檔支援
   - What's unclear: Phase 1 是否需要完整的 `mysd init` 還是只需要 config 讀取？
   - Recommendation: Phase 1 實作 config 讀取 + defaults（`internal/config/`）；`mysd init` 指令可以是 stub，打印 "Not yet implemented, create .claude/mysd.yaml manually" — 完整互動式初始化在 Phase 2

---

## Sources

### Primary (HIGH confidence)
- Go module registry (`go list -m ... @latest`) — 所有版本數字直接在環境驗證（2026-03-23）
- `github.com/Fission-AI/OpenSpec` — 直接讀取 changes 目錄的真實文件確認格式
- `.planning/research/STACK.md` — 技術棧研究，已引用官方文件
- `.planning/research/ARCHITECTURE.md` — 系統架構設計，三層架構確認
- `.planning/research/PITFALLS.md` — 已驗證的陷阱，含原始來源
- RFC 2119 (https://datatracker.ietf.org/doc/html/rfc2119) — MUST/SHOULD/MAY 規範
- `01-CONTEXT.md` — 使用者決策確認

### Secondary (MEDIUM confidence)
- OpenSpec `redreamality.com/garden/notes/openspec-guide/` — delta spec 語義確認（ADDED/MODIFIED/REMOVED）
- OpenSpec `docs/workflows.md` — workflow 模式確認
- OpenSpec `docs/concepts.md` — core concepts 確認

### Tertiary (LOW confidence)
- 無 — 所有關鍵發現均有 PRIMARY 或 SECONDARY 來源支撐

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — 所有版本在本機環境直接驗證（`go list -m`）
- Architecture: HIGH — 來自 ARCHITECTURE.md（已引用官方文件）+ OpenSpec 真實格式驗證
- OpenSpec format: HIGH — 直接讀取 Fission-AI/OpenSpec 真實 changes 目錄確認
- Pitfalls: HIGH — 來自 PITFALLS.md（已有原始來源）+ RFC 2119 直接確認

**Research date:** 2026-03-23
**Valid until:** 2026-04-23（30 天，技術棧穩定；若 OpenSpec 有重大格式更新需重驗證）
