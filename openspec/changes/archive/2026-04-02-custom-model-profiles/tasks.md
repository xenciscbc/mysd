## 1. Config 結構擴展

- [x] 1.1 在 `internal/config/defaults.go` 新增 `CustomProfile` struct（含 `Base` 和 `Models` 欄位），並在 `ProjectConfig` 新增 `CustomProfiles map[string]CustomProfile` 欄位（對應 design 的 CustomProfile 結構定義）
- [x] 1.2 新增 `ValidateCustomProfiles(knownRoles []string, customProfiles map[string]CustomProfile) []string` 函數於 `internal/config/`，檢查 invalid role name warning 和 invalid base profile warning，回傳 warning 訊息列表

## 2. Model 解析邏輯

- [x] 2.1 實作 custom profile model resolution（ResolveModel 擴展解析鏈）：擴展 `ResolveModel` 簽名，新增 `customProfiles map[string]CustomProfile` 參數，實作四層優先順序（ModelOverrides → custom models → base profile → fallback）
- [x] 2.2 更新所有 `ResolveModel` 呼叫端傳入 `customProfiles` 參數

## 3. CLI 命令更新

- [x] 3.1 實作 custom profile selection via CLI（model set 驗證擴展）：修改 `cmd/model.go` 的 `runModelSet`，驗證邏輯改為內建 profile → `cfg.CustomProfiles` → 報錯並列出所有可用 profile
- [x] 3.2 修改 `cmd/model.go` 的 `runModelRead`，支援 custom profile display，正確顯示自訂 profile 名稱與解析結果
- [x] 3.3 處理無效 role 警告時機：在 `runModelRead` 和 `runModelSet` 中呼叫 `ValidateCustomProfiles`，將 warning 輸出至 stderr

## 4. 測試

- [x] 4.1 為 `CustomProfile` struct 和 custom profile definition in configuration 撰寫單元測試：有效定義、空 models、缺少 base
- [x] 4.2 為 `ResolveModel` 擴展撰寫單元測試：role overridden in custom profile、role not overridden（fallback to base）、ModelOverrides takes precedence over custom profile
- [x] 4.3 為 `ValidateCustomProfiles` 撰寫單元測試：unknown role in custom profile、unknown base profile
- [x] 4.4 為 `runModelSet` 和 `runModelRead` 撰寫測試：set a custom profile、set an unknown profile、display custom profile
