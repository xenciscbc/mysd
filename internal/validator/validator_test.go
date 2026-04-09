package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupChangeDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "test-change")
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0644))
}

func TestValidate_MissingProposal(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")

	result := Validate(changeDir)

	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "proposal.md", result.Errors[0].Location)
	assert.Contains(t, result.Errors[0].Message, "file not found")
}

func TestValidate_MissingOpenspecYaml(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, "proposal.md", "# Proposal\n\nSome content\n")

	result := Validate(changeDir)

	assert.False(t, result.Valid)
	hasMetaError := false
	for _, e := range result.Errors {
		if e.Location == ".openspec.yaml" {
			hasMetaError = true
			break
		}
	}
	assert.True(t, hasMetaError, "expected error for missing .openspec.yaml")
}

func TestValidate_ValidChangeMinimal(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal

Some content
`)

	result := Validate(changeDir)

	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
}

func TestValidate_ProposalChangeNameMismatch(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: wrong-name
status: proposed
created: 2026-01-01
---

# Proposal
`)

	result := Validate(changeDir)

	assert.False(t, result.Valid)
	hasNameError := false
	for _, e := range result.Errors {
		if e.Location == "proposal.md" && e.Message != "" {
			if assert.ObjectsAreEqual("change name mismatch", "") {
				continue
			}
			if contains(e.Message, "change name mismatch") {
				hasNameError = true
			}
		}
	}
	assert.True(t, hasNameError, "expected change name mismatch error")
}

func TestValidate_ProposalMissingFields(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	// At least one field must be non-empty to avoid brownfield detection
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: ""
status: ""
created: ""
---

# Proposal
`)

	result := Validate(changeDir)

	assert.False(t, result.Valid)
	// Should have errors for: change, status, created (spec-version is set)
	assert.GreaterOrEqual(t, len(result.Errors), 3)
}

func TestValidate_SpecCapabilityMismatch(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)
	writeFile(t, filepath.Join(changeDir, "specs", "auth"), "spec.md", `---
spec-version: "1.0"
capability: authentication
delta: ADDED
status: pending
---

## Requirements
`)

	result := Validate(changeDir)

	assert.False(t, result.Valid)
	hasCapError := false
	for _, e := range result.Errors {
		if contains(e.Message, "capability mismatch") {
			hasCapError = true
			break
		}
	}
	assert.True(t, hasCapError, "expected capability mismatch error")
}

func TestValidate_SpecInvalidDelta(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)
	writeFile(t, filepath.Join(changeDir, "specs", "auth"), "spec.md", `---
spec-version: "1.0"
capability: auth
delta: INVALID
status: pending
---

## Requirements
`)

	result := Validate(changeDir)

	assert.False(t, result.Valid)
	hasDeltaError := false
	for _, e := range result.Errors {
		if contains(e.Message, "invalid delta") {
			hasDeltaError = true
			break
		}
	}
	assert.True(t, hasDeltaError, "expected invalid delta error")
}

func TestValidate_SpecValid(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)
	writeFile(t, filepath.Join(changeDir, "specs", "auth"), "spec.md", `---
spec-version: "1.0"
capability: auth
delta: ADDED
status: pending
---

## Requirements

The system MUST authenticate users.
`)

	result := Validate(changeDir)

	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
}

func TestValidate_TasksCountMismatch(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)
	writeFile(t, changeDir, "tasks.md", `---
spec-version: "1.0"
total: 5
completed: 0
---

- [ ] Task 1
- [ ] Task 2
- [ ] Task 3
`)

	result := Validate(changeDir)

	assert.True(t, result.Valid) // count mismatch is a warning, not error
	assert.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0].Message, "total (5) does not match actual task count (3)")
}

func TestValidate_TasksValid(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)
	writeFile(t, changeDir, "tasks.md", `---
spec-version: "1.0"
total: 3
completed: 0
---

- [ ] Task 1
- [ ] Task 2
- [ ] Task 3
`)

	result := Validate(changeDir)

	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
	assert.Empty(t, result.Warnings)
}

func TestValidate_ChangeMetaMissingSchema(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "created: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)

	result := Validate(changeDir)

	assert.False(t, result.Valid)
	hasSchemaError := false
	for _, e := range result.Errors {
		if e.Location == ".openspec.yaml" && contains(e.Message, "schema") {
			hasSchemaError = true
			break
		}
	}
	assert.True(t, hasSchemaError, "expected missing schema error")
}

func TestValidate_BrownfieldProposal(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", "# Proposal\n\nNo frontmatter here.\n")

	result := Validate(changeDir)

	// Brownfield proposal triggers a warning, not an error
	assert.True(t, result.Valid)
	assert.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0].Message, "brownfield")
}

func TestValidate_SpecDirWithoutSpecMd(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)
	// Create capability dir without spec.md
	require.NoError(t, os.MkdirAll(filepath.Join(changeDir, "specs", "orphan"), 0755))

	result := Validate(changeDir)

	assert.True(t, result.Valid) // missing spec.md in cap dir is a warning
	assert.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0].Message, "spec.md is missing")
}

func TestValidate_TasksV2Valid(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)
	writeFile(t, changeDir, "tasks.md", `---
spec-version: "1.0"
total: 2
completed: 0
tasks:
  - id: 1
    name: "Task one"
    description: "Do task one"
    spec: "core"
    status: pending
  - id: 2
    name: "Task two"
    description: "Do task two"
    spec: "core"
    status: pending
---

# Tasks
`)

	result := Validate(changeDir)

	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
	assert.Empty(t, result.Warnings)
}

func TestValidate_TasksV2CountMismatch(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, ".openspec.yaml", "schema: spec-driven\ncreated: 2026-01-01\n")
	writeFile(t, changeDir, "proposal.md", `---
spec-version: "1.0"
change: test-change
status: proposed
created: 2026-01-01
---

# Proposal
`)
	writeFile(t, changeDir, "tasks.md", `---
spec-version: "1.0"
total: 5
completed: 0
tasks:
  - id: 1
    name: "Task one"
    description: "Do task one"
    spec: "core"
    status: pending
  - id: 2
    name: "Task two"
    description: "Do task two"
    spec: "core"
    status: pending
---

# Tasks
`)

	result := Validate(changeDir)

	assert.True(t, result.Valid) // count mismatch is a warning
	assert.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0].Message, "total (5) does not match actual task count (2)")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
