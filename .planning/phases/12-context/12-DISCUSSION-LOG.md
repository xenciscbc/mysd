# Phase 12: 加入 context 的 % 數及色條 - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions captured in CONTEXT.md — this log preserves the analysis.

**Date:** 2026-03-27
**Phase:** 12-context
**Mode:** discuss
**Areas discussed:** 目標位置、內容組成、視覺格式、安裝機制、GSD 相容性

## Gray Areas Discussed

### 目標位置
| Question | User Answer |
|----------|-------------|
| 「加入 context % 色條」想加在哪裡？ | 改善現有 statusline（澄清後：把 GSD 功能移植到 mysd 自己的 plugin 中） |

**Key insight:** GSD statusline 已有完整 context bar 實作。Phase 12 的目標是讓 mysd 有自己獨立的 statusline hook，不依賴 GSD 安裝。

### 內容組成
| Question | User Answer |
|----------|-------------|
| 除了 context % 色條，還要顯示哪些資訊？ | model 名稱（multiSelect 只選了這個） |
| GSD 同時存在時怎處理？ | mysd statusline 從 GSD 讀取 ctx 資料（即：同樣的資料來源和邏輯） |

**Interpretation:** "從 GSD 讀取 ctx 資料" 意指使用相同的計算邏輯，並寫入 bridge file 保持 GSD context-monitor 相容性，而非字面上讀取 GSD 檔案。

### 視覺格式
| Question | User Answer |
|----------|-------------|
| mysd statusline 要用什麼格式？ | 加上當前 change → `claude-sonnet-4-5 │ my-change │ mysd │ █████░░░░░ 50%` |

### 安裝機制
| Question | User Answer |
|----------|-------------|
| 如何進入使用者的 Claude Code？ | mysd init 寫入 settings.json（Recommended 選項） |

## Decisions Made

- D-01 to D-02: 格式決策（含/不含 active change 的兩種格式）
- D-03 to D-05: Context 計算邏輯（direct read + normalization + 顏色閾值）
- D-06 to D-08: 安裝機制（mysd init 複製 hook + merge settings.json）
- D-09: Plugin manifest 不擴展（hooks 是 one-time init）
- D-10: .specs/state.yaml 讀取方式（regex line scan + silent fail）

## Deferred

- Task 名稱顯示（類似 GSD todos 讀取）— 未來 phase
- mysd update 自動更新 hook — 未來 milestone
