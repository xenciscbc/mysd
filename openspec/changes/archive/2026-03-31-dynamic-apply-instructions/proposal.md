## Why

目前 `mysd execute --context-only` 回傳結構化 JSON 資料（pending_tasks、execution_mode 等），但 orchestrator（apply SKILL.md）需要自行根據這些欄位做條件判斷來決定行為。這導致 SKILL.md 充斥分支邏輯，且每次新增情境都要改 prompt。

在 `--context-only` 輸出中新增一個動態生成的 `instruction` 自然語言欄位，由 Go binary 根據當前狀態組合指引，讓 orchestrator 直接遵循而不用自己推斷。這將狀態判斷邏輯集中到可測試的 Go code 裡，並讓 SKILL.md 變薄。

## What Changes

- `mysd execute --context-only` 的 JSON 輸出新增 `instruction` 字串欄位
- Go binary 根據 7 種狀態情境動態生成 instruction 內容：
  1. 首次執行（所有 task 皆 pending）
  2. 恢復中斷（部分 task 已完成）
  3. 單一 task 剩餘
  4. 全部完成
  5. 有失敗 task
  6. Preflight 警告（missing files）
  7. Stale artifacts（staleness 超門檻）
- `mysd/skills/apply/SKILL.md` 簡化，改為讀取並遵循 `instruction` 欄位

## Capabilities

### New Capabilities

- `dynamic-instruction`: 動態生成 execution instruction — 根據 change 狀態、task 進度、preflight 結果生成自然語言指引，輸出於 `--context-only` JSON

### Modified Capabilities

- `execution`: ExecutionContext JSON 新增 `instruction` 欄位；SKILL.md 改為依賴此欄位指導行為

## Impact

- Affected specs: `execution`（新增 instruction 欄位要求）、新建 `dynamic-instruction` spec
- Affected code:
  - `cmd/execute.go` — instruction 生成邏輯
  - `cmd/execute_test.go` — 各情境測試
  - `mysd/skills/apply/SKILL.md` — 簡化分支邏輯，改為讀取 instruction
