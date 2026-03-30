package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// capabilityPattern matches capability entries like: - `capability-name`: description
var capabilityPattern = regexp.MustCompile("^\\s*-\\s*`([a-z][a-z0-9-]*)`")

// CheckCoverage verifies that every capability listed in proposal has a corresponding spec file.
func CheckCoverage(changeDir string) []Finding {
	proposalPath := filepath.Join(changeDir, "proposal.md")
	data, err := os.ReadFile(proposalPath)
	if err != nil {
		return nil
	}

	capabilities := parseCapabilities(string(data))
	if len(capabilities) == 0 {
		return nil
	}

	specsDir := filepath.Join(changeDir, "specs")
	var findings []Finding
	counter := 1

	for _, cap := range capabilities {
		specPath := filepath.Join(specsDir, cap, "spec.md")
		if _, err := os.Stat(specPath); os.IsNotExist(err) {
			findings = append(findings, Finding{
				ID:             fmt.Sprintf("COV-%d", counter),
				Dimension:      string(DimensionCoverage),
				Severity:       string(SeverityCritical),
				Location:       "proposal.md",
				Summary:        fmt.Sprintf("Capability '%s' listed in proposal has no corresponding specs/%s/spec.md", cap, cap),
				Recommendation: fmt.Sprintf("Create specs/%s/spec.md or remove '%s' from proposal Capabilities", cap, cap),
			})
			counter++
		}
	}

	return findings
}

// parseCapabilities extracts capability names from the Capabilities section of a proposal.
func parseCapabilities(content string) []string {
	lines := strings.Split(content, "\n")
	inCapabilities := false
	var caps []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "## Capabilities") || strings.HasPrefix(trimmed, "### New Capabilities") || strings.HasPrefix(trimmed, "### Modified Capabilities") {
			inCapabilities = true
			continue
		}

		// Exit capabilities section on next ## that isn't a sub-heading of Capabilities
		if inCapabilities && strings.HasPrefix(trimmed, "## ") && !strings.HasPrefix(trimmed, "### ") {
			inCapabilities = false
			continue
		}

		if inCapabilities {
			matches := capabilityPattern.FindStringSubmatch(line)
			if len(matches) >= 2 {
				caps = append(caps, matches[1])
			}
		}
	}

	return caps
}
