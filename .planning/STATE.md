---
gsd_state_version: 1.0
milestone: v1.1
milestone_name: Interactive Discovery & Parallel Execution
status: Milestone complete
stopped_at: Completed 12-03-PLAN.md
last_updated: "2026-03-27T04:53:13.206Z"
progress:
  total_phases: 8
  completed_phases: 8
  total_plans: 31
  completed_plans: 31
---

# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-03-25)

**Core value:** Spec 和執行的緊密整合 — 規格驅動 AI 執行，驗證回饋到規格，形成完整閉環
**Current focus:** Phase 12 — context

## Current Position

Phase: 12
Plan: Not started
Next: Phase 08 — SKILL.md Orchestrators & Agent Definitions

## Performance Metrics

**Velocity (v1.0 reference):**

- Total plans completed (v1.0): 18
- Average duration: ~8 min/plan
- Total execution time: ~2.4 hours

**By Phase (v1.0):**

| Phase | Plans | Avg/Plan |
|-------|-------|----------|
| 1. Foundation | 3 | 6 min |
| 2. Execution Engine | 6 | 10 min |
| 3. Verification & Feedback Loop | 5 | 6 min |
| 4. Plugin Layer & Distribution | 4 | 14 min |

**v1.1 metrics:** Not yet started
| Phase 05-schema-foundation-plan-checker P01 | 18 | 3 tasks | 9 files |
| Phase 05 P02 | 20 | 3 tasks | 5 files |
| Phase 06 P01 | 12 | 2 tasks | 4 files |
| Phase 06-executor-wave-grouping-worktree-engine P02 | 3 | 2 tasks | 5 files |
| Phase 06 P03 | 2 | 2 tasks | 2 files |
| Phase 06-executor-wave-grouping-worktree-engine P04 | 3 | 2 tasks | 2 files |
| Phase 07 P01 | 8 | 1 tasks | 2 files |
| Phase 07 P02 | 6 | 2 tasks | 6 files |
| Phase 07 P03 | 2 | 2 tasks | 2 files |
| Phase 07 P04 | 7 | 1 tasks | 2 files |
| Phase 07 P05 | 5 | 3 tasks | 5 files |
| Phase 08 P02 | 4 | 2 tasks | 5 files |
| Phase 08-skill-md-orchestrators-agent-definitions P01 | 5 | 1 tasks | 7 files |
| Phase 08 P04 | 2 | 1 tasks | 2 files |
| Phase 08 P03 | 12 | 2 tasks | 10 files |
| Phase 08-skill-md-orchestrators-agent-definitions P05 | 10 | 2 tasks | 6 files |
| Phase 08-skill-md-orchestrators-agent-definitions P05 | 20 | 3 tasks | 7 files |
| Phase 09-interactive-discovery-integration P01 | 4 | 2 tasks | 5 files |
| Phase 09-interactive-discovery-integration P02 | 14 | 2 tasks | 2 files |
| Phase 09-interactive-discovery-integration P03 | 12 | 2 tasks | 4 files |
| Phase 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync P01 | 297 | 2 tasks | 6 files |
| Phase 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync P02 | 6 | 2 tasks | 5 files |
| Phase 10-self-update-command-mysd-update-binary-version-check-plugin-file-sync P03 | 3 | 2 tasks | 4 files |
| Phase 11-agent-doc P01 | 3 | 2 tasks | 6 files |
| Phase 11-agent-doc P02 | 11 | 2 tasks | 2 files |
| Phase 11-agent-doc P03 | 2 | 2 tasks | 3 files |
| Phase 11 P05 | 8 | 2 tasks | 11 files |
| Phase 11-agent-doc P04 | 8 | 2 tasks | 3 files |
| Phase 12-context P02 | 10 | 2 tasks | 2 files |
| Phase 12-context P01 | 258 | 2 tasks | 8 files |
| Phase 12-context P03 | 12 | 2 tasks | 4 files |

## Accumulated Context

### Decisions

Recent decisions affecting v1.1 work:

- [Phase 04-03]: Plugin manifest uses minimal schema — no nested arrays in plugin.json
- [Phase 03-01]: VerificationStatus sidecar pattern (not modifying spec files) — discovery-state.json should follow same pattern
- [Phase 02-05]: SKILL.md orchestrator pattern: thin files + agent delegation via Task tool
- [Phase 02-05]: Alignment gate enforced by prompt ordering — same pattern applies to new agents
- [v1.1 roadmap]: plan-checker uses deterministic Go string matching on satisfies IDs (not AI inference)
- [v1.1 roadmap]: subagent cannot spawn subagent — only top-level SKILL.md may use Task tool; manual audit required before Phase 8 closes
- [v1.1 roadmap]: worktree paths kept short as T{id} only (no change name in path) for Windows MAX_PATH mitigation
- [Phase 05-01]: New fields appended at END of structs (D-11/D-12) for stable YAML field order — additive-only extension pattern
- [Phase 05-01]: Budget profile new roles (researcher/advisor/proposal-writer/plan-checker) use sonnet-4-5, not haiku — new subagent roles require quality model (D-06)
- [Phase 05-01]: ReadOpenSpecConfig returns zero-value (not error) for absent file — convention-over-config pattern from Phase 1
- [Phase 05]: CheckCoverage is pure function — all filesystem I/O stays in cmd layer, package has zero I/O dependencies
- [Phase 05]: mysd-plan-checker agent excludes Task tool (D-03) and Bash tool — leaf agent resolves gaps via Edit only
- [Phase 06]: ErrCyclicDependency returned (not silent skip) — BuildContextFromParts ignores cycle error, WaveGroups nil triggers sequential fallback in SKILL.md
- [Phase 06-02]: kernel32.dll via syscall.NewLazyDLL for Windows disk space (avoids golang.org/x/sys dependency)
- [Phase 06-02]: mysd worktree create/remove outputs JSON to stdout for SKILL.md consumption (consistent with --context-only pattern)
- [Phase 06-03]: execute --context-only required zero cmd changes — Plan 01 already wired WaveGroups via BuildContextFromParts
- [Phase 06-03]: plan.go wave groups computed only when tasks.md exists (lazy load) — graceful nil on parse failure
- [Phase 06-04]: Mode selection per D-03: has_parallel_opportunity false skips prompt; auto_mode true uses wave without asking
- [Phase 06-04]: Merge loop in ascending task ID order with --no-ff; 3-retry AI conflict resolution with go build+test; continue-on-failure policy; failed worktrees preserved
- [Phase 07]: ScanContext 完全替換為語言無關通用 struct（primary_language/files/modules）— binary 只收集 metadata，LLM 處理任意語言（GSD 同樣模式）
- [Phase 07]: `mysd init` 內部直接展開為 `scan --scaffold-only`（無 warning，回展相容）
- [Phase 07]: Skills 推薦邏輯在 mysd-planner agent 層，確認流程在 SKILL.md 層，預設 accept-all
- [Phase 07]: `mysd model` table 輸出（lipgloss），`mysd model set` 直接寫 .claude/mysd.yaml
- [Phase 07]: FilterBlockedTasks uses BFS over adjacency map (same pattern as BuildWaveGroups) — no new data structures needed
- [Phase 07]: ScanContext replaced PackageInfo — language-agnostic struct with PrimaryLanguage/Files/Modules (D-02 no backward compat)
- [Phase 07]: scaffoldOpenSpecDir idempotent via os.MkdirAll — init removed --force flag, delegates to scaffold
- [Phase 07]: Plain text fmt.Fprintf for model table (not lipgloss) — satisfies D-11 without TTY dependency
- [Phase 07]: lang set uses defer rollback pattern (write mysd.yaml first, rollback if openspec write fails) — safer on Windows than write-then-rename
- [Phase 07]: Skills recommendation in planner agent layer, confirmation in SKILL.md layer, default accept-all (D-07 through D-10)
- [Phase 08-02]: mysd-executor: assigned_task is now the ONLY task input — no pending_tasks list, no execution_mode field; SKILL.md orchestrator handles the loop
- [Phase 08-02]: mysd-spec-writer: capability_area + auto_mode added; Discuss step and state transition removed — SKILL.md orchestrator responsibilities
- [Phase 08-01]: All 4 new/synced agent definitions have zero Task tool references in allowed-tools — enforces leaf agent constraint (D-17)
- [Phase 08-01]: Plugin sync pattern: .claude/agents/ is authoritative dev copy, plugin/agents/ is distribution copy with identical content
- [Phase 08]: auto_mode in discuss skips research entirely (FAUTO-02: ff/ffe-style auto = no interaction) — propagated to all spawned agents
- [Phase 08]: discuss source detection (D-06): change-name > file-path > dir-path > active change > auto-detect (gstack/context, not .claude/plans/) > create new
- [Phase 08]: /mysd:plan redesigned as 3-stage pipeline: researcher(x4 parallel) -> designer -> planner; execute renamed to apply at SKILL.md layer
- [Phase 08]: /mysd:apply spawns mysd-executor per task; single=sequential, wave=parallel within wave_groups; --auto flag parsed at SKILL.md layer and propagated as auto_mode
- [Phase 08-05]: fix uses safety valve: auto-detects path (conflict markers vs sidecar failure) but confirms with user before proceeding (D-08)
- [Phase 08-05]: ff/ffe do not use mysd-fast-forward agent — directly orchestrate designer+planner+executor pipeline with auto_mode hardcoded true (D-24/D-25/FAUTO-03)
- [Phase 08-05]: fix uses safety valve: auto-detects path (conflict markers vs sidecar failure) but confirms with user before proceeding (D-08)
- [Phase 08-05]: ff/ffe do not use mysd-fast-forward agent — directly orchestrate designer+planner+executor pipeline with auto_mode hardcoded true (D-24/D-25/FAUTO-03)
- [Phase 09-01]: DeferredStore stored as deferred.json in specDir root (not tied to active change) — scope-free, notes persist across change lifecycle
- [Phase 09-01]: noteCmd.RunE = runNoteList so mysd note without subcommand lists notes; status shows deferred count only when count > 0 (D-09)
- [Phase 09-02]: propose always loads deferred notes (D-02): cross-change context valuable for new proposals
- [Phase 09-02]: discuss conditionally loads deferred notes: active WIP = skip notes to avoid polluting focused WIP (D-02)
- [Phase 09-02]: dual-loop uses user-driven termination not numeric quota (D-01): binary choice per area sufficient termination signal
- [Phase 09]: D-04 fix: plan stage uses single mysd-researcher with architecture dimension (not 4 parallel researchers) — requirements already finalized at plan stage, only technical validation needed
- [Phase 09]: mysd-note SKILL.md is thin wrapper (Bash+Read only, no Task tool) — orchestrator pattern reserved for multi-agent flows; deferred notes count silent when zero
- [Phase 10]: CheckLatestVersionWithBase added as testable variant to enable httptest mocking without changing public API
- [Phase 10]: ApplyUpdate calls replaceExecutable internally; Rollback is public for cmd layer to call on post-update failure
- [Phase 10]: LoadManifest returns (nil, nil) for missing file — matches deferred.go convention-over-config pattern, represents pre-v1.1 installation without manifest
- [Phase 10]: DiffManifests with nil old manifest: all new files are add, zero deletes — backward compat per D-17 for pre-v1.1 installations
- [Phase 10]: SyncPlugins delete errors are non-fatal: appended to Errors slice, sync continues
- [Phase 10]: findClaudeDir walks up from cwd to filesystem root to locate .claude/ — supports running from any subdirectory
- [Phase 10]: Binary update only runs when --force AND update_available — SKILL.md layer handles user confirmation flow
- [Phase 11-agent-doc]: DocsToUpdate defaults to nil (not empty slice) — convention over config per D-14; omitempty ensures JSON omits the field when unconfigured
- [Phase 11-agent-doc]: mysd docs add/remove use viper read-modify-write to preserve other config fields — same pattern as runModelSet in cmd/model.go
- [Phase 11-agent-doc]: propose Step 11 auto-chains to mysd-spec-writer after proposal; --skip-spec flag bypasses
- [Phase 11-agent-doc]: apply Step 5 auto-chains to verifier after go build+test pass; auto_mode skips confirmation (D-05)
- [Phase 11-agent-doc]: On Failure path is alternative exit in executor — MUST NOT proceed to Mark Task Done or Atomic Commit after sidecar write (D-06, D-07)
- [Phase 11-agent-doc]: failure_context null fallback in fix agent — backward compat for pre-D-06 task states without sidecars (D-08)
- [Phase 11]: mysd-docs SKILL.md follows thin wrapper pattern (Bash+Read only) — consistent with mysd-note.md convention
- [Phase 11]: Plugin sync zero diff policy: plugin/ distribution must be byte-identical to .claude/ dev copies; mysd-lang.md and mysd-model.md gaps from Phase 7/9 closed
- [Phase 11-agent-doc]: archive Step 0 reads docs_to_update before archive runs — enables confirmation flow before irreversible action
- [Phase 11-agent-doc]: ff/ffe inline docs update always uses auto_mode=true (no confirmation) — consistent with ff/ffe being fully automatic pipelines
- [Phase 12-02]: Bridge file written only when gsd-context-monitor.js detected (D-04) — avoids /tmp pollution in non-GSD projects
- [Phase 12-02]: statusline_enabled=false suppresses output but bridge file still writes (D-12) — GSD context monitor must not lose data
- [Phase 12-01]: Go embed cannot use ../ path; workaround: copy JS to cmd/hooks/ subdirectory for embed directive
- [Phase 12-01]: runStatuslineInDir extracted for testability; MkdirAll before SafeWriteConfig required on fresh directories
- [Phase 12-context]: deleteResearchCache extracted as testable helper — enables direct unit testing without full runArchive scaffolding
- [Phase 12-context]: Interstitial step numbering (4.5, 6.5) for SKILL.md extension without renumbering existing steps

### Quick Tasks Completed

| # | Description | Date | Commit | Directory |
|---|-------------|------|--------|-----------|
| 260326-n23 | 全部現有 mysd 指令補上 argument-hint | 2026-03-26 | e391ab0 | [260326-n23-mysd-argument-hint](./quick/260326-n23-mysd-argument-hint/) |
| 260326-p0d | README.md v1.1 更新 — 涵蓋 self-update、wave execution、interactive discovery、6 個新指令 | 2026-03-26 | 626f275 | [260326-p0d-readme-md-v1-1-self-update-interactive-d](./quick/260326-p0d-readme-md-v1-1-self-update-interactive-d/) |

### Roadmap Evolution

- Phase 10 added: Self-Update Command — /mysd:update binary version check + plugin file sync
- Phase 11 added: 增強 agent 功能及增加 doc 維護流程
- Phase 12 added: 加入 context 的 % 數及色條

### Pending Todos

- [Phase 8] `/mysd:propose` 自動偵測輸入來源（參考 spectra:propose Step 1 設計）：
  - Case 1：偵測 `.planning/phases/{phase}/` 下的 CONTEXT.md / PLAN.md 作為 initial content
  - Case 2：從當前 conversation context 提取需求
  - 優先順序：argument > planning files > conversation context > 詢問使用者

- [Phase 8] `/mysd:status` SKILL.md 指令設計：
  - 顯示當前 workflow stage（propose → spec → plan → execute → verify），標示目前位置
  - Task 列表含編號、title、狀態符號（✓ done / ✗ failed / ⊘ skipped / ○ pending）
  - 最後一行：`Next: /mysd:{command}` 推薦下一步指令

- [Phase 8] `/mysd:fix` 流程設計決策（discuss-phase 8 時須納入）：
  - **Fix 兩條路徑**：(1) merge 衝突 → 解衝突 → merge → `-D` 強制刪 branch；(2) 實作問題 → 自動修正 task 內容 + 更新 spec → `-D` 刪舊 branch → 交 executor 子代理重新執行
  - **放棄路徑**：fix 無解時使用者可選擇放棄 → `-D` 強制刪 branch + worktree → task 回 `pending`
  - **失敗 context 保存**：task sidecar 記錄失敗原因、AI 嘗試解法、放棄理由，重新討論時直接讀取
  - **Execute 完成提示**：執行結束後自動列出需 fix 的 tasks（`T2 (setup-auth) — merge 衝突，執行 /mysd:fix T2`）
  - **Re-run 衝突處理**：偵測到 worktree 仍存在時 block 並指引使用者執行 fix（不允許直接重跑）
  - **Branch cleanup**：`Remove()` 保留 `git branch -d`（soft delete），但刪除失敗時提示使用者（已有 stderr warning，需改善訊息清晰度）
  - **Skipped tasks 恢復**：fix 成功 merge 後，所有因依賴此 task 而被 skipped 的下游 tasks 自動恢復為 `pending`，可排入下次執行（遞移性恢復）

### Blockers/Concerns

- Phase 6: Windows worktree MAX_PATH needs CI validation — `git config core.longpaths true` mitigation must be empirically verified
- Phase 9: Interactive Discovery dual-loop requires focused design review of termination conditions before writing agent prompts
- Phase 5: Verify whether `golang.org/x/term` IsTerminal is already imported in v1.0 binary (needed for Phase 7 interactive commands)
- Phase 8: All 9 agent definitions require manual audit for Task tool references before Phase 8 can close

## Session Continuity

Last session: 2026-03-27T04:48:30.261Z
Stopped at: Completed 12-03-PLAN.md
Resume file: None
