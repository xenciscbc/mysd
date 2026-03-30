package spec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// verificationStatusFile is the sidecar JSON file name for verification status.
const verificationStatusFile = "verification-status.json"

// VerificationStatus is the sidecar JSON tracking verification state per requirement.
// It is stored in {changeDir}/verification-status.json and does NOT modify spec.md.
type VerificationStatus struct {
	ChangeName   string                `json:"change_name"`
	VerifiedAt   time.Time             `json:"verified_at"`
	Requirements map[string]ItemStatus `json:"requirements"`
}

// ReadVerificationStatus reads verification-status.json from changeDir.
// Returns a zero-value VerificationStatus with an empty (non-nil) Requirements map
// if the file does not exist. Returns an error only on unexpected I/O failures.
func ReadVerificationStatus(changeDir string) (VerificationStatus, error) {
	path := filepath.Join(changeDir, verificationStatusFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return VerificationStatus{Requirements: map[string]ItemStatus{}}, nil
	}
	if err != nil {
		return VerificationStatus{}, fmt.Errorf("read verification-status: %w", err)
	}

	var vs VerificationStatus
	if err := json.Unmarshal(data, &vs); err != nil {
		return VerificationStatus{}, fmt.Errorf("parse verification-status: %w", err)
	}
	if vs.Requirements == nil {
		vs.Requirements = map[string]ItemStatus{}
	}
	return vs, nil
}

// WriteVerificationStatus writes a VerificationStatus to {changeDir}/verification-status.json.
func WriteVerificationStatus(changeDir string, vs VerificationStatus) error {
	path := filepath.Join(changeDir, verificationStatusFile)
	data, err := json.MarshalIndent(vs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal verification-status: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write verification-status: %w", err)
	}
	return nil
}

// UpdateItemStatus reads the verification-status.json sidecar, updates a single
// requirement's status by reqID, and writes it back. Creates the sidecar if missing.
// Per D-04: use StatusDone for pass, StatusBlocked for fail on MUST items.
func UpdateItemStatus(changeDir string, reqID string, newStatus ItemStatus) error {
	vs, err := ReadVerificationStatus(changeDir)
	if err != nil {
		return fmt.Errorf("load verification status: %w", err)
	}

	vs.Requirements[reqID] = newStatus
	vs.VerifiedAt = time.Now().UTC()

	return WriteVerificationStatus(changeDir, vs)
}

// ParseTasksV2 reads a tasks.md file, parses its YAML frontmatter into TasksFrontmatterV2,
// and returns the remaining body string. Enables YAML round-trip task status tracking.
func ParseTasksV2(tasksPath string) (TasksFrontmatterV2, string, error) {
	f, err := os.Open(tasksPath)
	if err != nil {
		return TasksFrontmatterV2{}, "", fmt.Errorf("open tasks file: %w", err)
	}
	defer f.Close()

	var fm TasksFrontmatterV2
	rest, err := frontmatter.Parse(f, &fm)
	if err != nil {
		// No valid frontmatter — return zero-value fm and full file as body
		content, readErr := os.ReadFile(tasksPath)
		if readErr != nil {
			return TasksFrontmatterV2{}, "", readErr
		}
		return TasksFrontmatterV2{}, string(content), nil
	}

	return fm, string(rest), nil
}

// UpdateTaskStatus updates a single task's status by ID, recomputes the Completed count,
// and writes the updated frontmatter + body back to the file.
func UpdateTaskStatus(tasksPath string, taskID int, newStatus ItemStatus) error {
	fm, body, err := ParseTasksV2(tasksPath)
	if err != nil {
		return fmt.Errorf("parse tasks: %w", err)
	}

	found := false
	for i := range fm.Tasks {
		if fm.Tasks[i].ID == taskID {
			fm.Tasks[i].Status = newStatus
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("task %d not found", taskID)
	}

	// Recompute Completed count: number of tasks with StatusDone
	completed := 0
	for _, t := range fm.Tasks {
		if t.Status == StatusDone {
			completed++
		}
	}
	fm.Completed = completed

	return WriteTasks(tasksPath, fm, body)
}

// MaxTaskID returns the maximum task ID in the frontmatter, or 0 if empty.
func MaxTaskID(fm TasksFrontmatterV2) int {
	maxID := 0
	for _, t := range fm.Tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID
}

// MergeTasksForSpec merges new tasks into an existing TasksFrontmatterV2.
// New tasks have their IDs renumbered starting from maxExistingID + 1.
// If specName is non-empty, existing tasks with the same spec are replaced.
func MergeTasksForSpec(existing TasksFrontmatterV2, newTasks []TaskEntry, specName string) TasksFrontmatterV2 {
	// Keep tasks that don't belong to the target spec
	var kept []TaskEntry
	if specName != "" {
		for _, t := range existing.Tasks {
			if t.Spec != specName {
				kept = append(kept, t)
			}
		}
	} else {
		kept = existing.Tasks
	}

	// Determine starting ID
	startID := MaxTaskID(existing) + 1

	// Renumber and append new tasks
	for i, t := range newTasks {
		t.ID = startID + i
		if specName != "" && t.Spec == "" {
			t.Spec = specName
		}
		kept = append(kept, t)
	}

	result := existing
	result.Tasks = kept
	result.Total = len(kept)

	// Recompute completed
	completed := 0
	for _, t := range kept {
		if t.Status == StatusDone {
			completed++
		}
	}
	result.Completed = completed

	return result
}

// WriteTasks serializes TasksFrontmatterV2 as YAML frontmatter, prepends/appends the
// `---` delimiters, and appends the original body content. Preserves markdown body unchanged.
func WriteTasks(tasksPath string, fm TasksFrontmatterV2, body string) error {
	yamlBytes, err := yaml.Marshal(fm)
	if err != nil {
		return fmt.Errorf("marshal frontmatter: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(yamlBytes)
	sb.WriteString("---\n")
	sb.WriteString(body)

	return os.WriteFile(tasksPath, []byte(sb.String()), 0644)
}
