## MODIFIED Requirements

### Requirement: Execution Context

The `mysd execute` command MUST build an `ExecutionContext` JSON for executor agents containing:
- ChangeName, MustItems, ShouldItems, MayItems
- Tasks, PendingTasks, WaveGroups
- AtomicCommits flag, ExecutionMode, AgentCount
- Instruction (dynamically generated guidance string)

The `--context-only` flag MUST output the JSON without triggering execution.

The `--context-only` path SHALL call `runPreflight` internally to obtain a `PreflightReport`, then pass both the `ExecutionContext` and the report to `GenerateInstruction` to populate the `Instruction` field before JSON serialization. The preflight data SHALL NOT appear in the `--context-only` JSON output — it is consumed only by the instruction generator.

#### Scenario: Single Agent Execution

- **WHEN** execute is called with agent-count=1
- **THEN** tasks are executed sequentially in dependency order

#### Scenario: Wave Parallel Execution

- **WHEN** execute is called with agent-count > 1
- **THEN** independent tasks are grouped into waves
- **AND** each wave's tasks run in parallel across agents

#### Scenario: Context-only includes instruction field

- **WHEN** `mysd execute --context-only` is run
- **THEN** the JSON output SHALL contain an `instruction` field
- **AND** the instruction SHALL reflect the current task state and any preflight issues
