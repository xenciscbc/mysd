---
spec-version: "1.0"
capability: Model Resolve Command
delta: ADDED
status: done
---

## Requirement: Model resolve subcommand returns single role model

The `mysd model resolve <role>` subcommand SHALL accept exactly one role name as a positional argument and print the resolved model short name to stdout, followed by a newline.

The command SHALL use the current profile (from config or `--profile` flag if present) and apply model overrides and custom profiles in the same order as `ResolveModel`.

The output SHALL contain only the model short name (e.g., `sonnet`, `opus`, `haiku`) with no additional formatting, headers, or whitespace beyond a trailing newline.

### Scenario: Resolve executor role under balanced profile

- **WHEN** the current profile is `balanced` and `mysd model resolve executor` is executed
- **THEN** stdout SHALL contain exactly `sonnet\n`

### Scenario: Resolve planner role under quality profile

- **WHEN** the current profile is `quality` and `mysd model resolve planner` is executed
- **THEN** stdout SHALL contain exactly `opus\n`

### Scenario: Resolve role with custom profile override

- **WHEN** the current profile is a custom profile with `base: balanced` and `models: { executor: opus }`, and `mysd model resolve executor` is executed
- **THEN** stdout SHALL contain exactly `opus\n`

<!-- @trace
source: unified-model-resolve
updated: 2026-04-09
code:
  - cmd/model.go
tests:
  - cmd/model_test.go
-->

---
## Requirement: Model resolve validates role name

The `mysd model resolve` subcommand SHALL exit with a non-zero exit code and print an error message to stderr when:
- No role argument is provided
- The provided role name does not exist in `DefaultModelMap`

### Scenario: No role argument provided

- **WHEN** `mysd model resolve` is executed without a role argument
- **THEN** the command SHALL exit with non-zero status
- **AND** stderr SHALL contain an error message indicating a role argument is required

### Scenario: Unknown role name

- **WHEN** `mysd model resolve nonexistent-role` is executed
- **THEN** the command SHALL exit with non-zero status
- **AND** stderr SHALL contain an error message indicating the role is unknown

<!-- @trace
source: unified-model-resolve
updated: 2026-04-09
code:
  - cmd/model.go
tests:
  - cmd/model_test.go
-->

---
## Requirement: Model resolve is the single source of truth for skill model queries

All workflow command skills that need to determine which model to use for spawning agents SHALL use `mysd model resolve <role>` instead of:
- Parsing the `mysd model` table output
- Reading `model` fields from `--context-only` JSON output

Each skill SHALL call `mysd model resolve <role>` once per role it needs, capturing the plain-text output as the model name.

### Scenario: Propose skill resolves three models via model resolve

- **WHEN** `/mysd:propose` needs models for researcher, advisor, and proposal-writer roles
- **THEN** the skill SHALL execute `mysd model resolve researcher`, `mysd model resolve advisor`, and `mysd model resolve proposal-writer` separately
- **AND** use each output as the model parameter when spawning the corresponding agent

### Scenario: Scan skill resolves scanner model via model resolve

- **WHEN** `/mysd:scan` needs the scanner model
- **THEN** the skill SHALL execute `mysd model resolve scanner`
- **AND** use the output as the model parameter when spawning the scanner agent

### Scenario: Apply skill resolves executor and verifier models via model resolve

- **WHEN** `/mysd:apply` needs executor and verifier models
- **THEN** the skill SHALL execute `mysd model resolve executor` and `mysd model resolve verifier` separately
- **AND** use each output as the model parameter when spawning the corresponding agent

<!-- @trace
source: unified-model-resolve
updated: 2026-04-09
code:
  - mysd/skills/propose/SKILL.md
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/plan/SKILL.md
  - mysd/skills/apply/SKILL.md
  - mysd/skills/ff/SKILL.md
  - mysd/skills/ffe/SKILL.md
  - mysd/skills/scan/SKILL.md
  - mysd/skills/uat/SKILL.md
  - mysd/skills/verify/SKILL.md
  - mysd/skills/fix/SKILL.md
-->

## Requirements

### Requirement: Model resolve subcommand returns single role model

The `mysd model resolve <role>` subcommand SHALL accept exactly one role name as a positional argument and print the resolved model short name to stdout, followed by a newline.

The command SHALL use the current profile (from config or `--profile` flag if present) and apply model overrides and custom profiles in the same order as `ResolveModel`.

The output SHALL contain only the model short name (e.g., `sonnet`, `opus`, `haiku`) with no additional formatting, headers, or whitespace beyond a trailing newline.

#### Scenario: Resolve executor role under balanced profile

- **WHEN** the current profile is `balanced` and `mysd model resolve executor` is executed
- **THEN** stdout SHALL contain exactly `sonnet\n`

#### Scenario: Resolve planner role under quality profile

- **WHEN** the current profile is `quality` and `mysd model resolve planner` is executed
- **THEN** stdout SHALL contain exactly `opus\n`

#### Scenario: Resolve role with custom profile override

- **WHEN** the current profile is a custom profile with `base: balanced` and `models: { executor: opus }`, and `mysd model resolve executor` is executed
- **THEN** stdout SHALL contain exactly `opus\n`


<!-- @trace
source: unified-model-resolve
updated: 2026-04-10
code:
  - mysd/skills/apply/SKILL.md
  - cmd/model.go
  - cmd/verify.go
  - mysd/skills/fix/SKILL.md
  - cmd/design.go
  - mysd/skills/uat/SKILL.md
  - cmd/plan.go
  - cmd/execute.go
  - mysd/skills/propose/SKILL.md
  - cmd/spec.go
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/scan/SKILL.md
  - internal/verifier/context.go
  - mysd/skills/verify/SKILL.md
  - mysd/skills/ff/SKILL.md
  - mysd/skills/ffe/SKILL.md
  - internal/executor/context.go
  - mysd/skills/plan/SKILL.md
tests:
  - cmd/model_test.go
  - cmd/plan_test.go
-->

---
### Requirement: Model resolve validates role name

The `mysd model resolve` subcommand SHALL exit with a non-zero exit code and print an error message to stderr when:
- No role argument is provided
- The provided role name does not exist in `DefaultModelMap`

#### Scenario: No role argument provided

- **WHEN** `mysd model resolve` is executed without a role argument
- **THEN** the command SHALL exit with non-zero status
- **AND** stderr SHALL contain an error message indicating a role argument is required

#### Scenario: Unknown role name

- **WHEN** `mysd model resolve nonexistent-role` is executed
- **THEN** the command SHALL exit with non-zero status
- **AND** stderr SHALL contain an error message indicating the role is unknown


<!-- @trace
source: unified-model-resolve
updated: 2026-04-10
code:
  - mysd/skills/apply/SKILL.md
  - cmd/model.go
  - cmd/verify.go
  - mysd/skills/fix/SKILL.md
  - cmd/design.go
  - mysd/skills/uat/SKILL.md
  - cmd/plan.go
  - cmd/execute.go
  - mysd/skills/propose/SKILL.md
  - cmd/spec.go
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/scan/SKILL.md
  - internal/verifier/context.go
  - mysd/skills/verify/SKILL.md
  - mysd/skills/ff/SKILL.md
  - mysd/skills/ffe/SKILL.md
  - internal/executor/context.go
  - mysd/skills/plan/SKILL.md
tests:
  - cmd/model_test.go
  - cmd/plan_test.go
-->

---
### Requirement: Model resolve is the single source of truth for skill model queries

All workflow command skills that need to determine which model to use for spawning agents SHALL use `mysd model resolve <role>` instead of:
- Parsing the `mysd model` table output
- Reading `model` fields from `--context-only` JSON output

Each skill SHALL call `mysd model resolve <role>` once per role it needs, capturing the plain-text output as the model name.

#### Scenario: Propose skill resolves three models via model resolve

- **WHEN** `/mysd:propose` needs models for researcher, advisor, and proposal-writer roles
- **THEN** the skill SHALL execute `mysd model resolve researcher`, `mysd model resolve advisor`, and `mysd model resolve proposal-writer` separately
- **AND** use each output as the model parameter when spawning the corresponding agent

#### Scenario: Scan skill resolves scanner model via model resolve

- **WHEN** `/mysd:scan` needs the scanner model
- **THEN** the skill SHALL execute `mysd model resolve scanner`
- **AND** use the output as the model parameter when spawning the scanner agent

#### Scenario: Apply skill resolves executor and verifier models via model resolve

- **WHEN** `/mysd:apply` needs executor and verifier models
- **THEN** the skill SHALL execute `mysd model resolve executor` and `mysd model resolve verifier` separately
- **AND** use each output as the model parameter when spawning the corresponding agent

<!-- @trace
source: unified-model-resolve
updated: 2026-04-10
code:
  - mysd/skills/apply/SKILL.md
  - cmd/model.go
  - cmd/verify.go
  - mysd/skills/fix/SKILL.md
  - cmd/design.go
  - mysd/skills/uat/SKILL.md
  - cmd/plan.go
  - cmd/execute.go
  - mysd/skills/propose/SKILL.md
  - cmd/spec.go
  - mysd/skills/discuss/SKILL.md
  - mysd/skills/scan/SKILL.md
  - internal/verifier/context.go
  - mysd/skills/verify/SKILL.md
  - mysd/skills/ff/SKILL.md
  - mysd/skills/ffe/SKILL.md
  - internal/executor/context.go
  - mysd/skills/plan/SKILL.md
tests:
  - cmd/model_test.go
  - cmd/plan_test.go
-->