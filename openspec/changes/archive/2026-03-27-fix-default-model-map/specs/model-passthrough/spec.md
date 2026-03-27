## MODIFIED Requirements

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

**budget profile** — spec-writer SHALL use `sonnet`, planner/verifier/researcher/advisor/proposal-writer/plan-checker SHALL use `sonnet`, designer/executor/fast-forward SHALL use `haiku`:

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
