# Phase 11: agent-doc — Context

**Gathered:** 2026-03-27
**Status:** Ready for planning

<domain>
## Phase Boundary

完善 mysd 工作流程的自動化串接（propose→spec、apply→verify）+ executor failure sidecar 寫入機制（讓 fix agent 有上下文可讀）+ archive 後的 doc 維護流程（可設定哪些文件需要更新）+ Phase 9-04 plugin sync 補足。

此 phase 不引入新 agent role，也不更動現有 spec/plan/archive 的核心邏輯。所有改動限於 SKILL.md 層和 agent prompt 層。

</domain>

<decisions>
## Implementation Decisions

### 工作流程自動化串接

- **D-01:** `propose` 完成後（Step 10 之後）自動呼叫 `mysd-spec-writer` agent，不需要使用者手動執行 `/mysd:spec`。spec 寫完後提示使用者「Next: `/mysd:plan`」。
- **D-02:** `apply` 所有 tasks 完成後（Step 4 之後）自動執行 verify — 先執行 `go build ./...` + `go test ./...`，再呼叫 `mysd-verifier` agent 做 MUST coverage check。
- **D-03:** `archive` 不增加額外 verify 步驟 — 維持現有邏輯（binary 已在 not-verified 狀態時 block archive）。apply 的自動 verify 已覆蓋覆蓋率檢查。
- **D-04:** `propose` 的自動 spec 呼叫可以被 `--skip-spec` flag 跳過（for cases where proposal is a skeleton only）。
- **D-05:** `apply` 的自動 verify 在 `--auto` 模式下不詢問確認，直接執行。

### Executor Failure Sidecar

- **D-06:** executor agent 在 task 執行失敗時（build error、test failure、中途退出）寫入 sidecar 文件：`.specs/changes/{change_name}/.sidecar/T{id}-failure.md`
- **D-07:** sidecar 包含：失敗時間戳、任務描述、build/test error output、AI 診斷嘗試（if any）
- **D-08:** `mysd-fix` Step 5B 讀取 T{id}-failure.md 作為初始 diagnosis context，若 sidecar 不存在則降級為無 context 診斷（backward compat）
- **D-09:** `.sidecar/` 目錄加入 `.gitignore`（failure context 是本地暫存，不應提交）

### Doc 維護流程

- **D-10:** `mysd.yaml` 新增 `docs_to_update` 欄位（字串陣列，指定需要在 archive 後更新的文件路徑，如 `["README.md", "CHANGELOG.md"]`）
- **D-11:** archive 完成後，SKILL.md 讀取 `docs_to_update` 配置，對每個指定文件呼叫 LLM 更新（讀取 change 的 proposal.md + specs + tasks 作為更新上下文）
- **D-12:** 文件更新在 SKILL.md 層實作，不需要新的 binary subcommand — 使用現有的 Bash/Read/Write/Edit 工具
- **D-13:** 更新前詢問使用者確認（顯示「這些文件將被更新：README.md」），`--auto` 模式跳過確認直接更新
- **D-14:** `docs_to_update` 未設定時，archive 後正常結束（不提示 doc 更新）— convention over config

### Plugin Sync 補足

- **D-15:** 完成 Phase 9-04 遺漏的 plugin sync：將 `.claude/commands/*.md` 同步到 `plugin/commands/`，`.claude/agents/*.md` 同步到 `plugin/agents/`
- **D-16:** 同步後執行 diff 確認兩側一致（`diff .claude/commands/X.md plugin/commands/X.md` 應為 zero output）

### Claude's Discretion

- Executor sidecar 的具體 Markdown 格式（結構自由，確保 fix agent 可讀即可）
- `docs_to_update` 中每個文件的更新策略（diff-based 增量 vs 全文重寫 — 選擇對 README 合適的方式）
- auto-verify 失敗時若 build 錯誤的 UX（提示用 `/mysd:fix` 還是直接顯示 build output）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### SKILL.md 修改目標

- `.claude/commands/mysd-propose.md` — Step 10 後需新增 auto-spec 呼叫（D-01/D-04）
- `.claude/commands/mysd-apply.md` — Step 4 後需新增 auto-verify 流程（D-02/D-05）
- `.claude/commands/mysd-archive.md` — Step 1 後需讀取 docs_to_update 並更新文件（D-11~D-14）
- `.claude/commands/mysd-ff.md` — 需同步加入 auto-verify 邏輯（D-02）
- `.claude/commands/mysd-ffe.md` — 需同步加入 auto-verify 邏輯（D-02）

### Agent 修改目標

- `plugin/agents/mysd-executor.md` — 新增 failure sidecar 寫入步驟（D-06/D-07）
- `.claude/agents/mysd-executor.md` — 同上（authoritative copy）
- `plugin/agents/mysd-fix.md` — Step 5B 新增 sidecar 讀取邏輯（D-08）（※ 注意：fix 是 SKILL.md，不在 plugin/agents/）
- `.claude/commands/mysd-fix.md` — Step 5B 新增 sidecar 讀取（D-08）

### 配置文件

- `.claude/mysd.yaml`（或 `.mysd.yaml`） — 新增 `docs_to_update` 欄位設計參考
- `internal/config/` — 現有 Viper 配置管理模式（新欄位在此讀取）

### Plugin Sync 來源

- `.claude/commands/` — 19 個 SKILL.md commands（authoritative dev copy）
- `.claude/agents/` — 12 個 agent definitions（authoritative dev copy）
- `plugin/commands/` — distribution copy（需與 .claude/commands/ 一致）
- `plugin/agents/` — distribution copy（需與 .claude/agents/ 一致）

### 先前相關決策

- `.planning/phases/08-skill-md-orchestrators-agent-definitions/08-CONTEXT.md` — Plugin sync pattern, auto_mode 行為定義
- `.planning/phases/09-interactive-discovery-integration/09-CONTEXT.md` — discuss 自動接 spec/re-plan 的參考模式

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets

- `mysd-discuss.md` Step 10-11 — **參考模式**：discuss 已實作「呼叫 spec-writer → 更新 specs → re-plan」串接，propose auto-spec 可沿用同樣的 Task 呼叫模式
- `mysd-spec-writer` agent — 已存在，接受 `change_name + capability_area + existing_spec_body + proposal + auto_mode`
- `mysd-verifier` agent — 已存在，接受 `change_name + must_items + should_items + may_items`
- `mysd apply` Step 4 後的 `mysd execute` 狀態轉換 — auto-verify 在此步後插入
- `internal/spec/deferred.go` — JSON 本地存取模式，sidecar 可參考同樣的讀寫慣例

### Established Patterns

- **SKILL.md auto_mode 傳播：** `--auto` flag 在 SKILL.md 層解析，`auto_mode: bool` 傳入 agent context — D-05/D-13 的跳過確認邏輯循此模式
- **Convention over config：** 未設定的配置欄位返回零值，不報錯（D-14 的 docs_to_update 未設定時靜默）
- **Plugin sync 雙目錄模式：** `.claude/` 是開發版本，`plugin/` 是 distribution copy，兩者保持同步

### Integration Points

- `mysd-propose.md` Step 10 → 新增 Step 11（auto invoke mysd-spec-writer）
- `mysd-apply.md` Step 4 → 新增 Step 5（auto verify: build + mysd-verifier）
- `mysd-archive.md` Step 1 後 → 新增 Step 2（read docs_to_update + update docs）
- `mysd-executor.md` Task Execution → 新增 on-failure sidecar 寫入
- `mysd-fix.md` Step 5B → 新增 sidecar 讀取
- `internal/config/` → 新增 `docs_to_update []string` 欄位 parse

</code_context>

<specifics>
## Specific Ideas

- sidecar 目錄路徑建議為 `.specs/changes/{change_name}/.sidecar/`，方便按 change 隔離，加入 `.gitignore`
- docs_to_update 的 LLM 更新上下文：讀取 proposal.md（what changed）+ tasks.md（what was done）+ 現有文件內容 → 生成更新後的文件
- propose auto-spec 的 Step 11 可以與 discuss Step 10 的格式完全一致（同樣呼叫 mysd-spec-writer）
- apply auto-verify 的 Step 5：先 go build，若 build 失敗直接顯示 build error 並跳過 verifier agent；build 成功才呼叫 verifier

</specifics>

<deferred>
## Deferred Ideas

- **全新 `/mysd:doc` 獨立指令** — 使用者手動觸發文件更新；目前優先實作 archive 後的自動提示（D-11~D-14），獨立指令可作為未來 quick task
- **CLAUDE.md 架構說明自動化** — 為各子目錄自動生成架構說明，與此 phase 的 doc 維護流程不同，留待下一 milestone
- **UAT 深入設計** — UAT 流程目前為 advisory（不阻止 archive），此 phase 不修改 UAT 邏輯；深入設計留待後續 phase

</deferred>

---

*Phase: 11-agent-doc*
*Context gathered: 2026-03-27*
