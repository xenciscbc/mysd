package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/spec"
)

var (
	instructionsChange string
	instructionsJSON   bool
)

var instructionsCmd = &cobra.Command{
	Use:   "instructions <artifact-id>",
	Short: "Output structured artifact instructions for agent consumption",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstructions,
}

func init() {
	instructionsCmd.Flags().StringVar(&instructionsChange, "change", "", "change name (required)")
	instructionsCmd.Flags().BoolVar(&instructionsJSON, "json", false, "output as JSON")
	_ = instructionsCmd.MarkFlagRequired("change")
	rootCmd.AddCommand(instructionsCmd)
}

// InstructionsDependency describes a dependency artifact's status.
type InstructionsDependency struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Done bool   `json:"done"`
}

// InstructionsOutput is the JSON structure returned by mysd instructions.
type InstructionsOutput struct {
	ArtifactID          string                   `json:"artifactId"`
	ChangeName          string                   `json:"changeName"`
	OutputPath          string                   `json:"outputPath"`
	Template            string                   `json:"template"`
	Rules               []string                 `json:"rules"`
	Instruction         string                   `json:"instruction"`
	Dependencies        []InstructionsDependency `json:"dependencies"`
	SelfReviewChecklist []string                 `json:"selfReviewChecklist"`
}

func runInstructions(cmd *cobra.Command, args []string) error {
	artifactID := args[0]

	if artifactID != "design" && artifactID != "tasks" {
		return fmt.Errorf("unknown artifact ID %q: supported values are \"design\" and \"tasks\"", artifactID)
	}

	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return fmt.Errorf("no spec directory: %w", err)
	}

	changeDir := filepath.Join(specDir, "changes", instructionsChange)
	if _, statErr := os.Stat(changeDir); statErr != nil {
		return fmt.Errorf("change %q not found at %s", instructionsChange, changeDir)
	}

	out := buildInstructions(artifactID, instructionsChange, specDir, changeDir)

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal instructions: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

func buildInstructions(artifactID, changeName, specDir, changeDir string) InstructionsOutput {
	switch artifactID {
	case "design":
		return buildDesignInstructions(changeName, specDir, changeDir)
	case "tasks":
		return buildTasksInstructions(changeName, specDir, changeDir)
	default:
		return InstructionsOutput{ArtifactID: artifactID, ChangeName: changeName}
	}
}

func buildDesignInstructions(changeName, specDir, changeDir string) InstructionsOutput {
	proposalDone := fileExists(filepath.Join(changeDir, "proposal.md"))
	specsDone := dirHasFiles(filepath.Join(changeDir, "specs"))

	return InstructionsOutput{
		ArtifactID: "design",
		ChangeName: changeName,
		OutputPath: filepath.Join(changeDir, "design.md"),
		Template: `## Context

<!-- Background and current state -->

## Goals / Non-Goals

**Goals:**

<!-- What this design aims to achieve -->

**Non-Goals:**

<!-- What is explicitly out of scope -->

## Decisions

### D-XX: Decision Title

<!-- Key design decisions with rationale. Include alternatives considered. -->

## Risks / Trade-offs

<!-- Known risks and trade-offs. Format: [Risk] → Mitigation -->`,
		Rules: []string{
			"No TBD/TODO/FIXME placeholders — fill in all sections with concrete content",
			"Every capability listed in proposal.md must have a corresponding section in the design",
			"Each decision must include at least one alternative considered with rationale for rejection",
			"File paths referenced in the design must be consistent with proposal Impact section",
		},
		Instruction: "Create the design document that explains HOW to implement the change. " +
			"Focus on architecture and approach, not line-by-line implementation. " +
			"Reference the proposal for motivation and specs for requirements. " +
			"Good design docs explain the 'why' behind technical decisions.",
		Dependencies: []InstructionsDependency{
			{ID: "proposal", Path: "proposal.md", Done: proposalDone},
			{ID: "specs", Path: "specs/", Done: specsDone},
		},
		SelfReviewChecklist: []string{
			"No TBD/TODO/FIXME placeholders in any section",
			"Every capability in the proposal has a corresponding design section",
			"Each decision rationale includes at least one alternative considered",
			"File paths referenced are consistent with proposal Impact section",
		},
	}
}

func buildTasksInstructions(changeName, specDir, changeDir string) InstructionsOutput {
	proposalDone := fileExists(filepath.Join(changeDir, "proposal.md"))
	specsDone := dirHasFiles(filepath.Join(changeDir, "specs"))
	designDone := fileExists(filepath.Join(changeDir, "design.md"))

	return InstructionsOutput{
		ArtifactID: "tasks",
		ChangeName: changeName,
		OutputPath: filepath.Join(changeDir, "tasks.md"),
		Template: `---
spec-version: "1.0"
total: {N}
completed: 0
tasks:
  - id: 1
    name: "{Task Name}"
    description: "{Brief description of what to implement}"
    spec: "{spec-directory-name}"
    status: pending
    depends: []
    files: ["{file1.go}", "{file2.go}"]
    satisfies: ["{REQ-ID}"]
---

# Tasks: {change_name}

{Optional markdown body with implementation notes or additional context}`,
		Rules: []string{
			"No TBD/TODO/FIXME placeholders — every task must have a concrete description",
			"Each task must include a spec field matching the spec directory name (e.g., material-selection)",
			"Tasks must be ordered by dependency — a task must not depend on a later task",
			"Each task should target at most 3 files to keep scope manageable",
			"Every MUST requirement must be covered by at least one task via the satisfies field",
			"Task dependencies must form a valid DAG — no circular references",
		},
		Instruction: "Create the task list that breaks the design into executable units. " +
			"Each task should be small enough to complete in one session. " +
			"Assign each task a spec field matching its capability area. " +
			"Reference specs for what needs to be built and design.md for how.",
		Dependencies: []InstructionsDependency{
			{ID: "proposal", Path: "proposal.md", Done: proposalDone},
			{ID: "specs", Path: "specs/", Done: specsDone},
			{ID: "design", Path: "design.md", Done: designDone},
		},
		SelfReviewChecklist: []string{
			"No TBD/TODO/FIXME placeholders in task descriptions",
			"Every MUST requirement has at least one task with a matching satisfies entry",
			"No single task targets more than 3 files",
			"All file paths referenced in tasks appear in proposal Impact or design",
			"Task dependencies form a valid DAG (no circular references)",
		},
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirHasFiles(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			return true // spec dirs contain subdirectories
		}
	}
	return false
}
