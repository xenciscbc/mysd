## ADDED Requirements

### Requirement: Standalone docs update skill

The system SHALL provide a `/mysd:docs update` skill that triggers documentation updates independently from the archive/ff/ffe workflows.

#### Scenario: Default scope — latest archived change

- **WHEN** the user invokes `/mysd:docs update` with no arguments
- **THEN** the system SHALL read the most recent archived change from `openspec/changes/archive/` and update all files in `docs_to_update` using that change as context

#### Scenario: Specified change scope

- **WHEN** the user invokes `/mysd:docs update --change <name>`
- **THEN** the system SHALL locate the archived change matching `<name>` in `openspec/changes/archive/` and update docs using that change as context

#### Scenario: Last N changes scope

- **WHEN** the user invokes `/mysd:docs update --last N`
- **THEN** the system SHALL read the N most recent archived changes (sorted by date prefix) and update docs using the combined context of all N changes

#### Scenario: Full codebase scan scope

- **WHEN** the user invokes `/mysd:docs update --full`
- **THEN** the system SHALL scan the current codebase (source files, commands, configuration) and update docs to reflect the actual state of the project, without relying on archived change context

#### Scenario: Free-text description scope

- **WHEN** the user invokes `/mysd:docs update "description text"`
- **THEN** the system SHALL use the provided description as the update context and update docs accordingly

#### Scenario: No docs_to_update configured

- **WHEN** the user invokes `/mysd:docs update` and `docs_to_update` is empty
- **THEN** the system SHALL inform the user that no files are configured and suggest using `mysd docs add <path>`

### Requirement: UAT prompt removal

The `mysd archive` command SHALL NOT display an interactive UAT prompt. The `Run UAT first? [y/N]` prompt and associated `isInteractive()` check SHALL be removed.

#### Scenario: Archive runs without UAT prompt

- **WHEN** the user runs `mysd archive` in an interactive terminal
- **THEN** the command SHALL proceed directly to gate checks without prompting about UAT

### Requirement: Archive SKILL.md path accuracy

The `/mysd:archive` SKILL.md SHALL reference the correct archive path format `openspec/changes/archive/YYYY-MM-DD-<changeName>/` in all path references and user-facing messages.

#### Scenario: SKILL.md reads archived change context

- **WHEN** the archive SKILL.md reads context for doc maintenance
- **THEN** it SHALL read from `openspec/changes/archive/YYYY-MM-DD-<changeName>/` instead of `.specs/archive/<changeName>/`

#### Scenario: Success message shows correct path

- **WHEN** the archive completes successfully
- **THEN** the SKILL.md SHALL display the archive location as `openspec/changes/archive/YYYY-MM-DD-<changeName>/`
