## ADDED Requirements

### Requirement: Profile model resolution uses short names

The `DefaultModelMap` in `internal/config/config.go` SHALL store model values as short names (`sonnet`, `opus`, `haiku`) instead of full model IDs (`claude-sonnet-4-5`).

The `ResolveModel` function SHALL return short names that are directly usable as the `model` parameter when spawning Claude Code agents via the Task tool.

#### Scenario: ResolveModel returns short name

- **WHEN** `ResolveModel("planner", "quality", nil)` is called
- **THEN** the return value SHALL be `"sonnet"` (not `"claude-sonnet-4-5"`)

#### Scenario: Budget profile returns haiku for executor

- **WHEN** `ResolveModel("executor", "budget", nil)` is called
- **THEN** the return value SHALL be `"haiku"`

#### Scenario: Override uses short name

- **WHEN** `ResolveModel("executor", "quality", {"executor": "opus"})` is called
- **THEN** the return value SHALL be `"opus"`

### Requirement: Binary context JSON includes model field

The `--context-only` output of workflow commands (`plan`, `execute`, `design`, `spec`) SHALL include a `"model"` field containing the resolved short name for the relevant agent role.

#### Scenario: Plan context includes planner model

- **WHEN** `mysd plan --context-only` is executed with profile set to `balanced`
- **THEN** the JSON output SHALL contain `"model": "sonnet"`

### Requirement: Command skills pass model to agents

Workflow command skills (`propose`, `discuss`, `plan`, `apply`, `ff`, `ffe`) SHALL read the `model` field from the binary's `--context-only` JSON output and pass it as the `model` parameter when spawning agent tasks.

#### Scenario: Plan command spawns designer with profile model

- **WHEN** `/mysd:plan` reads context JSON containing `"model": "opus"`
- **THEN** the command SHALL spawn `mysd-designer` with `model: opus`

### Requirement: Command skills display model on agent spawn

Workflow command skills SHALL display the model being used when spawning each agent, in the format: `Spawning {agent-name} ({model})...`

#### Scenario: Model display on spawn

- **WHEN** `/mysd:apply` spawns `mysd-executor` with model `sonnet`
- **THEN** the command SHALL display `Spawning mysd-executor (sonnet)...` before the spawn

### Requirement: Standalone utility commands specify fixed model

The following standalone commands SHALL specify `model: claude-sonnet-4-5` in their frontmatter: `status`, `lang`, `model`, `note`, `docs`, `statusline`, `update`.

The following standalone commands SHALL specify `model: claude-opus-4-6` in their frontmatter: `init`, `scan`, `fix`.

#### Scenario: Utility command uses fixed sonnet

- **WHEN** a user invokes `/mysd:status`
- **THEN** the command SHALL run using `claude-sonnet-4-5` regardless of profile settings

#### Scenario: Heavy standalone command uses fixed opus

- **WHEN** a user invokes `/mysd:fix`
- **THEN** the command SHALL run using `claude-opus-4-6` regardless of profile settings

### Requirement: Workflow commands and agents have no model frontmatter

All workflow command skills (`propose`, `discuss`, `plan`, `apply`, `archive`, `ff`, `ffe`, `uat`) SHALL NOT have a `model:` field in their frontmatter. They inherit the caller's model.

All agent definitions SHALL NOT have a `model:` field in their frontmatter. Their model is controlled by the profile system and passed by the calling command.

#### Scenario: Workflow command inherits caller model

- **WHEN** `/mysd:plan` is invoked from a session running `opus`
- **THEN** the plan command orchestrator SHALL run using `opus`
- **AND** the agents it spawns SHALL use the model resolved by the profile system
