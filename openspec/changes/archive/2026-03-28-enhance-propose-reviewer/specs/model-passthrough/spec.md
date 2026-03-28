## MODIFIED Requirements

### Requirement: Command skills pass model to agents

Workflow command skills (`propose`, `discuss`, `plan`, `apply`, `ff`, `ffe`) SHALL read the `model` field from the binary's `--context-only` JSON output and pass it as the `model` parameter when spawning agent tasks.

The `propose` skill SHALL additionally resolve a `reviewer_model` from the current profile's reviewer role mapping. Since `propose` uses `mysd model` (not `--context-only` JSON), the skill SHALL derive `reviewer_model` from the profile name:

| Profile | reviewer_model |
|---------|---------------|
| quality | opus |
| balanced | sonnet |
| budget | sonnet |

#### Scenario: Plan command spawns designer with profile model

- **WHEN** `/mysd:plan` reads context JSON containing `"model": "opus"`
- **THEN** the command SHALL spawn `mysd-designer` with `model: opus`

#### Scenario: Propose command resolves reviewer_model from profile

- **WHEN** `/mysd:propose` reads `Profile: balanced` from `mysd model` output
- **THEN** the skill SHALL set `reviewer_model` to `sonnet`
- **AND** spawn `mysd-reviewer` with `model: sonnet`

#### Scenario: Propose command resolves reviewer_model for quality profile

- **WHEN** `/mysd:propose` reads `Profile: quality` from `mysd model` output
- **THEN** the skill SHALL set `reviewer_model` to `opus`
- **AND** spawn `mysd-reviewer` with `model: opus`
