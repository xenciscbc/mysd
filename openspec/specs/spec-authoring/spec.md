---
spec-version: "1.0"
capability: Spec Authoring & Design
delta: ADDED
status: done
---

## Requirement: Proposal Creation

The `mysd propose` command MUST scaffold a change directory at `.specs/changes/{name}/` containing:
- `proposal.md` with YAML frontmatter (spec-version, change, status, created, updated)
- `.openspec.yaml` metadata file

The `mysd capture` command MUST extract conversation context into a proposal scaffold.

Proposal frontmatter MUST include fields: `spec-version`, `change`, `status`, `created`, `updated`.

## Requirement: Spec Writing

The `mysd spec` command MUST write detailed requirements using RFC 2119 keywords (MUST/SHOULD/MAY).

Spec files MUST be located at `.specs/changes/{name}/specs/{capability}/spec.md`.

Spec frontmatter MUST include: `spec-version`, `capability`, `delta` (ADDED/MODIFIED/REMOVED), `status`.

The `--auto` flag MUST allow non-interactive spec generation.

## Requirement: Design Capture

The `mysd design` command MUST generate `design.md` capturing technical decisions and architecture.

Design frontmatter MUST include: `spec-version`, `change`, `status`.

## Requirement: Spec Parsing

The `spec` package MUST parse YAML frontmatter from all artifact types (proposal, spec, design, tasks).

The parser MUST gracefully handle brownfield files without frontmatter.

`ParseRequirements()` MUST extract RFC 2119 keywords (case-sensitive, uppercase only) and group them under the nearest `## Requirement:` heading.

Requirements MUST be categorized by keyword level: MUST, SHOULD, MAY.

## Requirement: Workflow State

The `state` package MUST track workflow phases via `STATE.json`:
- Phases: Proposed → Specced → Designed → Planned → Executed → Verified → Archived

`ValidTransitions()` MUST enforce the legal phase transition sequence.

`LoadState()` MUST return zero-value on missing file (convention over configuration).

## Requirement: Change Tracking

The `roadmap` package MUST maintain `tracking.yaml` with change lifecycle records.

Each `ChangeRecord` MUST include: Name, Status, StartedAt, CompletedAt.

`GenerateMermaidTimeline()` MUST render a Mermaid diagram in `timeline.md`.

## Requirement: Deferred Notes

The `mysd note` command MUST support adding and listing deferred notes.

Notes MUST be stored in `.specs/deferred.json`.

## Requirement: Status Dashboard

The `mysd status` command MUST display current change name, phase, task progress, and next step recommendation.

### Scenario: New Change Proposal

WHEN a user runs `mysd propose my-feature`
THEN `.specs/changes/my-feature/` is created with proposal.md skeleton
AND STATE.json is updated to phase "proposed"

### Scenario: Spec with Requirements

WHEN spec.md contains "The system MUST validate input"
THEN ParseRequirements() extracts it as a MUST-level requirement

## Covered Packages

- `cmd/propose.go`, `cmd/capture.go`, `cmd/design.go`, `cmd/spec.go`, `cmd/note.go`, `cmd/status.go`
- `internal/spec/` — OpenSpec parsing, validation, YAML roundtrip
- `internal/state/` — workflow phase state machine
- `internal/roadmap/` — change lifecycle tracking and Mermaid timeline
