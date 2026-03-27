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

## Gray Areas Discussed (Session 2 — Update)

### D-01/D-02 Model 簡稱
| Question | User Answer |
|----------|-------------|
| model 顯示格式？ | 只要簡稱（sonnet / opus / haiku），不要完整 display_name |

### D-03/D-04/D-05 Bridge File 邏輯
| Question | User Answer |
|----------|-------------|
| statusline off 時要不要寫 bridge file？ | 選 B：只關顯示，繼續寫 bridge file（保護 GSD context-monitor） |
| 可以偵測 GSD 並存嗎？ | 是，偵測 gsd-context-monitor.js 存在與否 |
| 不並存時 on 狀態要不要寫？ | 不用，完全不寫 — bridge file 寫入只在並存時才做 |

### D-11/D-12/D-13 On/Off 控制
| Question | User Answer |
|----------|-------------|
| 切換方式？ | `/mysd:statusline` 指令：無參數 toggle，帶 on/off 直接設定 |

### Phase 範圍擴展
| Question | User Answer |
|----------|-------------|
| discuss research cache 放哪？ | 共同放入 Phase 12（改名為 context 管理） |
| cache 存放位置？ | `.specs/changes/{change}/discuss-research-cache.json` |
| cache 內容？ | 4 維度完整輸出 |

## Decisions Made

- **D-01/D-02:** model 用關鍵字比對提取簡稱（sonnet/opus/haiku）
- **D-03/D-04/D-05:** context 計算 + bridge file 只在 GSD 並存時寫入（與 on/off 無關）；🥵 取代 💀
- **D-06 to D-10:** 安裝、sync、state.yaml 讀取（不變）
- **D-11:** `statusline_enabled: true` 加入 ProjectConfig
- **D-12:** hook 讀 config 決定是否輸出（off 時仍檢查 GSD 並存寫 bridge file）
- **D-13:** `/mysd:statusline [on|off]` SKILL.md 指令
- **D-14:** cache 路徑 `.specs/changes/{change}/discuss-research-cache.json`
- **D-15:** cache 格式：JSON，含 change_name、cached_at、4 維度完整輸出
- **D-16:** discuss research 完成後立即寫入 cache（proactive）
- **D-17:** 啟動 discuss 時偵測 cache，詢問重用/重新 research/都不要
- **D-18:** archive 時自動刪除 cache 檔
- **D-19:** cache 檔加入 .gitignore

## Deferred

- Task 名稱顯示（類似 GSD todos 讀取）— 未來 phase
- mysd update 自動更新 hook — 未來 milestone
