---
phase: 2
slug: execution-engine
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-23
---

# Phase 2 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — standard Go test conventions |
| **Quick run command** | `go test ./internal/...` |
| **Full suite command** | `go test -v -count=1 ./...` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/...`
- **After every plan wave:** Run `go test -v -count=1 ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| *Populated during planning* | | | | | | | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] Verify existing test infrastructure from Phase 1 still passes
- [ ] `internal/task/` — new package for task management types and operations
- [ ] `internal/alignment/` — new package for alignment gate logic

*Existing test infrastructure from Phase 1 covers spec, state, config, output packages.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Alignment gate blocks AI code generation | EXEC-01 | Requires Claude Code agent interaction | Run `mysd execute` and verify alignment.md is created before any code changes |
| Plugin SKILL.md loads in Claude Code | WCMD-01 | Requires Claude Code runtime | Install plugin and verify `/mysd:` commands appear |
| Wave mode parallel execution | EXEC-04 | Requires Claude Code Task tool | Run `mysd execute --mode=wave --agents=2` and verify parallel agent spawning |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
