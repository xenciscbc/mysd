# Phase 2: Execution Engine - Context

**Gathered:** 2026-03-23
**Status:** Ready for planning

<domain>
## Phase Boundary

實作執行引擎核心：讓 AI 在寫 code 前強制對齊 spec（alignment gate），透過 Claude Code plugin skill + agent definition 編排完整的 spec-driven workflow（propose → spec → design → plan → execute），追蹤任務進度且支援中斷恢復。涵蓋所有 workflow 指令（spec/design/plan/execute/status/ff/init/capture）的 AI 互動實作。

</domain>

<decisions>
## Implementation Decisions

### AI 調用機制
- **D-01:** Plugin skill 反向調用 — `/mysd:execute` 等指令由 SKILL.md 作為入口觸發，Go binary 負責 spec 解析和狀態管理，Claude Code agent 負責實際 AI 執行
- **D-02:** Agent definition 編排 — 每個 workflow 階段有專屬 agent .md 檔案（仿 GSD 模式），如 mysd-spec-writer.md、mysd-designer.md、mysd-planner.md、mysd-executor.md、mysd-verifier.md，SKILL.md 編排調用
- **D-03:** Wave mode 使用 Claude Code 原生 Agent tool 生成多個平行 subagent，每個處理一個 task（GSD 也用此機制）
- **D-04:** Profile-based 模型管理 — 仿 GSD 的 resolve-model 機制，在 mysd.yaml 中配置 model profile（quality/balanced/budget），每個 agent 類型根據 profile 映射到具體模型，可在 mysd.yaml 覆蓋預設

### Alignment Gate 設計
- **D-05:** Prompt 注入 + 確認指令 — agent definition 中強制包含 spec 內容，要求 AI 在回應中明確列出它理解的 MUST/SHOULD/MAY 項目，然後才能開始寫 code
- **D-06:** 結構化摘要輸出 — AI 必須輸出 alignment summary：列出所有 MUST 項目、它理解的執行策略、任何疑問。此 summary 可被後續 verify 階段引用
- **D-07:** Alignment summary 寫入 `.specs/changes/{name}/alignment.md`，跟著 spec artifacts 走，版控可追蹤

### Workflow 指令深度
- **D-08:** spec / design / plan 三個指令都有完整的 AI 互動流程，各由專屬 agent 執行（類似 propose 的互動體驗）
- **D-09:** `mysd ff`（fast-forward）從 propose 一氣推進到 plan 完成，跳過互動確認（用預設值），結果是完整的 spec artifacts + plan，使用者可直接 execute
- **D-10:** `mysd capture` 分析當前 Claude Code 對話中討論過的變更，提取關鍵需求，然後帶預填內容自動進入 propose 流程（減少重複描述）
- **D-11:** `mysd status` 顯示綜合儀表板：當前 change name、workflow phase、任務完成率（X/Y tasks done）、MUST/SHOULD/MAY 達成狀態、上次執行時間
- **D-12:** Plan 階段可選管線 — 預設只有 plan（快速），可用 flag 啟用 research 和 plan-check（完整管線仿 GSD 的 research → plan → check）。Convention-over-configuration：快速為預設，完整為可選

### 任務進度與恢復
- **D-13:** Agent 回報 + binary 更新 — SKILL.md / agent definition 要求 AI 在開始和完成每個 task 時呼叫 `mysd task-update {id} {status}`，Go binary 更新 tasks.md frontmatter 和 STATE.json
- **D-14:** Task level 中斷恢復 — 從最後一個完成的 task 之後恢復，已完成的不重做。基於 tasks.md 中的 status 欄位
- **D-15:** TDD mode: test-first 指令注入 — 啟用 TDD 時，agent definition 增加指令要求先寫測試再寫實作，任務流程變成 RED → GREEN → REFACTOR。透過 prompt engineering 強制 TDD 行為
- **D-16:** Atomic commits 粒度為每個 task 一個 commit（--atomic-commits flag 啟用時）。每個 commit 對應一個有意義的變更單位

### Claude's Discretion
- Agent definition 的具體 prompt 措辭和結構
- Alignment summary 的具體 markdown 模板
- `mysd task-update` 的 CLI 介面設計（flag 名稱、輸出格式）
- Status 儀表板的 lipgloss 配色和排版
- ff 指令的預設值選擇策略
- Model profile 的具體模型映射表（哪個 profile 對應哪個模型）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### 專案架構
- `.planning/PROJECT.md` — 專案願景、約束條件、核心價值
- `.planning/REQUIREMENTS.md` — 完整 v1 需求清單，Phase 2 需覆蓋 EXEC-01~05, WCMD-01~05/08/10/11/13, TEST-01~03
- `.planning/ROADMAP.md` — Phase 2 goal 和 success criteria

### Phase 1 基礎
- `.planning/phases/01-foundation/01-CONTEXT.md` — Phase 1 決策記錄（目錄結構、CLI 設計、狀態管理等）
- `.planning/phases/01-foundation/01-RESEARCH.md` — Phase 1 技術研究（stack、architecture、pitfalls）

### Claude Code Plugin 文件
- Claude Code 官方文件：SKILL.md 格式、plugin 結構、agent definition 格式（已在 CLAUDE.md 中記錄）
- RFC 2119 (https://datatracker.ietf.org/doc/html/rfc2119) — MUST/SHOULD/MAY 關鍵字定義

### 現有程式碼
- `internal/spec/schema.go` — Spec 資料模型（Change, Requirement, Task, RFC2119Keyword, ItemStatus）
- `internal/spec/parser.go` — OpenSpec 解析器
- `internal/state/state.go` — WorkflowState 和 JSON 持久化
- `internal/state/transitions.go` — Phase transition 驗證
- `internal/config/` — Viper-based 專案設定管理
- `internal/output/` — TTY-aware lipgloss terminal printer
- `cmd/root.go` — Cobra CLI 根命令和 persistent flags
- `cmd/execute.go` — Execute 指令 stub（待實作）

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/spec.Change` struct: 完整的 change 結構（Proposal, Specs, Design, Tasks, Meta），executor agent 可直接消費
- `internal/spec.Task` struct: 有 ID, Name, Description, Status 欄位，支援 task-level 追蹤
- `internal/spec.Requirement` struct: 有 Keyword (MUST/SHOULD/MAY) 和 Status，alignment gate 可用
- `internal/state.Transition()`: 已有 phase 轉換驗證，execute 指令可直接使用
- `internal/config.Load()`: Viper 設定載入，可擴展支援 model profile
- `internal/output.Printer`: TTY-aware 輸出，status 儀表板可直接使用
- `cmd/root.go` persistent flags: 已有 --tdd, --atomic-commits, --execution-mode, --agent-count flags

### Established Patterns
- **Thin command layer**: cmd/*.go 不含業務邏輯，只做參數解析 + 呼叫 internal/ + 用 Printer 輸出
- **Convention-over-config**: 缺少設定檔時用 Defaults()，不報錯
- **Instance viper**: internal/config 用 viper.New() 而非 global viper（測試隔離）
- **TDD workflow**: Phase 1 建立了 RED → GREEN 的 commit 模式

### Integration Points
- `cmd/execute.go`, `cmd/spec.go`, `cmd/design.go`, `cmd/plan.go`, `cmd/status.go` — 都是 stub，待實作
- Plugin 層（SKILL.md、agent .md）— Phase 2 新建，Phase 4 打包發佈
- `mysd task-update` — 新 subcommand，需加入 cmd/

</code_context>

<specifics>
## Specific Ideas

- Agent definition 仿 GSD 模式 — 每個 workflow 階段有獨立的 agent .md，由 SKILL.md 編排調用
- Model profile 仿 GSD 的 resolve-model 機制 — 不自己發明新方式
- Plan 管線可選深度（研究/檢查）符合 convention-over-config — 預設快速，進階可選
- Capture 指令是 UX 創新重點 — 從對話脈絡直接進入 SDD 流程

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 02-execution-engine*
*Context gathered: 2026-03-23*
