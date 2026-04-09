---
description: Ad-hoc discussion with optional 4-dimension research, gray area exploration, and scope guardrail. Updates specs and triggers re-plan. Usage: /mysd:discuss [topic|change-name|file-path|dir-path] [--auto]
argument-hint: "[topic|change-name|file|dir] [--auto]"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Task
  - AskUserQuestion
---

# /mysd:discuss -- Ad-hoc Discussion & Research

You are the mysd discuss orchestrator. Your job is to facilitate structured discussion, optionally with parallel 4-dimension research, gray area exploration, and scope guardrail, then propagate conclusions to specs.

## Question Protocol

- Ask one question at a time. Wait for the user's answer before asking the next.
- When a question has concrete options, use the **AskUserQuestion tool** — do not list options as plain text.
- Open-ended questions may use plain text.

## Step 1: Parse Arguments

Check `$ARGUMENTS`:
- If `--auto` present: set `auto_mode = true`, remove from arguments
- Otherwise: `auto_mode = false`

## Step 2: Context Detection & Path Routing

Determine the active change and present path selection.

### Step 2a: Detect Active Change

1. If remaining arguments match a directory `.specs/changes/{name}/` → set `change_name = {name}`
2. Else if no argument → run `mysd status`. If it reports an active change → use that `change_name`
3. Else → `change_name = null` (no active change)

### Step 2b: Path Selection

**If `change_name` is set (active change exists):**

- **If `auto_mode` is true:** Enter source-driven path directly (no prompt).
- **If `auto_mode` is false:** Present path selection:
  ```
  Active change: {change_name}
  How would you like to proceed?
  1. Discuss an existing spec (spec-focused)
  2. Add content from external sources (source-driven)
  ```
  Wait for user selection.

**If `change_name` is null (no active change):**
→ Enter source-driven path directly (no prompt needed).

### Step 2c: Spec-Focused Path — Spec Listing

This step only executes when the user selected **spec-focused** path in Step 2b.

List all spec files under `.specs/changes/{change_name}/specs/*/spec.md`:
- For each spec, extract the `capability` field from frontmatter and count the number of `### Requirement:` headings
- Present as a numbered list:
  ```
  Specs in {change_name}:
  1. {capability_name} ({N} requirements)
  2. {capability_name} ({N} requirements)
  ...
  ```
- Wait for user to select one spec

**If the change has no specs:**
→ Inform the user: "No specs found in this change. Switching to source-driven path."
→ Fall back to source-driven path (continue to Step 2d).

After selection, set `selected_spec` = the chosen spec file path. Proceed to **Step 3** (gap analysis).

### Step 2d: Source-Driven Path — Material Selection

This step only executes when entering the **source-driven** path (user selection, auto mode, or no active change).

Proceed to Step 3b (Material Selection).

## Step 3: Spec Gap Analysis (spec-focused path only)

This step only executes when `selected_spec` is set (spec-focused path).
If this is the source-driven path, skip to Step 3b.

Read `selected_spec` and `.specs/changes/{change_name}/proposal.md`.

### Step 3a: Requirement Coverage Analysis

Compare the spec's `### Requirement:` headings against `proposal.md`'s Capabilities section:
- Identify capabilities listed in the proposal that have **no corresponding requirement** in the selected spec
- Identify requirements in the spec that do **not map to any capability** in the proposal (orphaned requirements)

### Step 3b-gap: Scenario Completeness Analysis

For each `### Requirement:` in the spec:
- Check if it has at least one `#### Scenario:` block
- Requirements with zero scenarios are flagged as **missing scenarios**

### Step 3c-gap: Boundary Condition Coverage Analysis

For each requirement's scenarios, check for:
- At least one **error/failure** scenario (keywords: "error", "fail", "invalid", "not found", "empty", "missing")
- At least one **edge case** scenario (keywords: "maximum", "minimum", "empty list", "single item", "concurrent", "timeout")

Requirements with only happy-path scenarios are flagged as **missing boundary coverage**.

### Step 3d: Present Gap Results

Present a structured summary grouped by dimension:

```
## Gap Analysis: {capability_name}

### Coverage Gaps
- {capability X} — no requirement covers this
- {requirement Y} — orphaned, not mapped to any proposal capability

### Scenario Gaps
- {requirement Z} — no scenarios defined

### Boundary Gaps
- {requirement W} — only happy-path scenarios, missing error/edge cases
```

**If gaps found:** Ask the user which gap to address first. Use the selected gap as the discussion starting point.

**If no gaps found:** Report the spec as complete and ask:
```
This spec has full coverage. Would you like to:
1. Discuss other aspects of this spec
2. End the discussion
```

Set `discussion_context` = gap analysis results. Proceed to Step 4.

## Step 3b: Material Selection (source-driven path only)

This step only executes for the source-driven path.

### Step 3b-i: Detect Sources

Apply the following detection logic for each source type:

| # | Source Type | Detection Method |
|---|-----------|-----------------|
| 1 | source_arg file/directory | Check if `source_arg` (from Step 1) is a valid file path or directory on disk. Skip if already matched as a change in Step 2a. |
| 2 | Conversation context | Check if the current conversation contains substantive requirement discussion (not just greetings or meta-talk). |
| 3 | Claude plan | Check if conversation system messages contain a plan file path matching `~/.claude/plans/<name>.md` and the file exists. |
| 4 | gstack plan | Scan `~/.gstack/projects/{project}/` for `.md` files (project derived from cwd name). |
| 5 | Active change | Run `mysd status`. If it reports an active change with `proposal.md`, record it. |
| 6 | Deferred notes | If `deferred_context` (from Step 4) is non-empty, record it. |

For each detected source, extract a brief content preview (first line, title, or summary).

### Step 3b-ii: Present Sources and Collect Selection

**If `auto_mode` is true:** Automatically aggregate content from all detected sources. If no sources detected, extract requirements from conversation context as best-effort. Go to Step 3b-iii.

**If `auto_mode` is false:**

If no sources detected:
→ Ask: "What would you like to discuss? Please describe the topic."
→ Use the user's description as `aggregated_content`. Go to Step 3b-iv.

If sources detected:
→ Present a numbered list with type labels and content previews.
→ Always include "Manual input" as the last option.
→ Allow multi-selection (e.g., "1,3").
→ Even with one source, still present for confirmation.

### Step 3b-iii: Aggregate Content

Read and combine content from all selected sources into `aggregated_content`.

### Step 3b-iv: Compare with Existing Specs

Scan specs under `.specs/changes/{change_name}/specs/` (if change exists) and `openspec/specs/`:
- Compare key concepts in `aggregated_content` against each spec's capability name and requirements
- Recommend:
  - **Clear match**: "This content maps to `{spec_name}`. Recommend merging."
  - **Multiple matches**: List related specs, let user choose which to extend
  - **No match**: "No existing spec matches. Recommend creating a new spec: `{suggested-name}`"
- Wait for user confirmation before proceeding to discussion

**If `change_name` is null (no active change):**
1. Derive a kebab-case change name from `aggregated_content`
2. Run `mysd propose {name}` to scaffold the change
3. Set `change_name` = derived name
4. Then execute the spec comparison above

Proceed to Step 4.

## Step 3c: Topic Identification

If `change_name` is set:
  - Read `.specs/changes/{change_name}/proposal.md` for context
  - Read `.specs/changes/{change_name}/specs/` for existing requirements

Extract topic:
- If arguments contained a topic string (not a path/change-name): use it directly
- If auto_mode: derive topic from the change context
- Otherwise: Ask "What topic would you like to discuss?"

## Step 3d: Resolve Agent Model

Run:
```
mysd model
```

Parse the output to find `Profile: {profile_name}`. The profile determines agent model:
- `quality` or `balanced` → model = `sonnet`
- `budget` → model = `haiku` (for researcher/advisor); model = `sonnet` (for proposal-writer/spec-writer/designer)

Use this `model` value when spawning agents in subsequent steps.

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

Show: "Spawning 4 mysd-researcher agents ({model})..."
Spawn 4 `mysd-researcher` agents in parallel using the Task tool, each with `model` parameter set to `{model}`:

For each dimension in ["codebase", "domain", "architecture", "pitfalls"]:
```
Task: Research {dimension} for topic: {topic}
Agent: mysd-researcher
Model: {model}
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

### Layer 2 — Unified Research Exit

After all identified gray areas from Step 7 are explored, present a unified summary of ALL gray area conclusions:

```
## Research Summary

### Conclusions
1. {gray_area_1}: {conclusion}
2. {gray_area_2}: {conclusion}
...

Would you like to:
1. Explore additional areas (describe what you'd like to investigate)
2. Continue discussing other aspects
3. Converge to conclusion and decide on spec updates
```

If user chooses "Explore additional areas":
- User describes new areas to investigate
- Spawn one `mysd-advisor` agent per new area (same pattern as Step 7)
- Run Layer 1 deep dive for each new area
- Then present updated unified summary again

If user chooses "Continue discussing": proceed to Step 9 with research conclusions as context.
If user chooses "Converge": proceed directly to Step 9's conclusion summary.

## Discussion Guidelines

Follow these rules throughout the discussion loop:

**一次一問。** 不要一次丟多個問題。問最重要的那個，聽完再追問。如果用戶的描述已經涵蓋了某個問題，直接跳過。

**提具體選項。** 探索方案時，給出 2–3 個有 trade-off 的具體選項，用比較表呈現：

| 方案 | 優點 | 缺點 |
|------|------|------|
| A    | ...  | ...  |
| B    | ...  | ...  |

**不說空話。** 禁止使用以下說法：
- ~~"這是個很有趣的想法"~~ → 說清楚有趣在哪、為什麼
- ~~"有很多方式可以思考"~~ → 直接列出 2–3 個方式和 trade-off
- ~~"這樣可能可行"~~ → 說清楚為什麼可行或不可行

**主動推薦。** 有意見就說。"我會選 B，因為..." 比 "各有優缺點" 有用。

**當用戶想快點結束時：**
1. 第一次：用一句話提醒重要的未解決問題。"在決定 X 之前，Y 可能影響 Z，要處理還是繼續？"
2. 若再催：尊重用戶的步調，直接跳到收斂，不再 push back。最多一次 nudge。

## Step 9: Discussion Loop

Facilitate discussion with the user. Follow the Discussion Guidelines above throughout.

Load context based on how this step was reached:

**If research was performed (Steps 6-8 executed):**
- Use the unified research summary from Step 8 Layer 2 as context
- Discuss remaining open questions or implementation decisions

**If no research + spec-focused path:**
- Use `discussion_context` (gap analysis results from Step 3d) as the discussion starting point
- Focus on addressing identified gaps

**If no research + source-driven path:**
- Use `aggregated_content` from Step 3b and the spec target recommendation as context
- Help clarify requirements, edge cases, trade-offs

Continue until a clear conclusion is reached. Then **proactively present the conclusion summary** — do not wait for the user to ask:

```
## Conclusion

**Decision**: [What was decided]
**Rationale**: [The key trade-off that drove this]
**Capture to**: [Which artifact: proposal.md / spec / design.md / tasks.md]
```

Say: "I'll capture this to {artifact} unless you'd rather not."

If the user tries to end without a conclusion, summarize what was discussed and state what remains unresolved. Do not let the discussion end without at least an explicit deferral (e.g., "We don't have enough information yet to decide X").

After presenting the conclusion summary, ask:
```
Would you like to:
1. Incorporate this conclusion into the spec
2. Continue discussing further
3. Done — end discussion without spec changes
```

If auto_mode: automatically choose "Incorporate" for all conclusions.

## Step 10: Spec Update with Confirmation

When user chooses to incorporate conclusions:

### Step 10a: Determine Affected Artifacts

Identify which artifacts are affected by the discussion conclusions:
- `proposal.md` — if scope change or motivation update
- `specs/{capability}/spec.md` — if requirement changes (one or more)
- `design.md` — if architecture changes

### Step 10b: User Confirmation

**If `auto_mode` is true:** Skip confirmation, update all affected artifacts.

**If `auto_mode` is false:** Present the affected artifacts as a checklist:

```
The following artifacts will be updated:
[x] 1. proposal.md — scope update
[x] 2. specs/material-selection/spec.md — new requirement
[x] 3. specs/discuss-path-routing/spec.md — scenario update

All items are selected. Deselect any you want to skip (e.g., "-2"):
```

Wait for user input:
- If user confirms without changes → update all
- If user deselects items (e.g., "-2") → remove those from the update list
- If user cancels → skip all updates, proceed to Step 11

### Step 10c: Execute Updates

Only spawn writer agents for the artifacts the user confirmed:

**If proposal layer confirmed:**
Show: "Spawning mysd-proposal-writer ({model})..."
```
Task: Update proposal with discussion conclusions
Agent: mysd-proposal-writer
Model: {model}
Context: {
  "change_name": "{change_name}",
  "conclusions": "{conclusions text}",
  "existing_proposal": "{current proposal body}",
  "auto_mode": {auto_mode}
}
```

**If specs/ layer confirmed:**
For each confirmed capability area, show: "Spawning mysd-spec-writer ({model})..."
```
Task: Update spec for {capability_area}
Agent: mysd-spec-writer
Model: {model}
Context: {
  "change_name": "{change_name}",
  "capability_area": "{area}",
  "existing_spec_body": "{current spec content}",
  "proposal": "{proposal body}",
  "auto_mode": {auto_mode}
}
```

**If design layer confirmed:**
Show: "Spawning mysd-designer ({model})..."
```
Task: Update design with discussion conclusions
Agent: mysd-designer
Model: {model}
Context: {
  "change_name": "{change_name}",
  "conclusions": "{conclusions text}",
  "auto_mode": {auto_mode}
}
```

## Step 11: Conditional Re-plan + Plan-Checker

After spec updates complete, check if a plan already exists:

Check whether `.specs/changes/{change_name}/tasks.md` exists (use the Read tool — if the file is not found, it does not exist).

**If `tasks.md` does NOT exist:** Skip Step 11 entirely. Proceed to Step 12.

**If `tasks.md` exists:** Execute the re-plan sequence:

1. Get new planning context:
   Run: `mysd plan --context-only`

2. Extract `model` from planning context JSON. Show: "Spawning mysd-planner ({model})..."
   Spawn planner with `model` parameter set to `{model}`:
   ```
   Task: Re-plan after discussion updates
   Agent: mysd-planner
   Model: {model}
   Context: {planning context JSON with auto_mode}
   ```

3. Run state transition:
   Run: `mysd plan`

4. Get check context:
   Run: `mysd plan --check --context-only`

5. Show: "Spawning mysd-plan-checker ({model})..."
   Spawn plan-checker with `model` parameter set to `{model}`:
   ```
   Task: Validate plan coverage after discussion updates
   Agent: mysd-plan-checker
   Model: {model}
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
