<!-- SPECTRA:START v1.0.1 -->

# Spectra Instructions

This project uses Spectra for Spec-Driven Development(SDD). Specs live in `openspec/specs/`, change proposals in `openspec/changes/`.

## Use `/spectra:*` skills when:

- A discussion needs structure before coding → `/spectra:discuss`
- User wants to plan, propose, or design a change → `/spectra:propose`
- Tasks are ready to implement → `/spectra:apply`
- There's an in-progress change to continue → `/spectra:ingest`
- User asks about specs or how something works → `/spectra:ask`
- Implementation is done → `/spectra:archive`

## Workflow

discuss? → propose → apply ⇄ ingest → archive

- `discuss` is optional — skip if requirements are clear
- Requirements change mid-work? Plan mode → `ingest` → resume `apply`

## Parked Changes

Changes can be parked（暫存）— temporarily moved out of `openspec/changes/`. Parked changes won't appear in `spectra list` but can be found with `spectra list --parked`. To restore: `spectra unpark <name>`. The `/spectra:apply` and `/spectra:ingest` skills handle parked changes automatically.

<!-- SPECTRA:END -->

<!-- GSD:project-start source:PROJECT.md -->
## Project

**my-ssd**

my-ssd 是一個用 Go 建造的 Claude Code plugin，將 OpenSpec 的 Spec-Driven Development（SDD）方法論與 GSD 級別的規劃/執行/驗證引擎整合為一個無縫系統。它讓獨立開發者（1 人 + N 個 AI agent）能以結構化規格驅動 AI 編程，確保 AI 在寫程式前先對齊需求，並在執行後自動驗證成果。

**Core Value:** **Spec 和執行的緊密整合** — 規格不只是文件，而是直接驅動 AI 執行和驗證的單一事實來源。AI 寫 code 前必須對齊 spec，寫完後自動驗證 spec 是否被滿足。

### Constraints

- **Tech stack**: Go — 單一 binary，跨平台編譯
- **相容性**: 必須能讀寫 OpenSpec 格式的 spec 檔案
- **Plugin 形式**: Claude Code slash commands + agent definitions
- **設計哲學**: Convention over configuration — 預設即好用，只在需要時才配置
<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->
## Technology Stack

## Recommended Stack
### Core Technologies
| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Go | 1.23+ | Primary language | Single-binary deployment, cross-platform compilation via GOOS/GOARCH, no runtime dependency — required by PROJECT.md constraint |
| github.com/spf13/cobra | v1.10.2 | CLI framework — subcommands, flags, help generation | De facto standard for multi-subcommand CLIs (used by kubectl, helm, GitHub CLI). 35k+ GitHub stars. Supports `propose`, `spec`, `design`, `plan`, `execute`, `verify`, `archive` as named subcommands. Pairs with pflag for POSIX-compliant flag handling. |
| gopkg.in/yaml.v3 | v3 (latest) | YAML parsing and serialization | Standard YAML library for Go. Required for OpenSpec frontmatter in `.specs/` files. Cobra v1.10.2 already migrated to go.yaml.in/yaml/v3 internally — use the same module to avoid dependency duplication. |
| github.com/adrg/frontmatter | v0.2.0 | YAML frontmatter extraction from Markdown | Standalone frontmatter parser — no dependency on a specific Markdown engine. Supports `---` delimited YAML, TOML, and JSON frontmatter. Best fit for parsing OpenSpec `proposal.md`, `specs/`, `design.md`, `tasks.md` files where only frontmatter + body split is needed, not full Markdown rendering. |
### Supporting Libraries
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| github.com/spf13/viper | v1.x (stable) | Configuration management | Reading `.mysd.yaml` project config and `~/.mysd/config.yaml` user config. Supports env override, file formats, and defaults. Viper v2 is unreleased — stay on v1. |
| github.com/stretchr/testify | v1.x (latest, Aug 2025) | Test assertions and mocks | All unit tests. Use `assert` for non-fatal checks, `require` for fatal early-exit assertions. The standard assertion library for Go CLI tools. |
| github.com/charmbracelet/lipgloss | v1.x | Terminal output styling | Colored status output for `verify` results, spec alignment summaries, and execution progress. Use only for terminal formatting — not bubbletea TUI (overkill for a CLI tool). |
| text/template (stdlib) | Go stdlib | Template engine for SKILL.md, agent Markdown generation | Generating Claude Code plugin files (`.claude/skills/*/SKILL.md`, `agents/*.md`). No external dependency needed — Go stdlib is sufficient for simple template rendering. Only add `github.com/Masterminds/sprig/v3` if advanced template functions (date formatting, string manipulation) become necessary. |
| github.com/yuin/goldmark | v1.x | Markdown parsing (body content) | Only needed if my-ssd must parse/render the body content of spec files (not just frontmatter). If the tool only reads frontmatter + passes body text through unchanged, skip this. Add it in Phase 2+ when spec content validation is needed. |
### Development Tools
| Tool | Purpose | Notes |
|------|---------|-------|
| goreleaser | v2.14+ | Cross-platform binary release automation | Creates GitHub Releases with Linux/macOS/Windows binaries, Homebrew formula, and checksums. Single `.goreleaser.yaml` config. Industry standard for Go CLI distribution. |
| golangci-lint | v2.x | Meta-linter (50+ linters in one) | v2 released March 2025, replaces v1 config format. Use `linters.default: standard` as baseline. Catches common Go anti-patterns, unused code, error handling issues. Run in CI pre-merge. |
| gopls | latest | Go Language Server | IDE/editor integration. Standard toolchain, no explicit version pinning needed. |
| cobra-cli | latest | Cobra project scaffolding | `cobra-cli add <command>` to generate new subcommand boilerplate. Install as dev tool only — not a library dependency. |
## Claude Code Plugin Integration
### Plugin Structure (what my-ssd generates/ships)
### SKILL.md Format (HIGH confidence — verified from official docs)
### Integration Architecture
### plugin.json Manifest (HIGH confidence)
### Key Constraint: No MCP Server Needed
## Installation
# Core dependencies
# Terminal styling (optional, Phase 2+)
# Markdown body parsing (optional, Phase 2+)
# Test dependencies
# Dev tools (install as binaries, not go.mod dependencies)
## Alternatives Considered
| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| cobra v1.10.2 | urfave/cli v3 | If the tool had a single command or simple flag parsing without deeply nested subcommands. my-ssd has 7+ subcommands with shared state — cobra's nested command model is the better fit. |
| adrg/frontmatter | goldmark-frontmatter | If my-ssd needed to render Markdown to HTML or perform AST traversal on spec bodies. For this tool, only frontmatter + body split is needed — goldmark-frontmatter forces a goldmark dependency unnecessarily. |
| gopkg.in/yaml.v3 | goccy/go-yaml | goccy/go-yaml has better anchor/alias support and faster parsing, but adds complexity. OpenSpec YAML files are simple key-value frontmatter — yaml.v3 is sufficient and already in the dependency tree via cobra. |
| text/template (stdlib) | Masterminds/sprig | Add sprig only if SKILL.md templates need advanced string manipulation (regex, date formatting). Start with stdlib — it handles Go struct rendering without extra dependencies. |
| spf13/viper v1 | koanf | koanf has a cleaner API and smaller dependency tree, but viper v1 is already a known quantity in the cobra ecosystem and most Go CLI tutorials use it. Use viper for familiarity unless dependency minimalism is critical. |
| golangci-lint v2 | staticcheck only | staticcheck alone is simpler to configure but misses many categories (unused variables, error wrapping, etc.). golangci-lint v2's `linters.default: standard` profile is low-friction and catches more issues. |
## What NOT to Use
| Avoid | Why | Use Instead |
|-------|-----|-------------|
| bubbletea (charmbracelet/bubbletea) | Full TUI framework with Elm architecture — massive overkill for a CLI that outputs structured text. Adds significant complexity to a tool that should be a simple command runner. | lipgloss for styled output only, stdlib fmt for structured text |
| MCP server pattern for plugin integration | MCP servers require an always-running process and add infrastructure complexity. The binary-called-from-skills pattern achieves the same result with zero additional infrastructure. | Bash invocation from SKILL.md files |
| encoding/json for spec files | OpenSpec format uses YAML frontmatter, not JSON. Using JSON would break OpenSpec compatibility. | gopkg.in/yaml.v3 |
| Viper v2 | Not yet released (as of 2026-03-23). The v2 API is discussed but unreleased — depending on pre-release code adds risk. | spf13/viper v1 (stable API) |
| cobra-cli in go.mod | cobra-cli is a scaffolding tool, not a runtime library. It generates boilerplate once and should not appear in go.mod. | Install globally: `go install github.com/spf13/cobra-cli@latest` |
| go-openapi/testify/v2 | Active fork, not the canonical stretchr/testify. The fork's API stability claims are unproven in production use. | stretchr/testify v1 (battle-tested, widely documented) |
## Stack Patterns by Variant
- adrg/frontmatter supports TOML (`+++` delimiters) out of the box — no library change needed
- yaml.v3 can be dropped if ALL spec files migrate to TOML
- Use Go's native `goroutines` + `sync.WaitGroup` for concurrent spec execution
- No additional concurrency library needed — stdlib is sufficient for this use case
- Implement as a custom Go parser over the spec body string — no external library needed
- Simple regex or string scanning against `regexp.MustCompile` is sufficient
- The binary must be committed to the plugin directory or downloaded via a `SessionStart` hook
- Use GoReleaser to build platform-specific binaries; the plugin's SessionStart hook can detect OS and download the correct binary from GitHub Releases
## Version Compatibility
| Package | Compatible With | Notes |
|---------|-----------------|-------|
| cobra v1.10.2 | pflag v1.0.9+ | cobra v1.10.0 introduced pflag v1.0.9 as required — earlier pflag versions break compilation |
| cobra v1.10.2 | go.yaml.in/yaml/v3 | cobra migrated from gopkg.in/yaml.v3 in v1.10.2 — both resolve to the same underlying implementation |
| golangci-lint v2 | Go 1.22+ | golangci-lint v2 requires Go 1.22 minimum; v1 supports older versions |
| goreleaser v2 | Go 1.22+ | goreleaser v2 builds require Go 1.22+; earlier Go versions need goreleaser v1 |
| lipgloss v1 | Go 1.21+ | lipgloss v1 uses range over integer (Go 1.22 feature) in some internal code — verify module's go directive |
## Sources
- [github.com/spf13/cobra releases](https://github.com/spf13/cobra/releases) — v1.10.2 confirmed latest (Dec 2024), HIGH confidence
- [code.claude.com/docs/en/slash-commands](https://code.claude.com/docs/en/slash-commands) — SKILL.md format, frontmatter fields, plugin structure, HIGH confidence (official Anthropic docs)
- [code.claude.com/docs/en/plugins-reference](https://code.claude.com/docs/en/plugins-reference) — Complete plugin directory structure, agent format, hooks, HIGH confidence (official Anthropic docs)
- [github.com/adrg/frontmatter](https://github.com/adrg/frontmatter) — v0.2.0 latest (Dec 2025), MEDIUM confidence (search result, not directly verified on releases page)
- [goreleaser.com](https://goreleaser.com/) — v2.14 latest confirmed, HIGH confidence
- [golangci-lint.run](https://golangci-lint.run/) — v2 released March 2025, HIGH confidence
- [github.com/stretchr/testify](https://github.com/stretchr/testify) — v1 stable, last updated Aug 2025, HIGH confidence
- WebSearch: Go CLI framework comparison 2025 — MEDIUM confidence (multiple sources agree on cobra for complex CLIs)
- WebSearch: viper v2 status — MEDIUM confidence (unreleased as of research date, multiple sources confirm)
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

### Release Versioning

Every release must keep these version sources in sync:
- **Git tag** (e.g. `v1.1.0`)
- `.claude-plugin/plugin.json` → `version`
- `.claude-plugin/marketplace.json` → `metadata.version` and `plugins[0].version`

All three locations must use the same version number (tag without the `v` prefix).
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd:quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd:debug` for investigation and bug fixing
- `/gsd:execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->

## Documentation Sync

When mysd project features, commands, configuration, or workflow change, **always update README.md** to reflect the changes. This includes:
- Adding/removing/renaming commands
- Changing workflow steps or flow
- Modifying model profile assignments
- Updating configuration options

<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd:profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
