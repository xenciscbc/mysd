# Phase 2: Execution Engine - Research

**Researched:** 2026-03-23
**Domain:** Go CLI execution engine, Claude Code SKILL.md/agent definition authoring, prompt-based alignment gate, task progress tracking
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Plugin skill 反向調用 — `/mysd:execute` 等指令由 SKILL.md 作為入口觸發，Go binary 負責 spec 解析和狀態管理，Claude Code agent 負責實際 AI 執行
- **D-02:** Agent definition 編排 — 每個 workflow 階段有專屬 agent .md 檔案（仿 GSD 模式），如 mysd-spec-writer.md、mysd-designer.md、mysd-planner.md、mysd-executor.md、mysd-verifier.md，SKILL.md 編排調用
- **D-03:** Wave mode 使用 Claude Code 原生 Agent tool 生成多個平行 subagent，每個處理一個 task（GSD 也用此機制）
- **D-04:** Profile-based 模型管理 — 仿 GSD 的 resolve-model 機制，在 mysd.yaml 中配置 model profile（quality/balanced/budget），每個 agent 類型根據 profile 映射到具體模型，可在 mysd.yaml 覆蓋預設
- **D-05:** Prompt 注入 + 確認指令 — agent definition 中強制包含 spec 內容，要求 AI 在回應中明確列出它理解的 MUST/SHOULD/MAY 項目，然後才能開始寫 code
- **D-06:** 結構化摘要輸出 — AI 必須輸出 alignment summary：列出所有 MUST 項目、它理解的執行策略、任何疑問。此 summary 可被後續 verify 階段引用
- **D-07:** Alignment summary 寫入 `.specs/changes/{name}/alignment.md`，跟著 spec artifacts 走，版控可追蹤
- **D-08:** spec / design / plan 三個指令都有完整的 AI 互動流程，各由專屬 agent 執行（類似 propose 的互動體驗）
- **D-09:** `mysd ff`（fast-forward）從 propose 一氣推進到 plan 完成，跳過互動確認（用預設值），結果是完整的 spec artifacts + plan，使用者可直接 execute
- **D-09b:** `mysd ffe`（fast-forward execute）從 propose 一氣推進到實作完成（propose → spec → design → plan → execute），跳過所有互動確認
- **D-10:** `mysd capture` 分析當前 Claude Code 對話中討論過的變更，提取關鍵需求，然後帶預填內容自動進入 propose 流程
- **D-11:** `mysd status` 顯示綜合儀表板：當前 change name、workflow phase、任務完成率（X/Y tasks done）、MUST/SHOULD/MAY 達成狀態、上次執行時間
- **D-12:** Plan 階段可選管線 — 預設只有 plan（快速），可用 flag 啟用 research 和 plan-check
- **D-13:** Agent 回報 + binary 更新 — SKILL.md / agent definition 要求 AI 在開始和完成每個 task 時呼叫 `mysd task-update {id} {status}`，Go binary 更新 tasks.md frontmatter 和 STATE.json
- **D-14:** Task level 中斷恢復 — 從最後一個完成的 task 之後恢復，已完成的不重做
- **D-15:** TDD mode: test-first 指令注入 — 啟用 TDD 時，agent definition 增加指令要求先寫測試再寫實作
- **D-16:** Atomic commits 粒度為每個 task 一個 commit（--atomic-commits flag 啟用時）

### Claude's Discretion

- Agent definition 的具體 prompt 措辭和結構
- Alignment summary 的具體 markdown 模板
- `mysd task-update` 的 CLI 介面設計（flag 名稱、輸出格式）
- Status 儀表板的 lipgloss 配色和排版
- ff 指令的預設值選擇策略
- Model profile 的具體模型映射表（哪個 profile 對應哪個模型）

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| EXEC-01 | User can execute spec tasks via `/mysd:execute` with pre-execution alignment gate | SKILL.md 入口設計 + Go binary `execute` subcommand + alignment gate prompt injection 模式 |
| EXEC-02 | Default execution mode is single-agent sequential | Go binary 讀取 tasks.md，輸出 SKILL.md 使用的 context block；single agent 依序處理每個 task |
| EXEC-03 | User can opt into multi-agent wave execution mode with configurable agent count | Claude Code 原生 `Task` tool 呼叫多個 subagent；`--mode=wave --agents=N` flag 傳遞到 SKILL.md |
| EXEC-04 | Atomic git commits per task is available as opt-in | agent definition 在每個 task 完成後執行 `git add -A && git commit -m "..."` |
| EXEC-05 | Execution engine tracks progress and can resume from interruption point | `mysd task-update` 寫入 tasks.md frontmatter；`mysd execute` 啟動時讀取 status 跳過已完成的 tasks |
| WCMD-01 | `/mysd:propose` — create new spec from user description | Phase 1 已實作 Go binary scaffold；Phase 2 補全 SKILL.md 入口和 AI 互動 agent |
| WCMD-02 | `/mysd:spec` — define detailed requirements | SKILL.md + mysd-spec-writer agent definition；Go binary `spec` subcommand 更新 specs/ 目錄 |
| WCMD-03 | `/mysd:design` — capture technical decisions | SKILL.md + mysd-designer agent definition；Go binary `design` subcommand 更新 design.md |
| WCMD-04 | `/mysd:plan` — break design into task list | SKILL.md + mysd-planner agent definition；Go binary `plan` subcommand 更新 tasks.md |
| WCMD-05 | `/mysd:execute` — run tasks with alignment and progress tracking | SKILL.md + mysd-executor agent definition；Go binary `execute` subcommand |
| WCMD-08 | `/mysd:status` — show current spec state and progress | Go binary `status` subcommand；lipgloss 儀表板 |
| WCMD-10 | `/mysd:ff` — fast-forward propose→plan | SKILL.md 鏈式呼叫；Go binary `ff` subcommand |
| WCMD-11 | `/mysd:init` — 初始化專案設定 | Go binary `init` subcommand（Phase 1 已有 stub）；互動式 viper 設定寫入 |
| WCMD-13 | `/mysd:capture` — 從對話中提取變更 | SKILL.md 分析對話歷史；帶預填內容進入 propose 流程 |
| WCMD-14 | `/mysd:ffe` — fast-forward propose→execute | SKILL.md 鏈式呼叫；Go binary `ffe` subcommand |
| TEST-01 | User can opt into TDD mode | agent definition 的 test-first prompt injection；`--tdd` flag 已在 root.go 定義 |
| TEST-02 | 執行後自動產生對應測試程式碼 | agent definition 在 implementation 完成後的 test generation step |
| TEST-03 | TDD 模式可在設定檔中設為預設 | `internal/config` 的 `TDD bool` 已存在；SKILL.md 讀取 binary context |
</phase_requirements>

---

## Summary

Phase 2 的核心挑戰是：**讓 Go binary 和 Claude Code plugin layer（SKILL.md + agent definitions）正確分工協作**。Go binary 負責所有狀態讀寫和 spec 解析；Claude Code plugin 負責 AI 互動。這兩層之間的溝通協議（binary 輸出 context、AI 呼叫 binary 更新狀態）是整個 phase 的關鍵設計點。

Alignment gate 的實作方式是 **prompt engineering**：在 mysd-executor agent definition 中，強制要求 AI 先輸出 alignment summary（列出所有 MUST/SHOULD/MAY 項目的理解），才能開始寫 code。這個 gate 完全依賴 agent definition 的 prompt 結構 — Go binary 不需要「知道」AI 是否對齊，只需提供 spec 內容給 agent。

任務進度追蹤採用雙寫策略：AI 呼叫 `mysd task-update <id> <status>`（CLI subcommand），Go binary 同步更新 tasks.md frontmatter 的 task 狀態和 STATE.json。中斷恢復時，`mysd execute` 重新讀取 tasks.md，跳過已標記為 `done` 的 tasks。

**Primary recommendation:** 分三個實作核心——(1) Go binary subcommands（task-update、status、spec/design/plan/execute/ff/ffe/init/capture 的 runX 實作）；(2) Claude Code plugin 層（SKILL.md x 10 + agent definition x 5）；(3) alignment gate 和進度追蹤的整合測試。

---

## Standard Stack

### Core（繼承自 Phase 1，無變更）

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.25.5 (env) | Primary language | 已鎖定；零 runtime 依賴 |
| github.com/spf13/cobra | v1.10.2 | CLI subcommands | Phase 1 已確認；`task-update`、`ff`、`ffe`、`capture` 為新 subcommands |
| gopkg.in/yaml.v3 | v3.0.1 | tasks.md frontmatter 讀寫 | Phase 1 已確認；task status 更新需要 YAML round-trip |
| github.com/adrg/frontmatter | v0.2.0 | tasks.md / spec frontmatter parsing | Phase 1 已確認 |
| github.com/spf13/viper | v1.21.0 | mysd.yaml model profile 讀取 | Phase 1 已確認；Phase 2 擴展 model_profile 欄位 |
| github.com/charmbracelet/lipgloss | v1.1.0 | status 儀表板樣式 | Phase 1 已確認；Phase 2 首次用於多行儀表板 |
| github.com/stretchr/testify | v1.11.1 | 測試 assertions | Phase 1 已確認 |

### Phase 2 新增（無需 go.mod 新依賴）

Phase 2 的所有新功能可用 Phase 1 已有依賴實現：

- `text/template` (stdlib) — alignment.md 模板渲染（tasks.md round-trip）
- `os/exec` (stdlib) — 如需在測試中 mock binary 呼叫（非必要）
- `bufio` + `strings` (stdlib) — tasks.md task status 行內更新

**Version verification:** 已從 go.mod 確認所有版本，無需 npm view 步驟。所有依賴已在 Phase 1 安裝。

---

## Architecture Patterns

### Plugin 目錄結構（Phase 2 新建）

```
.claude/
├── mysd.yaml                    # 專案設定（Phase 1 已確認路徑）
└── commands/
    ├── mysd-propose.md           # /mysd:propose SKILL.md
    ├── mysd-spec.md              # /mysd:spec SKILL.md
    ├── mysd-design.md            # /mysd:design SKILL.md
    ├── mysd-plan.md              # /mysd:plan SKILL.md
    ├── mysd-execute.md           # /mysd:execute SKILL.md
    ├── mysd-status.md            # /mysd:status SKILL.md
    ├── mysd-ff.md                # /mysd:ff SKILL.md
    ├── mysd-ffe.md               # /mysd:ffe SKILL.md
    ├── mysd-init.md              # /mysd:init SKILL.md
    └── mysd-capture.md           # /mysd:capture SKILL.md

.claude/agents/
    ├── mysd-spec-writer.md       # agent definition: spec 撰寫
    ├── mysd-designer.md          # agent definition: design 撰寫
    ├── mysd-planner.md           # agent definition: task 分解
    ├── mysd-executor.md          # agent definition: 實作執行（含 alignment gate）
    └── mysd-fast-forward.md      # agent definition: ff/ffe 無互動全流程

.specs/changes/{name}/
    └── alignment.md             # alignment gate 輸出（D-07）
```

### Go Binary 新增 subcommands

```
cmd/
├── task_update.go              # mysd task-update <id> <status>
├── ff.go                       # mysd ff [name]
├── ffe.go                      # mysd ffe [name]
├── capture.go                  # mysd capture
└── init_cmd.go                 # Phase 1 已有 stub，Phase 2 補完互動邏輯

internal/
├── spec/
│   └── updater.go              # Task status 更新（tasks.md round-trip 寫入）
└── executor/
    ├── context.go              # 為 SKILL.md 提供 spec context JSON
    ├── alignment.go            # alignment.md 生成和驗證
    └── progress.go             # 從 tasks.md 計算完成率
```

### Pattern 1: SKILL.md 入口 + Binary Context 協議

**What:** SKILL.md 是 Claude Code 的入口 slash command，它呼叫 Go binary 取得 spec context，再呼叫對應的 agent definition。

**When to use:** 所有需要 AI 互動的 workflow 指令（spec、design、plan、execute）。

**Example — mysd-execute.md SKILL.md 骨架：**
```markdown
---
model: claude-sonnet-4-5
description: Execute spec tasks with alignment gate
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:execute

## Step 1: Load spec context

```bash
mysd execute --context-only
```

This outputs JSON with: change_name, must_items[], should_items[], may_items[], tasks[], pending_tasks[], tdd_mode, atomic_commits, execution_mode, agent_count.

## Step 2: Spawn executor agent

Use the Task tool to invoke mysd-executor agent with the context JSON.
```

**Key insight:** `mysd execute --context-only` flag 讓 binary 只輸出 context JSON，不觸發實際執行。SKILL.md 消費這個 JSON，然後呼叫 agent。

### Pattern 2: Alignment Gate via Prompt Engineering

**What:** mysd-executor agent definition 的 prompt 中強制要求 AI 輸出 alignment summary 作為第一步，然後才能寫 code。

**When to use:** 所有觸發實際 code 寫作的 agent（executor）。

**Example — mysd-executor.md 的 alignment gate section：**
```markdown
## MANDATORY: Alignment Gate

Before writing any code, you MUST output an alignment summary in this exact format:

```markdown
## Alignment Summary

### MUST items I understand:
- [MUST-01]: [exact text] — My interpretation: [...]
- [MUST-02]: [exact text] — My interpretation: [...]

### Execution strategy:
[How you will implement each MUST item]

### Open questions (if any):
[Any ambiguity that needs clarification]
```

After outputting the alignment summary, write it to `.specs/changes/{change_name}/alignment.md` using the Write tool, then proceed with implementation.

DO NOT write any code before completing and saving the alignment summary.
```

**Key insight:** Alignment gate 是純 prompt engineering，不需要 Go binary 做任何 gate 邏輯。Binary 只需確保 spec content 被完整注入到 agent prompt 中。

### Pattern 3: `mysd task-update` Round-Trip

**What:** AI 呼叫 `mysd task-update <id> <status>` → binary 讀取 tasks.md → 更新對應 task 的 status → 寫回 tasks.md → 更新 STATE.json 的統計。

**Implementation challenge:** tasks.md 的 task 狀態在 frontmatter 中（TasksFrontmatter.Total/Completed），但個別 task 的 status 在 markdown body 中（以 markdown list 或 table 形式）。需要決定 tasks.md 的完整 schema。

**Recommended tasks.md schema（Claude's discretion）：**
```yaml
---
spec-version: "1"
total: 5
completed: 2
---

## Tasks

| ID | Name | Status |
|----|------|--------|
| 1 | Implement auth | done |
| 2 | Add tests | in_progress |
| 3 | Update docs | pending |
```

Binary 讀取整個 tasks.md，解析 markdown table（或 YAML block），更新指定 ID 的 status，重新序列化整個文件。使用 `bufio` + `strings` 做行掃描，比引入新的 markdown AST 庫更簡單。

**tasks.md 格式選擇：使用 YAML block（不用 markdown table）**

理由：markdown table 的行掃描脆弱、容易因格式化差異（空格）破壞。建議改用 YAML block 在 frontmatter 中內嵌 task list，或使用 YAML-only 的 tasks 區塊：

```yaml
---
spec-version: "1"
tasks:
  - id: 1
    name: "Implement auth"
    description: "..."
    status: done
  - id: 2
    name: "Add tests"
    description: "..."
    status: in_progress
---
```

這樣 `mysd task-update` 可以用 `yaml.v3` 做完整 round-trip，避免 markdown 解析的脆弱性。`internal/spec.Task` struct 已有 ID/Name/Description/Status，完全相容。

### Pattern 4: Wave Mode via Claude Code Task Tool

**What:** SKILL.md 在 `--mode=wave` 時，使用 Claude Code 原生 `Task` tool 生成多個 subagent，每個 subagent 處理一個 task。

**Implementation in SKILL.md：**
```markdown
## Wave Mode (when execution_mode == "wave")

Spawn {agent_count} parallel subagents using the Task tool.
Each subagent receives:
- task_id: the specific task ID to implement
- spec context: the full alignment summary and spec content
- Instructions: implement only this task, then run `mysd task-update {id} done`

Wait for all subagents to complete before proceeding.
```

**Key insight:** Wave mode 完全在 SKILL.md 層實作，Go binary 不需要知道 wave 的存在。Binary 只提供 context 和接受 `task-update` 呼叫。

### Pattern 5: Status Dashboard

**What:** `mysd status` 呼叫 Go binary，binary 讀取 STATE.json + tasks.md + spec files，用 lipgloss 渲染多行儀表板。

**Example output structure（Claude's discretion）：**
```
=== my-feature ===
Phase:    planned → execute (ready)
Progress: ████████░░ 4/5 tasks done

MUST items:  ✓ 3 done, ○ 1 pending
SHOULD items: ✓ 2 done
MAY items:   — (not tracked until verify)

Last run: 2026-03-23 14:30
```

**lipgloss pattern for multi-column layout：**
```go
// Source: charmbracelet/lipgloss v1.1.0 README
labelStyle := lipgloss.NewStyle().Bold(true).Width(12)
valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

row := lipgloss.JoinHorizontal(
    lipgloss.Top,
    labelStyle.Render("Phase:"),
    valueStyle.Render(phase),
)
```

### Pattern 6: ff / ffe SKILL.md 鏈式呼叫

**What:** `mysd ff` 和 `mysd ffe` 不是「新功能」，而是把 propose → spec → design → plan（→ execute）這條 workflow 鏈無互動地連接起來。

**Implementation:** SKILL.md 中順序呼叫各 agent（mysd-spec-writer → mysd-designer → mysd-planner → [mysd-executor]），每個 agent 完成後 binary 自動 transition state，不需要使用者確認。

**Flag 傳遞機制:** Go binary 的 `ff` subcommand 先執行 `mysd propose {name}`，然後輸出下一步所需的 context。SKILL.md 消費這個輸出後繼續呼叫下一個 agent。

### Anti-Patterns to Avoid

- **在 Go binary 中實作 alignment 邏輯**：Alignment gate 是 prompt engineering 的領域，binary 不應該嘗試「驗證」AI 是否輸出了 alignment summary。Binary 只做狀態管理。
- **讓 AI 直接寫入 tasks.md 的 YAML**：AI 應呼叫 `mysd task-update` binary，而不是直接 Edit tasks.md。直接 Edit 繞過了 binary 的狀態同步邏輯。
- **在 tasks.md 使用 markdown table 追蹤 task status**：Markdown table 解析脆弱。使用 YAML frontmatter 中的 task list（yaml.v3 round-trip 更安全）。
- **把 wave mode 邏輯放進 Go binary**：Wave mode 是 Claude Code plugin 層的責任（Task tool 呼叫），binary 不應依賴 Claude Code API。
- **SKILL.md 文件路徑放在 `.claude/commands/` 之外**：Claude Code 只識別 `.claude/commands/*.md` 作為 slash commands。

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| tasks.md YAML round-trip | Custom line scanner | `yaml.v3` unmarshal/marshal | YAML 有 anchors、多行字符串等邊界情況；yaml.v3 處理正確 |
| Terminal styling | ANSI escape codes | `lipgloss` (已在 go.mod) | TTY 偵測、顏色 profile、Windows 相容性已處理 |
| Subagent spawning | 自訂 goroutine pool | Claude Code `Task` tool（在 SKILL.md 中） | Wave mode 需要 Claude Code 的 Agent tool，不是 Go goroutines |
| Model selection | Hard-coded model names | viper config + mysd.yaml `model_profile` | 支援使用者配置，仿 GSD resolve-model 機制 |
| Markdown frontmatter write | 手動字串拼接 | `adrg/frontmatter`（已在 go.mod） | 處理 YAML 特殊字符的 marshaling |

**Key insight:** Phase 2 的「新依賴」幾乎是零 — 所有技術問題都由 Phase 1 已有的 stack 解決。主要工作量在 prompt engineering（agent definitions）和 Go binary 業務邏輯（tasks.md updater、status dashboard）。

---

## Common Pitfalls

### Pitfall 1: tasks.md Schema 不一致

**What goes wrong:** Phase 1 的 `TasksFrontmatter` struct 只有 `Total/Completed` 計數，沒有 per-task status 欄位。Phase 2 需要 per-task status，但如果 schema 設計不當（如在 body 用 markdown table），`mysd task-update` 的實作會非常脆弱。

**Why it happens:** Phase 1 設計的 frontmatter 是為了統計，Phase 2 需要讓 binary 更新個別 task status。

**How to avoid:** 把 tasks 列表移入 frontmatter（YAML block）。這需要更新 `TasksFrontmatter` struct，加入 `Tasks []TaskFrontmatterEntry`。現有 `spec.Task` struct 可直接重用。

**Warning signs:** 如果看到 `task-update` 的實作在用 `strings.Replace` 或 `bufio` 掃描 markdown table 行，這是錯誤的方向。

### Pitfall 2: SKILL.md 中 binary 路徑硬編碼

**What goes wrong:** SKILL.md 中呼叫 `mysd execute --context-only`，但不同機器的 binary 可能在不同路徑（PATH 中、專案目錄、~/bin 等）。

**Why it happens:** 開發期間 binary 通常在 `go run ./...` 或 `go build` 的輸出路徑。

**How to avoid:** SKILL.md 使用 `mysd`（依賴 PATH）。在文件中說明安裝前提（`go install` 或 `go build -o ~/bin/mysd`）。Phase 4 的 distribution 會處理安裝問題。

**Warning signs:** SKILL.md 中出現絕對路徑如 `/home/user/go/bin/mysd`。

### Pitfall 3: Alignment Gate 可被 AI 繞過

**What goes wrong:** AI agent 看到 alignment gate 的 prompt 要求，但仍然在輸出 alignment summary 之前就開始寫 code（例如因為 system prompt 衝突或 instruction following 問題）。

**Why it happens:** Prompt engineering 不是 100% 可靠。AI 可能在某些邊緣情況下不嚴格遵循指令順序。

**How to avoid:**
1. 在 agent definition 中把 alignment gate 放在最頂層，用大寫和 `MANDATORY` 標記
2. 讓 alignment gate 的輸出具體（指定格式、要求寫入檔案），增加「可驗證性」
3. 接受這個 gate 是「soft gate」（信任 AI），而不是 binary 層面的 hard gate

**Warning signs:** alignment.md 不存在但 task 被標記為 in_progress。

### Pitfall 4: tasks.md 的 YAML round-trip 破壞 body 內容

**What goes wrong:** `yaml.v3` unmarshal + marshal 會重新格式化 YAML，可能改變縮排、引號風格等，破壞 tasks.md 的人類可讀性。

**Why it happens:** `yaml.v3` 的 marshal 用標準化格式輸出，不保留原始格式。

**How to avoid:** 用 `adrg/frontmatter` 分離 frontmatter 和 body，只對 frontmatter 部分做 yaml.v3 round-trip，body 保持不變。寫回時重新組合 `---\n{yaml}\n---\n{body}`。

**Warning signs:** `mysd task-update` 執行後 tasks.md 的非 frontmatter 部分被清空或改變。

### Pitfall 5: ff 指令的 state transition 不正確

**What goes wrong:** `mysd ff` 依序執行 propose → spec → design → plan，但每步都需要正確的 state transition。如果 state 機器不允許跳步（例如直接從 PhaseProposed 到 PhasePlanned），ff 會失敗。

**Why it happens:** `internal/state/transitions.go` 的 `ValidTransitions` 要求逐步轉換，不允許跳躍。

**How to avoid:** `mysd ff` 必須依序呼叫 Transition()：PhaseNone → PhaseProposed → PhaseSpecced → PhaseDesigned → PhasePlanned，每步都更新 STATE.json。這已由現有 `state.Transition()` 支援，只需正確順序呼叫。

**Warning signs:** ff 執行時出現 `ErrInvalidTransition` 錯誤。

### Pitfall 6: `mysd capture` 的上下文訪問限制

**What goes wrong:** `mysd capture` 要「分析當前對話中討論過的變更」，但 Go binary 無法訪問 Claude Code 的對話歷史。

**Why it happens:** Binary 是一個獨立的 CLI 程式，它不知道 Claude Code 對話的內容。

**How to avoid:** `mysd capture` 的「分析對話」邏輯必須在 SKILL.md 層（由 Claude Code 讀取對話上下文），而不是在 Go binary 層。SKILL.md 負責分析對話並生成 pre-filled proposal 內容，然後呼叫 `mysd propose {name}` 並把內容寫入 proposal.md。

**Warning signs:** 試圖讓 binary 讀取 `.claude/` 目錄下的對話歷史文件。

---

## Code Examples

Verified patterns from existing codebase and Go stdlib:

### tasks.md 更新（YAML round-trip 模式）

```go
// Source: internal/spec/updater.go (to be created)
// Pattern: adrg/frontmatter split + yaml.v3 round-trip + body preserved

type TaskEntry struct {
    ID          int        `yaml:"id"`
    Name        string     `yaml:"name"`
    Description string     `yaml:"description"`
    Status      ItemStatus `yaml:"status"`
}

type TasksFrontmatterV2 struct {
    SpecVersion string      `yaml:"spec-version"`
    Total       int         `yaml:"total"`
    Completed   int         `yaml:"completed"`
    Tasks       []TaskEntry `yaml:"tasks"`
}

func UpdateTaskStatus(tasksPath string, taskID int, newStatus ItemStatus) error {
    f, err := os.Open(tasksPath)
    if err != nil {
        return err
    }
    var fm TasksFrontmatterV2
    body, err := frontmatter.Parse(f, &fm)
    f.Close()
    if err != nil {
        return err
    }

    // Update task status
    updated := false
    for i, t := range fm.Tasks {
        if t.ID == taskID {
            fm.Tasks[i].Status = newStatus
            updated = true
        }
    }
    if !updated {
        return fmt.Errorf("task %d not found", taskID)
    }

    // Recount completed
    fm.Completed = 0
    for _, t := range fm.Tasks {
        if t.Status == StatusDone {
            fm.Completed++
        }
    }

    // Serialize back
    fmBytes, err := yaml.Marshal(fm)
    if err != nil {
        return err
    }
    content := fmt.Sprintf("---\n%s---\n%s", fmBytes, body)
    return os.WriteFile(tasksPath, []byte(content), 0644)
}
```

### lipgloss Status Dashboard

```go
// Source: charmbracelet/lipgloss v1.1.0 official README pattern
// internal/executor/progress.go (to be created)

func RenderStatusDashboard(p *output.Printer, ws state.WorkflowState, tasks []TaskEntry) {
    // Only render in TTY mode — Printer.isTTY check is internal
    // Use p.Printf for raw output

    done := 0
    for _, t := range tasks {
        if t.Status == spec.StatusDone {
            done++
        }
    }

    labelStyle := lipgloss.NewStyle().Bold(true).Width(14)
    valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
    headerStyle := lipgloss.NewStyle().Bold(true).Underline(true)

    fmt.Println(headerStyle.Render("=== " + ws.ChangeName + " ==="))
    fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
        labelStyle.Render("Phase:"),
        valueStyle.Render(string(ws.Phase)),
    ))
    fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
        labelStyle.Render("Progress:"),
        valueStyle.Render(fmt.Sprintf("%d/%d tasks done", done, len(tasks))),
    ))
}
```

### SKILL.md Claude Code slash command format

```markdown
---
model: claude-sonnet-4-5
description: Execute spec tasks with pre-execution alignment gate
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:execute

## Step 1: Load execution context

```bash
mysd execute --context-only
```

Parse the JSON output. It contains:
- `change_name`: current change name
- `must_items`: array of MUST requirements
- `tasks`: full task list
- `pending_tasks`: tasks not yet done
- `tdd_mode`: boolean
- `execution_mode`: "single" or "wave"

## Step 2: Execute

[agent invocation using Task tool or inline]
```

### agent definition format (.claude/agents/mysd-executor.md)

```markdown
---
model: claude-sonnet-4-5
description: Implements spec tasks with mandatory alignment gate
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
---

You are the mysd executor agent. You implement spec tasks for the current change.

## MANDATORY: Alignment Gate

[alignment gate prompt section — see Pattern 2 above]

## Task Execution

For each pending task:
1. Run `mysd task-update {id} in_progress`
2. Implement the task
3. Run `mysd task-update {id} done`
[TDD injection section — conditional on tdd_mode]
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Hard-coded spec validation in binary | Prompt-based alignment gate in agent definition | GSD-era pattern (2024+) | Binary stays thin, AI does alignment |
| Global viper (viper.Get) | Instance viper (viper.New()) | Phase 1 decision | Full test isolation |
| Markdown table for task tracking | YAML frontmatter task list | Phase 2 decision | Safe round-trip via yaml.v3 |
| Single monolithic SKILL.md | Per-command SKILL.md + per-role agent .md | Claude Code plugin best practice | Separation of concerns, reuse |

**Deprecated/outdated:**

- `TasksFrontmatter.Total/Completed` only (無 per-task list)：Phase 2 必須擴展為含 `Tasks []TaskEntry`，但保持向後相容（Completed 計數仍存在）

---

## Open Questions

1. **`mysd execute --context-only` 的 output format**
   - What we know: SKILL.md 需要從 binary 取得 spec context（change name, MUST items, tasks, tdd mode 等）
   - What's unclear: 是輸出純 JSON（SKILL.md 用 bash 解析），還是輸出 markdown block（直接嵌入 agent prompt）？
   - Recommendation: 輸出 JSON，SKILL.md 在 bash 中用 jq 或直接把 JSON 嵌入 agent prompt 的 context block。JSON 更結構化，易於測試。

2. **mysd-executor agent definition 的調用方式**
   - What we know: SKILL.md 可以用 `Task` tool 呼叫 agent，或直接 inline 執行 agent 邏輯
   - What's unclear: 是否應該有獨立的 `.claude/agents/mysd-executor.md` 文件，還是把 executor 邏輯直接寫在 SKILL.md 中？
   - Recommendation: 保持 D-02 的決策，使用獨立 agent .md 文件。這讓每個 agent 可以被 SKILL.md 和其他 agent（如 ff）複用，也更易於測試和修改 prompt。

3. **alignment.md 的寫入時機**
   - What we know: D-07 說 alignment summary 寫入 `.specs/changes/{name}/alignment.md`
   - What's unclear: 誰負責創建這個檔案路徑？binary 預先創建空檔，還是 AI 直接用 Write tool 創建？
   - Recommendation: AI 用 Write tool 直接創建（SKILL.md 中已有 `Write` 在 allowed-tools）。Binary 不需要預先創建。這減少 binary 的責任範圍。

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go toolchain | Binary compilation | ✓ | 1.25.5 | — |
| `mysd` binary in PATH | SKILL.md `Bash` calls | ✗ (not yet built) | — | `go run ./...` (dev only) |
| Claude Code with Task tool | Wave mode subagents | ✓ (assumed) | — | Sequential inline execution |
| git | Atomic commits (D-16) | ✓ (assumed) | — | Skip atomic commits if not available |

**Missing dependencies with no fallback:**
- `mysd` binary must be built and in PATH before SKILL.md can function. Phase 2 implementation must include a `go build` step in the setup instructions.

**Missing dependencies with fallback:**
- Claude Code Task tool: if unavailable (non-Claude-Code runtime), SKILL.md falls back to sequential inline execution (documented in SKILL.md itself).

---

## Validation Architecture

nyquist_validation is enabled (no explicit false in config.json).

### Test Framework

| Property | Value |
|----------|-------|
| Framework | testify v1.11.1 + Go testing stdlib |
| Config file | none (go test ./...) |
| Quick run command | `go test ./internal/spec/... ./internal/executor/... -count=1` |
| Full suite command | `go test ./... -count=1 -race` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| EXEC-01 | `mysd execute --context-only` outputs valid JSON context | unit | `go test ./internal/executor/ -run TestContextOutput` | ❌ Wave 0 |
| EXEC-02 | Sequential mode: pending tasks returned in order | unit | `go test ./internal/executor/ -run TestSequentialTasks` | ❌ Wave 0 |
| EXEC-03 | Wave mode: flag parsing returns correct agent count | unit | `go test ./cmd/ -run TestExecuteFlags` | ❌ Wave 0 |
| EXEC-04 | Atomic commits: flag available in execute subcommand | unit | `go test ./cmd/ -run TestAtomicCommitsFlag` | ❌ Wave 0 |
| EXEC-05 | task-update: done tasks excluded from pending list | unit | `go test ./internal/executor/ -run TestResumeFromInterruption` | ❌ Wave 0 |
| WCMD-08 | status: reads STATE.json + tasks.md → correct counts | unit | `go test ./internal/executor/ -run TestStatusDashboard` | ❌ Wave 0 |
| D-13 | task-update: writes correct status to tasks.md YAML | unit | `go test ./internal/spec/ -run TestUpdateTaskStatus` | ❌ Wave 0 |
| D-07 | alignment.md written under .specs/changes/{name}/ | integration | `go test ./internal/executor/ -run TestAlignmentPath` | ❌ Wave 0 |
| TEST-01 | TDD flag available and passed to context JSON | unit | `go test ./cmd/ -run TestTDDFlag` | ❌ Wave 0 |
| TEST-03 | TDD default from config applied when flag absent | unit | `go test ./internal/config/ -run TestTDDDefault` | ❌ Wave 0 |

Plugin layer (SKILL.md / agent .md) は自動テスト不可 — manual-only（Claude Code 環境での手動検証）。

### Sampling Rate

- **Per task commit:** `go test ./internal/spec/... ./internal/executor/... -count=1`
- **Per wave merge:** `go test ./... -count=1 -race`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/executor/` package — 新規作成，含 context.go、alignment.go、progress.go
- [ ] `internal/executor/context_test.go` — covers EXEC-01, EXEC-02
- [ ] `internal/executor/progress_test.go` — covers EXEC-05, WCMD-08
- [ ] `internal/spec/updater.go` — TaskStatus round-trip 實作
- [ ] `internal/spec/updater_test.go` — covers D-13
- [ ] `cmd/task_update.go` — new subcommand
- [ ] `cmd/task_update_test.go` — covers EXEC-04 (flag), integration with updater
- [ ] `internal/config/config_test.go` 擴展 — covers TEST-03

---

## Project Constraints (from CLAUDE.md)

**Mandatory constraints the planner MUST comply with:**

| Directive | Source | Impact on Phase 2 |
|-----------|--------|-------------------|
| Tech stack: Go single binary | CLAUDE.md | All execution logic in Go binary; no Node.js scripts |
| Must be OpenSpec compatible | CLAUDE.md | tasks.md format changes must remain readable by OpenSpec tools; use schema versioning |
| Plugin form: Claude Code slash commands + agent definitions | CLAUDE.md | SKILL.md in `.claude/commands/`, agent .md in `.claude/agents/` |
| Convention over configuration | CLAUDE.md | All flags optional; defaults work out of the box; no required config before first run |
| Thin command layer | CLAUDE.md (Patterns) | cmd/*.go contains only arg parsing + internal/ calls + Printer output |
| Instance viper (not global) | CLAUDE.md (Patterns) | `internal/config.Load()` uses viper.New(); do NOT use global viper in new code |
| No bubbletea | CLAUDE.md | Status dashboard uses lipgloss only; no TUI framework |
| No MCP server | CLAUDE.md | Binary invoked from SKILL.md via Bash; no always-running process |
| GSD workflow required | CLAUDE.md | All implementation work must go through `/gsd:execute-phase` or `/gsd:quick` |

---

## Sources

### Primary (HIGH confidence)

- `D:/work_data/project/go/mysd/go.mod` — 所有依賴版本確認（yaml.v3 v3.0.1 via go.yaml.in, cobra v1.10.2, lipgloss v1.1.0, frontmatter v0.2.0）
- `D:/work_data/project/go/mysd/internal/spec/schema.go` — Task/Requirement/ItemStatus struct 已存在
- `D:/work_data/project/go/mysd/internal/state/transitions.go` — ValidTransitions 邏輯確認（ff 必須逐步 transition）
- `D:/work_data/project/go/mysd/internal/output/printer.go` — lipgloss TTY-aware pattern 確認
- `D:/work_data/project/go/mysd/internal/config/config.go` — TDD/AtomicCommits/ExecutionMode 設定欄位已存在
- `C:/Users/admin/.claude/skills/cbc-web-seed-init/SKILL.md` — SKILL.md frontmatter format (model, description, allowed-tools) 驗證
- `C:/Users/admin/.claude/commands/gsd/execute-phase.md` — Claude Code slash command format 驗證（name, description, argument-hint, allowed-tools）
- `C:/Users/admin/.claude/get-shit-done/workflows/execute-phase.md` — Wave mode Task tool invocation pattern (Task subagent_type 呼叫機制)

### Secondary (MEDIUM confidence)

- CLAUDE.md plugin 結構文件 — 記載了 `.claude/commands/` 和 `.claude/agents/` 路徑（已由現有 GSD 結構間接確認）
- CONTEXT.md D-01~D-16 — 設計決策，由討論 session 確認

### Tertiary (LOW confidence)

- None — all critical claims verified from source code or official files

---

## Metadata

**Confidence breakdown:**

- Standard stack: HIGH — all versions verified from go.mod; no new dependencies needed
- Architecture: HIGH — patterns verified from existing codebase (propose.go, printer.go, state.go) and GSD SKILL.md examples
- Pitfalls: HIGH — tasks.md round-trip and state transition pitfalls derived from existing code analysis
- Plugin format: HIGH — verified from actual SKILL.md files in ~/.claude/skills/ and ~/.claude/commands/

**Research date:** 2026-03-23
**Valid until:** 2026-04-22 (30 days — stable stack, no fast-moving dependencies)
