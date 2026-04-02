package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Phase represents the lifecycle phase of a spec workflow.
type Phase string

const (
	PhaseNone     Phase = ""
	PhaseProposed Phase = "proposed"
	PhaseSpecced  Phase = "specced"
	PhaseDesigned Phase = "designed"
	PhasePlanned  Phase = "planned"
	PhaseExecuted Phase = "executed"
	PhaseVerified Phase = "verified"
	PhaseArchived Phase = "archived"
)

// ErrInvalidTransition is returned when a state transition is not allowed.
var ErrInvalidTransition = errors.New("invalid state transition")

// WorkflowState holds the current state of a spec workflow.
type WorkflowState struct {
	ChangeName string    `json:"change_name"`
	Phase      Phase     `json:"phase"`
	LastRun    time.Time `json:"last_run"`
	VerifyPass *bool     `json:"verify_pass,omitempty"`
}

// MysdDir returns the .mysd directory path derived from specsDir.
// specsDir is always a child of the project root (e.g., "openspec/" or ".specs/"),
// so filepath.Dir(specsDir) gives the project root.
func MysdDir(specsDir string) string {
	return filepath.Join(filepath.Dir(specsDir), ".mysd")
}

// LoadState reads WorkflowState from .mysd/STATE.json (derived from specsDir).
// Falls back to legacy locations (specsDir/STATE.json) for backward compatibility.
// If no state file exists, it returns a zero-value WorkflowState (convention over config).
func LoadState(specsDir string) (WorkflowState, error) {
	// Primary location: .mysd/STATE.json
	statePath := filepath.Join(MysdDir(specsDir), "STATE.json")
	data, err := os.ReadFile(statePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// Fallback: legacy location specsDir/STATE.json
		statePath = filepath.Join(specsDir, "STATE.json")
		data, err = os.ReadFile(statePath)
	}
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return WorkflowState{}, nil
		}
		return WorkflowState{}, err
	}

	var ws WorkflowState
	if err := json.Unmarshal(data, &ws); err != nil {
		return WorkflowState{}, err
	}
	return ws, nil
}

// SaveState writes WorkflowState to .mysd/STATE.json (derived from specsDir).
// It creates the .mysd directory if it does not exist.
func SaveState(specsDir string, ws WorkflowState) error {
	dir := MysdDir(specsDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return err
	}

	statePath := filepath.Join(dir, "STATE.json")
	return os.WriteFile(statePath, data, 0644)
}
