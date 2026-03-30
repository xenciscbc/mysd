## ADDED Requirements

### Requirement: Propose workflow step ordering

The `mysd:propose` skill SHALL execute steps in the following order:

1. **Parse Arguments**: Parse `--auto` flag and `source_arg`
2. **Resolve Agent Model**: Run `mysd model` to determine per-role models
3. **Load Deferred Notes**: Run `mysd note list` to load cross-change context
4. **Material Selection**: Detect all available requirement sources, present to user for selection, aggregate selected content (see `material-selection` spec)
5. **Scan Existing Specs**: Scan `openspec/specs/*/spec.md` for related specs, retain content for interview step
6. **Requirement Interview**: Evaluate aggregated content completeness, ask clarifying questions as needed, produce structured `requirement_brief` (see `requirement-interview` spec)
7. **Derive Change Name + Classify Type**: Derive kebab-case change name from `requirement_brief`, classify as feature/bugfix/refactor
8. **Scaffold the Change**: Run `mysd propose {change-name}` to create directory structure
9. **Optional 4-Dimension Research**: Ask user whether to run research; if accepted, spawn 4 `mysd-researcher` agents in parallel
10. **Gray Area + Advisor** (research only): Identify gray areas from research output, spawn `mysd-advisor` per area
11. **Dual-Loop Exploration** (research only): Deep dive into gray areas with scope guardrail
12. **Invoke Proposal Writer**: Spawn `mysd-proposal-writer` with `requirement_brief` and research/exploration conclusions
13. **Auto-Invoke Spec Writer**: Spawn `mysd-spec-writer` per capability area
14. **Artifact Review**: Run `mysd validate` then spawn `mysd-reviewer`
15. **Final Summary**: Display results and next steps

Steps 10-11 SHALL be executed only when the user accepts 4-Dimension Research in Step 9. When research is declined, the workflow SHALL proceed directly from Step 9 to Step 12.

#### Scenario: Full workflow with research

- **WHEN** the user accepts 4-Dimension Research
- **THEN** the skill SHALL execute all 15 steps in order

#### Scenario: Workflow without research

- **WHEN** the user declines 4-Dimension Research
- **THEN** the skill SHALL skip Steps 10-11 and proceed from Step 9 directly to Step 12

### Requirement: Change name derived after requirement interview

The `mysd:propose` skill SHALL derive the change name from the completed `requirement_brief` (Step 7), not from the initial `source_arg` or raw source content.

If `source_arg` refers to an existing change directory (`.specs/changes/{source_arg}/`), the skill SHALL use that name directly without re-derivation.

#### Scenario: Change name from requirement_brief

- **WHEN** the user provides a vague description "improve the propose flow"
- **AND** the interview clarifies the scope to material selection and requirement interview
- **THEN** the derived change name SHALL reflect the clarified scope (e.g., `enhance-propose-workflow`)

#### Scenario: Existing change name preserved

- **WHEN** `source_arg` matches an existing change directory
- **THEN** the skill SHALL use the existing change name without re-derivation

### Requirement: Existing spec content fed into interview

The `mysd:propose` skill SHALL pass related existing spec content (from Step 5 scan) into the requirement interview step (Step 6).

The orchestrator SHALL use this spec content to detect overlap between the proposed change and existing capabilities, and SHALL ask the user whether to extend the existing spec or create a new capability when overlap is detected.

#### Scenario: Overlap detected with existing spec

- **WHEN** the proposed change overlaps with the existing `planning` spec
- **THEN** the orchestrator SHALL ask the user: extend the existing `planning` spec or create a new capability

#### Scenario: No overlap with existing specs

- **WHEN** no existing specs are related to the proposed change
- **THEN** the orchestrator SHALL proceed with the interview without spec-related questions
