---
phase: 12
slug: context
status: draft
nyquist_compliant: true
wave_0_complete: true
created: 2026-03-27
---

# Phase 12 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go testing (stdlib) + testify v1 |
| **Config file** | No standalone config, uses `go test ./...` |
| **Quick run command** | `go test ./cmd/... -run TestStatusline -v` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./cmd/... -run TestStatusline -v`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Wave 0 Strategy

Plans 12-01 and 12-03 use `tdd="true"` on their tasks. This means the executor follows a RED-GREEN-REFACTOR cycle: tests are written FIRST (RED), then implementation makes them pass (GREEN). This is functionally equivalent to a Wave 0 stub step — the test file is created as part of the TDD task's RED phase before any production code is written.

**Decision:** `tdd="true"` tasks satisfy the Nyquist Wave 0 requirement. No separate stub-creation tasks needed.

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | TDD | Status |
|---------|------|------|-------------|-----------|-------------------|-----|--------|
| 12-01-01 | 01 | 1 | ProjectConfig extension + statusline subcommand | unit | `go test ./cmd/... -run TestRunStatusline` | tdd=true | pending |
| 12-01-02 | 01 | 1 | mysd init hook install + settings merge | unit | `go test ./cmd/... -run "TestInitStatusline\|TestWriteSettings"` | tdd=true | pending |
| 12-02-01 | 02 | 1 | statusline hook JS | manual smoke | `node -c plugin/hooks/mysd-statusline.js && echo '...' \| node plugin/hooks/mysd-statusline.js` | N/A | pending |
| 12-02-02 | 02 | 1 | /mysd:statusline SKILL.md | grep verify | `grep -q "argument-hint" && grep -q "mysd statusline" plugin/commands/mysd-statusline.md` | N/A | pending |
| 12-03-01 | 03 | 2 | archive deletes cache | unit | `go test ./cmd/... -run TestArchiveDeletesResearchCache` | tdd=true | pending |
| 12-03-02 | 03 | 2 | discuss cache detection + write | grep verify | `grep -c "discuss-research-cache.json" plugin/commands/mysd-discuss.md` | N/A | pending |

*Status: pending / green / red / flaky*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| mysd-statusline.js outputs correct statusline format in Claude Code | D-01~D-10 | JS hook runs inside Claude Code process, cannot unit test end-to-end | Install, then observe statusline output visually |
| mysd-statusline.js disabled mode suppresses output | D-12 | Same as above | `mysd statusline off`, then observe statusline disappears |
| discuss cache can be reused across sessions | D-14~D-16 | SKILL.md orchestration, AI execution logic | Run /mysd:discuss, interrupt, restart, confirm cache detection prompt appears |

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or TDD mode (tdd=true is Wave 0 equivalent)
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covered by tdd=true task attribute (RED phase creates tests first)
- [x] No watch-mode flags
- [x] Feedback latency < 30s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** approved
