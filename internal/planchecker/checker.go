package planchecker

import "github.com/xenciscbc/mysd/internal/spec"

// CoverageResult holds the output of a plan coverage check (D-01, D-03).
type CoverageResult struct {
	TotalMust     int      `json:"total_must"`
	CoveredCount  int      `json:"covered_count"`
	UncoveredIDs  []string `json:"uncovered_ids,omitempty"`
	CoverageRatio float64  `json:"coverage_ratio"`
	Passed        bool     `json:"passed"`
}

// CheckCoverage validates that every MUST item ID in mustIDs appears in at least
// one task's Satisfies field. Uses exact string matching — no AI inference (D-01).
// Pure function: no filesystem I/O, no side effects.
func CheckCoverage(tasks []spec.TaskEntry, mustIDs []string) CoverageResult {
	if len(mustIDs) == 0 {
		return CoverageResult{Passed: true}
	}

	covered := make(map[string]bool)
	for _, t := range tasks {
		for _, id := range t.Satisfies {
			covered[id] = true
		}
	}

	var uncovered []string
	for _, id := range mustIDs {
		if !covered[id] {
			uncovered = append(uncovered, id)
		}
	}

	total := len(mustIDs)
	coveredCount := total - len(uncovered)
	ratio := float64(coveredCount) / float64(total)

	return CoverageResult{
		TotalMust:     total,
		CoveredCount:  coveredCount,
		UncoveredIDs:  uncovered,
		CoverageRatio: ratio,
		Passed:        len(uncovered) == 0,
	}
}
