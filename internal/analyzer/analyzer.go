package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
)

// Analyze performs cross-artifact structural analysis on a change directory.
// It analyzes whatever artifacts exist without requiring all to be present.
func Analyze(changeDir string) AnalysisResult {
	changeName := filepath.Base(changeDir)

	result := AnalysisResult{
		ChangeID: changeName,
	}

	// Detect available artifacts
	hasProposal := fileExists(filepath.Join(changeDir, "proposal.md"))
	hasSpecs := dirExists(filepath.Join(changeDir, "specs"))
	hasDesign := fileExists(filepath.Join(changeDir, "design.md"))
	hasTasks := fileExists(filepath.Join(changeDir, "tasks.md"))

	if hasProposal {
		result.ArtifactsAnalyzed = append(result.ArtifactsAnalyzed, "proposal")
	} else {
		result.ArtifactsMissing = append(result.ArtifactsMissing, "proposal")
	}
	if hasSpecs {
		result.ArtifactsAnalyzed = append(result.ArtifactsAnalyzed, "specs")
	} else {
		result.ArtifactsMissing = append(result.ArtifactsMissing, "specs")
	}
	if hasDesign {
		result.ArtifactsAnalyzed = append(result.ArtifactsAnalyzed, "design")
	} else {
		result.ArtifactsMissing = append(result.ArtifactsMissing, "design")
	}
	if hasTasks {
		result.ArtifactsAnalyzed = append(result.ArtifactsAnalyzed, "tasks")
	} else {
		result.ArtifactsMissing = append(result.ArtifactsMissing, "tasks")
	}

	// Run dimensions based on available artifacts
	var coverageFindings, ambiguityFindings, consistencyFindings, gapsFindings []Finding

	if hasProposal && hasSpecs {
		coverageFindings = CheckCoverage(changeDir)
	}

	if hasSpecs {
		ambiguityFindings = CheckAmbiguity(changeDir)
	}

	if hasProposal && hasDesign && hasTasks {
		consistencyFindings = CheckConsistency(changeDir)
	}

	if hasSpecs && hasTasks {
		gapsFindings = CheckGaps(changeDir)
	}

	// Build dimension results
	result.Dimensions = []DimensionResult{
		buildDimensionResult(DimensionCoverage, coverageFindings, hasProposal && hasSpecs),
		buildDimensionResult(DimensionConsistency, consistencyFindings, hasProposal && hasDesign && hasTasks),
		buildDimensionResult(DimensionAmbiguity, ambiguityFindings, hasSpecs),
		buildDimensionResult(DimensionGaps, gapsFindings, hasSpecs && hasTasks),
	}

	// Collect all findings
	result.Findings = append(result.Findings, coverageFindings...)
	result.Findings = append(result.Findings, consistencyFindings...)
	result.Findings = append(result.Findings, ambiguityFindings...)
	result.Findings = append(result.Findings, gapsFindings...)

	// Ensure non-nil slices for JSON
	if result.Findings == nil {
		result.Findings = []Finding{}
	}
	if result.ArtifactsAnalyzed == nil {
		result.ArtifactsAnalyzed = []string{}
	}
	if result.ArtifactsMissing == nil {
		result.ArtifactsMissing = []string{}
	}

	return result
}

func buildDimensionResult(dim Dimension, findings []Finding, applicable bool) DimensionResult {
	if !applicable {
		return DimensionResult{
			Dimension:    string(dim),
			Status:       "Skipped (insufficient artifacts)",
			FindingCount: 0,
		}
	}
	count := len(findings)
	status := "Clean"
	if count > 0 {
		status = fmt.Sprintf("%d issue(s) found", count)
	}
	return DimensionResult{
		Dimension:    string(dim),
		Status:       status,
		FindingCount: count,
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
