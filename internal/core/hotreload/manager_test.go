package hotreload

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestConfigType(t *testing.T) {
	types := []ConfigType{ConfigTypePlugin, ConfigTypeSystem, ConfigTypeService}
	for _, ct := range types {
		if string(ct) == "" {
			t.Error("ConfigType should not be empty")
		}
	}
}

func TestConfigChange(t *testing.T) {
	change := ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "test-plugin",
		NewConfig: map[string]interface{}{"key": "value"},
		Timestamp: time.Now(),
		Source:    "api",
	}

	if change.Type != ConfigTypePlugin {
		t.Error("Type should be plugin")
	}

	if change.Name != "test-plugin" {
		t.Error("Name should be test-plugin")
	}
}

func TestConfigChangeJSON(t *testing.T) {
	change := ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "test",
		NewConfig: map[string]interface{}{"key": "value"},
		Timestamp: time.Now(),
		Source:    "api",
	}

	json := change.ToJSON()
	if json == "" {
		t.Error("ToJSON should return non-empty string")
	}

	parsed, err := FromJSON(json)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if parsed.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", parsed.Name)
	}
}

func TestFromJSONInvalid(t *testing.T) {
	_, err := FromJSON("invalid json")
	if err == nil {
		t.Error("Should fail on invalid JSON")
	}
}

// Mock handler for testing
type mockHandler struct {
	name      string
	canReload bool
	reloadErr error
	rollbackErr error
}

func (h *mockHandler) Name() string { return h.name }
func (h *mockHandler) CanReload(change ConfigChange) bool { return h.canReload }
func (h *mockHandler) Reload(ctx context.Context, change ConfigChange) error { return h.reloadErr }
func (h *mockHandler) Rollback(ctx context.Context, change ConfigChange) error { return h.rollbackErr }

func TestNewManager(t *testing.T) {
	m := NewManager(nil)
	if m == nil {
		t.Fatal("Manager is nil")
	}
}

func TestManagerWithOptions(t *testing.T) {
	m := NewManager(nil,
		WithMaxHistory(50),
		WithConfigPath("config.yaml"),
	)

	if m.maxHistory != 50 {
		t.Error("MaxHistory option not applied")
	}

	if m.configPath != "config.yaml" {
		t.Error("ConfigPath option not applied")
	}
}

func TestManagerRegisterHandler(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	// Handler should be registered
}

func TestManagerUnregisterHandler(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})
	m.UnregisterHandler("test")

	// Handler should be unregistered
}

func TestManagerReload(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	change := ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "test",
		NewConfig: map[string]interface{}{"key": "value"},
		Timestamp: time.Now(),
	}

	err := m.Reload(context.Background(), change)
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}
}

func TestManagerReloadNoHandler(t *testing.T) {
	m := NewManager(nil)

	change := ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "nonexistent",
		NewConfig: map[string]interface{}{},
	}

	err := m.Reload(context.Background(), change)
	if err == nil {
		t.Error("Expected error for missing handler")
	}
}

func TestManagerReloadNotAllowed(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: false})

	change := ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "test",
		NewConfig: map[string]interface{}{},
	}

	err := m.Reload(context.Background(), change)
	if err == nil {
		t.Error("Expected error when reload not allowed")
	}
}

func TestManagerReloadError(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{
		name:      "test",
		canReload: true,
		reloadErr: errors.New("reload failed"),
	})

	change := ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "test",
		NewConfig: map[string]interface{}{},
	}

	err := m.Reload(context.Background(), change)
	if err == nil {
		t.Error("Expected error from reload")
	}
}

func TestManagerReloadWithRollback(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{
		name:        "test",
		canReload:   true,
		reloadErr:   errors.New("reload failed"),
		rollbackErr: nil,
	})

	change := ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "test",
		NewConfig: map[string]interface{}{},
	}

	err := m.Reload(context.Background(), change)
	if err == nil {
		t.Error("Expected error after rollback")
	}
}

func TestManagerRollback(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	change := ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "test",
		NewConfig: map[string]interface{}{},
	}

	err := m.Rollback(context.Background(), change)
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
}

func TestManagerApplyConfig(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	err := m.ApplyConfig(context.Background(), ConfigTypePlugin, "test", map[string]interface{}{"key": "value"})
	if err != nil {
		t.Fatalf("ApplyConfig failed: %v", err)
	}
}

func TestManagerGetHistory(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	m.ApplyConfig(context.Background(), ConfigTypePlugin, "test", map[string]interface{}{"key": "value"})

	history := m.GetHistory()
	if len(history) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(history))
	}
}

func TestManagerGetPending(t *testing.T) {
	m := NewManager(nil)
	pending := m.GetPending()

	if len(pending) != 0 {
		t.Error("Pending should be empty initially")
	}
}

func TestManagerStartStop(t *testing.T) {
	m := NewManager(nil)

	err := m.Start(context.Background())
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if !m.IsRunning() {
		t.Error("Should be running after start")
	}

	err = m.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if m.IsRunning() {
		t.Error("Should not be running after stop")
	}
}

func TestManagerDoubleStart(t *testing.T) {
	m := NewManager(nil)

	m.Start(context.Background())
	m.Start(context.Background()) // Should not panic

	m.Stop()
}

func TestManagerValidate(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	err := m.Validate(ConfigTypePlugin, "test", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

func TestManagerValidateNoHandler(t *testing.T) {
	m := NewManager(nil)

	err := m.Validate(ConfigTypePlugin, "nonexistent", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error for missing handler")
	}
}

func TestManagerValidateNotAllowed(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: false})

	err := m.Validate(ConfigTypePlugin, "test", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error when reload not allowed")
	}
}

func TestManagerExport(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})
	m.ApplyConfig(context.Background(), ConfigTypePlugin, "test", map[string]interface{}{"key": "value"})

	configs, err := m.Export()
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if len(configs) == 0 {
		t.Error("Export should return configs")
	}
}

func TestManagerImport(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})
	m.RegisterHandler(&mockHandler{name: "test2", canReload: true})

	configs := map[string]interface{}{
		"test":  map[string]interface{}{"key": "value"},
		"test2": map[string]interface{}{"key2": "value2"},
	}

	err := m.Import(context.Background(), configs)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}
}

func TestManagerImportError(t *testing.T) {
	m := NewManager(nil)
	// No handlers registered

	configs := map[string]interface{}{
		"test": map[string]interface{}{"key": "value"},
	}

	err := m.Import(context.Background(), configs)
	if err == nil {
		t.Error("Expected error for missing handler")
	}
}

func TestManagerOnChangeCallback(t *testing.T) {
	changes := []ConfigChange{}
	m := NewManager(nil, OnChange(func(change ConfigChange) {
		changes = append(changes, change)
	}))
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	m.ApplyConfig(context.Background(), ConfigTypePlugin, "test", map[string]interface{}{})

	if len(changes) != 1 {
		t.Error("OnChange callback should have been called")
	}
}

func TestManagerOnErrorCallback(t *testing.T) {
	errs := []error{}
	m := NewManager(nil, OnError(func(change ConfigChange, err error) {
		errs = append(errs, err)
	}))
	m.RegisterHandler(&mockHandler{
		name:      "test",
		canReload: true,
		reloadErr: errors.New("test error"),
	})

	m.ApplyConfig(context.Background(), ConfigTypePlugin, "test", map[string]interface{}{})

	if len(errs) == 0 {
		t.Error("OnError callback should have been called")
	}
}

func TestManagerMaxHistory(t *testing.T) {
	m := NewManager(nil, WithMaxHistory(3))
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	// Add more than max history
	for i := 0; i < 5; i++ {
		m.ApplyConfig(context.Background(), ConfigTypePlugin, "test", map[string]interface{}{"i": i})
	}

	history := m.GetHistory()
	if len(history) > 3 {
		t.Errorf("History should be limited to 3, got %d", len(history))
	}
}

// Mock watcher for testing
type mockWatcher struct {
	changes chan ConfigChange
}

func newMockWatcher() *mockWatcher {
	return &mockWatcher{
		changes: make(chan ConfigChange, 1),
	}
}

func (w *mockWatcher) Watch(ctx context.Context) (<-chan ConfigChange, error) {
	return w.changes, nil
}

func (w *mockWatcher) Close() error {
	close(w.changes)
	return nil
}

func (w *mockWatcher) Send(change ConfigChange) {
	w.changes <- change
}

func TestManagerWithWatcher(t *testing.T) {
	m := NewManager(nil)
	m.RegisterHandler(&mockHandler{name: "test", canReload: true})

	watcher := newMockWatcher()
	m.AddWatcher(watcher)

	ctx := context.Background()
	m.Start(ctx)

	// Send a change through the watcher
	watcher.Send(ConfigChange{
		Type:      ConfigTypePlugin,
		Name:      "test",
		NewConfig: map[string]interface{}{"watched": true},
	})

	time.Sleep(50 * time.Millisecond)

	m.Stop()
}

func TestErrors(t *testing.T) {
	if ErrHandlerNotFound == nil {
		t.Error("ErrHandlerNotFound should not be nil")
	}
	if ErrReloadNotAllowed == nil {
		t.Error("ErrReloadNotAllowed should not be nil")
	}
	if ErrRollbackFailed == nil {
		t.Error("ErrRollbackFailed should not be nil")
	}
}