## MODIFIED Requirements

### Requirement: Execution Context

The `mysd execute` command MUST build an `ExecutionContext` JSON for executor agents containing:
- ChangeName, MustItems, ShouldItems, MayItems
- Tasks, PendingTasks, WaveGroups
- AtomicCommits flag, ExecutionMode, AgentCount

The `--context-only` flag MUST output the JSON without triggering execution.

#### Scenario: Single Agent Execution

- **WHEN** execute is called with agent-count=1
- **THEN** tasks are executed sequentially in dependency order

#### Scenario: Wave Parallel Execution

- **WHEN** execute is called with agent-count > 1
- **THEN** independent tasks are grouped into waves
- **AND** each wave's tasks run in parallel across agents

### Requirement: Apply command verification is mandatory

The `/mysd:apply` command SHALL always run spec verification after task execution completes successfully. The verification step SHALL NOT be skippable by user interaction.

In auto mode, verification SHALL proceed without confirmation. In interactive mode, verification SHALL also proceed without confirmation — the user prompt asking whether to run verification SHALL be removed.

#### Scenario: Apply runs verification automatically

- **WHEN** `/mysd:apply` completes all tasks and build+tests pass
- **THEN** the verifier agent SHALL be invoked automatically without asking the user

#### Scenario: Apply skips verification only on build failure

- **WHEN** `/mysd:apply` completes tasks but `go build` or `go test` fails
- **THEN** verification SHALL be skipped
- **AND** the user SHALL be informed to run `/mysd:fix`

## REMOVED Requirements

### Requirement: Execute command skill

**Reason**: The `/mysd:execute` command skill has been renamed to `/mysd:apply`. The redirect stub is no longer needed as all references have been updated.

**Migration**: Use `/mysd:apply` for task execution. All agent and command references updated to use the new name.

#### Scenario: Execute command file removed

- **WHEN** a user invokes `/mysd:execute`
- **THEN** the command skill file SHALL NOT exist in `plugin/commands/`

### Requirement: Spec command skill

**Reason**: Spec writing is now embedded within `/mysd:propose` (Step 11) and `/mysd:discuss` (Step 10). A standalone `/mysd:spec` command is redundant.

**Migration**: Use `/mysd:propose` for initial spec generation or `/mysd:discuss` for spec refinement.

#### Scenario: Spec command file removed

- **WHEN** a user looks for `/mysd:spec`
- **THEN** the command skill file SHALL NOT exist in `plugin/commands/`

### Requirement: Design command skill

**Reason**: Design document creation is now embedded within `/mysd:plan` (Step 4). A standalone `/mysd:design` command is redundant.

**Migration**: Use `/mysd:plan` which includes the design phase automatically.

#### Scenario: Design command file removed

- **WHEN** a user looks for `/mysd:design`
- **THEN** the command skill file SHALL NOT exist in `plugin/commands/`

### Requirement: Capture command skill

**Reason**: Conversation capture is fully covered by `/mysd:discuss` which provides the same functionality plus research, gray area exploration, and spec updates.

**Migration**: Use `/mysd:discuss` to capture conversation context into structured proposals and specs.

#### Scenario: Capture command file removed

- **WHEN** a user looks for `/mysd:capture`
- **THEN** the command skill file SHALL NOT exist in `plugin/commands/`
