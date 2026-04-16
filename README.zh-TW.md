# mysd

**以規格驅動開發的內容智慧 skill。**

mysd 是一個 Claude Code plugin，提供 4 個獨立的 skill 來管理程式碼庫中的知識文件：研究決策、同步文件、撰寫規格、以及一鍵同步全部。

零依賴。純 SKILL.md 檔案。安裝即用。

## Skills

| 命令 | Skill | 功能 |
|------|-------|------|
| `/mysd:research` | Research | 灰區決策，附帶證據。Spec 健康檢查（4 維度分析）。 |
| `/mysd:doc` | Doc Writer | 偵測程式碼變更、識別受影響的文件、自動更新並支援多語言同步。 |
| `/mysd:spec` | Spec Writer | 撰寫/更新 [OpenSpec](https://github.com/openspec-dev/openspec) 格式的 spec 檔案。支援從程式碼反推 spec。 |
| `/mysd:sync` | Sync | 透過 subagent 依序執行三個 skill。一個命令，完整內容同步。 |

每個 skill 獨立運作。可以只用一個、兩個、或全部四個。

## 安裝

### Claude Code Plugin

```bash
claude plugin add https://github.com/xenciscbc/mysd
```

### 手動安裝

將 skill 複製到你的 Claude Code skills 目錄：

```bash
cp -r mysd/skills/research ~/.claude/skills/mysd-research
cp -r mysd/skills/doc ~/.claude/skills/mysd-doc
cp -r mysd/skills/spec ~/.claude/skills/mysd-spec
cp -r mysd/skills/sync ~/.claude/skills/mysd-sync
```

所有 `/mysd:*` 命令將在下次 Claude Code session 中可用。

## 使用方式

### Research：在灰區中做決策

```
/mysd:research 新的快取層應該用哪個資料庫？
```

- 分類問題（是否為灰區）
- 從程式碼、git 歷史、文件、網路蒐集上下文
- 整理 2-4 個選項，附帶證據、優缺點、工作量
- 產出 Decision Doc，含信心分數（1-10）

也可以執行 **Spec 健康檢查**（4 維度：覆蓋率、模糊性、一致性、缺口）：

```
/mysd:research 檢查 auth-refactor 變更的 spec 品質
```

### Doc：保持文件同步

```
/mysd:doc 根據最新的變更更新文件
```

- 透過 `git diff` 偵測變更
- 將變更類型對應到受影響的文件（新命令 -> README、bug 修復 -> CHANGELOG 等）
- 匹配目標文件的現有風格
- 多語言同步：更新 README.md 時自動同步 README.zh-TW.md
- 套用前逐一確認

### Spec：撰寫 OpenSpec 格式的規格

```
/mysd:spec 幫新的 auth middleware 寫 spec
```

- 產出正確的 YAML frontmatter（`spec-version`、`capability`、`delta`、`status`）
- 使用 RFC 2119 關鍵字（MUST、SHOULD、MAY）
- 撰寫 WHEN/THEN/AND 場景
- 可從程式碼反推 spec：讀取 Go 檔案，從函式簽名推斷需求

### Sync：執行完整流程

```
/mysd:sync 我剛完成新的快取功能，幫我同步所有內容
```

透過 subagent 串聯 research -> doc -> spec，每個 skill 擁有獨立的 context window。

## OpenSpec 相容性

spec writer 產出的檔案相容 [OpenSpec](https://github.com/openspec-dev/openspec) 格式：

- YAML frontmatter：`spec-version`、`capability`、`delta`（ADDED/MODIFIED/REMOVED/RENAMED）、`status`
- RFC 2119 關鍵字使用大寫
- WHEN/THEN/AND 場景格式
- 目錄結構：`openspec/specs/{capability}/spec.md`

## 系統需求

- Claude Code（支援 skill 的任何版本）
- 不需要 binary、不需要編譯、零執行時依賴

## 授權

MIT
