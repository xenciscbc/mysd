package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupChangeDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "test-change")
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0644))
}

func TestCoverage_MissingSpec(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, "proposal.md", "## Capabilities\n\n### New Capabilities\n\n- `foo-bar`: Foo bar feature\n")
	require.NoError(t, os.MkdirAll(filepath.Join(changeDir, "specs"), 0755))

	findings := CheckCoverage(changeDir)

	require.Len(t, findings, 1)
	assert.Equal(t, "Critical", findings[0].Severity)
	assert.Contains(t, findings[0].Summary, "foo-bar")
	assert.Equal(t, "COV-1", findings[0].ID)
}

func TestCoverage_SpecExists(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, "proposal.md", "## Capabilities\n\n### New Capabilities\n\n- `foo-bar`: Foo bar feature\n")
	writeFile(t, filepath.Join(changeDir, "specs", "foo-bar"), "spec.md", "### Requirement: Foo\n\nThe system SHALL foo.\n")

	findings := CheckCoverage(changeDir)
	assert.Empty(t, findings)
}

func TestAmbiguity_WeakLanguage(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, filepath.Join(changeDir, "specs", "test"), "spec.md",
		"### Requirement: Auth\n\nThe system should authenticate users.\n\n#### Scenario: Login\n\n- WHEN valid\n- THEN pass\n")

	findings := CheckAmbiguity(changeDir)

	require.Len(t, findings, 1)
	assert.Equal(t, "Suggestion", findings[0].Severity)
	assert.Contains(t, findings[0].Summary, "should")
}

func TestAmbiguity_CleanSpec(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, filepath.Join(changeDir, "specs", "test"), "spec.md",
		"### Requirement: Auth\n\nThe system SHALL authenticate users.\n\n#### Scenario: Login\n\n- WHEN valid\n- THEN pass\n")

	findings := CheckAmbiguity(changeDir)
	assert.Empty(t, findings)
}

func TestConsistency_DesignNotInTasks(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, "proposal.md", "## Capabilities\n\n### New Capabilities\n\n- `auth`: Auth\n")
	writeFile(t, changeDir, "design.md", "### D1: Authentication flow\n\nDesign details.\n")
	writeFile(t, changeDir, "tasks.md", "- [ ] 1.1 Something unrelated\n")

	findings := CheckConsistency(changeDir)

	require.Len(t, findings, 1)
	assert.Equal(t, "Warning", findings[0].Severity)
	assert.Contains(t, findings[0].Summary, "D1: Authentication flow")
}

func TestConsistency_DesignReferencedInTasks(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, "proposal.md", "## Capabilities\n\n### New Capabilities\n\n- `auth`: Auth\n")
	writeFile(t, changeDir, "design.md", "### D1: Authentication flow\n\nDesign details.\n")
	writeFile(t, changeDir, "tasks.md", "- [ ] 1.1 Implement D1: Authentication flow\n")

	findings := CheckConsistency(changeDir)
	assert.Empty(t, findings)
}

func TestGaps_RequirementWithoutScenario(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, filepath.Join(changeDir, "specs", "test"), "spec.md",
		"### Requirement: Missing scenario\n\nThe system SHALL do something.\n")

	findings := CheckGaps(changeDir)

	var gapFindings []Finding
	for _, f := range findings {
		if f.Summary == "Requirement 'Missing scenario' has no scenario" {
			gapFindings = append(gapFindings, f)
		}
	}
	require.NotEmpty(t, gapFindings)
	assert.Equal(t, "Warning", gapFindings[0].Severity)
}

func TestGaps_RequirementWithScenario(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, filepath.Join(changeDir, "specs", "test"), "spec.md",
		"### Requirement: Has scenario\n\nThe system SHALL do something.\n\n#### Scenario: It works\n\n- WHEN triggered\n- THEN success\n")
	writeFile(t, changeDir, "tasks.md", "- [ ] 1.1 Implement Has scenario\n")

	findings := CheckGaps(changeDir)
	assert.Empty(t, findings)
}

func TestAnalyze_AvailableArtifactsOnly(t *testing.T) {
	changeDir := setupChangeDir(t)
	writeFile(t, changeDir, "proposal.md", "## Summary\n\nTest\n")

	result := Analyze(changeDir)

	assert.Contains(t, result.ArtifactsAnalyzed, "proposal")
	assert.Contains(t, result.ArtifactsMissing, "specs")
	assert.Contains(t, result.ArtifactsMissing, "design")
	assert.Contains(t, result.ArtifactsMissing, "tasks")
	assert.Len(t, result.Dimensions, 4)
}
