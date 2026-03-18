package plugin

import (
	"context"
	"testing"
)

// Mock plugin for testing
type mockPlugin struct {
	name string
}

func (m *mockPlugin) Info() PluginInfo {
	return PluginInfo{Name: m.name}
}

func (m *mockPlugin) Init(ctx context.Context, config map[string]interface{}) error {
	return nil
}

func (m *mockPlugin) Start(ctx context.Context) error {
	return nil
}

func (m *mockPlugin) Stop(ctx context.Context) error {
	return nil
}

func (m *mockPlugin) HealthCheck(ctx context.Context) error {
	return nil
}

func (m *mockPlugin) Capabilities() []string {
	return []string{"test"}
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if r.factories == nil {
		t.Error("factories map not initialized")
	}
}

func TestRegistry_RegisterFactory(t *testing.T) {
	r := NewRegistry()
	r.RegisterFactory("test", func() Plugin {
		return &mockPlugin{name: "test"}
	})

	if len(r.factories) != 1 {
		t.Errorf("expected 1 factory, got %d", len(r.factories))
	}
}

func TestRegistry_GetFactory(t *testing.T) {
	r := NewRegistry()
	r.RegisterFactory("test", func() Plugin {
		return &mockPlugin{name: "test"}
	})

	f, ok := r.GetFactory("test")
	if !ok {
		t.Fatal("factory not found")
	}
	if f == nil {
		t.Fatal("factory is nil")
	}

	_, ok = r.GetFactory("nonexistent")
	if ok {
		t.Error("should not find nonexistent factory")
	}
}

func TestRegistry_ListFactories(t *testing.T) {
	r := NewRegistry()
	r.RegisterFactory("plugin1", func() Plugin { return &mockPlugin{} })
	r.RegisterFactory("plugin2", func() Plugin { return &mockPlugin{} })
	r.RegisterFactory("plugin3", func() Plugin { return &mockPlugin{} })

	names := r.ListFactories()
	if len(names) != 3 {
		t.Errorf("expected 3 factories, got %d", len(names))
	}
}

func TestRegistry_Create(t *testing.T) {
	r := NewRegistry()
	r.RegisterFactory("test", func() Plugin {
		return &mockPlugin{name: "test-plugin"}
	})

	p, err := r.Create("test")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if p == nil {
		t.Fatal("created plugin is nil")
	}

	info := p.Info()
	if info.Name != "test-plugin" {
		t.Errorf("expected name 'test-plugin', got %s", info.Name)
	}
}

func TestRegistry_Create_NotFound(t *testing.T) {
	r := NewRegistry()

	_, err := r.Create("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent factory")
	}
}

func TestGlobalRegistry(t *testing.T) {
	// Test global registry functions
	RegisterFactory("global-test", func() Plugin {
		return &mockPlugin{name: "global"}
	})

	f, ok := GetFactory("global-test")
	if !ok {
		t.Fatal("global factory not found")
	}
	if f == nil {
		t.Fatal("global factory is nil")
	}

	names := ListFactories()
	if len(names) == 0 {
		t.Error("no factories in global registry")
	}

	p, err := Create("global-test")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if p == nil {
		t.Fatal("created plugin is nil")
	}
}