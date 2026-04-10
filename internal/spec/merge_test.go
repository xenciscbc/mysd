package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeMainSpec(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "spec.md")
	require.NoError(t, os.WriteFile(p, []byte(content), 0644))
	return p
}

func TestMergeSpecs_Added(t *testing.T) {
	dir := t.TempDir()
	mainSpec := writeMainSpec(t, dir, "# Spec\n\nExisting content.\n")

	delta := `## ADDED Requirements

The system MUST support dark mode.
`
	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Contains(t, result, "Existing content.")
	assert.Contains(t, result, "The system MUST support dark mode.")
}

func TestMergeSpecs_Modified(t *testing.T) {
	dir := t.TempDir()
	mainSpec := writeMainSpec(t, dir, `### Requirement: Auth

The system MUST provide basic authentication.
`)

	delta := `## MODIFIED Requirements

The system MUST provide basic authentication.
`
	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Contains(t, result, "The system MUST provide basic authentication.")
}

func TestMergeSpecs_Removed(t *testing.T) {
	dir := t.TempDir()
	mainSpec := writeMainSpec(t, dir, `### Requirement: Auth

The system MUST provide authentication.

### Requirement: Logging

The system SHOULD log events.
`)

	delta := `## REMOVED Requirements

The system MUST provide authentication.
`
	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.NotContains(t, result, "provide authentication")
	assert.Contains(t, result, "log events")
}

func TestMergeSpecs_Renamed(t *testing.T) {
	dir := t.TempDir()
	mainSpec := writeMainSpec(t, dir, `### Requirement: Old Auth

The system MUST provide authentication.
`)

	delta := `## RENAMED Requirements

### FROM: Old Auth
### TO: New Auth
`
	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Contains(t, result, "### Requirement: New Auth")
	assert.NotContains(t, result, "Old Auth")
}

func TestMergeSpecs_NoMainSpecExists(t *testing.T) {
	dir := t.TempDir()
	mainSpec := filepath.Join(dir, "nonexistent", "spec.md")

	delta := `## ADDED Requirements

The system MUST support new feature.
`
	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Contains(t, result, "The system MUST support new feature.")
}

func TestMergeSpecs_HeadingMismatch(t *testing.T) {
	dir := t.TempDir()
	mainSpec := writeMainSpec(t, dir, `### Requirement: Auth

The system MUST provide authentication.
`)

	delta := `## MODIFIED Requirements

The system MUST provide NONEXISTENT requirement.
`
	_, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.NotEmpty(t, warnings)
	assert.Contains(t, warnings[0], "MODIFIED")
	assert.Contains(t, warnings[0], "not found")
}

func TestMergeSpecs_NewSpecGetsFrontmatter(t *testing.T) {
	dir := t.TempDir()
	mainSpec := filepath.Join(dir, "spec.md") // does not exist

	delta := `## ADDED Requirements

The system MUST support new feature.
`
	AppVersion = "1.2.3"
	defer func() { AppVersion = "dev" }()

	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Contains(t, result, "version: 1.0.0")
	assert.Contains(t, result, "generatedBy: mysd v1.2.3")
	assert.Contains(t, result, "---")
}

func TestMergeSpecs_ModifiedIncrementsVersion(t *testing.T) {
	dir := t.TempDir()
	mainContent := `---
version: 1.0.0
generatedBy: mysd v1.0.0
---
### Requirement: Auth

The system MUST provide authentication.
`
	mainSpec := writeMainSpec(t, dir, mainContent)

	delta := `## MODIFIED Requirements

The system MUST provide authentication.
`
	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Contains(t, result, "version: 1.1.0")
}

func TestIncrementMinorVersion(t *testing.T) {
	assert.Equal(t, "1.1.0", incrementMinorVersion("1.0.0"))
	assert.Equal(t, "1.3.0", incrementMinorVersion("1.2.5"))
	assert.Equal(t, "2.1.0", incrementMinorVersion("2.0.1"))
	assert.Equal(t, "1.1.0", incrementMinorVersion(""))
	assert.Equal(t, "1.1.0", incrementMinorVersion("invalid"))
}

func TestMergeSpecs_FallbackAdded(t *testing.T) {
	dir := t.TempDir()
	mainSpec := filepath.Join(dir, "spec.md") // does not exist

	delta := `---
spec-version: "1.0"
capability: new-cap
delta: ADDED
status: pending
---

### Requirement: New Feature

The system MUST support new feature.
`
	AppVersion = "1.0.0"
	defer func() { AppVersion = "dev" }()

	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Contains(t, result, "The system MUST support new feature.")
	assert.Contains(t, result, "version: 1.0.0")
	assert.Contains(t, result, "capability: new-cap")
}

func TestMergeSpecs_FallbackModified(t *testing.T) {
	dir := t.TempDir()
	mainContent := `---
spec-version: "1.0"
capability: auth
version: 1.0.0
---

### Requirement: Auth

The system MUST provide basic authentication.
`
	mainSpec := writeMainSpec(t, dir, mainContent)

	delta := `---
spec-version: "1.0"
capability: auth
delta: MODIFIED
status: pending
---

### Requirement: Auth

The system MUST provide OAuth2 authentication.
`
	result, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Contains(t, result, "The system MUST provide OAuth2 authentication.")
	assert.Contains(t, result, "version: 1.1.0")
	// Old content should be replaced
	assert.NotContains(t, result, "basic authentication")
}

func TestMergeSpecs_FallbackEmptyDelta(t *testing.T) {
	dir := t.TempDir()
	mainSpec := writeMainSpec(t, dir, "# Spec\n")

	delta := `---
spec-version: "1.0"
capability: auth
delta: ""
status: pending
---

Some content without delta headings.
`
	_, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.NotEmpty(t, warnings)
	assert.Contains(t, warnings[0], "no parseable operations")
}

func TestMergeSpecs_RenamedHeadingMismatch(t *testing.T) {
	dir := t.TempDir()
	mainSpec := writeMainSpec(t, dir, `### Requirement: Auth

The system MUST provide authentication.
`)

	delta := `## RENAMED Requirements

### FROM: Nonexistent
### TO: New Name
`
	_, warnings, err := MergeSpecs(mainSpec, delta)
	require.NoError(t, err)
	assert.NotEmpty(t, warnings)
	assert.Contains(t, warnings[0], "RENAMED")
	assert.Contains(t, warnings[0], "not found")
}
