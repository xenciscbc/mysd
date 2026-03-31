## Context

`mysd execute --context-only` 目前回傳 `ExecutionContext` JSON，包含結構化欄位（pending_tasks、execution_mode 等），由 apply SKILL.md 做條件判斷決定行為。所有狀態解讀邏輯分散在 prompt 裡，不可測試且每次修改需改 SKILL.md。

相關 code：
- `internal/executor/context.go` — `ExecutionContext` struct
- `cmd/execute.go` — `--context-only` 和 `--preflight` 輸出
- `mysd/skills/apply/SKILL.md` — orchestrator prompt

## Goals / Non-Goals

**Goals:**

- 在 `ExecutionContext` 新增 `Instruction` 欄位，由 Go binary 動態生成
- 涵蓋 7 種情境的自然語言指引
- Instruction 生成邏輯可單元測試
- 簡化 SKILL.md 的條件分支

**Non-Goals:**

- 不改變現有 `--preflight` flag 的行為（preflight 資訊融入 instruction 但 flag 本身不變）
- 不改變 executor agent 的 prompt（instruction 是給 orchestrator 的，不是給 executor 的）
- 不引入 template engine（用 Go 程式碼直接組合字串）

## Decisions

### D1: Instruction 生成函式放在 `internal/executor` package

Instruction 生成需要讀取 `ExecutionContext` 和 `PreflightReport`。將 `GenerateInstruction` 函式放在 `internal/executor` package，與 `BuildContext` 同層級。

替代方案：放在 `cmd/execute.go` — 但這會讓 cmd 層承擔業務邏輯，且不好單元測試。

### D2: Instruction 由多個 segment 組合

每種情境生成一個 segment（一句話），最後用換行合併。這允許同時觸發多個情境（例如「恢復中斷」+「有 stale artifacts」）。

```
segment 優先順序（由高到低）：
1. all_done     — "所有 task 已完成。建議執行 verify 或 archive。"
2. has_failed   — "T{id} 上次執行失敗。建議 retry 或 skip。"
3. resume       — "恢復中斷：{done}/{total} 已完成，從 T{next_id} 繼續。"
4. last_task    — "最後一個 task T{id}。完成後將進入 verify。"
5. first_run    — "{total} 個 task 待執行。從 T{first_id} 開始。"
6. stale        — "上次 plan 距今 {days} 天，建議重新 plan。"
7. missing_files — "偵測到 {count} 個檔案不存在，請確認再開始。"
```

Segment 1-5 互斥（task 狀態），segment 6-7 可與 1-5 疊加。

替代方案：單一 template 字串 — 但無法乾淨地組合多個條件。

### D3: Preflight 資料由 caller 傳入

`GenerateInstruction` 接收 `ExecutionContext` + 可選的 `*PreflightReport`。在 `cmd/execute.go` 中，`--context-only` 路徑同時呼叫 `runPreflight` 取得 report，傳入生成函式。

這避免 instruction 函式自己做 I/O（讀檔案、檢查 state），保持純函式可測試。

### D4: Instruction 語言為英文

Instruction 是給 LLM orchestrator 讀的，使用英文最不容易被誤解。SKILL.md 的 prompt 也是英文。

### D5: SKILL.md 簡化方式

在 SKILL.md Step 2 之後加一步：「讀取 `instruction` 欄位並顯示給使用者。遵循 instruction 中的指引。」移除 SKILL.md 中已被 instruction 涵蓋的硬編碼判斷邏輯（如 pending_tasks empty 檢查、preflight status 判斷）。

## Risks / Trade-offs

- [Instruction 可能過長] → segment 設計限制每個 segment 一句話，最多 7 句。實際上同時觸發超過 3 個 segment 很少見。
- [Instruction 與 SKILL.md 指引衝突] → SKILL.md 明確聲明「instruction 欄位優先」，逐步遷移硬編碼邏輯。
- [Preflight I/O 在 --context-only 路徑增加延遲] → preflight 只讀 tasks.md 和 STATE.json，I/O 開銷極低（< 10ms）。
