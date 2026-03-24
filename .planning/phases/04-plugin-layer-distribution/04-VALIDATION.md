---
phase: 4
slug: plugin-layer-distribution
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-24
---

# Phase 4 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — existing go test infrastructure |
| **Quick run command** | `go test ./internal/scanner/... ./internal/roadmap/... ./cmd/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/scanner/... ./internal/roadmap/... ./cmd/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 04-01-01 | 01 | 1 | WCMD-09 | unit | `go test ./internal/scanner/...` | ❌ W0 | ⬜ pending |
| 04-01-02 | 01 | 1 | WCMD-09 | unit | `go test ./cmd/... -run TestScan` | ❌ W0 | ⬜ pending |
| 04-02-01 | 02 | 1 | DIST-04 | unit | `go test ./internal/roadmap/...` | ❌ W0 | ⬜ pending |
| 04-02-02 | 02 | 1 | DIST-04 | file | `test -f plugin/plugin.json` | ❌ W0 | ⬜ pending |
| 04-03-01 | 03 | 2 | DIST-03 | file | `test -f .goreleaser.yaml` | ❌ W0 | ⬜ pending |
| 04-03-02 | 03 | 2 | DIST-03 | build | `go build ./...` | ✅ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/scanner/scanner_test.go` — stubs for WCMD-09 scan context builder
- [ ] `internal/roadmap/tracking_test.go` — stubs for DIST-04 roadmap tracking
- [ ] `cmd/scan_test.go` — stubs for scan command flags

*Existing go test infrastructure covers framework needs.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Plugin install in Claude Code | DIST-04 | Requires Claude Code runtime | Copy plugin/ to .claude/plugins/mysd/, start Claude Code, verify /mysd:* commands appear |
| go install cross-platform | DIST-03 | Requires CI/release pipeline | Run `go install` on macOS/Linux/Windows, verify binary works |
| SessionStart hook advisory | DIST-04 | Requires Claude Code session | Start session without mysd in PATH, verify warning displayed but session continues |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
