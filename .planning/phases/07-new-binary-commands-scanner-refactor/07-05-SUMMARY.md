---
phase: 07-new-binary-commands-scanner-refactor
plan: 05
subsystem: plugin
tags: [skill-md, claude-code-plugin, mysd-scan, mysd-init, mysd-model, mysd-lang, planner-agent, skills-recommendation]

# Dependency graph
requires:
  - phase: 07-02
    provides: Language-agnostic ScanContext struct with primary_language/files/modules
  - phase: 07-03
    provides: mysd model and model set commands
  - phase: 07-04
    provides: mysd lang and lang set commands with atomic dual-config write

provides:
  - Updated /mysd:scan SKILL.md referencing new language-agnostic JSON format
  - Rewritten /mysd:init SKILL.md with scaffold-only + mysd lang set locale flow
  - New /mysd:model SKILL.md for model profile display and switching
  - New /mysd:lang SKILL.md for language settings display and BCP47 locale configuration
  - Updated mysd-planner agent with Step 4.5 skills recommendation and Step 7.5 skills confirmation flow

affects:
  - Phase 08 (execute workflow uses skills field in tasks.md)
  - Any agent reading tasks.md YAML (skills field now present)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "SKILL.md orchestrator pattern: thin files invoke binary via Bash, delegate complex logic to agents"
    - "Skills recommendation heuristics: spec/design/code/verify/scan/capture task type detection"
    - "auto_mode flag in planner context controls interactive vs ffe skills confirmation"

key-files:
  created:
    - .claude/commands/mysd-model.md
    - .claude/commands/mysd-lang.md
  modified:
    - .claude/commands/mysd-scan.md
    - .claude/commands/mysd-init.md
    - .claude/agents/mysd-planner.md

key-decisions:
  - "Skills recommendation logic in mysd-planner agent layer, confirmation in SKILL.md layer, default accept-all (D-07, D-08, D-09)"
  - "auto_mode true skips skills confirmation entirely in ffe mode (D-10)"
  - "skills field added to tasks.md YAML as empty array default — backward compatible"
  - "init SKILL.md completely rewritten: scaffold-only + mysd lang set, removes old model_profile/execution_mode config explanation"
  - "scan SKILL.md: modules replace packages as the universal unit, primary_language shown prominently"
  - "locale setup prompt added to scan SKILL.md when config_exists is false (D-06/FSCAN-04)"

patterns-established:
  - "SKILL.md format: frontmatter model + allowed-tools + numbered steps + Bash invocation"
  - "BCP47 locale code used for both response_language and openspec locale config"

requirements-completed:
  - SKILL-01
  - SKILL-02
  - SKILL-03
  - SKILL-04

# Metrics
duration: 5min
completed: 2026-03-26
---

# Phase 7 Plan 05: SKILL.md Plugin Layer Update Summary

**Five plugin files updated: scan/init SKILL.md rewritten for language-agnostic scanner, model/lang SKILL.md created, planner agent gains skills recommendation with interactive/ffe confirmation flow**

## Performance

- **Duration:** ~5 min
- **Started:** 2026-03-26T01:57:10Z
- **Completed:** 2026-03-26T02:02:16Z
- **Tasks:** 3 (tasks 1-3 complete; task 4 = checkpoint:human-verify)
- **Files modified:** 5

## Accomplishments

- Updated mysd-scan.md to reference new language-agnostic ScanContext fields (primary_language, files, modules) — removes all Go-specific PackageInfo references
- Rewrote mysd-init.md to scaffold-only flow: `mysd init` creates scaffold, then `mysd lang set` configures locale — removes old interactive config editing
- Created mysd-model.md and mysd-lang.md as new SKILL.md files wiring to the binary commands from Plans 03 and 04
- Updated mysd-planner agent with Step 4.5 (skills recommendation heuristics) and Step 7.5 (interactive confirmation table with Accept all? Y/n default, ffe bypass)

## Task Commits

Each task was committed atomically:

1. **Task 1: Update mysd-scan.md and mysd-init.md** — `214a111` (feat)
2. **Task 2: Create mysd-model.md and mysd-lang.md** — `58ad3c2` (feat)
3. **Task 3: Update mysd-planner with skills recommendation** — `9e74e64` (feat)

**Plan metadata:** (pending — will commit with SUMMARY.md after human verify)

## Files Created/Modified

- `.claude/commands/mysd-scan.md` — Updated to language-agnostic format: primary_language, files map, modules array, config_exists locale prompt
- `.claude/commands/mysd-init.md` — Rewritten: scaffold-only init + mysd lang set + removed old config editing steps
- `.claude/commands/mysd-model.md` — New: show profile + 10-role table, offer mysd model set with profile explanation
- `.claude/commands/mysd-lang.md` — New: show language settings, BCP47 locale options, mysd lang set with atomic confirmation
- `.claude/agents/mysd-planner.md` — Added Step 4.5 skills recommendation + Step 7.5 skills confirmation (interactive/ffe)

## Decisions Made

- Skills recommendation heuristics defined: spec artifacts → /mysd:propose or /mysd:spec; design → /mysd:design; code → [] (none); testing → /mysd:verify; scanning → /mysd:scan; capturing → /mysd:capture
- Default accept-all behavior for skills confirmation (Enter = Y) per D-09
- ffe mode (auto_mode=true) skips confirmation entirely per D-10
- locale setup prompt in scan SKILL.md when config_exists=false per D-06/FSCAN-04

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## Next Phase Readiness

- Phase 7 is complete pending human verification (Task 4 checkpoint)
- Human verifier should: run `go test ./...`, `go build`, `./mysd.exe model`, `./mysd.exe scan --context-only`, and visually confirm SKILL.md files exist
- Phase 8 can proceed once verification passes

---
*Phase: 07-new-binary-commands-scanner-refactor*
*Completed: 2026-03-26*
