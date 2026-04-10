## 1. 新增 model resolve 子命令

- [x] 1.1 在 `cmd/model.go` 新增 `resolve` 子命令，接受一個 positional argument（role name），呼叫 `config.ResolveModel` 後將 short name 輸出至 stdout。實作 "model resolve subcommand returns single role model" requirement。
- [x] 1.2 在 resolve 子命令中加入 role 驗證：無 argument 時 exit non-zero 並輸出 error 至 stderr，unknown role 時同樣 exit non-zero。實作 "model resolve validates role name" requirement。
- [x] 1.3 在 `cmd/model_test.go` 新增測試：balanced profile resolve executor → `sonnet`、quality profile resolve planner → `opus`、custom profile override、無 argument error、unknown role error。

## 2. 移除 --context-only JSON 中的 model 欄位

- [x] 2.1 從 `cmd/plan.go` 的 `--context-only` JSON 輸出中移除 `model`、`reviewer_model`、`plan_checker_model` 欄位。實作移除 "binary context JSON includes model field" requirement。
- [x] 2.2 從 `cmd/execute.go` 的 `--context-only` JSON 輸出中移除 `Model` 和 `VerifierModel` 欄位。
- [x] 2.3 從 `cmd/verify.go` 的 `--context-only` JSON 輸出中移除 `Model` 欄位。
- [x] 2.4 從 `cmd/design.go` 的 `--context-only` JSON 輸出中移除 `model` 欄位。
- [x] 2.5 從 `cmd/spec.go` 的 `--context-only` JSON 輸出中移除 `model` 欄位。
- [x] 2.6 更新 `cmd/plan_test.go` 等相關測試，移除對 model 欄位的 assertion。

## 3. 更新 SKILL.md 統一使用 model resolve

- [x] 3.1 更新 `mysd/skills/plan/SKILL.md`：移除從 `--context-only` JSON 讀取 `model`、`reviewer_model`、`plan_checker_model` 的指示，改為呼叫 `mysd model resolve planner`、`mysd model resolve reviewer`、`mysd model resolve plan-checker`。實作 "command skills pass model to agents" requirement 及 "model resolve is the single source of truth for skill model queries" requirement。
- [x] 3.2 更新 `mysd/skills/propose/SKILL.md`：移除 `mysd model` 表格解析邏輯，改為呼叫 `mysd model resolve researcher`、`mysd model resolve advisor`、`mysd model resolve proposal-writer`、`mysd model resolve reviewer`。
- [x] 3.3 更新 `mysd/skills/discuss/SKILL.md`：移除 `mysd model` 表格解析邏輯，改為呼叫 `mysd model resolve advisor`（及其他需要的 role）。
- [x] 3.4 更新 `mysd/skills/scan/SKILL.md`：移除 `mysd model` 表格解析邏輯，改為呼叫 `mysd model resolve scanner`。
- [x] 3.5 更新 `mysd/skills/uat/SKILL.md`：移除 `mysd model` 表格解析邏輯，改為呼叫 `mysd model resolve uat-guide`。
- [x] 3.6 更新 `mysd/skills/ff/SKILL.md`：移除從 `--context-only` JSON 讀取 `model` 的指示，改為呼叫 `mysd model resolve fast-forward`（executor 和 verifier 各自 resolve）。
- [x] 3.7 更新 `mysd/skills/ffe/SKILL.md`：移除 `mysd model` 表格解析和 `--context-only` 中 model 讀取，改為呼叫 `mysd model resolve` 取得對應 role model。
- [x] 3.8 更新 `mysd/skills/fix/SKILL.md`：移除從 `mysd execute --context-only` 讀取 model 的指示，改為呼叫 `mysd model resolve executor`。
- [x] 3.9 更新 `mysd/skills/apply/SKILL.md`：移除從 `mysd execute --context-only` 讀取 `model`、`verifier_model` 的指示，改為呼叫 `mysd model resolve executor` 和 `mysd model resolve verifier`。同時更新已有的 `mysd model resolve spec-executor --json` 為 `mysd model resolve spec-executor`（移除 `--json`）。
- [x] 3.10 更新 `mysd/skills/verify/SKILL.md`：移除從 `mysd verify --context-only` 讀取 `model` 的指示，改為呼叫 `mysd model resolve verifier`。

## 4. 建置與驗證

- [x] 4.1 執行 `go build -o mysd.exe .` 建置 binary。
- [x] 4.2 執行 `go test ./cmd/...` 確認所有測試通過。
- [x] 4.3 手動驗證 `mysd model resolve executor` 輸出正確的 short name。
