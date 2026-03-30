---
spec-version: "1.0"
capability: Task Planning & Coverage Validation
delta: ADDED
status: draft
---

## ADDED Requirements

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
