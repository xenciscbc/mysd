---
phase: 5
slug: schema-foundation-plan-checker
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-25
---

# Phase 5 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — standard Go test tooling |
| **Quick run command** | `go test ./internal/spec/... ./internal/planchecker/... ./internal/config/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/spec/... ./internal/planchecker/... ./internal/config/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 05-01-01 | 01 | 1 | FSCHEMA-01 | unit | `go test ./internal/spec/...` | ❌ W0 | ⬜ pending |
| 05-01-02 | 01 | 1 | FSCHEMA-02 | unit | `go test ./internal/spec/...` | ❌ W0 | ⬜ pending |
| 05-01-03 | 01 | 1 | FSCHEMA-03 | unit | `go test ./internal/spec/...` | ❌ W0 | ⬜ pending |
| 05-02-01 | 02 | 1 | FSCHEMA-04 | unit | `go test ./internal/spec/...` | ❌ W0 | ⬜ pending |
| 05-02-02 | 02 | 1 | FSCHEMA-05 | unit | `go test ./internal/executor/...` | ❌ W0 | ⬜ pending |
| 05-03-01 | 03 | 2 | FSCHEMA-06 | unit | `go test ./internal/planchecker/...` | ❌ W0 | ⬜ pending |
| 05-03-02 | 03 | 2 | FSCHEMA-07 | unit | `go test ./internal/planchecker/...` | ❌ W0 | ⬜ pending |
| 05-04-01 | 04 | 1 | FMODEL-01 | unit | `go test ./internal/config/...` | ❌ W0 | ⬜ pending |
| 05-04-02 | 04 | 1 | FMODEL-02 | unit | `go test ./internal/config/...` | ❌ W0 | ⬜ pending |
| 05-04-03 | 04 | 1 | FMODEL-03 | unit | `go test ./internal/config/...` | ❌ W0 | ⬜ pending |
| 05-05-01 | 05 | 2 | FAGENT-04 | unit | `go test ./...` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/spec/schema_test.go` — tests for new TaskEntry fields (depends/files/satisfies/skills) serialization and backward compat
- [ ] `internal/planchecker/checker_test.go` — tests for MUST coverage checking
- [ ] `internal/config/config_test.go` — tests for new model role resolution
- [ ] `internal/config/openspec_test.go` — tests for config.yaml read/write

*Existing test infrastructure (go test) covers all phase requirements.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| `mysd plan --check` interactive prompt | FSCHEMA-07 | Requires terminal interaction | Run `mysd plan --check` on a change with incomplete satisfies coverage, verify gap list and prompt appear |
| `mysd plan --context-only` JSON output | FSCHEMA-06 | Integration with full CLI | Run command, verify JSON contains WaveGroups, WorktreeDir, AutoMode fields |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
