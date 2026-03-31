package executor

// PreflightReport is the JSON output of --preflight validation.
type PreflightReport struct {
	Status string         `json:"status"`
	Checks PreflightChecks `json:"checks"`
}

// PreflightChecks contains the individual check results.
type PreflightChecks struct {
	MissingFiles []string       `json:"missing_files"`
	Staleness    StalenessCheck `json:"staleness"`
}

// StalenessCheck reports artifact freshness.
type StalenessCheck struct {
	DaysSinceLastPlan int  `json:"days_since_last_plan"`
	IsStale           bool `json:"is_stale"`
}
