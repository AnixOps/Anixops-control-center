package state

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents a plugin state
type State int

const (
	StateNone State = iota
	StateCreated
	StateInitializing
	StateInitialized
	StateStarting
	StateRunning
	StateStopping
	StateStopped
	StateError
	StateDestroyed
)

func (s State) String() string {
	switch s {
	case StateNone:
		return "none"
	case StateCreated:
		return "created"
	case StateInitializing:
		return "initializing"
	case StateInitialized:
		return "initialized"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateStopping:
		return "stopping"
	case StateStopped:
		return "stopped"
	case StateError:
		return "error"
	case StateDestroyed:
		return "destroyed"
	default:
		return "unknown"
	}
}

// Transition represents a state transition
type Transition struct {
	From State
	To   State
}

// ValidTransitions defines allowed state transitions
var ValidTransitions = map[State][]State{
	StateNone:         {StateCreated},
	StateCreated:      {StateInitializing, StateDestroyed},
	StateInitializing: {StateInitialized, StateError},
	StateInitialized:  {StateStarting, StateStopped, StateDestroyed},
	StateStarting:     {StateRunning, StateError},
	StateRunning:      {StateStopping, StateError},
	StateStopping:     {StateStopped, StateError},
	StateStopped:      {StateStarting, StateDestroyed},
	StateError:        {StateInitializing, StateStarting, StateStopping, StateDestroyed},
	StateDestroyed:    {},
}

// CanTransition checks if a transition is valid
func CanTransition(from, to State) bool {
	allowed, exists := ValidTransitions[from]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// StateMachine manages state transitions
type StateMachine struct {
	mu          sync.RWMutex
	current     State
	previous    State
	enteredAt   time.Time
	transitions []TransitionRecord
	maxHistory  int
	callbacks   map[State][]StateCallback
}

// StateCallback is called on state entry
type StateCallback func(ctx context.Context, from, to State) error

// TransitionRecord records a state transition
type TransitionRecord struct {
	From      State
	To        State
	Timestamp time.Time
	Duration  time.Duration
	Error     error
}

// NewStateMachine creates a new state machine
func NewStateMachine() *StateMachine {
	return &StateMachine{
		current:    StateNone,
		previous:   StateNone,
		maxHistory: 100,
		callbacks:  make(map[State][]StateCallback),
	}
}

// Current returns the current state
func (sm *StateMachine) Current() State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current
}

// Previous returns the previous state
func (sm *StateMachine) Previous() State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.previous
}

// EnteredAt returns when the current state was entered
func (sm *StateMachine) EnteredAt() time.Time {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.enteredAt
}

// Duration returns how long in the current state
func (sm *StateMachine) Duration() time.Duration {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return time.Since(sm.enteredAt)
}

// OnEnter registers a callback for entering a state
func (sm *StateMachine) OnEnter(state State, callback StateCallback) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.callbacks[state] = append(sm.callbacks[state], callback)
}

// Transition attempts a state transition
func (sm *StateMachine) Transition(ctx context.Context, to State) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	from := sm.current

	// Validate transition
	if !CanTransition(from, to) {
		return fmt.Errorf("invalid transition from %s to %s", from, to)
	}

	// Record transition start
	start := time.Now()

	// Execute callbacks
	var callbackErr error
	for _, callback := range sm.callbacks[to] {
		if err := callback(ctx, from, to); err != nil {
			callbackErr = err
			break
		}
	}

	// Record transition
	record := TransitionRecord{
		From:      from,
		To:        to,
		Timestamp: start,
		Duration:  time.Since(start),
	}

	if callbackErr != nil {
		record.Error = callbackErr
		sm.transitions = append(sm.transitions, record)
		// Transition to error state
		sm.previous = from
		sm.current = StateError
		sm.enteredAt = time.Now()
		return fmt.Errorf("transition callback failed: %w", callbackErr)
	}

	// Update state
	sm.transitions = append(sm.transitions, record)
	sm.previous = from
	sm.current = to
	sm.enteredAt = time.Now()

	// Trim history if needed
	if len(sm.transitions) > sm.maxHistory {
		sm.transitions = sm.transitions[len(sm.transitions)-sm.maxHistory:]
	}

	return nil
}

// ForceTransition forces a state transition (use with caution)
func (sm *StateMachine) ForceTransition(to State) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.previous = sm.current
	sm.current = to
	sm.enteredAt = time.Now()
}

// History returns the transition history
func (sm *StateMachine) History() []TransitionRecord {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	result := make([]TransitionRecord, len(sm.transitions))
	copy(result, sm.transitions)
	return result
}

// CanGoTo checks if the machine can transition to a state
func (sm *StateMachine) CanGoTo(to State) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return CanTransition(sm.current, to)
}

// IsRunning returns true if in running state
func (sm *StateMachine) IsRunning() bool {
	return sm.Current() == StateRunning
}

// IsError returns true if in error state
func (sm *StateMachine) IsError() bool {
	return sm.Current() == StateError
}

// IsTerminal returns true if in a terminal state
func (sm *StateMachine) IsTerminal() bool {
	return sm.Current() == StateDestroyed
}

// Reset resets the state machine
func (sm *StateMachine) Reset() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.current = StateNone
	sm.previous = StateNone
	sm.enteredAt = time.Time{}
	sm.transitions = nil
}

// SetMaxHistory sets the maximum history size
func (sm *StateMachine) SetMaxHistory(max int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.maxHistory = max
}

// Stats returns state machine statistics
func (sm *StateMachine) Stats() TransitionStats {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats := TransitionStats{
		CurrentState:     sm.current,
		PreviousState:    sm.previous,
		TimeInState:      time.Since(sm.enteredAt),
		TransitionCount:  len(sm.transitions),
		TransitionCounts: make(map[State]int),
	}

	for _, t := range sm.transitions {
		stats.TransitionCounts[t.To]++
	}

	return stats
}

// TransitionStats contains state machine statistics
type TransitionStats struct {
	CurrentState     State
	PreviousState    State
	TimeInState      time.Duration
	TransitionCount  int
	TransitionCounts map[State]int
}

// TransitionError represents a state machine error
type TransitionError struct {
	Current State
	Target  State
	Message string
}

func (e *TransitionError) Error() string {
	return fmt.Sprintf("state error: %s (current: %s, target: %s)", e.Message, e.Current, e.Target)
}

// NewTransitionError creates a new transition error
func NewTransitionError(current, target State, message string) *TransitionError {
	return &TransitionError{
		Current: current,
		Target:  target,
		Message: message,
	}
}

// IsInvalidTransition checks if an error is an invalid transition error
func IsInvalidTransition(err error) bool {
	var transErr *TransitionError
	if errors.As(err, &transErr) {
		return true
	}
	return false
}