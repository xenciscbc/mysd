package verifier_test

import (
	"regexp"
	"testing"

	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/verifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStableID verifies CRC32-based ID format.
func TestStableID(t *testing.T) {
	r := spec.Requirement{
		ID:         "",
		Text:       "The system MUST validate input before processing",
		Keyword:    spec.Must,
		SourceFile: "user-auth/spec.md",
	}

	id := verifier.StableID(r)

	// ID must match format: {source_file}::{keyword_lower}-{hex_hash}
	re := regexp.MustCompile(`^.+::(must|should|may)-[0-9a-f]+$`)
	assert.Regexp(t, re, id, "StableID must match pattern {source}::{keyword}-{hex}")
	assert.Contains(t, id, "user-auth/spec.md::must-")
}

func TestStableID_DifferentTexts_DifferentIDs(t *testing.T) {
	r1 := spec.Requirement{Text: "System MUST do A", Keyword: spec.Must, SourceFile: "spec.md"}
	r2 := spec.Requirement{Text: "System MUST do B", Keyword: spec.Must, SourceFile: "spec.md"}

	id1 := verifier.StableID(r1)
	id2 := verifier.StableID(r2)

	assert.NotEqual(t, id1, id2, "different texts must produce different IDs")
}

func TestStableID_SameText_SameID(t *testing.T) {
	r1 := spec.Requirement{Text: "System MUST validate", Keyword: spec.Must, SourceFile: "spec.md"}
	r2 := spec.Requirement{Text: "System MUST validate", Keyword: spec.Must, SourceFile: "spec.md"}

	assert.Equal(t, verifier.StableID(r1), verifier.StableID(r2), "same text must produce same ID (stable)")
}

// TestBuildVerificationContextFromParts verifies pure function classification.
func TestBuildVerificationContextFromParts(t *testing.T) {
	reqs := []spec.Requirement{
		{Text: "System MUST validate input", Keyword: spec.Must, SourceFile: "core/spec.md"},
		{Text: "System MUST NOT allow SQL injection", Keyword: spec.Must, SourceFile: "core/spec.md"},
		{Text: "System SHOULD log all requests", Keyword: spec.Should, SourceFile: "core/spec.md"},
		{Text: "System MAY cache responses", Keyword: spec.May, SourceFile: "core/spec.md"},
	}
	tasks := []spec.Task{
		{ID: 1, Name: "Implement validation", Status: spec.StatusPending},
		{ID: 2, Name: "Write tests", Status: spec.StatusDone},
	}

	ctx := verifier.BuildVerificationContextFromParts("my-feature", "/specs/changes/my-feature", "/specs", reqs, tasks)

	assert.Equal(t, "my-feature", ctx.ChangeName)
	assert.Len(t, ctx.MustItems, 2, "should have 2 MUST items")
	assert.Len(t, ctx.ShouldItems, 1, "should have 1 SHOULD item")
	assert.Len(t, ctx.MayItems, 1, "should have 1 MAY item")
	assert.Len(t, ctx.TasksSummary, 2, "should have 2 task entries")
}

// TestVerificationContext_ItemClassification verifies keyword classification.
func TestVerificationContext_ItemClassification(t *testing.T) {
	reqs := []spec.Requirement{
		{Text: "System MUST do X", Keyword: spec.Must, SourceFile: "spec.md"},
		{Text: "System SHOULD do Y", Keyword: spec.Should, SourceFile: "spec.md"},
		{Text: "System MAY do Z", Keyword: spec.May, SourceFile: "spec.md"},
	}

	ctx := verifier.BuildVerificationContextFromParts("test", "/dir", "/specs", reqs, nil)

	require.Len(t, ctx.MustItems, 1)
	assert.Equal(t, "MUST", ctx.MustItems[0].Keyword)

	require.Len(t, ctx.ShouldItems, 1)
	assert.Equal(t, "SHOULD", ctx.ShouldItems[0].Keyword)

	require.Len(t, ctx.MayItems, 1)
	assert.Equal(t, "MAY", ctx.MayItems[0].Keyword)
}

// TestBuildVerificationContext_EmptySpecs verifies empty slices (not nil) are returned.
func TestBuildVerificationContext_EmptySpecs(t *testing.T) {
	ctx := verifier.BuildVerificationContextFromParts("empty", "/dir", "/specs", nil, nil)

	// Must return non-nil empty slices
	assert.NotNil(t, ctx.MustItems, "MustItems must not be nil")
	assert.NotNil(t, ctx.ShouldItems, "ShouldItems must not be nil")
	assert.NotNil(t, ctx.MayItems, "MayItems must not be nil")
	assert.NotNil(t, ctx.TasksSummary, "TasksSummary must not be nil")
	assert.Len(t, ctx.MustItems, 0)
	assert.Len(t, ctx.ShouldItems, 0)
	assert.Len(t, ctx.MayItems, 0)
}

// TestVerifyItem_IDFormat verifies StableID format in constructed VerifyItems.
func TestVerifyItem_IDFormat(t *testing.T) {
	reqs := []spec.Requirement{
		{Text: "System MUST handle errors gracefully", Keyword: spec.Must, SourceFile: "api/spec.md"},
	}

	ctx := verifier.BuildVerificationContextFromParts("my-change", "/dir", "/specs", reqs, nil)

	require.Len(t, ctx.MustItems, 1)
	id := ctx.MustItems[0].ID
	re := regexp.MustCompile(`^.+::(must|should|may)-[0-9a-f]+$`)
	assert.Regexp(t, re, id, "VerifyItem ID must match stable ID format")
}

// TestBuildVerificationContext_TasksSummary verifies tasks are populated correctly.
func TestBuildVerificationContext_TasksSummary(t *testing.T) {
	tasks := []spec.Task{
		{ID: 1, Name: "Task One", Status: spec.StatusPending},
		{ID: 2, Name: "Task Two", Status: spec.StatusDone},
		{ID: 3, Name: "Task Three", Status: spec.StatusInProgress},
	}

	ctx := verifier.BuildVerificationContextFromParts("change", "/dir", "/specs", nil, tasks)

	require.Len(t, ctx.TasksSummary, 3)
	assert.Equal(t, 1, ctx.TasksSummary[0].ID)
	assert.Equal(t, "Task One", ctx.TasksSummary[0].Name)
	assert.Equal(t, "pending", ctx.TasksSummary[0].Status)
	assert.Equal(t, 2, ctx.TasksSummary[1].ID)
	assert.Equal(t, "done", ctx.TasksSummary[1].Status)
}
