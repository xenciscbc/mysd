package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/xenciscbc/mysd/internal/spec"
	"gopkg.in/yaml.v3"
)

// validDeltas is the set of valid DeltaOp values.
var validDeltas = map[spec.DeltaOp]bool{
	spec.DeltaAdded:    true,
	spec.DeltaModified: true,
	spec.DeltaRemoved:  true,
	spec.DeltaRenamed:  true,
}

// Validate performs structural validation on a change directory.
// It checks file existence, frontmatter schemas, required fields, and cross-field consistency.
func Validate(changeDir string) ValidationResult {
	changeName := filepath.Base(changeDir)
	result := ValidationResult{
		ChangeID: changeName,
		Valid:    true,
	}

	// 1. File existence
	validateFileExistence(changeDir, &result)

	// 2. ChangeMeta (.openspec.yaml)
	validateChangeMeta(changeDir, &result)

	// 3. Proposal frontmatter
	validateProposal(changeDir, changeName, &result)

	// 4. Spec frontmatter (each specs/*/spec.md)
	validateSpecs(changeDir, &result)

	// 5. Tasks frontmatter (if tasks.md exists)
	validateTasks(changeDir, &result)

	// Ensure non-nil slices for JSON
	if result.Errors == nil {
		result.Errors = []ValidationFinding{}
	}
	if result.Warnings == nil {
		result.Warnings = []ValidationFinding{}
	}

	result.Valid = len(result.Errors) == 0
	return result
}

func addError(result *ValidationResult, location, message string) {
	result.Errors = append(result.Errors, ValidationFinding{
		Severity: "error",
		Location: location,
		Message:  message,
	})
}

func addWarning(result *ValidationResult, location, message string) {
	result.Warnings = append(result.Warnings, ValidationFinding{
		Severity: "warning",
		Location: location,
		Message:  message,
	})
}

func validateFileExistence(changeDir string, result *ValidationResult) {
	metaPath := filepath.Join(changeDir, ".openspec.yaml")
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		addError(result, ".openspec.yaml", "file not found")
	}

	proposalPath := filepath.Join(changeDir, "proposal.md")
	if _, err := os.Stat(proposalPath); os.IsNotExist(err) {
		addError(result, "proposal.md", "file not found")
	}
}

func validateChangeMeta(changeDir string, result *ValidationResult) {
	path := filepath.Join(changeDir, ".openspec.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return // already reported in file existence check
	}

	var meta spec.ChangeMeta
	if err := yaml.Unmarshal(data, &meta); err != nil {
		addError(result, ".openspec.yaml", fmt.Sprintf("invalid YAML: %v", err))
		return
	}

	if meta.Schema == "" {
		addError(result, ".openspec.yaml", "missing required field: schema")
	}
	if meta.Created == "" {
		addError(result, ".openspec.yaml", "missing required field: created")
	}
}

func validateProposal(changeDir, changeName string, result *ValidationResult) {
	path := filepath.Join(changeDir, "proposal.md")
	f, err := os.Open(path)
	if err != nil {
		return // already reported in file existence check
	}
	defer f.Close()

	var fm spec.ProposalFrontmatter
	_, err = frontmatter.Parse(f, &fm)
	if err != nil {
		addWarning(result, "proposal.md", "no valid frontmatter found (brownfield format)")
		return
	}

	// Detect brownfield: frontmatter.Parse succeeds but all fields are zero-value
	if fm.SpecVersion == "" && fm.ChangeName == "" && fm.Status == "" && fm.Created == "" {
		addWarning(result, "proposal.md", "no valid frontmatter found (brownfield format)")
		return
	}

	if fm.SpecVersion == "" {
		addError(result, "proposal.md", "missing required field: spec-version")
	}
	if fm.ChangeName == "" {
		addError(result, "proposal.md", "missing required field: change")
	} else if fm.ChangeName != changeName {
		addError(result, "proposal.md", fmt.Sprintf("change name mismatch: frontmatter has %q, directory is %q", fm.ChangeName, changeName))
	}
	if fm.Status == "" {
		addError(result, "proposal.md", "missing required field: status")
	}
	if fm.Created == "" {
		addError(result, "proposal.md", "missing required field: created")
	}
}

func validateSpecs(changeDir string, result *ValidationResult) {
	specsDir := filepath.Join(changeDir, "specs")
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return // specs dir is optional at early stages
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		capName := entry.Name()
		specPath := filepath.Join(specsDir, capName, "spec.md")
		relPath := fmt.Sprintf("specs/%s/spec.md", capName)

		if _, statErr := os.Stat(specPath); os.IsNotExist(statErr) {
			addWarning(result, fmt.Sprintf("specs/%s/", capName), "directory exists but spec.md is missing")
			continue
		}

		validateSpecFile(specPath, relPath, capName, result)
	}
}

func validateSpecFile(path, relPath, expectedCapability string, result *ValidationResult) {
	f, err := os.Open(path)
	if err != nil {
		addError(result, relPath, fmt.Sprintf("cannot read: %v", err))
		return
	}
	defer f.Close()

	var fm spec.SpecFrontmatter
	_, err = frontmatter.Parse(f, &fm)
	if err != nil {
		addWarning(result, relPath, "no valid frontmatter found (brownfield format)")
		return
	}

	// Detect brownfield: all fields zero-value
	if fm.SpecVersion == "" && fm.Capability == "" && fm.Delta == "" && fm.Status == "" {
		addWarning(result, relPath, "no valid frontmatter found (brownfield format)")
		return
	}

	if fm.SpecVersion == "" {
		addError(result, relPath, "missing required field: spec-version")
	}
	if fm.Capability == "" {
		addError(result, relPath, "missing required field: capability")
	} else if fm.Capability != expectedCapability {
		addError(result, relPath, fmt.Sprintf("capability mismatch: expected %q, got %q", expectedCapability, fm.Capability))
	}
	if fm.Delta == "" {
		addError(result, relPath, "missing required field: delta")
	} else if !validDeltas[fm.Delta] {
		addError(result, relPath, fmt.Sprintf("invalid delta value: %q (must be ADDED, MODIFIED, REMOVED, or RENAMED)", fm.Delta))
	}
	if fm.Status == "" {
		addError(result, relPath, "missing required field: status")
	}
}

func validateTasks(changeDir string, result *ValidationResult) {
	path := filepath.Join(changeDir, "tasks.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return // tasks.md is optional at early stages
	}

	f, err := os.Open(path)
	if err != nil {
		addError(result, "tasks.md", fmt.Sprintf("cannot read: %v", err))
		return
	}
	defer f.Close()

	var fm spec.TasksFrontmatter
	_, fmErr := frontmatter.Parse(f, &fm)
	if fmErr != nil {
		addWarning(result, "tasks.md", "no valid frontmatter found (brownfield format)")
		return
	}

	// Detect brownfield: all fields zero-value
	if fm.SpecVersion == "" && fm.Total == 0 && fm.Completed == 0 {
		addWarning(result, "tasks.md", "no valid frontmatter found (brownfield format)")
		return
	}

	if fm.SpecVersion == "" {
		addError(result, "tasks.md", "missing required field: spec-version")
	}

	// Count actual tasks from file
	tasks, _, parseErr := spec.ParseTasks(path)
	if parseErr != nil {
		addError(result, "tasks.md", fmt.Sprintf("parse error: %v", parseErr))
		return
	}

	actualCount := len(tasks)
	if fm.Total != actualCount {
		addWarning(result, "tasks.md", fmt.Sprintf("total (%d) does not match actual task count (%d)", fm.Total, actualCount))
	}

	// Check for empty task names
	for _, t := range tasks {
		if strings.TrimSpace(t.Name) == "" {
			addWarning(result, "tasks.md", fmt.Sprintf("task T%d has empty name", t.ID))
		}
	}
}
