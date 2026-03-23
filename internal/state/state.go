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

// LoadState reads WorkflowState from specsDir/STATE.json.
// If the file does not exist, it returns a zero-value WorkflowState (convention over config).
func LoadState(specsDir string) (WorkflowState, error) {
	statePath := filepath.Join(specsDir, "STATE.json")
	data, err := os.ReadFile(statePath)
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

// SaveState writes WorkflowState to specsDir/STATE.json.
// It creates specsDir if it does not exist.
func SaveState(specsDir string, ws WorkflowState) error {
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return err
	}

	statePath := filepath.Join(specsDir, "STATE.json")
	return os.WriteFile(statePath, data, 0644)
}
