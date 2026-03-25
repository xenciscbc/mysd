# Architecture Research

**Domain:** Go CLI Tool — Interactive Discovery, Worktree Parallel Execution, Subagent Orchestration (v1.1 Integration)
**Researched:** 2026-03-25
**Confidence:** HIGH (direct codebase analysis of v1.0 + proposal spec)

## Context: Subsequent Milestone

This is NOT a greenfield architecture document. v1.0 shipped with a clean 11-package architecture
(7,555 LOC). This document focuses specifically on how v1.1 features integrate with that
existing foundation — what is new, what is modified, and what must not be broken.

**Existing packages inventory:**

```
cmd/                  — 14 Cobra command files (thin layer, no business logic)
internal/config/      — ProjectConfig, ResolveModel, DefaultModelMap, Defaults, BindFlags
internal/spec/        — OpenSpec parser: schema.go, parser.go, updater.go, delta.go, detector.go, writer.go
internal/state/       — WorkflowState, ValidTransitions, Transition, LoadState, SaveState
internal/executor/    — ExecutionContext, BuildContext, BuildContextFromParts, alignment, progress, status
internal/scanner/     — Codebase scanner (currently Go-specific)
internal/verifier/    — Verification context and report
internal/output/      — Printer (lipgloss-backed terminal output)
internal/roadmap/     — tracking.yaml + Mermaid timeline
internal/uat/         — UAT checklist generation
plugin/commands/      — 14 SKILL.md orchestrator files
plugin/agents/        — 8 agent definition files
plugin/hooks/         — SessionStart hook
```

---

## Standard Architecture

### System Overview: v1.0 (Existing)

```
┌─────────────────────────────────────────────────────────────────────┐
│                       Claude Code Session                            │
│                                                                      │
│  plugin/commands/*.md     (SKILL.md orchestrators — 14 commands)    │
│           │ Task tool                                                │
│           ▼                                                          │
│  plugin/agents/*.md       (subagent definitions — 8 agents)         │
└─────────────────────┬───────────────────────────────────────────────┘
                      │ bash: mysd {cmd} --context-only → JSON
                      ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    mysd binary (Go)                                  │
│                                                                      │
│  cmd/ → internal/config → .claude/mysd.yaml                        │
│       → internal/spec   → .specs/changes/{name}/ filesystem        │
│       → internal/state  → STATE.json                                │
│       → internal/executor → ExecutionContext JSON (--context-only)  │
│       → internal/roadmap → tracking.yaml                           │
└─────────────────────────────────────────────────────────────────────┘
```

**Core principle (must be preserved in v1.1):** The Go binary is a state manager and
context provider. SKILL.md files are orchestrators. Agent .md files are workers. The
binary outputs structured JSON that SKILL.md files consume via `--context-only` flags.
The binary never directly invokes AI or prompts the user interactively.

---

### System Overview: v1.1 (Target)

```
┌─────────────────────────────────────────────────────────────────────┐
│                       Claude Code Session                            │
│                                                                      │
│  NEW COMMANDS:              MODIFIED COMMANDS:                       │
│  /mysd:discuss              /mysd:propose  (+ research mode)        │
│  /mysd:fix                  /mysd:spec     (+ research mode)        │
│  /mysd:model                /mysd:plan     (+ plan-checker trigger) │
│  /mysd:lang                 /mysd:execute  (+ wave groups, worktree)│
│  /mysd:scan (upgraded)      /mysd:ff       (+ auto, research-once)  │
│  /mysd:ffe (upgraded)                                                │
│           │ Task tool                                                │
│           ▼                                                          │
│  NEW AGENTS:                 MODIFIED AGENTS:                        │
│  mysd-researcher             mysd-spec-writer  (per-capability)     │
│  mysd-advisor x N            mysd-executor     (per-task, worktree) │
│  mysd-proposal-writer                                                │
│  mysd-plan-checker           UNCHANGED AGENTS:                      │
│                              mysd-designer, mysd-planner,           │
│                              mysd-verifier, mysd-fast-forward,      │
│                              mysd-scanner, mysd-uat-guide            │
└─────────────────────┬───────────────────────────────────────────────┘
                      │ bash: mysd {cmd} --context-only → JSON
                      ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    mysd binary (Go) — v1.1                           │
│                                                                      │
│  NEW COMMANDS:              NEW INTERNAL PACKAGES:                   │
│  cmd/model.go               internal/worktree/    (git worktree)    │
│  cmd/lang.go                internal/discovery/   (discovery state) │
│                             internal/planchecker/ (MUST coverage)   │
│  MODIFIED COMMANDS:                                                  │
│  cmd/execute.go  (+waves)   MODIFIED INTERNAL PACKAGES:             │
│  cmd/plan.go     (+checker) internal/executor/  (+depends, waves)   │
│  cmd/scan.go     (+lang)    internal/config/    (+new agents, auto) │
│  cmd/propose.go  (+auto)    internal/spec/      (+Depends, Files)   │
│  cmd/ff.go       (+auto)    internal/scanner/   (lang-agnostic)     │
│                                                                      │
│  UNCHANGED COMMANDS:        UNCHANGED INTERNAL PACKAGES:            │
│  cmd/spec.go                internal/state/                         │
│  cmd/design.go              internal/output/                        │
│  cmd/verify.go              internal/verifier/                      │
│  cmd/archive.go             internal/roadmap/                       │
│  cmd/status.go              internal/uat/                           │
│  cmd/task_update.go                                                  │
└─────────────────────────────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Filesystem State                                   │
│                                                                      │
│  EXISTING:                           NEW:                            │
│  .specs/changes/{name}/              .worktrees/T{id}/  (ephemeral) │
│    proposal.md                                                       │
│    specs/                            EXTENDED:                       │
│    design.md                         .specs/changes/{name}/         │
│    tasks.md        (+ depends/files) discovery-state.json (NEW)     │
│    STATE.json                        .claude/mysd.yaml   (extended) │
│    alignment.md                                                      │
│    verification-status.json                                          │
└─────────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

#### New Packages

| Package | Responsibility | Key Types/Functions |
|---------|----------------|---------------------|
| `internal/worktree/` | Git worktree lifecycle: create at `.worktrees/T{id}/`, create branch `mysd/{change}/T{id}-{slug}`, merge `--no-ff` in task ID order, cleanup on success, preserve on failure | `WorktreeManager`, `Create(id, slug)`, `Merge(id)`, `Remove(id)` |
| `internal/discovery/` | Builds discovery context JSON for SKILL.md consumption via `--context-only`. Persists discovery state (areas completed, deferred notes, research summary) as sidecar JSON. NOT responsible for AI questions — that is SKILL.md's job. | `DiscoveryContext`, `DiscoveryState`, `BuildDiscoveryContext(changeName, specDir)`, `LoadState`, `SaveState` |
| `internal/planchecker/` | Pure function: validate tasks.md MUST coverage. Receives task list + MUST items, returns uncovered IDs. No I/O side effects. | `CheckCoverage(tasks, mustItems) CoverageResult`, `UncoveredItem` |

#### Modified Packages

| Package | What Changes | Why |
|---------|-------------|-----|
| `internal/executor/` | Add `Depends []int` and `Files []string` to `TaskItem`. Add `BuildWaveGroups(tasks []TaskItem) [][]TaskItem` computing wave layers from depends + file overlap. Add `WaveGroups`, `WorktreeDir`, `AutoMode` to `ExecutionContext`. | Supports parallel task execution with dependency-ordered waves |
| `internal/config/` | Add new agent roles to `DefaultModelMap` (researcher, advisor, proposal-writer, plan-checker). Add `WorktreeDir string` and `AutoMode bool` to `ProjectConfig`. Update `Defaults()`. | New agents need model resolution; worktree path must be configurable |
| `internal/spec/` | Add `Depends []int` and `Files []string` to `TaskEntry` and `TasksFrontmatterV2`. Update `WriteTasks` to serialize new fields. | tasks.md carries dependency metadata written by planner agent |
| `internal/scanner/` | Refactor from Go-specific to language-agnostic. Add locale detection. Produce `openspec/config.yaml` + spec stubs. Accept `--scaffold-only` mode to replace init command logic. | /mysd:scan upgrade + /mysd:init replacement |

---

## Recommended Project Structure (v1.1 additions only)

```
mysd/
├── cmd/
│   ├── [existing 14 commands — unchanged]
│   ├── model.go           # NEW: mysd model / model set / model resolve
│   └── lang.go            # NEW: mysd lang (write locale to yaml files)
│   NOTE: discuss and fix have no binary counterparts — SKILL.md only
│
├── internal/
│   ├── config/            # MODIFIED: new agent roles, WorktreeDir, AutoMode
│   ├── executor/          # MODIFIED: wave grouping, depends/files, WaveGroups in context
│   ├── spec/              # MODIFIED: Depends/Files in TaskEntry schema
│   ├── scanner/           # MODIFIED: language-agnostic refactor
│   │
│   ├── worktree/          # NEW package
│   │   ├── worktree.go    # WorktreeManager, Create, Remove
│   │   ├── merge.go       # Merge, conflict detection helpers
│   │   └── worktree_test.go
│   │
│   ├── discovery/         # NEW package
│   │   ├── context.go     # DiscoveryContext, BuildDiscoveryContext
│   │   ├── state.go       # DiscoveryState persistence (sidecar JSON)
│   │   └── discovery_test.go
│   │
│   └── planchecker/       # NEW package
│       ├── checker.go     # CheckCoverage(tasks, mustItems) CoverageResult
│       └── checker_test.go
│
└── plugin/
    ├── commands/
    │   ├── [existing 14 commands — some modified]
    │   ├── mysd-discuss.md    # NEW
    │   ├── mysd-fix.md        # NEW
    │   ├── mysd-model.md      # NEW
    │   └── mysd-lang.md       # NEW
    └── agents/
        ├── [existing 8 agents — some modified]
        ├── mysd-researcher.md       # NEW
        ├── mysd-advisor.md          # NEW
        ├── mysd-proposal-writer.md  # NEW
        └── mysd-plan-checker.md     # NEW
```

### Structure Rationale

- **`internal/worktree/` separate package:** Git worktree operations are stateful and independently testable. Does NOT import `internal/executor/` — keeps the dependency graph acyclic. The executor calls worktree via context JSON passed through SKILL.md.
- **`internal/discovery/` separate package:** Discovery state has different lifecycle than execution state (spans propose/spec/discuss phases, not execution phase). Prevents `internal/executor/` from growing to include discovery concerns.
- **`internal/planchecker/` pure package:** No filesystem I/O, no config dependencies. `CheckCoverage` is a pure function. Can be unit-tested without any mocking. Intentionally minimal by design.
- **discuss and fix as SKILL.md-only commands:** These flows are AI orchestration with no binary state changes beyond what existing commands already handle. The binary cannot drive interactive discovery sessions or AI-powered conflict resolution. This preserves the core pattern: binary = state/context, SKILL.md = orchestration.

---

## Architectural Patterns

### Pattern 1: Binary-as-Context-Provider (existing — must be maintained)

**What:** Every SKILL.md orchestrator runs `mysd {cmd} --context-only` to get structured JSON, then passes that JSON to subagents via the Task tool. The binary never drives AI execution or prompts interactively.

**When to use:** All new binary commands (model, lang, execute with waves). The pattern MUST be preserved.

**Extended for v1.1 execute — wave context output:**
```json
{
  "change_name": "my-feature",
  "pending_tasks": [...],
  "wave_groups": [
    [{"id": 1, "name": "task-a", "depends": [], "files": ["pkg/auth.go"]},
     {"id": 2, "name": "task-b", "depends": [], "files": ["pkg/cache.go"]}],
    [{"id": 3, "name": "task-c", "depends": [1, 2], "files": ["main.go"]}]
  ],
  "worktree_dir": ".worktrees",
  "auto_mode": false
}
```

### Pattern 2: Orchestrator-Spawns-Subagents (existing — extended for parallelism)

**What:** SKILL.md orchestrators use the Claude Code Task tool to spawn named agents. New in v1.1: parallel spawn for advisors and parallel spawn for wave-mode executors.

**New parallel advisor pattern (in mysd-propose.md or mysd-discuss.md):**
```
For each gray area from researcher output, spawn parallel:
Task → mysd-advisor
  area: {one gray area}
  codebase_context: {from researcher}
  change_name: {change_name}
```

**New wave executor pattern (in mysd-execute.md):**
```
For each task in wave_groups[current_wave], spawn parallel:
Task → mysd-executor
  assigned_task: {task}
  isolation: worktree
  worktree_path: .worktrees/T{id}/
  branch: mysd/{change}/T{id}-{slug}
```

### Pattern 3: Sidecar JSON for Phase State (existing — extended)

**What:** Phase-specific state is stored as sidecar JSON files in the change directory. Existing examples: `STATE.json`, `verification-status.json`, `alignment.md`.

**New in v1.1 — discovery state sidecar format:**
```json
// .specs/changes/{name}/discovery-state.json
{
  "research_enabled": true,
  "areas_completed": ["auth-strategy", "error-handling"],
  "deferred_notes": ["consider oauth2 in v2"],
  "research_summary": "...",   // cached for ff/ffe reuse — avoid re-researching
  "updated_at": "2026-03-25T12:00:00Z"
}
```

### Pattern 4: Interactive Prompts as SKILL.md Responsibility (new — important constraint)

**What:** Whenever a feature requires interactive questions (research mode, locale choice, single vs wave), the SKILL.md asks the question and passes the answer as argument or config to the binary. The binary never prompts stdin.

**Why this matters:** The binary is called with `--context-only` flags by SKILL.md. If the binary reads stdin, it will hang when called non-interactively.

**Implementation example for mysd lang:**
- SKILL.md `/mysd:lang` asks: "What language for AI responses? (e.g. zh-TW, en, ja)"
- User answers, SKILL.md calls: `mysd lang set zh-TW`
- Binary `cmd/lang.go` receives the value as argument, writes to yaml files, exits

### Pattern 5: Worktree-Per-Task Isolation (new)

**What:** Each parallel task executor gets its own git worktree at `.worktrees/T{id}/`. After all tasks in a wave complete, merge back to main branch in ascending task ID order using `git merge --no-ff`. Delete worktree on success, preserve on failure for debugging.

**Wave boundary contract:** Create all worktrees for a wave from the same HEAD snapshot. Complete the entire wave (including merges) BEFORE creating worktrees for the next wave. This prevents divergence across waves.

**Conflict resolution:** AI attempts auto-resolve → run `go build` + `go test` → retry max 3 times → if still failing, notify user and preserve worktree.

---

## Data Flow

### Interactive Discovery Flow (new — propose/spec/discuss)

```
User: /mysd:propose {name}
    ↓
SKILL.md: "Use research mode? (y/n)"
    ↓ [If yes]
Task → mysd-researcher
  Reads codebase (grep/glob — no new binary command needed)
  Outputs: gray_areas[], codebase_summary
    ↓
[Parallel, for each gray area]
Task → mysd-advisor
  Input: one area + codebase_summary
  Output: analysis + recommendation with comparison table
    ↓
SKILL.md: Present areas one-by-one (discussion loop with user)
  Inner loop: deep-dive into one area
  Outer loop: after all areas, "Any new areas? (y/n)"
  Scope guardrail: off-topic → redirect to deferred_notes
    ↓
Task → mysd-proposal-writer
  Input: all discussion conclusions + user decisions
  Writes: proposal.md body
    ↓
Binary: mysd propose {name}   (scaffold + state transition only)
    ↓
SKILL.md: Show summary → next: /mysd:spec
```

### Worktree Parallel Execution Flow (new — execute wave mode)

```
SKILL.md: mysd execute --context-only
    ↓ JSON with wave_groups[][]
SKILL.md: Ask "Single or wave execution? (single/wave)"
          [or --auto: use config default]
    ↓ [wave mode selected]
For wave[0] — spawn parallel:
  Task → mysd-executor (worktree: .worktrees/T1/, branch: mysd/{change}/T1-{slug})
  Task → mysd-executor (worktree: .worktrees/T2/, branch: mysd/{change}/T2-{slug})
  Task → mysd-executor (worktree: .worktrees/T3/, branch: mysd/{change}/T3-{slug})
    ↓ [All wave[0] agents complete; failures do NOT stop others]
SKILL.md: Merge in ascending task ID order: T1 → T2 → T3
  For each merge:
    Conflict? → AI resolves → go build + go test → retry max 3x
    Success → delete worktree + branch
    Failure (after 3x) → preserve worktree, notify user
    ↓
[If all merges OK] mysd task-update {id} done for each task
    ↓
For wave[1] from updated HEAD — repeat...
```

### Plan-Checker Auto-Trigger Flow (new — post-plan)

```
mysd-planner writes tasks.md
    ↓
mysd-planner runs: mysd plan  (state transition)
    ↓
SKILL.md /mysd:plan orchestrator receives plan completion
    ↓
Task → mysd-plan-checker
  Input: must_items from spec + tasks from tasks.md
  Output: uncovered_must_ids[]
    ↓
[If uncovered_must_ids is empty]
  SKILL.md: "All MUST items covered."
[If not empty]
  SKILL.md: "Plan-checker: {N} uncovered MUST items.
             Auto-add suggested tasks? (y/n)"
  [--auto or user says yes] → mysd-planner adds tasks → re-run plan-checker
  [user says no] → show gaps, user decides
```

### Discuss Update Flow (new — re-entry after spec is set)

```
User: /mysd:discuss
    ↓
SKILL.md: Get topic from user OR offer research mode
    ↓ [Optional research]
Task → mysd-researcher (scoped to relevant area)
    ↓
SKILL.md: Discussion loop (same pattern as discovery)
    ↓
SKILL.md: Determine what changed (spec? design? tasks?)
    ↓ [Update affected files directly via SKILL.md Write tool]
mysd spec    (if specs changed → re-specced)
mysd design  (if design changed → re-designed)
mysd plan    (if tasks need replanning → re-planned)
    ↓
Auto-spawn mysd-plan-checker (same as post-plan flow)
    ↓
SKILL.md: "Updated: {list of changed files}. Tasks re-planned."
```

NOTE: The existing `ValidTransitions` map already supports these re-entries:
`PhaseSpecced → PhaseDesigned → PhasePlanned` are all valid transitions.
No state machine changes are needed.

---

## Integration Points: New vs Modified (explicit)

### Modified ExecutionContext (internal/executor/context.go)

```go
// TaskItem — add dependency tracking
type TaskItem struct {
    ID          int      `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    Status      string   `json:"status"`
    Depends     []int    `json:"depends,omitempty"`   // NEW: task IDs this depends on
    Files       []string `json:"files,omitempty"`     // NEW: files this task touches
}

// ExecutionContext — add wave groups and worktree support
type ExecutionContext struct {
    // ...all existing fields preserved unchanged...
    WaveGroups  [][]TaskItem `json:"wave_groups,omitempty"`   // NEW: computed wave layers
    WorktreeDir string       `json:"worktree_dir,omitempty"`  // NEW: default ".worktrees"
    AutoMode    bool         `json:"auto_mode,omitempty"`     // NEW: skip interactive prompts
}
```

### Modified ProjectConfig (internal/config/defaults.go)

```go
type ProjectConfig struct {
    // ...all existing fields preserved unchanged...
    WorktreeDir string `yaml:"worktree_dir" mapstructure:"worktree_dir"` // NEW, default ".worktrees"
    AutoMode    bool   `yaml:"auto_mode" mapstructure:"auto_mode"`       // NEW, default false
}
```

Extended `DefaultModelMap` — add new roles (config.go):
```go
// Add to each profile tier (quality/balanced/budget):
"researcher":      "claude-sonnet-4-5" / "claude-sonnet-4-5" / "claude-haiku-3-5"
"advisor":         "claude-sonnet-4-5" / "claude-sonnet-4-5" / "claude-haiku-3-5"
"proposal-writer": "claude-sonnet-4-5" / "claude-sonnet-4-5" / "claude-haiku-3-5"
"plan-checker":    "claude-sonnet-4-5" / "claude-sonnet-4-5" / "claude-haiku-3-5"
```

### Modified TasksFrontmatterV2 (internal/spec/schema.go)

```go
type TaskEntry struct {
    ID          int        `yaml:"id"`
    Name        string     `yaml:"name"`
    Description string     `yaml:"description,omitempty"`
    Status      ItemStatus `yaml:"status"`
    Depends     []int      `yaml:"depends,omitempty"`   // NEW
    Files       []string   `yaml:"files,omitempty"`     // NEW
}
```

### New Binary Commands

| Command | Binary Behavior | SKILL.md Role |
|---------|-----------------|---------------|
| `mysd model` | Output current profile JSON (`--context-only`), or print formatted text | Display/format output |
| `mysd model set {profile}` | Write `model_profile` to `.claude/mysd.yaml` | Ask user which profile to set |
| `mysd model resolve {agent}` | Call `config.ResolveModel(agent, profile, overrides)`, output JSON | Display resolved model name |
| `mysd lang set {locale}` | Write `response_language`/`document_language` to mysd.yaml AND `locale` to openspec/config.yaml | Ask user for locale choice interactively |

### State Machine: No Changes Required

The existing `ValidTransitions` map is sufficient. `discuss` is a re-entry point, not a new phase. It runs existing transitions (`specced → designed → planned`) without adding new states.

```
Current ValidTransitions already handle all v1.1 flows:
PhaseSpecced  → PhaseDesigned   (re-design after discuss)
PhaseDesigned → PhasePlanned    (re-plan after discuss or plan-checker adds tasks)
PhasePlanned  → PhaseExecuted   (normal execution path, single or wave)
PhaseVerified → PhaseExecuted   (re-execute after failed verify — unchanged)
```

---

## Scaling Considerations

This is a CLI tool on a single developer machine. Scaling = codebase size + task parallelism, not user traffic.

| Scale | Architecture Adjustments |
|-------|--------------------------|
| Small repo (< 10k LOC) | No changes needed — all features work as designed |
| Large repo (> 100k LOC) | mysd-researcher needs bounded glob/grep scope; scanning all files causes token overflow. Use targeted directory patterns instead of recursive full-scan. |
| Many parallel tasks (> 10) | Worktree creation at `.worktrees/` can consume significant disk space temporarily. `internal/worktree/` must clean up aggressively on success. |
| Long session with discuss | `discovery-state.json` research summary grows stale across days. Add `research_ttl` field or `--force-research` flag in a later iteration. |

### Scaling Priorities

1. **First bottleneck:** Git worktree creation on Windows has 260-char path limit. Mitigated by short `.worktrees/T{id}/` paths (already in design).
2. **Second bottleneck:** N parallel mysd-advisor agents × large gray area lists can exhaust context tokens. Cap advisor parallelism at `config.AgentCount` or a sensible default (5).

---

## Anti-Patterns

### Anti-Pattern 1: Binary Driving Interactive Prompts

**What people do:** Add `fmt.Scanln()` or `bufio.Scanner` interactive input to cmd/ functions.

**Why it is wrong:** Breaks the binary-as-context-provider pattern. The binary is called with `--context-only` by SKILL.md scripts; reading stdin causes it to hang. Also makes unit testing impossible.

**Do this instead:** All interactive prompts live in SKILL.md orchestrators. The binary receives user decisions as CLI arguments. Example: `mysd lang set zh-TW` — the locale was selected by the SKILL.md asking the user, then passed as an argument.

### Anti-Pattern 2: Plan-Checker Reading Planner Artifacts Directly

**What people do:** mysd-plan-checker reads mysd-planner's output `tasks.md` independently, then cross-checks against spec files.

**Why it is wrong:** Breaks verifier independence principle established in Phase 3 (mysd-verifier does not read executor artifacts). Creates tight coupling between agents.

**Do this instead:** `mysd plan --context-only` outputs both the tasks AND the coverage gaps as JSON. mysd-plan-checker receives this structured JSON as its input context, not raw filesystem paths.

### Anti-Pattern 3: Worktree Branches From Diverged HEAD

**What people do:** Create wave 1 worktrees while wave 0 merges are still in progress, resulting in branches starting from different HEAD snapshots.

**Why it is wrong:** When wave 0 merges complete, wave 1 branches are now behind. Merging them requires rebasing, which complicates the AI auto-resolve flow.

**Do this instead:** All worktrees in a wave are created from the same HEAD snapshot. Complete the full wave (all merges resolved) THEN create the next wave's worktrees from the updated HEAD.

### Anti-Pattern 4: Embedding Research Metadata in proposal.md Frontmatter

**What people do:** Add `research_enabled: true`, `gray_areas: [...]`, and `research_summary: "..."` to the YAML frontmatter of proposal.md.

**Why it is wrong:** Research state is ephemeral orchestration context, not spec content. Pollutes OpenSpec format and breaks compatibility with other OpenSpec tools.

**Do this instead:** Use the sidecar `discovery-state.json` pattern (same as existing `verification-status.json`). proposal.md body references conclusions from research, but not the process metadata.

### Anti-Pattern 5: New Cobra Commands for Pure-Orchestration Flows

**What people do:** Add `cmd/discuss.go` with a full cobra command that drives the discussion loop.

**Why it is wrong:** Discussion is an AI-driven interactive loop that cannot be expressed as a stateless CLI invocation. Attempting to do so produces a command that either runs non-interactively (useless) or reads stdin (breaks the context-provider pattern).

**Do this instead:** discuss and fix are SKILL.md-only flows. If the discuss flow needs to read binary state, it calls `mysd status --context-only` or `mysd execute --context-only` (which already exist). Binary state mutations use existing commands (`mysd spec`, `mysd plan`).

---

## Build Order Recommendation

Dependencies flow downward. Build in this order to avoid blocked work:

**Phase 1 — Schema Foundation (no cross-package deps)**
1. `internal/spec/` — add `Depends`/`Files` to `TaskEntry` and `TasksFrontmatterV2`
2. `internal/config/` — add new agent roles to `DefaultModelMap`, add `WorktreeDir`/`AutoMode`
3. `internal/planchecker/` — new package, depends only on `internal/spec/` types

**Phase 2 — Executor Extension**
4. `internal/executor/` — add `BuildWaveGroups`, extend `ExecutionContext` with wave data
5. `cmd/execute.go` — consume wave groups, add `--auto` flag

**Phase 3 — Worktree Support**
6. `internal/worktree/` — new package (depends only on `os/exec` stdlib for git)
7. `plugin/agents/mysd-executor.md` — add worktree isolation instructions

**Phase 4 — New Binary Commands**
8. `cmd/model.go` — model profile read/write/resolve
9. `cmd/lang.go` — locale setting (writes two yaml files)
10. `cmd/scan.go` + `internal/scanner/` — language-agnostic refactor, scaffold-only flag
11. `cmd/plan.go` — output coverage JSON after state transition (uses planchecker)

**Phase 5 — New SKILL.md Orchestrators**
12. `plugin/commands/mysd-discuss.md`
13. `plugin/commands/mysd-fix.md`
14. `plugin/commands/mysd-model.md`
15. `plugin/commands/mysd-lang.md`

**Phase 6 — New Agent Definitions**
16. `plugin/agents/mysd-researcher.md`
17. `plugin/agents/mysd-advisor.md`
18. `plugin/agents/mysd-proposal-writer.md`
19. `plugin/agents/mysd-plan-checker.md`

**Phase 7 — Discovery Integration (depends on Phase 6 complete)**
20. `internal/discovery/` — DiscoveryContext, DiscoveryState, BuildDiscoveryContext
21. Modify `plugin/commands/mysd-propose.md` — add research mode flow
22. Modify `plugin/commands/mysd-spec.md` — add research mode flow
23. Modify `plugin/commands/mysd-ff.md` + `mysd-ffe.md` — auto mode, research-once pattern

**Phase 8 — Plan-Checker Integration (depends on Phase 3 + Phase 6 complete)**
24. Modify `plugin/commands/mysd-plan.md` — auto-spawn plan-checker after plan completes

---

## Sources

- Direct codebase analysis — v1.0 source at `/d/work_data/project/go/mysd` (HIGH confidence)
  - `internal/executor/context.go` — ExecutionContext schema
  - `internal/spec/schema.go` — TaskEntry, TasksFrontmatterV2
  - `internal/state/transitions.go` — ValidTransitions map
  - `internal/config/defaults.go` + `config.go` — ProjectConfig, DefaultModelMap
  - `plugin/commands/mysd-execute.md` — SKILL.md orchestration pattern (binary-as-context-provider)
  - `plugin/agents/mysd-executor.md` — subagent input contract
  - `plugin/commands/mysd-ff.md` — fast-forward orchestration + Task tool pattern
  - `plugin/agents/mysd-planner.md` — planner agent input contract
- `.specs/changes/interactive-discovery/proposal.md` — v1.1 feature specification (HIGH confidence)
- `.planning/PROJECT.md` — milestone context, constraints, key decisions (HIGH confidence)

---
*Architecture research for: mysd v1.1 — Interactive Discovery, Worktree Parallel Execution, Subagent Orchestration*
*Researched: 2026-03-25*
