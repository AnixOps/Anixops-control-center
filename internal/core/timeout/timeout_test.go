package timeout

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("Manager is nil")
	}
}

func TestManagerSetGetDefault(t *testing.T) {
	m := NewManager()
	m.SetDefault("test", 5*time.Second)

	d, ok := m.GetDefault("test")
	if !ok {
		t.Error("Should find timeout")
	}
	if d != 5*time.Second {
		t.Errorf("Expected 5s, got %v", d)
	}
}

func TestManagerGetDefaultNotFound(t *testing.T) {
	m := NewManager()

	_, ok := m.GetDefault("nonexistent")
	if ok {
		t.Error("Should not find nonexistent timeout")
	}
}

func TestManagerContext(t *testing.T) {
	m := NewManager()
	m.SetDefault("test", 5*time.Second)

	ctx, cancel := m.Context(context.Background(), "test")
	defer cancel()

	if ctx == nil {
		t.Error("Context should not be nil")
	}

	// Check deadline is set
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Deadline should be set")
	}
	if time.Until(deadline) > 5*time.Second {
		t.Error("Deadline should be within 5 seconds")
	}
}

func TestManagerContextNotFound(t *testing.T) {
	m := NewManager()

	ctx, cancel := m.Context(context.Background(), "nonexistent")
	defer cancel()

	// Should return original context without deadline
	_, ok := ctx.Deadline()
	if ok {
		t.Error("Should not have deadline for nonexistent operation")
	}
}

func TestManagerContextWithFallback(t *testing.T) {
	m := NewManager()

	ctx, cancel := m.ContextWithFallback(context.Background(), "nonexistent", 3*time.Second)
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Deadline should be set with fallback")
	}
	if time.Until(deadline) > 3*time.Second {
		t.Error("Deadline should be within fallback duration")
	}
}

func TestManagerContextWithFallbackExisting(t *testing.T) {
	m := NewManager()
	m.SetDefault("test", 5*time.Second)

	ctx, cancel := m.ContextWithFallback(context.Background(), "test", 3*time.Second)
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Deadline should be set")
	}
	// Should use the registered timeout, not fallback
	if time.Until(deadline) > 5*time.Second {
		t.Error("Should use registered timeout, not fallback")
	}
}

func TestManagerExecute(t *testing.T) {
	m := NewManager()
	m.SetDefault("test", 5*time.Second)

	err := m.Execute(context.Background(), "test", func(ctx context.Context) error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManagerExecuteTimeout(t *testing.T) {
	m := NewManager()
	m.SetDefault("test", 50*time.Millisecond)

	err := m.Execute(context.Background(), "test", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			return nil
		}
	})

	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestRun(t *testing.T) {
	result := Run(context.Background(), 5*time.Second, func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Millisecond) // Ensure measurable duration
		return "success", nil
	})

	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}
	if result.Value != "success" {
		t.Errorf("Expected 'success', got %v", result.Value)
	}
	if result.TimedOut {
		t.Error("Should not be timed out")
	}
	if result.Duration <= 0 {
		t.Errorf("Duration should be positive, got %v", result.Duration)
	}
}

func TestRunTimeout(t *testing.T) {
	result := Run(context.Background(), 50*time.Millisecond, func(ctx context.Context) (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
			return "success", nil
		}
	})

	if !result.TimedOut {
		t.Error("Should be timed out")
	}
	if result.Error == nil {
		t.Error("Should have error")
	}
}

func TestRunError(t *testing.T) {
	expectedErr := errors.New("test error")
	result := Run(context.Background(), 5*time.Second, func(ctx context.Context) (interface{}, error) {
		return nil, expectedErr
	})

	if result.Error != expectedErr {
		t.Errorf("Expected test error, got %v", result.Error)
	}
	if result.TimedOut {
		t.Error("Should not be timed out")
	}
}

func TestNewPool(t *testing.T) {
	p := NewPool()
	if p == nil {
		t.Fatal("Pool is nil")
	}
}

func TestPoolRegister(t *testing.T) {
	p := NewPool()
	p.Register("test", 5*time.Second)

	d, ok := p.Get("test")
	if !ok {
		t.Error("Should find timeout")
	}
	if d != 5*time.Second {
		t.Errorf("Expected 5s, got %v", d)
	}
}

func TestPoolUnregister(t *testing.T) {
	p := NewPool()
	p.Register("test", 5*time.Second)
	p.Unregister("test")

	_, ok := p.Get("test")
	if ok {
		t.Error("Should not find unregistered timeout")
	}
}

func TestPoolContext(t *testing.T) {
	p := NewPool()
	p.Register("test", 5*time.Second)

	ctx, cancel, ok := p.Context(context.Background(), "test")
	defer cancel()

	if !ok {
		t.Error("Should find registered timeout")
	}
	if ctx == nil {
		t.Error("Context should not be nil")
	}
}

func TestPoolContextNotFound(t *testing.T) {
	p := NewPool()

	ctx, _, ok := p.Context(context.Background(), "nonexistent")

	if ok {
		t.Error("Should not find unregistered timeout")
	}
	_, hasDeadline := ctx.Deadline()
	if hasDeadline {
		t.Error("Should not have deadline")
	}
}

func TestPoolKeys(t *testing.T) {
	p := NewPool()
	p.Register("a", time.Second)
	p.Register("b", time.Second)
	p.Register("c", time.Second)

	keys := p.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
}

func TestFromConfig(t *testing.T) {
	cfg := Config{
		Default:    30 * time.Second,
		Operations: map[string]time.Duration{"test": 5 * time.Second},
	}

	m := FromConfig(cfg)
	if m == nil {
		t.Fatal("Manager is nil")
	}

	d, ok := m.GetDefault("test")
	if !ok {
		t.Error("Should find timeout from config")
	}
	if d != 5*time.Second {
		t.Errorf("Expected 5s, got %v", d)
	}
}

func TestTimeoutError(t *testing.T) {
	err := &TimeoutError{
		Operation: "test",
		Duration:  5 * time.Second,
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}

func TestIsTimeout(t *testing.T) {
	// TimeoutError
	if !IsTimeout(&TimeoutError{Operation: "test", Duration: time.Second}) {
		t.Error("TimeoutError should be timeout")
	}

	// context.DeadlineExceeded
	if !IsTimeout(context.DeadlineExceeded) {
		t.Error("DeadlineExceeded should be timeout")
	}

	// Other error
	if IsTimeout(errors.New("other error")) {
		t.Error("Other error should not be timeout")
	}

	// nil
	if IsTimeout(nil) {
		t.Error("nil should not be timeout")
	}
}

func TestConcurrentManager(t *testing.T) {
	m := NewManager()

	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(i int) {
			m.SetDefault("test", time.Duration(i)*time.Millisecond)
			m.GetDefault("test")
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestConcurrentPool(t *testing.T) {
	p := NewPool()

	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(i int) {
			p.Register("test", time.Duration(i)*time.Millisecond)
			p.Get("test")
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}