# Project Research Summary

**Project:** my-ssd
**Domain:** Go CLI tool — Claude Code plugin for Spec-Driven Development (SDD)
**Researched:** 2026-03-23
**Confidence:** HIGH

## Executive Summary

my-ssd is a Go binary + Claude Code plugin that bridges the gap between OpenSpec's spec-driven methodology and GSD's execution engine. Research confirms that no existing tool provides both (a) structured spec management with OpenSpec-compatible artifacts and (b) automated goal-backward verification that closes the feedback loop after execution. The recommended approach is a layered architecture: a single Go binary handles all business logic (spec parsing, state machine, verification pipeline), while a thin Claude Code plugin layer exposes slash commands and agent definitions as Markdown wrappers. The binary is invoked via `Bash` tool calls from within skill files — not as an MCP server — keeping the integration simple and portable.

The core differentiator is mandatory, non-bypassable spec alignment before execution combined with goal-backward verification after. Every competing tool (OpenSpec, GSD, GitHub Spec Kit, BMAD) either lacks the alignment gate, lacks automated post-execution verification, or lacks both. my-ssd's spec feedback loop — where verification results are written back into spec metadata — is unique in the landscape. The single Go binary with no Node.js dependency is a distinct distribution advantage.

The primary risks are (1) spec format lock-in from hardcoded parser assumptions, which must be addressed in Phase 1 with a schema-versioned tolerant parser; (2) AI self-verification blindness, where the same agent that executed also verifies — prevented by running the verifier with a fresh context and file-system-only evidence; and (3) Claude Code plugin API instability, mitigated by keeping all logic in the Go binary so the thin plugin layer can be rewritten cheaply if the API changes.

---

## Key Findings

### Recommended Stack

The stack is well-defined and low-risk. Go 1.23+ with Cobra v1.10.2 is the correct choice for a 7-subcommand CLI — it is the de facto standard (used by kubectl, helm, GitHub CLI) and handles the full propose→archive command tree cleanly. `adrg/frontmatter` v0.2.0 handles YAML frontmatter extraction without pulling in a full Markdown engine, which is the right scope for spec file parsing. `gopkg.in/yaml.v3` is already in the dependency tree via Cobra and serves OpenSpec's YAML frontmatter format. Viper v1 handles project and user config without requiring Viper v2 (unreleased). The critical integration insight: the Go binary is NOT a plugin itself — it is a shell command invoked by SKILL.md files, and the plugin is a directory of Markdown files with proper frontmatter.

**Core technologies:**
- Go 1.23+: primary language — single-binary deployment, cross-platform via GOOS/GOARCH, no runtime dependency
- github.com/spf13/cobra v1.10.2: CLI framework — de facto standard for multi-subcommand CLIs; handles all 7 commands
- gopkg.in/yaml.v3: YAML parsing — required for OpenSpec frontmatter; already in dependency tree via Cobra
- github.com/adrg/frontmatter v0.2.0: frontmatter extraction — no Markdown engine dependency; supports YAML/TOML/JSON delimiters
- github.com/spf13/viper v1.x: configuration — reads `.mysd.yaml` and `~/.mysd/config.yaml` with env override
- github.com/stretchr/testify v1.x: test assertions — standard library for Go CLI testing
- goreleaser v2.14+: release automation — GitHub Releases with Linux/macOS/Windows binaries; use casks NOT formulae (formulae deprecated June 2025)

**What NOT to use:**
- bubbletea: full TUI framework, overkill for a CLI that outputs structured text
- MCP server pattern: adds infrastructure complexity for no benefit over Bash invocation
- Viper v2: unreleased as of research date — stay on v1

### Expected Features

Research confirms a clear P1/P2/P3 structure. The MVP loop is: spec artifacts → workflow commands → pre-execution alignment gate → RFC 2119 parsing → goal-backward verification → spec feedback loop → session continuity → archive. This is the minimum set that validates whether spec-to-execution-to-verification is tighter than using OpenSpec + GSD separately.

**Must have (table stakes — v1):**
- Structured spec artifacts (proposal.md / specs/ / design.md / tasks.md) — every SDD tool organizes this way; missing = incomplete
- Linear workflow commands (propose → spec → design → plan → execute → verify → archive) — users trained to expect named phases
- Pre-execution alignment gate — mandatory spec read-and-acknowledge before any code is written; the single most important behavior change
- RFC 2119 keyword parsing (MUST / SHOULD / MAY) — required foundation for goal-backward verification
- Goal-backward verification — spec MUST items → verification checklist → pass/fail verdict; the core differentiator
- Spec feedback loop — verification results written back to spec status; makes spec a living document
- Session continuity (STATE.md pattern) — cross-session recovery; essential for multi-day projects
- Archive / history — completed specs move to `.specs/archive/`; project hygiene
- Single Go binary — the entire distribution story; without this, installation friction equals Node.js tools
- OpenSpec format compatibility — brownfield adoption by existing OpenSpec users

**Should have (competitive — v1.x):**
- Delta Specs (ADDED / MODIFIED / REMOVED) — scope verification narrowly to change type; add when users report broad verification
- Brownfield codebase onboarding (`/mysd:onboard`) — generates CONVENTIONS.md, ARCHITECTURE.md from existing code; add when users report AI ignoring existing patterns
- Atomic git commits per task — traceable task→commit mapping; low effort, high audit value
- `/mysd:design` and `/mysd:plan` commands — expand from 4-command MVP to full 7-command suite

**Defer (v2+):**
- Multi-agent wave execution — validate single-agent reliability first; wave orchestration adds significant complexity
- Multi-runtime support (Cursor, Gemini CLI) — abstract interface now, implement after Claude Code integration proven
- Spec templates / profiles — defer until user-generated spec patterns reveal what should be standardized

### Architecture Approach

The architecture is a clean three-layer system: Claude Code Plugin Layer (thin Markdown wrappers) → CLI Core (Cobra command tree in `cmd/`) → Internal Engine Layer (`internal/spec/`, `internal/engine/`, `internal/verify/`, `internal/state/`). The key design principles are: (1) thin commands, fat internal — zero business logic in `cmd/`; (2) spec as struct, not string — parse OpenSpec Markdown into typed Go structs at the boundary, all downstream logic operates on `spec.Change`, `spec.Requirement`, `spec.Task`; (3) explicit workflow state machine — `.ssd-state.json` owns the phase cursor, not filesystem presence inference; (4) plugin delegates to binary — the Claude Code plugin is purely a presentation layer, the Go binary owns all mutations.

**Major components:**
1. Claude Code Plugin Layer (`plugin/`) — slash commands + agent definitions as Markdown files; invokes binary via Bash tool; zero business logic
2. Cobra Command Tree (`cmd/`) — parse CLI arguments, route to internal functions, surface user-facing errors; one file per command
3. Spec Engine (`internal/spec/`) — parse/validate OpenSpec Markdown; generate artifact scaffolds; write delta specs; ALL format knowledge isolated here
4. Execution Engine (`internal/engine/`) — orchestrate single/multi-agent runs; build spec context injected before AI runs; manage wave-based execution in v2
5. Verification Pipeline (`internal/verify/`) — collect MUST items; goal-backward check per requirement; write verification report; feed results back to spec
6. State Manager (`internal/state/`) — track current phase in `.ssd-state.json`; enable crash recovery and `--resume`; the only authoritative phase record

**Build order (dependency-driven):**
Storage schema → internal/spec/ → internal/state/ → cmd/ skeleton → internal/engine/ → internal/verify/ → plugin/ → internal/engine/wave (v2)

### Critical Pitfalls

1. **Spec format lock-in** (Phase 1) — baking hardcoded heading strings like `"## Requirements"` into the parser; even a user renaming a heading silently breaks everything. Prevention: schema-driven parser with `spec-version` in frontmatter from day one; test against real OpenSpec projects, not synthetic fixtures.

2. **AI self-verification blindness** (Phase 3) — using the same agent that executed to verify its own work; results in 100% pass rate regardless of actual implementation. Prevention: verifier runs with fresh context, receives only spec MUST items + filesystem state; never the original execution transcript.

3. **Spec drift** (Phase 2-3) — AI execution makes undocumented architectural decisions; `.specs/` describes a system that no longer exists within weeks. Prevention: `archive` must fail if verification detected open MUST failures; verification is not optional or skippable.

4. **Context window budget overflow** (Phase 1) — too many or too large SKILL.md files; some are silently excluded from context; agent behavior becomes inconsistent across sessions. Prevention: keep each SKILL.md under 500 lines; use `disable-model-invocation: true` on task-specific skills; run `/context` to check for excluded skills warning.

5. **Claude Code plugin API instability** (Phase 1) — `commands/` deprecated in favor of `skills/`; `permissionMode` disallowed in plugin agents; these break silently on updated Claude Code versions. Prevention: use `skills/<name>/SKILL.md` format exclusively; keep all business logic in Go binary so plugin layer can be rewritten cheaply.

---

## Implications for Roadmap

Based on research, the build order is dictated by hard dependencies: spec parsing is required before execution, execution is required before verification, and the plugin layer is the last thing written (after the binary commands are stable). Suggested 4-phase structure:

### Phase 1: Foundation — Spec Data Model and CLI Skeleton

**Rationale:** Everything depends on the spec parser and state machine. Getting these right (schema-versioned, tolerant of OpenSpec variations, typed Go structs) prevents the most expensive pitfall (spec format lock-in). The plugin architecture decisions (thin wrappers, SKILL.md format, size discipline) must also be established here before content is written.

**Delivers:** Working spec artifact scaffold (`myssd propose`), spec parser that handles OpenSpec format, state machine with `.ssd-state.json`, Cobra command tree skeleton with all 7 commands registered, plugin directory with correct SKILL.md structure.

**Features from FEATURES.md:** Structured spec artifacts, Session continuity (STATE.md), Single Go binary, Convention over configuration, OpenSpec format compatibility foundation.

**Avoids from PITFALLS.md:** Spec format lock-in (schema-versioned parser), Claude Code plugin API breaking changes (establish thin wrapper discipline), RFC 2119 case sensitivity (parser unit tests), Windows binary naming (choose `myssd`, not `install`/`setup`).

**Research flag:** Standard patterns — Cobra CLI structure, Go project layout, and OpenSpec format are all well-documented. No phase-level research needed.

---

### Phase 2: Execution Engine

**Rationale:** With spec parsing and state management in place, the execution engine can be built on top of typed spec structs. The coordinator design (binary owns all `.specs/` state mutations; agents only read and report) must be established here to prevent the multi-agent shared state race condition pitfall later.

**Delivers:** Working `myssd execute` command, context builder (assembles spec context injected before AI runs), single-agent task runner, task completion tracking in `tasks.md`, pre-execution alignment gate (mandatory spec read-and-acknowledge), session continuity across runs.

**Features from FEATURES.md:** Pre-execution alignment gate, Core workflow commands (execute), Spec as source of truth, Atomic git commits per task (low effort, add here).

**Avoids from PITFALLS.md:** Multi-agent shared state races (establish coordinator ownership now), spec drift (hook for post-execute verification gate), loading all spec files into context on every command (select relevant files only).

**Research flag:** May benefit from phase research on Claude Code subagent invocation patterns and context injection mechanics — the agent-spawn boundary (Go binary prepares context, Claude Code plugin spawns the agent) is the most novel integration point.

---

### Phase 3: Goal-Backward Verification and Spec Feedback Loop

**Rationale:** This is the core differentiator. Goal-backward verification depends on RFC 2119 keyword parsing (from Phase 1) and execution state (from Phase 2). The feedback loop (writing results back to spec metadata) depends on the verifier. This phase delivers the feature that no other SDD tool provides.

**Delivers:** Working `myssd verify` command, MUST item collector, goal-backward verifier (fresh context, file-system evidence only), verification report writer, spec status update (feedback loop), `archive` gate that refuses if MUST failures exist.

**Features from FEATURES.md:** RFC 2119 keyword support, Goal-backward verification, Spec feedback loop, Archive / history.

**Avoids from PITFALLS.md:** AI self-verification blindness (fresh context for verifier), spec drift (archive gate on MUST failures), eager AI verification on every save (gate behind explicit `verify` command only).

**Research flag:** Goal-backward verification with AI agents is a niche pattern. Phase research recommended on verification prompting strategies and evidence collection approaches — specifically how to structure the "does this codebase satisfy this MUST requirement?" query to minimize false positives.

---

### Phase 4: Full Plugin Layer and Distribution

**Rationale:** The Claude Code plugin is written last, after binary commands are stable. This ensures the plugin Markdown files are accurate wrappers, not aspirational documentation. Distribution infrastructure (GoReleaser, Homebrew, code signing) is also here.

**Delivers:** Complete Claude Code plugin (`plugin/` directory with all skills, agents, hooks), `/mysd:propose` through `/mysd:archive` slash commands, `spec-writer`, `task-runner`, and `verifier` agents, GoReleaser configuration for cross-platform binary distribution, Homebrew cask (not formula), macOS code signing configuration.

**Features from FEATURES.md:** Claude Code slash commands, Complete 7-command workflow, Convention over configuration defaults.

**Avoids from PITFALLS.md:** GoReleaser Homebrew formula deprecation (use casks from day one), macOS Gatekeeper (configure signing in GoReleaser), Windows binary naming (confirm `myssd` does not trigger UAC), plugin skill context budget (run `/context` check before release).

**Research flag:** Distribution specifics (GoReleaser cask config, Apple Developer ID signing, GitHub Actions PAT for tap repo) warrant a focused lookup during planning. These are well-documented but have version-specific gotchas (formulae deprecated June 2025).

---

### Phase Ordering Rationale

- **Spec parsing before execution:** The execution engine operates on `spec.Change` structs, never on raw Markdown — it cannot be built until the parser produces those structs.
- **State machine before commands:** Every command validates phase transitions before running — the state machine must exist before command business logic.
- **Execution before verification:** The verification pipeline checks MUST items against implemented code — there is nothing to verify until execution has produced output.
- **Plugin last:** Plugin Markdown files describe what the binary commands do — writing them before the binary is stable produces inaccurate wrappers that mislead Claude.
- **Sequential-first, parallel-optional:** Multi-agent wave execution is deferred until single-agent execution is proven reliable. This avoids the shared state race condition pitfall and reduces Phase 2 scope significantly.

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 2:** Claude Code subagent invocation mechanics and context injection patterns — the agent-spawn boundary is the least-documented integration point.
- **Phase 3:** Goal-backward verification prompting strategies — how to query an AI to determine whether a spec requirement is satisfied in a codebase, without self-verification blindness.
- **Phase 4:** GoReleaser cask configuration, Apple Developer ID code signing workflow, GitHub Actions PAT setup for Homebrew tap — version-specific gotchas worth a focused lookup.

Phases with standard patterns (skip research-phase):
- **Phase 1:** Go project layout, Cobra CLI structure, OpenSpec format, YAML frontmatter parsing, and state machine patterns are all well-documented. Standard Go patterns apply directly.

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Core libraries verified from official sources (Cobra releases, goreleaser.com, golangci-lint.run). Claude Code plugin format verified from official Anthropic docs. Only `adrg/frontmatter` version is MEDIUM (search result, not releases page). |
| Features | HIGH | Primary sources: OpenSpec repo, GSD repo, GitHub Spec Kit announcement, BMAD docs, intent-driven.dev. Competitor feature matrix based on direct tool inspection. |
| Architecture | HIGH | Claude Code plugins reference (official), OpenSpec workflow docs, Go CLI structure patterns from multiple high-quality sources. Build order derived from hard dependency analysis, not opinion. |
| Pitfalls | HIGH (Critical), MEDIUM (Integration/Performance) | Critical pitfalls sourced from real incident reports (OpenSpec issue #666, SDD failure analyses, agentic engineering research). Integration gotchas from official Claude Code docs. Performance traps from general CLI patterns — MEDIUM because my-ssd's specific scale is unknown. |

**Overall confidence:** HIGH

### Gaps to Address

- **Claude Code plugin minimum version:** The plugin manifest should declare a minimum Claude Code version, but the exact version field and enforcement mechanism was not verified. Confirm during Phase 4 planning.
- **Agent invocation API from Go binary:** The exact mechanism by which the Go binary signals Claude Code to spawn a subagent (versus Claude Code doing it autonomously from plugin instructions) needs verification during Phase 2. The research describes the pattern but the API call specifics are not pinned.
- **adrg/frontmatter v0.2.0 stability:** Version confirmed via search result, not releases page. Verify before writing go.mod.
- **Viper v1 + Cobra v1.10.2 compatibility:** Both are confirmed stable, but the specific viper config loading behavior for `.mysd.yaml` in project root vs. `~/.mysd/config.yaml` should be tested early in Phase 1 to avoid surprises.

---

## Sources

### Primary (HIGH confidence)
- [Claude Code Plugins Reference](https://code.claude.com/docs/en/plugins-reference) — plugin directory structure, SKILL.md frontmatter schema, agent fields, hooks format
- [Claude Code Skills / Slash Commands](https://code.claude.com/docs/en/slash-commands) — SKILL.md format, context budget limits, disable-model-invocation, allowed-tools
- [OpenSpec GitHub Repository](https://github.com/Fission-AI/OpenSpec) — artifact model, workflow commands, brownfield support, delta specs, issue #666 (format lock-in)
- [GSD (get-shit-done) GitHub Repository](https://github.com/gsd-build/get-shit-done) — wave execution, STATE.md pattern, goal-backward verification, Nyquist Layer
- [goreleaser.com](https://goreleaser.com/) — v2.14 confirmed; cask vs formula deprecation (June 2025)
- [golangci-lint.run](https://golangci-lint.run/) — v2 released March 2025; linters.default: standard config
- [github.com/spf13/cobra releases](https://github.com/spf13/cobra/releases) — v1.10.2 confirmed latest
- [RFC 2119](https://datatracker.ietf.org/doc/html/rfc2119) — MUST/SHOULD/MAY definitions; case sensitivity requirement

### Secondary (MEDIUM confidence)
- [GitHub Spec Kit announcement](https://github.blog/ai-and-ml/generative-ai/spec-driven-development-with-ai-get-started-with-a-new-open-source-toolkit/) — spec/plan/tasks/implement workflow comparison
- [BMAD-METHOD Documentation](https://docs.bmad-method.org/) — specialized agent architecture, brownfield guide
- [SDD Brownfield Guide - intent-driven.dev](https://intent-driven.dev/blog/2026/03/10/spec-driven-development-brownfield/) — AI-generated spec inaccuracy risks
- [Agentic Engineering Part 6: Forensic Verification](https://www.sagarmandal.com/2026/03/15/agentic-engineering-part-6-forensic-verification-why-a-perfect-score-from-your-ai-agent-should-make-you-nervous/) — self-verification blindness evidence
- [How We Built True Parallel Agents With Git Worktrees](https://dev.to/getpochi/how-we-built-true-parallel-agents-with-git-worktrees-2580) — multi-agent isolation patterns
- [AI Workflow Patterns in Go](https://dasroot.net/posts/2026/02/ai-workflow-patterns-go-cli-tools-agents/) — goroutine-based orchestration patterns
- [Structuring Go Code for CLI Applications](https://www.bytesizego.com/blog/structure-go-cli-app) — cmd/internal/pkg organization
- [awesome-claude-code](https://github.com/hesreallyhim/awesome-claude-code) — Claude Code plugin ecosystem patterns

### Tertiary (LOW confidence / inferred)
- [github.com/adrg/frontmatter](https://github.com/adrg/frontmatter) — v0.2.0 version from search result; verify on releases page before use
- WebSearch: Go CLI framework comparison 2025 — multiple sources agree on Cobra for complex CLIs; no single authoritative source

---
*Research completed: 2026-03-23*
*Ready for roadmap: yes*
