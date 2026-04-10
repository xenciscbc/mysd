---
spec-version: "1.0"
capability: delta-spec-merge
delta: MODIFIED
status: pending
---

## MODIFIED Requirements

### Requirement: Delta spec merge on archive

The system SHALL merge delta specs from the change directory into main specs at `openspec/specs/<capability>/spec.md` during the archive process, before moving the change directory to the archive location.

When `ParseDelta` finds no delta section headings (`## ADDED Requirements`, `## MODIFIED Requirements`, `## REMOVED Requirements`, `## RENAMED Requirements`) in the delta spec body, the system SHALL fall back to frontmatter-based merge behavior:
- If the delta spec frontmatter `delta` field is `ADDED`, the system SHALL use the entire delta body as the new main spec content.
- If the delta spec frontmatter `delta` field is `MODIFIED`, the system SHALL replace the entire main spec body with the delta body content.
- If the delta spec frontmatter `delta` field is any other value or empty, the system SHALL emit a warning and skip the merge for that capability.

This fallback ensures that specs written in plain format (without delta section headings) are correctly merged instead of silently producing empty main specs.

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

#### Scenario: Fallback for ADDED delta without section headings

- **GIVEN** a delta spec with frontmatter `delta: ADDED` and body content in plain spec format (no `## ADDED Requirements` heading)
- **AND** no existing main spec for this capability
- **WHEN** the archive command executes
- **THEN** the system SHALL create a new main spec using the entire delta body as content

#### Scenario: Fallback for MODIFIED delta without section headings

- **GIVEN** a delta spec with frontmatter `delta: MODIFIED` and body content in plain spec format (no delta section headings)
- **AND** an existing main spec for this capability
- **WHEN** the archive command executes
- **THEN** the system SHALL replace the main spec body with the entire delta body content
- **AND** the main spec frontmatter version SHALL be incremented

#### Scenario: Fallback emits warning for unknown delta type

- **GIVEN** a delta spec with frontmatter `delta: ""` (empty) and no delta section headings
- **WHEN** the archive command executes
- **THEN** the system SHALL emit a warning indicating the delta spec has no parseable operations
- **AND** the system SHALL skip the merge for that capability
