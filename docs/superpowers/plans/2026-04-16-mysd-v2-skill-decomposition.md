# mysd v2 — Pure SKILL.md Decomposition

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Decompose the mysd Go CLI (23 commands) into 4 independent SKILL.md files — research, doc, spec, and orchestrator — with zero external dependencies.

**Architecture:** Each skill is a standalone SKILL.md that Claude Code loads on demand. The orchestrator chains the three content skills via subagents, giving each its own context window. All skills use only Claude's built-in tools (Read/Write/Edit/Grep/Glob/Bash). No Go binary needed.

**Tech Stack:** Claude Code SKILL.md format, YAML frontmatter, Markdown

**Design doc:** `~/.gstack/projects/xenciscbc-mysd/cbc-master-design-20260416-113346.md`

**Eng review decisions:**
1. Analyzer 4 dimensions → research skill (not spec writer)
2. Only support `openspec/` directory format (drop `.specs/`)
3. No Phase 0 PoC (Phase 2 validation is sufficient)
4. Explicit trigger boundaries in each SKILL.md description
5. Architecture: 3 skills + 1 orchestrator (subagent pattern)

---

## File Structure

```
mysd-skills/
  research/
    SKILL.md                     # Core flow + triggers (~150 lines)
    formats/
      decision-doc.md            # Decision Doc output template
      health-check.md            # Analyzer 4-dimension logic
  doc/
    SKILL.md                     # Core flow + triggers + impact mapping (~150 lines)
  spec/
    SKILL.md                     # Core flow + OpenSpec format reference (~200 lines)
  orchestrator/
    SKILL.md                     # Chains research → doc → spec via subagents (~100 lines)
  plugin.json                    # Claude Code plugin manifest
  README.md                      # Installation and usage instructions
```

---

### Task 1: Create the skill repo structure

**Files:**
- Create: `mysd-skills/research/SKILL.md` (placeholder)
- Create: `mysd-skills/research/formats/decision-doc.md`
- Create: `mysd-skills/research/formats/health-check.md`
- Create: `mysd-skills/doc/SKILL.md` (placeholder)
- Create: `mysd-skills/spec/SKILL.md` (placeholder)
- Create: `mysd-skills/orchestrator/SKILL.md` (placeholder)

- [ ] **Step 1: Create the directory structure**

```bash
mkdir -p mysd-skills/research/formats
mkdir -p mysd-skills/doc
mkdir -p mysd-skills/spec
mkdir -p mysd-skills/orchestrator
```

- [ ] **Step 2: Initialize git repo**

```bash
cd mysd-skills
git init
```

- [ ] **Step 3: Create .gitignore**

Write `mysd-skills/.gitignore`:

```
.DS_Store
```

- [ ] **Step 4: Commit skeleton**

```bash
cd mysd-skills
git add .
git commit -m "chore: initialize mysd-skills repo structure"
```

---

### Task 2: Write the research skill — formats/decision-doc.md

**Files:**
- Create: `mysd-skills/research/formats/decision-doc.md`

This companion file contains the Decision Doc output template. The SKILL.md will instruct Claude to `Read` this file at the output step.

- [ ] **Step 1: Write the Decision Doc template**

Write `mysd-skills/research/formats/decision-doc.md`:

```markdown
# Decision Doc Output Format

Use this exact structure when producing a Decision Doc.

## Template

---

# Decision: {title}

## Problem
{Describe the problem with full context. Include what triggered this decision, who is affected, and any constraints.}

## Gray Area Classification
{Explain why this qualifies as a gray area. Must be one of:}
{- Multiple viable approaches with no community consensus}
{- Best practice exists but does not apply to this specific context (explain why)}
{- Must decide with incomplete information (list what is unknown)}

## Options

### Option A: {name}
- **Evidence:** {Concrete evidence — link, benchmark, code reference, or observed behavior. No pure speculation.}
- **Pros:** {List}
- **Cons:** {List}
- **Effort:** {S/M/L}

### Option B: {name}
- **Evidence:** {Same standard}
- **Pros:** {List}
- **Cons:** {List}
- **Effort:** {S/M/L}

{Add Option C only if meaningfully different from A and B.}

## Recommendation
**{Option name}** — Confidence: {N}/10

**Reasoning:** {Step-by-step logic connecting evidence to conclusion}

**What would change my mind:** {Specific evidence that would reverse this recommendation}

## Open Questions
{Questions that remain unanswered. Each should have a suggested way to resolve it.}

---

## Confidence Scale

| Score | Meaning |
|-------|---------|
| 1-3 | Guess. No direct evidence. High uncertainty. |
| 4-6 | Partial evidence. Significant risks or unknowns remain. |
| 7-8 | Multiple evidence sources. Risks identified and manageable. |
| 9-10 | Strong evidence. Nearly certain. 10 = verified with actual results. |
```

- [ ] **Step 2: Commit**

```bash
cd mysd-skills
git add research/formats/decision-doc.md
git commit -m "feat(research): add Decision Doc output format template"
```

---

### Task 3: Write the research skill — formats/health-check.md

**Files:**
- Create: `mysd-skills/research/formats/health-check.md`

This companion file contains the Spec Health Check logic, translated from the Go analyzer package (`internal/analyzer/`). Claude reads this file when asked to analyze spec quality.

- [ ] **Step 1: Write the health check format**

Write `mysd-skills/research/formats/health-check.md`:

```markdown
# Spec Health Check — 4 Dimension Analysis

When asked to analyze spec quality, run these 4 checks against the `openspec/` directory.

## How to Run

1. Use `Glob` to find all spec files: `openspec/specs/*/spec.md` and `openspec/changes/*/specs/*/spec.md`
2. Use `Read` to read proposal.md, design.md, tasks.md in the change directory
3. Run each dimension below
4. Output findings in the format at the bottom

---

## Dimension 1: Coverage

**Question:** Does every capability listed in the proposal have a corresponding spec file?

**Steps:**
1. Read `proposal.md` in the change directory
2. Find the `## Capabilities`, `### New Capabilities`, or `### Modified Capabilities` section
3. Extract capability names from bullet points matching pattern: `- \`capability-name\`:`
4. For each capability name, check if `specs/{capability-name}/spec.md` exists
5. Missing spec file = CRITICAL finding

**Finding format:**
```
[COV-N] CRITICAL — proposal.md
Capability '{name}' listed in proposal has no corresponding specs/{name}/spec.md
→ Create specs/{name}/spec.md or remove '{name}' from proposal
```

## Dimension 2: Ambiguity

**Question:** Do spec files contain vague language that should be precise RFC 2119 keywords?

**Steps:**
1. Read each `spec.md` file
2. Skip lines inside code blocks (``` fences)
3. Search for these weak patterns (case-insensitive for words, exact for abbreviations):
   - `should` (lowercase — not RFC 2119 SHOULD)
   - `may` (lowercase — not RFC 2119 MAY)
   - `might`
   - `TBD`, `TODO`, `FIXME`, `TKTK`, `???`
4. One finding per line maximum

**Finding format:**
```
[AMB-N] SUGGESTION — specs/{capability}/spec.md:{line}
Vague language '{pattern}' found
→ Replace with MUST/MUST NOT/SHALL/SHALL NOT for clarity
```

## Dimension 3: Consistency

**Question:** Do design decisions align with tasks?

**Steps:**
1. Read `design.md` and extract all `###` headings (these are design decisions)
2. Read `tasks.md`
3. For each design heading, check if the heading text appears (case-insensitive) in tasks.md
4. Missing reference = WARNING finding

**Finding format:**
```
[CON-N] WARNING — design.md
Design topic '{heading}' not referenced in tasks
→ Verify tasks cover this design decision
```

## Dimension 4: Gaps

**Question:** Do requirements have scenarios? Do tasks reference requirements?

**Steps:**
1. Read each spec.md and find `### Requirement: {name}` headings
2. Check that each requirement has at least one `#### Scenario:` heading below it (before the next `### Requirement:`)
3. Read tasks.md and check that each requirement name appears somewhere in the tasks
4. Missing scenario = WARNING, missing task reference = WARNING

**Finding format:**
```
[GAP-N] WARNING — specs/{capability}/spec.md:{line}
Requirement '{name}' has no scenario
→ Add at least one #### Scenario: under this requirement

[GAP-N] WARNING — specs/{capability}/spec.md
Requirement '{name}' has no matching task
→ Add a task in tasks.md that references '{name}'
```

---

## Output Summary Format

```
Spec Health Check: {change-name}
──────────────────────────────
Coverage:      {Clean | N issue(s) found | Skipped (insufficient artifacts)}
Ambiguity:     {Clean | N issue(s) found | Skipped}
Consistency:   {Clean | N issue(s) found | Skipped}
Gaps:          {Clean | N issue(s) found | Skipped}

Findings:
  [COV-1] CRITICAL [proposal.md] Capability 'xyz' has no spec
  [AMB-1] SUGGESTION [specs/abc/spec.md:12] Vague language 'should'
  ...
```

**Skip rules:**
- Coverage: requires proposal.md + specs/ directory
- Ambiguity: requires specs/ directory
- Consistency: requires proposal.md + design.md + tasks.md
- Gaps: requires specs/ + tasks.md
```

- [ ] **Step 2: Commit**

```bash
cd mysd-skills
git add research/formats/health-check.md
git commit -m "feat(research): add Spec Health Check 4-dimension analysis format"
```

---

### Task 4: Write the research skill — SKILL.md

**Files:**
- Create: `mysd-skills/research/SKILL.md`

Core flow for the research + gray area decision skill. References companion files in `formats/`.

- [ ] **Step 1: Write the SKILL.md**

Write `mysd-skills/research/SKILL.md`:

```markdown
---
name: mysd:research
description: >
  Research ambiguous problems and make gray-area decisions with evidence.
  Use when facing technical choices with 2+ viable options and no clear consensus,
  or when analyzing spec quality. DO NOT use for documentation updates (use mysd:doc)
  or spec writing (use mysd:spec). DO NOT use for questions with clear best practices
  or official documentation answers.
---

# Research + Gray Area Decision Skill

## When to Use

Use this skill when:
- A question has 2+ reasonable answers and no community consensus
- A best practice exists but may not apply to the current context
- A decision must be made with incomplete information
- Someone asks to analyze spec quality or check spec health

Do NOT use this skill when:
- The answer is in official documentation → just answer directly
- It's a syntax/API usage question → just answer directly
- The user wants to update docs → use `/mysd:doc`
- The user wants to write or update specs → use `/mysd:spec`

## Flow

### Step 1: Classify the Problem

Read the user's question. Determine if it qualifies as a gray area:

**Gray area** = at least one of:
- (a) 2+ viable approaches, no community consensus
- (b) Best practice exists but does not apply here (must explain why)
- (c) Must decide with incomplete information

**Not gray area:**
- Clear best practice that applies → answer directly, skip this skill
- Official docs answer the question → answer directly, skip this skill
- Pure syntax/API question → answer directly, skip this skill

If not gray area, tell the user: "This isn't a gray area — here's the answer: ..."
and do not produce a Decision Doc.

### Step 2: Context Gathering

Collect information in this order:
1. **Codebase:** Read relevant spec files, source code, config files using Grep/Glob/Read
2. **Git history:** Run `git log --oneline -20` and `git diff --stat` to understand recent changes
3. **Project docs:** Read CLAUDE.md, README.md for project context
4. **External:** Use WebSearch for current best practices (if available; skip if not)

### Step 3: Option Framing

Structure the problem into distinct options. Each option MUST have:
- Concrete evidence (not speculation)
- Pros and cons
- Effort estimate (S/M/L)

Minimum 2 options, maximum 4.

### Step 4: Recommendation

Pick one option. State confidence level (1-10) with reasoning.
State what evidence would change your mind.

### Step 5: Output

Read `formats/decision-doc.md` in this skill's directory for the exact output format.
Write the Decision Doc and present it to the user.

## Spec Health Check Mode

When the user asks to analyze spec quality, check spec health, or asks
"what's wrong with these specs":

1. Identify the target: a specific change directory or the entire `openspec/specs/` directory
2. Read `formats/health-check.md` in this skill's directory for the analysis procedure
3. Run all 4 dimensions (Coverage, Ambiguity, Consistency, Gaps)
4. Present findings in the summary format defined in the health check file
```

- [ ] **Step 2: Verify the SKILL.md is under 80 lines (core flow only)**

```bash
wc -l mysd-skills/research/SKILL.md
```

Expected: under 80 lines. Detailed formats are in companion files.

- [ ] **Step 3: Commit**

```bash
cd mysd-skills
git add research/SKILL.md
git commit -m "feat(research): add research + gray area decision SKILL.md"
```

---

### Task 5: Write the doc skill — SKILL.md

**Files:**
- Create: `mysd-skills/doc/SKILL.md`

- [ ] **Step 1: Write the SKILL.md**

Write `mysd-skills/doc/SKILL.md`:

```markdown
---
name: mysd:doc
description: >
  Update documentation files based on code changes. Detects which docs need updating
  from git diff, generates content matching existing style, and applies changes with
  user confirmation. DO NOT use for spec files (use mysd:spec) or research/decisions
  (use mysd:research).
---

# Doc Writer Skill

## When to Use

Use this skill when:
- The user says "update docs", "sync README", "add to CHANGELOG"
- A feature or fix was completed and documentation needs to reflect it
- Multiple language versions of docs need to stay in sync (e.g., README.md + README.zh-TW.md)

Do NOT use this skill when:
- The user wants to update OpenSpec spec files → use `/mysd:spec`
- The user wants to research a decision → use `/mysd:research`
- The user wants to write new docs from scratch without a code change context → just write them directly

## Flow

### Step 1: Detect Changes

Determine the change scope. Accept these input formats:
- **Default:** `git diff --name-only HEAD~1`
- **User-specified range:** `HEAD~N`, `<sha1>..<sha2>`
- **Explicit file list:** user provides specific files

Run the appropriate git command and collect the list of changed files.

If no changes are detected (empty diff), tell the user:
"No changes detected. Specify a diff range or file list."
Do not proceed.

### Step 2: Impact Analysis

Use this mapping to determine which docs need updating:

| Change Type | Detection Pattern | Docs to Update |
|------------|-------------------|----------------|
| New/removed command | `cmd/*.go` added/removed | README.md, README.zh-TW.md, CLAUDE.md (if has command list) |
| API change | exported function signature changed | API docs, CHANGELOG.md |
| Config change | config struct or yaml schema changed | README (configuration section), example configs |
| Bug fix | any .go file modified with "fix" in commit message | CHANGELOG.md |
| Architecture change | new package directory or major refactor | ARCHITECTURE.md (if exists), README relevant sections |
| Dependency update | go.mod, package.json changed | README (installation section, if version requirements) |

**Fallback heuristic:** If the change doesn't match any row above, extract 3-5
significant keywords from the changed file names and content, then:

```bash
grep -rl "keyword1\|keyword2\|keyword3" *.md **/*.md 2>/dev/null
```

Any .md files that reference these keywords may need updating.

### Step 3: Style Matching

For each doc file that needs updating:
1. Read the first 50 lines to understand the style:
   - Heading levels (# vs ## vs ###)
   - List style (- vs * vs numbered)
   - Language/tone (formal vs casual, Chinese vs English)
   - Code block conventions
2. Generate content that matches the existing style

### Step 4: Multi-Language Sync

When updating a doc file, check for locale variants:
- `README.md` → also check `README.zh-TW.md`, `README.ja.md`, etc.
- Pattern: `{filename}.{locale}.{ext}` or `{filename}-{locale}.{ext}`

If a locale variant exists, generate the equivalent update in that language.

### Step 5: Apply and Confirm

For each doc update:
1. Show the user the proposed change using Edit tool's old_string/new_string format
2. Wait for user confirmation before applying
3. Apply confirmed changes

If multiple docs need updating, present them one at a time.

### Step 6: Summary

After all changes are applied, output:
```
Doc sync complete:
- README.md: updated (added command X)
- README.zh-TW.md: synced
- CHANGELOG.md: added entry for fix Y
```
```

- [ ] **Step 2: Verify line count**

```bash
wc -l mysd-skills/doc/SKILL.md
```

Expected: under 120 lines.

- [ ] **Step 3: Commit**

```bash
cd mysd-skills
git add doc/SKILL.md
git commit -m "feat(doc): add doc writer SKILL.md"
```

---

### Task 6: Write the spec skill — SKILL.md

**Files:**
- Create: `mysd-skills/spec/SKILL.md`

This is the largest skill because it embeds the OpenSpec format reference.

- [ ] **Step 1: Write the SKILL.md**

Write `mysd-skills/spec/SKILL.md`:

```markdown
---
name: mysd:spec
description: >
  Write and update OpenSpec format spec files based on code changes or plans.
  Generates correct YAML frontmatter, RFC 2119 requirements, and scenario definitions.
  DO NOT use for documentation updates (use mysd:doc) or research/decisions
  (use mysd:research).
---

# Spec Writer Skill

## When to Use

Use this skill when:
- The user says "update spec", "write spec", "sync spec"
- Implementation is done and specs need to reflect actual behavior
- A proposal was approved and formal specs need to be written

Do NOT use this skill when:
- The user wants to update README/CHANGELOG/docs → use `/mysd:doc`
- The user wants to research a decision → use `/mysd:research`
- The user wants to analyze spec quality → use `/mysd:research` (Spec Health Check mode)

## OpenSpec Format Reference

### Spec Frontmatter (required fields)

```yaml
---
spec-version: "1.0"
capability: "{capability-name}"
delta: ADDED | MODIFIED | REMOVED | RENAMED
status: pending | in_progress | done | blocked
---
```

Optional fields: `name`, `description`, `version`, `generatedBy`

### Proposal Frontmatter

```yaml
---
spec-version: "1.0"
change: "{change-name}"
status: "{status}"
created: "{ISO date}"
updated: "{ISO date}"
---
```

### Tasks Frontmatter (V2)

```yaml
---
spec-version: "1.0"
total: {N}
completed: {N}
tasks:
  - id: 1
    name: "{task name}"
    status: pending | in_progress | done | blocked
    spec: "{spec directory name}"
    depends: [2, 3]
    files: ["path/to/file.go"]
    satisfies: ["REQ-001"]
---
```

### RFC 2119 Keywords

- **MUST** — absolute requirement
- **MUST NOT** — absolute prohibition
- **SHOULD** — strongly recommended unless there is a compelling reason not to
- **SHOULD NOT** — strongly discouraged unless there is a compelling reason to
- **MAY** — optional

Use UPPERCASE only. Lowercase "should", "may" etc. are not RFC 2119 keywords.

### Scenario Format

```markdown
### Scenario: {descriptive name}

WHEN {trigger condition}
THEN {expected behavior}
AND {additional expectations}
```

### Directory Structure

```
openspec/
  specs/{capability}/spec.md          # Main specs (merged from changes)
  changes/{name}/
    proposal.md
    .openspec.yaml
    specs/{capability}/spec.md        # Delta specs for this change
    design.md
    tasks.md
  changes/archive/{date}-{name}/      # Archived completed changes
```

## Flow

### Step 1: Understand Change Context

Read the relevant sources:
- If a proposal exists: read `openspec/changes/{name}/proposal.md`
- If a plan exists: read `openspec/changes/{name}/tasks.md`
- If coming from code changes: run `git diff --name-only HEAD~1` and read changed files

### Step 2: Spec Discovery

Find related existing specs:
```bash
# Find all existing spec files
find openspec/specs -name "spec.md" -type f 2>/dev/null
```

Read each spec's frontmatter to find specs with matching `capability` names.

### Step 3: Gap Analysis

Compare the change context against existing specs:
- **New capability** → need new spec with `delta: ADDED`
- **Modified capability** → need updated spec with `delta: MODIFIED`
- **Removed capability** → need spec with `delta: REMOVED`
- **Renamed capability** → need spec with `delta: RENAMED`

### Step 4: Spec Generation

For each spec that needs to be written or updated:

1. **Determine frontmatter:**
   - `spec-version`: always `"1.0"`
   - `capability`: inferred from the change (see reverse-spec rules below)
   - `delta`: ADDED for new, MODIFIED for changed, REMOVED for deleted, RENAMED for renamed
   - `status`: `done` if implementation is complete, `pending` otherwise

2. **Write requirements** using RFC 2119 keywords:
   - Each requirement starts with `## Requirement: {name}`
   - Use MUST for absolute requirements
   - Use SHOULD for strong recommendations
   - Use MAY for optional features

3. **Write scenarios** under each requirement:
   - `### Scenario: {name}`
   - WHEN/THEN/AND format

### Step 5: Reverse-Spec Rules (from code)

When generating specs from code changes (no proposal/plan available):

1. Read the changed `.go` files
2. Identify exported functions, types, and constants
3. Map file paths to capabilities:
   - `cmd/` → command behavior capability
   - `internal/{package}/` → internal capability named after the package
4. From function signatures and doc comments, infer:
   - What the function MUST do (→ MUST requirement)
   - What inputs it accepts (→ WHEN clause in scenario)
   - What outputs it produces (→ THEN clause in scenario)
5. Generate frontmatter with `capability` = inferred capability name, `delta` = ADDED or MODIFIED

### Step 6: Validation

After generating the spec, verify:
- [ ] Frontmatter has all 4 required fields: `spec-version`, `capability`, `delta`, `status`
- [ ] All requirements use RFC 2119 keywords (UPPERCASE)
- [ ] Each requirement has at least one scenario
- [ ] File is at the correct path: `openspec/specs/{capability}/spec.md` or `openspec/changes/{name}/specs/{capability}/spec.md`

If validation fails, fix the issues before presenting to the user.

### Step 7: Output

Present the generated spec to the user. Use the Write tool to create the file at the correct path.
If modifying an existing spec, use the Edit tool and show the diff.
```

- [ ] **Step 2: Verify line count**

```bash
wc -l mysd-skills/spec/SKILL.md
```

Expected: under 200 lines.

- [ ] **Step 3: Commit**

```bash
cd mysd-skills
git add spec/SKILL.md
git commit -m "feat(spec): add spec writer SKILL.md with OpenSpec format reference"
```

---

### Task 7: Write the orchestrator skill — SKILL.md

**Files:**
- Create: `mysd-skills/orchestrator/SKILL.md`

- [ ] **Step 1: Write the SKILL.md**

Write `mysd-skills/orchestrator/SKILL.md`:

```markdown
---
name: mysd:run
description: >
  Orchestrate a full content intelligence workflow: research → doc → spec.
  Chains the three mysd skills using subagents, passing context between them.
  Each skill runs in its own context window. Use when you want the full pipeline,
  not just one skill. For individual skills, use mysd:research, mysd:doc, or mysd:spec directly.
---

# Orchestrator Skill

## When to Use

Use this skill when:
- The user wants a complete workflow: research a decision, update docs, and update specs
- The user says "run the full pipeline" or "do everything"
- A significant change was made and all content artifacts need updating

For individual tasks, use the specific skill directly:
- Decision needed → `/mysd:research`
- Docs out of date → `/mysd:doc`
- Specs need updating → `/mysd:spec`

## Flow

### Step 1: Understand the Scope

Ask the user what changed and what they need. Determine which skills to run:
- **Full pipeline:** research → doc → spec (all three)
- **Partial:** any subset (e.g., doc + spec only)

### Step 2: Run Research (if needed)

Dispatch a subagent using the Agent tool:

```
Agent({
  description: "Research and gray-area decision",
  prompt: "You are the mysd:research skill. [Read the SKILL.md at mysd-skills/research/SKILL.md and follow its instructions.] The user's question is: {question}. The project is at {repo_root}.",
  subagent_type: "general-purpose"
})
```

Collect the subagent's Decision Doc output. Pass relevant decisions to the next steps.

### Step 3: Run Doc Writer (if needed)

Dispatch a subagent:

```
Agent({
  description: "Update documentation",
  prompt: "You are the mysd:doc skill. [Read the SKILL.md at mysd-skills/doc/SKILL.md and follow its instructions.] Context from research: {research_output_summary}. The project is at {repo_root}. Diff range: {diff_range}.",
  subagent_type: "general-purpose"
})
```

### Step 4: Run Spec Writer (if needed)

Dispatch a subagent:

```
Agent({
  description: "Update specs",
  prompt: "You are the mysd:spec skill. [Read the SKILL.md at mysd-skills/spec/SKILL.md and follow its instructions.] Context from research: {research_output_summary}. The project is at {repo_root}. Change: {change_name}.",
  subagent_type: "general-purpose"
})
```

### Step 5: Summary

After all subagents complete, present a unified summary:

```
Pipeline complete:
- Research: {Decision Doc title, confidence N/10}
- Docs: {N files updated}
- Specs: {N specs written/updated}
```

## Error Handling

If a subagent times out or fails:
1. Report which step failed
2. Suggest running that skill directly: "The doc writer step failed. Try running `/mysd:doc` directly."
3. Continue with remaining steps (they are independent)
```

- [ ] **Step 2: Verify line count**

```bash
wc -l mysd-skills/orchestrator/SKILL.md
```

Expected: under 100 lines.

- [ ] **Step 3: Commit**

```bash
cd mysd-skills
git add orchestrator/SKILL.md
git commit -m "feat(orchestrator): add orchestrator SKILL.md with subagent chaining"
```

---

### Task 8: Write plugin.json and README

**Files:**
- Create: `mysd-skills/plugin.json`
- Create: `mysd-skills/README.md`

- [ ] **Step 1: Write plugin.json**

Write `mysd-skills/plugin.json`:

```json
{
  "name": "mysd-skills",
  "version": "1.0.0",
  "description": "Content intelligence skills for spec-driven development — research decisions, sync docs, write specs"
}
```

- [ ] **Step 2: Write README.md**

Write `mysd-skills/README.md`:

```markdown
# mysd-skills

Content intelligence skills for spec-driven development. Three independent skills + one orchestrator.

## Skills

| Skill | Command | Purpose |
|-------|---------|---------|
| Research | `/mysd:research` | Gray-area decisions with evidence. Spec health checks. |
| Doc Writer | `/mysd:doc` | Update docs based on code changes. Multi-language sync. |
| Spec Writer | `/mysd:spec` | Write/update OpenSpec format spec files. |
| Orchestrator | `/mysd:run` | Chain all three skills via subagents. |

## Install

```bash
# Option 1: Claude Code plugin
claude plugin add /path/to/mysd-skills

# Option 2: Manual copy
cp -r mysd-skills/research ~/.claude/skills/mysd-research
cp -r mysd-skills/doc ~/.claude/skills/mysd-doc
cp -r mysd-skills/spec ~/.claude/skills/mysd-spec
cp -r mysd-skills/orchestrator ~/.claude/skills/mysd-run
```

## Usage

Each skill works independently:

```
/mysd:research    — "Which database should we use for this feature?"
/mysd:doc         — "Update the README to reflect the new commands"
/mysd:spec        — "Write a spec for the new auth middleware"
/mysd:run         — "Run the full pipeline for this change"
```

## Requirements

- Claude Code (any version with skill support)
- No external dependencies — pure SKILL.md files

## OpenSpec Compatibility

The spec writer produces files compatible with the [OpenSpec](https://github.com/openspec) format:
- YAML frontmatter with `spec-version`, `capability`, `delta`, `status`
- RFC 2119 keywords (MUST/SHOULD/MAY)
- WHEN/THEN/AND scenario format
- Directory structure: `openspec/specs/{capability}/spec.md`
```

- [ ] **Step 3: Commit**

```bash
cd mysd-skills
git add plugin.json README.md
git commit -m "feat: add plugin.json and README"
```

---

### Task 9: Update design doc with eng review decisions

**Files:**
- Modify: `~/.gstack/projects/xenciscbc-mysd/cbc-master-design-20260416-113346.md`

- [ ] **Step 1: Update analyzer mapping**

In the design doc, find the section about "Skill 1: Research" and add a note that the analyzer 4 dimensions are part of this skill's Spec Health Check mode (not spec writer).

- [ ] **Step 2: Remove Phase 0**

Delete the Phase 0 section. Renumber Phase 1 → Phase 1, Phase 2 → Phase 2, etc.

- [ ] **Step 3: Add orchestrator to architecture**

In the "The Three Skills" section, add a fourth skill: orchestrator. Update the section title to "The Four Skills" or "Skills Architecture."

- [ ] **Step 4: Update Phase 2 validation table**

Add the 6 edge case test scenarios from the eng review:

| Skill | Test Scenario | Pass Condition |
|-------|--------------|----------------|
| research | (4) Spec Health Check | Runs 4 dimensions, produces summary |
| research | (5) Non-gray-area question | Correctly rejects, answers directly |
| doc | (4) Unmapped change type | Falls back to grep heuristic |
| doc | (5) Empty diff | Reports "no docs need updating" |
| spec | (4) RENAMED delta | Uses delta=RENAMED correctly |
| spec | (5) Incompatible spec-version | Handles gracefully |

- [ ] **Step 5: Add trigger boundary note**

In the Resolved Design Decisions section, add:
```
4. **Trigger boundaries:** Each SKILL.md's description field includes explicit
   排他語言 ("DO NOT use for X, use mysd:Y instead") to prevent Claude Code
   from triggering the wrong skill.
```

- [ ] **Step 6: Commit**

This file is outside the git repo (in ~/.gstack/), so no git commit needed. Just save.

---

### Task 10: Resolve TODOS

**Files:**
- Modify: `TODOS.md`

- [ ] **Step 1: Mark TODO #1 as done**

After Task 9 completes, update TODOS.md to mark the design doc update as done.

- [ ] **Step 2: Mark TODO #2 as done**

After Task 9 Step 4 completes (Phase 2 validation table updated), mark this as done.

- [ ] **Step 3: Commit**

```bash
git add TODOS.md
git commit -m "chore: resolve eng review TODOs"
```

---

## Self-Review Checklist

1. **Spec coverage:** Design doc decisions (5 items) → all covered by Tasks 4-9 ✓
2. **Placeholder scan:** No TBD/TODO in skill content. All code blocks contain actual content. ✓
3. **Type consistency:** `mysd:research`, `mysd:doc`, `mysd:spec`, `mysd:run` used consistently across all files. ✓
4. **Trigger boundaries:** Each SKILL.md description has explicit DO NOT USE clauses. ✓
5. **OpenSpec format:** Appendix A content correctly embedded in spec SKILL.md. ✓
6. **Analyzer logic:** All 4 dimensions (Coverage, Ambiguity, Consistency, Gaps) translated from Go to natural language in health-check.md. ✓
