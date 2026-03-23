# Pitfalls Research

**Domain:** AI-assisted Spec-Driven Development CLI tooling (Go binary + Claude Code plugin)
**Researched:** 2026-03-23
**Confidence:** HIGH (Critical pitfalls), MEDIUM (Integration gotchas), MEDIUM (Performance traps)

---

## Critical Pitfalls

### Pitfall 1: Spec Format Lock-In via Hardcoded Structure Expectations

**What goes wrong:**
The Go parser bakes in assumptions about spec file layout (heading names, section order, required fields). Any deviation — even a user renaming `## Requirements` to `## Functional Requirements` — silently fails validation or corrupts archive operations. OpenSpec itself has this exact bug (GitHub issue #666): custom schemas break because parsers hardcode `### Requirement:` and `#### Scenario:` inline format expectations.

**Why it happens:**
It is faster to write `strings.Contains(line, "## Requirements")` than to build a tolerant, schema-aware Markdown parser. Early prototypes hard-code headings; those choices solidify before anyone notices.

**How to avoid:**
Define the canonical spec format as a versioned schema document, not as parser assumptions. Parse by semantic meaning (frontmatter + content blocks), not by literal heading text. Build tolerance for common variations from day one. Store the spec schema version in the file's YAML frontmatter (`spec-version: 1`), enabling future migration.

**Warning signs:**
- Parser code contains literal heading strings like `"## Requirements"` or `"### Requirement:"`
- Brownfield OpenSpec projects fail to load on first run
- Any "must exactly match" language in internal spec documentation without a migration path

**Phase to address:**
Phase 1 (Core spec parsing) — this is a foundational decision that becomes prohibitively expensive to change later.

---

### Pitfall 2: AI Self-Verification Blindness — Verification Without Independence

**What goes wrong:**
The same agent that generated the spec or code is used to verify that the spec was satisfied. The agent is blind to its own assumptions: it marks all MUST requirements as "satisfied" because it hallucinated compliance rather than actually checking implementation. This was observed specifically in real-world SDD testing where agents marked test suites complete without writing a single test.

**Why it happens:**
It is natural to use the same AI context for both execution and verification. The verification step becomes a rubber-stamp: the agent confirms what it already believes to be true, not what the code actually does.

**How to avoid:**
Goal-backward verification must run with a fresh context, given only the spec and the file system state — never the original execution transcript. The verifier reads the spec's MUST items, searches for evidence in the codebase, and produces a structured result. The Go binary owns the verification orchestration and provides the file list to check — the AI provides judgment on what it finds.

**Warning signs:**
- Verification agent receives the original execution context as input
- Verification pass rate is consistently 100% after AI execution
- Verify step completes faster than it took to read all files it should be checking

**Phase to address:**
Phase 3 (Verification engine) — but the verification architecture decision must be made in Phase 1 before the execution model is locked in.

---

### Pitfall 3: Spec Drift — Specs Become Retrospective Documentation

**What goes wrong:**
AI execution makes unanticipated architectural decisions that the spec does not reflect. Each iteration accumulates undocumented choices. Within weeks, the `.specs/` directory describes a system that no longer exists. Future AI executions misfire because they read stale intent.

**Why it happens:**
There is no mechanism that forces spec updates when implementation diverges. The spec is written before execution; nothing enforces reconciliation after. This is the same failure mode as API documentation drift, which causes integration failures costing teams days per incident in enterprise environments.

**How to avoid:**
After each `execute` phase, run an automatic diff between spec MUST/SHOULD items and the actual implementation. Surface deltas as first-class output, not buried in logs. Make the `verify` command produce a structured diff that the `archive` command can use to update the spec state. Never allow `archive` to succeed if verification detected open MUST failures.

**Warning signs:**
- The `verify` command is optional or skippable in the workflow
- Archived specs contain unchecked MUST items marked as "done"
- Users routinely run `execute` multiple times without running `verify` in between

**Phase to address:**
Phase 2 (Execute engine) for the hook, Phase 3 (Verification) for the mechanism.

---

### Pitfall 4: Context Window Budget Overflow Silently Degrades Quality

**What goes wrong:**
Claude Code's skill system has a dynamic character budget (2% of context window, fallback 16,000 characters). When the plugin ships many large skill files, some are silently excluded from context. Agents execute without the full workflow instructions. Research shows accuracy drops measurably once context utilization exceeds 60–70% of the window, particularly for items in the middle of the context.

**Why it happens:**
Plugin developers write comprehensive skills to improve AI behavior. Each skill is added incrementally. No single skill causes the overflow; the combination does. The failure is silent — no error, just degraded behavior.

**How to avoid:**
Keep each `SKILL.md` under 500 lines (this is the official Claude Code recommendation). Use progressive disclosure: put navigation instructions in `SKILL.md` and detailed reference in `reference.md` files that Claude loads on demand. Use `disable-model-invocation: true` on task-specific skills so they only load when explicitly invoked, not on every session. Run `/context` during development to check for excluded skills warning.

**Warning signs:**
- Any single `SKILL.md` file exceeds 400 lines
- The plugin ships more than 6–8 skills without measuring total description character count
- AI agent behavior becomes inconsistent across sessions without explanation
- `/context` output shows "excluded skills" warning

**Phase to address:**
Phase 1 (Plugin architecture) — establish skill size discipline before writing content.

---

### Pitfall 5: Multi-Agent Shared State Race Conditions

**What goes wrong:**
When parallel agents execute in the same working directory, they compete for the same files. Agent A reads `tasks.md`, Agent B writes `tasks.md` with updated status, Agent A writes an overwrite based on stale data. The result is lost progress and corrupted spec state. This is described as "not parallelism, that is a race condition" by practitioners who have built parallel AI systems.

**Why it happens:**
The default multi-agent model is a natural extension of single-agent execution. Developers add parallelism for speed without first establishing isolation boundaries.

**How to avoid:**
For parallel mode, use git worktree isolation per agent — each agent gets its own branch and working directory. Shared state (spec status, task completion) must be written through the Go binary coordinator, not by agents directly. The Go binary owns `.specs/` state mutations; agents only read state and report results. Sequential-by-default (the project's stated design) avoids this pitfall for v1; parallel mode should be explicitly scoped to Phase 4+.

**Warning signs:**
- Parallel agents are given write access to `.specs/` files
- No locking or coordinator pattern exists for spec state mutations
- Agent outputs are merged by string concatenation rather than structured diff

**Phase to address:**
Phase 2 (Execute engine coordinator design), Phase 4 (Parallel mode) if implemented.

---

### Pitfall 6: Claude Code Plugin Architecture Breaking Changes

**What goes wrong:**
Claude Code's plugin API has evolved rapidly in 2025–2026: `commands/` became `skills/`, slash commands merged into skills, the `context: fork` frontmatter field appeared, and `hooks`, `mcpServers`, `permissionMode` are now explicitly disallowed in plugin-shipped agents for security reasons. A plugin built on early-2025 assumptions silently breaks on mid-2025+ Claude Code versions.

**Why it happens:**
Claude Code is a fast-moving product. The plugin system is not yet stable API. Features added one month may have their interface changed the next.

**How to avoid:**
Pin to documented stable interfaces only. Avoid undocumented plugin behaviors. Use the `SKILL.md` + frontmatter approach (not raw `commands/` directory) since skills are the forward-declared standard. Test against Claude Code updates in a separate test environment before releasing. Document the minimum Claude Code version in the plugin manifest. Structure the plugin so the Go binary contains all business logic; the Claude Code layer is thin wrappers that are easy to rewrite.

**Warning signs:**
- Plugin code uses `commands/` directory without a migration plan to `skills/`
- Plugin relies on `permissionMode` in agent frontmatter (disallowed in plugins)
- No minimum version requirement documented for Claude Code

**Phase to address:**
Phase 1 (Plugin architecture) — the separation of concerns between Go binary and Claude Code plugin layer.

---

## Technical Debt Patterns

Shortcuts that seem reasonable but create long-term problems.

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Hard-code spec heading names in parser | 2 hours saved | Rewrite required when any real project uses different headings; brownfield support breaks | Never — use schema-driven parsing from day one |
| Use same AI agent for execute + verify | No extra orchestration code | Self-verification blindness; silent false positives on verification | Never for MUST-item verification |
| Skip spec version in file frontmatter | Simpler initial format | Cannot migrate existing specs when format evolves; brownfield users stranded | Never — add `spec-version` in initial spec format |
| Write all workflow logic in SKILL.md | Easier to prototype | Exceeds context budget; business logic not testable outside Claude Code | Prototype phase only — move logic to Go binary before Phase 2 |
| Single global spec state file | Simpler read/write | Race condition surface for parallel agents; large files slow incremental updates | Acceptable in v1 sequential mode only |
| GoReleaser formulae (not casks) for Homebrew | Auto-generated | Homebrew deprecated formulae for precompiled binaries as of June 2025; distribution breaks | Never — use casks from day one |
| Store API keys in plugin config files | Convenient for testing | Security violation; keys persist in version-controlled `.claude/settings.json` | Never |
| Verification as an optional CLI flag | Simpler initial UX | Users skip verification; spec drift accumulates silently | Never — verify should be integrated into archive workflow |

---

## Integration Gotchas

Common mistakes when connecting to external services.

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| Claude Code skill system | Putting task skills in `commands/` legacy format without frontmatter | Use `skills/<name>/SKILL.md` with proper frontmatter; `commands/` is deprecated |
| Claude Code agent definitions | Using `permissionMode` or `hooks` in plugin-shipped agent frontmatter | These fields are disallowed in plugin agents for security; use only `name`, `description`, `model`, `effort`, `maxTurns`, `tools`, `disallowedTools`, `skills`, `memory`, `background`, `isolation` |
| Claude Code skills context budget | Writing all instructions in one skill file | Keep `SKILL.md` under 500 lines; move reference material to supporting files loaded on demand |
| Claude Code plugin caching | Referencing files outside the plugin directory with `../` paths | Installed plugins are copied to `~/.claude/plugins/cache/`; path traversal is blocked; use symlinks if needed |
| GoReleaser + Homebrew | Creating formulae for a precompiled Go binary | Use `brews` config for casks, not formulae; formulae were disabled June 2025 |
| GoReleaser + GitHub Actions | Using default `GITHUB_TOKEN` to push to separate tap repository | Default token only has access to the current repo; create a dedicated PAT with content write access to the tap repo |
| Windows binary naming | Binary name contains "install", "setup", "patch", or "update" | Windows UAC interprets these as installer names and triggers elevation prompts; choose a name like `myssd` or `ssd` |
| OpenSpec Delta Specs (MODIFIED) | Using MODIFIED to add a new concern without including previous spec text | Previous text is lost at archive time; use ADDED for new concerns, MODIFIED only for changing existing requirement text |
| RFC 2119 keyword parsing | Case-insensitive matching of `must`/`should`/`may` in spec body text | Only UPPERCASE keywords are normative per RFC 2119; lowercase "must" is natural English, not a requirement; parser must be case-sensitive |

---

## Performance Traps

Patterns that work at small scale but fail as usage grows.

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Loading all spec files into AI context on every command | Slow invocations; context budget exceeded | Load only the spec(s) relevant to the current task/command; use Go binary to select relevant files | Projects with 10+ spec files |
| Synchronous spec file scanning at CLI startup | High startup latency for simple commands | Defer spec loading to commands that need it; keep startup path to config + arg parsing only | Projects with 50+ files in `.specs/` |
| Monolithic spec state file (single tasks.md for all tasks) | Concurrent read/write conflicts; merge conflicts in git history | Partition state by spec ID; one status file per spec or per milestone | 20+ active specs in parallel |
| Eager AI verification on every save | Runaway costs; slow feedback loop | Gate verification behind explicit `verify` command; never run AI verification automatically | Any size project |
| Building parallel multi-agent orchestration in Phase 1 | Over-engineering; delayed delivery; complex bugs | Sequential-first, parallel-optional later | Not a performance issue at small scale; an engineering complexity trap at any scale |

---

## Security Mistakes

Domain-specific security issues beyond general web security.

| Mistake | Risk | Prevention |
|---------|------|------------|
| Storing Anthropic API keys or tokens in `.specs/` or project config files committed to git | Credential leakage in version control; financial exposure | Never read API keys from project files; read from environment variables or OS keychain only |
| Plugin hooks executing arbitrary user-provided shell commands from spec files | Code injection via maliciously crafted spec files | Validate and sanitize any spec content before passing to shell commands; use `allowed-tools` to restrict hook capabilities |
| Plugin agent with `isolation: worktree` writing back to main branch without review | Unreviewed AI-generated code in production branch | Make worktree agents create PRs, never merge to main directly |
| Shipping a Go binary without code signing on macOS | Gatekeeper blocks execution; users get "unidentified developer" warning and bypass security prompt | Sign with Apple Developer ID; use GoReleaser's `sign` and `notarize` configuration; Apple charges yearly fee but is required for distribution |
| Exposing full project file tree to AI agents during execute phase | Over-broad context leaks sensitive files (`.env`, credentials) | Use Go binary to construct a minimal file list; pass only files relevant to the current spec task |

---

## UX Pitfalls

Common user experience mistakes in this domain.

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| Verbose spec output by default | Users spend most time reading long Markdown instead of acting — the primary SDD complaint | Default to summary view; use `--verbose` or `--full` flags for detailed spec output |
| Blocking permission prompts during long agentic runs | User goes to do something else, returns 10 minutes later to a stalled prompt | Use Claude Code's `allowedTools` in agent definitions to pre-authorize expected operations; gate only genuinely destructive actions |
| Sequential workflow with no progress visibility | Users have no idea if the system is working or stuck during AI execution | Stream structured progress output to stdout; log AI turn starts/ends with elapsed time |
| Making `verify` a separate manual step users can skip | Spec drift accumulates silently; verification is skipped 80% of the time in practice | Integrate verification as the mandatory gate before `archive`; fail `archive` if verification has not been run since last `execute` |
| Ambiguous command names that clash with OpenSpec (`propose`, `spec`, etc.) | Users who know OpenSpec expect identical behavior; differences cause confusion | Map command names to OpenSpec equivalents where possible; where behavior differs, document the difference explicitly in the command's help text |
| Convention-over-config that silently picks wrong defaults | Users assume the tool did the right thing; errors are invisible | Make "what defaults were applied" visible in output; log `Using .specs/ directory (auto-detected)` not silence |

---

## "Looks Done But Isn't" Checklist

Things that appear complete but are missing critical pieces.

- [ ] **Spec parser:** Verify it correctly handles OpenSpec Delta Specs (ADDED/MODIFIED/REMOVED) — the most common real-world spec pattern. Run against an actual OpenSpec project's `specs/` directory, not a synthetic test fixture.
- [ ] **RFC 2119 verification:** Verify the keyword extractor is CASE-SENSITIVE — `must` (lowercase) in prose is not a MUST requirement. A test with `"you must not"` vs `"you MUST NOT"` is the minimum.
- [ ] **Archive workflow:** Verify that `archive` fails when MUST items are unsatisfied — not just when verification has not been run. Test with a deliberate MUST failure.
- [ ] **Plugin skill budget:** Run `/context` in Claude Code with all plugin skills installed. Check for the "excluded skills" warning before any real use.
- [ ] **Go binary installation:** Verify the binary name does not contain Windows-triggering keywords (install/setup/patch/update). Test installation on Windows 11 without admin rights.
- [ ] **Cross-platform file paths:** Verify `.specs/` path handling uses `filepath.Join` not string concatenation. Test on Windows with a project in a directory containing spaces.
- [ ] **Brownfield compatibility:** Load an actual OpenSpec project (`openspec/` directory, not `.specs/`) and verify all commands work without manual migration.
- [ ] **Agent isolation:** Verify that the Go binary — not the AI agent — is the only writer to `.specs/` status files. Review all file write operations in agent skill definitions.
- [ ] **Homebrew distribution:** Confirm GoReleaser config uses `brews` as a cask, not a formula. Verify the tap repository uses a dedicated PAT, not `GITHUB_TOKEN`.
- [ ] **macOS code signing:** Verify the binary passes Gatekeeper without the "unidentified developer" bypass prompt. Test on a machine that has never opened the binary before.

---

## Recovery Strategies

When pitfalls occur despite prevention, how to recover.

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Spec format lock-in discovered after 3 phases | HIGH | Introduce schema version in frontmatter; write a migration command (`myssd migrate --from v1 --to v2`); parse both formats for one release cycle |
| Spec drift with diverged implementation | MEDIUM | Run `verify --forensic` to produce a full MUST/SHOULD status report; treat failures as new spec tasks; re-archive with current state documented |
| Context budget overflow causing silent skill exclusion | LOW | Split large skill files; add `disable-model-invocation: true` to task skills; re-test with `/context` |
| Plugin API breaking change from Claude Code update | MEDIUM | Maintain a thin plugin layer with all logic in Go binary; update plugin wrapper only; no business logic to rewrite |
| Binary blocked by Windows UAC elevation prompt | LOW | Rename binary to non-triggering name; rebuild and redistribute; update all distribution channels |
| GoReleaser Homebrew formula deprecated | LOW | Switch `brews` config from formula to cask; bump version; re-release; existing users need `brew upgrade` |
| Parallel agent shared state corruption | HIGH | Abort parallel run; restore from last clean git state; re-run in sequential mode; redesign with explicit coordinator ownership of state |

---

## Pitfall-to-Phase Mapping

How roadmap phases should address these pitfalls.

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| Spec format lock-in | Phase 1 (Core data model) | Integration test: load an unmodified OpenSpec project's `specs/` directory without errors |
| AI self-verification blindness | Phase 1 (Architecture decision) + Phase 3 (Verify engine) | Test: verification with deliberately planted MUST failure must return failure, not success |
| Spec drift | Phase 2 (Execute engine) + Phase 3 (Archive gate) | Test: archive command refuses when last verify predates last execute |
| Context budget overflow | Phase 1 (Plugin architecture) | Manual check: `/context` shows no excluded skills warning with all plugin skills installed |
| Multi-agent shared state | Phase 2 (Coordinator design) | Test: no `.specs/` writes in agent skill files; all state mutations go through Go binary |
| Claude Code API breaking changes | Phase 1 (Thin plugin layer design) | Measure: plugin skill files contain zero business logic; all decisions in Go binary |
| Homebrew distribution | Phase 4 (Distribution) | Test: `brew install` from tap on clean macOS machine without errors |
| Windows binary naming | Phase 4 (Distribution) | Test: install binary on Windows 11 without admin rights, no UAC prompt |
| macOS Gatekeeper | Phase 4 (Distribution) | Test: open binary on fresh macOS without right-click override |
| RFC 2119 case sensitivity | Phase 1 (Spec parser) | Unit test: `"you must"` (lowercase) extracts 0 requirements; `"you MUST"` extracts 1 |
| OpenSpec Delta Spec MODIFIED misuse | Phase 1 (Spec parser) + Docs | Integration test: round-trip ADDED/MODIFIED/REMOVED through archive and verify content is correct |

---

## Sources

- OpenSpec GitHub Issue #666: Hardcoded spec format prevents custom schemas — https://github.com/Fission-AI/OpenSpec/issues/666
- "Why Spec-Driven Development Fails" — https://dev.to/casamia918/why-spec-driven-development-fails-and-what-we-can-learn-from-it-2pec
- "Spec-Driven Development: The Waterfall Strikes Back" (marmelab, 2025) — https://marmelab.com/blog/2025/11/12/spec-driven-development-waterfall-strikes-back.html
- Claude Code Skills documentation (official, 2026) — https://code.claude.com/docs/en/slash-commands
- Claude Code Plugins Reference (official, 2026) — https://code.claude.com/docs/en/plugins-reference
- Claude Code Subagents: Common Mistakes — https://claudekit.cc/blog/vc-04-subagents-from-basic-to-deep-dive-i-misunderstood
- Context Management with Subagents in Claude Code — https://www.richsnapp.com/article/2025/10-05-context-management-with-subagents-in-claude-code
- Multi-Agent Orchestration: AI Agent Orchestration is Broken (builder.io) — https://www.builder.io/blog/ai-agent-orchestration
- How We Built True Parallel Agents With Git Worktrees — https://dev.to/getpochi/how-we-built-true-parallel-agents-with-git-worktrees-2580
- Git worktrees for parallel AI coding agents (Upsun) — https://devcenter.upsun.com/posts/git-worktrees-for-parallel-ai-coding-agents/
- Agentic Engineering Part 6: Forensic Verification — https://www.sagarmandal.com/2026/03/15/agentic-engineering-part-6-forensic-verification-why-a-perfect-score-from-your-ai-agent-should-make-you-nervous/
- GoReleaser Homebrew Casks (formulae deprecated June 2025) — https://goreleaser.com/customization/homebrew_casks/
- GoReleaser Homebrew formulae discussion #5563 — https://github.com/orgs/goreleaser/discussions/5563
- Go binary Windows elevation issue — https://github.com/golang/go/issues/8711
- OpenSpec Deep Dive guide — https://redreamality.com/garden/notes/openspec-guide/
- RFC 2119 (MUST/SHOULD/MAY specification) — https://datatracker.ietf.org/doc/html/rfc2119
- Multi-Agent Systems Orchestration 2026 — https://www.codebridge.tech/articles/mastering-multi-agent-orchestration-coordination-is-the-new-scale-frontier
- AI Agent Context Window Research (Redis) — https://redis.io/blog/ai-agent-orchestration/

---
*Pitfalls research for: AI-assisted Spec-Driven Development CLI (my-ssd)*
*Researched: 2026-03-23*
