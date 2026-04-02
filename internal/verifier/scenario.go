package verifier

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

// reScenarioHeading matches #### Scenario: <name> headings.
var reScenarioHeading = regexp.MustCompile(`^####\s+Scenario:\s*(.+)$`)

// ValidateScenarioFormat checks that each #### Scenario: block in specBody
// contains **GIVEN**, **WHEN**, and **THEN** keywords.
// Returns a list of warnings for scenarios missing any keyword.
func ValidateScenarioFormat(specBody string) []string {
	var warnings []string
	scanner := bufio.NewScanner(strings.NewReader(specBody))

	var currentScenario string
	var hasGiven, hasWhen, hasThen bool

	flush := func() {
		if currentScenario == "" {
			return
		}
		var missing []string
		if !hasGiven {
			missing = append(missing, "GIVEN")
		}
		if !hasWhen {
			missing = append(missing, "WHEN")
		}
		if !hasThen {
			missing = append(missing, "THEN")
		}
		if len(missing) > 0 {
			warnings = append(warnings, fmt.Sprintf(
				"Scenario %q is missing: %s",
				currentScenario, strings.Join(missing, ", "),
			))
		}
	}

	for scanner.Scan() {
		line := scanner.Text()

		if m := reScenarioHeading.FindStringSubmatch(line); m != nil {
			flush()
			currentScenario = strings.TrimSpace(m[1])
			hasGiven = false
			hasWhen = false
			hasThen = false
			continue
		}

		if currentScenario != "" {
			if strings.Contains(line, "**GIVEN**") {
				hasGiven = true
			}
			if strings.Contains(line, "**WHEN**") {
				hasWhen = true
			}
			if strings.Contains(line, "**THEN**") {
				hasThen = true
			}
		}
	}

	flush()
	return warnings
}
