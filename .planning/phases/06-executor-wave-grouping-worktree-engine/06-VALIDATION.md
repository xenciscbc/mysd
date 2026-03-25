---
phase: 6
slug: executor-wave-grouping-worktree-engine
status: draft
nyquist_compliant: false
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
| 6-01-01 | 01 | 1 | FEXEC-01 | unit | `go test ./internal/executor/... -run TestBuildWaveGroups` | ❌ W0 | ⬜ pending |
| 6-01-02 | 01 | 1 | FEXEC-02 | unit | `go test ./internal/executor/... -run TestFileOverlapSplit` | ❌ W0 | ⬜ pending |
| 6-02-01 | 02 | 1 | FEXEC-03 | unit | `go test ./internal/worktree/... -run TestCreate` | ❌ W0 | ⬜ pending |
| 6-02-02 | 02 | 1 | FEXEC-04 | unit | `go test ./internal/worktree/... -run TestRemove` | ❌ W0 | ⬜ pending |
| 6-02-03 | 02 | 1 | FEXEC-05 | unit | `go test ./internal/worktree/... -run TestPathConvention` | ❌ W0 | ⬜ pending |
| 6-03-01 | 03 | 2 | FEXEC-06 | unit | `go test ./internal/worktree/... -run TestDiskSpaceCheck` | ❌ W0 | ⬜ pending |
| 6-03-02 | 03 | 2 | FEXEC-07 | unit | `go test ./internal/worktree/... -run TestWindowsLongPaths` | ❌ W0 | ⬜ pending |
| 6-04-01 | 04 | 2 | FEXEC-08 | unit | `go test ./internal/executor/... -run TestMergeLoop` | ❌ W0 | ⬜ pending |
| 6-04-02 | 04 | 2 | FEXEC-09 | unit | `go test ./internal/executor/... -run TestConflictResolution` | ❌ W0 | ⬜ pending |
| 6-04-03 | 04 | 2 | FEXEC-10 | unit | `go test ./internal/executor/... -run TestContinueOnFailure` | ❌ W0 | ⬜ pending |
| 6-05-01 | 05 | 3 | FEXEC-11 | integration | `go test ./... -run TestExecuteWaveIntegration` | ❌ W0 | ⬜ pending |
| 6-05-02 | 05 | 3 | FEXEC-12 | integration | `go test ./... -run TestWorktreeCleanup` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/executor/waves_test.go` — stubs for FEXEC-01, FEXEC-02
- [ ] `internal/worktree/worktree_test.go` — stubs for FEXEC-03, FEXEC-04, FEXEC-05, FEXEC-06, FEXEC-07
- [ ] `internal/executor/merge_test.go` — stubs for FEXEC-08, FEXEC-09, FEXEC-10
- [ ] `internal/executor/integration_test.go` — stubs for FEXEC-11, FEXEC-12

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| AI conflict resolution (3 retries) | FEXEC-09 | Requires LLM interaction, non-deterministic | Create conflicting changes in two tasks, run wave execution, observe AI resolution attempts |
| Progress display output | FEXEC-12 | Visual terminal output validation | Run `mysd execute` on wave tasks, observe lipgloss-styled progress |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 30s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
