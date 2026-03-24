package verifier

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// VerifierReport is the structured output produced by the verifier agent.
// It contains overall pass/fail status and per-requirement results.
type VerifierReport struct {
	ChangeName  string               `json:"change_name"`
	OverallPass bool                 `json:"overall_pass"`
	MustPass    bool                 `json:"must_pass"`
	Results     []VerifierResultItem `json:"results"`
	HasUIItems  bool                 `json:"has_ui_items"`
	UIItems     []UIItem             `json:"ui_items"`
}

// VerifierResultItem represents the verifier agent's assessment of a single requirement.
type VerifierResultItem struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Keyword    string `json:"keyword"`
	Pass       bool   `json:"pass"`
	Evidence   string `json:"evidence"`
	Suggestion string `json:"suggestion"`
}

// UIItem represents a UI verification item that requires visual inspection.
type UIItem struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	TestSteps string `json:"test_steps"`
}

// ParseVerifierReport deserializes JSON from the verifier agent into a VerifierReport.
// Returns an error if the JSON is invalid or Results is nil.
func ParseVerifierReport(data []byte) (VerifierReport, error) {
	var report VerifierReport
	if err := json.Unmarshal(data, &report); err != nil {
		return VerifierReport{}, fmt.Errorf("unmarshal verifier report: %w", err)
	}
	// Ensure Results is not nil (zero-length slice is acceptable)
	if report.Results == nil {
		report.Results = []VerifierResultItem{}
	}
	if report.UIItems == nil {
		report.UIItems = []UIItem{}
	}
	return report, nil
}

// WriteGapReport writes a gap-report.md to changeDir for all failed MUST items.
// If there are no MUST failures, no file is written (returns nil).
// The report contains YAML frontmatter with failed_must_ids and a markdown body
// describing each failure with evidence and suggested fix.
func WriteGapReport(changeDir string, report VerifierReport) error {
	// Collect failed MUST items
	var failedMust []VerifierResultItem
	for _, r := range report.Results {
		if r.Keyword == "MUST" && !r.Pass {
			failedMust = append(failedMust, r)
		}
	}

	// No failures — skip writing
	if len(failedMust) == 0 {
		return nil
	}

	var sb strings.Builder

	// YAML frontmatter
	sb.WriteString("---\n")
	sb.WriteString("failed_task_ids: []\n") // task ID mapping is done at CLI layer
	sb.WriteString("failed_must_ids:\n")
	for _, item := range failedMust {
		sb.WriteString(fmt.Sprintf("  - %s\n", item.ID))
	}
	sb.WriteString("---\n\n")

	// Markdown body — one section per failed MUST item
	sb.WriteString("# Gap Report\n\n")
	sb.WriteString("The following MUST requirements were not satisfied:\n\n")
	for _, item := range failedMust {
		sb.WriteString(fmt.Sprintf("### %s\n\n", item.Text))
		sb.WriteString(fmt.Sprintf("**Evidence:** %s\n\n", item.Evidence))
		sb.WriteString(fmt.Sprintf("**Suggested Fix:** %s\n\n", item.Suggestion))
	}

	gapPath := filepath.Join(changeDir, "gap-report.md")
	if err := os.WriteFile(gapPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("write gap-report.md: %w", err)
	}

	return nil
}

// WriteVerificationReport writes a verification.md to changeDir summarizing
// the full verifier report. Sections are ordered: MUST, SHOULD, MAY.
// Per D-01 and D-02 design decisions.
func WriteVerificationReport(changeDir string, report VerifierReport) error {
	now := time.Now().UTC()

	// Compute counts by keyword
	var mustTotal, mustPassed, shouldTotal, shouldPassed, mayTotal int
	for _, r := range report.Results {
		switch r.Keyword {
		case "MUST":
			mustTotal++
			if r.Pass {
				mustPassed++
			}
		case "SHOULD":
			shouldTotal++
			if r.Pass {
				shouldPassed++
			}
		case "MAY":
			mayTotal++
		}
	}

	var sb strings.Builder

	// YAML frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("verified_at: %s\n", now.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("overall_pass: %v\n", report.OverallPass))
	sb.WriteString(fmt.Sprintf("must_pass: %v\n", report.MustPass))
	sb.WriteString(fmt.Sprintf("must_total: %d\n", mustTotal))
	sb.WriteString(fmt.Sprintf("must_passed: %d\n", mustPassed))
	sb.WriteString(fmt.Sprintf("should_total: %d\n", shouldTotal))
	sb.WriteString(fmt.Sprintf("should_passed: %d\n", shouldPassed))
	sb.WriteString(fmt.Sprintf("may_total: %d\n", mayTotal))
	sb.WriteString("---\n\n")

	// Markdown body
	sb.WriteString(fmt.Sprintf("# Verification Report: %s\n\n", report.ChangeName))

	// MUST Items section
	sb.WriteString("## MUST Items\n\n")
	for _, r := range report.Results {
		if r.Keyword != "MUST" {
			continue
		}
		badge := "PASS"
		if !r.Pass {
			badge = "FAIL"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", badge, r.Text))
		if r.Evidence != "" {
			sb.WriteString(fmt.Sprintf("  - Evidence: %s\n", r.Evidence))
		}
		if !r.Pass && r.Suggestion != "" {
			sb.WriteString(fmt.Sprintf("  - Suggestion: %s\n", r.Suggestion))
		}
	}
	sb.WriteString("\n")

	// SHOULD Items section
	sb.WriteString("## SHOULD Items\n\n")
	for _, r := range report.Results {
		if r.Keyword != "SHOULD" {
			continue
		}
		badge := "PASS"
		if !r.Pass {
			badge = "FAIL"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", badge, r.Text))
		if r.Evidence != "" {
			sb.WriteString(fmt.Sprintf("  - Evidence: %s\n", r.Evidence))
		}
	}
	sb.WriteString("\n")

	// MAY Items section
	sb.WriteString("## MAY Items\n\n")
	for _, r := range report.Results {
		if r.Keyword != "MAY" {
			continue
		}
		badge := "PASS"
		if !r.Pass {
			badge = "SKIP"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", badge, r.Text))
	}

	verPath := filepath.Join(changeDir, "verification.md")
	if err := os.WriteFile(verPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("write verification.md: %w", err)
	}

	return nil
}
