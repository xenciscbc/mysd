## ADDED Requirements

### Requirement: Analyze command performs cross-artifact structural analysis

The `mysd analyze [change-name]` command SHALL perform deterministic structural analysis across all artifacts in the change directory. If `change-name` is omitted, the command SHALL use the active change from workflow state.

The command SHALL analyze 4 dimensions:

1. **Coverage**: Every capability listed in proposal's Capabilities section SHALL have a corresponding `specs/<name>/spec.md`
2. **Consistency**: File paths and component names SHALL be consistent across proposal, specs, design, and tasks. Design SHALL reference only capabilities from the proposal. Tasks SHALL cover all design decisions.
3. **Ambiguity**: Spec files SHALL NOT contain weak language patterns (`should`, `may`, `might`, `TBD`, `TODO`, `FIXME`, `TKTK`, `???`) where normative language (`SHALL`, `MUST`) is expected
4. **Gaps**: Every requirement in spec files SHALL have at least one scenario. Tasks SHALL reference spec requirements.

#### Scenario: Coverage detects missing spec

- **WHEN** proposal lists capability `foo-bar` in Capabilities section
- **AND** no `specs/foo-bar/spec.md` file exists in the change directory
- **THEN** the command SHALL report a Critical finding with dimension "Coverage"

#### Scenario: Ambiguity detects weak language in spec

- **WHEN** a spec file contains the word `should` outside of a quoted string
- **THEN** the command SHALL report a Suggestion finding with dimension "Ambiguity"

#### Scenario: Gaps detects requirement without scenario

- **WHEN** a spec file contains `### Requirement:` without a following `#### Scenario:`
- **THEN** the command SHALL report a Warning finding with dimension "Gaps"

### Requirement: Analyze command outputs structured JSON

The `mysd analyze` command SHALL support a `--json` flag that outputs analysis results in JSON format.

The JSON output SHALL contain:
- `change_id`: the change name
- `dimensions`: array of `{dimension, status, finding_count}` for each of the 4 dimensions
- `findings`: array of `{id, dimension, severity, location, summary, recommendation}` for each finding
- `artifacts_analyzed`: array of artifact types that were found and analyzed
- `artifacts_missing`: array of artifact types that were not found

Severity levels SHALL be: `Critical`, `Warning`, `Suggestion`.

Finding IDs SHALL use the format: `COV-N` (Coverage), `CON-N` (Consistency), `AMB-N` (Ambiguity), `GAP-N` (Gaps).

#### Scenario: JSON output with no findings

- **WHEN** `mysd analyze my-change --json` finds no issues
- **THEN** the JSON output SHALL contain an empty `findings` array
- **AND** all dimensions SHALL have `status: "Clean"` and `finding_count: 0`

#### Scenario: JSON output with findings

- **WHEN** `mysd analyze my-change --json` finds 2 Coverage issues
- **THEN** the `findings` array SHALL contain 2 entries with `dimension: "Coverage"`
- **AND** the Coverage dimension SHALL have `finding_count: 2`

### Requirement: Analyze command provides styled terminal output

When the `--json` flag is NOT provided, `mysd analyze` SHALL output lipgloss-styled terminal output showing each dimension's status and any findings with severity coloring.

#### Scenario: Styled output shows dimension summary

- **WHEN** `mysd analyze my-change` is run without `--json`
- **THEN** the output SHALL display each dimension name, status, and finding count
- **AND** Critical findings SHALL be displayed in red, Warning in yellow, Suggestion in gray

### Requirement: Analyze operates on available artifacts only

The `mysd analyze` command SHALL analyze whatever artifacts exist in the change directory without requiring all artifacts to be present.

- If only `proposal.md` and `specs/` exist (propose phase): analyze Coverage and Ambiguity
- If `design.md` and `tasks.md` also exist (plan phase): additionally analyze Consistency and Gaps

Artifacts that do not exist SHALL be listed in `artifacts_missing` but SHALL NOT cause the command to fail.

#### Scenario: Propose phase analysis with missing design

- **WHEN** a change has `proposal.md` and `specs/` but no `design.md` or `tasks.md`
- **THEN** analyze SHALL check Coverage and Ambiguity
- **AND** Consistency checks requiring design/tasks SHALL be skipped
- **AND** `artifacts_missing` SHALL contain `"design"` and `"tasks"`
