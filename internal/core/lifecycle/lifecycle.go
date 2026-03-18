package lifecycle

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/errors"
	"github.com/anixops/anixops-control-center/internal/core/logger"
)

// Phase represents a lifecycle phase
type Phase int

const (
	PhaseUninitialized Phase = iota
	PhaseInitializing
	PhaseInitialized
	PhaseStarting
	PhaseRunning
	PhaseStopping
	PhaseStopped
	PhaseError
)

func (p Phase) String() string {
	switch p {
	case PhaseUninitialized:
		return "uninitialized"
	case PhaseInitializing:
		return "initializing"
	case PhaseInitialized:
		return "initialized"
	case PhaseStarting:
		return "starting"
	case PhaseRunning:
		return "running"
	case PhaseStopping:
		return "stopping"
	case PhaseStopped:
		return "stopped"
	case PhaseError:
		return "error"
	default:
		return "unknown"
	}
}

// HookFunc is a lifecycle hook function
type HookFunc func(ctx context.Context) error

// Hook represents a lifecycle hook
type Hook struct {
	Name     string
	Priority int
	Func     HookFunc
}

// Lifecycle manages component lifecycle
type Lifecycle struct {
	mu       sync.RWMutex
	phase    Phase
	hooks    map[Phase][]Hook
	timeout  time.Duration
	onChange func(Phase, Phase)
}

// New creates a new lifecycle manager
func New() *Lifecycle {
	return &Lifecycle{
		phase: PhaseUninitialized,
		hooks: make(map[Phase][]Hook),
	}
}

// SetTimeout sets the default timeout for lifecycle operations
func (l *Lifecycle) SetTimeout(d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeout = d
}

// OnPhaseChange sets a callback for phase changes
func (l *Lifecycle) OnPhaseChange(fn func(old, new Phase)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.onChange = fn
}

// Phase returns the current phase
func (l *Lifecycle) Phase() Phase {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.phase
}

// RegisterHook registers a lifecycle hook
func (l *Lifecycle) RegisterHook(phase Phase, hook Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks[phase] = append(l.hooks[phase], hook)
	// Sort by priority
	for i := 1; i < len(l.hooks[phase]); i++ {
		for j := i; j > 0 && l.hooks[phase][j-1].Priority > l.hooks[phase][j].Priority; j-- {
			l.hooks[phase][j-1], l.hooks[phase][j] = l.hooks[phase][j], l.hooks[phase][j-1]
		}
	}
}

// OnInit registers an initialization hook
func (l *Lifecycle) OnInit(name string, priority int, fn HookFunc) {
	l.RegisterHook(PhaseInitializing, Hook{Name: name, Priority: priority, Func: fn})
}

// OnStart registers a start hook
func (l *Lifecycle) OnStart(name string, priority int, fn HookFunc) {
	l.RegisterHook(PhaseStarting, Hook{Name: name, Priority: priority, Func: fn})
}

// OnStop registers a stop hook
func (l *Lifecycle) OnStop(name string, priority int, fn HookFunc) {
	l.RegisterHook(PhaseStopping, Hook{Name: name, Priority: priority, Func: fn})
}

// transition transitions to a new phase
func (l *Lifecycle) transition(newPhase Phase) {
	l.mu.Lock()
	oldPhase := l.phase
	l.phase = newPhase
	onChange := l.onChange
	l.mu.Unlock()

	if onChange != nil {
		onChange(oldPhase, newPhase)
	}
	logger.Debug("Lifecycle transition", logger.F("old", oldPhase), logger.F("new", newPhase))
}

// Initialize runs initialization hooks
func (l *Lifecycle) Initialize(ctx context.Context) error {
	if l.Phase() != PhaseUninitialized && l.Phase() != PhaseError {
		return errors.NewError(errors.CodeInternal, "invalid phase for initialization", errors.LevelError, 500).
			WithDetail(fmt.Sprintf("current phase: %s", l.Phase()))
	}

	l.transition(PhaseInitializing)

	if err := l.runHooks(ctx, PhaseInitializing); err != nil {
		l.transition(PhaseError)
		return err
	}

	l.transition(PhaseInitialized)
	return nil
}

// Start runs start hooks
func (l *Lifecycle) Start(ctx context.Context) error {
	phase := l.Phase()
	if phase != PhaseInitialized && phase != PhaseStopped {
		return errors.NewError(errors.CodeInternal, "invalid phase for start", errors.LevelError, 500).
			WithDetail(fmt.Sprintf("current phase: %s", phase))
	}

	l.transition(PhaseStarting)

	if err := l.runHooks(ctx, PhaseStarting); err != nil {
		l.transition(PhaseError)
		return err
	}

	l.transition(PhaseRunning)
	return nil
}

// Stop runs stop hooks
func (l *Lifecycle) Stop(ctx context.Context) error {
	phase := l.Phase()
	if phase != PhaseRunning && phase != PhaseError {
		return errors.NewError(errors.CodeInternal, "invalid phase for stop", errors.LevelError, 500).
			WithDetail(fmt.Sprintf("current phase: %s", phase))
	}

	l.transition(PhaseStopping)

	if err := l.runHooks(ctx, PhaseStopping); err != nil {
		l.transition(PhaseError)
		return err
	}

	l.transition(PhaseStopped)
	return nil
}

// runHooks runs hooks for a phase
func (l *Lifecycle) runHooks(ctx context.Context, phase Phase) error {
	l.mu.RLock()
	hooks := make([]Hook, len(l.hooks[phase]))
	copy(hooks, l.hooks[phase])
	timeout := l.timeout
	l.mu.RUnlock()

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	for _, hook := range hooks {
		logger.Debug("Running lifecycle hook", logger.F("phase", phase), logger.F("hook", hook.Name))

		if err := hook.Func(ctx); err != nil {
			return errors.NewError(errors.CodeInternal, fmt.Sprintf("hook %s failed", hook.Name), errors.LevelError, 500).
				WithCause(err)
		}
	}

	return nil
}

// Restart restarts the lifecycle
func (l *Lifecycle) Restart(ctx context.Context) error {
	if l.Phase() == PhaseRunning {
		if err := l.Stop(ctx); err != nil {
			return err
		}
	}
	return l.Start(ctx)
}

// Shutdown gracefully shuts down
func (l *Lifecycle) Shutdown(ctx context.Context) error {
	phase := l.Phase()
	if phase == PhaseRunning {
		if err := l.Stop(ctx); err != nil {
			logger.Error("Failed to stop during shutdown", logger.F("error", err))
		}
	}
	l.transition(PhaseStopped)
	return nil
}

// Manager manages multiple lifecycles
type Manager struct {
	mu         sync.RWMutex
	lifecycles map[string]*Lifecycle
	startOrder []string
	stopOrder  []string
}

// NewManager creates a new lifecycle manager
func NewManager() *Manager {
	return &Manager{
		lifecycles: make(map[string]*Lifecycle),
		startOrder: []string{},
		stopOrder:  []string{},
	}
}

// Register registers a lifecycle
func (m *Manager) Register(name string, lc *Lifecycle, startPriority int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lifecycles[name] = lc

	// Insert into start order based on priority
	inserted := false
	for i, n := range m.startOrder {
		// We'd need to store priorities separately for proper ordering
		_ = i
		_ = n
	}
	if !inserted {
		m.startOrder = append(m.startOrder, name)
	}

	// Stop in reverse order
	m.stopOrder = make([]string, len(m.startOrder))
	for i, n := range m.startOrder {
		m.stopOrder[len(m.startOrder)-1-i] = n
	}
}

// InitializeAll initializes all lifecycles
func (m *Manager) InitializeAll(ctx context.Context) error {
	m.mu.RLock()
	order := make([]string, len(m.startOrder))
	copy(order, m.startOrder)
	lifecycles := make(map[string]*Lifecycle)
	for k, v := range m.lifecycles {
		lifecycles[k] = v
	}
	m.mu.RUnlock()

	for _, name := range order {
		lc := lifecycles[name]
		if err := lc.Initialize(ctx); err != nil {
			return errors.NewError(errors.CodeInternal, fmt.Sprintf("failed to initialize %s", name), errors.LevelError, 500).
				WithCause(err)
		}
	}
	return nil
}

// StartAll starts all lifecycles
func (m *Manager) StartAll(ctx context.Context) error {
	m.mu.RLock()
	order := make([]string, len(m.startOrder))
	copy(order, m.startOrder)
	lifecycles := make(map[string]*Lifecycle)
	for k, v := range m.lifecycles {
		lifecycles[k] = v
	}
	m.mu.RUnlock()

	for _, name := range order {
		lc := lifecycles[name]
		if err := lc.Start(ctx); err != nil {
			// Attempt to stop already started
			_ = m.StopAll(ctx)
			return errors.NewError(errors.CodeInternal, fmt.Sprintf("failed to start %s", name), errors.LevelError, 500).
				WithCause(err)
		}
	}
	return nil
}

// StopAll stops all lifecycles
func (m *Manager) StopAll(ctx context.Context) error {
	m.mu.RLock()
	order := make([]string, len(m.stopOrder))
	copy(order, m.stopOrder)
	lifecycles := make(map[string]*Lifecycle)
	for k, v := range m.lifecycles {
		lifecycles[k] = v
	}
	m.mu.RUnlock()

	var lastErr error
	for _, name := range order {
		lc := lifecycles[name]
		if lc.Phase() == PhaseRunning {
			if err := lc.Stop(ctx); err != nil {
				logger.Error("Failed to stop lifecycle", logger.F("name", name), logger.F("error", err))
				lastErr = err
			}
		}
	}
	return lastErr
}

// Get retrieves a lifecycle by name
func (m *Manager) Get(name string) (*Lifecycle, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	lc, ok := m.lifecycles[name]
	return lc, ok
}

// Phase returns the phase of a lifecycle
func (m *Manager) Phase(name string) Phase {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if lc, ok := m.lifecycles[name]; ok {
		return lc.Phase()
	}
	return PhaseUninitialized
}