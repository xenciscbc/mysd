## 1. 更新 Planner Agent Template

- [x] 1.1 修正 `mysd/agents/mysd-planner.md` Step 4 的 tasks.md template，加入 `satisfies`、`depends`、`files` 欄位，使其與 `cmd/instructions.go` 的 template 和 `TasksFrontmatterV2` struct 完全一致。實作 "Task Planning" requirement 中「tasks 陣列 MUST 包含 satisfies 欄位」的要求
- [x] 1.2 在 planner agent 的 Step 4 Key fields 說明中，補充 `satisfies`、`depends`、`files`、`skills` 欄位的描述，確保 agent 理解每個欄位的用途。實作 "Task Planning" requirement 中 MAY 欄位的說明

## 2. 驗證 Go 代碼鏈路

- [x] 2.1 確認 `internal/spec/updater.go` 的 `ParseTasksV2()` 能正確 parse 包含完整欄位的 tasks.md（含 `satisfies`、`depends`、`files`）。跑 `go test ./internal/spec/...` 確認現有測試通過。實作 "Task Planning" requirement 中 "Task status round-trip via UpdateTaskStatus" scenario
- [x] 2.2 確認 `internal/planchecker/checker.go` 的 `CheckCoverage()` 在有 `satisfies` 欄位的 tasks.md 上能正確計算覆蓋率。跑 `go test ./internal/planchecker/...` 確認現有測試通過。實作 "Plan Coverage Checking" requirement
- [x] 2.3 確認 `internal/executor/` 的 wave grouping 在有 `depends` 欄位的 tasks.md 上能正確分組。跑 `go test ./internal/executor/...` 確認現有測試通過。實作 "Wave Grouping" requirement

## 3. 建置與端對端驗證

- [x] 3.1 執行 `go build -o mysd.exe .` 建置 binary
- [x] 3.2 用更新後的 planner agent template 作為參考，手動建立一個包含完整 V2 frontmatter 的 tasks.md 測試檔案，驗證 `mysd task-update` 能正確更新 task status 並 round-trip YAML frontmatter
