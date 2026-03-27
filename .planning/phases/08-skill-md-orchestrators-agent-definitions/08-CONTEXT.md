# Phase 8: SKILL.md Orchestrators & Agent Definitions - Context

**Gathered:** 2026-03-26
**Status:** Ready for planning

<domain>
## Phase Boundary

建立 `mysd-researcher`、`mysd-advisor`、`mysd-proposal-writer` 三個新 agent definitions；重構現有 agents（executor per task spawn、spec-writer per spec spawn）；新增 `/mysd:discuss` 和 `/mysd:fix` SKILL.md 指令；改寫 `/mysd:plan`（加入 researcher → designer → planner 流程）、`/mysd:ff`（plan + apply + archive）、`/mysd:ffe`（research + plan + apply + archive）；`execute` 改名為 `apply`；`--auto` flag 跨 propose/spec/discuss/plan 運作；ff/ffe 隱含 auto mode。

</domain>

<decisions>
## Implementation Decisions

### `/mysd:discuss` 指令設計

- **D-01:** Topic 輸入：`$ARGUMENTS` 有 topic 字串就直接使用；無 argument 時先問「你想討論什麼主題？」
- **D-02:** Research 可選：確定 topic 後詢問「要啟動 4 維度 research 嗎？」—— research 不是強制的，使用者可選跳過
- **D-03:** Research 執行：spawn 4 個 `mysd-researcher` agents 並行，各自從 Codebase / Domain / Architecture / Pitfalls 維度研究，結果逐題與使用者確認
- **D-04:** 結論決策：結論達成後問「要納入 spec 還是繼續深入討論？」—— 使用者主導，可多輪
- **D-05:** Spec 更新：結論確認納入後，依影響的 spec 層級委派給原本負責的 agent：
  - proposal 層 → spawn `mysd-proposal-writer`
  - specs/ 層 → spawn `mysd-spec-writer`（per spec file）
  - design 層 → spawn `mysd-designer`
  - 更新完後 → re-plan（spawn `mysd-planner`）→ plan-checker
- **D-06:** Source 解析規則（discuss 與 propose 共用相同邏輯）：
  1. `$ARGUMENTS` = change name（匹配 `.specs/changes/{name}/`）→ mysd change 模式
  2. `$ARGUMENTS` = 檔案路徑 → 單檔模式
  3. `$ARGUMENTS` = 目錄路徑 → 選擇模式（列出 .md，複選）
  4. 無 argument + 有活躍 change → 用當前 change
  5. 無 argument + 無活躍 change → 自動偵測（對話上下文 + 支援路徑），複選，預設勾最新 .md
  6. 都找不到 → 新建（類似 propose 流程）
- **D-07:** 三種支援的文件路徑（無 argument 時的自動偵測來源）：
  1. 使用者自定位置（明確指定的任意路徑）
  2. `~/.gstack/projects/{project}/`（含 design、test plan 等所有 .md）
  3. 對話上下文中提到的計畫文件
  - 注意：`.claude/plans/` 因檔名為隨機 hash 無 project 資訊，**不納入**自動偵測

### `/mysd:fix` 指令設計

- **D-08:** 路徑自動偵測：fix 讀取 task sidecar 失敗記錄 + 檢查 worktree 狀態自動判斷路徑（有 conflict markers → merge 衝突路徑；build/test 失敗但無 conflict → 實作問題路徑），偵測結果告知使用者確認
- **D-09:** Argument 格式（兩種都支援）：
  - `/mysd:fix T2` — 操作當前活躍 change 的 T2
  - `/mysd:fix {change-name} T2` — 指定 change 的 T2
- **D-10:** 無 argument 時：列出當前 change 中 failed/blocked 的 tasks，讓使用者選擇要 fix 哪個
- **D-11:** Research 觸發：research 選項只在**實作問題路徑**出現（merge 衝突路徑解法明確，無需 research）
- **D-12:** 實作問題診斷：AI 自動讀 task sidecar 的失敗原因、AI 嘗試解法、test 輸出，自動診斷並說明後執行修復
- **D-13:** 純 SKILL.md + agent 實作，不需新 Go binary subcommand（現有 `mysd execute --context-only`、`mysd task-update`、`mysd worktree` 已足夠）
- **D-14:** 兩條路徑（STATE.md 設計為定案）：
  - **Merge 衝突路徑**：解衝突 → build + test 驗證 → merge → `-D` 強制刪 branch
  - **實作問題路徑**：修正 task 內容 + 更新 spec → `-D` 刪舊 branch → spawn executor 重新執行
  - **放棄路徑**：task 回 `pending`，`-D` 刪 branch + worktree
  - **Skipped tasks 恢復**：fix 成功 merge 後，下游 skipped tasks 遞移性恢復為 `pending`

### Agent 重構（FAGENT-05/06/07）

- **D-15:** FAGENT-07（executor per task）：無論 single 或 wave mode，每個 task 都是獨立的新 agent 實例；single mode 改為序列 spawn（一個接一個），wave mode 維持並行 spawn
- **D-16:** FAGENT-06（spec-writer per spec）：每個 spec 檔案一個新 `mysd-spec-writer` agent 實例，不再一個 agent 處理所有 spec areas
- **D-17:** FAGENT-05（agent audit）：執行時人工審計，確認所有 9 個 agent definitions 無 Task tool 呼叫；Task tool 只允許在頂層 SKILL.md orchestrator 層使用

### `--auto` Flag 傳遞機制

- **D-18:** SKILL.md 層解析 `$ARGUMENTS` 中的 `--auto` flag，spawn agent 時將 `auto_mode: true` 放入 context JSON；Go binary 不需新 flag
- **D-19:** ff/ffe 隱含 auto mode — agent 在每個決策點自主選擇最佳方案，不詢問使用者
- **D-20:** ff/ffe 不使用 research（FAUTO-04）；單獨加 `--auto` 的指令（如 `/mysd:propose --auto`）仍可有 research，只是跳過互動確認

### 工作流程調整

- **D-21:** 新工作流程：`propose/discuss` → `discuss`（選擇性）→ `plan` → `apply` → `archive`
- **D-22:** `execute` 改名為 `apply`（SKILL.md 指令 `/mysd:execute` → `/mysd:apply`，agent `mysd-executor` 仍保持原名）
- **D-23:** `/mysd:plan` 內部流程改寫：
  1. [選擇性] spawn `mysd-researcher` × 4 並行（4 維度 research，輸出作為 design 輸入）
  2. spawn `mysd-designer` → 產出 `design.md`（架構決策）
  3. spawn `mysd-planner` → 讀 spec + design.md → 產出 `tasks.md`
- **D-24:** `/mysd:ff` 改寫 = plan（無 research）+ apply + archive（假設 propose/discuss 已完成，spec 就緒）
- **D-25:** `/mysd:ffe` 改寫 = plan（4 維度 research）+ apply + archive（research 目的是為 design 提供技術依據）
- **D-26:** design.md 由 plan 階段的 `mysd-designer` sub-agent 產生，不再是獨立 SKILL.md 步驟

### `/mysd:propose` 自動偵測輸入來源

- **D-27:** source 偵測邏輯與 discuss 完全相同（D-06/D-07）：argument > 活躍 change > 自動偵測（gstack/對話上下文/使用者自定）> 新建
- **D-28:** 原 STATE.md 中提到的 `.planning/phases/{phase}/` 路徑移除（GSD planning 路徑因 naming 問題無法自動識別 project）

### `/mysd:status` SKILL.md 設計

- **D-29:** 顯示新版 workflow stage（`propose` → `plan` → `apply` → `archive`），標示目前所在位置
- **D-30:** Task 列表含編號、title、狀態符號（✓ done / ✗ failed / ⊘ skipped / ○ pending）
- **D-31:** 最後一行：`Next: /mysd:{command}` 推薦下一步指令

### Claude's Discretion

- 4 個新 agent definitions 的 allowed-tools 列表細節
- `mysd-researcher` / `mysd-advisor` / `mysd-proposal-writer` 的 prompt 設計細節
- `/mysd:fix` SKILL.md 中偵測 conflict markers 的具體方式
- `mysd-advisor` 的比較表格格式設計
- apply（原 execute）命令的 state transition binary call 名稱（`mysd apply` 或保留 `mysd execute`）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### 現有 SKILL.md（需改寫）
- `.claude/commands/mysd-plan.md` — 現有 plan orchestrator（需加入 researcher → designer → planner 流程）
- `.claude/commands/mysd-ff.md` — 現有 ff（需改寫為 plan + apply + archive）
- `.claude/commands/mysd-ffe.md` — 現有 ffe（需改寫為 research + plan + apply + archive）
- `.claude/commands/mysd-execute.md` — 現有 execute（改名為 apply）
- `.claude/commands/mysd-propose.md` — 現有 propose（加入 D-27 source 偵測邏輯）
- `.claude/commands/mysd-status.md` — 現有 status（改寫為 D-29~31 設計）

### 現有 Agent Definitions（需審計/重構）
- `.claude/agents/mysd-executor.md` — 改為 per task spawn（D-15）
- `.claude/agents/mysd-spec-writer.md` — 改為 per spec file spawn（D-16）
- `.claude/agents/mysd-planner.md` — 確認無 Task tool 呼叫
- `.claude/agents/mysd-designer.md` — 升格為 plan 流程必要子 agent
- `.claude/agents/mysd-verifier.md` — 確認無 Task tool 呼叫
- `.claude/agents/mysd-scanner.md` — 確認無 Task tool 呼叫
- `.claude/agents/mysd-uat-guide.md` — 確認無 Task tool 呼叫
- `.claude/agents/mysd-fast-forward.md` — 確認無 Task tool 呼叫

### v1.1 需求
- `.planning/REQUIREMENTS.md` — FCMD-01, FCMD-02, FAGENT-01~03, FAGENT-05~07, FAUTO-01~04 完整需求定義

### Phase 6 成果（worktree 操作基礎）
- `.planning/phases/06-executor-wave-grouping-worktree-engine/06-CONTEXT.md` — worktree lifecycle 設計（fix 路徑操作的基礎）

### Phase 7 成果（skills 推薦流程）
- `.planning/phases/07-new-binary-commands-scanner-refactor/07-CONTEXT.md` — D-07~D-10 skills 推薦與確認設計

### 現有 Specs
- `.specs/changes/interactive-discovery/proposal.md` — v1.1 完整功能規格

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `.claude/commands/mysd-plan.md` — 現有 Task tool spawn 模式（plan → planner agent）可延伸為三段式
- `.claude/commands/mysd-execute.md` — 現有 wave/single spawn 邏輯供 apply 參考
- `.claude/agents/mysd-planner.md` — Step 4.5 skills 推薦、Step 7.5 確認流程已實作

### Established Patterns
- **SKILL.md orchestrator pattern**：thin orchestrator + Task tool spawn agents（Phase 2 決策）
- **Per-unit spawn**：wave mode executor 已是 per task spawn，Phase 8 將 single mode 統一
- **auto_mode context field**：execution context JSON 已有此欄位，Phase 8 擴展到所有 SKILL.md
- **subagent 不 spawn subagent**：v1.1 roadmap 決策，Phase 8 audit 確認

### Integration Points
- `mysd execute --context-only` → fix 讀取 task 失敗狀態
- `mysd task-update {id} {status}` → fix 更新 task 狀態
- `mysd worktree remove/list` → fix 清理或保留 worktree
- `mysd plan --context-only` → plan orchestrator 讀取 context 並傳給 researcher/designer/planner

</code_context>

<specifics>
## Specific Ideas

- discuss 和 propose 的 source 偵測共用相同邏輯（D-06/D-07）——這是一個可以抽取為共用流程說明的模式
- ffe 的 research 目的是「為 design 提供技術依據」，不是「探索需求」——這個語義差異很重要，影響 researcher 的 prompt 設計
- fix 的「偵測後告知使用者確認」是刻意設計的安全閥——自動偵測不等於靜默執行
- execute → apply 的改名只影響 SKILL.md 層，binary subcommand 名稱由 Claude's Discretion 決定

</specifics>

<deferred>
## Deferred Ideas

- `/mysd:propose` 與 `/mysd:discuss` 的合併討論（保持獨立，discuss 新建 change 功能重疊但不取代 propose）
- `design` 步驟是否完全消失（目前定案：design.md 由 plan 內部 sub-agent 產生，`/mysd:design` SKILL.md 是否保留待評估）
- Phase 9：propose/discuss 完整互動式雙層循環（DISC-01~09）

</deferred>

---

*Phase: 08-skill-md-orchestrators-agent-definitions*
*Context gathered: 2026-03-26*
