# Phase 5: Schema Foundation & Plan-Checker - Research

**Researched:** 2026-03-25
**Domain:** Go struct extension (yaml.v3), pure-function package design, agent definition authoring, openspec/config.yaml writer
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Plan-checker 架構**
- D-01: MUST 覆蓋率檢查在 Go binary 層實作（`internal/planchecker/` 新 package），使用 structured ID matching（satisfies 欄位 vs MUST IDs），不使用 AI 語意推測
- D-02: 觸發方式沿用現有 `mysd plan --check` flag，plan 完成後自動執行檢查
- D-03: 檢查失敗時輸出 uncovered MUST IDs 清單 + 覆蓋率比例（簡潔模式），agent 層負責渲染互動 UI 讓使用者選擇自動補齊或手動調整（對應 FSCHEMA-06）
- D-04: PlanningContext JSON 擴展 WaveGroups、WorktreeDir、AutoMode 欄位，WaveGroups 由 Go binary 對 tasks.md 的 depends 欄位做 topological sort 計算得出

**Model profile 對應表**
- D-05: 4 個新 agent roles（researcher, advisor, proposal-writer, plan-checker）在 quality/balanced/budget 三層全部映射到 sonnet，與現有 6 roles 保持一致
- D-06: budget 層新 roles 也用 sonnet（不降為 haiku），保持整體一致性

**openspec/config.yaml 內容**
- D-07: 最小必要欄位：project name, locale (BCP47), spec_dir, created — 其餘由 convention-over-config 推導
- D-08: locale 欄位使用 BCP47 標準格式（zh-TW, en-US, ja-JP），與 Go 的 golang.org/x/text/language 直接相容
- D-09: openspec/config.yaml 的 locale 為 source of truth，mysd.yaml 的 response_language/document_language 讀取時參考。/mysd:lang 修改時兩者原子同步
- D-10: config.yaml 是 OpenSpec 標準格式（project-level），mysd.yaml 是 mysd 專用配置。兩者職責分離

**Schema 向後相容策略**
- D-11: 新欄位使用 `omitempty` YAML tag（`yaml:"depends,omitempty"` 等），舊 tasks.md 讀取時新欄位為 nil/empty，寫回時不輸出空欄位。零遷移、零 migration
- D-12: TaskItem（executor/context.go JSON 輸出）同步擴展 Depends/Files/Satisfies/Skills 欄位，讓 agent 執行時能看到完整資訊，Phase 6 wave grouping 直接可用

### Claude's Discretion
- plan-checker 的具體 JSON 輸出結構（遵循簡潔模式：uncovered IDs + 覆蓋率）
- topological sort 演算法選擇（Kahn's vs DFS-based）
- config.yaml writer 的錯誤處理細節
- 新欄位的測試案例覆蓋範圍

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| FSCHEMA-01 | TaskEntry 支援 `depends` 欄位標記 task 間依賴關係 | 直接擴展 `internal/spec/schema.go` TaskEntry struct，加 `yaml:"depends,omitempty"` |
| FSCHEMA-02 | TaskEntry 支援 `files` 欄位標記 task 會修改的檔案 | 同 FSCHEMA-01，加 `yaml:"files,omitempty"` |
| FSCHEMA-03 | TaskEntry 支援 `satisfies` 欄位對應 MUST requirement IDs | 同 FSCHEMA-01，加 `yaml:"satisfies,omitempty"` — plan-checker 比對的關鍵欄位 |
| FSCHEMA-04 | TaskEntry 支援 `skills` 欄位標記執行時建議使用的 slash commands | 同 FSCHEMA-01，加 `yaml:"skills,omitempty"` |
| FSCHEMA-05 | Plan-checker 自動驗證所有 MUST items 都有 task 的 `satisfies` 對應 | 新 `internal/planchecker/` package，純函數 `CheckCoverage(tasks, mustIDs) CoverageResult` |
| FSCHEMA-06 | Plan-checker 未通過時顯示缺口，互動式詢問自動補齊或手動調整 | planchecker 輸出 JSON 給 `mysd-plan-checker` agent；agent 負責互動 UI |
| FSCHEMA-07 | openspec/config.yaml writer 可產生/讀取 OpenSpec config（含 project metadata + locale） | 新 writer 函數於 `internal/spec/` 或獨立 writer.go，使用 yaml.v3 序列化 |
| FAGENT-04 | 新增 `mysd-plan-checker` agent definition（驗證 MUST 覆蓋率） | 參考 mysd-verifier.md 格式；agent 接收 checker JSON 輸出並渲染互動 UI |
| FMODEL-01 | Model profile 分層表涵蓋所有新 agents（researcher, advisor, proposal-writer, plan-checker） | 擴展 `DefaultModelMap` 於 `internal/config/config.go` |
| FMODEL-02 | Orchestrator（SKILL.md）動態指定 model 參數給每個 spawned agent | `ResolveModel()` 已存在且可直接使用；`mysd plan --context-only` 輸出需加入 model 欄位 |
| FMODEL-03 | quality/balanced/budget 三層完整對應表 | 4 個新 roles 全部映射到 sonnet（D-05/D-06） |
</phase_requirements>

---

## Summary

Phase 5 是 v1.1 里程碑的 foundation layer。它的工作範圍窄且清晰：在現有程式碼上做 additive-only 修改，不改變任何現有行為。所有改動都屬於三類之一：(1) Go struct 欄位擴展（`TaskEntry`、`TaskItem`、`ExecutionContext`、`ProjectConfig`）；(2) 新增一個純函數 package（`internal/planchecker/`）；(3) 新增一個 agent definition 文件（`plugin/agents/mysd-plan-checker.md`）。

現有測試套件全部通過（10 packages，無 test failures）。所有要修改的套件都有充足的測試覆蓋。Phase 5 的工作只需按照現有測試模式新增測試案例，不需要建立新的 test infrastructure。

**Primary recommendation:** 依照 `internal/spec/schema.go` → `internal/config/config.go` → `internal/executor/context.go` → `internal/planchecker/` → `cmd/plan.go` → `plugin/agents/mysd-plan-checker.md` 的順序實作，每個步驟都有獨立的測試可驗證。openspec/config.yaml writer 可在任意位置插入，因為它沒有跨套件依賴。

---

## Standard Stack

### Core（全部已在 go.mod 中）

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `gopkg.in/yaml.v3` | v3.0.1 | TaskEntry YAML round-trip；openspec/config.yaml 序列化 | 已在 go.mod；cobra 內部也使用相同 module |
| `github.com/adrg/frontmatter` | v0.2.0 | tasks.md frontmatter 解析（`ParseTasksV2`） | 已在 go.mod；現有 `updater.go` 依賴此套件 |
| `github.com/stretchr/testify` | v1.11.1 | 所有新 package 的單元測試 | 已在 go.mod；所有 21 個測試檔案都使用此套件 |

### No New Dependencies Required

Phase 5 不需要任何新的外部依賴。所有功能使用 Go stdlib + 現有 go.mod 套件即可完成。

### Topological Sort

- **使用 Go stdlib 實作** — `internal/planchecker/` 或 `internal/executor/` 中的 Kahn's algorithm
- **不需要** 第三方圖形套件（`gonum/graph` 等）
- Kahn's algorithm 優於 DFS-based（Claude's Discretion）因為：(1) 同時自然偵測 cycle，(2) 結果 deterministic（依 task ID 排序同層），(3) 實作簡單（約 30 行）

---

## Architecture Patterns

### 現有程式碼結構（需了解的整合點）

```
internal/
├── spec/
│   ├── schema.go          ← TaskEntry 在此 (L76-90)；本 phase 擴展
│   └── updater.go         ← ParseTasksV2, WriteTasks；yaml round-trip 已建立
├── config/
│   ├── config.go          ← DefaultModelMap (L15-40)；本 phase 擴展
│   └── defaults.go        ← ProjectConfig struct；本 phase 新增 WorktreeDir/AutoMode
├── executor/
│   └── context.go         ← TaskItem (L32-38), ExecutionContext；本 phase 擴展
└── planchecker/           ← NEW package（Phase 5 建立）
    ├── checker.go
    └── checker_test.go
cmd/
└── plan.go                ← --check flag 已定義但未接線；本 phase 實作
plugin/
└── agents/
    └── mysd-plan-checker.md  ← NEW agent（Phase 5 建立）
```

### Pattern 1: Additive Struct Extension with omitempty

**What:** 在現有 struct 末尾加新欄位，全部使用 `omitempty` tag。

**Why it works:** `adrg/frontmatter` 使用 `yaml.v3` 解析，yaml.v3 的 `omitempty` 在 unmarshal 時忽略缺少的欄位（零值），在 marshal 時跳過零值欄位。現有 tasks.md 讀取時新欄位為 `nil`（slice）或 `""` （string），`WriteTasks` 寫回時不輸出這些欄位。

**Example（schema.go 修改）:**
```go
// Source: direct codebase analysis — internal/spec/schema.go L76-81
type TaskEntry struct {
    ID          int        `yaml:"id"`
    Name        string     `yaml:"name"`
    Description string     `yaml:"description,omitempty"`
    Status      ItemStatus `yaml:"status"`
    // Phase 5 additions — all omitempty for backward compatibility
    Depends     []int      `yaml:"depends,omitempty"`    // FSCHEMA-01
    Files       []string   `yaml:"files,omitempty"`      // FSCHEMA-02
    Satisfies   []string   `yaml:"satisfies,omitempty"`  // FSCHEMA-03
    Skills      []string   `yaml:"skills,omitempty"`     // FSCHEMA-04
}
```

**TaskItem（executor/context.go）同步擴展（D-12）:**
```go
type TaskItem struct {
    ID          int      `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    Status      string   `json:"status"`
    // Phase 5 additions
    Depends     []int    `json:"depends,omitempty"`
    Files       []string `json:"files,omitempty"`
    Satisfies   []string `json:"satisfies,omitempty"`
    Skills      []string `json:"skills,omitempty"`
}
```

### Pattern 2: Pure Function Package Design

**What:** `internal/planchecker/` 是一個 pure function package：接受輸入，回傳輸出，無 I/O 副作用，無外部依賴（只 import `internal/spec`）。

**Why:** Pure functions 可以在沒有任何 mock 的情況下單元測試。這是 verifier 的已建立模式（`BuildVerificationContextFromParts` 就是 pure function）。

**CoverageResult 設計（Claude's Discretion — 遵循簡潔模式）:**
```go
// Source: Claude's Discretion — follows verifier pattern
// File: internal/planchecker/checker.go
package planchecker

import "github.com/xenciscbc/mysd/internal/spec"

// CoverageResult holds the output of a plan coverage check.
type CoverageResult struct {
    TotalMust     int      `json:"total_must"`
    CoveredCount  int      `json:"covered_count"`
    UncoveredIDs  []string `json:"uncovered_ids"`   // MUST IDs with no satisfies match
    CoverageRatio float64  `json:"coverage_ratio"`  // 0.0-1.0
    Passed        bool     `json:"passed"`          // true when UncoveredIDs is empty
}

// CheckCoverage validates that every MUST item ID in mustIDs appears in at least
// one task's Satisfies field. Uses exact string matching — no AI inference.
// Pure function: no filesystem I/O, no side effects.
func CheckCoverage(tasks []spec.TaskEntry, mustIDs []string) CoverageResult
```

**Matching logic:** 對每個 mustID，遍歷 tasks，若任何 task 的 `Satisfies` 包含該 ID（exact string match），即為 covered。未 covered 的 ID 加入 `UncoveredIDs`。

### Pattern 3: ExecutionContext Extension for Plan Context

**What:** `cmd/plan.go --context-only` 的 JSON 輸出需要新增三個欄位（D-04）。

**Current structure（cmd/plan.go L65-74）:**
```go
ctx := map[string]interface{}{
    "change_name":      ws.ChangeName,
    "phase":            ws.Phase,
    "specs":            reqTexts,
    "design":           change.Design.Body,
    "model":            config.ResolveModel("planner", cfg.ModelProfile, cfg.ModelOverrides),
    "research_enabled": planResearch,
    "check_enabled":    planCheck,
    "test_generation":  cfg.TestGeneration,
}
```

**Phase 5 additions:**
```go
// Add to ctx map:
"wave_groups":    [][]interface{}{},   // WaveGroups: topological sort output (empty for Phase 5, populated in Phase 6)
"worktree_dir":   cfg.WorktreeDir,     // WorktreeDir: default ".worktrees"
"auto_mode":      cfg.AutoMode,        // AutoMode: from ProjectConfig
"coverage":       planchecker.CheckCoverage(fm.Tasks, mustIDs), // when check_enabled
```

**Note:** WaveGroups 計算邏輯（topological sort）屬於 Phase 6（FEXEC-01），Phase 5 只需輸出空 slice 或直接傳遞 tasks 不分組。`WorktreeDir` 和 `AutoMode` 是從 ProjectConfig 讀取的。

### Pattern 4: openspec/config.yaml Writer

**What:** 新函數讀寫 `openspec/config.yaml`（D-07 最小欄位集）。

**OpenSpecConfig struct:**
```go
// File: internal/spec/openspec_config.go (or extend writer.go)
type OpenSpecConfig struct {
    Project string `yaml:"project"`
    Locale  string `yaml:"locale"`   // BCP47: zh-TW, en-US, ja-JP
    SpecDir string `yaml:"spec_dir"` // convention default: "openspec/specs"
    Created string `yaml:"created"`  // RFC3339 timestamp
}

func WriteOpenSpecConfig(projectRoot string, cfg OpenSpecConfig) error
func ReadOpenSpecConfig(projectRoot string) (OpenSpecConfig, error)
```

**File location:** `{projectRoot}/openspec/config.yaml`

**Error handling（Claude's Discretion）:** 讀取不存在時回傳零值 + `nil` error（convention-over-config），這與 `config.Load()` 的現有模式一致。寫入時若 `openspec/` 目錄不存在則 `os.MkdirAll` 自動建立。

### Pattern 5: Agent Definition Format

**Agent 格式（參考 mysd-verifier.md）:**
```yaml
---
description: {one-line description}
allowed-tools:
  - Read
  - Write
  - Bash
  - Grep
---
```

**mysd-plan-checker 設計原則（FAGENT-04）:**
- 接收 `CoverageResult` JSON（來自 `mysd plan --context-only` 的 `coverage` 欄位）
- 顯示未覆蓋 IDs + 覆蓋率
- 互動式詢問：「自動補齊（invoke mysd-planner add tasks）或手動調整？」
- **禁用 Task tool**（subagent cannot spawn subagent 規則）
- agent 透過 `Write` 工具直接修改 tasks.md，或輸出建議讓 SKILL.md orchestrator 決定下一步

### Anti-Patterns to Avoid

- **不要在 planchecker 加 I/O**：`CheckCoverage` 必須是 pure function。讀取 tasks.md 和 spec 的責任在 `cmd/plan.go`，把已解析的資料傳入 checker。
- **不要在 TaskEntry 新欄位使用非 omitempty tag**：舊 tasks.md 讀取會設置空 slice，`yaml.Marshal` 會輸出 `depends: []` 這類不乾淨的 YAML。必須用 `omitempty`。
- **不要用 global viper instance**：`internal/config` 已建立 `viper.New()` 模式（instance viper for test isolation）；openspec config writer 若需要 viper 也要用 instance。
- **不要在 mysd-plan-checker agent 中使用 Task tool**：Claude Code subagent 不能 spawn 其他 subagent。agent definition 不得包含 Task tool 在 allowed-tools 中。
- **不要在 config.go 更改現有 budget 層的 haiku mapping**：D-05/D-06 只加新 roles，不改現有 6 roles 的對應（executor/spec-writer/designer 在 budget 層仍是 haiku）。

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| YAML struct serialization | 自訂 YAML formatter | `gopkg.in/yaml.v3` `yaml.Marshal/Unmarshal` | 已在 go.mod；omitempty 語意完整支援；round-trip 已驗證 |
| Frontmatter parsing | 自訂 `---` parser | `github.com/adrg/frontmatter` | 已在 go.mod；現有 `ParseTasksV2` 依賴此套件，不用換 |
| String set intersection | 自訂 map/loop | 直接 `for range` + `strings.Contains` or map lookup | 資料量小（tasks 數量 < 100），不需要複雜資料結構 |
| Config file location | 自訂 discovery | 直接 `filepath.Join(projectRoot, "openspec", "config.yaml")` | convention-over-config，路徑固定 |

**Key insight:** Phase 5 的所有問題都是 standard Go patterns。套件齊備，不需要新的 abstraction。

---

## Common Pitfalls

### Pitfall 1: omitempty on Slice Fields — nil vs empty slice

**What goes wrong:** yaml.v3 的 `omitempty` 對 slice 的行為是：`nil` slice 被 omit，`[]int{}` empty slice 也被 omit（因為長度為 0）。但如果某段程式碼用 `make([]int, 0)` 初始化而非讓欄位保持 nil，寫出的 YAML 仍然乾淨。

**Why it happens:** Go 的 `append` 呼叫會把 nil slice 轉換為非 nil slice，但長度為 0 時 omitempty 仍然 omit。

**How to avoid:** 新欄位初始化時不要 `make(..., 0)`，讓 unmarshal 時未出現的欄位保持 nil。在 `BuildContextFromParts` 中，從 `TaskEntry` 映射到 `TaskItem` 時直接複製 slice 欄位（nil 保持 nil）。

**Warning signs:** tasks.md 出現 `depends: []` 或 `satisfies: []` 這類空陣列輸出。

### Pitfall 2: plan-checker 接受錯誤的 MUST ID 格式

**What goes wrong:** spec 檔案中的 MUST item ID 有兩種來源：(1) `Requirement.ID` 欄位（planner 設定，如 `REQ-01`）；(2) `verifier.StableID(r)` 生成的 CRC32 hash ID（如 `spec.md::must-a1b2c3d4`）。如果 plan-checker 比對的 ID 格式與 planner 寫入 `satisfies` 欄位的 ID 格式不一致，永遠無法 match。

**Why it happens:** verifier 用 CRC32 hash 確保 stability；planner 可能用自訂 short ID 或直接引用 requirement text。

**How to avoid:** `cmd/plan.go` 傳入 `CheckCoverage` 的 mustIDs 必須來源於 `change.Specs`（`Requirement.ID` 欄位，非 StableID hash）。planner agent 寫入 `satisfies` 欄位時也用相同 ID 格式（即 `Requirement.ID`）。在 plan-checker agent definition 中明確說明 satisfies 欄位的 ID 格式。

**Warning signs:** `UncoveredIDs` 包含全部 MUST IDs，即使 tasks.md 有 satisfies 欄位。

### Pitfall 3: WriteTasks 序列化順序變化

**What goes wrong:** `yaml.Marshal` 對 map 輸出是按 key 排序，但對 struct 是按欄位定義順序。如果在 struct 中間插入新欄位（而非末尾），既有 tasks.md 讀後寫的輸出欄位順序會改變，造成不必要的 git diff。

**Why it happens:** yaml.v3 struct marshaling 按 struct field declaration order 輸出。

**How to avoid:** 新欄位必須加在 struct 末尾（現有欄位順序：id, name, description, status）。D-11 明確要求，且 `omitempty` 確保舊檔案寫回後不會出現新欄位。

**Warning signs:** 讀取舊 tasks.md 再 `WriteTasks`，git diff 顯示現有欄位順序改變。

### Pitfall 4: ProjectConfig 新欄位的 mapstructure tag 遺漏

**What goes wrong:** viper 使用 `mapstructure` tag 進行 `Unmarshal`；`yaml` tag 只用於直接 yaml 解析。`ProjectConfig` 兩個 tag 都需要（見 `defaults.go` 現有欄位均有兩者）。若只加 `yaml` tag 而漏了 `mapstructure`，新欄位在透過 viper 讀取設定檔時永遠是零值。

**How to avoid:** 新欄位複製現有欄位的 tag 模式：`yaml:"worktree_dir" mapstructure:"worktree_dir"`。

**Warning signs:** `TestLoad_WithConfigFile` 對新欄位的測試失敗（讀到零值）。

### Pitfall 5: openspec/config.yaml 目錄不存在

**What goes wrong:** 第一次呼叫 `WriteOpenSpecConfig` 時，`openspec/` 目錄不存在，`os.WriteFile` 直接失敗。

**How to avoid:** writer 函數在寫入前呼叫 `os.MkdirAll(filepath.Dir(configPath), 0755)`。這與 `internal/spec/updater.go WriteTasks` 的 implicit pattern 一致。

---

## Code Examples

Verified patterns from existing codebase:

### yaml.v3 omitempty Round-Trip Verification

```go
// Source: direct codebase analysis — internal/spec/updater.go L134-147
// WriteTasks already handles round-trip; new fields will serialize correctly with omitempty.
func WriteTasks(tasksPath string, fm TasksFrontmatterV2, body string) error {
    yamlBytes, err := yaml.Marshal(fm)
    // ...
}
// With omitempty, a TaskEntry{ID:1, Name:"A", Status:"pending"} (no Depends/Files/Satisfies/Skills)
// serializes to:
// - id: 1
//   name: A
//   status: pending
// (new fields absent — backward compatible)
```

### planchecker Pure Function Pattern

```go
// Source: pattern derived from internal/verifier/context.go BuildVerificationContextFromParts
// Pure function — receives pre-loaded data, returns result, no I/O
func CheckCoverage(tasks []spec.TaskEntry, mustIDs []string) CoverageResult {
    covered := make(map[string]bool)
    for _, t := range tasks {
        for _, id := range t.Satisfies {
            covered[id] = true
        }
    }

    var uncovered []string
    for _, id := range mustIDs {
        if !covered[id] {
            uncovered = append(uncovered, id)
        }
    }

    total := len(mustIDs)
    coveredCount := total - len(uncovered)
    ratio := 0.0
    if total > 0 {
        ratio = float64(coveredCount) / float64(total)
    }

    return CoverageResult{
        TotalMust:     total,
        CoveredCount:  coveredCount,
        UncoveredIDs:  uncovered,
        CoverageRatio: ratio,
        Passed:        len(uncovered) == 0,
    }
}
```

### DefaultModelMap Extension Pattern

```go
// Source: direct codebase analysis — internal/config/config.go L15-40
// Add to each tier; existing entries UNCHANGED
var DefaultModelMap = map[string]map[string]string{
    "quality": {
        // existing 6 roles unchanged...
        "researcher":      "claude-sonnet-4-5",   // D-05
        "advisor":         "claude-sonnet-4-5",   // D-05
        "proposal-writer": "claude-sonnet-4-5",   // D-05
        "plan-checker":    "claude-sonnet-4-5",   // D-05
    },
    "balanced": {
        // existing 6 roles unchanged...
        "researcher":      "claude-sonnet-4-5",   // D-05
        "advisor":         "claude-sonnet-4-5",   // D-05
        "proposal-writer": "claude-sonnet-4-5",   // D-05
        "plan-checker":    "claude-sonnet-4-5",   // D-05
    },
    "budget": {
        // existing 6 roles unchanged (executor/spec-writer etc stay haiku)...
        "researcher":      "claude-sonnet-4-5",   // D-06: budget新roles也用sonnet
        "advisor":         "claude-sonnet-4-5",   // D-06
        "proposal-writer": "claude-sonnet-4-5",   // D-06
        "plan-checker":    "claude-sonnet-4-5",   // D-06
    },
}
```

### plan.go --context-only Extension

```go
// Source: direct codebase analysis — cmd/plan.go L55-80
// Current output: change_name, phase, specs, design, model, research_enabled, check_enabled, test_generation
// Phase 5 adds: wave_groups (empty), worktree_dir, auto_mode, coverage (when check_enabled)

if planContextOnly {
    // ... existing code ...

    // Load tasks for coverage calculation
    tasksPath := filepath.Join(changeDir, "tasks.md")
    fm, _, _ := spec.ParseTasksV2(tasksPath)

    // Extract MUST IDs from specs
    var mustIDs []string
    for _, r := range change.Specs {
        if r.Keyword == spec.Must && r.ID != "" {
            mustIDs = append(mustIDs, r.ID)
        }
    }

    ctx := map[string]interface{}{
        // existing fields...
        "wave_groups":  [][]interface{}{},    // populated in Phase 6
        "worktree_dir": cfg.WorktreeDir,
        "auto_mode":    cfg.AutoMode,
    }
    if planCheck {
        ctx["coverage"] = planchecker.CheckCoverage(fm.Tasks, mustIDs)
    }
}
```

### openspec/config.yaml Writer

```go
// File: internal/spec/openspec_config.go
type OpenSpecConfig struct {
    Project string `yaml:"project"`
    Locale  string `yaml:"locale"`
    SpecDir string `yaml:"spec_dir"`
    Created string `yaml:"created"`
}

func WriteOpenSpecConfig(projectRoot string, cfg OpenSpecConfig) error {
    dir := filepath.Join(projectRoot, "openspec")
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("create openspec dir: %w", err)
    }
    data, err := yaml.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("marshal openspec config: %w", err)
    }
    return os.WriteFile(filepath.Join(dir, "config.yaml"), data, 0644)
}

func ReadOpenSpecConfig(projectRoot string) (OpenSpecConfig, error) {
    path := filepath.Join(projectRoot, "openspec", "config.yaml")
    data, err := os.ReadFile(path)
    if os.IsNotExist(err) {
        return OpenSpecConfig{}, nil  // convention-over-config: absent = zero value
    }
    if err != nil {
        return OpenSpecConfig{}, fmt.Errorf("read openspec config: %w", err)
    }
    var cfg OpenSpecConfig
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return OpenSpecConfig{}, fmt.Errorf("parse openspec config: %w", err)
    }
    return cfg, nil
}
```

### mysd-plan-checker Agent Definition Pattern

```markdown
---
description: Plan-checker agent. Receives MUST coverage result and presents uncovered items, then guides user to resolve gaps.
allowed-tools:
  - Read
  - Write
  - Edit
---
```

**Key design points:**
- 不使用 Task tool（subagent constraint）
- 不使用 Bash（不需要執行命令）
- 接收 `coverage` JSON（來自 `mysd plan --check --context-only`）
- 輸出：列出 uncovered IDs + coverage ratio，詢問自動補齊或手動調整
- 自動補齊：透過 Write/Edit 工具直接修改 tasks.md 的 `satisfies` 欄位或新增 tasks

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | `testing` (stdlib) + `github.com/stretchr/testify` v1.11.1 |
| Config file | none — Go test discovery via `_test.go` convention |
| Quick run command | `go test ./internal/spec/... ./internal/config/... ./internal/executor/... ./internal/planchecker/...` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FSCHEMA-01 | TaskEntry.Depends round-trip (YAML marshal/unmarshal) | unit | `go test ./internal/spec/... -run TestTaskEntry` | ❌ Wave 0 |
| FSCHEMA-02 | TaskEntry.Files round-trip | unit | `go test ./internal/spec/... -run TestTaskEntry` | ❌ Wave 0 |
| FSCHEMA-03 | TaskEntry.Satisfies round-trip | unit | `go test ./internal/spec/... -run TestTaskEntry` | ❌ Wave 0 |
| FSCHEMA-04 | TaskEntry.Skills round-trip | unit | `go test ./internal/spec/... -run TestTaskEntry` | ❌ Wave 0 |
| FSCHEMA-01~04 | 舊 tasks.md（無新欄位）讀取後寫回不輸出空欄位 | unit | `go test ./internal/spec/... -run TestBackwardCompat` | ❌ Wave 0 |
| FSCHEMA-05 | CheckCoverage — all covered → Passed=true | unit | `go test ./internal/planchecker/... -run TestCheckCoverage` | ❌ Wave 0 |
| FSCHEMA-05 | CheckCoverage — partial coverage → correct UncoveredIDs | unit | `go test ./internal/planchecker/... -run TestCheckCoverage` | ❌ Wave 0 |
| FSCHEMA-05 | CheckCoverage — empty mustIDs → Passed=true | unit | `go test ./internal/planchecker/... -run TestCheckCoverage_Empty` | ❌ Wave 0 |
| FSCHEMA-06 | plan-checker agent definition file exists | manual | `ls plugin/agents/mysd-plan-checker.md` | ❌ Wave 0 |
| FSCHEMA-07 | WriteOpenSpecConfig creates file with correct YAML | unit | `go test ./internal/spec/... -run TestOpenSpecConfig` | ❌ Wave 0 |
| FSCHEMA-07 | ReadOpenSpecConfig absent → zero value, no error | unit | `go test ./internal/spec/... -run TestOpenSpecConfig` | ❌ Wave 0 |
| FAGENT-04 | mysd-plan-checker.md has correct frontmatter (no Task tool) | manual | inspect `plugin/agents/mysd-plan-checker.md` | ❌ Wave 0 |
| FMODEL-01~03 | ResolveModel returns sonnet for all 4 new roles × 3 profiles | unit | `go test ./internal/config/... -run TestResolveModel_NewRoles` | ❌ Wave 0 |
| D-04 | plan --context-only JSON includes wave_groups, worktree_dir, auto_mode | unit | `go test ./cmd/... -run TestPlanContextOnly` | ❌ Wave 0 |
| D-11 | ProjectConfig defaults: WorktreeDir=".worktrees", AutoMode=false | unit | `go test ./internal/config/... -run TestDefaults` | ❌ Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./internal/spec/... ./internal/config/... ./internal/executor/... ./internal/planchecker/...`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/planchecker/checker.go` + `checker_test.go` — 新 package，Wave 0 建立
- [ ] `internal/planchecker/` directory — 需建立
- [ ] `internal/spec/openspec_config.go` + test cases in `schema_test.go` — covers FSCHEMA-07
- [ ] 擴展 `internal/spec/schema_test.go` — backward compat + new field round-trip cases
- [ ] 擴展 `internal/config/config_test.go` — new agent role model resolution tests
- [ ] 擴展 `internal/executor/context_test.go` — new TaskItem fields in JSON output

---

## Environment Availability

Step 2.6: 僅需要 Go toolchain（已確認可用，`go test ./...` 通過）。Phase 5 是純程式碼修改，無外部服務依賴。

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go toolchain | 全部 | ✓ | go 1.25.5 (from go.mod) | — |
| `gopkg.in/yaml.v3` | TaskEntry YAML round-trip | ✓ | v3.0.1 (in go.mod) | — |
| `github.com/adrg/frontmatter` | ParseTasksV2 | ✓ | v0.2.0 (in go.mod) | — |
| `github.com/stretchr/testify` | All unit tests | ✓ | v1.11.1 (in go.mod) | — |

**Missing dependencies with no fallback:** None.

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Plan-checker 作為 agent 語意推測 | Pure Go string matching on `satisfies` IDs（D-01） | Phase 5 設計決策 | Deterministic、可測試、速度快 |
| Tasks.md 無依賴欄位 | TaskEntry 加 depends/files/satisfies/skills（D-11 omitempty） | Phase 5 | Phase 6 wave grouping 可直接讀取，零 migration |
| Model profile 只有 6 roles | 擴展為 10 roles（加 researcher/advisor/proposal-writer/plan-checker） | Phase 5 | Phase 8/9 新 agent 直接可用 |

---

## Open Questions

1. **MUST ID 格式標準化**
   - What we know: `Requirement.ID` 欄位由 spec parser 設定，格式可能是 `REQ-01` 或空字串（如果 spec 沒有顯式 ID）
   - What's unclear: 若 `Requirement.ID` 為空，plan-checker 無法比對。需要確認 spec parser 是否總是設置 ID，或者是否需要 fallback 到 StableID
   - Recommendation: 在 `cmd/plan.go` 傳入 mustIDs 前，過濾掉空 ID；同時在 planner agent 的 context JSON 中只輸出有 ID 的 MUST items 給 satisfies 欄位參考

2. **WorktreeDir 和 AutoMode 的 viper default 設置**
   - What we know: `internal/config/config.go` 的 `Load()` 函數使用 `v.SetDefault()` 設置所有 defaults（L75-82）
   - What's unclear: 新欄位需要加入 `Load()` 的 `SetDefault` 呼叫，且 `Defaults()` 也要更新，但 CONTEXT.md 沒有明確指定 `WorktreeDir` 的默認值
   - Recommendation: 使用 `".worktrees"` 作為 `WorktreeDir` default（與 ARCHITECTURE.md 一致），`AutoMode` default 為 `false`

3. **mysd-plan-checker 自動補齊的確切行為**
   - What we know: D-03 說「agent 層負責渲染互動 UI 讓使用者選擇自動補齊或手動調整」；FSCHEMA-06 要求互動式詢問
   - What's unclear: 自動補齊是讓 agent 直接修改 tasks.md 的 satisfies 欄位，還是新增 task？
   - Recommendation: Phase 5 的 mysd-plan-checker.md 先實作「顯示缺口 + 詢問使用者」的互動流程，自動補齊邏輯（修改 tasks.md）可在 Phase 5 實作基本版（補 satisfies 欄位到最近匹配的 task），完整 AI 補齊邏輯留給 mysd-planner 處理

---

## Sources

### Primary (HIGH confidence)
- Direct codebase analysis — `internal/spec/schema.go` — TaskEntry struct，YAML tag 模式
- Direct codebase analysis — `internal/spec/updater.go` — ParseTasksV2, WriteTasks，YAML round-trip 已驗證
- Direct codebase analysis — `internal/executor/context.go` — TaskItem, ExecutionContext，JSON output 模式
- Direct codebase analysis — `internal/config/config.go` — DefaultModelMap，ResolveModel
- Direct codebase analysis — `internal/config/defaults.go` — ProjectConfig struct，tag 模式
- Direct codebase analysis — `cmd/plan.go` — `--check` flag 定義，`--context-only` JSON 輸出
- Direct codebase analysis — `internal/verifier/context.go` — pure function 設計模式參考
- Direct codebase analysis — `plugin/agents/mysd-verifier.md` — agent definition 格式參考
- Direct codebase analysis — `go.mod` — 確認所有依賴已存在，無需新增
- `.planning/phases/05-schema-foundation-plan-checker/05-CONTEXT.md` — 鎖定決策
- `.planning/research/ARCHITECTURE.md` — v1.1 架構研究（`internal/planchecker/` package 設計，型別定義）
- `.planning/research/SUMMARY.md` — Gap: `satisfies []string` 欄位需加入（SUMMARY.md Gaps to Address）

### Secondary (MEDIUM confidence)
- `go test ./...` 執行結果 — 確認現有測試全部通過，無 regression baseline 問題
- gopkg.in/yaml.v3 omitempty slice behavior — stdlib behavior，HIGH confidence（yaml.v3 是 canonical Go YAML library）

---

## Project Constraints (from CLAUDE.md)

| Directive | Impact on Phase 5 |
|-----------|-------------------|
| Tech stack: Go — 單一 binary | 所有新功能在 Go binary 層實作；不引入外部 service |
| 必須能讀寫 OpenSpec 格式的 spec 檔案 | openspec/config.yaml writer 使用標準 YAML（yaml.v3），符合 OpenSpec 格式 |
| Plugin 形式: Claude Code slash commands + agent definitions | mysd-plan-checker.md 要符合 Claude Code agent definition 格式 |
| Convention over configuration | openspec/config.yaml 不存在時回傳零值（不報錯）；WorktreeDir default 為 ".worktrees" |
| GSD Workflow Enforcement | 所有程式碼修改透過 GSD 工作流程進行 |

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — 所有依賴已在 go.mod，無新增依賴
- Architecture: HIGH — 直接分析現有程式碼，extension points 清晰
- Pitfalls: HIGH (structural), MEDIUM (ID format) — yaml.v3 行為為文件化事實；ID format 問題需執行時驗證

**Research date:** 2026-03-25
**Valid until:** 2026-06-25（Go stdlib patterns 穩定，30 天 TTL 為保守估計）
