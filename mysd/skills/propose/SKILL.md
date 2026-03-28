---
description: Create a new spec change with proposal scaffolding. Supports 4-dimension research, gray area exploration, and scope guardrail. Usage: /mysd:propose [change-name|file-path|dir-path] [--auto]
argument-hint: "[change-name|file|dir] [--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:propose — Create a New Change Proposal

You are the mysd propose orchestrator. Your job is to scaffold a new change, detect the input source, run optional 4-dimension research, facilitate gray area exploration, and invoke the proposal writer agent.

## Step 1: Parse Arguments

Check `$ARGUMENTS` for `--auto`. Remove it from the arguments list.
Set `auto_mode` = true if `--auto` is present, false otherwise.

The remaining arguments (after removing `--auto`) are the `source_arg`.

## Step 2: Source Detection

Apply the following priority order to determine the input source and change name:

**Priority 1:** If `source_arg` matches a directory `.specs/changes/{source_arg}/`
→ Use `source_arg` as the change name (mysd change mode)
→ Read existing proposal.md if present as initial content

**Priority 2:** If `source_arg` is a file path (ends with `.md` or file exists on disk)
→ Single file mode: read the file as initial content
→ Derive change name from filename (strip extension, kebab-case)

**Priority 3:** If `source_arg` is a directory path (directory exists on disk)
→ Selection mode: list all `.md` files in the directory
→ If `auto_mode` is true: use all files as initial content
→ If `auto_mode` is false: present list and let user multi-select

**Priority 4:** If no `source_arg` and there is an active change (check `mysd status` output)
→ Use the current active change

**Priority 5:** If no `source_arg` and no active change → auto-detect from known sources:
→ Check `~/.gstack/projects/{project}/` for `.md` files (design docs, test plans, etc.)
→ Check conversation context for mentioned plan documents or design files
→ Do NOT check `.claude/plans/` (hash filenames have no project info)
→ If `auto_mode` is true: use first detected source
→ If `auto_mode` is false: present detected sources and let user choose

**Priority 6:** If nothing found
→ If `auto_mode` is true: auto-generate change name from conversation context
→ If `auto_mode` is false: ask user for change name and brief description

## Step 3: Scaffold the Change

Run:
```
mysd propose {change-name}
```

This creates `.specs/changes/{change-name}/` with a template `proposal.md`.

If source content was detected in Step 2 (file/directory mode), read that content now.

## Step 3b: Resolve Agent Model

Run:
```
mysd model
```

Parse the output to find `Profile: {profile_name}`. The profile determines agent model:
- `quality` or `balanced` → model = `sonnet`
- `budget` → model = `haiku` (for researcher/advisor); model = `sonnet` (for proposal-writer/spec-writer)

Use this `model` value when spawning agents in subsequent steps.

## Step 4: Load Deferred Notes (D-02)

Per D-02: propose ALWAYS loads deferred notes (cross-change context is valuable for new proposals).

Run:
```bash
mysd note list
```

- If output shows notes: include them as `deferred_context` for research and proposal writing
- If no notes or command returns empty: set `deferred_context` to empty string

## Step 5: Optional Research (DISC-01, DISC-04, D-06)

If `auto_mode` is true: skip research entirely (FAUTO-02 — auto mode means no interaction). Go directly to Step 9.

If `auto_mode` is false: Ask user:
```
Would you like to run 4-dimension research on this proposal?
(Codebase / Domain / Architecture / Pitfalls) [y/N]
```

- If user declines: go to Step 9.
- If user accepts: proceed to Step 6.

## Step 6: Parallel Research Spawning

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
  "topic": "{source content or user description}",
  "spec_files": [".specs/changes/{change_name}/proposal.md"],
  "auto_mode": false
}
```

Collect all 4 research outputs. Present organized summary by dimension to the user.

## Step 7: Gray Area Identification + Advisor Spawning (DISC-06)

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

## Step 8: Dual-Loop Exploration (DISC-05, DISC-07, D-01, D-07, D-08)

### Layer 1 — Per-Area Deep Dive

For each gray area with its advisor analysis:

1. Present the advisor's comparison table
2. Facilitate discussion (DISC-05 dual-mode):
   - AI presents findings and asks clarifying questions (AI-led)
   - User can answer or ask their own questions (user-led)
   - This is natural conversation flow — no explicit mode switch needed
3. **Scope Guardrail (D-08):** During discussion, if a suggestion expands beyond the current proposal scope:
   - Acknowledge the idea
   - State: "This is outside the current proposal scope."
   - Run: `mysd note add "{idea summary}"` to save to deferred notes
   - Continue exploration without incorporating the out-of-scope idea
   - Scope boundary is determined by reading the proposal.md's **In Scope / Out of Scope** sections
4. After the area discussion concludes, ask (D-01 — user-driven, no quota):
   ```
   This area is resolved. Would you like to:
   1. Continue to the next area
   2. Finish exploration
   ```
   If user chooses "Finish exploration": exit Layer 1 and go directly to Step 9.

### Layer 2 — New Area Discovery

After all identified gray areas from Step 7 are explored:
```
All identified areas have been explored.
Would you like to:
1. Explore additional areas (describe what you'd like to investigate)
2. Finish exploration and proceed to proposal writing
```

If user chooses "Explore additional areas":
- User describes new areas to investigate
- Spawn one `mysd-advisor` agent per new area (same pattern as Step 7)
- Run Layer 1 deep dive for each new area

If user chooses "Finish exploration": proceed to Step 9.

## Step 9: Invoke Proposal Writer

Show: "Spawning mysd-proposal-writer ({model})..."
Use the Task tool to invoke `mysd-proposal-writer` with `model` parameter set to `{model}`:

```
Task: Write proposal for {change_name}
Agent: mysd-proposal-writer
Model: {model}
Context: {
  "change_name": "{change_name}",
  "conclusions": "{research findings + exploration conclusions + source content}",
  "existing_proposal": "{current proposal.md body if exists, else null}",
  "deferred_context": "{deferred notes from Step 4}",
  "auto_mode": {auto_mode}
}
```

The proposal writer will fill in the proposal.md with structured content based on the source material, research findings, and exploration conclusions.

## Step 10: Confirm

Show the user:
1. The proposal file path: `.specs/changes/{change_name}/proposal.md`
2. A brief summary of what was written
3. Whether 4-dimension research was performed
4. Number of gray areas explored (if research was run)
5. Number of ideas deferred to notes via scope guardrail (if any)
6. Proceeding to spec generation...

## Step 11: Auto-Invoke Spec Writer (D-01, D-04)

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

After spec-writer completes, show:
1. Spec summary: number of MUST / SHOULD / MAY requirements generated
2. Spec file paths created/updated
3. Next steps:
   - `/mysd:plan` — Create execution plan from specs
   - `/mysd:discuss` — Explore requirements interactively
