# plan-conversation-context Specification

## Purpose

TBD - created by archiving change 'state-files-and-plan-context'. Update Purpose after archive.

## Requirements

### Requirement: Conversation context option in plan spec selector

The plan SKILL.md spec selector (Step 2b) SHALL include a "From conversation context" option alongside the existing per-spec options and "All" option when no `--spec` flag is given.

#### Scenario: Spec selector shows conversation context option

- **WHEN** the user runs `/mysd:plan` without `--spec` flag and multiple specs exist
- **THEN** the interactive spec selector SHALL display individual specs, an "All" option, and a "From conversation context" option

#### Scenario: User selects conversation context option

- **WHEN** the user selects "From conversation context" in the spec selector
- **THEN** the orchestrator SHALL extract relevant requirements and task descriptions from the current conversation context
- **AND** write the extracted content to a temporary file `conversation-context.md` in the change directory
- **AND** pass the file path via the existing `--from` flag to `mysd plan --context-only`

#### Scenario: Conversation context temp file cleanup

- **WHEN** the plan pipeline completes (regardless of success or failure)
- **THEN** the orchestrator SHALL delete the `conversation-context.md` temporary file from the change directory

#### Scenario: All option does not include conversation context

- **WHEN** the user selects "All" in the spec selector
- **THEN** the planner SHALL plan all specs without extracting conversation context


<!-- @trace
source: state-files-and-plan-context
updated: 2026-04-02
code:
  - .omc/state/last-tool-error.json
  - .omc/state/mission-state.json
  - .omc/sessions/c3e46307-af7b-4b85-b578-57c643e1d6f6.json
  - .omc/sessions/c9f8e296-c5c6-4b63-9077-1bef396f2cf5.json
  - .omc/project-memory.json
  - .omc/sessions/d620ac15-32ae-4049-892b-0c0e9feeb48c.json
  - .omc/state/agent-replay-39a9d699-2d3f-42fc-b2da-fa419882ce6b.jsonl
  - .omc/sessions/61c9a8ed-7bb4-4d6a-9108-6c2835c4560a.json
  - .omc/sessions/5dab8e8e-7dec-40f8-b680-f2f8f554ad4b.json
  - .omc/state/idle-notif-cooldown.json
  - .omc/state/subagent-tracking.json
-->

---
### Requirement: State files stored in .mysd directory

The mysd binary SHALL store workflow state files (`STATE.json`, `roadmap-tracking.json`, `roadmap-timeline.md`) in the `.mysd/` directory instead of the `openspec/` directory.

#### Scenario: STATE.json location

- **WHEN** any mysd command writes or reads workflow state
- **THEN** the state file SHALL be located at `.mysd/STATE.json` relative to the project root

#### Scenario: Roadmap files location

- **WHEN** the roadmap tracking system writes tracking or timeline files
- **THEN** the files SHALL be located at `.mysd/roadmap-tracking.json` and `.mysd/roadmap-timeline.md`

#### Scenario: Backward compatibility with existing STATE.json

- **WHEN** mysd reads state and `.mysd/STATE.json` does not exist but `openspec/STATE.json` does
- **THEN** mysd SHALL read from the legacy location and continue normally


<!-- @trace
source: state-files-and-plan-context
updated: 2026-04-02
code:
  - .omc/state/last-tool-error.json
  - .omc/state/mission-state.json
  - .omc/sessions/c3e46307-af7b-4b85-b578-57c643e1d6f6.json
  - .omc/sessions/c9f8e296-c5c6-4b63-9077-1bef396f2cf5.json
  - .omc/project-memory.json
  - .omc/sessions/d620ac15-32ae-4049-892b-0c0e9feeb48c.json
  - .omc/state/agent-replay-39a9d699-2d3f-42fc-b2da-fa419882ce6b.jsonl
  - .omc/sessions/61c9a8ed-7bb4-4d6a-9108-6c2835c4560a.json
  - .omc/sessions/5dab8e8e-7dec-40f8-b680-f2f8f554ad4b.json
  - .omc/state/idle-notif-cooldown.json
  - .omc/state/subagent-tracking.json
-->

---
### Requirement: STATE.json cleanup after archive

The `mysd archive` command SHALL delete `.mysd/STATE.json` after a successful archive operation.

#### Scenario: STATE.json deleted on successful archive

- **WHEN** the archive command completes successfully
- **THEN** `.mysd/STATE.json` SHALL be deleted

#### Scenario: Archive succeeds even if STATE.json deletion fails

- **WHEN** the archive command completes successfully but STATE.json deletion fails
- **THEN** the archive operation SHALL still be considered successful
- **AND** a warning SHALL be printed to stderr

<!-- @trace
source: state-files-and-plan-context
updated: 2026-04-02
code:
  - .omc/state/last-tool-error.json
  - .omc/state/mission-state.json
  - .omc/sessions/c3e46307-af7b-4b85-b578-57c643e1d6f6.json
  - .omc/sessions/c9f8e296-c5c6-4b63-9077-1bef396f2cf5.json
  - .omc/project-memory.json
  - .omc/sessions/d620ac15-32ae-4049-892b-0c0e9feeb48c.json
  - .omc/state/agent-replay-39a9d699-2d3f-42fc-b2da-fa419882ce6b.jsonl
  - .omc/sessions/61c9a8ed-7bb4-4d6a-9108-6c2835c4560a.json
  - .omc/sessions/5dab8e8e-7dec-40f8-b680-f2f8f554ad4b.json
  - .omc/state/idle-notif-cooldown.json
  - .omc/state/subagent-tracking.json
-->