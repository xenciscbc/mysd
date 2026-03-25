# Phase 7: New Binary Commands & Scanner Refactor - Context

**Gathered:** 2026-03-25
**Status:** Ready for planning

<domain>
## Phase Boundary

使用者可透過 `/mysd:model` 顯示/切換 model profile，`/mysd:lang` 互動式同步 locale；`mysd scan` 升級為語言無關通用掃描器（language-agnostic via LLM），`mysd init` 展開為 `scan --scaffold-only`；plan 完成後 SKILL.md 層呈現 task↔skills 對應表供使用者批次確認。

</domain>

<decisions>
## Implementation Decisions

### Scanner 語言偵測架構
- **D-01:** 語言偵測使用 file-based markers：`go.mod` → Go、`package.json` → Node.js、`requirements.txt`/`pyproject.toml` → Python。未匹配到已知 marker 時 primary_language 標為 `unknown`
- **D-02:** ScanContext struct 完全替換為語言無關的通用 struct：`primary_language`（string）、`files`（副檔名統計）、`modules`（偵測到的 module/package 列表）、`existing_specs`。移除 Go-specific `PackageInfo` 陣列
- **D-03:** Scanner 的核心設計原則：binary 只收集 metadata，LLM agent 負責理解任意語言並生成 spec。Go binary 不做語言特定的 spec 決策（GSD 同樣模式—語言無關性來自 LLM）
- **D-04:** 保持 `--context-only` 執行模式。scan 輸出 JSON metadata，spec 寫入由 SKILL.md agent 負責。Binary 不新增 `--write-specs` flag

### init → scan --scaffold-only 遷移
- **D-05:** `mysd init` 內部直接展開為 `scan --scaffold-only` 執行（不顯示 deprecation warning，完全回展相容）。FSCAN-05 達成
- **D-06:** 首次建立 `openspec/config.yaml` 時的互動式 locale 詢問（FSCAN-04）在 **SKILL.md agent 層**發生：agent 詢問使用者後呼叫 `mysd lang set {locale}` 寫入。Go binary 的 scaffold-only 只建立空結構

### Skills 推薦 UX
- **D-07:** SKILL-01 的 skills 推薦邏輯在 **mysd-planner agent 層**（LLM 根據 task 內容推斷），Go binary 不實作規則式 skills 對映
- **D-08:** SKILL-02/03 的表格顯示與使用者確認流程在 **SKILL.md 層**：plan 完成後 SKILL.md 讀取 `--context-only` JSON，Claude 直接呈現 task↔skills 對應表並互動確認
- **D-09:** 批次同意 UX（SKILL-03）：呈現完整對應表後詢問 `Accept all recommended? Y/n`，預設 accept（Enter 即同意）。使用者選 n 才進入逐一調整流程
- **D-10:** `ffe` 模式（SKILL-04）跳過互動，直接使用 planner 推薦值

### /mysd:model 指令設計
- **D-11:** `mysd model`（讀）輸出 lipgloss table 格式：第一行顯示 `Profile: {name}`，接著 Role｜Model 兩欄表格，列出所有 agent role 的 resolved model（含 Phase 5/7 新增的 researcher/advisor/proposal-writer/plan-checker）
- **D-12:** `mysd model set <profile>` 在 **Go binary 層**直接寫入 `.claude/mysd.yaml` 的 `model_profile` 欄位（不透過 SKILL.md 中繼）

### 延續自前面階段的決策（適用 Phase 7）
- **Phase 5 D-09**：`/mysd:lang` 修改 locale 時，`openspec/config.yaml` 和 `mysd.yaml` 原子同步更新（兩者同時成功或同時不變）
- **Phase 5 D-08**：locale 使用 BCP47 標準格式（zh-TW, en-US, ja-JP）
- **Phase 5 D-09**：`openspec/config.yaml` 的 locale 為 source of truth

### Claude's Discretion
- ScanContext 的 `modules` 欄位具體結構（per-language module metadata 細節）
- `unknown` language 時 scan 回傳的 fallback metadata 格式
- `mysd model` table 的 lipgloss 樣式細節（顏色、對齊）
- `scan --scaffold-only` 建立的空結構具體目錄和文件列表
- `mysd model set` 的 profile 驗證邏輯（無效 profile 時的錯誤訊息）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### 現有 Scanner
- `internal/scanner/scanner.go` — 現有 Go-only ScanContext + BuildScanContext（需完全重構）
- `internal/scanner/scanner_test.go` — 現有測試（需同步更新）
- `cmd/scan.go` — scan subcommand（保持 --context-only 模式，更新 scanner 呼叫）

### 現有 Init 指令
- `cmd/init_cmd.go` — 現有 `mysd init` 實作（需改為呼叫 scan --scaffold-only）

### Config & Model Profile
- `internal/config/config.go` — DefaultModelMap（所有 roles）、ResolveModel function
- `internal/config/defaults.go` — ProjectConfig struct（含 ModelProfile、WorktreeDir、AutoMode）

### CLI Layer
- `cmd/plan.go` — --context-only JSON 輸出（SKILL.md 讀取 skills recommendations 的來源）
- `cmd/worktree.go` — worktree subcommand（`model` subcommand 應遵循相同 CLI 設計模式）

### v1.1 需求
- `.planning/REQUIREMENTS.md` — FCMD-03, FCMD-04, FCMD-05, FSCAN-01~05, SKILL-01~04 完整需求定義
- `.specs/changes/interactive-discovery/proposal.md` — v1.1 完整功能規格

### Phase 5 成果（lang 原子同步的基礎）
- `.planning/phases/05-schema-foundation-plan-checker/05-CONTEXT.md` — D-07~D-10 openspec/config.yaml 設計決策

### Agent Definitions
- `.claude/agents/mysd-planner.md` — 現有 planner agent（SKILL-01 skills 推薦在此層擴展）

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/output.Printer` — lipgloss table 輸出（`mysd model` table 直接使用）
- `internal/config.ResolveModel` — 已有完整 role→model resolve 邏輯，`mysd model` 直接呼叫
- `internal/config.Defaults()` — 現有 convention-over-config 模式，scaffold-only 遵循
- `cmd/worktree.go` — JSON stdout 模式可參考設計 `mysd model set` 的輸出

### Established Patterns
- **Binary-as-state-manager**: cmd/ 做參數解析，internal/ 業務邏輯，output/ 呈現
- **--context-only pattern**: scan 保持此模式，agent 消費 JSON
- **Viper config write**: `mysd model set` 寫 .claude/mysd.yaml 需用 viper.WriteConfig()
- **Atomic file write**: lang 同步兩個 config 文件需要 atomic（write-then-rename 或 defer rollback）

### Integration Points
- `cmd/root.go` — 新增 `model` subcommand（與 worktree 同層）
- `internal/scanner/scanner.go` — 完整替換（ScanContext struct + BuildScanContext function）
- `cmd/init_cmd.go` — runInit 改為呼叫 runScan with scaffold-only flag

</code_context>

<specifics>
## Specific Ideas

- Scanner 的語言無關性來自「binary 不做語意決策」這個原則，跟 GSD 工作流語言無關的原因相同（LLM 理解程式碼）
- `mysd model` table 格式已由使用者確認（Role | Model 兩欄，Profile 標題行）
- `accept all` 預設 accept 是關鍵 UX 決策：減少 ffe 以外場景的摩擦

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 07-new-binary-commands-scanner-refactor*
*Context gathered: 2026-03-25*
