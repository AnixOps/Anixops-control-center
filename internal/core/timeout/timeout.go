package timeout

import (
	"context"
	"sync"
	"time"
)

// Manager manages timeouts for operations
type Manager struct {
	mu       sync.RWMutex
	defaults map[string]time.Duration
}

// NewManager creates a new timeout manager
func NewManager() *Manager {
	return &Manager{
		defaults: make(map[string]time.Duration),
	}
}

// SetDefault sets the default timeout for an operation
func (m *Manager) SetDefault(operation string, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaults[operation] = d
}

// GetDefault gets the default timeout for an operation
func (m *Manager) GetDefault(operation string) (time.Duration, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.defaults[operation]
	return d, ok
}

// Context creates a context with timeout for an operation
func (m *Manager) Context(ctx context.Context, operation string) (context.Context, context.CancelFunc) {
	d, ok := m.GetDefault(operation)
	if !ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, d)
}

// ContextWithFallback creates a context with timeout, falling back to a default
func (m *Manager) ContextWithFallback(ctx context.Context, operation string, fallback time.Duration) (context.Context, context.CancelFunc) {
	d, ok := m.GetDefault(operation)
	if !ok {
		d = fallback
	}
	return context.WithTimeout(ctx, d)
}

// Execute runs a function with a timeout
func (m *Manager) Execute(ctx context.Context, operation string, fn func(context.Context) error) error {
	ctx, cancel := m.Context(ctx, operation)
	defer cancel()
	return fn(ctx)
}

// ExecuteWithFallback runs a function with a timeout and fallback
func (m *Manager) ExecuteWithFallback(ctx context.Context, operation string, fallback time.Duration, fn func(context.Context) error) error {
	ctx, cancel := m.ContextWithFallback(ctx, operation, fallback)
	defer cancel()
	return fn(ctx)
}

// Result represents the result of a timed operation
type Result struct {
	Value      interface{}
	Error      error
	TimedOut   bool
	Duration   time.Duration
	FinishedAt time.Time
}

// Run executes a function with timeout and returns detailed result
func Run(ctx context.Context, timeout time.Duration, fn func(context.Context) (interface{}, error)) Result {
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultCh := make(chan Result, 1)
	go func() {
		val, err := fn(ctx)
		resultCh <- Result{
			Value:      val,
			Error:      err,
			Duration:   time.Since(start),
			FinishedAt: time.Now(),
		}
	}()

	select {
	case result := <-resultCh:
		return result
	case <-ctx.Done():
		return Result{
			Error:      ctx.Err(),
			TimedOut:   true,
			Duration:   time.Since(start),
			FinishedAt: time.Now(),
		}
	}
}

// Pool manages a pool of timeouts
type Pool struct {
	mu       sync.RWMutex
	timeouts map[string]time.Duration
}

// NewPool creates a new timeout pool
func NewPool() *Pool {
	return &Pool{
		timeouts: make(map[string]time.Duration),
	}
}

// Register registers a timeout for a key
func (p *Pool) Register(key string, timeout time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.timeouts[key] = timeout
}

// Unregister removes a timeout for a key
func (p *Pool) Unregister(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.timeouts, key)
}

// Get gets a timeout for a key
func (p *Pool) Get(key string) (time.Duration, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	d, ok := p.timeouts[key]
	return d, ok
}

// Context creates a context with the registered timeout
func (p *Pool) Context(ctx context.Context, key string) (context.Context, context.CancelFunc, bool) {
	d, ok := p.Get(key)
	if !ok {
		return ctx, func() {}, false
	}
	ctx, cancel := context.WithTimeout(ctx, d)
	return ctx, cancel, true
}

// Keys returns all registered keys
func (p *Pool) Keys() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	keys := make([]string, 0, len(p.timeouts))
	for k := range p.timeouts {
		keys = append(keys, k)
	}
	return keys
}

// Config represents timeout configuration
type Config struct {
	Default     time.Duration            `json:"default" yaml:"default"`
	Operations  map[string]time.Duration `json:"operations" yaml:"operations"`
	GracePeriod time.Duration            `json:"grace_period" yaml:"grace_period"`
}

// FromConfig creates a Manager from configuration
func FromConfig(cfg Config) *Manager {
	m := NewManager()
	for op, d := range cfg.Operations {
		m.SetDefault(op, d)
	}
	return m
}

// Common timeout presets
const (
	DefaultTimeout     = 30 * time.Second
	ShortTimeout       = 5 * time.Second
	LongTimeout        = 5 * time.Minute
	ConnectionTimeout  = 10 * time.Second
	ReadTimeout        = 15 * time.Second
	WriteTimeout       = 15 * time.Second
	ShutdownTimeout    = 30 * time.Second
	HealthCheckTimeout = 5 * time.Second
)

// TimeoutError represents a timeout error
type TimeoutError struct {
	Operation string
	Duration  time.Duration
}

func (e *TimeoutError) Error() string {
	return "operation " + e.Operation + " timed out after " + e.Duration.String()
}

// IsTimeout checks if an error is a timeout error
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(*TimeoutError); ok {
		return true
	}
	if err == context.DeadlineExceeded {
		return true
	}
	return false
}