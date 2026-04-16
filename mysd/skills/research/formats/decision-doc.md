# Decision Doc Template

This file is a template for producing a Decision Doc — the primary output of the research skill when a question qualifies as a "gray area" decision. Load this file when you need to format a final recommendation.

Fill in each section as described. Replace `{placeholders}` with actual content. Remove instructional notes (lines beginning with `>`).

---

## Template

```markdown
# Decision: {title}

> One-line statement of the decision to be made. Be specific — name the exact choice, not the general topic.

## Problem

> Describe the situation requiring a decision. Include:
> - What you are trying to accomplish
> - Why this decision matters (consequences of choosing poorly)
> - Any constraints, deadlines, or prior context that shape the choice

{describe the problem and context here}

## Gray Area Classification

> State which gray area category this falls into and why. Choose one:
> - **Multiple viable approaches with no consensus** — well-understood problem space, but practitioners disagree on the right answer
> - **Best practice exists but doesn't apply** — a standard recommendation exists, but specific constraints make it a poor fit here
> - **Must decide with incomplete info** — decision is time-sensitive; the information needed to be certain doesn't yet exist or can't be obtained

This qualifies as gray area because: {explain which category and why the standard playbook doesn't resolve it}

## Options

> List 2–4 options. Fewer than 2 options means there is no real decision. More than 4 usually means options need consolidation.
> For each option: provide concrete evidence (what you actually found — docs, benchmarks, real-world reports), not speculation.

### Option A: {name}

**Evidence**
> Concrete sources: official docs, benchmarks, known production usage, measured behavior. Do not list unverified claims.

- {evidence item}
- {evidence item}

**Pros**
- {pro}
- {pro}

**Cons**
- {con}
- {con}

**Effort:** {S / M / L}
> S = hours, M = days, L = weeks or requires structural changes

---

### Option B: {name}

**Evidence**
- {evidence item}
- {evidence item}

**Pros**
- {pro}

**Cons**
- {con}

**Effort:** {S / M / L}

---

### Option C: {name}  *(if applicable)*

**Evidence**
- {evidence item}

**Pros**
- {pro}

**Cons**
- {con}

**Effort:** {S / M / L}

---

## Recommendation

**Choose: Option {X} — {name}**

**Confidence:** {N}/10

> See Confidence Scale at the bottom of this template file.

**Reasoning**

> Explain why this option wins. Reference specific evidence from the Options section. State the key trade-off explicitly — what you are giving up and why that trade-off is acceptable.

{reasoning here}

**What would change my mind**

> Name specific new information or conditions that would cause you to pick a different option. This keeps the recommendation honest and makes it easier to revisit later.

- If {condition}, reconsider Option {X}
- If {condition}, the effort estimate for Option {X} changes significantly

## Open Questions

> List anything that remains unresolved after the recommendation. Include a suggested resolution path for each — who should answer it, what experiment would resolve it, or by what date it must be decided.

| Question | Suggested Resolution |
|----------|----------------------|
| {question} | {how to resolve — experiment, doc to read, person to ask, deadline} |
| {question} | {how to resolve} |
```

---

## Confidence Scale

Use this scale when setting the **Confidence** value in the Recommendation section.

| Score | Meaning |
|-------|---------|
| 1–3 | Guess. No direct evidence. High uncertainty. |
| 4–6 | Partial evidence. Significant risks or unknowns remain. |
| 7–8 | Multiple evidence sources. Risks identified and manageable. |
| 9–10 | Strong evidence. Nearly certain. 10 = verified with actual results. |

Avoid scores of exactly 5 or 10 by default — 5 signals you haven't done enough research, 10 means the question was never gray area. Most well-researched gray area decisions land at 6–8.
