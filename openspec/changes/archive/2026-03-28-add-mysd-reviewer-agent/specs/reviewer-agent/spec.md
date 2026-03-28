## ADDED Requirements

### Requirement: Reviewer agent performs artifact quality checks

The `mysd-reviewer` agent SHALL scan all artifacts for the given phase and fix quality issues inline.

The agent SHALL accept the following input context:
- `change_name`: The change to review
- `phase`: `"propose"` (proposal + specs) or `"plan"` (all 4 artifacts)
- `validate_output`: Output from `mysd validate` (empty string if unavailable)
- `auto_mode`: Boolean — if true, fix silently; if false, include issues in summary

#### Scenario: Plan phase loads all 4 artifacts

- **WHEN** `mysd-reviewer` is invoked with `phase: "plan"`
- **THEN** it SHALL read `proposal.md`, all `specs/*/spec.md`, `design.md`, and `tasks.md`

#### Scenario: Propose phase loads 2 artifacts

- **WHEN** `mysd-reviewer` is invoked with `phase: "propose"`
- **THEN** it SHALL read `proposal.md` and all `specs/*/spec.md` only

### Requirement: Reviewer checks for placeholder content

The `mysd-reviewer` agent SHALL detect and fix placeholder content in all loaded artifacts.

Placeholder patterns that SHALL be detected and fixed:
- Literal strings: `TBD`, `TODO`, `FIXME`, `implement later`, `details to follow`
- Vague instructions without specifics: "Add appropriate error handling", "Handle edge cases"
- Empty template sections left unfilled
- Weasel quantities: "some", "various", "several" when a specific list is needed

#### Scenario: TBD placeholder is fixed inline

- **WHEN** an artifact contains the text `TBD` in a requirement field
- **THEN** the reviewer SHALL replace it with a specific value using context from other artifacts
- **AND** count it as a fixed issue in the summary

### Requirement: Reviewer checks internal consistency

The `mysd-reviewer` agent SHALL verify cross-artifact consistency.

For `phase: "propose"`:
- Every capability in `proposal.md` SHALL have a corresponding `specs/<capability>/spec.md`
- Specs SHALL reference only capabilities described in the proposal
- File paths and component names SHALL be consistent across proposal and specs

For `phase: "plan"` (additional):
- Design SHALL reference only capabilities from the proposal
- Tasks SHALL cover all design decisions, and nothing outside proposal scope
- File paths SHALL be consistent across proposal Impact, design, and tasks

#### Scenario: Missing spec for proposal capability is flagged

- **WHEN** `proposal.md` lists a capability `foo-bar` with no corresponding `specs/foo-bar/spec.md`
- **THEN** the reviewer SHALL flag this as a cannot-auto-fix issue in the summary

### Requirement: Reviewer checks scope

The `mysd-reviewer` agent SHALL flag over-scoped artifacts that cannot be auto-fixed.

- For `phase: "propose"`: more than 15 MUST requirements SHALL trigger a decomposition warning
- For `phase: "plan"`: more than 15 pending tasks SHALL trigger a decomposition warning
- Any single item touching more than 3 unrelated subsystems SHALL be flagged

#### Scenario: Over-scoped plan triggers warning

- **WHEN** `tasks.md` contains more than 15 pending tasks
- **THEN** the reviewer SHALL include a cannot-auto-fix warning recommending decomposition

### Requirement: Reviewer checks ambiguity

The `mysd-reviewer` agent SHALL detect and fix ambiguous requirements.

- Success/failure conditions SHALL be testable and specific
- Boundary conditions SHALL be defined (empty input, max limits, error cases)
- "The system" SHALL NOT refer ambiguously to multiple components — reviewer SHALL replace with specific component names

#### Scenario: Ambiguous system reference is fixed

- **WHEN** an artifact contains "the system should handle X" without specifying which component
- **THEN** the reviewer SHALL replace with the specific component name inferred from context

### Requirement: Reviewer returns structured summary

The `mysd-reviewer` agent SHALL return a structured summary after completing all checks.

Summary format:
```
## Review Results
- Phase: {phase}
- Issues fixed: N
- Fixed: {list of what was fixed, or "None"}
- Cannot auto-fix (structural): {list if any, or "None"}
```

#### Scenario: Summary reports zero issues

- **WHEN** no issues are found across all 4 checks
- **THEN** the summary SHALL report `Issues fixed: 0` and `Fixed: None`

#### Scenario: Summary reports structural issues

- **WHEN** a scope or consistency issue cannot be auto-fixed
- **THEN** the summary SHALL list it under `Cannot auto-fix (structural)`
- **AND** SHALL NOT block the calling workflow from proceeding
