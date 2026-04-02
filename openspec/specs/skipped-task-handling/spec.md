# skipped-task-handling Specification

## Purpose

TBD - created by archiving change 'openspec-compliance'. Update Purpose after archive.

## Requirements

### Requirement: Skipped task marker

The system SHALL recognize `- [~]` as a skipped task marker in tasks.md. A skipped task MUST include a reason after the task description, separated by a colon or parenthetical notation.

#### Scenario: Parse skipped task with reason

- **GIVEN** a tasks.md line containing `- [~] 3.2 Implement caching（跳過：需求變更）`
- **WHEN** the task parser processes this line
- **THEN** the system SHALL mark the task as skipped with reason "需求變更"

#### Scenario: Skipped task without reason

- **GIVEN** a tasks.md line containing `- [~] 3.2 Implement caching` with no reason
- **WHEN** the task parser processes this line
- **THEN** the system SHALL emit a warning indicating a skipped task requires a reason


<!-- @trace
source: openspec-compliance
updated: 2026-04-02
code:
  - .omc/project-memory.json
  - .omc/state/last-tool-error.json
  - cmd/root.go
  - .omc/sessions/c3e46307-af7b-4b85-b578-57c643e1d6f6.json
  - .omc/sessions/c9f8e296-c5c6-4b63-9077-1bef396f2cf5.json
  - cmd/model.go
  - cmd/plan.go
  - .omc/sessions/61c9a8ed-7bb4-4d6a-9108-6c2835c4560a.json
  - internal/config/config.go
  - cmd/archive.go
  - internal/config/defaults.go
  - .omc/sessions/5dab8e8e-7dec-40f8-b680-f2f8f554ad4b.json
  - .omc/state/idle-notif-cooldown.json
  - cmd/spec.go
  - cmd/design.go
  - internal/verifier/scenario.go
  - .omc/state/mission-state.json
tests:
  - internal/config/config_test.go
  - internal/spec/delta.go
  - cmd/archive_test.go
  - cmd/integration_test.go
  - internal/spec/parser.go
  - internal/spec/delta_test.go
  - internal/spec/merge.go
  - cmd/model_test.go
  - internal/verifier/scenario_test.go
  - internal/spec/parser_test.go
  - internal/spec/merge_test.go
  - internal/spec/schema.go
-->

---
### Requirement: Archive gate accepts skipped tasks

The archive gate SHALL treat `[~]` (skipped) tasks the same as `[x]` (completed) tasks. The gate SHALL only block on `[ ]` (incomplete) tasks.

#### Scenario: All tasks completed or skipped

- **GIVEN** tasks.md with 8 `[x]` tasks and 2 `[~]` tasks
- **WHEN** the archive command checks task completion
- **THEN** the archive gate SHALL pass

#### Scenario: Incomplete tasks remain

- **GIVEN** tasks.md with 8 `[x]` tasks, 1 `[~]` task, and 1 `[ ]` task
- **WHEN** the archive command checks task completion
- **THEN** the archive gate SHALL block with an error indicating 1 incomplete task


<!-- @trace
source: openspec-compliance
updated: 2026-04-02
code:
  - .omc/project-memory.json
  - .omc/state/last-tool-error.json
  - cmd/root.go
  - .omc/sessions/c3e46307-af7b-4b85-b578-57c643e1d6f6.json
  - .omc/sessions/c9f8e296-c5c6-4b63-9077-1bef396f2cf5.json
  - cmd/model.go
  - cmd/plan.go
  - .omc/sessions/61c9a8ed-7bb4-4d6a-9108-6c2835c4560a.json
  - internal/config/config.go
  - cmd/archive.go
  - internal/config/defaults.go
  - .omc/sessions/5dab8e8e-7dec-40f8-b680-f2f8f554ad4b.json
  - .omc/state/idle-notif-cooldown.json
  - cmd/spec.go
  - cmd/design.go
  - internal/verifier/scenario.go
  - .omc/state/mission-state.json
tests:
  - internal/config/config_test.go
  - internal/spec/delta.go
  - cmd/archive_test.go
  - cmd/integration_test.go
  - internal/spec/parser.go
  - internal/spec/delta_test.go
  - internal/spec/merge.go
  - cmd/model_test.go
  - internal/verifier/scenario_test.go
  - internal/spec/parser_test.go
  - internal/spec/merge_test.go
  - internal/spec/schema.go
-->

---
### Requirement: Skipped task spec impact analysis output

The system SHALL provide a `--analyze-skipped` flag on the archive command that outputs the relationship between skipped tasks and their corresponding spec requirements in JSON format, without performing the archive.

#### Scenario: Analyze skipped tasks

- **GIVEN** a change with 2 skipped tasks that reference spec requirements
- **WHEN** the user runs `mysd archive --analyze-skipped`
- **THEN** the system SHALL output JSON listing each skipped task with its associated requirement names and skip reasons

#### Scenario: No skipped tasks

- **GIVEN** a change with all tasks completed (no `[~]` markers)
- **WHEN** the user runs `mysd archive --analyze-skipped`
- **THEN** the system SHALL output an empty JSON array

<!-- @trace
source: openspec-compliance
updated: 2026-04-02
code:
  - .omc/project-memory.json
  - .omc/state/last-tool-error.json
  - cmd/root.go
  - .omc/sessions/c3e46307-af7b-4b85-b578-57c643e1d6f6.json
  - .omc/sessions/c9f8e296-c5c6-4b63-9077-1bef396f2cf5.json
  - cmd/model.go
  - cmd/plan.go
  - .omc/sessions/61c9a8ed-7bb4-4d6a-9108-6c2835c4560a.json
  - internal/config/config.go
  - cmd/archive.go
  - internal/config/defaults.go
  - .omc/sessions/5dab8e8e-7dec-40f8-b680-f2f8f554ad4b.json
  - .omc/state/idle-notif-cooldown.json
  - cmd/spec.go
  - cmd/design.go
  - internal/verifier/scenario.go
  - .omc/state/mission-state.json
tests:
  - internal/config/config_test.go
  - internal/spec/delta.go
  - cmd/archive_test.go
  - cmd/integration_test.go
  - internal/spec/parser.go
  - internal/spec/delta_test.go
  - internal/spec/merge.go
  - cmd/model_test.go
  - internal/verifier/scenario_test.go
  - internal/spec/parser_test.go
  - internal/spec/merge_test.go
  - internal/spec/schema.go
-->