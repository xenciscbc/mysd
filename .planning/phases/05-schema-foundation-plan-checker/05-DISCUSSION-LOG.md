# Phase 5: Schema Foundation & Plan-Checker - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-25
**Phase:** 05-schema-foundation-plan-checker
**Areas discussed:** Plan-checker 架構, 新 agent model 對應表, openspec/config.yaml 內容, Schema 向後相容策略

---

## Plan-checker 架構

| Option | Description | Selected |
|--------|-------------|----------|
| Go binary 實作 | satisfies 欄位對 MUST ID 的 structured matching 放在 Go。確定性高、快速、可測試。agent 只負責顯示結果 + 互動補齊。 | ✓ |
| Agent 層實作 | 全部交給 mysd-plan-checker agent 用 AI 判斷。彈性高但結果不可預測。 | |
| 混合式 | Go 做 structured matching + agent 做語意分析（檢查 satisfies ID 之外的覆蓋度） | |

**User's choice:** Go binary 實作
**Notes:** 無額外備註

### 觸發方式

| Option | Description | Selected |
|--------|-------------|----------|
| 沿用 --check flag | `mysd plan --check` 在 plan 完成後自動執行檢查。與現有 CLI 介面一致。 | ✓ |
| 新增獨立 subcommand | `mysd check` 作為獨立指令，可在任何時候執行。 | |

**User's choice:** 沿用 --check flag

### 失敗行為

| Option | Description | Selected |
|--------|-------------|----------|
| 顯示缺口 + 互動補齊 | 列出未覆蓋的 MUST IDs，讓使用者選擇：自動補齊新 tasks、手動調整、或忽略繼續。 | ✓ |
| 單純報告 | 只列缺口清單，不做互動補齊。 | |
| You decide | Claude 自行決定互動方式。 | |

**User's choice:** 顯示缺口 + 互動補齊

### 輸出格式

| Option | Description | Selected |
|--------|-------------|----------|
| 簡潔模式 | 輸出 uncovered MUST IDs 清單 + 覆蓋率比例。agent 負責渲染互動 UI。 | ✓ |
| 完整模式 | 輸出每個 MUST ID 的 covered/uncovered 狀態、對應 task IDs、建議補齊方式。 | |
| You decide | Claude 自行決定。 | |

**User's choice:** 簡潔模式

### PlanningContext JSON 擴展

| Option | Description | Selected |
|--------|-------------|----------|
| 從 tasks.md depends 欄位計算 | WaveGroups 由 Go binary 做 topological sort 得到，WorktreeDir 用預設 .worktrees/，AutoMode 從 config 讀取。 | ✓ |
| 純靜態屬性 | Phase 5 只先加佔位（空值或預設值），實際計算等 Phase 6。 | |

**User's choice:** 從 tasks.md depends 欄位計算

---

## 新 agent 的 model 對應表

| Option | Description | Selected |
|--------|-------------|----------|
| 全部 sonnet（沿用現有） | 與現有 6 roles 一致，quality/balanced 全 sonnet，budget 層 advisor/proposal-writer 也用 sonnet。 | ✓ |
| 依重要性分化 | plan-checker 和 researcher 用 sonnet，advisor 和 proposal-writer 在 budget 用 haiku。 | |
| You decide | Claude 自行決定。 | |

**User's choice:** 全部 sonnet（沿用現有）

---

## openspec/config.yaml 內容

### 欄位設計

| Option | Description | Selected |
|--------|-------------|----------|
| 最小必要 | project name, locale (BCP47), spec_dir, created — 其餘由 convention-over-config 推導。 | ✓ |
| 完整 metadata | project name, locale, spec_dir, language, framework, description, version, created。 | |
| You decide | Claude 自行決定。 | |

**User's choice:** 最小必要

### Locale 格式

| Option | Description | Selected |
|--------|-------------|----------|
| BCP47 | zh-TW, en-US, ja-JP — 國際標準，與 Go golang.org/x/text/language 相容。 | ✓ |
| 簡化格式 | zh, en, ja — 更簡潔但失去區域區分。 | |

**User's choice:** BCP47
**Notes:** 確認 OpenSpec 未定義 config.yaml 標準格式，由 mysd 自行設計。

### Locale 同步

| Option | Description | Selected |
|--------|-------------|----------|
| config.yaml 為主 | openspec/config.yaml 的 locale 是 source of truth，mysd.yaml 參考。/mysd:lang 修改時兩者同步。 | ✓ |
| 各自獨立 | 各自維護 locale 欄位，/mysd:lang 負責同時更新兩者。 | |
| You decide | Claude 自行決定。 | |

**User's choice:** config.yaml 為主

---

## Schema 向後相容策略

### 新欄位策略

| Option | Description | Selected |
|--------|-------------|----------|
| omitempty 零值過渡 | 新欄位加 `yaml:"depends,omitempty"` 等 tag。舊 tasks.md 讀取時新欄位為 nil/empty，寫回時不輸出空欄位。 | ✓ |
| 明確 migration | 新版 binary 首次讀取舊 tasks.md 時自動加入新欄位預設值並寫回。 | |
| You decide | Claude 自行決定。 | |

**User's choice:** omitempty 零值過渡

### TaskItem JSON 擴展

| Option | Description | Selected |
|--------|-------------|----------|
| 同步擴展 | TaskItem 加 Depends/Files/Satisfies/Skills，Phase 6 wave grouping 直接可用。 | ✓ |
| 延後到 Phase 6 | Phase 5 只擴展 schema 層，TaskItem JSON 等 Phase 6 再加。 | |
| You decide | Claude 自行決定。 | |

**User's choice:** 同步擴展

---

## Claude's Discretion

- plan-checker 的具體 JSON 輸出結構
- topological sort 演算法選擇
- config.yaml writer 的錯誤處理細節
- 新欄位的測試案例覆蓋範圍

## Deferred Ideas

None — discussion stayed within phase scope
