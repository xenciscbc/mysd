# Phase 11: agent-doc — Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions captured in CONTEXT.md — this log preserves the reasoning behind each decision.

**Date:** 2026-03-27
**Phase:** 11-agent-doc
**Mode:** discuss (update — existing CONTEXT.md revised)
**Areas discussed:** Plugin sync scope, propose auto-spec UX, docs_to_update reading, executor sidecar trigger

---

## Session Context

Existing CONTEXT.md from previous session had D-01 through D-16. User chose to update it. Codebase scout revealed 4 gray areas requiring clarification before planning. 0 pending todos matched Phase 11 scope.

---

## Gray Areas Presented

| # | Area | Rationale for Discussion |
|---|------|--------------------------|
| ① | Plugin sync scope | diff showed gsd-*.md agents, CLAUDE.md, subdirectories differ — scope unclear |
| ② | propose auto-spec UX | Step 11 interaction detail: what does user see after spec-writer runs? |
| ③ | docs_to_update reading | D-12 said SKILL.md layer but YAML parsing from bash is fragile |
| ④ | executor sidecar trigger | Who writes the sidecar — executor agent or apply orchestrator? |

All 4 areas selected by user for discussion.

---

## Decisions Made

### ① Plugin sync scope

**Decision:** Only sync `mysd-*.md` files. Exclude `gsd-*.md`, `CLAUDE.md`, and subdirectories (`gsd/`, `spectra/`).

**Evidence from codebase:**
- `.claude/agents/` has 18 `gsd-*.md` files not in `plugin/agents/` — these are GSD framework agents, not mysd agents
- `plugin/commands/` is missing `mysd-lang.md`, `mysd-model.md` (these ARE mysd commands)
- `mysd-designer.md` differs between the two sides (needs alignment)
- CLAUDE.md files differ — managed separately per directory

**Impact on D-15:** Explicitly scoped to `mysd-*.md` pattern. Adds specificity about which files need attention.

---

### ② propose auto-spec UX

**User preference:** After spec-writer completes, show spec content summary and prompt the user with available next step commands including their purpose descriptions.

**Decision:** Step 11 flow:
1. Invoke `mysd-spec-writer` via Task tool (same pattern as `mysd-discuss.md` Step 10)
2. Show generated spec content summary (MUST/SHOULD/MAY requirement counts + key points)
3. List available next commands with descriptions (e.g., `/mysd:plan` — 建立執行計劃、`/mysd:design` — 補充設計決策)

**Impact on D-01:** Added UX detail to the completion step. Spec summary display + next command menu.

---

### ③ docs_to_update reading mechanism

**Context:** ProjectConfig struct in `internal/config/defaults.go` does NOT currently have `docs_to_update` field. D-12 said SKILL.md layer with no binary subcommand.

**User decision:** Add `DocsToUpdate []string` to `ProjectConfig` struct.

**Rationale:** Binary config pattern is already established. Reading YAML directly from SKILL.md bash would be fragile (multi-line array parsing). Exposing via `mysd execute --context-only` JSON maintains consistency with how SKILL.md reads other config values.

**Impact on D-12:** Revised to clarify binary struct change is included. The doc update *logic* (LLM calls, file edits) stays in SKILL.md layer — only the config reading goes through binary.

**Scope note:** This is a small config struct extension, not a new subcommand. Phase 11 restriction "no binary subcommand" remains intact. Phase scope updated to "改動限於 SKILL.md 層、agent prompt 層，以及一個小幅 binary config struct 擴充".

---

### ④ executor sidecar trigger location

**Decision:** Sidecar writing logic belongs in `mysd-executor.md` agent prompt itself, in an on-failure paragraph after the Task Execution section.

**Rationale:**
- Executor agent has direct access to build/test error output in its context
- Apply orchestrator only sees "task returned failure" — no access to inner error details
- On-failure logic in agent prompt = most complete error context captured
- Consistent with executor being the "single responsible unit" for each task

**Impact on D-06:** Adds explicit "in agent prompt inner" clarification. Implementation site is `mysd-executor.md`, not `mysd-apply.md`.

---

## Codebase State (Scout findings)

| File | Current State | What Phase 11 Must Add |
|------|--------------|------------------------|
| mysd-propose.md | 10 steps, no auto-spec | Step 11: auto-spec + summary + next cmds |
| mysd-apply.md | 4 steps, no auto-verify | Step 5: go build + go test + verifier |
| mysd-archive.md | 2 steps, no doc update | Step 2: docs_to_update + LLM update |
| mysd-executor.md | 227 lines, no sidecar | on-failure paragraph + sidecar write |
| mysd-fix.md | Has sidecar read framework | Align path format with D-06 |
| plugin/commands/ | Missing lang, model files | Sync mysd-lang.md, mysd-model.md |
| internal/config/defaults.go | No DocsToUpdate | Add DocsToUpdate []string field |

---

## Unchanged Decisions

D-02, D-03, D-04, D-05 (apply auto-verify flow) — confirmed unchanged.
D-07, D-08, D-09 (sidecar format, fix reading, .gitignore) — confirmed unchanged.
D-10, D-11, D-13, D-14 (docs_to_update config and archive flow) — confirmed unchanged.
D-16 (plugin sync diff verification) — confirmed unchanged.

---

## Prior Session (Assumptions Mode — 2026-03-27 earlier)

Original CONTEXT.md was created via assumptions analyzer with 16 decisions D-01 through D-16. All assumptions were broadly correct; this session added implementation-level precision to D-01, D-06, D-12, and D-15.

---

## Session 3 (2026-03-27 — ff/ffe pipeline + docs config)

**Areas discussed:** ff/ffe inline verify+docs, archive.md config read mechanism, docs_to_update configuration UX

### New Gray Areas Identified

Codebase scout revealed that `ff.md` and `ffe.md` are inline orchestrators — they directly call `mysd archive` binary (not archive SKILL.md), meaning any logic added to archive SKILL.md will NOT be inherited by ff/ffe. Two implementation gaps:

1. **ff/ffe don't inherit auto-verify** — apply.md's new Step 5 won't propagate to ff/ffe
2. **ff/ffe don't inherit docs_to_update** — archive SKILL.md changes won't trigger in ff/ffe pipeline

### Decisions Made

**D-17 (new):** ff.md and ffe.md must inline both verify and docs logic independently:
- Verify: inserted after `mysd execute` state transition, before `mysd archive`
- docs_to_update: after `mysd archive` binary, call `mysd execute --context-only` to read `docs_to_update`, update inline with auto_mode=true

**D-18 (new):** archive.md reads docs_to_update by calling `mysd execute --context-only` at the start — consistent with D-12's mechanism, no new binary command needed.

**D-19 (new):** Add `mysd docs` command (thin wrapper, analogous to `mysd note`):
- User requested CLI management for docs_to_update, citing `/mysd:note` as the preferred pattern
- `mysd docs` / `mysd docs add <path>` / `mysd docs remove <path>`
- binary: `cmd/docs.go`, reads/writes `.claude/mysd.yaml` docs_to_update field
- SKILL.md: `.claude/commands/mysd-docs.md` thin wrapper (Bash only)

**Context usage statusline (deferred):** User asked about GSD-style context usage display with color bars. Out of Phase 11 scope — deferred to next milestone or quick task.

### Codebase State Updated

| File | Status |
|------|--------|
| mysd-ff.md | 5 steps, missing inline verify (D-17) + docs update (D-17) |
| mysd-ffe.md | 6 steps, missing inline verify (D-17) + docs update (D-17) |
| cmd/docs.go | Does not exist — new binary command (D-19) |
| .claude/commands/mysd-docs.md | Does not exist — new SKILL.md (D-19) |
