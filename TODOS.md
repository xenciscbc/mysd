# TODOS

## ~~Update design doc with eng review decisions~~ ✅ DONE (2026-04-16)
- **What:** Update `~/.gstack/projects/xenciscbc-mysd/cbc-master-design-20260416-113346.md` to reflect 5 eng review decisions
- **Why:** Design doc is the single source of truth for implementation. Must reflect latest decisions.
- **Context:** 5 decisions from eng review:
  1. Analyzer 4 dimensions → absorbed into research skill (not spec writer)
  2. Only support `openspec/` directory format (drop `.specs/` compatibility)
  3. Delete Phase 0 PoC (Phase 2 validation covers this)
  4. Add explicit trigger boundaries in each SKILL.md description (排他邊界)
  5. Architecture change: 3 skills + 1 orchestrator (orchestrator uses subagents to chain skills)
- **Depends on:** Nothing — can be done immediately
- **Added:** 2026-04-16 by /plan-eng-review

## ~~Add 6 edge case test scenarios to Phase 2 validation table~~ ✅ DONE (2026-04-16)
- **What:** Add 6 boundary test scenarios to the Phase 2 validation matrix in the design doc
- **Why:** Current validation table covers happy paths only (60%). Edge cases are untested.
- **Context:** 6 gaps:
  1. research: Spec Health Check (analyzer 4 dimensions)
  2. research: Non-gray-area question (should reject or redirect)
  3. doc: Unmapped change type (falls back to grep heuristic)
  4. doc: Empty diff (no changes to process)
  5. spec: RENAMED delta operation
  6. spec: Incompatible spec-version handling
- **Depends on:** Design doc update (TODO #1) should be done first
- **Added:** 2026-04-16 by /plan-eng-review
