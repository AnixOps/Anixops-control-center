package container

import (
	"testing"
)

func TestNewContainer(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("Container is nil")
	}
}

func TestRegisterAndResolve(t *testing.T) {
	c := New()

	// Register a service
	service := "test-service"
	c.Register("test", service)

	// Resolve the service
	resolved, err := c.Resolve("test")
	if err != nil {
		t.Fatalf("Failed to resolve: %v", err)
	}

	if resolved != service {
		t.Errorf("Expected %v, got %v", service, resolved)
	}
}

func TestResolveNotFound(t *testing.T) {
	c := New()

	_, err := c.Resolve("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent service")
	}
}

func TestMustResolve(t *testing.T) {
	c := New()
	c.Register("test", "value")

	// Should not panic
	value := c.MustResolve("test")
	if value != "value" {
		t.Errorf("Expected 'value', got %v", value)
	}
}

func TestMustResolvePanic(t *testing.T) {
	c := New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nonexistent service")
		}
	}()

	c.MustResolve("nonexistent")
}

func TestRegisterFactory(t *testing.T) {
	c := New()
	c.RegisterFactory("test", func() interface{} {
		return "factory-value"
	})

	value, err := c.Resolve("test")
	if err != nil {
		t.Fatalf("Failed to resolve: %v", err)
	}

	if value != "factory-value" {
		t.Errorf("Expected 'factory-value', got %v", value)
	}
}

func TestRegisterSingleton(t *testing.T) {
	c := New()

	callCount := 0
	c.RegisterSingleton("test", func() interface{} {
		callCount++
		return "singleton-value"
	})

	// Resolve multiple times
	value1, _ := c.Resolve("test")
	value2, _ := c.Resolve("test")

	if value1 != value2 {
		t.Error("Singleton should return same instance")
	}

	if callCount != 1 {
		t.Errorf("Factory should be called once, called %d times", callCount)
	}
}

func TestHas(t *testing.T) {
	c := New()

	if c.Has("test") {
		t.Error("Should not have unregistered service")
	}

	c.Register("test", "value")

	if !c.Has("test") {
		t.Error("Should have registered service")
	}
}

func TestRemove(t *testing.T) {
	c := New()
	c.Register("test", "value")

	c.Remove("test")

	if c.Has("test") {
		t.Error("Service should be removed")
	}
}

func TestClear(t *testing.T) {
	c := New()
	c.Register("test1", "value1")
	c.Register("test2", "value2")

	c.Clear()

	if c.Has("test1") || c.Has("test2") {
		t.Error("All services should be cleared")
	}
}

func TestList(t *testing.T) {
	c := New()
	c.Register("service1", "value1")
	c.Register("service2", "value2")

	list := c.List()

	if len(list) != 2 {
		t.Errorf("Expected 2 services, got %d", len(list))
	}
}

func TestResolveInto(t *testing.T) {
	c := New()
	c.Register("test", "value")

	var target string
	err := c.ResolveInto("test", &target)
	if err != nil {
		t.Fatalf("Failed to resolve into: %v", err)
	}

	if target != "value" {
		t.Errorf("Expected 'value', got '%s'", target)
	}
}

func TestResolveIntoNonPointer(t *testing.T) {
	c := New()
	c.Register("test", "value")

	var target string
	err := c.ResolveInto("test", target) // not a pointer
	if err == nil {
		t.Error("Expected error for non-pointer target")
	}
}

func TestResolveIntoTypeMismatch(t *testing.T) {
	c := New()
	c.Register("test", "value")

	var target int
	err := c.ResolveInto("test", &target)
	if err == nil {
		t.Error("Expected error for type mismatch")
	}
}

// Global container tests
func TestGlobalRegister(t *testing.T) {
	Register("global-test", "value")

	if !Has("global-test") {
		t.Error("Global container should have service")
	}

	Remove("global-test")
}

func TestGlobalResolve(t *testing.T) {
	Register("global-test", "value")

	resolved, err := Resolve("global-test")
	if err != nil {
		t.Fatalf("Failed to resolve: %v", err)
	}

	if resolved != "value" {
		t.Errorf("Expected 'value', got %v", resolved)
	}

	Remove("global-test")
}

func TestGlobalMustResolve(t *testing.T) {
	Register("global-test", "value")

	value := MustResolve("global-test")
	if value != "value" {
		t.Errorf("Expected 'value', got %v", value)
	}

	Remove("global-test")
}

func TestGlobalList(t *testing.T) {
	Register("test1", "value1")
	Register("test2", "value2")

	list := List()
	if len(list) < 2 {
		t.Errorf("Expected at least 2 services, got %d", len(list))
	}

	Remove("test1")
	Remove("test2")
}

func TestGlobalClear(t *testing.T) {
	Register("clear-test", "value")
	Clear()

	if Has("clear-test") {
		t.Error("Service should be cleared")
	}
}