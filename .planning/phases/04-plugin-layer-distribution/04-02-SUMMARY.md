---
phase: 04-plugin-layer-distribution
plan: 02
subsystem: roadmap-tracking
tags: [roadmap, tracking, mermaid, yaml, best-effort]
dependency_graph:
  requires: [internal/state, gopkg.in/yaml.v3]
  provides: [internal/roadmap.UpdateTracking, internal/roadmap.GenerateMermaid, internal/roadmap.ReadTracking]
  affects: [cmd/propose, cmd/spec, cmd/design, cmd/plan, cmd/verify, cmd/archive, cmd/ff, cmd/ffe, cmd/task_update]
tech_stack:
  added: []
  patterns: [best-effort tracking, zero-value on missing file, text/template for Mermaid generation]
key_files:
  created:
    - internal/roadmap/roadmap.go
    - internal/roadmap/roadmap_test.go
    - internal/roadmap/mermaid.go
    - internal/roadmap/mermaid_test.go
  modified:
    - cmd/propose.go
    - cmd/spec.go
    - cmd/design.go
    - cmd/plan.go
    - cmd/verify.go
    - cmd/archive.go
    - cmd/ff.go
    - cmd/ffe.go
    - cmd/task_update.go
decisions:
  - "UpdateTracking derives project root from filepath.Dir(specsDir) — handles both .specs/ and openspec/ conventions"
  - "timeline.md wrapped in ```mermaid code fence for direct rendering in GitHub/GitLab"
  - "task_update SaveState also gets tracking — records last activity even without phase change"
metrics:
  duration: 5m
  completed: "2026-03-24"
  tasks: 2
  files: 13
---

# Phase 4 Plan 02: Roadmap Tracking Package Summary

**One-liner:** YAML-based change lifecycle tracker with auto-regenerated Mermaid gantt chart, integrated best-effort into all 9 state-transitioning command call sites.

## What Was Built

### internal/roadmap package

**`roadmap.go`** — Core tracking logic:
- `UpdateTracking(specsDir string, ws state.WorkflowState) error` — reads/upserts `.mysd/roadmap/tracking.yaml`, regenerates `timeline.md`
- `ReadTracking(roadmapDir string) (TrackingFile, error)` — returns zero-value on missing file (convention over config, mirrors `uat.ReadUAT` pattern)
- `TrackingFile` / `ChangeRecord` structs with YAML tags
- Project root derived from `filepath.Dir(specsDir)` — supports both `.specs/` and `openspec/` layouts

**`mermaid.go`** — Mermaid chart generation:
- `GenerateMermaid(tf TrackingFile) string` — produces Mermaid gantt chart using `text/template` stdlib
- `ganttStatus` mapping: archived/verified → "done", executed → "active", others → no modifier
- Empty changes list produces valid gantt header with no sections

### Command Integration

All 9 `state.SaveState` call sites now have a best-effort `roadmap.UpdateTracking` immediately after:
- `cmd/propose.go`, `cmd/spec.go`, `cmd/design.go`, `cmd/plan.go`
- `cmd/verify.go`, `cmd/archive.go`
- `cmd/ff.go`, `cmd/ffe.go` (per transition, inside the loop)
- `cmd/task_update.go`

Pattern matches Phase 3's ARCHIVED-STATE.json best-effort decision: warning to stderr, never blocks the command.

## Test Results

All 8 TDD test cases pass:
- `TestUpdateTracking_NewFile` — creates tracking.yaml with correct schema_version and change record
- `TestUpdateTracking_UpsertExisting` — updates without duplicating existing records
- `TestUpdateTracking_MultipleChanges` — tracks multiple distinct changes
- `TestUpdateTracking_CompletedAt` — sets CompletedAt on PhaseArchived
- `TestUpdateTracking_TimelineMdGenerated` — timeline.md exists and contains "gantt"
- `TestUpdateTracking_ProjectRootDerivation` — tracking.yaml at `{root}/.mysd/roadmap/`, NOT inside `.specs/`
- `TestGenerateMermaid_BasicChart` — output contains "gantt", "dateFormat", change names
- `TestGenerateMermaid_EmptyChanges` — valid gantt header with no sections

All existing cmd tests pass without modification.

## Deviations from Plan

None — plan executed exactly as written, with one minor addition: `cmd/task_update.go` was included (plan listed it under "Files to check") and it calls `state.SaveState`, so tracking was added there too.

## Known Stubs

None.

## Self-Check

### Created Files

- `internal/roadmap/roadmap.go` — exists
- `internal/roadmap/roadmap_test.go` — exists
- `internal/roadmap/mermaid.go` — exists
- `internal/roadmap/mermaid_test.go` — exists
- `.planning/phases/04-plugin-layer-distribution/04-02-SUMMARY.md` — this file

### Commits

- `5cf13d0` — feat(04-02): create internal/roadmap package with tracking and Mermaid generation
- `f07b090` — feat(04-02): integrate roadmap tracking into all state-transitioning commands

## Self-Check: PASSED
