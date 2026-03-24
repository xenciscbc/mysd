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

# Define detailed requirements with RFC 2119 keywords
/mysd:spec

# Capture technical decisions
/mysd:design

# Break design into executable tasks
/mysd:plan

# Execute tasks (AI reads spec first, then codes)
/mysd:execute

# Verify all MUST items are satisfied
/mysd:verify

# Archive completed spec
/mysd:archive
```

Or fast-forward the entire workflow:

```bash
# propose -> spec -> design -> plan in one shot
/mysd:ff my-feature

# propose -> spec -> design -> plan -> execute in one shot
/mysd:ffe my-feature
```

## Commands

| Command | Description |
|---------|-------------|
| `/mysd:propose` | Create a new spec change with proposal scaffolding |
| `/mysd:spec` | Define requirements with RFC 2119 keywords (MUST/SHOULD/MAY) |
| `/mysd:design` | Capture technical decisions and architecture |
| `/mysd:plan` | Break design into an executable task list |
| `/mysd:execute` | Run tasks with mandatory alignment gate |
| `/mysd:verify` | Goal-backward verification of all MUST items |
| `/mysd:archive` | Archive verified spec to `.specs/archive/` |
| `/mysd:status` | Show current workflow state and progress |
| `/mysd:scan` | Scan existing codebase and generate specs |
| `/mysd:capture` | Extract changes from current conversation |
| `/mysd:ff` | Fast-forward: propose through plan |
| `/mysd:ffe` | Fast-forward: propose through execute |
| `/mysd:init` | Initialize project configuration |
| `/mysd:uat` | Interactive user acceptance testing |

## How It Works

mysd follows a structured lifecycle for every code change:

```
propose -> spec -> design -> plan -> execute -> verify -> archive
```

1. **Propose** — Scaffold spec artifacts (proposal.md, specs/, design.md, tasks.md) in `.specs/changes/`
2. **Spec** — Define requirements using RFC 2119 keywords. MUST items become verification gates.
3. **Design** — Record architecture decisions and technical approach
4. **Plan** — Break the design into ordered, executable tasks
5. **Execute** — AI reads the spec (alignment gate), then implements each task
6. **Verify** — Independent verifier agent checks every MUST item against the actual codebase
7. **Archive** — Move completed spec to `.specs/archive/` (blocked if any MUST item fails)

### Key Architecture

- **Go binary** handles state management, spec parsing, config, and structured JSON output
- **SKILL.md files** orchestrate the AI workflow (invoke binary, present results, delegate to agents)
- **Agent definitions** perform the actual AI work (spec writing, execution, verification)
- **Reverse-calling pattern** — Claude Code calls the binary, not the other way around. No MCP server needed.

## OpenSpec Compatibility

mysd reads and writes [OpenSpec](https://github.com/openspec-dev/openspec) format natively:

- Supports both `.specs/` and `openspec/` directory structures
- Parses YAML frontmatter with `spec-version` field
- Handles RFC 2119 keywords (MUST, SHOULD, MAY) and Delta Specs (ADDED, MODIFIED, REMOVED)
- Point mysd at an existing OpenSpec project and run `/mysd:execute` or `/mysd:verify` without migration

## Configuration

Project config lives in `.claude/mysd.yaml`:

```yaml
execution_mode: single    # single | wave
agent_count: 1
atomic_commits: false
tdd_mode: false
model_profile: balanced   # quality | balanced | budget
```

All options can be overridden per-command via flags.

## Tech Stack

- **Go 1.25+** — single binary, zero runtime dependencies
- **Cobra** — CLI framework
- **Viper** — configuration management
- **lipgloss** — terminal output styling
- **yaml.v3** — YAML parsing for OpenSpec frontmatter
- **GoReleaser** — cross-platform binary distribution

## License

MIT
