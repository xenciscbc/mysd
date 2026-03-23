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

### Active

- [ ] 提供類似 OpenSpec 的指令集（propose → spec → design → plan → execute → verify → archive）
- [ ] Spec 作為 AI 執行的事實來源，執行前強制對齊
- [ ] 多 agent 執行引擎（預設單 agent，複雜任務可選平行模式）
- [ ] Goal-backward 驗證機制（驗證 spec 中所有 MUST 條目是否被滿足）
- [ ] 驗證結果回饋到 spec（自動更新 spec 狀態）
- [ ] Claude Code plugin 整合（slash commands + agents）

### Out of Scope

- 支援 Claude Code 以外的 AI 工具 — v1 專注 Claude Code，之後再擴展
- GUI / Web 介面 — CLI-first，不做視覺化儀表板
- 團隊協作功能（code review、多人 spec 審核）— 專注單人 + AI 場景
- GSD 的完整 57 個指令集 — 只取核心流程，精簡設計

## Context

- 現有工具的缺口：OpenSpec 有完整的 SDD 方法論但沒有執行引擎；GSD 有強大的執行引擎但沒有 spec 管理。my-ssd 填補這個缺口。
- 技術棧：Go 語言，單一 binary 部署。Claude Code plugin 整合層透過 slash commands 和 agent definitions 實現。
- 目標用戶：獨立開發者（solo developer），使用 AI 輔助開發，希望有結構化的 spec 驅動流程而非 vibecoding。
- 設計參考：OpenSpec 的 artifact-guided workflow（proposal → specs → design → tasks）、GSD 的 multi-agent orchestration 和 wave-based execution。
- Spec 存放位置：專案內的 `.specs/` 目錄（相容 OpenSpec 的 `openspec/` 結構）。

## Constraints

- **Tech stack**: Go — 單一 binary，跨平台編譯
- **相容性**: 必須能讀寫 OpenSpec 格式的 spec 檔案
- **Plugin 形式**: Claude Code slash commands + agent definitions
- **設計哲學**: Convention over configuration — 預設即好用，只在需要時才配置

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| 用 Go 而非 Node.js | 單一 binary 部署，不依賴 runtime，安裝更簡單 | — Pending |
| 完全相容 OpenSpec 格式 | 讓現有 OpenSpec 用戶無縫遷移，不重新發明格式 | — Pending |
| Spec 存放在 .specs/ 目錄 | 跟著專案走，版控追蹤，相容 OpenSpec 的目錄結構 | — Pending |
| 混合執行模式（預設單 agent） | 平衡簡單性和效能，convention over config | — Pending |
| 全新系統而非基於 GSD/OpenSpec 擴展 | 避免繼承兩者的技術債，從零設計更精簡的架構 | — Pending |

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
*Last updated: 2026-03-23 after Phase 1 completion*
