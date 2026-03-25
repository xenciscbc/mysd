# Project Research Summary

**Project:** mysd v1.1 — Interactive Discovery, Parallel Execution, Subagent Architecture, Language-Agnostic Scan
**Domain:** Go CLI tool — subsequent milestone adding orchestration intelligence to an existing spec-driven development system
**Researched:** 2026-03-25
**Confidence:** HIGH

## Executive Summary

mysd v1.1 is a targeted capability expansion on a shipped v1.0 system (11 packages, 7,555 LOC). The milestone closes three major gaps left in v1.0: no interactive discovery before proposal writing, no dependency-safe parallel task execution, and no path for mid-project spec updates. The existing architecture — binary as context-provider, SKILL.md as orchestrator, agents as workers — is sound and must be preserved without alteration. All new features either add Go packages (`internal/worktree/`, `internal/discovery/`, `internal/planchecker/`), extend existing packages (executor, spec schema, config, scanner), or add new SKILL.md + agent definitions. No new external dependencies are required; stdlib additions (`os/exec`, `filepath`, `syscall.Statfs`) are sufficient for all new functionality.

The highest-risk area is worktree-based parallel execution: it introduces a hard prerequisite (task dependency + file overlap detection via topological sort) before worktrees can be created safely, and carries two platform hazards — Windows MAX_PATH truncation and disk space exhaustion from full working-directory copies. Both have documented mitigations but require verification in CI before the phase is closed. The second risk area is Interactive Discovery's dual-loop conversation model, which requires explicit hard termination conditions (max 3 areas, max 3 questions per area) to prevent user-trapping spirals. A critical structural constraint affects all agent design: Claude Code subagents cannot spawn other subagents — only the top-level SKILL.md orchestrator may use the Task tool. This constraint shapes every agent definition in the milestone.

The recommended build order runs schema changes first (enabling safe extension of existing packages), then executor wave grouping, then worktree isolation, then new binary commands, then SKILL.md orchestrators and agent definitions, then Interactive Discovery integration. This ordering ensures each layer builds on verified foundations. Particularly important: `/mysd:fix` reuses worktree machinery (cannot be built before Phase 2) and Interactive Discovery requires all 4 new agent definitions (cannot be built before Phase 4).

---

## Key Findings

### Recommended Stack

No new external dependencies are required for v1.1. The existing stack (Go 1.23+, Cobra v1.10.2, Viper v1.x, lipgloss v1.x, yaml.v3, adrg/frontmatter v0.2.0, testify v1.x, GoReleaser v2.14+) is sufficient for all new features. New capabilities use Go stdlib exclusively.

**Core technologies (all existing — no changes):**
- Go 1.23+: single binary deployment, cross-platform — required by PROJECT.md constraint
- github.com/spf13/cobra v1.10.2: CLI framework for all 16+ subcommands (14 existing + model, lang)
- gopkg.in/yaml.v3: YAML parsing — required for tasks.md `depends`/`files` schema extension
- github.com/adrg/frontmatter v0.2.0: spec file frontmatter parsing — no changes needed
- github.com/charmbracelet/lipgloss v1.x: wave execution progress table rendering
- github.com/stretchr/testify v1.x: assertions for new packages (worktree, planchecker, discovery)

**One stdlib addition to verify:**
- `golang.org/x/term.IsTerminal` — must be checked before any interactive prompt to prevent stdin hangs in non-TTY contexts (CI, pipes). Confirm this is already imported; add to `cmd/lang.go` and any other commands that add interactive prompts.

**Note on STACK.md:** No STACK.md was produced for this milestone — the existing v1.0 stack is confirmed sufficient. All findings above are derived from direct codebase analysis.

### Expected Features

Research (based on authoritative `proposal.md` + direct codebase analysis) produces a clear priority split. All P1 items are required for v1.1 milestone closure.

**Must have — P1 (v1.1 milestone blockers):**
- Task dependency + file overlap detection (`depends[]`, `files[]` in TaskItem; topological sort + file-overlap wave layering) — prerequisite for all parallel execution
- Worktree parallel execution with AI conflict resolution, 3-retry policy, partial failure handling
- Interactive Discovery dual-mode (research + general) in propose/spec/discuss
- `/mysd:discuss` with auto re-plan + plan-checker trigger
- Plan-checker (`mysd-plan-checker` agent, auto-triggered after every plan, ID-based `satisfies` field matching)
- New subagent definitions: `mysd-researcher`, `mysd-advisor`, `mysd-proposal-writer`, `mysd-plan-checker`
- `/mysd:fix` with worktree isolation
- Scan upgrade to language-agnostic (top-3 languages by file count + LOC, directory mapping, openspec/config.yaml writer)
- `/mysd:model` CLI (show/set/resolve — surfaces existing config capability)
- `/mysd:lang` CLI (interactive locale + atomic two-file sync)
- Model profile extended for 4 new agent roles (researcher, advisor, proposal-writer, plan-checker)

**Should have — P2 (add when P1 is stable):**
- `--auto` flag polish across propose/spec/discuss
- Codebase Scout refinement (smarter grep patterns per detected language)

**Defer — v2+:**
- Multi-language model profile differentiation
- Worktree progress streaming to terminal (bubbletea TUI explicitly out of scope)
- Support for AI tools beyond Claude Code

**Anti-features to avoid:**
- MCP server for subagent coordination (breaks single-binary deployment)
- Interactive discovery in design stage (creates spec/design feedback loops)
- Auto-push worktree branches to remote (mysd is a single-developer tool)
- Auto-rewrite spec intent from code output (inverts spec-as-source-of-truth)
- Real-time streaming output from parallel tasks (interleaved output from N worktrees is unreadable)

### Architecture Approach

v1.1 extends the established binary-as-context-provider pattern without modifying it. Three new Go packages are added — each with a single, narrow responsibility and no circular imports. Four existing packages receive additive-only schema changes. The core principle is preserved: the Go binary manages state and produces structured JSON; SKILL.md files are orchestrators; agent `.md` files are workers. The binary never prompts stdin interactively.

**Major components added in v1.1:**
1. `internal/worktree/` — git worktree lifecycle: create at `.worktrees/T{id}/`, branch naming `mysd/{change}/T{id}-{slug}`, merge in ascending task ID order with `--no-ff`, cleanup-on-success/preserve-on-failure; Windows MAX_PATH mitigation; disk space pre-flight check
2. `internal/executor/` (modified) — `BuildWaveGroups`: topological sort by `Depends` + file-overlap check produces `[][]TaskItem`; extended `ExecutionContext` with `WaveGroups`, `WorktreeDir`, `AutoMode`
3. `internal/planchecker/` — pure function `CheckCoverage(tasks, mustItems) CoverageResult`; deterministic Go string matching on `satisfies` field IDs, not AI semantic inference
4. `internal/discovery/` — `DiscoveryContext` and `DiscoveryState`; persisted as `discovery-state.json` sidecar (same pattern as existing `verification-status.json`); research summary cached for ff/ffe reuse
5. `internal/scanner/` (refactored) — language-agnostic: top-3 languages by file count + LOC with directory mapping; module file detection (`go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`); shared `openspec/config.yaml` writer used by both scan and lang commands
6. 4 new plugin agents + 4 new SKILL.md orchestrators — all spawned by the top-level orchestrator only; nested agent spawning is structurally forbidden

**State machine:** No changes required. Existing `ValidTransitions` already handles all v1.1 re-entry flows (specced → designed → planned). `/mysd:discuss` is a re-entry point, not a new phase state.

### Critical Pitfalls

1. **Subagent cannot spawn subagents** (Phase 1/4) — `mysd-researcher` must return a gray area list; the SKILL.md orchestrator spawns `mysd-advisor x N` in parallel. No subagent definition may reference the Task tool. Manual audit of all 9 agent definitions required before any end-to-end test.

2. **Interactive discovery loop has no exit condition** (Phase 1/5) — Hard limits required: max 3 gray areas, max 3 depth questions per area. After each area, present binary "proceed or continue?" with visible countdown. `--auto` must skip all loops entirely and use AI's first recommendation.

3. **Worktree disk explosion on large codebases** (Phase 2) — Each worktree is a full working-directory copy. Pre-flight disk check required (codebase size × task count × 1.5). Default parallel cap: 4 tasks. `GOCACHE` redirected per worktree to a path outside `.worktrees/`, cleaned on removal.

4. **Windows MAX_PATH silently truncates worktree paths** (Phase 2) — Run `git config core.longpaths true` as the first git operation in every new worktree when `GOOS == "windows"`. Keep worktree root names as `T{id}` only — no change name in the path. Verify with Windows CI runner before closing Phase 2.

5. **Plan-checker false negatives from fuzzy MUST matching** (Phase 1/2) — Add `satisfies: [REQ-001, REQ-003]` field to `TaskEntry`. Plan-checker uses deterministic Go string matching on requirement IDs. AI only suggests `satisfies` values when planner writes tasks; it never performs coverage verification by inference.

6. **Orphaned worktrees after interrupted execution** (Phase 2) — `defer` cleanup does not run on SIGKILL. Write state marker before worktree creation. Run `git worktree prune` + orphan scan at every `mysd execute` startup. Never silently overwrite an existing worktree path.

7. **Locale config desync between mysd.yaml and openspec/config.yaml** (Phase 3) — Use atomic write-both-or-neither: prepare both new contents in memory, write to temp files, rename atomically. If either rename fails, surface the error explicitly. Do not rely on sequential writes.

---

## Implications for Roadmap

Based on combined research, the 8-phase build order in ARCHITECTURE.md maps to a 5-phase roadmap grouped by logical delivery unit. The ordering is driven by hard dependencies: schema before executor extension, executor before worktree, worktree before `/mysd:fix`, agents before Interactive Discovery integration.

### Phase 1: Schema Foundation + Plan-Checker Infrastructure

**Rationale:** All downstream phases depend on the extended TaskItem schema (`Depends`, `Files`, `Satisfies`) and the new agent model mappings. `internal/planchecker/` is a pure function with no external dependencies — the safest first package to build and easiest to prove correct. Getting the schema right here avoids costly migrations later.

**Delivers:** Extended `TaskEntry` with `Depends []int`, `Files []string`, `Satisfies []string`; extended `ExecutionContext` with `WaveGroups`, `WorktreeDir`, `AutoMode`; updated `DefaultModelMap` with 4 new agent roles; `WorktreeDir`/`AutoMode` in `ProjectConfig`; fully tested `planchecker.CheckCoverage` function; updated spec `WriteTasks` serialization.

**Addresses:** Task dependency detection (prerequisite), plan-checker infrastructure (P1), model profile extension
**Avoids:** Plan-checker false negatives (ID-based `satisfies` field designed before the checker is built)
**Research flag:** Standard Go struct extension with yaml.v3 — no deeper research needed.

---

### Phase 2: Executor Wave Grouping + Worktree Execution Engine

**Rationale:** `BuildWaveGroups` must be proven correct before worktrees are created around it. Worktree lifecycle is the highest-risk component; building it early (while other dependencies are minimal) surfaces failures when they are cheapest to fix.

**Delivers:** `internal/executor.BuildWaveGroups` (topological sort + file-overlap); extended `ExecutionContext` with wave groups output via `--context-only`; `internal/worktree/` package (create/merge/cleanup); `mysd-executor` agent updated with worktree isolation instructions; pre-flight disk space check; Windows MAX_PATH mitigation (`git config core.longpaths true`); startup orphan scan + state marker for crash recovery.

**Addresses:** Worktree parallel execution (P1), partial failure policy, 3-retry AI conflict resolution
**Avoids:** Disk explosion (pre-flight disk check), Windows MAX_PATH (longpaths config), orphaned worktrees (startup orphan scan), wave divergence (complete-wave-before-next-wave boundary contract)
**Research flag:** Git worktree on Windows needs CI validation — add Windows runner before closing this phase. The `core.longpaths` mitigation is documented but needs empirical verification against the actual project path length.

---

### Phase 3: New Binary Commands + Scanner Refactor

**Rationale:** These are independent of worktree machinery (can proceed in parallel with Phase 2 if capacity allows). `openspec/config.yaml` writer is shared by scan and lang — implement once and reuse. Low-complexity items (`/mysd:model`, `/mysd:lang`) deliver immediate user value while Phase 2 stabilizes.

**Delivers:** `cmd/model.go` (profile show/set/resolve); `cmd/lang.go` (locale interactive set with atomic two-file sync); refactored `internal/scanner/` (language-agnostic, top-3 languages with directory mapping, module file detection); `openspec/config.yaml` writer (shared); `cmd/plan.go` extended with plan-checker context output; locale BCP 47 validation before writing.

**Addresses:** `/mysd:model` CLI (P2), `/mysd:lang` CLI (P2), scan upgrade (P1)
**Avoids:** Locale config desync (atomic rename pattern), language detection misclassification (multi-language fixture tests required), non-TTY stdin hang (`term.IsTerminal` check before any interactive prompt)
**Research flag:** Standard Cobra/Viper patterns — no deeper research needed.

---

### Phase 4: New SKILL.md Orchestrators + Agent Definitions

**Rationale:** SKILL.md files can only be written after the binary commands they invoke are stable (Phases 1-3 complete). All new agent definitions are prerequisite for Interactive Discovery integration (Phase 5). Grouping all new SKILL.md + agents in one phase ensures orchestration patterns are consistent.

**Delivers:** `mysd-discuss.md`, `mysd-fix.md`, `mysd-model.md`, `mysd-lang.md` SKILL.md files; `mysd-researcher.md`, `mysd-advisor.md`, `mysd-proposal-writer.md`, `mysd-plan-checker.md` agent definitions; plan-checker auto-trigger wired into `/mysd:plan` orchestrator; interactive loop termination structure (max 3 areas, max 3 questions per area) established as SKILL.md frontmatter conventions.

**Addresses:** `/mysd:discuss` (P1), `/mysd:fix` (P1), all new agent definitions (P1), plan-checker integration
**Avoids:** Nested subagent spawning (every agent definition explicitly forbids Task tool usage; manual audit of all 9 definitions before phase closes), subagent context overload (paths-not-content rule enforced; Task prompts capped at 300 words)
**Research flag:** Claude Code subagent constraints are HIGH confidence (official docs + issue tracker). No deeper research needed. Manual audit of all 9 agent SKILL.md files for Task tool references is required verification gate before closing this phase.

---

### Phase 5: Interactive Discovery Integration

**Rationale:** The most complex orchestration work in the milestone. Depends on all 4 new agent definitions (Phase 4) and the `internal/discovery/` package. Saved for last so it integrates only proven components.

**Delivers:** `internal/discovery/` package (`DiscoveryContext`, `DiscoveryState`, sidecar JSON persistence with `research_summary` for ff/ffe reuse); modified `mysd-propose.md` and `mysd-spec.md` with research mode flow and hard question limits; modified `mysd-ff.md` + `mysd-ffe.md` with research-once pattern and auto mode; Codebase Scout orchestrator-level grep/glob steps; scope guardrail (deferred notes redirect for scope expansion).

**Addresses:** Interactive Discovery dual-mode (P1 — primary user-facing value of v1.1), `--auto` flag (P2), shared research result for ff/ffe (P2), Codebase Scout (P2), scope guardrail
**Avoids:** Discovery loop with no exit condition (hard max 3 areas/3 questions in SKILL.md; binary "proceed?" prompt with visible countdown), research mode on every run (opt-in for existing changes, default off; only auto-triggered for new proposals)
**Research flag:** The dual-loop conversation model with scope guardrails is non-trivial SKILL.md orchestration. A focused design review of loop termination conditions, deferred notes format, and research summary schema is recommended before writing discovery-related prompts.

---

### Phase Ordering Rationale

- Schema changes first because every other package reads/writes `TaskEntry` and `ExecutionContext`
- Worktree in Phase 2 (not later) because it is the highest-risk component — surface failures early when changes are easy
- Binary commands in Phase 3 because SKILL.md files must invoke working binaries; writing wrappers before the binary is stable produces inaccurate documentation that misleads the AI
- SKILL.md/agents in Phase 4 because Interactive Discovery integration requires proven, audited agent definitions
- Discovery last because it is the most complex orchestration work and depends on everything else
- Platform hazard prevention is mandatory within its phase (disk guard and Windows mitigation in Phase 2; locale atomicity in Phase 3) — never deferred to a follow-up phase

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 2 (Worktree on Windows):** Pre-implementation CI validation is needed to confirm `git config core.longpaths true` is sufficient given the actual project root path length of the target user base, or whether additional mitigations are required.
- **Phase 5 (Interactive Discovery):** The dual-loop conversation model with scope guardrails is novel for this codebase. A focused design session on loop termination conditions, deferred notes format, and auto-mode behavior is recommended before writing agent prompts.

Phases with standard patterns (skip research-phase):
- **Phase 1 (Schema Foundation):** Pure Go struct extensions with yaml.v3 — well-documented, established patterns.
- **Phase 3 (Binary Commands):** Cobra/Viper command patterns are established in 14 existing commands; scanner refactor follows standard Go refactoring patterns.
- **Phase 4 (Agents/SKILL.md):** Claude Code agent format is HIGH confidence (official docs); patterns established by 8 existing agents; constraint (no nested spawning) is definitively documented.

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | No new dependencies; existing stack verified in production (v1.0 shipped). STACK.md not produced — not needed. Stdlib additions are standard Go patterns. |
| Features | HIGH | Derived from authoritative `proposal.md` (primary source) + direct codebase analysis of v1.0. P1/P2/P3 split is unambiguous. Anti-features documented with concrete rationale. |
| Architecture | HIGH | Direct analysis of v1.0 source code for all integration points. Package boundaries are explicit. Build order derived from hard dependency analysis, not opinion. All new type signatures specified in ARCHITECTURE.md. |
| Pitfalls | HIGH (critical), MEDIUM (integration gotchas) | Critical pitfalls verified via official Anthropic docs + Claude Code issue tracker (GitHub issues #4182, #4277). Windows MAX_PATH via Microsoft Learn. Integration gotchas from multiple practitioner sources. Performance traps from Git worktree engineering posts. |

**Overall confidence:** HIGH

### Gaps to Address

- **`satisfies` field in TaskEntry:** PITFALLS.md requires this for deterministic plan-checker matching. ARCHITECTURE.md's modified `TaskEntry` schema does not include it. Must be added to Phase 1 schema work. Treat `satisfies []string` as a required addition alongside `depends` and `files`.

- **Advisor parallelism cap:** PITFALLS.md recommends capping at 4 advisors maximum. ARCHITECTURE.md mentions `config.AgentCount` but this field does not currently exist in `ProjectConfig`. Decide during Phase 1/4 planning whether to add `AgentCount` to config schema or hard-code the cap in SKILL.md.

- **`golang.org/x/term` import status:** Required for non-TTY detection before interactive prompts. Unclear whether already imported in v1.0. Verify during Phase 3 when `cmd/lang.go` is written.

- **`discovery-state.json` staleness:** ARCHITECTURE.md notes research summaries grow stale across days and defers a `--force-research` flag. Roadmapper must decide whether v1.1 ships with a TTL check or treats this as a v1.1.x follow-up. If deferred, document the known limitation in the Phase 5 spec.

---

## Sources

### Primary (HIGH confidence)
- `.specs/changes/interactive-discovery/proposal.md` — authoritative v1.1 feature list and design decisions
- `.planning/PROJECT.md` — milestone scope, constraints, out-of-scope items
- `internal/executor/context.go` — ExecutionContext schema (direct analysis)
- `internal/spec/schema.go` — TaskEntry, TasksFrontmatterV2 (direct analysis)
- `internal/state/transitions.go` — ValidTransitions map (direct analysis)
- `internal/config/defaults.go`, `internal/config/config.go` — ProjectConfig, DefaultModelMap (direct analysis)
- `plugin/commands/mysd-execute.md` — SKILL.md orchestration pattern (direct analysis)
- `plugin/agents/mysd-executor.md` — subagent input contract (direct analysis)
- Claude Code Sub-Agents official docs — https://code.claude.com/docs/en/sub-agents
- GitHub Issue #4182: Sub-Agent Task Tool Not Exposed — https://github.com/anthropics/claude-code/issues/4182
- GitHub Issue #4277: Claude Code agentic loop detection — https://github.com/anthropics/claude-code/issues/4277
- Windows MAX_PATH Limitation — https://learn.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation
- golang.org/x/term IsTerminal — https://pkg.go.dev/golang.org/x/term

### Secondary (MEDIUM confidence)
- Git Worktrees for Parallel AI Agents (Upsun) — https://devcenter.upsun.com/posts/git-worktrees-for-parallel-ai-coding-agents/
- Git worktrees for parallel development (nrmitchi) — https://www.nrmitchi.com/2025/10/using-git-worktrees-for-multi-feature-development-with-ai-agents/
- Context Management with Subagents in Claude Code — https://www.richsnapp.com/article/2025/10-05-context-management-with-subagents-in-claude-code
- Claude Code Subagents: Common Mistakes — https://claudekit.cc/blog/vc-04-subagents-from-basic-to-deep-dive-i-misunderstood
- Token Cost Trap in AI agents — https://medium.com/@klaushofenbitzer/token-cost-trap-why-your-ai-agents-roi-breaks-at-scale
- Interactive CLI prompts in Go — https://dev.to/tidalcloud/interactive-cli-prompts-in-go-3bj9
- Go i18n with golang.org/x/text/language — https://phrase.com/blog/posts/internationalization-i18n-go/

### Tertiary (LOW confidence)
- When agents learn to ask: Active questioning in agentic AI — https://medium.com/@milesk_33/when-agents-learn-to-ask — single source; use as inspiration only for discovery loop design

---
*Research completed: 2026-03-25*
*Ready for roadmap: yes*
