---
phase: 3
slug: verification-feedback-loop
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-24
---

# Phase 3 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — standard Go test conventions |
| **Quick run command** | `go test ./internal/... ./cmd/...` |
| **Full suite command** | `go test -v -count=1 ./...` |
| **Estimated runtime** | ~15 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/... ./cmd/...`
- **After every plan wave:** Run `go test -v -count=1 ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 15 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| *Populated during planning* | | | | | | | pending |

*Status: pending -- green -- red -- flaky*

---

## Wave 0 Requirements

- [ ] `internal/verifier/` — new package for verification context and report types
- [ ] `internal/archiver/` — new package for archive operations
- [ ] Fix Requirement.ID empty string issue (parser.go needs SourceFile + CRC32 ID generation)
- [ ] Verify existing test infrastructure from Phase 1+2 still passes

*Wave 0 gaps are addressed by the first plan in Wave 1.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Verifier agent independence from executor | VRFY-02 | Requires Claude Code agent interaction | Run verify after execute, confirm verifier has no executor context |
| UAT interactive walkthrough quality | WCMD-12 | Requires human interaction with agent | Run /mysd:uat and confirm agent guides through checklist items |
| Archive 'Run UAT first?' prompt behavior | WCMD-07 | Requires interactive CLI session | Run mysd archive with UAT available, verify prompt appears |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 15s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
