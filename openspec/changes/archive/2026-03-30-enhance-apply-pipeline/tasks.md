## 1. spec-executor Model Role (D-02, per-spec-execution spec)

- [x] 1.1 Add `spec-executor` role to `DefaultModelMap` in `internal/config/config.go`: quality=opus, balanced=opus, budget=sonnet. This implements the "spec-executor model role" requirement.
- [x] 1.2 Write test in `internal/config/config_test.go` for `spec-executor` role: verify ResolveModel returns opus for quality/balanced, sonnet for budget. Add to TestResolveModel_AllRoles expected map.

## 2. Preflight Check CLI (D-05, execution-preflight spec)

- [x] 2.1 Add `--preflight` flag to `cmd/execute.go`. When set with `--context-only`, read tasks.md `files` fields, check file existence (skip files where task description contains "create" or "add"), read STATE.json `last_run` for staleness (>7 days = warning, >30 days = critical). Output JSON with `status`, `checks.missing_files`, `checks.staleness`. This implements the "Preflight check CLI command" requirement.
- [x] 2.2 Write tests in `cmd/execute_test.go` for `--preflight`: all files exist (status ok), missing file detected (status warning), new file task excluded, stale artifacts warning (>7 days), critical staleness (>30 days).

## 3. Per-spec Execution Mode (D-01, D-06, per-spec-execution spec)

- [x] 3.1 Update `mysd/skills/apply/SKILL.md` Step 3: add "Spec Mode" section alongside Single/Wave. In spec mode, group pending tasks by `spec` field, resolve model via `spec-executor` role, spawn one executor per spec group with `assigned_tasks` array. Change-level tasks (empty spec) grouped separately and executed last. This implements the "Per-spec execution mode" requirement.
- [x] 3.2 Update `mysd/skills/apply/SKILL.md` Step 2: add Step 2c (Preflight Check) after Step 2b. Call `mysd execute --preflight --json`, display warnings, ask confirmation if warning/critical (skip in auto_mode). This implements the "Apply orchestrator preflight step" requirement.

## 4. Executor Pre-task Checks and Pause Conditions (D-03, D-04, execution spec)

- [x] 4.1 Update `mysd/agents/mysd-executor.md` Task Execution section: add Step 1b (Pre-Task Checks) after Step 1 (Mark In Progress). Define 4 checks: Reuse (search adjacent modules), Quality (use existing types/constants), Efficiency (parallelize async), No Placeholders (read spec section, pause if TBD found). This implements the "Executor pre-task checks" requirement.
- [x] 4.2 Update `mysd/agents/mysd-executor.md`: add "Pause Conditions" section defining 4 conditions (task unclear, design issue, error/blocker, user interrupt). When triggered, agent outputs issue description + resolution options + waits. This implements the "Executor pause conditions" requirement.
- [x] 4.3 Update `mysd/agents/mysd-executor.md` Input section: document `assigned_tasks` (array) as alternative to `assigned_task` (single). When array is present, agent executes each task in order with full cycle (pre-task checks, implement, mark done). Add explicit instruction to re-read spec/design before each task. This implements the "Executor handles assigned_tasks array" requirement.

## 5. Execution Mode Config (execution spec)

- [x] 5.1 Update `internal/executor/context.go` or `internal/config/config.go`: ensure `ExecutionMode` validation accepts "spec" in addition to "single" and "wave". Update any switch/if statements that check execution mode values. This implements the "Execution Modes" MODIFIED requirement.
- [x] 5.2 Write test: verify `execution_mode: "spec"` is accepted in config parsing and context building.
