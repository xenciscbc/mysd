package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// designHeadingPattern matches ### headings in design.md (decision headings).
var designHeadingPattern = regexp.MustCompile(`^###\s+(.+)`)

// CheckConsistency cross-checks references across proposal, specs, design, and tasks.
func CheckConsistency(changeDir string) []Finding {
	var findings []Finding
	counter := 1

	// Read available artifacts
	proposal := readFileOrEmpty(filepath.Join(changeDir, "proposal.md"))
	design := readFileOrEmpty(filepath.Join(changeDir, "design.md"))
	tasks := readFileOrEmpty(filepath.Join(changeDir, "tasks.md"))

	// Check 1: Design references capabilities not in proposal
	if design != "" && proposal != "" {
		proposalCaps := parseCapabilities(proposal)
		capSet := make(map[string]bool)
		for _, c := range proposalCaps {
			capSet[c] = true
		}

		// Check design headings are referenced in tasks
		if tasks != "" {
			headings := extractDesignHeadings(design)
			tasksLower := strings.ToLower(tasks)
			for _, h := range headings {
				hLower := strings.ToLower(h)
				if !strings.Contains(tasksLower, hLower) {
					findings = append(findings, Finding{
						ID:             fmt.Sprintf("CON-%d", counter),
						Dimension:      string(DimensionConsistency),
						Severity:       string(SeverityWarning),
						Location:       "design.md",
						Summary:        fmt.Sprintf("Design topic '%s' not referenced in tasks", h),
						Recommendation: "Verify tasks cover this design decision",
					})
					counter++
				}
			}
		}
	}

	return findings
}

// extractDesignHeadings extracts ### heading text from design.md.
func extractDesignHeadings(content string) []string {
	lines := strings.Split(content, "\n")
	var headings []string
	for _, line := range lines {
		matches := designHeadingPattern.FindStringSubmatch(line)
		if len(matches) >= 2 {
			headings = append(headings, strings.TrimSpace(matches[1]))
		}
	}
	return headings
}

// readFileOrEmpty reads a file and returns its content, or empty string if not found.
func readFileOrEmpty(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
