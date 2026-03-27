# Phase 11: agent-doc — Context

**Gathered:** 2026-03-27 (updated)
**Status:** Ready for planning

<domain>
## Phase Boundary

完善 mysd 工作流程的自動化串接（propose→spec、apply→verify）+ executor failure sidecar 寫入機制（讓 fix agent 有上下文可讀）+ archive 後的 doc 維護流程（可設定哪些文件需要更新）+ Phase 9-04 plugin sync 補足。

此 phase 不引入新 agent role，也不更動現有 spec/plan/archive 的核心邏輯。改動限於 SKILL.md 層、agent prompt 層，以及一個小幅 binary config struct 擴充（`DocsToUpdate` 欄位）。

</domain>

<decisions>
## Implementation Decisions

### 工作流程自動化串接

- **D-01:** `propose` 完成後（Step 10 之後）自動呼叫 `mysd-spec-writer` agent，不需要使用者手動執行 `/mysd:spec`。spec-writer 完成後：(a) 顯示生成的 spec 內容摘要（MUST/SHOULD/MAY 要求數量及關鍵點），(b) 提示使用者可選的後續指令清單及各指令用途說明（如 `/mysd:plan` — 建立執行計劃、`/mysd:design` — 補充設計決策等）。
- **D-02:** `apply` 所有 tasks 完成後（Step 4 之後）自動執行 verify — 先執行 `go build ./...` + `go test ./...`，再呼叫 `mysd-verifier` agent 做 MUST coverage check。
- **D-03:** `archive` 不增加額外 verify 步驟 — 維持現有邏輯（binary 已在 not-verified 狀態時 block archive）。apply 的自動 verify 已覆蓋覆蓋率檢查。
- **D-04:** `propose` 的自動 spec 呼叫可以被 `--skip-spec` flag 跳過（for cases where proposal is a skeleton only）。
- **D-05:** `apply` 的自動 verify 在 `--auto` 模式下不詢問確認，直接執行。

### Executor Failure Sidecar

- **D-06:** executor agent 在 task 執行失敗時（build error、test failure、中途退出）在 agent prompt 內的 on-failure 段落寫入 sidecar 文件：`.specs/changes/{change_name}/.sidecar/T{id}-failure.md`。寫入邏輯在 `mysd-executor.md` agent prompt 內，executor agent 自己偵測並寫入（不由外部 orchestrator 注入）。
- **D-07:** sidecar 包含：失敗時間戳、任務描述、build/test error output（直接在 agent context 中可取得）、AI 診斷嘗試（if any）
- **D-08:** `mysd-fix` Step 5B 讀取 T{id}-failure.md 作為初始 diagnosis context，若 sidecar 不存在則降級為無 context 診斷（backward compat）
- **D-09:** `.sidecar/` 目錄加入專案根目錄 `.gitignore`（failure context 是本地暫存，不應提交）

### Doc 維護流程

- **D-10:** `mysd.yaml` 新增 `docs_to_update` 欄位（字串陣列，指定需要在 archive 後更新的文件路徑，如 `["README.md", "CHANGELOG.md"]`）
- **D-11:** archive 完成後，SKILL.md 讀取 `docs_to_update` 配置，對每個指定文件呼叫 LLM 更新。更新策略依文件類型自適應：(a) `CHANGELOG.md` → LLM 生成新條目後 **prepend 到頂部**，保留舊內容不動；(b) `README.md` → **全文重寫**；(c) 其他自訂文件 → LLM 從檔名和現有內容推斷更新方式（全文重寫為預設）。
- **D-11b:** LLM 更新每個文件時讀入的 context：`proposal.md`（what + why）+ `tasks.md`（what was done）+ `specs/` 目錄（MUST/SHOULD/MAY 要求）+ **現有文件內容**（不能不看就重寫）。
- **D-12:** `DocsToUpdate []string` 加進 `internal/config/defaults.go` 的 `ProjectConfig` struct，並透過 `mysd execute --context-only` JSON 輸出暴露。SKILL.md 從 binary 讀取此值，不直接解析 YAML。文件更新本身（Read/Edit/Write）在 SKILL.md 層實作，不需要新的 binary subcommand。
- **D-13:** 更新前向使用者展示文件清單確認（顯示「將更新以下文件：README.md, CHANGELOG.md」），使用者按 Enter 確認或輸入 `n` 跳過。`--auto` 模式跳過確認直接更新。
- **D-14:** `docs_to_update` 未設定時，archive 後正常結束（不提示 doc 更新）— convention over config

### Plugin Sync 補足

- **D-15:** 完成 Phase 9-04 遺漏的 plugin sync：同步範圍限定為 `mysd-*.md` 檔案 — 將 `.claude/commands/mysd-*.md` 同步到 `plugin/commands/`，`.claude/agents/mysd-*.md` 同步到 `plugin/agents/`。排除 `gsd-*.md`（GSD 框架的 agents，不屬於 mysd distribution）、`CLAUDE.md`（兩側各自管理）、`gsd/`/`spectra/` 等子目錄。目前缺少的 `mysd-lang.md`、`mysd-model.md` 需補入 `plugin/commands/`。`mysd-designer.md` 兩側有差異需對齊。
- **D-16:** 同步後執行 diff 確認兩側一致（`diff .claude/commands/mysd-X.md plugin/commands/mysd-X.md` 應為 zero output）

### Claude's Discretion

- Executor sidecar 的具體 Markdown 格式（結構自由，確保 fix agent 可讀即可）
- auto-verify 失敗時若 build 錯誤的 UX（提示用 `/mysd:fix` 還是直接顯示 build output）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### SKILL.md 修改目標

- `.claude/commands/mysd-propose.md` — Step 10 後需新增 Step 11 auto-spec 呼叫（D-01/D-04）；Step 11 完成後顯示 spec 摘要 + 後續指令清單
- `.claude/commands/mysd-apply.md` — Step 4 後需新增 Step 5 auto-verify 流程（D-02/D-05）；build 失敗跳過 verifier
- `.claude/commands/mysd-archive.md` — Step 1 後需讀取 docs_to_update 並更新文件（D-11~D-14）
- `.claude/commands/mysd-ff.md` — 需同步加入 auto-verify 邏輯（D-02）
- `.claude/commands/mysd-ffe.md` — 需同步加入 auto-verify 邏輯（D-02）

### Agent 修改目標

- `.claude/agents/mysd-executor.md` — 新增 on-failure 段落（D-06/D-07）：build/test 失敗時寫 .sidecar/T{id}-failure.md
- `plugin/agents/mysd-executor.md` — 同上（distribution copy，與 .claude/agents/ 同步）
- `.claude/commands/mysd-fix.md` — Step 5B 確認 sidecar 讀取路徑正確（D-08；已有框架，需對齊 D-06 的路徑格式）

### Binary 配置修改

- `internal/config/defaults.go` — 新增 `DocsToUpdate []string` 欄位（D-12）
- `internal/executor/` 或對應的 `--context-only` 輸出路徑 — 需包含 docs_to_update 在 JSON 輸出（D-12）

### Plugin Sync 來源

- `.claude/commands/mysd-*.md` — 19+ 個 SKILL.md commands（authoritative dev copy）
- `.claude/agents/mysd-*.md` — 12 個 mysd agent definitions（authoritative dev copy）
- `plugin/commands/mysd-*.md` — distribution copy（需對齊 .claude/commands/）
- `plugin/agents/mysd-*.md` — distribution copy（需對齊 .claude/agents/）
- **特別注意：** `mysd-lang.md`、`mysd-model.md` 需新增到 `plugin/commands/`；`mysd-designer.md` 需對齊

### 先前相關決策

- `.planning/phases/08-skill-md-orchestrators-agent-definitions/08-CONTEXT.md` — Plugin sync pattern, auto_mode 行為定義
- `.planning/phases/09-interactive-discovery-integration/09-CONTEXT.md` — discuss 自動接 spec/re-plan 的參考模式（propose auto-spec 的 UX 參考）

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets

- `mysd-discuss.md` Step 10-11 — **最佳參考模式**：discuss 已實作「呼叫 spec-writer → 更新 specs → re-plan 提示」串接，propose auto-spec Step 11 可沿用完全相同的 Task 呼叫模式
- `mysd-spec-writer` agent — 已存在，接受 `change_name + capability_area + existing_spec_body + proposal + auto_mode`
- `mysd-verifier` agent — 已存在，接受 `change_name + must_items + should_items + may_items`
- `mysd-apply.md` Step 4 後的狀態轉換 — auto-verify 在此步後插入
- `internal/config/defaults.go` — `ProjectConfig` struct，需新增 `DocsToUpdate []string`

### Established Patterns

- **SKILL.md auto_mode 傳播：** `--auto` flag 在 SKILL.md 層解析，`auto_mode: bool` 傳入 agent context — D-05/D-13 的跳過確認邏輯循此模式
- **Convention over config：** 未設定的配置欄位返回零值，不報錯（D-14 的 docs_to_update 未設定時靜默）
- **Plugin sync 雙目錄模式：** `.claude/` 是開發版本，`plugin/` 是 distribution copy，兩者保持同步（但只同步 mysd-*.md）
- **`--context-only` JSON 暴露模式：** binary 透過 `mysd execute --context-only` 輸出 JSON，SKILL.md 解析取值 — docs_to_update 跟隨此模式

### Integration Points

- `mysd-propose.md` Step 10 → 新增 Step 11（auto invoke mysd-spec-writer + 顯示 spec 摘要 + 後續指令清單）
- `mysd-apply.md` Step 4 → 新增 Step 5（auto verify: go build + go test + mysd-verifier）
- `mysd-archive.md` Step 1 後 → 新增 Step 2（read docs_to_update from binary JSON + update docs）
- `mysd-executor.md` Task Execution → 新增 on-failure sidecar 寫入段落（建議在 Step 3/TDD Step 後）
- `mysd-fix.md` Step 5B → 確認已有的 sidecar 讀取路徑與 D-06 格式對齊
- `internal/config/defaults.go` → 新增 `DocsToUpdate []string` 欄位

### Current State (codebase scout 2026-03-27)

- `mysd-propose.md` — 10 steps，D-01 的 auto-spec 尚未實作（Step 11 missing）
- `mysd-apply.md` — 4 steps，D-02 的 auto-verify 尚未實作（Step 5 missing）
- `mysd-archive.md` — 2 steps，D-10~D-14 的 docs_to_update 尚未實作
- `mysd-executor.md` — 227 lines，無 sidecar 寫入邏輯
- `mysd-fix.md` — 已有 "task sidecar" 讀取框架（Steps 3+4），但路徑待對齊 D-06
- `plugin/commands/` 缺 `mysd-lang.md`、`mysd-model.md`；`mysd-designer.md` 兩側有差異

</code_context>

<specifics>
## Specific Ideas

- sidecar 目錄路徑：`.specs/changes/{change_name}/.sidecar/T{id}-failure.md`，加入專案根 `.gitignore`
- docs_to_update 的 LLM 更新 context：proposal.md + tasks.md + specs/ + 現有文件內容（全讀）
- CHANGELOG.md 特殊處理：LLM 只生成新條目，SKILL.md 用 prepend 插入頂部，不修改其餘內容
- 確認 UX：「將更新以下文件：X, Y」→ Enter 確認 / n 跳過；--auto 直接執行
- propose auto-spec 的 Step 11 格式與 discuss Step 10 完全一致（同樣呼叫 mysd-spec-writer via Task tool）
- apply auto-verify 的 Step 5：先 `go build ./...`，若 build 失敗直接顯示 build error 並跳過 verifier agent；build 成功才呼叫 verifier
- executor sidecar on-failure 段落位置：建議在 Task Execution 的 "如果上述步驟失敗" 捕捉點插入，寫入後再 `mysd task-update {id} failed`

</specifics>

<deferred>
## Deferred Ideas

- **全新 `/mysd:doc` 獨立指令** — 使用者手動觸發文件更新；目前優先實作 archive 後的自動提示（D-11~D-14），獨立指令可作為未來 quick task
- **CLAUDE.md 架構說明自動化** — 為各子目錄自動生成架構說明，與此 phase 的 doc 維護流程不同，留待下一 milestone
- **UAT 深入設計** — UAT 流程目前為 advisory（不阻止 archive），此 phase 不修改 UAT 邏輯；深入設計留待後續 phase
- **GSD agents plugin sync** — `.claude/agents/` 的 gsd-*.md 不是 mysd 的 distribution 範圍，如有需要留待另一個 plugin sync task 討論

</deferred>

---

*Phase: 11-agent-doc*
*Context gathered: 2026-03-27 (updated after discussion)*
