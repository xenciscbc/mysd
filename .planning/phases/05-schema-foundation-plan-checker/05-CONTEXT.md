# Phase 5: Schema Foundation & Plan-Checker - Context

**Gathered:** 2026-03-25
**Status:** Ready for planning

<domain>
## Phase Boundary

擴展 TaskEntry schema 支援依賴追蹤（depends）、檔案衝突偵測（files）、需求可追溯性（satisfies）、技能建議（skills），新增 plan-checker 基礎設施驗證 MUST 覆蓋率，擴展 model profile 支援 4 個新 agent role，並實作 openspec/config.yaml 的產生/讀取。

</domain>

<decisions>
## Implementation Decisions

### Plan-checker 架構
- **D-01:** MUST 覆蓋率檢查在 Go binary 層實作（`internal/planchecker/` 新 package），使用 structured ID matching（satisfies 欄位 vs MUST IDs），不使用 AI 語意推測
- **D-02:** 觸發方式沿用現有 `mysd plan --check` flag，plan 完成後自動執行檢查
- **D-03:** 檢查失敗時輸出 uncovered MUST IDs 清單 + 覆蓋率比例（簡潔模式），agent 層負責渲染互動 UI 讓使用者選擇自動補齊或手動調整（對應 FSCHEMA-06）
- **D-04:** PlanningContext JSON 擴展 WaveGroups、WorktreeDir、AutoMode 欄位，WaveGroups 由 Go binary 對 tasks.md 的 depends 欄位做 topological sort 計算得出

### Model profile 對應表
- **D-05:** 4 個新 agent roles（researcher, advisor, proposal-writer, plan-checker）在 quality/balanced/budget 三層全部映射到 sonnet，與現有 6 roles 保持一致
- **D-06:** budget 層新 roles 也用 sonnet（不降為 haiku），保持整體一致性

### openspec/config.yaml 內容
- **D-07:** 最小必要欄位：project name, locale (BCP47), spec_dir, created — 其餘由 convention-over-config 推導
- **D-08:** locale 欄位使用 BCP47 標準格式（zh-TW, en-US, ja-JP），與 Go 的 golang.org/x/text/language 直接相容
- **D-09:** openspec/config.yaml 的 locale 為 source of truth，mysd.yaml 的 response_language/document_language 讀取時參考。/mysd:lang 修改時兩者原子同步
- **D-10:** config.yaml 是 OpenSpec 標準格式（project-level），mysd.yaml 是 mysd 專用配置。兩者職責分離

### Schema 向後相容策略
- **D-11:** 新欄位使用 `omitempty` YAML tag（`yaml:"depends,omitempty"` 等），舊 tasks.md 讀取時新欄位為 nil/empty，寫回時不輸出空欄位。零遮修、零 migration
- **D-12:** TaskItem（executor/context.go JSON 輸出）同步擴展 Depends/Files/Satisfies/Skills 欄位，讓 agent 執行時能看到完整資訊，Phase 6 wave grouping 直接可用

### Claude's Discretion
- plan-checker 的具體 JSON 輸出結構（遵循簡潔模式：uncovered IDs + 覆蓋率）
- topological sort 演算法選擇（Kahn's vs DFS-based）
- config.yaml writer 的錯誤處理細節
- 新欄位的測試案例覆蓋範圍

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### 現有 Schema 定義
- `internal/spec/schema.go` — TaskEntry struct（L76-90）、TasksFrontmatterV2（L85-90）、ItemStatus enum（L24-32）、RFC2119Keyword（L5-12）
- `internal/spec/updater.go` — ParseTasksV2（L79-98）、WriteTasks（L134-147）、UpdateTaskStatus（L102-130）

### Executor & Verifier Context
- `internal/executor/context.go` — ExecutionContext + TaskItem structs（L13-39）、BuildContext function（L92-112）
- `internal/verifier/context.go` — VerificationContext + VerifyItem + StableID generation（L17-50）

### Config & Model Profile
- `internal/config/config.go` — DefaultModelMap（L15-40）、ResolveModel function（L45-57）
- `internal/config/defaults.go` — ProjectConfig struct（L3-29）

### CLI Layer
- `cmd/plan.go` — --check flag 定義（L18-35）、PlanningContext JSON 輸出（L55-81）

### v1.1 Proposal & Research
- `.specs/changes/interactive-discovery/proposal.md` — v1.1 完整功能規格
- `.planning/research/ARCHITECTURE.md` — v1.1 架構研究（含 internal/planchecker/ 建議）
- `.planning/research/SUMMARY.md` — v1.1 研究摘要

### Agent Definitions
- `.claude/agents/mysd-planner.md` — 現有 planner agent（plan-check 觸發點）

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/spec.TaskEntry` struct: 直接擴展加入新欄位，YAML round-trip 已建立
- `internal/spec.ParseTasksV2`: 支援 brownfield 格式，新欄位 omitempty 自動向後相容
- `internal/config.ResolveModel`: 直接擴展 DefaultModelMap 加入新 roles
- `internal/config.ProjectConfig`: 已有 ModelProfile + ModelOverrides 機制
- `internal/verifier/report.go`: VerifierReport 模式可參考設計 plan-checker 報告結構
- `internal/verifier/context.go StableID`: CRC32 based ID 生成模式可參考

### Established Patterns
- **Thin command layer**: cmd/*.go 做參數解析 + 呼叫 internal/ + 用 Printer 輸出
- **Convention-over-config**: 缺少設定檔時用 Defaults()
- **Instance viper**: internal/config 用 viper.New()（測試隔離）
- **Sidecar pattern**: VerificationStatus 用 JSON sidecar（verification-status.json），不修改 spec 檔案

### Integration Points
- `cmd/plan.go --check`: 已有 flag，需要實際呼叫 planchecker 邏輯
- `cmd/plan.go --context-only`: JSON 輸出需擴展 WaveGroups/WorktreeDir/AutoMode
- `internal/executor/context.go BuildContext`: 需要讀取新的 TaskEntry 欄位並映射到 TaskItem

</code_context>

<specifics>
## Specific Ideas

- Plan-checker 放在 Go binary 而非 agent 層 — 確保檢查結果 deterministic、可測試、速度快
- config.yaml 的 locale 作為 source of truth — 避免 mysd.yaml 和 config.yaml 不同步的問題
- 新欄位用 omitempty 而非 migration — 最簡單的向後相容方式，符合 convention-over-config 哲學

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 05-schema-foundation-plan-checker*
*Context gathered: 2026-03-25*
