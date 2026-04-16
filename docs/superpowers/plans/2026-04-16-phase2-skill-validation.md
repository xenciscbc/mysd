# Phase 2: Skill Validation Testing

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Validate each of the 4 mysd skills against 15 test scenarios (9 happy path + 6 edge cases), confirming output format compliance and behavioral correctness.

**Architecture:** Each test scenario creates a fixture (or uses existing project state), invokes the skill by reading its SKILL.md and following its instructions, then validates the output against structural pass criteria. Tests run against the real mysd project, not mocks.

**Tech Stack:** Claude Code SKILL.md, OpenSpec fixtures, git

**Validation matrix source:** Design doc at `~/.gstack/projects/xenciscbc-mysd/cbc-master-design-20260416-113346.md`

---

## File Structure

```
mysd-skills/                           # Skills under test (already created)
testdata/
  validation/
    research/                          # Research skill test fixtures
      tech-decision/question.md        # Test scenario input
      health-check-change/             # Simulated change dir for health check
        proposal.md
        specs/test-cap/spec.md
        design.md
        tasks.md
    doc/
      new-command-diff.patch           # Simulated git diff for new command scenario
    spec/
      reverse-spec-input/             # Go files for reverse-spec test
        cmd/sample.go
    results.md                         # Validation results log
```

---

### Task 1: Create test fixtures

**Files:**
- Create: `testdata/validation/research/tech-decision/question.md`
- Create: `testdata/validation/research/health-check-change/proposal.md`
- Create: `testdata/validation/research/health-check-change/specs/test-feature/spec.md`
- Create: `testdata/validation/research/health-check-change/specs/missing-feature/` (empty dir, no spec.md)
- Create: `testdata/validation/research/health-check-change/design.md`
- Create: `testdata/validation/research/health-check-change/tasks.md`
- Create: `testdata/validation/spec/reverse-spec-input/cmd/sample.go`
- Create: `testdata/validation/results.md`

- [ ] **Step 1: Create the fixture directories**

```bash
mkdir -p testdata/validation/research/tech-decision
mkdir -p testdata/validation/research/health-check-change/specs/test-feature
mkdir -p testdata/validation/research/health-check-change/specs/missing-feature
mkdir -p testdata/validation/doc
mkdir -p testdata/validation/spec/reverse-spec-input/cmd
```

- [ ] **Step 2: Create the tech decision fixture**

Write `testdata/validation/research/tech-decision/question.md`:

```markdown
# Test Scenario: Tech Decision

## Question
Should the mysd-skills plugin use a single README.md or per-skill README files?

## Context
- 4 skills in one plugin
- Each skill is independently installable
- Users may only install 1-2 skills, not all 4
- A single README covers everything but is long
- Per-skill READMEs are discoverable but duplicated

## Expected from research skill
- Classify as gray area (multiple viable approaches, no consensus)
- Frame 2+ options with evidence
- Produce a Decision Doc with confidence score
```

- [ ] **Step 3: Create the health check change fixture — proposal.md**

Write `testdata/validation/research/health-check-change/proposal.md`:

```markdown
---
spec-version: "1.0"
change: "test-validation"
status: "proposed"
created: "2026-04-16"
updated: "2026-04-16"
---

## Summary
Test change for validation of Spec Health Check.

## Capabilities

- `test-feature`: A feature that has a spec file
- `missing-feature`: A feature that intentionally has NO spec file (should trigger COV finding)
- `orphan-feature`: A feature listed here but no spec dir exists at all
```

- [ ] **Step 4: Create the test-feature spec**

Write `testdata/validation/research/health-check-change/specs/test-feature/spec.md`:

```markdown
---
spec-version: "1.0"
capability: "test-feature"
delta: ADDED
status: done
---

## Requirement: Test Feature Works

The test feature MUST produce correct output.

It should also handle edge cases gracefully.

### Scenario: Basic Operation

WHEN the user invokes the test feature
THEN it produces output
```

Note: Line "It should also handle edge cases gracefully." uses lowercase "should" — this should trigger an AMB finding.

- [ ] **Step 5: Create design.md with headings**

Write `testdata/validation/research/health-check-change/design.md`:

```markdown
---
spec-version: "1.0"
change: "test-validation"
status: "designed"
---

### Data Model
Use a simple struct to hold test data.

### Caching Strategy
Cache results in memory for the session.

### Orphan Design Topic
This heading is intentionally NOT referenced in tasks.md — should trigger a CON finding.
```

- [ ] **Step 6: Create tasks.md referencing some (not all) items**

Write `testdata/validation/research/health-check-change/tasks.md`:

```markdown
---
spec-version: "1.0"
total: 2
completed: 0
tasks:
  - id: 1
    name: "Implement data model"
    status: pending
    spec: "test-feature"
  - id: 2
    name: "Implement caching strategy"
    status: pending
    spec: "test-feature"
---

## Task 1: Implement data model
Implement the Data Model as described in the design.

## Task 2: Implement caching strategy
Implement the Caching Strategy per the design.
```

Note: "Orphan Design Topic" is NOT referenced → should trigger CON finding. "Test Feature Works" requirement name is NOT in tasks → should trigger GAP finding.

- [ ] **Step 7: Create reverse-spec Go file**

Write `testdata/validation/spec/reverse-spec-input/cmd/sample.go`:

```go
package cmd

// SampleCommand executes the sample action.
// It validates input, processes the data, and returns a result.
func SampleCommand(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input must not be empty")
	}
	return "processed: " + input, nil
}

// SampleConfig holds configuration for the sample command.
type SampleConfig struct {
	Verbose bool   `yaml:"verbose"`
	Output  string `yaml:"output"`
}
```

- [ ] **Step 8: Create results template**

Write `testdata/validation/results.md`:

```markdown
# Phase 2 Validation Results

Date: 2026-04-16
Skills under test: mysd-skills v1.0.0

## Research Skill

| # | Scenario | Pass Criteria | Result | Notes |
|---|----------|---------------|--------|-------|
| 1 | Tech decision | Decision Doc format, confidence score | | |
| 2 | Architecture decision | Decision Doc format, 2+ options with evidence | | |
| 3 | Vague requirement | Structured options from vague input | | |
| 4 | Spec Health Check | 4 dimensions, summary format | | |
| 5 | Non-gray-area question | Rejects gray area, answers directly | | |

## Doc Skill

| # | Scenario | Pass Criteria | Result | Notes |
|---|----------|---------------|--------|-------|
| 1 | New command | Identifies README.md + README.zh-TW.md | | |
| 2 | Bug fix | Identifies CHANGELOG.md | | |
| 3 | Multi-language sync | Both READMEs updated identically | | |
| 4 | Unmapped change type | Falls back to grep heuristic | | |
| 5 | Empty diff | Reports "no changes detected" | | |

## Spec Skill

| # | Scenario | Pass Criteria | Result | Notes |
|---|----------|---------------|--------|-------|
| 1 | New feature spec | 4 required frontmatter fields, RFC 2119 | | |
| 2 | Modify existing spec | delta=MODIFIED, existing content preserved | | |
| 3 | Reverse spec from code | Infers capability from file path, scenarios | | |
| 4 | RENAMED delta | delta=RENAMED used correctly | | |
| 5 | Incompatible spec-version | Warns or handles gracefully | | |

## Summary

- Total: 15 scenarios
- Pass: _/15
- Fail: _/15
- Skills ready for Phase 3: [ ] research [ ] doc [ ] spec
```

- [ ] **Step 9: Commit fixtures**

```bash
git add testdata/validation/
git commit -m "test: add Phase 2 validation fixtures"
```

---

### Task 2: Validate research skill (5 scenarios)

**Files:**
- Read: `mysd-skills/research/SKILL.md`
- Read: `mysd-skills/research/formats/decision-doc.md`
- Read: `mysd-skills/research/formats/health-check.md`
- Read: `testdata/validation/research/tech-decision/question.md`
- Read: `testdata/validation/research/health-check-change/*`
- Modify: `testdata/validation/results.md`

- [ ] **Step 1: Test Scenario R1 — Tech Decision (happy path)**

1. Read `mysd-skills/research/SKILL.md`
2. Read `testdata/validation/research/tech-decision/question.md` for the scenario
3. Follow the SKILL.md flow with this question: "Should the mysd-skills plugin use a single README.md or per-skill README files?"
4. Validate the output:
   - [ ] Contains `# Decision:` heading
   - [ ] Contains `## Gray Area Classification` with one of the 3 categories
   - [ ] Contains `## Options` with 2+ options
   - [ ] Each option has Evidence, Pros, Cons, Effort
   - [ ] Contains `## Recommendation` with Confidence N/10
   - [ ] Contains `What would change my mind`
5. Record result in `testdata/validation/results.md` row R1

**Pass criteria:** Output matches Decision Doc template structure. Confidence score is 1-10 integer.

- [ ] **Step 2: Test Scenario R2 — Architecture Decision (happy path)**

1. Follow the research SKILL.md flow with this question: "For the orchestrator, should subagents read the SKILL.md themselves or should the orchestrator embed the instructions in the prompt?"
2. Validate same criteria as R1
3. Additionally verify: options reference actual tradeoffs (context window size vs. instruction freshness)
4. Record result in results.md row R2

**Pass criteria:** Decision Doc format, 2+ options with concrete evidence.

- [ ] **Step 3: Test Scenario R3 — Vague Requirement (happy path)**

1. Follow the research SKILL.md flow with this vague input: "Make the spec writer better."
2. Validate:
   - [ ] Skill classifies this as gray area (category a: multiple approaches, no consensus)
   - [ ] Frames concrete options (e.g., "add more validation", "improve reverse-spec", "add templates")
   - [ ] Each option has evidence from the actual codebase
3. Record result in results.md row R3

**Pass criteria:** Transforms vague input into structured options with evidence.

- [ ] **Step 4: Test Scenario R4 — Spec Health Check (edge case)**

1. Read `mysd-skills/research/SKILL.md` — follow Spec Health Check Mode
2. Read `mysd-skills/research/formats/health-check.md`
3. Run the 4-dimension analysis against `testdata/validation/research/health-check-change/`
4. Validate expected findings:
   - [ ] **Coverage:** COV-1 CRITICAL for `missing-feature` (listed in proposal, no spec.md in its dir), COV-2 CRITICAL for `orphan-feature` (listed in proposal, no spec dir at all)
   - [ ] **Ambiguity:** AMB-1 SUGGESTION for lowercase "should" in test-feature/spec.md
   - [ ] **Consistency:** CON-1 WARNING for "Orphan Design Topic" not in tasks.md
   - [ ] **Gaps:** GAP-1 WARNING for "Test Feature Works" requirement not in tasks.md
5. Validate output format matches the summary template in health-check.md
6. Record result in results.md row R4

**Pass criteria:** All 4 dimensions produce findings. Output matches summary format. Finding IDs are sequential.

- [ ] **Step 5: Test Scenario R5 — Non-Gray-Area Question (edge case)**

1. Follow the research SKILL.md flow with this question: "What is the Go fmt package used for?"
2. Validate:
   - [ ] Step 1 (Classify) determines this is NOT gray area (clear official docs answer)
   - [ ] Skill answers directly WITHOUT producing a Decision Doc
   - [ ] Response says something like "This isn't a gray area" or just answers the question
3. Record result in results.md row R5

**Pass criteria:** Skill correctly rejects gray area classification. No Decision Doc produced.

- [ ] **Step 6: Commit research results**

```bash
git add testdata/validation/results.md
git commit -m "test: research skill validation complete"
```

---

### Task 3: Validate doc skill (5 scenarios)

**Files:**
- Read: `mysd-skills/doc/SKILL.md`
- Read: `README.md`, `README.zh-TW.md`
- Modify: `testdata/validation/results.md`

**Important:** This task tests the doc skill's ANALYSIS capabilities only (Steps 1-2: change detection and impact analysis). Do NOT actually apply doc changes to the repo — just validate that the skill correctly identifies which files need updating.

- [ ] **Step 1: Test Scenario D1 — New Command (happy path)**

1. Read `mysd-skills/doc/SKILL.md`
2. Simulate: a commit that adds `cmd/newfeature.go` to the repo. Use the actual git log to find a commit that added a cmd/ file, or use this diff context: "new file cmd/newfeature.go with 50 lines"
3. Follow the SKILL.md Step 2 (Impact Analysis) with this change
4. Validate:
   - [ ] Impact mapping identifies change type as "New/removed command"
   - [ ] Lists `README.md` as needing update
   - [ ] Lists `README.zh-TW.md` as needing update
   - [ ] Lists `CLAUDE.md` as needing update (it has a command list)
5. Record result in results.md row D1

**Pass criteria:** Correctly identifies all 3 docs that need updating for a new command.

- [ ] **Step 2: Test Scenario D2 — Bug Fix (happy path)**

1. Simulate: a commit modifying `internal/spec/parser.go` with commit message "fix: correct frontmatter parsing for empty files"
2. Follow the SKILL.md Step 2 with this change
3. Validate:
   - [ ] Impact mapping identifies change type as "Bug fix" (pattern: .go file + "fix" in commit)
   - [ ] Lists `CHANGELOG.md` as needing update
4. Record result in results.md row D2

**Pass criteria:** Identifies CHANGELOG.md as the doc to update for a bug fix.

- [ ] **Step 3: Test Scenario D3 — Multi-Language Sync (happy path)**

1. Simulate: the skill has identified README.md needs updating (from D1 scenario)
2. Follow the SKILL.md Step 4 (Multi-Language Sync)
3. Validate:
   - [ ] Detects `README.zh-TW.md` as a locale variant of `README.md`
   - [ ] Plans to generate equivalent update in Traditional Chinese
5. Record result in results.md row D3

**Pass criteria:** Correctly discovers and plans sync for README.zh-TW.md.

- [ ] **Step 4: Test Scenario D4 — Unmapped Change Type (edge case)**

1. Simulate: a commit modifying `Makefile` (not in the impact mapping table)
2. Follow the SKILL.md Step 2 with this change
3. Validate:
   - [ ] Change type does NOT match any row in the impact mapping table
   - [ ] Falls back to grep heuristic: searches for "Makefile" or "make" in .md files
   - [ ] Identifies any .md files that reference Makefile/make commands
4. Record result in results.md row D4

**Pass criteria:** Falls back to grep heuristic when mapping table doesn't match. Reports what it found.

- [ ] **Step 5: Test Scenario D5 — Empty Diff (edge case)**

1. Run the SKILL.md flow with no git changes (clean working directory)
2. Step 1 runs `git diff --name-only HEAD~1` — if this returns files, use `HEAD..HEAD` instead (guaranteed empty)
3. Validate:
   - [ ] Skill detects empty diff
   - [ ] Reports "No changes detected" or similar
   - [ ] Does NOT proceed to Step 2
4. Record result in results.md row D5

**Pass criteria:** Gracefully handles empty diff. Does not attempt to update any docs.

- [ ] **Step 6: Commit doc results**

```bash
git add testdata/validation/results.md
git commit -m "test: doc skill validation complete"
```

---

### Task 4: Validate spec skill (5 scenarios)

**Files:**
- Read: `mysd-skills/spec/SKILL.md`
- Read: `openspec/specs/*/spec.md` (existing specs)
- Read: `testdata/validation/spec/reverse-spec-input/cmd/sample.go`
- Modify: `testdata/validation/results.md`

**Important:** Write spec output to `testdata/validation/spec/output/` — do NOT modify the actual `openspec/` directory.

- [ ] **Step 1: Test Scenario S1 — New Feature Spec (happy path)**

1. Read `mysd-skills/spec/SKILL.md`
2. Ask the skill to write a spec for a new capability called "batch-processing" that allows running multiple commands in sequence
3. Follow the SKILL.md flow (Steps 1-6)
4. Write output to `testdata/validation/spec/output/batch-processing/spec.md`
5. Validate:
   - [ ] Frontmatter has `spec-version: "1.0"`
   - [ ] Frontmatter has `capability: "batch-processing"`
   - [ ] Frontmatter has `delta: ADDED`
   - [ ] Frontmatter has `status:` field (pending or done)
   - [ ] Body has at least one `## Requirement:` heading
   - [ ] Requirements use RFC 2119 UPPERCASE keywords (MUST, SHOULD, MAY)
   - [ ] At least one `### Scenario:` with WHEN/THEN format
6. Record result in results.md row S1

**Pass criteria:** All 4 required frontmatter fields present and valid. RFC 2119 keywords uppercase. Scenario format correct.

- [ ] **Step 2: Test Scenario S2 — Modify Existing Spec (happy path)**

1. Read an existing spec: `openspec/specs/spec-authoring/spec.md`
2. Ask the skill to add a new requirement: "Proposal frontmatter MUST include a `type` field indicating Feature/BugFix/Refactor"
3. Follow the SKILL.md flow
4. Write output to `testdata/validation/spec/output/spec-authoring-modified/spec.md`
5. Validate:
   - [ ] Frontmatter has `delta: MODIFIED`
   - [ ] Original requirements are preserved (not deleted)
   - [ ] New requirement is added with RFC 2119 MUST keyword
   - [ ] New requirement has a scenario
6. Record result in results.md row S2

**Pass criteria:** delta=MODIFIED, existing content preserved, new requirement properly formatted.

- [ ] **Step 3: Test Scenario S3 — Reverse Spec from Code (happy path)**

1. Read `testdata/validation/spec/reverse-spec-input/cmd/sample.go`
2. Follow the SKILL.md Step 5 (Reverse-Spec Rules):
   - Read the Go file
   - Identify exported symbols: `SampleCommand`, `SampleConfig`
   - Map path `cmd/` → command behavior capability
   - Infer requirements from function signatures and doc comments
3. Write output to `testdata/validation/spec/output/sample-command/spec.md`
4. Validate:
   - [ ] Inferred capability name relates to "sample" or "sample-command"
   - [ ] delta is ADDED (new spec from code)
   - [ ] Requirements inferred: input validation (empty string check), output format ("processed: ...")
   - [ ] Scenario has WHEN/THEN that matches the function behavior
5. Record result in results.md row S3

**Pass criteria:** Correctly infers capability from file path, generates requirements from code behavior, produces scenarios.

- [ ] **Step 4: Test Scenario S4 — RENAMED Delta (edge case)**

1. Ask the skill to handle a rename: capability "artifact-analysis" was renamed to "spec-health-check"
2. Follow the SKILL.md flow
3. Write output to `testdata/validation/spec/output/spec-health-check-renamed/spec.md`
4. Validate:
   - [ ] Frontmatter has `delta: RENAMED`
   - [ ] Content references the rename (old name → new name)
   - [ ] Status and other fields are valid
5. Record result in results.md row S4

**Pass criteria:** delta=RENAMED used correctly. Rename is documented in the spec content.

- [ ] **Step 5: Test Scenario S5 — Incompatible spec-version (edge case)**

1. Create a spec file with `spec-version: "0.5"` (non-standard version)
2. Ask the skill to modify this spec (add a requirement)
3. Follow the SKILL.md flow
4. Validate:
   - [ ] Skill detects the non-standard version
   - [ ] Either: warns the user about version incompatibility, OR upgrades to "1.0", OR handles gracefully
   - [ ] Does NOT silently produce output with `spec-version: "0.5"` — must acknowledge the discrepancy
5. Record result in results.md row S5

**Pass criteria:** Handles non-standard spec-version gracefully. Does not silently pass through.

- [ ] **Step 6: Commit spec results**

```bash
mkdir -p testdata/validation/spec/output
git add testdata/validation/
git commit -m "test: spec skill validation complete"
```

---

### Task 5: Compile results and fix issues

**Files:**
- Read: `testdata/validation/results.md`
- Modify: `testdata/validation/results.md` (fill in summary)
- Potentially modify: any SKILL.md that failed a test

- [ ] **Step 1: Read all validation results**

Read `testdata/validation/results.md` and count pass/fail across all 15 scenarios.

- [ ] **Step 2: Fill in the summary section**

Update the Summary at the bottom of results.md:
```markdown
## Summary

- Total: 15 scenarios
- Pass: {N}/15
- Fail: {N}/15
- Skills ready for Phase 3: [x] research [x] doc [x] spec
```

- [ ] **Step 3: If any failures, identify the root cause**

For each failed scenario:
1. Determine if the failure is in the SKILL.md instructions (unclear, ambiguous, wrong) or in the test fixture (bad setup)
2. If SKILL.md issue: describe the fix needed
3. If fixture issue: describe the correction

- [ ] **Step 4: Fix any SKILL.md issues found**

For each SKILL.md fix:
1. Edit the specific SKILL.md file
2. Re-run the failed test scenario
3. Update results.md with the re-test result

- [ ] **Step 5: Commit final results**

```bash
git add testdata/validation/ mysd-skills/
git commit -m "test: Phase 2 validation complete — all scenarios pass"
```

If some scenarios still fail:

```bash
git commit -m "test: Phase 2 validation complete — N/15 pass, issues documented"
```

---

## Self-Review Checklist

1. **Spec coverage:** All 15 scenarios from the design doc validation matrix are covered (R1-R5, D1-D5, S1-S5) ✓
2. **Placeholder scan:** Every step has concrete validation criteria (checkbox lists), not "verify it works" ✓
3. **Type consistency:** Fixture file paths match across tasks. `results.md` row references (R1, D1, S1) are consistent ✓
4. **Test isolation:** Tests write to `testdata/validation/` only, never modify actual `openspec/` or project docs ✓
5. **Fixtures are deterministic:** Health check fixture has intentional errors (lowercase "should", missing spec, orphan heading) that produce predictable findings ✓
