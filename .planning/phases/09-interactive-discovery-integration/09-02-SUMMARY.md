---
phase: 09-interactive-discovery-integration
plan: 02
subsystem: plugin-commands
tags: [skill-md, orchestrator, discovery, research, gray-areas, advisor, dual-loop, scope-guardrail]
dependency_graph:
  requires: [mysd-researcher agent, mysd-advisor agent, mysd-proposal-writer agent]
  provides: [propose discovery pipeline, discuss discovery pipeline]
  affects: [.claude/commands/mysd-propose.md, .claude/commands/mysd-discuss.md]
tech_stack:
  added: []
  patterns: [parallel Task tool spawn, dual-loop exploration, scope guardrail with deferred notes]
key_files:
  created: []
  modified:
    - .claude/commands/mysd-propose.md
    - .claude/commands/mysd-discuss.md
decisions:
  - "propose always loads deferred notes (D-02) — cross-change context valuable for new proposals"
  - "discuss conditionally loads deferred notes based on active WIP check (D-02) — avoid polluting focused WIP discussion"
  - "dual-loop uses user-driven termination not numeric quota (D-01) — binary choice per area, no artificial limits"
  - "advisors spawned at orchestrator layer only — FAGENT-05 compliance, not inside researcher agents"
  - "auto_mode skips research entirely in both commands (FAUTO-02) — ff/ffe-style auto means no interaction"
metrics:
  duration: 14min
  completed: "2026-03-26T06:45:46Z"
  tasks_completed: 2
  files_modified: 2
---

# Phase 09 Plan 02: Propose & Discuss Discovery Pipeline Summary

Rewrote `/mysd:propose` and `/mysd:discuss` SKILL.md orchestrators with the full interactive discovery pipeline: opt-in 4-dimension parallel research, gray area identification, parallel advisor spawning, dual-loop exploration with user-driven termination, and scope guardrail writing to deferred notes.

## Tasks Completed

### Task 1: Rewrite mysd-propose.md with discovery pipeline
**Commit:** `1793354`

Added 6 new steps (Steps 4-9 in the new structure) to the existing 5-step propose orchestrator:

- **Step 4** — Load deferred notes via `mysd note list` (D-02: always loads for new proposals)
- **Step 5** — Optional research prompt: `Would you like to run 4-dimension research? [y/N]` — skipped entirely when `auto_mode=true`
- **Step 6** — Parallel spawn of 4 `mysd-researcher` agents (codebase / domain / architecture / pitfalls)
- **Step 7** — Gray area identification from research output + parallel `mysd-advisor` spawn per area at orchestrator layer
- **Step 8** — Dual-loop exploration: Layer 1 (per-area deep dive with scope guardrail) + Layer 2 (new area discovery after all areas explored)
- **Step 9** — Updated proposal writer invocation to include `deferred_context` field
- **Step 10** — Updated confirm to show research/exploration/deferred summary

### Task 2: Rewrite mysd-discuss.md with discovery pipeline
**Commit:** `686e66d`

Restructured the 8-step discuss orchestrator to 12 steps with identical discovery pipeline pattern:

- **Step 4** — Conditional deferred notes loading: `mysd status` check first — active WIP = skip notes, no active WIP = load notes (D-02 context-aware)
- **Steps 5-8** — Same discovery pipeline as propose: optional research prompt, parallel researchers, gray area + advisor spawning, dual-loop
- **Step 9** — Updated discussion loop to incorporate research/exploration findings
- **Steps 10-11** — Preserved existing spec update and re-plan logic (renumbered from old Steps 6-7)
- **Step 12** — Updated confirm with research/exploration/deferred summary

## Deviations from Plan

None - plan executed exactly as written.

## Decisions Made

| Decision | Rationale |
|----------|-----------|
| propose always loads deferred notes (D-02) | Cross-change context is valuable when starting a new proposal — past ideas can inform new direction |
| discuss conditionally loads deferred notes (D-02) | Active WIP discussion should stay focused — deferred notes from past changes would pollute current scope |
| User-driven dual-loop termination (D-01) | Binary choice per area is sufficient termination signal — no quota counter needed |
| Advisors at orchestrator layer only | FAGENT-05 compliance — leaf agents (researchers) must not spawn other agents |
| auto_mode skips research entirely | FAUTO-02 — ff/ffe-style automation means zero interaction, research is inherently interactive |

## Key Links Implemented

| Source | To | Via | Pattern |
|--------|----|-----|---------|
| mysd-propose.md | mysd-researcher | Task tool x4 parallel | `Agent: mysd-researcher` |
| mysd-propose.md | mysd-advisor | Task tool per gray area | `Agent: mysd-advisor` |
| mysd-discuss.md | mysd-researcher | Task tool x4 parallel | `Agent: mysd-researcher` |
| mysd-discuss.md | mysd-advisor | Task tool per gray area | `Agent: mysd-advisor` |

## Known Stubs

None — both orchestrators reference real agent definitions that were created in Phase 8.

## Self-Check: PASSED

- `1793354` — feat(09-02): rewrite mysd-propose.md with full discovery pipeline
- `686e66d` — feat(09-02): rewrite mysd-discuss.md with full discovery pipeline
- `.claude/commands/mysd-propose.md` — exists, contains mysd-researcher (2), mysd-advisor (3), gray area (9), mysd note add (1), Layer 1/2, In Scope/Out of Scope
- `.claude/commands/mysd-discuss.md` — exists, contains mysd-researcher (2), mysd-advisor (3), gray area (10), mysd note (2), Layer 1/2, mysd status (D-02 conditional check)
