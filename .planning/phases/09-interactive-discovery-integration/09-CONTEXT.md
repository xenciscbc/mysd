# Phase 9: Interactive Discovery Integration - Context

**Gathered:** 2026-03-26
**Status:** Ready for planning

<domain>
## Phase Boundary

在 Phase 8 建立的 SKILL.md orchestrator 和 agent definitions 基礎上，為 `/mysd:propose`、`/mysd:spec`、`/mysd:discuss` 三個指令加入互動式探索能力。包含：4 維度並行 research 整合、gray areas 識別與 advisor agent 分析、雙層探索循環（area 內深挖 + area 間發現新議題）、scope guardrail（deferred notes 機制）、`/mysd:note` 新指令。同時修正 Phase 8 的 `/mysd:plan` 研究架構（4 parallel → single researcher）。

</domain>

<decisions>
## Implementation Decisions

### D-01: 探索循環終止條件 — 使用者驅動而非硬性數字上限

- DISC-07 原始規格指定「每輪最多 3 個 gray areas、每個 area 最多 3 個深挖問題」+ 顯示剩餘配額
- **決策：改用使用者驅動終止**。每個 area 完成後呈現「繼續探索/完成」二元選擇，不設人工上限
- **Why:** 硬性數字限制會在有價值的討論進行中強制截斷，而使用者選擇本身就是終止條件，數字配額不增加價值且增加實作複雜度（配額追蹤）
- **How to apply:** SKILL.md 中 area 完成後顯示二元選擇即可，不需 quota counter

### D-02: Deferred notes 使用 project-level 儲存 + context-aware 載入

- Deferred notes 儲存於 `.specs/deferred.json`（project root 層級），跨 change 共享
- **propose** 指令載入 deferred notes（因為是建立新 change，跨 change context 有價值）
- **discuss** 指令在有活躍未完成 change 時**不載入** deferred notes（避免污染聚焦中的 WIP）
- **Why:** 平衡兩個需求 — 跨 change 知識重用 vs 當前 change 範圍聚焦
- **How to apply:** propose 啟動時讀取 deferred.json 作為 context；discuss 檢查 active change 狀態再決定是否載入

### D-03: /mysd:note 指令 — 手動 idea 捕捉

- `/mysd:note` 無參數 → 列出所有 deferred notes（含 ID 編號）
- `/mysd:note add {content}` → 新增 note
- `/mysd:note delete {id}` → 刪除指定 note
- 手動新增的 notes 成為未來 change proposal 的起點
- **Claude's Discretion:** 是否需要 Go binary subcommand 或可用純 SKILL.md orchestrator 實作；deferred.json 欄位結構（預期：id, content, created_at）

### D-04: Plan stage research 修正為單一 researcher

- Phase 8 的 `/mysd:plan` 錯誤實作了 4 維度並行 research，違反 DISC-03 規格
- **正確架構：plan stage 使用單一 `mysd-researcher` agent**，讀取 spec files + design.md + proposal，補充實作細節和識別技術風險
- 輸出供 designer agent 消費
- DISC-04 互動式 opt-in：plan stage 開始時詢問是否執行 research
- **架構分離：**
  - propose/discuss stages = 4 維度並行 research → 探索需求模糊地帶和 gray areas（discovery phase）
  - plan stage = 單一 focused researcher → 驗證技術可行性和補充實作細節（requirements 已定案後）
- **Phase 8 mysd-plan.md 需更新：** 移除 3 個多餘的 parallel researchers，實作 single-agent 模式

### D-05: discovery-state.json 取消

- DISC-09 success criterion 5 原指定 discovery state 持久化為 `discovery-state.json`
- **決策：取消 discovery-state.json**
- **Why:** research summary 已寫入 spec 檔案（proposal.md / specs/ / design.md），不需額外持久化層；ff/ffe 直接讀取已有的 spec 內容即可
- **How to apply:** DISC-09 success criteria 重新解讀 — `--auto` 跳過探索循環直接使用 AI 第一推薦，research summary 不需獨立 JSON 儲存

### D-06: propose/discuss 4 維度 research 整合

- Research 為可選（DISC-04），確定 topic 後互動式詢問「要啟動 4 維度 research 嗎？」
- 選擇啟用時，spawn 4 個 `mysd-researcher` agents 並行（Codebase / Domain / Architecture / Pitfalls）
- Research 結果由 SKILL.md orchestrator 識別 gray areas
- Gray areas 由 orchestrator 並行 spawn `mysd-advisor` agents 分析（產出比較表）
- **subagent 不 spawn subagent** — advisor spawning 必須在 SKILL.md orchestrator 層

### D-07: 雙層探索循環設計

- **Layer 1（area 內深挖）：** AI 或使用者提出深挖問題，逐步釐清單一 gray area
- **Layer 2（area 間發現）：** 所有目前 gray areas 完成後，使用者可選擇發現新 areas 或結束探索
- **雙模式（DISC-05）：** AI 研究後主導提問 + 使用者主導提問，兩者切換
- **終止：** 使用者主導（D-01），每個 area 完成後二元選擇，無硬性數字上限

### D-08: Scope guardrail 機制

- 探索中超出目前 spec 範圍的建議，由 AI 識別並 redirect 到 deferred notes
- 不修改當前 spec 內容 — deferred notes 是分離的儲存（`.specs/deferred.json`）
- 使用者可透過 `/mysd:note` 瀏覽和管理 deferred notes
- Deferred notes 在下次 `/mysd:propose` 時載入作為 context

### D-09: /mysd:status 增強 — 顯示 deferred notes 數量

- `/mysd:status` 底部增加 deferred notes 計數提示
- 格式：`Deferred notes: {N} — run /mysd:note to browse`
- **Why:** 讓使用者隨時知道有待處理的 ideas

### Claude's Discretion

- `/mysd:note` Go binary subcommand vs 純 SKILL.md 實作（deferred.json 讀寫可能需要 binary 支援 JSON CRUD）
- propose workflow 的步驟順序（scaffold-first vs research-first）
- spec stage 單一 researcher 的 prompt 方向設計（DISC-02）
- deferred.json 完整欄位結構（基本預期：id, content, created_at）
- 雙模式切換的 UX 設計（AI-led vs user-led 提問模式如何在 CLI 呈現）
- ROADMAP.md success criteria 如何更新以反映 D-01（使用者驅動終止）和 D-05（取消 discovery-state.json）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### 現有 SKILL.md（Phase 9 需修改）
- `.claude/commands/mysd-propose.md` — 加入 4 維度 research + gray areas + advisor + dual-loop
- `.claude/commands/mysd-discuss.md` — 加入 gray areas identification + advisor + dual-loop + scope guardrail
- `.claude/commands/mysd-plan.md` — 修正 4 parallel researchers → single researcher（D-04）
- `.claude/commands/mysd-status.md` — 加入 deferred notes count（D-09）

### 現有 Agent Definitions（Phase 9 用到的）
- `.claude/agents/mysd-researcher.md` — 4 維度 research 的執行者（propose/discuss）+ 單一 focused research（plan/spec）
- `.claude/agents/mysd-advisor.md` — gray area 分析，產出比較表
- `.claude/agents/mysd-proposal-writer.md` — propose 流程中寫 proposal.md
- `.claude/agents/mysd-spec-writer.md` — spec 更新（per spec file spawn）

### Phase 8 CONTEXT.md（基礎設計）
- `.planning/phases/08-skill-md-orchestrators-agent-definitions/08-CONTEXT.md` — discuss source detection D-06/D-07、auto_mode 機制、agent audit constraint

### v1.1 需求
- `.planning/REQUIREMENTS.md` — DISC-01~09 完整需求定義

### OpenSpec Specs
- `.specs/changes/interactive-discovery/proposal.md` — v1.1 完整功能規格（如存在）

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `.claude/commands/mysd-discuss.md` — Phase 8 的 8-step orchestration 流程（source detection, optional research, discussion loop, spec update, re-plan）是 Phase 9 的擴展基礎
- `.claude/commands/mysd-propose.md` — Phase 8 的 source detection 邏輯（D-06 6-priority chain）可直接沿用
- `.claude/commands/mysd-plan.md` — 現有 3-stage pipeline（researcher → designer → planner），Phase 9 需修正 researcher 為 single agent

### Established Patterns
- **SKILL.md orchestrator pattern**: thin orchestrator + Task tool spawn agents — Phase 9 所有新 spawning 都必須在此層
- **Per-unit spawn**: researcher/advisor/spec-writer 都採用 per-instance spawn，Phase 9 沿用
- **auto_mode context field**: 已在所有 SKILL.md 中實作，Phase 9 確保 `--auto` 跳過探索循環
- **subagent 不 spawn subagent**: advisor spawning 必須在 orchestrator 層（不在 researcher 內部）

### Integration Points
- `mysd execute --context-only` → 讀取 context JSON（含 auto_mode）
- `.specs/deferred.json` → 新增檔案，Phase 9 的 scope guardrail + /mysd:note 寫入/讀取
- `mysd task-update {id} {status}` → discuss 結論更新 tasks 後的狀態追蹤

</code_context>

<specifics>
## Specific Ideas

- propose/discuss 的 4 維度 research → gray areas → advisor 是一條連續 pipeline，SKILL.md 中要清楚劃分階段邊界
- plan stage research 修正是 Phase 8 的 bug fix，但影響範圍在 Phase 9（因為 Phase 8 已 ship，需要在 Phase 9 修正）
- `/mysd:note` 的 Go binary 實作可能最合適（deferred.json 的 JSON CRUD 在 shell 中不易可靠操作）
- 雙模式切換（AI-led / user-led）可能簡化為：research 完成後 AI 先呈現 findings + 提問 → 使用者可回答或自己提問 → 形成自然對話
- Scope guardrail 的「識別超出範圍」由 AI 判斷 — prompt engineering 是關鍵，需在 advisor agent 或 orchestrator prompt 中明確指示

</specifics>

<deferred>
## Deferred Ideas

- ROADMAP.md success criteria 更新（反映 D-01 和 D-05 的偏離）— 可在 planning 後或 verification 前處理
- `/mysd:propose` 與 `/mysd:discuss` 的進一步整合（Phase 8 deferred，繼續延後）
- `/mysd:design` 指令是否完全移除（Phase 8 deferred）
- 多語言 deferred notes 支援（locale-aware content）

</deferred>

---

*Phase: 09-interactive-discovery-integration*
*Context gathered: 2026-03-26 via discuss-phase session reconstruction*
