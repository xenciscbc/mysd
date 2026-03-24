# Roadmap: my-ssd

## Overview

my-ssd 從零開始建構一個 Go binary + Claude Code plugin 的 Spec-Driven Development 工具。建構順序由硬依賴關係決定：spec 解析器必須先於執行引擎，執行引擎必須先於驗證管道，plugin 層最後才寫（在 binary 指令穩定後）。四個 phases 各自交付一個完整、可驗證的能力，依序建立：資料模型與 CLI 骨架 → 執行引擎 → 目標反推驗證 → Plugin 層與發佈。

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Foundation** - Spec 資料模型、CLI 骨架、OpenSpec 格式解析器、狀態機 (completed 2026-03-23)
- [ ] **Phase 2: Execution Engine** - 執行引擎、pre-execution alignment gate、workflow 指令核心
- [ ] **Phase 3: Verification & Feedback Loop** - 全自動目標反推驗證、spec 狀態回寫、自動產出 UAT 文件（非阻塞）
- [x] **Phase 4: Plugin Layer & Distribution** - 完整 Claude Code plugin、所有 slash commands、GoReleaser 發佈 (completed 2026-03-24)

## Phase Details

### Phase 1: Foundation
**Goal**: 開發者可以用 `mysd propose` 建立結構化 spec artifacts，CLI skeleton 可被執行，spec 解析器能讀寫 OpenSpec 格式，狀態機追蹤 spec 的生命週期
**Depends on**: Nothing (first phase)
**Requirements**: SPEC-01, SPEC-02, SPEC-03, SPEC-04, SPEC-07, OPSX-01, OPSX-02, OPSX-03, OPSX-04, STAT-01, STAT-02, STAT-03, CONF-01, CONF-02, CONF-03, CONF-04, DIST-01, DIST-02
**Success Criteria** (what must be TRUE):
  1. User can run `mysd propose "feature description"` and get scaffolded proposal.md, specs/, design.md, tasks.md files in `.specs/`
  2. Spec files contain RFC 2119 keywords (MUST / SHOULD / MAY) that the parser correctly identifies and categorises by priority level
  3. Parser can read an existing OpenSpec `openspec/` directory without modification and produce typed Go structs
  4. State machine enforces valid transitions (proposed → specced → designed → planned → executed → verified → archived) and blocks invalid ones
  5. Project config file `.claude/mysd.yaml` is created via `mysd init` and persists user preferences across sessions
**Plans**: 3 plans

Plans:
- [x] 01-01-PLAN.md — Go module + spec schema types + parser/writer/delta/detector
- [x] 01-02-PLAN.md — State machine + config management + terminal output
- [x] 01-03-PLAN.md — CLI skeleton (Cobra) + propose/init commands

### Phase 2: Execution Engine
**Goal**: 開發者可以用 `mysd execute` 執行 spec 任務，AI 在寫 code 前必須通過 alignment gate（強制讀取並確認 spec），執行進度被追蹤且可從中斷點恢復
**Depends on**: Phase 1
**Requirements**: EXEC-01, EXEC-02, EXEC-03, EXEC-04, EXEC-05, WCMD-01, WCMD-02, WCMD-03, WCMD-04, WCMD-05, WCMD-08, WCMD-10, WCMD-11, WCMD-13, WCMD-14, TEST-01, TEST-02, TEST-03
**Success Criteria** (what must be TRUE):
  1. User runs `mysd execute` and AI cannot write code until it has explicitly acknowledged the spec content (alignment gate is non-bypassable)
  2. Execution runs single-agent sequential by default; user can opt into wave mode with `--mode=wave --agents=N`
  3. Each task in tasks.md is marked IN_PROGRESS / DONE as execution proceeds; interrupted session can resume from last completed task
  4. User can run `mysd status` and see current spec state, completed tasks, and any pending items
  5. TDD mode is available as opt-in: test code is generated before implementation when `--tdd` flag is set
**Plans**: 6 plans

Plans:
- [x] 02-01-PLAN.md — Execution engine core: tasks.md updater, executor context builder, progress tracker, alignment path
- [x] 02-02-PLAN.md — Config model profiles, status dashboard with lipgloss
- [x] 02-03-PLAN.md — CLI commands: execute, task-update, status, ff, ffe, capture, init
- [x] 02-04-PLAN.md — CLI commands: spec, design, plan with --context-only
- [x] 02-05-PLAN.md — Claude Code plugin: 10 SKILL.md + 5 agent definitions
- [x] 02-06-PLAN.md — Integration tests: execute, status, ff end-to-end verification

### Phase 3: Verification & Feedback Loop
**Goal**: 開發者可以用 `mysd verify` 觸發全自動的目標反推驗證，驗證結果自動寫回 spec 狀態，archive 指令在 MUST items 有未解決的失敗時拒絕執行。驗證過程中自動產出 UAT 文件（若 spec 有 UI 相關項目），但不阻塞任何流程。
**Depends on**: Phase 2
**Requirements**: VRFY-01, VRFY-02, VRFY-03, VRFY-04, VRFY-05, SPEC-05, SPEC-06, WCMD-06, WCMD-07, WCMD-12, UAT-01, UAT-02, UAT-03, UAT-04, UAT-05
**Success Criteria** (what must be TRUE):
  1. User runs `mysd verify` and receives a structured report listing every MUST item with PASS / FAIL verdict
  2. Verifier agent runs with fresh context (not the same agent that executed) and evaluates only filesystem evidence and spec MUST items
  3. SHOULD items appear in the report with lower priority; MAY items are noted but do not affect the overall verdict
  4. Failed MUST items generate a gap report that user can feed back into re-execution
  5. User runs `mysd archive` after all MUST items pass — command succeeds; running it with open MUST failures returns an error
  6. If spec contains UI-related items, `mysd verify` auto-generates UAT checklist files in `.mysd/uat/` without blocking the verification flow
  7. User can run `/mysd:uat` independently at any time (before archive, after archive, or much later) to interactively walk through the UAT checklist; `mysd archive` optionally prompts "Run UAT first?" but proceeds regardless of answer
**Plans**: 5 plans

Plans:
- [x] 03-01-PLAN.md — Verification engine core: VerificationContext builder, VerifierReport parser, gap report writer, spec status sidecar
- [x] 03-02-PLAN.md — UAT checklist package: data model, read/write with history preservation
- [x] 03-03-PLAN.md — CLI commands: verify (--context-only, --write-results) + archive (double gate, directory move)
- [x] 03-04-PLAN.md — Plugin layer: 3 SKILL.md (verify, archive, uat) + 2 agents (verifier, uat-guide)
- [ ] 03-05-PLAN.md — Integration tests: verify pipeline, archive pipeline, UAT round-trip

### Phase 4: Plugin Layer & Distribution
**Goal**: 完整的 Claude Code plugin 可被安裝，所有 `/mysd:*` slash commands 在 Claude Code 中可用，預編譯 binary 可透過 `go install` 和 GitHub Releases 取得
**Depends on**: Phase 3
**Requirements**: DIST-03, DIST-04, WCMD-09, RMAP-01, RMAP-02, RMAP-03
**Success Criteria** (what must be TRUE):
  1. User installs plugin by placing the `plugin/` directory under `.claude/plugins/` and all `/mysd:*` slash commands appear in Claude Code
  2. User runs `go install github.com/[owner]/mysd@latest` and gets a working binary on macOS, Linux, and Windows
  3. User can download precompiled binaries from GitHub Releases page for macOS/Linux/Windows (per D-07: no Homebrew in v1)
  4. User runs `/mysd:scan` on an existing codebase and gets OpenSpec-format spec documents generated in `.specs/`
**Plans**: 3 plans

Plans:
- [x] 04-01-PLAN.md — Scanner package + scan CLI command + version wiring
- [x] 04-02-PLAN.md — Roadmap tracking package + integration into state transitions
- [x] 04-03-PLAN.md — Plugin directory structure + GoReleaser config + scan SKILL.md/agent

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 3/3 | Complete   | 2026-03-23 |
| 2. Execution Engine | 5/6 | In Progress|  |
| 3. Verification & Feedback Loop | 2/5 | In Progress|  |
| 4. Plugin Layer & Distribution | 3/3 | Complete   | 2026-03-24 |
