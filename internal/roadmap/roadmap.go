package roadmap

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mysd/internal/state"
	"gopkg.in/yaml.v3"
)

// UpdateTracking reads/upserts tracking.yaml in {projectRoot}/.mysd/roadmap/,
// then regenerates timeline.md. The function is best-effort: callers should
// log warnings on error but must not block state transitions.
//
// projectRoot is derived from specsDir by taking its parent directory,
// which handles both ".specs/" and "openspec/" conventions.
func UpdateTracking(specsDir string, ws state.WorkflowState) error {
	projectRoot := filepath.Dir(specsDir)
	roadmapDir := filepath.Join(projectRoot, ".mysd", "roadmap")

	if err := os.MkdirAll(roadmapDir, 0755); err != nil {
		return fmt.Errorf("create roadmap dir: %w", err)
	}

	trackingPath := filepath.Join(roadmapDir, "tracking.yaml")

	// Read existing tracking file; zero-value on missing file (convention over config).
	tf, err := ReadTracking(roadmapDir)
	if err != nil {
		return fmt.Errorf("read tracking: %w", err)
	}
	if tf.SchemaVersion == "" {
		tf.SchemaVersion = "1"
	}

	// Find existing record for this change name.
	now := time.Now()
	found := false
	for i, cr := range tf.Changes {
		if cr.Name == ws.ChangeName {
			tf.Changes[i].Status = string(ws.Phase)
			if ws.Phase == state.PhaseArchived && tf.Changes[i].CompletedAt == nil {
				completedAt := now
				tf.Changes[i].CompletedAt = &completedAt
			}
			found = true
			break
		}
	}

	if !found {
		startedAt := now
		rec := ChangeRecord{
			Name:      ws.ChangeName,
			Status:    string(ws.Phase),
			StartedAt: &startedAt,
		}
		if ws.Phase == state.PhaseArchived {
			completedAt := now
			rec.CompletedAt = &completedAt
		}
		tf.Changes = append(tf.Changes, rec)
	}

	tf.UpdatedAt = now

	// Marshal and write tracking.yaml.
	data, err := yaml.Marshal(tf)
	if err != nil {
		return fmt.Errorf("marshal tracking: %w", err)
	}
	if err := os.WriteFile(trackingPath, data, 0644); err != nil {
		return fmt.Errorf("write tracking.yaml: %w", err)
	}

	// Regenerate timeline.md.
	mermaid := GenerateMermaid(tf)
	timelinePath := filepath.Join(roadmapDir, "timeline.md")
	timelineContent := "```mermaid\n" + mermaid + "```\n"
	if err := os.WriteFile(timelinePath, []byte(timelineContent), 0644); err != nil {
		return fmt.Errorf("write timeline.md: %w", err)
	}

	return nil
}

// ReadTracking reads tracking.yaml from roadmapDir.
// If the file does not exist, it returns a zero-value TrackingFile without error.
// This follows the same convention as uat.ReadUAT.
func ReadTracking(roadmapDir string) (TrackingFile, error) {
	trackingPath := filepath.Join(roadmapDir, "tracking.yaml")
	data, err := os.ReadFile(trackingPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return TrackingFile{}, nil
		}
		return TrackingFile{}, fmt.Errorf("read tracking.yaml: %w", err)
	}

	var tf TrackingFile
	if err := yaml.Unmarshal(data, &tf); err != nil {
		return TrackingFile{}, fmt.Errorf("parse tracking.yaml: %w", err)
	}
	return tf, nil
}
