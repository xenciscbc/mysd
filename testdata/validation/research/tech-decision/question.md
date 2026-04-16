# Research Question: Plugin README Structure

## Question

Should the mysd-skills plugin use a single README.md or per-skill README files?

## Context

The mysd-skills plugin currently contains 4 skills:
- `research` — structured research and analysis
- `doc` — documentation generation and validation
- `spec` — reverse spec generation and health checks
- `plan` — planning and task breakdown

Each skill has its own `SKILL.md` file that Claude Code loads. The question is whether to maintain:

**Option A: Single top-level README.md**
- One file documents all 4 skills
- Easier for users discovering the plugin for the first time
- Risk: becomes very long as skills are added
- Updates require editing one large file

**Option B: Per-skill README files**
- Each skill directory has its own `README.md`
- Example: `skills/research/README.md`, `skills/doc/README.md`
- More modular — each skill is independently understandable
- Independent installability: users can copy a single skill directory and get all docs with it
- Risk: harder to get an overview of all skills at once

## Independent Installability Requirement

The plugin design goal states that individual skills should be independently installable. A user should be able to copy just the `skills/research/` directory into their own plugin and have a fully functional, documented skill.

This requirement weights toward Option B (per-skill README files).

## Additional Considerations

- The plugin is distributed as a single unit today, but may be split in the future
- Marketplace listings typically show one README per plugin
- Developer experience: contributors only need to update one file when modifying a skill
- 4 skills today, potentially 10+ in the future

## Desired Output

A recommendation with rationale, covering:
1. Which option to choose
2. Whether a top-level README is still needed alongside per-skill READMEs
3. Migration path if the current structure needs to change
