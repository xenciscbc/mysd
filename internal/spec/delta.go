package spec

import (
	"bufio"
	"regexp"
	"strings"
)

// reDeltaHeading matches headings that declare a delta operation section.
// Matches ## ADDED, ## MODIFIED, ## REMOVED (and ### equivalents).
var reDeltaHeading = regexp.MustCompile(`^#{1,3}\s+(ADDED|MODIFIED|REMOVED)\b`)

// DetectDeltaOp extracts the delta operation from a section heading string.
// Returns DeltaNone if no ADDED/MODIFIED/REMOVED keyword is found at the start.
func DetectDeltaOp(heading string) DeltaOp {
	m := reDeltaHeading.FindStringSubmatch(heading)
	if m == nil {
		return DeltaNone
	}
	switch m[1] {
	case "ADDED":
		return DeltaAdded
	case "MODIFIED":
		return DeltaModified
	case "REMOVED":
		return DeltaRemoved
	default:
		return DeltaNone
	}
}

// reAnyHeading matches any markdown heading (# to ###).
var reAnyHeading = regexp.MustCompile(`^#{1,3}\s+`)

// ParseDelta parses a delta spec body into categorized requirement slices.
// Sections are identified by ## ADDED / ## MODIFIED / ## REMOVED headings.
// Lines with RFC 2119 keywords within each section are extracted as requirements.
func ParseDelta(body string) (added []Requirement, modified []Requirement, removed []Requirement) {
	var currentOp DeltaOp
	scanner := bufio.NewScanner(strings.NewReader(body))

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line is a delta-section heading
		if reAnyHeading.MatchString(line) {
			op := DetectDeltaOp(line)
			if op != DeltaNone {
				currentOp = op
				continue
			}
		}

		if currentOp == DeltaNone {
			continue
		}

		// Extract RFC 2119 keyword from the line
		kw := extractKeyword(line)
		if kw == "" {
			continue
		}

		req := Requirement{
			Text:    strings.TrimSpace(line),
			Keyword: kw,
			DeltaOp: currentOp,
			Status:  StatusPending,
		}

		switch currentOp {
		case DeltaAdded:
			added = append(added, req)
		case DeltaModified:
			modified = append(modified, req)
		case DeltaRemoved:
			removed = append(removed, req)
		}
	}

	return added, modified, removed
}
