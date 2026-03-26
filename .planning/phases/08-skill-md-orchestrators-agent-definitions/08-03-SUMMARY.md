---
phase: 08-skill-md-orchestrators-agent-definitions
plan: "03"
subsystem: plugin-commands
tags: [skill-md, orchestrators, auto-mode, apply, propose, status]
dependency_graph:
  requires: ["08-01", "08-02"]
  provides: ["mysd-plan-3stage", "mysd-apply-cmd", "mysd-execute-redirect", "mysd-propose-source-detection", "mysd-status-workflow"]
  affects: [".claude/commands/", "plugin/commands/"]
tech_stack:
  added: []
  patterns: ["3-stage pipeline orchestration", "per-task agent spawn", "source auto-detection", "auto_mode propagation"]
key_files:
  created:
    - .claude/commands/mysd-apply.md
    - plugin/commands/mysd-apply.md
  modified:
    - .claude/commands/mysd-plan.md
    - .claude/commands/mysd-execute.md
    - .claude/commands/mysd-propose.md
    - .claude/commands/mysd-status.md
    - plugin/commands/mysd-plan.md
    - plugin/commands/mysd-execute.md
    - plugin/commands/mysd-propose.md
    - plugin/commands/mysd-status.md
decisions:
  - "/mysd:plan redesigned as 3-stage pipeline: researcher(x4 parallel) -> designer -> planner, with optional research and check phases"
  - "execute renamed to apply at SKILL.md layer; /mysd:execute preserved as redirect file pointing to /mysd:apply"
  - "/mysd:apply spawns mysd-executor per task (not all tasks to one executor); single=sequential, wave=parallel within wave_groups"
  - "--auto flag parsed at SKILL.md layer; auto_mode propagated as context JSON field to all spawned agents"
  - "propose source detection uses 6-priority order: change-name > file-path > dir-path > active-change > gstack/conversation > create-new"
metrics:
  duration: "~12 min"
  completed_date: "2026-03-26"
  tasks_completed: 2
  files_modified: 10
---

# Phase 08 Plan 03: SKILL.md Orchestrators Rewrite Summary

Rewrote 4 existing SKILL.md orchestrators (plan, apply/execute, propose, status) and added --auto flag support across all commands.

## What Was Built

### /mysd:plan — 3-Stage Pipeline (D-23)

Redesigned the plan orchestrator from a single-agent spawn to a 3-stage pipeline:

1. **Research Phase** (optional, `--research` flag or `research_enabled`): Spawns 4 `mysd-researcher` agents in parallel, one per dimension (codebase, domain, architecture, pitfalls). Research outputs feed into the design phase.
2. **Design Phase**: Spawns `mysd-designer` via Task tool to produce `design.md`, then runs `mysd design` state transition.
3. **Planning Phase**: Spawns `mysd-planner` via Task tool with full context (specs + research + design), then runs `mysd plan` state transition.
4. **Check Phase** (optional, `--check` flag): Spawns `mysd-plan-checker` to validate MUST coverage.

`--auto` flag parsed in Step 1, propagated as `auto_mode` field in all agent contexts.

### /mysd:apply — Per-Task Spawn Orchestrator (D-22, FAGENT-07)

Created new `/mysd:apply` command implementing per-task agent spawning:

- **Single Mode**: Sequential spawn — one `mysd-executor` per task, wait for completion before next.
- **Wave Mode**: Parallel spawn within each wave group, sequential across waves, with merge step after each wave (ascending task ID, `--no-ff`, 3-retry AI conflict resolution).

Binary subcommand remains `mysd execute --context-only` per RESEARCH.md recommendation.

### /mysd:execute — Redirect

Preserved `mysd-execute.md` as a redirect file pointing users to `/mysd:apply` with explanation of the rename.

### /mysd:propose — Source Detection + proposal-writer Spawn (D-27, D-06/D-07)

Redesigned propose with 6-priority source detection:
1. `$ARGUMENTS` matches `.specs/changes/{name}/` → change mode
2. `$ARGUMENTS` is file path → single file mode
3. `$ARGUMENTS` is directory path → selection mode (multi-select or auto-all if auto_mode)
4. No arg + active change → use current change
5. No arg + no active change → auto-detect from `~/.gstack/projects/`, conversation context (excludes `.claude/plans/`)
6. Nothing found → create new (ask or auto-generate if auto_mode)

Spawns `mysd-proposal-writer` via Task tool with source content as `conclusions`.

### /mysd:status — Workflow Dashboard (D-29/30/31)

Redesigned status display with:
- Workflow stage indicator: `propose > plan > apply > archive` with current position marker
- Task list with status symbols: `done` / `failed` / `skipped` / `pending` / `in_progress`
- Next step recommendation based on current workflow stage

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None — all commands are fully specified orchestrators. Actual execution depends on the binary (`mysd execute --context-only`, `mysd plan --context-only`, `mysd status`) and spawned agent definitions.

## Commits

- `1b8829a`: feat(08-03): rewrite /mysd:plan with 3-stage pipeline and --auto support
- `b1d3794`: feat(08-03): create apply, redirect execute, rewrite propose and status

## Self-Check: PASSED

All created/modified files verified to exist and contain required content.
