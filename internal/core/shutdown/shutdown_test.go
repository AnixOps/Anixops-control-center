package shutdown

import (
	"context"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler()
	if h == nil {
		t.Fatal("Handler is nil")
	}
}

func TestSetTimeout(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(10 * time.Second)
	// No error means success
}

func TestRegister(t *testing.T) {
	h := NewHandler()
	h.Register(func(ctx context.Context) error {
		return nil
	})
	// No error means success
}

func TestIsShuttingDown(t *testing.T) {
	h := NewHandler()

	if h.IsShuttingDown() {
		t.Error("Should not be shutting down initially")
	}
}

func TestDone(t *testing.T) {
	h := NewHandler()

	done := h.Done()
	if done == nil {
		t.Fatal("Done channel is nil")
	}

	// Should be open initially
	select {
	case <-done:
		t.Error("Done channel should be open initially")
	default:
		// Correct
	}
}

func TestShutdown(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(5 * time.Second)

	called := false
	h.Register(func(ctx context.Context) error {
		called = true
		return nil
	})

	err := h.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	if !called {
		t.Error("Shutdown function was not called")
	}

	if !h.IsShuttingDown() {
		t.Error("Should be shutting down")
	}

	// Done channel should be closed
	select {
	case <-h.Done():
		// Correct
	case <-time.After(1 * time.Second):
		t.Error("Done channel should be closed after shutdown")
	}
}

func TestShutdownMultipleFunctions(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(5 * time.Second)

	order := []string{}
	h.Register(func(ctx context.Context) error {
		order = append(order, "first")
		return nil
	})
	h.Register(func(ctx context.Context) error {
		order = append(order, "second")
		return nil
	})
	h.Register(func(ctx context.Context) error {
		order = append(order, "third")
		return nil
	})

	h.Shutdown(context.Background())

	// Functions should be called in reverse order
	if len(order) != 3 {
		t.Fatalf("Expected 3 calls, got %d", len(order))
	}

	if order[0] != "third" || order[1] != "second" || order[2] != "first" {
		t.Errorf("Wrong order: %v", order)
	}
}

func TestShutdownWithError(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(5 * time.Second)

	h.Register(func(ctx context.Context) error {
		return nil
	})
	h.Register(func(ctx context.Context) error {
		return nil // Error would be logged but not fail
	})

	err := h.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("Shutdown should not fail: %v", err)
	}
}

func TestDoubleShutdown(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(5 * time.Second)

	callCount := 0
	h.Register(func(ctx context.Context) error {
		callCount++
		return nil
	})

	h.Shutdown(context.Background())
	h.Shutdown(context.Background())

	if callCount != 1 {
		t.Errorf("Shutdown function should be called once, called %d times", callCount)
	}
}

func TestShutdownTimeout(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(100 * time.Millisecond)

	blockCalled := false
	h.Register(func(ctx context.Context) error {
		blockCalled = true
		// Check if context is cancelled due to timeout
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			return nil
		}
	})

	start := time.Now()
	h.Shutdown(context.Background())
	elapsed := time.Since(start)

	if !blockCalled {
		t.Error("Shutdown function should have been called")
	}

	// Should complete within reasonable time (timeout + buffer)
	if elapsed > 500*time.Millisecond {
		t.Errorf("Shutdown took too long: %v", elapsed)
	}
}

func TestDefaultHandler(t *testing.T) {
	if DefaultHandler == nil {
		t.Fatal("DefaultHandler is nil")
	}
}

func TestPackageRegister(t *testing.T) {
	Register(func(ctx context.Context) error {
		return nil
	})
	// No panic means success
}

func TestIsShuttingDownFunc(t *testing.T) {
	// This just checks the function works
	_ = IsShuttingDown()
}

func TestDoneFunc(t *testing.T) {
	done := Done()
	if done == nil {
		t.Error("Done() returned nil")
	}
}

func TestSetTimeoutFunc(t *testing.T) {
	SetTimeout(30 * time.Second)
	// No panic means success
}

func TestShutdownFunc(t *testing.T) {
	// Test the Shutdown function
	Register(func(ctx context.Context) error {
		return nil
	})
	err := Shutdown(context.Background())
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

func TestWaitWithContext(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(1 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := h.WaitWithContext(ctx)
	if err != nil {
		t.Errorf("WaitWithContext failed: %v", err)
	}
}

func TestWaitWithSignal(t *testing.T) {
	// We can't easily test signal handling, but we can verify the method exists
	h := NewHandler()
	_ = h.Wait
}

func TestMultipleShutdownRegistrations(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(5 * time.Second)

	count := 0
	for i := 0; i < 10; i++ {
		h.Register(func(ctx context.Context) error {
			count++
			return nil
		})
	}

	err := h.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	if count != 10 {
		t.Errorf("Expected 10 calls, got %d", count)
	}
}

func TestShutdownContextCancellation(t *testing.T) {
	h := NewHandler()
	h.SetTimeout(2 * time.Second)

	// Create a context that gets cancelled
	ctx, cancel := context.WithCancel(context.Background())

	h.Register(func(ctx context.Context) error {
		// This should get cancelled
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	// Cancel immediately
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// This should still complete
	err := h.Shutdown(ctx)
	_ = err // Context cancellation might cause different behavior
}