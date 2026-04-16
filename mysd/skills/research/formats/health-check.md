# Spec Health Check — Analysis Format

This file defines how to perform a Spec Health Check across 4 analysis dimensions. Load this file when the user asks about spec quality, health, or completeness for a change. Follow each dimension's steps exactly and produce findings in the specified formats.

---

## How to Run

1. Identify the change directory. It is typically `openspec/changes/{change-name}/`.
2. Use Glob to discover spec files: `openspec/changes/{change-name}/specs/*/spec.md`
3. Use Read to read each artifact: `proposal.md`, `design.md`, `tasks.md`, and each `specs/{name}/spec.md`
4. Run all 4 dimensions in order: Coverage → Consistency → Ambiguity → Gaps
5. Skip any dimension whose required artifacts are missing (see skip rules per dimension)
6. Collect all findings, then produce the output summary

---

## Dimension 1: Coverage

**Purpose:** Every capability declared in `proposal.md` must have a corresponding spec file. A proposal that lists a capability with no spec is an unfulfilled promise.

**Required artifacts:** `proposal.md` + `specs/` directory

**Skip rule:** If `proposal.md` is missing or the `specs/` directory does not exist, mark this dimension as `Skipped (insufficient artifacts)` and produce no findings.

### Steps

1. Read `proposal.md`.
2. Locate capability sections. These are headings that match any of:
   - `## Capabilities`
   - `### New Capabilities`
   - `### Modified Capabilities`
3. Within each capability section, collect lines that match the pattern: `- \`capability-name\`: description`
   - The capability name is the text inside the backticks: a lowercase identifier containing only `a-z`, `0-9`, and `-`, starting with a letter.
   - Example: `- \`spec-health-check\`: Runs structural analysis on change artifacts` → name is `spec-health-check`
4. Stop collecting capability names when you reach the next `##`-level heading that is not a sub-heading of a capability section (i.e., does not start with `###`).
5. For each collected capability name, check whether `specs/{capability-name}/spec.md` exists.
6. If the file does not exist, emit a CRITICAL finding.

### Finding Format

```
[COV-N] CRITICAL — proposal.md: Capability '{name}' listed in proposal has no corresponding specs/{name}/spec.md
  Recommendation: Create specs/{name}/spec.md or remove '{name}' from proposal Capabilities
```

**N** increments sequentially across all COV findings (COV-1, COV-2, …).

---

## Dimension 2: Ambiguity

**Purpose:** Spec files must use RFC 2119 keywords (SHALL, SHALL NOT, MUST, MUST NOT) for requirements. Weak language like "should" or "TBD" leaves requirements open to interpretation and should be flagged.

**Required artifacts:** `specs/` directory with at least one `spec.md`

**Skip rule:** If the `specs/` directory does not exist or contains no `spec.md` files, mark this dimension as `Skipped (insufficient artifacts)` and produce no findings.

### Steps

1. For each `specs/{capability}/spec.md`:
   a. Read the file and split it into lines.
   b. Track code block state: when you encounter a line whose trimmed content starts with ` ``` `, toggle the code-block flag. Skip all lines inside a code block (lines between ` ``` ` fences).
   c. For each non-code-block line, check for the following weak patterns:

      | Pattern | Match rule |
      |---------|-----------|
      | `should` | Case-insensitive whole-word match: `/\bshould\b/i` |
      | `may` | Case-insensitive whole-word match: `/\bmay\b/i` |
      | `might` | Case-insensitive whole-word match: `/\bmight\b/i` |
      | `TBD` | Exact string match (case-sensitive): `TBD` |
      | `TODO` | Exact string match (case-sensitive): `TODO` |
      | `FIXME` | Exact string match (case-sensitive): `FIXME` |
      | `TKTK` | Exact string match (case-sensitive): `TKTK` |
      | `???` | Exact string match: `???` |

   d. Check patterns in the order listed above. Stop at the first match per line — emit at most one finding per line.
   e. Record the 1-based line number of the match.

2. Emit a SUGGESTION finding for each matched line.

### Finding Format

```
[AMB-N] SUGGESTION — specs/{capability}/spec.md:{line}: Vague language '{pattern}' found
  Recommendation: Replace '{pattern}' with SHALL/SHALL NOT for clarity
```

**N** increments sequentially across all AMB findings (AMB-1, AMB-2, …). `{line}` is the 1-based line number.

---

## Dimension 3: Consistency

**Purpose:** Every design decision documented in `design.md` (as a `###` heading) must be traceable to at least one task in `tasks.md`. Design decisions with no corresponding task are likely to be forgotten during implementation.

**Required artifacts:** `proposal.md` + `design.md` + `tasks.md`

**Skip rule:** If any of `proposal.md`, `design.md`, or `tasks.md` is missing, mark this dimension as `Skipped (insufficient artifacts)` and produce no findings.

### Steps

1. Read `design.md`.
2. Extract all `###`-level headings. A `###` heading is a line that starts with exactly `### ` (three hashes followed by a space). The heading text is everything after `### `, trimmed of leading/trailing whitespace.
   - Example: `### Error Handling Strategy` → heading text is `Error Handling Strategy`
   - Do not extract `####` or deeper headings (those belong to Gaps dimension).
3. Read `tasks.md` and convert its entire content to lowercase.
4. For each extracted heading text, convert it to lowercase and check whether it appears anywhere in the lowercased `tasks.md` content (substring match).
5. If the heading text is not found in `tasks.md`, emit a WARNING finding.

### Finding Format

```
[CON-N] WARNING — design.md: Design topic '{heading}' not referenced in tasks
  Recommendation: Verify tasks cover this design decision
```

**N** increments sequentially across all CON findings (CON-1, CON-2, …). `{heading}` is the original (non-lowercased) heading text.

---

## Dimension 4: Gaps

**Purpose:** Every requirement in a spec must have at least one concrete scenario, and every requirement name must appear somewhere in `tasks.md`. Requirements without scenarios cannot be verified; requirements without tasks will never be implemented.

**Required artifacts:** `specs/` directory + `tasks.md`

**Skip rule:** If the `specs/` directory does not exist or `tasks.md` is missing, mark this dimension as `Skipped (insufficient artifacts)` and produce no findings.

### Steps

#### Part A — Requirements must have scenarios

For each `specs/{capability}/spec.md`:

1. Read the file and split into lines.
2. Scan lines sequentially, tracking the current requirement:
   - A line matching `### Requirement: {name}` (pattern: `^###\s+Requirement:\s+(.+)`) opens a new requirement. The requirement name is the captured text, trimmed.
   - A line matching `#### Scenario: {name}` (pattern: `^####\s+Scenario:\s+(.+)`) marks that the current requirement has at least one scenario.
   - When a new `### Requirement:` line is encountered, first check whether the previous requirement had any scenario. If not, emit a WARNING.
3. After processing all lines, check the last open requirement. If it has no scenario, emit a WARNING.
4. Record the 1-based line number of the `### Requirement:` line for the finding location.

**Finding format (missing scenario):**
```
[GAP-N] WARNING — specs/{capability}/spec.md:{line}: Requirement '{name}' has no scenario
  Recommendation: Add at least one #### Scenario: under this requirement
```

#### Part B — Requirements must appear in tasks

1. Collect all `### Requirement: {name}` entries from every `specs/{capability}/spec.md` (same pattern as Part A).
2. Read `tasks.md` and convert its entire content to lowercase.
3. For each requirement name, convert it to lowercase and check whether it appears anywhere in the lowercased `tasks.md` content (substring match).
4. If the requirement name is not found in `tasks.md`, emit a WARNING.

**Finding format (missing task reference):**
```
[GAP-N] WARNING — specs/{capability}/spec.md: Requirement '{name}' has no matching task
  Recommendation: Add a task in tasks.md that references '{name}'
```

**N** for GAP findings is a single shared counter across Part A and Part B, incrementing globally (GAP-1, GAP-2, …). Process Part A for all specs first, then Part B.

---

## Output Summary Format

After running all 4 dimensions, produce a summary in this format:

```
## Spec Health Check — {change-name}

**Artifacts analyzed:** {comma-separated list of present artifacts}
**Artifacts missing:** {comma-separated list of missing artifacts, or "none"}

### Dimension Results

| Dimension   | Status                      | Findings |
|-------------|-----------------------------|----------|
| Coverage    | {Clean / N issue(s) found / Skipped (insufficient artifacts)} | {N} |
| Consistency | {Clean / N issue(s) found / Skipped (insufficient artifacts)} | {N} |
| Ambiguity   | {Clean / N issue(s) found / Skipped (insufficient artifacts)} | {N} |
| Gaps        | {Clean / N issue(s) found / Skipped (insufficient artifacts)} | {N} |

### Findings

{If no findings across all dimensions:}
No issues found. All analyzed dimensions are clean.

{If there are findings, list them grouped by dimension in order: COV, CON, AMB, GAP:}

#### Coverage
- [COV-1] CRITICAL — ...
  Recommendation: ...

#### Consistency
- [CON-1] WARNING — ...
  Recommendation: ...

#### Ambiguity
- [AMB-1] SUGGESTION — ...
  Recommendation: ...

#### Gaps
- [GAP-1] WARNING — ...
  Recommendation: ...
```

**Status values:**
- `Clean` — dimension ran and found zero issues
- `N issue(s) found` — dimension ran and found N issues (use exact count)
- `Skipped (insufficient artifacts)` — required artifacts were missing; dimension did not run

---

## Severity Reference

| Severity | Code prefix | Meaning |
|----------|-------------|---------|
| CRITICAL | COV | Spec is structurally broken — a declared capability has no spec file |
| WARNING | CON, GAP | A gap or inconsistency that will likely cause implementation problems |
| SUGGESTION | AMB | Weak language that should be strengthened but does not block progress |

Address CRITICAL findings before proceeding with implementation. WARNING findings should be resolved before review. SUGGESTION findings can be addressed incrementally.
