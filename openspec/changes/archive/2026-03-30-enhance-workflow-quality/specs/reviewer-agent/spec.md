## MODIFIED Requirements

### Requirement: Reviewer agent performs artifact quality checks

The `mysd-reviewer` agent SHALL scan all artifacts for the given phase and fix quality issues inline.

The agent SHALL accept the following input context:
- `change_name`: The change to review
- `phase`: `"propose"` (proposal + specs) or `"plan"` (all 4 artifacts)
- `validate_output`: Output from `mysd validate` (empty string if unavailable)
- `auto_mode`: Boolean — if true, fix silently; if false, include issues in summary
- `change_type`: Optional string (`"feature"`, `"bugfix"`, `"refactor"`) — when provided, the reviewer SHALL verify that the proposal uses the correct template structure for the given type

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

#### Scenario: Reviewer validates proposal template matches change type

- **WHEN** `change_type` is `"bugfix"` and proposal uses Feature template (Why/What Changes/Capabilities)
- **THEN** the reviewer SHALL flag this as a cannot-auto-fix issue recommending the correct template (Problem/Root Cause/Proposed Solution/Success Criteria)
