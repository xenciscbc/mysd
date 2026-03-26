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

# Create a new spec from a feature description
/mysd:propose add-user-auth

# (Optional) Explore before spec writing — 4-dimension interactive research
/mysd:discuss

# Define detailed requirements with RFC 2119 keywords
/mysd:spec

# Capture technical decisions
/mysd:design

# Break design into executable tasks (includes plan-checker MUST coverage verification)
/mysd:plan

# Execute tasks — AI reads spec first, then codes
# Tasks with no dependency overlap run in parallel waves automatically
/mysd:apply

# Verify all MUST items are satisfied
/mysd:verify

# Archive completed spec
/mysd:archive
```

Or fast-forward the entire workflow:

```bash
# propose -> spec -> design -> plan in one shot (with auto research)
/mysd:ff my-feature

# propose -> spec -> design -> plan -> execute in one shot (fully automated)
/mysd:ffe my-feature
```

## Commands

### Spec Workflow

| Command | Description |
|---------|-------------|
| `/mysd:propose` | Create a new spec change with proposal scaffolding |
| `/mysd:discuss` | Interactive exploration with 4-dimension research and advisor agents |
| `/mysd:spec` | Define requirements with RFC 2119 keywords (MUST/SHOULD/MAY) |
| `/mysd:design` | Capture technical decisions and architecture |
| `/mysd:plan` | Break design into executable tasks with plan-checker MUST coverage verification |
| `/mysd:apply` | Run tasks with mandatory alignment gate; supports wave parallel execution |
| `/mysd:verify` | Goal-backward verification of all MUST items |
| `/mysd:archive` | Archive verified spec to `.specs/archive/` |

### Utility

| Command | Description |
|---------|-------------|
| `/mysd:status` | Show current workflow state and progress |
| `/mysd:scan` | Scan existing codebase and generate specs |
| `/mysd:capture` | Extract changes from current conversation |
| `/mysd:fix` | Fix failed tasks in worktree isolation with optional research mode |
| `/mysd:note` | Manage deferred notes — capture out-of-scope ideas without interrupting work |
| `/mysd:model` | View or set model profile (quality / balanced / budget) |
| `/mysd:lang` | Set response language and OpenSpec locale |
| `/mysd:update` | Check for updates and install new mysd binary + sync plugin files |
| `/mysd:init` | Initialize project configuration |
| `/mysd:uat` | Interactive user acceptance testing |

### Fast-Forward

| Command | Description |
|---------|-------------|
| `/mysd:ff` | Fast-forward: propose through plan (with auto research) |
| `/mysd:ffe` | Fast-forward: propose through execute (with auto research, auto execution) |

## How It Works

mysd follows a structured lifecycle for every code change:

```
propose -> [discuss] -> spec -> design -> plan -> apply -> verify -> archive
```

1. **Propose** — Scaffold spec artifacts (proposal.md, specs/, design.md, tasks.md) in `.specs/changes/`
2. **Discuss** *(optional)* — Run 4-dimension interactive research (Codebase/Domain/Architecture/Pitfalls) to explore unknowns before committing to requirements
3. **Spec** — Define requirements using RFC 2119 keywords. MUST items become verification gates.
4. **Design** — Record architecture decisions and technical approach
5. **Plan** — Break the design into ordered, executable tasks; plan-checker verifies every MUST item has a corresponding task
6. **Apply** — AI reads the spec (alignment gate), then implements each task. Tasks with no file overlap are grouped into waves and run in parallel git worktrees.
7. **Verify** — Independent verifier agent checks every MUST item against the actual codebase
8. **Archive** — Move completed spec to `.specs/archive/` (blocked if any MUST item fails)

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
- `model_profile` controls AI model selection across all agent roles (executor, verifier, researcher, etc.)
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
