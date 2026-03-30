---
spec-version: "1.0"
capability: Execution Preflight Check
delta: ADDED
status: draft
---

## ADDED Requirements

### Requirement: Preflight check CLI command

The `mysd execute --preflight` flag SHALL perform pre-execution validation and output a JSON report.

The report SHALL include:
- `status`: "ok", "warning", or "critical"
- `checks.missing_files`: Array of file paths referenced in tasks that do not exist on disk
- `checks.staleness.days_since_last_plan`: Number of days since STATE.json `last_run`
- `checks.staleness.is_stale`: Boolean (true if days > 7)

Missing file detection SHALL skip files where the task description contains "create" or "add" keywords (indicating the file is expected to not exist yet).

Staleness SHALL be computed from the `last_run` field in STATE.json. If STATE.json does not exist, staleness SHALL be reported as critical with `days_since_last_plan: -1`.

The overall `status` SHALL be:
- "ok" if no missing files and not stale
- "warning" if stale (> 7 days) or has missing files
- "critical" if stale (> 30 days)

#### Scenario: All files exist and artifacts are fresh

- **WHEN** `mysd execute --preflight --json` is executed
- **AND** all task files exist on disk
- **AND** STATE.json last_run is 2 days ago
- **THEN** the output SHALL have `status: "ok"` and empty `missing_files`

#### Scenario: Missing file detected

- **WHEN** a task references `internal/foo.go` in its files field
- **AND** `internal/foo.go` does not exist
- **AND** the task description does not contain "create" or "add"
- **THEN** `missing_files` SHALL include `internal/foo.go`

#### Scenario: New file creation task excluded

- **WHEN** a task references `internal/new.go` in its files field
- **AND** `internal/new.go` does not exist
- **AND** the task description contains "Create internal/new.go"
- **THEN** `missing_files` SHALL NOT include `internal/new.go`

#### Scenario: Stale artifacts warning

- **WHEN** STATE.json last_run is 10 days ago
- **THEN** `staleness.is_stale` SHALL be true
- **AND** `status` SHALL be "warning"

### Requirement: Apply orchestrator preflight step

The `mysd:apply` skill SHALL call `mysd execute --preflight --json` after getting execution context (Step 2) and before task execution (Step 3).

If preflight status is "warning" or "critical":
- The orchestrator SHALL display the issues found
- The orchestrator SHALL ask the user for confirmation to continue
- In `auto_mode`, the orchestrator SHALL display warnings but proceed without confirmation

If preflight status is "ok": proceed silently.

#### Scenario: Preflight warning prompts user

- **WHEN** preflight returns `status: "warning"` with 1 missing file
- **AND** `auto_mode` is false
- **THEN** the orchestrator SHALL display the missing file and ask to continue or abort

#### Scenario: Preflight auto-mode continues

- **WHEN** preflight returns `status: "warning"`
- **AND** `auto_mode` is true
- **THEN** the orchestrator SHALL display the warning and proceed without asking
