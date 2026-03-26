---
model: claude-sonnet-4-5
description: Goal-backward verification of all MUST items. Invokes independent verifier agent.
argument-hint: ""
allowed-tools:
  - Bash
  - Read
  - Write
  - Task
---

# /mysd:verify — Verify Spec Change

You are the mysd verify orchestrator. Your job is to gather verification context and invoke the independent verifier agent, then write the results back.

## Step 1: Get Verification Context

Run:
```
mysd verify --context-only
```

Parse the JSON output. It contains:
- `change_name`: The current change
- `change_dir`: Path to `.specs/changes/{change_name}/`
- `specs_dir`: Path to the specs directory
- `must_items`: Array of MUST requirements (id, text, keyword, source_file)
- `should_items`: Array of SHOULD requirements (id, text, keyword, source_file)
- `may_items`: Array of MAY requirements (id, text, keyword, source_file)
- `tasks_summary`: Array of tasks with their current status (id, name, status)

If this returns an error such as "not in executed phase", guide the user to complete `/mysd:execute` first before verifying.

If `must_items` is empty, inform the user: "No MUST requirements found. Add MUST requirements to your spec files before verifying."

## Step 2: Invoke Independent Verifier Agent

Use the Task tool to invoke the mysd-verifier agent with the full context:

```
Task: Invoke mysd-verifier agent for independent spec verification
Agent: mysd-verifier
Context: {full context JSON from Step 1}
```

The verifier agent will:
1. Read all spec files independently (NOT alignment.md)
2. Find concrete filesystem evidence for each MUST item
3. Run tests and check builds as needed
4. Write the report to `{change_dir}/verifier-report.json`

Wait for the verifier agent to complete before proceeding.

## Step 3: Write Verification Results

After the verifier agent completes and writes `verifier-report.json`, run:

```
mysd verify --write-results {change_dir}/verifier-report.json
```

This command:
- Reads the verifier report
- Writes `{change_dir}/verification.md` (full report)
- Writes `{change_dir}/gap-report.md` if any MUST items failed (with fix suggestions)
- Updates `{change_dir}/verification-status.json` sidecar
- Transitions state from `executed` to `verified` if all MUST items pass

## Step 4: Report Results to User

After Step 3 completes, present a clear summary:

**If all MUST items pass (overall_pass == true):**
```
Verification PASSED for change: {change_name}

MUST items: {count} / {count} passed
SHOULD items: {should_pass_count} / {should_total} passed (warnings only)

Full report written to: {change_dir}/verification.md

Next step: Run `/mysd:archive` to archive this verified change.
```

If `has_ui_items` is true, also show:
```
UAT checklist generated: .mysd/uat/{change_name}-uat.md
Run `/mysd:uat` to conduct interactive user acceptance testing before archiving.
```

**If any MUST items fail (overall_pass == false):**
```
Verification FAILED for change: {change_name}

MUST items: {must_pass_count} / {must_total} passed
Failed items:
  - {id}: {text}
    Evidence: {evidence}
    Suggestion: {suggestion}

Gap report written to: {change_dir}/gap-report.md

Next step: Run `/mysd:execute` to fix the failing items, then re-run `/mysd:verify`.
```
