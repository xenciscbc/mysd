# mysd

**Spec-Driven Development for AI Programming**

mysd is a Go CLI tool + Claude Code plugin that integrates [OpenSpec](https://github.com/openspec-dev/openspec)'s Spec-Driven Development (SDD) methodology with [GSD (Get Shit Done)](https://github.com/bfra-me/get-shit-done)'s planning/execution/verification engine into one seamless system.

It lets solo developers (1 human + N AI agents) drive AI coding with structured specs — ensuring AI reads and aligns with requirements before writing code, and automatically verifies results after execution.

## Why mysd?

OpenSpec gives you a complete SDD methodology but no execution engine. GSD gives you a powerful execution engine but no spec management. mysd fills the gap:

- **Specs are the single source of truth** — not just documentation, but what directly drives AI execution and verification
- **Alignment gate** — AI must read and acknowledge the spec before writing any code (non-bypassable)
- **Goal-backward verification** — an independent AI agent checks every MUST item against filesystem evidence
- **Convention over configuration** — works out of the box, configure only when needed
- **Wave parallel execution** — tasks with no dependency overlap run in parallel git worktrees, cutting execution time significantly
- **Interactive discovery** — 4-dimension research (Codebase / Domain / Architecture / Pitfalls) via advisor agents before committing to a spec
- **Self-update** — `mysd update` checks GitHub Releases, downloads the correct platform binary, and syncs plugin files automatically
- **Deferred notes** — scope guardrail that captures out-of-scope ideas without interrupting the current change

## Installation

```bash
go install github.com/xenciscbc/mysd@latest
```

Or download precompiled binaries from [GitHub Releases](https://github.com/xenciscbc/mysd/releases).

### Claude Code Plugin

```bash
claude plugin add --marketplace https://github.com/xenciscbc/mysd
```

Or manually copy the `plugin/` directory:

```bash
cp -r plugin/ ~/.claude/plugins/mysd/
```

All `/mysd:*` slash commands will be available in your next Claude Code session.

## Quick Start

```bash
# Initialize mysd in your project
/mysd:init

# Create a change proposal — auto-generates spec, design, and tasks
/mysd:propose add-user-auth

# (Optional) Explore before committing — 4-dimension interactive research
/mysd:discuss

# Break design into executable tasks (includes plan-checker MUST coverage verification)
/mysd:plan

# Execute tasks — AI reads spec first, then codes
# Verification is mandatory and runs automatically after execution
/mysd:apply

# Archive completed change
/mysd:archive
```

Or fast-forward the entire workflow:

```bash
# propose -> plan in one shot (with auto research)
/mysd:ff my-feature

# propose -> plan -> apply -> verify -> archive in one shot (fully automated)
/mysd:ffe my-feature
```

## Commands

### Workflow

| Command | Description |
|---------|-------------|
| `/mysd:propose` | Create a change proposal with auto-generated spec, design, and tasks |
| `/mysd:discuss` | Interactive exploration with 4-dimension research and advisor agents |
| `/mysd:plan` | Break design into executable tasks with plan-checker MUST coverage verification |
| `/mysd:apply` | Execute tasks with mandatory spec alignment gate and built-in verification; supports wave parallel execution |
| `/mysd:archive` | Archive completed change to `openspec/changes/archive/` |

### Fast-Forward

| Command | Description |
|---------|-------------|
| `/mysd:ff` | Fast-forward: propose through plan (with auto research) |
| `/mysd:ffe` | Full fast-forward: propose → plan → apply → verify → archive (fully automated) |

### Utility

| Command | Description |
|---------|-------------|
| `/mysd:status` | Show current workflow state and progress |
| `/mysd:scan` | Scan existing codebase and generate specs |
| `/mysd:fix` | Fix failed tasks in worktree isolation with optional research mode |
| `/mysd:note` | Manage deferred notes — capture out-of-scope ideas without interrupting work |
| `/mysd:model` | View or set model profile (quality / balanced / budget) |
| `/mysd:lang` | Set response language and OpenSpec locale |
| `/mysd:update` | Check for updates and install new mysd binary + sync plugin files |
| `/mysd:init` | Initialize project configuration |
| `/mysd:uat` | Interactive user acceptance testing |

## How It Works

mysd follows a structured lifecycle for every code change:

```
propose -> [discuss] -> plan -> apply (with verification) -> archive
```

1. **Propose** — Create a change proposal and auto-generate all artifacts (proposal.md, specs/, design.md, tasks.md) in `openspec/changes/`
2. **Discuss** *(optional)* — Run 4-dimension interactive research (Codebase/Domain/Architecture/Pitfalls) to explore unknowns before committing to requirements
3. **Plan** — Break the design into ordered, executable tasks; plan-checker verifies every MUST item has a corresponding task
4. **Apply** — AI reads the spec (alignment gate), then implements each task. Verification is mandatory and runs automatically after execution. Tasks with no file overlap are grouped into waves and run in parallel git worktrees.
5. **Archive** — Move completed change to `openspec/changes/archive/` with delta spec sync

### Key Architecture

- **Go binary** handles state management, spec parsing, config, and structured JSON output
- **SKILL.md files** orchestrate the AI workflow (invoke binary, present results, delegate to agents)
- **Agent definitions** perform the actual AI work (spec writing, execution, verification)
- **Reverse-calling pattern** — Claude Code calls the binary, not the other way around. No MCP server needed.

### Wave Parallel Execution

When a spec's tasks.md contains `depends` fields, mysd performs dependency analysis:

- Tasks are topologically sorted into waves — each wave contains tasks with no inter-dependencies
- Same-wave tasks with no file overlap run in parallel git worktrees, cutting execution time
- After all worktrees complete, branches are merged in task-ID order (`--no-ff`)
- AI conflict resolution with 3 retries (`go build` + test); failed worktrees are preserved for `/mysd:fix`
- Mode selection: `sequential` (safe, default) or `wave` (parallel). `auto_mode` skips the mode prompt.

### Interactive Discovery

`/mysd:discuss` runs a dual-loop research session before you commit to writing requirements:

- **4 research dimensions**: Codebase (existing patterns), Domain (requirements & constraints), Architecture (approach options), Pitfalls (known failure modes)
- Advisor agents surface unknowns per dimension; you terminate each loop when you have enough context
- Output feeds directly into `/mysd:spec` so nothing is lost

## Self-Update

mysd can update itself from GitHub Releases:

```bash
/mysd:update           # Check for updates and install interactively
mysd update --check    # Check only (JSON output, no install)
mysd update --force    # Update without confirmation
```

Updates include:

- Binary replacement with SHA256 checksum verification
- Automatic rollback on failure
- Plugin file sync (commands + agents) via manifest diff — only changed files are written

## OpenSpec Compatibility

mysd reads and writes [OpenSpec](https://github.com/openspec-dev/openspec) format natively:

- Supports both `.specs/` and `openspec/` directory structures
- Parses YAML frontmatter with `spec-version` field
- Handles RFC 2119 keywords (MUST, SHOULD, MAY) and Delta Specs (ADDED, MODIFIED, REMOVED)
- Point mysd at an existing OpenSpec project and run `/mysd:apply` or `/mysd:verify` without migration

## Model Profiles

mysd uses a profile system to control which AI model each agent role uses. Thinking-heavy roles (spec writing, planning, verification) use stronger models; execution roles use efficient ones.

```bash
/mysd:model              # View current profile and role-to-model mapping
/mysd:model set quality  # Switch profile
```

Three profiles are available:

### quality (8 opus / 2 sonnet)

Maximum capability. All thinking roles use opus.

| Role | Model | Purpose |
|------|-------|---------|
| spec-writer | opus | Write requirements specs |
| designer | opus | Architecture and technical design |
| planner | opus | Task breakdown and dependency analysis |
| executor | sonnet | Implement tasks from plan |
| verifier | opus | Verify spec satisfaction |
| fast-forward | sonnet | Orchestrate accelerated workflows |
| researcher | opus | 4-dimension codebase research |
| advisor | opus | Trade-off analysis for gray areas |
| proposal-writer | opus | Write change proposals |
| plan-checker | opus | Verify plan covers all MUST items |

### balanced (6 opus / 4 sonnet) — default

Opus for judgment/design/gating roles, sonnet for execution and research.

| Role | Model |
|------|-------|
| spec-writer | opus |
| designer | opus |
| planner | opus |
| executor | sonnet |
| verifier | opus |
| fast-forward | sonnet |
| researcher | sonnet |
| advisor | opus |
| proposal-writer | sonnet |
| plan-checker | opus |

### budget (7 sonnet / 3 haiku)

Minimize cost. Spec-writer uses sonnet as the quality floor.

| Role | Model |
|------|-------|
| spec-writer | sonnet |
| designer | haiku |
| planner | sonnet |
| executor | haiku |
| verifier | sonnet |
| fast-forward | haiku |
| researcher | sonnet |
| advisor | sonnet |
| proposal-writer | sonnet |
| plan-checker | sonnet |

Standalone commands use fixed models regardless of profile: `init`, `scan`, `fix` always use opus; `status`, `lang`, `model`, `note`, `docs`, `update` always use sonnet.

## Configuration

Project config lives in `.claude/mysd.yaml`:

```yaml
execution_mode: single      # single | wave
agent_count: 1
atomic_commits: false
tdd_mode: false
model_profile: balanced     # quality | balanced | budget
response_language: en       # BCP 47 language tag, e.g. zh-TW, ja, fr
```

- `execution_mode: wave` enables parallel worktree execution for tasks with no dependency overlap
- `model_profile` controls AI model selection across all agent roles — see [Model Profiles](#model-profiles) for the full mapping
- `response_language` sets the language for all agent responses and OpenSpec locale

All options can be overridden per-command via flags.

## Tech Stack

- **Go 1.23+** — single binary, zero runtime dependencies
- **Cobra** — CLI framework
- **Viper** — configuration management
- **lipgloss** — terminal output styling
- **yaml.v3** — YAML parsing for OpenSpec frontmatter
- **GoReleaser** — cross-platform binary distribution

## License

MIT
