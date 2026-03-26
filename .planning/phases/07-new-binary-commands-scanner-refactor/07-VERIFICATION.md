---
phase: 07-new-binary-commands-scanner-refactor
verified: 2026-03-26T11:00:00Z
status: passed
score: 7/7 must-haves verified
re_verification: false
human_verification:
  - test: "mysd lang set zh-TW 的 atomic rollback 在非 Windows 系統上實際執行"
    expected: "若 openspec/config.yaml 寫入失敗，.claude/mysd.yaml 回滾到原始值"
    why_human: "TestLangSet_AtomicRollback 在 Windows 上以 SKIP 執行，無法在此環境驗證。邏輯程式碼已確認存在（rollback pattern），但需要在 Linux/macOS 環境手動確認"
---

# Phase 7: New Binary Commands & Scanner Refactor Verification Report

**Phase Goal:** 使用者可透過 `/mysd:model` 和 `/mysd:lang` 管理設定，scan 支援任意語言，plan 完成後可確認 skills 對應
**Verified:** 2026-03-26T11:00:00Z
**Status:** passed
**Re-verification:** No — initial verification

## Goal Achievement

### Observable Truths

| #   | Truth | Status | Evidence |
|-----|-------|--------|----------|
| 1   | `mysd model` 顯示目前 profile 名稱及 10 個 agent role 與 model 的對應表 | ✓ VERIFIED | `./mysd.exe model` 輸出 "Profile: balanced" + 10 行 Role/Model 表格，spot-check 通過 |
| 2   | `mysd model set quality` 寫入 model_profile 到 .claude/mysd.yaml | ✓ VERIFIED | `runModelSet` 使用 viper ReadInConfig + Set + WriteConfig，TestModelSet_ValidProfile 通過 |
| 3   | `mysd model set invalidname` 回傳含合法 profile 清單的錯誤訊息 | ✓ VERIFIED | `./mysd.exe model set invalidname` 輸出 `unknown profile "invalidname"; valid profiles: quality, balanced, budget` |
| 4   | `mysd lang set zh-TW` 同時更新 .claude/mysd.yaml 的 response_language 及 openspec/config.yaml 的 locale | ✓ VERIFIED | TestLangSet_UpdatesBothConfigs 通過，rollback 邏輯確認存在於 lang.go:80-95 |
| 5   | scan 支援語言無關掃描（Go/Node.js/Python/unknown），輸出新 JSON 格式 | ✓ VERIFIED | `./mysd.exe scan --context-only` 輸出含 primary_language, files, modules 欄位的 JSON；所有 10 個 scanner 測試通過 |
| 6   | mysd init 呼叫 scaffold-only 邏輯建立 openspec/ 結構 | ✓ VERIFIED | init_cmd.go:27 呼叫 `scaffoldOpenSpecDir(".")`，init SKILL.md 明確指示執行 `mysd lang set` 設定 locale |
| 7   | plan 完成後顯示 task-skills 對應表並詢問確認；ffe 模式跳過互動 | ✓ VERIFIED | mysd-planner.md 包含 Step 4.5 (Recommend Skills)、Step 7.5 (Skills Confirmation)、Accept all? Y/n 提示、auto_mode ffe 跳過邏輯 |

**Score:** 7/7 truths verified

### Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `internal/executor/waves.go` | FilterBlockedTasks function with BFS | ✓ VERIFIED | Line 140: `func FilterBlockedTasks(tasks []TaskItem, failedIDs []int) []TaskItem`，完整 BFS 實作 |
| `internal/executor/waves_test.go` | 6 FilterBlockedTasks 測試案例 | ✓ VERIFIED | TestFilterBlockedTasks_EmptyFailedIDs, DirectDependency, TransitivePropagation, MultipleFailures, FailedExcluded, NoDependencies 全數通過 |
| `internal/scanner/scanner.go` | 語言無關 ScanContext + BuildScanContext | ✓ VERIFIED | ScanContext 含 PrimaryLanguage, Files, Modules；無 PackageInfo 殘留；detectPrimaryLanguage 支援 go/nodejs/python/unknown |
| `internal/scanner/scanner_test.go` | 10 個語言偵測測試 | ✓ VERIFIED | TestBuildScanContext_GoProject 等 10 個測試全數通過 |
| `cmd/scan.go` | --context-only 及 --scaffold-only flags | ✓ VERIFIED | 兩個 flag 均已登記；scanner.BuildScanContext 呼叫在 line 43；scaffoldOpenSpecDir 在 line 68 |
| `cmd/init_cmd.go` | init 委派給 scaffold logic | ✓ VERIFIED | line 27 呼叫 `scaffoldOpenSpecDir(".")`；不含 PLAN 要求移除的 `yaml.Marshal(cfg)` 用途（僅建立 mysd.yaml 預設值） |
| `cmd/model.go` | model 及 model set 子命令 | ✓ VERIFIED | modelCmd, modelSetCmd, knownRoles (10 roles), runModelRead, runModelSet 均存在 |
| `cmd/model_test.go` | TestModelRead 及 TestModelSet 測試 | ✓ VERIFIED | TestModelRead_DefaultProfile, ContainsAllRoles, NonTTY, TestModelSet_ValidProfile, InvalidProfile, PreservesOtherConfig 通過 |
| `cmd/lang.go` | lang 及 lang set 子命令，含 atomic rollback | ✓ VERIFIED | langCmd, langSetCmd, runLangRead, runLangSet；rollback pattern 在 line 93：`v.Set("response_language", oldResponseLang)` |
| `cmd/lang_test.go` | TestLangSet 測試含 atomic rollback | ✓ VERIFIED | TestLangSet_UpdatesBothConfigs, AtomicRollback (SKIP on Windows), CreatesOpenSpecConfig, PreservesOtherFields |
| `.claude/commands/mysd-scan.md` | 更新為新 JSON 欄位 | ✓ VERIFIED | 含 primary_language, modules, config_exists；不含舊的 go_files/test_files/PackageInfo |
| `.claude/commands/mysd-init.md` | 委派給 scaffold + lang set | ✓ VERIFIED | 含 `mysd lang set {user_choice}` 呼叫；無舊 model_profile/execution_mode 說明 |
| `.claude/commands/mysd-model.md` | /mysd:model 新 SKILL.md | ✓ VERIFIED | 含 frontmatter (model: claude-sonnet-4-5, allowed-tools: Bash/Read)；含 `mysd model` 及 `mysd model set` 步驟 |
| `.claude/commands/mysd-lang.md` | /mysd:lang 新 SKILL.md | ✓ VERIFIED | 含 frontmatter；含 `mysd lang set`、BCP47 參照、zh-TW/en-US 選項 |
| `.claude/agents/mysd-planner.md` | 含 skills 推薦及確認流程 | ✓ VERIFIED | Step 4.5 Recommend Skills、Step 7.5 Skills Confirmation、"Accept all recommended skills? (Y/n)"、auto_mode ffe 跳過 |

### Key Link Verification

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `cmd/scan.go` | `internal/scanner/scanner.go` | `scanner.BuildScanContext` | ✓ WIRED | line 43: `ctx, err := scanner.BuildScanContext(root, exclude)` |
| `cmd/init_cmd.go` | `cmd/scan.go` | `scaffoldOpenSpecDir` | ✓ WIRED | line 27: `scaffoldOpenSpecDir(".")` — 共用同一 package 函式 |
| `cmd/model.go` | `internal/config/config.go` | `config.ResolveModel` | ✓ WIRED | line 56: `model := config.ResolveModel(role, profile, cfg.ModelOverrides)` |
| `cmd/model.go` | `internal/output/printer.go` | `output.NewPrinter` | ✓ WIRED | line 50, 90: `output.NewPrinter(...)` |
| `cmd/lang.go` | `internal/spec/openspec_config.go` | `spec.WriteOpenSpecConfig` | ✓ WIRED | line 91: `spec.WriteOpenSpecConfig(".", osCfg)` |
| `cmd/lang.go` | `internal/config/config.go` | viper response_language | ✓ WIRED | line 80: `v.Set("response_language", locale)`；rollback line 93 |
| `.claude/commands/mysd-scan.md` | `cmd/scan.go` | `mysd scan --context-only` | ✓ WIRED | line 19 of SKILL.md 呼叫 `mysd scan --context-only` |
| `.claude/agents/mysd-planner.md` | `.claude/commands/mysd-plan.md` | skills 欄位整合 | ✓ WIRED | Step 4.5 寫入 tasks.md skills 欄位；Step 7.5 確認流程 |

### Data-Flow Trace (Level 4)

此階段的動態資料元件為 CLI 命令（非 React/Vue 元件），資料流從 binary flag → Go function → stdout。關鍵路徑已驗證：

| Artifact | Data Variable | Source | Produces Real Data | Status |
|----------|---------------|--------|-------------------|--------|
| `cmd/model.go` runModelRead | `cfg.ModelProfile` | `config.Load(".")` 讀取 .claude/mysd.yaml | 是（實際讀取 Viper config） | ✓ FLOWING |
| `cmd/scan.go` runScanContextOnly | `ctx` ScanContext | `scanner.BuildScanContext` 走訪目錄樹 | 是（os.Stat + filepath.WalkDir） | ✓ FLOWING |
| `cmd/lang.go` runLangSet | `osCfg.Locale` | `spec.ReadOpenSpecConfig` + 寫入雙 config | 是（實際讀寫 YAML 檔案） | ✓ FLOWING |

### Behavioral Spot-Checks

| Behavior | Command | Result | Status |
|----------|---------|--------|--------|
| mysd model 顯示 balanced profile + 10 roles | `./mysd.exe model` | "Profile: balanced" + 10 行 Role/Model 表格 | ✓ PASS |
| mysd scan 輸出含 primary_language 的新 JSON | `./mysd.exe scan --context-only` | JSON 含 primary_language: "go", files, modules 欄位 | ✓ PASS |
| mysd model set 無效 profile 回傳清楚錯誤 | `./mysd.exe model set invalidname` | `unknown profile "invalidname"; valid profiles: quality, balanced, budget` (exit 1) | ✓ PASS |
| mysd lang 顯示目前設定 | `./mysd.exe lang` | "Language settings: mysd.yaml response_language: (not set) / openspec/config.yaml locale: tw" | ✓ PASS |
| go test ./... 全數通過 | `go test ./...` | 所有 13 個 packages OK（含 cmd, scanner, executor） | ✓ PASS |
| go build ./... 成功 | `go build ./...` | 無輸出（編譯成功） | ✓ PASS |

### Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
|-------------|-------------|-------------|--------|---------|
| FCMD-03 | 07-03 | `/mysd:model` 顯示/切換 model profile + resolve 特定 agent model | ✓ SATISFIED | cmd/model.go 實作完整；mysd-model.md SKILL.md 存在 |
| FCMD-04 | 07-04 | `/mysd:lang` 互動式設定 response_language，同步 mysd.yaml 和 openspec/config.yaml | ✓ SATISFIED | cmd/lang.go 雙 config 寫入；mysd-lang.md SKILL.md 存在 |
| FCMD-05 | 07-04 | `/mysd:lang` 使用者可選擇或輸入語言，自動轉換為合法 locale 值 | ✓ SATISFIED | mysd-lang.md 提供 1~4 選項 + 自訂 BCP47 輸入；lang set 接受任意 locale 字串 |
| FSCAN-01 | 07-02 | `/mysd:scan` 升級為語言無關通用掃描器 | ✓ SATISFIED | scanner.go 以 detectPrimaryLanguage 支援 go/nodejs/python/unknown；舊 PackageInfo 已完全移除 |
| FSCAN-02 | 07-02 | Scan 偵測語言/模組結構，產生 openspec/specs/ 下 spec 文件 | ✓ SATISFIED | scanner.go 輸出 Modules；mysd-scan.md Step 3 呼叫 mysd-scanner agent 生成 openspec/specs/{module}/spec.md |
| FSCAN-03 | 07-02 | 已存在 openspec/config.yaml 時只增量更新，不覆蓋 config | ✓ SATISFIED | mysd-scan.md line 89: "CRITICAL: Modules with existing specs MUST be skipped — Never overwrite an existing spec." |
| FSCAN-04 | 07-02 | 首次建立 config.yaml 時互動式詢問 locale | ✓ SATISFIED | mysd-scan.md Step 5: 若 config_exists=false，提示使用者執行 /mysd:lang；mysd-init.md Step 2 互動詢問語言 |
| FSCAN-05 | 07-02 | `/mysd:init` 改為 scan --scaffold-only，只建空結構 + 互動式設定 locale | ✓ SATISFIED | init_cmd.go 呼叫 scaffoldOpenSpecDir；mysd-init.md Step 2 互動詢問並呼叫 `mysd lang set` |
| SKILL-01 | 07-05 | Planner 自動依 task 內容推薦 skills 欄位 | ✓ SATISFIED | mysd-planner.md Step 4.5 含 6 項 heuristic 規則並寫入 tasks.md skills 欄位 |
| SKILL-02 | 07-05 | Plan 完成後列出 task-skills 對應表，互動式讓使用者確認 | ✓ SATISFIED | mysd-planner.md Step 7.5 含完整表格輸出及 "Accept all recommended skills? (Y/n)" |
| SKILL-03 | 07-05 | 使用者可逐一調整或批次同意推薦的 skills | ✓ SATISFIED | mysd-planner.md Step 7.5 若 n 則逐 task 詢問 "T1 [{skill}] — change to:" |
| SKILL-04 | 07-05 | ffe 模式跳過互動，直接使用推薦值 | ✓ SATISFIED | mysd-planner.md Step 7.5: "If auto_mode is true (ffe mode): Skip confirmation entirely." |

注意：07-01-PLAN.md 宣告 FCMD-05，但 FilterBlockedTasks 是 executor wave 功能，不直接對應 `/mysd:lang` 的 FCMD-05 語義。此為 Plan 命名不一致問題——FCMD-05 在 REQUIREMENTS.md 中對應 `/mysd:lang` 語言選擇功能，由 07-04 實作；07-01 的 FCMD-05 標記可能是計畫撰寫時的錯誤。FilterBlockedTasks 本身已完整實作並測試通過。

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| `cmd/init_cmd.go` | 42-53 | 仍包含 `yaml.Marshal(cfg)` 及 `os.WriteFile` 寫入 .claude/mysd.yaml | ℹ️ Info | Plan 07-02 要求「移除舊 init 行為」，但此程式碼保留了建立預設 mysd.yaml 的邏輯。SUMMARY 顯示這是刻意保留（保持現有 init 行為）。不影響目標，屬於 scope interpretation 差異 |
| `cmd/lang_test.go` | 106-112 | TestLangSet_AtomicRollback 在 Windows 上 SKIP | ⚠️ Warning | rollback 邏輯的測試在此開發環境無法驗證。程式碼邏輯存在（lang.go:93），但缺乏 Windows 等效的原子性測試 |

### Human Verification Required

#### 1. Atomic Rollback 行為驗證（Linux/macOS）

**Test:** 在 Linux 或 macOS 環境執行 `go test ./cmd/... -run TestLangSet_AtomicRollback -v`
**Expected:** 測試通過，確認若 openspec/ 目錄設為 read-only 導致 config.yaml 寫入失敗時，.claude/mysd.yaml 的 response_language 維持原始值
**Why human:** TestLangSet_AtomicRollback 在 Windows 上因 `os.Chmod` 不阻止寫入而必然 SKIP，無法在此環境自動驗證

### Gaps Summary

無 blocking gap。所有 7 個 observable truths 均已驗證，12 個 requirements 均有程式碼實作佐證。

唯一需要人工確認的項目為 atomic rollback 在 Unix 環境的行為，但此屬 test coverage 問題，非功能缺失——rollback 程式碼邏輯已確認存在於 lang.go:93。

---

_Verified: 2026-03-26T11:00:00Z_
_Verifier: Claude (gsd-verifier)_
