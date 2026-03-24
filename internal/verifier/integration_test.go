package verifier_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/verifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestChange creates a temp change directory with a fixture spec containing
// 3 MUST items, 2 SHOULD items, and 1 MAY item. Returns specsDir and changeName.
func setupTestChange(t *testing.T) (specsDir, changeName string) {
	t.Helper()
	tmp := t.TempDir()
	changeName = "integration-test-change"
	specsDir = tmp

	changeDir := filepath.Join(specsDir, "changes", changeName)
	specsSubDir := filepath.Join(changeDir, "specs", "capability-a")
	require.NoError(t, os.MkdirAll(specsSubDir, 0755))

	// .openspec.yaml
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, ".openspec.yaml"),
		[]byte("schema: \"1\"\n"),
		0644,
	))

	// proposal.md
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte("---\ntitle: Integration test proposal\n---\n\n# Proposal\nTest change.\n"),
		0644,
	))

	// specs/capability-a/spec.md — 3 MUST, 2 SHOULD, 1 MAY
	specContent := `# Specification: Capability A

## Requirement: Authentication

The system MUST validate user credentials before granting access.
The system MUST reject invalid credentials with a 401 response.
The system MUST store passwords using bcrypt hashing.

## Requirement: Logging

The system SHOULD log all authentication attempts with timestamps.
The system SHOULD notify administrators of repeated failures.

## Requirement: Extras

The system MAY support single sign-on via OAuth2.
`
	require.NoError(t, os.WriteFile(
		filepath.Join(specsSubDir, "spec.md"),
		[]byte(specContent),
		0644,
	))

	// design.md
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "design.md"),
		[]byte("# Design\nAuthentication system design.\n"),
		0644,
	))

	// tasks.md — V2 frontmatter + 2 done tasks
	tasksContent := `---
schema: "v2"
total: 2
completed: 2
tasks:
  - id: 1
    name: Implement auth handler
    status: done
  - id: 2
    name: Write tests
    status: done
---

- [x] Implement auth handler
- [x] Write tests
`
	require.NoError(t, os.WriteFile(
		filepath.Join(changeDir, "tasks.md"),
		[]byte(tasksContent),
		0644,
	))

	return specsDir, changeName
}

// buildPassingReport constructs a VerifierReport where all items pass.
func buildPassingReport(ctx verifier.VerificationContext) verifier.VerifierReport {
	report := verifier.VerifierReport{
		ChangeName:  ctx.ChangeName,
		OverallPass: true,
		MustPass:    true,
		Results:     []verifier.VerifierResultItem{},
		UIItems:     []verifier.UIItem{},
	}
	for _, item := range ctx.MustItems {
		report.Results = append(report.Results, verifier.VerifierResultItem{
			ID:       item.ID,
			Text:     item.Text,
			Keyword:  "MUST",
			Pass:     true,
			Evidence: "Verified in code",
		})
	}
	for _, item := range ctx.ShouldItems {
		report.Results = append(report.Results, verifier.VerifierResultItem{
			ID:       item.ID,
			Text:     item.Text,
			Keyword:  "SHOULD",
			Pass:     true,
			Evidence: "Logging implemented",
		})
	}
	for _, item := range ctx.MayItems {
		report.Results = append(report.Results, verifier.VerifierResultItem{
			ID:      item.ID,
			Text:    item.Text,
			Keyword: "MAY",
			Pass:    false, // MAY items can be skipped
		})
	}
	return report
}

// buildFailingReport constructs a VerifierReport where the first MUST item fails.
func buildFailingReport(ctx verifier.VerificationContext) verifier.VerifierReport {
	report := buildPassingReport(ctx)
	report.OverallPass = false
	report.MustPass = false
	// Fail the first MUST item
	for i := range report.Results {
		if report.Results[i].Keyword == "MUST" {
			report.Results[i].Pass = false
			report.Results[i].Evidence = "No authentication handler found"
			report.Results[i].Suggestion = "Implement validateCredentials function"
			break
		}
	}
	return report
}

// TestVerificationPipeline_AllPass tests the full pipeline when all MUST items pass.
func TestVerificationPipeline_AllPass(t *testing.T) {
	specsDir, changeName := setupTestChange(t)
	changeDir := filepath.Join(specsDir, "changes", changeName)

	// Step 1: Build verification context
	ctx, err := verifier.BuildVerificationContext(specsDir, changeName)
	require.NoError(t, err)

	// Assert item counts: 3 MUST, 2 SHOULD, 1 MAY
	assert.Len(t, ctx.MustItems, 3, "should have 3 MUST items")
	assert.Len(t, ctx.ShouldItems, 2, "should have 2 SHOULD items")
	assert.Len(t, ctx.MayItems, 1, "should have 1 MAY item")

	// Step 2: Construct a passing VerifierReport
	report := buildPassingReport(ctx)

	// Step 3: WriteVerificationReport
	err = verifier.WriteVerificationReport(changeDir, report)
	require.NoError(t, err)

	verPath := filepath.Join(changeDir, "verification.md")
	assert.FileExists(t, verPath, "verification.md must be created")
	content, err := os.ReadFile(verPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "## MUST Items", "verification.md must contain MUST Items section")

	// Step 4: WriteVerificationStatus — all MUST = done
	vs := spec.VerificationStatus{
		ChangeName:   changeName,
		VerifiedAt:   time.Now().UTC(),
		Requirements: map[string]spec.ItemStatus{},
	}
	for _, item := range ctx.MustItems {
		vs.Requirements[item.ID] = spec.StatusDone
	}
	err = spec.WriteVerificationStatus(changeDir, vs)
	require.NoError(t, err)

	// Step 5: ReadVerificationStatus — assert all 3 MUST = done
	loaded, err := spec.ReadVerificationStatus(changeDir)
	require.NoError(t, err)
	assert.Len(t, loaded.Requirements, 3, "should have 3 requirements tracked")
	for _, item := range ctx.MustItems {
		status, ok := loaded.Requirements[item.ID]
		assert.True(t, ok, "MUST item %s should be tracked", item.ID)
		assert.Equal(t, spec.StatusDone, status, "MUST item %s should be done", item.ID)
	}
}

// TestVerificationPipeline_MustFailure tests the pipeline when a MUST item fails.
func TestVerificationPipeline_MustFailure(t *testing.T) {
	specsDir, changeName := setupTestChange(t)
	changeDir := filepath.Join(specsDir, "changes", changeName)

	// Build context
	ctx, err := verifier.BuildVerificationContext(specsDir, changeName)
	require.NoError(t, err)
	require.Len(t, ctx.MustItems, 3)

	// Build a report with first MUST failing
	report := buildFailingReport(ctx)

	// WriteGapReport — assert gap-report.md exists and contains failed item text
	err = verifier.WriteGapReport(changeDir, report)
	require.NoError(t, err)

	gapPath := filepath.Join(changeDir, "gap-report.md")
	assert.FileExists(t, gapPath, "gap-report.md must be created when MUST items fail")

	gapContent, err := os.ReadFile(gapPath)
	require.NoError(t, err)
	gapStr := string(gapContent)

	// Gap report should contain the failed item text
	assert.Contains(t, gapStr, "MUST validate user credentials", "gap report should contain failing MUST item text")

	// Gap report should NOT contain passing MUST items
	assert.NotContains(t, gapStr, "MUST reject invalid credentials", "passing MUST items must not appear in gap report")
	assert.NotContains(t, gapStr, "MUST store passwords", "passing MUST items must not appear in gap report")

	// WriteVerificationStatus with mix of done/blocked
	vs := spec.VerificationStatus{
		ChangeName:   changeName,
		VerifiedAt:   time.Now().UTC(),
		Requirements: map[string]spec.ItemStatus{},
	}
	firstMust := true
	for _, item := range ctx.MustItems {
		if firstMust {
			vs.Requirements[item.ID] = spec.StatusBlocked
			firstMust = false
		} else {
			vs.Requirements[item.ID] = spec.StatusDone
		}
	}
	err = spec.WriteVerificationStatus(changeDir, vs)
	require.NoError(t, err)

	// ReadVerificationStatus — assert 1 blocked
	loaded, err := spec.ReadVerificationStatus(changeDir)
	require.NoError(t, err)

	blockedCount := 0
	for _, status := range loaded.Requirements {
		if status == spec.StatusBlocked {
			blockedCount++
		}
	}
	assert.Equal(t, 1, blockedCount, "should have exactly 1 blocked MUST item")
}

// TestStableID_Consistency tests that BuildVerificationContext produces consistent IDs.
func TestStableID_Consistency(t *testing.T) {
	specsDir, changeName := setupTestChange(t)

	// Call BuildVerificationContext twice
	ctx1, err := verifier.BuildVerificationContext(specsDir, changeName)
	require.NoError(t, err)

	ctx2, err := verifier.BuildVerificationContext(specsDir, changeName)
	require.NoError(t, err)

	// IDs must be identical across both calls
	require.Equal(t, len(ctx1.MustItems), len(ctx2.MustItems), "both calls must produce same number of MUST items")
	for i := range ctx1.MustItems {
		assert.Equal(t, ctx1.MustItems[i].ID, ctx2.MustItems[i].ID,
			"MUST item %d ID must be identical across two BuildVerificationContext calls", i)
	}
}

// TestVerificationReport_Ordering tests that sections appear in MUST -> SHOULD -> MAY order.
func TestVerificationReport_Ordering(t *testing.T) {
	specsDir, changeName := setupTestChange(t)
	changeDir := filepath.Join(specsDir, "changes", changeName)

	ctx, err := verifier.BuildVerificationContext(specsDir, changeName)
	require.NoError(t, err)

	// Build a mixed report
	report := buildPassingReport(ctx)

	err = verifier.WriteVerificationReport(changeDir, report)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(changeDir, "verification.md"))
	require.NoError(t, err)
	contentStr := string(content)

	mustIdx := strings.Index(contentStr, "## MUST Items")
	shouldIdx := strings.Index(contentStr, "## SHOULD Items")
	mayIdx := strings.Index(contentStr, "## MAY Items")

	assert.GreaterOrEqual(t, mustIdx, 0, "## MUST Items section must exist")
	assert.GreaterOrEqual(t, shouldIdx, 0, "## SHOULD Items section must exist")
	assert.GreaterOrEqual(t, mayIdx, 0, "## MAY Items section must exist")

	assert.Less(t, mustIdx, shouldIdx, "## MUST Items must appear before ## SHOULD Items")
	assert.Less(t, shouldIdx, mayIdx, "## SHOULD Items must appear before ## MAY Items")
}
