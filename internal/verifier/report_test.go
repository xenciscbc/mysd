package verifier_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mysd/internal/verifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sampleVerifierReport is a valid JSON verifier report for testing.
const sampleVerifierReport = `{
  "change_name": "add-dark-mode",
  "overall_pass": false,
  "must_pass": false,
  "results": [
    {
      "id": "spec.md::must-aabbccdd",
      "text": "System MUST support dark mode",
      "keyword": "MUST",
      "pass": false,
      "evidence": "No theme toggle found in UI",
      "suggestion": "Add ThemeToggle component to header"
    },
    {
      "id": "spec.md::should-11223344",
      "text": "System SHOULD remember user preference",
      "keyword": "SHOULD",
      "pass": true,
      "evidence": "localStorage saves theme preference",
      "suggestion": ""
    }
  ],
  "has_ui_items": false,
  "ui_items": []
}`

const samplePassingReport = `{
  "change_name": "add-dark-mode",
  "overall_pass": true,
  "must_pass": true,
  "results": [
    {
      "id": "spec.md::must-aabbccdd",
      "text": "System MUST support dark mode",
      "keyword": "MUST",
      "pass": true,
      "evidence": "Theme toggle implemented in header",
      "suggestion": ""
    }
  ],
  "has_ui_items": false,
  "ui_items": []
}`

// TestParseVerifierReport verifies successful JSON deserialization.
func TestParseVerifierReport(t *testing.T) {
	report, err := verifier.ParseVerifierReport([]byte(sampleVerifierReport))
	require.NoError(t, err)

	assert.Equal(t, "add-dark-mode", report.ChangeName)
	assert.False(t, report.OverallPass)
	assert.False(t, report.MustPass)
	require.Len(t, report.Results, 2)
	assert.Equal(t, "spec.md::must-aabbccdd", report.Results[0].ID)
	assert.Equal(t, "MUST", report.Results[0].Keyword)
	assert.False(t, report.Results[0].Pass)
	assert.Equal(t, "No theme toggle found in UI", report.Results[0].Evidence)
	assert.Equal(t, "Add ThemeToggle component to header", report.Results[0].Suggestion)
}

// TestParseVerifierReport_InvalidJSON verifies error on malformed JSON.
func TestParseVerifierReport_InvalidJSON(t *testing.T) {
	_, err := verifier.ParseVerifierReport([]byte(`{invalid json`))
	require.Error(t, err)
}

// TestWriteGapReport verifies gap-report.md is created with expected content.
func TestWriteGapReport(t *testing.T) {
	dir := t.TempDir()
	report, err := verifier.ParseVerifierReport([]byte(sampleVerifierReport))
	require.NoError(t, err)

	err = verifier.WriteGapReport(dir, report)
	require.NoError(t, err)

	gapPath := filepath.Join(dir, "gap-report.md")
	assert.FileExists(t, gapPath)

	content, err := os.ReadFile(gapPath)
	require.NoError(t, err)

	// YAML frontmatter must contain failed_must_ids
	assert.Contains(t, string(content), "failed_must_ids:")
	assert.Contains(t, string(content), "spec.md::must-aabbccdd")

	// Body must contain the failed MUST item
	assert.Contains(t, string(content), "System MUST support dark mode")
	assert.Contains(t, string(content), "No theme toggle found in UI")
	assert.Contains(t, string(content), "Add ThemeToggle component to header")

	// SHOULD pass item must NOT appear in gap report (only MUST failures)
	assert.NotContains(t, string(content), "System SHOULD remember user preference")
}

// TestWriteGapReport_NoFailures verifies no file is created when all MUST items pass.
func TestWriteGapReport_NoFailures(t *testing.T) {
	dir := t.TempDir()
	report, err := verifier.ParseVerifierReport([]byte(samplePassingReport))
	require.NoError(t, err)

	err = verifier.WriteGapReport(dir, report)
	require.NoError(t, err)

	// No gap-report.md when there are no failures
	gapPath := filepath.Join(dir, "gap-report.md")
	assert.NoFileExists(t, gapPath)
}

// TestWriteVerificationReport verifies verification.md structure.
func TestWriteVerificationReport(t *testing.T) {
	dir := t.TempDir()
	report, err := verifier.ParseVerifierReport([]byte(sampleVerifierReport))
	require.NoError(t, err)

	err = verifier.WriteVerificationReport(dir, report)
	require.NoError(t, err)

	verPath := filepath.Join(dir, "verification.md")
	assert.FileExists(t, verPath)

	content, err := os.ReadFile(verPath)
	require.NoError(t, err)

	// YAML frontmatter fields
	assert.Contains(t, string(content), "overall_pass:")
	assert.Contains(t, string(content), "must_pass:")
	assert.Contains(t, string(content), "must_total:")
	assert.Contains(t, string(content), "verified_at:")

	// Section ordering: MUST before SHOULD before MAY
	mustIdx := indexOf(string(content), "## MUST Items")
	shouldIdx := indexOf(string(content), "## SHOULD Items")
	assert.True(t, mustIdx < shouldIdx, "## MUST Items must appear before ## SHOULD Items")

	// Both items should appear
	assert.Contains(t, string(content), "System MUST support dark mode")
	assert.Contains(t, string(content), "System SHOULD remember user preference")
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
