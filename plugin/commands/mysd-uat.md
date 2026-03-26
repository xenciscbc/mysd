---
model: claude-sonnet-4-5
description: Interactive User Acceptance Testing walkthrough.
argument-hint: ""
allowed-tools:
  - Bash
  - Read
  - Write
  - Task
---

# /mysd:uat — Interactive UAT Walkthrough

You are the mysd UAT orchestrator. Your job is to load the UAT checklist and invoke the interactive UAT guide agent.

## Step 1: Find Current Change and UAT File

Run:
```
mysd status
```

Parse the output to find `change_name`. Then check if the UAT file exists:
```
.mysd/uat/{change_name}-uat.md
```

**If the UAT file does not exist:**
```
No UAT checklist found for change: {change_name}

UAT checklists are generated automatically during verification when UI-related MUST or SHOULD items are detected.

Next step: Run `/mysd:verify` first. If your spec includes UI-related requirements, a UAT checklist will be generated automatically.

If your spec has no UI items, UAT is not required — you can proceed directly to `/mysd:archive`.
```

Stop here if the file does not exist.

**If the UAT file exists:**

Read the file contents:
```
.mysd/uat/{change_name}-uat.md
```

## Step 2: Invoke UAT Guide Agent

Use the Task tool to invoke the mysd-uat-guide agent with the UAT file content:

```
Task: Invoke mysd-uat-guide agent for interactive UAT walkthrough
Agent: mysd-uat-guide
Context: {
  "change_name": "{change_name}",
  "uat_file_path": ".mysd/uat/{change_name}-uat.md",
  "uat_file_content": "{full content of the UAT file}"
}
```

The UAT guide agent will:
1. Walk the user through each UAT item interactively
2. Ask for pass/fail/skip on each item
3. Record notes on failures
4. Save results and update run_history in the UAT file

Wait for the UAT guide agent to complete.

## Step 3: Post-UAT Summary

After the agent completes, summarize the outcome:

**If UAT completed with all items passed:**
```
UAT Complete for: {change_name}

All {count} acceptance tests passed.

Next step: Run `/mysd:archive` to archive this change.
```

**If UAT completed with some failures:**
```
UAT Complete for: {change_name}

Results: {pass_count} passed, {fail_count} failed, {skip_count} skipped

Failed items require attention before archiving. Review the failure notes in:
  .mysd/uat/{change_name}-uat.md

Options:
- Fix the issues and run `/mysd:execute` + `/mysd:verify` again
- Or run `/mysd:archive` to archive anyway (UAT failures are advisory, not blocking)
```

**If UAT was stopped early (user requested early exit):**
```
UAT paused for: {change_name}

Progress has been saved. Run `/mysd:uat` again to continue from where you left off.
```
