# Phase 3: Verification & Feedback Loop - Context

**Gathered:** 2026-03-24
**Status:** Ready for planning

<domain>
## Phase Boundary

實作 goal-backward 驗證引擎：`mysd verify` 觸發獨立 verifier agent 驗證所有 MUST 項目（SHOULD/MAY 以較低優先級呈現），結果自動寫回 spec status，`mysd archive` 在 MUST 全通過後歸檔至 `.specs/archive/`。驗證過程中自動偵測 UI 相關項目並產生 UAT checklist，`/mysd:uat` 提供互動式驗收流程。

</domain>

<decisions>
## Implementation Decisions

### Verification Report 設計
- **D-01:** Terminal styled + markdown 檔雙輸出 — Terminal 用 lipgloss 顯示摘要（MUST 5/7 passed），同時寫入 `.specs/changes/{name}/verification.md` 完整報告。仿 GSD 的 VERIFICATION.md 模式
- **D-02:** 分級顯示 — MUST 先、SHOULD 次、MAY 尾。MUST 全通過才算 overall pass，SHOULD 有警告但不 block，MAY 只備註
- **D-03:** AI agent 綜合判斷 pass/fail — Verifier agent 檢查 filesystem evidence（檔案存在、grep 關鍵字、test 通過），結合 spec 描述做出判斷。彈性最大
- **D-04:** 驗證結果寫回 spec frontmatter — 通過的 MUST 設為 DONE，失敗的設為 BLOCKED。用 Go binary 批量更新（verify 完成後一次性讀取 report 並更新），僅更新 MUST 項目

### Gap Report 與 Re-execute 迴路
- **D-05:** 失敗項目 + 原因 + 建議修正 — 每個失敗的 MUST 列出：spec 描述、失敗原因（evidence）、AI 建議的修正方向。寫入 `.specs/changes/{name}/gap-report.md` 供 re-execute 使用
- **D-06:** 手動觸發 re-execute — 使用者看完 gap report 後手動執行 `/mysd:execute`，執行時自動讀取 gap report 並只修正失敗項目
- **D-07:** Re-execute scope 僅限失敗項目 — 只 re-execute 失敗的 MUST items 對應的 tasks，已通過的不重做

### Archive 行為與 UAT 關係
- **D-08:** archive 移至 `.specs/archive/{name}/` — 將整個 change 目錄移到 archive，更新 STATE.json phase 為 archived，保留完整歷史紀錄
- **D-09:** archive 時提示但不強制 UAT — 顯示 'Run UAT first?'，使用者可選 yes/no，無論選什麼都會完成 archive。UAT 永遠是 optional
- **D-10:** UAT 清單在 verify 過程中自動產生 — verifier agent 偵測 spec 中有 UI 相關 MUST/SHOULD 時，自動產生 UAT checklist 到 `.mysd/uat/`，不阻塞 verify 流程
- **D-11:** UAT 互動由專屬 agent 引導 — `/mysd:uat` 由 mysd-uat-guide.md agent 引導，逐項顯示測試步驟，使用者回報 pass/fail/skip，結果寫回 UAT 檔案

### Verifier Agent 獨立性
- **D-12:** 全新 agent，只讀 spec 和 filesystem — mysd-verifier.md 不讀 executor 的 alignment.md 或執行歷史，僅根據 spec MUST 項目 + 實際檔案系統證據判斷。最強獨立性
- **D-13:** 多層次證據 — 檔案存在 + grep 關鍵字碼 + 執行 test suite + 檢查 build 通過。多維度驗證確保可靠性
- **D-14:** Go binary 輸出 verification context JSON — `mysd verify --context-only` 輸出 spec MUST/SHOULD/MAY 清單、當前狀態、tasks 完成情況的 JSON，SKILL.md 傳給 verifier agent。與 execute --context-only 同模式

### UAT UI 偵測邏輯
- **D-15:** AI agent 判斷 UI 相關性 — Verifier agent 在驗證過程中用 AI 判斷哪些 MUST/SHOULD 涉及使用者可見的行為（UI 顯示、互動、排版），最彈性，無需特殊標註

### State Transition 設計
- **D-16:** MUST 全通過才 transition — MUST 全通過：executed → verified。有失敗：維持 executed 狀態（不新增 intermediate state），用戶可 re-execute 後再 verify。用現有 state machine
- **D-17:** archive gate 雙重檢查 — Go binary 在 archive 前檢查：(1) state == verified，(2) 所有 MUST items status == DONE。兩個條件都滿足才允許 archive

### Claude's Discretion
- Verification report 的具體 markdown 模板
- Gap report 的詳細程度和建議修正的具體寫法
- UAT checklist 的格式和測試步驟描述
- mysd-verifier.md 的具體 prompt 措辭
- mysd-uat-guide.md 的互動 prompt 設計
- archive 時 'Run UAT first?' 的呈現方式

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### 專案架構
- `.planning/PROJECT.md` — 專案願景、約束條件、核心價值
- `.planning/REQUIREMENTS.md` — 完整需求清單，Phase 3 需覆蓋 VRFY-01~05, SPEC-05, SPEC-06, WCMD-06, WCMD-07, WCMD-12, UAT-01~05
- `.planning/ROADMAP.md` — Phase 3 goal 和 success criteria

### Phase 2 基礎（直接依賴）
- `.planning/phases/02-execution-engine/02-CONTEXT.md` — Phase 2 架構決策（alignment gate、plugin reverse-calling、model profiles）
- `internal/spec/schema.go` — Requirement struct 含 Keyword (MUST/SHOULD/MAY) 和 ItemStatus (PENDING/IN_PROGRESS/DONE/BLOCKED)
- `internal/spec/updater.go` — TasksFrontmatterV2 round-trip 更新機制（可延伸為 spec item status 更新）
- `internal/executor/context.go` — ExecutionContext JSON 結構（verification context 需參考此模式）
- `internal/executor/alignment.go` — AlignmentPath 和 AlignmentTemplate（verifier 不讀但需知其存在）
- `internal/state/transitions.go` — ValidTransitions：PhaseExecuted → PhaseVerified → PhaseArchived（已支援 re-execute）
- `cmd/verify.go` — verify 指令 stub（待實作）
- `cmd/archive.go` — archive 指令 stub（待實作）

### Claude Code Plugin 文件
- `.claude/agents/mysd-executor.md` — executor agent definition（Phase 3 新增 mysd-verifier.md 和 mysd-uat-guide.md 需參考此格式）
- `.claude/commands/mysd-execute.md` — execute SKILL.md（verify SKILL.md 需參考此模式）

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/spec/schema.go` Requirement.Keyword — 直接可用於 MUST/SHOULD/MAY 分類
- `internal/spec/updater.go` UpdateTaskStatus — 可延伸為 UpdateItemStatus 批量更新 spec requirement status
- `internal/executor/context.go` BuildContext — verification context JSON 可參考此模式建立 VerificationContext
- `internal/executor/status.go` RenderStatus — lipgloss 渲染模式可複用於 verification terminal output
- `internal/state/transitions.go` — PhaseExecuted → PhaseVerified 轉換已內建

### Established Patterns
- `cmd/execute.go` --context-only 模式 — verify 和 archive 應遵循相同的 thin-command-layer 模式
- Plugin reverse-calling — SKILL.md 入口 → Go binary context → agent definition 執行
- `internal/spec/writer.go` Scaffold — archive 可參考此模式建立 archive 目錄結構

### Integration Points
- `cmd/verify.go` stub — 已註冊在 rootCmd，需實作 RunE
- `cmd/archive.go` stub — 已註冊在 rootCmd，需實作 RunE
- `.claude/commands/` — 需新增 mysd-verify.md 和 mysd-uat.md SKILL.md
- `.claude/agents/` — 需新增 mysd-verifier.md 和 mysd-uat-guide.md agent definitions

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches

</specifics>

<deferred>
## Deferred Ideas

- **Playwright 腳本自動產生** — UAT 清單如果偵測到 web UI 相關項目，可自動產生 Playwright e2e test 腳本。屬於 test automation 工具整合，超出 Phase 3 核心驗證流程範圍，可作為未來 phase 功能

</deferred>

---

*Phase: 03-verification-feedback-loop*
*Context gathered: 2026-03-24*
