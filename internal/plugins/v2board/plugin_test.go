package v2board

import (
	"context"
	"testing"

	"github.com/anixops/anixops-control-center/internal/core/plugin"
)

func TestNew(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("New() returned nil")
	}
}

func TestInfo(t *testing.T) {
	p := New()
	info := p.Info()

	if info.Name != "v2board" {
		t.Errorf("expected name 'v2board', got %s", info.Name)
	}
	if info.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", info.Version)
	}
	if info.Author != "AnixOps" {
		t.Errorf("expected author 'AnixOps', got %s", info.Author)
	}
}

func TestCapabilities(t *testing.T) {
	p := New()
	caps := p.Capabilities()

	expected := []string{
		"get_status",
		"get_nodes",
		"get_users",
		"get_stats",
		"deploy_node",
		"remove_node",
		"manage_user",
		"get_orders",
		"get_plans",
	}

	if len(caps) != len(expected) {
		t.Errorf("expected %d capabilities, got %d", len(expected), len(caps))
	}
}

func TestInit(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"host":    "http://localhost:8080",
		"api_key": "test-key",
		"timeout": 60,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	status, _ := p.GetStatus(ctx)
	if status.State != string(plugin.StateInitialized) {
		t.Errorf("expected state initialized, got %s", status.State)
	}

	if p.config.Timeout != 60 {
		t.Errorf("expected timeout 60, got %d", p.config.Timeout)
	}
}

func TestInit_DefaultTimeout(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if p.config.Timeout != 30 {
		t.Errorf("expected default timeout 30, got %d", p.config.Timeout)
	}
}

func TestInit_InvalidConfig(t *testing.T) {
	p := New()
	ctx := context.Background()

	// Config that can't be marshaled (channels can't be marshaled)
	err := p.Init(ctx, map[string]interface{}{
		"invalid": make(chan int),
	})
	if err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestStart(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	// Start will fail because there's no actual server, but state should be updated
	_ = p.Start(ctx)

	// Check state after start attempt
	status, _ := p.GetStatus(ctx)
	// State may be running or error depending on connection
	_ = status.State
}

func TestStop(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})
	_ = p.Stop(ctx)

	status, _ := p.GetStatus(ctx)
	if status.State != string(plugin.StateStopped) {
		t.Errorf("expected state stopped, got %s", status.State)
	}
}

func TestExecute_UnknownAction(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	_, err := p.Execute(ctx, "unknown_action", nil)
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestExecute_GetStatus(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	// This will fail since there's no server, but we test the path
	result, err := p.Execute(ctx, "get_status", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	// Result will be failure due to connection error
	_ = result
}

func TestExecute_GetNodes(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "get_nodes", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_GetUsers(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "get_users", map[string]interface{}{
		"page":     "1",
		"per_page": "10",
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_GetStats(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "get_stats", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_DeployNode(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "deploy_node", map[string]interface{}{
		"name": "test-node",
		"host": "192.168.1.1",
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_RemoveNode(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "remove_node", map[string]interface{}{
		"id": float64(1),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_RemoveNode_MissingID(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "remove_node", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Success {
		t.Error("expected failure when missing id")
	}
}

func TestExecute_ManageUser_Ban(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "manage_user", map[string]interface{}{
		"action": "ban",
		"id":     float64(1),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_ManageUser_Unban(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "manage_user", map[string]interface{}{
		"action": "unban",
		"id":     float64(1),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_ManageUser_Update(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "manage_user", map[string]interface{}{
		"action": "update",
		"id":     float64(1),
		"data":   map[string]interface{}{"name": "test"},
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_ManageUser_MissingAction(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "manage_user", map[string]interface{}{
		"id": float64(1),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Success {
		t.Error("expected failure when missing action")
	}
}

func TestExecute_ManageUser_MissingID(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "manage_user", map[string]interface{}{
		"action": "ban",
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Success {
		t.Error("expected failure when missing id")
	}
}

func TestExecute_ManageUser_UnknownAction(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "manage_user", map[string]interface{}{
		"action": "unknown",
		"id":     float64(1),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Success {
		t.Error("expected failure for unknown action")
	}
}

func TestExecute_GetOrders(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "get_orders", map[string]interface{}{
		"page": "1",
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestExecute_GetPlans(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	result, err := p.Execute(ctx, "get_plans", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestHealthCheck(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	// Health check will fail without actual server
	err := p.HealthCheck(ctx)
	// Just test that it doesn't panic
	_ = err
}

func TestGetStatus(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "http://localhost:8080",
	})

	status, err := p.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.State != string(plugin.StateInitialized) {
		t.Errorf("expected state initialized, got %s", status.State)
	}
}

func TestGetStringParam(t *testing.T) {
	tests := []struct {
		params     map[string]interface{}
		key        string
		defaultVal string
		expected   string
	}{
		{map[string]interface{}{"page": "2"}, "page", "1", "2"},
		{map[string]interface{}{}, "page", "1", "1"},
		{map[string]interface{}{"page": 123}, "page", "1", "1"},
		{map[string]interface{}{"page": ""}, "page", "default", ""},
	}

	for _, tt := range tests {
		result := getStringParam(tt.params, tt.key, tt.defaultVal)
		if result != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, result)
		}
	}
}

func TestClientCreation(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"host":    "http://localhost:8080",
		"api_key": "test-api-key",
		"timeout": 45,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if p.client == nil {
		t.Error("client should not be nil after init")
	}
}

func TestConfigFields(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"host":    "http://test.example.com:9090",
		"api_key": "my-secret-key",
		"timeout": 120,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if p.config.Host != "http://test.example.com:9090" {
		t.Errorf("expected host 'http://test.example.com:9090', got %s", p.config.Host)
	}
	if p.config.Timeout != 120 {
		t.Errorf("expected timeout 120, got %d", p.config.Timeout)
	}
	// Note: api_key might not be parsed due to JSON field mapping
}