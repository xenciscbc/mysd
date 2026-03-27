## Why

`DefaultModelMap` 中三個 profile（quality / balanced / budget）的 model 分配不正確。quality profile 全部使用 sonnet，但「最佳品質」profile 的思考型 role（spec-writer、planner、verifier 等）應使用 opus。balanced profile 也應在關鍵的判斷/設計/把關 role 使用 opus。budget profile 的 spec-writer 使用 haiku 品質不足，應提升為 sonnet。這是歷史遺留問題——原始設計時只考慮 sonnet/haiku 二選一，未將 opus 納入選項。

## What Changes

- quality profile：8 個思考型 role 從 sonnet 改為 opus（spec-writer、designer、planner、verifier、researcher、advisor、proposal-writer、plan-checker），executor 和 fast-forward 維持 sonnet
- balanced profile：6 個判斷/設計/把關 role 從 sonnet 改為 opus（spec-writer、designer、planner、verifier、advisor、plan-checker），executor、fast-forward、researcher、proposal-writer 維持 sonnet
- budget profile：spec-writer 從 haiku 改為 sonnet，其餘不變
- 更新對應的 unit test 以反映新的 model 分配

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `model-passthrough`: DefaultModelMap 的 model 分配邏輯變更，三個 profile 的 role-to-model 映射重新定義

## Impact

- 受影響程式碼：`internal/config/config.go`（DefaultModelMap 變數）、`internal/config/config_test.go`（對應測試）
- 受影響行為：所有透過 profile system 解析 model 的 workflow command（plan、apply、propose、discuss、ff、ffe）將為思考型 role 使用更強的 model
