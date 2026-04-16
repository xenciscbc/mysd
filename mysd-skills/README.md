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

The spec writer produces files compatible with the OpenSpec format:
- YAML frontmatter with `spec-version`, `capability`, `delta`, `status`
- RFC 2119 keywords (MUST/SHOULD/MAY)
- WHEN/THEN/AND scenario format
- Directory structure: `openspec/specs/{capability}/spec.md`
