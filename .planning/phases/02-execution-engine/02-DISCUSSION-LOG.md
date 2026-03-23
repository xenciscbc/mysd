# Phase 2: Execution Engine - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-23
**Phase:** 02-execution-engine
**Areas discussed:** AI 調用機制, Alignment Gate 設計, Workflow 指令深度, 任務進度與恢復

---

## AI 調用機制

### Q1: mysd execute 如何讓 AI 實際執行任務？

| Option | Description | Selected |
|--------|-------------|----------|
| Prompt 檔案產生器 | Go binary 產生結構化 prompt 檔案，透過 Claude Code --print 或 pipe 送入 | |
| Claude Code CLI exec | Go binary 直接 os.Exec 呼叫 claude CLI | |
| Plugin skill 反向調用 | 由 /mysd:execute 觸發，SKILL.md 定義流程，Go binary 負責 spec 解析和狀態管理 | ✓ |
| 混合模式 | Plugin skill 入口 + Go binary 結構化指令 + Claude Code agent 執行 | |

**User's choice:** Plugin skill 反向調用
**Notes:** Claude Code agent 是執行主體（有完整 tool access），Go binary 是 spec 引擎

### Q2: 執行流程的控制權在誰手上？

| Option | Description | Selected |
|--------|-------------|----------|
| SKILL.md 編排 | SKILL.md 定義完整流程，透過 bash 呼叫 mysd binary | |
| Go binary 編排 | Go binary 透過 JSON 輸出結構化指令 | |
| Agent definition 編排 | 專用 agent .md 定義執行行為，SKILL.md 只是入口 | ✓ |

**User's choice:** Agent definition 編排
**Notes:** 使用者問了 GSD 的做法，確認後選擇仿 GSD 模式

### Q3: Agent 的粒度如何劃分？

| Option | Description | Selected |
|--------|-------------|----------|
| 單一執行 agent | 一個 mysd-executor.md 處理所有執行邏輯 | |
| 每階段專屬 agent | 仿 GSD：spec-writer, designer, planner, executor, verifier | ✓ |
| 混合粒度 | 執行有專屬 agent，其他共用輕量 agent | |

**User's choice:** 每階段專屬 agent
**Notes:** 使用者先詢問 GSD 的做法（每階段專屬 agent），了解後選擇相同模式

### Q4: Wave mode 多 agent 平行機制？

| Option | Description | Selected |
|--------|-------------|----------|
| Claude Code 原生 Task | 用 Agent tool 生成平行 subagent | ✓ |
| 多次順序呼叫 | 自動推進但不平行 | |
| Go binary goroutine | goroutine 同時啟動多個 Claude Code CLI | |

**User's choice:** Claude Code 原生 Task
**Notes:** 無

### Q5: 執行時的模型選擇策略？

| Option | Description | Selected |
|--------|-------------|----------|
| 使用者當前模型 | 直接用 Claude Code 當前 session 模型 | |
| 可配置模型對應 | mysd.yaml 配置每階段模型 | |
| 智能分派 | 根據任務複雜度自動分派 | |
| 仿 GSD profile-based | Profile (quality/balanced/budget) + per-agent 覆蓋 | ✓ |

**User's choice:** 仿 GSD 的選擇方式
**Notes:** Profile-based 模型管理，可在 mysd.yaml 配置預設 profile

---

## Alignment Gate 設計

### Q1: Alignment gate 的強制機制？

| Option | Description | Selected |
|--------|-------------|----------|
| Prompt 注入 + 確認指令 | Agent definition 強制包含 spec，要求 AI 列出理解的 MUST/SHOULD/MAY | ✓ |
| 兩階段驗證 | AI 產生理解摘要 → Go binary 檢查是否涵蓋所有 MUST | |
| Spec 檔案強制讀取 | @file 引用 + 每個 task 必須引用 spec ID | |

**User's choice:** Prompt 注入 + 確認指令
**Notes:** 無

### Q2: 確認輸出格式？

| Option | Description | Selected |
|--------|-------------|----------|
| 結構化摘要 | 列出所有 MUST 項目、執行策略、疑問 | ✓ |
| 簡單確認 | 只輸出確認訊息 + MUST 數量 | |
| 任務對映表 | 每個 task 對映到 requirement IDs | |

**User's choice:** 結構化摘要
**Notes:** 無

### Q3: 摘要存在哪裡？

| Option | Description | Selected |
|--------|-------------|----------|
| 寫入 .specs/ 目錄 | .specs/changes/{name}/alignment.md | ✓ |
| 只在 session 中 | 不寫檔 | |
| 寫入 STATE.json | WorkflowState 一部分 | |

**User's choice:** 寫入 .specs/ 目錄
**Notes:** 無

---

## Workflow 指令深度

### Q1: spec / design / plan 的互動深度？

| Option | Description | Selected |
|--------|-------------|----------|
| 完整 AI 互動流程 | 每個指令由專屬 agent 執行完整互動 | ✓ |
| 狀態更新 + 檔案編輯 | 只更新 state 和開啟檔案 | |
| 混合深度 | spec/design 完整互動，plan 只做拆解 | |

**User's choice:** 完整 AI 互動流程
**Notes:** 無

### Q2: mysd ff 跳過哪些步驟？

| Option | Description | Selected |
|--------|-------------|----------|
| propose → plan 一氣完成 | 跳過互動確認，用預設值 | ✓ |
| 全自動到 execute | 不只停在 plan | |
| 可配置終點 | --to=execute 指定推進到哪 | |

**User's choice:** propose → plan 一氣完成
**Notes:** 無

### Q3: mysd capture 的運作方式？

| Option | Description | Selected |
|--------|-------------|----------|
| 分析對話 + 進入 propose | 提取需求後帶預填進 propose 互動 | ✓ |
| 直接產生 spec | 跳過 propose 互動 | |
| Phase 2 延後 | 太複雜，延到後面 | |

**User's choice:** 分析對話 + 進入 propose
**Notes:** 無

### Q4: mysd status 顯示什麼？

| Option | Description | Selected |
|--------|-------------|----------|
| 綜合儀表板 | change name, phase, 完成率, MUST 達成, 時間 | ✓ |
| 簡潔狀態 | 只有 phase 和下一步建議 | |

**User's choice:** 綜合儀表板
**Notes:** 無

### Q5: Plan 階段是否含 research + check？

| Option | Description | Selected |
|--------|-------------|----------|
| 完整管線 | research → plan → check | |
| 簡化版 | 只有 plan | |
| 可選管線 | 預設只 plan，可啟用 research/check | ✓ |

**User's choice:** 可選管線（仿 GSD 可以選擇）
**Notes:** 使用者明確要求「仿 gsd 可以選擇，可以用完整的也可以只 plan 就好」

---

## 任務進度與恢復

### Q1: Task 狀態更新機制？

| Option | Description | Selected |
|--------|-------------|----------|
| Agent 回報 + binary 更新 | AI 呼叫 mysd task-update，binary 更新 tasks.md 和 STATE.json | ✓ |
| Agent 直接更新檔案 | Agent 編輯 tasks.md | |
| Binary 主動偵測 | Go binary 檢查檔案系統變更 | |

**User's choice:** Agent 回報 + binary 更新
**Notes:** 無

### Q2: 中斷恢復的粒度？

| Option | Description | Selected |
|--------|-------------|----------|
| Task level | 從最後完成的 task 之後恢復 | ✓ |
| Sub-task level | 更細粒度追蹤 | |
| Phase level | 從頭開始 | |

**User's choice:** Task level
**Notes:** 無

### Q3: TDD 模式改變什麼？

| Option | Description | Selected |
|--------|-------------|----------|
| Test-first 指令注入 | Agent definition 增加 RED → GREEN → REFACTOR 指令 | ✓ |
| 分離的 test agent | 專屬 test-writer agent 先執行 | |
| Post-implementation test | 執行後自動產生測試 | |

**User's choice:** Test-first 指令注入
**Notes:** 無

### Q4: Atomic commits 粒度？

| Option | Description | Selected |
|--------|-------------|----------|
| 每個 task 一個 commit | 完成一個 task 就 commit | ✓ |
| 每個 TDD cycle | RED/GREEN/REFACTOR 各 commit | |
| 整個 execute 一個 commit | 所有 tasks 完成後 commit | |

**User's choice:** 每個 task 一個 commit
**Notes:** 無

---

## Claude's Discretion

- Agent definition 具體 prompt 措辭
- Alignment summary markdown 模板
- mysd task-update CLI 介面設計
- Status 儀表板配色和排版
- ff 預設值選擇策略
- Model profile 具體映射表

## Deferred Ideas

None — discussion stayed within phase scope
