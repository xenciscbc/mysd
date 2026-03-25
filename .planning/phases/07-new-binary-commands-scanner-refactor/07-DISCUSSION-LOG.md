# Phase 7: Discussion Log

**Session:** 2026-03-25
**Phase:** 07 — New Binary Commands & Scanner Refactor

---

## Area Selection

**Question:** Phase 7 有哪些灰色地帶要討論？
**Selected:** Scanner 語言偵測架構、init → scan --scaffold-only 遷移、Skills 推薦 UX、/mysd:model 輸出格式（全選）

---

## Scanner 語言偵測架構

**Q:** Scanner 語言偵測的主要指標是什麼？
**A:** File-based markers（推薦）— go.mod → Go, package.json → Node.js, requirements.txt/pyproject.toml → Python

**Q:** 現有 ScanContext struct（Go-specific）如何重構？
**A:** 替換為通用 struct（推薦）— 新的 ScanContext 包含 primary_language, files/ext stats, modules，移除 PackageInfo

**Q:** scan 的 spec 產生策略？
**A:** Agent 自由發揮（推薦）— binary 只提供 metadata，SKILL.md agent 根據 context 自行決定 spec 粒度

**User clarification:** GSD 工作流也是語言無關的（取決於 LLM 而非工具本身）。mysd:scan 應遵循相同原則——binary 不做語意決策，LLM 理解任意語言。

**Q:** scan 的執行模式：重構後如何寫入 openspec/specs/？
**A:** 保持 --context-only（推薦）— binary 只輸出 JSON，spec 寫入由 SKILL.md agent 負責

---

## init → scan --scaffold-only 遷移

**Q:** `mysd init` 遷移為 `scan --scaffold-only`，舊指令如何處理？
**A:** init 直接展開為 scaffold-only scan（推薦）— 內部呼叫，無 warning，完全回展相容

**Q:** FSCAN-04：首次建立 config.yaml 時互動詢問 locale，在哪一層？
**A:** SKILL.md agent 層（推薦）— agent 詢問後呼叫 `mysd lang set {locale}`，符合現有模式

---

## Skills 推薦 UX

**Q:** SKILL-01：planner 推薦 skills 的邏輯在哪一層？
**A:** mysd-planner agent 層（推薦）— LLM 根據 task 內容推斷，Go binary 不做規則式對映

**Q:** SKILL-02/03：表格顯示 + 使用者確認流程在哪一層？
**A:** SKILL.md 層互動（推薦）— plan 完成後 SKILL.md 讀取 context-only JSON 呈現表格

**Q:** 「批次同意」的具體互動形式？
**A:** 單一按鍵 accept-all（推薦）— `Accept all recommended? Y/n`，預設 accept（Enter 即同意）

---

## /mysd:model 輸出格式

**Q:** `mysd model`（讀）的輸出格式？
**A:** Table 格式（推薦）— lipgloss 表格：Profile: xxx 標題行 + Role | Model 兩欄

```
Profile: quality

Role             Model
───────────────────────────────
 executor        claude-opus-4-6
 planner         claude-opus-4-6
 verifier        claude-sonnet-4-6
 researcher      claude-sonnet-4-6
 advisor         claude-sonnet-4-6
 plan-checker    claude-sonnet-4-6
```

**Q:** `mysd model set <profile>` 的實作層次？
**A:** Go binary 寫入 mysd.yaml（推薦）— 直接更新 .claude/mysd.yaml 的 model_profile 欄位

---

## Outcome

CONTEXT.md created: `.planning/phases/07-new-binary-commands-scanner-refactor/07-CONTEXT.md`
Next step: `/gsd:plan-phase 7`
