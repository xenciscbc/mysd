---
phase: 7
slug: new-binary-commands-scanner-refactor
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-26
---

# Phase 7 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — uses Go built-in testing |
| **Quick run command** | `go test ./internal/executor/... ./internal/scanner/... ./internal/config/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/executor/... ./internal/scanner/... ./internal/config/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 7-01-01 | 01 | 1 | FCMD-05 (D-13) | unit | `go test ./internal/executor/... -run TestFilterBlockedTasks` | ❌ W0 | ⬜ pending |
| 7-01-02 | 01 | 1 | FCMD-05 (D-13) | unit | `go test ./internal/executor/...` | ✅ | ⬜ pending |
| 7-02-01 | 02 | 1 | FSCAN-01~03 (D-01~04) | unit | `go test ./internal/scanner/... -run TestBuildScanContext` | ✅ | ⬜ pending |
| 7-02-02 | 02 | 1 | FSCAN-01~03 | unit | `go test ./internal/scanner/...` | ✅ | ⬜ pending |
| 7-02-03 | 02 | 1 | FSCAN-05 (D-05) | unit | `go test ./cmd/... -run TestInitUsesScaffoldOnly` | ❌ W0 | ⬜ pending |
| 7-03-01 | 03 | 2 | FCMD-03 (D-11,D-12) | unit | `go test ./cmd/... -run TestModelCommand` | ❌ W0 | ⬜ pending |
| 7-03-02 | 03 | 2 | FCMD-03 | integration | `go build ./... && ./mysd model` | ✅ | ⬜ pending |
| 7-04-01 | 04 | 2 | FCMD-04 (D-06, Phase5 D-09) | unit | `go test ./internal/config/... -run TestAtomicLangSync` | ❌ W0 | ⬜ pending |
| 7-04-02 | 04 | 2 | FCMD-04 | integration | `go build ./... && ./mysd lang set zh-TW` | ✅ | ⬜ pending |
| 7-05-01 | 05 | 3 | SKILL-01~04 (D-07~10) | manual | Verify SKILL.md / agent definition renders task↔skills table | ✅ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/executor/waves_test.go` — stubs for `TestFilterBlockedTasks` (transitive closure test cases)
- [ ] `cmd/model_test.go` — stubs for `TestModelCommand` (read / set subcommand)
- [ ] `internal/config/lang_test.go` — stubs for `TestAtomicLangSync` (rollback on partial write failure)
- [ ] `cmd/init_test.go` — stubs for `TestInitUsesScaffoldOnly` (init delegates to scan --scaffold-only)

*Wave 0 must run `go build ./...` green before any other tasks.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| task↔skills 表格顯示與批次確認互動 | SKILL-02, SKILL-03 | 純 SKILL.md/agent 層 UX，Go binary 不參與 | 執行 `/mysd:plan`，plan 完成後確認 Claude 呈現 task↔skills table，輸入 Y 測試批次同意，輸入 n 測試逐一調整 |
| ffe 模式跳過確認 | SKILL-04 | agent 行為，需人工觀察 SKILL.md 執行流程 | 執行 `/mysd:ffe`，確認 plan 完成後不出現確認對話，直接使用推薦 skills |
| `/mysd:scan` 在 Node.js 專案偵測語言 | FSCAN-01 | 需要真實非 Go 專案目錄環境 | 在含 `package.json` 的目錄執行 `mysd scan --context-only`，確認輸出 JSON `primary_language: "nodejs"` |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
