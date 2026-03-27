## 1. Profile model resolution uses short names

- [x] 1.1 更新 `internal/config/config.go` 的 `DefaultModelMap`：所有 model 值從 `claude-sonnet-4-5` / `claude-haiku-3-5` 改為 `sonnet` / `haiku`，fallback 也改為 `sonnet`（profile model resolution uses short names）
- [x] 1.2 更新 `internal/config/config_test.go`：所有 assert 的 model 期望值改為短名（`sonnet`、`haiku`、`opus`）
- [x] 1.3 確認 `cmd/plan.go`、`cmd/design.go`、`cmd/spec.go` 的 `--context-only` JSON 輸出自動使用短名（binary context JSON includes model field）

## 2. Execute/Spec/Design/Capture command skill removal

- [x] 2.1 刪除 `plugin/commands/mysd-execute.md`（execute command skill removed）
- [x] 2.2 刪除 `plugin/commands/mysd-capture.md`（capture command skill removed）
- [x] 2.3 刪除 `plugin/commands/mysd-design.md`（design command skill removed）
- [x] 2.4 刪除 `plugin/commands/mysd-spec.md`（spec command skill removed）

## 3. Execution Context 引用清理 — execute → apply

- [x] 3.1 更新 `plugin/agents/mysd-planner.md` :130 的 `/mysd:execute` → `/mysd:apply`（execution context references）
- [x] 3.2 更新 `plugin/agents/mysd-plan-checker.md` :55, :170 的 `/mysd:execute` → `/mysd:apply`
- [x] 3.3 更新 `plugin/agents/mysd-fast-forward.md` :157 的 `/mysd:execute` → `/mysd:apply`（fast-forward ff references updated）
- [x] 3.4 更新 `plugin/commands/mysd-archive.md` :67 的 `/mysd:execute` → `/mysd:apply`
- [x] 3.5 更新 `plugin/commands/mysd-uat.md` :93 的 `/mysd:execute` → `/mysd:apply`
- [x] 3.6 更新 `plugin/commands/mysd-verify.md` :32, :103 的 `/mysd:execute` → `/mysd:apply`
- [x] 3.7 更新 `cmd/execute.go` :71 提示訊息從 `/mysd:execute` → `/mysd:apply`
- [x] 3.8 更新 `internal/spec/schema_test.go` :74 和 `internal/executor/context_test.go` :109, :121 的 `/mysd:execute` → `/mysd:apply`

## 4. 引用清理 — spec/design/verify 指向修正

- [x] 4.1 更新 `plugin/commands/mysd-propose.md`：移除 `--skip-spec` 相關邏輯（:3, :4, :19, :21, :23, :203），next steps 的 `/mysd:design` 和 `/mysd:spec` 改為 `/mysd:discuss` 和 `/mysd:plan`（:237, :239）
- [x] 4.2 更新 `plugin/commands/mysd-status.md` :87-88：移除 `specced` 狀態引用，`proposed` 的 next step 改為 `/mysd:plan`
- [x] 4.3 更新 `plugin/commands/mysd-scan.md` :108-111：移除 next steps 整段
- [x] 4.4 更新 `plugin/commands/mysd-archive.md` :55, :66 的 `/mysd:verify` 引用改為 `/mysd:apply`
- [x] 4.5 更新 `plugin/commands/mysd-uat.md` :34 的 `/mysd:verify` 引用改為 `/mysd:apply`
- [x] 4.6 更新 `plugin/agents/mysd-uat-guide.md` :104 的 `/mysd:verify` 引用改為 `/mysd:apply`
- [x] 4.7 更新 `plugin/agents/mysd-proposal-writer.md` :126 next step 從 `/mysd:spec` 改為 `/mysd:discuss` 或 `/mysd:plan`
- [x] 4.8 更新 `cmd/spec.go` :79 和 `cmd/design.go` :86 的提示訊息，說明功能已整合進 plan/propose

## 5. Apply command verification is mandatory

- [x] 5.1 更新 `plugin/commands/mysd-apply.md` Step 5b：移除使用者確認提示，verify 在 build+test 通過後自動執行；移除 `/mysd:verify` 的 fallback 提示（apply command verification is mandatory）

## 6. Workflow commands and agents have no model frontmatter / Standalone utility commands specify fixed model

- [x] 6.1 移除工作流 command 的 model frontmatter：`propose`、`discuss`、`plan`、`apply`、`archive`、`ff`、`ffe`、`uat`（workflow commands and agents have no model frontmatter）
- [x] 6.2 移除所有 agent 的 model frontmatter：`advisor`、`executor`、`fast-forward`、`proposal-writer`、`researcher`、`spec-writer`
- [x] 6.3 將 `init`、`scan`、`fix` 的 model 改為 `claude-opus-4-6`（standalone utility commands specify fixed model — opus）
- [x] 6.4 確認 `status`、`lang`、`model`、`note`、`docs`、`statusline`、`update` 保持 `model: claude-sonnet-4-5`（standalone utility commands specify fixed model — sonnet）

## 7. Command skills pass model to agents / Command skills display model on agent spawn

- [x] 7.1 更新 `plugin/commands/mysd-plan.md`：讀取 context JSON 的 model 欄位，command skills pass model to agents，spawn 前 command skills display model on agent spawn `Spawning {agent} ({model})...`
- [x] 7.2 更新 `plugin/commands/mysd-apply.md`：讀取 context JSON 的 model 欄位，spawn executor 和 verifier agent 時傳入 model 參數並顯示
- [x] 7.3 更新 `plugin/commands/mysd-propose.md`：spawn researcher、advisor、proposal-writer、spec-writer agent 時傳入 model 並顯示
- [x] 7.4 更新 `plugin/commands/mysd-discuss.md`：spawn researcher、advisor、proposal-writer、spec-writer、designer agent 時傳入 model 並顯示
- [x] 7.5 更新 `plugin/commands/mysd-ff.md`（fast-forward ff）和 `mysd-ffe.md`（extended fast-forward ffe）：spawn agent 時傳入 model 並顯示
- [x] 7.6 更新 `plugin/commands/mysd-archive.md`：確認不需要 — archive 不 spawn agent，doc 更新由 command 直接執行
