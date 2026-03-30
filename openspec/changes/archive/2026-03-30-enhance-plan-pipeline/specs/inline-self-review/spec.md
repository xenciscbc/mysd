---
spec-version: "1.0"
capability: Inline Self-Review
delta: ADDED
status: draft
---

## ADDED Requirements

### Requirement: Plan orchestrator inline self-review step

The `mysd:plan` skill SHALL execute an inline self-review step (Step 5a) after the planner completes (Step 5) and before the reviewer step (Step 5b).

The self-review SHALL be performed by the orchestrator directly (no agent spawn) using Read and Edit tools.

The self-review SHALL execute the following 4 checks in order:

1. **Placeholder check**: Scan tasks.md and design.md for TBD, TODO, FIXME, "implement later", "details to follow", and empty template sections. Fix each occurrence inline.
2. **Consistency check**: Verify that every capability in proposal.md has corresponding tasks, and that all file paths in tasks.md appear in proposal Impact or design.md. Fix mismatches inline.
3. **Scope check**: Warn if total tasks exceed 15, or if any single task description references more than 3 files. Do not auto-fix — display warning for user awareness.
4. **Ambiguity check**: Verify that task descriptions include specific verifiable conditions (not vague like "handle edge cases" or "add error handling"). Fix vague descriptions inline by adding specifics from the spec.

The self-review SHALL display a summary of fixes applied and warnings raised.

#### Scenario: Placeholders found and fixed

- **WHEN** the planner produces tasks.md containing "TBD" in a task description
- **THEN** the orchestrator SHALL replace "TBD" with specific content derived from the spec
- **AND** display "Fixed 1 placeholder in tasks.md"

#### Scenario: Consistency mismatch found

- **WHEN** proposal.md lists capability `material-selection` but no task has `spec: "material-selection"`
- **THEN** the orchestrator SHALL flag the mismatch
- **AND** add a task for the missing capability

#### Scenario: Scope warning raised

- **WHEN** tasks.md contains 18 tasks
- **THEN** the orchestrator SHALL display a warning: "18 tasks exceed recommended maximum of 15 — consider splitting the change"
- **AND** SHALL NOT auto-fix (proceed normally)

#### Scenario: All checks pass

- **WHEN** all 4 checks find no issues
- **THEN** the orchestrator SHALL display "Self-review passed" and proceed to the reviewer step

### Requirement: Self-review uses instructions checklist

The inline self-review step SHALL read the `selfReviewChecklist` from `mysd instructions tasks --change <name> --json` output to determine the checks to perform.

This ensures the self-review checklist is maintained in a single location (the instructions command) and is consistent between the agent-layer prevention (D-03) and the orchestrator-layer verification (D-02).

#### Scenario: Checklist loaded from instructions

- **WHEN** the orchestrator reaches the self-review step
- **THEN** it SHALL call `mysd instructions tasks --change <name> --json`
- **AND** use the `selfReviewChecklist` field to guide the 4-check process
