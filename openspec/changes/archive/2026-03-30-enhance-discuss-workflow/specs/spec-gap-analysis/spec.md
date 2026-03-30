---
spec-version: "1.0"
capability: Spec Gap Analysis
delta: ADDED
status: draft
---

## ADDED Requirements

### Requirement: Gap analysis evaluates requirement coverage against proposal

When a spec is selected for focused discussion, the `mysd:discuss` skill SHALL analyze requirement coverage by comparing the spec's requirements against the change's `proposal.md` Capabilities section.

The analysis SHALL identify:
- Capabilities listed in the proposal that have no corresponding requirement in the selected spec
- Requirements in the spec that do not map to any capability in the proposal (orphaned requirements)

The results SHALL be presented as a coverage summary listing each gap.

#### Scenario: Proposal capability missing from spec

- **WHEN** the proposal lists capability "source detection" under the spec's capability area
- **AND** the spec has no requirement covering source detection
- **THEN** the gap analysis SHALL report "source detection" as an uncovered capability

#### Scenario: Full coverage

- **WHEN** every capability in the proposal has at least one corresponding requirement in the spec
- **THEN** the gap analysis SHALL report full coverage with no gaps

### Requirement: Gap analysis evaluates scenario completeness

The `mysd:discuss` skill SHALL check that every requirement in the selected spec has at least one scenario defined.

A requirement with zero scenarios SHALL be reported as a scenario gap.

#### Scenario: Requirement without scenario

- **WHEN** the spec contains "Requirement: Auto mode behavior" with no scenario blocks
- **THEN** the gap analysis SHALL report this requirement as missing scenarios

#### Scenario: All requirements have scenarios

- **WHEN** every requirement in the spec has at least one `#### Scenario:` block
- **THEN** the gap analysis SHALL report scenario completeness as satisfied

### Requirement: Gap analysis evaluates boundary condition coverage

The `mysd:discuss` skill SHALL evaluate whether scenarios cover boundary conditions beyond the happy path.

For each requirement, the analysis SHALL check for the presence of:
- At least one error or failure scenario (keywords: "error", "fail", "invalid", "not found", "empty", "missing")
- At least one edge case scenario (keywords: "maximum", "minimum", "empty list", "single item", "concurrent", "timeout")

Requirements with only happy-path scenarios SHALL be flagged as missing boundary coverage.

#### Scenario: Only happy path scenarios exist

- **WHEN** a requirement has 2 scenarios but neither contains error or edge case keywords
- **THEN** the gap analysis SHALL flag it as "missing boundary condition coverage"

#### Scenario: Error and edge case scenarios present

- **WHEN** a requirement has scenarios covering both a success path and an error condition
- **THEN** the gap analysis SHALL report boundary coverage as satisfied for that requirement

### Requirement: Gap analysis results drive the discussion starting point

After completing the three-dimension analysis, the `mysd:discuss` skill SHALL present the results as a structured summary and use it as the discussion starting point.

The summary SHALL list gaps grouped by dimension (coverage, scenario, boundary) with specific gap descriptions. The skill SHALL then ask the user which gap to address first, rather than asking an open-ended "what do you want to discuss?"

#### Scenario: Gaps found drive focused discussion

- **WHEN** the gap analysis finds 2 coverage gaps and 1 scenario gap
- **THEN** the skill SHALL present all 3 gaps and ask the user which to address first

#### Scenario: No gaps found

- **WHEN** the gap analysis finds no gaps in any dimension
- **THEN** the skill SHALL report the spec as complete and ask if the user wants to discuss other aspects or end the discussion
