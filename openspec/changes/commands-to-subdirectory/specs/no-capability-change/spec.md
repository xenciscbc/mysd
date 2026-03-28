## ADDED Requirements

### Requirement: Plugin command files use subdirectory structure

The mysd plugin command files SHALL be organized under `plugin/commands/mysd/` subdirectory instead of flat `plugin/commands/mysd-*.md` files.

#### Scenario: Command file location

- **WHEN** the plugin is installed
- **THEN** command files SHALL exist at `plugin/commands/mysd/<name>.md` (e.g., `plugin/commands/mysd/apply.md`)

#### Scenario: Slash command invocation unchanged

- **WHEN** a user invokes `/mysd:apply`
- **THEN** Claude Code SHALL resolve the command from `plugin/commands/mysd/apply.md`
