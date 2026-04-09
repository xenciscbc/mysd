---
description: Plan orchestrator. Optional single-agent research, then design, then task planning. Usage: /mysd:plan [--research] [--check] [--auto]
argument-hint: "[--research] [--check] [--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
  - AskUserQuestion
---

# /mysd:plan -- Create Executable Task List

You are the mysd plan orchestrator. Your job is to run the planning pipeline: optional research, then design, then task planning.

## Question Protocol

- Ask one question at a time. Wait for the user's answer before asking the next.
- When a question has concrete options, use the **AskUserQuestion tool** — do not list options as plain text.
- Open-ended questions may use plain text.

## Step 1: Parse Arguments

Check `$ARGUMENTS` for flags:
- `--research`: Enable focused research before design
- `--check`: Enable plan-checker validation after planning
- `--auto`: Auto mode — skip interactive prompts, use AI recommendations

Set `auto_mode` = true if `--auto` is present, false otherwise.

## Step 2: Get Planning Context

Run:
```
mysd plan --context-only [--research] [--check]
```

Parse the JSON output. It contains:
- `spec_dir`, `change_name`, `phase`, `specs`, `design`, `model`
- `reviewer_model`, `plan_checker_model`
- `research_enabled`, `check_enabled`, `test_generation`

Extract the following fields:
- `spec_dir`: the detected spec directory (`.specs` or `openspec`) — pass to all agents
- `model`: profile-resolved model for designer and planner agents
- `reviewer_model`: profile-resolved model for mysd-reviewer
- `plan_checker_model`: profile-resolved model for mysd-plan-checker
- `spec`: (optional) the `--spec` value if passed, restricts planning to this spec
- `external_input`: (optional) content from `--from` flag, used as planner context

If error (not in designed/specced phase), guide user to complete prerequisites.

## Step 2b: Per-Spec Selection

If `--spec` was NOT passed in the original command:

1. Read the change's specs directory to list all capability names
2. Read tasks.md (if exists) and group existing tasks by `spec` field
3. Identify which specs have no corresponding tasks yet
4. If `auto_mode` is true: plan all specs (skip selection)
5. If `auto_mode` is false: present an interactive selection list:
   ```
   Specs to plan:
   1. material-selection (no tasks)
   2. planning (5 tasks)
   3. execution (3 tasks)
   4. [All]
   5. [From conversation context]

   Select:
   ```
6. If user selects a specific spec: re-run `mysd plan --spec {selected} --context-only` to get filtered context
7. If user selects "All": proceed with full context (all specs)
8. If user selects "From conversation context":
   a. Extract relevant requirements, task descriptions, and design decisions from the current conversation history
   b. Write the extracted content to `{changeDir}/conversation-context.md`
   c. Re-run `mysd plan --from {changeDir}/conversation-context.md --context-only` to get context with external input
   d. Proceed with the enriched context

If `--spec` WAS passed: skip this step (context already filtered).
If `--from` WAS passed: skip this step (external input already provided).

## Step 2c: External Input Display

If `external_input` is present in context JSON:
- Show: "External input loaded from {path}: {first line or title}"
- This content will be passed to the planner as additional context.

## Step 3: Research Phase (if research_enabled)

If `research_enabled` is true (from context JSON) or `--research` flag present:

  If `auto_mode` is false:
    Ask: "Would you like to run focused research on implementation details? [y/N]"
    If user declines: skip to Step 4.

  If `auto_mode` is true:
    Skip research entirely. Go to Step 4.

  Show: "Spawning mysd-researcher ({model})..."
  Spawn ONE `mysd-researcher` agent (single, NOT parallel), with `model` parameter set to `{model}`:

  Task: Research implementation details for {change_name}
  Agent: mysd-researcher
  Model: {model}
  Context: {
    "spec_dir": "{spec_dir}",
    "change_name": "{change_name}",
    "dimension": "architecture",
    "topic": "implementation of {change_name} — validate technical feasibility and supplement implementation details",
    "spec_files": [{all spec file paths from context + design.md path}],
    "auto_mode": {auto_mode}
  }

  Collect research output. This becomes additional input for the designer in Step 4.

  If `auto_mode` is false:
    Present research summary and ask: "Research complete. Proceed to design? (Y/n)"

## Step 3b: Evaluate Design Skip

Before the design phase, evaluate whether design.md is needed.

Design MAY be skipped when ALL of the following conditions are met:
1. Proposal Impact section lists 2 or fewer affected files
2. Proposal has no New Capabilities (only Modified Capabilities or none)
3. Proposal does not contain keywords: "cross-cutting", "migration", "architecture", "new pattern"

Read `proposal.md` and check these conditions.

If conditions are met:
- If `auto_mode` is true: skip design, show "⊘ Skipped design (not needed for this change)", go to Step 5.
- If `auto_mode` is false: show the assessment and ask: "Design appears optional for this change. Skip design? [Y/n]"
  - If user confirms skip: go to Step 5.
  - If user declines: proceed to Step 4.

If any condition is NOT met: proceed to Step 4.

## Step 4: Design Phase

First, get structured instructions for the designer:
```
mysd instructions design --change {change_name} --json
```
Parse the JSON output to get `template`, `rules`, `instruction`, `selfReviewChecklist`, and `dependencies`.

Show: "Spawning mysd-designer ({model})..."
Use the Task tool to invoke `mysd-designer` with `model` parameter set to `{model}`:

Task: Create design document for {change_name}
Agent: mysd-designer
Model: {model}
Context: {
  "spec_dir": "{spec_dir}",
  "change_name": "{change_name}",
  "specs": [{spec content}],
  "research_findings": [{from Step 3, or empty if no research}],
  "instructions": {instructions JSON from mysd instructions},
  "auto_mode": {auto_mode}
}

The designer produces `design.md`.

After designer completes, run state transition:
```
mysd design
```

## Step 5: Planning Phase

First, get structured instructions for the planner:
```
mysd instructions tasks --change {change_name} --json
```
Parse the JSON output to get `template`, `rules`, `instruction`, `selfReviewChecklist`, and `dependencies`.

Show: "Spawning mysd-planner ({model})..."
Use the Task tool to invoke `mysd-planner` with `model` parameter set to `{model}`:

Task: Create task list for {change_name}
Agent: mysd-planner
Model: {model}
Context: {full context JSON from Step 2 (including spec_dir), plus research_findings and design content, plus auto_mode, plus:
  "spec_dir": "{spec_dir}",
  "instructions": {instructions JSON from mysd instructions},
  "external_input": {external_input from context, if present},
  "target_spec": {spec name from Step 2b selection, if per-spec planning}
}

After planner completes, run state transition:
```
mysd plan
```

## Step 5a: Inline Self-Review

After the planner completes and before the reviewer, perform an inline quality check on the produced artifacts. The orchestrator executes this step directly — no agent is spawned.

First, load the self-review checklist:
```
mysd instructions tasks --change {change_name} --json
```
Use the `selfReviewChecklist` field as a guide. Then execute these 4 checks in order:

### Check 1: Placeholder Scan

Read `tasks.md` and `design.md`. Scan for:
- Literal strings: "TBD", "TODO", "FIXME", "implement later", "details to follow"
- Empty template sections (e.g., `<!-- ... -->` placeholders left unfilled)

For each occurrence:
1. Read the relevant spec to find the concrete information
2. Use the Edit tool to replace the placeholder with specific content
3. Count fixes applied

Show: "Placeholder check: fixed {N} placeholder(s)" (or "Placeholder check: clean ✓")

### Check 2: Consistency Check

Read `proposal.md` to extract capability names from the Capabilities section.

For each capability listed in the proposal:
1. Check if at least one task in `tasks.md` has `spec: "{capability-name}"` (or references the capability in its description)
2. If a capability has no corresponding task, add a task to cover it

Read file paths mentioned in `tasks.md` task descriptions. Verify each path appears in either:
- Proposal's Impact section, OR
- Design document

Flag and fix any mismatches.

Show: "Consistency check: fixed {N} mismatch(es)" (or "Consistency check: clean ✓")

### Check 3: Scope Check (warning only)

Count total tasks in `tasks.md`:
- If total > 15: show warning "⚠ {N} tasks exceed recommended maximum of 15 — consider splitting the change"

For each task, count the number of distinct file references in its description:
- If any task references > 3 files: show warning "⚠ Task T{id} references {N} files — consider splitting"

Do NOT auto-fix scope issues — display warnings only and proceed.

### Check 4: Ambiguity Check

Scan task descriptions for vague phrases:
- "handle edge cases"
- "add appropriate error handling"
- "implement the flow"
- "add tests for the above"
- "similar to Task N" (without repeating specifics)

For each vague phrase:
1. Read the relevant spec to find the specific conditions
2. Use the Edit tool to replace the vague phrase with concrete details
3. Count fixes applied

Show: "Ambiguity check: fixed {N} vague phrase(s)" (or "Ambiguity check: clean ✓")

### Summary

After all 4 checks, show a one-line summary:
```
Self-review: {total_fixes} fix(es), {total_warnings} warning(s)
```

Proceed to the reviewer step.

## Step 5b: Invoke Reviewer

Run:
```
mysd validate {change_name}
```
Capture output as `validate_output` (empty string if command not found or fails).

Show: "Spawning mysd-reviewer ({reviewer_model})..."
Use the Task tool to invoke `mysd-reviewer` with `model` parameter set to `{reviewer_model}`:

Task: Review artifacts for {change_name}
Agent: mysd-reviewer
Model: {reviewer_model}
Context: {
  "spec_dir": "{spec_dir}",
  "change_name": "{change_name}",
  "phase": "plan",
  "validate_output": "{validate_output}",
  "auto_mode": {auto_mode}
}

Collect the reviewer summary. Include it in Step 7 output.

## Step 5c: Analyze-Fix Loop

Run cross-artifact structural analysis and fix Critical/Warning findings (max 2 iterations).

### Iteration Loop (max 2):

1. Run:
   ```
   mysd analyze {change_name} --json
   ```

2. Parse the JSON output. Filter findings to **Critical and Warning severity only** (ignore Suggestion).

3. If no Critical/Warning findings: show "Artifacts look consistent ✓" and proceed to Step 6.

4. If Critical/Warning findings exist:
   - Show: "Found N issue(s), fixing... (attempt M/2)"
   - For each finding: read the affected artifact, apply the recommended fix using the Edit tool
   - Re-run `mysd analyze {change_name} --json`
   - If still has findings and iteration < 2: repeat
   - If iteration reaches 2 and findings remain: show remaining findings as a summary, proceed to Step 6 (do NOT block)

## Step 6: Plan Check (if check_enabled)

If `check_enabled` is true:

Run:
```
mysd plan --check --context-only
```

Show: "Spawning mysd-plan-checker ({plan_checker_model})..."
Use the Task tool to invoke `mysd-plan-checker` with `model` parameter set to `{plan_checker_model}`:

Task: Validate plan coverage for {change_name}
Agent: mysd-plan-checker
Model: {model}
Context: {check output JSON}

## Step 7: Cleanup

Delete temporary files created during the plan pipeline:
- If `conversation-context.md` exists in the change directory, delete it (best-effort, ignore errors)

## Step 8: Confirm

Show:
1. "Planning complete. Pipeline: {research if enabled} -> design {or skipped} -> plan -> reviewer -> analyze {-> check if enabled}."
2. Reviewer summary from Step 5b (issues fixed, cannot-auto-fix items if any)
3. Next: `/mysd:apply`
