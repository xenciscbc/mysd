package spec

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadDeferredStore_NotExist verifies zero-value DeferredStore returned for missing file.
func TestLoadDeferredStore_NotExist(t *testing.T) {
	specDir := t.TempDir()

	store, err := LoadDeferredStore(specDir)
	require.NoError(t, err)
	assert.Equal(t, DeferredStore{}, store, "missing file should return zero-value DeferredStore")
	assert.Empty(t, store.Notes, "notes should be empty for missing file")
}

// TestDeferredStore_AddAutoIncrement verifies Add assigns auto-increment IDs starting at 1.
func TestDeferredStore_AddAutoIncrement(t *testing.T) {
	store := DeferredStore{}

	note1 := store.Add("first note")
	assert.Equal(t, 1, note1.ID, "first note should have ID 1")
	assert.Equal(t, "first note", note1.Content)
	assert.NotEmpty(t, note1.CreatedAt, "CreatedAt should be set")

	note2 := store.Add("second note")
	assert.Equal(t, 2, note2.ID, "second note should have ID 2")

	note3 := store.Add("third note")
	assert.Equal(t, 3, note3.ID, "third note should have ID 3")
	assert.Len(t, store.Notes, 3)
}

// TestDeferredStore_AddNoReuseDeletedIDs verifies IDs are not reused after deletion.
func TestDeferredStore_AddNoReuseDeletedIDs(t *testing.T) {
	store := DeferredStore{}

	store.Add("note 1") // ID 1
	store.Add("note 2") // ID 2
	store.Add("note 3") // ID 3

	deleted := store.Delete(2)
	assert.True(t, deleted, "should return true when deleting existing note")
	assert.Len(t, store.Notes, 2, "store should have 2 notes after delete")

	note4 := store.Add("note 4")
	assert.Equal(t, 4, note4.ID, "new note should have ID 4, not reuse deleted ID 2")
}

// TestDeferredStore_DeleteFoundAndNotFound verifies Delete returns correct bool.
func TestDeferredStore_DeleteFoundAndNotFound(t *testing.T) {
	store := DeferredStore{}
	store.Add("note 1") // ID 1
	store.Add("note 2") // ID 2

	// Delete existing
	found := store.Delete(1)
	assert.True(t, found, "should return true for existing note")
	assert.Len(t, store.Notes, 1, "should have 1 note remaining")
	assert.Equal(t, 2, store.Notes[0].ID, "remaining note should be ID 2")

	// Delete non-existent
	notFound := store.Delete(999)
	assert.False(t, notFound, "should return false for non-existent ID")
	assert.Len(t, store.Notes, 1, "notes count should remain unchanged")
}

// TestSaveAndLoadDeferredStore_Roundtrip verifies all fields preserved after save+load.
func TestSaveAndLoadDeferredStore_Roundtrip(t *testing.T) {
	specDir := t.TempDir()
	store := DeferredStore{}
	store.Add("implement caching")
	store.Add("add rate limiting")

	err := SaveDeferredStore(specDir, store)
	require.NoError(t, err)

	loaded, err := LoadDeferredStore(specDir)
	require.NoError(t, err)
	assert.Len(t, loaded.Notes, 2)
	assert.Equal(t, 1, loaded.Notes[0].ID)
	assert.Equal(t, "implement caching", loaded.Notes[0].Content)
	assert.NotEmpty(t, loaded.Notes[0].CreatedAt)
	assert.Equal(t, 2, loaded.Notes[1].ID)
	assert.Equal(t, "add rate limiting", loaded.Notes[1].Content)
}

// TestCountDeferredNotes_MissingAndPopulated verifies count for missing file and populated file.
func TestCountDeferredNotes_MissingAndPopulated(t *testing.T) {
	specDir := t.TempDir()

	// Missing file -> 0
	count, err := CountDeferredNotes(specDir)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "missing file should return count 0")

	// Populate store
	store := DeferredStore{}
	store.Add("note A")
	store.Add("note B")
	store.Add("note C")
	require.NoError(t, SaveDeferredStore(specDir, store))

	count, err = CountDeferredNotes(specDir)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

// TestDeferredPath verifies the path construction.
func TestDeferredPath(t *testing.T) {
	specDir := "/some/spec/dir"
	expected := filepath.Join(specDir, "deferred.json")
	assert.Equal(t, expected, DeferredPath(specDir))
}

// TestDeferredStore_AddCreatedAtRFC3339 verifies CreatedAt is valid RFC3339.
func TestDeferredStore_AddCreatedAtRFC3339(t *testing.T) {
	store := DeferredStore{}
	before := time.Now().Add(-time.Second)
	note := store.Add("test note")
	after := time.Now().Add(time.Second)

	parsed, err := time.Parse(time.RFC3339, note.CreatedAt)
	require.NoError(t, err, "CreatedAt should be valid RFC3339")
	assert.True(t, parsed.After(before), "CreatedAt should be after test start")
	assert.True(t, parsed.Before(after), "CreatedAt should be before test end")
}

// TestSaveAndLoadDeferredStore_FileCreated verifies deferred.json file is created.
func TestSaveAndLoadDeferredStore_FileCreated(t *testing.T) {
	specDir := t.TempDir()
	store := DeferredStore{}
	store.Add("note")

	require.NoError(t, SaveDeferredStore(specDir, store))

	path := DeferredPath(specDir)
	assert.FileExists(t, path)

	// Verify it's valid JSON
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var raw map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &raw), "file should contain valid JSON")
}
