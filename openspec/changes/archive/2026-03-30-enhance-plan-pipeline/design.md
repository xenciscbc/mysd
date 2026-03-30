---
spec-version: "1"
change: enhance-plan-pipeline
status: designed
---

## Context

mysd plan 目前的 pipeline 是：research? → design → planner → reviewer → analyze-fix → plan-checker。Agent 收到的是一包 context JSON，自由發揮產出 artifact，沒有結構化的 template 或 rules 約束。品質把關只有後端的 reviewer + analyze-fix loop，缺少 planner 完成後的即時自檢。

此外，planning 和 execution 都以整個 change 為單位運作。使用者無法針對單一 spec 獨立 plan/apply，也無法接入外部 planning 工具（如 gstack plan-eng-review）的產出。

現有結構：
- `cmd/plan.go` — plan CLI，已有 `--context-only`、`--research`、`--check` flags
- `cmd/execute.go` — execute CLI，已有 `--context-only` flag
- `internal/executor/context.go` — `ExecutionContext` 結構，`TaskItem` 有 `id`、`name`、`description`、`status`、`depends`、`files`、`satisfies`、`skills`
- `mysd/skills/plan/SKILL.md` — plan orchestrator，7 步驟 pipeline
- `mysd/agents/mysd-planner.md` — planner agent，產出 `TasksFrontmatterV2` 格式

## Goals / Non-Goals

**Goals:**

- Agent 產出 artifact 前有結構化指引（template + rules + dependencies）
- Planner 完成後有即時自檢（4 項品質檢查）
- 支援 per-spec 的 plan 和 apply（透過 task 的 `spec` tag）
- 支援外部 plan/tasks 作為 planner 的 context 輸入

**Non-Goals:**

- 不改 Spectra 的 `spectra instructions` API — 這是 mysd 自己的 `mysd instructions` 指令
- 不改 propose workflow — propose 已有自己的 reviewer + analyze-fix
- 不做 artifact 類型的動態發現 — mysd 的 artifact 固定為 design + tasks
- 不實作 tasks.md 的 per-spec 拆分（多檔案）— 維持單一 tasks.md + spec tag

## Decisions

### D-01: `mysd instructions` CLI 指令

新增 `cmd/instructions.go`，提供 `mysd instructions <artifact-id> --change <name> --json` 指令。

支援的 artifact-id：`design`、`tasks`。

輸出 JSON 結構：

```json
{
  "artifactId": "tasks",
  "changeName": "enhance-plan-pipeline",
  "outputPath": ".specs/changes/enhance-plan-pipeline/tasks.md",
  "template": "--- yaml frontmatter template ---",
  "rules": ["rule 1", "rule 2"],
  "instruction": "artifact-specific guidance text",
  "dependencies": [
    {"id": "design", "path": "design.md", "done": true},
    {"id": "specs", "path": "specs/", "done": true}
  ],
  "selfReviewChecklist": [
    "No TBD/TODO/FIXME placeholders",
    "Every MUST requirement has a task with matching satisfies field",
    "No single task targets more than 3 files",
    "All file paths in tasks exist in proposal Impact or design"
  ]
}
```

template 和 rules 內容內嵌在 Go code 中（不從外部檔案讀取），因為 artifact 類型固定。

**替代方案**：從 YAML 設定檔讀取 template/rules — 棄選，因為只有 2 個 artifact，hardcode 更簡單且避免檔案找不到的問題。

### D-02: Inline Self-Review（orchestrator 層）

在 plan SKILL.md 中，planner 完成（Step 5）之後、reviewer（Step 5b）之前，新增 Step 5a: Inline Self-Review。

Orchestrator 直接讀取 design.md 和 tasks.md，執行 4 項檢查：

1. **Placeholders** — 掃描 TBD、TODO、FIXME、"implement later"、empty sections
2. **Internal Consistency** — proposal capabilities ↔ specs ↔ design decisions ↔ tasks 互相對齊
3. **Scope** — tasks > 15 則警告、單 task 描述中涵蓋 > 3 個檔案則警告
4. **Ambiguity** — success/failure 條件是否可測試、boundary conditions 是否定義

發現問題時 orchestrator 直接用 Edit tool 修復，不再 spawn 額外 agent。

### D-03: Self-Review Checklist 內嵌於 instructions rules（agent 層）

`mysd instructions tasks` 的 `selfReviewChecklist` 欄位提供 checklist 給 planner agent。Agent 在寫作時應遵守這些規則，orchestrator 在 Step 5a 再驗一次。

這形成雙重保障：agent 寫作時預防 + orchestrator 交付前驗證。

### D-04: TasksFrontmatterV2 新增 `spec` 欄位

在 `TaskItem` struct 中新增 `Spec` 欄位：

```go
type TaskItem struct {
    ID          int      `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    Status      string   `json:"status"`
    Spec        string   `json:"spec,omitempty"`      // NEW
    Depends     []int    `json:"depends,omitempty"`
    Files       []string `json:"files,omitempty"`
    Satisfies   []string `json:"satisfies,omitempty"`
    Skills      []string `json:"skills,omitempty"`
}
```

YAML frontmatter 對應格式：

```yaml
tasks:
  - id: 1
    name: "Implement source detection"
    spec: "material-selection"
    status: pending
```

`spec` 欄位值對應 `specs/<name>/spec.md` 中的 `<name>`。未標記 spec 的 task 視為 change-level task（如 integration test）。

### D-05: Per-spec plan（`--spec` flag）

`cmd/plan.go` 新增 `--spec` flag。行為：

- `mysd plan --spec material-selection --context-only`：context JSON 只包含該 spec 的 requirements 和相關 design sections
- planner agent 只為該 spec 產出 tasks
- 寫入 tasks.md 時，merge 進既有 tasks（不覆蓋其他 spec 的 tasks），新 task ID 從現有最大 ID + 1 開始

Interactive 模式（無 `--spec`）：
1. 讀取 change 的所有 specs
2. 檢查哪些 specs 尚未有對應 tasks（`spec` 欄位未出現在 tasks.md 中）
3. 列出選單讓使用者選擇，或選 "All"
4. `--auto` 模式不指定 `--spec` 則跑全部

### D-06: Per-spec apply（`--spec` flag）

`cmd/execute.go` 新增 `--spec` flag。行為：

- `mysd execute --spec material-selection --context-only`：`pending_tasks` 過濾為 `spec == "material-selection"` 的 tasks
- wave grouping 只計算過濾後的 tasks

Interactive 模式（無 `--spec`）：
1. 讀取 tasks.md 中所有 pending tasks
2. 按 `spec` 分組顯示
3. 列出選單讓使用者選擇，或選 "All"
4. `--auto` 模式不指定 `--spec` 則跑全部

### D-07: External input（`--from` flag）

`cmd/plan.go` 新增 `--from <path>` flag。行為：

- 讀取指定路徑的檔案內容（支援 `.md` 格式）
- 將內容作為 `external_input` 欄位加入 plan context JSON
- planner agent 將此內容作為 context 參考（等同 research findings 的地位）

不做格式轉換 — planner agent 負責理解外部內容並轉換為 TasksFrontmatterV2。

## Risks / Trade-offs

- [Task merge 衝突] 多次 per-spec plan 可能產生 ID 衝突或 depends 引用斷裂 → Mitigation: ID 從最大值 +1 遞增，merge 前驗證 depends 引用
- [Self-review 過度修改] Orchestrator 自檢可能改壞 planner 的正確產出 → Mitigation: 自檢只修 4 類明確問題，不做風格調整
- [External input 品質] gstack 產出可能包含 mysd 無法理解的格式 → Mitigation: 作為 context 參考而非直接導入，planner 自行判斷
- [`spec` 欄位空值] 既有 tasks.md 沒有 `spec` 欄位，per-spec 過濾會漏選 → Mitigation: 空值 task 視為 change-level，per-spec 過濾不影響它們（"All" 選項包含）
