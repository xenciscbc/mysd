# Phase 6: Executor Wave Grouping & Worktree Engine - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-25

---

## Gray Area 1: Worktree 管理主體

**Question:** Git worktree 的 create/cleanup/pre-check 由誰負責？

**Options presented:**
- Go binary（`internal/worktree/` package 負責完整 lifecycle）
- SKILL.md bash commands（直接 `git worktree add/remove`）
- Claude Code `isolation: "worktree"` 全權管（路徑由 Claude Code 控制）
- 混合：binary 做 pre/post，Claude Code 做 execution

**Selected:** Go binary

**Rationale:** 符合 v1.0 binary-as-state-manager 架構原則，邏輯可測試、deterministic。Claude Code isolation 的路徑命名不符合 FEXEC-05 的 `.worktrees/T{id}/` 規格。

---

## Gray Area 2: 執行模式切換 UX

**Question (1/2):** Single sequential 和 wave parallel 怎麼決定？

**Options presented:**
- 自動判斷（tasks 有 depends/files 就自動 wave）
- 每次互動詢問
- flag 指定

**Selected:** 每次互動詢問

---

**Question (2/2):** 互動詢問執行模式時，選項設計是？

**Options presented:**
- Binary 選項（Sequential / Wave）永遠詢問
- 只顯示有意義的選擇（有 depends/files 才問）

**Selected:** 只顯示有意義的選擇

**Rationale:** Tasks 無 depends/files 時詢問 wave 毫無意義，直接 sequential。只在有實際並行機會時才詢問，降低使用者認知負擔。

---

## Gray Area 3: Merge 衝突失敗 UX

**Question:** AI 自動解衝突 3 次失敗後，使用者面對什麼？

**Options presented:**
- 保留 worktree + 明確指引（路徑 + 建議下一步）
- 保留 + resume 機制（`mysd execute --resume`）
- 自動 rollback + 通知（丟失 task 的修改）

**Selected:** 保留 worktree + 明確指引

**Rationale:** Resume 機制增加實作複雜度；rollback 會丟失 AI 已完成的工作。保留 worktree 是最直接的 debug 入口，加上清楚的指引就夠了。失敗是罕見 case，人工處理可接受。

---

## Gray Area 4: Wave 執行進度顯示

**Question:** Wave 並行執行時，使用者在 terminal 看到什麼？

**Options presented:**
- 每個 task 的 inline status（Wave 開始時列出 tasks，完成時 inline 更新）
- 結果流（Wave 執行間無輸出，完成後統一顯示）
- 詳細 JSON log（worktree path + branch + stdout/stderr）

**Selected:** 每個 task 的 inline status

**Rationale:** 平衡資訊量與可讀性。使用者知道哪些在跑、哪些完成，但不被 JSON 淹沒。符合現有 lipgloss Printer 的輸出模式。

---

*Discussion completed: 2026-03-25*
