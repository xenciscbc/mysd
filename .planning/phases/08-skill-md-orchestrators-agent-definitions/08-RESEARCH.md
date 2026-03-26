# Phase 8: SKILL.md Orchestrators & Agent Definitions - Research

**Researched:** 2026-03-26
**Domain:** Claude Code plugin layer (SKILL.md orchestrators, agent definitions, --auto flag propagation)
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**`/mysd:discuss` 指令設計**
- D-01: Topic 輸入：`$ARGUMENTS` 有 topic 字串就直接使用；無 argument 時先問「你想討論什麼主題？」
- D-02: Research 可選：確定 topic 後詢問「要啟動 4 維度 research 嗎？」—— research 不是強制的，使用者可選跳過
- D-03: Research 執行：spawn 4 個 `mysd-researcher` agents 並行，各自從 Codebase / Domain / Architecture / Pitfalls 維度研究，結果逐題與使用者確認
- D-04: 結論決策：結論達成後問「要納入 spec 還是繼續深入討論？」—— 使用者主導，可多輪
- D-05: Spec 更新：結論確認納入後，依影響的 spec 層級委派給原本負責的 agent：proposal 層 → `mysd-proposal-writer`；specs/ 層 → `mysd-spec-writer`（per spec file）；design 層 → `mysd-designer`；更新完後 → re-plan（spawn `mysd-planner`）→ plan-checker
- D-06: Source 解析規則（discuss 與 propose 共用相同邏輯）：(1) `$ARGUMENTS` = change name → mysd change 模式；(2) `$ARGUMENTS` = 檔案路徑 → 單檔模式；(3) `$ARGUMENTS` = 目錄路徑 → 選擇模式；(4) 無 argument + 有活躍 change → 用當前 change；(5) 無 argument + 無活躍 change → 自動偵測；(6) 都找不到 → 新建
- D-07: 三種支援的文件路徑（自動偵測來源）：使用者自定位置、`~/.gstack/projects/{project}/`、對話上下文中提到的計畫文件。`.claude/plans/` 因 hash 檔名不納入

**`/mysd:fix` 指令設計**
- D-08: 路徑自動偵測：fix 讀取 task sidecar 失敗記錄 + 檢查 worktree 狀態自動判斷路徑
- D-09: Argument 格式：`/mysd:fix T2`（當前 change）或 `/mysd:fix {change-name} T2`（指定 change）
- D-10: 無 argument 時：列出當前 change 中 failed/blocked 的 tasks，讓使用者選擇
- D-11: Research 觸發：research 選項只在實作問題路徑出現
- D-12: 實作問題診斷：AI 自動讀 task sidecar 失敗原因、AI 嘗試解法、test 輸出，自動診斷並說明後執行修復
- D-13: 純 SKILL.md + agent 實作，不需新 Go binary subcommand
- D-14: 兩條路徑：merge 衝突路徑（解衝突 → build + test 驗證 → merge → `-D` 強制刪 branch）；實作問題路徑（修正 task 內容 + 更新 spec → `-D` 刪舊 branch → spawn executor 重新執行）；放棄路徑（task 回 `pending`，`-D` 刪 branch + worktree）；fix 成功 merge 後，下游 skipped tasks 遞移性恢復為 `pending`

**Agent 重構（FAGENT-05/06/07）**
- D-15: FAGENT-07（executor per task）：無論 single 或 wave mode，每個 task 都是獨立的新 agent 實例；single mode 改為序列 spawn（一個接一個），wave mode 維持並行 spawn
- D-16: FAGENT-06（spec-writer per spec）：每個 spec 檔案一個新 `mysd-spec-writer` agent 實例
- D-17: FAGENT-05（agent audit）：執行時人工審計，確認所有 9 個 agent definitions 無 Task tool 呼叫；Task tool 只允許在頂層 SKILL.md orchestrator 層使用

**`--auto` Flag 傳遞機制**
- D-18: SKILL.md 層解析 `$ARGUMENTS` 中的 `--auto` flag，spawn agent 時將 `auto_mode: true` 放入 context JSON；Go binary 不需新 flag
- D-19: ff/ffe 隱含 auto mode — agent 在每個決策點自主選擇最佳方案，不詢問使用者
- D-20: ff/ffe 不使用 research（FAUTO-04）；單獨加 `--auto` 的指令仍可有 research，只是跳過互動確認

**工作流程調整**
- D-21: 新工作流程：`propose/discuss` → `discuss`（選擇性）→ `plan` → `apply` → `archive`
- D-22: `execute` 改名為 `apply`（SKILL.md 指令 `/mysd:execute` → `/mysd:apply`，agent `mysd-executor` 仍保持原名）
- D-23: `/mysd:plan` 內部流程改寫：[選擇性] spawn `mysd-researcher` × 4 並行；spawn `mysd-designer` → 產出 `design.md`；spawn `mysd-planner` → 讀 spec + design.md → 產出 `tasks.md`
- D-24: `/mysd:ff` 改寫 = plan（無 research）+ apply + archive
- D-25: `/mysd:ffe` 改寫 = plan（4 維度 research）+ apply + archive
- D-26: design.md 由 plan 階段的 `mysd-designer` sub-agent 產生，不再是獨立 SKILL.md 步驟

**`/mysd:propose` 自動偵測輸入來源**
- D-27: source 偵測邏輯與 discuss 完全相同（D-06/D-07）
- D-28: `.planning/phases/{phase}/` 路徑移除

**`/mysd:status` SKILL.md 設計**
- D-29: 顯示新版 workflow stage（`propose` → `plan` → `apply` → `archive`），標示目前所在位置
- D-30: Task 列表含編號、title、狀態符號（✓ done / ✗ failed / ⊘ skipped / ○ pending）
- D-31: 最後一行：`Next: /mysd:{command}` 推薦下一步指令

### Claude's Discretion

- 4 個新 agent definitions 的 allowed-tools 列表細節
- `mysd-researcher` / `mysd-advisor` / `mysd-proposal-writer` 的 prompt 設計細節
- `/mysd:fix` SKILL.md 中偵測 conflict markers 的具體方式
- `mysd-advisor` 的比較表格格式設計
- apply（原 execute）命令的 state transition binary call 名稱（`mysd apply` 或保留 `mysd execute`）

### Deferred Ideas (OUT OF SCOPE)

- `/mysd:propose` 與 `/mysd:discuss` 的合併討論（保持獨立）
- `design` 步驟是否完全消失（`/mysd:design` SKILL.md 是否保留待評估）
- Phase 9：propose/discuss 完整互動式雙層循環（DISC-01~09）
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| FCMD-01 | `/mysd:discuss` 隨時補充討論，支援 4 維度並行 research | D-01 to D-07 決定 discuss orchestrator 的完整行為；`mysd-researcher` agent 需新建（FAGENT-01）|
| FCMD-02 | `/mysd:fix` 互動式修復，可選 research，spawn executor subagent（worktree 隔離），只改 code | D-08 to D-14 決定 fix orchestrator 邏輯；現有 `mysd execute --context-only`、`mysd task-update`、`mysd worktree` binary subcommands 已足夠 |
| FAGENT-01 | 新增 `mysd-researcher` agent definition（研究 codebase/domain） | 需新建 `.claude/agents/mysd-researcher.md` 和 `plugin/agents/mysd-researcher.md`，prompt 設計在 Claude's Discretion |
| FAGENT-02 | 新增 `mysd-advisor` agent definition（gray area 分析，帶比較表） | 需新建 `.claude/agents/mysd-advisor.md` 和 `plugin/agents/mysd-advisor.md` |
| FAGENT-03 | 新增 `mysd-proposal-writer` agent definition（寫 proposal.md） | 需新建 `.claude/agents/mysd-proposal-writer.md` 和 `plugin/agents/mysd-proposal-writer.md` |
| FAGENT-05 | 所有 agent definitions 確認無 Task tool 呼叫 | 9 個 agent 需人工 audit（verifier、scanner、uat-guide、planner、designer、executor、fast-forward、spec-writer、plan-checker）；audit 結果是人工動作 + 文件更新 |
| FAGENT-06 | `mysd-spec-writer` 改為 per capability area spawn | 改寫 `mysd-spec-writer.md`：移除一次性處理所有 specs 的邏輯，改為接收單一 capability area 並輸出一個 spec 檔案 |
| FAGENT-07 | `mysd-executor` 改為 per task spawn | 改寫 `mysd-executor.md`：移除 sequential 模式中「迴圈處理所有 pending tasks」邏輯；改為每次只處理 `assigned_task`；SKILL.md apply 改為序列 spawn |
| FAUTO-01 | `--auto` flag 支援 propose/spec/discuss/plan | SKILL.md 層解析 `$ARGUMENTS`，注入 `auto_mode: true` 到 context JSON（D-18）；不需 binary 改動 |
| FAUTO-02 | `--auto` 跳過互動提問，自動選推薦方案 | agent 收到 `auto_mode: true` 時跳過所有確認步驟 |
| FAUTO-03 | ff/ffe 隱含 `--auto` | `/mysd:ff` 和 `/mysd:ffe` 改寫時固定傳入 `auto_mode: true` |
| FAUTO-04 | ff/ffe 不使用 research，直接用 subagent 依照既有 spec 內容完成 | ff/ffe 改寫時不包含 research 步驟（D-20） |
</phase_requirements>

---

## Summary

Phase 8 是一個純 plugin layer 的改寫——不涉及 Go binary 的新 subcommand（D-13），所有工作集中在 `.claude/commands/`（SKILL.md orchestrators）和 `.claude/agents/`（agent definitions）這兩個目錄，以及對應的 `plugin/commands/` 和 `plugin/agents/`。

核心工作分四類：

1. **新 agent definitions（3 個）**：`mysd-researcher`、`mysd-advisor`、`mysd-proposal-writer`。這三個 agent 在 Phase 5 的 model profile 表中已有對應（research/advisory 角色），但 `.md` 檔案尚未建立。
2. **新 SKILL.md 指令（2 個）**：`/mysd:discuss` 和 `/mysd:fix`，以及對應的 plugin 版本。
3. **現有 SKILL.md 改寫（6 個）**：`plan`、`ff`、`ffe`、`execute → apply`、`propose`、`status`。
4. **現有 agent 改寫/審計（9 個）**：`executor`（per-task）、`spec-writer`（per-spec）；7 個 agent 需要確認無 Task tool 呼叫。

最重要的架構原則（已在 STATE.md 確立）：**subagent 不 spawn subagent**——Task tool 只允許在 SKILL.md orchestrator 層使用。這是 Phase 8 agent audit 的核心。

**Primary recommendation:** SKILL.md 層負責所有 Task tool 呼叫；agent 層只做實際工作（讀、寫、分析），不呼叫其他 agent。

---

## Project Constraints (from CLAUDE.md)

- **Tech stack**: Go 單一 binary + Claude Code plugin layer。Phase 8 完全在 plugin layer，不改 Go binary。
- **Plugin 形式**: Claude Code slash commands + agent definitions，按 SKILL.md 格式
- **Convention over configuration**: 預設即好用，`auto_mode: true` 由 SKILL.md 注入，不需 binary flag
- **相容性**: 必須能讀寫 OpenSpec 格式。討論結論更新 spec 後，binary 的 spec 解析邏輯不受影響
- **GSD Workflow Enforcement**: 必須透過 GSD 指令進行，不直接修改 repo

---

## Standard Stack

### Core (SKILL.md Layer)

| Element | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| SKILL.md orchestrator pattern | - | 薄 orchestrator + Task tool spawn agents | Phase 2 決策，STATE.md 確立 |
| Task tool | Claude Code built-in | 在 SKILL.md 層 spawn agent 執行工作 | 唯一允許呼叫 subagent 的工具 |
| Context JSON pattern | - | `mysd {cmd} --context-only` 輸出 JSON → SKILL.md 解析後傳給 agent | v1.0 建立的 binary-as-state-manager 模式 |
| `auto_mode` context field | - | 在 context JSON 中傳遞 auto 模式旗標 | ExecutionContext 已有此欄位（`internal/executor/context.go:28`）|

### Existing Binary Subcommands (Phase 8 Fix 所需)

| Command | Purpose | Phase 8 Usage |
|---------|---------|---------------|
| `mysd execute --context-only` | 輸出執行 context JSON（含 failed tasks 資訊）| fix 讀取 failed tasks 狀態 |
| `mysd task-update {id} {status}` | 更新 task 狀態（pending/in_progress/done/blocked）| fix 恢復 task 為 pending；skipped tasks 恢復 |
| `mysd worktree remove {id} {branch}` | 刪除 worktree | fix 清理失敗的 worktree |
| `mysd plan --context-only --check` | 輸出 coverage 結果 | discuss 結論更新後的 re-plan + plan-checker |

**Version verification:** 以上 subcommands 均已存在於 `cmd/` 下，Phase 8 不需新增。

### Plugin File Structure

```
.claude/
  commands/           ← 使用者直接執行的 SKILL.md（主要）
    mysd-discuss.md   ← NEW
    mysd-fix.md       ← NEW
    mysd-apply.md     ← RENAMED from mysd-execute.md
    mysd-plan.md      ← REWRITE
    mysd-ff.md        ← REWRITE
    mysd-ffe.md       ← REWRITE
    mysd-propose.md   ← REWRITE（加 source 偵測）
    mysd-status.md    ← REWRITE
  agents/             ← Task tool spawn 的 agent definitions（主要）
    mysd-researcher.md    ← NEW
    mysd-advisor.md       ← NEW
    mysd-proposal-writer.md ← NEW
    mysd-executor.md      ← REWRITE（per-task）
    mysd-spec-writer.md   ← REWRITE（per-spec）
    mysd-planner.md       ← AUDIT only
    mysd-designer.md      ← AUDIT only
    mysd-verifier.md      ← AUDIT only
    mysd-scanner.md       ← AUDIT only
    mysd-uat-guide.md     ← AUDIT only
    mysd-fast-forward.md  ← AUDIT only
    mysd-plan-checker.md  ← 存在於 plugin/agents/，需確認位置

plugin/
  commands/           ← plugin distribution 版本（需與 .claude/commands/ 同步）
  agents/             ← plugin distribution 版本（需與 .claude/agents/ 同步）
```

**重要發現：** `mysd-plan-checker.md` 目前只在 `plugin/agents/` 存在，**不在** `.claude/agents/`。這是 Phase 5 的產物。Phase 8 需要確認此 agent 是否需要出現在 `.claude/agents/` 以供 SKILL.md orchestrator 呼叫。

---

## Architecture Patterns

### Pattern 1: SKILL.md Orchestrator（固定模式，不可更改）

**What:** SKILL.md 是薄層 orchestrator：讀 context → 做 minimal 決策 → 用 Task tool spawn agent → 展示結果
**When to use:** 所有 SKILL.md 檔案都遵循此模式
**Key constraint:** SKILL.md 可以 spawn agents，但 agents **絕不** spawn agents

```markdown
## Step N: [Something]

Run:
```
mysd {cmd} --context-only
```

Use the Task tool to invoke `mysd-{agent}`:
```
Task: [description]
Agent: mysd-{agent}
Context: {JSON from previous step}
```
```

### Pattern 2: Context JSON 傳遞

**What:** 所有 agent input 都是 context JSON，包含執行所需的所有資訊
**Source:** `mysd {cmd} --context-only` 輸出，或 SKILL.md 手動構建
**Pattern for `auto_mode`:**

```markdown
## Step 1: Parse Arguments

Check `$ARGUMENTS`:
- If `--auto` is present: set `auto_mode: true` in context JSON
- Otherwise: set `auto_mode: false`
```

然後在構建 context JSON 時加入 `auto_mode` 欄位。

### Pattern 3: Per-Unit Spawn（Phase 8 的核心改變）

**舊模式（FAGENT-06/07 被廢棄）:**
- executor: 一個 agent 處理所有 pending tasks（sequential loop）
- spec-writer: 一個 agent 處理所有 capability areas

**新模式（Phase 8 後）:**
- executor: 每個 task 一個新 agent 實例（sequential spawn = 一個接一個）
- spec-writer: 每個 capability area 一個新 agent 實例

**SKILL.md orchestrator 如何做 sequential spawn:**
```markdown
For each task in pending_tasks:
  Use the Task tool to invoke mysd-executor:
    Agent: mysd-executor
    Context: {... assigned_task: task ...}
  Wait for completion before spawning next.
```

### Pattern 4: `--auto` Flag 傳遞

**Flow:**
1. SKILL.md 解析 `$ARGUMENTS` 中的 `--auto`
2. 構建 context JSON 時加入 `"auto_mode": true`
3. Agent 收到 context 後，遇到互動決策點時檢查 `auto_mode`：true → 自動選推薦；false → 詢問使用者

**ff/ffe 的處理：** ff/ffe SKILL.md 不解析 `--auto`，直接在 context JSON 中硬寫 `"auto_mode": true`（因為 ff/ffe 本身就隱含 auto mode）。

### Pattern 5: discuss 的 Spec 更新後 Re-plan

**Flow（D-05）:**
1. 結論達成 → 確認要納入 spec
2. SKILL.md 判斷影響的 spec 層級
3. Spawn 對應 agent 更新 spec（proposal-writer / spec-writer / designer）
4. 更新完成後，呼叫 `mysd plan --context-only` 取得新 planning context
5. Spawn `mysd-planner` 重新 plan
6. 呼叫 `mysd plan --check --context-only` 取得 coverage
7. Spawn `mysd-plan-checker` 驗證覆蓋率

### Pattern 6: fix 的兩條路徑

**Merge 衝突路徑（conflict markers detected）:**
1. 讀取 worktree path（從 task 狀態/worktree list）
2. 在 worktree 中解衝突（Edit 工具移除 conflict markers）
3. `git add` + `git commit`（完成 merge）
4. `go build ./...` + `go test ./...` 驗證
5. `mysd worktree remove {id} {branch}`（clean up）
6. Spawn downstream skipped tasks recovery（`mysd task-update {id} pending`）

**實作問題路徑（implementation failure）:**
1. 讀 task sidecar 失敗原因
2. 可選：spawn `mysd-researcher` 研究
3. 修正 task description + 更新相關 spec
4. `git branch -D {branch}` 刪舊 branch
5. `mysd worktree remove {id} {branch}` 清理
6. Spawn `mysd-executor` 重新執行 task（worktree isolation）

### Anti-Patterns to Avoid

- **Agent spawns Agent:** 絕對禁止。Verifier、planner、spec-writer 不可有 Task tool 呼叫。
- **SKILL.md 直接做實作工作:** SKILL.md 只做 orchestration，實際工作交 agent。
- **ff/ffe 包含 research:** FAUTO-04 明確禁止。ff/ffe 假設 spec 已就緒。
- **apply 改名影響 binary:** D-22 只改 SKILL.md 檔名，binary subcommand `mysd execute` 保持不變（或由 Claude's Discretion 決定）。

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Conflict marker detection | 自己 grep `<<<<<<<` 符號 | SKILL.md 讓 Claude 自然讀取文件並識別衝突 | Claude 本身能識別 conflict markers，無需 binary 輔助 |
| Task failure state lookup | 新 binary subcommand | `mysd execute --context-only` | 已有 ExecutionContext 含 task status |
| Skipped task list | 手動追蹤 | `mysd execute --context-only` pending_tasks 過濾 | binary 已計算 |
| Source detection logic | 複雜自定義邏輯 | SKILL.md 條件判斷（D-06 有明確 6 步優先順序）| 邏輯在 SKILL.md 層即可，不需 binary |

**Key insight:** Phase 8 的所有工作都在 SKILL.md 和 agent `.md` 文字層，Go binary 不需任何改動（D-13）。

---

## Runtime State Inventory

> Phase 8 是純 plugin layer 改寫，無 rename/refactor scope，此 section 不適用。

**SKIPPED** — Phase 8 是 SKILL.md 和 agent definition 的新建/改寫，不涉及任何 rename 或 migration。

---

## Environment Availability

> Step 2.6: Phase 8 是純文字檔案改寫（Markdown），無外部 CLI/service 依賴。

**SKIPPED** — Phase 8 只操作 `.claude/commands/*.md`、`.claude/agents/*.md`、`plugin/commands/*.md`、`plugin/agents/*.md` 等文字檔。唯一的 runtime 依賴是 `mysd` binary，已存在（`mysd.exe` 在 repo root）。

---

## Common Pitfalls

### Pitfall 1: Agent 呼叫 Task Tool（違反 FAGENT-05）
**What goes wrong:** agent definition `.md` 的 allowed-tools 包含 Task，或 agent prompt 描述中有呼叫其他 agent 的指示
**Why it happens:** 複製既有 SKILL.md 的模式時不小心帶入 Task tool
**How to avoid:** agent audit 時明確檢查 allowed-tools frontmatter 是否包含 Task；prompt body 是否描述 spawn subagent 行為
**Warning signs:** allowed-tools 列表中有 "Task"；prompt 中有 "Use the Task tool to invoke" 字樣

### Pitfall 2: .claude/agents/ 和 plugin/agents/ 不同步
**What goes wrong:** 更新 `.claude/agents/mysd-executor.md` 但忘記同步 `plugin/agents/mysd-executor.md`，或反之
**Why it happens:** 兩個目錄平行存在，容易遺漏其中一個
**How to avoid:** 每次更新 agent/command 定義，兩個目錄同步更新；plan 中明確列出兩個路徑
**Warning signs:** `.claude/agents/` 和 `plugin/agents/` 的 md5 不同

### Pitfall 3: mysd-executor 仍有 sequential loop（FAGENT-07 未完成）
**What goes wrong:** executor 改寫後 prompt 中還有「For each task in pending_tasks...」的迴圈描述
**Why it happens:** executor 原本的設計就是 sequential loop，改寫時容易遺漏刪除
**How to avoid:** 改寫後確認 executor `.md` 只處理 `assigned_task`（singular），沒有 `pending_tasks` 的迴圈
**Warning signs:** prompt 中有 `for each` 或 `pending_tasks` 陣列的迭代描述

### Pitfall 4: apply SKILL.md 仍呼叫 `mysd execute`
**What goes wrong:** 改名為 `mysd-apply.md` 後，內部仍呼叫 `mysd execute --context-only`
**Why it happens:** D-22 說「agent `mysd-executor` 仍保持原名」但沒說 binary subcommand 如何處理；Claude's Discretion 尚未決定
**How to avoid:** 確認 binary subcommand 名稱決策（apply 或 execute）後再更新 SKILL.md；暫時保留 `mysd execute --context-only` 呼叫是安全的（binary 不變）
**Warning signs:** 出現 `mysd apply --context-only` 但 binary 沒有對應 subcommand

### Pitfall 5: discuss re-plan 忘記觸發 plan-checker
**What goes wrong:** discuss 結論更新 spec 後只 spawn `mysd-planner`，沒有後續的 plan-checker 驗證
**Why it happens:** D-05 列出了完整流程（更新 spec → re-plan → plan-checker），但容易只實作前半段
**How to avoid:** discuss SKILL.md 的 Step 5（spec 更新後）需要明確包含 `mysd plan --check --context-only` → spawn `mysd-plan-checker` 兩步
**Warning signs:** discuss SKILL.md 的 spec 更新流程末尾沒有 plan-checker 呼叫

### Pitfall 6: fix 的 skipped tasks 恢復遺漏遞移性
**What goes wrong:** fix 成功 merge 後只恢復直接依賴此 task 的下游 tasks，沒有遞移性恢復（T1 失敗 → T3 depends T2 depends T1，T3 也要恢復）
**Why it happens:** 簡單實作只看第一層 depends
**How to avoid:** 使用 `mysd execute --context-only` 的 wave_groups 資訊（binary 已做 topological sort），識別所有被此 task 影響的下游 tasks
**Warning signs:** 恢復後只有直接 depends 的 tasks 變 pending，更深層的 skipped tasks 沒有恢復

### Pitfall 7: mysd-plan-checker.md 位置問題
**What goes wrong:** `plugin/agents/mysd-plan-checker.md` 存在但 `.claude/agents/mysd-plan-checker.md` 不存在，導致 SKILL.md 呼叫時找不到 agent
**Why it happens:** Phase 5 在 `plugin/agents/` 建立了 plan-checker，但 `.claude/agents/` 沒有同步
**How to avoid:** Phase 8 開始時確認 `.claude/agents/` 目錄下是否有 `mysd-plan-checker.md`；若無，複製過去
**Warning signs:** 執行 `/mysd:plan --check` 時 Task tool 找不到 `mysd-plan-checker` agent

---

## Code Examples

### 現有 SKILL.md Orchestrator 模式（plan.md plugin 版本）

```markdown
# Source: plugin/commands/mysd-plan.md

## Step 3: Invoke Planner Agent

Use the Task tool to invoke the `mysd-planner` agent with the full context JSON:

Task: Invoke mysd-planner agent
Agent: mysd-planner
Context: {context JSON from Step 2}
```

### 現有 per-task spawn 模式（execute.md plugin 版本，wave mode）

```markdown
# Source: plugin/commands/mysd-execute.md — Step 3B-3

Use the Task tool to spawn ONE executor agent PER task in the current wave (all in parallel):

For each task in the wave:
  Task: Invoke mysd-executor agent for wave task T{task.id}
  Agent: mysd-executor
  Context: {
    ...
    execution_mode: "wave",
    assigned_task: {task},
    worktree_path: {path},
    branch: {branch},
    isolation: "worktree"
  }
```

### 現有 `auto_mode` 讀取模式（mysd-planner.md）

```markdown
# Source: .claude/agents/mysd-planner.md — Step 7.5

Check the `auto_mode` flag in the input context (set to true when running in ffe mode).

If `auto_mode` is false (interactive mode):
  Present task-skills mapping table to the user...

If `auto_mode` is true (ffe mode, per D-10):
  Skip confirmation entirely. Use the recommended skills as-is.
```

### 現有 agent allowed-tools 格式（合規範例）

```markdown
# Source: plugin/agents/mysd-plan-checker.md（合規 — 無 Task tool）

---
description: Plan-checker agent. ...
allowed-tools:
  - Read
  - Write
  - Edit
  - Glob
  - Grep
---
```

**注意：** `mysd-executor.md`、`mysd-verifier.md` 等所有 agent 都有 Bash 工具但**沒有 Task**。

### ExecutionContext.auto_mode 欄位（binary）

```go
// Source: internal/executor/context.go:27-28
// Wave grouping fields (Phase 06 extension — additive only per D-11)
AutoMode       bool         `json:"auto_mode,omitempty"`
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| executor 處理所有 pending tasks（sequential loop）| per-task spawn（每個 task 獨立 agent）| Phase 8（D-15）| SKILL.md apply 需要改寫為序列 spawn loop |
| spec-writer 處理所有 capability areas | per-spec-file spawn | Phase 8（D-16）| SKILL.md spec 需要改寫為逐 spec 串行 spawn |
| fast-forward 是單一 agent 做全部 pipeline | 分解為 SKILL.md orchestrator 呼叫 designer + planner | Phase 8（D-23/24/25）| mysd-fast-forward agent 角色改變或廢棄 |
| `/mysd:execute` | `/mysd:apply`（SKILL.md 改名）| Phase 8（D-22）| 舊檔案刪除，新檔案建立；binary 不動 |
| design 是獨立的 SKILL.md 步驟 | design.md 由 plan 內部 `mysd-designer` sub-agent 產生 | Phase 8（D-26）| `/mysd:plan` orchestrator 需 spawn designer |

**Deprecated/outdated:**
- `mysd-execute.md`（.claude/commands/ 和 plugin/commands/）: 被 `mysd-apply.md` 取代（D-22）
- `mysd-fast-forward` agent 的「Phase 4: Execute Tasks」邏輯: ff/ffe 改寫後由 apply SKILL.md 負責執行（D-24/25）

---

## Open Questions

1. **`mysd apply` vs `mysd execute` binary subcommand**
   - What we know: D-22 說 SKILL.md 層改名（execute → apply），agent `mysd-executor` 保持原名
   - What's unclear: apply SKILL.md 內部呼叫的 binary subcommand 名稱是 `mysd execute --context-only` 還是 `mysd apply --context-only`？
   - Recommendation: Claude's Discretion 決定。保守做法是 apply SKILL.md 仍呼叫 `mysd execute --context-only`（binary 不變）；進取做法是新增 `mysd apply` 作為 `mysd execute` 的 alias。建議保守做法，Phase 8 不改 binary。

2. **mysd-plan-checker.md 在 .claude/agents/ 的位置**
   - What we know: `plugin/agents/mysd-plan-checker.md` 存在（Phase 5 成果），`.claude/agents/` 目錄下只有 8 個 agent，沒有 `mysd-plan-checker.md`
   - What's unclear: Phase 5 計畫只建立了 plugin 版本還是兩個都建了？
   - Recommendation: Phase 8 開工前先確認 `.claude/agents/mysd-plan-checker.md` 是否存在；若不存在，Phase 8 第一個 plan 需要補建。

3. **mysd-fast-forward agent 的命運**
   - What we know: D-24/25 改寫了 ff/ffe 的流程（plan + apply + archive），不再呼叫 `mysd-fast-forward` agent
   - What's unclear: `mysd-fast-forward.md` agent 是保留（供其他用途）還是廢棄？
   - Recommendation: Phase 8 的 ff/ffe 改寫後，`mysd-fast-forward` agent 實際上不再被任何 SKILL.md 呼叫。保留為 dead code 或在 Phase 8 結束後移除。建議保留但不主動更新（不是 Phase 8 的 scope）。

4. **discuss 新建 change 的邏輯（D-06 的 Case 6）**
   - What we know: 當 discuss 找不到任何來源時，走「新建」流程（類似 propose）
   - What's unclear: 新建流程是直接呼叫 `mysd propose {name}` binary 還是 spawn `mysd-proposal-writer` agent？
   - Recommendation: D-05 描述了結論納入後 spawn `mysd-proposal-writer`，所以「新建」應該也走相同路徑：先 `mysd propose {name}` 建立 change 目錄，再 spawn `mysd-proposal-writer` 填寫內容。

---

## Validation Architecture

> nyquist_validation 在 .planning/config.json 中為 true，此 section 必須包含。

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go testing (stdlib) + testify v1 |
| Config file | none (go test ./... 直接執行) |
| Quick run command | `go test ./... -count=1` |
| Full suite command | `go test ./... -count=1 -race` |

### Phase 8 Validation Notes

Phase 8 的工作全部是 Markdown `.md` 文字檔（SKILL.md 和 agent definitions），**沒有 Go 程式碼改動**，因此：

- Go 測試框架對 Phase 8 **不適用**（無可測試的 Go 程式碼）
- Phase 8 的 validation 是人工審計（FAGENT-05）和功能測試（在 Claude Code 中執行指令）
- Phase 8 成功標準（ROADMAP.md）是可執行的功能測試，不是自動化測試

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FCMD-01 | `/mysd:discuss` 進入討論流程 | manual | 在 Claude Code 執行 `/mysd:discuss` | ❌ Wave 0（無自動化測試）|
| FCMD-02 | `/mysd:fix` 修復失敗 task | manual | 在 Claude Code 執行 `/mysd:fix` | ❌ Wave 0（無自動化測試）|
| FAGENT-01 | `mysd-researcher.md` 存在且格式正確 | file-check | `ls .claude/agents/mysd-researcher.md` | ❌ |
| FAGENT-02 | `mysd-advisor.md` 存在且格式正確 | file-check | `ls .claude/agents/mysd-advisor.md` | ❌ |
| FAGENT-03 | `mysd-proposal-writer.md` 存在且格式正確 | file-check | `ls .claude/agents/mysd-proposal-writer.md` | ❌ |
| FAGENT-05 | 所有 9 agent 無 Task tool | audit | `grep -r "Task" .claude/agents/ --include="*.md"` | N/A（audit）|
| FAGENT-06 | spec-writer per-spec | manual | `/mysd:spec` 觀察 spawn 次數 | ❌ |
| FAGENT-07 | executor per-task | manual | `/mysd:apply` 觀察 spawn 次數 | ❌ |
| FAUTO-01/02 | `--auto` 跳過互動 | manual | `/mysd:plan --auto` | ❌ |
| FAUTO-03/04 | ff/ffe 隱含 auto 且無 research | manual | `/mysd:ff` 觀察行為 | ❌ |

### Wave 0 Gaps

- [ ] `.claude/agents/mysd-plan-checker.md` — 確認是否需要從 `plugin/agents/` 複製
- [ ] FAGENT-05 audit 腳本：`grep -rn "Task" .claude/agents/ --include="*.md" | grep -v "^Binary\|description\|#\|Task tool"` — 用於初步篩選可疑行

**Note:** Phase 8 沒有 Go 程式碼改動，Wave 0 的「test infrastructure gaps」主要是確認 `.claude/agents/` 目錄的完整性，以及建立一個 agent audit checklist。

---

## Sources

### Primary (HIGH confidence)

- `.claude/commands/mysd-execute.md`（.claude 版，Phase 4 原始）— 現有 execute orchestrator 模式
- `plugin/commands/mysd-execute.md`（Phase 6 改寫版）— 最新 wave mode 模式，含 per-task spawn
- `.claude/agents/mysd-executor.md` — 現有 executor agent（需改寫為 per-task）
- `.claude/agents/mysd-planner.md` — 現有 planner（含 auto_mode 讀取模式，Step 7.5）
- `plugin/agents/mysd-plan-checker.md` — Phase 5 建立的 plan-checker（leaf agent 模式，無 Task tool）
- `internal/executor/context.go` — ExecutionContext struct，含 auto_mode 欄位（line 27-28）
- `internal/config/defaults.go` — ProjectConfig struct（auto_mode 欄位定義）
- `.planning/phases/08-skill-md-orchestrators-agent-definitions/08-CONTEXT.md` — Phase 8 所有設計決策（D-01 to D-31）

### Secondary (MEDIUM confidence)

- `.planning/REQUIREMENTS.md` — Phase 8 需求（FCMD-01/02, FAGENT-01~07, FAUTO-01~04）— 官方需求文件
- `.planning/STATE.md` — 「subagent 不 spawn subagent」原則確立記錄
- `.specs/changes/interactive-discovery/proposal.md` — v1.1 完整功能規格，subagent 架構表格

### Tertiary (LOW confidence)

- 無

---

## Metadata

**Confidence breakdown:**
- Standard Stack: HIGH — 完全基於現有 codebase 的直接觀察
- Architecture Patterns: HIGH — 基於現有 SKILL.md 和 agent `.md` 檔案的直接閱讀
- Pitfalls: HIGH — 基於 CONTEXT.md 的 Decisions 和 STATE.md 的已知 blockers

**Research date:** 2026-03-26
**Valid until:** 2026-04-25（plugin layer 設計穩定，有效期 30 天）
