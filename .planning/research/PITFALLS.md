# Pitfalls Research

**Domain:** Adding interactive discovery, worktree-based parallel execution, subagent orchestration, and language-agnostic scanning to an existing Go CLI tool (my-ssd v1.1)
**Researched:** 2026-03-25
**Confidence:** HIGH (Critical pitfalls — verified via official docs + multiple sources), MEDIUM (Integration gotchas — multiple community sources), MEDIUM (Performance traps)

> This file extends the v1.0 pitfalls research. It focuses specifically on pitfalls when **adding** these features to an already-shipped system. Pitfalls from v1.0 research (spec format lock-in, AI self-verification blindness, spec drift, context budget overflow, multi-agent shared state, Claude Code API breaking changes) remain valid and are not repeated here.

---

## Critical Pitfalls

### Pitfall 1: Subagent Cannot Spawn Subagents — Nested Orchestration Fails Silently

**What goes wrong:**
The proposal defines a hierarchy where orchestrator (skill) manages interaction and subagents do work. But subagents in Claude Code cannot themselves spawn other subagents. If `mysd-researcher` attempts to spawn `mysd-advisor` instances internally, the spawn silently fails or throws an error inside the subagent's isolated context — the parent orchestrator receives no error signal, just a truncated or incorrect result.

**Why it happens:**
The proposal describes `mysd-advisor x N` running in parallel for each gray area. This appears to require `mysd-researcher` to spawn multiple `mysd-advisor` instances — but that is nested spawning. Claude Code strictly forbids subagents from spawning subagents. Only the top-level orchestrator (the SKILL.md session) can spawn subagents.

**How to avoid:**
The orchestrator (SKILL.md) must be the exclusive spawner. Design the flow as: SKILL.md spawns `mysd-researcher` → researcher returns a list of gray areas → SKILL.md then spawns `mysd-advisor x N` in parallel for each gray area. The researcher must not attempt to spawn advisors itself. Document this constraint explicitly in each subagent's SKILL.md with: "You MUST NOT use the Task tool to spawn additional agents."

**Warning signs:**
- Any subagent's instructions mention "spawn", "launch", or "use the Task tool" for another agent
- The researcher agent's instructions say it produces advisors or delegates to them
- Proposal language like "researcher spawns advisors in parallel" — this must be re-read as "researcher produces gray area list, orchestrator spawns advisors"

**Phase to address:**
Phase 1 (Subagent architecture design) — this is a structural constraint that shapes every agent definition in the milestone.

---

### Pitfall 2: Interactive Discovery Loop Has No Deterministic Exit Condition

**What goes wrong:**
The dual-loop design (area deepening + new area discovery) can become a user-trapping cycle with no clear exit. The AI discovers new areas, the user explores them, the AI discovers more areas from the exploration — the loop never terminates. Users abandon the session mid-way, leaving proposal.md in a half-written state. Alternatively, the AI asks one more clarifying question for every answer, producing a questioning spiral that frustrates users before the proposal is written.

**Why it happens:**
Adaptive questioning without an explicit termination condition is a known failure mode. The feature description says "until the user is satisfied" — but AI agents do not naturally know when "satisfied" is reached. Without a hard limit or a structured satisfaction signal, the loop defaults to generating more questions.

**How to avoid:**
Implement a hard maximum: no more than 3 gray areas per research session; no more than 3 depth questions per area. After each area completes, ask explicitly "Are you ready to proceed, or explore another area? (remaining: N)" — make the counter visible. After all areas are explored once, present a summary and ask "Anything else, or shall I write the proposal?" — binary choice, not open-ended. The `--auto` flag must skip all loops and use the AI's first recommendations directly.

**Warning signs:**
- The SKILL.md for `propose` or `discuss` contains no maximum question count
- The exit condition is described as "when the user says they are done" with no backup limit
- Test run with `--auto` results in any interactive prompts being shown

**Phase to address:**
Phase 1 (Interactive discovery design) — the loop termination contract must be in the SKILL.md frontmatter and verified before any end-to-end test.

---

### Pitfall 3: Git Worktree Disk Explosion on Large Codebases

**What goes wrong:**
Each worktree is a full working-directory copy. A 2 GB codebase running 4 parallel executor tasks creates 8+ GB of additional disk usage in `.worktrees/`. Go's build cache (`$GOPATH/pkg/mod`, `GOCACHE`) is duplicated per worktree. On developer machines with limited SSD space (256 GB is common), a wave execution can fill the disk mid-run, causing the Go binary or Git to fail with cryptic I/O errors rather than clean "out of space" messages.

**Why it happens:**
The proposal correctly chose "complete copy (Plan A)" for worktrees over symlink-based approaches. Plan A is the right call for path correctness, but it multiplies disk usage proportionally to parallel task count. Developers building on large monorepos (node_modules, vendor/, generated code) do not anticipate this.

**How to avoid:**
Before creating worktrees, check available disk space; abort with a clear error if free space is less than (codebase size × task count × 1.5). Exclude `.git/` (already shared by worktree design), `vendor/`, `node_modules/`, and Go build cache from worktree disk estimates. Create a `go env GOCACHE` redirect: set `GOCACHE` per worktree to a temp directory, then clean it on worktree removal. Cap the default parallel task count at 4 regardless of available CPUs. Log disk usage after each worktree is created.

**Warning signs:**
- Worktree creation code does not call `syscall.Statfs` (or `golang.org/x/sys/windows.GetDiskFreeSpaceEx` on Windows) before proceeding
- The executor cleanup routine does not verify that the worktree directory was actually removed
- CI runs fail with "no space left on device" after adding wave execution tests

**Phase to address:**
Phase 2 (Worktree execution engine) — disk guard must be the first check in the worktree creation flow.

---

### Pitfall 4: Windows MAX_PATH Silently Truncates Worktree Paths

**What goes wrong:**
Windows enforces a 260-character `MAX_PATH` limit by default. The worktree path `.worktrees/T{id}/{change-name}/{project-path}` compounds quickly: `C:\Users\developer\projects\my-company-webapp\.worktrees\T12\interactive-discovery\internal\executor\alignment_test.go` can easily reach 280+ characters. Git operations inside the worktree fail with "Filename too long" — but Go's `exec.Command` running `git` reports this as exit code 1 with a non-descriptive error, not a path-limit error. The user sees "execution failed" with no indication that a Windows path setting is to blame.

**Why it happens:**
The proposal specifically chose short paths (`.worktrees/T{id}/`) to mitigate this — that is correct thinking. But the fix only shortens the worktree root; it does not control the file paths inside the checked-out code. Files with long package paths (e.g., `internal/spec/parser/testdata/openspec_v2/changes/...`) can still breach the limit when combined with any project root longer than ~100 characters.

**How to avoid:**
On Windows, detect `GOOS == "windows"` at worktree creation time and run `git config core.longpaths true` in the worktree before any file operations. Additionally, check if the Windows `LongPathsEnabled` registry key is set; if not, print a one-time warning with instructions (not a hard failure — the user may have already enabled it system-wide). Keep worktree root names as `T{id}` with no change name appended (the proposal already does this — enforce it strictly, do not let the change name appear in the worktree path).

**Warning signs:**
- Worktree creation does not run `git config core.longpaths true` on Windows
- The error message when git fails inside a worktree does not mention Windows path limits as a possible cause
- Integration tests on Windows fail with "exit status 1" from git commands but pass on macOS/Linux

**Phase to address:**
Phase 2 (Worktree execution engine) — Windows path handling must be verified in CI with a Windows runner before the phase is marked complete.

---

### Pitfall 5: Orphaned Worktrees and Branches After Interrupted Execution

**What goes wrong:**
If the parent process is killed (Ctrl+C, OOM, CI timeout), the cleanup routine in `defer` blocks may not run. The `.worktrees/T{id}/` directories remain on disk. The `mysd/{change-name}/T{id}-{task-slug}` branches remain in the repository. After 5 interrupted executions, the developer has 20 stale worktrees and branches they do not know exist. `git worktree list` shows them as "prunable" but developers never run this command. Subsequent runs may attempt to create a worktree at a path that already exists and fail.

**Why it happens:**
Go's `defer` cleanup works within a normal process lifecycle, but does not handle `SIGKILL`, `os.Exit()` from nested calls, or CI runner teardowns. Even `os.Signal` handlers don't get to run on hard kills.

**How to avoid:**
At the start of each `execute` command, scan for existing `.worktrees/T*` directories and offer to clean them. Record all created worktrees in a state file (`tasks.md` frontmatter or a dedicated `.mysd-worktrees.yaml`) before creating them — this enables recovery scanning. Use `git worktree prune` as part of startup cleanup. Before creating a worktree at path P, check if P exists and fail with a clear "stale worktree detected, run `mysd cleanup` to remove it." Do not silently overwrite.

**Warning signs:**
- The worktree creation code does not check if the target path already exists
- No `.mysd-worktrees.yaml` or equivalent state tracking exists for in-flight worktrees
- `mysd execute` has no startup orphan scan

**Phase to address:**
Phase 2 (Worktree execution engine) — cleanup and orphan detection must be implemented in the same phase as worktree creation, not deferred.

---

### Pitfall 6: Subagent Context Overload — Passing Full Spec Content Instead of File Paths

**What goes wrong:**
When the orchestrator spawns `mysd-researcher`, `mysd-advisor`, or `mysd-spec-writer`, it is tempting to embed the full content of proposal.md, all spec files, and the entire task list in the subagent's prompt. This inflates the subagent's input token count. A `mysd-advisor x N` spawning pattern with 4 advisors each receiving 8,000 tokens of spec content costs 32,000 input tokens before any work is done. Multiplied across a wave execution, token costs become unpredictable.

**Why it happens:**
Subagents have isolated context windows. The orchestrator's natural instinct is to give the subagent everything it might need. Without deliberate content minimization, this is the default pattern.

**How to avoid:**
Pass file paths, not file content. The orchestrator should pass `specPath`, `proposalPath`, `taskListPath` as strings and instruct the subagent to read them. The subagent reads only what it needs. For `mysd-advisor`, pass only the specific gray area description (one paragraph), not the entire proposal. For `mysd-researcher`, pass the directory to scan, not pre-read files. Establish a rule in SKILL.md preambles: "You will receive file paths. Read only the files relevant to your assigned task."

**Warning signs:**
- Any SKILL.md invocation includes `$(cat proposal.md)` or equivalent content injection
- The orchestrator's Task tool calls have prompts longer than 500 words
- Cost per `propose` run is more than 3x the cost of a single `mysd-planner` run

**Phase to address:**
Phase 1 (Subagent architecture design) — the "paths not content" rule must be established before writing any agent SKILL.md.

---

### Pitfall 7: Plan-Checker False Negatives from Fuzzy MUST Matching

**What goes wrong:**
The plan-checker verifies that all MUST items have corresponding tasks. It does this by asking an AI to match requirement text against task names and descriptions. The AI matches loosely — "The system MUST validate input encoding" is matched against a task called "Build form handler" because the AI infers that form handling implies validation. The plan-checker reports 100% coverage when critical MUST items have no explicit tasks. Users trust the coverage report and ship with gaps.

**Why it happens:**
Fuzzy semantic matching is what LLMs do naturally. Without a strict matching discipline (e.g., requiring task descriptions to explicitly quote the MUST text they satisfy), the checker becomes a rubber-stamp.

**How to avoid:**
The plan-checker must use structured matching: each task in tasks.md should have a `satisfies` field listing MUST item IDs (e.g., `satisfies: [REQ-001, REQ-003]`). The plan-checker verifies that every MUST item ID appears in at least one task's `satisfies` field — this is a deterministic Go check, not an AI inference. The AI's role is limited to suggesting `satisfies` fields when the planner creates tasks; the checker itself uses string matching on IDs, not semantic analysis.

**Warning signs:**
- The plan-checker prompt says "determine if this task covers this requirement" without an explicit ID linkage requirement
- The planner's tasks.md format has no `satisfies` field
- Plan-checker pass rate is consistently 100% on the first run

**Phase to address:**
Phase 2 (Plan-checker implementation) — the `satisfies` field must be added to the tasks.md schema before the plan-checker is built, not retrofitted.

---

### Pitfall 8: Language Detection in Codebase Scout Misclassifies Mixed-Language Projects

**What goes wrong:**
The Codebase Scout scans the project to find integration points and reusable patterns. Language detection based on file extension counts misidentifies mixed-language projects. A Go backend with a React frontend reports as "Go project" if `.go` files outnumber `.tsx` files, causing the scanner to miss the entire frontend. A Python data pipeline with SQL migration scripts is scanned only for Python patterns, missing that the schema files define the real integration contract. The AI executes changes against the wrong layer.

**Why it happens:**
File-extension counting is the simplest detection approach and the first thing developers implement. It is fast and works for single-language projects. Multi-language projects are more common in real-world codebases than in synthetic test fixtures.

**How to avoid:**
Do not report a single primary language — report the top 3 languages by file count AND by line count. Include the directory each language is concentrated in (e.g., `Go: cmd/, internal/ — 7,555 lines; TypeScript: frontend/src/ — 12,000 lines`). When scanning, pass the language distribution and directory map to the research subagent, not just the primary language. Use directory heuristics: `go.mod` implies Go, `package.json` implies JS/TS, `requirements.txt` or `pyproject.toml` implies Python — check for all of these simultaneously rather than stopping at the first match.

**Warning signs:**
- Scanner logic has an early-return after finding the first language indicator
- Test fixtures are all single-language Go projects
- Scanner output has a "primary language" field but no "other detected languages" field

**Phase to address:**
Phase 3 (Codebase Scout + Scan refactor) — multi-language support must be in the initial scanner design, not a follow-up feature.

---

### Pitfall 9: Locale Config Desync Between mysd.yaml and openspec/config.yaml

**What goes wrong:**
The `/mysd:lang` command updates both `mysd.yaml` (response_language, document_language) and `openspec/config.yaml` (locale). If either write fails (permission error, file doesn't exist yet, YAML parsing error in one file), the two configs end up with inconsistent locale values. Subsequent AI runs read from different config files and produce output in different languages mid-session. Documents generated by `mysd-spec-writer` are in Japanese while the orchestrator's progress messages are in English.

**Why it happens:**
Two-file updates are an atomicity problem. The Go code updates file A, then file B. If B fails, A is already updated. There is no rollback. YAML serialization for one format (mysd.yaml) may succeed while the other (openspec/config.yaml with different schema) fails, leaving the system in a split state.

**How to avoid:**
Use a write-both-or-neither pattern: read both files, prepare both new contents in memory, write both to temporary files, then rename (which is atomic on POSIX and Windows NTFS). If either rename fails, attempt to restore from the originals. Log explicitly when both files are updated: "Updated mysd.yaml locale: ja AND openspec/config.yaml locale: ja". On first scan (when openspec/config.yaml doesn't exist), create it as part of `mysd:scan` initialization rather than deferring to `mysd:lang`.

**Warning signs:**
- The `lang` command writes to mysd.yaml first, then openspec/config.yaml, without error recovery between the two writes
- `openspec/config.yaml` may not exist on first run, but the code does not create it proactively
- No test covers the case where the second file write fails after the first succeeds

**Phase to address:**
Phase 3 (lang command + init/scan refactor) — the atomic two-file update pattern must be in the initial implementation.

---

## Technical Debt Patterns

Shortcuts that seem reasonable but create long-term problems.

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Embed full spec content in subagent prompts instead of passing file paths | No need to design file-path contract; works immediately | Token cost quadruples with each additional subagent spawned; orchestrator context fills with redundant content | Never — paths-not-content is cheaper from day one |
| Use AI semantic matching for plan-checker coverage | No tasks.md schema change needed; feels "smarter" | False negatives accepted as true; users trust a broken coverage report; critical MUST items shipped without tasks | Never for MUST-item verification |
| Skip disk space check before worktree creation | Simpler code; most machines have space | Cryptic I/O error mid-execution on developer's full SSD; corrupted worktree state requires manual cleanup | Never — the check is 5 lines of code |
| Open-ended interactive loop with no question count limit | More "natural" conversation feel | Users get trapped in questioning spirals; `--auto` flag cannot reliably skip loops without a defined count | Never — always establish a hard maximum before implementing the loop |
| Single locale field in only one config file | Simpler; avoids sync problem | mysd and openspec tools read different configs; document language diverges; brownfield openspec users see wrong locale | Only acceptable as a transitional state for one phase; must be unified before feature ships |
| Hard-code Go as the scanner's primary language | This is a Go CLI tool written for Go developers | Breaks completely for any user running mysd on a TypeScript, Python, or Rust project | Never — the proposal explicitly requires language-agnostic scanning |
| Detect worktree completion via process exit code only | Simple; works when processes exit cleanly | Process kill / CI timeout leaves no completion marker; orphan detection requires directory scanning instead of state file | Never — always write a state marker before work begins |

---

## Integration Gotchas

Common mistakes when connecting to external services or existing systems.

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| Claude Code subagent spawning | Subagent spawns another subagent (nested delegation) | Only the top-level SKILL.md session can use the Task tool; subagents return results, they do not delegate further |
| Claude Code Task tool prompt | Sending the full conversation history or all spec file contents in the Task prompt | Pass file paths and a concise task description; keep Task prompts under 300 words |
| Git worktree + Windows | Not running `git config core.longpaths true` in the new worktree before any file operations | Detect GOOS at runtime; run longpaths config as the first git operation in each worktree |
| Git worktree + Go build cache | GOCACHE inside the worktree directory multiplies disk usage | Set `GOCACHE` env var to a path outside the worktree (`/tmp/mysd-cache/T{id}`) before running go commands |
| YAML frontmatter in tasks.md | Adding `satisfies` field to tasks without updating the existing `spec.ParseTasksV2` parser | The parser must handle unknown fields gracefully (yaml.v3 does this by default with `omitempty`; verify it does not error on unknown keys) |
| openspec/config.yaml locale field | Writing the BCP 47 locale tag (e.g., `zh-TW`) without validating it against Go's `golang.org/x/text/language` package | Use `language.Parse()` to validate user input before writing; map common aliases ("Chinese Traditional", "繁體中文") to canonical BCP 47 tags |
| Interactive prompts in Go CLI | Using `fmt.Scan` or `bufio.Reader` for prompts without checking `term.IsTerminal(syscall.Stdin)` | Always check `golang.org/x/term.IsTerminal` before presenting interactive prompts; fall back to `--auto` behavior in non-TTY contexts (pipes, CI) |
| Worktree merge sequence | Merging worktrees in arbitrary order when tasks touch the same file | Merge strictly in task ID order; compute the merge sequence before any worktree is created and log it |

---

## Performance Traps

Patterns that work at small scale but fail as usage grows.

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Spawning one `mysd-advisor` per gray area without a maximum cap | 8 advisors spawn for a complex proposal; 8 parallel AI turns consume the context window and token budget simultaneously | Cap advisor parallelism at 4; process remaining areas sequentially | Projects with more than 4 gray areas (common for real proposals) |
| Research mode on every propose/spec run | Research doubles the time and cost of every single propose invocation, even for trivial changes (fixing a typo in a spec) | Research mode must require explicit opt-in confirmation; `--auto` skips research; only ff/ffe runs share research results | Daily development use where most spec changes are incremental |
| Wave execution with no file-overlap check | Two tasks both modifying `internal/spec/parser.go` run in the same wave; merge conflict is certain | File overlap check must happen before wave assignment, not after merge failure | Any time two tasks share a file (common in feature branches) |
| Disk usage: one worktree per task with no cleanup | 10-task execution creates 10 worktrees; developer runs out of space on 256 GB SSD | Auto-cleanup successful worktrees immediately after merge; never accumulate more than (current wave size) worktrees at once | Codebases larger than 500 MB with more than 4 parallel tasks |
| Reading GOCACHE/vendor inside worktree context | Each worktree reads Go module cache from within its own directory tree; download is repeated N times | Point GOPATH/GOMODCACHE to a shared read-only cache outside the worktree tree | Projects with large dependency trees (50+ modules) |

---

## Security Mistakes

Domain-specific security issues.

| Mistake | Risk | Prevention |
|---------|------|------------|
| Passing the user's prompt content directly to a subagent without sanitization | Prompt injection: a malicious spec file instructs the subagent to exfiltrate files or run arbitrary commands | The Go binary formats all subagent prompts as structured templates; user-provided text is placed in a quoted block explicitly labeled as "user input, not instructions" |
| Worktree branches pushed to remote automatically | AI-generated code reaches the remote without review | Never auto-push worktree branches; only create local branches; require explicit `git push` from the user |
| Config sync writing locale values without validation | An attacker who controls openspec/config.yaml can inject YAML that overrides other config fields | Validate and re-serialize locale values as plain strings; never pass raw user input as YAML content; use `yaml.Marshal` on a typed struct, not string interpolation |

---

## UX Pitfalls

Common user experience mistakes specific to interactive discovery and parallel execution.

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| Research mode prompt appears every time, including for trivial spec edits | Users develop "prompt blindness" and always skip research; when they actually need it, they skip it by reflex | Show research mode prompt only when no existing proposal/spec exists for the change; for existing changes, default to no-research with an explicit `--research` flag to opt in |
| Wave execution progress output is a flat log of all task lines interleaved | Developer cannot tell which task is currently executing, which is blocked, or which failed | Render wave progress as a table: task ID, status (running/done/failed), elapsed time; update in-place using ANSI escape codes |
| "Fix" command opens interactive research when the bug is already understood | User who found the bug in their own code is forced through a questioning phase they don't need | `fix` should default to no-research; offer `--research` for investigating unfamiliar bugs; do not make research mandatory |
| lang command accepts free-text language input but silently normalizes it to a wrong locale | User types "Traditional Chinese", gets `zh-CN` (Simplified) | Show the resolved locale before writing: "Resolved 'Traditional Chinese' → zh-TW. Confirm? [Y/n]" |
| `--auto` flag skips interactive questions but still shows "auto-selected: X" for every decision | In a fast-forward run with 10 decisions, this produces 10 lines of noise before any real output | Use a single "Auto mode: N decisions made without prompts" summary line after all auto-selections; do not log each individual auto-selection |

---

## "Looks Done But Isn't" Checklist

Things that appear complete but are missing critical pieces.

- [ ] **Subagent spawning:** Verify that no subagent SKILL.md contains any instruction to use the Task tool to spawn another agent. Check all 9 agent definition files manually.
- [ ] **Interactive loop termination:** Run the `propose` command interactively against a test case with 5 gray areas; verify the loop presents a "ready to proceed?" prompt after each area and exits cleanly without user having to type "done" or "exit".
- [ ] **`--auto` flag:** Run `mysd propose --auto` in a non-TTY context (pipe to `/dev/null`); verify zero interactive prompts appear and the command completes with a written proposal.
- [ ] **Worktree cleanup:** Kill the `mysd execute` process with SIGKILL mid-execution; verify that running `mysd execute` again detects and offers to clean up orphaned worktrees before starting.
- [ ] **Windows path limit:** On Windows (or in CI with a Windows runner), create a worktree for a change in a project at `C:\Users\longusername\projects\my-company-webapp\`; verify no "Filename too long" errors occur.
- [ ] **Disk space guard:** Mock `syscall.Statfs` to return insufficient free space; verify that worktree creation aborts with a human-readable "insufficient disk space" error, not a cryptic I/O error.
- [ ] **Plan-checker MUST coverage:** Create a tasks.md with one MUST item that has no corresponding `satisfies` reference; verify that plan-checker reports a gap rather than 100% coverage.
- [ ] **Language detection:** Run Codebase Scout against the `mysd` repo itself (Go + Markdown) AND against a test fixture with Go + TypeScript files; verify both languages are reported, not just Go.
- [ ] **Locale sync:** After running `mysd lang`, verify that both `mysd.yaml` and `openspec/config.yaml` contain the same locale value; manually corrupt one file mid-write and verify the tool reports an error rather than completing silently.
- [ ] **Non-TTY interactive prompts:** Run any interactive command with stdin redirected from a file (`< /dev/null`); verify the command either uses `--auto` defaults or exits with a clear error message, not a hang waiting for input that never comes.
- [ ] **Merge order:** Run a 3-task wave where tasks T1 and T3 conflict; verify that merges happen in T1→T2→T3 order, not T3→T1→T2, and that the AI conflict resolution runs before the next merge starts.

---

## Recovery Strategies

When pitfalls occur despite prevention, how to recover.

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Nested subagent spawn failure discovered after agent definitions are written | MEDIUM | Restructure orchestration flow: move all Task tool calls to the top-level SKILL.md; update each subagent to remove any Task tool references; retest end-to-end flow |
| Interactive loop trapping users (no exit condition) | LOW | Add hard question count limits to SKILL.md; deploy updated plugin; users can always Ctrl+C and re-run with `--auto` |
| Worktree disk explosion filling developer's SSD | MEDIUM | Run `git worktree prune`; manually delete `.worktrees/` directories; run `git branch -d` for each orphaned branch; add disk guard to next release |
| Plan-checker rubber-stamp (false 100% coverage) | HIGH | Audit all MUST items against tasks manually; add missing tasks; retrofit `satisfies` fields to all existing tasks; rebuild plan-checker with ID-based matching |
| Locale desync between mysd.yaml and openspec/config.yaml | LOW | Run `mysd lang` again to force both files to resync; check for YAML parse errors in both files; fix manually if needed |
| Windows path limit errors in worktrees mid-execution | MEDIUM | Enable longpaths in system registry or git config; remove and recreate affected worktrees; resume execution from the failed task |
| Language detection misclassification causing wrong scaffold | LOW | Rerun scan with `--lang` override flag; update generated specs to reflect correct language targets; delete and regenerate any misclassified spec files |

---

## Pitfall-to-Phase Mapping

How roadmap phases should address these pitfalls.

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| Nested subagent spawning | Phase 1 (Subagent architecture) | Manual audit: no subagent SKILL.md references the Task tool |
| Interactive loop no exit condition | Phase 1 (Interactive discovery design) | Integration test: `--auto` runs without any interactive prompt; manual test: loop presents binary "proceed or continue" choice |
| Worktree disk explosion | Phase 2 (Worktree execution engine) | Unit test: disk space check with mocked filesystem; integration test: verify disk usage after 4-task wave on a known-size codebase |
| Windows MAX_PATH in worktrees | Phase 2 (Worktree execution engine) | CI: Windows runner test with longpaths both enabled and disabled |
| Orphaned worktrees after interrupt | Phase 2 (Worktree execution engine) | Integration test: SIGKILL mid-execution + verify startup orphan detection fires on next run |
| Subagent context overload | Phase 1 (Subagent architecture) | Code review: no SKILL.md invocation embeds file contents; all Task prompts under 300 words |
| Plan-checker false negatives | Phase 2 (Plan-checker) | Unit test: tasks.md with deliberate MUST gap must return gap report; 100% coverage only when all MUST IDs are present in satisfies fields |
| Language detection misclassification | Phase 3 (Codebase Scout) | Integration test: mixed Go+TypeScript fixture reports both languages with directory mapping |
| Locale config desync | Phase 3 (lang command + init/scan refactor) | Unit test: mock second file write failure; verify first file is restored or error is surfaced cleanly |
| Research mode on every run (cost trap) | Phase 1 (Interactive discovery design) | Verify: research prompt appears only for new proposals, not for existing change updates |
| Non-TTY stdin hang | Phase 1 (Interactive discovery design) | Integration test: all interactive commands run with `< /dev/null` without hanging |

---

## Sources

- Claude Code Sub-Agents official docs — https://code.claude.com/docs/en/sub-agents (HIGH confidence, official Anthropic docs)
- Claude Code Subagents: Common Mistakes — https://claudekit.cc/blog/vc-04-subagents-from-basic-to-deep-dive-i-misunderstood (MEDIUM confidence, practitioner blog)
- Context Management with Subagents in Claude Code — https://www.richsnapp.com/article/2025/10-05-context-management-with-subagents-in-claude-code (MEDIUM confidence, practitioner blog)
- GitHub Issue: Sub-Agent Task Tool Not Exposed When Launching Nested Agents — https://github.com/anthropics/claude-code/issues/4182 (HIGH confidence, official issue tracker)
- Git Worktrees for Parallel AI Agents (Upsun) — https://devcenter.upsun.com/posts/git-worktrees-for-parallel-ai-coding-agents/ (MEDIUM confidence, engineering blog with measured data)
- Git worktrees for parallel development (nrmitchi) — https://www.nrmitchi.com/2025/10/using-git-worktrees-for-multi-feature-development-with-ai-agents/ (MEDIUM confidence)
- Windows MAX_PATH Limitation (Microsoft Learn) — https://learn.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation (HIGH confidence, official Microsoft docs)
- Git Windows long path fix — https://www.shadynagy.com/solving-windows-path-length-limitations-in-git/ (MEDIUM confidence)
- Token Cost Trap in AI agents — https://medium.com/@klaushofenbitzer/token-cost-trap-why-your-ai-agents-roi-breaks-at-scale-and-how-to-fix-it-4e4a9f6f5b9a (MEDIUM confidence)
- Taming agent sprawl (CIO) — https://www.cio.com/article/4132287/taming-agent-sprawl-3-pillars-of-ai-orchestration.html (MEDIUM confidence)
- Claude Code agentic loop detection issue — https://github.com/anthropics/claude-code/issues/4277 (HIGH confidence, official issue tracker)
- When agents learn to ask: Active questioning in agentic AI — https://medium.com/@milesk_33/when-agents-learn-to-ask-active-questioning-in-agentic-ai-f9088e249cf7 (LOW confidence, single source)
- Interactive CLI prompts in Go — https://dev.to/tidalcloud/interactive-cli-prompts-in-go-3bj9 (MEDIUM confidence)
- golang.org/x/term IsTerminal — https://pkg.go.dev/golang.org/x/term (HIGH confidence, official Go package docs)
- Go i18n with golang.org/x/text/language — https://phrase.com/blog/posts/internationalization-i18n-go/ (MEDIUM confidence)
- Go static analysis tools 2026 — https://analysis-tools.dev/tag/go (MEDIUM confidence)

---
*Pitfalls research for: my-ssd v1.1 — Interactive Discovery & Parallel Execution milestone*
*Researched: 2026-03-25*
