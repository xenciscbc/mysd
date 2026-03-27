## MODIFIED Requirements

### Requirement: Fast-Forward (ff)

The `mysd ff` command MUST orchestrate the sequence: plan (including design) → apply (including verify) → archive (including docs) in a single invocation.

The command MUST skip interactive confirmations (implies --auto).

The command MUST accept a change name as argument.

If any stage fails, the command MUST stop and report the failure point.

All references within the ff command skill and its spawned agents SHALL use `/mysd:apply` instead of `/mysd:execute`.

#### Scenario: Fast-Forward Happy Path

- **WHEN** `mysd ff my-change` runs with a clear spec
- **THEN** design, plan, apply, verify, and archive are executed without user interaction
- **AND** state progresses to "archived"

### Requirement: Extended Fast-Forward (ffe)

The `mysd ffe` command MUST orchestrate the full sequence: research → plan (including design) → apply (including verify) → archive (including docs) in a single invocation.

The command MUST skip interactive confirmations (implies --auto).

The command MUST accept a change name as argument.

If any stage fails, the command MUST stop and report the failure point.

All references within the ffe command skill and its spawned agents SHALL use `/mysd:apply` instead of `/mysd:execute`.

#### Scenario: Extended Fast-Forward Failure

- **WHEN** `mysd ffe my-change` fails during execution
- **THEN** the command stops at the failed stage
- **AND** state reflects the last successful phase
