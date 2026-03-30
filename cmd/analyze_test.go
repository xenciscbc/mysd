package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/xenciscbc/mysd/internal/analyzer"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAnalyzeTestChange(t *testing.T) (specsDir, changeName string) {
	t.Helper()
	tmp := t.TempDir()
	specsDir = tmp
	changeName = "test-analyze"
	changeDir := filepath.Join(specsDir, "changes", changeName)

	require.NoError(t, os.MkdirAll(changeDir, 0755))

	// proposal with one capability
	proposal := "## Capabilities\n\n### New Capabilities\n\n- `auth`: Authentication\n\n## Impact\n\n- Affected code: cmd/auth.go\n"
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposal), 0644))

	// spec for auth
	specDir := filepath.Join(changeDir, "specs", "auth")
	require.NoError(t, os.MkdirAll(specDir, 0755))
	specContent := "### Requirement: User authentication\n\nThe system SHALL authenticate users.\n\n#### Scenario: Valid login\n\n- **WHEN** valid credentials\n- **THEN** access granted\n"
	require.NoError(t, os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0644))

	// tasks
	tasks := "- [ ] 1.1 Implement User authentication — D1: auth design\n"
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasks), 0644))

	// design
	design := "### D1: auth design\n\nDesign details.\n"
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "design.md"), []byte(design), 0644))

	// state
	ws := state.WorkflowState{ChangeName: changeName, Phase: state.PhaseProposed}
	require.NoError(t, state.SaveState(specsDir, ws))

	return specsDir, changeName
}

func TestAnalyzeJSON_CleanChange(t *testing.T) {
	specsDir, changeName := setupAnalyzeTestChange(t)

	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(filepath.Dir(specsDir)))
	defer os.Chdir(origDir)

	// Symlink or adjust: runAnalyze uses spec.DetectSpecDir(".")
	// Instead, call analyzer directly
	changeDir := filepath.Join(specsDir, "changes", changeName)
	result := analyzer.Analyze(changeDir)

	data, err := json.MarshalIndent(result, "", "  ")
	require.NoError(t, err)

	// Verify JSON structure
	var parsed analyzer.AnalysisResult
	require.NoError(t, json.Unmarshal(data, &parsed))

	assert.Equal(t, changeName, parsed.ChangeID)
	assert.Len(t, parsed.Dimensions, 4)
	assert.Contains(t, parsed.ArtifactsAnalyzed, "proposal")
	assert.Contains(t, parsed.ArtifactsAnalyzed, "specs")
}

func TestAnalyzeJSON_MissingSpec(t *testing.T) {
	tmp := t.TempDir()
	changeDir := filepath.Join(tmp, "changes", "missing-spec")
	require.NoError(t, os.MkdirAll(changeDir, 0755))

	// proposal references a capability with no spec
	proposal := "## Capabilities\n\n### New Capabilities\n\n- `missing-cap`: Something\n\n## Impact\n\n- None\n"
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposal), 0644))

	// Create empty specs dir (no spec file)
	require.NoError(t, os.MkdirAll(filepath.Join(changeDir, "specs"), 0755))

	result := analyzer.Analyze(changeDir)

	// Should have a Coverage Critical finding
	var coverageFindings []analyzer.Finding
	for _, f := range result.Findings {
		if f.Dimension == "Coverage" {
			coverageFindings = append(coverageFindings, f)
		}
	}
	require.Len(t, coverageFindings, 1)
	assert.Equal(t, "Critical", coverageFindings[0].Severity)
	assert.Contains(t, coverageFindings[0].Summary, "missing-cap")
}

func TestAnalyzeStyledOutput(t *testing.T) {
	tmp := t.TempDir()
	changeDir := filepath.Join(tmp, "changes", "styled-test")
	require.NoError(t, os.MkdirAll(changeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("## Summary\n\nTest\n"), 0644))

	result := analyzer.Analyze(changeDir)

	var buf bytes.Buffer
	printAnalyzeSummary(&buf, result)

	output := buf.String()
	assert.Contains(t, output, "styled-test")
	assert.Contains(t, output, "Coverage")
}
