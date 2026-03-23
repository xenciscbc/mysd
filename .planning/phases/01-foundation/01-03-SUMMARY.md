---
phase: 01-foundation
plan: 03
subsystem: cli
tags: [cobra, viper, go, cli, spec, state]

requires:
  - phase: 01-01
    provides: spec.Scaffold, spec.DetectSpecDir, spec.Change type
  - phase: 01-02
    provides: state.Transition, state.LoadState, state.SaveState, config.Defaults, output.NewPrinter

provides:
  - Executable mysd binary with cobra CLI framework
  - mysd propose command: scaffolds .specs/changes/{name}/ and transitions state to PhaseProposed
  - mysd init command: creates .claude/mysd.yaml from config.Defaults()
  - Stub commands for Phase 2/3: spec, design, plan, execute, verify, archive, status
  - Persistent flags bound to viper: --execution-mode, --agent-count, --lang, --doc-lang, --tdd, --atomic-commits

affects:
  - 02-spec-commands
  - 03-execute-engine
  - 04-distribution

tech-stack:
  added: []
  patterns:
    - "Thin cmd layer: all business logic delegated to internal/ packages"
    - "Cobra PersistentFlags bound to global viper instance for flag override of config"
    - "TDD: test RED commit before implementation GREEN commit"
    - "cmd/init_cmd.go naming convention to avoid collision with Go init() function semantics"

key-files:
  created:
    - main.go
    - cmd/root.go
    - cmd/propose.go
    - cmd/init_cmd.go
    - cmd/spec.go
    - cmd/design.go
    - cmd/plan.go
    - cmd/execute.go
    - cmd/verify.go
    - cmd/archive.go
    - cmd/status.go
    - cmd/root_test.go
    - cmd/propose_test.go
    - cmd/init_cmd_test.go
    - Makefile
  modified: []

key-decisions:
  - "File named init_cmd.go (not init.go) to avoid confusion with Go's init() function convention"
  - "propose uses .specs as default specDir when DetectSpecDir returns ErrNoSpecDir — new project bootstrapping"
  - "Global viper instance used in root.go for cobra PersistentFlags binding (not instance viper) — root cmd is singleton"
  - "state.Transition error on re-propose handled gracefully: set phase directly to allow idempotent re-runs"

patterns-established:
  - "Pattern 1 (Thin commands): cmd/*.go files are pure wiring — parse args, call internal/, print results"
  - "Pattern 2 (TDD): test RED commit with test(01-03) prefix before GREEN feat(01-03) commit"
  - "Pattern 3 (Printer): always use output.NewPrinter(cmd.OutOrStdout()) for testable output"

requirements-completed:
  - SPEC-01
  - CONF-01
  - CONF-04
  - DIST-01
  - DIST-02
  - STAT-01

duration: 4min
completed: 2026-03-23
---

# Phase 01 Plan 03: CLI Skeleton Summary

**Cobra CLI binary with mysd propose (spec scaffold + state transition) and mysd init (config bootstrap), plus 7 Phase 2/3 stub commands — all wired to internal/ packages via thin cmd layer**

## Performance

- **Duration:** ~4 min
- **Started:** 2026-03-23T08:58:18Z
- **Completed:** 2026-03-23T08:59:09Z
- **Tasks:** 2 (+ 1 TDD test commit)
- **Files modified:** 15

## Accomplishments

- Single Go binary compiles cross-platform (`go build -o mysd.exe .` exits 0)
- `mysd propose test-feature` creates complete .specs/changes/test-feature/ with proposal.md, design.md, tasks.md, specs/ and saves STATE.json with phase="proposed"
- `mysd init` creates .claude/mysd.yaml with all preference fields; warns on existing (--force to overwrite)
- All 7 Phase 2/3 commands registered and return "not yet implemented"
- Persistent flags (--execution-mode, --lang, etc.) bound to viper for config override
- Full test suite: 5 packages, all pass

## Task Commits

Each task was committed atomically:

1. **Task 1: Cobra CLI skeleton** - `75b9fd4` (feat)
2. **Task 2 RED: Failing tests** - `1a8ac99` (test)
3. **Task 2 GREEN: propose + init implementation** - `c8efa2a` (feat)

**Plan metadata:** (docs commit follows)

_Note: TDD tasks have separate test (RED) and implementation (GREEN) commits_

## Files Created/Modified

- `main.go` - Binary entry point calling cmd.Execute()
- `cmd/root.go` - Root cobra command, Execute(), persistent flags, viper binding, initConfig
- `cmd/propose.go` - mysd propose: calls spec.Scaffold, state.Transition, output.Printer
- `cmd/init_cmd.go` - mysd init: writes .claude/mysd.yaml from config.Defaults() with --force flag
- `cmd/spec.go` - Stub: "not yet implemented"
- `cmd/design.go` - Stub: "not yet implemented"
- `cmd/plan.go` - Stub: "not yet implemented"
- `cmd/execute.go` - Stub: "not yet implemented"
- `cmd/verify.go` - Stub: "not yet implemented"
- `cmd/archive.go` - Stub: "not yet implemented"
- `cmd/status.go` - Stub: "not yet implemented"
- `cmd/root_test.go` - Tests: help contains propose/init
- `cmd/propose_test.go` - Tests: scaffold, success msg, no-args error, STATE.json phase
- `cmd/init_cmd_test.go` - Tests: config creation, no-overwrite, --force overwrite
- `Makefile` - build, test, lint, clean targets

## Decisions Made

- `init_cmd.go` (not `init.go`) to avoid confusion with Go's `init()` function semantics
- Default specDir falls back to `.specs` when `DetectSpecDir` returns `ErrNoSpecDir` — enables first-time `mysd propose` without prior `mysd init`
- Global viper used in root.go (not instance viper) — root cmd is a package-level singleton, so global viper is appropriate here
- `state.Transition` error on re-propose handled gracefully by setting phase directly, allowing idempotent re-runs

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

- Phase 02 can immediately wire `mysd spec`, `mysd design`, `mysd plan` using the same thin-command pattern established here
- `rootCmd` is package-level and accessible to all cmd/*.go files via the same package
- Blocker noted from Phase 01-02 still applies: Claude Code subagent invocation API from Go binary not yet pinned

## Self-Check: PASSED

- FOUND: main.go
- FOUND: cmd/root.go
- FOUND: cmd/propose.go
- FOUND: cmd/init_cmd.go
- FOUND: Makefile
- FOUND commit: 75b9fd4 (Task 1 — CLI skeleton)
- FOUND commit: 1a8ac99 (Task 2 RED — failing tests)
- FOUND commit: c8efa2a (Task 2 GREEN — propose + init)

---
*Phase: 01-foundation*
*Completed: 2026-03-23*
