## 1. Proposal Template 按 Type 切換

- [x] 1.1 修改 propose SKILL.md Step 9，將 `change_type` 加入 proposal-writer invocation context — D4: Proposal template 按 type 切換
- [x] 1.2 修改 mysd-proposal-writer.md，加入三種 template（Feature: Why/What Changes/Capabilities/Impact、Bug Fix: Problem/Root Cause/Proposed Solution/Success Criteria/Impact、Refactor: Summary/Motivation/Proposed Solution/Alternatives Considered/Impact），根據 `change_type` 選用 — D4: Proposal template 按 type 切換
- [x] 1.3 修改 mysd-reviewer.md，新增 `change_type` input context 欄位，增加 proposal template 驗證邏輯（Check 5: Template Match）— Reviewer agent performs artifact quality checks、Reviewer validates proposal template matches change type

## 2. mysd analyze CLI Command

- [x] 2.1 建立 `internal/analyzer/types.go`：定義 Finding、AnalysisResult、Severity、DimensionResult 型別 — D3: analyze output 的 JSON 格式、Analyze command outputs structured JSON
- [x] 2.2 建立 `internal/analyzer/coverage.go`：實作 Coverage dimension — parse proposal Capabilities section，比對 specs/ 目錄下的 spec 檔 — D1: mysd analyze 的 4 個 dimension、Analyze command performs cross-artifact structural analysis
- [x] 2.3 建立 `internal/analyzer/ambiguity.go`：實作 Ambiguity dimension — scan spec 檔的弱語言 pattern（should/may/might/TBD/TODO/FIXME） — D1: mysd analyze 的 4 個 dimension、Analyze command performs cross-artifact structural analysis
- [x] 2.4 建立 `internal/analyzer/consistency.go`：實作 Consistency dimension — cross-check proposal/specs/design/tasks 的 file paths 和 capability references — D1: mysd analyze 的 4 個 dimension、Analyze command performs cross-artifact structural analysis
- [x] 2.5 建立 `internal/analyzer/gaps.go`：實作 Gaps dimension — check requirements 有 scenario、tasks 引用 spec requirements — D1: mysd analyze 的 4 個 dimension、Analyze command performs cross-artifact structural analysis
- [x] 2.6 建立 `internal/analyzer/analyzer.go`：public API `Analyze(changeDir string) AnalysisResult`，串接 4 個 dimension — D2: analyze 的 Go package 結構、Analyze operates on available artifacts only
- [x] 2.7 建立 `cmd/analyze.go`：Cobra command `mysd analyze [change-name] [--json]`，wiring CLI args → analyzer.Analyze() → JSON or lipgloss styled output — D2: analyze 的 Go package 結構、Analyze command provides styled terminal output
- [x] 2.8 建立 `cmd/analyze_test.go`：測試 JSON output 格式、styled output、missing change name error handling
- [x] 2.9 建立 `internal/analyzer/analyzer_test.go`：unit tests for 4 dimensions — Coverage missing spec、Ambiguity weak language、Consistency cross-check、Gaps missing scenario

## 3. Plan Pipeline 增強

- [x] 3.1 修改 plan SKILL.md Step 4 前加入 design skip 判斷邏輯：檢查 proposal Impact 檔案數 ≤ 2、無 New Capabilities、無 cross-cutting 關鍵字 — D5: Plan optional design skip、Plan pipeline supports optional design skip
- [x] 3.2 修改 plan SKILL.md Step 5b 之後加入 Step 5c analyze-fix loop：`mysd analyze --json` → 過濾 Critical/Warning → fix → 最多 2 輪 — D6: Plan analyze-fix loop、Plan pipeline includes analyze-fix loop after reviewer
