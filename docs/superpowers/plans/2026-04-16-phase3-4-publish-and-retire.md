# Phase 3+4: Publish mysd-skills Plugin & Retire Go Binary

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Publish mysd-skills as a standalone Claude Code plugin, then deprecate the Go binary and old plugin in the original repo.

**Architecture:** The new mysd-skills plugin lives in `mysd-skills/` subdirectory for now (same repo). Phase 4 adds deprecation notices to the old binary and plugin, pointing users to the new skills. The old binary stays installable but is no longer maintained.

**Tech Stack:** Claude Code plugin format (plugin.json + SKILL.md), Git, Markdown

---

## File Structure

```
Changes to existing files:
  .claude-plugin/plugin.json          # Add deprecation notice
  .claude-plugin/marketplace.json     # Add deprecation notice
  README.md                           # Add deprecation banner + migration guide
  README.zh-TW.md                     # Same in Traditional Chinese
  CLAUDE.md                           # Remove binary command references

Changes to mysd-skills/:
  mysd-skills/plugin.json             # Complete with author/repository/keywords
  mysd-skills/README.md               # Already complete
```

---

### Task 1: Finalize mysd-skills plugin.json

**Files:**
- Modify: `mysd-skills/plugin.json`

- [ ] **Step 1: Read the current plugin.json**

Read `mysd-skills/plugin.json`. Current content:
```json
{
  "name": "mysd-skills",
  "version": "1.0.0",
  "description": "Content intelligence skills for spec-driven development — research decisions, sync docs, write specs"
}
```

- [ ] **Step 2: Update with complete fields**

Write `mysd-skills/plugin.json`:

```json
{
  "name": "mysd-skills",
  "version": "1.0.0",
  "description": "Content intelligence skills for spec-driven development — research decisions, sync docs, write specs",
  "author": {
    "name": "xenciscbc"
  },
  "repository": "https://github.com/xenciscbc/mysd",
  "license": "MIT",
  "keywords": [
    "content-intelligence",
    "spec-driven",
    "openspec",
    "documentation",
    "research",
    "decision-making"
  ]
}
```

- [ ] **Step 3: Commit**

```bash
git add mysd-skills/plugin.json
git commit -m "chore(mysd-skills): finalize plugin.json with author and metadata"
```

---

### Task 2: Add deprecation notice to old plugin manifests

**Files:**
- Modify: `.claude-plugin/plugin.json`
- Modify: `.claude-plugin/marketplace.json`

- [ ] **Step 1: Update .claude-plugin/plugin.json**

Read `.claude-plugin/plugin.json`, then update the description to include deprecation:

```json
{
  "name": "mysd",
  "version": "1.0.5",
  "description": "[DEPRECATED] Use mysd-skills instead. Original Spec-Driven Development Go CLI — replaced by pure SKILL.md plugin.",
  "author": {
    "name": "my-ssd contributors"
  },
  "repository": "https://github.com/xenciscbc/mysd"
}
```

Note: version bumped to 1.0.5 (deprecation release).

- [ ] **Step 2: Update .claude-plugin/marketplace.json**

Read `.claude-plugin/marketplace.json`, then update:

```json
{
  "name": "mysd",
  "owner": {
    "name": "xenciscbc"
  },
  "metadata": {
    "description": "[DEPRECATED] Use mysd-skills instead. Original Spec-Driven Development Go CLI — replaced by pure SKILL.md plugin.",
    "version": "1.0.5"
  },
  "plugins": [
    {
      "name": "mysd",
      "source": "./mysd",
      "description": "[DEPRECATED] Use mysd-skills instead — pure SKILL.md plugin with research, doc, and spec skills",
      "version": "1.0.5",
      "author": {
        "name": "xenciscbc"
      },
      "repository": "https://github.com/xenciscbc/mysd",
      "license": "MIT",
      "keywords": [
        "deprecated",
        "sdd",
        "spec-driven",
        "openspec"
      ],
      "category": "development"
    }
  ]
}
```

- [ ] **Step 3: Commit**

```bash
git add .claude-plugin/
git commit -m "chore: mark old mysd plugin as deprecated (v1.0.5)"
```

---

### Task 3: Add deprecation banner to README.md

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Read the current README.md header**

Read `README.md` lines 1-40 to see the current structure.

- [ ] **Step 2: Add deprecation banner after the title**

Replace the existing testing notice with a deprecation banner. Find this text:

```markdown
> **Testing / 測試中** — This project is under active development. APIs and workflows may change.
```

Replace with:

```markdown
> **⚠️ DEPRECATED** — This Go CLI binary is no longer maintained. Use **[mysd-skills](mysd-skills/)** instead — a pure SKILL.md plugin with zero dependencies.
>
> The new plugin provides 3 content intelligence skills: `/mysd:research` (gray-area decisions), `/mysd:doc` (doc sync), `/mysd:spec` (spec writing), plus `/mysd:run` (orchestrator).
>
> **Migration:** Copy `mysd-skills/` to your Claude Code skills directory. No binary needed.
```

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: add deprecation banner to README.md"
```

---

### Task 4: Add deprecation banner to README.zh-TW.md

**Files:**
- Modify: `README.zh-TW.md`

- [ ] **Step 1: Read the current README.zh-TW.md header**

Read `README.zh-TW.md` lines 1-40.

- [ ] **Step 2: Add deprecation banner**

Find:

```markdown
> **測試中** — 本專案正在積極開發中，API 和工作流程可能會變動。
```

Replace with:

```markdown
> **⚠️ 已棄用** — 此 Go CLI binary 不再維護。請改用 **[mysd-skills](mysd-skills/)** — 純 SKILL.md plugin，零依賴。
>
> 新的 plugin 提供 3 個內容智慧 skill：`/mysd:research`（灰區決策）、`/mysd:doc`（文件同步）、`/mysd:spec`（規格撰寫），以及 `/mysd:run`（協調器）。
>
> **遷移方式：** 將 `mysd-skills/` 複製到你的 Claude Code skills 目錄即可。不需要 binary。
```

- [ ] **Step 3: Commit**

```bash
git add README.zh-TW.md
git commit -m "docs: add deprecation banner to README.zh-TW.md"
```

---

### Task 5: Update CLAUDE.md — remove binary-dependent instructions

**Files:**
- Modify: `CLAUDE.md`

This is the most sensitive edit. CLAUDE.md controls how Claude Code interacts with this project. We need to:
1. Remove references to Go binary commands
2. Remove the GSD workflow enforcement section (no longer applicable)
3. Update the project description
4. Keep OpenSpec/Spectra references (those are still valid)
5. Keep build/docs conventions that still apply

- [ ] **Step 1: Read the full CLAUDE.md**

Read `CLAUDE.md` to understand all sections.

- [ ] **Step 2: Update the project description**

Find the `## Project` section. Update the description to reflect the new pure SKILL.md architecture. Replace references to "Go CLI binary" with "Claude Code SKILL.md plugin". Keep the core value statement about spec-driven development.

- [ ] **Step 3: Remove GSD Workflow Enforcement section**

Find and remove this entire section:

```markdown
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd:quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd:debug` for investigation and bug fixing
- `/gsd:execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
```

This section enforced the old binary-based workflow. The new skills don't need it.

- [ ] **Step 4: Remove Build After Changes section**

Find and remove:

```markdown
## Build After Changes

When Go source code is modified (cmd/, internal/, etc.), **always rebuild the binary** before testing or committing:
```bash
go build -o mysd.exe .
```
This ensures the installed binary matches the source code. Skipping this step causes `mysd --help` and runtime behavior to be out of sync with the code.
```

The Go binary is deprecated — no need to rebuild.

- [ ] **Step 5: Update Documentation Sync section**

Find the `## Documentation Sync` section. Update to mention the new `/mysd:doc` skill instead of manual sync:

```markdown
## Documentation Sync

When project features, commands, configuration, or workflow change, use `/mysd:doc` to automatically detect and update affected documentation files. This handles both README.md and README.zh-TW.md.
```

- [ ] **Step 6: Remove Technology Stack section if it references Go-specific tools**

The Technology Stack section describes Go dependencies (cobra, viper, lipgloss, etc.). Since the binary is deprecated, this section is no longer relevant for new development. Add a note:

```markdown
> **Note:** The Technology Stack section below describes the deprecated Go binary. New development uses pure SKILL.md files with no compiled dependencies. See `mysd-skills/README.md` for the current architecture.
```

- [ ] **Step 7: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: update CLAUDE.md for mysd-skills architecture, remove binary references"
```

---

### Task 6: Final verification and tag

**Files:**
- None modified (verification only)

- [ ] **Step 1: Verify all deprecation notices are in place**

```bash
grep -l "DEPRECATED\|已棄用\|deprecated" README.md README.zh-TW.md .claude-plugin/plugin.json .claude-plugin/marketplace.json
```

Expected: all 4 files listed.

- [ ] **Step 2: Verify mysd-skills is complete**

```bash
ls mysd-skills/*/SKILL.md mysd-skills/research/formats/*.md mysd-skills/plugin.json mysd-skills/README.md
```

Expected: 8 files (4 SKILL.md + 2 format files + plugin.json + README.md).

- [ ] **Step 3: Verify CLAUDE.md no longer references `go build` or GSD workflow**

```bash
grep -c "go build\|gsd:quick\|gsd:debug\|gsd:execute" CLAUDE.md
```

Expected: 0 matches.

- [ ] **Step 4: Check version alignment**

Verify versions are consistent:
- `.claude-plugin/plugin.json` → `1.0.5` (deprecation release)
- `.claude-plugin/marketplace.json` → `1.0.5`
- `mysd-skills/plugin.json` → `1.0.0`

```bash
grep '"version"' .claude-plugin/plugin.json .claude-plugin/marketplace.json mysd-skills/plugin.json
```

- [ ] **Step 5: Commit verification results**

If any verification fails, fix before proceeding.

```bash
git status
```

If clean: no action needed. If dirty: commit remaining changes.

---

## Self-Review Checklist

1. **Spec coverage:** Design doc Phase 3 (publish) → Tasks 1. Phase 4 (retire) → Tasks 2-5. Verification → Task 6. All covered. ✓
2. **Placeholder scan:** All edits show exact content. No "TBD" or "fill in later". ✓
3. **Type consistency:** Version numbers consistent: old plugin 1.0.5, new plugin 1.0.0. `mysd-skills` naming used throughout. ✓
4. **Both READMEs updated:** EN and zh-TW both get deprecation banners (design doc requirement). ✓
5. **CLAUDE.md references cleaned:** GSD enforcement, build command, and tech stack all addressed. ✓
