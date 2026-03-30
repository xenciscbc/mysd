package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// requirementPattern matches ### Requirement: headings in spec files.
var requirementPattern = regexp.MustCompile(`^###\s+Requirement:\s+(.+)`)

// scenarioPattern matches #### Scenario: headings in spec files.
var scenarioPattern = regexp.MustCompile(`^####\s+Scenario:\s+(.+)`)

// CheckGaps checks for requirements without scenarios and tasks without spec references.
func CheckGaps(changeDir string) []Finding {
	specsDir := filepath.Join(changeDir, "specs")
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil
	}

	var findings []Finding
	counter := 1

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		specPath := filepath.Join(specsDir, entry.Name(), "spec.md")
		data, err := os.ReadFile(specPath)
		if err != nil {
			continue
		}

		relPath := fmt.Sprintf("specs/%s/spec.md", entry.Name())
		reqFindings := checkRequirementsHaveScenarios(string(data), relPath, &counter)
		findings = append(findings, reqFindings...)
	}

	// Check tasks reference spec requirements
	tasks := readFileOrEmpty(filepath.Join(changeDir, "tasks.md"))
	if tasks != "" {
		allReqs := collectAllRequirementNames(specsDir)
		tasksLower := strings.ToLower(tasks)
		for _, req := range allReqs {
			reqLower := strings.ToLower(req.name)
			if !strings.Contains(tasksLower, reqLower) {
				findings = append(findings, Finding{
					ID:             fmt.Sprintf("GAP-%d", counter),
					Dimension:      string(DimensionGaps),
					Severity:       string(SeverityWarning),
					Location:       req.location,
					Summary:        fmt.Sprintf("Requirement '%s' has no matching task", req.name),
					Recommendation: fmt.Sprintf("Add a task in tasks.md that references '%s'", req.name),
				})
				counter++
			}
		}
	}

	return findings
}

type reqRef struct {
	name     string
	location string
}

// collectAllRequirementNames gathers all requirement names from all spec files.
func collectAllRequirementNames(specsDir string) []reqRef {
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil
	}

	var refs []reqRef
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		specPath := filepath.Join(specsDir, entry.Name(), "spec.md")
		data, err := os.ReadFile(specPath)
		if err != nil {
			continue
		}

		relPath := fmt.Sprintf("specs/%s/spec.md", entry.Name())
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			matches := requirementPattern.FindStringSubmatch(line)
			if len(matches) >= 2 {
				refs = append(refs, reqRef{
					name:     strings.TrimSpace(matches[1]),
					location: relPath,
				})
			}
		}
	}
	return refs
}

// checkRequirementsHaveScenarios checks that each ### Requirement has at least one #### Scenario.
func checkRequirementsHaveScenarios(content, relPath string, counter *int) []Finding {
	lines := strings.Split(content, "\n")
	var findings []Finding
	var currentReq string
	var currentReqLine int
	hasScenario := false

	for i, line := range lines {
		reqMatch := requirementPattern.FindStringSubmatch(line)
		if len(reqMatch) >= 2 {
			// Check previous requirement
			if currentReq != "" && !hasScenario {
				findings = append(findings, Finding{
					ID:             fmt.Sprintf("GAP-%d", *counter),
					Dimension:      string(DimensionGaps),
					Severity:       string(SeverityWarning),
					Location:       fmt.Sprintf("%s:%d", relPath, currentReqLine+1),
					Summary:        fmt.Sprintf("Requirement '%s' has no scenario", currentReq),
					Recommendation: "Add at least one #### Scenario: under this requirement",
				})
				*counter++
			}
			currentReq = strings.TrimSpace(reqMatch[1])
			currentReqLine = i
			hasScenario = false
			continue
		}

		if scenarioPattern.MatchString(line) {
			hasScenario = true
		}
	}

	// Check last requirement
	if currentReq != "" && !hasScenario {
		findings = append(findings, Finding{
			ID:             fmt.Sprintf("GAP-%d", *counter),
			Dimension:      string(DimensionGaps),
			Severity:       string(SeverityWarning),
			Location:       fmt.Sprintf("%s:%d", relPath, currentReqLine+1),
			Summary:        fmt.Sprintf("Requirement '%s' has no scenario", currentReq),
			Recommendation: "Add at least one #### Scenario: under this requirement",
		})
		*counter++
	}

	return findings
}
