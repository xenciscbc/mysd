package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
)

// TestNoteList_EmptyNoFile verifies `mysd note` with no file prints "No deferred notes."
func TestNoteList_EmptyNoFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Create a .specs dir so DetectSpecDir works
	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, ".specs"), 0755))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"note"})

	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No deferred notes", "empty list should say no deferred notes")
}

// TestNoteAdd_CreatesNoteWithID verifies `mysd note add` creates deferred.json with note.
func TestNoteAdd_CreatesNoteWithID(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, ".specs"), 0755))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"note", "add", "implement caching"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Added note #1", "output should confirm note added with ID 1")
	assert.Contains(t, output, "implement caching", "output should echo the content")

	// Verify file was created
	deferredPath := filepath.Join(tmpDir, ".specs", "deferred.json")
	assert.FileExists(t, deferredPath)

	data, err := os.ReadFile(deferredPath)
	require.NoError(t, err)
	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))
	notes := raw["notes"].([]interface{})
	assert.Len(t, notes, 1)
}

// TestNoteAdd_MultipleWordsJoined verifies multiple args are joined as content.
func TestNoteAdd_MultipleWordsJoined(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, ".specs"), 0755))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"note", "add", "add", "rate", "limiting"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "add rate limiting", "multi-word args should be joined")
}

// TestNoteDelete_RemovesNote verifies `mysd note delete <id>` removes note and confirms.
func TestNoteDelete_RemovesNote(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, ".specs"), 0755))

	// First add a note
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"note", "add", "to be deleted"})
	require.NoError(t, rootCmd.Execute())

	// Now delete it
	buf.Reset()
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"note", "delete", "1"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Deleted note #1", "output should confirm deletion")

	// Verify file has empty notes
	deferredPath := filepath.Join(tmpDir, ".specs", "deferred.json")
	data, err := os.ReadFile(deferredPath)
	require.NoError(t, err)
	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))
	notes, ok := raw["notes"]
	if ok && notes != nil {
		assert.Empty(t, notes.([]interface{}), "notes should be empty after delete")
	}
}

// TestNoteDelete_NotFound verifies `mysd note delete <id>` returns error for missing ID.
func TestNoteDelete_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, ".specs"), 0755))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"note", "delete", "999"})

	err = rootCmd.Execute()
	assert.Error(t, err, "deleting non-existent note should return error")
	assert.Contains(t, err.Error(), "not found", "error should mention not found")
}

// TestNoteList_ShowsNotesAfterAdd verifies listing shows notes with ID and content.
func TestNoteList_ShowsNotesAfterAdd(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, ".specs"), 0755))

	// Add two notes
	for _, content := range []string{"first note", "second note"} {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)
		rootCmd.SetArgs([]string{"note", "add", content})
		require.NoError(t, rootCmd.Execute())
	}

	// List
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"note"})
	require.NoError(t, rootCmd.Execute())

	output := buf.String()
	assert.Contains(t, output, "[1]", "should show ID 1")
	assert.Contains(t, output, "first note", "should show first note content")
	assert.Contains(t, output, "[2]", "should show ID 2")
	assert.Contains(t, output, "second note", "should show second note content")
}

// TestStatusDeferred_ShowsCountWhenNotesExist verifies status shows deferred count line.
func TestStatusDeferred_ShowsCountWhenNotesExist(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	// Setup full change structure for status
	tasks := spec.TasksFrontmatterV2{
		SpecVersion: "1",
		Total:       1,
		Completed:   0,
		Tasks: []spec.TaskEntry{
			{ID: 1, Name: "Task One", Status: spec.StatusPending},
		},
	}
	setupTestChange(t, tmpDir, tasks, state.PhasePlanned)

	// Write deferred.json with 2 notes
	specsDir := filepath.Join(tmpDir, ".specs")
	deferredJSON := `{"notes":[{"id":1,"content":"note A","created_at":"2026-03-26T00:00:00Z"},{"id":2,"content":"note B","created_at":"2026-03-26T00:00:00Z"}]}`
	require.NoError(t, os.WriteFile(filepath.Join(specsDir, "deferred.json"), []byte(deferredJSON), 0644))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"status"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, strings.Contains(output, "Deferred notes: 2"),
		"status output should show 'Deferred notes: 2', got: %s", output)
}

// TestStatusDeferred_HidesLineWhenNoNotes verifies status does NOT show deferred line when empty.
func TestStatusDeferred_HidesLineWhenNoNotes(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()
	require.NoError(t, os.Chdir(tmpDir))

	tasks := spec.TasksFrontmatterV2{
		SpecVersion: "1",
		Total:       1,
		Completed:   0,
		Tasks: []spec.TaskEntry{
			{ID: 1, Name: "Task One", Status: spec.StatusPending},
		},
	}
	setupTestChange(t, tmpDir, tasks, state.PhasePlanned)
	// No deferred.json written

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"status"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.NotContains(t, output, "Deferred notes:", "status should NOT show deferred line when no notes")
}
