## Context

現有的 `mysd/skills/propose/SKILL.md` 包含 13 個步驟的 propose 流程。Step 2 Source Detection 有 6 個 Priority，混合了來源偵測與 change name 推導。使用者的需求只在 Priority 6 時透過一個簡單問題收集，缺乏深度挖掘。4 維度 research 聚焦於 codebase/domain/architecture/pitfalls，沒有需求釐清維度。

此變更僅修改 `mysd/skills/propose/SKILL.md`（SKILL markdown），不涉及 Go 程式碼或 agent 定義變更。

## Goals / Non-Goals

**Goals:**

- 讓 propose 流程能主動挖掘需求，而非被動接收
- 整合分散的來源偵測邏輯到單一互動點
- 讓使用者明確選擇素材來源
- 在訪談階段偵測與既有 spec 的重疊
- 將 change name 推導延後到需求充分理解後

**Non-Goals:**

- 不新增 agent — 訪談由 orchestrator 執行
- 不改變 4 維度 research 的維度（不加第 5 維度）
- 不修改 mysd-researcher、mysd-advisor、mysd-proposal-writer 等 agent 定義
- 不修改 Go CLI 程式碼（mysd binary）

## Decisions

### D-01: 素材選擇整合原 Source Detection

將現有 Step 2 的 Priority 1-6 來源偵測邏輯拆分：
- Step 1 僅做 argument parsing（`--auto`、`source_arg`）
- 新 Step 4 集中處理所有來源偵測與使用者選擇

**替代方案**：保留原 Priority 結構，在旁邊加素材選擇 → 拒絕，因為兩套偵測邏輯會重疊且難以維護。

偵測的 6 種來源：

| 來源 | 偵測方式 |
|------|---------|
| source_arg 指定的檔案/目錄 | 檢查 source_arg 是否為有效路徑 |
| 對話 context | 檢查 conversation 中是否有實質需求討論 |
| Claude plan | 從 conversation system messages 找 plan file path |
| gstack plan | 掃描 `~/.gstack/projects/{project}/` 下 `.md` 檔 |
| Active change | `mysd status` 輸出 |
| Deferred notes | `mysd note list` 輸出 |

呈現規則：
- 只列出有內容的來源
- 手動輸入永遠作為最後一個選項
- 沒偵測到來源時直接進手動輸入
- auto mode：自動用所有偵測到的來源

### D-02: Orchestrator 自主決定訪談次數

訪談不硬編碼輪次上限。Orchestrator 讀完彙整素材後，對照三個面向評估：

| 面向 | 檢查項 |
|------|--------|
| Problem | 是否清楚描述了要解決的問題（而非只說想要的解法） |
| Boundary | 是否明確什麼做、什麼不做 |
| Success Criteria | 是否有具體可驗證的成功條件 |

規則：
- 已有充足資訊的面向不問
- 每次只問一個問題
- 若與既有 spec 有重疊，主動問「擴充還是新建」
- auto mode：跳過訪談，用現有 context best-effort 推測

**替代方案**：固定 3 輪訪談 → 拒絕，因為資訊充足時強制問 3 輪會浪費使用者時間。

### D-03: requirement_brief 結構化格式

訪談完成後產出結構化摘要，作為後續 research 和 proposal-writer 的輸入：

```
## Problem
{要解決的問題}

## Boundary
{做什麼 / 不做什麼}

## Success Criteria
{怎麼算完成}

## Source
{素材來源標記，方便追溯}
```

此格式不寫入檔案，僅作為流程中的中間資料傳遞給後續步驟。

### D-04: Change name 延後推導

Change name 從 requirement_brief 推導（Step 7），而非從初始 source content 推導。此時 orchestrator 已完成訪談，對需求有更完整的理解，推導的名稱品質更高。

若 source_arg 已指向既有 change（Priority 1），則直接使用該名稱，不重新推導。

### D-05: Gray area 步驟條件化

Step 10-11（Gray Area + Advisor、Dual-Loop Exploration）僅在使用者於 Step 9 接受 4-Dimension Research 後才執行。未跑 research 時無 research 輸出可分析，直接跳到 Step 12 Proposal Writer。

## Risks / Trade-offs

- **[流程步驟增加]** → 素材選擇和訪談增加了互動步驟。Mitigation：auto mode 完全跳過這些步驟；資訊充足時訪談為 0 問，只增加一個素材選擇的互動點。
- **[requirement_brief 非持久化]** → brief 不寫檔，conversation context 壓縮後可能遺失。Mitigation：brief 的核心內容會被 proposal-writer 寫入 proposal.md，不會真正遺失。
- **[既有 spec 重疊判斷可能不準]** → orchestrator 用名稱比對判斷重疊，可能漏判或誤判。Mitigation：這是 informational 輔助，最終決策權在使用者手上。
