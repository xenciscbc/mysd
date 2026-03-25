package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRFC2119Keywords(t *testing.T) {
	assert.Equal(t, RFC2119Keyword("MUST"), Must, "Must should equal 'MUST'")
	assert.Equal(t, RFC2119Keyword("SHOULD"), Should, "Should should equal 'SHOULD'")
	assert.Equal(t, RFC2119Keyword("MAY"), May, "May should equal 'MAY'")
}

func TestDeltaOpConstants(t *testing.T) {
	assert.Equal(t, DeltaOp("ADDED"), DeltaAdded, "DeltaAdded should equal 'ADDED'")
	assert.Equal(t, DeltaOp("MODIFIED"), DeltaModified, "DeltaModified should equal 'MODIFIED'")
	assert.Equal(t, DeltaOp("REMOVED"), DeltaRemoved, "DeltaRemoved should equal 'REMOVED'")
	assert.Equal(t, DeltaOp(""), DeltaNone, "DeltaNone should equal ''")
}

func TestItemStatusConstants(t *testing.T) {
	assert.Equal(t, ItemStatus("pending"), StatusPending, "StatusPending should equal 'pending'")
	assert.Equal(t, ItemStatus("in_progress"), StatusInProgress, "StatusInProgress should equal 'in_progress'")
	assert.Equal(t, ItemStatus("done"), StatusDone, "StatusDone should equal 'done'")
	assert.Equal(t, ItemStatus("blocked"), StatusBlocked, "StatusBlocked should equal 'blocked'")
}

func TestSpecDirFlavorConstants(t *testing.T) {
	assert.Equal(t, SpecDirFlavor("mysd"), FlavorMySD, "FlavorMySD should equal 'mysd'")
	assert.Equal(t, SpecDirFlavor("openspec"), FlavorOpenSpec, "FlavorOpenSpec should equal 'openspec'")
	assert.Equal(t, SpecDirFlavor("none"), FlavorNone, "FlavorNone should equal 'none'")
}

func TestChangeStruct(t *testing.T) {
	c := Change{
		Name: "test-change",
		Dir:  "/some/dir",
	}
	assert.Equal(t, "test-change", c.Name)
	assert.Equal(t, "/some/dir", c.Dir)
}

func TestRequirementStruct(t *testing.T) {
	r := Requirement{
		ID:      "REQ-01",
		Text:    "System MUST validate input",
		Keyword: Must,
		DeltaOp: DeltaAdded,
		Status:  StatusPending,
	}
	assert.Equal(t, "REQ-01", r.ID)
	assert.Equal(t, Must, r.Keyword)
	assert.Equal(t, DeltaAdded, r.DeltaOp)
}

func TestErrVariables(t *testing.T) {
	assert.NotNil(t, ErrNoSpecDir)
	assert.NotNil(t, ErrInvalidTransition)
	assert.Contains(t, ErrNoSpecDir.Error(), "spec")
}

// TestTaskEntryNewFields_YAMLRoundTrip verifies TaskEntry with new fields marshals/unmarshals correctly.
func TestTaskEntryNewFields_YAMLRoundTrip(t *testing.T) {
	entry := TaskEntry{
		ID:     1,
		Name:   "Build auth",
		Status: StatusPending,
		Depends:   []int{2, 3},
		Files:     []string{"auth.go", "handler.go"},
		Satisfies: []string{"REQ-01", "REQ-02"},
		Skills:    []string{"/mysd:execute"},
	}

	data, err := yaml.Marshal(entry)
	require.NoError(t, err)

	var got TaskEntry
	require.NoError(t, yaml.Unmarshal(data, &got))

	assert.Equal(t, entry.ID, got.ID)
	assert.Equal(t, entry.Name, got.Name)
	assert.Equal(t, entry.Depends, got.Depends)
	assert.Equal(t, entry.Files, got.Files)
	assert.Equal(t, entry.Satisfies, got.Satisfies)
	assert.Equal(t, entry.Skills, got.Skills)
}

// TestTaskEntryNewFields_OmitEmpty verifies TaskEntry without new fields does NOT emit those keys.
func TestTaskEntryNewFields_OmitEmpty(t *testing.T) {
	entry := TaskEntry{
		ID:     1,
		Name:   "A",
		Status: StatusPending,
	}

	data, err := yaml.Marshal(entry)
	require.NoError(t, err)

	output := string(data)
	assert.NotContains(t, output, "depends:")
	assert.NotContains(t, output, "files:")
	assert.NotContains(t, output, "satisfies:")
	assert.NotContains(t, output, "skills:")
}
