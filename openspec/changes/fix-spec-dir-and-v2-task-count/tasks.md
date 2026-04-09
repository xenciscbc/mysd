## 1. Context-only JSON includes spec_dir (D1: spec_dir 傳遞機制, D4: 受影響的 Go cmd 檔案)

- [x] 1.1 cmd/spec.go — 加 `spec_dir` 欄位到 context-only JSON 輸出
- [x] 1.2 cmd/plan.go — 加 `spec_dir` 欄位到 context-only JSON 輸出
- [x] 1.3 cmd/design.go — 加 `spec_dir` 欄位到 context-only JSON 輸出
- [x] 1.4 cmd/execute.go — 加 `spec_dir` 欄位到 ExecutionContext struct 和 JSON 輸出
- [x] 1.5 cmd/scan.go — 加 `spec_dir` 欄位到 ScanContext struct 和 JSON 輸出
- [x] 1.6 cmd/verify.go — 加 `spec_dir` 欄位到 VerificationContext struct 和 JSON 輸出

## 2. Validator supports TasksFrontmatterV2 task count (D3: Validator V2 task count)

- [x] 2.1 internal/validator/validator.go — validateTasks 改用 ParseTasksV2，支援 frontmatter tasks 陣列計數
- [x] 2.2 internal/validator/validator_test.go — 加 V2 task count 測試案例

## 3. Agent definitions use dynamic spec_dir (D2: Agent 路徑替換策略)

- [x] 3.1 mysd/agents/mysd-designer.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.2 mysd/agents/mysd-planner.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.3 mysd/agents/mysd-executor.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.4 mysd/agents/mysd-proposal-writer.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.5 mysd/agents/mysd-spec-writer.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.6 mysd/agents/mysd-reviewer.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.7 mysd/agents/mysd-fast-forward.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.8 mysd/agents/mysd-plan-checker.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.9 mysd/agents/mysd-advisor.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.10 mysd/agents/mysd-verifier.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明
- [x] 3.11 mysd/agents/mysd-scanner.md — 替換 `.specs/` 為 `{spec_dir}/`，Input 段加 spec_dir 說明

## 4. Orchestrators pass spec_dir to agents (D1: spec_dir 傳遞機制)

- [ ] 4.1 mysd/skills/plan/SKILL.md — 從 plan --context-only JSON 提取 spec_dir，傳入 designer/planner/reviewer/plan-checker agent context
- [ ] 4.2 mysd/skills/propose/SKILL.md — 新增 spec_dir 偵測步驟，傳入 researcher/advisor/proposal-writer/spec-writer/reviewer agent context
- [ ] 4.3 mysd/skills/apply/SKILL.md — 從 execute --context-only JSON 提取 spec_dir，傳入 executor/verifier agent context
- [ ] 4.4 mysd/skills/discuss/SKILL.md — 從 plan --context-only JSON 提取 spec_dir，傳入 agent context
- [ ] 4.5 mysd/skills/ff/SKILL.md — 從 plan/execute --context-only JSON 提取 spec_dir，傳入 fast-forward agent context
- [ ] 4.6 mysd/skills/ffe/SKILL.md — 從 plan/execute --context-only JSON 提取 spec_dir，傳入 fast-forward agent context
- [ ] 4.7 mysd/skills/fix/SKILL.md — 從 execute --context-only JSON 提取 spec_dir，傳入 agent context
- [ ] 4.8 mysd/skills/scan/SKILL.md — 從 scan --context-only JSON 提取 spec_dir，傳入 scanner agent context
- [ ] 4.9 mysd/skills/verify/SKILL.md — 從 verify --context-only JSON 提取 spec_dir，傳入 verifier agent context
- [ ] 4.10 mysd/skills/archive/SKILL.md — 從 execute --context-only JSON 提取 spec_dir

## 5. Build 和驗證

- [ ] 5.1 跑測試：go test ./internal/validator/ ./cmd/ -v
- [ ] 5.2 Build binary：go build -o mysd.exe .
- [ ] 5.3 在 test_mysd（openspec/ 專案）跑 mysd validate 驗證無誤報
