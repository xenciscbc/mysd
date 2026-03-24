package uat_test

import (
	"testing"
	"time"

	"github.com/mysd/internal/uat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUATRoundTrip tests NewUATChecklist -> WriteUAT -> ReadUAT -> update -> WriteUAT -> ReadUAT
// and verifies that run_history is preserved across writes.
func TestUATRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	changeName := "round-trip-change"
	filePath := uat.UATFilePath(tmp, changeName)

	// Step 1: Create initial checklist with 3 items
	items := []uat.UATItem{
		{Description: "User can log in"},
		{Description: "User can log out"},
		{Description: "User sees error on invalid credentials"},
	}
	checklist := uat.NewUATChecklist(changeName, items)

	// Verify initial state
	assert.Equal(t, "1", checklist.SpecVersion)
	assert.Equal(t, changeName, checklist.Change)
	assert.Len(t, checklist.Results, 3)
	for _, item := range checklist.Results {
		assert.Equal(t, "pending", item.Status)
	}

	// Step 2: WriteUAT to temp path (first write — no history yet)
	err := uat.WriteUAT(filePath, checklist)
	require.NoError(t, err)

	// Step 3: ReadUAT — assert 3 items, all status "pending"
	loaded, err := uat.ReadUAT(filePath)
	require.NoError(t, err)
	assert.Equal(t, changeName, loaded.Change)
	assert.Len(t, loaded.Results, 3)
	for _, item := range loaded.Results {
		assert.Equal(t, "pending", item.Status)
	}
	assert.Len(t, loaded.RunHistory, 0, "no history on first write")

	// Step 4: Update items — item 1 passes, item 2 fails with notes
	loaded.Results[0].Status = "pass"
	loaded.Results[0].RunAt = time.Now().Format(time.RFC3339)
	loaded.Results[1].Status = "fail"
	loaded.Results[1].Notes = "Button not found in DOM"
	loaded.Results[1].RunAt = time.Now().Format(time.RFC3339)
	loaded.Summary.Pass = 1
	loaded.Summary.Fail = 1

	// Step 5: WriteUAT again — should preserve history from first write
	err = uat.WriteUAT(filePath, loaded)
	require.NoError(t, err)

	// Step 6: ReadUAT — assert RunHistory has 1 entry with original summary
	updated, err := uat.ReadUAT(filePath)
	require.NoError(t, err)
	assert.Len(t, updated.RunHistory, 1, "second write should produce 1 history entry")
	assert.Equal(t, loaded.Results[0].Status, updated.Results[0].Status)
}

// TestUATRoundTrip_MultipleRuns tests that multiple write cycles accumulate history (UAT-05).
func TestUATRoundTrip_MultipleRuns(t *testing.T) {
	tmp := t.TempDir()
	changeName := "multi-run-change"
	filePath := uat.UATFilePath(tmp, changeName)

	// Create initial checklist
	items := []uat.UATItem{
		{Description: "Feature A works"},
		{Description: "Feature B works"},
	}
	checklist := uat.NewUATChecklist(changeName, items)

	// === Run 1: initial write (no history) ===
	err := uat.WriteUAT(filePath, checklist)
	require.NoError(t, err)

	// Read and modify for run 1 completion
	run1, err := uat.ReadUAT(filePath)
	require.NoError(t, err)
	run1.Results[0].Status = "pass"
	run1.Results[1].Status = "fail"
	run1.Summary.Pass = 1
	run1.Summary.Fail = 1

	// === Run 2: write with run1 results (creates 1 history entry) ===
	err = uat.WriteUAT(filePath, run1)
	require.NoError(t, err)

	// Read and modify for run 2 completion
	run2, err := uat.ReadUAT(filePath)
	require.NoError(t, err)
	assert.Len(t, run2.RunHistory, 1, "after 2nd write, history should have 1 entry")

	run2.Results[0].Status = "pass"
	run2.Results[1].Status = "pass"
	run2.Summary.Pass = 2
	run2.Summary.Fail = 0

	// === Run 3: write with run2 results (creates 2 history entries) ===
	err = uat.WriteUAT(filePath, run2)
	require.NoError(t, err)

	// Final read: assert RunHistory has 2 entries
	final, err := uat.ReadUAT(filePath)
	require.NoError(t, err)
	assert.Len(t, final.RunHistory, 2, "after 3rd write, history should have 2 entries (two completed runs before current)")

	// Each history entry should have a timestamp
	for i, entry := range final.RunHistory {
		assert.NotEmpty(t, entry.RunAt, "history entry %d should have a timestamp", i)
	}

	// Current results should reflect the latest run
	assert.Equal(t, 2, final.Summary.Pass)
	assert.Equal(t, 0, final.Summary.Fail)
}
