---
spec-version: "1.0"
capability: Project Setup & Configuration
delta: ADDED
status: done
---

## Requirement: Project Initialization

The `mysd init` command MUST scaffold the following directory structure:
- `openspec/config.yaml` with default schema and locale
- `.specs/` directory for active changes
- `.claude/` directory with plugin hooks and settings

The init command MUST be idempotent — re-running it MUST NOT overwrite existing files.

The init command MUST install the statusline hook to `.claude/hooks/mysd-statusline.js`.

The init command MUST merge statusline configuration into `.claude/settings.json` without overwriting existing keys.

## Requirement: Configuration Management

The system MUST read project configuration from `.claude/mysd.yaml` via Viper.

`ProjectConfig` MUST support the following fields:
- `model_profile`: AI model profile selection (quality/balanced/budget)
- `model_overrides`: Per-agent-role model overrides
- `statusline_enabled`: Toggle for statusline display (pointer-to-bool, defaults to true)

The `mysd model` command MUST display the current model profile and MUST support setting a new profile.

The `mysd lang` command MUST set both `response_language` and `document_language` in config.

The `mysd docs` command MUST manage the `docs_to_update` list for post-archive documentation updates.

## Requirement: Codebase Scanning

The `mysd scan` command MUST analyze the project directory and output a `ScanContext` JSON containing:
- `root_dir`, `primary_language`, `files` (extension counts), `modules` (discovered packages)
- `existing_specs` (specs already present), `total_files`, `config_exists`

The scanner MUST auto-detect primary language from marker files (`go.mod`, `package.json`, `pyproject.toml`).

The `--context-only` flag MUST output JSON without generating any files.

The `--exclude` flag MUST accept comma-separated directory names to skip.

## Requirement: Terminal Output

The `output.Printer` MUST auto-detect TTY and apply colored output (via lipgloss) only when connected to a terminal.

Non-TTY output MUST use plain text with `[PREFIX]` markers (e.g., `[OK]`, `[ERR]`, `[WARN]`).

## Requirement: Statusline Hook

The statusline hook (`mysd-statusline.js`) MUST read JSON from stdin and produce a formatted statusline showing:
- AI model shortname
- Active change name (from `.specs/state.yaml`)
- Context usage visualization with color-coded progress bar

The `mysd statusline` command MUST toggle `statusline_enabled` in `.claude/mysd.yaml`.

### Scenario: First-time Setup

WHEN a user runs `mysd init` in a fresh project
THEN the openspec directory structure is created
AND `.claude/hooks/mysd-statusline.js` is installed
AND `.claude/settings.json` is updated with statusline configuration

### Scenario: Re-initialization

WHEN a user runs `mysd init` in an already-initialized project
THEN no existing files are overwritten
AND any missing files are created

## Covered Packages

- `cmd/init.go`, `cmd/scan.go`, `cmd/model.go`, `cmd/lang.go`, `cmd/docs.go`, `cmd/statusline.go`
- `internal/config/` — configuration loading and model resolution
- `internal/output/` — TTY-aware terminal styling
- `internal/scanner/` — codebase analysis and language detection
- `plugin/hooks/mysd-statusline.js` — statusline Node.js hook
