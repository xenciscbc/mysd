---
phase: 6
slug: executor-wave-grouping-worktree-engine
status: draft
nyquist_compliant: true
wave_0_complete: false
created: 2026-03-25
---

# Phase 6 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — stdlib go test |
| **Quick run command** | `go test ./internal/worktree/... ./internal/executor/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~30 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/worktree/... ./internal/executor/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 6-01-01 | 01 | 1 | FEXEC-01 | unit | `go test ./internal/executor/... -run TestBuildWaveGroups` | W0 | pending |
| 6-01-02 | 01 | 1 | FEXEC-02 | unit | `go test ./internal/executor/... -run TestFileOverlapSplit` | W0 | pending |
| 6-02-01 | 02 | 1 | FEXEC-03 | unit | `go test ./internal/worktree/... -run TestCreate` | W0 | pending |
| 6-02-02 | 02 | 1 | FEXEC-04 | unit | `go test ./internal/worktree/... -run TestRemove` | W0 | pending |
| 6-02-03 | 02 | 1 | FEXEC-05 | unit | `go test ./internal/worktree/... -run TestPathConvention` | W0 | pending |
| 6-03-01 | 03 | 2 | FEXEC-03 | unit | `go test ./cmd/... -run TestExecute -v -count=1` | W0 | pending |
| 6-03-02 | 03 | 2 | FEXEC-03 | unit | `go build ./... && go test ./cmd/... -run TestPlan -v -count=1` | W0 | pending |
| 6-04-01 | 04 | 2 | FEXEC-06, FEXEC-07, FEXEC-09 | manual-only | N/A — SKILL.md Markdown, not Go code | N/A | pending |
| 6-04-02 | 04 | 2 | FEXEC-12 | manual-only | N/A — SKILL.md Markdown, not Go code | N/A | pending |

*Status: pending / green / red / flaky*

---

## Wave 0 Requirements

- [ ] `internal/executor/waves_test.go` — stubs for FEXEC-01, FEXEC-02
- [ ] `internal/worktree/worktree_test.go` — stubs for FEXEC-03, FEXEC-04, FEXEC-05

*Note: Plans 01 and 02 are TDD plans that create their own test files as part of execution. Wave 0 stubs are created inline by the TDD RED phase, not as separate pre-requisite tasks.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Merge loop ascending ID order | FEXEC-06 | Logic lives in SKILL.md Markdown instructions, not executable Go code. Verified by grep for keyword presence in SKILL.md. | Review `plugin/commands/mysd-execute.md` for merge loop section; confirm ascending task ID sort and `--no-ff` flag. |
| AI conflict resolution (3 retries) | FEXEC-07 | Requires LLM interaction, non-deterministic. Logic is SKILL.md instructions, not Go code. | Create conflicting changes in two tasks, run wave execution, observe AI resolution attempts (up to 3). |
| Continue-on-failure policy | FEXEC-09 | SKILL.md orchestration logic — one task failure does not abort others. Cannot be unit tested. | Simulate one task failure in a multi-task wave, verify other tasks still merge. |
| Worktree isolation in executor agent | FEXEC-12 | Agent behavior defined in Markdown (`mysd-executor.md`), not testable Go code. | Spawn executor with `isolation: "worktree"`, verify all operations stay within worktree path. |
| Progress display output | D-07 | Visual terminal output validation (lipgloss styling). | Run `mysd execute` on wave tasks, observe styled progress output. |

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or are explicitly marked Manual-Only
- [x] Sampling continuity: Go-testable plans (01, 02, 03) have full automated coverage; Plan 04 is Manual-Only (SKILL.md Markdown)
- [x] No watch-mode flags
- [x] Feedback latency < 30s for automated tests
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
