# no-capability-change Specification

## Purpose

TBD - created by archiving change 'commands-to-subdirectory'. Update Purpose after archive.

## Requirements

### Requirement: Plugin command files use subdirectory structure

The mysd plugin command files SHALL be organized under `plugin/commands/mysd/` subdirectory instead of flat `plugin/commands/mysd-*.md` files.

#### Scenario: Command file location

- **WHEN** the plugin is installed
- **THEN** command files SHALL exist at `plugin/commands/mysd/<name>.md` (e.g., `plugin/commands/mysd/apply.md`)

#### Scenario: Slash command invocation unchanged

- **WHEN** a user invokes `/mysd:apply`
- **THEN** Claude Code SHALL resolve the command from `plugin/commands/mysd/apply.md`

<!-- @trace
source: commands-to-subdirectory
updated: 2026-03-30
code:
  - internal/config/config.go
  - plugin/commands/mysd-docs.md
  - mysd/agents/mysd-researcher.md
  - .claude-plugin/marketplace.json
  - mysd/skills/ffe/SKILL.md
  - mysd/skills/note/SKILL.md
  - mysd/skills/statusline/SKILL.md
  - mysd/agents/mysd-plan-checker.md
  - mysd/hooks/hooks.json
  - mysd/skills/scan/SKILL.md
  - plugin/commands/mysd-plan.md
  - cmd/analyze.go
  - plugin/agents/mysd-advisor.md
  - mysd/skills/update/SKILL.md
  - plugin/commands/mysd-archive.md
  - mysd/agents/mysd-advisor.md
  - plugin/agents/mysd-researcher.md
  - mysd/agents/mysd-executor.md
  - plugin/commands/mysd-uat.md
  - mysd/agents/mysd-uat-guide.md
  - mysd/skills/propose/SKILL.md
  - plugin/commands/mysd-ff.md
  - mysd/skills/init/SKILL.md
  - mysd/skills/lang/SKILL.md
  - plugin/agents/mysd-executor.md
  - plugin/agents/mysd-verifier.md
  - plugin/commands/mysd-statusline.md
  - internal/analyzer/analyzer.go
  - mysd/skills/docs/SKILL.md
  - plugin/agents/mysd-proposal-writer.md
  - mysd/agents/mysd-verifier.md
  - plugin/commands/mysd-lang.md
  - mysd/skills/status/SKILL.md
  - mysd/skills/fix/SKILL.md
  - mysd/agents/mysd-reviewer.md
  - plugin/agents/mysd-scanner.md
  - mysd/skills/apply/SKILL.md
  - mysd/skills/model/SKILL.md
  - plugin/agents/mysd-fast-forward.md
  - mysd/agents/mysd-scanner.md
  - plugin/agents/mysd-plan-checker.md
  - mysd/agents/mysd-spec-writer.md
  - plugin/commands/mysd-init.md
  - plugin/commands/mysd-verify.md
  - internal/analyzer/consistency.go
  - internal/analyzer/ambiguity.go
  - plugin/commands/mysd-ffe.md
  - mysd/agents/mysd-planner.md
  - internal/analyzer/types.go
  - plugin/agents/mysd-uat-guide.md
  - plugin/commands/mysd-note.md
  - mysd/agents/mysd-designer.md
  - internal/analyzer/gaps.go
  - plugin/commands/mysd-apply.md
  - plugin/commands/mysd-model.md
  - plugin/commands/mysd-status.md
  - mysd/skills/plan/SKILL.md
  - mysd/agents/mysd-proposal-writer.md
  - cmd/plan.go
  - mysd/skills/verify/SKILL.md
  - plugin/commands/mysd-discuss.md
  - mysd/skills/ff/SKILL.md
  - plugin/agents/mysd-planner.md
  - plugin/agents/mysd-designer.md
  - mysd/hooks/mysd-statusline.js
  - mysd/skills/archive/SKILL.md
  - plugin/hooks/hooks.json
  - plugin/hooks/mysd-statusline.js
  - mysd/agents/mysd-fast-forward.md
  - plugin/commands/mysd-update.md
  - plugin/commands/mysd-propose.md
  - mysd/skills/discuss/SKILL.md
  - plugin/commands/mysd-scan.md
  - internal/analyzer/coverage.go
  - mysd/skills/uat/SKILL.md
  - plugin/agents/mysd-spec-writer.md
  - plugin/commands/mysd-fix.md
tests:
  - cmd/analyze_test.go
  - internal/config/config_test.go
  - internal/analyzer/analyzer_test.go
  - cmd/plan_test.go
-->