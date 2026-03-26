---
phase: 08-skill-md-orchestrators-agent-definitions
plan: 01
subsystem: plugin
tags: [claude-code, agent-definitions, researcher, advisor, proposal-writer, plan-checker]

# Dependency graph
requires:
  - phase: 05-schema-foundation-plan-checker
    provides: mysd-plan-checker agent in plugin/agents/ (source for sync)
  - phase: 07-new-binary-commands-scanner-refactor
    provides: skills recommendation pattern in mysd-planner.md (auto_mode context field)
provides:
  - mysd-researcher agent definition (4-dimension research: codebase/domain/architecture/pitfalls)
  - mysd-advisor agent definition (gray area analysis with comparison table output)
  - mysd-proposal-writer agent definition (creates/updates proposal.md from conclusions)
  - mysd-plan-checker synced to .claude/agents/ (was missing from dev directory)
affects: [08-02, 08-03, 08-04, 08-05]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "leaf agent constraint: no Task tool in allowed-tools, no subagent spawning"
    - "auto_mode: boolean context field pattern for skipping interactive confirmations"
    - "plugin distribution pattern: .claude/agents/ and plugin/agents/ hold identical content"

key-files:
  created:
    - .claude/agents/mysd-researcher.md
    - .claude/agents/mysd-advisor.md
    - .claude/agents/mysd-proposal-writer.md
    - .claude/agents/mysd-plan-checker.md
    - plugin/agents/mysd-researcher.md
    - plugin/agents/mysd-advisor.md
    - plugin/agents/mysd-proposal-writer.md
  modified: []

key-decisions:
  - "mysd-researcher uses WebFetch for domain/pitfalls dimensions, Bash for codebase/architecture — tool selection matches dimension type"
  - "mysd-advisor produces MECE comparison table (2-4 options) with explicit trade-off acceptance rationale"
  - "mysd-proposal-writer handles both create (new) and update (merge) modes based on existing_proposal field"
  - "mysd-plan-checker synced verbatim from plugin/agents/ — no model field (uses default) to preserve existing behavior"
  - "All 4 agents: zero Task tool references, zero subagent spawning — enforces D-17 leaf agent constraint"

patterns-established:
  - "Leaf agent pattern: allowed-tools lists only direct-action tools (Read/Write/Edit/Grep/Glob/Bash/WebFetch), never Task"
  - "auto_mode field: agents receive boolean in context JSON and branch behavior (interactive vs direct) accordingly"
  - "Plugin sync: .claude/agents/ is the authoritative copy; plugin/agents/ is the distribution copy with identical content"

requirements-completed:
  - FAGENT-01
  - FAGENT-02
  - FAGENT-03

# Metrics
duration: 8min
completed: 2026-03-26
---

# Phase 08 Plan 01: Agent Definitions (Researcher, Advisor, Proposal-Writer, Plan-Checker Sync) Summary

**4 agent definitions created — mysd-researcher (4D research), mysd-advisor (trade-off comparison tables), mysd-proposal-writer (proposal.md writer), and mysd-plan-checker synced to .claude/agents/ — all with zero Task tool references**

## Performance

- **Duration:** 8 min
- **Started:** 2026-03-26T04:29:58Z
- **Completed:** 2026-03-26T04:38:00Z
- **Tasks:** 1
- **Files modified:** 7

## Accomplishments

- Created mysd-researcher agent with 4-dimension research capability (codebase/domain/architecture/pitfalls), each dimension using appropriate tools (WebFetch for external research, Bash/Grep for codebase)
- Created mysd-advisor agent with structured comparison table output format (MECE options, explicit trade-off acceptance, auto_mode branch for ff/ffe workflows)
- Created mysd-proposal-writer agent handling both create-new and update-existing modes, with state transition integration (`mysd propose {change_name}`)
- Synced mysd-plan-checker from plugin/agents/ to .claude/agents/ (verbatim copy, fixing gap identified in RESEARCH.md Open Question 2)
- Distributed all 3 new agents to plugin/agents/ with identical content (plugin distribution copies)

## Task Commits

1. **Task 1: Create 3 new agent definitions + sync plan-checker** - `b7f9efb` (feat)

## Files Created/Modified

- `.claude/agents/mysd-researcher.md` - Research agent for 4 dimensions (codebase/domain/architecture/pitfalls)
- `.claude/agents/mysd-advisor.md` - Trade-off analysis agent with comparison table output
- `.claude/agents/mysd-proposal-writer.md` - Proposal.md writer (create + update modes)
- `.claude/agents/mysd-plan-checker.md` - Synced from plugin/agents/ (MUST coverage agent)
- `plugin/agents/mysd-researcher.md` - Distribution copy (identical to .claude/agents/)
- `plugin/agents/mysd-advisor.md` - Distribution copy (identical to .claude/agents/)
- `plugin/agents/mysd-proposal-writer.md` - Distribution copy (identical to .claude/agents/)

## Decisions Made

- mysd-researcher: WebFetch included in allowed-tools to support domain and pitfalls research dimensions (requires web lookup for best practices and CVEs); not included in advisor since analysis is internal
- mysd-plan-checker: synced verbatim without adding `model:` field — preserves existing behavior and avoids silent semantic change
- Researcher body is read-only (no Write/Edit) — enforces research-only contract
- Advisor body is analysis-only (no Write/Edit) — ensures output is always user-reviewed before acting
- Proposal-writer runs `mysd propose {change_name}` state transition only when creating new (not updating) — avoids double-transition

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- All 3 new agents (researcher, advisor, proposal-writer) are ready to be referenced by SKILL.md orchestrators in Plans 02-05
- mysd-plan-checker is now available in both .claude/agents/ and plugin/agents/
- Plan 02 (SKILL.md orchestrators) can now spawn these agents via Task tool in the SKILL.md layer
- Blocker from STATE.md: "All 9 agent definitions require manual audit for Task tool references" — this plan adds 4 agents all confirmed clean; remaining 5 existing agents need audit in Plan 03 or 04

---
*Phase: 08-skill-md-orchestrators-agent-definitions*
*Completed: 2026-03-26*
