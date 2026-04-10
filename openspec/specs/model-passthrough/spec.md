---
spec-version: "1.0"
capability: Model Passthrough
delta: ADDED
status: done
---

## Requirement: Profile model resolution uses short names

The `DefaultModelMap` in `internal/config/config.go` SHALL store model values as short names (`sonnet`, `opus`, `haiku`) instead of full model IDs (`claude-sonnet-4-5`).

The `ResolveModel` function SHALL return short names that are directly usable as the `model` parameter when spawning Claude Code agents via the Task tool.

### Scenario: ResolveModel returns short name

WHEN `ResolveModel("planner", "quality", nil)` is called
THEN the return value SHALL be `"sonnet"` (not `"claude-sonnet-4-5"`)

### Scenario: Budget profile returns haiku for executor

WHEN `ResolveModel("executor", "budget", nil)` is called
THEN the return value SHALL be `"haiku"`

### Scenario: Override uses short name

WHEN `ResolveModel("executor", "quality", {"executor": "opus"})` is called
THEN the return value SHALL be `"opus"`

## Requirement: DefaultModelMap includes reviewer role

The `DefaultModelMap` in `internal/config/config.go` SHALL include a `reviewer` role with the following assignments:

| Profile | Model |
|---------|-------|
| quality | opus |
| balanced | sonnet |
| budget | sonnet |

Rationale: reviewer is a judgment role requiring reasoning capability. Budget profile uses `sonnet` (not `haiku`) because quality assessment requires sufficient model capability.

### Scenario: Quality profile returns opus for reviewer

- **WHEN** `ResolveModel("reviewer", "quality", nil)` is called
- **THEN** the return value SHALL be `"opus"`

### Scenario: Balanced profile returns sonnet for reviewer

- **WHEN** `ResolveModel("reviewer", "balanced", nil)` is called
- **THEN** the return value SHALL be `"sonnet"`

### Scenario: Budget profile returns sonnet for reviewer

- **WHEN** `ResolveModel("reviewer", "budget", nil)` is called
- **THEN** the return value SHALL be `"sonnet"`

<!-- Requirement "Binary context JSON includes model field" REMOVED by unified-model-resolve (2026-04-09). See model-resolve-command spec. -->

## Requirement: Command skills pass model to agents

Workflow command skills (`propose`, `discuss`, `plan`, `apply`, `ff`, `ffe`, `scan`, `uat`, `verify`, `fix`) SHALL use `mysd model resolve <role>` to obtain the model short name for each agent role they spawn, and pass it as the `model` parameter when spawning agent tasks.

Skills SHALL NOT parse the `mysd model` table output or read model fields from `--context-only` JSON for model resolution.

### Scenario: Plan command resolves models via model resolve

WHEN `/mysd:plan` needs to spawn mysd-designer, mysd-planner, mysd-reviewer, and mysd-plan-checker
THEN the skill SHALL execute `mysd model resolve designer`, `mysd model resolve planner`, `mysd model resolve reviewer`, and `mysd model resolve plan-checker`
AND use each output as the model parameter for the corresponding agent

### Scenario: Propose command resolves models via model resolve

WHEN `/mysd:propose` needs to spawn mysd-researcher, mysd-advisor, mysd-proposal-writer, and mysd-reviewer
THEN the skill SHALL execute `mysd model resolve researcher`, `mysd model resolve advisor`, `mysd model resolve proposal-writer`, and `mysd model resolve reviewer`
AND use each output as the model parameter for the corresponding agent

## Requirement: Command skills display model on agent spawn

Workflow command skills SHALL display the model being used when spawning each agent, in the format: `Spawning {agent-name} ({model})...`

### Scenario: Model display on spawn

WHEN `/mysd:apply` spawns `mysd-executor` with model `sonnet`
THEN the command SHALL display `Spawning mysd-executor (sonnet)...` before the spawn

## Requirement: Standalone utility commands specify fixed model

The following standalone commands SHALL specify `model: claude-sonnet-4-5` in their frontmatter: `status`, `lang`, `model`, `note`, `docs`, `statusline`, `update`.

The following standalone commands SHALL specify `model: claude-opus-4-6` in their frontmatter: `init`, `scan`, `fix`.

### Scenario: Utility command uses fixed sonnet

WHEN a user invokes `/mysd:status`
THEN the command SHALL run using `claude-sonnet-4-5` regardless of profile settings

### Scenario: Heavy standalone command uses fixed opus

WHEN a user invokes `/mysd:fix`
THEN the command SHALL run using `claude-opus-4-6` regardless of profile settings

## Requirement: Workflow commands and agents have no model frontmatter

All workflow command skills (`propose`, `discuss`, `plan`, `apply`, `archive`, `ff`, `ffe`, `uat`) SHALL NOT have a `model:` field in their frontmatter. They inherit the caller's model.

All agent definitions SHALL NOT have a `model:` field in their frontmatter. Their model is controlled by the profile system and passed by the calling command.

### Scenario: Workflow command inherits caller model

WHEN `/mysd:plan` is invoked from a session running `opus`
THEN the plan command orchestrator SHALL run using `opus`
AND the agents it spawns SHALL use the model resolved by the profile system

## Covered Packages

- `internal/config/config.go` — DefaultModelMap, ResolveModel
- `plugin/commands/` — all command skill frontmatter
- `plugin/agents/` — all agent definition frontmatter

## Requirements

### Requirement: Profile model resolution uses short names

The `DefaultModelMap` in `internal/config/config.go` SHALL store model values as short names (`sonnet`, `opus`, `haiku`) instead of full model IDs (`claude-sonnet-4-5`).

The `ResolveModel` function SHALL return short names that are directly usable as the `model` parameter when spawning Claude Code agents via the Task tool.

The `DefaultModelMap` SHALL assign models per profile as follows:

**quality profile** — thinking roles SHALL use `opus`, execution roles SHALL use `sonnet`:

| Role | Model |
|------|-------|
| spec-writer | opus |
| designer | opus |
| planner | opus |
| executor | sonnet |
| verifier | opus |
| fast-forward | sonnet |
| researcher | opus |
| advisor | opus |
| proposal-writer | opus |
| plan-checker | opus |
| reviewer | opus |

**balanced profile** — judgment/design/gating roles SHALL use `opus`, others SHALL use `sonnet`:

| Role | Model |
|------|-------|
| spec-writer | opus |
| designer | opus |
| planner | opus |
| executor | sonnet |
| verifier | opus |
| fast-forward | sonnet |
| researcher | sonnet |
| advisor | opus |
| proposal-writer | sonnet |
| plan-checker | opus |
| reviewer | sonnet |

**budget profile** — spec-writer SHALL use `sonnet`, planner/verifier/researcher/advisor/proposal-writer/plan-checker/reviewer SHALL use `sonnet`, designer/executor/fast-forward SHALL use `haiku`:

| Role | Model |
|------|-------|
| spec-writer | sonnet |
| designer | haiku |
| planner | sonnet |
| executor | haiku |
| verifier | sonnet |
| fast-forward | haiku |
| researcher | sonnet |
| advisor | sonnet |
| proposal-writer | sonnet |
| plan-checker | sonnet |
| reviewer | sonnet |

#### Scenario: Quality profile returns opus for thinking roles

- **WHEN** `ResolveModel("planner", "quality", nil)` is called
- **THEN** the return value SHALL be `"opus"`

#### Scenario: Quality profile returns sonnet for execution roles

- **WHEN** `ResolveModel("executor", "quality", nil)` is called
- **THEN** the return value SHALL be `"sonnet"`

#### Scenario: Balanced profile returns opus for gating roles

- **WHEN** `ResolveModel("verifier", "balanced", nil)` is called
- **THEN** the return value SHALL be `"opus"`

#### Scenario: Balanced profile returns sonnet for non-gating roles

- **WHEN** `ResolveModel("researcher", "balanced", nil)` is called
- **THEN** the return value SHALL be `"sonnet"`

#### Scenario: Budget profile returns sonnet for spec-writer

- **WHEN** `ResolveModel("spec-writer", "budget", nil)` is called
- **THEN** the return value SHALL be `"sonnet"`

#### Scenario: Budget profile returns haiku for executor

- **WHEN** `ResolveModel("executor", "budget", nil)` is called
- **THEN** the return value SHALL be `"haiku"`

#### Scenario: Override takes precedence over profile

- **WHEN** `ResolveModel("executor", "quality", {"executor": "opus"})` is called
- **THEN** the return value SHALL be `"opus"`


<!-- @trace
source: fix-default-model-map
updated: 2026-03-27
code:
  - plugin/commands/mysd-execute.md
  - plugin/agents/mysd-fast-forward.md
  - plugin/commands/mysd-uat.md
  - plugin/agents/mysd-planner.md
  - plugin/commands/mysd-design.md
  - cmd/design.go
  - plugin/agents/mysd-plan-checker.md
  - internal/config/config.go
  - plugin/commands/mysd-propose.md
  - plugin/agents/mysd-executor.md
  - plugin/commands/mysd-ff.md
  - plugin/agents/mysd-advisor.md
  - plugin/commands/mysd-capture.md
  - cmd/spec.go
  - plugin/commands/mysd-fix.md
  - plugin/commands/mysd-init.md
  - .specs/deferred.json
  - plugin/agents/mysd-spec-writer.md
  - plugin/commands/mysd-spec.md
  - plugin/commands/mysd-plan.md
  - plugin/commands/mysd-verify.md
  - plugin/commands/mysd-scan.md
  - plugin/agents/mysd-proposal-writer.md
  - plugin/commands/mysd-ffe.md
  - cmd/execute.go
  - plugin/agents/mysd-uat-guide.md
  - plugin/commands/mysd-archive.md
  - plugin/commands/mysd-discuss.md
  - .spectra.yaml
  - plugin/agents/mysd-researcher.md
  - plugin/commands/mysd-status.md
  - plugin/commands/mysd-apply.md
tests:
  - internal/executor/context_test.go
  - internal/spec/schema_test.go
  - internal/config/config_test.go
-->

<!-- Requirement "Binary context JSON includes model field" REMOVED by unified-model-resolve (2026-04-09).
     Reason: Model resolution responsibility transferred to dedicated `mysd model resolve <role>` subcommand.
     Embedding model fields in --context-only JSON mixed configuration concerns (model selection) with
     workflow context (spec_dir, tasks, change_name). Skills now use `mysd model resolve <role>` instead of
     reading model, verifier_model, reviewer_model, or plan_checker_model from --context-only JSON. -->

---
### Requirement: Command skills pass model to agents

Workflow command skills (`propose`, `discuss`, `plan`, `apply`, `ff`, `ffe`, `scan`, `uat`, `verify`, `fix`) SHALL use `mysd model resolve <role>` to obtain the model short name for each agent role they spawn, and pass it as the `model` parameter when spawning agent tasks.

Skills SHALL NOT parse the `mysd model` table output or read model fields from `--context-only` JSON for model resolution.

#### Scenario: Plan command resolves models via model resolve

- **WHEN** `/mysd:plan` needs to spawn mysd-designer, mysd-planner, mysd-reviewer, and mysd-plan-checker
- **THEN** the skill SHALL execute `mysd model resolve designer`, `mysd model resolve planner`, `mysd model resolve reviewer`, and `mysd model resolve plan-checker`
- **AND** use each output as the model parameter for the corresponding agent

#### Scenario: Propose command resolves models via model resolve

- **WHEN** `/mysd:propose` needs to spawn mysd-researcher, mysd-advisor, mysd-proposal-writer, and mysd-reviewer
- **THEN** the skill SHALL execute `mysd model resolve researcher`, `mysd model resolve advisor`, `mysd model resolve proposal-writer`, and `mysd model resolve reviewer`
- **AND** use each output as the model parameter for the corresponding agent

#### Scenario: Apply command resolves executor model via model resolve

- **WHEN** `/mysd:apply` needs to spawn mysd-executor
- **THEN** the skill SHALL execute `mysd model resolve executor`
- **AND** use the output as the model parameter when spawning mysd-executor


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
### Requirement: Command skills display model on agent spawn

Workflow command skills SHALL display the model being used when spawning each agent, in the format: `Spawning {agent-name} ({model})...`

#### Scenario: Model display on spawn

- **WHEN** `/mysd:apply` spawns `mysd-executor` with model `sonnet`
- **THEN** the command SHALL display `Spawning mysd-executor (sonnet)...` before the spawn


<!-- @trace
source: command-model-cleanup
updated: 2026-03-27
code:
  - plugin/commands/mysd-propose.md
  - .spectra.yaml
  - plugin/commands/mysd-verify.md
  - cmd/spec.go
  - plugin/commands/mysd-discuss.md
  - plugin/commands/mysd-ffe.md
  - plugin/agents/mysd-executor.md
  - plugin/commands/mysd-apply.md
  - plugin/agents/mysd-spec-writer.md
  - plugin/commands/mysd-status.md
  - plugin/agents/mysd-advisor.md
  - plugin/agents/mysd-uat-guide.md
  - plugin/commands/mysd-uat.md
  - plugin/agents/mysd-plan-checker.md
  - plugin/commands/mysd-spec.md
  - cmd/execute.go
  - plugin/agents/mysd-planner.md
  - plugin/commands/mysd-design.md
  - plugin/commands/mysd-execute.md
  - plugin/agents/mysd-researcher.md
  - internal/config/config.go
  - cmd/design.go
  - plugin/agents/mysd-proposal-writer.md
  - plugin/commands/mysd-archive.md
  - .specs/deferred.json
  - plugin/commands/mysd-ff.md
  - plugin/commands/mysd-scan.md
  - plugin/commands/mysd-init.md
  - plugin/commands/mysd-capture.md
  - plugin/commands/mysd-fix.md
  - plugin/agents/mysd-fast-forward.md
  - plugin/commands/mysd-plan.md
tests:
  - internal/config/config_test.go
  - internal/executor/context_test.go
  - internal/spec/schema_test.go
-->

---
### Requirement: Standalone utility commands specify fixed model

The following standalone commands SHALL specify `model: claude-sonnet-4-5` in their frontmatter: `status`, `lang`, `model`, `note`, `docs`, `statusline`, `update`.

The following standalone commands SHALL specify `model: claude-opus-4-6` in their frontmatter: `init`, `scan`, `fix`.

#### Scenario: Utility command uses fixed sonnet

- **WHEN** a user invokes `/mysd:status`
- **THEN** the command SHALL run using `claude-sonnet-4-5` regardless of profile settings

#### Scenario: Heavy standalone command uses fixed opus

- **WHEN** a user invokes `/mysd:fix`
- **THEN** the command SHALL run using `claude-opus-4-6` regardless of profile settings


<!-- @trace
source: command-model-cleanup
updated: 2026-03-27
code:
  - plugin/commands/mysd-propose.md
  - .spectra.yaml
  - plugin/commands/mysd-verify.md
  - cmd/spec.go
  - plugin/commands/mysd-discuss.md
  - plugin/commands/mysd-ffe.md
  - plugin/agents/mysd-executor.md
  - plugin/commands/mysd-apply.md
  - plugin/agents/mysd-spec-writer.md
  - plugin/commands/mysd-status.md
  - plugin/agents/mysd-advisor.md
  - plugin/agents/mysd-uat-guide.md
  - plugin/commands/mysd-uat.md
  - plugin/agents/mysd-plan-checker.md
  - plugin/commands/mysd-spec.md
  - cmd/execute.go
  - plugin/agents/mysd-planner.md
  - plugin/commands/mysd-design.md
  - plugin/commands/mysd-execute.md
  - plugin/agents/mysd-researcher.md
  - internal/config/config.go
  - cmd/design.go
  - plugin/agents/mysd-proposal-writer.md
  - plugin/commands/mysd-archive.md
  - .specs/deferred.json
  - plugin/commands/mysd-ff.md
  - plugin/commands/mysd-scan.md
  - plugin/commands/mysd-init.md
  - plugin/commands/mysd-capture.md
  - plugin/commands/mysd-fix.md
  - plugin/agents/mysd-fast-forward.md
  - plugin/commands/mysd-plan.md
tests:
  - internal/config/config_test.go
  - internal/executor/context_test.go
  - internal/spec/schema_test.go
-->

---
### Requirement: Workflow commands and agents have no model frontmatter

All workflow command skills (`propose`, `discuss`, `plan`, `apply`, `archive`, `ff`, `ffe`, `uat`) SHALL NOT have a `model:` field in their frontmatter. They inherit the caller's model.

All agent definitions SHALL NOT have a `model:` field in their frontmatter. Their model is controlled by the profile system and passed by the calling command.

#### Scenario: Workflow command inherits caller model

- **WHEN** `/mysd:plan` is invoked from a session running `opus`
- **THEN** the plan command orchestrator SHALL run using `opus`
- **AND** the agents it spawns SHALL use the model resolved by the profile system

<!-- @trace
source: command-model-cleanup
updated: 2026-03-27
code:
  - plugin/commands/mysd-propose.md
  - .spectra.yaml
  - plugin/commands/mysd-verify.md
  - cmd/spec.go
  - plugin/commands/mysd-discuss.md
  - plugin/commands/mysd-ffe.md
  - plugin/agents/mysd-executor.md
  - plugin/commands/mysd-apply.md
  - plugin/agents/mysd-spec-writer.md
  - plugin/commands/mysd-status.md
  - plugin/agents/mysd-advisor.md
  - plugin/agents/mysd-uat-guide.md
  - plugin/commands/mysd-uat.md
  - plugin/agents/mysd-plan-checker.md
  - plugin/commands/mysd-spec.md
  - cmd/execute.go
  - plugin/agents/mysd-planner.md
  - plugin/commands/mysd-design.md
  - plugin/commands/mysd-execute.md
  - plugin/agents/mysd-researcher.md
  - internal/config/config.go
  - cmd/design.go
  - plugin/agents/mysd-proposal-writer.md
  - plugin/commands/mysd-archive.md
  - .specs/deferred.json
  - plugin/commands/mysd-ff.md
  - plugin/commands/mysd-scan.md
  - plugin/commands/mysd-init.md
  - plugin/commands/mysd-capture.md
  - plugin/commands/mysd-fix.md
  - plugin/agents/mysd-fast-forward.md
  - plugin/commands/mysd-plan.md
tests:
  - internal/config/config_test.go
  - internal/executor/context_test.go
  - internal/spec/schema_test.go
-->

---
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