# Phase 1: Foundation - Context

**Gathered:** 2026-03-23
**Status:** Ready for planning

<domain>
## Phase Boundary

建立 my-ssd 的基礎架構：Go CLI 骨架（Cobra）、OpenSpec 格式 parser（讀寫 proposal.md / specs/ / design.md / tasks.md）、RFC 2119 keyword parser、spec 狀態機（proposed → specced → designed → planned → executed → verified → archived）、專案設定檔管理（`.claude/mysd.yaml`）、以及基本的 `mysd propose` 和 `mysd init` 指令。

</domain>

<decisions>
## Implementation Decisions

### 目錄結構
- **D-01:** Spec 檔案存放在專案根目錄的 `.specs/` 目錄下（相容 OpenSpec 的 `openspec/` 目錄）
- **D-02:** 全域追蹤資訊（roadmap 歷史、UAT 結果等）存放在 `.mysd/` 目錄
- **D-03:** `.mysd/roadmap/` 目錄記錄每個 change 的名稱、狀態、完成日期時間，格式需可被第三方工具解析（用於 roadmap 視覺化）
- **D-04:** 設定檔存放在 `.claude/mysd.yaml`

### CLI 指令設計
- **D-05:** 指令風格為 `mysd <verb>` — 簡潔直觀（如 `mysd propose`, `mysd verify`, `mysd ff`）
- **D-06:** 使用 Cobra CLI 框架
- **D-07:** 輸出使用彩色終端輸出（lipgloss），TTY 偵測自動降級為純文字

### 狀態管理
- **D-08:** Spec 狀態存放在每個 spec 檔案的 YAML frontmatter 中（跟著檔案走）
- **D-09:** Frontmatter 包含 `spec-version` 欄位用於 schema 版本控制（forward compatibility）

### Propose 指令行為
- **D-10:** `mysd propose` 先進行 GSD 式的互動提問（了解使用者想做什麼）
- **D-11:** 提問結束後，可選擇是否使用 agent 進行領域研究（像 GSD 的研究模式）
- **D-12:** 最後產出完整的 spec artifacts（proposal.md / specs/ / design.md / tasks.md）
- **D-13:** 實作完成後，自動產生或更新 `.mysd/roadmap/` 下的文件

### 新增指令
- **D-14:** 新增 `/mysd:capture` 指令 — 從當前對話中分析並提取要做的變更，然後自動進入 propose 的討論模式（不需重新描述需求）

### OpenSpec 相容性
- **D-15:** Parser 自動偵測 `openspec/` 或 `.specs/` 目錄
- **D-16:** 支援讀寫 OpenSpec 的完整 artifact 結構（proposal.md, specs/, design.md, tasks.md）
- **D-17:** Delta Specs 語義（ADDED / MODIFIED / REMOVED）在 parser 層就被識別

### Claude's Discretion
- 具體的 frontmatter schema 欄位設計
- lipgloss 的配色方案
- 錯誤訊息的措辭風格
- `.mysd/roadmap/` 的具體檔案格式（JSON / YAML / Markdown，只要能被第三方工具讀取即可）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### OpenSpec 格式規範
- 需參考 OpenSpec 官方文件了解 artifact 結構、frontmatter 格式、delta spec 語法
- RFC 2119 (https://datatracker.ietf.org/doc/html/rfc2119) — MUST/SHOULD/MAY 關鍵字定義

### 研究文件
- `.planning/research/STACK.md` — Go 技術棧建議（Cobra, adrg/frontmatter, yaml.v3, lipgloss）
- `.planning/research/ARCHITECTURE.md` — 系統架構設計（三層架構、component boundaries）
- `.planning/research/PITFALLS.md` — 關鍵陷阱（spec format hardcoding、skill context budget）

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- 無 — 全新 greenfield 專案

### Established Patterns
- 無既有模式 — Phase 1 將建立所有基礎模式

### Integration Points
- Claude Code plugin 層（`commands/`, `agents/`）將在 Phase 4 建立
- Go binary 在 Phase 1 只需要 CLI 和 spec I/O 能力

</code_context>

<specifics>
## Specific Ideas

- Propose 指令要有 GSD 的 deep questioning 體驗 — 不是表單填寫，而是對話式探索
- Roadmap 追蹤需可被第三方工具讀取 — 例如能匯出成 Mermaid gantt chart 或類似格式
- `/mysd:capture` 是一個重要的 UX 創新 — 從對話脈絡直接進入 SDD 流程，減少重複描述

</specifics>

<deferred>
## Deferred Ideas

- Debug session 功能 — 可考慮在 v1.x 加入
- Todo / Notes 管理 — 可考慮在 v1.x 加入
- Session pause/resume — 可考慮在 v1.x 加入

</deferred>

---

*Phase: 01-foundation*
*Context gathered: 2026-03-23*
