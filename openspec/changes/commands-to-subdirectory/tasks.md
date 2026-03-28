## 1. 搬移 command 檔案到子目錄

- [x] 1.1 建立 `plugin/commands/mysd/` 子目錄
- [x] 1.2 將 19 個 `plugin/commands/mysd-*.md` 檔案搬到 `plugin/commands/mysd/`，去掉 `mysd-` prefix（例如 `mysd-apply.md` → `mysd/apply.md`）
- [x] 1.3 確認 `plugin/commands/CLAUDE.md` 維持在原位不搬移

## 2. 更新 marketplace plugin

- [x] 2.1 同步更新 `C:/Users/admin/.claude/plugins/marketplaces/mysd/plugin/commands/` 的結構，使其與開發版一致

## 3. 驗證

- [x] 3.1 確認 `plugin/commands/mysd/` 下有 19 個 `.md` 檔案，`plugin/commands/` 下不再有 `mysd-*.md` 檔案
- [x] 3.2 確認 `go build -o mysd.exe .` 編譯成功（確認 Go binary 不受影響）
