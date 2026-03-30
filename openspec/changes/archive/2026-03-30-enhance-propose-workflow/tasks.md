## 1. 重構 Step 順序與 Argument Parsing

- [x] 1.1 重寫 Step 1 (Parse Arguments)：僅解析 `--auto` flag 和 `source_arg`，移除原 Source Detection 的 Priority 1-6 邏輯
- [x] 1.2 將 Resolve Agent Model 移至 Step 2（原 Step 3b），Load Deferred Notes 移至 Step 3（原 Step 4）

## 2. Material Selection（素材選擇）

- [x] 2.1 實作 Step 4 來源偵測邏輯：偵測 source_arg 檔案/目錄、對話 context、Claude plan、gstack plan、active change、deferred notes 共 6 種來源（D-01: 素材選擇整合原 Source Detection）
- [x] 2.2 實作來源列表呈現：僅列有內容的來源，手動輸入永遠作為最後選項，支援多選（propose skill detects all available requirement sources）
- [x] 2.3 實作 auto mode 行為：自動用所有偵測到的來源，無來源時從 conversation context 萃取（auto mode uses all detected sources without interaction）
- [x] 2.4 實作 aggregated_content 彙整：讀取並合併所有選中來源的內容（user selects requirement sources interactively）

## 3. Scan Existing Specs

- [x] 3.1 將 Scan Existing Specs 移至 Step 5，保留相關 spec 內容供訪談步驟使用（existing spec content fed into interview）

## 4. Requirement Interview（需求訪談）

- [x] 4.1 實作 Step 6 完整度評估邏輯：對照 Problem / Boundary / Success Criteria 三個面向檢查 aggregated_content（orchestrator evaluates requirement completeness）
- [x] 4.2 實作既有 spec 重疊偵測：比對 aggregated_content 與相關既有 spec，重疊時詢問擴充或新建（existing spec content fed into interview）
- [x] 4.3 實作動態訪談迴圈：D-02: orchestrator 自主決定訪談次數，每次只問一個問題，問完重新評估完整度，不足處繼續問（interview asks one question at a time、interview question count is dynamic）
- [x] 4.4 實作 requirement_brief 產出：結構化格式包含 Problem / Boundary / Success Criteria / Source（D-03: requirement_brief 結構化格式、interview produces structured requirement_brief）
- [x] 4.5 實作 auto mode 訪談跳過：best-effort 推測填入各維度，不留空（auto mode skips interview with best-effort brief）

## 5. Change Name 延後推導

- [x] 5.1 實作 Step 7：從 requirement_brief 推導 change name 和分類 change type，既有 change 直接使用原名（D-04: Change name 延後推導、change name derived after requirement interview）

## 6. Gray Area 條件化

- [x] 6.1 修改 Step 10-11（Gray Area + Advisor、Dual-Loop Exploration）：加入條件判斷，僅在使用者接受 research 後才執行（D-05: Gray area 步驟條件化、propose workflow step ordering）

## 7. 整合驗證

- [x] 7.1 驗證完整的 15 步流程順序符合 spec 定義（propose workflow step ordering）
- [x] 7.2 驗證 auto mode 全流程：素材自動選擇 → 跳過訪談 → 跳過 research → proposal writer
- [x] 7.3 驗證 research 路徑：素材選擇 → 訪談 → research → gray area → exploration → proposal writer
