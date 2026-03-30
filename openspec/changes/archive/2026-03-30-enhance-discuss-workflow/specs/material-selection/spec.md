---
spec-version: "1.0"
capability: Material Selection
delta: MODIFIED
status: draft
---

## MODIFIED Requirements

### Requirement: Propose skill detects all available requirement sources

The `mysd:propose` and `mysd:discuss` (source-driven path) skills SHALL detect the following requirement sources during the material selection step:

1. **source_arg file/directory**: If `source_arg` is a valid file or directory path on disk
2. **Conversation context**: If the current conversation contains substantive requirement discussion (not greetings or meta-talk)
3. **Claude plan**: If conversation system messages contain a plan file path matching `~/.claude/plans/<name>.md` and the file exists
4. **gstack plan**: If `.md` files exist under `~/.gstack/projects/{project}/`
5. **Active change**: If `mysd status` reports an active change with existing `proposal.md`
6. **Deferred notes**: If `mysd note list` returns non-empty output

Each detected source SHALL be identified by type and a brief content preview (first line or title).

#### Scenario: Multiple sources detected

- **WHEN** the conversation contains prior discussion AND a Claude plan file exists
- **THEN** the skill SHALL list both sources with their type labels and content previews

#### Scenario: No sources detected

- **WHEN** no requirement sources are detected from any of the 6 types
- **THEN** the skill SHALL proceed directly to manual input without displaying an empty list

#### Scenario: Discuss source-driven path uses same detection

- **WHEN** the user enters the source-driven path in `mysd:discuss`
- **THEN** the skill SHALL detect sources using the same 6-type detection logic as `mysd:propose`

### Requirement: User selects requirement sources interactively

The `mysd:propose` and `mysd:discuss` (source-driven path) skills SHALL present all detected sources as a numbered list and allow the user to select one or more sources.

The list SHALL:
- Display only sources that have content (empty sources are omitted)
- Include "Manual input" as the last option, always present
- Allow multi-selection (e.g., "1,3" to select sources 1 and 3)

After selection, the skill SHALL read and aggregate the content from all selected sources into a single `aggregated_content` string.

#### Scenario: User selects multiple sources

- **WHEN** the user enters "1,3"
- **THEN** the skill SHALL read content from sources 1 and 3 and combine them into `aggregated_content`

#### Scenario: User selects manual input only

- **WHEN** the user selects the "Manual input" option
- **THEN** the skill SHALL prompt the user to describe their requirement and use that text as `aggregated_content`

#### Scenario: Only one source detected

- **WHEN** only one source is detected (plus manual input)
- **THEN** the skill SHALL still present the list for user confirmation, not auto-select

### Requirement: Auto mode uses all detected sources without interaction

When `auto_mode` is true, the `mysd:propose` and `mysd:discuss` skills SHALL automatically aggregate content from all detected sources without presenting a selection prompt.

If no sources are detected in auto mode, the skill SHALL extract requirements from conversation context as best-effort.

#### Scenario: Auto mode with detected sources

- **WHEN** `auto_mode` is true AND 3 sources are detected
- **THEN** the skill SHALL aggregate all 3 sources without prompting the user

#### Scenario: Auto mode with no sources

- **WHEN** `auto_mode` is true AND no sources are detected
- **THEN** the skill SHALL extract requirements from conversation context
