---
spec-version: "1.0"
capability: Requirement Interview
delta: ADDED
status: draft
---

### Requirement: Orchestrator evaluates requirement completeness

After material selection, the `mysd:propose` orchestrator SHALL evaluate the `aggregated_content` against three completeness dimensions:

| Dimension | Complete when |
|-----------|--------------|
| Problem | The content describes the problem being solved, not just a desired solution |
| Boundary | The content states what is in scope and what is explicitly excluded |
| Success Criteria | The content includes specific, verifiable conditions for success |

The orchestrator SHALL also compare `aggregated_content` against related existing specs (from the spec scan step). If overlap is detected between the proposed change and an existing spec, the orchestrator SHALL flag it as an additional interview question.

#### Scenario: All dimensions are covered

- **WHEN** the aggregated content clearly addresses problem, boundary, and success criteria
- **THEN** the orchestrator SHALL skip the interview and produce the requirement_brief directly

#### Scenario: Problem dimension is missing

- **WHEN** the aggregated content describes a desired solution but not the underlying problem
- **THEN** the orchestrator SHALL ask a clarifying question about the problem

#### Scenario: Overlap with existing spec detected

- **WHEN** the aggregated content overlaps with an existing spec (e.g., `planning`)
- **THEN** the orchestrator SHALL ask whether to extend the existing spec or create a new capability


<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
-->

### Requirement: Interview asks one question at a time

The `mysd:propose` orchestrator SHALL ask at most one clarifying question per turn during the requirement interview.

The orchestrator SHALL NOT present a list of multiple questions. Each question SHALL target the single most important gap in the current understanding.

#### Scenario: Two dimensions are incomplete

- **WHEN** both Problem and Boundary dimensions are incomplete
- **THEN** the orchestrator SHALL ask about Problem first (higher priority), then Boundary in a subsequent turn


<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
-->

### Requirement: Interview question count is dynamic

The number of interview questions SHALL NOT be hardcoded. The orchestrator SHALL determine the number of questions based on the completeness evaluation:

- 0 questions if all dimensions are sufficiently covered
- 1+ questions if gaps exist, one per turn until all dimensions are addressed

The orchestrator SHALL re-evaluate completeness after each user answer to determine if further questions are needed.

#### Scenario: User's first answer covers remaining gaps

- **WHEN** the user's answer to the first question also addresses the remaining incomplete dimensions
- **THEN** the orchestrator SHALL recognize this and proceed to requirement_brief generation without further questions


<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
-->

### Requirement: Interview produces structured requirement_brief

After the interview completes (0 or more questions), the orchestrator SHALL produce a `requirement_brief` with the following structure:

```
## Problem
{The problem being solved}

## Boundary
{What is in scope / what is explicitly excluded}

## Success Criteria
{Specific, verifiable conditions for success}

## Source
{List of source types used, for traceability}
```

The `requirement_brief` SHALL be passed as input to subsequent steps (research, proposal-writer) but SHALL NOT be written to disk as a separate file.

#### Scenario: requirement_brief from full interview

- **WHEN** the orchestrator asked 2 clarifying questions and received answers
- **THEN** the requirement_brief SHALL incorporate both the original aggregated_content and the interview answers

#### Scenario: requirement_brief from sufficient content

- **WHEN** the aggregated content was sufficient (0 questions asked)
- **THEN** the requirement_brief SHALL be synthesized directly from the aggregated_content


<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
-->

### Requirement: Auto mode skips interview with best-effort brief

When `auto_mode` is true, the orchestrator SHALL skip the interview entirely and produce a `requirement_brief` using best-effort inference from the `aggregated_content`.

Dimensions that cannot be inferred SHALL be filled with reasonable defaults based on the available context, not left empty or marked as incomplete.

#### Scenario: Auto mode produces complete brief

- **WHEN** `auto_mode` is true AND aggregated content covers Problem but not Boundary
- **THEN** the requirement_brief SHALL infer Boundary from the problem context and source material

<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
tests: []
-->

## Requirements


<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
-->

### Requirement: Orchestrator evaluates requirement completeness

After material selection, the `mysd:propose` orchestrator SHALL evaluate the `aggregated_content` against three completeness dimensions:

| Dimension | Complete when |
|-----------|--------------|
| Problem | The content describes the problem being solved, not just a desired solution |
| Boundary | The content states what is in scope and what is explicitly excluded |
| Success Criteria | The content includes specific, verifiable conditions for success |

The orchestrator SHALL also compare `aggregated_content` against related existing specs (from the spec scan step). If overlap is detected between the proposed change and an existing spec, the orchestrator SHALL flag it as an additional interview question.

#### Scenario: All dimensions are covered

- **WHEN** the aggregated content clearly addresses problem, boundary, and success criteria
- **THEN** the orchestrator SHALL skip the interview and produce the requirement_brief directly

#### Scenario: Problem dimension is missing

- **WHEN** the aggregated content describes a desired solution but not the underlying problem
- **THEN** the orchestrator SHALL ask a clarifying question about the problem

#### Scenario: Overlap with existing spec detected

- **WHEN** the aggregated content overlaps with an existing spec (e.g., `planning`)
- **THEN** the orchestrator SHALL ask whether to extend the existing spec or create a new capability

---
### Requirement: Interview asks one question at a time

The `mysd:propose` orchestrator SHALL ask at most one clarifying question per turn during the requirement interview.

The orchestrator SHALL NOT present a list of multiple questions. Each question SHALL target the single most important gap in the current understanding.

#### Scenario: Two dimensions are incomplete

- **WHEN** both Problem and Boundary dimensions are incomplete
- **THEN** the orchestrator SHALL ask about Problem first (higher priority), then Boundary in a subsequent turn

---
### Requirement: Interview question count is dynamic

The number of interview questions SHALL NOT be hardcoded. The orchestrator SHALL determine the number of questions based on the completeness evaluation:

- 0 questions if all dimensions are sufficiently covered
- 1+ questions if gaps exist, one per turn until all dimensions are addressed

The orchestrator SHALL re-evaluate completeness after each user answer to determine if further questions are needed.

#### Scenario: User's first answer covers remaining gaps

- **WHEN** the user's answer to the first question also addresses the remaining incomplete dimensions
- **THEN** the orchestrator SHALL recognize this and proceed to requirement_brief generation without further questions

---
### Requirement: Interview produces structured requirement_brief

After the interview completes (0 or more questions), the orchestrator SHALL produce a `requirement_brief` with the following structure:

```
## Problem
{The problem being solved}

## Boundary
{What is in scope / what is explicitly excluded}

## Success Criteria
{Specific, verifiable conditions for success}

## Source
{List of source types used, for traceability}
```

The `requirement_brief` SHALL be passed as input to subsequent steps (research, proposal-writer) but SHALL NOT be written to disk as a separate file.

#### Scenario: requirement_brief from full interview

- **WHEN** the orchestrator asked 2 clarifying questions and received answers
- **THEN** the requirement_brief SHALL incorporate both the original aggregated_content and the interview answers

#### Scenario: requirement_brief from sufficient content

- **WHEN** the aggregated content was sufficient (0 questions asked)
- **THEN** the requirement_brief SHALL be synthesized directly from the aggregated_content

---
### Requirement: Auto mode skips interview with best-effort brief

When `auto_mode` is true, the orchestrator SHALL skip the interview entirely and produce a `requirement_brief` using best-effort inference from the `aggregated_content`.

Dimensions that cannot be inferred SHALL be filled with reasonable defaults based on the available context, not left empty or marked as TBD.

#### Scenario: Auto mode produces complete brief

- **WHEN** `auto_mode` is true AND aggregated content covers Problem but not Boundary
- **THEN** the requirement_brief SHALL infer Boundary from the problem context and source material