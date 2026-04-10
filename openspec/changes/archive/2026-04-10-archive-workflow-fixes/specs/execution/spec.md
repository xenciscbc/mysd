---
spec-version: "1.0"
capability: execution
delta: MODIFIED
status: pending
---

## MODIFIED Requirements

### Requirement: Task Status Updates

The `task_update` helper MUST update task status in `tasks.md` (pending → in_progress → done/blocked).

Status updates MUST preserve the YAML frontmatter and markdown structure.

After updating a task's status, `task_update` SHALL check whether all tasks in the change are in a terminal state (done or skipped). If all tasks are terminal AND the current workflow phase is `planned`, `task_update` SHALL automatically transition the phase to `executed`.

The auto-transition SHALL only fire on the `planned → executed` edge. If the phase is already `executed` or later, no transition SHALL occur.

`task_update` SHALL print a message when auto-transitioning: "All tasks complete — phase advanced to executed".

#### Scenario: Auto-transition on last task completion

- **WHEN** `mysd task-update 3 done` is called
- **AND** task 3 is the last pending task (all others are done or skipped)
- **AND** the current phase is `planned`
- **THEN** the phase SHALL transition to `executed`
- **AND** stdout SHALL include "All tasks complete — phase advanced to executed"

#### Scenario: No auto-transition when tasks remain

- **WHEN** `mysd task-update 2 done` is called
- **AND** task 3 is still pending
- **THEN** the phase SHALL remain `planned`

#### Scenario: No auto-transition when already executed

- **WHEN** `mysd task-update 3 done` is called
- **AND** all tasks are terminal
- **AND** the current phase is `executed`
- **THEN** no state transition SHALL occur (idempotent)
