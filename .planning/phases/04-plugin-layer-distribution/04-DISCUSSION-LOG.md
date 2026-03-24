# Phase 4: Plugin Layer & Distribution - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-24
**Phase:** 04-plugin-layer-distribution
**Areas discussed:** Scan command strategy, Distribution & installation, Roadmap tracking format, Plugin packaging

---

## Scan Command Strategy

### Execution Architecture

| Option | Description | Selected |
|--------|-------------|----------|
| Context-only pattern | Same as execute/verify: binary outputs codebase metadata JSON, AI agent generates specs | ✓ |
| Binary-only approach | Go binary does all analysis independently, no AI agent | |
| Hybrid approach | Binary does structural analysis, AI adds human-readable descriptions | |

**User's choice:** Context-only pattern
**Notes:** Follows established architecture from Phase 2-3

### Scan Scope

| Option | Description | Selected |
|--------|-------------|----------|
| Whole project default | Default scans whole project, --path for subdirectory | |
| Explicit directory required | Must specify directory | |
| Interactive selection | AI lists modules, user selects | |

**User's choice:** Custom — interactive exclusion step. Binary lists codebase structure, user confirms what to exclude, then proceeds.
**Notes:** User wants to discuss and confirm exclusions before scanning, not just blindly scan everything.

### Existing Spec Handling

| Option | Description | Selected |
|--------|-------------|----------|
| Skip existing | Skip if .specs/changes/{name}/ exists | ✓ |
| Overwrite with flag | Default skip, --force to overwrite | |
| Merge/update | Read existing and update | |

**User's choice:** Skip existing

### Granularity

| Option | Description | Selected |
|--------|-------------|----------|
| Per-module/package | Each Go package = one change | ✓ |
| Per-file | Each file = one spec | |
| Per-feature (AI decides) | AI groups files into features | |

**User's choice:** Per-module/package

---

## Distribution & Installation

### Installation Methods

| Option | Description | Selected |
|--------|-------------|----------|
| go install | go install github.com/owner/mysd@latest | ✓ |
| GitHub Releases binaries | GoReleaser precompiled binaries | (not selected as primary) |
| Homebrew | brew install via tap/cask | |
| SessionStart auto-download | Hook auto-downloads binary | |

**User's choice:** go install as primary method
**Notes:** GoReleaser still produces binaries for GitHub Releases, but go install is the recommended path

### GoReleaser Configuration

| Option | Description | Selected |
|--------|-------------|----------|
| Standard 3-platform | Linux + macOS + Windows (amd64, arm64) | ✓ |
| Minimal (macOS + Linux) | Skip Windows binaries | |
| Full matrix | Include 386, arm, etc. | |

**User's choice:** Standard 3-platform

---

## Roadmap Tracking Format

### File Format

| Option | Description | Selected |
|--------|-------------|----------|
| Single YAML file | .mysd/roadmap/tracking.yaml | ✓ |
| Per-change markdown | One .md per change | |
| JSON for tooling | Pure JSON | |

**User's choice:** Single YAML file

### Update Trigger

| Option | Description | Selected |
|--------|-------------|----------|
| On state transitions | Every propose→spec→...→archive transition | ✓ |
| On command completion | Every /mysd:* command | |
| Manual via /mysd:status | Only when status queried | |

**User's choice:** On state transitions

### Tracked Fields

**User's choice:** Multiple fields selected:
- Name + status + dates (core)
- Task count + completion (progress)
- Verification status (quality)

### Mermaid Chart

| Option | Description | Selected |
|--------|-------------|----------|
| Auto-generate in YAML | Mermaid block appended to YAML | |
| Separate timeline.md | .mysd/roadmap/timeline.md | ✓ |
| CLI output only | Dynamic generation on status | |

**User's choice:** Separate timeline.md — user suggested separating data (YAML) from visualization (Mermaid) for cleanliness
**Notes:** timeline.md auto-regenerated whenever tracking.yaml updates

---

## Plugin Packaging

### Installation Experience

| Option | Description | Selected |
|--------|-------------|----------|
| Plugin directory + SessionStart hook | Copy plugin/, hook checks binary | ✓ |
| Git clone + auto-setup | Clone repo, hook builds binary | |
| NPM-style install | npx install-mysd | |

**User's choice:** Plugin directory + SessionStart hook

### SessionStart Hook Behavior

| Option | Description | Selected |
|--------|-------------|----------|
| Check binary exists + version | Verify mysd in PATH, version check | ✓ |
| Auto-download binary | Detect OS, download from releases | |
| No hook | User manages binary manually | |

**User's choice:** Check binary exists + version

### plugin.json Format

| Option | Description | Selected |
|--------|-------------|----------|
| Standard Claude Code format | name, version, description, commands[], agents[], hooks[] | ✓ |
| Extended with metadata | Standard + author, license, homepage | |
| You decide | Claude's discretion | |

**User's choice:** Standard Claude Code format

### Upgrade Strategy

| Option | Description | Selected |
|--------|-------------|----------|
| Version check in SessionStart | Compare versions, show upgrade instructions | ✓ |
| Auto-update binary | Auto-download new version | |
| Manual only | No version checking | |

**User's choice:** Version check in SessionStart — advisory only, never blocks

---

## Claude's Discretion

- GoReleaser configuration details
- Scan agent prompt wording
- tracking.yaml schema design
- timeline.md Mermaid styling
- Plugin directory organization

## Deferred Ideas

None — discussion stayed within phase scope
