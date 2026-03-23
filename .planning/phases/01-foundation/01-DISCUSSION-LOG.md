# Phase 1: Foundation - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-23
**Phase:** 01-foundation
**Areas discussed:** 目錄結構, CLI 指令設計, 狀態管理, Propose 行為, GSD 功能比對

---

## 目錄結構

| Option | Description | Selected |
|--------|-------------|----------|
| .specs/ (OpenSpec 相容) | 跟 OpenSpec 一樣用頂層目錄，簡影直觀 | ✓ |
| .mysd/ 統一管理 | 所有 my-ssd 相關檔案都在 .mysd/ 下 | |
| 你決定 | 讓 Claude 選擇最合適的結構 | |

**User's choice:** .specs/ (OpenSpec 相容)
**Notes:** 使用者明確希望保持 OpenSpec 相容性

---

## CLI 指令風格

| Option | Description | Selected |
|--------|-------------|----------|
| mysd \<verb\> | 簡潔風格：mysd propose, mysd verify, mysd ff | ✓ |
| mysd \<noun\> \<verb\> | Git 風格：mysd spec create, mysd change verify | |
| mysd \<verb\> --flags | cobra 標準風格，加 subcommands 和 flags | |

**User's choice:** mysd \<verb\>
**Notes:** 簡潔直觀

---

## 狀態存放方式

| Option | Description | Selected |
|--------|-------------|----------|
| Frontmatter 內嵌 | 狀態存在每個 spec 檔案的 YAML frontmatter 中 | ✓ |
| .mysd/state.json | 集中式 JSON 檔案 | |
| 兩者並用 | Frontmatter 存單一 spec 狀態 + state.json 存全域索引 | |

**User's choice:** Frontmatter 內嵌
**Notes:** 跟著檔案走

---

## CLI 輸出風格

| Option | Description | Selected |
|--------|-------------|----------|
| 彩色終端輸出 | 用 lipgloss 等工具做彩色、表格、進度條 | ✓ |
| 純文字優先 | 簡潔純文字，可被 pipe 處理 | |
| 你決定 | Claude 選擇最合適的輸出方式 | |

**User's choice:** 彩色終端輸出

---

## GSD 功能比對

**User's input (freeform):**
1. `propose` 指令要有 GSD 式互動提問 + 可選研究
2. 完成後更新 `.mysd/roadmap/` 下的文件（含 change 名稱和完成日期時間）
3. 新增 `/mysd:capture` 指令 — 從對話中分析要做的變更

**Notes:** 使用者對 debug session、todo/notes、session management 暫不加入 v1

---

## Claude's Discretion

- frontmatter schema 具體欄位設計
- lipgloss 配色方案
- 錯誤訊息措辭
- `.mysd/roadmap/` 檔案格式

## Deferred Ideas

- Debug session 功能 → v1.x
- Todo / Notes 管理 → v1.x
- Session pause/resume → v1.x
