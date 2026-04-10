## REMOVED Requirements

### Requirement: Binary context JSON includes model field

**Reason**: Model resolution responsibility is transferred to the dedicated `mysd model resolve <role>` subcommand. Embedding model fields in `--context-only` JSON mixed configuration concerns (model selection) with workflow context (spec_dir, tasks, change_name). Removing model fields from `--context-only` enforces single responsibility: `--context-only` provides workflow context, `model resolve` provides model selection.

**Migration**: All skills that previously read `model`, `verifier_model`, `reviewer_model`, or `plan_checker_model` from `--context-only` JSON SHALL use `mysd model resolve <role>` instead. Specifically:
- `model` → `mysd model resolve <role>` (where role depends on the command: planner, executor, verifier, designer, spec-writer)
- `verifier_model` → `mysd model resolve verifier`
- `reviewer_model` → `mysd model resolve reviewer`
- `plan_checker_model` → `mysd model resolve plan-checker`

#### Scenario: Context-only JSON no longer contains model fields

- **WHEN** `mysd plan --context-only` is executed after this change
- **THEN** the JSON output SHALL NOT contain `model`, `reviewer_model`, or `plan_checker_model` fields

#### Scenario: Execute context-only JSON no longer contains model fields

- **WHEN** `mysd execute --context-only` is executed after this change
- **THEN** the JSON output SHALL NOT contain `model` or `verifier_model` fields

## MODIFIED Requirements

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
