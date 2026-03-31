## ADDED Requirements

### Requirement: GenerateInstruction function

The `internal/executor` package SHALL export a `GenerateInstruction` function that accepts an `ExecutionContext` and an optional `*PreflightReport`, and returns a string containing one or more instruction segments.

The function SHALL be a pure function with no I/O side effects.

#### Scenario: GenerateInstruction returns non-empty string

- **WHEN** `GenerateInstruction` is called with a valid `ExecutionContext`
- **THEN** the returned string SHALL contain at least one instruction segment

### Requirement: Task state instruction segments

The `GenerateInstruction` function SHALL produce exactly one of the following mutually exclusive task state segments:

1. **all_done**: WHEN all tasks have status "done", the instruction SHALL state that all tasks are complete and suggest verify or archive
2. **has_failed**: WHEN any task has status "blocked" or "failed", the instruction SHALL identify the failed task IDs and suggest retry or skip
3. **resume**: WHEN at least one task is "done" and at least one is "pending", the instruction SHALL state the progress (done/total) and identify the next task to continue from
4. **last_task**: WHEN exactly one task is "pending" and all others are "done", the instruction SHALL identify the final task and note that verify follows
5. **first_run**: WHEN all tasks are "pending", the instruction SHALL state the total count and identify the first task

The segments SHALL be evaluated in the priority order listed above (1 highest, 5 lowest). The first matching segment SHALL be used.

#### Scenario: All tasks complete

- **WHEN** `ExecutionContext` has 5 tasks all with status "done"
- **THEN** the instruction SHALL contain "All 5 tasks complete"
- **AND** SHALL mention verify or archive

#### Scenario: Resume from interruption

- **WHEN** `ExecutionContext` has 3 done tasks and 2 pending tasks
- **AND** the first pending task is T4
- **THEN** the instruction SHALL contain "3/5 complete" and "continue from T4"

#### Scenario: Last remaining task

- **WHEN** `ExecutionContext` has 4 done tasks and 1 pending task T5
- **THEN** the instruction SHALL contain "T5" and mention verify

#### Scenario: First execution

- **WHEN** `ExecutionContext` has 5 pending tasks starting with T1
- **THEN** the instruction SHALL contain "5 tasks pending" and "start from T1"

#### Scenario: Failed task detected

- **WHEN** `ExecutionContext` has task T3 with status "blocked"
- **THEN** the instruction SHALL contain "T3" and suggest retry or skip

### Requirement: Preflight instruction segments

The `GenerateInstruction` function SHALL append preflight-related segments when a `PreflightReport` is provided and contains issues. Preflight segments are additive — they combine with the task state segment.

1. **stale**: WHEN `PreflightReport.Checks.Staleness.IsStale` is true and `DaysSinceLastPlan` > 0, the instruction SHALL state the number of days since last plan and suggest re-planning
2. **missing_files**: WHEN `PreflightReport.Checks.MissingFiles` has length > 0, the instruction SHALL state the count of missing files and advise review before starting

#### Scenario: Stale artifacts warning

- **WHEN** `PreflightReport` has `DaysSinceLastPlan: 15` and `IsStale: true`
- **THEN** the instruction SHALL contain "15 days since last plan"

#### Scenario: Missing files warning

- **WHEN** `PreflightReport` has 2 entries in `MissingFiles`
- **THEN** the instruction SHALL contain "2 missing files"

#### Scenario: Combined task state and preflight

- **WHEN** `ExecutionContext` has a resume state (3/5 done)
- **AND** `PreflightReport` has stale artifacts
- **THEN** the instruction SHALL contain both the resume segment and the stale segment

#### Scenario: No preflight issues

- **WHEN** `PreflightReport` is nil or has status "ok"
- **THEN** no preflight segments SHALL be appended

### Requirement: Instruction output in ExecutionContext JSON

The `ExecutionContext` struct SHALL include an `Instruction` field with JSON tag `json:"instruction"`.

When `--context-only` is used, the `instruction` field SHALL be populated by calling `GenerateInstruction` before JSON serialization.

#### Scenario: Instruction field present in JSON output

- **WHEN** `mysd execute --context-only` is run
- **THEN** the JSON output SHALL contain an `instruction` field with a non-empty string value
