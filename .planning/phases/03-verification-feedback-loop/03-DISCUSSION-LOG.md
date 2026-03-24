# Phase 3: Verification & Feedback Loop - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-24
**Phase:** 03-verification-feedback-loop
**Areas discussed:** Verification report 設計, Gap report 與 re-execute 迴路, Archive 行為與 UAT 關係, Verifier agent 獨立性, Spec status 回寫策略, UAT 清單的 UI 偵測邏輯, Verification 與 execute 的 state transition

---

## Verification Report 設計

| Option | Description | Selected |
|--------|-------------|----------|
| Terminal styled + markdown 檔 | Terminal 用 lipgloss 顯示摘要，同時寫入 verification.md 完整報告 | ✓ |
| JSON 純機器可讀 | 輸出結構化 JSON 供下游工具消費 | |
| 僅 terminal output | lipgloss styled terminal 輸出，不寫檔案 | |

**User's choice:** Terminal styled + markdown 檔 (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| 分級顯示：MUST 先、SHOULD 次、MAY 尾 | MUST 全通過才算 overall pass | ✓ |
| 混合顯示（按 spec 文件順序） | 按 spec 中出現順序顯示 | |
| 僅 MUST（最簡化） | 報告只列 MUST 項目 | |

**User's choice:** 分級顯示 (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| AI agent 綜合判斷 | 檢查 filesystem evidence，結合 spec 描述做出判斷 | ✓ |
| 純自動化檢查 | 每個 MUST 對應一個可執行的 check command | |
| Checklist 人工確認 | 產生 checklist 讓使用者手動確認 | |

**User's choice:** AI agent 綜合判斷 (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| 更新 spec frontmatter status | 通過的 MUST 設為 DONE，失敗的設 BLOCKED | ✓ |
| 另建 verification overlay 檔 | 不修改原 spec，在旁邊建立結果檔 | |
| Claude decides | | |

**User's choice:** 更新 spec frontmatter status (Recommended)

---

## Gap Report 與 Re-execute 迴路

| Option | Description | Selected |
|--------|-------------|----------|
| 失敗項目 + 原因 + 建議修正 | 每個失敗 MUST 列出描述、原因、修正建議 | ✓ |
| 僅失敗清單 | 只列哪些 MUST 失敗 | |
| Claude decides | | |

**User's choice:** 失敗項目 + 原因 + 建議修正 (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| 手動觸發 /mysd:execute | 使用者看完 gap report 後手動觸發 | ✓ |
| verify --fix 自動修正 | 失敗時自動觸發 re-execute | |
| 兩者皆支援 | 預設手動，--fix 可自動 | |

**User's choice:** 手動觸發 (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| 僅失敗項目 | 只 re-execute 失敗的 MUST items | ✓ |
| 全部重新執行 | 全部 tasks 重新執行 | |
| Claude decides | | |

**User's choice:** 僅失敗項目 (Recommended)

---

## Archive 行為與 UAT 關係

| Option | Description | Selected |
|--------|-------------|----------|
| 移至 .specs/archive/{name}/ | 整個 change 目錄移到 archive | ✓ |
| 刪除 change 目錄 | 不保留歷史 | |
| Claude decides | | |

**User's choice:** 移至 .specs/archive/ (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| 提示但不強制 | 顯示 'Run UAT first?'，無論選什麼都完成 archive | ✓ |
| 無互動，直接 archive | 完全不提及 UAT | |
| 強制 UAT 完成後才能 archive | UAT 必須全部通過 | |

**User's choice:** 提示但不強制 (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| verify 過程中自動產生 | 偵測 UI 相關 MUST/SHOULD 自動產生 UAT checklist | ✓ |
| 僅手動觸發 /mysd:uat | 永遠不自動產生 | |
| Claude decides | | |

**User's choice:** verify 過程中自動產生 (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| Agent 引導使用者逐項確認 | /mysd:uat 由專屬 agent 引導互動 | ✓ |
| 純 checklist 檔案 | 產生靜態 markdown 讓使用者自行編輯 | |
| Claude decides | | |

**User's choice:** Agent 引導 (Recommended)

---

## Spec Status 回寫策略

| Option | Description | Selected |
|--------|-------------|----------|
| Go binary 批量更新 | verify 完成後一次性讀取 report 並更新 spec status | ✓ |
| Verifier agent 逐項寫回 | 每確認一項就呼叫 binary 更新 | |
| Claude decides | | |

**User's choice:** Go binary 批量更新 (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| 只更新 MUST | SHOULD/MAY 不寫回 spec status | ✓ |
| 全部更新 | MUST/SHOULD/MAY 全部寫回 | |
| Claude decides | | |

**User's choice:** 只更新 MUST (Recommended)

---

## UAT 清單的 UI 偵測邏輯

| Option | Description | Selected |
|--------|-------------|----------|
| AI agent 判斷 | Verifier agent 用 AI 判斷哪些項目涉及使用者可見行為 | ✓ |
| 關鍵字匹配 | 用 'UI', 'display' 等關鍵字自動過濾 | |
| Spec 中明確標註 | 使用者在 spec 中用 tag 標註 'uat: true' | |

**User's choice:** AI agent 判斷 (Recommended)

---

## Verification 與 Execute 的 State Transition

| Option | Description | Selected |
|--------|-------------|----------|
| 全通過才 transition | MUST 全通過：executed → verified | ✓ |
| 新增 'verified_with_gaps' state | 部分通過轉換到 intermediate state | |
| Claude decides | | |

**User's choice:** 全通過才 transition (Recommended)

| Option | Description | Selected |
|--------|-------------|----------|
| Binary 檢查 verified state + MUST 狀態 | 雙重 gate：state 和 item status 都要滿足 | ✓ |
| 僅檢查 state == verified | 只看 state，信任 verify 結果 | |
| Claude decides | | |

**User's choice:** 雙重檢查 (Recommended)

---

## Claude's Discretion

- Verification report 的具體 markdown 模板
- Gap report 的詳細程度和建議修正的具體寫法
- UAT checklist 的格式和測試步驟描述
- mysd-verifier.md 和 mysd-uat-guide.md 的具體 prompt 措辭
- archive 時 'Run UAT first?' 的呈現方式

## Deferred Ideas

- **Playwright 腳本自動產生** — UAT 若偵測到 web UI 項目，可自動產生 Playwright e2e test 腳本。屬未來 phase 功能
