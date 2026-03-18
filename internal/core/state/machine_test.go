package state

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestStateString(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateNone, "none"},
		{StateCreated, "created"},
		{StateInitializing, "initializing"},
		{StateInitialized, "initialized"},
		{StateStarting, "starting"},
		{StateRunning, "running"},
		{StateStopping, "stopping"},
		{StateStopped, "stopped"},
		{StateError, "error"},
		{StateDestroyed, "destroyed"},
	}

	for _, test := range tests {
		if test.state.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.state.String())
		}
	}
}

func TestCanTransition(t *testing.T) {
	// Valid transitions
	validTransitions := []struct {
		from, to State
	}{
		{StateNone, StateCreated},
		{StateCreated, StateInitializing},
		{StateInitializing, StateInitialized},
		{StateInitialized, StateStarting},
		{StateStarting, StateRunning},
		{StateRunning, StateStopping},
		{StateStopping, StateStopped},
		{StateStopped, StateStarting},
		{StateError, StateInitializing},
	}

	for _, test := range validTransitions {
		if !CanTransition(test.from, test.to) {
			t.Errorf("Expected valid transition from %s to %s", test.from, test.to)
		}
	}

	// Invalid transitions
	invalidTransitions := []struct {
		from, to State
	}{
		{StateNone, StateRunning},
		{StateRunning, StateInitialized},
		{StateDestroyed, StateRunning},
		{StateCreated, StateRunning},
	}

	for _, test := range invalidTransitions {
		if CanTransition(test.from, test.to) {
			t.Errorf("Expected invalid transition from %s to %s", test.from, test.to)
		}
	}
}

func TestNewStateMachine(t *testing.T) {
	sm := NewStateMachine()
	if sm == nil {
		t.Fatal("StateMachine is nil")
	}

	if sm.Current() != StateNone {
		t.Errorf("Expected StateNone, got %v", sm.Current())
	}
}

func TestStateMachineTransition(t *testing.T) {
	sm := NewStateMachine()

	// None -> Created
	err := sm.Transition(context.Background(), StateCreated)
	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	if sm.Current() != StateCreated {
		t.Errorf("Expected StateCreated, got %v", sm.Current())
	}

	if sm.Previous() != StateNone {
		t.Errorf("Expected Previous to be StateNone, got %v", sm.Previous())
	}
}

func TestStateMachineInvalidTransition(t *testing.T) {
	sm := NewStateMachine()

	// Try to skip to running
	err := sm.Transition(context.Background(), StateRunning)
	if err == nil {
		t.Error("Expected error for invalid transition")
	}
}

func TestStateMachineCallbacks(t *testing.T) {
	sm := NewStateMachine()

	callbackCalled := false
	sm.OnEnter(StateCreated, func(ctx context.Context, from, to State) error {
		callbackCalled = true
		return nil
	})

	err := sm.Transition(context.Background(), StateCreated)
	if err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	if !callbackCalled {
		t.Error("Callback was not called")
	}
}

func TestStateMachineCallbackError(t *testing.T) {
	sm := NewStateMachine()

	sm.OnEnter(StateCreated, func(ctx context.Context, from, to State) error {
		return errors.New("callback error")
	})

	err := sm.Transition(context.Background(), StateCreated)
	if err == nil {
		t.Error("Expected error from callback")
	}

	// Should transition to error state
	if sm.Current() != StateError {
		t.Errorf("Expected StateError, got %v", sm.Current())
	}
}

func TestStateMachineHistory(t *testing.T) {
	sm := NewStateMachine()

	sm.Transition(context.Background(), StateCreated)
	sm.Transition(context.Background(), StateInitializing)
	sm.Transition(context.Background(), StateInitialized)

	history := sm.History()
	if len(history) != 3 {
		t.Errorf("Expected 3 transitions, got %d", len(history))
	}
}

func TestStateMachineCanGoTo(t *testing.T) {
	sm := NewStateMachine()

	if !sm.CanGoTo(StateCreated) {
		t.Error("Should be able to go to StateCreated")
	}

	if sm.CanGoTo(StateRunning) {
		t.Error("Should not be able to go to StateRunning")
	}
}

func TestStateMachineIsRunning(t *testing.T) {
	sm := NewStateMachine()

	if sm.IsRunning() {
		t.Error("Should not be running initially")
	}

	sm.Transition(context.Background(), StateCreated)
	sm.Transition(context.Background(), StateInitializing)
	sm.Transition(context.Background(), StateInitialized)
	sm.Transition(context.Background(), StateStarting)
	sm.Transition(context.Background(), StateRunning)

	if !sm.IsRunning() {
		t.Error("Should be running")
	}
}

func TestStateMachineIsError(t *testing.T) {
	sm := NewStateMachine()

	if sm.IsError() {
		t.Error("Should not be in error initially")
	}

	sm.OnEnter(StateCreated, func(ctx context.Context, from, to State) error {
		return errors.New("force error")
	})
	sm.Transition(context.Background(), StateCreated)

	if !sm.IsError() {
		t.Error("Should be in error state")
	}
}

func TestStateMachineIsTerminal(t *testing.T) {
	sm := NewStateMachine()

	if sm.IsTerminal() {
		t.Error("Should not be terminal initially")
	}

	sm.Transition(context.Background(), StateCreated)
	sm.Transition(context.Background(), StateDestroyed)

	if !sm.IsTerminal() {
		t.Error("Destroyed should be terminal")
	}
}

func TestStateMachineDuration(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(context.Background(), StateCreated)

	time.Sleep(10 * time.Millisecond)

	if sm.Duration() < 10*time.Millisecond {
		t.Error("Duration should be at least 10ms")
	}
}

func TestStateMachineEnteredAt(t *testing.T) {
	sm := NewStateMachine()
	before := time.Now()
	sm.Transition(context.Background(), StateCreated)
	after := time.Now()

	enteredAt := sm.EnteredAt()
	if enteredAt.Before(before) || enteredAt.After(after) {
		t.Error("EnteredAt should be between before and after")
	}
}

func TestStateMachineReset(t *testing.T) {
	sm := NewStateMachine()

	sm.Transition(context.Background(), StateCreated)
	sm.Transition(context.Background(), StateInitializing)
	sm.Reset()

	if sm.Current() != StateNone {
		t.Errorf("Expected StateNone after reset, got %v", sm.Current())
	}

	if len(sm.History()) != 0 {
		t.Error("History should be empty after reset")
	}
}

func TestStateMachineSetMaxHistory(t *testing.T) {
	sm := NewStateMachine()
	sm.SetMaxHistory(2)

	sm.Transition(context.Background(), StateCreated)
	sm.Transition(context.Background(), StateInitializing)
	sm.Transition(context.Background(), StateInitialized)

	history := sm.History()
	if len(history) > 2 {
		t.Errorf("History should be limited to 2, got %d", len(history))
	}
}

func TestStateMachineForceTransition(t *testing.T) {
	sm := NewStateMachine()

	// Force transition (skip validation)
	sm.ForceTransition(StateRunning)

	if sm.Current() != StateRunning {
		t.Errorf("Expected StateRunning, got %v", sm.Current())
	}
}

func TestStateMachineStats(t *testing.T) {
	sm := NewStateMachine()

	sm.Transition(context.Background(), StateCreated)
	sm.Transition(context.Background(), StateInitializing)
	sm.Transition(context.Background(), StateInitialized)

	stats := sm.Stats()

	if stats.CurrentState != StateInitialized {
		t.Errorf("Expected StateInitialized, got %v", stats.CurrentState)
	}

	if stats.TransitionCount != 3 {
		t.Errorf("Expected 3 transitions, got %d", stats.TransitionCount)
	}
}

func TestStateMachineFullLifecycle(t *testing.T) {
	sm := NewStateMachine()
	ctx := context.Background()

	// Full lifecycle
	steps := []State{
		StateCreated,
		StateInitializing,
		StateInitialized,
		StateStarting,
		StateRunning,
		StateStopping,
		StateStopped,
	}

	for _, state := range steps {
		err := sm.Transition(ctx, state)
		if err != nil {
			t.Errorf("Transition to %s failed: %v", state, err)
		}
	}

	if sm.Current() != StateStopped {
		t.Errorf("Expected StateStopped, got %v", sm.Current())
	}
}

func TestStateMachineErrorRecovery(t *testing.T) {
	sm := NewStateMachine()
	ctx := context.Background()

	// Go through normal states
	sm.Transition(ctx, StateCreated)
	sm.Transition(ctx, StateInitializing)
	sm.Transition(ctx, StateInitialized)
	sm.Transition(ctx, StateStarting)

	// Force error
	sm.ForceTransition(StateError)

	// Can recover from error
	if !sm.CanGoTo(StateInitializing) {
		t.Error("Should be able to recover from error")
	}

	sm.Transition(ctx, StateInitializing)
	sm.Transition(ctx, StateInitialized)
	sm.Transition(ctx, StateStarting)
	sm.Transition(ctx, StateRunning)

	if sm.Current() != StateRunning {
		t.Errorf("Expected StateRunning after recovery, got %v", sm.Current())
	}
}

func TestNewTransitionError(t *testing.T) {
	err := NewTransitionError(StateRunning, StateStopped, "test error")

	if err.Current != StateRunning {
		t.Error("Current should be StateRunning")
	}

	if err.Target != StateStopped {
		t.Error("Target should be StateStopped")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}

func TestIsInvalidTransition(t *testing.T) {
	transErr := NewTransitionError(StateRunning, StateStopped, "invalid")
	if !IsInvalidTransition(transErr) {
		t.Error("Should be invalid transition error")
	}

	normalErr := errors.New("normal error")
	if IsInvalidTransition(normalErr) {
		t.Error("Should not be invalid transition error")
	}
}