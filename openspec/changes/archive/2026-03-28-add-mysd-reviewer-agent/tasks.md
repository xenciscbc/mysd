## 1. Config 與後端：Reviewer role 與 per-role model

- [x] 1.1 在 `internal/config/config.go` 的 `DefaultModelMap` 新增 `reviewer` role（reviewer role model assignments：quality=opus, balanced=sonnet, budget=sonnet）
- [x] 1.2 在 `internal/config/config_test.go` 加入 reviewer role 斷言，驗證 DefaultModelMap includes reviewer role 三個 profile 的回傳值
- [x] 1.3 修改 `cmd/plan.go`（或 plan context 相關程式碼），在 `--context-only` JSON 輸出中新增 `reviewer_model` 和 `plan_checker_model` 欄位（Per-role model fields in plan context JSON：呼叫 `ResolveModel("reviewer", profile, overrides)` 和 `ResolveModel("plan-checker", profile, overrides)`）
- [x] 1.4 更新對應的 plan context JSON 測試，驗證 Binary context JSON includes model field 新欄位存在且值正確

## 2. mysd-reviewer Agent 定義

- [x] 2.1 建立 `mysd/agents/mysd-reviewer.md`，定義 agent 基本結構：frontmatter（description, allowed-tools），input context 格式（change_name, phase, validate_output, auto_mode），以及 Reviewer agent performs artifact quality checks 的 phase 對應 artifact 載入邏輯（mysd-reviewer agent 的 phase 參數設計）
- [x] 2.2 實作 Check 1：掃描並修復 placeholder content（TBD, TODO, FIXME, 空白欄位, 模糊數量），對應 Reviewer checks for placeholder content
- [x] 2.3 實作 Check 2：跨 artifact 一致性驗證，包含 propose phase 的 proposal↔spec 對應，以及 plan phase 額外的 design↔tasks 一致性，對應 Reviewer checks internal consistency
- [x] 2.4 實作 Check 3：範圍檢查（超過 15 個 MUST requirements 或 pending tasks，以及觸及 3 個以上無關子系統），對應 Reviewer checks scope
- [x] 2.5 實作 Check 4：歧義檢查（可測試的成功條件、邊界條件定義、"the system" 指向明確），對應 Reviewer checks ambiguity
- [x] 2.6 實作結構化 summary 輸出格式（Phase, Issues fixed, Fixed, Cannot auto-fix），對應 Reviewer returns structured summary

## 3. plan/SKILL.md：Reviewer step 與 per-role model

- [x] 3.1 修改 `mysd/skills/plan/SKILL.md` Step 2，從 plan context JSON 解析 `reviewer_model` 和 `plan_checker_model` 欄位（Plan pipeline uses per-role models for reviewer and plan-checker）
- [x] 3.2 在 Step 5 之後、Step 6 之前插入 Step 5b：執行 `mysd validate`、顯示 "Spawning mysd-reviewer ({reviewer_model})..."、用 Task tool 呼叫 mysd-reviewer（Plan pipeline includes reviewer step after planner）
- [x] 3.3 修改 Step 6（plan-checker）改用 `{plan_checker_model}` 而非共用的 `{model}`

## 4. discuss/SKILL.md：品質規範與收斂機制

- [x] 4.1 在 `mysd/skills/discuss/SKILL.md` 的討論 loop 前加入品質規範區塊（一次一問、具體選項、禁止空話、主動推薦、用戶催促處理方式），對應 Discuss skill enforces discussion quality guidelines
- [x] 4.2 在討論 loop 中加入強制收斂機制：達到明確結論時主動提出 Conclusion summary（Decision/Rationale/Capture to 格式），以及不允許無結論結束的處理邏輯，對應 Discuss skill enforces convergence and conclusion capture
