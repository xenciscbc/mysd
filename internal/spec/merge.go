package spec

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// AppVersion is set by the CLI layer to inject the current mysd version.
var AppVersion = "dev"

// reReqHeading matches requirement headings: ### Requirement: <name>
var reReqHeading = regexp.MustCompile(`(?i)^###\s+Requirement:\s*(.+)$`)

// MergeSpecs merges a delta spec body into the main spec file at mainSpecPath.
// It applies operations in OpenSpec order: RENAMED → REMOVED → MODIFIED → ADDED.
// Returns the merged content and any warnings encountered.
// If mainSpecPath does not exist, a new spec is created from the ADDED requirements.
func MergeSpecs(mainSpecPath string, deltaBody string) (string, []string, error) {
	added, modified, removed, renamed := ParseDelta(deltaBody)

	var warnings []string

	// Fallback: when ParseDelta finds no delta headings, use frontmatter delta field
	if len(added) == 0 && len(modified) == 0 && len(removed) == 0 && len(renamed) == 0 {
		return mergeFallback(mainSpecPath, deltaBody)
	}

	// Read existing main spec content (or start empty if file doesn't exist)
	var content string
	var fm SpecFrontmatter
	isNew := false
	data, err := os.ReadFile(mainSpecPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", nil, fmt.Errorf("read main spec: %w", err)
		}
		// No existing main spec — will create new
		content = ""
		isNew = true
	} else {
		// Parse frontmatter from existing spec
		rest, fmErr := frontmatter.Parse(bytes.NewReader(data), &fm)
		if fmErr != nil {
			// No frontmatter — use full content as body
			content = string(data)
		} else {
			content = string(rest)
		}
	}

	hasModified := len(modified) > 0

	// 1. RENAMED: find heading and rename it
	for _, r := range renamed {
		newContent, ok := renameRequirement(content, r.From, r.To)
		if !ok {
			warnings = append(warnings, fmt.Sprintf("RENAMED: heading %q not found, skipping", r.From))
			continue
		}
		content = newContent
	}

	// 2. REMOVED: find heading and delete the entire requirement block
	for _, r := range removed {
		newContent, ok := removeRequirement(content, r.Text)
		if !ok {
			warnings = append(warnings, fmt.Sprintf("REMOVED: heading %q not found, skipping", r.Text))
			continue
		}
		content = newContent
	}

	// 3. MODIFIED: find heading and replace the entire requirement block
	for _, r := range modified {
		newContent, ok := modifyRequirement(content, r.Text)
		if !ok {
			warnings = append(warnings, fmt.Sprintf("MODIFIED: heading %q not found, skipping", r.Text))
			continue
		}
		content = newContent
	}

	// 4. ADDED: append to the end of the Requirements section or file
	for _, r := range added {
		content = addRequirement(content, r.Text)
	}

	// Update frontmatter
	if isNew {
		fm.Version = "1.0.0"
		fm.GeneratedBy = "mysd v" + AppVersion
	} else if hasModified {
		fm.Version = incrementMinorVersion(fm.Version)
	}

	// Serialize frontmatter + body
	result := renderWithFrontmatter(fm, content)
	return result, warnings, nil
}

// mergeFallback handles delta specs that have no delta section headings.
// It reads the frontmatter delta field to decide the merge strategy:
// ADDED → use delta body as new spec content; MODIFIED → replace main spec body + increment version.
func mergeFallback(mainSpecPath string, deltaBody string) (string, []string, error) {
	var deltaFM SpecFrontmatter
	rest, fmErr := frontmatter.Parse(bytes.NewReader([]byte(deltaBody)), &deltaFM)
	if fmErr != nil {
		return "", []string{"fallback: could not parse delta frontmatter, skipping merge"}, nil
	}
	body := string(rest)

	switch deltaFM.Delta {
	case DeltaAdded:
		fm := SpecFrontmatter{
			SpecVersion: deltaFM.SpecVersion,
			Capability:  deltaFM.Capability,
			Version:     "1.0.0",
			GeneratedBy: "mysd v" + AppVersion,
		}
		return renderWithFrontmatter(fm, body), nil, nil

	case DeltaModified:
		// Read existing main spec frontmatter
		var mainFM SpecFrontmatter
		data, err := os.ReadFile(mainSpecPath)
		if err != nil {
			if os.IsNotExist(err) {
				// No existing spec — treat as ADDED
				fm := SpecFrontmatter{
					SpecVersion: deltaFM.SpecVersion,
					Capability:  deltaFM.Capability,
					Version:     "1.0.0",
					GeneratedBy: "mysd v" + AppVersion,
				}
				return renderWithFrontmatter(fm, body), nil, nil
			}
			return "", nil, fmt.Errorf("read main spec: %w", err)
		}
		if _, fErr := frontmatter.Parse(bytes.NewReader(data), &mainFM); fErr != nil {
			mainFM = SpecFrontmatter{}
		}
		mainFM.Version = incrementMinorVersion(mainFM.Version)
		mainFM.GeneratedBy = "mysd v" + AppVersion
		if mainFM.Capability == "" {
			mainFM.Capability = deltaFM.Capability
		}
		if mainFM.SpecVersion == "" {
			mainFM.SpecVersion = deltaFM.SpecVersion
		}
		return renderWithFrontmatter(mainFM, body), nil, nil

	default:
		return "", []string{fmt.Sprintf("fallback: delta spec has no parseable operations (delta=%q), skipping merge", deltaFM.Delta)}, nil
	}
}

// renameRequirement finds "### Requirement: <from>" and renames it to "### Requirement: <to>".
func renameRequirement(content, from, to string) (string, bool) {
	lines := strings.Split(content, "\n")
	found := false
	for i, line := range lines {
		if matchReqHeading(line, from) {
			lines[i] = "### Requirement: " + to
			found = true
			break
		}
	}
	if !found {
		return content, false
	}
	return strings.Join(lines, "\n"), true
}

// removeRequirement finds "### Requirement: <name>" and removes the entire block
// (from heading to next requirement heading or EOF).
func removeRequirement(content, reqText string) (string, bool) {
	// For REMOVED, the reqText is the RFC2119 line. We need to find which requirement
	// block contains this text. We'll search by requirement name extracted from the text.
	// However, the design says REMOVED section contains requirement names.
	// The ParseDelta extracts lines with RFC2119 keywords as requirement Text.
	// We need to match by searching the content for a block containing this text.
	lines := strings.Split(content, "\n")
	blockStart := -1
	blockEnd := len(lines)

	for i, line := range lines {
		if reReqHeading.MatchString(line) {
			if blockStart >= 0 {
				// Found the next heading — end of block
				blockEnd = i
				break
			}
			// Check if the block under this heading contains the reqText
			if blockContainsText(lines, i, reqText) {
				blockStart = i
			}
		}
	}

	if blockStart < 0 {
		return content, false
	}

	// Remove lines from blockStart to blockEnd
	result := make([]string, 0, len(lines))
	result = append(result, lines[:blockStart]...)
	result = append(result, lines[blockEnd:]...)
	return strings.Join(result, "\n"), true
}

// modifyRequirement finds the requirement block containing reqText and replaces
// the requirement content line with the new text.
func modifyRequirement(content, reqText string) (string, bool) {
	lines := strings.Split(content, "\n")

	// Find the requirement block and replace the RFC2119 line
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == reqText {
			lines[i] = reqText
			return strings.Join(lines, "\n"), true
		}
	}

	return content, false
}

// addRequirement appends a new requirement line to the end of the content.
func addRequirement(content, reqText string) string {
	if content == "" {
		return reqText + "\n"
	}
	// Ensure content ends with newline before appending
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + "\n" + reqText + "\n"
}

// matchReqHeading checks if a line is "### Requirement: <name>" (case-insensitive match on name).
func matchReqHeading(line, name string) bool {
	m := reReqHeading.FindStringSubmatch(line)
	if m == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(m[1]), strings.TrimSpace(name))
}

// blockContainsText checks if the block starting at headingIdx contains reqText.
func blockContainsText(lines []string, headingIdx int, reqText string) bool {
	for i := headingIdx + 1; i < len(lines); i++ {
		if reReqHeading.MatchString(lines[i]) {
			break
		}
		if strings.TrimSpace(lines[i]) == strings.TrimSpace(reqText) {
			return true
		}
	}
	return false
}

// incrementMinorVersion takes a semver string like "1.0.0" and returns "1.1.0".
// If the version is empty or unparseable, returns "1.1.0" as a safe default.
func incrementMinorVersion(version string) string {
	if version == "" {
		return "1.1.0"
	}
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return "1.1.0"
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "1.1.0"
	}
	parts[1] = strconv.Itoa(minor + 1)
	// Reset patch to 0 on minor bump
	if len(parts) >= 3 {
		parts[2] = "0"
	}
	return strings.Join(parts, ".")
}

// renderWithFrontmatter serializes SpecFrontmatter as YAML frontmatter
// and prepends it to the body content.
func renderWithFrontmatter(fm SpecFrontmatter, body string) string {
	fmBytes, err := yaml.Marshal(&fm)
	if err != nil {
		// Fallback: return body only
		return body
	}
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(fmBytes)
	sb.WriteString("---\n")
	sb.WriteString(body)
	return sb.String()
}
