# mysd

**Content intelligence skills for spec-driven development.**

mysd is a Claude Code plugin that provides 4 independent skills for managing knowledge artifacts in your codebase: researching decisions, syncing documentation, writing specs, and orchestrating all three.

Zero dependencies. Pure SKILL.md files. Install and use immediately.

## Skills

| Command | Skill | What It Does |
|---------|-------|-------------|
| `/mysd:research` | Research | Gray-area decisions with evidence. Spec health checks (4-dimension analysis). |
| `/mysd:doc` | Doc Writer | Detect code changes, identify affected docs, update with style matching and multi-language sync. |
| `/mysd:spec` | Spec Writer | Write/update [OpenSpec](https://github.com/openspec-dev/openspec) format spec files. Reverse-spec from code. |
| `/mysd:sync` | Sync | Run all three skills in sequence via subagents. One command, full content sync. |

Each skill works independently. Use one, two, or all four.

## Install

### Claude Code Plugin

```bash
claude plugin add https://github.com/xenciscbc/mysd
```

### Manual

Copy the skills into your Claude Code skills directory:

```bash
cp -r mysd/skills/research ~/.claude/skills/mysd-research
cp -r mysd/skills/doc ~/.claude/skills/mysd-doc
cp -r mysd/skills/spec ~/.claude/skills/mysd-spec
cp -r mysd/skills/sync ~/.claude/skills/mysd-sync
```

All `/mysd:*` commands will be available in your next Claude Code session.

## Usage

### Research: Make decisions in gray areas

```
/mysd:research Which database should we use for the new caching layer?
```

- Classifies the question (gray area or not)
- Gathers context from codebase, git history, docs, and web
- Frames 2-4 options with evidence, pros/cons, effort
- Produces a Decision Doc with confidence score (1-10)

Also runs **Spec Health Checks** with 4 dimensions (Coverage, Ambiguity, Consistency, Gaps):

```
/mysd:research Check the spec health for the auth-refactor change
```

### Doc: Keep documentation in sync

```
/mysd:doc Update docs after the latest changes
```

- Detects changes via `git diff`
- Maps change types to affected docs (new command -> README, bug fix -> CHANGELOG, etc.)
- Matches the existing style of each doc
- Multi-language sync: updates README.zh-TW.md automatically when README.md changes
- Confirms each change before applying

### Spec: Write OpenSpec format specs

```
/mysd:spec Write a spec for the new auth middleware
```

- Generates correct YAML frontmatter (`spec-version`, `capability`, `delta`, `status`)
- Uses RFC 2119 keywords (MUST, SHOULD, MAY)
- Writes WHEN/THEN/AND scenarios
- Can reverse-spec from code: reads Go files, infers requirements from function signatures

### Sync: Run the full pipeline

```
/mysd:sync I just finished the new caching feature, sync everything
```

Chains research -> doc -> spec via subagents, each with its own context window.

## OpenSpec Compatibility

The spec writer produces files compatible with [OpenSpec](https://github.com/openspec-dev/openspec):

- YAML frontmatter: `spec-version`, `capability`, `delta` (ADDED/MODIFIED/REMOVED/RENAMED), `status`
- RFC 2119 keywords in UPPERCASE
- WHEN/THEN/AND scenario format
- Directory structure: `openspec/specs/{capability}/spec.md`

## Requirements

- Claude Code (any version with skill support)
- No binary, no compilation, no runtime dependencies

## License

MIT
