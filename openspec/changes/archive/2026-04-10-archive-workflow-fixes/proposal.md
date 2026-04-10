## Why

Archive 流程有三個互相關聯的 bug，導致從 execute → verify → archive 的 pipeline 無法正常完成：

1. **State transition 斷裂**：`mysd execute --context-only` 不做 state transition（by design），但 `task-update` 在所有 tasks 完成時也不推進到 `executed`，導致 phase 永遠卡在 `planned`，後續 verify 和 archive 都被擋住。
2. **Verifier ID 格式不匹配**：verifier agent 產出 `MUST-01` 格式的 ID，但 archive 的 `checkMustItemsDone()` 期望 hash-based ID（如 `spec.md::must-5451802d`），導致 verification 結果無法被 archive 辨識。
3. **Delta spec merge 靜默失敗**：當 delta spec 使用普通 spec 格式（無 `## ADDED Requirements` 等 heading）時，`ParseDelta` 回傳空結果，`MergeSpecs` 只寫入空 frontmatter，主 spec body 消失。

## What Changes

- `mysd task-update` 在最後一個 task 完成（所有 tasks 為 done/skipped）時，自動執行 `planned → executed` state transition
- `mysd verify --write-results` 加入 safety net：若 phase 仍為 `planned` 但所有 tasks 已完成，先自動推進到 `executed` 再處理驗證結果
- Verifier agent prompt（`mysd-verifier.md`）明確要求使用 `--context-only` 輸出的 hash-based ID 原樣回傳
- `MergeSpecs` 加入 fallback：當 `ParseDelta` 未找到任何 delta heading 時，檢查 frontmatter 的 `delta` 欄位決定合併策略（`MODIFIED` → 全量替換 body，`ADDED` → 使用 delta body 作為新 spec 內容）
- Spec writer agent prompt（`mysd-spec-writer.md`）明確要求產出的 delta spec 使用 delta heading 格式

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `execution`: task-update 完成時自動推進 state 到 executed
- `verification`: verify --write-results 加入 phase safety net；verifier agent 使用正確的 ID 格式
- `delta-spec-merge`: MergeSpecs 加入非 delta 格式的 fallback 處理

## Impact

- 受影響程式碼：`cmd/task_update.go`、`cmd/verify.go`、`internal/spec/merge.go`
- 受影響 agent prompts：`mysd/agents/mysd-verifier.md`、`mysd/agents/mysd-spec-writer.md`
- 受影響測試：`cmd/task_update_test.go`、`cmd/verify_test.go`、`internal/spec/merge_test.go`
