---
spec-version: "1.0"
capability: Reviewer Agent
delta: ADDED
status: draft
---

## Requirement: Reviewer includes rationalization table

The `mysd-reviewer` agent SHALL include a Rationalization Table before executing quality checks. The table SHALL list common anti-patterns where the reviewer might skip or weaken checks, paired with the correct action to take.

The table SHALL include entries for at least:
- Skipping checks because "requirements are clear enough"
- Accepting placeholders because "they'll be filled in later"
- Skipping boundary conditions because "the requirement is obvious"
- Skipping scope check because "the change is small"
- Ignoring minor inconsistencies
- Skimming validate output

### Scenario: Rationalization table is present in agent definition

- **WHEN** the `mysd-reviewer` agent definition is loaded
- **THEN** the Rationalization Table SHALL appear between Step 1 (Load Artifacts) and Step 2 (Check 1 — No Placeholders)

### Scenario: Reviewer applies rationalization checks

- **WHEN** the reviewer encounters an artifact with a placeholder marked "TBD"
- **THEN** the reviewer SHALL NOT skip the fix because "it will be filled in later"
- **AND** SHALL replace it with specific content inferred from context

## Requirement: Reviewer agent performs artifact quality checks

The `mysd-reviewer` agent SHALL scan all artifacts for the given phase and fix quality issues inline.

The agent SHALL accept the following input context:
- `change_name`: The change to review
- `phase`: `"propose"` (proposal + specs) or `"plan"` (all 4 artifacts)
- `validate_output`: Output from `mysd validate` (empty string if unavailable)
- `auto_mode`: Boolean — if true, fix silently; if false, include issues in summary

The agent SHALL be invocable from both the `propose` and `plan` pipelines. The `propose` pipeline SHALL invoke the reviewer after spec generation (Step 12), and the `plan` pipeline SHALL invoke it after task generation (Step 5b).

### Scenario: Plan phase loads all 4 artifacts

- **WHEN** `mysd-reviewer` is invoked with `phase: "plan"`
- **THEN** it SHALL read `proposal.md`, all `specs/*/spec.md`, `design.md`, and `tasks.md`

### Scenario: Propose phase loads 2 artifacts

- **WHEN** `mysd-reviewer` is invoked with `phase: "propose"`
- **THEN** it SHALL read `proposal.md` and all `specs/*/spec.md` only

### Scenario: Propose pipeline invokes reviewer with validate output

- **WHEN** the propose skill reaches Step 12
- **THEN** it SHALL run `mysd validate {change_name}` and capture the output
- **AND** spawn `mysd-reviewer` with `phase: "propose"` and the captured `validate_output`
- **AND** use the `reviewer_model` resolved from the current profile

## Requirement: Reviewer checks for placeholder content

The `mysd-reviewer` agent SHALL detect and fix placeholder content in all loaded artifacts.

Placeholder patterns that SHALL be detected and fixed:
- Literal strings: `TBD`, `TODO`, `FIXME`, `implement later`, `details to follow`
- Vague instructions without specifics: "Add appropriate error handling", "Handle edge cases"
- Empty template sections left unfilled
- Weasel quantities: "some", "various", "several" when a specific list is needed

### Scenario: TBD placeholder is fixed inline

- **WHEN** an artifact contains the text `TBD` in a requirement field
- **THEN** the reviewer SHALL replace it with a specific value using context from other artifacts
- **AND** count it as a fixed issue in the summary

## Requirement: Reviewer checks internal consistency

The `mysd-reviewer` agent SHALL verify cross-artifact consistency.

For `phase: "propose"`:
- Every capability in `proposal.md` SHALL have a corresponding `specs/<capability>/spec.md`
- Specs SHALL reference only capabilities described in the proposal
- File paths and component names SHALL be consistent across proposal and specs

For `phase: "plan"` (additional):
- Design SHALL reference only capabilities from the proposal
- Tasks SHALL cover all design decisions, and nothing outside proposal scope
- File paths SHALL be consistent across proposal Impact, design, and tasks

### Scenario: Missing spec for proposal capability is flagged

- **WHEN** `proposal.md` lists a capability `foo-bar` with no corresponding `specs/foo-bar/spec.md`
- **THEN** the reviewer SHALL flag this as a cannot-auto-fix issue in the summary

## Requirement: Reviewer checks scope

The `mysd-reviewer` agent SHALL flag over-scoped artifacts that cannot be auto-fixed.

- For `phase: "propose"`: more than 15 MUST requirements SHALL trigger a decomposition warning
- For `phase: "plan"`: more than 15 pending tasks SHALL trigger a decomposition warning
- Any single item touching more than 3 unrelated subsystems SHALL be flagged

### Scenario: Over-scoped plan triggers warning

- **WHEN** `tasks.md` contains more than 15 pending tasks
- **THEN** the reviewer SHALL include a cannot-auto-fix warning recommending decomposition

## Requirement: Reviewer checks ambiguity

The `mysd-reviewer` agent SHALL detect and fix ambiguous requirements.

- Success/failure conditions SHALL be testable and specific
- Boundary conditions SHALL be defined (empty input, max limits, error cases)
- "The system" SHALL NOT refer ambiguously to multiple components — reviewer SHALL replace with specific component names

### Scenario: Ambiguous system reference is fixed

- **WHEN** an artifact contains "the system should handle X" without specifying which component
- **THEN** the reviewer SHALL replace with the specific component name inferred from context

## Requirement: Reviewer returns structured summary

The `mysd-reviewer` agent SHALL return a structured summary after completing all checks.

Summary format:
```
## Review Results
- Phase: {phase}
- Issues fixed: N
- Fixed: {list of what was fixed, or "None"}
- Cannot auto-fix (structural): {list if any, or "None"}
```

### Scenario: Summary reports zero issues

- **WHEN** no issues are found across all 4 checks
- **THEN** the summary SHALL report `Issues fixed: 0` and `Fixed: None`

### Scenario: Summary reports structural issues

- **WHEN** a scope or consistency issue cannot be auto-fixed
- **THEN** the summary SHALL list it under `Cannot auto-fix (structural)`
- **AND** SHALL NOT block the calling workflow from proceeding

## Requirements

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

<!-- @trace
source: enhance-propose-reviewer
updated: 2026-03-28
code: []
tests: []
-->


<!-- @trace
source: enhance-propose-reviewer
updated: 2026-03-28
code:
  - mysd/agents/mysd-reviewer.md
  - mysd/skills/propose/SKILL.md
-->

---
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


<!-- @trace
source: enhance-propose-reviewer
updated: 2026-03-28
code:
  - mysd/agents/mysd-reviewer.md
  - mysd/skills/propose/SKILL.md
-->

---
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


<!-- @trace
source: add-mysd-reviewer-agent
updated: 2026-03-28
code:
  - mysd/skills/plan/SKILL.md
  - mysd/skills/propose/SKILL.md
  - internal/config/config.go
  - mysd/skills/model/SKILL.md
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/lang/SKILL.md
  - cmd/plan.go
  - mysd/skills/init/SKILL.md
  - mysd/agents/mysd-reviewer.md
tests:
  - cmd/plan_test.go
  - internal/config/config_test.go
-->

---
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


<!-- @trace
source: add-mysd-reviewer-agent
updated: 2026-03-28
code:
  - mysd/skills/plan/SKILL.md
  - mysd/skills/propose/SKILL.md
  - internal/config/config.go
  - mysd/skills/model/SKILL.md
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/lang/SKILL.md
  - cmd/plan.go
  - mysd/skills/init/SKILL.md
  - mysd/agents/mysd-reviewer.md
tests:
  - cmd/plan_test.go
  - internal/config/config_test.go
-->

---
### Requirement: Reviewer checks scope

The `mysd-reviewer` agent SHALL flag over-scoped artifacts that cannot be auto-fixed.

- For `phase: "propose"`: more than 15 MUST requirements SHALL trigger a decomposition warning
- For `phase: "plan"`: more than 15 pending tasks SHALL trigger a decomposition warning
- Any single item touching more than 3 unrelated subsystems SHALL be flagged

#### Scenario: Over-scoped plan triggers warning

- **WHEN** `tasks.md` contains more than 15 pending tasks
- **THEN** the reviewer SHALL include a cannot-auto-fix warning recommending decomposition


<!-- @trace
source: add-mysd-reviewer-agent
updated: 2026-03-28
code:
  - mysd/skills/plan/SKILL.md
  - mysd/skills/propose/SKILL.md
  - internal/config/config.go
  - mysd/skills/model/SKILL.md
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/lang/SKILL.md
  - cmd/plan.go
  - mysd/skills/init/SKILL.md
  - mysd/agents/mysd-reviewer.md
tests:
  - cmd/plan_test.go
  - internal/config/config_test.go
-->

---
### Requirement: Reviewer checks ambiguity

The `mysd-reviewer` agent SHALL detect and fix ambiguous requirements.

- Success/failure conditions SHALL be testable and specific
- Boundary conditions SHALL be defined (empty input, max limits, error cases)
- "The system" SHALL NOT refer ambiguously to multiple components — reviewer SHALL replace with specific component names

#### Scenario: Ambiguous system reference is fixed

- **WHEN** an artifact contains "the system should handle X" without specifying which component
- **THEN** the reviewer SHALL replace with the specific component name inferred from context


<!-- @trace
source: add-mysd-reviewer-agent
updated: 2026-03-28
code:
  - mysd/skills/plan/SKILL.md
  - mysd/skills/propose/SKILL.md
  - internal/config/config.go
  - mysd/skills/model/SKILL.md
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/lang/SKILL.md
  - cmd/plan.go
  - mysd/skills/init/SKILL.md
  - mysd/agents/mysd-reviewer.md
tests:
  - cmd/plan_test.go
  - internal/config/config_test.go
-->

---
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

<!-- @trace
source: add-mysd-reviewer-agent
updated: 2026-03-28
code:
  - mysd/skills/plan/SKILL.md
  - mysd/skills/propose/SKILL.md
  - internal/config/config.go
  - mysd/skills/model/SKILL.md
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/lang/SKILL.md
  - cmd/plan.go
  - mysd/skills/init/SKILL.md
  - mysd/agents/mysd-reviewer.md
tests:
  - cmd/plan_test.go
  - internal/config/config_test.go
-->