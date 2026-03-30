## ADDED Requirements

### Requirement: Plan pipeline includes analyze-fix loop after reviewer

The `mysd:plan` skill SHALL include an analyze-fix loop (Step 5c) after the reviewer step (Step 5b) and before the plan-checker step (Step 6).

The analyze-fix loop SHALL:
1. Run `mysd analyze {change_name} --json`
2. Filter findings to Critical and Warning severity only (ignore Suggestion)
3. If no Critical/Warning findings: display "Artifacts look consistent" and proceed
4. If Critical/Warning findings exist: fix each finding in the affected artifact, then re-run analyze
5. Repeat up to 2 total iterations
6. After 2 attempts, if findings remain: display remaining findings as a summary and proceed (SHALL NOT block the workflow)

#### Scenario: Analyze finds no issues

- **WHEN** `mysd analyze` returns no Critical or Warning findings
- **THEN** the plan skill SHALL display "Artifacts look consistent" and proceed to the next step

#### Scenario: Analyze-fix loop fixes issues within 2 iterations

- **WHEN** `mysd analyze` returns Critical/Warning findings on first run
- **AND** the findings are fixed and re-analysis passes
- **THEN** the plan skill SHALL display the fix count and proceed

#### Scenario: Analyze-fix loop exhausts 2 iterations

- **WHEN** `mysd analyze` still returns Critical/Warning findings after 2 fix iterations
- **THEN** the plan skill SHALL display remaining findings as a summary
- **AND** SHALL proceed to the next step without blocking

### Requirement: Plan pipeline supports optional design skip

The `mysd:plan` skill SHALL evaluate whether design.md is needed before the Design Phase (Step 4).

Design SHALL be skipped when ALL of the following conditions are met:
- Proposal Impact section lists 2 or fewer affected files
- Proposal has no New Capabilities (only Modified Capabilities or none)
- Proposal does not contain keywords: "cross-cutting", "migration", "architecture", "new pattern"

If `auto_mode` is false: the skill SHALL display the skip assessment and ask the user to confirm.
If `auto_mode` is true: the skill SHALL skip automatically when conditions are met.

When design is skipped, the plan skill SHALL proceed directly to the Planning Phase (Step 5) and pass an empty design content to the planner agent.

#### Scenario: Small change skips design in auto mode

- **WHEN** auto_mode is true
- **AND** proposal Impact lists 1 affected file
- **AND** proposal has no New Capabilities
- **THEN** the plan skill SHALL skip the Design Phase and proceed to Planning

#### Scenario: Cross-cutting change requires design

- **WHEN** proposal contains the word "cross-cutting" or "architecture"
- **THEN** the plan skill SHALL NOT skip the Design Phase regardless of file count

#### Scenario: User overrides design skip in interactive mode

- **WHEN** auto_mode is false
- **AND** skip conditions are met
- **AND** user chooses not to skip
- **THEN** the plan skill SHALL execute the Design Phase normally
