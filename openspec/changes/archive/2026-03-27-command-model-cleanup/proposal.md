## Why

mysd 的 command 和 agent 檔案目前全部 hardcode `model: claude-sonnet-4-5`，忽略了已有的 profile 系統（`/mysd:model`）。這導致 profile 設定無法實際影響工作流中 agent 的 model 選擇。同時，多個已被合併到其他命令中的獨立 command（execute、capture、design、spec）仍然存在，造成引用混亂和使用者困惑。

## What Changes

- **移除工作流 command 和 agent 的 model frontmatter**，讓 profile 系統控制 agent model
- **輕量工具命令指定 sonnet**：status、lang、model、note、docs、statusline、update
- **重度獨立命令指定 opus**：init、scan、fix
- **刪除已廢棄的 command 檔案**：execute（redirect 空殼）、capture（被 discuss 取代）、design（內嵌於 plan）、spec（內嵌於 propose/discuss）
- **清理所有過時引用**：`/mysd:execute` → `/mysd:apply`，移除 `/mysd:spec`、`/mysd:design` 引用
- **移除 propose 的 `--skip-spec` flag**，spec 生成是 propose 的必要步驟
- **apply 的 verify 改為不可跳過**
- **scan 移除 next steps 段落**
- **Profile 短名修正**：`DefaultModelMap` 改為 `sonnet`/`opus`/`haiku` 而非 `claude-sonnet-4-5`
- **Command 傳遞 profile model 給 agent**：spawn 時指定 model 參數並顯示使用的 model

## Capabilities

### New Capabilities

- `model-passthrough`: Command skill 讀取 binary context JSON 中的 model 欄位，spawn agent 時傳入並顯示

### Modified Capabilities

- `execution`: execute 命令已改名為 apply，清理所有殘留引用；apply 的 verify 步驟改為強制執行
- `fast-forward`: ff/ffe 中的 `/mysd:execute` 引用更新為 `/mysd:apply`

## Impact

- Affected specs: `execution`, `fast-forward`
- Affected code:
  - `plugin/commands/*.md` — 所有 command 的 model frontmatter 及引用修正
  - `plugin/agents/*.md` — 所有 agent 的 model frontmatter 移除及引用修正
  - `internal/config/config.go` — DefaultModelMap 值改為短名
  - `internal/config/config_test.go` — 對應測試更新
  - `cmd/execute.go`, `cmd/spec.go`, `cmd/design.go` — 提示訊息更新
  - `internal/spec/schema_test.go`, `internal/executor/context_test.go` — 測試中的引用更新
