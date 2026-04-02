## Why

目前 mysd 僅提供三個內建 model profile（quality、balanced、budget），無法滿足不同團隊或專案對 agent role 與 model 配對的特定需求。使用者需要能透過設定檔定義自己的 profile，只寫與 base profile 的差異部分，即可快速切換。

## What Changes

- 新增 `custom_profiles` 設定區塊於 `mysd.yaml`，允許使用者定義自訂 profile
- 每個自訂 profile 可指定 `base`（繼承某個內建 profile）和 `models`（差異覆蓋）
- `mysd model set <name>` 擴展查找範圍：內建 profile → `custom_profiles` → 報錯
- `mysd model` 顯示時能正確解析自訂 profile
- 載入時若 `models` 中出現不在 `knownRoles` 裡的 role name，發出 warning

## Capabilities

### New Capabilities

- `custom-model-profile`: 透過 `mysd.yaml` 的 `custom_profiles` 區塊定義自訂 model profile，支援 base 繼承與差異覆蓋

### Modified Capabilities

(none)

## Impact

- Affected code:
  - `internal/config/defaults.go` — `ProjectConfig` 新增 `CustomProfiles` 欄位，`ResolveModel` 擴展解析邏輯
  - `internal/config/config.go` — `Load` 函數處理新欄位的 unmarshal
  - `cmd/model.go` — `model set` 驗證邏輯擴展、`model` 顯示邏輯支援自訂 profile
