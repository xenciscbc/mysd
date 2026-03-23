---
phase: 1
slug: foundation
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-23
---

# Phase 1 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test (stdlib) + testify v1.11.1 |
| **Config file** | none — Wave 0 installs |
| **Quick run command** | `go test ./...` |
| **Full suite command** | `go test -v -count=1 ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./...`
- **After every plan wave:** Run `go test -v -count=1 ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 01-01-01 | 01 | 1 | SPEC-07 | unit | `go test ./internal/spec/...` | ❌ W0 | ⬜ pending |
| 01-01-02 | 01 | 1 | SPEC-01 | unit | `go test ./internal/spec/...` | ❌ W0 | ⬜ pending |
| 01-02-01 | 02 | 1 | SPEC-02 | unit | `go test ./internal/spec/...` | ❌ W0 | ⬜ pending |
| 01-02-02 | 02 | 1 | OPSX-01 | unit | `go test ./internal/spec/...` | ❌ W0 | ⬜ pending |
| 01-03-01 | 03 | 2 | STAT-01 | unit | `go test ./internal/state/...` | ❌ W0 | ⬜ pending |
| 01-03-02 | 03 | 2 | STAT-02 | unit | `go test ./internal/state/...` | ❌ W0 | ⬜ pending |
| 01-04-01 | 04 | 2 | CONF-01 | unit | `go test ./internal/config/...` | ❌ W0 | ⬜ pending |
| 01-04-02 | 04 | 2 | CONF-03 | unit | `go test ./internal/config/...` | ❌ W0 | ⬜ pending |
| 01-05-01 | 05 | 3 | DIST-01 | integration | `go build ./cmd/mysd` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/spec/parser_test.go` — stubs for SPEC-01, SPEC-02, SPEC-07, OPSX-01, OPSX-02
- [ ] `internal/spec/rfc2119_test.go` — stubs for RFC 2119 keyword parsing
- [ ] `internal/state/machine_test.go` — stubs for STAT-01, STAT-02, STAT-03
- [ ] `internal/config/config_test.go` — stubs for CONF-01, CONF-02, CONF-03, CONF-04
- [ ] `internal/spec/testdata/` — OpenSpec fixture directories
- [ ] `go.mod` + `go.sum` — module initialization with dependencies

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Colored terminal output | DIST-01 | Visual verification of lipgloss styling | Run `mysd status` in TTY and verify colored output |
| Cross-platform binary | DIST-02 | Requires multiple OS environments | Build for linux/darwin/windows and verify execution |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
