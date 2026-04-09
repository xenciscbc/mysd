package validator

// ValidationFinding represents a single validation issue.
type ValidationFinding struct {
	Severity string `json:"severity"` // "error" or "warning"
	Location string `json:"location"` // e.g. "proposal.md", "specs/auth/spec.md"
	Message  string `json:"message"`
}

// ValidationResult is the complete output of artifact validation.
type ValidationResult struct {
	ChangeID string              `json:"change_id"`
	Valid    bool                `json:"valid"` // true if no errors (warnings OK)
	Errors   []ValidationFinding `json:"errors"`
	Warnings []ValidationFinding `json:"warnings"`
}
