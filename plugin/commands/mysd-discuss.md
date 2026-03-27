---
model: claude-sonnet-4-5
description: Ad-hoc discussion with optional 4-dimension research, gray area exploration, and scope guardrail. Updates specs and triggers re-plan. Usage: /mysd:discuss [topic|change-name|file-path|dir-path] [--auto]
argument-hint: "[topic|change-name|file|dir] [--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
---

# /mysd:discuss -- Ad-hoc Discussion & Research

You are the mysd discuss orchestrator. Your job is to facilitate structured discussion, optionally with parallel 4-dimension research, gray area exploration, and scope guardrail, then propagate conclusions to specs.

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

## Step 4: Conditional Deferred Notes Loading (D-02)

Per D-02: discuss checks for an active WIP change before loading deferred notes (to avoid polluting a focused WIP discussion).

Run:
```bash
mysd status
```

- If output shows an active change in a non-archived state (proposed / spec-ready / planned / executing): do NOT load deferred notes. Set `deferred_context` to empty string.
- If NO active WIP change exists, or the active change is archived: run `mysd note list` and include the output as `deferred_context`.

## Step 4.5: Research Cache Detection (D-17)

If `cache_action` is not yet set and `change_name` is set (mode = "change"), check for cached research:

Read file `.specs/changes/{change_name}/discuss-research-cache.json` using the Read tool.

**If file exists and contains valid JSON:**
- Extract `cached_at` field
- If `auto_mode` is true: set `cache_action = "fresh"` (always run fresh research in auto mode)
- If `auto_mode` is false: Ask user:
  ```
  Found cached research from {cached_at}.
  1. Reuse cached research (skip research step)
  2. Run fresh research (overwrite cache)
  3. Skip research entirely (cache unchanged)
  ```
  - If user chooses 1 (reuse): set `cache_action = "reuse"`, load `research` object from cache, skip Steps 5-8, go directly to Step 9 with cached research as context
  - If user chooses 2 (fresh): set `cache_action = "fresh"`, proceed to Step 5 normally
  - If user chooses 3 (skip): set `cache_action = "skip"`, skip Steps 5-8, go to Step 9 without research

**If file does not exist or is invalid JSON (`discuss-research-cache.json` missing or unparseable):**
- Set `cache_action = "none"` (no cache, proceed normally to Step 5)

## Step 5: Optional Research (DISC-04, D-06)

If `cache_action` is "reuse" or "skip": skip this step entirely (already handled in Step 4.5).

If `auto_mode` is true: skip research entirely (FAUTO-02 — auto means no interaction). Go directly to Step 9.

If `auto_mode` is false: Ask user:
```
Would you like to run 4-dimension research on this topic?
(Codebase / Domain / Architecture / Pitfalls) [y/N]
```

- If user declines: go to Step 9 (discussion without research).
- If user accepts: proceed to Step 6.

## Step 6: Parallel Research Spawning

Spawn 4 `mysd-researcher` agents in parallel using the Task tool:

For each dimension in ["codebase", "domain", "architecture", "pitfalls"]:
```
Task: Research {dimension} for topic: {topic}
Agent: mysd-researcher
Context: {
  "change_name": "{change_name}",
  "dimension": "{dimension}",
  "topic": "{topic}",
  "spec_files": [{spec file paths from Step 3}],
  "auto_mode": false
}
```

Collect all 4 research outputs. Present organized summary by dimension to the user.

## Step 6.5: Write Research Cache (D-16)

After collecting all 4 research outputs from Step 6, write the cache file:

Use the Write tool to create `.specs/changes/{change_name}/discuss-research-cache.json` with content:
```json
{
  "change_name": "{change_name}",
  "cached_at": "{current ISO8601 timestamp}",
  "research": {
    "architecture": "{architecture research output}",
    "codebase": "{codebase research output}",
    "ux": "{ux/domain research output}",
    "security": "{security/pitfalls research output}"
  }
}
```

IMPORTANT: The research dimension values must be properly escaped JSON strings. Use the Write tool which handles escaping automatically.

Set `cache_action = "written"`.

If write fails: continue silently (cache is best-effort, do not interrupt the discussion flow).

## Step 7: Gray Area Identification + Advisor Spawning (DISC-06)

From the 4 research outputs, identify gray areas: ambiguous design decisions where multiple valid approaches exist, conflicting recommendations between dimensions, or areas needing user input.

For each gray area, spawn one `mysd-advisor` agent in parallel using the Task tool:
```
Task: Analyze gray area: {gray_area_description}
Agent: mysd-advisor
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
2. Finish exploration and proceed to discussion
```

If user chooses "Explore additional areas":
- User describes new areas to investigate
- Spawn one `mysd-advisor` agent per new area (same pattern as Step 7)
- Run Layer 1 deep dive for each new area

If user chooses "Finish exploration": proceed to Step 9.

## Step 9: Discussion Loop

Facilitate discussion with the user:

If research was performed (Steps 6-8 executed):
- Present key findings from each dimension
- Highlight conclusions from gray area exploration
- Discuss remaining open questions or implementation decisions

If no research:
- Discuss the topic based on existing spec context
- Help clarify requirements, edge cases, trade-offs

Continue discussion until user reaches a conclusion.

After each conclusion point, ask:
```
Would you like to:
1. Incorporate this conclusion into the spec
2. Continue discussing further
3. Done — end discussion without spec changes
```

If auto_mode: automatically choose "Incorporate" for all conclusions.

## Step 10: Spec Update

When user chooses to incorporate conclusions:

Determine which spec layer(s) are affected:

**If proposal layer** (scope change, motivation update):
```
Task: Update proposal with discussion conclusions
Agent: mysd-proposal-writer
Context: {
  "change_name": "{change_name}",
  "conclusions": "{conclusions text}",
  "existing_proposal": "{current proposal body}",
  "auto_mode": {auto_mode}
}
```

**If specs/ layer** (requirement changes):
For each affected capability area:
```
Task: Update spec for {capability_area}
Agent: mysd-spec-writer
Context: {
  "change_name": "{change_name}",
  "capability_area": "{area}",
  "existing_spec_body": "{current spec content}",
  "proposal": "{proposal body}",
  "auto_mode": {auto_mode}
}
```

**If design layer** (architecture changes):
```
Task: Update design with discussion conclusions
Agent: mysd-designer
Context: {
  "change_name": "{change_name}",
  "conclusions": "{conclusions text}",
  "auto_mode": {auto_mode}
}
```

## Step 11: Re-plan + Plan-Checker

After spec updates complete:

1. Get new planning context:
   Run: `mysd plan --context-only`

2. Spawn planner:
   ```
   Task: Re-plan after discussion updates
   Agent: mysd-planner
   Context: {planning context JSON with auto_mode}
   ```

3. Run state transition:
   Run: `mysd plan`

4. Get check context:
   Run: `mysd plan --check --context-only`

5. Spawn plan-checker:
   ```
   Task: Validate plan coverage after discussion updates
   Agent: mysd-plan-checker
   Context: {check output JSON}
   ```

## Step 12: Confirm

Show summary:
- Topic discussed
- Research performed (4-dimension research: yes/no)
- Research cache: {written/reused/skipped/not applicable}
- Number of gray areas explored (if research was run)
- Number of ideas deferred via scope guardrail (if any)
- Spec files updated
- Plan-checker results
- Next: `/mysd:apply` to execute updated plan
