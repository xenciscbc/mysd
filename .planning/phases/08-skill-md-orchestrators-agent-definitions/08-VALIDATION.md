---
phase: 8
slug: skill-md-orchestrators-agent-definitions
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-26
---

# Phase 8 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Manual verification (plugin layer — no Go unit tests) |
| **Config file** | none |
| **Quick run command** | `claude -p "test mysd skill invocation"` (manual) |
| **Full suite command** | Manual audit of all 9 agent definitions |
| **Estimated runtime** | ~5 minutes (manual walkthrough) |

---

## Sampling Rate

- **After every task commit:** Verify SKILL.md syntax is valid markdown with correct frontmatter
- **After every plan wave:** Manually invoke affected skill to confirm it runs without errors
- **Before `/gsd:verify-work`:** Full audit of all agent definitions for nested Task tool usage
- **Max feedback latency:** 5 minutes

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 8-01-01 | 01 | 1 | FAGENT-01 | manual | `ls .claude/agents/mysd-discuss*.md` | ✅ / ❌ W0 | ⬜ pending |
| 8-01-02 | 01 | 1 | FAGENT-02 | manual | `ls .claude/agents/mysd-fix*.md` | ✅ / ❌ W0 | ⬜ pending |
| 8-02-01 | 02 | 1 | FCMD-01 | manual | `ls .claude/commands/mysd-discuss.md` | ✅ / ❌ W0 | ⬜ pending |
| 8-02-02 | 02 | 1 | FCMD-02 | manual | `ls .claude/commands/mysd-fix.md` | ✅ / ❌ W0 | ⬜ pending |
| 8-03-01 | 03 | 2 | FAGENT-03 | manual | `grep -c "Task tool" .claude/agents/*.md` | ✅ / ❌ | ⬜ pending |
| 8-04-01 | 04 | 2 | FAUTO-01 | manual | `grep "auto_mode" .claude/commands/mysd-plan.md` | ✅ / ❌ | ⬜ pending |
| 8-05-01 | 05 | 3 | FAUTO-03 | manual | `grep "ff\|ffe" .claude/commands/mysd-ff.md` | ✅ / ❌ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- Existing infrastructure covers all phase requirements (plugin layer only — no new test files needed).

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| /mysd:discuss invokes 4-dim parallel research | FCMD-01 | Claude Code skills cannot be unit tested | Run `/mysd:discuss` with a test change, observe agent spawning |
| /mysd:fix detects worktree path correctly | FCMD-02 | Path detection is runtime behavior | Run `/mysd:fix` in a worktree context, verify correct path |
| No nested Task tool in agent definitions | FAGENT-05 | Static audit of markdown files | `grep -n "Task(" .claude/agents/*.md` must return 0 results |
| ff implies --auto | FAUTO-03 | Behavioral test | Run `/mysd:ff` and verify it skips research/interactive prompts |
| ffe adds research before ff pipeline | FAUTO-04 | Behavioral test | Run `/mysd:ffe` and verify research step runs before planning |
| auto_mode propagates through propose→spec→plan | FAUTO-01 | End-to-end behavioral | Run propose with --auto, verify no interactive prompts throughout |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 300s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
