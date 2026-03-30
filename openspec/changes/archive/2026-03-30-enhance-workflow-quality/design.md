## Context

mysd 目前缺少 CLI 層面的 cross-artifact 分析。spectra 的 `spectra analyze` 提供 4 個 dimension 的靜態分析（Coverage、Consistency、Ambiguity、Gaps），在 propose 流程的 self-review 和 reviewer agent 之後執行，作為第三道品質防線。

現有品質防線：
1. AI inline self-review（已移除，改為 reviewer agent）
2. `mysd-reviewer` agent — AI self-check（Check 1-4）
3. `mysd validate` — 檢查 artifact 存在性和結構完整性

缺少的：CLI 靜態分析（不依賴 AI 判斷的結構性檢查）

## Goals / Non-Goals

**Goals:**

- `mysd analyze` 做 deterministic 的結構分析，不依賴 AI 判斷
- 4 個 dimension 的 findings 分級：Critical / Warning / Suggestion
- JSON output 供 SKILL.md 消費（plan pipeline analyze-fix loop）
- proposal-writer 根據 change_type 使用對應 template
- plan 可以跳過 design（optional artifact）

**Non-Goals:**

- 不做語意分析（那是 reviewer agent 的工作）
- 不做 code 層面的 lint（那是 golangci-lint 的工作）
- analyze 不修改 artifact（只報告，修改由 SKILL.md 的 fix loop 執行）

## Decisions

### D1: mysd analyze 的 4 個 Dimension

仿照 spectra analyze 的結構，但簡化為 mysd 的 artifact 格式：

| Dimension | 檢查內容 | 嚴重度 |
|-----------|---------|--------|
| **Coverage** | proposal 列出的每個 capability 都有對應 `specs/<name>/spec.md` | Critical |
| **Consistency** | proposal 的 file paths / component names 與 specs 一致；design 引用的 capability 都在 proposal 中；tasks 涵蓋所有 design decisions | Warning |
| **Ambiguity** | spec 中使用弱語言（should/may/might/TBD/TODO/FIXME）而非 SHALL/MUST | Suggestion |
| **Gaps** | spec 的 requirement 缺少 scenario；tasks 中有 task 但沒有對應的 spec requirement | Warning |

### D2: analyze 的 Go package 結構

```
internal/analyzer/
  analyzer.go      — public API: Analyze(changeDir string) → AnalysisResult
  coverage.go      — Coverage dimension
  consistency.go   — Consistency dimension
  ambiguity.go     — Ambiguity dimension
  gaps.go          — Gaps dimension
  types.go         — Finding, AnalysisResult, Severity 等共用型別
```

`cmd/analyze.go` 只做 CLI wiring：parse args → 呼叫 `analyzer.Analyze()` → JSON/styled output。

### D3: analyze output 的 JSON 格式

仿照 spectra analyze 的 output：

```json
{
  "change_id": "my-feature",
  "dimensions": [
    {"dimension": "Coverage", "status": "Clean", "finding_count": 0},
    {"dimension": "Consistency", "status": "2 issue(s) found", "finding_count": 2}
  ],
  "findings": [
    {
      "id": "CON-1",
      "dimension": "Consistency",
      "severity": "Warning",
      "location": "specs/foo/spec.md:15",
      "summary": "Capability 'bar' referenced in spec but not in proposal",
      "recommendation": "Add 'bar' to proposal Capabilities or remove from spec"
    }
  ],
  "artifacts_analyzed": ["proposal", "specs", "design", "tasks"],
  "artifacts_missing": []
}
```

### D4: Proposal template 按 type 切換

propose SKILL.md Step 9 傳入 `change_type`（已在 Step 2d 分類），proposal-writer agent 收到後選用對應 template：

| Type | Template 結構 |
|------|-------------|
| Feature | Why / What Changes / Capabilities / Impact |
| Bug Fix | Problem / Root Cause / Proposed Solution / Success Criteria / Impact |
| Refactor | Summary / Motivation / Proposed Solution / Alternatives Considered / Impact |

三種 template 直接寫在 proposal-writer agent definition 中，替換現有的固定 template。

### D5: Plan optional design skip

plan SKILL.md 在 Step 4（Design Phase）前加入判斷：

如果 change 滿足以下所有條件，跳過 design：
- proposal 的 Impact 只涉及 2 個以下檔案
- 沒有 New Capabilities（只有 Modified）
- proposal 沒有 "cross-cutting"、"migration"、"architecture" 等關鍵字

如果 auto_mode 為 false，顯示判斷結果並讓用戶確認是否跳過。
如果 auto_mode 為 true，自動跳過。

### D6: Plan analyze-fix loop

plan SKILL.md 在 Step 5b（reviewer）之後加入 Step 5c：

1. 執行 `mysd analyze {change_name} --json`
2. 過濾 Critical 和 Warning findings（忽略 Suggestion）
3. 如果有 findings：修復 → 重跑 analyze → 最多 2 輪
4. 如果 2 輪後仍有 findings：顯示 summary，不 block workflow

## Risks / Trade-offs

- **[Risk] Coverage 檢查需要 parse proposal 的 Capabilities section** → 用 regex 或 simple string parsing，不依賴完整的 markdown AST
- **[Risk] Optional design skip 的判斷邏輯可能太簡單，跳過不該跳的** → 非 auto_mode 下讓用戶確認，降低風險
- **[Trade-off] analyze 不做語意分析** → 刻意保持 deterministic，語意部分交給 reviewer agent
