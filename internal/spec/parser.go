package spec

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// RFC 2119 case-sensitive regex patterns.
// Only uppercase keywords are RFC 2119; lowercase is natural language.
var (
	reMust   = regexp.MustCompile(`\bMUST\b|\bMUST NOT\b|\bREQUIRED\b|\bSHALL\b|\bSHALL NOT\b`)
	reShould = regexp.MustCompile(`\bSHOULD\b|\bSHOULD NOT\b|\bRECOMMENDED\b`)
	reMay    = regexp.MustCompile(`\bMAY\b|\bOPTIONAL\b`)
)

// reRequirementHeading matches headings that introduce a requirement block.
// Avoids Pitfall 1: not hardcoded to a specific heading level.
var reRequirementHeading = regexp.MustCompile(`^#{1,3}\s+Requirement:?\s*(.*)$`)

// extractKeyword performs case-sensitive RFC 2119 keyword detection on a single line.
// Returns empty string if no RFC 2119 keyword is found.
// Note: lowercase "must", "should", "may" are intentionally NOT matched.
func extractKeyword(line string) RFC2119Keyword {
	if reMust.MatchString(line) {
		return Must
	}
	if reShould.MatchString(line) {
		return Should
	}
	if reMay.MatchString(line) {
		return May
	}
	return ""
}

// ParseProposal reads a proposal.md file and returns a ProposalDoc.
// Gracefully handles brownfield files with no frontmatter by returning
// zero-value frontmatter and the full file content as body.
func ParseProposal(path string) (ProposalDoc, error) {
	f, err := os.Open(path)
	if err != nil {
		return ProposalDoc{}, err
	}
	defer f.Close()

	var fm ProposalFrontmatter
	rest, err := frontmatter.Parse(f, &fm)
	if err != nil {
		// Brownfield: no valid frontmatter — re-read raw content as body
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return ProposalDoc{}, readErr
		}
		return ProposalDoc{Body: string(content)}, nil
	}

	return ProposalDoc{
		Frontmatter: fm,
		Body:        string(rest),
	}, nil
}

// ParseSpec reads a spec.md file and extracts all Requirement entries.
// Handles both brownfield (no frontmatter) and native (with frontmatter) formats.
// Extracts RFC 2119 keywords from all lines, grouping under the nearest heading.
func ParseSpec(path string) ([]Requirement, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var fm SpecFrontmatter
	rest, err := frontmatter.Parse(f, &fm)
	if err != nil {
		// Brownfield: treat full file as body
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil, readErr
		}
		rest = content
	}

	reqs := parseRequirementsFromBody(string(rest), fm.Delta)
	baseName := filepath.Base(path)
	for i := range reqs {
		reqs[i].SourceFile = baseName
	}
	return reqs, nil
}

// parseRequirementsFromBody scans the body text for RFC 2119 keywords
// and groups them under the nearest Requirement heading.
func parseRequirementsFromBody(body string, defaultDeltaOp DeltaOp) []Requirement {
	var reqs []Requirement
	scanner := bufio.NewScanner(strings.NewReader(body))

	var currentHeading string
	reqID := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Check for a Requirement: heading
		if m := reRequirementHeading.FindStringSubmatch(line); m != nil {
			currentHeading = strings.TrimSpace(m[1])
			continue
		}

		// Extract RFC 2119 keyword from any line
		kw := extractKeyword(line)
		if kw == "" {
			continue
		}

		reqID++
		reqs = append(reqs, Requirement{
			ID:      "",
			Text:    strings.TrimSpace(line),
			Keyword: kw,
			DeltaOp: defaultDeltaOp,
			Status:  StatusPending,
			// Note: currentHeading is available but not stored in ID by default
			// The heading context is implicit through ordering
		})
		_ = currentHeading // suppress unused warning; available for future use
	}

	return reqs
}

// ParseTasks reads a tasks.md file and returns parsed Task slices and frontmatter.
// Handles both brownfield (no frontmatter) and native (with frontmatter) formats.
func ParseTasks(path string) ([]Task, TasksFrontmatter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, TasksFrontmatter{}, err
	}
	defer f.Close()

	var fm TasksFrontmatter
	rest, err := frontmatter.Parse(f, &fm)
	if err != nil {
		// Brownfield: treat full file as body
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil, TasksFrontmatter{}, readErr
		}
		rest = content
	}

	tasks := parseTaskLines(string(rest))
	return tasks, fm, nil
}

// reTaskLine matches markdown task list items: `- [ ] text`, `- [x] text`, or `- [~] text`
var reTaskLine = regexp.MustCompile(`^-\s+\[([xX ~])\]\s+(.+)$`)

// reSkipReason extracts skip reason from parenthetical or colon notation.
// Matches: "（跳過：reason）", "(跳過：reason)", "（reason）", or trailing "：reason" / ": reason"
var reSkipReason = regexp.MustCompile(`[（(](?:跳過[：:]?\s*)?(.+?)[）)]`)

// parseTaskLines extracts Task entries from a task list body.
func parseTaskLines(body string) []Task {
	var tasks []Task
	scanner := bufio.NewScanner(strings.NewReader(body))
	id := 0

	for scanner.Scan() {
		line := scanner.Text()
		m := reTaskLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		id++
		marker := m[1]
		text := strings.TrimSpace(m[2])

		task := Task{
			ID:   id,
			Name: text,
		}

		switch {
		case marker == "x" || marker == "X":
			task.Status = StatusDone
		case marker == "~":
			task.Status = StatusDone
			task.Skipped = true
			// Extract skip reason
			if rm := reSkipReason.FindStringSubmatch(text); rm != nil {
				task.SkipReason = strings.TrimSpace(rm[1])
			}
		default:
			task.Status = StatusPending
		}

		tasks = append(tasks, task)
	}

	return tasks
}

// ParseChangeMeta reads a .openspec.yaml file and returns ChangeMeta.
func ParseChangeMeta(path string) (ChangeMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ChangeMeta{}, err
	}

	var meta ChangeMeta
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return ChangeMeta{}, err
	}
	return meta, nil
}

// ParseChange assembles a complete Change from a change directory by orchestrating
// all parse functions. It handles both brownfield and native formats gracefully.
func ParseChange(changeDir string) (Change, error) {
	name := filepath.Base(changeDir)

	// Parse .openspec.yaml
	meta, err := ParseChangeMeta(filepath.Join(changeDir, ".openspec.yaml"))
	if err != nil {
		return Change{}, err
	}

	// Parse proposal.md
	proposal, err := ParseProposal(filepath.Join(changeDir, "proposal.md"))
	if err != nil {
		return Change{}, err
	}

	// Parse design.md (always plain markdown, no frontmatter)
	var design DesignDoc
	designContent, err := os.ReadFile(filepath.Join(changeDir, "design.md"))
	if err == nil {
		design.Body = string(designContent)
	}

	// Parse tasks.md
	tasks, _, err := ParseTasks(filepath.Join(changeDir, "tasks.md"))
	if err != nil {
		return Change{}, err
	}

	// Parse all spec files under specs/
	var allReqs []Requirement
	specsDir := filepath.Join(changeDir, "specs")
	specEntries, err := os.ReadDir(specsDir)
	if err == nil {
		for _, capDir := range specEntries {
			if !capDir.IsDir() {
				continue
			}
			specFile := filepath.Join(specsDir, capDir.Name(), "spec.md")
			if _, statErr := os.Stat(specFile); statErr != nil {
				continue
			}
			reqs, parseErr := ParseSpec(specFile)
			if parseErr == nil {
				allReqs = append(allReqs, reqs...)
			}
		}
	}

	return Change{
		Name:     name,
		Dir:      changeDir,
		Proposal: proposal,
		Specs:    allReqs,
		Design:   design,
		Tasks:    tasks,
		Meta:     meta,
	}, nil
}
