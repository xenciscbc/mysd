---
spec-version: "1.0"
capability: Artifact Instructions
delta: ADDED
status: draft
---

## ADDED Requirements

### Requirement: mysd instructions CLI command

The `mysd instructions` command SHALL accept an artifact ID and change name, and output structured JSON containing template, rules, instruction, dependencies, and self-review checklist for the specified artifact.

Supported artifact IDs: `design`, `tasks`.

The command SHALL accept the following flags:
- `--change <name>`: The change name (required)
- `--json`: Output as JSON (required)

The output JSON SHALL include these fields:
- `artifactId`: The requested artifact ID
- `changeName`: The change name
- `outputPath`: The file path where the artifact should be written
- `template`: The structural template for the artifact
- `rules`: Array of constraint strings the agent must follow
- `instruction`: Artifact-specific guidance text
- `dependencies`: Array of dependency objects with `id`, `path`, and `done` fields
- `selfReviewChecklist`: Array of checklist items for the agent to verify before completing

The command SHALL read the current change state to determine which dependencies are satisfied (`done: true` vs `done: false`).

The command SHALL return an error if the change does not exist or the artifact ID is not recognized.

#### Scenario: Instructions for tasks artifact

- **WHEN** `mysd instructions tasks --change my-change --json` is executed
- **AND** the change `my-change` exists with completed design and specs
- **THEN** the output SHALL include `artifactId: "tasks"`, the tasks template, rules for task writing, and dependencies showing design and specs as done

#### Scenario: Instructions for design artifact

- **WHEN** `mysd instructions design --change my-change --json` is executed
- **AND** the change `my-change` exists with completed proposal
- **THEN** the output SHALL include `artifactId: "design"`, the design template, and proposal as a completed dependency

#### Scenario: Unknown artifact ID

- **WHEN** `mysd instructions unknown --change my-change --json` is executed
- **THEN** the command SHALL return an error indicating the artifact ID is not recognized

### Requirement: Self-review checklist in instructions output

The `selfReviewChecklist` field in the instructions output SHALL contain artifact-specific quality checks that the consuming agent is expected to verify.

For the `tasks` artifact, the checklist SHALL include:
- No TBD/TODO/FIXME placeholders in task descriptions
- Every MUST requirement has at least one task with a matching `satisfies` entry
- No single task targets more than 3 files
- All file paths referenced in tasks appear in proposal Impact or design
- Task dependencies form a valid DAG (no circular references)

For the `design` artifact, the checklist SHALL include:
- No TBD/TODO/FIXME placeholders
- Every capability in the proposal has a corresponding section in the design
- Decision rationale includes at least one alternative considered
- File paths referenced in the design are consistent with proposal Impact

#### Scenario: Tasks checklist includes coverage check

- **WHEN** instructions for the `tasks` artifact are requested
- **THEN** the `selfReviewChecklist` SHALL include an item about MUST requirement coverage

#### Scenario: Design checklist includes consistency check

- **WHEN** instructions for the `design` artifact are requested
- **THEN** the `selfReviewChecklist` SHALL include an item about proposal-to-design consistency
