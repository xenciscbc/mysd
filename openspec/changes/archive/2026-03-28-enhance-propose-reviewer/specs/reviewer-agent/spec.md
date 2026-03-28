## ADDED Requirements

### Requirement: Reviewer includes rationalization table

The `mysd-reviewer` agent SHALL include a Rationalization Table before executing quality checks. The table SHALL list common anti-patterns where the reviewer might skip or weaken checks, paired with the correct action to take.

The table SHALL include entries for at least:
- Skipping checks because "requirements are clear enough"
- Accepting placeholders because "they'll be filled in later"
- Skipping boundary conditions because "the requirement is obvious"
- Skipping scope check because "the change is small"
- Ignoring minor inconsistencies
- Skimming validate output

#### Scenario: Rationalization table is present in agent definition

- **WHEN** the `mysd-reviewer` agent definition is loaded
- **THEN** the Rationalization Table SHALL appear between Step 1 (Load Artifacts) and Step 2 (Check 1 — No Placeholders)

#### Scenario: Reviewer applies rationalization checks

- **WHEN** the reviewer encounters an artifact with a placeholder marked "TBD"
- **THEN** the reviewer SHALL NOT skip the fix because "it will be filled in later"
- **AND** SHALL replace it with specific content inferred from context

## MODIFIED Requirements

### Requirement: Reviewer agent performs artifact quality checks

The `mysd-reviewer` agent SHALL scan all artifacts for the given phase and fix quality issues inline.

The agent SHALL accept the following input context:
- `change_name`: The change to review
- `phase`: `"propose"` (proposal + specs) or `"plan"` (all 4 artifacts)
- `validate_output`: Output from `mysd validate` (empty string if unavailable)
- `auto_mode`: Boolean — if true, fix silently; if false, include issues in summary

The agent SHALL be invocable from both the `propose` and `plan` pipelines. The `propose` pipeline SHALL invoke the reviewer after spec generation (Step 12), and the `plan` pipeline SHALL invoke it after task generation (Step 5b).

#### Scenario: Plan phase loads all 4 artifacts

- **WHEN** `mysd-reviewer` is invoked with `phase: "plan"`
- **THEN** it SHALL read `proposal.md`, all `specs/*/spec.md`, `design.md`, and `tasks.md`

#### Scenario: Propose phase loads 2 artifacts

- **WHEN** `mysd-reviewer` is invoked with `phase: "propose"`
- **THEN** it SHALL read `proposal.md` and all `specs/*/spec.md` only

#### Scenario: Propose pipeline invokes reviewer with validate output

- **WHEN** the propose skill reaches Step 12
- **THEN** it SHALL run `mysd validate {change_name}` and capture the output
- **AND** spawn `mysd-reviewer` with `phase: "propose"` and the captured `validate_output`
- **AND** use the `reviewer_model` resolved from the current profile
