## Context

`mysd:plan` pipeline 目前在 planner 完成後直接進入 plan-checker（optional）或確認，沒有品質閘道。現有的 inline self-review 邏輯內嵌於 `mysd:propose` Step 12，無法被 plan 複用，也無法透過 model profile 控制品質層級。

此外，`mysd:plan` 的所有 agents（designer, planner, plan-checker）目前共用同一個從 `mysd plan --context-only` 取得的 `model` 欄位，但 `config.go` 的 `DefaultModelMap` 已為每個 role 提供獨立設定——SKILL 層沒有用到這個能力。

`mysd:discuss` 缺乏明確的討論品質規範和強制收斂機制，導致討論可能在沒有清楚結論的情況下結束。

## Goals / Non-Goals

**Goals:**

- 在 plan pipeline 加入獨立的 `mysd-reviewer` agent 作為品質閘道
- 讓 reviewer 和 plan-checker 能各自使用 profile-resolved 的獨立 model
- 為 discuss skill 加入品質規範和強制收斂機制
- propose 的 inline self-review 保持不變

**Non-Goals:**

- 不在 propose 階段引入 reviewer（inline review 已足夠）
- 不修改 plan-checker 的內部邏輯
- 不改變其他 workflow commands 的 model 解析方式

## Decisions

### Per-role model fields in plan context JSON

**決策**: `mysd plan --context-only` 的 JSON 輸出新增 `reviewer_model` 和 `plan_checker_model` 兩個欄位。

**理由**: 現有的單一 `model` 欄位對所有 agents 共用，但 `DefaultModelMap` 已支援 per-role 解析。在 JSON 輸出層加欄位是最小侵入性的修改——`cmd/plan.go` 呼叫 `ResolveModel` 時傳入對應 role 即可。不需要改 `ResolveModel` 的介面，也不影響其他 commands。

**替代方案考慮**: 讓 SKILL.md 呼叫 `mysd model --role reviewer` — 被拒絕，因為這需要新增 CLI 子命令，改動範圍更大。

### reviewer role model assignments

**決策**: quality=opus, balanced=sonnet, budget=sonnet

**理由**: reviewer 是判斷型 role，需要足夠的推理能力。budget profile 使用 sonnet（非 haiku），因為 reviewer 需要跨 artifacts 分析一致性。plan-checker 在 balanced 下是 opus；reviewer 承擔類似但更廣的職責，balanced 下設定為 sonnet 在品質與成本間取得平衡。

### mysd-reviewer agent 的 phase 參數設計

**決策**: phase="propose" vs phase="plan" 決定載入哪些 artifacts。

**理由**: 兩種 phase 的 scope 不同——propose 只有 proposal + specs（2 個檔案），plan 有 4 個完整 artifacts。用 phase 而非傳入 artifact 路徑清單，讓 reviewer 自己決定載入範圍，保持 SKILL 層的呼叫介面簡潔。

## Risks / Trade-offs

- [Risk] plan context JSON 新增欄位可能影響現有 config_test.go 的測試斷言 → Mitigation: 更新測試加入新欄位驗證，舊欄位 `model` 保留不變
- [Risk] reviewer agent 在 plan pipeline 增加延遲 → Mitigation: reviewer 是可接受的品質成本；不提供跳過旗標（reviewer 比 plan-checker 更基礎）
- [Trade-off] discuss 強制收斂可能讓部分用戶感覺限制過強 → 設計保留「明確暫緩」作為合法結論類型，用戶可說「先暫緩」結束討論

## Migration Plan

1. 新增 `mysd/agents/mysd-reviewer.md`
2. 修改 `internal/config/config.go` — `DefaultModelMap` 加入 `reviewer` role
3. 修改 `cmd/plan.go`（或對應的 plan context 程式碼）— 加入 `reviewer_model`、`plan_checker_model` 欄位
4. 更新 `internal/config/config_test.go` — 加入 reviewer role 斷言
5. 修改 `mysd/skills/plan/SKILL.md` — 加入 Step 5b，Step 6 改用 `{plan_checker_model}`
6. 修改 `mysd/skills/discuss/SKILL.md` — 加入品質規範與收斂機制

無 migration 風險：所有變更向後相容，`model` 欄位保留，plan context 只是新增欄位。

## Open Questions

（無）
