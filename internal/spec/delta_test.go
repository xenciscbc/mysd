package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectDeltaOp_Added(t *testing.T) {
	op := DetectDeltaOp("## ADDED Requirements")
	assert.Equal(t, DeltaAdded, op)
}

func TestDetectDeltaOp_Modified(t *testing.T) {
	op := DetectDeltaOp("## MODIFIED Requirements")
	assert.Equal(t, DeltaModified, op)
}

func TestDetectDeltaOp_Removed(t *testing.T) {
	op := DetectDeltaOp("## REMOVED Requirements")
	assert.Equal(t, DeltaRemoved, op)
}

func TestDetectDeltaOp_None(t *testing.T) {
	op := DetectDeltaOp("## Requirements")
	assert.Equal(t, DeltaNone, op)
}

func TestDetectDeltaOp_H3(t *testing.T) {
	op := DetectDeltaOp("### ADDED New Feature")
	assert.Equal(t, DeltaAdded, op)
}

func TestDetectDeltaOp_Renamed(t *testing.T) {
	op := DetectDeltaOp("## RENAMED Requirements")
	assert.Equal(t, DeltaRenamed, op)
}

func TestParseDelta_Renamed(t *testing.T) {
	body := `## RENAMED Requirements

### FROM: Old Name
### TO: New Name
`
	_, _, _, renamed := ParseDelta(body)
	assert.Len(t, renamed, 1)
	assert.Equal(t, "Old Name", renamed[0].From)
	assert.Equal(t, "New Name", renamed[0].To)
}

func TestParseDelta_RenamedMultiple(t *testing.T) {
	body := `## RENAMED Requirements

### FROM: First Old
### TO: First New
### FROM: Second Old
### TO: Second New
`
	_, _, _, renamed := ParseDelta(body)
	assert.Len(t, renamed, 2)
	assert.Equal(t, "First Old", renamed[0].From)
	assert.Equal(t, "First New", renamed[0].To)
	assert.Equal(t, "Second Old", renamed[1].From)
	assert.Equal(t, "Second New", renamed[1].To)
}

func TestParseDelta_MixedAddedRenamed(t *testing.T) {
	body := `## ADDED Requirements

The system MUST support new feature.

## RENAMED Requirements

### FROM: Old Auth
### TO: New Auth
`
	added, _, _, renamed := ParseDelta(body)
	assert.Len(t, added, 1)
	assert.Len(t, renamed, 1)
	assert.Equal(t, "Old Auth", renamed[0].From)
	assert.Equal(t, "New Auth", renamed[0].To)
}

func TestParseDelta_MultiSection(t *testing.T) {
	body := `## ADDED Requirements

The system MUST support dark mode.
The system SHOULD detect OS preference.

## MODIFIED Requirements

The system MUST update the existing color scheme implementation.

## REMOVED Requirements

The old hardcoded colors MUST be removed.
`
	added, modified, removed, _ := ParseDelta(body)
	assert.NotEmpty(t, added)
	assert.NotEmpty(t, modified)
	assert.NotEmpty(t, removed)
}

func TestParseDelta_OnlyAdded(t *testing.T) {
	body := `## ADDED Requirements

The system MUST support new feature.
`
	added, modified, removed, _ := ParseDelta(body)
	assert.NotEmpty(t, added)
	assert.Empty(t, modified)
	assert.Empty(t, removed)
}
