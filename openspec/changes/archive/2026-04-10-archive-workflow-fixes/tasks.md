## 1. Task Status Updates Auto-Transition

- [x] 1.1 在 `cmd/task_update.go` 的 `runTaskUpdate` 中，實作 Task Status Updates 的自動推進邏輯：更新 task status 後檢查是否所有 tasks 都為 terminal state（done/skipped）。若是且 phase 為 `planned`，呼叫 `state.Transition(&ws, state.PhaseExecuted)` 並印出 "All tasks complete — phase advanced to executed"
- [x] 1.2 在 `cmd/task_update_test.go` 加入 Task Status Updates auto-transition 測試：(a) 最後一個 task 完成時自動推進到 executed (b) 還有 pending tasks 時不推進 (c) phase 已經是 executed 時不重複推進

## 2. Goal-Backward Verification Safety Net

- [x] 2.1 在 `cmd/verify.go` 的 `--write-results` 處理邏輯中，實作 Goal-Backward Verification 的 safety net：於解析 report 前檢查 phase。若 phase 為 `planned` 且所有 tasks 為 terminal，先呼叫 `state.Transition(&ws, state.PhaseExecuted)` 再繼續
- [x] 2.2 在 `cmd/verify_test.go` 加入 Goal-Backward Verification safety net 測試：(a) phase 為 planned + 全部 tasks done 時自動推進到 executed 再處理結果 (b) phase 已為 executed 時不做額外推進

## 3. Verification Context ID Format

- [x] 3.1 在 `mysd/agents/mysd-verifier.md` Phase 6 report format 的 Rules 區塊加入 Verification Context stable ID 規則：「`id` 欄位必須使用 input context 中 must_items/should_items/may_items 提供的 id 值原樣複製，不可自行編號（如 MUST-01）。例如 context 給 `spec.md::must-5451802d`，report 就寫 `spec.md::must-5451802d`」

## 4. Delta spec merge on archive Fallback

- [x] 4.1 實作 Delta spec merge on archive fallback：在 `internal/spec/merge.go` 的 `MergeSpecs` 中，當 `ParseDelta` 回傳的 added/modified/removed/renamed 全部為空時，解析 deltaBody 的 frontmatter 取得 `delta` 欄位。根據 delta 值決定行為：`ADDED` → 使用 delta body 作為新 spec 內容；`MODIFIED` → 用 delta body 全量替換主 spec body 並遞增版本號；其他 → 回傳 warning
- [x] 4.2 修改 `mergeDeltaSpecs`（`cmd/archive.go`），將 delta spec 的完整檔案內容（含 frontmatter）傳給 `MergeSpecs`。確保 `MergeSpecs` 能正確處理含 frontmatter 的輸入（目前已在做 frontmatter parse）
- [x] 4.3 在 `internal/spec/merge_test.go` 加入測試：(a) delta body 無 delta heading + frontmatter delta=ADDED → 整個 body 作為新 spec (b) delta body 無 delta heading + frontmatter delta=MODIFIED → 全量替換 body (c) delta body 無 delta heading + delta 為空 → 回傳 warning 並跳過

## 5. Spec Writer Agent Prompt 修正

- [x] 5.1 修改 `mysd/agents/mysd-spec-writer.md` Step 2 的範例格式，將 `## Requirements` / `### MUST` 格式改為使用 delta heading 格式（`## ADDED Requirements` / `### Requirement: <name>`），與 `spectra instructions specs` 的 template 保持一致
