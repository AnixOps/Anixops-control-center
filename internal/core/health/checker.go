package health

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Status represents health status
type Status int

const (
	StatusUnknown Status = iota
	StatusHealthy
	StatusDegraded
	StatusUnhealthy
)

func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusDegraded:
		return "degraded"
	case StatusUnhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

// CheckResult represents a health check result
type CheckResult struct {
	Name      string                 `json:"name"`
	Status    Status                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
	Error     error                  `json:"error,omitempty"`
}

// Checker performs health checks
type Checker interface {
	Name() string
	Check(ctx context.Context) CheckResult
}

// CheckFunc is a function that performs a health check
type CheckFunc func(ctx context.Context) CheckResult

// FuncChecker wraps a function as a Checker
type FuncChecker struct {
	name string
	fn   CheckFunc
}

// NewFuncChecker creates a new function-based checker
func NewFuncChecker(name string, fn CheckFunc) *FuncChecker {
	return &FuncChecker{name: name, fn: fn}
}

func (c *FuncChecker) Name() string { return c.name }
func (c *FuncChecker) Check(ctx context.Context) CheckResult {
	return c.fn(ctx)
}

// HealthChecker manages health checks
type HealthChecker struct {
	mu         sync.RWMutex
	checkers   map[string]Checker
	results    map[string]CheckResult
	interval   time.Duration
	timeout    time.Duration
	lastCheck  time.Time
	onChange   func(name string, old, new Status)
	running    bool
	cancel     context.CancelFunc
}

// HealthCheckerOption configures a health checker
type HealthCheckerOption func(*HealthChecker)

// WithInterval sets the check interval
func WithInterval(d time.Duration) HealthCheckerOption {
	return func(h *HealthChecker) { h.interval = d }
}

// WithTimeout sets the check timeout
func WithTimeout(d time.Duration) HealthCheckerOption {
	return func(h *HealthChecker) { h.timeout = d }
}

// OnChange sets the status change callback
func OnChange(fn func(name string, old, new Status)) HealthCheckerOption {
	return func(h *HealthChecker) { h.onChange = fn }
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(opts ...HealthCheckerOption) *HealthChecker {
	h := &HealthChecker{
		checkers: make(map[string]Checker),
		results:  make(map[string]CheckResult),
		interval: 30 * time.Second,
		timeout:  10 * time.Second,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Register registers a health checker
func (h *HealthChecker) Register(checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers[checker.Name()] = checker
}

// RegisterFunc registers a function as a health checker
func (h *HealthChecker) RegisterFunc(name string, fn CheckFunc) {
	h.Register(NewFuncChecker(name, fn))
}

// Unregister removes a health checker
func (h *HealthChecker) Unregister(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.checkers, name)
	delete(h.results, name)
}

// Check performs all health checks
func (h *HealthChecker) Check(ctx context.Context) map[string]CheckResult {
	h.mu.RLock()
	checkers := make(map[string]Checker, len(h.checkers))
	for k, v := range h.checkers {
		checkers[k] = v
	}
	h.mu.RUnlock()

	results := make(map[string]CheckResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, checker := range checkers {
		wg.Add(1)
		go func(name string, checker Checker) {
			defer wg.Done()

			// Create timeout context
			checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
			defer cancel()

			start := time.Now()
			result := checker.Check(checkCtx)
			result.Duration = time.Since(start)
			result.Timestamp = time.Now()

			mu.Lock()
			// Check for status change
			if old, exists := h.results[name]; exists && h.onChange != nil {
				if old.Status != result.Status {
					h.onChange(name, old.Status, result.Status)
				}
			}
			results[name] = result
			mu.Unlock()
		}(name, checker)
	}

	wg.Wait()

	h.mu.Lock()
	h.results = results
	h.lastCheck = time.Now()
	h.mu.Unlock()

	return results
}

// CheckOne performs a single health check
func (h *HealthChecker) CheckOne(ctx context.Context, name string) (CheckResult, error) {
	h.mu.RLock()
	checker, exists := h.checkers[name]
	h.mu.RUnlock()

	if !exists {
		return CheckResult{}, fmt.Errorf("checker %s not found", name)
	}

	checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	start := time.Now()
	result := checker.Check(checkCtx)
	result.Duration = time.Since(start)
	result.Timestamp = time.Now()

	h.mu.Lock()
	if old, exists := h.results[name]; exists && h.onChange != nil {
		if old.Status != result.Status {
			h.onChange(name, old.Status, result.Status)
		}
	}
	h.results[name] = result
	h.mu.Unlock()

	return result, nil
}

// GetResults returns the last check results
func (h *HealthChecker) GetResults() map[string]CheckResult {
	h.mu.RLock()
	defer h.mu.RUnlock()
	results := make(map[string]CheckResult, len(h.results))
	for k, v := range h.results {
		results[k] = v
	}
	return results
}

// GetResult returns the last check result for a specific checker
func (h *HealthChecker) GetResult(name string) (CheckResult, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result, ok := h.results[name]
	return result, ok
}

// OverallStatus returns the overall health status
func (h *HealthChecker) OverallStatus() Status {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.results) == 0 {
		return StatusUnknown
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, result := range h.results {
		switch result.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return StatusUnhealthy
	}
	if hasDegraded {
		return StatusDegraded
	}
	return StatusHealthy
}

// Start starts periodic health checks
func (h *HealthChecker) Start(ctx context.Context) {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return
	}
	h.running = true
	ctx, h.cancel = context.WithCancel(ctx)
	h.mu.Unlock()

	go func() {
		ticker := time.NewTicker(h.interval)
		defer ticker.Stop()

		// Initial check
		h.Check(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.Check(ctx)
			}
		}
	}()
}

// Stop stops periodic health checks
func (h *HealthChecker) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.running && h.cancel != nil {
		h.cancel()
	}
	h.running = false
}

// LastCheck returns the time of the last check
func (h *HealthChecker) LastCheck() time.Time {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.lastCheck
}

// List returns all registered checker names
func (h *HealthChecker) List() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	names := make([]string, 0, len(h.checkers))
	for name := range h.checkers {
		names = append(names, name)
	}
	return names
}

// IsRunning returns true if periodic checks are running
func (h *HealthChecker) IsRunning() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.running
}

// ReadyChecker checks if the component is ready
type ReadyChecker struct {
	checkers map[string]func() bool
	mu       sync.RWMutex
}

// NewReadyChecker creates a new readiness checker
func NewReadyChecker() *ReadyChecker {
	return &ReadyChecker{
		checkers: make(map[string]func() bool),
	}
}

// Register registers a readiness check
func (r *ReadyChecker) Register(name string, check func() bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checkers[name] = check
}

// IsReady returns true if all checks pass
func (r *ReadyChecker) IsReady() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, check := range r.checkers {
		if !check() {
			return false
		}
	}
	return len(r.checkers) > 0
}

// LiveChecker checks if the component is alive
type LiveChecker struct {
	alive    bool
	mu       sync.RWMutex
	checkers []func() bool
}

// NewLiveChecker creates a new liveness checker
func NewLiveChecker() *LiveChecker {
	return &LiveChecker{
		alive:    true,
		checkers: make([]func() bool, 0),
	}
}

// Register registers a liveness check
func (l *LiveChecker) Register(check func() bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.checkers = append(l.checkers, check)
}

// IsAlive returns true if alive
func (l *LiveChecker) IsAlive() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if !l.alive {
		return false
	}
	for _, check := range l.checkers {
		if !check() {
			return false
		}
	}
	return true
}

// Kill marks the checker as not alive
func (l *LiveChecker) Kill() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.alive = false
}

// Revive marks the checker as alive
func (l *LiveChecker) Revive() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.alive = true
}