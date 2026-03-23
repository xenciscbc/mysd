package state

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanTransition_ValidTransitions(t *testing.T) {
	tests := []struct {
		from Phase
		to   Phase
	}{
		{PhaseNone, PhaseProposed},
		{PhaseProposed, PhaseSpecced},
		{PhaseSpecced, PhaseDesigned},
		{PhaseDesigned, PhasePlanned},
		{PhasePlanned, PhaseExecuted},
		{PhaseExecuted, PhaseVerified},
		{PhaseVerified, PhaseArchived},
		{PhaseVerified, PhaseExecuted}, // re-execute after verify fail
		{PhaseArchived, PhaseProposed}, // new change
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			assert.True(t, CanTransition(tt.from, tt.to), "expected valid transition: %s -> %s", tt.from, tt.to)
		})
	}
}

func TestCanTransition_InvalidTransitions(t *testing.T) {
	tests := []struct {
		from Phase
		to   Phase
		desc string
	}{
		{PhaseProposed, PhaseExecuted, "skip not allowed"},
		{PhaseProposed, PhaseVerified, "skip multiple phases"},
		{PhaseArchived, PhaseExecuted, "archived cannot go to executed"},
		{PhaseNone, PhaseSpecced, "none cannot skip to specced"},
		{PhaseExecuted, PhaseProposed, "reverse not allowed"},
		{PhaseVerified, PhaseProposed, "reverse to proposed not allowed from verified"},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.False(t, CanTransition(tt.from, tt.to), "expected invalid transition: %s -> %s (%s)", tt.from, tt.to, tt.desc)
		})
	}
}

func TestTransition_UpdatesPhaseAndLastRun(t *testing.T) {
	ws := &WorkflowState{
		ChangeName: "test",
		Phase:      PhaseNone,
		LastRun:    time.Time{},
	}
	before := time.Now()
	err := Transition(ws, PhaseProposed)
	after := time.Now()

	require.NoError(t, err)
	assert.Equal(t, PhaseProposed, ws.Phase, "Phase should be updated")
	assert.True(t, ws.LastRun.After(before) || ws.LastRun.Equal(before), "LastRun should be updated")
	assert.True(t, ws.LastRun.Before(after) || ws.LastRun.Equal(after), "LastRun should not be in the future")
}

func TestTransition_InvalidReturnsErrInvalidTransition(t *testing.T) {
	ws := &WorkflowState{
		ChangeName: "test",
		Phase:      PhaseProposed,
	}
	err := Transition(ws, PhaseExecuted)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidTransition), "error should wrap ErrInvalidTransition")
	assert.Equal(t, PhaseProposed, ws.Phase, "Phase should not change on invalid transition")
}

func TestTransition_ValidTransitionsAllPhases(t *testing.T) {
	ws := &WorkflowState{Phase: PhaseNone}

	phases := []Phase{PhaseProposed, PhaseSpecced, PhaseDesigned, PhasePlanned, PhaseExecuted, PhaseVerified, PhaseArchived}
	for _, phase := range phases {
		err := Transition(ws, phase)
		require.NoError(t, err, "transition to %s should succeed", phase)
		assert.Equal(t, phase, ws.Phase)
	}
}
