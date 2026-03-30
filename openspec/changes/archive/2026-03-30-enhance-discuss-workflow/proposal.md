## Why

目前 mysd discuss 的入口只能以整個 change 為單位進行討論，粒度太粗。無法聚焦到單一 spec 做深度缺口分析，也無法從外部來源內容驅動新增或擴充 spec。此外 Step 8 與 Step 9 的出口存在重複詢問，Step 10 缺少用戶確認機制，Step 11 在沒有 plan 時仍會嘗試 re-plan。

## What Changes

- 在 Step 2 新增路徑分流：進入 discuss 後先判斷 context，詢問用戶要走「討論既有 spec」還是「從來源增加內容」
- 路徑 1（Spec-focused）：列出 change 下的 specs 供用戶選擇，自動執行三面向缺口分析（requirement 覆蓋度、scenario 完整度、邊界條件），以缺口列表驅動討論
- 路徑 2（來源驅動）：復用 propose 的 6 來源偵測 + material selection 機制，聚合內容後自動比對既有 specs 推薦新建或併入；無 active change 時自動 scaffold
- 合併 Step 8 Layer 2 出口與 Step 9 入口，消除重複詢問：research 完成後統一呈現結論摘要，提供「繼續討論 / 收斂」選項
- Step 10 新增用戶確認清單：僅列出有影響的 artifact，預設全勾選，用戶可取消個別項目
- Step 11 改為條件觸發：檢查是否已有 plan（tasks.md 存在），有才執行 re-plan

## Capabilities

### New Capabilities

- `discuss-path-routing`: discuss 技能的路徑分流機制，根據 context 和用戶選擇決定走 spec-focused 或來源驅動路徑
- `spec-gap-analysis`: 針對單一 spec 的三面向缺口分析能力（requirement 覆蓋度、scenario 完整度、邊界條件檢查）

### Modified Capabilities

- `material-selection`: 將 material selection 的 6 來源偵測機制從 propose 專用擴展為 discuss 也能使用
- `planning`: Step 11 re-plan 改為條件觸發，僅在 tasks.md 存在時執行

## Impact

- Affected code: `mysd/skills/discuss/SKILL.md`（主要修改目標）
- Affected specs: `material-selection`（擴展使用範圍）、`planning`（條件觸發）
- 新增 specs: `discuss-path-routing`、`spec-gap-analysis`
