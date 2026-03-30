## 1. mysd instructions CLI 指令 (D-01, artifact-instructions spec)

- [x] 1.1 Create `cmd/instructions.go` implementing the "mysd instructions CLI command" requirement. Register `design` and `tasks` as valid artifact IDs with `--change <name> --json` flags. Return error for unknown IDs. Read change state via `state.LoadState()` to determine dependency completion status.
- [x] 1.2 Implement "Self-review checklist in instructions output" for `design` artifact (D-03: Self-Review Checklist): template string (Context / Goals-NonGoals / Decisions / Risks sections), rules array (no placeholders, capability coverage, decision rationale with alternatives), selfReviewChecklist (4 items).
- [x] 1.3 Implement "Self-review checklist in instructions output" for `tasks` artifact (D-03: Self-Review Checklist): template string (TasksFrontmatterV2 YAML + markdown body), rules array (small tasks, spec field required, dependency ordering), selfReviewChecklist (5 items: no placeholders, MUST coverage, max 3 files per task, file path consistency, DAG validation).
- [x] 1.4 Write `cmd/instructions_test.go` with tests for: design output structure, tasks output structure, unknown artifact ID error, dependency done/not-done detection, JSON output format validation.

## 2. TasksFrontmatterV2 新增 spec 欄位 (D-04, "TaskItem includes spec field")

- [x] 2.1 Add `Spec string` field to `executor.TaskItem` struct with `json:"spec,omitempty"` tag per "TaskItem includes spec field" requirement. Update `BuildContextFromParts()` to populate the field from `spec.TaskEntry`.
- [x] 2.2 Add `Spec string` field to `spec.TaskEntry` struct (in `internal/spec/` package) and update the YAML frontmatter parser to read/write the `spec` field in TasksFrontmatterV2 format.
- [x] 2.3 Update `mysd-planner.md` agent definition: add `spec` field to the TasksFrontmatterV2 template example, instruct the planner to assign each task a `spec` value matching the spec directory name.
- [x] 2.4 Write tests for TaskItem.Spec field: YAML parsing with spec present, YAML parsing without spec (empty string), JSON serialization with omitempty behavior.

## 3. Per-spec plan --spec flag (D-05, D-07 external input --from flag, planning spec)

- [x] 3.1 Extend "Task Planning" requirement: add `--spec` and `--from` flags to `cmd/plan.go` (D-07: External input). When `--spec` is provided, filter the plan context to include only requirements from the specified spec. When `--from` is provided, read the file and include content as `external_input` in the context JSON.
- [x] 3.2 Implement task merge logic in `cmd/plan.go` or `internal/spec/`: when `--spec` is used and tasks.md already exists, append new tasks with IDs starting from max existing ID + 1, preserving existing tasks for other specs.
- [x] 3.3 Update `mysd/skills/plan/SKILL.md` Step 2: after getting context, if no `--spec` flag, detect unplanned specs (specs without corresponding tasks), present interactive selection list with spec names + task counts + "All" option. Pass selected spec to subsequent steps.
- [x] 3.4 Update `mysd/skills/plan/SKILL.md` Step 2: when `--from` flag is present (D-07: External input), display the external input source path and inform that it will be passed to the planner as context.
- [x] 3.5 Implement "Plan pipeline uses mysd instructions for agent guidance" in `mysd/skills/plan/SKILL.md` Step 5 (Planning Phase): before spawning mysd-planner, call `mysd instructions tasks --change <name> --json` and include the instructions output in the planner's context. Also pass `external_input` if present.
- [x] 3.6 Implement "Plan pipeline uses mysd instructions for agent guidance" in `mysd/skills/plan/SKILL.md` Step 4 (Design Phase): before spawning mysd-designer, call `mysd instructions design --change <name> --json` and include the instructions output in the designer's context.
- [x] 3.7 Write tests for `--spec` flag: context JSON filtering by spec, task merge with existing tasks (ID continuation, no overwrite), error when spec name doesn't match any spec file.
- [x] 3.8 Write tests for `--from` flag: file content included in context JSON, error when file doesn't exist.

## 4. Per-spec apply --spec flag (D-06, execution spec)

- [x] 4.1 Extend "Execution Context" requirement: add `--spec` flag to `cmd/execute.go`. When provided, filter `PendingTasks` to only tasks matching the specified `spec` field. Recompute `WaveGroups` from filtered tasks.
- [x] 4.2 Update `mysd/skills/apply/SKILL.md` Step 2: after getting context, if no `--spec` flag and `auto_mode` is false, group pending tasks by `spec` field, present interactive selection list with spec names + pending task counts + "All" option.
- [x] 4.3 Write tests for `--spec` execution filtering: pending tasks filtered by spec, wave groups recomputed, change-level tasks (empty spec) excluded from per-spec filter but included in "All".

## 5. Inline Self-Review orchestrator 層 (D-02, D-03 Self-Review Checklist 內嵌於 instructions rules, inline-self-review spec)

- [x] 5.1 Implement "Plan orchestrator inline self-review step" in `mysd/skills/plan/SKILL.md`: add Step 5a between planner completion (Step 5) and reviewer (Step 5b). Step 5a implements "Self-review uses instructions checklist" by calling `mysd instructions tasks --change <name> --json` to load the selfReviewChecklist, then executes 4 checks.
- [x] 5.2 Define placeholder check logic in Step 5a: scan tasks.md and design.md for TBD, TODO, FIXME, "implement later", "details to follow", empty sections. Use Edit tool to replace each with specific content from the spec. Display fix count.
- [x] 5.3 Define consistency check logic in Step 5a: read proposal.md capabilities, verify each has a task with matching `spec` field. Read tasks.md file paths, verify they appear in proposal Impact or design.md. Flag and fix mismatches.
- [x] 5.4 Define scope check logic in Step 5a: count total tasks, warn if > 15. Check each task for file references, warn if any task references > 3 files. Display warnings only (no auto-fix).
- [x] 5.5 Define ambiguity check logic in Step 5a: scan task descriptions for vague phrases ("handle edge cases", "add error handling", "implement the flow"). Replace with specific conditions from the spec using Edit tool.

## 6. Agent Definition Updates

- [x] 6.1 Update `mysd/agents/mysd-designer.md`: add a section instructing the agent to expect and use `instructions` context (template, rules, selfReviewChecklist). Agent SHALL use `template` as output structure and verify against `selfReviewChecklist` before completing.
- [x] 6.2 Update `mysd/agents/mysd-planner.md`: add a section instructing the agent to expect and use `instructions` context. Agent SHALL use `template` for TasksFrontmatterV2 format, follow `rules`, and verify against `selfReviewChecklist`. Agent SHALL assign `spec` field to each task.
