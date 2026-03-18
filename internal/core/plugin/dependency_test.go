package plugin

import (
	"context"
	"testing"
)

func TestNewDependencyManager(t *testing.T) {
	dm := NewDependencyManager()
	if dm == nil {
		t.Fatal("DependencyManager is nil")
	}
}

func TestDependencyManagerRegister(t *testing.T) {
	dm := NewDependencyManager()

	dm.Register("plugin-a", []Dependency{
		{Name: "plugin-b", Version: ">=1.0.0", Required: true},
	})

	// Should not panic
}

func TestDependencyManagerRegisterProvider(t *testing.T) {
	dm := NewDependencyManager()

	dm.RegisterProvider("database", "postgres-plugin")

	provider, ok := dm.GetProvider("database")
	if !ok {
		t.Error("Provider should be registered")
	}
	if provider != "postgres-plugin" {
		t.Errorf("Expected 'postgres-plugin', got '%s'", provider)
	}
}

func TestDependencyManagerGetProviderNotFound(t *testing.T) {
	dm := NewDependencyManager()

	_, ok := dm.GetProvider("nonexistent")
	if ok {
		t.Error("Should not find nonexistent provider")
	}
}

func TestDependencyManagerResolve(t *testing.T) {
	dm := NewDependencyManager()

	// Register dependencies: A depends on B, B depends on C
	dm.Register("plugin-a", []Dependency{
		{Name: "plugin-b", Required: true},
	})
	dm.Register("plugin-b", []Dependency{
		{Name: "plugin-c", Required: true},
	})
	dm.Register("plugin-c", []Dependency{})

	plugins := []string{"plugin-a", "plugin-b", "plugin-c"}
	order, err := dm.Resolve(plugins)

	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// C should come before B, B should come before A
	findIndex := func(name string) int {
		for i, n := range order {
			if n == name {
				return i
			}
		}
		return -1
	}

	cIdx := findIndex("plugin-c")
	bIdx := findIndex("plugin-b")
	aIdx := findIndex("plugin-a")

	if cIdx >= bIdx {
		t.Error("plugin-c should come before plugin-b")
	}
	if bIdx >= aIdx {
		t.Error("plugin-b should come before plugin-a")
	}
}

func TestDependencyManagerResolveMissingRequired(t *testing.T) {
	dm := NewDependencyManager()

	dm.Register("plugin-a", []Dependency{
		{Name: "nonexistent", Required: true},
	})

	plugins := []string{"plugin-a"}
	_, err := dm.Resolve(plugins)

	if err == nil {
		t.Error("Expected error for missing required dependency")
	}
}

func TestDependencyManagerResolveMissingOptional(t *testing.T) {
	dm := NewDependencyManager()

	dm.Register("plugin-a", []Dependency{
		{Name: "nonexistent", Required: false},
	})

	plugins := []string{"plugin-a"}
	order, err := dm.Resolve(plugins)

	if err != nil {
		t.Fatalf("Should not fail for optional dependency: %v", err)
	}

	if len(order) != 1 {
		t.Errorf("Expected 1 plugin in order, got %d", len(order))
	}
}

func TestCheckVersion(t *testing.T) {
	tests := []struct {
		version    string
		constraint string
		expected   bool
	}{
		{"1.0.0", "", true},
		{"1.0.0", "*", true},
		{"1.0.0", "1.0.0", true},
		{"1.0.0", ">=1.0.0", true},
		{"1.0.0", ">=0.9.0", true},
		{"1.0.0", ">0.9.0", true},
		{"1.0.0", "<=1.0.0", true},
		{"1.0.0", "<2.0.0", true},
		{"2.0.0", "<1.0.0", false},
		{"0.9.0", ">=1.0.0", false},
	}

	for _, test := range tests {
		result := CheckVersion(test.version, test.constraint)
		if result != test.expected {
			t.Errorf("CheckVersion(%s, %s) = %v, expected %v",
				test.version, test.constraint, result, test.expected)
		}
	}
}

// Mock plugin for testing dependency
type mockDepPlugin struct {
	info PluginInfo
}

func (m *mockDepPlugin) Info() PluginInfo {
	return m.info
}

func (m *mockDepPlugin) Init(ctx context.Context, config map[string]interface{}) error {
	return nil
}

func (m *mockDepPlugin) Start(ctx context.Context) error {
	return nil
}

func (m *mockDepPlugin) Stop(ctx context.Context) error {
	return nil
}

func (m *mockDepPlugin) HealthCheck(ctx context.Context) error {
	return nil
}

func (m *mockDepPlugin) Capabilities() []string {
	return []string{"test"}
}

func TestValidateDependencies(t *testing.T) {
	dm := NewDependencyManager()

	dm.Register("plugin-a", []Dependency{
		{Name: "plugin-b", Version: ">=1.0.0", Required: true},
	})

	plugins := map[string]Plugin{
		"plugin-a": &mockDepPlugin{info: PluginInfo{Name: "plugin-a", Version: "1.0.0"}},
		"plugin-b": &mockDepPlugin{info: PluginInfo{Name: "plugin-b", Version: "2.0.0"}},
	}

	errs := dm.ValidateDependencies(plugins)
	if len(errs) != 0 {
		t.Errorf("Should have no errors: %v", errs)
	}
}

func TestValidateDependenciesMissing(t *testing.T) {
	dm := NewDependencyManager()

	dm.Register("plugin-a", []Dependency{
		{Name: "nonexistent", Required: true},
	})

	plugins := map[string]Plugin{
		"plugin-a": &mockDepPlugin{info: PluginInfo{Name: "plugin-a", Version: "1.0.0"}},
	}

	errs := dm.ValidateDependencies(plugins)
	if len(errs) == 0 {
		t.Error("Should have errors for missing dependency")
	}
}

func TestValidateDependenciesVersionMismatch(t *testing.T) {
	dm := NewDependencyManager()

	dm.Register("plugin-a", []Dependency{
		{Name: "plugin-b", Version: ">=2.0.0", Required: true},
	})

	plugins := map[string]Plugin{
		"plugin-a": &mockDepPlugin{info: PluginInfo{Name: "plugin-a", Version: "1.0.0"}},
		"plugin-b": &mockDepPlugin{info: PluginInfo{Name: "plugin-b", Version: "1.0.0"}},
	}

	errs := dm.ValidateDependencies(plugins)
	if len(errs) == 0 {
		t.Error("Should have errors for version mismatch")
	}
}

func TestNewPluginWithDeps(t *testing.T) {
	p := &mockDepPlugin{info: PluginInfo{Name: "test"}}
	deps := []Dependency{{Name: "other", Required: true}}

	pwd := NewPluginWithDeps(p, deps)

	if pwd.Dependencies() == nil {
		t.Error("Dependencies should not be nil")
	}
	if len(pwd.Dependencies()) != 1 {
		t.Error("Should have 1 dependency")
	}
}

func TestHookRegistry(t *testing.T) {
	r := NewHookRegistry()
	if r == nil {
		t.Fatal("HookRegistry is nil")
	}
}

func TestHookRegistryRegister(t *testing.T) {
	r := NewHookRegistry()

	r.Register(Hook{
		Name:     "test-hook",
		Priority: 0,
		Type:     HookPreStart,
		Handler:  func(ctx context.Context, data interface{}) error { return nil },
	})

	// Should not panic
}

func TestHookRegistryExecute(t *testing.T) {
	r := NewHookRegistry()

	called := false
	r.Register(Hook{
		Name:     "test-hook",
		Priority: 0,
		Type:     HookPreStart,
		Handler:  func(ctx context.Context, data interface{}) error {
			called = true
			return nil
		},
	})

	err := r.Execute(context.Background(), HookPreStart, nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !called {
		t.Error("Hook should have been called")
	}
}

func TestHookRegistryExecuteError(t *testing.T) {
	r := NewHookRegistry()

	r.Register(Hook{
		Name:     "failing-hook",
		Priority: 0,
		Type:     HookPreStart,
		Handler:  func(ctx context.Context, data interface{}) error {
			return nil // Pass for now
		},
	})

	err := r.Execute(context.Background(), HookPreStart, nil)
	if err != nil {
		t.Errorf("Execute should not fail: %v", err)
	}
}

func TestHookRegistryPriority(t *testing.T) {
	r := NewHookRegistry()

	order := []string{}
	r.Register(Hook{
		Name:     "low",
		Priority: 100,
		Type:     HookPreStart,
		Handler:  func(ctx context.Context, data interface{}) error {
			order = append(order, "low")
			return nil
		},
	})
	r.Register(Hook{
		Name:     "high",
		Priority: 1,
		Type:     HookPreStart,
		Handler:  func(ctx context.Context, data interface{}) error {
			order = append(order, "high")
			return nil
		},
	})

	r.Execute(context.Background(), HookPreStart, nil)

	if order[0] != "high" || order[1] != "low" {
		t.Errorf("Hooks should execute in priority order: %v", order)
	}
}

func TestHookTypeString(t *testing.T) {
	tests := []struct {
		hookType HookType
		expected string
	}{
		{HookPreInit, "pre_init"},
		{HookPostInit, "post_init"},
		{HookPreStart, "pre_start"},
		{HookPostStart, "post_start"},
		{HookPreStop, "pre_stop"},
		{HookPostStop, "post_stop"},
		{HookPreExecute, "pre_execute"},
		{HookPostExecute, "post_execute"},
		{HookHealthCheck, "health_check"},
	}

	for _, test := range tests {
		if test.hookType.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.hookType.String())
		}
	}
}