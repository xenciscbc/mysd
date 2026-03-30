---
description: Create a new spec change with proposal scaffolding. Supports material selection, requirement interview, 4-dimension research, gray area exploration, and scope guardrail. Usage: /mysd:propose [change-name|file-path|dir-path] [--auto]
argument-hint: "[change-name|file|dir] [--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
  - AskUserQuestion
---

# /mysd:propose — Create a New Change Proposal

You are the mysd propose orchestrator. Your job is to collect and understand requirements through material selection and interview, scaffold a new change, run optional 4-dimension research, facilitate gray area exploration, and invoke the proposal writer agent.

## Step 1: Parse Arguments

Check `$ARGUMENTS` for `--auto`. Remove it from the arguments list.
Set `auto_mode` = true if `--auto` is present, false otherwise.

The remaining arguments (after removing `--auto`) are the `source_arg`.

If `source_arg` matches an existing directory `.specs/changes/{source_arg}/`:
→ Set `existing_change` = true, `change_name` = `source_arg`
→ Read existing proposal.md if present

Otherwise: set `existing_change` = false, `change_name` = null (will be derived later in Step 7).

## Step 2: Resolve Agent Model

Run:
```
mysd model
```

Parse the output to find `Profile: {profile_name}`. The profile determines agent models:

| Role | quality | balanced | budget |
|------|---------|----------|--------|
| researcher / advisor | `sonnet` | `sonnet` | `haiku` |
| proposal-writer / spec-writer | `sonnet` | `sonnet` | `sonnet` |
| reviewer | `opus` | `sonnet` | `sonnet` |

Set `model` for researcher/advisor/proposal-writer/spec-writer, and `reviewer_model` for the reviewer agent.

## Step 3: Load Deferred Notes

Per D-02 from prior design: propose ALWAYS loads deferred notes (cross-change context is valuable for new proposals).

Run:
```bash
mysd note list
```

- If output shows notes: store as `deferred_context` and record as a detected source for Step 4
- If no notes or command returns empty: set `deferred_context` to empty string

## Step 4: Material Selection

Detect all available requirement sources and let the user choose which to use.

### Step 4a: Detect Sources

Apply the following detection logic for each source type:

| # | Source Type | Detection Method |
|---|-----------|-----------------|
| 1 | source_arg file/directory | Check if `source_arg` is a valid file path (ends with `.md` or exists on disk) or a valid directory (exists on disk with `.md` files). Skip if `source_arg` was already matched as an existing change in Step 1. |
| 2 | Conversation context | Check if the current conversation contains substantive requirement discussion (not just greetings or meta-talk). Look for problem descriptions, feature requests, or design discussions. |
| 3 | Claude plan | Check if conversation system messages contain a plan file path matching `~/.claude/plans/<name>.md`. If found, verify the file exists on disk. |
| 4 | gstack plan | Scan `~/.gstack/projects/{project}/` for `.md` files. The `{project}` is derived from the current working directory name. |
| 5 | Active change | Run `mysd status`. If it reports an active change with an existing `proposal.md`, record it. |
| 6 | Deferred notes | Use the result from Step 3. If `deferred_context` is non-empty, record it as a source. |

For each detected source, extract a brief content preview (first line, title, or summary).

### Step 4b: Present Sources and Collect Selection

**If `auto_mode` is true:** Automatically aggregate content from all detected sources. If no sources detected, extract requirements from conversation context as best-effort. Go to Step 4c.

**If `auto_mode` is false:**

If no sources are detected from any of the 6 types:
→ Skip the list, go directly to manual input: ask "What would you like to change or build? Please describe the goal."
→ Use the user's description as `aggregated_content`. Go to Step 5.

If sources are detected:
→ Present a numbered list of detected sources with type labels and content previews.
→ Always include "Manual input" as the last option.
→ Allow multi-selection (e.g., "1,3" to select sources 1 and 3).
→ Even if only one source is detected, still present the list for user confirmation.

Example display:
```
Detected requirement sources:
1. [Conversation] Discussion about improving propose workflow depth
2. [Claude Plan] ~/.claude/plans/enhance-propose.md
3. [Deferred Notes] 2 notes from prior changes
4. [Manual input] Describe your requirement

Select sources to use (comma-separated, e.g., 1,3):
```

Wait for the user's selection.

If user selects "Manual input":
→ Ask "What would you like to change or build? Please describe the goal."
→ Use the user's description as additional content.

### Step 4c: Aggregate Content

Read and combine the content from all selected sources into a single `aggregated_content` string.

For each source type, read as follows:
- **source_arg file**: Read the file content
- **source_arg directory**: If `auto_mode`, read all `.md` files; otherwise present list for multi-select
- **Conversation context**: Extract the substantive discussion portions
- **Claude plan**: Read the plan file, extract H1 heading as title, Context section as motivation, implementation stages as structure
- **gstack plan**: Read selected `.md` files from the project directory
- **Active change**: Read the existing `proposal.md`
- **Deferred notes**: Include the `deferred_context` string

## Step 5: Scan Existing Specs

Use Glob to list `openspec/specs/*/spec.md`. Extract directory names as spec identifiers.

Compare against the `aggregated_content` to identify related specs (max 5 candidates). For each candidate (max 3), read the first 10 lines to retrieve the Purpose section.

Store the related spec names and their content as `related_specs` — this will be passed to the interview step.

If related specs are found, display as an informational note — do NOT stop or ask for confirmation, continue automatically.
If none found, silently proceed.

## Step 6: Requirement Interview

Evaluate the completeness of `aggregated_content` and ask clarifying questions as needed. The orchestrator performs this step directly — no agent is spawned.

### Step 6a: Evaluate Completeness

Assess `aggregated_content` against three dimensions:

| Dimension | Complete when |
|-----------|--------------|
| **Problem** | The content describes the problem being solved, not just a desired solution |
| **Boundary** | The content states what is in scope and what is explicitly excluded |
| **Success Criteria** | The content includes specific, verifiable conditions for success |

Additionally, check `related_specs` from Step 5. If the proposed change overlaps with an existing spec (similar capability area, same subsystem), flag it as a question to ask.

### Step 6b: Interview Loop

**If `auto_mode` is true:** Skip the interview entirely. Produce `requirement_brief` using best-effort inference from `aggregated_content`. Dimensions that cannot be inferred are filled with reasonable defaults based on available context — never leave them empty. Go to Step 6c.

**If `auto_mode` is false:**

If all three dimensions are sufficiently covered AND no spec overlap is detected:
→ Skip the interview, go directly to Step 6c.

Otherwise, for each gap or overlap issue:
1. Ask **one** clarifying question targeting the most important gap:
   - If Problem is missing: ask what problem this change solves (not what the user wants to build)
   - If Boundary is missing: ask what is explicitly out of scope
   - If Success Criteria is missing: ask how the user will know the change is successful
   - If spec overlap is detected: ask "This overlaps with the existing `{spec_name}` spec. Do you want to extend it or create a new capability?"
2. Wait for the user's answer
3. Incorporate the answer into `aggregated_content`
4. Re-evaluate completeness
5. If gaps remain, repeat from (1). If all dimensions are covered, proceed to Step 6c.

Priority order when multiple gaps exist: Problem > Boundary > Success Criteria > Spec overlap.

### Step 6c: Produce requirement_brief

Synthesize a structured `requirement_brief` from the aggregated content and interview answers:

```
## Problem
{The problem being solved}

## Boundary
{What is in scope / what is explicitly excluded}

## Success Criteria
{Specific, verifiable conditions for success}

## Source
{List of source types used, e.g., "Conversation context, Claude plan, Deferred notes"}
```

This is an in-memory intermediate artifact — do NOT write it to disk.

## Step 7: Derive Change Name + Classify Type

**If `existing_change` is true (from Step 1):** Use the existing `change_name`. Skip name derivation.

**Otherwise:** From the `requirement_brief`, auto-derive a kebab-case change name:
- Short (2–4 words), lowercase, hyphen-separated
- Examples: "add dark mode" → `add-dark-mode`, "fix login crash" → `fix-login-crash`

Set `change_name` to the derived name. Do not ask the user to confirm or choose unless the name is ambiguous.

Classify into one of:

| Type | When to use |
|------|-------------|
| Feature | New functionality, new capabilities |
| Bug Fix | Fixing existing behavior, resolving errors |
| Refactor | Architecture improvements, performance, reorganization |

Set `change_type`.

## Step 8: Scaffold the Change

**If `existing_change` is true:** Skip scaffolding (directory already exists).

**Otherwise:** Run:
```
mysd propose {change-name}
```

This creates `.specs/changes/{change-name}/` with a template `proposal.md`.

## Step 9: Optional Research

If `auto_mode` is true: skip research entirely. Go directly to Step 12.

If `auto_mode` is false: Ask user:
```
Would you like to run 4-dimension research on this proposal?
(Codebase / Domain / Architecture / Pitfalls) [y/N]
```

- If user declines: go to Step 12.
- If user accepts: proceed to Step 9a.

### Step 9a: Parallel Research Spawning

Show: "Spawning 4 mysd-researcher agents ({model})..."
Spawn 4 `mysd-researcher` agents in parallel using the Task tool, each with `model` parameter set to `{model}`:

For each dimension in ["codebase", "domain", "architecture", "pitfalls"]:
```
Task: Research {dimension} for proposal: {change_name}
Agent: mysd-researcher
Model: {model}
Context: {
  "change_name": "{change_name}",
  "dimension": "{dimension}",
  "topic": "{requirement_brief}",
  "spec_files": [".specs/changes/{change_name}/proposal.md"],
  "auto_mode": false
}
```

Collect all 4 research outputs. Present organized summary by dimension to the user.

Then proceed to Step 10.

## Step 10: Gray Area Identification + Advisor Spawning (research only)

This step is ONLY executed if the user accepted research in Step 9.

From the 4 research outputs, identify gray areas: ambiguous design decisions where multiple valid approaches exist, conflicting recommendations between dimensions, or areas needing user input.

For each gray area, show: "Spawning mysd-advisor ({model})..." and spawn one `mysd-advisor` agent in parallel using the Task tool with `model` parameter set to `{model}`:
```
Task: Analyze gray area: {gray_area_description}
Agent: mysd-advisor
Model: {model}
Context: {
  "change_name": "{change_name}",
  "gray_area": "{gray_area_description}",
  "research_findings": "{all 4 researcher outputs combined}",
  "auto_mode": false
}
```

CRITICAL: Advisors MUST be spawned at this orchestrator layer, NOT inside any researcher agent.

Collect all advisor comparison tables.

## Step 11: Dual-Loop Exploration (research only)

This step is ONLY executed if the user accepted research in Step 9.

### Layer 1 — Per-Area Deep Dive

For each gray area with its advisor analysis:

1. Present the advisor's comparison table
2. Facilitate discussion:
   - AI presents findings and asks clarifying questions (AI-led)
   - User can answer or ask their own questions (user-led)
   - This is natural conversation flow — no explicit mode switch needed
3. **Scope Guardrail:** During discussion, if a suggestion expands beyond the current proposal scope:
   - Acknowledge the idea
   - State: "This is outside the current proposal scope."
   - Run: `mysd note add "{idea summary}"` to save to deferred notes
   - Continue exploration without incorporating the out-of-scope idea
   - Scope boundary is determined by reading the proposal.md's **In Scope / Out of Scope** sections (or the requirement_brief's Boundary section)
4. After the area discussion concludes, ask:
   ```
   This area is resolved. Would you like to:
   1. Continue to the next area
   2. Finish exploration
   ```
   If user chooses "Finish exploration": exit Layer 1 and go directly to Step 12.

### Layer 2 — New Area Discovery

After all identified gray areas from Step 10 are explored:
```
All identified areas have been explored.
Would you like to:
1. Explore additional areas (describe what you'd like to investigate)
2. Finish exploration and proceed to proposal writing
```

If user chooses "Explore additional areas":
- User describes new areas to investigate
- Spawn one `mysd-advisor` agent per new area (same pattern as Step 10)
- Run Layer 1 deep dive for each new area

If user chooses "Finish exploration": proceed to Step 12.

## Step 12: Invoke Proposal Writer

Show: "Spawning mysd-proposal-writer ({model})..."
Use the Task tool to invoke `mysd-proposal-writer` with `model` parameter set to `{model}`:

```
Task: Write proposal for {change_name}
Agent: mysd-proposal-writer
Model: {model}
Context: {
  "change_name": "{change_name}",
  "change_type": "{change_type}",
  "conclusions": "{requirement_brief + research findings + exploration conclusions (if any)}",
  "existing_proposal": "{current proposal.md body if exists, else null}",
  "deferred_context": "{deferred notes from Step 3}",
  "auto_mode": {auto_mode}
}
```

The proposal writer will fill in the proposal.md with structured content based on the requirement_brief, research findings, and exploration conclusions.

## Step 13: Auto-Invoke Spec Writer

Automatically invoke the spec-writer agent to generate specs from the proposal.

Read the proposal body:
```
Read .specs/changes/{change_name}/proposal.md
```

Read existing spec files (if any):
```
ls .specs/changes/{change_name}/specs/
```

For each capability area found in the proposal (or a single "core" area if not structured by capability):

Show: "Spawning mysd-spec-writer ({model})..."
Use the Task tool to invoke `mysd-spec-writer` with `model` parameter set to `{model}`:
```
Task: Generate specs for {change_name} — {capability_area}
Agent: mysd-spec-writer
Model: {model}
Context: {
  "change_name": "{change_name}",
  "capability_area": "{capability_area}",
  "existing_spec_body": "{existing spec content if any, else null}",
  "proposal": "{proposal.md body}",
  "auto_mode": {auto_mode}
}
```

After spec-writer completes, proceed to Step 14.

## Step 14: Artifact Review (mysd-reviewer)

Run validation first, then spawn the reviewer agent to scan and fix quality issues.

### Step 14a: Run Validation

Run:
```
mysd validate {change_name}
```

Capture output as `validate_output` (empty string if command not found or fails).

### Step 14b: Spawn Reviewer

Show: "Spawning mysd-reviewer ({reviewer_model})..."
Use the Task tool to invoke `mysd-reviewer` with `model` parameter set to `{reviewer_model}`:

```
Task: Review artifacts for {change_name}
Agent: mysd-reviewer
Model: {reviewer_model}
Context: {
  "change_name": "{change_name}",
  "phase": "propose",
  "change_type": "{change_type}",
  "validate_output": "{validate_output}",
  "auto_mode": {auto_mode}
}
```

Collect the reviewer summary. Include it in Step 15 output.

## Step 15: Final Summary

Show the user:
1. The proposal file path: `.specs/changes/{change_name}/proposal.md`
2. A brief summary of what was written
3. Whether 4-dimension research was performed
4. Number of gray areas explored (if research was run)
5. Number of ideas deferred to notes via scope guardrail (if any)
6. Spec summary: number of MUST / SHOULD / MAY requirements generated
7. Spec file paths created/updated
8. Reviewer result: issues fixed and cannot-auto-fix items (from Step 14b)
9. Next steps:
   - `/mysd:plan` — Create execution plan from specs
   - `/mysd:discuss` — Explore requirements interactively
