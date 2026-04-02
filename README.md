# mysd

> **Testing / 測試中** — This project is under active development. APIs and workflows may change.

**Spec-Driven Development for AI Programming**

mysd is a Go CLI tool + Claude Code plugin that integrates [OpenSpec](https://github.com/openspec-dev/openspec)'s Spec-Driven Development (SDD) methodology with a planning/execution/verification engine into one seamless system.

It lets solo developers (1 human + N AI agents) drive AI coding with structured specs — ensuring AI reads and aligns with requirements before writing code, and automatically verifies results after execution.

## Why mysd?

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

## Usage Modes

mysd supports three usage modes — pick whichever fits your situation.

### Step-by-Step (Interactive)

Full control over each phase. Best for complex features or when you want to review between steps.

```bash
/mysd:init                        # One-time project setup
/mysd:propose add-user-auth       # Create change proposal + artifacts
/mysd:discuss                     # (Optional) 4-dimension research
/mysd:plan                        # Break design into tasks
/mysd:apply                       # Execute tasks (alignment gate + verification)
/mysd:archive                     # Archive completed change
```

### Fast-Forward (Semi-Automated)

Skip interactive confirmations. Requires an active change (run `/mysd:propose` first).

```bash
/mysd:propose my-feature   # Create change first
/mysd:ff my-feature        # plan → apply → archive (no research, auto mode)
```

### Full Fast-Forward (Fully Automated)

Same as fast-forward but with a research phase before planning.

```bash
/mysd:propose my-feature   # Create change first
/mysd:ffe my-feature       # research → plan → apply → archive (auto mode)
```

### Execution-Only Mode

Already have specs or design docs from another tool? Use `/mysd:plan` to convert them into executable tasks, then run:

```bash
# Convert external docs into tasks
/mysd:plan --from design.md          # Load external file as planner context
/mysd:plan --spec auth --from notes  # Plan a specific spec with external input

# If you already have tasks.md, jump straight to execution
/mysd:apply             # Execute pending tasks from existing tasks.md
/mysd:verify            # Verify MUST items independently
/mysd:archive           # Archive when done
```

## Commands

### Core Workflow

| Command | Description | Arguments |
|---------|-------------|-----------|
| `/mysd:propose` | Create a change proposal with spec, design, and tasks | `[change-name\|file\|dir] [--auto]` |
| `/mysd:discuss` | Ad-hoc research with 4-dimension exploration and advisor agents | `[topic\|change-name\|file\|dir] [--auto]` |
| `/mysd:plan` | Break design into executable tasks with MUST coverage check | `[--research] [--check] [--spec <name>] [--from <file>] [--auto]` |
| `/mysd:apply` | Execute tasks with spec alignment gate; supports single/wave/spec modes | `[--auto]` |
| `/mysd:verify` | Goal-backward verification of all MUST items by independent verifier | |
| `/mysd:archive` | Archive completed change to `openspec/changes/archive/` with delta spec sync | `[--auto]` |

### Fast-Forward

| Command | Description | Arguments |
|---------|-------------|-----------|
| `/mysd:ff` | Fast-forward: plan → apply → archive (assumes spec ready, no research, auto mode) | `[change-name]` |
| `/mysd:ffe` | Full fast-forward: research → plan → apply → archive (with research, auto mode) | `[change-name]` |

### Documentation

| Command | Description | Arguments |
|---------|-------------|-----------|
| `/mysd:docs` | Manage `docs_to_update` list (files auto-updated after archive) | `[add <path> \| remove <path>]` |
| `/mysd:docs-update` | Trigger doc updates independently — supports multiple scopes | `[--change <name> \| --last N \| --full \| "text"]` |

`/mysd:docs-update` scopes:
- **No arguments** — update from the most recent archived change
- `--change <name>` — update from a specific archived change
- `--last N` — update from the N most recent archived changes
- `--full` — scan the codebase and update docs to reflect actual project state
- `"free text"` — use the provided description as update context

### Utility

| Command | Description | Arguments |
|---------|-------------|-----------|
| `/mysd:status` | Show workflow state, task progress, and next step recommendation | |
| `/mysd:scan` | Scan existing codebase and generate OpenSpec-format specs | |
| `/mysd:fix` | Fix failed tasks in worktree isolation with optional research | `[change-name] [T{id}]` |
| `/mysd:note` | Manage deferred notes — capture out-of-scope ideas | `[add {content} \| delete {id}]` |
| `/mysd:model` | View or set model profile (quality / balanced / budget) | |
| `/mysd:lang` | Set response language and OpenSpec locale | |
| `/mysd:update` | Check for updates and install new binary + sync plugin files | `[--check] [--force]` |
| `/mysd:init` | Initialize project configuration and openspec structure | |
| `/mysd:uat` | Interactive user acceptance testing walkthrough | |
| `/mysd:statusline` | Toggle statusline display | `[on\|off]` |

## How It Works

### Lifecycle

Every code change follows a structured lifecycle:

```
propose → [discuss] → plan → apply (with verification) → archive
```

1. **Propose** — Create a change proposal and auto-generate all artifacts (proposal.md, specs/, design.md, tasks.md) in `openspec/changes/<name>/`
2. **Discuss** *(optional)* — Run 4-dimension interactive research (Codebase/Domain/Architecture/Pitfalls) to explore unknowns before committing to requirements
3. **Plan** — Break the design into ordered, executable tasks; plan-checker verifies every MUST item has a corresponding task
4. **Apply** — AI reads the spec (alignment gate), then implements each task. Verification runs automatically after execution.
5. **Archive** — Move completed change to `openspec/changes/archive/YYYY-MM-DD-<name>/` with delta spec sync back to main specs

### Architecture

- **Go binary** handles state management, spec parsing, config, and structured JSON output
- **SKILL.md files** orchestrate the AI workflow (invoke binary, present results, delegate to agents)
- **Agent definitions** (13 agents) perform the actual AI work (spec writing, execution, verification, research, etc.)
- **Reverse-calling pattern** — Claude Code calls the binary, not the other way around. No MCP server needed.

### Agent Roles

| Agent | Role |
|-------|------|
| mysd-proposal-writer | Write change proposals from user descriptions |
| mysd-spec-writer | Write requirements specs with RFC 2119 keywords |
| mysd-designer | Architecture and technical design |
| mysd-researcher | 4-dimension codebase research |
| mysd-advisor | Trade-off analysis for gray areas |
| mysd-planner | Task breakdown and dependency analysis |
| mysd-plan-checker | Verify plan covers all MUST items |
| mysd-executor | Implement tasks from plan |
| mysd-verifier | Goal-backward verification of MUST items |
| mysd-reviewer | Code review |
| mysd-scanner | Codebase scanning for spec generation |
| mysd-uat-guide | User acceptance testing guidance |
| mysd-fast-forward | Orchestrate accelerated workflows |

### Planning with Context

`/mysd:plan` supports multiple ways to control the planning process:

- **Per-spec planning** (`--spec <name>`) — restrict planning to a single spec capability. Useful when you have a multi-spec change and want to plan incrementally.
- **External input** (`--from <file>`) — load a file (e.g., a design doc, meeting notes, or existing plan from another tool) as additional planner context. The planner uses this alongside the spec artifacts to generate tasks.
- **Interactive spec selection** — when multiple specs exist and no `--spec` flag is given, an interactive picker lets you choose which specs to plan for.
- **Conversation context** — select "From conversation context" in the spec picker to extract requirements and task ideas from the current conversation. The planner writes a temp file and feeds it via `--from` automatically.
- **Research phase** (`--research`) — runs a focused architecture research pass before planning, useful for complex or unfamiliar areas.
- **Plan checker** (`--check`) — after planning, an independent agent verifies every MUST item in the spec has a corresponding task.
- **Design skip evaluation** — for simple changes (few files, no new capabilities), the planner can automatically skip the design phase.

The plan pipeline also includes automated self-review (placeholder detection, consistency checks, scope warnings, ambiguity fixes) and a reviewer agent pass before finalizing.

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
- Output feeds directly into spec artifacts so nothing is lost

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
