---
spec-version: "1.0"
capability: Discuss Path Routing
status: draft
---

### Requirement: Discuss skill presents path selection at entry

The `mysd:discuss` skill SHALL present the user with a path selection prompt after determining context:

- **When an active change exists** (from argument or `mysd status`): the skill SHALL ask the user to choose between:
  1. Discuss an existing spec within the change (spec-focused path)
  2. Add content from external sources (source-driven path)

- **When no active change exists**: the skill SHALL enter the source-driven path directly.

In `auto_mode`, the skill SHALL default to the source-driven path without prompting.

#### Scenario: Active change with path selection

- **WHEN** the user invokes `mysd:discuss` with an active change
- **THEN** the skill SHALL present the two path options and wait for user selection

#### Scenario: No active change defaults to source-driven

- **WHEN** the user invokes `mysd:discuss` with no active change and no change argument
- **THEN** the skill SHALL enter the source-driven path directly

#### Scenario: Auto mode defaults to source-driven

- **WHEN** `auto_mode` is true and an active change exists
- **THEN** the skill SHALL enter the source-driven path without prompting

<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
tests: []
-->


<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
-->

---
### Requirement: Spec-focused path lists specs and accepts user selection

When the user selects the spec-focused path, the `mysd:discuss` skill SHALL:

1. List all spec files under `.specs/changes/{change_name}/specs/*/spec.md` with their capability name and requirement count
2. Allow the user to select one spec for focused discussion
3. After selection, proceed to gap analysis (see `spec-gap-analysis` capability)

If the change has no specs, the skill SHALL inform the user and fall back to the source-driven path.

#### Scenario: Multiple specs available for selection

- **WHEN** the change has 3 specs under `specs/`
- **THEN** the skill SHALL list all 3 with capability name and requirement count, and wait for user selection

#### Scenario: Change has no specs

- **WHEN** the change has no spec files under `specs/`
- **THEN** the skill SHALL display a message and fall back to the source-driven path

<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
tests: []
-->


<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
-->

---
### Requirement: Source-driven path uses material selection and recommends spec target

When the user selects the source-driven path (or enters it by default), the `mysd:discuss` skill SHALL:

1. Execute the same 6-source detection and material selection process as `mysd:propose` (see `material-selection` spec)
2. After content aggregation, compare the aggregated content against existing specs under `.specs/changes/{change_name}/specs/` (if a change exists) and `openspec/specs/`
3. Present a recommendation:
   - **Clear match**: recommend merging into the matching spec, naming the spec file
   - **Multiple matches**: list related specs and let the user choose which to extend
   - **No match**: recommend creating a new spec with a suggested kebab-case name
4. Wait for user confirmation before proceeding to discussion

When no active change exists, the skill SHALL:
1. Complete material selection and content aggregation
2. Derive a change name from the aggregated content
3. Run `mysd propose {name}` to scaffold a new change
4. Then execute the spec comparison and recommendation step above

#### Scenario: Aggregated content matches one existing spec

- **WHEN** the aggregated content is highly related to `specs/material-selection/spec.md`
- **THEN** the skill SHALL recommend merging into that spec

#### Scenario: Aggregated content matches no existing spec

- **WHEN** the aggregated content does not relate to any existing spec
- **THEN** the skill SHALL recommend creating a new spec with a suggested name

#### Scenario: No active change triggers automatic scaffold

- **WHEN** the source-driven path is entered with no active change
- **AND** the user confirms material selection
- **THEN** the skill SHALL derive a change name and run `mysd propose {name}` to scaffold

<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
tests: []
-->


<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
-->

---
### Requirement: Discussion loop merges research conclusion and exit into unified flow

After research steps (Steps 6-8) complete, the `mysd:discuss` skill SHALL present a unified summary of all gray area conclusions from Step 7-8, then offer:

1. Continue discussing other aspects
2. Converge to conclusion and decide on spec updates

This unified exit SHALL replace the separate Layer 2 discovery prompt (Step 8) and the research-recap in Step 9. There SHALL be only one exit point from research into discussion, not two sequential prompts.

When no research was performed, the skill SHALL enter the discussion loop directly using the context from path routing (gap analysis results for spec-focused path, or aggregated content for source-driven path).

#### Scenario: Research completed with unified exit

- **WHEN** research Steps 6-8 complete with 2 gray area conclusions
- **THEN** the skill SHALL present both conclusions in a single summary and offer the two options

#### Scenario: No research enters discussion with path context

- **WHEN** research is skipped
- **AND** the user entered via spec-focused path with gap analysis results
- **THEN** the discussion loop SHALL use the gap analysis results as the discussion starting point

<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
tests: []
-->

## Requirements


<!-- @trace
source: enhance-discuss-workflow
updated: 2026-03-30
code:
  - mysd/skills/discuss/SKILL.md
-->

### Requirement: Discuss skill presents path selection at entry

The `mysd:discuss` skill SHALL present the user with a path selection prompt after determining context:

- **When an active change exists** (from argument or `mysd status`): the skill SHALL ask the user to choose between:
  1. Discuss an existing spec within the change (spec-focused path)
  2. Add content from external sources (source-driven path)

- **When no active change exists**: the skill SHALL enter the source-driven path directly.

In `auto_mode`, the skill SHALL default to the source-driven path without prompting.

#### Scenario: Active change with path selection

- **WHEN** the user invokes `mysd:discuss` with an active change
- **THEN** the skill SHALL present the two path options and wait for user selection

#### Scenario: No active change defaults to source-driven

- **WHEN** the user invokes `mysd:discuss` with no active change and no change argument
- **THEN** the skill SHALL enter the source-driven path directly

#### Scenario: Auto mode defaults to source-driven

- **WHEN** `auto_mode` is true and an active change exists
- **THEN** the skill SHALL enter the source-driven path without prompting

---
### Requirement: Spec-focused path lists specs and accepts user selection

When the user selects the spec-focused path, the `mysd:discuss` skill SHALL:

1. List all spec files under `.specs/changes/{change_name}/specs/*/spec.md` with their capability name and requirement count
2. Allow the user to select one spec for focused discussion
3. After selection, proceed to gap analysis (see `spec-gap-analysis` capability)

If the change has no specs, the skill SHALL inform the user and fall back to the source-driven path.

#### Scenario: Multiple specs available for selection

- **WHEN** the change has 3 specs under `specs/`
- **THEN** the skill SHALL list all 3 with capability name and requirement count, and wait for user selection

#### Scenario: Change has no specs

- **WHEN** the change has no spec files under `specs/`
- **THEN** the skill SHALL display a message and fall back to the source-driven path

---
### Requirement: Source-driven path uses material selection and recommends spec target

When the user selects the source-driven path (or enters it by default), the `mysd:discuss` skill SHALL:

1. Execute the same 6-source detection and material selection process as `mysd:propose` (see `material-selection` spec)
2. After content aggregation, compare the aggregated content against existing specs under `.specs/changes/{change_name}/specs/` (if a change exists) and `openspec/specs/`
3. Present a recommendation:
   - **Clear match**: recommend merging into the matching spec, naming the spec file
   - **Multiple matches**: list related specs and let the user choose which to extend
   - **No match**: recommend creating a new spec with a suggested kebab-case name
4. Wait for user confirmation before proceeding to discussion

When no active change exists, the skill SHALL:
1. Complete material selection and content aggregation
2. Derive a change name from the aggregated content
3. Run `mysd propose {name}` to scaffold a new change
4. Then execute the spec comparison and recommendation step above

#### Scenario: Aggregated content matches one existing spec

- **WHEN** the aggregated content is highly related to `specs/material-selection/spec.md`
- **THEN** the skill SHALL recommend merging into that spec

#### Scenario: Aggregated content matches no existing spec

- **WHEN** the aggregated content does not relate to any existing spec
- **THEN** the skill SHALL recommend creating a new spec with a suggested name

#### Scenario: No active change triggers automatic scaffold

- **WHEN** the source-driven path is entered with no active change
- **AND** the user confirms material selection
- **THEN** the skill SHALL derive a change name and run `mysd propose {name}` to scaffold

---
### Requirement: Discussion loop merges research conclusion and exit into unified flow

After research steps (Steps 6-8) complete, the `mysd:discuss` skill SHALL present a unified summary of all gray area conclusions from Step 7-8, then offer:

1. Continue discussing other aspects
2. Converge to conclusion and decide on spec updates

This unified exit SHALL replace the separate Layer 2 discovery prompt (Step 8) and the research-recap in Step 9. There SHALL be only one exit point from research into discussion, not two sequential prompts.

When no research was performed, the skill SHALL enter the discussion loop directly using the context from path routing (gap analysis results for spec-focused path, or aggregated content for source-driven path).

#### Scenario: Research completed with unified exit

- **WHEN** research Steps 6-8 complete with 2 gray area conclusions
- **THEN** the skill SHALL present both conclusions in a single summary and offer the two options

#### Scenario: No research enters discussion with path context

- **WHEN** research is skipped
- **AND** the user entered via spec-focused path with gap analysis results
- **THEN** the discussion loop SHALL use the gap analysis results as the discussion starting point