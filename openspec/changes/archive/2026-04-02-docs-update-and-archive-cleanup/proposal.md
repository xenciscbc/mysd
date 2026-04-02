## Why

mysd 的 archive SKILL.md 路徑過時（仍引用 `.specs/archive/`），互動式 UAT prompt 不再需要，且缺少獨立觸發 doc 更新的機制 — 目前只能透過 archive/ff/ffe 間接觸發。

## What Changes

- 移除 `cmd/archive.go` 的互動式 UAT prompt（`isInteractive()` + `Run UAT first?` 區塊）
- 修正 `mysd/skills/archive/SKILL.md` 的 archive 路徑引用，從 `.specs/archive/{change_name}/` 改為 `openspec/changes/archive/YYYY-MM-DD-{change_name}/`
- 新增 `/mysd:docs update` skill，支援多種 scope 觸發 doc 更新：
  - 預設：最近一次 archived change
  - `--change <name>`：指定某個 change
  - `--last N`：最近 N 次 changes
  - `--full`：掃描 codebase 全面更新
  - 自由描述文字：使用者自述更新範圍

## Non-Goals

- 不修改 `mysd docs add/remove` 的行為
- 不改變 archive/ff/ffe 中既有的 doc 更新觸發邏輯
- 不新增 binary 層面的 `mysd docs update` 子命令 — 這是純 SKILL.md 層面的功能

## Capabilities

### New Capabilities

- `docs-update`: 獨立觸發 doc 更新的 SKILL.md，支援多種 scope（change、last N、full codebase scan、自由描述）

### Modified Capabilities

(none)

## Impact

- Affected code:
  - `cmd/archive.go` — 移除 UAT prompt 相關程式碼
  - `cmd/archive_test.go` — 移除 UAT 相關測試
  - `mysd/skills/archive/SKILL.md` — 修正 archive 路徑引用
  - `mysd/skills/docs/SKILL.md` — 新增 update 子命令說明
- Affected skills:
  - 新增 `mysd/skills/docs-update/SKILL.md`（或擴展 `mysd/skills/docs/SKILL.md`）
