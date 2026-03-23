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

func TestParseDelta_MultiSection(t *testing.T) {
	body := `## ADDED Requirements

The system MUST support dark mode.
The system SHOULD detect OS preference.

## MODIFIED Requirements

The system MUST update the existing color scheme implementation.

## REMOVED Requirements

The old hardcoded colors MUST be removed.
`
	added, modified, removed := ParseDelta(body)
	assert.NotEmpty(t, added)
	assert.NotEmpty(t, modified)
	assert.NotEmpty(t, removed)
}

func TestParseDelta_OnlyAdded(t *testing.T) {
	body := `## ADDED Requirements

The system MUST support new feature.
`
	added, modified, removed := ParseDelta(body)
	assert.NotEmpty(t, added)
	assert.Empty(t, modified)
	assert.Empty(t, removed)
}
