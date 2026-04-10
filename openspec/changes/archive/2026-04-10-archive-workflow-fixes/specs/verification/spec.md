---
spec-version: "1.0"
capability: verification
delta: MODIFIED
status: pending
---

## MODIFIED Requirements

### Requirement: Verification Context

The `mysd verify` command MUST build a `VerificationContext` JSON for verifier agents containing:
- ChangeName, MustItems, ShouldItems, MayItems, TasksSummary

The `--context-only` flag MUST output JSON without triggering verification.

Each item in MustItems, ShouldItems, and MayItems MUST include an `id` field containing the stable hash-based identifier (format: `{source_file}::{keyword}-{crc32_hex}`). The verifier agent MUST use these exact IDs in its report — no custom numbering schemes (e.g., `MUST-01`) are permitted.

#### Scenario: Context-only output includes stable IDs

- **WHEN** `mysd verify --context-only` is called
- **THEN** each item in `must_items` SHALL have an `id` field like `spec.md::must-5451802d`
- **AND** the verifier agent report SHALL use the same `id` values verbatim

### Requirement: Goal-Backward Verification

Verification MUST check each MUST requirement against the implemented code (goal-backward, not task-forward).

The verifier agent MUST produce a structured report with pass/fail per requirement.

`ParseVerifierReport()` MUST parse the agent's JSON output into a typed result.

The `--write-results` flag MUST persist the verification report to disk.

Before processing verification results, `--write-results` SHALL check the current phase. If the phase is `planned` and all tasks are in a terminal state (done or skipped), `--write-results` SHALL auto-transition to `executed` before applying the verification outcome. This serves as a safety net for cases where `task-update` auto-transition was missed.

#### Scenario: Write-results auto-advances from planned

- **WHEN** `mysd verify --write-results report.json` is called
- **AND** the current phase is `planned`
- **AND** all tasks are done or skipped
- **THEN** the phase SHALL first transition to `executed`
- **AND** then proceed with normal verification result processing (potentially advancing to `verified`)

#### Scenario: Write-results with phase already executed

- **WHEN** `mysd verify --write-results report.json` is called
- **AND** the current phase is `executed`
- **THEN** no extra transition SHALL occur before processing results
