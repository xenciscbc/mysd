package executor

import (
	"path/filepath"
	"strings"
)

// AlignmentPath returns the canonical path for the alignment summary file (per D-07).
// The alignment.md lives at specsDir/changes/{changeName}/alignment.md.
// Uses forward slashes for spec path conventions (cross-platform compatibility).
func AlignmentPath(changeName string, specsDir string) string {
	p := filepath.Join(specsDir, "changes", changeName, "alignment.md")
	return strings.ReplaceAll(p, "\\", "/")
}

// AlignmentTemplate returns the markdown template for AI alignment summaries (per D-06).
// The AI fills this in before starting execution to confirm spec understanding.
const alignmentTemplateStr = `## Alignment Summary

### MUST items I understand:

<!-- List each MUST requirement and your interpretation -->

### Execution strategy:

<!-- Describe how you will approach the tasks, in what order, and why -->

### Open questions (if any):

<!-- List any ambiguities or clarifications needed before proceeding -->
`

// AlignmentTemplate returns the markdown template string for alignment summaries.
func AlignmentTemplate() string {
	return alignmentTemplateStr
}
