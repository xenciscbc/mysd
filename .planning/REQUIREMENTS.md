# Requirements: my-ssd

**Defined:** 2026-03-23
**Core Value:** Spec 和執行的緊密整合 — 規格驅動 AI 執行，驗證回饋到規格，形成完整閉環

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Spec Management

- [x] **SPEC-01**: User can create structured spec artifacts (proposal.md, specs/, design.md, tasks.md) via `/mysd:propose` command
- [x] **SPEC-02**: Spec files support RFC 2119 semantic keywords (MUST / SHOULD / MAY) with machine-parseable priority levels
- [x] **SPEC-03**: User can use Delta Specs semantics (ADDED / MODIFIED / REMOVED) to describe changes to existing specs
- [x] **SPEC-04**: Spec status is tracked per-item (PENDING / IN_PROGRESS / DONE / BLOCKED) in spec metadata
- [ ] **SPEC-05**: Verification results are automatically written back to spec status (spec feedback loop)
- [ ] **SPEC-06**: Completed specs can be archived to `.specs/archive/` via `/mysd:archive` command
- [x] **SPEC-07**: Spec format uses schema-versioned frontmatter (`spec-version` field) for forward compatibility

### Execution Engine

- [ ] **EXEC-01**: User can execute spec tasks via `/mysd:execute` with pre-execution alignment gate (AI must read and acknowledge spec before writing code)
- [ ] **EXEC-02**: Default execution mode is single-agent sequential
- [ ] **EXEC-03**: User can opt into multi-agent wave execution mode with configurable agent count
- [ ] **EXEC-04**: Atomic git commits per task is available as an opt-in option
- [ ] **EXEC-05**: Execution engine tracks progress and can resume from interruption point

### Verification

- [ ] **VRFY-01**: Goal-backward verification parses all MUST items from spec and generates verification checklist
- [ ] **VRFY-02**: Verification uses an independent fresh-context agent (not the same agent that executed)
- [ ] **VRFY-03**: SHOULD items are verified with lower priority; MAY items are noted but not required
- [ ] **VRFY-04**: Verification produces structured pass/fail report per MUST/SHOULD/MAY item
- [ ] **VRFY-05**: Failed MUST items trigger gap report that can feed back into re-execution

### Workflow Commands

- [ ] **WCMD-01**: `/mysd:propose` — create new spec from user description
- [ ] **WCMD-02**: `/mysd:spec` — define detailed requirements with RFC 2119 keywords and scenarios (Given/When/Then)
- [ ] **WCMD-03**: `/mysd:design` — capture technical decisions and architecture choices
- [ ] **WCMD-04**: `/mysd:plan` — break design into executable task list with dependency analysis
- [ ] **WCMD-05**: `/mysd:execute` — run tasks with pre-execution alignment and progress tracking
- [ ] **WCMD-06**: `/mysd:verify` — goal-backward verification of all MUST items
- [ ] **WCMD-07**: `/mysd:archive` — archive completed spec to history
- [ ] **WCMD-08**: `/mysd:status` — show current spec state, progress, and verification results
- [ ] **WCMD-09**: `/mysd:scan` — scan existing project codebase and generate OpenSpec-format spec documents
- [ ] **WCMD-10**: `/mysd:ff` — fast-forward 指令，從 propose 快速推進到 plan 完成（跳過互動確認），讓使用者可直接進入實作階段
- [ ] **WCMD-11**: `/mysd:init` — 初始化專案設定檔（`.claude/mysd.yaml`），互動式設定預設偏好
- [ ] **WCMD-12**: `/mysd:uat` — 產生互動式使用者驗收測試清單（從 spec 中有 UI 相關的項目衍生），可與使用者互動逐項確認
- [ ] **WCMD-13**: `/mysd:capture` — 從當前對話中分析並提取要做的變更，自動進入 propose 的討論模式

### Roadmap Tracking

- [ ] **RMAP-01**: 實作完成後自動產生或更新 `.mysd/roadmap/` 下的追蹤文件
- [ ] **RMAP-02**: 追蹤文件記錄每個 change 的名稱、狀態、開始/完成日期時間
- [ ] **RMAP-03**: 追蹤文件格式可被第三方工具讀取（支援 roadmap 視覺化，如 Mermaid gantt chart）

### UAT Acceptance

- [ ] **UAT-01**: 驗證階段可選擇產生互動式 UAT 驗收清單（從 spec 的 UI 相關 MUST/SHOULD 項目衍生）
- [ ] **UAT-02**: UAT 清單為可選步驟，不是 archive 的前提條件
- [ ] **UAT-03**: UAT 清單存放於 `.mysd/uat/` 目錄，可跨 session 保留
- [ ] **UAT-04**: 使用者可透過 `/mysd:uat` 獨立觸發 UAT 流程，可重複執行
- [ ] **UAT-05**: UAT 清單記錄每次執行的結果（通過/未通過/跳過）與時間戳

### Testing

- [ ] **TEST-01**: User can opt into TDD mode — 先產生測試程式碼，再執行實作
- [ ] **TEST-02**: 執行完成後可選擇自動產生對應的測試程式碼（如果語言/框架支援）
- [ ] **TEST-03**: TDD 模式為可選設定，可在專案設定檔中設為預設

### Configuration

- [ ] **CONF-01**: 專案設定檔存放於 `.claude/mysd.yaml`，記憶使用者的偏好預設值
- [ ] **CONF-02**: 設定檔支援：執行模式（single/wave）、agent 數量、atomic commits、TDD 模式、測試產出等可選項目的預設值
- [ ] **CONF-03**: 設定檔支援預設回應語言（response_language）和文件產出語言（document_language）
- [ ] **CONF-04**: 所有可選項目在指令執行時可被 flag 覆蓋（flag 優先於設定檔）

### OpenSpec Compatibility

- [x] **OPSX-01**: Parser can read existing OpenSpec `openspec/` directory structure
- [x] **OPSX-02**: Parser can read and write OpenSpec's proposal.md / specs/ / design.md / tasks.md format
- [x] **OPSX-03**: Delta Specs support matches OpenSpec's ADDED / MODIFIED / REMOVED semantics
- [x] **OPSX-04**: User can point my-ssd at an existing OpenSpec project and run execute/verify without migration

### CLI & Distribution

- [x] **DIST-01**: Single Go binary with zero runtime dependencies
- [x] **DIST-02**: Cross-platform support (macOS / Linux / Windows)
- [ ] **DIST-03**: Install via `go install` and GitHub releases (precompiled binaries)
- [ ] **DIST-04**: Claude Code plugin integration via slash commands and agent definitions

### State & Session

- [ ] **STAT-01**: Project state tracked in `.specs/STATE.md` for cross-session continuity
- [ ] **STAT-02**: State machine enforces valid transitions (proposed → specced → designed → planned → executed → verified → archived)
- [ ] **STAT-03**: User can resume interrupted workflow from last valid state

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Multi-Runtime Support

- **MRUN-01**: Abstract plugin interface for supporting other AI tools (Cursor, Gemini CLI, OpenCode)
- **MRUN-02**: Plugin generator for each supported runtime

### Advanced Features

- **ADVN-01**: Spec templates / profiles for common project types
- **ADVN-02**: Brownfield codebase onboarding with AI-assisted architecture mapping (`/mysd:onboard`)
- **ADVN-03**: Spec diff visualization in terminal

## Out of Scope

| Feature | Reason |
|---------|--------|
| GUI / Web dashboard | CLI-first tool; developer audience comfortable with terminal; massive scope expansion |
| Full reverse-engineering of codebase into specs | AI-generated specs from existing code are often inaccurate; incremental spec authoring is safer |
| Team collaboration features (shared review, multi-user approval) | Solo developer + AI is the target use case; git handles multi-person collaboration |
| Real-time spec sync (auto-update specs as code changes) | Specs describe intent, code describes implementation; auto-sync inverts causality |
| 57-command surface (GSD-style) | Every additional command is maintenance burden; minimal command set covering core loop |
| Configuration-heavy setup | Convention over config; all defaults work out of the box |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| SPEC-01 | Phase 1 | Complete |
| SPEC-02 | Phase 1 | Complete |
| SPEC-03 | Phase 1 | Complete |
| SPEC-04 | Phase 1 | Complete |
| SPEC-05 | Phase 3 | Pending |
| SPEC-06 | Phase 3 | Pending |
| SPEC-07 | Phase 1 | Complete |
| EXEC-01 | Phase 2 | Pending |
| EXEC-02 | Phase 2 | Pending |
| EXEC-03 | Phase 2 | Pending |
| EXEC-04 | Phase 2 | Pending |
| EXEC-05 | Phase 2 | Pending |
| VRFY-01 | Phase 3 | Pending |
| VRFY-02 | Phase 3 | Pending |
| VRFY-03 | Phase 3 | Pending |
| VRFY-04 | Phase 3 | Pending |
| VRFY-05 | Phase 3 | Pending |
| WCMD-01 | Phase 2 | Pending |
| WCMD-02 | Phase 2 | Pending |
| WCMD-03 | Phase 2 | Pending |
| WCMD-04 | Phase 2 | Pending |
| WCMD-05 | Phase 2 | Pending |
| WCMD-06 | Phase 3 | Pending |
| WCMD-07 | Phase 3 | Pending |
| WCMD-08 | Phase 2 | Pending |
| WCMD-09 | Phase 4 | Pending |
| WCMD-10 | Phase 2 | Pending |
| WCMD-11 | Phase 2 | Pending |
| TEST-01 | Phase 2 | Pending |
| TEST-02 | Phase 2 | Pending |
| TEST-03 | Phase 2 | Pending |
| CONF-01 | Phase 1 | Pending |
| CONF-02 | Phase 1 | Pending |
| CONF-03 | Phase 1 | Pending |
| CONF-04 | Phase 1 | Pending |
| OPSX-01 | Phase 1 | Complete |
| OPSX-02 | Phase 1 | Complete |
| OPSX-03 | Phase 1 | Complete |
| OPSX-04 | Phase 1 | Complete |
| DIST-01 | Phase 1 | Complete |
| DIST-02 | Phase 1 | Complete |
| DIST-03 | Phase 4 | Pending |
| DIST-04 | Phase 4 | Pending |
| STAT-01 | Phase 1 | Pending |
| STAT-02 | Phase 1 | Pending |
| STAT-03 | Phase 1 | Pending |
| WCMD-12 | Phase 3 | Pending |
| UAT-01 | Phase 3 | Pending |
| UAT-02 | Phase 3 | Pending |
| UAT-03 | Phase 3 | Pending |
| UAT-04 | Phase 3 | Pending |
| UAT-05 | Phase 3 | Pending |
| WCMD-13 | Phase 2 | Pending |
| RMAP-01 | Phase 1 | Pending |
| RMAP-02 | Phase 1 | Pending |
| RMAP-03 | Phase 1 | Pending |

**Coverage:**
- v1 requirements: 56 total
- Mapped to phases: 56
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-23*
*Last updated: 2026-03-23 after roadmap creation*
