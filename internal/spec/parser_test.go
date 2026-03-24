package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// extractKeyword tests — RFC 2119 case-sensitive matching (Pitfall 2)
func TestExtractKeyword_Must(t *testing.T) {
	kw := extractKeyword("System MUST validate input")
	assert.Equal(t, Must, kw)
}

func TestExtractKeyword_Should(t *testing.T) {
	kw := extractKeyword("System SHOULD log events")
	assert.Equal(t, Should, kw)
}

func TestExtractKeyword_May(t *testing.T) {
	kw := extractKeyword("System MAY cache results")
	assert.Equal(t, May, kw)
}

func TestExtractKeyword_LowercaseMust_NotRFC2119(t *testing.T) {
	// Lowercase "must" must NOT match RFC 2119 — it's a common English word
	kw := extractKeyword("you must not do this")
	assert.Equal(t, RFC2119Keyword(""), kw, "lowercase 'must' should NOT match RFC 2119")
}

func TestExtractKeyword_NoKeyword(t *testing.T) {
	kw := extractKeyword("no keywords here")
	assert.Equal(t, RFC2119Keyword(""), kw)
}

func TestExtractKeyword_MustNot(t *testing.T) {
	kw := extractKeyword("The system MUST NOT allow unauthenticated access")
	assert.Equal(t, Must, kw)
}

func TestExtractKeyword_ShouldNot(t *testing.T) {
	kw := extractKeyword("The system SHOULD NOT store passwords in plaintext")
	assert.Equal(t, Should, kw)
}

// ParseProposal tests
func TestParseProposal_Brownfield(t *testing.T) {
	doc, err := ParseProposal("../../testdata/fixtures/openspec-project/openspec/changes/sample-change/proposal.md")
	require.NoError(t, err)
	assert.NotEmpty(t, doc.Body)
	assert.Contains(t, doc.Body, "Summary")
	// Brownfield: no frontmatter, so SpecVersion should be empty
	assert.Empty(t, doc.Frontmatter.SpecVersion)
}

func TestParseProposal_Native(t *testing.T) {
	doc, err := ParseProposal("../../testdata/fixtures/mysd-project/.specs/changes/add-dark-mode/proposal.md")
	require.NoError(t, err)
	assert.Equal(t, "1", doc.Frontmatter.SpecVersion)
	assert.Equal(t, "add-dark-mode", doc.Frontmatter.ChangeName)
	assert.Equal(t, "proposed", doc.Frontmatter.Status)
	assert.NotEmpty(t, doc.Body)
}

// ParseSpec tests
func TestParseSpec_BrownfieldRFC2119(t *testing.T) {
	reqs, err := ParseSpec("../../testdata/fixtures/openspec-project/openspec/changes/sample-change/specs/user-auth/spec.md")
	require.NoError(t, err)
	assert.NotEmpty(t, reqs)

	// Should find requirements with RFC 2119 keywords
	var mustReqs []Requirement
	var shouldReqs []Requirement
	for _, r := range reqs {
		if r.Keyword == Must {
			mustReqs = append(mustReqs, r)
		}
		if r.Keyword == Should {
			shouldReqs = append(shouldReqs, r)
		}
	}
	assert.NotEmpty(t, mustReqs, "Should find MUST requirements")
	assert.NotEmpty(t, shouldReqs, "Should find SHOULD requirements")
}

func TestParseSpec_NativeWithFrontmatter(t *testing.T) {
	reqs, err := ParseSpec("../../testdata/fixtures/mysd-project/.specs/changes/add-dark-mode/specs/theme-support/spec.md")
	require.NoError(t, err)
	assert.NotEmpty(t, reqs)
}

// ParseTasks tests
func TestParseTasks_Brownfield(t *testing.T) {
	tasks, fm, err := ParseTasks("../../testdata/fixtures/openspec-project/openspec/changes/sample-change/tasks.md")
	require.NoError(t, err)
	assert.NotEmpty(t, tasks)
	// Brownfield: no frontmatter
	assert.Empty(t, fm.SpecVersion)

	// Check that done tasks are detected
	var doneTasks []Task
	for _, t := range tasks {
		if t.Status == StatusDone {
			doneTasks = append(doneTasks, t)
		}
	}
	assert.NotEmpty(t, doneTasks, "Should detect [x] tasks as done")
}

func TestParseTasks_NativeWithFrontmatter(t *testing.T) {
	tasks, fm, err := ParseTasks("../../testdata/fixtures/mysd-project/.specs/changes/add-dark-mode/tasks.md")
	require.NoError(t, err)
	assert.Equal(t, "1", fm.SpecVersion)
	assert.Equal(t, 3, fm.Total)
	assert.Equal(t, 0, fm.Completed)
	assert.Len(t, tasks, 3)
}

// SourceFile field tests — ensure parser fills SourceFile after Task 1 changes
func TestParseSpec_FillsSourceFile(t *testing.T) {
	reqs, err := ParseSpec("../../testdata/fixtures/mysd-project/.specs/changes/add-dark-mode/specs/theme-support/spec.md")
	require.NoError(t, err)
	require.NotEmpty(t, reqs, "Expected at least one requirement")
	for _, r := range reqs {
		assert.Equal(t, "spec.md", r.SourceFile, "SourceFile must be filepath.Base(path)")
	}
}

func TestParseRequirementsFromBody_RegressionKeywords(t *testing.T) {
	body := "The system MUST validate input.\nThe system SHOULD log events.\nThe system MAY cache."
	reqs := parseRequirementsFromBody(body, DeltaNone)
	require.Len(t, reqs, 3)
	assert.Equal(t, Must, reqs[0].Keyword)
	assert.Equal(t, Should, reqs[1].Keyword)
	assert.Equal(t, May, reqs[2].Keyword)
}

// ParseChangeMeta tests
func TestParseChangeMeta_MySD(t *testing.T) {
	meta, err := ParseChangeMeta("../../testdata/fixtures/mysd-project/.specs/changes/add-dark-mode/.openspec.yaml")
	require.NoError(t, err)
	assert.Equal(t, "spec-driven", meta.Schema)
	assert.Equal(t, "2026-03-23", meta.Created)
}

// ParseChange integration test
func TestParseChange_Brownfield(t *testing.T) {
	c, err := ParseChange("../../testdata/fixtures/openspec-project/openspec/changes/sample-change")
	require.NoError(t, err)
	assert.Equal(t, "sample-change", c.Name)
	assert.NotEmpty(t, c.Specs)
	assert.Equal(t, "spec-driven", c.Meta.Schema)
}

func TestParseChange_Native(t *testing.T) {
	c, err := ParseChange("../../testdata/fixtures/mysd-project/.specs/changes/add-dark-mode")
	require.NoError(t, err)
	assert.Equal(t, "add-dark-mode", c.Name)
	assert.Equal(t, "1", c.Proposal.Frontmatter.SpecVersion)
	assert.NotEmpty(t, c.Specs)
	assert.Equal(t, "spec-driven", c.Meta.Schema)
}
