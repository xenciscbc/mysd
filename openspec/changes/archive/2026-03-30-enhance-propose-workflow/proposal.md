## Why

現有的 propose 流程對需求來源的處理偏「被動接收」— 使用者輸入什麼就用什麼，缺乏主動挖掘需求的能力。Step 2b 只問一句 "What would you like to change?"，不會追問動機、邊界或成功標準。4 維度 research 也沒有專門做需求釐清的維度。此外，來源偵測邏輯散落在 Step 2 的 6 個 Priority 中，混合了「找 change name」和「找 source content」兩件不同的職責。

這導致：
- proposal-writer 必須自行推測功能邊界，品質取決於運氣
- 既有 spec 的掃描只是資訊性的，不會主動偵測功能重疊或衝突
- change name 在素材不足時就被推導，品質不佳

## What Changes

- **素材選擇步驟**：整合原 Source Detection 的來源偵測邏輯，偵測 6 種來源（對話 context、Claude plan、gstack plan、active change、deferred notes、手動輸入），列出讓使用者選擇
- **需求評估 + 訪談步驟**：orchestrator 根據彙整素材評估完整度（problem / boundary / success criteria），針對不足處逐一提問（0-N 次），產出結構化 `requirement_brief`
- **既有 spec 對齊**：訪談階段將相關既有 spec 內容餵入，orchestrator 主動詢問功能重疊時是擴充還是新建
- **Change name 延後推導**：移到訪談完成後，從完整的 requirement_brief 推導，品質更好
- **Gray area 條件化**：Step 11-12（gray area + dual-loop exploration）僅在使用者接受 research 後才執行

## Capabilities

### New Capabilities

- `material-selection`: 素材來源偵測與使用者選擇機制，整合原 Source Detection 的分散邏輯
- `requirement-interview`: orchestrator 主導的需求評估與訪談流程，產出結構化 requirement_brief

### Modified Capabilities

- `planning`: propose 流程步驟重新排序，change name 延後推導，gray area 步驟條件化

## Impact

- Affected specs: `material-selection`（新）、`requirement-interview`（新）、`planning`（修改 propose 相關流程）
- Affected code: `mysd/skills/propose/SKILL.md`（主要變更）
