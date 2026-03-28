## 1. propose SKILL.md — Model Resolution 擴充

- [x] 1.1 修改 propose SKILL.md Step 3b，將 model resolution 從單一 `model` 擴充為包含 `reviewer_model`，加入 profile-to-reviewer-model mapping table（quality=opus, balanced=sonnet, budget=sonnet）— 對應 spec: Command skills pass model to agents

## 2. propose SKILL.md — Step 12 改為呼叫 mysd-reviewer

- [x] 2.1 移除 propose SKILL.md Step 12 的 inline self-review 邏輯（Check 1-4）
- [x] 2.2 新增 Step 12a：執行 `mysd validate {change_name}` 並捕獲 output — 對應 spec: Reviewer agent performs artifact quality checks（propose pipeline invokes reviewer with validate output）
- [x] 2.3 新增 Step 12b：spawn `mysd-reviewer` with `phase: "propose"`、`reviewer_model`、`validate_output` — 對應 spec: Reviewer agent performs artifact quality checks
- [x] 2.4 更新 Step 13 Final Summary：將 "Self-review result" 改為 "Reviewer result"，引用 Step 12b 的 reviewer summary

## 3. mysd-reviewer agent — 加入 Rationalization Table

- [x] 3.1 在 mysd-reviewer.md 的 Step 1（Load Artifacts）和 Step 2（Check 1）之間插入 Rationalization Table section，包含 6 條 anti-pattern 對照表 — 對應 spec: Reviewer includes rationalization table
