# Phase 6: Executor Wave Grouping & Worktree Engine - Context

**Gathered:** 2026-03-25
**Status:** Ready for planning

<domain>
## Phase Boundary

依 `depends` 欄位對 tasks 做 topological sort 分 wave；同 wave 內 `files` 無 overlap 的 tasks 在 Go binary 管理的獨立 git worktree 中並行執行；merge 後自動清理，衝突交由 AI 最多嘗試 3 次，失敗保留 worktree 供人工處理。

</domain>

<decisions>
## Implementation Decisions

### Worktree 管理主體
- **D-01:** Go binary 的 `internal/worktree/` 新 package 負責所有 git worktree lifecycle：create（`git worktree add`）、remove（`git worktree remove`）、disk space check（FEXEC-10）、Windows longpaths 設定（`git config core.longpaths true`，FEXEC-11）。SKILL.md 透過呼叫 `mysd worktree` subcommand 委派。符合 v1.0 binary-as-state-manager 架構原則，邏輯可測試、deterministic。
- **D-02:** Worktree 路徑：`.worktrees/T{id}/`（短路徑，Windows MAX_PATH 相容，沿用 ROADMAP 規格）。Branch 命名：`mysd/{change-name}/T{id}-{slug}`（沿用 ROADMAP 規格）。

### 執行模式切換 UX
- **D-03:** 只在有實際並行機會時才詢問執行模式：
  - Tasks 都無 `depends` 且都無 `files` 欄位 → 直接 sequential，不詢問
  - Tasks 有 `depends` 或 `files` 欄位（有並行機會）→ 詢問「Sequential（安全穩定）/ Wave parallel（N 個 tasks 並行）」
- **D-04:** `ffe` 和 `--auto` 跳過詢問：有 `depends`/`files` 時用 wave mode，否則用 sequential。符合 FAUTO-03/FAUTO-04 的 auto 語義。

### Merge 衝突失敗 UX
- **D-05:** AI 自動解衝突最多嘗試 3 次（每次：解衝突 → build + test 驗證），失敗後：
  - 保留 worktree 在 `.worktrees/T{id}/`
  - 顯示清楚的錯誤訊息：失敗原因 + worktree 路徑 + branch 名稱 + 建議下一步（`cd .worktrees/T{id}` 手動解衝突）
  - 該 wave 其他已成功的 tasks 照常 merge（continue-on-failure policy，沿用 Phase 5 討論決策）
  - 無自動 resume 機制，人工解決後再次執行
- **D-06:** Merge 成功的 worktree 自動刪除（FEXEC-08），失敗的保留。

### Wave 執行進度顯示
- **D-07:** 使用 lipgloss Printer 輸出 inline status：
  - Wave 開始：`Wave 1/3: T1, T2, T3 並行執行中...`
  - Task 完成：inline 輸出 `T1 ✓`（成功）或 `T1 ✗`（失敗）
  - Wave 結束：摘要行 `Wave 1 complete: 2 succeeded, 1 failed`
  - 格式簡潔，符合現有 lipgloss Printer 模式

### Claude's Discretion
- `internal/worktree/` 的 Go interface 設計細節（struct vs function-based API）
- disk space check 的臨界值（多少 MB 算不足）
- topological sort 演算法選擇（Kahn's vs DFS — Phase 5 D-04 已標示 Claude discretion）
- wave grouping 時 files overlap 的比較邏輯（exact match vs prefix match）
- `mysd worktree` subcommand 的 CLI 介面細節

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### 現有 Executor 基礎
- `internal/executor/context.go` — TaskItem 結構（含 Depends/Files/Skills 欄位，Phase 5 新增）、BuildContextFromParts、ExecutionContext
- `cmd/execute.go` — 現有 execute command，Phase 6 需要擴展支援 wave mode

### Schema（Phase 5 成果）
- `internal/spec/schema.go` — TaskEntry（含 Depends/Files 欄位，L76-85）
- `internal/config/defaults.go` — ProjectConfig（含 WorktreeDir/AutoMode，Phase 5 新增）

### Architecture Research
- `.planning/research/ARCHITECTURE.md` — v1.1 架構研究，含 `internal/worktree/` package 建議和 wave grouping 整合點
- `.planning/phases/05-schema-foundation-plan-checker/05-CONTEXT.md` — Phase 5 decisions（D-04: WaveGroups 在 Go binary 計算；D-12: TaskItem 已有完整欄位）

### v1.1 Spec
- `.specs/changes/interactive-discovery/proposal.md` §5-6 — Task 依賴 + 並行執行設計、Worktree 並行執行規格
- `.planning/REQUIREMENTS.md` — FEXEC-01 ~ FEXEC-12 完整需求定義

### 輸出格式
- `internal/output/` — lipgloss Printer（Wave 進度 inline status 的輸出工具）

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/executor/context.go TaskItem.Depends` + `.Files`：Phase 5 已有，Phase 6 wave grouping 直接讀取
- `internal/config.ProjectConfig.WorktreeDir`：Phase 5 新增，已可存放 worktree root path
- `internal/config.ProjectConfig.AutoMode`：Phase 5 新增，控制 ffe/--auto 行為
- `internal/output.Printer`：lipgloss 輸出，用於 wave 進度顯示

### Established Patterns
- **Binary-as-state-manager**: cmd/ 只做參數解析和輸出，業務邏輯在 internal/ package
- **Convention-over-config**: WorktreeDir 預設 `.worktrees/`，不存在時用 default
- **Sidecar pattern**: 進度/狀態用 JSON sidecar，不修改 tasks.md（可沿用追蹤 wave 狀態）
- **Pure function packages**: planchecker 是純函數，worktree 邏輯應盡量 pure（I/O 集中在 cmd 層）

### Integration Points
- `cmd/execute.go` — 需要擴展：讀取 WaveGroups → 詢問模式 → 呼叫 internal/worktree → spawn executors
- `internal/executor/context.go BuildContext` — PendingTasks 已篩選，wave grouping 從 pending tasks 的 Depends/Files 計算
- `cmd/plan.go --context-only` — WaveGroups 欄位已預留（Phase 5 D-04），Phase 6 實際填入計算結果

</code_context>

<specifics>
## Specific Ideas

- `internal/worktree/` 作為 Go binary 新 package，負責所有 git worktree 生命週期，保持業務邏輯可測試
- 執行模式詢問只在有意義時出現（有 depends/files 才問），避免每次都要回答無意義的選擇
- 失敗的 worktree 保留是刻意設計的 debug 機制，而非僅僅容錯

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 06-executor-wave-grouping-worktree-engine*
*Context gathered: 2026-03-25*
