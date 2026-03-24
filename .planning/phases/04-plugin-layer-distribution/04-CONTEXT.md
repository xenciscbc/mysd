# Phase 4: Plugin Layer & Distribution - Context

**Gathered:** 2026-03-24
**Status:** Ready for planning

<domain>
## Phase Boundary

Phase 4 delivers three capabilities:
1. **Distribution** — GoReleaser cross-platform binary builds, `go install` support, GitHub Releases
2. **Scan command** — `/mysd:scan` scans existing codebase and generates OpenSpec-format spec documents
3. **Roadmap tracking** — Automatic `.mysd/roadmap/` tracking files updated on state transitions
4. **Plugin packaging** — plugin.json manifest, SessionStart hook, formal plugin directory structure

This is the final phase — all core functionality (spec management, execution, verification) is complete. Phase 4 wraps everything into a distributable, installable package.

</domain>

<decisions>
## Implementation Decisions

### Scan Command (WCMD-09)
- **D-01:** Scan uses the context-only pattern (same as execute/verify) — binary produces codebase metadata JSON, AI agent analyzes and generates spec files
- **D-02:** Before scanning, binary lists codebase structure and presents to user for confirmation — user selects which directories/files to exclude before AI analysis begins
- **D-03:** Scan skips existing specs — if `.specs/changes/{name}/` already exists, that module is skipped (no overwrite). No `--force` flag in v1.
- **D-04:** Scan granularity is per-module/package — each Go package or major module produces one change/spec. Aligns with OpenSpec "one spec per capability" philosophy.

### Distribution (DIST-03)
- **D-05:** Primary installation method is `go install github.com/owner/mysd@latest`
- **D-06:** GoReleaser configured for standard 3-platform matrix: Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64). Produces GitHub Releases with precompiled binaries and checksums.
- **D-07:** No Homebrew tap/cask in v1 — `go install` and direct binary download are sufficient for initial release.

### Roadmap Tracking (RMAP-01~03)
- **D-08:** Tracking data stored in `.mysd/roadmap/tracking.yaml` — single YAML file recording all changes with name, status, dates, task counts, and verification status
- **D-09:** Mermaid gantt chart generated as separate `.mysd/roadmap/timeline.md` — auto-regenerated whenever tracking.yaml is updated. Keeps data (YAML) and visualization (Mermaid) cleanly separated.
- **D-10:** Tracking updates triggered on state transitions — integrated with existing SaveState flow. Every propose→spec→design→plan→execute→verify→archive transition updates the tracking file.
- **D-11:** Tracked fields: change name, current status, start/completion dates, task total/completed count, verification MUST pass/fail statistics

### Plugin Packaging (DIST-04)
- **D-12:** Plugin installed by copying plugin directory to `.claude/plugins/mysd/`. Standard Claude Code plugin directory structure.
- **D-13:** SessionStart hook checks binary existence and version — if `mysd` not in PATH or version below minimum, displays installation instructions (not auto-download). Non-blocking — session continues even if binary missing.
- **D-14:** plugin.json follows standard Claude Code format: name, version, description, commands[], agents[], hooks[]. No extended metadata in v1.
- **D-15:** Plugin upgrade via SessionStart version check — compares `mysd --version` against plugin.json min_version. Displays upgrade instructions when outdated. No auto-update.

### Claude's Discretion
- GoReleaser `.goreleaser.yaml` specific configuration details (archive format, naming convention, checksum algorithm)
- Scan agent's exact prompt wording and analysis depth
- tracking.yaml schema field ordering and naming conventions
- timeline.md Mermaid chart styling and section grouping
- plugin/ directory internal file organization

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Plugin Structure
- `CLAUDE.md` §Claude Code Plugin Integration — plugin.json format, SKILL.md format, agent format, hooks
- `CLAUDE.md` §Technology Stack — GoReleaser v2.14+, golangci-lint v2, recommended libraries

### Existing Plugin Files (pattern reference)
- `.claude/commands/mysd-execute.md` — SKILL.md orchestrator pattern to replicate for scan
- `.claude/agents/mysd-executor.md` — Agent definition pattern
- `.claude/commands/mysd-verify.md` — Context-only → agent → write-results pattern

### Existing CLI Commands (pattern reference)
- `cmd/verify.go` — --context-only and --write-results flag pattern
- `cmd/archive.go` — State gate checking pattern
- `internal/state/state.go` — SaveState integration point for roadmap tracking

### Requirements
- `.planning/REQUIREMENTS.md` — WCMD-09, DIST-03, DIST-04, RMAP-01~03

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- 13 SKILL.md files — pattern well-established for scan command SKILL.md
- 7 agent definitions — pattern for scan agent (mysd-scanner.md)
- `cmd/verify.go` — --context-only + --write-results flag pattern to replicate for scan
- `internal/state/state.go` SaveState function — integration point for roadmap tracking hook
- `go.mod` — all dependencies already present, no new libraries needed

### Established Patterns
- Context-only pattern: binary outputs JSON → SKILL.md invokes agent via Task → agent produces output → binary processes result
- Thin command layer: cobra command with flags, delegates to internal packages
- Testability: io.Writer injection, temp directory fixtures, testify assert/require
- State transitions: LoadState → gate check → action → Transition → SaveState

### Integration Points
- `internal/state/state.go` SaveState — hook roadmap tracking update here
- `cmd/root.go` — register new `scan` subcommand
- `.claude/commands/` — add mysd-scan.md SKILL.md
- `.claude/agents/` — add mysd-scanner.md agent
- Project root — add `.goreleaser.yaml`, `plugin/` directory, `plugin.json`

</code_context>

<specifics>
## Specific Ideas

- Scan's interactive exclusion step: binary lists directory tree, user marks directories to skip (e.g., vendor/, node_modules/, test fixtures), then proceeds with AI analysis
- tracking.yaml + timeline.md separation ensures machine-readable data stays clean while Mermaid visualization is human-friendly
- SessionStart hook is purely advisory — never blocks session start, just displays warnings

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 04-plugin-layer-distribution*
*Context gathered: 2026-03-24*
