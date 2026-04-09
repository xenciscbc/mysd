---
spec-version: "1.0"
capability: execution
delta: MODIFIED
status: pending
---

## ADDED Requirements

### Requirement: Context-only JSON includes spec_dir

All `--context-only` JSON outputs MUST include a `spec_dir` field containing the detected spec directory name (`.specs` or `openspec`).

The following commands MUST include `spec_dir`:
- `mysd spec --context-only`
- `mysd plan --context-only`
- `mysd design --context-only`
- `mysd execute --context-only`
- `mysd scan --context-only`
- `mysd verify --context-only`

#### Scenario: openspec project outputs correct spec_dir

- **WHEN** a project uses `openspec/` as its spec directory
- **THEN** all `--context-only` JSON outputs MUST contain `"spec_dir": "openspec"`

#### Scenario: mysd-flavor project outputs correct spec_dir

- **WHEN** a project uses `.specs/` as its spec directory
- **THEN** all `--context-only` JSON outputs MUST contain `"spec_dir": ".specs"`

### Requirement: Agent definitions use dynamic spec_dir

All agent definition files (mysd/agents/*.md) MUST reference artifact paths using `{spec_dir}` placeholder instead of hardcoded `.specs`.

Agent paths MUST follow the pattern: `{spec_dir}/changes/{change_name}/`

#### Scenario: Agent reads proposal in openspec project

- **WHEN** an agent receives `spec_dir: "openspec"` in its context
- **THEN** the agent MUST read `openspec/changes/{change_name}/proposal.md`

### Requirement: Orchestrators pass spec_dir to agents

All orchestrator SKILL.md files that spawn agents MUST extract `spec_dir` from `--context-only` JSON output and include it in the agent's Task context JSON.

#### Scenario: Plan orchestrator passes spec_dir to designer

- **WHEN** the plan orchestrator spawns mysd-designer
- **THEN** the Task context JSON MUST include `"spec_dir": "{spec_dir}"`

### Requirement: Validator supports TasksFrontmatterV2 task count

The `mysd validate` command MUST correctly count tasks from TasksFrontmatterV2 format where tasks are defined in the YAML frontmatter `tasks` array, not as markdown checkboxes in the body.

When `tasks` array is present and non-empty in frontmatter, the validator MUST use `len(tasks)` for the actual task count instead of counting `- [ ]` lines.

When `tasks` array is absent or empty, the validator MUST fall back to counting markdown checkbox lines (V1 behavior).

#### Scenario: V2 tasks.md validates correctly

- **WHEN** tasks.md has `total: 3` in frontmatter and 3 entries in the `tasks` array
- **THEN** `mysd validate` MUST NOT report a task count mismatch warning

#### Scenario: V1 tasks.md still validates correctly

- **WHEN** tasks.md has `total: 3` in frontmatter and 3 `- [ ]` lines in body with no `tasks` array
- **THEN** `mysd validate` MUST NOT report a task count mismatch warning
