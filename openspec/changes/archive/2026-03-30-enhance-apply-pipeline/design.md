---
spec-version: "1"
change: enhance-apply-pipeline
status: designed
---

## Context

mysd apply 目前有兩種執行模式：single（per-task sequential）和 wave（per-wave parallel）。Executor agent 收到 task 後經過 alignment gate 直接執行，沒有執行前的 reuse/quality/efficiency 檢查。遇到問題時的暫停行為未明確定義。此外，執行前沒有 preflight 機制來偵測 missing files 或 stale artifacts。

現有結構：
- `cmd/execute.go` — execute CLI，已有 `--context-only`、`--spec` flags
- `internal/executor/context.go` — `ExecutionContext` struct，含 `ExecutionMode`（"single"/"wave"）
- `internal/config/config.go` — `DefaultModelMap`，executor role 用 sonnet（quality/balanced）或 haiku（budget）
- `mysd/skills/apply/SKILL.md` — apply orchestrator，Step 3 依 mode 分流
- `mysd/agents/mysd-executor.md` — executor agent，已有 alignment gate + worktree isolation

## Goals / Non-Goals

**Goals:**

- Per-spec 執行模式：agent 在同一 spec 的多個 tasks 間保持 context 連續性
- Executor 寫 code 前有結構化的品質檢查
- 明確的暫停條件避免 executor 在不確定時猜測
- 執行前 preflight 偵測 missing files 和 stale artifacts

**Non-Goals:**

- 不改 alignment gate 結構 — 它已經做得好，只在前面加 pre-task checks
- 不改 wave mode — wave 和 spec mode 是獨立的選擇
- 不做 drift detection（語義層面的檔案變更偵測）— 成本太高

## Decisions

### D-01: Per-spec execution mode

新增 `execution_mode: "spec"`。行為：

1. Orchestrator 讀取 pending tasks，按 `spec` field 分組
2. 每個 spec group 各 spawn 一個 executor agent
3. Agent 收到該 spec 的 **所有 tasks**（不只一個），依序執行
4. Agent 在 tasks 間保持 context 連續性（不會重新 spawn）

Orchestrator 層面：apply SKILL.md Step 3 新增 `spec` mode 分支，與 single/wave 平行。

```
execution_mode:
  "single" → per-task, sequential
  "wave"   → per-wave, parallel with worktree
  "spec"   → per-spec, one agent per spec
```

Per-spec 不使用 worktree — agent 直接在 repo root 執行（同 single mode）。原因：per-spec 是 sequential 的（一次一個 spec agent），不需要 isolation。

如果使用者同時指定 `--spec X`（來自上一個 change），orchestrator 只需 spawn 一個 agent。

**替代方案**：per-spec + parallel（多個 spec agents 同時跑，各用 worktree）— 太複雜，留待未來需求。

### D-02: `spec-executor` model role

在 `DefaultModelMap` 新增 `spec-executor` role：

| Profile | executor | spec-executor |
|---------|----------|---------------|
| quality | sonnet | opus |
| balanced | sonnet | opus |
| budget | haiku | sonnet |

Apply SKILL.md 根據 `execution_mode` 選擇 role：
- `single` / `wave` → `ResolveModel("executor", ...)`
- `spec` → `ResolveModel("spec-executor", ...)`

Orchestrator 在 Step 2 取得 context 後，根據 mode 決定用哪個 role 查 model。

### D-03: Executor pre-task checks

在 `mysd-executor.md` 的 Task Execution section，Step 1（Mark In Progress）之後、Step 2（TDD）之前，新增 **Step 1b: Pre-Task Checks**。

4 項檢查：

1. **Reuse** — 搜索 adjacent modules 和 shared utilities，確認沒有既有實作可重用
2. **Quality** — 確認使用既有 types 和 constants，不重複定義
3. **Efficiency** — 確認 async operations 被正確平行化，scope 匹配需求
4. **No Placeholders** — 讀取此 task 對應的 spec/design section，確認沒有 TBD/TODO/vague language。如果發現 placeholder，暫停並通報，不實作模糊需求

Per-spec mode 的 agent 在處理每個 task 前都執行這 4 項。

### D-04: Executor pause conditions

在 `mysd-executor.md` 新增明確的 Pause Conditions section。Agent 遇到以下 4 種情況時 **必須暫停並通報**，不猜測：

1. **Task unclear** — task description 模糊或矛盾，無法確定預期行為
2. **Design issue discovered** — 實作過程中發現 design 有遺漏或矛盾
3. **Error/blocker** — build 失敗、dependency 問題、或其他技術障礙（現有 On Failure 路徑處理）
4. **User interrupt** — 使用者中斷執行

暫停時 agent 必須：輸出問題描述 + 建議的解決選項 + 等待指引。

### D-05: Preflight check CLI

`cmd/execute.go` 新增 `--preflight` flag。行為：

```json
{
  "status": "ok|warning|critical",
  "checks": {
    "missing_files": ["internal/foo.go"],
    "staleness": {
      "days_since_last_plan": 14,
      "is_stale": true
    }
  }
}
```

**Missing files check**：讀取 tasks.md 中所有 task 的 `files` 欄位，檢查每個檔案是否存在（新檔案除外 — 如果 task description 包含 "create" 或 "add" 則跳過）。

**Staleness check**：比較 STATE.json 的 `last_run` 和現在時間。超過 7 天為 `warning`，超過 30 天為 `critical`。

Apply SKILL.md 在 Step 2 後（Step 2b 之後）新增 Step 2c: Preflight Check。呼叫 `mysd execute --preflight --json`，有 warning/critical 時顯示並請使用者確認。

### D-06: Spec-executor agent reuse

Per-spec mode 不需要新的 agent definition。`mysd-executor.md` 已有 multi-task 的能力基礎（alignment gate 讀全部 specs），只需修改 apply SKILL.md 的 context 傳遞方式：

- Single mode：`assigned_task` = 單一 task
- Spec mode：`assigned_tasks` = 該 spec 的所有 tasks（array），agent 依序執行

Executor agent 需判斷 `assigned_tasks`（array）vs `assigned_task`（single）來決定行為。

## Risks / Trade-offs

- [Context 壓縮] Per-spec agent 處理多個 tasks 時，後期 tasks 可能因 context 壓縮失真 → Mitigation: 每個 task 前強制 re-read spec/design；用 opus 提供更好的長 context 理解
- [Preflight 誤報] 檔案不存在但 task 是要新建它 → Mitigation: 檢查 task description 包含 "create"/"add" 關鍵字時跳過
- [Staleness 閾值爭議] 7 天 warning 可能太短 → Mitigation: 可配置閾值，但先 hardcode 合理預設值
