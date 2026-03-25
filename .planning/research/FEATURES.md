# Feature Research

**Domain:** AI-powered CLI tool for Spec-Driven Development — v1.1 milestone: interactive discovery, parallel execution, subagent orchestration, language-agnostic scanning
**Researched:** 2026-03-25
**Confidence:** HIGH (proposal.md is authoritative; existing codebase verified directly)

---

## Context: What Already Exists (v1.0 — Do Not Rebuild)

The following are already shipped and must be treated as stable foundations:

| Existing Feature | Location | Notes |
|-----------------|----------|-------|
| propose → spec → design → plan → execute → verify → archive | `cmd/` | Full pipeline with phase gates |
| OpenSpec-compatible parser (brownfield) | `internal/spec/` | Delta specs, RFC 2119 keywords |
| Single/wave execution modes | `cmd/execute.go`, `internal/executor/` | Wave is basic — no `depends`/`files` handling |
| Model profile system (quality/balanced/budget) | `internal/config/config.go` | `ResolveModel` + `DefaultModelMap`; config only, no CLI |
| ff/ffe fast-forward commands | `cmd/ff.go`, `cmd/ffe.go` | Already imply auto behavior |
| /mysd:scan (Go packages only) | `internal/scanner/scanner.go` | `PackageInfo` is Go-specific (`GoFiles`, `TestFiles`) |
| status dashboard, capture, init | `cmd/` | status is read-only; init creates scaffold |
| 14 SKILL.md commands + 8 agent definitions | `plugin/` | mysd-executor, mysd-planner, etc. |

---

## Feature Landscape

### Table Stakes (Users Expect These)

Non-negotiable for v1.1. Missing these = milestone feels incomplete.

| Feature | Why Expected | Complexity | Depends On |
|---------|--------------|------------|------------|
| `/mysd:discuss` command | Users need ad-hoc spec updates outside the main workflow; gap is visible in v1.0 — no path to update spec mid-project without re-running the full pipeline | MEDIUM | Existing spec/design/tasks parsers; re-plan trigger; plan-checker |
| `/mysd:fix` command | Bug fixes are the most common unplanned work; no current path for code-only targeted repairs | MEDIUM | Worktree isolation; existing mysd-executor agent |
| `/mysd:model` CLI | Model profile exists in config but has no CLI surface; editing `mysd.yaml` manually to change profile is friction | LOW | Existing `ResolveModel` + `config.Load`; no new agents |
| `/mysd:lang` CLI | `response_language` / `document_language` in config but no CLI surface; locale sync with `openspec/config.yaml` is manual | LOW | Existing config system; new `openspec/config.yaml` writer |
| Task dependency + file overlap detection | Wave mode currently distributes tasks without checking `depends` or `files` — causes ordering bugs and race conditions on shared files | HIGH | Extends `TaskItem` schema + `tasks.md` format; planner must write `depends`/`files`; topological sort + overlap algorithm |
| Plan-checker auto-trigger | Every plan should be validated for MUST coverage; easy to miss requirements silently | MEDIUM | New `mysd-plan-checker` agent definition; existing MUST parser |
| Scan upgrade to language-agnostic | `/mysd:scan` only detects Go files; non-Go projects get no output at all | HIGH | Rewrite `internal/scanner/scanner.go`; language detection by extension + module file; `openspec/config.yaml` writer |

### Differentiators (Competitive Advantage)

Features that set mysd apart from generic AI coding assistants and from v1.0.

| Feature | Value Proposition | Complexity | Depends On |
|---------|-------------------|------------|------------|
| Interactive Discovery — dual-mode (research + general) | Most AI tools ask clarifying questions linearly; mysd spawns parallel `mysd-advisor` agents per gray area and brings back comparison tables — significantly higher spec quality; dual-loop conversation model (deep dive within area + discover new areas) | HIGH | New `mysd-researcher` + `mysd-advisor` agent definitions; Codebase Scout (orchestrator grep/glob); scope guardrail in SKILL.md conversation flow |
| Worktree-based parallel execution with AI conflict resolution | Git worktrees give each task a fully isolated working copy; AI resolves merge conflicts and validates with `go build`/`go test` (or language-appropriate equivalent); 3-retry before user notification; partial failure policy (one wave task failing does not block others) | HIGH | Task dependency + file overlap detection (prerequisite); `git worktree` CLI; branch naming `mysd/{change}/T{id}-{slug}`; `.worktrees/T{id}/` short paths for Windows compat |
| Shared research result across auto/ff/ffe | Research done once in `propose` stage and reused by `spec`, `plan`, etc. — avoids repeated AI calls for the same context; ff/ffe imply `--auto` and inherit the shared research | MEDIUM | Interactive Discovery (research mode); existing ff/ffe commands |
| Subagent architecture with per-stage orchestration | Each workflow stage has a dedicated orchestrator (SKILL.md) + specialized subagents; clear separation of concerns vs. monolithic agents; `mysd-proposal-writer` and `mysd-spec-writer` per capability area improve output quality | MEDIUM | New agent `.md` definitions under `plugin/agents/`; existing Task tool pattern in SKILL.md |
| Codebase Scout (zero new Go code) | Before propose/spec/discuss, orchestrator does targeted grep/glob to surface reusable patterns and integration points — grounds AI proposals in actual code; no new subagent needed | LOW | Orchestrator-level grep/glob steps in SKILL.md; no changes to Go binary |
| Scope guardrail in discovery | Prevents scope creep during gray-area exploration; redirects expansion ideas to deferred notes instead of silently widening scope — enforces discipline that most AI tools lack | LOW | Discovery loop logic in SKILL.md; no new Go code |
| `--auto` flag for non-interactive workflows | `propose`/`spec`/`discuss` can run without prompting; ff/ffe implicitly set `--auto`; combined with research reuse this enables fully automated propose→archive runs | MEDIUM | Interactive Discovery (research mode must produce shareable output); ff/ffe integration |

### Anti-Features (Commonly Requested, Often Problematic)

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Interactive discovery in design stage | Users want AI involvement at every phase | `design.md` records architectural decisions made after spec is settled; adding discovery here creates feedback loops between spec and design, causing spec churn | Keep design as decision-recording only; all discovery in propose/spec/discuss |
| GUI / Web dashboard for worktree status | Visual appeal, easier monitoring | Violates CLI-first constraint; adds runtime server dependency; out-of-scope per PROJECT.md | lipgloss-styled terminal summary; `mysd status` extended for worktree state |
| MCP server for subagent coordination | Seems like the "AI-native" integration pattern | Always-running process, infrastructure complexity, breaks single-binary deployment; already proven unnecessary in v1.0 | Binary-called-from-SKILL.md pattern (validated v1.0 architecture) |
| Real-time streaming output from parallel tasks | Feels more responsive during execution | Interleaved output from N worktrees is unreadable; requires TUI event loop (bubbletea) which is explicitly out of scope | Per-task progress written to `.worktrees/T{id}/.progress`; lipgloss summary printed after wave completes |
| Model selection per individual task | Fine-grained cost control | Cognitive overhead; defeats the purpose of model profiles; hard to reason about token costs | `ModelOverrides` map in config for per-agent-role overrides (already in `ProjectConfig`) |
| Auto-push worktree branches to remote | Convenient for team review | Out of scope — mysd is a single-developer tool; adds git remote as execution dependency | Local merge only; user pushes when ready |
| Auto-rewrite spec intent from code output | Prevent spec drift by keeping spec in sync with code | Inverts causality — spec drives code, not vice versa; auto-sync from code to spec destroys the spec-as-source-of-truth guarantee | After execution, verify that MUST items were satisfied; update spec STATUS field only, never intent |

---

## Feature Dependencies

```
[Interactive Discovery]
    └──requires──> [mysd-researcher agent definition]
    └──requires──> [mysd-advisor agent definition (x N parallel)]
    └──requires──> [Codebase Scout (grep/glob in orchestrator SKILL.md)]
    └──enhances──> [/mysd:propose SKILL.md]
    └──enhances──> [/mysd:spec SKILL.md]
    └──enhances──> [/mysd:discuss SKILL.md]
    └──produces──> [Shared Research Result]
                        └──consumed-by──> [Auto Mode / ff / ffe]

[/mysd:discuss]
    └──requires──> [existing spec/design/tasks parsers]
    └──triggers──> [re-plan: mysd plan re-run]
    └──triggers──> [Plan-checker auto-run]
    └──can-use──> [Interactive Discovery (optional research mode)]

[Plan-checker]
    └──requires──> [mysd-plan-checker agent definition]
    └──requires──> [existing MUST requirement parser]
    └──triggered-by──> [/mysd:plan completion]
    └──triggered-by──> [/mysd:discuss → re-plan path]

[Task Dependency + File Overlap Detection]
    └──requires──> [TaskItem schema extended: depends[] + files[]]
    └──requires──> [tasks.md format extended (planner writes depends/files)]
    └──requires──> [topological sort + file-overlap algorithm in execute orchestrator]
    └──enables──> [Worktree Parallel Execution (correct wave layering)]
    └──MUST precede──> [Worktree Parallel Execution]

[Worktree Parallel Execution]
    └──requires──> [Task Dependency + File Overlap Detection]
    └──requires──> [git worktree CLI (system dependency, standard git)]
    └──requires──> [mysd-executor agent (existing, extended for worktree context)]
    └──extends──> [existing wave execution in /mysd:execute SKILL.md]
    └──provides-isolation-for──> [/mysd:fix]

[/mysd:fix]
    └──requires──> [Worktree Parallel Execution isolation machinery]
    └──reuses──> [mysd-executor agent (existing)]
    └──redirects-spec-issues-to──> [/mysd:discuss]

[--auto flag]
    └──requires──> [Interactive Discovery (research mode must produce shareable output)]
    └──integrates-with──> [ff/ffe (already imply --auto behavior)]

[Scan Upgrade (language-agnostic)]
    └──replaces──> [existing Go-only scanner.go]
    └──requires──> [language detection: extension mapping + module file detection]
    └──requires──> [openspec/config.yaml writer (new)]
    └──enables──> [/mysd:init refactored as scan --scaffold-only]

[/mysd:model]
    └──requires──> [existing ResolveModel + config.Load]
    └──independent (no new agents, no schema change)

[/mysd:lang]
    └──requires──> [existing config system]
    └──requires──> [openspec/config.yaml writer (shared with Scan Upgrade)]

[Subagent Definitions (researcher, advisor, proposal-writer, plan-checker)]
    └──required-by──> [Interactive Discovery]
    └──required-by──> [Plan-checker]
    └──enhances──> [/mysd:propose (proposal-writer replaces inline writing)]
    └──enhances──> [/mysd:spec (spec-writer becomes per-capability-area spawn)]
```

### Dependency Notes

- **Task Dependency + File Overlap Detection must be built before Worktree Parallel Execution.** Worktree isolation without correct topological ordering will produce incorrect merge sequences and corrupt shared state even with AI conflict resolution.
- **Interactive Discovery requires Codebase Scout.** The researcher agent needs grounded starting points (actual code patterns found by grep/glob) before formulating questions; without this, researcher output hallucinates integration points.
- **/mysd:discuss triggers plan-checker.** Because discuss updates spec and auto-triggers re-plan, plan-checker must be complete before /mysd:discuss is considered done.
- **Scan Upgrade and /mysd:lang share the `openspec/config.yaml` writer.** Both features write to `openspec/config.yaml`; the writer should be implemented once (in Go binary) and reused.
- **Auto mode depends on Interactive Discovery producing a shareable research result.** Without this, ff/ffe cannot reuse research across stages and would either re-run research (wasteful) or skip it (lower quality).
- **/mysd:fix reuses worktree isolation.** Fix is not an independent worktree implementation — it calls the same machinery as parallel execution but with a single task and interactive repair loop.

---

## MVP Definition

This is a subsequent milestone (v1.1), not greenfield. "Launch with" = required for milestone closure.

### Launch With (v1.1 milestone)

- [ ] Task dependency + file overlap detection (`depends` + `files` in `TaskItem`; planner writes them; topological sort + overlap algorithm in execute orchestrator) — foundation for correct parallel execution
- [ ] Worktree parallel execution with AI conflict resolution, 3-retry policy, partial failure handling — core new execution capability
- [ ] Interactive Discovery dual-mode (research + general) in propose/spec/discuss — primary user-facing value of v1.1
- [ ] `/mysd:discuss` with auto re-plan + plan-checker trigger — closes the spec-update gap left in v1.0
- [ ] Plan-checker (`mysd-plan-checker` agent, auto-triggered after every plan) — mandatory MUST coverage quality gate
- [ ] New subagent definitions: `mysd-researcher`, `mysd-advisor`, `mysd-proposal-writer`, `mysd-plan-checker` — required by discovery and plan-checker
- [ ] `/mysd:fix` with worktree isolation — targeted code-only repair path
- [ ] Scan upgrade to language-agnostic (extension detection + module file detection + `openspec/config.yaml` writer) — unblocks non-Go users
- [ ] `/mysd:model` CLI (show/set/resolve) — surfaces existing config capability
- [ ] `/mysd:lang` CLI (interactive locale + `openspec/config.yaml` sync) — surfaces existing config + locale sync
- [ ] Model profile table extended for new agents (researcher, advisor, proposal-writer, plan-checker) — required for model profile system to be complete

### Add After Validation (v1.1.x)

- [ ] `--auto` flag polish across propose/spec/discuss — trigger: interactive discovery + research reuse works stably in v1.1
- [ ] Codebase Scout refinement (smarter grep patterns per detected language) — trigger: scanner upgrade shipped and language is known
- [ ] Incremental scan update mode (add new specs without overwriting existing ones) — trigger: users report scan overwriting manual edits

### Future Consideration (v2+)

- [ ] Multi-language model profile differentiation (different defaults for Python vs Go projects) — defer: language-agnostic scan baseline is enough for v1.1
- [ ] Worktree progress streaming to terminal — defer: per-task `.progress` file pattern is sufficient; TUI is explicitly out of scope
- [ ] Support for AI tools beyond Claude Code — explicitly out of scope for v1.x

---

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Task dependency + file overlap detection | HIGH | HIGH | P1 — prerequisite for correct parallel execution |
| Worktree parallel execution | HIGH | HIGH | P1 — core v1.1 capability |
| Interactive Discovery (propose/spec/discuss) | HIGH | HIGH | P1 — primary user-facing value |
| /mysd:discuss + auto re-plan | HIGH | MEDIUM | P1 — closes major v1.0 workflow gap |
| Plan-checker | HIGH | LOW | P1 — every plan needs MUST coverage check |
| New subagent definitions (researcher, advisor, proposal-writer, plan-checker) | HIGH | MEDIUM | P1 — required by discovery and plan-checker |
| /mysd:fix | MEDIUM | MEDIUM | P1 — unblocks targeted repairs without full re-execute |
| Scan upgrade (language-agnostic) | HIGH | HIGH | P1 — non-Go users get nothing from scan currently |
| /mysd:model CLI | MEDIUM | LOW | P2 — config works; CLI is convenience surface |
| /mysd:lang CLI | MEDIUM | LOW | P2 — config works; CLI + locale sync is convenience |
| --auto flag | MEDIUM | MEDIUM | P2 — ff/ffe already work; --auto is an enhancement |
| Codebase Scout | MEDIUM | LOW | P2 — orchestrator grep/glob, no new Go code |
| Model profile extended for new agents | LOW | LOW | P3 — existing defaults are usable as fallback |

**Priority key:**
- P1: Must have for v1.1 milestone closure
- P2: Should have, add when P1 is stable
- P3: Nice to have, future consideration

---

## Complexity Notes by Feature Group

### HIGH Complexity — Requires Careful Phase Ordering

**Task Dependency + Worktree Parallel Execution** (one logical unit across two features):
- Schema change: `TaskItem` and `tasks.md` format must carry `depends []int` and `files []string`
- Planner agent must be updated to write these fields; existing tasks.md files are backward-compatible (empty = no deps)
- Wave layering algorithm: topological sort of tasks by `depends`, then file-overlap check within each layer to split tasks into separate waves
- Git worktree lifecycle: `git worktree add .worktrees/T{id} -b mysd/{change}/T{id}-{slug}` → execute → merge by task ID order (`--no-ff`) → cleanup
- Merge conflict resolution in SKILL.md: AI reads both versions, resolves, runs build+test, retries up to 3x before notifying user
- Partial failure policy: one failing task must not block others in the same wave
- Windows path constraint: `.worktrees/T{id}/` short path convention (git has 260-char path limit on Windows without long path enabled)

**Interactive Discovery** (research mode):
- Spawning `mysd-researcher` to scan codebase AND domain context, then `mysd-advisor` x N in parallel per gray area
- Dual-loop conversation model: within-area deep dive (can drill further) + cross-area new area discovery loop
- Scope guardrail: redirect scope expansion to deferred notes, not inline discussion
- Shared research result: proposal orchestrator stores research output for reuse by spec, plan stages in same ff/ffe run
- Each command that supports research (propose/spec/discuss) must interactively ask "use research mode?" unless `--auto` is set

**Scan Upgrade (language-agnostic)**:
- Replace Go-specific `PackageInfo` struct; new `FileInfo` or `ModuleInfo` abstraction
- Language detection: file extension mapping (`.py`, `.ts`, `.js`, `.rb`, `.rs`, `.java`, `.cs`, etc.) + module file detection (`package.json`, `Cargo.toml`, `pyproject.toml`, `go.mod`, `pom.xml`, `build.gradle`)
- `openspec/config.yaml` writer: new Go function, shared by scan and lang commands
- Incremental update: if `config.yaml` exists, update specs; never overwrite config
- Interactive locale prompt on first `config.yaml` creation (only when building new config)

### MEDIUM Complexity

**/mysd:discuss**:
- Must update spec/design/tasks — three separate file writers, three parsers already exist
- Must trigger re-plan (invoke plan command path or instruct orchestrator to re-run `/mysd:plan`)
- Must trigger plan-checker after re-plan
- Optional research mode (same interactive prompt as propose/spec)

**/mysd:fix**:
- Isolated from spec workflow — code-only repair, no spec writes
- Reuses worktree isolation machinery (must be built after worktree execution)
- Redirect: if user raises a spec-level issue, fix must route them to `/mysd:discuss`
- Interactive repair loop: AI attempts fix, validates, retries up to 3x

**New Subagent Definitions** (`plugin/agents/*.md`):
- `mysd-proposal-writer`: replaces inline proposal writing in `/mysd:propose` SKILL.md; receives research context
- `mysd-spec-writer`: changes from one shared agent to per-capability-area spawn (already implied by proposal)
- `mysd-advisor`: stateless per gray-area analysis; receives one gray area question; outputs comparison table + recommendation
- `mysd-researcher`: receives change context + codebase scout output; produces structured question list for discovery
- `mysd-plan-checker`: receives tasks.md + all MUST items; outputs gap report (uncovered MUST items); offers auto-fill

**--auto flag**:
- `--auto bool` cobra flag on propose/spec/discuss commands
- SKILL.md checks flag value: if auto, skip all interactive prompts and use recommended options
- ff/ffe set `--auto` implicitly (already implied by their fast-forward semantics)

### LOW Complexity

**/mysd:model**: New `cmd/model.go` + SKILL.md; reads/writes existing config via `config.Load`/viper; no schema changes
**/mysd:lang**: New `cmd/lang.go` + SKILL.md; reads/writes config + writes `openspec/config.yaml`; interactive locale prompt
**Codebase Scout**: Pure SKILL.md change — orchestrator adds grep/glob steps before spawning researcher; no Go binary changes
**Scope guardrail**: SKILL.md conversation flow addition in discovery loop; no Go code needed

---

## Competitor Feature Analysis

| Feature | OpenSpec CLI | GSD system | mysd v1.0 | mysd v1.1 target |
|---------|-------------|------------|-----------|------------------|
| Interactive gray area discovery | Manual Q&A | discuss-phase with subagents + research | None (linear propose) | Dual-mode: research + general; parallel advisors per gray area |
| Parallel task execution | None | Wave mode | Basic wave (no dep check) | Wave with deps + file overlap + worktree isolation |
| Merge conflict resolution | Manual | Manual | Manual | AI-assisted, 3-retry, build-validated |
| Language-agnostic scan | None | N/A | Go-only | All major languages via extension + module detection |
| Ad-hoc spec updates | Manual file edit | /gsd:quick | None | /mysd:discuss with auto re-plan |
| Targeted bug fix | Manual | /gsd:debug | None | /mysd:fix with worktree isolation |
| Model profile CLI | None | None | Config-only | /mysd:model show/set/resolve |
| Locale management | Manual | None | Config-only | /mysd:lang interactive + openspec sync |
| Plan MUST coverage check | None | Manual | None | Auto plan-checker after every plan |
| Scope creep prevention | None | Partial | None | Scope guardrail in discovery loop |

---

## Sources

- `.specs/changes/interactive-discovery/proposal.md` — authoritative v1.1 feature list, HIGH confidence (primary source)
- `.planning/PROJECT.md` — v1.1 milestone scope, constraints, out-of-scope items, HIGH confidence
- `internal/scanner/scanner.go` — current Go-only scanner implementation verified directly, HIGH confidence
- `internal/config/config.go`, `internal/config/defaults.go` — current model profile and config system, HIGH confidence
- `internal/executor/context.go`, `cmd/execute.go` — current TaskItem schema and execute orchestration, HIGH confidence
- `plugin/commands/mysd-execute.md` — current wave execution SKILL.md (basic, no dep ordering), HIGH confidence
- `plugin/agents/` — existing 8 agent definitions (designer, executor, fast-forward, planner, scanner, spec-writer, uat-guide, verifier), HIGH confidence
- GSD discuss-phase pattern (referenced in proposal motivation as proven approach) — MEDIUM confidence (described in proposal, pattern origin in GSD system)

---
*Feature research for: mysd v1.1 — Interactive Discovery, Parallel Execution, Subagent Architecture, Language-Agnostic Scan*
*Researched: 2026-03-25*
