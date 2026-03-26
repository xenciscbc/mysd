---
phase: 9
slug: interactive-discovery-integration
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-26
---

# Phase 9 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — uses Go stdlib testing |
| **Quick run command** | `go test ./internal/... ./cmd/...` |
| **Full suite command** | `go test -v ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/... ./cmd/...`
- **After every plan wave:** Run `go test -v ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 09-01-01 | 01 | 1 | DISC-03, DISC-08 | unit | `go test ./internal/deferred/...` | ❌ W0 | ⬜ pending |
| 09-01-02 | 01 | 1 | DISC-03, DISC-08 | unit | `go test ./cmd/... -run TestNote` | ❌ W0 | ⬜ pending |
| 09-02-01 | 02 | 1 | DISC-04 | manual | SKILL.md prompt verification | N/A | ⬜ pending |
| 09-02-02 | 02 | 1 | DISC-01, DISC-06 | manual | SKILL.md gray areas flow | N/A | ⬜ pending |
| 09-02-03 | 02 | 1 | DISC-07 | manual | SKILL.md dual-loop cycle | N/A | ⬜ pending |
| 09-03-01 | 03 | 2 | DISC-02 | manual | SKILL.md spec research | N/A | ⬜ pending |
| 09-03-02 | 03 | 2 | DISC-03 | manual | SKILL.md plan single researcher | N/A | ⬜ pending |
| 09-03-03 | 03 | 2 | DISC-09 | manual | SKILL.md status deferred count | N/A | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/deferred/deferred_test.go` — stubs for DISC-08 deferred JSON CRUD
- [ ] `cmd/note_test.go` — stubs for note subcommand

*Note: SKILL.md orchestrator changes are Markdown-only and verified manually via `/mysd:verify`.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| 4-dimension research spawning | DISC-01, DISC-06 | SKILL.md orchestrator, no binary code | Run `/mysd:propose` with research, verify 4 researcher agents spawn |
| Dual-loop exploration | DISC-07 | Interactive UI flow | Run `/mysd:discuss` with research, verify area deep-dive + new area discovery |
| Scope guardrail redirect | DISC-08 | AI judgment in prompt | During discuss, mention out-of-scope feature, verify redirect to deferred notes |
| Auto mode skip | DISC-05 | Flag propagation | Run `/mysd:propose --auto`, verify exploration loop is skipped |
| Spec auto-update + re-plan | DISC-09 | Multi-agent orchestration | Run `/mysd:discuss`, incorporate conclusion, verify spec update + plan-checker trigger |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
