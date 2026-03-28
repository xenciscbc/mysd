---
spec-version: "1.0"
capability: Task Planning & Coverage Validation
delta: ADDED
status: done
---

## Requirement: Task Planning

The `mysd plan` command MUST break a design into executable tasks stored in `tasks.md`.

Tasks frontmatter MUST include: `spec-version`, `change`, `status`.

Each task MUST have: ID (T{n}), description, status, and `satisfies` field mapping to requirement IDs.

The `--research` flag MUST spawn a researcher agent before planning.

The `--context-only` flag MUST output planning context as JSON without writing files.

## Requirement: Plan Coverage Checking

The `planchecker` package MUST validate that all MUST-level requirements are covered by at least one task.

`CheckCoverage()` MUST return a `CoverageResult` with: TotalMust, CoveredCount, UncoveredIDs, CoverageRatio, Passed.

The `--check` flag on `mysd plan` MUST invoke the plan checker after task generation.

Coverage MUST be computed by mapping each task's `satisfies[]` field to MUST requirement IDs.

## Requirement: Wave Grouping

The `executor` package MUST compute parallel execution waves via `ComputeWaves()`.

Tasks with no inter-dependencies MUST be grouped into the same wave for parallel execution.

`BuildAlignmentReport()` MUST validate task dependencies and ordering before execution.

### Scenario: Full Coverage

WHEN all MUST requirements have at least one task with a matching `satisfies` entry
THEN CheckCoverage() returns Passed=true

### Scenario: Missing Coverage

WHEN a MUST requirement has no matching task
THEN CheckCoverage() returns Passed=false with the uncovered ID in UncoveredIDs

## Requirement: Plan pipeline includes reviewer step after planner

The `mysd:plan` skill SHALL invoke `mysd-reviewer` (phase="plan") after the planner completes and before `mysd-plan-checker` runs.

The reviewer step (Step 5b) SHALL:
1. Run `mysd validate {change_name}` and capture output as `validate_output` (empty string if command not found or fails)
2. Display: `Spawning mysd-reviewer ({reviewer_model})...`
3. Spawn `mysd-reviewer` via the Task tool with `model` set to `{reviewer_model}`

The reviewer SHALL be invoked with:
```json
{
  "change_name": "{change_name}",
  "phase": "plan",
  "validate_output": "{validate_output}",
  "auto_mode": {auto_mode}
}
```

### Scenario: Reviewer is invoked after planner in plan pipeline

- **WHEN** `mysd:plan` completes the planner step (Step 5)
- **THEN** it SHALL run `mysd validate` and spawn `mysd-reviewer` with `phase: "plan"` before proceeding to plan-checker

### Scenario: Reviewer uses reviewer_model from plan context

- **WHEN** `mysd:plan` spawns `mysd-reviewer`
- **THEN** it SHALL use `{reviewer_model}` from the plan context JSON (not the shared `{model}`)

## Requirement: Plan pipeline uses per-role models for reviewer and plan-checker

The `mysd:plan` skill SHALL read `reviewer_model` and `plan_checker_model` from the plan context JSON (Step 2) and use them when spawning the respective agents.

`mysd-plan-checker` SHALL be spawned with `model: {plan_checker_model}` (not the shared `{model}`).

### Scenario: Plan-checker uses plan_checker_model

- **WHEN** `mysd:plan` spawns `mysd-plan-checker` in Step 6
- **THEN** it SHALL use `{plan_checker_model}` from the plan context JSON

## Requirement: Discuss skill enforces discussion quality guidelines

The `mysd:discuss` skill SHALL follow discussion quality guidelines throughout the conversation loop:

- Ask one question at a time — the most important unresolved question first
- Present 2–3 concrete options with trade-offs when exploring approaches
- Prohibit empty validation phrases ("That's interesting", "There are many ways", "That could work")
- Provide a direct recommendation when one exists

When the user signals impatience ("let's just go with X", "move on"):
- First occurrence: flag one important unresolved question in one sentence, then offer to continue
- Second occurrence: respect user's pace, skip to convergence without further pushback

### Scenario: Single question per turn

- **WHEN** the discuss skill has multiple unresolved questions
- **THEN** it SHALL ask only the single most important question in each turn

### Scenario: User impatience is respected on second occurrence

- **WHEN** the user pushes to move forward a second time
- **THEN** the discuss skill SHALL proceed directly to convergence without additional pushback

## Requirement: Discuss skill enforces convergence and conclusion capture

The `mysd:discuss` skill SHALL converge to an explicit conclusion before the discussion ends.

When a clear conclusion is reached, the skill SHALL proactively present a conclusion summary without waiting for the user to ask:

```
## Conclusion

**Decision**: [What was decided]
**Rationale**: [The key trade-off that drove this]
**Capture to**: [Which artifact: proposal.md / spec / design.md / tasks.md]
```

The skill SHALL say: "I'll capture this to {artifact} unless you'd rather not."

If the user attempts to end without a conclusion, the skill SHALL summarize the current state and state what remains unresolved. The discussion SHALL NOT end without at least an explicit deferral.

### Scenario: Proactive conclusion summary is presented

- **WHEN** the discuss skill determines a clear conclusion has been reached
- **THEN** it SHALL present the conclusion summary format before the user asks
- **AND** offer to capture the conclusion to the appropriate artifact

### Scenario: No-conclusion ending is blocked

- **WHEN** the user tries to end the discussion without a conclusion
- **THEN** the discuss skill SHALL summarize what was discussed and state what is unresolved
- **AND** SHALL NOT end the discussion without at least an explicit deferral statement

## Covered Packages

- `cmd/plan.go`
- `internal/planchecker/` — MUST coverage validation
- `internal/executor/` — wave computation and alignment (shared with execution spec)

## Requirements

### Requirement: Plan pipeline includes reviewer step after planner

The `mysd:plan` skill SHALL invoke `mysd-reviewer` (phase="plan") after the planner completes and before `mysd-plan-checker` runs.

The reviewer step (Step 5b) SHALL:
1. Run `mysd validate {change_name}` and capture output as `validate_output` (empty string if command not found or fails)
2. Display: `Spawning mysd-reviewer ({reviewer_model})...`
3. Spawn `mysd-reviewer` via the Task tool with `model` set to `{reviewer_model}`

The reviewer SHALL be invoked with:
```json
{
  "change_name": "{change_name}",
  "phase": "plan",
  "validate_output": "{validate_output}",
  "auto_mode": {auto_mode}
}
```

#### Scenario: Reviewer is invoked after planner in plan pipeline

- **WHEN** `mysd:plan` completes the planner step (Step 5)
- **THEN** it SHALL run `mysd validate` and spawn `mysd-reviewer` with `phase: "plan"` before proceeding to plan-checker

#### Scenario: Reviewer uses reviewer_model from plan context

- **WHEN** `mysd:plan` spawns `mysd-reviewer`
- **THEN** it SHALL use `{reviewer_model}` from the plan context JSON (not the shared `{model}`)


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
### Requirement: Plan pipeline uses per-role models for reviewer and plan-checker

The `mysd:plan` skill SHALL read `reviewer_model` and `plan_checker_model` from the plan context JSON (Step 2) and use them when spawning the respective agents.

`mysd-plan-checker` SHALL be spawned with `model: {plan_checker_model}` (not the shared `{model}`).

#### Scenario: Plan-checker uses plan_checker_model

- **WHEN** `mysd:plan` spawns `mysd-plan-checker` in Step 6
- **THEN** it SHALL use `{plan_checker_model}` from the plan context JSON


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
### Requirement: Discuss skill enforces discussion quality guidelines

The `mysd:discuss` skill SHALL follow discussion quality guidelines throughout the conversation loop:

- Ask one question at a time — the most important unresolved question first
- Present 2–3 concrete options with trade-offs when exploring approaches
- Prohibit empty validation phrases ("That's interesting", "There are many ways", "That could work")
- Provide a direct recommendation when one exists

When the user signals impatience ("let's just go with X", "move on"):
- First occurrence: flag one important unresolved question in one sentence, then offer to continue
- Second occurrence: respect user's pace, skip to convergence without further pushback

#### Scenario: Single question per turn

- **WHEN** the discuss skill has multiple unresolved questions
- **THEN** it SHALL ask only the single most important question in each turn

#### Scenario: User impatience is respected on second occurrence

- **WHEN** the user pushes to move forward a second time
- **THEN** the discuss skill SHALL proceed directly to convergence without additional pushback


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
### Requirement: Discuss skill enforces convergence and conclusion capture

The `mysd:discuss` skill SHALL converge to an explicit conclusion before the discussion ends.

When a clear conclusion is reached, the skill SHALL proactively present a conclusion summary without waiting for the user to ask:

```
## Conclusion

**Decision**: [What was decided]
**Rationale**: [The key trade-off that drove this]
**Capture to**: [Which artifact: proposal.md / spec / design.md / tasks.md]
```

The skill SHALL say: "I'll capture this to {artifact} unless you'd rather not."

If the user attempts to end without a conclusion, the skill SHALL summarize the current state and state what remains unresolved. The discussion SHALL NOT end without at least an explicit deferral.

#### Scenario: Proactive conclusion summary is presented

- **WHEN** the discuss skill determines a clear conclusion has been reached
- **THEN** it SHALL present the conclusion summary format before the user asks
- **AND** offer to capture the conclusion to the appropriate artifact

#### Scenario: No-conclusion ending is blocked

- **WHEN** the user tries to end the discussion without a conclusion
- **THEN** the discuss skill SHALL summarize what was discussed and state what is unresolved
- **AND** SHALL NOT end the discussion without at least an explicit deferral statement

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