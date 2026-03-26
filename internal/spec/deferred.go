package spec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DeferredNote represents a single deferred note with auto-assigned ID.
type DeferredNote struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// DeferredStore holds the collection of deferred notes.
type DeferredStore struct {
	Notes []DeferredNote `json:"notes"`
}

// DeferredPath returns the path to the deferred.json file within specDir.
func DeferredPath(specDir string) string {
	return filepath.Join(specDir, "deferred.json")
}

// LoadDeferredStore reads deferred.json from specDir.
// Returns zero-value DeferredStore and nil error if the file does not exist (convention-over-config).
func LoadDeferredStore(specDir string) (DeferredStore, error) {
	path := DeferredPath(specDir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DeferredStore{}, nil
	}
	if err != nil {
		return DeferredStore{}, fmt.Errorf("deferred: %w", err)
	}
	var store DeferredStore
	if err := json.Unmarshal(data, &store); err != nil {
		return DeferredStore{}, fmt.Errorf("deferred: %w", err)
	}
	return store, nil
}

// SaveDeferredStore writes the store to deferred.json in specDir with indented JSON.
func SaveDeferredStore(specDir string, store DeferredStore) error {
	path := DeferredPath(specDir)
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("deferred: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("deferred: %w", err)
	}
	return nil
}

// Add appends a new note with auto-incremented ID (max existing + 1, never reuses deleted IDs).
// Returns the newly created DeferredNote.
func (s *DeferredStore) Add(content string) DeferredNote {
	maxID := 0
	for _, n := range s.Notes {
		if n.ID > maxID {
			maxID = n.ID
		}
	}
	note := DeferredNote{
		ID:        maxID + 1,
		Content:   content,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	s.Notes = append(s.Notes, note)
	return note
}

// Delete removes the note with the given id. Returns true if found and removed, false otherwise.
func (s *DeferredStore) Delete(id int) bool {
	for i, n := range s.Notes {
		if n.ID == id {
			s.Notes = append(s.Notes[:i], s.Notes[i+1:]...)
			return true
		}
	}
	return false
}

// CountDeferredNotes returns the number of notes in specDir's deferred.json.
// Returns 0 and nil error when the file does not exist.
func CountDeferredNotes(specDir string) (int, error) {
	store, err := LoadDeferredStore(specDir)
	if err != nil {
		return 0, err
	}
	return len(store.Notes), nil
}
