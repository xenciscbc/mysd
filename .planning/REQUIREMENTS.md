# Requirements: mysd

**Defined:** 2026-03-23
**Core Value:** Spec 和執行的緊密整合 — 規格驅動 AI 執行，驗證回饋到規格，形成完整閉環

## v1.0 Requirements (Complete)

<details>
<summary>57 requirements — all shipped 2026-03-24</summary>

### Spec Management

- [x] **SPEC-01**: User can create structured spec artifacts (proposal.md, specs/, design.md, tasks.md) via `/mysd:propose` command
- [x] **SPEC-02**: Spec files support RFC 2119 semantic keywords (MUST / SHOULD / MAY) with machine-parseable priority levels
- [x] **SPEC-03**: User can use Delta Specs semantics (ADDED / MODIFIED / REMOVED) to describe changes to existing specs
- [x] **SPEC-04**: Spec status is tracked per-item (PENDING / IN_PROGRESS / DONE / BLOCKED) in spec metadata
- [x] **SPEC-05**: Verification results are automatically written back to spec status (spec feedback loop)
- [x] **SPEC-06**: Completed specs can be archived to `.specs/archive/` via `/mysd:archive` command
- [x] **SPEC-07**: Spec format uses schema-versioned frontmatter (`spec-version` field) for forward compatibility

### Execution Engine

- [x] **EXEC-01**: User can execute spec tasks via `/mysd:execute` with pre-execution alignment gate (AI must read and acknowledge spec before writing code)
- [x] **EXEC-02**: Default execution mode is single-agent sequential
- [x] **EXEC-03**: User can opt into multi-agent wave execution mode with configurable agent count
- [x] **EXEC-04**: Atomic git commits per task is available as an opt-in option
- [x] **EXEC-05**: Execution engine tracks progress and can resume from interruption point

### Verification

- [x] **VRFY-01**: Goal-backward verification parses all MUST items from spec and generates verification checklist
- [x] **VRFY-02**: Verification uses an independent fresh-context agent (not the same agent that executed)
- [x] **VRFY-03**: SHOULD items are verified with lower priority; MAY items are noted but not required
- [x] **VRFY-04**: Verification produces structured pass/fail report per MUST/SHOULD/MAY item
- [x] **VRFY-05**: Failed MUST items trigger gap report that can feed back into re-execution

### Workflow Commands

- [x] **WCMD-01** ~ **WCMD-14**: 14 slash commands (propose, spec, design, plan, execute, verify, archive, status, scan, ff, init, uat, capture, ffe)

### Roadmap Tracking

- [x] **RMAP-01** ~ **RMAP-03**: Tracking files with Mermaid gantt chart support

### UAT Acceptance

- [x] **UAT-01** ~ **UAT-05**: Interactive UAT checklist from spec UI items

### Testing

- [x] **TEST-01** ~ **TEST-03**: Optional TDD mode

### Configuration

- [x] **CONF-01** ~ **CONF-04**: Project config with convention-over-config defaults

### OpenSpec Compatibility

- [x] **OPSX-01** ~ **OPSX-04**: Brownfield-compatible OpenSpec parser

### CLI & Distribution

- [x] **DIST-01** ~ **DIST-04**: Cross-platform binary + Claude Code plugin

### State & Session

- [x] **STAT-01** ~ **STAT-03**: State machine with cross-session continuity

</details>

## v1.1 Requirements

Requirements for Interactive Discovery & Parallel Execution milestone.

### Schema & Foundation

- [x] **FSCHEMA-01**: TaskEntry 支援 `depends` 欄位標記 task 間依賴關係
- [x] **FSCHEMA-02**: TaskEntry 支援 `files` 欄位標記 task 會修改的檔案
- [x] **FSCHEMA-03**: TaskEntry 支援 `satisfies` 欄位對應 MUST requirement IDs
- [x] **FSCHEMA-04**: TaskEntry 支援 `skills` 欄位標記執行時建議使用的 slash commands
- [x] **FSCHEMA-05**: Plan-checker 自動驗證所有 MUST items 都有 task 的 `satisfies` 對應（structured ID matching）
- [x] **FSCHEMA-06**: Plan-checker 未通過時顯示缺口，互動式詢問自動補齊或手動調整
- [x] **FSCHEMA-07**: openspec/config.yaml writer 可產生/讀取 OpenSpec config（含 project metadata + locale）

### Research & Discovery

- [ ] **DISC-01**: propose 階段支援 4 維度並行 research（Codebase, Domain, Architecture, Pitfalls）
- [ ] **DISC-02**: spec 階段支援單一 researcher，專注「如何實作 spec」
- [ ] **DISC-03**: plan 階段支援單一 researcher，整合 spec + design 內容並補充實作細節
- [ ] **DISC-04**: 每個支援 research 的階段（propose/spec/plan/discuss）在開始時互動式詢問是否使用 research
- [ ] **DISC-05**: Research 模式支援雙模式 — AI 研究後主導提問 + 使用者主導提問
- [ ] **DISC-06**: propose/discuss 的 research 產出 gray areas，由 SKILL.md orchestrator 並行 spawn advisor agents 分析（subagent 不 spawn subagent）
- [ ] **DISC-07**: 雙層循環 — area 內可深挖 + 全部 areas 完成後可發現新 areas，直到使用者滿意
- [ ] **DISC-08**: Scope guardrail — 防止 scope creep，超出範圍的想法 redirect 到 deferred notes
- [ ] **DISC-09**: discuss 結論自動更新 spec/design/tasks，更新後自動 re-plan + plan-checker

### Execution Engine

- [x] **FEXEC-01**: Wave grouping 演算法依 `depends` 做 topological sort 分層
- [x] **FEXEC-02**: 同層 tasks 檢查 `files` overlap，有 overlap 拆到不同 wave
- [ ] **FEXEC-03**: 每個並行 task spawn executor with `isolation: "worktree"`
- [ ] **FEXEC-04**: Worktree branch 命名 `mysd/{change-name}/T{id}-{task-slug}`
- [ ] **FEXEC-05**: Worktree 建在 `.worktrees/T{id}/`（短路徑，Windows 相容）
- [ ] **FEXEC-06**: 合併依 task ID 順序，`git merge --no-ff`
- [ ] **FEXEC-07**: AI 自動解衝突 → build + test 驗證 → 失敗 AI 修復 → 最多 3 次 → 仍失敗通知使用者
- [ ] **FEXEC-08**: 成功自動刪除 worktree + branch；失敗保留供檢查
- [ ] **FEXEC-09**: Wave 中一個 task 失敗，其他繼續跑完
- [ ] **FEXEC-10**: Worktree 建立前檢查磁碟空間（disk space guard）
- [ ] **FEXEC-11**: Windows worktree 自動設定 `git config core.longpaths true`
- [ ] **FEXEC-12**: Executor 遵守 task 的 `skills` 欄位，執行時使用指定的 slash commands

### Skills Alignment

- [ ] **SKILL-01**: Planner 自動依 task 內容推薦 `skills` 欄位
- [ ] **SKILL-02**: Plan 完成後列出所有 task 與推薦 skills 的對應表，互動式讓使用者確認
- [ ] **SKILL-03**: 使用者可逐一調整或批次同意推薦的 skills
- [ ] **SKILL-04**: ffe 模式跳過互動，直接使用推薦值

### New Commands

- [ ] **FCMD-01**: `/mysd:discuss` 隨時補充討論，支援 4 維度並行 research
- [ ] **FCMD-02**: `/mysd:fix` 互動式修復，可選 research，spawn executor subagent（worktree 隔離），只改 code
- [ ] **FCMD-03**: `/mysd:model` 顯示/切換 model profile + resolve 特定 agent model
- [ ] **FCMD-04**: `/mysd:lang` 互動式設定 response_language 和 document_language，同步 mysd.yaml 和 openspec/config.yaml
- [ ] **FCMD-05**: `/mysd:lang` 使用者可選擇或輸入語言，自動轉換為合法 locale 值

### Scan & Init

- [ ] **FSCAN-01**: `/mysd:scan` 升級為語言無關通用掃描器（不再限 Go）
- [ ] **FSCAN-02**: Scan 偵測專案語言/模組結構，產生 `openspec/config.yaml` + `openspec/specs/` 下的 spec 文件
- [ ] **FSCAN-03**: 已存在 `openspec/config.yaml` 時只增量更新 specs，不覆蓋 config
- [ ] **FSCAN-04**: 首次建立 config.yaml 時互動式詢問 locale
- [ ] **FSCAN-05**: `/mysd:init` 改為 `scan --scaffold-only`，只建空結構 + 互動式設定 locale

### Subagent Architecture

- [ ] **FAGENT-01**: 新增 `mysd-researcher` agent definition（研究 codebase/domain）
- [ ] **FAGENT-02**: 新增 `mysd-advisor` agent definition（gray area 分析，帶比較表）
- [ ] **FAGENT-03**: 新增 `mysd-proposal-writer` agent definition（寫 proposal.md）
- [x] **FAGENT-04**: 新增 `mysd-plan-checker` agent definition（驗證 MUST 覆蓋率）
- [ ] **FAGENT-05**: 所有 agent definitions 確認無 Task tool 呼叫（subagent 不 spawn subagent）
- [ ] **FAGENT-06**: `mysd-spec-writer` 改為 per capability area spawn
- [ ] **FAGENT-07**: `mysd-executor` 改為 per task spawn

### Model Profile

- [x] **FMODEL-01**: Model profile 分層表涵蓋所有新 agents（researcher, advisor, proposal-writer, plan-checker）
- [x] **FMODEL-02**: Orchestrator（SKILL.md）動態指定 model 參數給每個 spawned agent
- [x] **FMODEL-03**: quality/balanced/budget 三層完整對應表

### Auto Mode

- [ ] **FAUTO-01**: `--auto` flag 支援 propose/spec/discuss/plan
- [ ] **FAUTO-02**: `--auto` 跳過互動提問，自動選推薦方案
- [ ] **FAUTO-03**: ff/ffe 隱含 `--auto`
- [ ] **FAUTO-04**: ff/ffe 不使用 research，直接用 subagent 依照既有 spec 內容完成

## Future Requirements

### Advanced Features

- **FUT-01**: Spec diff — 視覺化比較 spec 變更前後差異
- **FUT-02**: Multi-change orchestration — 同時進行多個 changes
- **FUT-03**: Change dependency graph — changes 之間的依賴追蹤

### Multi-Runtime Support

- **MRUN-01**: Abstract plugin interface for supporting other AI tools (Cursor, Gemini CLI, OpenCode)
- **MRUN-02**: Plugin generator for each supported runtime

## Out of Scope

| Feature | Reason |
|---------|--------|
| GUI / Web 介面 | CLI-first，不做視覺化儀表板 |
| 團隊協作功能 | 專注單人 + AI 場景 |
| 支援 Claude Code 以外的 AI 工具 | v1 專注 Claude Code |
| design 階段的 interactive discovery | propose/spec 已夠深入，design 只記錄決策 |
| GSD 的完整指令集 | 只取核心流程，精簡設計 |

## Traceability

### v1.0 (Complete)

All 57 requirements mapped and shipped. See [v1.0 archive](milestones/v1.0-ROADMAP.md).

### v1.1

| Requirement | Phase | Status |
|-------------|-------|--------|
| FSCHEMA-01 | Phase 5 | Complete |
| FSCHEMA-02 | Phase 5 | Complete |
| FSCHEMA-03 | Phase 5 | Complete |
| FSCHEMA-04 | Phase 5 | Complete |
| FSCHEMA-05 | Phase 5 | Complete |
| FSCHEMA-06 | Phase 5 | Complete |
| FSCHEMA-07 | Phase 5 | Complete |
| FAGENT-04 | Phase 5 | Complete |
| FMODEL-01 | Phase 5 | Complete |
| FMODEL-02 | Phase 5 | Complete |
| FMODEL-03 | Phase 5 | Complete |
| FEXEC-01 | Phase 6 | Complete |
| FEXEC-02 | Phase 6 | Complete |
| FEXEC-03 | Phase 6 | Pending |
| FEXEC-04 | Phase 6 | Pending |
| FEXEC-05 | Phase 6 | Pending |
| FEXEC-06 | Phase 6 | Pending |
| FEXEC-07 | Phase 6 | Pending |
| FEXEC-08 | Phase 6 | Pending |
| FEXEC-09 | Phase 6 | Pending |
| FEXEC-10 | Phase 6 | Pending |
| FEXEC-11 | Phase 6 | Pending |
| FEXEC-12 | Phase 6 | Pending |
| FCMD-03 | Phase 7 | Pending |
| FCMD-04 | Phase 7 | Pending |
| FCMD-05 | Phase 7 | Pending |
| FSCAN-01 | Phase 7 | Pending |
| FSCAN-02 | Phase 7 | Pending |
| FSCAN-03 | Phase 7 | Pending |
| FSCAN-04 | Phase 7 | Pending |
| FSCAN-05 | Phase 7 | Pending |
| SKILL-01 | Phase 7 | Pending |
| SKILL-02 | Phase 7 | Pending |
| SKILL-03 | Phase 7 | Pending |
| SKILL-04 | Phase 7 | Pending |
| FCMD-01 | Phase 8 | Pending |
| FCMD-02 | Phase 8 | Pending |
| FAGENT-01 | Phase 8 | Pending |
| FAGENT-02 | Phase 8 | Pending |
| FAGENT-03 | Phase 8 | Pending |
| FAGENT-05 | Phase 8 | Pending |
| FAGENT-06 | Phase 8 | Pending |
| FAGENT-07 | Phase 8 | Pending |
| FAUTO-01 | Phase 8 | Pending |
| FAUTO-02 | Phase 8 | Pending |
| FAUTO-03 | Phase 8 | Pending |
| FAUTO-04 | Phase 8 | Pending |
| DISC-01 | Phase 9 | Pending |
| DISC-02 | Phase 9 | Pending |
| DISC-03 | Phase 9 | Pending |
| DISC-04 | Phase 9 | Pending |
| DISC-05 | Phase 9 | Pending |
| DISC-06 | Phase 9 | Pending |
| DISC-07 | Phase 9 | Pending |
| DISC-08 | Phase 9 | Pending |
| DISC-09 | Phase 9 | Pending |

**Coverage:**
- v1.1 requirements: 45 total
- Mapped to phases: 45
- Unmapped: 0

---
*Requirements defined: 2026-03-23*
*Last updated: 2026-03-25 — v1.1 roadmap created, all 45 requirements mapped to Phases 5-9*
