package analyzer

// Severity represents the severity level of an analysis finding.
type Severity string

const (
	SeverityCritical   Severity = "Critical"
	SeverityWarning    Severity = "Warning"
	SeveritySuggestion Severity = "Suggestion"
)

// Dimension represents an analysis dimension name.
type Dimension string

const (
	DimensionCoverage    Dimension = "Coverage"
	DimensionConsistency Dimension = "Consistency"
	DimensionAmbiguity   Dimension = "Ambiguity"
	DimensionGaps        Dimension = "Gaps"
)

// Finding represents a single analysis finding.
type Finding struct {
	ID             string   `json:"id"`
	Dimension      string   `json:"dimension"`
	Severity       string   `json:"severity"`
	Location       string   `json:"location"`
	Summary        string   `json:"summary"`
	Recommendation string   `json:"recommendation"`
}

// DimensionResult represents the result of analyzing one dimension.
type DimensionResult struct {
	Dimension    string `json:"dimension"`
	Status       string `json:"status"`
	FindingCount int    `json:"finding_count"`
}

// AnalysisResult represents the complete analysis output.
type AnalysisResult struct {
	ChangeID          string            `json:"change_id"`
	Dimensions        []DimensionResult `json:"dimensions"`
	Findings          []Finding         `json:"findings"`
	ArtifactsAnalyzed []string          `json:"artifacts_analyzed"`
	ArtifactsMissing  []string          `json:"artifacts_missing"`
}
