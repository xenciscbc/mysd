package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAlignmentPath_ReturnsCorrectPath verifies the path resolves per D-07.
func TestAlignmentPath_ReturnsCorrectPath(t *testing.T) {
	path := AlignmentPath("my-feature", ".specs")
	assert.Equal(t, ".specs/changes/my-feature/alignment.md", path)
}

// TestAlignmentTemplate_NonEmptyAndContainsHeading verifies the template is valid.
func TestAlignmentTemplate_NonEmptyAndContainsHeading(t *testing.T) {
	tmpl := AlignmentTemplate()

	assert.NotEmpty(t, tmpl)
	assert.Contains(t, tmpl, "## Alignment Summary")
	assert.Contains(t, tmpl, "### MUST items I understand:")
	assert.Contains(t, tmpl, "### Execution strategy:")
	assert.Contains(t, tmpl, "### Open questions (if any):")
}
