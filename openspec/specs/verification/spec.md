---
spec-version: "1.0"
capability: Verification & Archive
delta: ADDED
status: done
---

## Requirement: Verification Context

The `mysd verify` command MUST build a `VerificationContext` JSON for verifier agents containing:
- ChangeName, MustItems, ShouldItems, MayItems, TasksSummary

The `--context-only` flag MUST output JSON without triggering verification.

## Requirement: Goal-Backward Verification

Verification MUST check each MUST requirement against the implemented code (goal-backward, not task-forward).

The verifier agent MUST produce a structured report with pass/fail per requirement.

`ParseVerifierReport()` MUST parse the agent's JSON output into a typed result.

## Requirement: Gap Reporting

`WriteGapReport()` MUST document any uncovered or failed requirements.

The `--write-results` flag MUST persist the verification report to disk.

## Requirement: Archive

The `mysd archive` command MUST move a verified change to `openspec/changes/archive/`.

Archive MUST only proceed if the change has passed verification (state = Verified).

Archive MUST update `tracking.yaml` with CompletedAt timestamp.

Archive MUST delete `discuss-research-cache.json` if present (best-effort, silent on failure).

## Requirement: UAT Checklist

The `uat` package MUST manage User Acceptance Testing checklists at `.mysd/uat/{change}-uat.md`.

Each `UATItem` MUST have: ID, Description, Status (pass/fail/skip), Notes, RunAt.

`UATSummary` MUST compute: Total, Pass, Fail, Skip counts.

### Scenario: All MUST Requirements Pass

WHEN the verifier report marks all MUST items as passed
THEN state transitions to "verified"
AND `mysd archive` is unblocked

### Scenario: MUST Requirement Fails

WHEN the verifier report marks a MUST item as failed
THEN state remains "executed"
AND a gap report is generated

### Scenario: Archive Lifecycle

WHEN `mysd archive` completes
THEN the change directory moves to `openspec/changes/archive/`
AND tracking.yaml records the completion timestamp
AND research cache is cleaned up

## Covered Packages

- `cmd/verify.go`, `cmd/archive.go`
- `internal/verifier/` — VerificationContext building, report parsing, gap reporting
- `internal/uat/` — UAT checklist management
