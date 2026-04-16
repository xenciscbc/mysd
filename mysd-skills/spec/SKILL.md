---
name: mysd:spec
description: >
  Write and update OpenSpec format spec files based on code changes or plans.
  Generates correct YAML frontmatter, RFC 2119 requirements, and scenario definitions.
  DO NOT use for documentation updates (use mysd:doc) or research/decisions
  (use mysd:research).
---

## When to Use

USE this skill when:
- The user says "update spec", "write spec", "sync spec"
- Post-implementation spec reflection — code is done, spec needs to catch up
- Converting a proposal into a spec file
- A new capability was added and has no spec yet

DO NOT USE when:
- The file is README, CHANGELOG, or any prose documentation → use `mysd:doc`
- The task is research, analysis, or a technical decision → use `mysd:research`
- The user asks about spec quality or completeness → use `mysd:research` (health check mode)

---

## OpenSpec Format Reference

### Spec File (`openspec/specs/{capability}/spec.md` or `openspec/changes/{name}/specs/{capability}/spec.md`)

```yaml
---
spec-version: "1.0"          # required
capability: {name}           # required — kebab-case identifier
delta: ADDED                 # required — ADDED | MODIFIED | REMOVED | RENAMED
status: pending              # required — pending | in_progress | done | blocked
name: Human Name             # optional
description: One-liner.      # optional
version: "0.1.0"             # optional
generatedBy: mysd:spec       # optional
---
```

### Proposal File (`openspec/changes/{name}/proposal.md`)

```yaml
---
spec-version: "1.0"
change: {change-name}
status: pending
created: YYYY-MM-DD
updated: YYYY-MM-DD
---
```

### Tasks File (`openspec/changes/{name}/tasks.md`) — V2

```yaml
---
spec-version: "1.0"
total: N
completed: N
tasks:
  - id: 1
    name: Task name
    status: pending          # pending | in_progress | done | blocked
    spec: capability-name    # optional — spec directory this task belongs to
    depends: [2, 3]          # optional — task IDs this depends on
    files: [path/to/file]    # optional — files touched
    satisfies: [REQ-01]      # optional — requirement IDs satisfied
---
```

### RFC 2119 Keywords

Use UPPERCASE only: **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, **MAY**

```markdown
## Requirements

### REQ-01: Descriptive title
The system MUST do X when Y.

### REQ-02: Another requirement
The system SHOULD do Z unless W.
```

### Scenario Format

```markdown
### Scenario: {name}

WHEN {trigger or condition}
THEN {expected outcome}
AND {additional assertion}
```

Each requirement MUST have at least one scenario.

---

## Flow

### Step 1: Change Context

Read the available context — in priority order:
1. `openspec/changes/{name}/proposal.md` — stated intent
2. `openspec/changes/{name}/tasks.md` — task breakdown
3. Staged/committed `.go` file diffs — actual implementation
4. User's message — ad-hoc description

If none of the above are present, ask the user to describe the change before continuing.

### Step 2: Spec Discovery

Run `find openspec/specs -name "spec.md"` and, if a change directory exists,
`find openspec/changes/{name}/specs -name "spec.md"` as well.

Read the frontmatter of each result. Match on `capability` field to determine which
specs already cover the change's scope.

### Step 3: Gap Analysis

For each capability touched by the change, classify:

| Signal | Delta |
|---|---|
| No existing spec.md for this capability | ADDED |
| Spec exists and requirements are changing | MODIFIED |
| Capability is being deleted | REMOVED |
| Capability is being renamed | RENAMED |

If RENAMED, record both the old and new capability names.

### Step 4: Spec Generation

For each spec to write or update:

1. **Frontmatter** — fill all four required fields; set `generatedBy: mysd:spec`
2. **Requirements** — number sequentially (REQ-01, REQ-02, …); use RFC 2119 UPPERCASE
3. **Scenarios** — at least one `### Scenario:` per requirement; use WHEN/THEN/AND

Keep language declarative and implementation-neutral unless the spec is intentionally
implementation-specific.

### Step 5: Reverse-Spec Rules (no proposal/plan — code-first)

When context comes only from changed Go files:

1. Read each changed `.go` file; identify exported functions, types, and constants
2. Map by path:
   - `cmd/` → command behavior (user-facing)
   - `internal/{pkg}/` → internal capability named after the package
3. From function signatures and doc comments, infer:
   - MUST requirements for error conditions and guaranteed behaviors
   - WHEN/THEN scenarios from call sites and test cases
4. Set delta to ADDED (new export) or MODIFIED (changed signature/behavior)
5. Flag inferred requirements with a `<!-- inferred -->` comment so reviewers can verify

### Step 6: Validation Checklist

Before writing, confirm each spec satisfies:

- [ ] Frontmatter has all four required fields (`spec-version`, `capability`, `delta`, `status`)
- [ ] All RFC 2119 keywords are UPPERCASE
- [ ] Every requirement has at least one scenario
- [ ] File path matches the correct directory convention (main spec vs. delta spec)
- [ ] `capability` value is kebab-case and matches the directory name

Fail any spec that misses a required field — do not silently omit it.

### Step 6.5: Confirm Scope

Use AskUserQuestion to present the specs to be written/updated and let the user choose:

> 根據分析，以下 spec 需要處理：
>
> [列出每個 spec 的 capability name、delta type、目標路徑]
>
> A) 全部寫入
> B) 只處理部分 — [讓使用者選擇要處理哪些]
> C) 跳過 — 這次不寫 spec

如果使用者選 C，結束。選 A 或 B，帶著選定的 spec 清單繼續到 Step 7。

### Step 7: Output

- **New spec**: use Write tool to create the file
- **Existing spec**: use Edit tool to apply targeted changes; show old → new blocks
- After all writes, list each file path and its delta type (ADDED / MODIFIED / etc.)
