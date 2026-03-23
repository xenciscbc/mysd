# Architecture Research

**Domain:** AI-assisted Spec-Driven Development CLI (Go binary + Claude Code plugin)
**Researched:** 2026-03-23
**Confidence:** HIGH

## Standard Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Claude Code Plugin Layer                          │
│  commands/          agents/              hooks/         .mcp.json   │
│  ┌────────────┐  ┌───────────────┐  ┌──────────────┐  ┌──────────┐ │
│  │ /ssd:new   │  │ spec-writer   │  │ PostToolUse  │  │ MCP      │ │
│  │ /ssd:spec  │  │ task-runner   │  │ SessionStart │  │ Server   │ │
│  │ /ssd:apply │  │ verifier      │  │ (Go binary)  │  │(optional)│ │
│  │ /ssd:verify│  │ orchestrator  │  └──────────────┘  └──────────┘ │
│  └─────┬──────┘  └──────┬────────┘                                  │
└────────┼───────────────┼──────────────────────────────────────────┘
         │  invokes       │  spawns
         ▼               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    CLI Core (Go binary: myssd)                       │
├─────────────────────────────────────────────────────────────────────┤
│  cmd/                  (Cobra command tree)                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ propose  │ │  spec    │ │  design  │ │  plan    │ │ execute  │ │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ │
│       │            │            │             │            │        │
│  ┌──────────┐ ┌──────────┐                                          │
│  │  verify  │ │ archive  │                                          │
│  └────┬─────┘ └────┬─────┘                                          │
├───────┴────────────┴────────────────────────────────────────────────┤
│                    Internal Engine Layer                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌───────────┐  │
│  │ Spec Engine │  │  Execution  │  │ Verification│  │   State   │  │
│  │  (parser +  │  │   Engine    │  │  Pipeline   │  │ Manager   │  │
│  │  generator) │  │ (orchestr.) │  │ (goal-back) │  │           │  │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └─────┬─────┘  │
└─────────┼───────────────┼────────────────┼───────────────┼─────────┘
          │               │                │               │
          ▼               ▼                ▼               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Storage Layer                                      │
│  .specs/                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │
│  │ changes/     │  │ specs/       │  │ .ssd-state.json          │   │
│  │ [name]/      │  │ (source of   │  │ (workflow state, current  │   │
│  │  proposal.md │  │  truth)      │  │  phase, last run result)  │   │
│  │  specs/      │  │              │  └──────────────────────────┘   │
│  │  design.md   │  └──────────────┘                                 │
│  │  tasks.md    │                                                    │
│  └──────────────┘                                                    │
└─────────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Typical Implementation |
|-----------|----------------|------------------------|
| Claude Code Plugin Layer | Expose slash commands and agents to Claude Code; invoke Go binary via hooks or MCP | Markdown files in `commands/`, `agents/`; `hooks/hooks.json`; optional `.mcp.json` |
| Cobra Command Tree (`cmd/`) | Parse CLI arguments; route to correct engine function; surface user-facing errors | `cobra.Command` structs, one file per top-level command |
| Spec Engine | Parse/validate OpenSpec Markdown; generate artifact scaffolds; write delta specs | `internal/spec/` — parser, schema validator, writer |
| Execution Engine | Orchestrate single/multi-agent runs; manage wave-based task execution; pass spec context | `internal/engine/` — orchestrator, agent runner, context builder |
| Verification Pipeline | Goal-backward check: iterate MUST items in specs; compare to evidence; emit verdict | `internal/verify/` — must-collector, verifier, report writer |
| State Manager | Track current workflow phase, last command run, spec change name; enable resume | `internal/state/` — read/write `.ssd-state.json` |
| Storage Layer | `.specs/` directory on disk; Markdown files as single source of truth | No DB; plain files + `encoding/json` for state blob |

## Recommended Project Structure

```
mysd/                              # repository root
├── main.go                        # entry point — cobra root cmd
├── cmd/                           # one file per top-level command
│   ├── root.go                    # persistent flags, version
│   ├── propose.go                 # ssd propose
│   ├── spec.go                    # ssd spec
│   ├── design.go                  # ssd design
│   ├── plan.go                    # ssd plan
│   ├── execute.go                 # ssd execute
│   ├── verify.go                  # ssd verify
│   └── archive.go                 # ssd archive
├── internal/
│   ├── spec/
│   │   ├── parser.go              # parse proposal.md / specs/*.md / design.md / tasks.md
│   │   ├── schema.go              # OpenSpec schema structs (RFC 2119 keywords, delta types)
│   │   ├── validator.go           # validate spec completeness and format
│   │   ├── writer.go              # scaffold new change dirs, write artifacts
│   │   └── delta.go               # ADDED / MODIFIED / REMOVED delta logic
│   ├── engine/
│   │   ├── orchestrator.go        # coordinate multi-step execution
│   │   ├── agent.go               # single agent runner (Claude Code subagent call)
│   │   ├── context_builder.go     # assemble spec context injected before AI runs
│   │   └── wave.go                # wave-based parallel task execution (future)
│   ├── verify/
│   │   ├── must_collector.go      # extract all MUST items from specs/
│   │   ├── verifier.go            # goal-backward: check each MUST against codebase
│   │   ├── reporter.go            # write verification report, update spec status
│   │   └── feedback.go            # feed results back into spec metadata
│   ├── state/
│   │   ├── state.go               # WorkflowState struct + read/write
│   │   └── lock.go                # prevent concurrent runs (file lock)
│   └── config/
│       ├── config.go              # project config (.ssd.toml or convention defaults)
│       └── defaults.go            # convention-over-config defaults
├── plugin/                        # Claude Code plugin directory (installed separately)
│   ├── .claude-plugin/
│   │   └── plugin.json            # plugin manifest
│   ├── commands/
│   │   ├── new.md                 # /ssd:new
│   │   ├── spec.md                # /ssd:spec
│   │   ├── design.md              # /ssd:design
│   │   ├── plan.md                # /ssd:plan
│   │   ├── apply.md               # /ssd:apply (wraps execute)
│   │   ├── verify.md              # /ssd:verify
│   │   └── archive.md             # /ssd:archive
│   ├── agents/
│   │   ├── spec-writer.md         # generates spec artifacts
│   │   ├── task-runner.md         # executes tasks from tasks.md
│   │   └── verifier.md            # runs goal-backward verification
│   └── hooks/
│       └── hooks.json             # e.g. PreToolUse to enforce spec alignment check
├── testdata/
│   └── fixtures/                  # sample .specs/ trees for unit tests
├── go.mod
├── go.sum
└── Makefile
```

### Structure Rationale

- **`cmd/`**: Each command file is one `cobra.Command`. Commands are thin — they only parse flags and call `internal/` functions. No business logic in `cmd/`.
- **`internal/spec/`**: OpenSpec format knowledge is isolated here. If the format changes, only this package changes.
- **`internal/engine/`**: Execution concerns (orchestration, agent invocation, context building) are separated from spec concerns. The engine reads parsed spec structs, not raw Markdown.
- **`internal/verify/`**: Goal-backward verification is a distinct pipeline, not mixed into execution. This makes it independently testable and callable as a standalone step.
- **`internal/state/`**: Explicit state machine file (`.ssd-state.json`) rather than inferring state from filesystem presence. Enables crash recovery and the `--resume` flag.
- **`plugin/`**: Claude Code plugin lives in a separate directory. It is distributed and installed independently from the Go binary. It invokes the binary via `${CLAUDE_PLUGIN_ROOT}/../bin/myssd` or a PATH lookup.

## Architectural Patterns

### Pattern 1: Thin Commands, Fat Internal

**What:** `cmd/` files contain zero business logic. Every command function calls one `internal/` function, then formats output.
**When to use:** Always, from day one.
**Trade-offs:** Slightly more indirection, but enables unit-testing all business logic without CLI ceremony.

**Example:**
```go
// cmd/verify.go
func runVerify(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load(".")
    if err != nil {
        return err
    }
    report, err := verify.RunPipeline(cfg)
    if err != nil {
        return err
    }
    printer.PrintVerifyReport(cmd.OutOrStdout(), report)
    return nil
}
```

### Pattern 2: Spec as Struct, Not String

**What:** Parse OpenSpec Markdown into typed Go structs at the boundary (`internal/spec/parser.go`). All downstream code operates on `spec.Change`, `spec.Requirement`, `spec.Task` — never on raw strings.
**When to use:** From the first parser implementation.
**Trade-offs:** Parser complexity upfront, but all engine/verify logic becomes straightforward struct traversal.

**Example:**
```go
// internal/spec/schema.go
type Requirement struct {
    ID       string
    Text     string
    Keyword  RFC2119Keyword  // MUST | SHOULD | MAY
    DeltaOp  DeltaOp         // ADDED | MODIFIED | REMOVED | (none)
}

type Change struct {
    Name     string
    Proposal ProposalDoc
    Specs    []Requirement
    Design   DesignDoc
    Tasks    []Task
}
```

### Pattern 3: Explicit Workflow State Machine

**What:** Store current phase (`proposed | specced | designed | planned | executed | verified | archived`) in `.ssd-state.json`. Commands validate state transitions before running.
**When to use:** From the first multi-step command.
**Trade-offs:** Small overhead writing state file, but prevents running `verify` before `execute` and enables `--resume`.

**Example:**
```go
// internal/state/state.go
type Phase string
const (
    PhaseProposed Phase = "proposed"
    PhaseSpecced  Phase = "specced"
    PhaseDesigned Phase = "designed"
    PhasePlanned  Phase = "planned"
    PhaseExecuted Phase = "executed"
    PhaseVerified Phase = "verified"
    PhaseArchived Phase = "archived"
)

type WorkflowState struct {
    ChangeName  string    `json:"change_name"`
    Phase       Phase     `json:"phase"`
    LastRun     time.Time `json:"last_run"`
    VerifyPass  *bool     `json:"verify_pass,omitempty"`
}
```

### Pattern 4: Plugin Delegates to Binary

**What:** Claude Code plugin commands are thin Markdown wrappers. They construct and run the `myssd` binary via shell command. Plugin agents read spec files directly but trigger the binary for mutations.
**When to use:** This is the integration model — not optional.
**Trade-offs:** Requires Go binary to be on PATH when Claude Code invokes hooks. Installation step needed. Benefit: no logic duplication between CLI and plugin.

**Example (commands/apply.md):**
```markdown
---
description: Execute implementation tasks from tasks.md
---
Run the spec-driven development execution engine.

Execute: `myssd execute --change $ARGUMENTS`

After execution completes, report the results and ask whether to proceed to verify.
```

## Data Flow

### Primary Workflow: propose → verify

```
Developer: myssd propose "add auth"
    │
    ▼
cmd/propose.go
    │ config.Load()
    ▼
internal/spec/writer.go
    │ scaffold .specs/changes/add-auth/
    │   proposal.md (template)
    │   specs/      (empty)
    │   design.md   (empty)
    │   tasks.md    (empty)
    ▼
internal/state/state.go
    │ write .ssd-state.json { phase: "proposed", change: "add-auth" }
    ▼
[Claude Code reads proposal.md, human reviews and edits]

Developer: myssd spec  (or /ssd:spec in Claude Code)
    │
    ▼
cmd/spec.go
    │ state.Load() → assert phase == "proposed"
    ▼
internal/engine/context_builder.go
    │ read proposal.md → build AI prompt context
    ▼
internal/engine/agent.go
    │ invoke spec-writer agent (Claude subagent or direct AI call)
    │ agent writes specs/*.md with RFC 2119 MUST/SHOULD/MAY
    ▼
internal/spec/validator.go
    │ validate generated specs against schema
    ▼
internal/state/state.go
    │ update phase → "specced"
    ▼
[repeat pattern for design → plan]

Developer: myssd execute
    │
    ▼
cmd/execute.go
    │ state.Load() → assert phase == "planned"
    ▼
internal/spec/parser.go
    │ parse tasks.md → []Task (ordered, with spec references)
    ▼
internal/engine/orchestrator.go
    │ for each Task:
    │   context_builder: inject task + relevant spec requirements
    │   agent.Run(task, specContext)  ← task-runner agent
    │   update task completion status in tasks.md
    ▼
internal/state/state.go
    │ update phase → "executed"
    ▼

Developer: myssd verify
    │
    ▼
cmd/verify.go
    │ state.Load() → assert phase == "executed"
    ▼
internal/spec/parser.go
    │ parse specs/*.md → []Requirement where Keyword == MUST
    ▼
internal/verify/must_collector.go
    │ collect all MUST requirements
    ▼
internal/verify/verifier.go
    │ for each MUST requirement:
    │   build verification query (goal-backward)
    │   invoke verifier agent with codebase context
    │   collect PASS / FAIL verdict
    ▼
internal/verify/reporter.go
    │ write verification-report.md
    │ update spec metadata (status fields)
    ▼
internal/verify/feedback.go
    │ if any FAIL → set state.VerifyPass = false
    │ else → set state.VerifyPass = true
    ▼
internal/state/state.go
    │ update phase → "verified"
    ▼

Developer: myssd archive (only if VerifyPass == true)
    │
    ▼
internal/spec/delta.go
    │ merge ADDED/MODIFIED/REMOVED delta specs into .specs/specs/
    │ move .specs/changes/[name]/ → .specs/changes/archive/[name]/
    ▼
internal/state/state.go
    │ update phase → "archived"
    │ clear current change context
```

### Claude Code Plugin Integration Flow

```
User types /ssd:apply in Claude Code
    │
    ▼
commands/apply.md is loaded as context
    │ instructs Claude to run: myssd execute --change [current-change]
    ▼
Claude Code executes Bash tool: myssd execute
    │
    ▼
myssd binary (Go) runs execute command
    │ reads .ssd-state.json for change context
    │ orchestrates task-runner agent (agents/task-runner.md)
    ▼
Subagent (agents/task-runner.md) receives:
    │ - task description
    │ - spec context (relevant MUST requirements)
    │ - design constraints
    │ Implements code, calls Write/Edit tools
    ▼
myssd binary receives subagent completion signal
    │ validates task marked complete in tasks.md
    ▼
Claude Code reports results to user
```

### State Transitions

```
(none)
  │ myssd propose
  ▼
proposed
  │ myssd spec
  ▼
specced
  │ myssd design
  ▼
designed
  │ myssd plan
  ▼
planned
  │ myssd execute
  ▼
executed
  │ myssd verify
  ▼
verified ──(FAIL)──► executed  (can re-execute and re-verify)
  │ (PASS)
  │ myssd archive
  ▼
archived
```

## Scaling Considerations

This is a local developer CLI tool. "Scaling" means complexity growth, not user count.

| Scale | Architecture Adjustments |
|-------|--------------------------|
| v1 (single agent, sequential tasks) | Current design — orchestrator runs tasks one at a time, single AI call per task |
| v2 (parallel agent waves) | `engine/wave.go` dispatches independent tasks concurrently via goroutines; tasks with no shared spec sections run in parallel |
| v3 (multi-project / workspace) | `config.go` supports workspace-level `.ssd.toml`; state manager tracks multiple active changes |

### Scaling Priorities

1. **First bottleneck:** Sequential task execution is slow for large plans. Fix: implement wave-based parallelism in `engine/wave.go` (goroutines + channels). Tasks with no overlapping spec dependencies run concurrently.
2. **Second bottleneck:** Verification of many MUST items is slow. Fix: batch verification queries; run independent MUST checks in parallel goroutines.

## Anti-Patterns

### Anti-Pattern 1: Embed Spec Logic in Commands

**What people do:** Parse `proposal.md` directly inside `cmd/propose.go`; string-match MUST/SHOULD inline.
**Why it's wrong:** Commands become untestable monoliths. Changing OpenSpec format requires touching command files.
**Do this instead:** `cmd/propose.go` calls `spec.ParseChange(dir)`. All format knowledge lives in `internal/spec/`.

### Anti-Pattern 2: Infer Phase from Filesystem

**What people do:** Check "if design.md exists, phase is designed." No explicit state file.
**Why it's wrong:** Partial writes, interrupted commands, and manual edits create ambiguous states. No way to resume.
**Do this instead:** Maintain explicit `.ssd-state.json` as the authoritative phase record. Filesystem content is data; state file is the workflow cursor.

### Anti-Pattern 3: Plugin Contains Business Logic

**What people do:** Write spec parsing and verification logic inside `commands/*.md` agent instructions.
**Why it's wrong:** Logic duplicated between plugin and binary. Plugin updates require re-deploying Markdown files; not testable.
**Do this instead:** Plugin commands are wrappers that invoke `myssd` binary. Binary owns all logic. Plugin owns presentation and Claude Code integration.

### Anti-Pattern 4: Single Agent for All Tasks

**What people do:** One monolithic "do everything" agent that reads all tasks and implements them in one shot.
**Why it's wrong:** Large context leads to incomplete implementations. No task-level auditability. Cannot resume from partial failure.
**Do this instead:** Orchestrator dispatches one agent invocation per task. Each agent receives task + relevant spec context only. Completion tracked per-task in `tasks.md`.

### Anti-Pattern 5: Verify by Re-running Code

**What people do:** Verification = run tests. PASS if tests pass.
**Why it's wrong:** Tests may not cover all MUST requirements. Spec requirements that aren't tested (yet) are silently skipped.
**Do this instead:** Goal-backward verification starts from spec MUST items and asks "is this requirement satisfied in the codebase?" Tests are evidence, not the only evidence. AI-driven requirement tracing covers gaps.

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| Claude Code (slash commands) | Plugin Markdown files in `commands/` directory; invoke binary via Bash tool | Plugin installed to `~/.claude/plugins/` or `.claude/` scoped |
| Claude Code (subagents) | Agent Markdown files in `agents/` directory with frontmatter; orchestrator spawns them | Agents receive injected spec context via system prompt additions |
| Claude Code (hooks) | `hooks/hooks.json` with PostToolUse or PreToolUse matchers | Optional: enforce spec-alignment check before Write/Edit tools |
| MCP Server (optional v2) | Go binary exposes MCP stdio server; Claude Code connects via `.mcp.json` | Enables richer tool calls (read spec, check state) as MCP tools instead of shell invocation |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| `cmd/` ↔ `internal/engine/` | Direct Go function calls; cmd passes parsed config | No interface abstraction needed at v1; add if multiple execution backends needed |
| `cmd/` ↔ `internal/spec/` | Direct Go function calls; returns typed structs | Parser is the only consumer of raw Markdown |
| `internal/engine/` ↔ `internal/spec/` | Engine receives `spec.Change` struct; never reads files | Dependency direction: engine depends on spec, not reverse |
| `internal/engine/` ↔ `internal/verify/` | Separate pipeline invocations; both read parsed spec | No direct coupling; both are orchestrated by cmd layer |
| `internal/state/` ↔ all | State is read/written by cmd layer; internal packages do not write state | Prevents hidden state transitions inside engine or verify |
| Plugin ↔ Binary | Shell exec: `myssd <command> --change <name>` | Binary reads `.ssd-state.json` for context; no in-process coupling |

## Suggested Build Order

Dependencies drive this order — each layer depends on the one before it.

```
1. Storage schema         .specs/ directory conventions, file layout docs
        │
        ▼
2. internal/spec/         Parser + schema structs + writer
   (parse, validate,      Foundation: everything depends on spec types
    scaffold)
        │
        ▼
3. internal/state/        WorkflowState struct + read/write + lock
   (phase tracking)       Needed before any command can enforce transitions
        │
        ▼
4. cmd/ skeleton          Cobra root + all commands registered (thin stubs)
   (CLI wiring)           Enables integration testing of CLI surface early
        │
        ▼
5. internal/engine/       Context builder + single-agent runner
   (execution)            Depends on spec (for context) and state (for phase)
        │
        ▼
6. internal/verify/       MUST collector + verifier + reporter
   (verification)         Depends on spec (for requirements), engine patterns
        │
        ▼
7. plugin/                Slash commands + agents (Markdown)
   (Claude Code layer)    Written after binary commands are stable; thin wrappers
        │
        ▼
8. internal/engine/wave   Parallel task dispatch (v2 — after v1 proven)
   (parallelism)
```

## Sources

- [Claude Code Plugins Reference](https://code.claude.com/docs/en/plugins-reference) — authoritative plugin directory structure, agent frontmatter schema, hook events (HIGH confidence)
- [OpenSpec Workflow](https://openspec.pro/workflow/) — proposal/spec/design/tasks artifact model (HIGH confidence)
- [OpenSpec GitHub Workflows doc](https://github.com/Fission-AI/OpenSpec/blob/main/docs/workflows.md) — command set and flow patterns (HIGH confidence)
- [Structuring Go Code for CLI Applications](https://www.bytesizego.com/blog/structure-go-cli-app) — cmd/internal/pkg organization patterns (HIGH confidence)
- [Go CLI Applications with Cobra and Viper](https://www.glukhov.org/post/2025/11/go-cli-applications-with-cobra-and-viper/) — Cobra command hierarchy, flag inheritance (HIGH confidence)
- [AI Workflow Patterns in Go](https://dasroot.net/posts/2026/02/ai-workflow-patterns-go-cli-tools-agents/) — goroutine-based orchestration, agent dispatch patterns (MEDIUM confidence)
- [mcp-go](https://github.com/mark3labs/mcp-go) — Go MCP server stdlib stdio integration (MEDIUM confidence, for v2 MCP path)
- [GitHub spec-kit](https://github.com/github/spec-kit) — parallel SDD toolkit for comparison (MEDIUM confidence)

---
*Architecture research for: AI-assisted Spec-Driven Development CLI (my-ssd)*
*Researched: 2026-03-23*
