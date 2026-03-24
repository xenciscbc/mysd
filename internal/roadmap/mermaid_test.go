package roadmap_test

import (
	"testing"
	"time"

	"github.com/mysd/internal/roadmap"
	"github.com/stretchr/testify/assert"
)

// TestGenerateMermaid_BasicChart verifies the gantt chart contains expected sections.
func TestGenerateMermaid_BasicChart(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	completed := now.Add(-1 * time.Hour)

	tf := roadmap.TrackingFile{
		SchemaVersion: "1",
		UpdatedAt:     now,
		Changes: []roadmap.ChangeRecord{
			{
				Name:        "auth",
				Status:      "archived",
				StartedAt:   &yesterday,
				CompletedAt: &completed,
			},
			{
				Name:      "payments",
				Status:    "proposed",
				StartedAt: &yesterday,
			},
		},
	}

	output := roadmap.GenerateMermaid(tf)

	assert.Contains(t, output, "gantt")
	assert.Contains(t, output, "dateFormat")
	assert.Contains(t, output, "auth")
	assert.Contains(t, output, "payments")
}

// TestGenerateMermaid_EmptyChanges verifies that an empty change list still produces a valid gantt header.
func TestGenerateMermaid_EmptyChanges(t *testing.T) {
	tf := roadmap.TrackingFile{
		SchemaVersion: "1",
		UpdatedAt:     time.Now(),
		Changes:       []roadmap.ChangeRecord{},
	}

	output := roadmap.GenerateMermaid(tf)

	assert.Contains(t, output, "gantt")
	assert.Contains(t, output, "dateFormat")
}
