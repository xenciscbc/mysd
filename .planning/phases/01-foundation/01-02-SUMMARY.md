---
phase: 01-foundation
plan: 02
subsystem: state
tags: [go, state-machine, viper, lipgloss, tty-detection, config]

requires: []
provides:
  - WorkflowState struct with JSON persistence to .specs/STATE.json
  - Phase constants (None/Proposed/Specced/Designed/Planned/Executed/Verified/Archived)
  - ValidTransitions map enforcing spec lifecycle order
  - CanTransition / Transition functions with ErrInvalidTransition error
  - ProjectConfig struct with all preference fields
  - Load() with instance viper, .claude/mysd.yaml, convention-over-config defaults
  - BindFlags() for cobra flag override integration
  - Printer with TTY detection via charmbracelet/x/term
  - Styled lipgloss output in TTY, plain prefixed text in non-TTY
affects: [cmd, integration, verification, execution-engine]

tech-stack:
  added:
    - github.com/charmbracelet/x/term v0.2.1 (TTY detection)
    - github.com/spf13/viper v1.21.0 (config loading, already in go.mod)
    - github.com/charmbracelet/lipgloss v1.1.0 (styling, already in go.mod)
  patterns:
    - Instance viper (not global) for testable config loading
    - Convention-over-config: missing STATE.json or config file returns zero-value/defaults, never error
    - TTY degradation: lipgloss in terminal, plain prefixed text in pipes/CI
    - ErrInvalidTransition wrapped with fmt.Errorf for errors.Is compatibility

key-files:
  created:
    - internal/state/state.go
    - internal/state/transitions.go
    - internal/state/state_test.go
    - internal/state/transitions_test.go
    - internal/config/config.go
    - internal/config/defaults.go
    - internal/config/config_test.go
    - internal/output/printer.go
    - internal/output/colors.go
    - internal/output/printer_test.go
  modified: []

key-decisions:
  - "Instance viper (viper.New()) instead of global viper for full test isolation — each Load() call gets fresh state"
  - "charmbracelet/x/term for TTY detection instead of golang.org/x/term — already in go.mod as lipgloss transitive dependency"
  - "ErrInvalidTransition defined in state.go (not transitions.go) to keep error definitions co-located with the type"
  - "BindFlags exported as standalone function, not called from Load() — allows cmd/root.go to wire flag overrides at init time"

patterns-established:
  - "Pattern: Convention-over-config — missing files return zero-values/defaults, not errors"
  - "Pattern: Wrap sentinel errors with fmt.Errorf + %w for errors.Is compatibility across callers"
  - "Pattern: Use instance viper (viper.New()) in tests to avoid global state pollution between test cases"
  - "Pattern: TTY detection by interface assertion (Fd() uintptr) — no hard dependency on os.File"

requirements-completed: [STAT-01, STAT-02, STAT-03, CONF-01, CONF-02, CONF-03, CONF-04]

duration: 5min
completed: 2026-03-23
---

# Phase 01 Plan 02: State Machine, Config Management, and Terminal Printer Summary

**Go state machine enforcing 8-phase spec lifecycle with ValidTransitions map, Viper-backed ProjectConfig with convention-over-config defaults, and TTY-aware Printer using lipgloss styles**

## Performance

- **Duration:** ~5 min
- **Started:** 2026-03-23T08:47:23Z
- **Completed:** 2026-03-23T08:51:54Z
- **Tasks:** 2
- **Files modified:** 10

## Accomplishments

- State machine (`internal/state`) enforces the 8-phase OpenSpec lifecycle (None → Proposed → Specced → Designed → Planned → Executed → Verified → Archived), blocks invalid transitions, persists to `.specs/STATE.json`
- Config package (`internal/config`) loads `.claude/mysd.yaml` via instance Viper with convention-over-config defaults, and exports `BindFlags` for cobra flag override wiring
- Output package (`internal/output`) provides TTY-detecting Printer — lipgloss-styled in terminal, plain prefixed text (`[OK]`, `[ERROR]`, `[WARN]`, `[INFO]`, `===`) in pipe/CI mode

## Task Commits

Each task was committed atomically:

1. **TDD RED — State/Config tests** - `a397aa1` (test)
2. **Task 1: State machine + config** - `f2106e6` (feat)
3. **Task 2: Terminal output printer** - `cae831d` (feat)

_Note: Task 1 used TDD — test commit (RED) followed by implementation commit (GREEN)_

## Files Created/Modified

- `internal/state/state.go` - Phase constants, WorkflowState struct, LoadState/SaveState with JSON
- `internal/state/transitions.go` - ValidTransitions map, CanTransition, Transition with ErrInvalidTransition
- `internal/state/state_test.go` - Zero-state on missing file, roundtrip, JSON marshaling, field naming
- `internal/state/transitions_test.go` - All valid transitions, 6 invalid transitions, phase/LastRun updates
- `internal/config/defaults.go` - ProjectConfig struct with all 7 preference fields, Defaults() function
- `internal/config/config.go` - Load() with instance viper, BindFlags() for cobra integration
- `internal/config/config_test.go` - Defaults values, no-file defaults, full override, partial override
- `internal/output/colors.go` - 5 lipgloss color constants and 6 style variables
- `internal/output/printer.go` - Printer struct, NewPrinter with TTY detection, 7 output methods
- `internal/output/printer_test.go` - Non-TTY prefix tests, Printf format, no ANSI escapes in non-TTY

## Decisions Made

- Used instance `viper.New()` instead of global viper to ensure test isolation (each test call to `Load()` starts clean)
- Used `github.com/charmbracelet/x/term` for `IsTerminal()` — already a transitive dependency via lipgloss, avoids adding `golang.org/x/term` as a new direct dependency
- `ErrInvalidTransition` is defined in `state.go` alongside the `WorkflowState` type, keeping error definitions co-located with the domain type they describe
- `BindFlags` is exported as a standalone function (not called automatically in `Load()`) so `cmd/root.go` can wire cobra flags to viper at the appropriate init-time hook

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - all packages compiled and all tests passed on first run.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- `internal/state`, `internal/config`, `internal/output` packages are complete and fully tested
- `BindFlags` is ready for `cmd/root.go` to call during cobra init (Plan 03)
- State machine phase constants are the authoritative source for all workflow state throughout the binary
- Printer can be instantiated with any `io.Writer` — compatible with both `os.Stdout` and test `bytes.Buffer`

---
*Phase: 01-foundation*
*Completed: 2026-03-23*

## Self-Check: PASSED

All 11 files confirmed present. All 3 task commits confirmed in git log.
