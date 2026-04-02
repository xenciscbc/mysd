# delta-spec-merge Specification

## Purpose

TBD - created by archiving change 'openspec-compliance'. Update Purpose after archive.

## Requirements

### Requirement: Delta spec merge on archive

The system SHALL merge delta specs from the change directory into main specs at `openspec/specs/<capability>/spec.md` during the archive process, before moving the change directory to the archive location.

#### Scenario: Merge ADDED requirements

- **GIVEN** a delta spec with `## ADDED Requirements` containing a new requirement
- **WHEN** the archive command executes
- **THEN** the system SHALL append the new requirement to the corresponding main spec file

#### Scenario: Merge MODIFIED requirements

- **GIVEN** a delta spec with `## MODIFIED Requirements` containing an updated requirement
- **WHEN** the archive command executes
- **THEN** the system SHALL replace the matching requirement in the main spec file with the updated content

#### Scenario: Merge REMOVED requirements

- **GIVEN** a delta spec with `## REMOVED Requirements` containing a requirement name
- **WHEN** the archive command executes
- **THEN** the system SHALL delete the matching requirement block from the main spec file

#### Scenario: Merge RENAMED requirements

- **GIVEN** a delta spec with `## RENAMED Requirements` containing FROM/TO entries
- **WHEN** the archive command executes
- **THEN** the system SHALL rename the matching requirement heading in the main spec file

#### Scenario: No matching main spec exists

- **GIVEN** a delta spec for a capability that has no existing main spec
- **WHEN** the archive command executes
- **THEN** the system SHALL create a new main spec file at `openspec/specs/<capability>/spec.md` containing the ADDED requirements


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
### Requirement: Delta merge operation order

The system SHALL apply delta operations in the following strict order: RENAMED first, then REMOVED, then MODIFIED, then ADDED.

#### Scenario: RENAMED applied before MODIFIED

- **GIVEN** a delta spec with both RENAMED (A → B) and MODIFIED (B) operations
- **WHEN** the merge executes
- **THEN** the system SHALL rename A to B first, then apply the modification to B


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
### Requirement: Archive path with date prefix

The `mysd archive` command SHALL write archived changes to `openspec/changes/archive/YYYY-MM-DD-<changeName>/` instead of `openspec/archive/<changeName>/`.

#### Scenario: Archive creates date-prefixed directory

- **GIVEN** a change named `add-auth` archived on 2026-04-02
- **WHEN** the archive command executes
- **THEN** the change directory SHALL be moved to `openspec/changes/archive/2026-04-02-add-auth/`


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
### Requirement: Spec frontmatter on merge

The system SHALL ensure main spec files contain YAML frontmatter with `name`, `description`, `version`, and `generatedBy` fields after merge. When a MODIFIED operation is applied, the `version` field SHALL be incremented (minor version +1).

#### Scenario: Version increment on MODIFIED

- **GIVEN** a main spec with `version: 1.0.0` and a delta spec with MODIFIED requirements
- **WHEN** the merge executes
- **THEN** the main spec's `version` field SHALL be updated to `1.1.0`

#### Scenario: New spec gets initial frontmatter

- **GIVEN** a delta spec for a capability with no existing main spec
- **WHEN** the merge creates the new main spec
- **THEN** the spec SHALL have frontmatter with `version: 1.0.0` and `generatedBy: mysd v<current-version>`


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
### Requirement: RENAMED delta operation support

The delta spec parser SHALL support `## RENAMED Requirements` sections with `### FROM: <old>` and `### TO: <new>` heading pairs.

#### Scenario: Parse RENAMED section

- **GIVEN** a delta spec body containing `## RENAMED Requirements` with `### FROM: Old Name` and `### TO: New Name`
- **WHEN** `ParseDelta` is called
- **THEN** the function SHALL return a renamed entry with From="Old Name" and To="New Name"


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
### Requirement: Scenario GIVEN validation

The system SHALL validate that spec scenarios include **GIVEN**, **WHEN**, and **THEN** keywords. Missing any keyword SHALL produce a warning.

#### Scenario: Scenario missing GIVEN

- **GIVEN** a spec file with a scenario containing only WHEN and THEN
- **WHEN** scenario validation runs
- **THEN** the system SHALL emit a warning indicating the scenario is missing GIVEN

#### Scenario: Complete scenario passes validation

- **GIVEN** a spec file with a scenario containing GIVEN, WHEN, and THEN
- **WHEN** scenario validation runs
- **THEN** the system SHALL emit no warning for that scenario


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
### Requirement: Merge failure handling

The system SHALL NOT block the archive process when a delta spec merge fails for a specific requirement. The system SHALL emit a warning and continue with the remaining operations.

#### Scenario: Heading mismatch during MODIFIED

- **GIVEN** a MODIFIED requirement whose heading does not match any requirement in the main spec
- **WHEN** the merge executes
- **THEN** the system SHALL emit a warning and skip that modification without blocking the archive

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