## ADDED Requirements

### Requirement: DefaultModelMap includes reviewer role

The `DefaultModelMap` in `internal/config/config.go` SHALL include a `reviewer` role with the following assignments:

| Profile | Model |
|---------|-------|
| quality | opus |
| balanced | sonnet |
| budget | sonnet |

Rationale: reviewer is a judgment role requiring reasoning capability. Budget profile uses `sonnet` (not `haiku`) because quality assessment requires sufficient model capability.

#### Scenario: Quality profile returns opus for reviewer

- **WHEN** `ResolveModel("reviewer", "quality", nil)` is called
- **THEN** the return value SHALL be `"opus"`

#### Scenario: Balanced profile returns sonnet for reviewer

- **WHEN** `ResolveModel("reviewer", "balanced", nil)` is called
- **THEN** the return value SHALL be `"sonnet"`

#### Scenario: Budget profile returns sonnet for reviewer

- **WHEN** `ResolveModel("reviewer", "budget", nil)` is called
- **THEN** the return value SHALL be `"sonnet"`

## MODIFIED Requirements

### Requirement: Binary context JSON includes model field

The `--context-only` output of workflow commands (`plan`, `execute`, `design`, `spec`) SHALL include a `"model"` field containing the resolved short name for the relevant agent role.

For the `plan` command specifically, the `--context-only` output SHALL additionally include:
- `"reviewer_model"`: resolved model for the `reviewer` role under the current profile
- `"plan_checker_model"`: resolved model for the `plan-checker` role under the current profile

#### Scenario: Plan context includes planner model

- **WHEN** `mysd plan --context-only` is executed with profile set to `balanced`
- **THEN** the JSON output SHALL contain `"model": "sonnet"`

#### Scenario: Plan context includes reviewer_model

- **WHEN** `mysd plan --context-only` is executed with profile set to `balanced`
- **THEN** the JSON output SHALL contain `"reviewer_model": "sonnet"`

#### Scenario: Plan context includes plan_checker_model for quality profile

- **WHEN** `mysd plan --context-only` is executed with profile set to `quality`
- **THEN** the JSON output SHALL contain `"plan_checker_model": "opus"`
