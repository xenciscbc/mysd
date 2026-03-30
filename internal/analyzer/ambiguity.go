package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// weakPatterns are words that should be SHALL/MUST in spec files.
var weakPatterns = []string{"should", "may", "might", "TBD", "TODO", "FIXME", "TKTK", "???"}

// weakWordRegexes maps each weak pattern to a compiled regex that matches whole words.
var weakWordRegexes map[string]*regexp.Regexp

func init() {
	weakWordRegexes = make(map[string]*regexp.Regexp, len(weakPatterns))
	for _, p := range weakPatterns {
		// Case-insensitive for words, exact for abbreviations
		if p == strings.ToUpper(p) || p == "???" {
			weakWordRegexes[p] = regexp.MustCompile(regexp.QuoteMeta(p))
		} else {
			weakWordRegexes[p] = regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(p) + `\b`)
		}
	}
}

// CheckAmbiguity scans spec files for weak language patterns.
func CheckAmbiguity(changeDir string) []Finding {
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

		lines := strings.Split(string(data), "\n")
		relPath := fmt.Sprintf("specs/%s/spec.md", entry.Name())

		for lineNum, line := range lines {
			// Skip lines that are inside code blocks
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") {
				continue
			}

			for _, pattern := range weakPatterns {
				re := weakWordRegexes[pattern]
				if re.MatchString(line) {
					findings = append(findings, Finding{
						ID:             fmt.Sprintf("AMB-%d", counter),
						Dimension:      string(DimensionAmbiguity),
						Severity:       string(SeveritySuggestion),
						Location:       fmt.Sprintf("%s:%d", relPath, lineNum+1),
						Summary:        fmt.Sprintf("Vague language '%s' found", pattern),
						Recommendation: fmt.Sprintf("Replace '%s' with SHALL/SHALL NOT for clarity", pattern),
					})
					counter++
					break // One finding per line max
				}
			}
		}
	}

	return findings
}
