package hotreload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/config"
)

// ConfigType represents the type of configuration
type ConfigType string

const (
	ConfigTypePlugin  ConfigType = "plugin"
	ConfigTypeSystem  ConfigType = "system"
	ConfigTypeService ConfigType = "service"
)

// ConfigChange represents a configuration change
type ConfigChange struct {
	Type      ConfigType         `json:"type"`
	Name      string             `json:"name"`
	OldConfig map[string]interface{} `json:"old_config,omitempty"`
	NewConfig map[string]interface{} `json:"new_config"`
	Timestamp time.Time          `json:"timestamp"`
	Source    string             `json:"source"` // file, api, etc.
}

// ConfigWatcher watches for configuration changes
type ConfigWatcher interface {
	Watch(ctx context.Context) (<-chan ConfigChange, error)
	Close() error
}

// ReloadHandler handles configuration reload
type ReloadHandler interface {
	Name() string
	CanReload(change ConfigChange) bool
	Reload(ctx context.Context, change ConfigChange) error
	Rollback(ctx context.Context, change ConfigChange) error
}

// Manager manages hot reload of configurations
type Manager struct {
	mu          sync.RWMutex
	watchers    []ConfigWatcher
	handlers    map[string]ReloadHandler
	history     []ConfigChange
	maxHistory  int
	pending     map[string]ConfigChange
	onChange    func(change ConfigChange)
	onError     func(change ConfigChange, err error)
	running     bool
	cancel      context.CancelFunc
	config      *config.Config
	configPath  string
}

// ManagerOption configures the manager
type ManagerOption func(*Manager)

// WithMaxHistory sets the maximum history size
func WithMaxHistory(max int) ManagerOption {
	return func(m *Manager) { m.maxHistory = max }
}

// OnChange sets the change callback
func OnChange(fn func(change ConfigChange)) ManagerOption {
	return func(m *Manager) { m.onChange = fn }
}

// OnError sets the error callback
func OnError(fn func(change ConfigChange, err error)) ManagerOption {
	return func(m *Manager) { m.onError = fn }
}

// WithConfigPath sets the configuration file path
func WithConfigPath(path string) ManagerOption {
	return func(m *Manager) { m.configPath = path }
}

// NewManager creates a new hot reload manager
func NewManager(cfg *config.Config, opts ...ManagerOption) *Manager {
	m := &Manager{
		handlers:   make(map[string]ReloadHandler),
		history:    make([]ConfigChange, 0),
		maxHistory: 100,
		pending:    make(map[string]ConfigChange),
		config:     cfg,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// RegisterHandler registers a reload handler
func (m *Manager) RegisterHandler(handler ReloadHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[handler.Name()] = handler
}

// UnregisterHandler removes a reload handler
func (m *Manager) UnregisterHandler(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, name)
}

// AddWatcher adds a configuration watcher
func (m *Manager) AddWatcher(watcher ConfigWatcher) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.watchers = append(m.watchers, watcher)
}

// Reload reloads configuration for a specific component
func (m *Manager) Reload(ctx context.Context, change ConfigChange) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find handler
	handler, exists := m.handlers[change.Name]
	if !exists {
		return fmt.Errorf("no handler for %s", change.Name)
	}

	// Check if reload is possible
	if !handler.CanReload(change) {
		return fmt.Errorf("reload not allowed for %s", change.Name)
	}

	// Store in pending
	m.pending[change.Name] = change

	// Execute reload
	if err := handler.Reload(ctx, change); err != nil {
		// Attempt rollback
		rollbackErr := handler.Rollback(ctx, change)
		if m.onError != nil {
			m.onError(change, err)
		}
		if rollbackErr != nil {
			return fmt.Errorf("reload failed: %w, rollback also failed: %v", err, rollbackErr)
		}
		return fmt.Errorf("reload failed, rolled back: %w", err)
	}

	// Record in history
	m.history = append(m.history, change)
	if len(m.history) > m.maxHistory {
		m.history = m.history[len(m.history)-m.maxHistory:]
	}

	// Remove from pending
	delete(m.pending, change.Name)

	// Notify callback
	if m.onChange != nil {
		m.onChange(change)
	}

	return nil
}

// Rollback rolls back a configuration change
func (m *Manager) Rollback(ctx context.Context, change ConfigChange) error {
	m.mu.RLock()
	handler, exists := m.handlers[change.Name]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handler for %s", change.Name)
	}

	return handler.Rollback(ctx, change)
}

// ApplyConfig applies new configuration
func (m *Manager) ApplyConfig(ctx context.Context, configType ConfigType, name string, newConfig map[string]interface{}) error {
	change := ConfigChange{
		Type:      configType,
		Name:      name,
		NewConfig: newConfig,
		Timestamp: time.Now(),
		Source:    "api",
	}

	// Get old config if available
	m.mu.RLock()
	for _, c := range m.history {
		if c.Type == configType && c.Name == name {
			change.OldConfig = c.NewConfig
			break
		}
	}
	m.mu.RUnlock()

	return m.Reload(ctx, change)
}

// GetHistory returns configuration change history
func (m *Manager) GetHistory() []ConfigChange {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]ConfigChange, len(m.history))
	copy(result, m.history)
	return result
}

// GetPending returns pending configuration changes
func (m *Manager) GetPending() map[string]ConfigChange {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]ConfigChange)
	for k, v := range m.pending {
		result[k] = v
	}
	return result
}

// Start starts watching for configuration changes
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return nil
	}
	m.running = true
	ctx, m.cancel = context.WithCancel(ctx)
	m.mu.Unlock()

	// Start watchers
	for _, watcher := range m.watchers {
		go func(w ConfigWatcher) {
			ch, err := w.Watch(ctx)
			if err != nil {
				return
			}
			for {
				select {
				case <-ctx.Done():
					return
				case change, ok := <-ch:
					if !ok {
						return
					}
					m.Reload(ctx, change)
				}
			}
		}(watcher)
	}

	return nil
}

// Stop stops watching for configuration changes
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running && m.cancel != nil {
		m.cancel()
	}
	m.running = false

	// Close watchers
	var lastErr error
	for _, watcher := range m.watchers {
		if err := watcher.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// IsRunning returns true if the manager is running
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// Validate validates configuration before applying
func (m *Manager) Validate(configType ConfigType, name string, newConfig map[string]interface{}) error {
	m.mu.RLock()
	handler, exists := m.handlers[name]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handler for %s", name)
	}

	change := ConfigChange{
		Type:      configType,
		Name:      name,
		NewConfig: newConfig,
	}

	if !handler.CanReload(change) {
		return fmt.Errorf("configuration not valid for %s", name)
	}

	return nil
}

// Export exports current configuration
func (m *Manager) Export() (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{})

	// Export from handlers
	for name := range m.handlers {
		// Find last config for this handler
		for i := len(m.history) - 1; i >= 0; i-- {
			if m.history[i].Name == name {
				result[name] = m.history[i].NewConfig
				break
			}
		}
	}

	return result, nil
}

// Import imports configuration
func (m *Manager) Import(ctx context.Context, configs map[string]interface{}) error {
	for name, cfg := range configs {
		configMap, ok := cfg.(map[string]interface{})
		if !ok {
			continue
		}
		if err := m.ApplyConfig(ctx, ConfigTypePlugin, name, configMap); err != nil {
			return fmt.Errorf("failed to apply config for %s: %w", name, err)
		}
	}
	return nil
}

// ToJSON converts config change to JSON
func (c *ConfigChange) ToJSON() string {
	b, _ := json.Marshal(c)
	return string(b)
}

// FromJSON parses config change from JSON
func FromJSON(data string) (ConfigChange, error) {
	var change ConfigChange
	err := json.Unmarshal([]byte(data), &change)
	return change, err
}

// Common errors
var (
	ErrHandlerNotFound = errors.New("handler not found")
	ErrReloadNotAllowed = errors.New("reload not allowed")
	ErrRollbackFailed = errors.New("rollback failed")
)