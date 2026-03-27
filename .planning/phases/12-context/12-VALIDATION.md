---
phase: 12
slug: context
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-27
---

# Phase 12 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go testing (stdlib) + testify v1 |
| **Config file** | 無獨立 config，使用 `go test ./...` |
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

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 12-01-01 | 01 | 1 | ProjectConfig extension | unit | `go test ./internal/config/...` | ❌ W0 | ⬜ pending |
| 12-01-02 | 01 | 1 | mysd statusline toggle | unit | `go test ./cmd/... -run TestRunStatusline` | ❌ W0 | ⬜ pending |
| 12-01-03 | 01 | 1 | mysd init hook install | unit | `go test ./cmd/... -run TestInitStatuslineInstall` | ❌ W0 | ⬜ pending |
| 12-01-04 | 01 | 1 | settings.json merge | unit | `go test ./cmd/... -run TestWriteSettingsStatusLine` | ❌ W0 | ⬜ pending |
| 12-02-01 | 02 | 1 | statusline hook JS | manual smoke | N/A | N/A | ⬜ pending |
| 12-03-01 | 03 | 2 | archive deletes cache | unit | `go test ./cmd/... -run TestArchiveDeletesResearchCache` | ❌ W0 | ⬜ pending |
| 12-03-02 | 03 | 2 | discuss cache write | manual smoke | N/A | N/A | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `cmd/statusline_test.go` — stubs for TestRunStatusline, TestInitStatuslineInstall, TestWriteSettingsStatusLine
- [ ] `internal/config/config_test.go` — stub for ProjectConfig statusline_enabled field test
- [ ] `cmd/archive_test.go` — stub for TestArchiveDeletesResearchCache

*Existing `go test` infrastructure covers all phase requirements — no new framework installation needed.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| mysd-statusline.js 正確輸出 statusline 格式 | D-01~D-10 | JS hook 在 Claude Code process 內執行，無法單元測試 | 安裝後手動觀察 statusline 輸出 |
| mysd-statusline.js disabled 時不輸出 | D-12 | 同上 | `mysd statusline off` 後觀察 statusline 消失 |
| discuss cache 可被 reuse | D-14~D-16 | SKILL.md orchestration，AI 執行邏輯 | 跑 /mysd:discuss，中斷後重啟，確認 cache 偵測提示出現 |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 30s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
