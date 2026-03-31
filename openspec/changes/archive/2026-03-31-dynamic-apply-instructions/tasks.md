## 1. GenerateInstruction function 核心實作

- [x] 1.1 在 `internal/executor/instruction.go` 新增 GenerateInstruction function，實作 task state instruction segments（5 個互斥 segment：all_done / has_failed / resume / last_task / first_run），按優先順序評估。函式放在 internal/executor package（D1: Instruction 生成函式放在 internal/executor package），instruction 語言為英文（D4: Instruction 語言為英文）（D2: Instruction 由多個 segment 組合）
- [x] 1.2 在 GenerateInstruction function 中實作 preflight instruction segments（stale / missing_files），當 preflight 非 nil 且有 issue 時附加到 task state segment 之後（D3: Preflight 資料由 caller 傳入）
- [x] 1.3 新增 `internal/executor/instruction_test.go`，測試所有 task state instruction segments 和 preflight instruction segments 的 7 種情境及 combined 情境

## 2. ExecutionContext 整合

- [x] 2.1 在 `internal/executor/context.go` 的 `ExecutionContext` struct 新增 `Instruction string` 欄位（json:"instruction"），滿足 Instruction output in ExecutionContext JSON 要求
- [x] 2.2 將 `cmd/execute.go` 的 `PreflightReport` 和 `PreflightChecks` 和 `StalenessCheck` type 搬到 `internal/executor` package（讓 GenerateInstruction 可以引用，不引入循環依賴）
- [x] 2.3 修改 `cmd/execute.go` 的 `--context-only` 路徑：在 JSON 序列化前呼叫 `runPreflight` 取得 report，再呼叫 `GenerateInstruction` 填入 `ctx.Instruction`（D3: Preflight 資料由 caller 傳入，修改 Execution Context 要求）
- [x] 2.4 新增 `cmd/execute_test.go` 測試：驗證 `--context-only` JSON 輸出包含非空 `instruction` 欄位

## 3. SKILL.md 簡化

- [x] 3.1 修改 `mysd/skills/apply/SKILL.md`：在 Step 2 之後新增一步讀取並顯示 `instruction` 欄位，移除 Step 2c Preflight Check 中已被 instruction 涵蓋的硬編碼判斷邏輯，以及 Step 2 中 pending_tasks empty 的硬編碼檢查（D5: SKILL.md 簡化方式）
