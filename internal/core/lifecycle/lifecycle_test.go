package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewLifecycle(t *testing.T) {
	lc := New()
	if lc == nil {
		t.Fatal("Lifecycle is nil")
	}

	if lc.Phase() != PhaseUninitialized {
		t.Errorf("Expected PhaseUninitialized, got %v", lc.Phase())
	}
}

func TestPhaseString(t *testing.T) {
	tests := []struct {
		phase    Phase
		expected string
	}{
		{PhaseUninitialized, "uninitialized"},
		{PhaseInitializing, "initializing"},
		{PhaseInitialized, "initialized"},
		{PhaseStarting, "starting"},
		{PhaseRunning, "running"},
		{PhaseStopping, "stopping"},
		{PhaseStopped, "stopped"},
		{PhaseError, "error"},
	}

	for _, test := range tests {
		if test.phase.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.phase.String())
		}
	}
}

func TestSetTimeout(t *testing.T) {
	lc := New()
	lc.SetTimeout(5 * time.Second)
	// No error means success
}

func TestOnPhaseChange(t *testing.T) {
	lc := New()

	changes := []Phase{}
	lc.OnPhaseChange(func(old, new Phase) {
		changes = append(changes, new)
	})

	ctx := context.Background()
	lc.Initialize(ctx)

	if len(changes) == 0 {
		t.Error("Expected phase change callback")
	}
}

func TestInitialize(t *testing.T) {
	lc := New()
	ctx := context.Background()

	err := lc.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if lc.Phase() != PhaseInitialized {
		t.Errorf("Expected PhaseInitialized, got %v", lc.Phase())
	}
}

func TestInitializeWithHooks(t *testing.T) {
	lc := New()
	hookCalled := false

	lc.OnInit("test-hook", 0, func(ctx context.Context) error {
		hookCalled = true
		return nil
	})

	ctx := context.Background()
	lc.Initialize(ctx)

	if !hookCalled {
		t.Error("Init hook was not called")
	}
}

func TestInitializeHookError(t *testing.T) {
	lc := New()
	hookErr := errors.New("hook error")

	lc.OnInit("failing-hook", 0, func(ctx context.Context) error {
		return hookErr
	})

	ctx := context.Background()
	err := lc.Initialize(ctx)

	if err == nil {
		t.Error("Expected error from failing hook")
	}

	if lc.Phase() != PhaseError {
		t.Errorf("Expected PhaseError, got %v", lc.Phase())
	}
}

func TestStart(t *testing.T) {
	lc := New()
	ctx := context.Background()

	lc.Initialize(ctx)
	err := lc.Start(ctx)

	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if lc.Phase() != PhaseRunning {
		t.Errorf("Expected PhaseRunning, got %v", lc.Phase())
	}
}

func TestStartWithHooks(t *testing.T) {
	lc := New()
	hookCalled := false

	lc.OnStart("test-hook", 0, func(ctx context.Context) error {
		hookCalled = true
		return nil
	})

	ctx := context.Background()
	lc.Initialize(ctx)
	lc.Start(ctx)

	if !hookCalled {
		t.Error("Start hook was not called")
	}
}

func TestStop(t *testing.T) {
	lc := New()
	ctx := context.Background()

	lc.Initialize(ctx)
	lc.Start(ctx)
	err := lc.Stop(ctx)

	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if lc.Phase() != PhaseStopped {
		t.Errorf("Expected PhaseStopped, got %v", lc.Phase())
	}
}

func TestStopWithHooks(t *testing.T) {
	lc := New()
	hookCalled := false

	lc.OnStop("test-hook", 0, func(ctx context.Context) error {
		hookCalled = true
		return nil
	})

	ctx := context.Background()
	lc.Initialize(ctx)
	lc.Start(ctx)
	lc.Stop(ctx)

	if !hookCalled {
		t.Error("Stop hook was not called")
	}
}

func TestRestart(t *testing.T) {
	lc := New()
	ctx := context.Background()

	lc.Initialize(ctx)
	lc.Start(ctx)
	err := lc.Restart(ctx)

	if err != nil {
		t.Fatalf("Restart failed: %v", err)
	}

	if lc.Phase() != PhaseRunning {
		t.Errorf("Expected PhaseRunning after restart, got %v", lc.Phase())
	}
}

func TestShutdown(t *testing.T) {
	lc := New()
	ctx := context.Background()

	lc.Initialize(ctx)
	lc.Start(ctx)
	err := lc.Shutdown(ctx)

	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
}

func TestHookPriority(t *testing.T) {
	lc := New()
	order := []string{}

	lc.OnInit("low", 100, func(ctx context.Context) error {
		order = append(order, "low")
		return nil
	})
	lc.OnInit("high", 1, func(ctx context.Context) error {
		order = append(order, "high")
		return nil
	})
	lc.OnInit("medium", 50, func(ctx context.Context) error {
		order = append(order, "medium")
		return nil
	})

	ctx := context.Background()
	lc.Initialize(ctx)

	if len(order) != 3 {
		t.Fatalf("Expected 3 hooks, got %d", len(order))
	}

	// High priority (low number) should be first
	if order[0] != "high" {
		t.Errorf("Expected 'high' first, got '%s'", order[0])
	}
}

func TestInvalidPhaseTransitions(t *testing.T) {
	lc := New()
	ctx := context.Background()

	// Can't start from uninitialized
	err := lc.Start(ctx)
	if err == nil {
		t.Error("Expected error when starting from uninitialized")
	}

	// Can't stop when not running
	err = lc.Stop(ctx)
	if err == nil {
		t.Error("Expected error when stopping from uninitialized")
	}
}

// Manager tests
func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("Manager is nil")
	}
}

func TestManagerRegister(t *testing.T) {
	m := NewManager()
	lc := New()
	m.Register("test", lc, 0)

	_, ok := m.Get("test")
	if !ok {
		t.Error("Lifecycle should be registered")
	}
}

func TestManagerGet(t *testing.T) {
	m := NewManager()
	lc := New()
	m.Register("test", lc, 0)

	retrieved, ok := m.Get("test")
	if !ok {
		t.Fatal("Lifecycle not found")
	}

	if retrieved != lc {
		t.Error("Retrieved lifecycle should be same")
	}
}

func TestManagerGetNotFound(t *testing.T) {
	m := NewManager()

	_, ok := m.Get("nonexistent")
	if ok {
		t.Error("Should not find nonexistent lifecycle")
	}
}

func TestManagerPhase(t *testing.T) {
	m := NewManager()
	lc := New()
	m.Register("test", lc, 0)

	phase := m.Phase("test")
	if phase != PhaseUninitialized {
		t.Errorf("Expected PhaseUninitialized, got %v", phase)
	}
}

func TestManagerInitializeAll(t *testing.T) {
	m := NewManager()
	lc1 := New()
	lc2 := New()
	m.Register("test1", lc1, 0)
	m.Register("test2", lc2, 0)

	ctx := context.Background()
	err := m.InitializeAll(ctx)

	if err != nil {
		t.Fatalf("InitializeAll failed: %v", err)
	}

	if lc1.Phase() != PhaseInitialized || lc2.Phase() != PhaseInitialized {
		t.Error("All lifecycles should be initialized")
	}
}

func TestManagerStartAll(t *testing.T) {
	m := NewManager()
	lc1 := New()
	lc2 := New()
	m.Register("test1", lc1, 0)
	m.Register("test2", lc2, 0)

	ctx := context.Background()
	m.InitializeAll(ctx)
	err := m.StartAll(ctx)

	if err != nil {
		t.Fatalf("StartAll failed: %v", err)
	}

	if lc1.Phase() != PhaseRunning || lc2.Phase() != PhaseRunning {
		t.Error("All lifecycles should be running")
	}
}

func TestManagerStopAll(t *testing.T) {
	m := NewManager()
	lc1 := New()
	lc2 := New()
	m.Register("test1", lc1, 0)
	m.Register("test2", lc2, 0)

	ctx := context.Background()
	m.InitializeAll(ctx)
	m.StartAll(ctx)
	err := m.StopAll(ctx)

	if err != nil {
		t.Fatalf("StopAll failed: %v", err)
	}

	if lc1.Phase() != PhaseStopped || lc2.Phase() != PhaseStopped {
		t.Error("All lifecycles should be stopped")
	}
}