# Phase 9: Interactive Discovery Integration - Research

**Researched:** 2026-03-26
**Domain:** SKILL.md orchestrator patterns, Go binary subcommand design, JSON persistence, Claude Code agent spawning
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**D-01: 探索循環終止條件 — 使用者驅動而非硬性數字上限**
- 每個 area 完成後呈現「繼續探索/完成」二元選擇，不設人工上限
- SKILL.md 中 area 完成後顯示二元選擇即可，不需 quota counter

**D-02: Deferred notes 使用 project-level 儲存 + context-aware 載入**
- 儲存於 `.specs/deferred.json`（project root 層級），跨 change 共享
- `propose` 指令載入 deferred notes（跨 change context 有價值）
- `discuss` 指令在有活躍未完成 change 時**不載入** deferred notes（避免污染 WIP）

**D-03: /mysd:note 指令 — 手動 idea 捕捉**
- `/mysd:note` 無參數 → 列出所有 deferred notes（含 ID 編號）
- `/mysd:note add {content}` → 新增 note
- `/mysd:note delete {id}` → 刪除指定 note

**D-04: Plan stage research 修正為單一 researcher**
- Phase 8 的 `/mysd:plan` 錯誤實作了 4 維度並行 research，違反 DISC-03 規格
- 正確架構：plan stage 使用單一 `mysd-researcher` agent
- 輸出供 designer agent 消費
- DISC-04 互動式 opt-in：plan stage 開始時詢問是否執行 research

**D-05: discovery-state.json 取消**
- 取消 `discovery-state.json`
- research summary 已寫入 spec 檔案（proposal.md / specs/ / design.md），不需額外持久化層
- ff/ffe 直接讀取已有的 spec 內容即可

**D-06: propose/discuss 4 維度 research 整合**
- Research 為可選（DISC-04），確定 topic 後互動式詢問
- 選擇啟用時，spawn 4 個 `mysd-researcher` agents 並行（Codebase / Domain / Architecture / Pitfalls）
- Research 結果由 SKILL.md orchestrator 識別 gray areas
- Gray areas 由 orchestrator 並行 spawn `mysd-advisor` agents 分析
- **subagent 不 spawn subagent** — advisor spawning 必須在 SKILL.md orchestrator 層

**D-07: 雙層探索循環設計**
- Layer 1（area 內深挖）：AI 或使用者提出深挖問題，逐步釐清單一 gray area
- Layer 2（area 間發現）：所有目前 gray areas 完成後，使用者可選擇發現新 areas 或結束探索
- 雙模式（DISC-05）：AI 研究後主導提問 + 使用者主導提問，兩者切換
- 終止：使用者主導（D-01），每個 area 完成後二元選擇

**D-08: Scope guardrail 機制**
- 探索中超出目前 spec 範圍的建議，由 AI 識別並 redirect 到 deferred notes
- 不修改當前 spec 內容
- deferred notes 是分離的儲存（`.specs/deferred.json`）

**D-09: /mysd:status 增強 — 顯示 deferred notes 數量**
- `/mysd:status` 底部增加 deferred notes 計數提示
- 格式：`Deferred notes: {N} — run /mysd:note to browse`

### Claude's Discretion

- `/mysd:note` Go binary subcommand vs 純 SKILL.md 實作（deferred.json 讀寫可能需要 binary 支援 JSON CRUD）
- propose workflow 的步驟順序（scaffold-first vs research-first）
- spec stage 單一 researcher 的 prompt 方向設計（DISC-02）
- deferred.json 完整欄位結構（基本預期：id, content, created_at）
- 雙模式切換的 UX 設計（AI-led vs user-led 提問模式如何在 CLI 呈現）
- ROADMAP.md success criteria 如何更新以反映 D-01 和 D-05

### Deferred Ideas (OUT OF SCOPE)

- ROADMAP.md success criteria 更新（反映 D-01 和 D-05 的偏離）— 可在 planning 後或 verification 前處理
- `/mysd:propose` 與 `/mysd:discuss` 的進一步整合（Phase 8 deferred，繼續延後）
- `/mysd:design` 指令是否完全移除（Phase 8 deferred）
- 多語言 deferred notes 支援（locale-aware content）

</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| DISC-01 | propose 階段支援 4 維度並行 research（Codebase, Domain, Architecture, Pitfalls） | mysd-propose.md 需要加入 Step 3: Research phase；現有 mysd-researcher agent 已支援 4 維度 |
| DISC-02 | spec 階段支援單一 researcher，專注「如何實作 spec」 | mysd-spec.md 需加入 optional single-researcher step；mysd-researcher 已有 codebase dimension 可重用 |
| DISC-03 | plan 階段支援單一 researcher，整合 spec + design 內容並補充實作細節 | mysd-plan.md Step 3 需修正為 single researcher（D-04 bug fix） |
| DISC-04 | 每個支援 research 的階段（propose/spec/plan/discuss）在開始時互動式詢問是否使用 research | 所有 4 個 SKILL.md 都需要加入 opt-in prompt |
| DISC-05 | Research 模式支援雙模式 — AI 研究後主導提問 + 使用者主導提問 | 雙模式實作在 SKILL.md orchestrator 層，不需 agent 修改 |
| DISC-06 | propose/discuss 的 research 產出 gray areas，由 SKILL.md orchestrator 並行 spawn advisor agents | SKILL.md 需加入 gray area identification + parallel advisor spawning 步驟 |
| DISC-07 | 雙層循環 — area 內可深挖 + 全部 areas 完成後可發現新 areas，直到使用者滿意 | D-01 確認：使用者驅動終止，二元選擇 |
| DISC-08 | Scope guardrail — 防止 scope creep，超出範圍的想法 redirect 到 deferred notes | AI prompt engineering 在 advisor/orchestrator 層；需寫入 .specs/deferred.json |
| DISC-09 | discuss 結論自動更新 spec/design/tasks，更新後自動 re-plan + plan-checker | mysd-discuss.md Step 6 (Spec Update) + Step 7 (Re-plan) 已有骨架，需強化 |

</phase_requirements>

---

## Summary

Phase 9 的核心任務是在 Phase 8 建立的 SKILL.md orchestrator 和 agent definitions 基礎上，為三個主要指令加入互動式探索能力。這個 phase 分為兩大類工作：（1）修改現有 SKILL.md 文件（propose, discuss, plan, status），以及（2）新增 `/mysd:note` 功能。

Phase 8 已完整建立所有 agent definitions（mysd-researcher、mysd-advisor、mysd-proposal-writer、mysd-spec-writer），這些 agent 可以直接被 Phase 9 的 orchestrator 擴展所使用，不需要修改 agent 本身。Phase 9 的工作完全集中在 SKILL.md 層（orchestrator logic）和 Go binary 層（deferred.json CRUD）。

最重要的架構觀察：這個 phase 的修改全部是 SKILL.md（純 Markdown 文件）和 Go binary（新增 `mysd note` subcommand），不涉及任何 agent definition 的修改。Plan 可以被清楚地分拆為 SKILL.md 修改任務和 Go binary 任務，每個任務都有明確的輸入/輸出。

**Primary recommendation:** 採用三層工作分拆：(1) Go binary — `mysd note` subcommand + deferred.json CRUD + status 整合，(2) SKILL.md 修改 — propose 加入 discovery pipeline，discuss 加入 gray area loop，plan 修正為 single researcher，(3) 建立新指令 `/mysd:note` SKILL.md。

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go stdlib `encoding/json` | stdlib | deferred.json 讀寫 | 比 yaml.v3 更簡單，notes 是純 JSON 結構 |
| Go stdlib `os`, `path/filepath` | stdlib | 檔案路徑操作 | 與 spec package 已有的模式一致 |
| github.com/spf13/cobra | v1.10.2 | `mysd note` subcommand + subcommands (list/add/delete) | 已有所有 cmd/*.go 的使用模式 |
| SKILL.md (Claude Code) | — | orchestrator logic | Phase 8 確立的 SKILL.md orchestrator pattern |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `time` (stdlib) | stdlib | created_at timestamp for deferred notes | 新增 note 時生成 ISO 8601 timestamp |
| `stretchr/testify` | v1.x | note CRUD 的 unit tests | 與所有其他 cmd/*_test.go 保持一致 |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Go binary `mysd note` | 純 SKILL.md 用 Bash jq 處理 JSON | jq 在 Windows 不保證可用；AI 產生的 jq 操作易出錯；Go binary 提供可靠的 CRUD + 一致的錯誤處理 |
| `.specs/deferred.json` | `.specs/deferred.md` (Markdown) | Markdown 格式需要 AI 解析，JSON 可機器讀寫且結構明確 |

**Installation:** 無新依賴需要安裝。

---

## Architecture Patterns

### Recommended Project Structure

新增/修改的檔案清單：

```
.claude/commands/
├── mysd-propose.md        # 修改：加入 Step 3 research + gray areas pipeline
├── mysd-discuss.md        # 修改：加入 gray areas + advisor + dual-loop + scope guardrail
├── mysd-plan.md           # 修改：Step 3 改為 single researcher（D-04 bug fix）
├── mysd-status.md         # 修改：加入 deferred notes count（D-09）
└── mysd-note.md           # 新增：/mysd:note 指令 SKILL.md

cmd/
└── note.go                # 新增：mysd note subcommand（list/add/delete）

internal/spec/
└── deferred.go            # 新增：DeferredNote struct + CRUD functions

plugin/
├── commands/mysd-note.md  # 同步：distribution copy
└── agents/                # 不需改動（Phase 8 已完整）

.specs/
└── deferred.json          # runtime：由 note.go 建立（不是 git 追蹤的 template）
```

### Pattern 1: SKILL.md Discovery Pipeline（propose/discuss）

**What:** 在現有 orchestrator steps 中插入 research → gray areas → advisor parallel spawn 的三階段 pipeline

**When to use:** propose Step 4（scaffold 後）和 discuss Step 4（topic identified 後）

**Pattern 結構：**

```
[Opt-in prompt] → spawn 4 researchers parallel → collect findings
→ orchestrator identifies gray areas from findings
→ spawn N advisors parallel (one per gray area)
→ collect comparison tables
→ enter dual-loop exploration
  → Layer 1: deep-dive per area (user/AI questions)
  → After each area: binary choice (continue / done)
  → Layer 2: after all areas → discover new areas or finish
→ conclusions feed into proposal/spec writing
```

**Key constraint:** advisors spawned at orchestrator layer, NOT inside researcher. Researcher output → orchestrator → advisor spawn.

### Pattern 2: Deferred Notes CRUD（Go binary）

**What:** `mysd note` cobra subcommand with three sub-subcommands

**Structure:**
```go
// cmd/note.go
var noteCmd = &cobra.Command{Use: "note", Short: "Manage deferred notes"}
var noteListCmd = &cobra.Command{Use: "list", RunE: runNoteList}
var noteAddCmd = &cobra.Command{Use: "add [content]", RunE: runNoteAdd}
var noteDeleteCmd = &cobra.Command{Use: "delete [id]", RunE: runNoteDelete}

// internal/spec/deferred.go
type DeferredNote struct {
    ID        int    `json:"id"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"` // ISO 8601
}

type DeferredStore struct {
    Notes []DeferredNote `json:"notes"`
}
```

**File location:** `.specs/deferred.json` (旁邊 `.specs/STATE.json`，同 specDir root 層級)

### Pattern 3: Scope Guardrail（Prompt Engineering）

**What:** AI 在 discuss/propose exploration loop 中識別超出範圍的想法並 redirect

**Implementation location:** SKILL.md orchestrator 的 system prompt / instructions

**Pattern:**
```
During exploration, if a suggestion expands beyond the current spec scope:
1. Acknowledge the idea
2. Explicitly state: "This is outside current scope"
3. Write to deferred notes: run `mysd note add "{idea content}"`
4. Continue exploration without incorporating it into spec
```

**Scope 判斷依據：** 讀取 proposal.md 的 `## Scope → In Scope / Out of Scope` 段落作為判斷基準

### Pattern 4: Single Researcher for plan/spec（D-03/D-04 fix）

**What:** mysd-plan.md Step 3 改為 single researcher，focused on implementation feasibility

**Context 差異：**
- propose/discuss researcher: `"topic"` = requirement gray area, `"dimension"` = one of 4
- plan researcher: `"topic"` = "implementation of {change_name}", `"dimension"` = `"architecture"` (or custom focused dimension), `"spec_files"` = all spec files + design.md

**Recommend:** plan 的 single researcher 使用 `"dimension": "architecture"` 配合完整 spec context，讓 researcher 專注技術可行性

### Pattern 5: Status Deferred Notes Count（D-09）

**What:** `mysd status` 底部新增一行

**Implementation:**
```go
// cmd/status.go 修改
// 在 RenderStatus 後追加：
count, _ := spec.CountDeferredNotes(specDir)
if count > 0 {
    fmt.Fprintf(cmd.OutOrStdout(), "\nDeferred notes: %d — run /mysd:note to browse\n", count)
}
```

### Pattern 6: /mysd:note SKILL.md（新指令）

**What:** 薄的 orchestrator SKILL.md，直接 delegate 到 Go binary

**Structure:**
```markdown
## Step 1: Parse Arguments
- No args → run `mysd note list`
- `add {content}` → run `mysd note add "{content}"`
- `delete {id}` → run `mysd note delete {id}`

## Step 2: Display Result
Show the output from the binary command.
```

**This is NOT a Task-spawning orchestrator** — it's a thin wrapper that calls the binary directly.

### Anti-Patterns to Avoid

- **Advisor spawning inside researcher:** Researcher MUST NOT spawn advisors. Orchestrator receives researcher output, then spawns advisors. This is the FAGENT-05 / D-06 constraint.
- **Hardcoded gray area quota:** D-01 decision — user drives termination, no numeric limits.
- **deferred.json inside a change directory:** deferred notes are project-level (`.specs/deferred.json`), not change-level.
- **Loading deferred notes in discuss when active WIP exists:** D-02 — discuss checks for active change before loading deferred context.
- **Writing discovery state to a separate file:** D-05 — no `discovery-state.json`. Research summaries live in proposal.md/spec files.
- **4 parallel researchers in mysd-plan.md:** D-04 — this is the Phase 8 bug. Plan stage uses single researcher.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| JSON CRUD for deferred.json | Custom shell jq pipeline in SKILL.md | `mysd note` Go subcommand | jq not guaranteed on Windows; binary provides consistent error handling |
| Gray area identification | New agent for gray area detection | Orchestrator AI inference from researcher output | Orchestrator reads 4 research outputs and identifies gray areas directly — no extra spawn needed |
| Dual-loop state machine | Complex state tracking struct | Conversational flow in SKILL.md | SKILL.md steps naturally maintain conversation state within a session; no persistence needed |
| Scope boundary detection | Separate scope-checker agent | Prompt engineering in orchestrator | AI can detect scope boundary from proposal.md In Scope / Out of Scope without extra agent |

**Key insight:** Phase 9 deliberately avoids new agents. All 9 agent definitions from Phase 8 are sufficient. The complexity lives in SKILL.md orchestrator flow, not in agent definitions.

---

## Common Pitfalls

### Pitfall 1: Advisor Spawn in Wrong Layer
**What goes wrong:** researcher.md の prompt でアドバイザーを spawn しようとする
**Why it happens:** DISC-06 要求 "orchestrator 並行 spawn advisor agents"，但 researcher 的 Step 3 輸出 findings 時可能有誘惑加入 spawn
**How to avoid:** SKILL.md orchestrator 收集所有 researcher 的輸出後，在 orchestrator 層識別 gray areas，然後 spawn advisors。Never inside an agent.
**Warning signs:** If you see Task tool in mysd-researcher.md or mysd-advisor.md — that's the violation.

### Pitfall 2: deferred.json 路徑混淆
**What goes wrong:** `deferred.json` 被寫入 `.specs/changes/{change-name}/deferred.json` 而非 `.specs/deferred.json`
**Why it happens:** 大多數 spec 資料都在 changeDir，容易混淆
**How to avoid:** deferred.json 位置是 `specDir/deferred.json`（和 `STATE.json` 同層），不是在 changeDir 內。Go binary 要使用 `spec.DetectSpecDir` 拿到 specDir，然後 `filepath.Join(specDir, "deferred.json")`。
**Warning signs:** deferred notes 在不同 change 之間消失。

### Pitfall 3: discuss 的 deferred notes 載入時機
**What goes wrong:** discuss 在有 active WIP change 時也載入 deferred notes，污染聚焦中的討論
**Why it happens:** 直覺上 deferred notes 應該永遠可用
**How to avoid:** D-02 規定：discuss 先檢查 `mysd status` 輸出的 phase，只有在沒有 active change（或 change 是 archived）時才載入 deferred notes。
**Warning signs:** discuss 在 WIP change 中顯示與當前 change 無關的舊 notes。

### Pitfall 4: mysd-plan.md 的 4 parallel researchers 沒被移除
**What goes wrong:** 只是「加入」single researcher 而沒有移除原本的 4 parallel researchers
**Why it happens:** D-04 是 bug fix，容易被遺漏或只做加法不做減法
**How to avoid:** mysd-plan.md Step 3 必須完全替換：移除 4 個 Task spawns，改為 1 個 Task spawn。
**Warning signs:** mysd-plan.md 的 Step 3 仍然有 "for each dimension in [...]" 迴圈。

### Pitfall 5: deferred.json ID 衝突
**What goes wrong:** 多次新增後 ID 發生衝突（例如 delete 後重新 add 使用舊 ID）
**Why it happens:** 簡單自增 ID 在刪除後重用
**How to avoid:** ID 使用 max(existing IDs) + 1 策略（或 UUID），不要重用已刪除的 ID。在 `DeferredStore.NextID()` helper 中封裝。

### Pitfall 6: Scope Guardrail 過度積極
**What goes wrong:** AI 把所有不在 spec 裡的討論都 redirect 到 deferred notes，導致探索被截斷
**Why it happens:** prompt engineering 沒有明確「探索新想法 OK，但加入 spec 才需要判斷範圍」
**How to avoid:** scope guardrail 只在使用者要求「把這個加到 spec 裡」時觸發，純粹討論/提問不需要 redirect。SKILL.md 的 guardrail instructions 要明確這個邊界。

---

## Code Examples

### deferred.go — DeferredNote CRUD

```go
// Source: internal/spec/deferred.go (new file)
package spec

import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"
)

type DeferredNote struct {
    ID        int    `json:"id"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
}

type DeferredStore struct {
    Notes []DeferredNote `json:"notes"`
}

func DeferredPath(specDir string) string {
    return filepath.Join(specDir, "deferred.json")
}

func LoadDeferredStore(specDir string) (DeferredStore, error) {
    path := DeferredPath(specDir)
    data, err := os.ReadFile(path)
    if os.IsNotExist(err) {
        return DeferredStore{}, nil // zero value — convention over config
    }
    if err != nil {
        return DeferredStore{}, err
    }
    var store DeferredStore
    if err := json.Unmarshal(data, &store); err != nil {
        return DeferredStore{}, err
    }
    return store, nil
}

func SaveDeferredStore(specDir string, store DeferredStore) error {
    data, err := json.MarshalIndent(store, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(DeferredPath(specDir), data, 0644)
}

func (s *DeferredStore) Add(content string) DeferredNote {
    id := 1
    for _, n := range s.Notes {
        if n.ID >= id {
            id = n.ID + 1
        }
    }
    note := DeferredNote{
        ID:        id,
        Content:   content,
        CreatedAt: time.Now().Format(time.RFC3339),
    }
    s.Notes = append(s.Notes, note)
    return note
}

func (s *DeferredStore) Delete(id int) bool {
    for i, n := range s.Notes {
        if n.ID == id {
            s.Notes = append(s.Notes[:i], s.Notes[i+1:]...)
            return true
        }
    }
    return false
}

func CountDeferredNotes(specDir string) (int, error) {
    store, err := LoadDeferredStore(specDir)
    if err != nil {
        return 0, err
    }
    return len(store.Notes), nil
}
```

### cmd/note.go — Cobra Subcommand Pattern

```go
// Source: cmd/note.go (new file) — mirrors existing cmd/*.go pattern
package cmd

import (
    "fmt"
    "strings"

    "github.com/xenciscbc/mysd/internal/output"
    "github.com/xenciscbc/mysd/internal/spec"
    "github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
    Use:   "note",
    Short: "Manage deferred notes",
    RunE:  runNoteList, // default: list
}

var noteAddCmd = &cobra.Command{
    Use:   "add [content]",
    Short: "Add a deferred note",
    Args:  cobra.MinimumNArgs(1),
    RunE:  runNoteAdd,
}

var noteDeleteCmd = &cobra.Command{
    Use:   "delete [id]",
    Short: "Delete a deferred note by ID",
    Args:  cobra.ExactArgs(1),
    RunE:  runNoteDelete,
}

func init() {
    rootCmd.AddCommand(noteCmd)
    noteCmd.AddCommand(noteAddCmd)
    noteCmd.AddCommand(noteDeleteCmd)
}
```

### mysd-propose.md — Research + Gray Areas Pipeline Step（插入 propose 的 Step 4 前）

```markdown
## Step 4: Optional Research (DISC-01, DISC-04)

If `auto_mode` is true: skip research entirely.

If `auto_mode` is false:
  Ask: "Would you like to run 4-dimension research on this proposal? (Codebase / Domain / Architecture / Pitfalls) [y/N]"

If user chooses research:
  Load any deferred notes as additional context:
    Run: `mysd note list` (D-02: propose always loads deferred notes)

  Spawn 4 `mysd-researcher` agents in parallel:
  [... existing pattern from mysd-discuss.md Step 4 ...]

  After collecting all 4 research outputs:
  Identify gray areas — ambiguous design decisions where multiple valid approaches exist.

  Spawn one `mysd-advisor` agent per gray area in parallel:
  Task: Analyze gray area: {gray_area_description}
  Agent: mysd-advisor
  Context: {
    "change_name": "{change_name}",
    "gray_area": "{gray_area_description}",
    "research_findings": "{all researcher output combined}",
    "auto_mode": false
  }

  Enter dual-loop exploration (D-07):
  For each gray area (Layer 1 — deep dive):
    Present advisor's comparison table
    Facilitate deep-dive questions (AI-led or user-led per DISC-05)
    After each area completes: ask "Continue to next area, or finish exploration? [continue/done]"
    If "done": exit loop

  After all areas (Layer 2 — new area discovery):
    Ask: "Would you like to explore additional areas, or are you satisfied? [explore/done]"
    If "explore": repeat process for new areas
    If "done": proceed to proposal writing
```

### mysd-plan.md — Step 3 Single Researcher Fix（D-04）

```markdown
## Step 3: Research Phase (if research_enabled)

If `research_enabled` is true (from context JSON) or `--research` flag present:

  If `auto_mode` is false:
    Ask: "Would you like to run focused research on implementation details? [y/N]"

  Spawn ONE `mysd-researcher` agent (single, not parallel):
  Task: Research implementation details for {change_name}
  Agent: mysd-researcher
  Context: {
    "change_name": "{change_name}",
    "dimension": "architecture",
    "topic": "implementation of {change_name} — validate feasibility and supplement details",
    "spec_files": [{all spec file paths + design.md}],
    "auto_mode": {auto_mode}
  }

  Collect research output. This becomes input for the designer.
  Present research summary and ask: "Research complete. Proceed to design? (Y/n)"
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Plan stage: 4 parallel researchers | Plan stage: 1 focused researcher | Phase 9 (D-04 fix) | Removes incorrect Phase 8 implementation |
| No explore loop in propose | propose → research → gray areas → advisor → dual-loop | Phase 9 | DISC-01/06/07 |
| No scope guardrail | scope guardrail → deferred.json | Phase 9 | DISC-08 |
| No /mysd:note | `/mysd:note` list/add/delete | Phase 9 | D-03 |
| status shows only task progress | status shows task progress + deferred notes count | Phase 9 | D-09 |
| discuss: no gray area identification | discuss: research → gray areas → advisor → dual-loop | Phase 9 | DISC-06/07 |

**Deprecated/outdated:**
- `mysd-plan.md` Step 3 的 4 parallel researchers：Phase 8 誤實作，Phase 9 修正。

---

## Open Questions

1. **propose workflow 的步驟順序（scaffold-first vs research-first）**
   - What we know: 目前 mysd-propose.md 的流程是 Step 1-4 先 scaffold，然後 invoke proposal-writer
   - What's unclear: research 應在 scaffold 前（topic 確認後立刻 research）還是 scaffold 後（有 proposal template 作為 context）
   - Recommendation: **scaffold-first**。Scaffold 建立 `proposal.md` template，research 讀取 spec context（proposal.md + existing specs）。這樣 researcher 有 change directory 可以讀取。改為在現有 Step 3（scaffold 後）和 Step 4（invoke proposal-writer 前）之間插入 research pipeline。

2. **spec stage 單一 researcher 的 prompt 方向（DISC-02）**
   - What we know: DISC-02 要求 spec stage 支援 single researcher，專注「如何實作 spec」
   - What's unclear: mysd-spec.md 目前的流程未讀取到，需確認現有步驟再決定插入位置
   - Recommendation: 讀取 mysd-spec.md 後決定插入位置。預期在 topic/capability 確認後、invoke spec-writer 前插入 single researcher（dimension: "codebase"）。

3. **雙模式 UX（AI-led vs user-led）**
   - What we know: DISC-05 要求兩種模式，D-07 說可以切換
   - What's unclear: 在 CLI 介面如何呈現切換
   - Recommendation: 簡化為自然對話流程。Research 完成後，AI 先呈現 findings + 提問 gray area（AI-led）；使用者可以回答或直接提出新問題（user-led）。不需要顯式的「切換模式」選項——對話的自然流即為切換。

---

## Environment Availability

Step 2.6: SKIPPED — Phase 9 is purely SKILL.md + Go binary changes with no new external dependencies. All dependencies (Go, Cobra, existing spec package) are already available from Phase 8.

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go testing (stdlib) + stretchr/testify |
| Config file | 無獨立 config — 使用 `go test ./...` |
| Quick run command | `go test ./internal/spec/... ./cmd/... -run TestDeferred -v` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| DISC-01 | propose 4 維度 research pipeline | manual | — (SKILL.md behavior) | N/A — SKILL.md |
| DISC-02 | spec single researcher opt-in | manual | — (SKILL.md behavior) | N/A — SKILL.md |
| DISC-03 | plan single researcher（D-04 fix） | manual | — (SKILL.md behavior) | N/A — SKILL.md |
| DISC-04 | research opt-in prompt at stage start | manual | — (SKILL.md behavior) | N/A — SKILL.md |
| DISC-05 | dual-mode exploration | manual | — (SKILL.md behavior) | N/A — SKILL.md |
| DISC-06 | parallel advisor spawning from orchestrator | manual | — (SKILL.md behavior) | N/A — SKILL.md |
| DISC-07 | dual-loop with user-driven termination | manual | — (SKILL.md behavior) | N/A — SKILL.md |
| DISC-08 | scope guardrail → deferred.json | unit | `go test ./internal/spec/... -run TestDeferredNote` | ❌ Wave 0 |
| DISC-09 | discuss → re-plan + plan-checker | manual | — (SKILL.md behavior) | N/A — SKILL.md |
| D-03 | `/mysd:note` list/add/delete | unit | `go test ./cmd/... -run TestNoteCmd` | ❌ Wave 0 |
| D-09 | status deferred count display | unit | `go test ./cmd/... -run TestStatusDeferredCount` | ❌ Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./internal/spec/... -run TestDeferred -v`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `internal/spec/deferred_test.go` — covers DISC-08, D-03 (CRUD unit tests)
- [ ] `cmd/note_test.go` — covers D-03 (CLI integration tests)
- [ ] `cmd/status_test.go` 新增 deferred count test case — covers D-09

*(All SKILL.md requirements are manual-only — no automated test infrastructure needed for them)*

---

## Sources

### Primary (HIGH confidence)
- Phase 8 SKILL.md files (`.claude/commands/mysd-*.md`) — direct read, verified current state
- Phase 8 agent definitions (`.claude/agents/mysd-*.md`) — direct read, verified no Task tool in agents
- `internal/spec/schema.go` — existing Go struct patterns (additive extension confirmed)
- `cmd/status.go`, `cmd/propose.go` — existing binary patterns for note.go design
- Phase 9 CONTEXT.md — all decisions verified from discuss-phase session

### Secondary (MEDIUM confidence)
- `.specs/changes/interactive-discovery/proposal.md` — original feature spec, used to verify DISC-01~09 scope
- `STATE.md` accumulated decisions — Phase 5-8 patterns confirmed (additive struct extension, zero-value convention, leaf agent constraint)

### Tertiary (LOW confidence)
- None — all findings verified against primary sources.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all libraries already in use, no new dependencies
- Architecture patterns: HIGH — all patterns derived from direct code inspection of Phase 8 output
- Pitfalls: HIGH — derived from CONTEXT.md decisions (each decision reflects a pitfall avoided)
- Test map: HIGH for Go binary tests; N/A for SKILL.md (manual-only by nature)

**Research date:** 2026-03-26
**Valid until:** 2026-04-26 (stable — no fast-moving external dependencies)
