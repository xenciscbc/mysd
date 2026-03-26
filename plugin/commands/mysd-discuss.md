---
model: claude-sonnet-4-5
description: Ad-hoc discussion with optional 4-dimension research. Updates specs and triggers re-plan. Usage: /mysd:discuss [topic|change-name|file-path|dir-path] [--auto]
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:discuss -- Ad-hoc Discussion & Research

You are the mysd discuss orchestrator. Your job is to facilitate structured discussion, optionally with parallel research, and propagate conclusions to specs.

## Step 1: Parse Arguments

Check `$ARGUMENTS`:
- If `--auto` present: set `auto_mode = true`, remove from arguments
- Otherwise: `auto_mode = false`

## Step 2: Source Detection (D-06)

Apply source detection in priority order:

1. If remaining arguments match a directory `.specs/changes/{name}/` -> set `change_name = {name}`, mode = "change"
2. If arguments is a file path (exists as file) -> mode = "file", read file content as context
3. If arguments is a directory path -> mode = "directory", list `.md` files for selection
4. If no argument + run `mysd status` shows active change -> use that change_name, mode = "change"
5. If no argument + no active change:
   - Check `~/.gstack/projects/` for project directory with `.md` files
   - Check conversation context for mentioned documents
   - Do NOT check `.claude/plans/` (D-07)
   - If auto_mode: use first detected; else: present options
6. If nothing found:
   - Ask: "No existing change found. Create a new one? (provide change name)"
   - Run `mysd propose {name}` to scaffold
   - Set mode = "change", change_name = {name}

## Step 3: Topic Identification (D-01)

If mode is "change":
  - Read `.specs/changes/{change_name}/proposal.md` for context
  - Read `.specs/changes/{change_name}/specs/` for existing requirements

Extract topic:
- If arguments contained a topic string (not a path/change-name): use it directly
- If auto_mode: derive topic from the change context
- Otherwise: Ask "What topic would you like to discuss?"

## Step 4: Optional Research (D-02, D-03)

If `auto_mode` is false:
  Ask: "Would you like to run 4-dimension research on this topic? (Codebase / Domain / Architecture / Pitfalls) [y/N]"

If `auto_mode` is true:
  Skip research (per FAUTO-02 — auto means no interactive, and discuss auto skips research like ff)

If user chooses research:
  Spawn 4 `mysd-researcher` agents in parallel:

  For each dimension in ["codebase", "domain", "architecture", "pitfalls"]:
    Task: Research {dimension} for topic: {topic}
    Agent: mysd-researcher
    Context: {
      "change_name": "{change_name}",
      "dimension": "{dimension}",
      "topic": "{topic}",
      "spec_files": [{spec file paths}],
      "auto_mode": false
    }

  Collect all 4 research outputs.
  Present findings summary to user, organized by dimension.

## Step 5: Discussion Loop (D-04)

Facilitate discussion with the user:

If research was performed:
  - Present key findings from each dimension
  - Highlight conflicts or gray areas between dimensions
  - Ask user for their perspective on each finding

If no research:
  - Discuss the topic based on existing spec context
  - Help clarify requirements, edge cases, trade-offs

Continue discussion until user reaches a conclusion.

After each conclusion point, ask (D-04):
"Would you like to:
  1. **Incorporate** this conclusion into the spec
  2. **Continue** discussing further
  3. **Done** — end discussion without spec changes"

If auto_mode: automatically choose "Incorporate" for all conclusions.

## Step 6: Spec Update (D-05)

When user chooses to incorporate conclusions:

Determine which spec layer(s) are affected:

**If proposal layer** (scope change, motivation update):
  Task: Update proposal with discussion conclusions
  Agent: mysd-proposal-writer
  Context: {
    "change_name": "{change_name}",
    "conclusions": "{conclusions text}",
    "existing_proposal": "{current proposal body}",
    "auto_mode": {auto_mode}
  }

**If specs/ layer** (requirement changes):
  For each affected capability area:
    Task: Update spec for {capability_area}
    Agent: mysd-spec-writer
    Context: {
      "change_name": "{change_name}",
      "capability_area": "{area}",
      "existing_spec_body": "{current spec content}",
      "proposal": "{proposal body}",
      "auto_mode": {auto_mode}
    }

**If design layer** (architecture changes):
  Task: Update design with discussion conclusions
  Agent: mysd-designer
  Context: {
    "change_name": "{change_name}",
    "conclusions": "{conclusions text}",
    "auto_mode": {auto_mode}
  }

## Step 7: Re-plan + Plan-Checker (D-05, D-09)

After spec updates complete:

1. Get new planning context:
   Run: `mysd plan --context-only`

2. Spawn planner:
   Task: Re-plan after discussion updates
   Agent: mysd-planner
   Context: {planning context JSON with auto_mode}

3. Run state transition:
   Run: `mysd plan`

4. Get check context:
   Run: `mysd plan --check --context-only`

5. Spawn plan-checker:
   Task: Validate plan coverage after discussion updates
   Agent: mysd-plan-checker
   Context: {check output JSON}

## Step 8: Confirm

Show summary:
- Topic discussed
- Research performed (if any)
- Spec files updated
- Plan-checker results
- Next: `/mysd:apply` to execute updated plan
