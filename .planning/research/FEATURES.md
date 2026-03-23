# Feature Research

**Domain:** AI-assisted Spec-Driven Development (SDD) CLI tooling — Claude Code plugin
**Researched:** 2026-03-23
**Confidence:** HIGH (primary sources: OpenSpec, GSD repos, GitHub Spec Kit, BMAD-METHOD docs, intent-driven.dev)

---

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist. Missing these = product feels incomplete.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **Structured spec artifacts** (proposal + specs/ + design + tasks) | Every SDD tool (OpenSpec, GitHub Spec Kit, BMAD) organizes changes this way — users trained to expect this folder structure | MEDIUM | Must map to OpenSpec's `proposal.md / specs/ / design.md / tasks.md`; my-ssd stores under `.specs/` |
| **Linear workflow commands** (propose → spec → design → plan → execute → verify → archive) | Users of OpenSpec, GitHub Spec Kit, BMAD all expect a named command per phase — the phase gate is the UX | MEDIUM | GSD maps to new-project/plan-phase/execute-phase/verify-work; OpenSpec maps to propose/apply/archive |
| **Single-agent execution** (default) | Most users start simple; multi-agent is a power feature — if default is complex, users abandon | LOW | Convention over config: run sequentially unless parallelism is explicitly triggered |
| **Spec as source of truth** | The core claim of every SDD tool; AI must read and acknowledge spec before writing code | LOW | Enforced via pre-execution alignment step — not optional |
| **Claude Code slash commands** | This is a Claude Code plugin; slash commands ARE the UX — no slash commands = no plugin | LOW | `/mysd:propose`, `/mysd:execute`, etc. — the primary interface |
| **Spec status tracking** (PENDING / IN_PROGRESS / DONE / BLOCKED) | Users need to know what state a spec is in across sessions | LOW | Stored in spec metadata header; machine-readable |
| **Archive / history** | Completed specs must be retrievable; OpenSpec has `/opsx:archive` and bulk-archive | LOW | Move completed specs to `.specs/archive/`; retain full artifact folder |
| **RFC 2119 keyword support** (MUST / SHOULD / MAY) | Required for goal-backward verification — verification logic parses MUST items | MEDIUM | Parser must distinguish MUST (required) from SHOULD (preferred) from MAY (optional) |
| **Brownfield support** (run on existing projects) | Most real projects are not greenfield; OpenSpec explicitly calls this out as a design goal | MEDIUM | `/mysd:onboard` equivalent — reads existing code, generates baseline context without reverse-engineering full specs |
| **Session continuity** (cross-session state) | Long projects span days; tool must recover current position without user explaining history | MEDIUM | GSD's `STATE.md` + `HANDOFF.json` pattern; at minimum a `.specs/STATE.md` |
| **Convention over configuration** | Power users want zero-config defaults; forcing config before first run kills adoption | LOW | Sensible defaults for all paths, naming, and behavior — config optional |

---

### Differentiators (Competitive Advantage)

Features that set the product apart. Not required, but valued.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Goal-backward verification** | Unlike OpenSpec (no built-in verifier) and GitHub Spec Kit (no automated post-execution check), my-ssd verifies that every MUST item in the spec was satisfied after execution — closes the feedback loop | HIGH | Parses spec for MUST items → generates checklist → runs verification agent → updates spec status automatically |
| **Spec feedback loop** (verification results written back to spec) | Spec is a living document, not a one-shot prompt; failed MUST items stay open, passed items get marked done — no other tool does this automatically | HIGH | Requires bidirectional spec I/O: read before execution, write after verification |
| **Delta Specs** (ADDED / MODIFIED / REMOVED semantics) | Uniquely OpenSpec-derived — lets AI understand the *change type* not just the requirement, enabling smarter diffs and rollback reasoning | MEDIUM | Parse delta markers; use them to scope execution and verification narrowly |
| **Single Go binary, no runtime dependency** | Every competing tool (OpenSpec, GSD, BMAD) requires Node.js; Go binary = `curl | install` UX; zero friction for non-Node projects | MEDIUM | Cross-compile for macOS/Linux/Windows; distribute via GitHub releases |
| **Tight spec-execution coupling** (pre-execution alignment gate) | Tools like OpenSpec and GitHub Spec Kit are spec-first but leave execution to ad-hoc AI prompting; my-ssd makes the alignment step non-bypassable | LOW | Before any execute command, agent must confirm it has read and understood the spec — structured checkpoint |
| **Multi-agent wave execution** (opt-in parallel) | GSD has this but OpenSpec doesn't; for complex tasks, parallel execution cuts wall-clock time significantly | HIGH | Dependency analysis → wave grouping → dispatch parallel agents with isolated contexts; default OFF |
| **OpenSpec format compatibility** (zero-migration path) | Existing OpenSpec users can point my-ssd at their `openspec/` dir and get execution + verification for free — no re-tooling | MEDIUM | Format parser must handle both OpenSpec's directory layout and my-ssd's `.specs/` layout |
| **Brownfield codebase onboarding** (map-then-spec) | `/mysd:onboard` generates `STACK.md`, `ARCHITECTURE.md`, `CONVENTIONS.md` from existing code before any spec is written — prevents AI from ignoring existing patterns | HIGH | Run static analysis + AI summarization on existing codebase; output feeds all subsequent spec generation |
| **Atomic git commits per task** | GSD does this; OpenSpec doesn't enforce it; traceable task→commit mapping makes bisect and rollback practical | LOW | Each executed task triggers `git commit` with structured message; requires git in PATH |

---

### Anti-Features (Commonly Requested, Often Problematic)

Features that seem good but create problems.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| **GUI / web dashboard** | Visualizing spec status and history looks appealing | Massive scope expansion; the target user is a developer comfortable with CLI and file-based workflows; adds maintenance burden for a non-CLI surface | Rich terminal output (`/mysd:status`) with color-coded phase progress; let files be the UI |
| **Full reverse-engineering of entire codebase into specs** | Users want "generate specs for everything I have" on brownfield projects | AI-generated specs from existing code are often inaccurate (intent-driven.dev research confirms this); specs that diverge from reality are worse than no specs — false confidence | Incremental spec authoring: write specs for upcoming changes only; use `/mysd:onboard` only for conventions/stack, not for full reverse-engineering |
| **Multi-AI-tool support** (Cursor, Copilot, Gemini CLI) | Broader compatibility sounds like more value | Each tool has different context injection mechanics; supporting N tools in v1 multiplies complexity without validating the core spec-execution loop | v1 is Claude Code only; abstract the plugin interface so other runtimes can be added after core is proven |
| **Team collaboration features** (shared spec review, multi-user approval) | Looks like enterprise value | Solo developer + AI is the validated use case; multi-user workflows require auth, concurrency control, conflict resolution — entirely different product | Version control (git) handles multi-person collaboration on spec files; the tool doesn't need to own it |
| **Real-time spec sync** (auto-update specs as code changes) | Spec drift is real; auto-sync seems like the fix | Specs describe intent; code describes implementation; auto-sync from code back to specs inverts the causality — the spec should drive code, not vice versa | After execution, run verification to mark MUST items as satisfied or failed; never auto-rewrite spec intent from code output |
| **57-command surface** (GSD-style comprehensive command set) | More commands = more power | Every additional command is surface area to maintain, document, and keep consistent; GSD's breadth is also its onboarding problem | Minimal command set covering the core loop: `propose → spec → design → plan → execute → verify → archive`; add commands only when validated need exists |
| **Configuration-heavy setup** | Power users want to customize everything | Forces decisions before any value is delivered; kills the "try it in 5 minutes" experience | Convention over configuration: all defaults work out of the box; `mysd.yaml` exists but is never required |

---

## Feature Dependencies

```
[Spec Artifacts (proposal/specs/design/tasks)]
    └──required by──> [Workflow Commands]
                           └──required by──> [Pre-Execution Alignment Gate]
                                                  └──required by──> [Goal-Backward Verification]
                                                                         └──required by──> [Spec Feedback Loop]

[RFC 2119 Keyword Support]
    └──required by──> [Goal-Backward Verification]
                           └──required by──> [Spec Feedback Loop]

[Delta Specs (ADDED/MODIFIED/REMOVED)]
    └──enhances──> [Goal-Backward Verification]  (scope verification to the change type)
    └──enhances──> [Spec Feedback Loop]          (mark delta items individually)

[Brownfield Onboarding (map-then-spec)]
    └──required by──> [OpenSpec Format Compatibility]  (must understand existing layout)
    └──enhances──> [Workflow Commands]                 (context feeds spec generation quality)

[Session Continuity (STATE.md)]
    └──required by──> [Multi-Agent Wave Execution]  (wave state must survive restarts)
    └──enhances──> [Workflow Commands]              (resume mid-workflow)

[Multi-Agent Wave Execution]
    └──requires──> [Session Continuity]
    └──requires──> [Atomic Git Commits]   (rollback must be per-task, not per-wave)
    └──conflicts with──> [Single Go Binary constraint]  (spawning sub-agents requires Claude Code's agent API, not a Go subprocess — integration boundary must be clear)

[Atomic Git Commits]
    └──enhances──> [Multi-Agent Wave Execution]
    └──independent of──> [Goal-Backward Verification]  (can verify without commits, but commits help audit trail)

[Single Go Binary]
    └──conflicts with──> [Node.js ecosystem tooling]  (cannot use npm packages; MCP/plugin integration is via file I/O and slash commands, not Node modules)
```

### Dependency Notes

- **Goal-Backward Verification requires RFC 2119 Keyword Support:** The verifier needs to know which items are MUST (non-negotiable) vs SHOULD (preferred) to build the correct checklist; without this parsing, verification is guesswork.
- **Spec Feedback Loop requires Goal-Backward Verification:** Writing results back to the spec is only meaningful after verification has produced structured pass/fail outcomes; the feedback loop is the output of verification, not a separate feature.
- **Multi-Agent Wave Execution conflicts with Single Go Binary at the agent-spawn boundary:** Go binary manages orchestration state and file I/O; actual agent spawning uses Claude Code's native subagent API (`.claude/agents/`); the Go binary cannot spawn Claude agents directly — it prepares context and dispatches instructions via Claude Code's plugin mechanism.
- **Brownfield Onboarding enhances all Workflow Commands:** Once `CONVENTIONS.md` and `ARCHITECTURE.md` are generated by onboarding, all subsequent spec-generation steps use them as grounding context, preventing the AI from inventing patterns that contradict the existing codebase.
- **OpenSpec Format Compatibility requires Brownfield Onboarding logic:** Recognizing and reading an existing OpenSpec `openspec/` directory is a specialization of the brownfield problem — same parser, different entry path.

---

## MVP Definition

### Launch With (v1)

Minimum viable product — what's needed to validate that spec-to-execution-to-verification is tighter than using OpenSpec + GSD separately.

- [ ] **Spec Artifacts** (proposal.md, specs/, design.md, tasks.md) — the container for all work; without this, there is nothing to drive execution from
- [ ] **Core Workflow Commands** (`/mysd:propose`, `/mysd:execute`, `/mysd:verify`, `/mysd:archive`) — the minimum loop; propose and execute are necessary to produce output; verify is the entire differentiating value; archive closes the loop
- [ ] **Pre-Execution Alignment Gate** — forces AI to read and acknowledge the spec before writing any code; this is the single most important behavior change over plain AI coding
- [ ] **RFC 2119 Keyword Parsing** (MUST / SHOULD / MAY) — required for goal-backward verification; the parser is simple but must be correct
- [ ] **Goal-Backward Verification** — generates checklist from spec MUST items, runs post-execution check, marks items pass/fail; this is the core differentiator vs OpenSpec
- [ ] **Spec Feedback Loop** — writes verification results back into spec status; makes spec a living document rather than a one-shot artifact
- [ ] **Session Continuity** (STATE.md) — without this, multi-session projects lose progress; low complexity but high user value
- [ ] **Archive / History** — completed specs move to `.specs/archive/`; essential for project cleanliness
- [ ] **Single Go Binary** — the deployment story; without this, installation friction is identical to OpenSpec/GSD
- [ ] **OpenSpec Format Compatibility** — required for brownfield adoption by existing OpenSpec users; parser for existing `openspec/` layouts

### Add After Validation (v1.x)

Features to add once core loop is working and users are reporting specific pain points.

- [ ] **Delta Specs** (ADDED / MODIFIED / REMOVED) — add when users report that verification scope is too broad for incremental changes; MEDIUM complexity
- [ ] **Brownfield Codebase Onboarding** (`/mysd:onboard`) — add when users report that AI is ignoring existing patterns during spec execution; HIGH complexity
- [ ] **Atomic Git Commits per Task** — add when users report difficulty bisecting failures; LOW complexity once execution scaffolding is in place
- [ ] **`/mysd:design` and `/mysd:plan` commands** — expand the workflow from the 4-command MVP to the full 7-command suite once the core loop is validated

### Future Consideration (v2+)

Features to defer until product-market fit is established.

- [ ] **Multi-Agent Wave Execution** — defer until single-agent execution is reliable; wave orchestration adds significant complexity; validate that users hit single-agent performance limits before building this
- [ ] **Multi-Runtime Support** (OpenCode, Gemini CLI) — defer until Claude Code integration is proven; abstracting the plugin interface should be designed for early but not implemented until v1 is stable
- [ ] **Spec Templates / Profiles** — defer until patterns in user-generated specs reveal what should be standardized

---

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Spec Artifacts (4-artifact structure) | HIGH | MEDIUM | P1 |
| Core Workflow Commands (propose/execute/verify/archive) | HIGH | MEDIUM | P1 |
| Pre-Execution Alignment Gate | HIGH | LOW | P1 |
| RFC 2119 Keyword Parsing | HIGH | LOW | P1 |
| Goal-Backward Verification | HIGH | HIGH | P1 |
| Spec Feedback Loop | HIGH | MEDIUM | P1 |
| Session Continuity (STATE.md) | HIGH | LOW | P1 |
| Archive / History | MEDIUM | LOW | P1 |
| Single Go Binary | HIGH | MEDIUM | P1 |
| OpenSpec Format Compatibility | HIGH | MEDIUM | P1 |
| Delta Specs | MEDIUM | MEDIUM | P2 |
| Atomic Git Commits per Task | MEDIUM | LOW | P2 |
| Brownfield Codebase Onboarding | HIGH | HIGH | P2 |
| `/mysd:design` + `/mysd:plan` commands | MEDIUM | LOW | P2 |
| Multi-Agent Wave Execution | MEDIUM | HIGH | P3 |
| Multi-Runtime Support | LOW | HIGH | P3 |
| Spec Templates / Profiles | LOW | MEDIUM | P3 |

**Priority key:**
- P1: Must have for launch
- P2: Should have, add when possible
- P3: Nice to have, future consideration

---

## Competitor Feature Analysis

| Feature | OpenSpec | GSD | GitHub Spec Kit | BMAD-METHOD | my-ssd Approach |
|---------|----------|-----|-----------------|-------------|-----------------|
| Spec artifacts (4-artifact structure) | YES (proposal/specs/design/tasks) | NO (planning files only) | YES (spec/plan/tasks) | YES (PRD/arch/stories) | YES — full 4-artifact, OpenSpec-compatible |
| Workflow commands (named phases) | YES (propose/apply/archive) | YES (57 commands) | YES (/specify/plan/tasks/implement) | YES (agent personas) | YES — minimal 7-command set |
| Pre-execution alignment gate | Partial (AI reads spec, but not enforced) | NO — execution proceeds from plan, not spec | NO — task list drives execution | Partial — agent personas have role constraints | YES — mandatory, non-bypassable |
| Goal-backward verification | NO | YES (Nyquist Layer + goal-backward planning) | NO | Partial (QA agent persona) | YES — spec MUST items → verification checklist |
| Spec feedback loop (write results back) | NO | NO | NO | NO | YES — differentiator |
| RFC 2119 keyword support | NO (plain language specs) | NO | NO | NO | YES — MUST/SHOULD/MAY parsed and acted on |
| Delta Specs | NO | NO | NO | NO | YES (ported from OpenSpec design) |
| Brownfield support | YES (/opsx:onboard, custom profiles) | YES (/gsd:map-codebase) | Partial (incremental spec authoring) | YES (explicit brownfield guide) | YES — /mysd:onboard generates conventions context |
| Session continuity | Partial (archive preserves history) | YES (STATE.md, HANDOFF.json) | NO | Partial (git-versioned artifacts) | YES — STATE.md pattern from GSD |
| Multi-agent orchestration | NO | YES (wave execution, parallel agents) | NO | YES (12+ specialized agent personas) | YES — opt-in wave execution (v1.x+) |
| Archive / history | YES (/opsx:archive, bulk-archive) | YES (/gsd:complete-milestone) | NO | YES (git-versioned artifacts) | YES — archive to `.specs/archive/` |
| Zero runtime dependency | NO (Node.js required) | NO (Node.js required) | NO (Node.js required) | NO (Node.js required) | YES — single Go binary |
| Atomic git commits per task | NO | YES | NO | Partial (git-versioned artifacts) | YES (v1.x) |

---

## Sources

- [OpenSpec GitHub Repository](https://github.com/Fission-AI/OpenSpec) — PRIMARY: workflow commands, artifact structure, brownfield support, philosophy
- [GSD (get-shit-done) GitHub Repository](https://github.com/gsd-build/get-shit-done) — PRIMARY: wave execution, multi-agent orchestration, goal-backward verification, STATE.md, context engineering
- [GSD User Guide](https://github.com/gsd-build/get-shit-done/blob/main/docs/USER-GUIDE.md) — wave execution mechanics, Nyquist Layer, verification workflows
- [GitHub Spec Kit announcement](https://github.blog/ai-and-ml/generative-ai/spec-driven-development-with-ai-get-started-with-a-new-open-source-toolkit/) — spec/plan/tasks/implement workflow, multi-agent support, living artifacts
- [BMAD-METHOD Documentation](https://docs.bmad-method.org/) — specialized agent architecture, PRD/architecture/stories, brownfield guide, adversarial review
- [SDD Brownfield Guide - intent-driven.dev](https://intent-driven.dev/blog/2026/03/10/spec-driven-development-brownfield/) — brownfield challenges: AI-generated spec inaccuracy, spec drift, incremental authoring recommendation
- [Spec-Driven Development Is Eating Software Engineering - Medium](https://medium.com/@visrow/spec-driven-development-is-eating-software-engineering-a-map-of-30-agentic-coding-frameworks-6ac0b5e2b484) — 30+ framework landscape map, living vs static spec platforms
- [SDD 2026: Future of AI Coding or Waterfall? - alexcloudstar.com](https://www.alexcloudstar.com/blog/spec-driven-development-2026/) — anti-features: over-specification, spec drift, agent non-compliance pitfalls
- [awesome-claude-code](https://github.com/hesreallyhim/awesome-claude-code) — Claude Code plugin ecosystem patterns: skills, slash commands, hooks, context engineering
- [Stop Context Rot: GSD Explained](https://hoangyell.com/get-shit-done-explained/) — context window management, atomic commits, XML task format details
- [RFC 2119](https://datatracker.ietf.org/doc/html/rfc2119) — MUST/SHOULD/MAY definitions and usage guidelines

---
*Feature research for: AI-assisted Spec-Driven Development CLI (my-ssd)*
*Researched: 2026-03-23*
