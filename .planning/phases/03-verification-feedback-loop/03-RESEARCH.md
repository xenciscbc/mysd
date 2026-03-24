# Phase 3: Verification & Feedback Loop - Research

**Researched:** 2026-03-24
**Domain:** Go CLI — goal-backward verification engine, spec status write-back, archive gate, UAT generation
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Verification Report 設計**
- D-01: Terminal styled + markdown 檔雙輸出 — Terminal 用 lipgloss 顯示摘要（MUST 5/7 passed），同時寫入 `.specs/changes/{name}/verification.md` 完整報告。仿 GSD 的 VERIFICATION.md 模式
- D-02: 分級顯示 — MUST 先、SHOULD 次、MAY 尾。MUST 全通過才算 overall pass，SHOULD 有警告但不 block，MAY 只備註
- D-03: AI agent 綜合判斷 pass/fail — Verifier agent 檢查 filesystem evidence（檔案存在、grep 關鍵字、test 通過），結合 spec 描述做出判斷
- D-04: 驗證結果寫回 spec frontmatter — 通過的 MUST 設為 DONE，失敗的設為 BLOCKED。Go binary 批量更新，僅更新 MUST 項目

**Gap Report 與 Re-execute 迴路**
- D-05: 失敗項目 + 原因 + 建議修正 — 每個失敗的 MUST 列出：spec 描述、失敗原因（evidence）、AI 建議的修正方向。寫入 `.specs/changes/{name}/gap-report.md`
- D-06: 手動觸發 re-execute — 使用者看完 gap report 後手動執行 `/mysd:execute`，執行時自動讀取 gap report 並只修正失敗項目
- D-07: Re-execute scope 僅限失敗項目 — 只 re-execute 失敗的 MUST items 對應的 tasks，已通過的不重做

**Archive 行為與 UAT 關係**
- D-08: archive 移至 `.specs/archive/{name}/` — 將整個 change 目錄移到 archive，更新 STATE.json phase 為 archived
- D-09: archive 時提示但不強制 UAT — 顯示 'Run UAT first?'，使用者可選 yes/no，無論選什麼都會完成 archive
- D-10: UAT 清單在 verify 過程中自動產生 — verifier agent 偵測 spec 中有 UI 相關 MUST/SHOULD 時，自動產生 UAT checklist 到 `.mysd/uat/`
- D-11: UAT 互動由專屬 agent 引導 — `/mysd:uat` 由 mysd-uat-guide.md agent 引導，逐項顯示測試步驟，使用者回報 pass/fail/skip

**Verifier Agent 獨立性**
- D-12: 全新 agent，只讀 spec 和 filesystem — mysd-verifier.md 不讀 executor 的 alignment.md 或執行歷史
- D-13: 多層次證據 — 檔案存在 + grep 關鍵字碼 + 執行 test suite + 檢查 build 通過
- D-14: Go binary 輸出 verification context JSON — `mysd verify --context-only` 輸出 spec MUST/SHOULD/MAY 清單 JSON，SKILL.md 傳給 verifier agent
- D-15: AI agent 判斷 UI 相關性 — Verifier agent 在驗證過程中用 AI 判斷哪些 MUST/SHOULD 涉及使用者可見的行為

**State Transition 設計**
- D-16: MUST 全通過才 transition — MUST 全通過：executed → verified。有失敗：維持 executed 狀態
- D-17: archive gate 雙重檢查 — Go binary 在 archive 前檢查：(1) state == verified，(2) 所有 MUST items status == DONE

### Claude's Discretion

- Verification report 的具體 markdown 模板
- Gap report 的詳細程度和建議修正的具體寫法
- UAT checklist 的格式和測試步驟描述
- mysd-verifier.md 的具體 prompt 措辭
- mysd-uat-guide.md 的互動 prompt 設計
- archive 時 'Run UAT first?' 的呈現方式

### Deferred Ideas (OUT OF SCOPE)

- **Playwright 腳本自動產生** — UAT 清單如果偵測到 web UI 相關項目，可自動產生 Playwright e2e test 腳本。超出 Phase 3 核心驗證流程範圍，可作為未來 phase 功能

</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| VRFY-01 | Goal-backward verification parses all MUST items from spec and generates verification checklist | `spec.ParseChange()` 已能解析所有 MUST items；VerificationContext JSON 結構由 `internal/verifier/context.go` 提供（新建，仿 executor/context.go 模式） |
| VRFY-02 | Verification uses an independent fresh-context agent (not the same agent that executed) | mysd-verifier.md 新 agent，不讀 alignment.md；由 mysd-verify.md SKILL.md 透過 Task tool 啟動 |
| VRFY-03 | SHOULD items verified with lower priority; MAY items noted but not required | VerificationContext 已包含分層 must_items / should_items / may_items；verifier agent 依優先級依序驗證 |
| VRFY-04 | Verification produces structured pass/fail report per MUST/SHOULD/MAY item | verifier agent 輸出結構化 JSON report；Go binary `mysd verify --write-results` 讀取並寫入 verification.md |
| VRFY-05 | Failed MUST items trigger gap report that can feed back into re-execution | `gap-report.md` 在 verify 完成後由 Go binary 依 verifier report JSON 產生；execute SKILL.md 讀取此檔案縮小 scope |
| SPEC-05 | Verification results are automatically written back to spec status | `internal/spec/updater.go` 的 `UpdateItemStatus`（新建，仿 UpdateTaskStatus 模式）批量更新 spec requirement status |
| SPEC-06 | Completed specs can be archived to `.specs/archive/` via `/mysd:archive` command | `cmd/archive.go` stub 實作：double gate check → os.Rename changeDir → archiveDir → SaveState |
| WCMD-06 | `/mysd:verify` — goal-backward verification of all MUST items | `mysd-verify.md` SKILL.md + `cmd/verify.go` 實作，遵循 execute 的 thin-command-layer 模式 |
| WCMD-07 | `/mysd:archive` — archive completed spec to history | `mysd-archive.md` SKILL.md + `cmd/archive.go` 實作，含 UAT prompt |
| WCMD-12 | `/mysd:uat` — 互動式 UAT 流程 | `mysd-uat.md` SKILL.md + `mysd-uat-guide.md` agent definition |
| UAT-01 | 驗證階段可選擇產生互動式 UAT 驗收清單 | verifier agent 在驗證過程中偵測 UI 相關項目，寫入 `.mysd/uat/{change_name}-uat.md` |
| UAT-02 | UAT 清單為可選步驟，不是 archive 的前提條件 | archive gate 只檢查 state==verified + MUST all DONE，UAT 狀態不在 gate 條件內 |
| UAT-03 | UAT 清單存放於 `.mysd/uat/` 目錄，可跨 session 保留 | 用 `os.MkdirAll(".mysd/uat", 0755)` 確保目錄存在；YAML frontmatter 記錄 last_run |
| UAT-04 | 使用者可透過 `/mysd:uat` 獨立觸發 UAT 流程，可重複執行 | mysd-uat-guide.md agent 讀取現有 UAT 檔案並覆寫結果 |
| UAT-05 | UAT 清單記錄每次執行的結果（通過/未通過/跳過）與時間戳 | UAT 檔案 YAML frontmatter 含 `results[]` 陣列，每項含 status + timestamp |

</phase_requirements>

---

## Summary

Phase 3 在已有的 spec parsing、state machine、lipgloss output 基礎上，增加三個主要系統：

1. **Verification Engine**（`internal/verifier/`）— 仿 `internal/executor/` 模式，新建 `VerificationContext` struct 和 `BuildVerificationContext()` 函式，輸出 JSON 給 mysd-verifier agent。核心挑戰是 AI verifier agent 的獨立性（不讀 alignment.md）與多層次 evidence 收集（file existence + grep + test run + build check）。

2. **Spec Status Write-back**（`internal/spec/updater.go` 擴充）— 新增 `UpdateItemStatus()` 仿現有 `UpdateTaskStatus()` 模式，但操作 spec 檔案中的 requirement status，而非 tasks.md frontmatter。關鍵挑戰：spec 的 requirement 沒有 stable ID（見 parser.go），需要以 text + keyword 作為匹配 key 或引入 ID 方案。

3. **Archive Gate**（`cmd/archive.go`）— double-gate 邏輯：state == verified AND all MUST items status == DONE。使用 `os.Rename()` 做原子性目錄移動，更新 STATE.json。UAT prompt 以互動式 confirm 呈現，不阻塞流程。

**Primary recommendation:** 優先解決 requirement ID 方案（目前 `Requirement.ID` 是空字串），因為這是 spec status write-back（SPEC-05）的前提。建議採用 `{capability-name}::{line-hash}` 或在 spec 檔案中加入 `req-id` frontmatter 欄位。

---

## Standard Stack

### Core (已在 go.mod 中，無需新增)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `github.com/mysd/internal/spec` | (internal) | Requirement 解析、spec 檔案讀寫 | ParseChange()、ParseSpec()、UpdateTaskStatus() 模式均已驗證可用 |
| `github.com/mysd/internal/state` | (internal) | Phase transition、STATE.json 讀寫 | ValidTransitions 已含 PhaseExecuted → PhaseVerified → PhaseArchived |
| `github.com/charmbracelet/lipgloss` | v1.1.0 | Terminal styled output | 已在 executor/status.go 使用，顏色風格已建立 |
| `encoding/json` | stdlib | VerificationContext JSON 序列化 | 與 ExecutionContext 同模式，無需新增依賴 |
| `os` stdlib | stdlib | 目錄移動（Rename）、檔案讀寫 | archive 的 os.Rename() 是最簡潔的 atomic directory move |
| `gopkg.in/yaml.v3` | v3.0.1 | UAT 檔案 frontmatter 讀寫 | 已在 go.mod，spec updater 使用同款 |

### 無需新增依賴

Phase 3 的所有功能均可用現有依賴實現。UAT 檔案格式使用與 spec 相同的 YAML frontmatter + markdown body。驗證 report 是純 markdown 檔案寫入。

---

## Architecture Patterns

### 新增 Package 結構

```
internal/
├── verifier/
│   ├── context.go       # VerificationContext struct + BuildVerificationContext()
│   ├── context_test.go
│   ├── report.go        # ParseVerifierReport()、UpdateSpecFromReport()
│   └── report_test.go
├── spec/
│   ├── updater.go       # 擴充：新增 UpdateItemStatus() + WriteSpec()
│   └── updater_test.go  # 新增對應測試
└── uat/
    ├── checklist.go     # UAT 檔案格式、ReadUAT()、WriteUAT()
    └── checklist_test.go

cmd/
├── verify.go            # 實作 --context-only flag 和 --write-results flag
└── archive.go           # 實作 double gate + os.Rename + UAT prompt

.claude/
├── commands/
│   ├── mysd-verify.md   # SKILL.md：context-only → verifier agent → write-results
│   ├── mysd-archive.md  # SKILL.md：UAT prompt → archive binary call
│   └── mysd-uat.md      # SKILL.md：讀取 UAT 檔案 → uat-guide agent
└── agents/
    ├── mysd-verifier.md     # 獨立驗證 agent
    └── mysd-uat-guide.md    # UAT 互動引導 agent
```

### Pattern 1: Thin Command Layer（與 Phase 2 一致）

**What:** cmd/verify.go 和 cmd/archive.go 不含業務邏輯，所有邏輯在 internal/ packages。
**When to use:** 所有 cmd/ 層實作。

```go
// cmd/verify.go
var verifyCmd = &cobra.Command{
    Use:   "verify [change-name]",
    Short: "Goal-backward verification of MUST items",
    RunE: func(cmd *cobra.Command, args []string) error {
        contextOnly, _ := cmd.Flags().GetBool("context-only")
        writeResults, _ := cmd.Flags().GetString("write-results")

        specsDir, _ := spec.DetectSpecDir(".")
        ws, _ := state.LoadState(specsDir)

        if contextOnly {
            return runVerifyContextOnly(cmd.OutOrStdout(), specsDir, ws)
        }
        if writeResults != "" {
            return runVerifyWriteResults(cmd.OutOrStdout(), specsDir, ws, writeResults)
        }
        return fmt.Errorf("use --context-only or --write-results; direct execution via /mysd:verify")
    },
}
```

### Pattern 2: VerificationContext JSON（仿 ExecutionContext）

**What:** Go binary 輸出 verification context JSON，SKILL.md 傳給 verifier agent。

```go
// internal/verifier/context.go
type VerificationContext struct {
    ChangeName   string            `json:"change_name"`
    MustItems    []VerifyItem      `json:"must_items"`
    ShouldItems  []VerifyItem      `json:"should_items"`
    MayItems     []VerifyItem      `json:"may_items"`
    TasksSummary []TaskSummaryItem `json:"tasks_summary"`
    SpecsDir     string            `json:"specs_dir"`
    ChangeDir    string            `json:"change_dir"`
}

type VerifyItem struct {
    ID      string `json:"id"`
    Text    string `json:"text"`
    Status  string `json:"status"`   // current spec status
    Keyword string `json:"keyword"`  // MUST / SHOULD / MAY
}
```

### Pattern 3: Verifier Report JSON（agent → binary 的回傳格式）

**What:** mysd-verifier agent 完成驗證後輸出結構化 JSON，binary 讀取並執行 write-back。

```json
{
  "change_name": "my-feature",
  "overall_pass": false,
  "must_pass": false,
  "results": [
    {
      "id": "req-1",
      "text": "System MUST validate input",
      "keyword": "MUST",
      "pass": true,
      "evidence": "Found input_validator.go:45, tests pass",
      "suggestion": ""
    },
    {
      "id": "req-2",
      "text": "System MUST log errors",
      "keyword": "MUST",
      "pass": false,
      "evidence": "No logger usage found in error handlers",
      "suggestion": "Add zerolog or stdlib log calls in error branches"
    }
  ],
  "has_ui_items": false,
  "ui_items": []
}
```

### Pattern 4: Spec Requirement ID 方案（關鍵設計決策）

**Problem:** `Requirement.ID` 目前在 parser.go 中總是空字串（`ID: ""`）。verification write-back 需要穩定 ID 來匹配 spec 項目。

**Solution（建議）:** 以 `{spec-file-basename}::req-{sequential-number}` 作為 runtime ID，在 `ParseSpec()` 時填入（不修改 .md 檔案，不需要新 frontmatter 欄位）：

```go
// internal/spec/parser.go 中修改 parseRequirementsFromBody
reqID++
reqs = append(reqs, Requirement{
    ID:      fmt.Sprintf("%s::req-%d", specBasename, reqID),
    Text:    strings.TrimSpace(line),
    ...
})
```

這個方案：
- 不需要修改任何 spec .md 格式
- ID 在同一個 parse run 中穩定
- verify binary 和 verifier agent 用同一個 context JSON，ID 在 JSON 中傳遞給 agent，agent 在 report 中回傳同一個 ID

### Pattern 5: Spec Status Write-back

**What:** 新增 `UpdateItemStatus()` 到 `internal/spec/updater.go`，仿 `UpdateTaskStatus()` 模式。

**核心挑戰:** spec 的 requirements 在 spec.md body 中（非 frontmatter），無法用 YAML round-trip 更新。需要在 spec 的 SpecFrontmatter 或新欄位記錄 status，或在另一個 sidecar 檔案記錄。

**建議方案（sidecar JSON）:** 在 changeDir 建立 `.specs/changes/{name}/verification-status.json`，記錄每個 requirement ID 的最新 status。`mysd status` 讀取此 sidecar 來顯示 MUST done/pending 計數。這樣不修改 spec.md，符合「spec 是唯讀事實來源」的設計哲學。

```json
{
  "change_name": "my-feature",
  "verified_at": "2026-03-24T10:00:00Z",
  "requirements": {
    "capability-a/spec.md::req-1": "done",
    "capability-a/spec.md::req-2": "blocked"
  }
}
```

### Pattern 6: Archive Double Gate

```go
// cmd/archive.go RunE
func runArchive(...) error {
    // Gate 1: state must be verified
    if ws.Phase != state.PhaseVerified {
        return fmt.Errorf("cannot archive: phase is %s, must be verified", ws.Phase)
    }

    // Gate 2: all MUST items must be DONE
    mustItems := filterByKeyword(reqs, spec.Must)
    for _, r := range mustItems {
        if statusMap[r.ID] != spec.StatusDone {
            return fmt.Errorf("cannot archive: MUST item %q is not done", r.Text[:50])
        }
    }

    // UAT prompt
    if isInteractive {
        fmt.Print("Run UAT first? [y/N] ")
        // read response, if yes: print hint to run /mysd:uat
        // regardless of answer, continue to archive
    }

    // Archive: move directory
    archiveDir := filepath.Join(specsDir, "archive", changeName)
    if err := os.MkdirAll(filepath.Dir(archiveDir), 0755); err != nil { ... }
    if err := os.Rename(changeDir, archiveDir); err != nil { ... }

    // Update state
    ws.Phase = state.PhaseArchived
    return state.SaveState(specsDir, ws)
}
```

### Pattern 7: UAT 檔案格式

```markdown
---
spec-version: "1"
change: my-feature
generated: 2026-03-24
last_run: 2026-03-24T10:30:00Z
summary:
  total: 5
  pass: 3
  fail: 1
  skip: 1
results:
  - id: "uat-1"
    description: "用戶可以在首頁看到登入按鈕"
    status: "pass"
    run_at: "2026-03-24T10:30:00Z"
  - id: "uat-2"
    description: "點擊登入按鈕跳轉至登入頁"
    status: "fail"
    notes: "按鈕跳轉到錯誤頁面"
    run_at: "2026-03-24T10:31:00Z"
---

## UAT Checklist: my-feature

_Generated by mysd-verifier. Run `/mysd:uat` to step through interactively._

### UI Acceptance Tests

- [ ] 用戶可以在首頁看到登入按鈕
- [ ] 點擊登入按鈕跳轉至登入頁
...
```

### Pattern 8: SKILL.md 流程（mysd-verify.md）

與 mysd-execute.md 完全相同的 orchestrator pattern：

```markdown
## Step 1: Get Verification Context
Run: mysd verify --context-only
Parse the JSON output.

## Step 2: Invoke Verifier Agent
Use the Task tool to invoke mysd-verifier agent with context JSON.

## Step 3: Write Results Back
After verifier completes, run:
  mysd verify --write-results /tmp/verifier-report.json
This updates spec status and transitions state if all MUST pass.
```

### Anti-Patterns to Avoid

- **Verifier reads alignment.md:** 違反 D-12 的獨立性原則。verifier 只讀 spec + filesystem。
- **archive gate 只檢查 state，不檢查 MUST items:** state == verified 不代表 re-verify 後沒有回退，需雙重 gate（D-17）。
- **os.Rename() 跨 volume:** 在 Windows 上，如果 `.specs/` 和 `.specs/archive/` 在不同 drive letter，Rename 會失敗。應先嘗試 Rename，失敗再 fallback 到 copy + delete。
- **Requirement ID 不穩定:** 如果 spec.md 內容被編輯，sequential ID 會漂移。sidecar JSON 中的 key 應使用 text hash 而非單純序號。
- **直接修改 spec.md 的 MUST item status:** spec.md 是使用者編寫的文件，不應被 binary 靜默修改。用 sidecar JSON 記錄 verification status 更安全。

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Terminal color output | 自訂 ANSI escape 碼 | `lipgloss`（已安裝） | lipgloss 處理了 TTY detection、Windows ANSI、color profile |
| YAML frontmatter parsing | 自訂 --- 分割器 | `adrg/frontmatter`（已安裝） | 已處理邊界 case（多行 value、嵌套 YAML） |
| Directory copy for archive | 遞迴 os.ReadDir + copy | `os.Rename()`（stdlib） | 在同 volume 時是 atomic 操作；不需要額外依賴 |
| AI evidence collection format | 複雜自訂格式 | 結構化 JSON（同 ExecutionContext 模式） | 已驗證的模式，binary ↔ agent 通訊一致 |

---

## Common Pitfalls

### Pitfall 1: Requirement ID 在 verify cycle 中漂移
**What goes wrong:** 用戶在 execute 之後、verify 之前編輯了 spec.md，新增了一個 MUST item。sequential ID 從 req-2 開始，原本的 req-2 變成 req-3，verification status JSON 的 key 對不上。
**Why it happens:** parser.go 用 sequential counter，沒有 content-based ID。
**How to avoid:** sidecar JSON 的 key 使用 `{filename}::{keyword}::{text-truncated-hash}` 而非純序號。hash 用 `fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(r.Text)))[:6]` 即可，無需引入新依賴。
**Warning signs:** verify 後 MUST count 與 context-only 輸出的 count 不一致。

### Pitfall 2: os.Rename 在 Windows 跨 volume 失敗
**What goes wrong:** `os.Rename(".specs/changes/foo", ".specs/archive/foo")` 在 Windows 上若兩個路徑在不同 drive 會回傳 `invalid cross-device link`。
**Why it happens:** Windows os.Rename 底層 MoveFile 不允許跨 volume。
**How to avoid:** archive.go 捕捉 rename error，檢查是否為 `syscall.EXDEV`（errno 18），若是則 fallback 到 recursive copy + delete。
**Warning signs:** 在 Windows 開發機上 archive 失敗，Linux/macOS 正常。

### Pitfall 3: Verifier agent 自我驗證盲點（STATE.md 中已記錄的 blocker）
**What goes wrong:** verifier agent 傾向於「相信」executor 已完成的工作，不做實際 evidence 收集，直接 pass 所有項目。
**Why it happens:** LLM training bias — 助手傾向於確認，而非質疑。
**How to avoid:** mysd-verifier.md prompt 明確指定：「你必須為每個 MUST item 找到至少一個具體 evidence（檔案路徑 + 行號、test output、build output）。找不到 evidence 即為 FAIL，不得假設已完成。」加入 evidence 格式要求：`evidence: "internal/foo/bar.go:42 — function validateInput found"`。
**Warning signs:** 所有項目都 pass，但 gap report 為空。

### Pitfall 4: UAT 檔案被覆蓋而非 append
**What goes wrong:** 每次 verify 都重新產生 UAT 檔案，覆蓋了之前手動記錄的 pass/fail 結果（UAT-05 要求記錄歷史）。
**Why it happens:** 簡單 os.WriteFile 覆蓋模式。
**How to avoid:** `internal/uat/checklist.go` 的 WriteUAT 在檔案已存在時，保留 results 歷史，只更新 `last_run` 和 summary。新增 `run_history` 陣列記錄每次完整執行結果。
**Warning signs:** UAT-05 測試失敗。

### Pitfall 5: archive 後 STATE.json 仍在 changes/ 目錄
**What goes wrong:** archive 把整個 change 目錄移到 archive/，但 STATE.json 在 `.specs/STATE.json`（changeDir 的上層），不會被移動。STATE.json 仍顯示 phase=verified，但 change 目錄已不存在。
**Why it happens:** STATE.json 和 changeDir 是不同層級。
**How to avoid:** archive RunE 在 os.Rename 後，先 transition ws.Phase = PhaseArchived，再 SaveState。同時在 archive 目錄內複製一份 STATE.json snapshot（`archive/{name}/ARCHIVED-STATE.json`）作為歷史記錄。
**Warning signs:** archive 後 `mysd status` 報錯找不到 change 目錄。

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | go test (stdlib) + testify v1.11.1 |
| Config file | none — `go test ./...` |
| Quick run command | `go test ./internal/verifier/... ./internal/uat/... ./cmd/ -count=1` |
| Full suite command | `go test ./... -count=1` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|--------------|
| VRFY-01 | BuildVerificationContext() 從 change 目錄正確解析 MUST/SHOULD/MAY | unit | `go test ./internal/verifier/ -run TestBuildVerificationContext -v` | Wave 0 |
| VRFY-02 | verifier agent definition 獨立（不讀 alignment.md） | manual | 觀察 agent output 無 alignment.md 引用 | manual-only |
| VRFY-03 | context JSON 的 should_items/may_items 分層正確 | unit | `go test ./internal/verifier/ -run TestVerificationContext_ItemClassification -v` | Wave 0 |
| VRFY-04 | ParseVerifierReport() 正確解析 agent 輸出的 JSON | unit | `go test ./internal/verifier/ -run TestParseVerifierReport -v` | Wave 0 |
| VRFY-05 | WriteGapReport() 在失敗 MUST items 時產生 gap-report.md | unit | `go test ./internal/verifier/ -run TestWriteGapReport -v` | Wave 0 |
| SPEC-05 | UpdateItemStatus() 更新 verification-status.json | unit | `go test ./internal/spec/ -run TestUpdateItemStatus -v` | Wave 0 |
| SPEC-06 | archive 指令移動目錄並更新 STATE.json | unit | `go test ./cmd/ -run TestArchiveCmd -v` | Wave 0 |
| WCMD-06 | `mysd verify --context-only` 輸出合法 JSON | integration | `go test ./cmd/ -run TestVerifyContextOnly -v` | Wave 0 |
| WCMD-07 | `mysd archive` 在非 verified 狀態回傳錯誤 | unit | `go test ./cmd/ -run TestArchiveGate -v` | Wave 0 |
| WCMD-12 | mysd-uat.md SKILL.md 格式合法 | smoke | `grep -n "allowed-tools" .claude/commands/mysd-uat.md` | Wave 0 |
| UAT-01 | verifier context JSON 含 has_ui_items + ui_items 欄位 | unit | `go test ./internal/verifier/ -run TestUIDetection -v` | Wave 0 |
| UAT-02 | archive gate 不檢查 UAT status | unit | `go test ./cmd/ -run TestArchiveGateNoUAT -v` | Wave 0 |
| UAT-03 | WriteUAT() 寫入正確目錄並 cross-session 保留 | unit | `go test ./internal/uat/ -run TestWriteUAT -v` | Wave 0 |
| UAT-04 | ReadUAT() 可讀取現有 UAT 檔案 | unit | `go test ./internal/uat/ -run TestReadUAT -v` | Wave 0 |
| UAT-05 | WriteUAT() 保留歷史 results，不覆蓋舊記錄 | unit | `go test ./internal/uat/ -run TestUATHistoryPreservation -v` | Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./internal/verifier/... ./internal/uat/... -count=1`
- **Per wave merge:** `go test ./... -count=1`
- **Phase gate:** `go test ./... -count=1` full suite green before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/verifier/context.go` + `context_test.go` — covers VRFY-01, VRFY-03, UAT-01
- [ ] `internal/verifier/report.go` + `report_test.go` — covers VRFY-04, VRFY-05
- [ ] `internal/uat/checklist.go` + `checklist_test.go` — covers UAT-03, UAT-04, UAT-05
- [ ] `internal/spec/updater.go` 新增 `UpdateItemStatus()` + 對應測試 — covers SPEC-05
- [ ] `cmd/verify_test.go` — covers WCMD-06
- [ ] `cmd/archive_test.go` — covers SPEC-06, WCMD-07, UAT-02

---

## Code Examples

### BuildVerificationContext (新建，仿 BuildContextFromParts)

```go
// internal/verifier/context.go
// Source: 仿 internal/executor/context.go BuildContextFromParts 模式

type VerificationContext struct {
    ChangeName   string        `json:"change_name"`
    ChangeDir    string        `json:"change_dir"`
    SpecsDir     string        `json:"specs_dir"`
    MustItems    []VerifyItem  `json:"must_items"`
    ShouldItems  []VerifyItem  `json:"should_items"`
    MayItems     []VerifyItem  `json:"may_items"`
    TasksSummary []TaskItem    `json:"tasks_summary"`
}

type VerifyItem struct {
    ID      string `json:"id"`      // e.g. "capability-a/spec.md::must-a1b2c3"
    Text    string `json:"text"`
    Keyword string `json:"keyword"` // "MUST" / "SHOULD" / "MAY"
    Status  string `json:"status"`  // current status from verification-status.json
}

func BuildVerificationContext(specsDir, changeName string) (VerificationContext, error) {
    changeDir := filepath.Join(specsDir, "changes", changeName)
    change, err := spec.ParseChange(changeDir)
    if err != nil {
        return VerificationContext{}, fmt.Errorf("parse change: %w", err)
    }

    ctx := VerificationContext{
        ChangeName: changeName,
        ChangeDir:  changeDir,
        SpecsDir:   specsDir,
    }

    for _, r := range change.Specs {
        item := VerifyItem{
            ID:      stableID(r),
            Text:    r.Text,
            Keyword: string(r.Keyword),
            Status:  string(r.Status),
        }
        switch r.Keyword {
        case spec.Must:
            ctx.MustItems = append(ctx.MustItems, item)
        case spec.Should:
            ctx.ShouldItems = append(ctx.ShouldItems, item)
        case spec.May:
            ctx.MayItems = append(ctx.MayItems, item)
        }
    }
    return ctx, nil
}

// stableID 用 text CRC32 hash 產生穩定 ID，避免 sequential counter 漂移
func stableID(r spec.Requirement) string {
    h := crc32.ChecksumIEEE([]byte(r.Text))
    return fmt.Sprintf("%s::%s-%x", r.SourceFile, strings.ToLower(string(r.Keyword)), h)
}
```

### UpdateItemStatus (新建，仿 UpdateTaskStatus)

```go
// internal/spec/updater.go 擴充
// Source: 仿現有 UpdateTaskStatus 模式

// VerificationStatus holds per-requirement verification results as a sidecar JSON.
type VerificationStatus struct {
    ChangeName   string                 `json:"change_name"`
    VerifiedAt   time.Time              `json:"verified_at"`
    Requirements map[string]ItemStatus  `json:"requirements"` // key = stableID
}

func ReadVerificationStatus(changeDir string) (VerificationStatus, error) { ... }
func WriteVerificationStatus(changeDir string, vs VerificationStatus) error { ... }
func UpdateItemStatus(changeDir string, reqID string, newStatus ItemStatus) error { ... }
```

### os.Rename fallback for cross-volume (Windows)

```go
// cmd/archive.go
func moveDir(src, dst string) error {
    if err := os.Rename(src, dst); err == nil {
        return nil
    }
    // Fallback: copy + remove (handles cross-volume on Windows)
    if err := copyDir(src, dst); err != nil {
        return fmt.Errorf("copy to archive: %w", err)
    }
    return os.RemoveAll(src)
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Requirement.ID 空字串 | stableID() 用 CRC32 hash | Phase 3 引入 | 啟用 spec status write-back（SPEC-05） |
| ExecutionContext 無驗證資料 | 新增 VerificationContext（獨立 struct） | Phase 3 | 清晰分離 execute 與 verify 的 context |
| STATE.json 無驗證結果欄位 | `VerifyPass *bool` 已在 WorkflowState（Phase 2 已加） | Phase 2 | archive gate 可讀取此欄位 |

---

## Open Questions

1. **Requirement.SourceFile 欄位是否需要加到 schema.go？**
   - What we know: 現有 `Requirement` struct 沒有 `SourceFile` 欄位；stableID() 需要知道是哪個 spec.md 產生的
   - What's unclear: 是在 `parser.go` 的 `ParseSpec()` 時填入，還是在 `ParseChange()` 時 post-process？
   - Recommendation: 在 `ParseSpec(path string)` 時填入 `SourceFile: filepath.Base(path)`，最簡單

2. **gap-report.md 的格式 vs 現有 alignment.md 格式的一致性**
   - What we know: executor SKILL.md 在 re-execute 時需要讀取 gap-report.md 並縮小 scope（D-06, D-07）；alignment.md 是 executor 自己產生的
   - What's unclear: gap-report.md 應是 binary 產生的結構化 markdown，還是給 agent 讀的自由格式文件？
   - Recommendation: binary 產生結構化 markdown（含 YAML frontmatter 標記哪些 task ID 需要重做），executor agent 在讀到 gap-report.md 後只執行 `failed_task_ids` 中的 tasks

3. **mysd-verify.md SKILL.md 如何讓 verifier agent 輸出 JSON？**
   - What we know: execute SKILL.md 透過 Task tool 啟動 executor agent，executor agent 輸出文字 summary；Phase 3 需要 verifier agent 輸出可被 binary 解析的 JSON
   - What's unclear: agent 能否被指示「輸出到檔案」vs「輸出到 stdout」？
   - Recommendation: verifier agent 被指示將 report JSON 寫入 `.specs/changes/{name}/verifier-report.json`（Write tool），SKILL.md 在 agent 完成後呼叫 `mysd verify --write-results .specs/changes/{name}/verifier-report.json`

---

## Environment Availability

Phase 3 是純 Go code + Claude Code plugin 檔案的新增，不依賴外部服務或新工具。

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| go | Build | ✓（per go.mod go 1.25.5） | 1.25.5 | — |
| lipgloss v1.1.0 | Terminal output | ✓（in go.mod） | v1.1.0 | — |
| adrg/frontmatter v0.2.0 | UAT YAML parsing | ✓（in go.mod） | v0.2.0 | — |
| gopkg.in/yaml.v3 | VerificationStatus JSON | ✓（in go.mod） | v3.0.1 | — |

**Missing dependencies:** None — all requirements can be satisfied with existing go.mod dependencies.

---

## Sources

### Primary (HIGH confidence)

- `internal/executor/context.go` — ExecutionContext 模式，VerificationContext 直接參考
- `internal/spec/updater.go` — UpdateTaskStatus 模式，UpdateItemStatus 直接參考
- `internal/state/transitions.go` — ValidTransitions 確認 PhaseExecuted → PhaseVerified → PhaseArchived 已支援
- `internal/state/state.go` — WorkflowState.VerifyPass *bool 欄位已存在
- `.claude/agents/mysd-executor.md` — verifier agent 格式的直接參考範本
- `.claude/commands/mysd-execute.md` — SKILL.md orchestrator pattern 的直接參考範本
- `internal/spec/parser.go` — 確認 Requirement.ID 目前為空字串，需要解決

### Secondary (MEDIUM confidence)

- STATE.md 中的 blocker 記錄 — 「Phase 3: Verification prompting strategy to avoid AI self-verification blindness needs phase research」

### Tertiary (LOW confidence)

- Go `os.Rename` cross-volume behavior on Windows — 基於 Go stdlib 文件和 syscall 知識，未在此環境實際測試跨 volume

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — 全部使用現有依賴，無新引入
- Architecture: HIGH — 基於 Phase 2 已驗證的 thin-command-layer + context JSON pattern
- Pitfalls: MEDIUM — Windows cross-volume Rename 為 MEDIUM（未在此環境實測）；其他基於程式碼審查

**Research date:** 2026-03-24
**Valid until:** 2026-06-24（穩定 Go stdlib 模式，無 fast-moving 依賴）
