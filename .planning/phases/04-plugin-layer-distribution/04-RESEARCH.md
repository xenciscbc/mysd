# Phase 4: Plugin Layer & Distribution - Research

**Researched:** 2026-03-24
**Domain:** Go binary distribution (GoReleaser), Claude Code plugin packaging, scan command (context-only pattern), roadmap tracking (YAML + Mermaid)
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Scan Command (WCMD-09)**
- D-01: Scan uses the context-only pattern (same as execute/verify) — binary produces codebase metadata JSON, AI agent analyzes and generates spec files
- D-02: Before scanning, binary lists codebase structure and presents to user for confirmation — user selects which directories/files to exclude before AI analysis begins
- D-03: Scan skips existing specs — if `.specs/changes/{name}/` already exists, that module is skipped (no overwrite). No `--force` flag in v1.
- D-04: Scan granularity is per-module/package — each Go package or major module produces one change/spec. Aligns with OpenSpec "one spec per capability" philosophy.

**Distribution (DIST-03)**
- D-05: Primary installation method is `go install github.com/owner/mysd@latest`
- D-06: GoReleaser configured for standard 3-platform matrix: Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64). Produces GitHub Releases with precompiled binaries and checksums.
- D-07: No Homebrew tap/cask in v1 — `go install` and direct binary download are sufficient for initial release.

**Roadmap Tracking (RMAP-01~03)**
- D-08: Tracking data stored in `.mysd/roadmap/tracking.yaml` — single YAML file recording all changes with name, status, dates, task counts, and verification status
- D-09: Mermaid gantt chart generated as separate `.mysd/roadmap/timeline.md` — auto-regenerated whenever tracking.yaml is updated. Keeps data (YAML) and visualization (Mermaid) cleanly separated.
- D-10: Tracking updates triggered on state transitions — integrated with existing SaveState flow. Every propose→spec→design→plan→execute→verify→archive transition updates the tracking file.
- D-11: Tracked fields: change name, current status, start/completion dates, task total/completed count, verification MUST pass/fail statistics

**Plugin Packaging (DIST-04)**
- D-12: Plugin installed by copying plugin directory to `.claude/plugins/mysd/`. Standard Claude Code plugin directory structure.
- D-13: SessionStart hook checks binary existence and version — if `mysd` not in PATH or version below minimum, displays installation instructions (not auto-download). Non-blocking — session continues even if binary missing.
- D-14: plugin.json follows standard Claude Code format: name, version, description, commands[], agents[], hooks[]. No extended metadata in v1.
- D-15: Plugin upgrade via SessionStart version check — compares `mysd --version` against plugin.json min_version. Displays upgrade instructions when outdated. No auto-update.

### Claude's Discretion
- GoReleaser `.goreleaser.yaml` specific configuration details (archive format, naming convention, checksum algorithm)
- Scan agent's exact prompt wording and analysis depth
- tracking.yaml schema field ordering and naming conventions
- timeline.md Mermaid chart styling and section grouping
- plugin/ directory internal file organization

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| DIST-03 | Install via `go install` and GitHub releases (precompiled binaries) | GoReleaser v2 config patterns, version injection via ldflags, go module path requirements |
| DIST-04 | Claude Code plugin integration via slash commands and agent definitions | Plugin manifest schema (HIGH confidence — verified from official docs), hooks.json SessionStart format, commands/ and agents/ directory structure |
| WCMD-09 | `/mysd:scan` — scan existing project codebase and generate OpenSpec-format spec documents | Context-only pattern from verify.go + execute.go, new internal/scanner package, SKILL.md orchestrator pattern from existing commands |
| RMAP-01 | Auto-generate/update `.mysd/roadmap/` tracking files after state transitions | SaveState hook point in state.go, gopkg.in/yaml.v3 already in go.mod, Mermaid gantt chart format |
| RMAP-02 | Track change name, status, start/completion dates | tracking.yaml schema design, WorkflowState fields already capture phase + timestamp |
| RMAP-03 | Tracking format readable by third-party tools (supports Mermaid gantt chart) | Mermaid gantt syntax verified, YAML as machine-readable format |
</phase_requirements>

---

## Summary

Phase 4 is the final packaging and distribution phase. All core business logic is complete in Phases 1-3. This phase wraps the existing binary into a distributable package and adds two new capabilities: the `scan` command (WCMD-09) and roadmap tracking (RMAP-01~03).

The scan command follows the established context-only pattern already used in `execute` and `verify`. The binary outputs codebase structure metadata as JSON; the AI agent (invoked via Task tool from SKILL.md) analyzes the code and generates spec files. No new architectural patterns are needed — this is a direct replication of the verify pattern.

Distribution via GoReleaser v2 is straightforward: the go.mod module path (`github.com/mysd`) needs to match the actual GitHub repository for `go install` to work. The `.goreleaser.yaml` config produces cross-platform binaries; Windows arm64 should be excluded (not commonly requested and adds complexity).

The Claude Code plugin structure has been **significantly updated** from the old CLAUDE.md research. The current plugin system (verified from official Anthropic docs, 2026-03) uses `.claude-plugin/plugin.json` as the manifest location, with `commands/`, `agents/`, `skills/`, and `hooks/hooks.json` at the plugin root level — NOT inside `.claude-plugin/`. The current project already ships commands in `.claude/commands/` and agents in `.claude/agents/` (standalone mode). The plugin package will mirror this structure.

**Primary recommendation:** Follow the existing context-only pattern for scan. For the plugin directory, use the verified current structure: `.claude-plugin/plugin.json` manifest + `commands/` + `agents/` + `hooks/hooks.json` at plugin root. Hooks use `hooks.json` format, not inline in `plugin.json`.

---

## Standard Stack

### Core (all already in go.mod — no new dependencies needed)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| github.com/spf13/cobra | v1.10.2 | CLI — new `scan` subcommand | Already used for all commands |
| gopkg.in/yaml.v3 | v3.0.1 | tracking.yaml read/write | Already in go.mod; OpenSpec standard |
| go.yaml.in/yaml/v3 | v3.0.4 | Same underlying YAML impl | Cobra's transitive dep; both resolve to same |
| golang.org/x/sys | v0.30.0 | filepath.WalkDir for scan | stdlib sufficient; already present |
| github.com/charmbracelet/lipgloss | v1.1.0 | Scan progress output styling | Already used in verify command |
| encoding/json | stdlib | Scan context JSON output | Used throughout codebase |

### No New Dependencies Required

All required functionality is available in the current dependency tree. The scan command uses `filepath.WalkDir` (stdlib), JSON encoding (stdlib), and YAML (already in go.mod).

**Version verification (current as of 2026-03-24):**
- go.mod confirmed: all versions above are already resolved in go.sum
- No `go get` commands needed for Phase 4

---

## Architecture Patterns

### Recommended Project Structure (new files only)

```
.goreleaser.yaml                     # GoReleaser config — new file at project root
plugin/                              # Distributable plugin directory — new
├── .claude-plugin/
│   └── plugin.json                  # Plugin manifest
├── commands/                        # Slash command definitions
│   ├── mysd-scan.md                 # new — /mysd:scan SKILL.md
│   └── [copy of .claude/commands/]  # all existing commands
├── agents/                          # Agent definitions
│   ├── mysd-scanner.md              # new — scanner agent
│   └── [copy of .claude/agents/]    # all existing agents
└── hooks/
    └── hooks.json                   # SessionStart version check hook

cmd/scan.go                          # New cobra command — context-only pattern
internal/scanner/
├── scanner.go                       # WalkDir + exclusion logic + JSON output
└── scanner_test.go
internal/roadmap/
├── roadmap.go                       # tracking.yaml read/write + Mermaid generation
└── roadmap_test.go
```

### Pattern 1: Context-Only for Scan (HIGH confidence — verified from existing code)

The scan command mirrors `cmd/verify.go` exactly. The binary outputs JSON for AI consumption; it does NOT invoke AI itself.

```go
// cmd/scan.go — mirrors cmd/verify.go structure
var scanCmd = &cobra.Command{
    Use:   "scan",
    Short: "Scan codebase and output metadata JSON for spec generation",
    RunE:  runScan,
}

func init() {
    scanCmd.Flags().Bool("context-only", false, "Output scan context as JSON (for /mysd:scan agent consumption)")
    scanCmd.Flags().StringSlice("exclude", nil, "Directories/files to exclude from scan")
    rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
    contextOnly, _ := cmd.Flags().GetBool("context-only")
    exclude, _ := cmd.Flags().GetStringSlice("exclude")

    if !contextOnly {
        return fmt.Errorf("usage: mysd scan --context-only [--exclude dir1,dir2]")
    }

    return runScanContextOnly(cmd.OutOrStdout(), ".", exclude)
}
```

### Pattern 2: Scanner Internal Package (HIGH confidence — stdlib patterns)

```go
// internal/scanner/scanner.go
type ScanContext struct {
    RootDir      string         `json:"root_dir"`
    Packages     []PackageInfo  `json:"packages"`
    ExistingSpecs []string      `json:"existing_specs"` // already has .specs/changes/{name}/
    ExcludedDirs []string       `json:"excluded_dirs"`
    TotalFiles   int            `json:"total_files"`
}

type PackageInfo struct {
    Name       string   `json:"name"`
    ImportPath string   `json:"import_path"`
    Dir        string   `json:"dir"`
    GoFiles    []string `json:"go_files"`
    TestFiles  []string `json:"test_files"`
    HasSpec    bool     `json:"has_spec"`  // true if .specs/changes/{name}/ exists
}

// BuildScanContext walks the directory tree and returns ScanContext JSON.
// exclude is a list of directory names to skip (e.g., vendor, testdata, .git).
func BuildScanContext(root string, exclude []string) (ScanContext, error) {
    // Use filepath.WalkDir — stdlib, no deps
    // Skip: hidden dirs (prefix "."), exclude list, non-Go dirs
    // For each Go package found: populate PackageInfo
    // Check .specs/changes/{pkg.Name}/ to set HasSpec
}
```

### Pattern 3: Roadmap Tracking (HIGH confidence — verified YAML + SaveState patterns)

The roadmap package hooks into `state.SaveState`. The cleanest integration: after every call to `state.SaveState` in cmd/ files, also call `roadmap.UpdateTracking`.

```go
// internal/roadmap/roadmap.go
type TrackingFile struct {
    SchemaVersion string         `yaml:"schema_version"`
    UpdatedAt     time.Time      `yaml:"updated_at"`
    Changes       []ChangeRecord `yaml:"changes"`
}

type ChangeRecord struct {
    Name            string     `yaml:"name"`
    Status          string     `yaml:"status"`
    StartedAt       *time.Time `yaml:"started_at,omitempty"`
    CompletedAt     *time.Time `yaml:"completed_at,omitempty"`
    TotalTasks      int        `yaml:"total_tasks"`
    CompletedTasks  int        `yaml:"completed_tasks"`
    MustTotal       int        `yaml:"must_total"`
    MustPassed      int        `yaml:"must_passed"`
    VerifyPassed    *bool      `yaml:"verify_passed,omitempty"`
}

// UpdateTracking reads tracking.yaml, upserts the record for ws.ChangeName,
// saves tracking.yaml, then regenerates timeline.md.
func UpdateTracking(specsDir string, ws state.WorkflowState) error
```

Integration point — call after SaveState in each command:
```go
// In cmd/propose.go, cmd/spec.go, cmd/execute.go, cmd/verify.go, cmd/archive.go:
if err := state.SaveState(specsDir, ws); err != nil {
    return err
}
// New: also update roadmap tracking
_ = roadmap.UpdateTracking(specsDir, ws) // best-effort, never fatal
```

**IMPORTANT: `specsDir` is `.specs/`, but tracking lives in `.mysd/roadmap/`.** The roadmap package must resolve `.mysd/` relative to the project root (one level up from `.specs/`), not inside specsDir.

### Pattern 4: Mermaid Gantt Generation (HIGH confidence — string template)

```go
// internal/roadmap/mermaid.go
func GenerateMermaid(tf TrackingFile) string {
    // Use text/template (stdlib) — no goldmark needed
    // Output format:
    // gantt
    //   title my-ssd Roadmap
    //   dateFormat YYYY-MM-DD
    //   section Changes
    //     change-name :done, 2026-01-01, 2026-01-15
    //     another-change :active, 2026-01-16, 2026-02-01
}
```

Mermaid gantt date format: `YYYY-MM-DD`. Status values: `done`, `active`, `crit` (for blocked/failed), or no modifier for future items.

### Pattern 5: Plugin Manifest Structure (HIGH confidence — verified from official Anthropic docs)

**CRITICAL FINDING:** The Claude Code plugin system has been updated. The current structure (March 2026) differs from earlier documentation referenced in CLAUDE.md:

| CLAUDE.md (old assumption) | Current official format (verified) |
|---------------------------|-------------------------------------|
| `plugin.json` at plugin root | `.claude-plugin/plugin.json` manifest directory |
| inline hooks in plugin.json | `hooks/hooks.json` separate file |
| `commands/` path was implicit | `commands/` at plugin root (NOT inside `.claude-plugin/`) |

Current verified plugin directory layout:
```
plugin/
├── .claude-plugin/
│   └── plugin.json          # metadata only: name, version, description, author
├── commands/                # .md files → /mysd:* commands
│   ├── mysd-propose.md
│   ├── mysd-spec.md
│   ├── mysd-design.md
│   ├── mysd-plan.md
│   ├── mysd-execute.md
│   ├── mysd-verify.md
│   ├── mysd-archive.md
│   ├── mysd-status.md
│   ├── mysd-ff.md
│   ├── mysd-ffe.md
│   ├── mysd-init.md
│   ├── mysd-capture.md
│   └── mysd-scan.md         # new in Phase 4
├── agents/
│   ├── mysd-spec-writer.md
│   ├── mysd-designer.md
│   ├── mysd-planner.md
│   ├── mysd-executor.md
│   ├── mysd-fast-forward.md
│   ├── mysd-verifier.md
│   └── mysd-scanner.md      # new in Phase 4
└── hooks/
    └── hooks.json           # SessionStart version check
```

plugin.json (minimal — name is required, rest is metadata):
```json
{
  "name": "mysd",
  "version": "1.0.0",
  "description": "Spec-Driven Development for AI programming — integrates OpenSpec SDD with a planning/execution/verification engine",
  "author": {
    "name": "my-ssd contributors"
  },
  "repository": "https://github.com/[owner]/mysd"
}
```

SessionStart hook (advisory only — non-blocking):
```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "mysd --version >/dev/null 2>&1 || echo 'mysd binary not found. Install: go install github.com/[owner]/mysd@latest'"
          }
        ]
      }
    ]
  }
}
```

**CRITICAL WARNING:** Hook command exit codes matter. If the hook command exits non-zero, it may block session start. Use `|| echo ...` pattern to ensure the hook always exits 0 (advisory display only). Verify this behavior from hooks documentation before finalizing.

### Pattern 6: GoReleaser Configuration (HIGH confidence — verified from goreleaser.com)

```yaml
# .goreleaser.yaml
version: 2
project_name: mysd

before:
  hooks:
    - go mod tidy

builds:
  - id: mysd
    main: .
    binary: mysd
    env:
      - CGO_ENABLED=0
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64  # Windows arm64 excluded per D-06

archives:
  - id: mysd
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "checksums.txt"
  algorithm: sha256

release:
  github:
    owner: "{{ .Env.GITHUB_REPOSITORY_OWNER }}"
    name: mysd

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
```

**IMPORTANT: `main.go` version variable.** The current `main.go` has no `version` variable declared. The `Makefile` passes `-X main.version=$(VERSION)` but there is no `var version string` in `main.go`. This must be added:

```go
// main.go — add these variables
var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)
```

Cobra's `rootCmd.Version` must also be set. Add to `cmd/root.go`:
```go
// SetVersion allows main.go to inject build-time version into the root command
func SetVersion(v string) {
    rootCmd.Version = v
}
```

Then in `main.go`:
```go
func main() {
    cmd.SetVersion(version)
    cmd.Execute()
}
```

**`go install` support:** For `go install github.com/[owner]/mysd@latest` to work, the `module` declaration in `go.mod` must match the GitHub repository path (`github.com/[owner]/mysd`). Currently `go.mod` declares `module github.com/mysd` — this needs to be updated to the actual repo path before `go install` can function. This is a prerequisite for DIST-03.

### Anti-Patterns to Avoid

- **Plugin files inside `.claude-plugin/`:** Only `plugin.json` goes in `.claude-plugin/`. Commands, agents, hooks go at plugin root.
- **Blocking SessionStart hooks:** Exit codes on hooks matter. The version check hook MUST always exit 0 — use `|| true` or `|| echo ...` fallback to ensure advisory-only behavior.
- **Roadmap update as fatal error:** roadmap.UpdateTracking must never fail a state transition. Always treat as best-effort (`_ = roadmap.UpdateTracking(...)`).
- **Using `os.Getwd()` in scanner for root detection:** Use the passed `root` parameter and `spec.DetectSpecDir` pattern for consistency.
- **Hardcoding `.specs/` to derive `.mysd/`:** Use `filepath.Join(filepath.Dir(specsDir), ".mysd")` to get `.mysd/` from specsDir. Do not assume specsDir is always `.specs/` — it can be `openspec/`.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| YAML serialization for tracking.yaml | Custom string builder | gopkg.in/yaml.v3 (already in go.mod) | Handles time.Time marshaling, omitempty, nested structs correctly |
| Mermaid chart | External template library | text/template (stdlib) | Simple string interpolation only; no complex logic needed |
| Plugin manifest | Custom JSON format | Claude Code `.claude-plugin/plugin.json` schema (verified) | Native integration with claude --plugin-dir, /plugin commands |
| Cross-platform binary builds | Custom CI matrix | goreleaser v2 | Handles archive formats, checksums, GitHub Releases creation atomically |
| Go package discovery | Custom AST parser | filepath.WalkDir + go/build.Context | WalkDir finds all .go files; go/build.Context resolves package names without full AST |
| Version variable injection | Runtime version detection | ldflags -X at build time (goreleaser) | Static linking at compile time; no runtime dependency on git or file reads |

**Key insight:** The scan command's job is NOT to analyze code semantics — that's the AI agent's job. The binary only needs to walk the directory tree and emit structured metadata (file lists, package names, import paths). Keep the scanner simple: WalkDir + JSON output.

---

## Common Pitfalls

### Pitfall 1: go.mod Module Path Mismatch
**What goes wrong:** `go install github.com/[owner]/mysd@latest` returns "no such module" or installs wrong binary.
**Why it happens:** `go.mod` currently declares `module github.com/mysd` which is not a real importable path. The module path must exactly match the GitHub URL.
**How to avoid:** Update `go.mod` module declaration to `module github.com/[actual-owner]/mysd` and run `go mod tidy`. All internal imports (`import "github.com/mysd/internal/..."`) must also be updated to the new path.
**Warning signs:** `go install` fails with "404" or version mismatch errors.

### Pitfall 2: Plugin Structure — Components Inside .claude-plugin/ (HIGH)
**What goes wrong:** `/mysd:scan` and other slash commands don't appear in Claude Code after plugin installation.
**Why it happens:** Placing `commands/` or `agents/` inside `.claude-plugin/` — only `plugin.json` belongs there. This is a common mistake per the official docs.
**How to avoid:** Keep `commands/`, `agents/`, `hooks/` at plugin root level. Verify with `claude --plugin-dir ./plugin` before finalizing.
**Warning signs:** Plugin loads (no error) but `/mysd:*` commands missing from autocomplete.

### Pitfall 3: SessionStart Hook Exit Code Blocks Session
**What goes wrong:** Claude Code fails to start (or starts with an error banner) when `mysd` binary is not installed.
**Why it happens:** Hook commands that exit non-zero can interfere with session startup behavior (exact behavior depends on Claude Code version).
**How to avoid:** Always use `command || true` or `command 2>/dev/null || echo "..."` so the hook exits 0 regardless of binary presence.
**Warning signs:** Users without `mysd` installed report Claude Code failing to start after plugin installation.

### Pitfall 4: Roadmap .mysd/ Directory Derivation
**What goes wrong:** `roadmap.UpdateTracking` writes to `.specs/.mysd/roadmap/` instead of `.mysd/roadmap/` at project root.
**Why it happens:** Passing `specsDir` (`.specs/`) and appending `.mysd/` to it.
**How to avoid:** Derive project root as `filepath.Dir(specsDir)`, then `filepath.Join(projectRoot, ".mysd", "roadmap")`. Test with both `.specs/` and `openspec/` spec dirs.
**Warning signs:** tracking.yaml created inside specs directory instead of project root.

### Pitfall 5: Scan Agent Over-Generates Specs
**What goes wrong:** AI agent creates specs for every individual file, overwhelming the user.
**Why it happens:** Unclear granularity instructions in mysd-scanner.md agent prompt.
**How to avoid:** SKILL.md and agent prompt must clearly state: one spec per Go package/major module. The `PackageInfo.HasSpec` field from context JSON signals which packages already have specs (skip them per D-03).
**Warning signs:** Dozens of near-identical spec files generated in `.specs/changes/`.

### Pitfall 6: GoReleaser v2 `version:` Field Required
**What goes wrong:** GoReleaser v2 fails with "unknown field" or version error on `.goreleaser.yaml`.
**Why it happens:** GoReleaser v2 changed the config format from v1. The `version: 2` field at the top of `.goreleaser.yaml` is required for v2 syntax.
**How to avoid:** Always start `.goreleaser.yaml` with `version: 2` as the first line.
**Warning signs:** `goreleaser check` reports validation errors.

### Pitfall 7: Windows arm64 Binary Combination
**What goes wrong:** GoReleaser fails building `windows/arm64` combination.
**Why it happens:** While Go supports `GOOS=windows GOARCH=arm64`, it's an uncommon target and may require specific Go toolchain setup.
**How to avoid:** Add `ignore: - {goos: windows, goarch: arm64}` to builds config per D-06.
**Warning signs:** GoReleaser CI fails on `windows/arm64` build step.

---

## Code Examples

Verified patterns from official sources and existing codebase:

### Scanner WalkDir Pattern (stdlib)
```go
// Source: Go stdlib filepath.WalkDir — established pattern
func buildPackageList(root string, excludeDirs map[string]bool) ([]PackageInfo, error) {
    var packages []PackageInfo
    seen := map[string]bool{} // deduplicate package dirs

    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            name := d.Name()
            // Skip hidden dirs, excluded dirs, and common non-source dirs
            if strings.HasPrefix(name, ".") || excludeDirs[name] {
                return filepath.SkipDir
            }
            return nil
        }
        if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
            return nil
        }
        dir := filepath.Dir(path)
        if !seen[dir] {
            seen[dir] = true
            rel, _ := filepath.Rel(root, dir)
            packages = append(packages, PackageInfo{Dir: rel})
        }
        return nil
    })
    return packages, err
}
```

### Mermaid Gantt Template (stdlib text/template)
```go
// Source: Mermaid.js gantt syntax — verified from mermaid.js.org
const ganttTemplate = `gantt
    title my-ssd Roadmap
    dateFormat YYYY-MM-DD
    {{- range .Changes}}
    section {{.Name}}
        {{.Name}} :{{ganttStatus .}}, {{formatDate .StartedAt}}, {{formatDate .CompletedAt}}
    {{- end}}`

// Status mapping:
// archived  → done
// verified  → done
// executed  → active
// proposed/specced/designed/planned → (no modifier = future/scheduled)
```

### hooks.json SessionStart (verified from official Anthropic docs)
```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "mysd --version >/dev/null 2>&1 && echo 'mysd: OK' || echo 'WARNING: mysd binary not found. Install: go install github.com/[owner]/mysd@latest'"
          }
        ]
      }
    ]
  }
}
```

### UpdateTracking Integration Point (derived from existing SaveState pattern)
```go
// Best-effort hook — insert after every state.SaveState call
// Pattern matches Phase 03's best-effort ARCHIVED-STATE.json (Pitfall 5 in Phase 03)
if err := state.SaveState(specsDir, ws); err != nil {
    return fmt.Errorf("save state: %w", err)
}
// Roadmap tracking — best-effort, never blocks state transition
if trackErr := roadmap.UpdateTracking(specsDir, ws); trackErr != nil {
    fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `plugin.json` at plugin root | `.claude-plugin/plugin.json` (manifest in subdirectory) | 2025 (Claude Code plugin system redesign) | Plugin directory structure is different from CLAUDE.md research — must use new format |
| Inline hooks in plugin.json | Separate `hooks/hooks.json` file | Same redesign | Hook config is in dedicated file, not embedded in manifest |
| `commands/` = slash commands (old naming) | `skills/` (SKILL.md structure) preferred for new work, `commands/` still supported | 2025 | Either `commands/` or `skills/` works; project already uses `commands/` pattern |
| plugin.json min_version field (D-15 assumption) | No `min_version` field in official schema | Verified 2026-03 | Version check must be done in SessionStart hook script, not declarative config |

**Deprecated/outdated from CLAUDE.md research:**
- `plugin.json` at root (old format): replaced by `.claude-plugin/plugin.json`
- Inline `hooks` array in `plugin.json`: replaced by `hooks/hooks.json`
- `min_version` field (CONTEXT.md D-15): this field does not exist in the official schema. Version check logic must live in the hooks.json SessionStart command script.

---

## Open Questions

1. **Hook exit code behavior for SessionStart**
   - What we know: hooks/hooks.json supports SessionStart event; `type: command` runs shell commands
   - What's unclear: Whether a non-zero exit from a SessionStart hook blocks session startup or merely logs a warning
   - Recommendation: Use `command || true` defensively regardless; test manually with `claude --plugin-dir ./plugin` before finalizing

2. **go.mod module path update scope**
   - What we know: Current `go.mod` declares `module github.com/mysd`; all 20+ internal import paths use `github.com/mysd/internal/...`
   - What's unclear: The actual target GitHub repository owner/name (not specified in CONTEXT.md)
   - Recommendation: Plan a "module path update" task that uses `grep -r github.com/mysd` + sed to update all import paths atomically; leave `[owner]` as a placeholder in research, fill in at implementation

3. **`mysd --version` flag currently not wired**
   - What we know: `rootCmd` in cmd/root.go does not set `.Version`; Makefile passes `-ldflags "-X main.version=$(VERSION)"` but `main.go` has no `var version string`
   - What's unclear: Whether `rootCmd.Version = ""` causes `--version` to output empty string or panic
   - Recommendation: Plan task to add `var version = "dev"` to main.go and `cmd.SetVersion(version)` call; test with `go build && ./mysd --version`

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|-------------|-----------|---------|---------|
| Go toolchain | Building scan command | Yes | go1.25.5 | — |
| goreleaser | DIST-03 binary releases | No | — | Install: `go install github.com/goreleaser/goreleaser/v2@latest` |
| git | GoReleaser release tagging | Yes (assumed) | — | — |
| claude CLI | Plugin testing `--plugin-dir` | Unknown | — | Manual directory inspection |

**Missing dependencies with no fallback:**
- goreleaser: required for DIST-03. Plan must include installation step or CI-only (GitHub Actions goreleaser action).

**Missing dependencies with fallback:**
- claude CLI for plugin testing: if not available locally, plugin structure can be validated with `claude plugin validate` once installed, or via directory structure inspection.

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go testing stdlib + testify v1.11.1 |
| Config file | none (go test ./...) |
| Quick run command | `go test ./internal/scanner/... ./internal/roadmap/... ./cmd/... -count=1` |
| Full suite command | `go test ./... -count=1` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| WCMD-09 | `mysd scan --context-only` outputs valid JSON with packages | unit | `go test ./internal/scanner/... -run TestBuildScanContext` | No — Wave 0 |
| WCMD-09 | Existing specs are skipped (HasSpec=true) | unit | `go test ./internal/scanner/... -run TestSkipExistingSpecs` | No — Wave 0 |
| WCMD-09 | Excluded dirs not included in output | unit | `go test ./internal/scanner/... -run TestExcludeDirs` | No — Wave 0 |
| DIST-03 | `mysd --version` outputs version string | unit | `go test ./cmd/... -run TestVersion` | No — Wave 0 |
| RMAP-01 | tracking.yaml created on state transition | unit | `go test ./internal/roadmap/... -run TestUpdateTracking` | No — Wave 0 |
| RMAP-02 | tracking.yaml contains correct fields | unit | `go test ./internal/roadmap/... -run TestTrackingFields` | No — Wave 0 |
| RMAP-03 | timeline.md contains valid Mermaid gantt | unit | `go test ./internal/roadmap/... -run TestMermaidGeneration` | No — Wave 0 |
| DIST-04 | plugin directory structure is valid | manual | `claude --plugin-dir ./plugin` (manual validation) | N/A |
| DIST-04 | plugin.json manifest schema valid | manual | `claude plugin validate` (if CLI available) | N/A |

### Sampling Rate
- **Per task commit:** `go test ./internal/scanner/... ./internal/roadmap/... -count=1`
- **Per wave merge:** `go test ./... -count=1`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `internal/scanner/scanner_test.go` — covers WCMD-09 unit tests
- [ ] `internal/scanner/scanner.go` — new package (no production code yet)
- [ ] `internal/roadmap/roadmap_test.go` — covers RMAP-01~03 unit tests
- [ ] `internal/roadmap/roadmap.go` — new package (no production code yet)
- [ ] `cmd/scan_test.go` — covers scan command JSON output
- [ ] Version variable wiring in `main.go` and `cmd/root.go` — prerequisite for DIST-03 test

---

## Sources

### Primary (HIGH confidence)
- Official Anthropic Claude Code docs — https://code.claude.com/docs/en/plugins — plugin directory structure, `.claude-plugin/plugin.json` manifest, hooks.json format, SessionStart event, commands/agents/hooks directory layout (verified 2026-03-24)
- Official Anthropic Claude Code docs — https://code.claude.com/docs/en/plugins-reference — complete plugin manifest schema, environment variables `${CLAUDE_PLUGIN_ROOT}`, hook types, component location table (verified 2026-03-24)
- goreleaser.com quick start — https://goreleaser.com/quick-start/ — minimal .goreleaser.yaml structure, archive format, checksum config (verified 2026-03-24)
- goreleaser.com Go builder — https://goreleaser.com/customization/builds/builders/go/ — ldflags version injection, CGO_ENABLED=0, goos/goarch matrix, ignore combinations (verified 2026-03-24)
- goreleaser.com ldflags cookbook — https://goreleaser.com/cookbooks/using-main.version/ — main.version variable pattern (verified 2026-03-24)
- Existing codebase: `cmd/verify.go`, `cmd/execute.go`, `cmd/archive.go`, `internal/state/state.go` — context-only pattern, SaveState integration points (direct code read)

### Secondary (MEDIUM confidence)
- WebSearch: GoReleaser v2 ldflags syntax — multiple sources agree on `version: 2` header requirement and `-X main.version={{.Version}}` pattern
- Mermaid.js gantt chart syntax — https://mermaid.js.org/syntax/gantt.html — dateFormat, section, status modifiers (done/active/crit)

### Tertiary (LOW confidence)
- Hook exit code behavior (blocking vs non-blocking): derived from common plugin hook patterns — needs empirical verification with `claude --plugin-dir`

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all deps already in go.mod, no new libraries
- Architecture: HIGH — scan follows verified context-only pattern from existing code; plugin structure verified from official docs
- GoReleaser config: HIGH — verified from goreleaser.com official docs
- Plugin structure: HIGH — verified from official Anthropic Claude Code docs (2026-03)
- Pitfalls: HIGH for structural pitfalls (verified from docs); MEDIUM for hook exit code behavior (not empirically tested)

**Research date:** 2026-03-24
**Valid until:** 2026-04-24 (plugin API may change; GoReleaser config stable)
