---
phase: 2
slug: execution-engine
status: draft
nyquist_compliant: true
wave_0_complete: true
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
| 02-01-T1 | 02-01 | 1 | EXEC-05, D-13 | unit | `go test ./internal/spec/ -run "TestParseTasks\|TestUpdateTaskStatus\|TestWriteTasks" -count=1` | internal/spec/updater_test.go | pending |
| 02-01-T2 | 02-01 | 1 | EXEC-01, EXEC-02, TEST-03 | unit | `go test ./internal/executor/... -count=1` | internal/executor/context_test.go, internal/executor/progress_test.go, internal/executor/alignment_test.go | pending |
| 02-02-T1 | 02-02 | 1 | WCMD-11 | unit | `go test ./internal/config/... -count=1` | internal/config/config_test.go | pending |
| 02-02-T2 | 02-02 | 1 | WCMD-08 | unit | `go test ./internal/executor/ -run "TestRenderStatus\|TestBuildStatusSummary" -count=1` | internal/executor/status_test.go | pending |
| 02-03-T1 | 02-03 | 2 | EXEC-01, WCMD-05, WCMD-08 | unit | `go test ./cmd/ -run "TestTaskUpdate" -count=1` | cmd/task_update_test.go | pending |
| 02-03-T2 | 02-03 | 2 | WCMD-10, WCMD-14 | build | `go build ./...` | cmd/ff.go, cmd/ffe.go, cmd/capture.go | pending |
| 02-04-T1 | 02-04 | 2 | WCMD-02, WCMD-03, WCMD-04, TEST-02 | build | `go build ./... && go vet ./cmd/...` | cmd/spec.go, cmd/design.go, cmd/plan.go | pending |
| 02-05-T1 | 02-05 | 3 | WCMD-01~05, WCMD-08~14 | content | `grep -l "mysd execute --context-only" .claude/commands/mysd-execute.md && test $(ls .claude/commands/mysd-*.md \| wc -l) -eq 10` | .claude/commands/mysd-*.md | pending |
| 02-05-T2 | 02-05 | 3 | EXEC-01, TEST-02 | content | `grep -l "MANDATORY: Alignment Gate" .claude/agents/mysd-executor.md && grep -l "test_generation" .claude/agents/mysd-executor.md && test $(ls .claude/agents/mysd-*.md \| wc -l) -eq 5` | .claude/agents/mysd-*.md | pending |
| 02-06-T1 | 02-06 | 3 | EXEC-01, EXEC-05, TEST-01 | integration | `go test ./cmd/ -run "TestExecute\|TestStatus" -count=1` | cmd/execute_test.go, cmd/status_test.go | pending |
| 02-06-T2 | 02-06 | 3 | WCMD-10, WCMD-14, EXEC-04 | integration | `go test ./cmd/ -run "TestFF" -count=1` | cmd/ff_test.go | pending |

*Status: pending -- green -- red -- flaky*

---

## Wave 0 Requirements

- [x] Verify existing test infrastructure from Phase 1 still passes
- [x] `internal/executor/` package — new package created in Plan 02-01 (context.go, alignment.go, progress.go)
- [x] `internal/spec/updater.go` — TaskStatus round-trip created in Plan 02-01
- [x] `internal/spec/updater_test.go` — covers D-13, created in Plan 02-01
- [x] `internal/executor/context_test.go` — covers EXEC-01, EXEC-02, created in Plan 02-01
- [x] `internal/executor/progress_test.go` — covers EXEC-05, WCMD-08, created in Plan 02-01
- [x] `cmd/task_update.go` — new subcommand created in Plan 02-03
- [x] `cmd/task_update_test.go` — covers EXEC-04 (flag), created in Plan 02-03
- [x] `internal/config/config_test.go` extension — covers TEST-03, created in Plan 02-02

*All Wave 0 gaps are addressed by Plan 02-01 through 02-03. No orphan test references remain.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Alignment gate blocks AI code generation | EXEC-01 | Requires Claude Code agent interaction | Run `mysd execute` and verify alignment.md is created before any code changes |
| Plugin SKILL.md loads in Claude Code | WCMD-01 | Requires Claude Code runtime | Install plugin and verify `/mysd:` commands appear |
| Wave mode parallel execution | EXEC-04 | Requires Claude Code Task tool | Run `mysd execute --mode=wave --agents=2` and verify parallel agent spawning |

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or Wave 0 dependencies
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [x] No watch-mode flags
- [x] Feedback latency < 10s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** ready
