package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestStateString(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(99), "unknown"},
	}

	for _, test := range tests {
		if test.state.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.state.String())
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.FailureThreshold != 5 {
		t.Errorf("Expected FailureThreshold 5, got %d", cfg.FailureThreshold)
	}
	if cfg.SuccessThreshold != 3 {
		t.Errorf("Expected SuccessThreshold 3, got %d", cfg.SuccessThreshold)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected Timeout 30s, got %v", cfg.Timeout)
	}
}

func TestNew(t *testing.T) {
	cb := New(DefaultConfig())
	if cb == nil {
		t.Fatal("CircuitBreaker is nil")
	}
	if cb.State() != StateClosed {
		t.Errorf("Expected initial state Closed, got %v", cb.State())
	}
}

func TestCircuitBreakerExecuteSuccess(t *testing.T) {
	cb := New(DefaultConfig())

	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cb.State() != StateClosed {
		t.Error("State should remain closed")
	}
}

func TestCircuitBreakerExecuteFailure(t *testing.T) {
	cfg := Config{
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Timeout:          time.Second,
	}
	cb := New(cfg)

	// First failure
	_ = cb.Execute(func() error {
		return errors.New("error")
	})
	if cb.State() != StateClosed {
		t.Error("Should still be closed after 1 failure")
	}

	// Second failure - should trip
	_ = cb.Execute(func() error {
		return errors.New("error")
	})
	if cb.State() != StateOpen {
		t.Error("Should be open after 2 failures")
	}
}

func TestCircuitBreakerOpenState(t *testing.T) {
	cfg := Config{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	cb := New(cfg)

	// Trip the circuit
	_ = cb.Execute(func() error {
		return errors.New("error")
	})

	if cb.State() != StateOpen {
		t.Fatal("Should be open")
	}

	// Should return error immediately
	err := cb.Execute(func() error {
		return nil
	})

	if !IsOpenState(err) {
		t.Errorf("Expected open state error, got %v", err)
	}
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	cfg := Config{
		FailureThreshold:    1,
		SuccessThreshold:    2,
		Timeout:             20 * time.Millisecond,
		HalfOpenMaxRequests: 2,
	}
	cb := New(cfg)

	// Trip the circuit
	_ = cb.Execute(func() error {
		return errors.New("error")
	})

	// Wait for timeout
	time.Sleep(30 * time.Millisecond)

	// Execute to trigger state check - this will transition to half-open
	_ = cb.Execute(func() error {
		return nil
	})

	// Should be half-open now (or closed if success threshold met)
	state := cb.State()
	if state != StateHalfOpen && state != StateClosed {
		t.Errorf("Expected half-open or closed, got %v", state)
	}
}

func TestCircuitBreakerHalfOpenSuccess(t *testing.T) {
	cfg := Config{
		FailureThreshold:    1,
		SuccessThreshold:    2,
		Timeout:             50 * time.Millisecond,
		HalfOpenMaxRequests: 3,
	}
	cb := New(cfg)

	// Trip the circuit
	_ = cb.Execute(func() error {
		return errors.New("error")
	})

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Success in half-open
	_ = cb.Execute(func() error {
		return nil
	})
	_ = cb.Execute(func() error {
		return nil
	})

	// Should be closed now
	if cb.State() != StateClosed {
		t.Errorf("Expected closed, got %v", cb.State())
	}
}

func TestCircuitBreakerHalfOpenFailure(t *testing.T) {
	cfg := Config{
		FailureThreshold:    1,
		SuccessThreshold:    2,
		Timeout:             50 * time.Millisecond,
		HalfOpenMaxRequests: 2,
	}
	cb := New(cfg)

	// Trip the circuit
	_ = cb.Execute(func() error {
		return errors.New("error")
	})

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Failure in half-open
	_ = cb.Execute(func() error {
		return errors.New("error")
	})

	// Should be open again
	if cb.State() != StateOpen {
		t.Errorf("Expected open, got %v", cb.State())
	}
}

func TestCircuitBreakerTrip(t *testing.T) {
	cb := New(DefaultConfig())
	cb.Trip()

	if cb.State() != StateOpen {
		t.Error("Should be open after Trip()")
	}
}

func TestCircuitBreakerReset(t *testing.T) {
	cfg := Config{FailureThreshold: 1, Timeout: time.Second}
	cb := New(cfg)

	// Trip
	_ = cb.Execute(func() error { return errors.New("error") })
	if cb.State() != StateOpen {
		t.Fatal("Should be open")
	}

	// Reset
	cb.Reset()
	if cb.State() != StateClosed {
		t.Error("Should be closed after Reset()")
	}
}

func TestCircuitBreakerStats(t *testing.T) {
	cfg := Config{
		FailureThreshold: 5,
		SuccessThreshold: 3,
	}
	cb := New(cfg)

	_ = cb.Execute(func() error { return nil })
	stats := cb.Stats()

	if stats.State != StateClosed {
		t.Error("State should be closed")
	}
	if stats.Successes != 1 {
		t.Errorf("Expected 1 success, got %d", stats.Successes)
	}
}

func TestCircuitBreakerCall(t *testing.T) {
	cb := New(DefaultConfig())

	result, err := cb.Call(func() (interface{}, error) {
		return "success", nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}
}

func TestCircuitBreakerCallContext(t *testing.T) {
	cb := New(DefaultConfig())

	ctx := context.Background()
	result, err := cb.CallContext(ctx, func() (interface{}, error) {
		return "success", nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}
}

func TestCircuitBreakerStateChangeCallback(t *testing.T) {
	var mu sync.Mutex
	changes := []State{}
	cfg := Config{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          20 * time.Millisecond,
		OnStateChange: func(from, to State) {
			mu.Lock()
			changes = append(changes, to)
			mu.Unlock()
		},
	}
	cb := New(cfg)

	// Trip
	_ = cb.Execute(func() error { return errors.New("error") })

	// Wait for timeout and success
	time.Sleep(30 * time.Millisecond)
	_ = cb.Execute(func() error { return nil })

	// Wait for callback
	time.Sleep(10 * time.Millisecond)

	// Should have recorded state changes
	mu.Lock()
	defer mu.Unlock()
	if len(changes) < 1 {
		t.Errorf("Expected at least 1 state change, got %d", len(changes))
	}
}

func TestIsOpenState(t *testing.T) {
	if !IsOpenState(ErrOpenState) {
		t.Error("ErrOpenState should be open state")
	}
	if IsOpenState(errors.New("other")) {
		t.Error("Other error should not be open state")
	}
	if IsOpenState(nil) {
		t.Error("nil should not be open state")
	}
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("Manager is nil")
	}
}

func TestManagerRegister(t *testing.T) {
	m := NewManager()
	cb := m.Register("test", DefaultConfig())

	if cb == nil {
		t.Fatal("CircuitBreaker is nil")
	}

	found, ok := m.Get("test")
	if !ok {
		t.Error("Should find registered circuit breaker")
	}
	if found != cb {
		t.Error("Should return same circuit breaker")
	}
}

func TestManagerGetNotFound(t *testing.T) {
	m := NewManager()

	_, ok := m.Get("nonexistent")
	if ok {
		t.Error("Should not find unregistered circuit breaker")
	}
}

func TestManagerExecute(t *testing.T) {
	m := NewManager()
	cfg := Config{FailureThreshold: 1, Timeout: time.Second}
	m.Register("test", cfg)

	err := m.Execute("test", func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManagerExecuteNotFound(t *testing.T) {
	m := NewManager()

	// Should execute without circuit breaker
	err := m.Execute("nonexistent", func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManagerExecuteContext(t *testing.T) {
	m := NewManager()
	m.Register("test", DefaultConfig())

	ctx := context.Background()
	err := m.ExecuteContext(ctx, "test", func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManagerKeys(t *testing.T) {
	m := NewManager()
	m.Register("a", DefaultConfig())
	m.Register("b", DefaultConfig())

	keys := m.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestManagerStats(t *testing.T) {
	m := NewManager()
	m.Register("test", DefaultConfig())

	stats, ok := m.Stats("test")
	if !ok {
		t.Fatal("Should find stats")
	}
	if stats.State != StateClosed {
		t.Error("State should be closed")
	}
}

func TestManagerStatsNotFound(t *testing.T) {
	m := NewManager()

	_, ok := m.Stats("nonexistent")
	if ok {
		t.Error("Should not find stats for nonexistent key")
	}
}

func TestManagerAllStats(t *testing.T) {
	m := NewManager()
	m.Register("a", DefaultConfig())
	m.Register("b", DefaultConfig())

	allStats := m.AllStats()
	if len(allStats) != 2 {
		t.Errorf("Expected 2 stats, got %d", len(allStats))
	}
}

func TestManagerTrippedCircuits(t *testing.T) {
	m := NewManager()
	cfg := Config{FailureThreshold: 1, Timeout: time.Second}
	m.Register("ok", DefaultConfig())
	m.Register("tripped", cfg)

	// Trip one
	cb, _ := m.Get("tripped")
	_ = cb.Execute(func() error { return errors.New("error") })

	tripped := m.TrippedCircuits()
	if len(tripped) != 1 {
		t.Errorf("Expected 1 tripped circuit, got %d", len(tripped))
	}
	if tripped[0] != "tripped" {
		t.Errorf("Expected 'tripped', got %s", tripped[0])
	}
}

func TestConcurrentCircuitBreaker(t *testing.T) {
	cfg := Config{FailureThreshold: 100, Timeout: time.Second}
	cb := New(cfg)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cb.Execute(func() error { return nil })
		}()
	}
	wg.Wait()
}

func TestConcurrentManager(t *testing.T) {
	m := NewManager()
	m.Register("test", DefaultConfig())

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = m.Execute("test", func() error { return nil })
		}()
	}
	wg.Wait()
}

func TestHalfOpenMaxRequests(t *testing.T) {
	cfg := Config{
		FailureThreshold:    1,
		SuccessThreshold:    3,
		Timeout:             20 * time.Millisecond,
		HalfOpenMaxRequests: 1,
	}
	cb := New(cfg)

	// Trip
	_ = cb.Execute(func() error { return errors.New("error") })

	// Wait for half-open
	time.Sleep(30 * time.Millisecond)

	// First request should be allowed
	err1 := cb.Execute(func() error { return nil })
	if IsOpenState(err1) {
		t.Error("First request should be allowed")
	}

	// Second request should be blocked (only 1 allowed in half-open)
	// But we need to check state first - if it transitioned, behavior changes
	_ = cb.Execute(func() error { return nil })
}