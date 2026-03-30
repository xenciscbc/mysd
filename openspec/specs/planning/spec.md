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

## Requirement: Discuss skill re-plan is conditional on existing plan

The `mysd:discuss` skill SHALL execute the re-plan and plan-checker steps (Step 11) only when a plan already exists for the current change.

The skill SHALL determine plan existence by checking whether `.specs/changes/{change_name}/tasks.md` exists.

- If `tasks.md` exists: execute re-plan (run `mysd plan --context-only`, spawn planner, run `mysd plan`, run plan-checker)
- If `tasks.md` does not exist: skip Step 11 entirely and proceed to Step 12 (Confirm)

### Scenario: Re-plan executes when tasks.md exists

- **WHEN** the discuss skill reaches Step 11
- **AND** `.specs/changes/{change_name}/tasks.md` exists
- **THEN** the skill SHALL execute the full re-plan and plan-checker sequence

### Scenario: Re-plan skipped when no tasks.md

- **WHEN** the discuss skill reaches Step 11
- **AND** `.specs/changes/{change_name}/tasks.md` does not exist
- **THEN** the skill SHALL skip Step 11 and proceed directly to Step 12

<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
tests: []
-->

---
## Requirement: Discuss skill spec update confirmation

The `mysd:discuss` skill SHALL present a confirmation list before executing spec updates (Step 10).

The list SHALL:
- Include only artifacts that are affected by the discussion conclusions (do not list unaffected artifacts)
- Default all items to selected (checked)
- Allow the user to deselect individual items

After confirmation, the skill SHALL execute updates only for the items that remain selected.

In `auto_mode`, the skill SHALL execute all affected updates without presenting the confirmation list.

### Scenario: User confirms all updates

- **WHEN** the discuss skill presents 2 affected specs
- **AND** the user confirms without changes
- **THEN** the skill SHALL update both specs

### Scenario: User deselects one update

- **WHEN** the discuss skill presents 2 affected specs
- **AND** the user deselects 1 spec
- **THEN** the skill SHALL update only the remaining selected spec

### Scenario: Auto mode skips confirmation

- **WHEN** `auto_mode` is true
- **AND** there are 3 affected artifacts
- **THEN** the skill SHALL update all 3 without presenting the confirmation list

<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
tests: []
-->

---
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

---
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


<!-- @trace
source: enhance-workflow-quality
updated: 2026-03-30
code:
  - internal/analyzer/ambiguity.go
  - cmd/analyze.go
  - internal/analyzer/consistency.go
  - mysd/skills/propose/SKILL.md
  - internal/analyzer/analyzer.go
  - mysd/skills/plan/SKILL.md
  - mysd/agents/mysd-reviewer.md
  - internal/analyzer/types.go
  - internal/analyzer/gaps.go
  - mysd/agents/mysd-proposal-writer.md
  - internal/analyzer/coverage.go
tests:
  - internal/analyzer/analyzer_test.go
  - cmd/analyze_test.go
-->

---
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

---
### Requirement: Propose workflow step ordering

The `mysd:propose` skill SHALL execute steps in the following order:

1. **Parse Arguments**: Parse `--auto` flag and `source_arg`
2. **Resolve Agent Model**: Run `mysd model` to determine per-role models
3. **Load Deferred Notes**: Run `mysd note list` to load cross-change context
4. **Material Selection**: Detect all available requirement sources, present to user for selection, aggregate selected content (see `material-selection` spec)
5. **Scan Existing Specs**: Scan `openspec/specs/*/spec.md` for related specs, retain content for interview step
6. **Requirement Interview**: Evaluate aggregated content completeness, ask clarifying questions as needed, produce structured `requirement_brief` (see `requirement-interview` spec)
7. **Derive Change Name + Classify Type**: Derive kebab-case change name from `requirement_brief`, classify as feature/bugfix/refactor
8. **Scaffold the Change**: Run `mysd propose {change-name}` to create directory structure
9. **Optional 4-Dimension Research**: Ask user whether to run research; if accepted, spawn 4 `mysd-researcher` agents in parallel
10. **Gray Area + Advisor** (research only): Identify gray areas from research output, spawn `mysd-advisor` per area
11. **Dual-Loop Exploration** (research only): Deep dive into gray areas with scope guardrail
12. **Invoke Proposal Writer**: Spawn `mysd-proposal-writer` with `requirement_brief` and research/exploration conclusions
13. **Auto-Invoke Spec Writer**: Spawn `mysd-spec-writer` per capability area
14. **Artifact Review**: Run `mysd validate` then spawn `mysd-reviewer`
15. **Final Summary**: Display results and next steps

Steps 10-11 SHALL be executed only when the user accepts 4-Dimension Research in Step 9. When research is declined, the workflow SHALL proceed directly from Step 9 to Step 12.

#### Scenario: Full workflow with research

- **WHEN** the user accepts 4-Dimension Research
- **THEN** the skill SHALL execute all 15 steps in order

#### Scenario: Workflow without research

- **WHEN** the user declines 4-Dimension Research
- **THEN** the skill SHALL skip Steps 10-11 and proceed from Step 9 directly to Step 12


<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
-->

---
### Requirement: Change name derived after requirement interview

The `mysd:propose` skill SHALL derive the change name from the completed `requirement_brief` (Step 7), not from the initial `source_arg` or raw source content.

If `source_arg` refers to an existing change directory (`.specs/changes/{source_arg}/`), the skill SHALL use that name directly without re-derivation.

#### Scenario: Change name from requirement_brief

- **WHEN** the user provides a vague description "improve the propose flow"
- **AND** the interview clarifies the scope to material selection and requirement interview
- **THEN** the derived change name SHALL reflect the clarified scope (e.g., `enhance-propose-workflow`)

#### Scenario: Existing change name preserved

- **WHEN** `source_arg` matches an existing change directory
- **THEN** the skill SHALL use the existing change name without re-derivation


<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
-->

---
### Requirement: Existing spec content fed into interview

The `mysd:propose` skill SHALL pass related existing spec content (from Step 5 scan) into the requirement interview step (Step 6).

The orchestrator SHALL use this spec content to detect overlap between the proposed change and existing capabilities, and SHALL ask the user whether to extend the existing spec or create a new capability when overlap is detected.

#### Scenario: Overlap detected with existing spec

- **WHEN** the proposed change overlaps with the existing `planning` spec
- **THEN** the orchestrator SHALL ask the user: extend the existing `planning` spec or create a new capability

#### Scenario: No overlap with existing specs

- **WHEN** no existing specs are related to the proposed change
- **THEN** the orchestrator SHALL proceed with the interview without spec-related questions

---
### Requirement: Discuss skill re-plan is conditional on existing plan

The `mysd:discuss` skill SHALL execute the re-plan and plan-checker steps (Step 11) only when a plan already exists for the current change.

The skill SHALL determine plan existence by checking whether `.specs/changes/{change_name}/tasks.md` exists.

- If `tasks.md` exists: execute re-plan (run `mysd plan --context-only`, spawn planner, run `mysd plan`, run plan-checker)
- If `tasks.md` does not exist: skip Step 11 entirely and proceed to Step 12 (Confirm)

#### Scenario: Re-plan executes when tasks.md exists

- **WHEN** the discuss skill reaches Step 11
- **AND** `.specs/changes/{change_name}/tasks.md` exists
- **THEN** the skill SHALL execute the full re-plan and plan-checker sequence

#### Scenario: Re-plan skipped when no tasks.md

- **WHEN** the discuss skill reaches Step 11
- **AND** `.specs/changes/{change_name}/tasks.md` does not exist
- **THEN** the skill SHALL skip Step 11 and proceed directly to Step 12


<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
-->

---
### Requirement: Discuss skill spec update confirmation

The `mysd:discuss` skill SHALL present a confirmation list before executing spec updates (Step 10).

The list SHALL:
- Include only artifacts that are affected by the discussion conclusions (do not list unaffected artifacts)
- Default all items to selected (checked)
- Allow the user to deselect individual items

After confirmation, the skill SHALL execute updates only for the items that remain selected.

In `auto_mode`, the skill SHALL execute all affected updates without presenting the confirmation list.

#### Scenario: User confirms all updates

- **WHEN** the discuss skill presents 2 affected specs
- **AND** the user confirms without changes
- **THEN** the skill SHALL update both specs

#### Scenario: User deselects one update

- **WHEN** the discuss skill presents 2 affected specs
- **AND** the user deselects 1 spec
- **THEN** the skill SHALL update only the remaining selected spec

#### Scenario: Auto mode skips confirmation

- **WHEN** `auto_mode` is true
- **AND** there are 3 affected artifacts
- **THEN** the skill SHALL update all 3 without presenting the confirmation list

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


<!-- @trace
source: enhance-workflow-quality
updated: 2026-03-30
code:
  - internal/analyzer/ambiguity.go
  - cmd/analyze.go
  - internal/analyzer/consistency.go
  - mysd/skills/propose/SKILL.md
  - internal/analyzer/analyzer.go
  - mysd/skills/plan/SKILL.md
  - mysd/agents/mysd-reviewer.md
  - internal/analyzer/types.go
  - internal/analyzer/gaps.go
  - mysd/agents/mysd-proposal-writer.md
  - internal/analyzer/coverage.go
tests:
  - internal/analyzer/analyzer_test.go
  - cmd/analyze_test.go
-->


<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
-->


<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
-->

---
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

<!-- @trace
source: enhance-workflow-quality
updated: 2026-03-30
code: []
tests: []
-->

---
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

<!-- @trace
source: enhance-workflow-quality
updated: 2026-03-30
code: []
tests: []
-->

---
### Requirement: Propose workflow step ordering

The `mysd:propose` skill SHALL execute steps in the following order:

1. **Parse Arguments**: Parse `--auto` flag and `source_arg`
2. **Resolve Agent Model**: Run `mysd model` to determine per-role models
3. **Load Deferred Notes**: Run `mysd note list` to load cross-change context
4. **Material Selection**: Detect all available requirement sources, present to user for selection, aggregate selected content (see `material-selection` spec)
5. **Scan Existing Specs**: Scan `openspec/specs/*/spec.md` for related specs, retain content for interview step
6. **Requirement Interview**: Evaluate aggregated content completeness, ask clarifying questions as needed, produce structured `requirement_brief` (see `requirement-interview` spec)
7. **Derive Change Name + Classify Type**: Derive kebab-case change name from `requirement_brief`, classify as feature/bugfix/refactor
8. **Scaffold the Change**: Run `mysd propose {change-name}` to create directory structure
9. **Optional 4-Dimension Research**: Ask user whether to run research; if accepted, spawn 4 `mysd-researcher` agents in parallel
10. **Gray Area + Advisor** (research only): Identify gray areas from research output, spawn `mysd-advisor` per area
11. **Dual-Loop Exploration** (research only): Deep dive into gray areas with scope guardrail
12. **Invoke Proposal Writer**: Spawn `mysd-proposal-writer` with `requirement_brief` and research/exploration conclusions
13. **Auto-Invoke Spec Writer**: Spawn `mysd-spec-writer` per capability area
14. **Artifact Review**: Run `mysd validate` then spawn `mysd-reviewer`
15. **Final Summary**: Display results and next steps

Steps 10-11 SHALL be executed only when the user accepts 4-Dimension Research in Step 9. When research is declined, the workflow SHALL proceed directly from Step 9 to Step 12.

#### Scenario: Full workflow with research

- **WHEN** the user accepts 4-Dimension Research
- **THEN** the skill SHALL execute all 15 steps in order

#### Scenario: Workflow without research

- **WHEN** the user declines 4-Dimension Research
- **THEN** the skill SHALL skip Steps 10-11 and proceed from Step 9 directly to Step 12

### Requirement: Change name derived after requirement interview

The `mysd:propose` skill SHALL derive the change name from the completed `requirement_brief` (Step 7), not from the initial `source_arg` or raw source content.

If `source_arg` refers to an existing change directory (`.specs/changes/{source_arg}/`), the skill SHALL use that name directly without re-derivation.

#### Scenario: Change name from requirement_brief

- **WHEN** the user provides a vague description "improve the propose flow"
- **AND** the interview clarifies the scope to material selection and requirement interview
- **THEN** the derived change name SHALL reflect the clarified scope (e.g., `enhance-propose-workflow`)

#### Scenario: Existing change name preserved

- **WHEN** `source_arg` matches an existing change directory
- **THEN** the skill SHALL use the existing change name without re-derivation

### Requirement: Existing spec content fed into interview

The `mysd:propose` skill SHALL pass related existing spec content (from Step 5 scan) into the requirement interview step (Step 6).

The orchestrator SHALL use this spec content to detect overlap between the proposed change and existing capabilities, and SHALL ask the user whether to extend the existing spec or create a new capability when overlap is detected.

#### Scenario: Overlap detected with existing spec

- **WHEN** the proposed change overlaps with the existing `planning` spec
- **THEN** the orchestrator SHALL ask the user: extend the existing `planning` spec or create a new capability

#### Scenario: No overlap with existing specs

- **WHEN** no existing specs are related to the proposed change
- **THEN** the orchestrator SHALL proceed with the interview without spec-related questions

<!-- @trace
source: enhance-propose-workflow
updated: 2026-03-30
code:
  - mysd/skills/propose/SKILL.md
tests: []
-->