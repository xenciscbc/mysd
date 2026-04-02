## Context

mysd 目前透過 `DefaultModelMap` 提供三個內建 model profile（quality、balanced、budget），每個 profile 定義了所有 agent role 對應的 model。使用者可透過 `ModelOverrides` 覆蓋個別 role，但無法定義完整的自訂 profile。

現有解析順序：`ModelOverrides[role]` → `DefaultModelMap[profile][role]` → `"sonnet"`

設定檔位置：`.claude/mysd.yaml`（project-level）或 `~/.claude/mysd.yaml`（user-level），由 Viper 處理。

## Goals / Non-Goals

**Goals:**

- 使用者可在 `mysd.yaml` 中定義自訂 profile，指定 base 繼承與 role-model 差異覆蓋
- `mysd model set` 和 `mysd model` 正確處理自訂 profile
- 載入時警告無效的 role name

**Non-Goals:**

- 不支援自訂 profile 之間的鏈式繼承（自訂 profile 的 base 只能是內建 profile）
- 不支援自訂 model name（值仍限制為 "opus"、"sonnet"、"haiku"）
- 不新增 CLI 子命令來管理 custom profiles（透過直接編輯 YAML 管理）

## Decisions

### CustomProfile 結構定義

新增 `CustomProfile` struct 和 `CustomProfiles` map 至 `ProjectConfig`：

```go
type CustomProfile struct {
    Base   string            `yaml:"base" mapstructure:"base"`
    Models map[string]string `yaml:"models" mapstructure:"models"`
}
```

`ProjectConfig` 新增欄位：
```go
CustomProfiles map[string]CustomProfile `yaml:"custom_profiles" mapstructure:"custom_profiles"`
```

選擇 struct 而非 `map[string]map[string]string`，因為需要明確的 `base` 欄位與 `models` 分離。

### ResolveModel 擴展解析鏈

新的解析順序：

```
ModelOverrides[role]
  → CustomProfiles[profile].Models[role]
    → DefaultModelMap[CustomProfiles[profile].Base][role]
      → DefaultModelMap[profile][role]  (profile 本身是內建 profile 時)
        → "sonnet"
```

`ResolveModel` 需要接收 `CustomProfiles` 參數。為保持向後相容，擴展簽名：

```go
func ResolveModel(agentRole, profile string, overrides map[string]string, customProfiles map[string]CustomProfile) string
```

### model set 驗證擴展

`runModelSet` 驗證邏輯從只查 `DefaultModelMap` 改為：

1. 查 `DefaultModelMap[profile]` → 找到則為內建 profile
2. 查 `cfg.CustomProfiles[profile]` → 找到則為自訂 profile
3. 都找不到 → 報錯，列出所有可用 profile（內建 + 自訂）

### 無效 role 警告時機

在 `config.Load` 完成後，由呼叫端（`runModelRead`、`runModelSet`）檢查 `CustomProfiles` 中的 role name 是否存在於 `knownRoles`。不在 `Load` 內部做，因為 `Load` 不應依賴 `cmd` 層的 `knownRoles` 定義。

新增 `ValidateCustomProfiles` 函數於 `internal/config`，接收 `knownRoles []string` 參數，回傳 warning 列表。

## Risks / Trade-offs

- [Risk] `ResolveModel` 簽名變更影響所有呼叫端 → 呼叫端數量有限（`cmd/model.go` 和 SKILL.md template），影響可控
- [Risk] base 指向不存在的內建 profile → `ValidateCustomProfiles` 一併檢查，載入時警告
- [Trade-off] 不支援鏈式繼承簡化了實作，但限制了靈活性 → 目前內建 profile 僅三個，鏈式需求不明確，YAGNI
