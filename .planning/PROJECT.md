# my-ssd

## What This Is

my-ssd 是一個用 Go 建造的 Claude Code plugin，將 OpenSpec 的 Spec-Driven Development（SDD）方法論與 GSD 級別的規劃/執行/驗證引擎整合為一個無縫系統。它讓獨立開發者（1 人 + N 個 AI agent）能以結構化規格驅動 AI 編程，確保 AI 在寫程式前先對齊需求，並在執行後自動驗證成果。

## Core Value

**Spec 和執行的緊密整合** — 規格不只是文件，而是直接驅動 AI 執行和驗證的單一事實來源。AI 寫 code 前必須對齊 spec，寫完後自動驗證 spec 是否被滿足。

## Requirements

### Validated

- ✓ 完全相容 OpenSpec 格式（proposal.md / specs/ / design.md / tasks.md） — Phase 1
- ✓ 支援 OpenSpec 的 Delta Specs（ADDED / MODIFIED / REMOVED） — Phase 1
- ✓ 支援 RFC 2119 語義關鍵字（MUST / SHOULD / MAY） — Phase 1
- ✓ 單一 Go binary 發佈，無需 Node.js runtime — Phase 1
- ✓ Convention over configuration 設計哲學 — Phase 1
- ✓ 能在既有 OpenSpec 專案上直接運作（brownfield 支援） — Phase 1

### Validated (Phase 2)

- ✓ 提供類似 OpenSpec 的指令集（propose → spec → design → plan → execute → verify → archive） — Phase 2
- ✓ Spec 作為 AI 執行的事實來源，執行前強制對齊（alignment gate） — Phase 2
- ✓ 多 agent 執行引擎（預設單 agent，--mode=wave 可選平行模式） — Phase 2
- ✓ Claude Code plugin 整合（10 slash commands + 5 agent definitions） — Phase 2

### Validated (Phase 3)

- ✓ Goal-backward 驗證機制（驗證 spec 中所有 MUST 條目是否被滿足） — Phase 3
- ✓ 驗證結果回饋到 spec（自動更新 spec 狀態） — Phase 3
- ✓ 自動產出 UAT checklist（若 spec 有 UI 相關項目），非阻塞 — Phase 3
- ✓ Archive 雙重閘門（state==verified + all MUST done）確保品質 — Phase 3
- ✓ Verifier agent 獨立性（不讀 executor artifacts） — Phase 3

### Validated (Phase 4)

- ✓ /mysd:scan 掃描既有 codebase 自動產生 OpenSpec spec 文件 — Phase 4
- ✓ 完整 Claude Code plugin（14 commands + 8 agents + hooks）— Phase 4
- ✓ GoReleaser 跨平台 binary 發佈（Linux/macOS/Windows）— Phase 4
- ✓ Roadmap tracking 自動記錄 change lifecycle（tracking.yaml + Mermaid timeline）— Phase 4

### Active

## Current Milestone: v1.1 Interactive Discovery & Parallel Execution

**Goal:** 讓 mysd 具備互動式需求探索、model profile 分層、並行執行及修復機制

**Target features:**
- Interactive Discovery — propose/spec 階段的 adaptive questioning + research 模式
- 新指令 — /mysd:discuss, /mysd:fix, /mysd:model, /mysd:lang
- Subagent 架構升級 — researcher, advisor, proposal-writer, plan-checker 等新 agents
- Model Profile 分層 — quality/balanced/budget 動態指定
- Task 依賴 + 並行執行 — depends/files 標記 + wave 自動分層
- Worktree 並行執行 — isolation: worktree + AI 自動解衝突
- Auto mode — --auto flag + ff/ffe 隱含 auto
- Plan-checker — 自動驗證 MUST 覆蓋率
- Codebase Scout — 輕量掃描現有 code
- Scan 重構 + Init 整合 — 語言無關通用掃描 + locale 互動設定

### Out of Scope

- 支援 Claude Code 以外的 AI 工具 — v1 專注 Claude Code，之後再擴展
- GUI / Web 介面 — CLI-first，不做視覺化儀表板
- 團隊協作功能（code review、多人 spec 審核）— 專注單人 + AI 場景
- GSD 的完整 57 個指令集 — 只取核心流程，精簡設計

## Context

- **v1.0 shipped** (2026-03-24): 7,555 lines Go, 11 packages, 57 requirements, 18 plans across 4 phases in 2 days
- 技術棧：Go 1.25, Cobra CLI, Viper config, lipgloss output, yaml.v3, adrg/frontmatter
- Module path: `github.com/xenciscbc/mysd`
- Distribution: `go install github.com/xenciscbc/mysd@latest` + GitHub Releases via GoReleaser
- Plugin: 14 SKILL.md commands, 8 agent definitions, SessionStart hook
- 目標用戶：獨立開發者（solo developer），使用 AI 輔助開發，希望有結構化的 spec 驅動流程而非 vibecoding
- Spec 存放位置：專案內的 `.specs/` 目錄（相容 OpenSpec 的 `openspec/` 結構）

## Constraints

- **Tech stack**: Go — 單一 binary，跨平台編譯
- **相容性**: 必須能讀寫 OpenSpec 格式的 spec 檔案
- **Plugin 形式**: Claude Code slash commands + agent definitions
- **設計哲學**: Convention over configuration — 預設即好用，只在需要時才配置

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| 用 Go 而非 Node.js | 單一 binary 部署，不依賴 runtime，安裝更簡單 | ✓ v1.0 — 7,555 LOC, cross-platform binary |
| 完全相容 OpenSpec 格式 | 讓現有 OpenSpec 用戶無縫遷移，不重新發明格式 | ✓ v1.0 — brownfield parser tested with real OpenSpec dirs |
| Spec 存放在 .specs/ 目錄 | 跟著專案走，版控追蹤，相容 OpenSpec 的目錄結構 | ✓ v1.0 — DetectSpecDir handles both .specs/ and openspec/ |
| 混合執行模式（預設單 agent） | 平衡簡單性和效能，convention over config | ✓ Phase 2 — wave mode via --mode=wave --agents=N |
| 全新系統而非基於 GSD/OpenSpec 擴展 | 避免繼承兩者的技術債，從零設計更精簡的架構 | ✓ v1.0 — clean 11-package architecture, no legacy debt |
| Plugin reverse-calling 架構 | Go binary 提供狀態管理，SKILL.md 觸發 agent definitions | ✓ Phase 2+4 — 14 commands + 8 agents |
| Alignment gate 純 prompt engineering | binary 輸出 spec 內容，agent definition 強制 AI 確認 | ✓ Phase 2 — MANDATORY section in mysd-executor.md |
| Model profile 系統 | quality/balanced/budget 配置，per-agent override | ✓ Phase 2 — ResolveModel in config package |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-25 — v1.1 milestone started*
