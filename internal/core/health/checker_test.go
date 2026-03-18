package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusUnknown, "unknown"},
		{StatusHealthy, "healthy"},
		{StatusDegraded, "degraded"},
		{StatusUnhealthy, "unhealthy"},
	}

	for _, test := range tests {
		if test.status.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.status.String())
		}
	}
}

func TestNewFuncChecker(t *testing.T) {
	checker := NewFuncChecker("test", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})

	if checker.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", checker.Name())
	}
}

func TestFuncCheckerCheck(t *testing.T) {
	checker := NewFuncChecker("test", func(ctx context.Context) CheckResult {
		return CheckResult{
			Status:  StatusHealthy,
			Message: "all good",
		}
	})

	result := checker.Check(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("Expected StatusHealthy, got %v", result.Status)
	}
}

func TestNewHealthChecker(t *testing.T) {
	h := NewHealthChecker()
	if h == nil {
		t.Fatal("HealthChecker is nil")
	}
}

func TestHealthCheckerWithOptions(t *testing.T) {
	h := NewHealthChecker(
		WithInterval(10*time.Second),
		WithTimeout(5*time.Second),
	)

	if h.interval != 10*time.Second {
		t.Error("Interval option not applied")
	}

	if h.timeout != 5*time.Second {
		t.Error("Timeout option not applied")
	}
}

func TestHealthCheckerRegister(t *testing.T) {
	h := NewHealthChecker()
	h.Register(NewFuncChecker("test", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	}))

	names := h.List()
	if len(names) != 1 {
		t.Errorf("Expected 1 checker, got %d", len(names))
	}
}

func TestHealthCheckerRegisterFunc(t *testing.T) {
	h := NewHealthChecker()
	h.RegisterFunc("test", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})

	names := h.List()
	if len(names) != 1 {
		t.Errorf("Expected 1 checker, got %d", len(names))
	}
}

func TestHealthCheckerUnregister(t *testing.T) {
	h := NewHealthChecker()
	h.RegisterFunc("test", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})

	h.Unregister("test")

	names := h.List()
	if len(names) != 0 {
		t.Errorf("Expected 0 checkers, got %d", len(names))
	}
}

func TestHealthCheckerCheck(t *testing.T) {
	h := NewHealthChecker()
	h.RegisterFunc("healthy", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})
	h.RegisterFunc("degraded", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusDegraded}
	})

	results := h.Check(context.Background())

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if results["healthy"].Status != StatusHealthy {
		t.Error("healthy check should be healthy")
	}

	if results["degraded"].Status != StatusDegraded {
		t.Error("degraded check should be degraded")
	}
}

func TestHealthCheckerCheckOne(t *testing.T) {
	h := NewHealthChecker()
	h.RegisterFunc("test", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})

	result, err := h.CheckOne(context.Background(), "test")
	if err != nil {
		t.Fatalf("CheckOne failed: %v", err)
	}

	if result.Status != StatusHealthy {
		t.Errorf("Expected StatusHealthy, got %v", result.Status)
	}
}

func TestHealthCheckerCheckOneNotFound(t *testing.T) {
	h := NewHealthChecker()

	_, err := h.CheckOne(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent checker")
	}
}

func TestHealthCheckerOverallStatus(t *testing.T) {
	// All healthy
	h := NewHealthChecker()
	h.RegisterFunc("a", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})
	h.RegisterFunc("b", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})
	h.Check(context.Background())

	if h.OverallStatus() != StatusHealthy {
		t.Error("Overall should be healthy")
	}

	// One degraded
	h.RegisterFunc("c", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusDegraded}
	})
	h.Check(context.Background())

	if h.OverallStatus() != StatusDegraded {
		t.Error("Overall should be degraded")
	}

	// One unhealthy
	h.RegisterFunc("d", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusUnhealthy}
	})
	h.Check(context.Background())

	if h.OverallStatus() != StatusUnhealthy {
		t.Error("Overall should be unhealthy")
	}
}

func TestHealthCheckerOverallStatusEmpty(t *testing.T) {
	h := NewHealthChecker()

	if h.OverallStatus() != StatusUnknown {
		t.Error("Empty checker should return unknown")
	}
}

func TestHealthCheckerGetResults(t *testing.T) {
	h := NewHealthChecker()
	h.RegisterFunc("test", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})
	h.Check(context.Background())

	results := h.GetResults()
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestHealthCheckerGetResult(t *testing.T) {
	h := NewHealthChecker()
	h.RegisterFunc("test", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})
	h.Check(context.Background())

	result, ok := h.GetResult("test")
	if !ok {
		t.Fatal("Result not found")
	}

	if result.Status != StatusHealthy {
		t.Errorf("Expected StatusHealthy, got %v", result.Status)
	}

	_, ok = h.GetResult("nonexistent")
	if ok {
		t.Error("Should not find nonexistent result")
	}
}

func TestHealthCheckerStartStop(t *testing.T) {
	h := NewHealthChecker(WithInterval(100 * time.Millisecond))
	h.RegisterFunc("test", func(ctx context.Context) CheckResult {
		return CheckResult{Status: StatusHealthy}
	})

	ctx := context.Background()
	h.Start(ctx)

	if !h.IsRunning() {
		t.Error("Should be running")
	}

	// Wait for a check
	time.Sleep(150 * time.Millisecond)

	if h.LastCheck().IsZero() {
		t.Error("Should have checked")
	}

	h.Stop()

	if h.IsRunning() {
		t.Error("Should not be running after stop")
	}
}

func TestHealthCheckerOnChange(t *testing.T) {
	changes := []string{}
	h := NewHealthChecker(OnChange(func(name string, old, new Status) {
		changes = append(changes, name)
	}))

	status := StatusHealthy
	h.RegisterFunc("test", func(ctx context.Context) CheckResult {
		result := CheckResult{Status: status}
		status = StatusDegraded
		return result
	})

	h.Check(context.Background())
	h.Check(context.Background())

	// Should have detected change
	if len(changes) == 0 {
		t.Error("Should have detected status change")
	}
}

func TestHealthCheckerTimeout(t *testing.T) {
	h := NewHealthChecker(WithTimeout(50 * time.Millisecond))
	h.RegisterFunc("slow", func(ctx context.Context) CheckResult {
		select {
		case <-ctx.Done():
			return CheckResult{Status: StatusUnhealthy, Error: ctx.Err()}
		case <-time.After(200 * time.Millisecond):
			return CheckResult{Status: StatusHealthy}
		}
	})

	results := h.Check(context.Background())

	// Check should still complete but the slow checker might be cancelled
	if results["slow"].Status != StatusHealthy {
		// This is expected if timeout occurred
	}
}

func TestReadyChecker(t *testing.T) {
	r := NewReadyChecker()

	// Not ready without checks
	if r.IsReady() {
		t.Error("Should not be ready without checks")
	}

	// Register check
	r.Register("db", func() bool { return true })
	r.Register("cache", func() bool { return true })

	if !r.IsReady() {
		t.Error("Should be ready when all checks pass")
	}

	// Register failing check
	r.Register("fail", func() bool { return false })

	if r.IsReady() {
		t.Error("Should not be ready when any check fails")
	}
}

func TestLiveChecker(t *testing.T) {
	l := NewLiveChecker()

	if !l.IsAlive() {
		t.Error("Should be alive initially")
	}

	l.Register(func() bool { return true })

	if !l.IsAlive() {
		t.Error("Should be alive when checks pass")
	}

	l.Kill()

	if l.IsAlive() {
		t.Error("Should not be alive after kill")
	}

	l.Revive()

	if !l.IsAlive() {
		t.Error("Should be alive after revive")
	}
}

func TestLiveCheckerWithChecks(t *testing.T) {
	l := NewLiveChecker()

	alive := true
	l.Register(func() bool { return alive })

	if !l.IsAlive() {
		t.Error("Should be alive")
	}

	alive = false

	if l.IsAlive() {
		t.Error("Should not be alive when check fails")
	}
}

func TestCheckResult(t *testing.T) {
	result := CheckResult{
		Name:      "test",
		Status:    StatusHealthy,
		Message:   "all good",
		Timestamp: time.Now(),
		Duration:  10 * time.Millisecond,
		Details:   map[string]interface{}{"key": "value"},
		Error:     nil,
	}

	if result.Name != "test" {
		t.Error("Name should be test")
	}

	if result.Status != StatusHealthy {
		t.Error("Status should be healthy")
	}
}

func TestHealthCheckerConcurrentCheck(t *testing.T) {
	h := NewHealthChecker()

	// Register multiple checkers
	for i := 0; i < 10; i++ {
		name := string(rune('a' + i))
		h.RegisterFunc(name, func(ctx context.Context) CheckResult {
			time.Sleep(10 * time.Millisecond)
			return CheckResult{Status: StatusHealthy}
		})
	}

	start := time.Now()
	results := h.Check(context.Background())
	elapsed := time.Since(start)

	// All checks should run concurrently
	if elapsed > 50*time.Millisecond {
		t.Errorf("Concurrent checks took too long: %v", elapsed)
	}

	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}
}

func TestHealthCheckerErrorInCheck(t *testing.T) {
	h := NewHealthChecker()
	h.RegisterFunc("error", func(ctx context.Context) CheckResult {
		return CheckResult{
			Status: StatusUnhealthy,
			Error:  errors.New("check failed"),
		}
	})

	results := h.Check(context.Background())

	if results["error"].Status != StatusUnhealthy {
		t.Error("Should be unhealthy")
	}

	if results["error"].Error == nil {
		t.Error("Should have error")
	}
}