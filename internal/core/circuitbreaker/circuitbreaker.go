package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State represents the state of a circuit breaker
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Config holds circuit breaker configuration
type Config struct {
	// Failure threshold to trip the circuit
	FailureThreshold int
	// Success threshold to close the circuit
	SuccessThreshold int
	// Timeout before attempting to close
	Timeout time.Duration
	// HalfOpenMaxRequests is max requests in half-open state
	HalfOpenMaxRequests int
	// OnStateChange is called when state changes
	OnStateChange func(from, to State)
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		FailureThreshold:    5,
		SuccessThreshold:    3,
		Timeout:             30 * time.Second,
		HalfOpenMaxRequests: 1,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu          sync.Mutex
	config      Config
	state       State
	failures    int
	successes   int
	lastFailure time.Time
	requests    int
}

// New creates a new circuit breaker
func New(config Config) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// State returns the current state
func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Execute runs a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	return cb.ExecuteContext(context.Background(), fn)
}

// ExecuteContext runs a function with circuit breaker protection
func (cb *CircuitBreaker) ExecuteContext(ctx context.Context, fn func() error) error {
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	err := fn()
	cb.afterRequest(err)

	return err
}

// Call runs a function with circuit breaker protection and returns result
func (cb *CircuitBreaker) Call(fn func() (interface{}, error)) (interface{}, error) {
	return cb.CallContext(context.Background(), fn)
}

// CallContext runs a function with circuit breaker protection and returns result
func (cb *CircuitBreaker) CallContext(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	if err := cb.beforeRequest(); err != nil {
		return nil, err
	}

	result, err := fn()
	cb.afterRequest(err)

	return result, err
}

func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailure) > cb.config.Timeout {
			cb.setState(StateHalfOpen)
			cb.requests = 0
			return nil
		}
		return ErrOpenState
	case StateHalfOpen:
		if cb.requests >= cb.config.HalfOpenMaxRequests {
			return ErrOpenState
		}
		cb.requests++
		return nil
	case StateClosed:
		return nil
	}
	return nil
}

func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}
}

func (cb *CircuitBreaker) onFailure() {
	cb.failures++
	cb.lastFailure = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.config.FailureThreshold {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		cb.setState(StateOpen)
	}
}

func (cb *CircuitBreaker) onSuccess() {
	cb.failures = 0

	switch cb.state {
	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.SuccessThreshold {
			cb.setState(StateClosed)
		}
	case StateClosed:
		cb.successes++
	}
}

func (cb *CircuitBreaker) setState(state State) {
	if cb.state == state {
		return
	}
	oldState := cb.state
	cb.state = state
	cb.failures = 0
	cb.successes = 0
	cb.requests = 0

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(oldState, state)
	}
}

// Trip manually trips the circuit breaker
func (cb *CircuitBreaker) Trip() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.setState(StateOpen)
	cb.lastFailure = time.Now()
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.setState(StateClosed)
}

// Stats returns current statistics
func (cb *CircuitBreaker) Stats() Stats {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return Stats{
		State:           cb.state,
		Failures:        cb.failures,
		Successes:       cb.successes,
		FailureThreshold: cb.config.FailureThreshold,
		SuccessThreshold: cb.config.SuccessThreshold,
	}
}

// Stats holds circuit breaker statistics
type Stats struct {
	State            State
	Failures         int
	Successes        int
	FailureThreshold int
	SuccessThreshold int
}

// Errors
var (
	ErrOpenState = errors.New("circuit breaker is open")
)

// IsOpenState checks if error is open state error
func IsOpenState(err error) bool {
	return errors.Is(err, ErrOpenState)
}

// Manager manages multiple circuit breakers
type Manager struct {
	mu      sync.RWMutex
	breakers map[string]*CircuitBreaker
}

// NewManager creates a new circuit breaker manager
func NewManager() *Manager {
	return &Manager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// Register registers a circuit breaker for a key
func (m *Manager) Register(key string, config Config) *CircuitBreaker {
	m.mu.Lock()
	defer m.mu.Unlock()
	cb := New(config)
	m.breakers[key] = cb
	return cb
}

// Get gets a circuit breaker for a key
func (m *Manager) Get(key string) (*CircuitBreaker, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cb, ok := m.breakers[key]
	return cb, ok
}

// Execute executes a function with the named circuit breaker
func (m *Manager) Execute(key string, fn func() error) error {
	cb, ok := m.Get(key)
	if !ok {
		return fn()
	}
	return cb.Execute(fn)
}

// ExecuteContext executes a function with the named circuit breaker
func (m *Manager) ExecuteContext(ctx context.Context, key string, fn func() error) error {
	cb, ok := m.Get(key)
	if !ok {
		return fn()
	}
	return cb.ExecuteContext(ctx, fn)
}

// Keys returns all registered keys
func (m *Manager) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.breakers))
	for k := range m.breakers {
		keys = append(keys, k)
	}
	return keys
}

// Stats returns stats for a key
func (m *Manager) Stats(key string) (Stats, bool) {
	cb, ok := m.Get(key)
	if !ok {
		return Stats{}, false
	}
	return cb.Stats(), true
}

// AllStats returns stats for all circuit breakers
func (m *Manager) AllStats() map[string]Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]Stats)
	for k, cb := range m.breakers {
		result[k] = cb.Stats()
	}
	return result
}

// TrippedCircuits returns keys of all tripped circuits
func (m *Manager) TrippedCircuits() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]string, 0)
	for k, cb := range m.breakers {
		if cb.State() == StateOpen {
			result = append(result, k)
		}
	}
	return result
}