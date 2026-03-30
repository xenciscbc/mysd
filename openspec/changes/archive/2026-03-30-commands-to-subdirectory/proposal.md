## Summary

將 `plugin/commands/mysd-*.md` 平面結構搬到 `plugin/commands/mysd/*.md` 子目錄結構。

## Motivation

目前 20 個 command 檔案全部平放在 `plugin/commands/` 下，每個檔名都帶 `mysd-` prefix。Claude Code 支援 `commands/<namespace>/<name>.md` 子目錄結構，子目錄名稱自動成為 namespace prefix（`/mysd:apply`）。搬到子目錄後檔名去掉重複的 `mysd-` prefix，結構更乾淨。

## Proposed Solution

- 建立 `plugin/commands/mysd/` 子目錄
- 將所有 `plugin/commands/mysd-*.md` 移到 `plugin/commands/mysd/*.md`，去掉 `mysd-` prefix（例如 `mysd-apply.md` → `mysd/apply.md`）
- `plugin/commands/CLAUDE.md` 維持在 `plugin/commands/` 不搬移
- 觸發方式不變：`/mysd:apply`、`/mysd:propose` 等

## Non-Goals

- 不將 commands 轉換為 skills 格式（維持現有 frontmatter 格式）
- 不搬移 agents 目錄結構
- 不修改任何 command 的內容，只搬移檔案

## Capabilities

### New Capabilities

- `no-capability-change`: Plugin command 檔案的目錄組織結構定義

### Modified Capabilities

（無）

## Impact

- 受影響檔案：`plugin/commands/` 下所有 `mysd-*.md` 檔案（約 20 個）
- 受影響系統：Claude Code plugin 載入路徑（但觸發方式不變）
- 不影響 Go binary、agents、或任何程式碼
