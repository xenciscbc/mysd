## 1. 路徑分流入口（discuss-path-routing）

- [x] 1.1 重構 SKILL.md Step 2：移除現有 6 層 fallback source detection，改為先判斷是否有 active change（從 argument 或 `mysd status`），然後呈現路徑選擇提示（spec-focused / source-driven），實作 "discuss skill presents path selection at entry" requirement
- [x] 1.2 實作 spec-focused 路徑：列出 `.specs/changes/{change_name}/specs/*/spec.md` 下所有 spec，顯示 capability name 和 requirement 數量，等待用戶選擇，實作 "spec-focused path lists specs and accepts user selection" requirement
- [x] 1.3 實作無 active change 時自動進入 source-driven 路徑的邏輯，以及 auto_mode 預設走 source-driven 路徑

## 2. 來源驅動路徑 + Material Selection 擴展

- [x] 2.1 將 propose SKILL.md 的 6 來源偵測邏輯（source_arg file/directory、conversation context、Claude plan、gstack plan、active change、deferred notes）複製到 discuss SKILL.md 的 source-driven 路徑中，確保 "Propose skill detects all available requirement sources" requirement 在 discuss 中同樣適用，實作 "discuss source-driven path uses same detection" scenario
- [x] 2.2 實作 material selection 互動流程（numbered list、multi-select、aggregation），確保 "user selects requirement sources interactively" requirement 在 discuss 中運作
- [x] 2.3 實作 auto mode 在 discuss 中自動聚合所有偵測到的來源，實作 "auto mode uses all detected sources without interaction" requirement
- [x] 2.4 實作聚合內容與既有 specs 的比對邏輯：掃描 change 下的 specs 和 `openspec/specs/`，比較關鍵概念，推薦新建或併入，實作 "source-driven path uses material selection and recommends spec target" requirement
- [x] 2.5 實作無 active change 時從聚合內容推導 change name 並執行 `mysd propose {name}` scaffold 的流程

## 3. Spec 缺口分析（spec-gap-analysis）

- [x] 3.1 實作 requirement 覆蓋度分析：比對選中 spec 的 requirements 與 proposal.md Capabilities，識別未覆蓋的 capability 和孤立的 requirement，實作 "gap analysis evaluates requirement coverage against proposal" requirement
- [x] 3.2 實作 scenario 完整度分析：檢查每條 requirement 是否有至少一個 `#### Scenario:` block，實作 "gap analysis evaluates scenario completeness" requirement
- [x] 3.3 實作邊界條件覆蓋度分析：檢查 scenario 是否包含 error/failure 和 edge case 關鍵字，實作 "gap analysis evaluates boundary condition coverage" requirement
- [x] 3.4 實作缺口結果呈現：按三個維度分組列出缺口，詢問用戶要先處理哪個缺口，實作 "gap analysis results drive the discussion starting point" requirement

## 4. Discussion Loop 與 Research 出口合併

- [x] 4.1 合併 Step 8 Layer 2 出口與 Step 9 入口：research 完成後呈現所有 gray area 結論的統一摘要，提供「繼續討論 / 收斂」兩個選項，移除原本重複的提示，實作 "discussion loop merges research conclusion and exit into unified flow" requirement
- [x] 4.2 調整無 research 時的 Step 9 入口：根據路徑類型載入對應 context（spec-focused 用缺口分析結果、source-driven 用聚合內容），確保不會出現空白開場

## 5. Spec Update 確認機制（planning delta）

- [x] 5.1 實作 Step 10 用戶確認清單：收斂結論後列出有影響的 artifact（spec files、proposal.md、design.md），預設全勾選，允許用戶取消個別項目，auto_mode 跳過確認，實作 "discuss skill spec update confirmation" requirement
- [x] 5.2 修改 Step 10 的 writer agent 呼叫邏輯：只對用戶確認的項目 spawn 對應的 writer agent

## 6. Re-plan 條件觸發（planning delta）

- [x] 6.1 修改 Step 11：在執行 re-plan 前檢查 `.specs/changes/{change_name}/tasks.md` 是否存在，存在才執行 re-plan + plan-checker，不存在則跳過到 Step 12，實作 "discuss skill re-plan is conditional on existing plan" requirement

## 7. 驗證

- [x] 7.1 驗證路徑分流：模擬有 active change 時兩條路徑的入口行為正確
- [x] 7.2 驗證 source-driven 路徑的 material selection 流程與 propose 一致
- [x] 7.3 驗證 spec-focused 路徑的三面向缺口分析輸出正確
- [x] 7.4 驗證 Step 9 unified exit 在有/無 research 兩種情況下均正常運作
- [x] 7.5 驗證 Step 10 確認清單和 Step 11 條件觸發的行為正確
