---
spec-version: "1"
change: interactive-discovery
status: proposed
created: 2026-03-24
updated: 2026-03-25
---

## Summary

v1.1 milestone — 讓 mysd 具備 GSD 級別的互動式需求探索、model profile 分層、並行執行及修復機制。

## Motivation

v1.0 的 propose/spec 階段缺乏深度探索能力，產出的 proposal 和 spec 品質不足。GSD 的 discuss-phase 證明了互動式 gray area 探索 + subagent 研究的模式能顯著提升需求品質。

## Features

### 1. Interactive Discovery（propose/spec/discuss）
- **雙模式**：Research 模式（AI 研究後主導提問）+ 一般模式（使用者主導提問）
- **Research 模式**：spawn `mysd-researcher` 掃描 codebase/domain → 產出必要問題（數量不固定）→ `mysd-advisor` 並行分析每個 gray area → 帶比較表討論
- **一般模式**：使用者自己提出想討論的主題 → AI 給建議
- **雙層循環**：area 內可深挖 + 全部 areas 完成後可發現新 areas → 直到使用者滿意
- **Scope guardrail**：防止 scope creep，redirect 到 deferred notes
- 每個支援 research 的指令（propose/spec/discuss）在開始時互動式詢問是否使用 research

### 2. 新指令
- **`/mysd:discuss`**：隨時補充討論，結論自動更新 spec/design/tasks，更新 spec 後自動 re-plan + plan-checker
- **`/mysd:fix`**：互動式修復，可選 research，spawn 獨立 executor subagent（worktree 隔離），修完自動合併。fix 只改 code，spec 問題引導用 discuss
- **`/mysd:model`**：切換 model profile + resolve 機制
  - `mysd model` — 顯示目前 profile
  - `mysd model set quality` — 切換 profile
  - `mysd model resolve {agent}` — 解析特定 agent 的 model
- **`/mysd:lang`**：互動式設定語系
  - 設定 `response_language`（AI 回應語言）和 `document_language`（文件語言 / locale）
  - 同時更新 `mysd.yaml` 和 `openspec/config.yaml` 的 `locale`，保持同步

### 3. Subagent 架構
所有階段由 orchestrator (skill) 管互動，subagent 做實際工作：

| Subagent | 用途 | Model (balanced) | 狀態 |
|----------|------|-----------------|------|
| `mysd-proposal-writer` | 寫 proposal.md | sonnet | 新增 |
| `mysd-spec-writer` x N | 每個 capability area 一個 | sonnet | 改為 per spec spawn |
| `mysd-designer` | 寫 design.md | sonnet | 已有 |
| `mysd-planner` | 寫 tasks.md | opus | 已有 |
| `mysd-plan-checker` | 驗證 MUST 覆蓋率 | sonnet | 新增 |
| `mysd-researcher` | 研究 codebase/domain | sonnet | 新增 |
| `mysd-advisor` x N | 每個 gray area 並行分析 | sonnet | 新增 |
| `mysd-executor` x N | 每個 task 獨立執行 | sonnet | 改為 per task spawn |
| `mysd-verifier` | 驗證 MUST items | sonnet | 已有 |

### 4. Model Profile 分層
quality/balanced/budget 三層，orchestrator 動態指定 model：

| Agent | quality | balanced | budget |
|-------|---------|----------|--------|
| mysd-planner | opus | opus | sonnet |
| mysd-spec-writer | opus | sonnet | haiku |
| mysd-designer | opus | sonnet | haiku |
| mysd-executor | opus | sonnet | sonnet |
| mysd-verifier | sonnet | sonnet | haiku |
| mysd-scanner | sonnet | haiku | haiku |
| mysd-fast-forward | opus | sonnet | sonnet |
| mysd-uat-guide | sonnet | sonnet | haiku |
| mysd-researcher | opus | sonnet | haiku |
| mysd-advisor | opus | sonnet | haiku |
| mysd-proposal-writer | opus | sonnet | haiku |
| mysd-plan-checker | sonnet | sonnet | haiku |

### 5. Task 依賴 + 並行執行
- Planner 標記每個 task 的 `depends` 和 `files` 欄位
- Executor orchestrator 依 depends 分層，同層檢查 files overlap
- 有 overlap → 拆到不同 wave；無 overlap → 同 wave 並行
- `/mysd:execute` 互動式詢問 single/wave；ffe 和 --auto 用 config 設定

### 6. Worktree 並行執行
- 每個並行 task spawn executor with `isolation: "worktree"`
- worktree 用完整副本（方案 A），路徑自然正確
- Branch 命名：`mysd/{change-name}/T{id}-{task-slug}`
- worktree 建在 `.worktrees/T{id}/`（短路徑，Windows 相容）
- 合併：依 task ID 順序，`git merge --no-ff`
- 衝突：AI 自動解衝突 → go build + go test 驗證 → 失敗 AI 修復 → 最多 3 次 → 仍失敗通知使用者
- Cleanup：成功自動刪除 worktree + branch；失敗保留供檢查
- Wave 中一個 task 失敗，其他繼續跑完

### 7. Auto mode
- `--auto` flag 支援 propose/spec/discuss
- 跳過互動提問，自動選推薦方案
- ff/ffe 隱含 --auto
- ff/ffe 時 research 只在 propose 做一次，成果共享後續階段

### 8. Plan-checker
- `/mysd:plan` 完成後自動 spawn plan-checker
- 驗證所有 MUST items 都有對應 task
- 未通過 → 顯示缺口，問使用者自動補齊或手動調整

### 9. Codebase Scout
- propose/spec/discuss 前輕量掃描現有 codebase
- 找可重用的 patterns、相關 code、integration points
- 不需要新 subagent，orchestrator 自己 grep/glob

### 10. Scan 重構 + Init 整合
- **`/mysd:scan`** 升級為語言無關的通用掃描器（不再限 Go）：
  - 掃描專案結構 → 偵測語言/模組 → 產生 `openspec/config.yaml`（含 project metadata + locale）+ `openspec/specs/` 下的 spec 文件
  - 如果 `openspec/config.yaml` 已存在 → 增量更新 specs，不覆蓋 config
  - 首次建立 `config.yaml` 時，互動式詢問 locale（使用者選擇或輸入語言，自動轉換為合法 locale 值）
- **`/mysd:init`** 改為 `scan --scaffold-only`：
  - 只建空結構：`openspec/config.yaml` + `openspec/specs/` + `openspec/archive/`
  - 同樣互動式詢問 locale

## Key Decisions

| Decision | Rationale |
|----------|-----------|
| design 維持現狀不加 discovery | propose/spec 已夠深入，design 只記錄決策 |
| plan-checker 自動觸發 | 每次 plan 都應檢查 |
| discuss 更新 spec 後自動 re-plan | 保持 tasks 與 spec 同步 |
| 並行 task 失敗時其他繼續 | 獨立 worktree 不受影響 |
| fix 只改 code | 職責分離，spec 問題用 discuss |
| worktree 完整副本(A) | 路徑自然正確，不需特殊處理 |
| AI 自動解衝突 3 次 | 平衡自動化與人工介入 |
| ff/ffe research 做一次 | 避免重複浪費 |
| model 不寫在 agent .md | 繼承主 session 或由 orchestrator 動態指定 |

## Scope

**In scope:** 以上所有 features

**Out of scope:**
- GUI / Web 介面
- 團隊協作功能
- 支援 Claude Code 以外的 AI 工具
- design 階段的 interactive discovery
