package state

import (
	"fmt"
	"time"
)

// ValidTransitions defines the allowed state transitions for the spec workflow.
var ValidTransitions = map[Phase][]Phase{
	PhaseNone:     {PhaseProposed},
	PhaseProposed: {PhaseSpecced},
	PhaseSpecced:  {PhaseDesigned},
	PhaseDesigned: {PhasePlanned},
	PhasePlanned:  {PhaseExecuted},
	PhaseExecuted: {PhaseVerified},
	PhaseVerified: {PhaseArchived, PhaseExecuted}, // FAIL can re-execute
	PhaseArchived: {PhaseProposed},               // new change
}

// CanTransition returns true if transitioning from `from` to `to` is valid.
func CanTransition(from Phase, to Phase) bool {
	allowed, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == to {
			return true
		}
	}
	return false
}

// Transition updates ws to the new phase if the transition is valid.
// Returns ErrInvalidTransition (wrapped) if the transition is not allowed.
func Transition(ws *WorkflowState, to Phase) error {
	if !CanTransition(ws.Phase, to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, ws.Phase, to)
	}
	ws.Phase = to
	ws.LastRun = time.Now()
	return nil
}
