package v2bx

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

	if info.Name != "v2bx" {
		t.Errorf("expected name 'v2bx', got %s", info.Name)
	}
	if info.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", info.Version)
	}
}

func TestCapabilities(t *testing.T) {
	p := New()
	caps := p.Capabilities()

	expected := []string{
		"get_nodes",
		"get_node_status",
		"get_users",
		"get_stats",
		"restart_node",
		"get_config",
		"get_logs",
	}

	if len(caps) != len(expected) {
		t.Errorf("expected %d capabilities, got %d", len(expected), len(caps))
	}
}

func TestInit(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "node1",
				"host": "192.168.1.1",
				"port": 8080,
			},
		},
		"timeout": 60,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if len(p.clients) != 1 {
		t.Errorf("expected 1 client, got %d", len(p.clients))
	}

	status, _ := p.GetStatus(ctx)
	if status.State != string(plugin.StateInitialized) {
		t.Errorf("expected state initialized, got %s", status.State)
	}
}

func TestInit_DefaultTimeout(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
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

	// Config with invalid type
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
		"hosts": []interface{}{},
	})
	_ = p.Start(ctx)

	status, _ := p.GetStatus(ctx)
	if status.State != string(plugin.StateRunning) {
		t.Errorf("expected state running, got %s", status.State)
	}
}

func TestStop(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
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
		"hosts": []interface{}{},
	})

	_, err := p.Execute(ctx, "unknown_action", nil)
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestExecute_GetNodes(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "node1",
				"host": "192.168.1.1",
				"port": 8080,
			},
			map[string]interface{}{
				"id":   2,
				"name": "node2",
				"host": "192.168.1.2",
				"port": 8080,
			},
		},
	})

	result, err := p.Execute(ctx, "get_nodes", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Error("expected success for get_nodes")
	}

	if result.Metrics["total"] != 2 {
		t.Errorf("expected 2 nodes, got %v", result.Metrics["total"])
	}
}

func TestExecute_GetNodeStatus_MissingNodeID(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_node_status", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when missing node_id")
	}
}

func TestExecute_GetNodeStatus_NodeNotFound(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_node_status", map[string]interface{}{
		"node_id": float64(999),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when node not found")
	}
}

func TestExecute_GetUsers_MissingNodeID(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_users", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when missing node_id")
	}
}

func TestExecute_GetUsers_NodeNotFound(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_users", map[string]interface{}{
		"node_id": float64(999),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when node not found")
	}
}

func TestExecute_GetStats_AllNodes(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_stats", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Error("expected success for get_stats")
	}
}

func TestExecute_GetStats_SingleNode(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_stats", map[string]interface{}{
		"node_id": float64(999),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when node not found")
	}
}

func TestExecute_RestartNode_MissingNodeID(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "restart_node", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when missing node_id")
	}
}

func TestExecute_RestartNode_NodeNotFound(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "restart_node", map[string]interface{}{
		"node_id": float64(999),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when node not found")
	}
}

func TestExecute_GetConfig_MissingNodeID(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_config", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when missing node_id")
	}
}

func TestExecute_GetConfig_NodeNotFound(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_config", map[string]interface{}{
		"node_id": float64(999),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when node not found")
	}
}

func TestExecute_GetLogs_MissingNodeID(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_logs", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when missing node_id")
	}
}

func TestExecute_GetLogs_NodeNotFound(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	result, err := p.Execute(ctx, "get_logs", map[string]interface{}{
		"node_id": float64(999),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when node not found")
	}
}

func TestHealthCheck(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})
	_ = p.Start(ctx)

	// Should return error if unhealthy (no nodes)
	err := p.HealthCheck(ctx)
	if err == nil {
		t.Error("expected error for unhealthy plugin")
	}
}

func TestGetStatus(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})

	status, err := p.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.State != string(plugin.StateInitialized) {
		t.Errorf("expected state initialized, got %s", status.State)
	}
}

func TestMultipleNodes(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "node1",
				"host": "192.168.1.1",
				"port": 8080,
			},
			map[string]interface{}{
				"id":   2,
				"name": "node2",
				"host": "192.168.1.2",
				"port": 8080,
			},
			map[string]interface{}{
				"id":   3,
				"name": "node3",
				"host": "192.168.1.3",
				"port": 8080,
			},
		},
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if len(p.clients) != 3 {
		t.Errorf("expected 3 clients, got %d", len(p.clients))
	}

	// Verify all clients are created
	for _, id := range []int{1, 2, 3} {
		if _, ok := p.clients[id]; !ok {
			t.Errorf("client %d not found", id)
		}
	}
}

func TestNodeConfig(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "node1",
				"host": "192.168.1.1",
				"port": 8080,
			},
		},
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if len(p.config.Hosts) != 1 {
		t.Fatalf("expected 1 host config, got %d", len(p.config.Hosts))
	}

	node := p.config.Hosts[0]
	if node.ID != 1 {
		t.Errorf("expected ID 1, got %d", node.ID)
	}
	if node.Name != "node1" {
		t.Errorf("expected name 'node1', got %s", node.Name)
	}
	if node.Host != "192.168.1.1" {
		t.Errorf("expected host '192.168.1.1', got %s", node.Host)
	}
	if node.Port != 8080 {
		t.Errorf("expected port 8080, got %d", node.Port)
	}
	// Note: api_key and grpc_port might not be parsed due to JSON field mapping
}

func TestStartHealthStatus(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{},
	})
	_ = p.Start(ctx)

	status, _ := p.GetStatus(ctx)

	// With no nodes, should be unhealthy
	if status.Health != string(plugin.HealthUnhealthy) {
		t.Errorf("expected unhealthy status with no nodes, got %s", status.Health)
	}
}

func TestMetricsUpdated(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"hosts": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "node1",
				"host": "192.168.1.1",
				"port": 8080,
			},
		},
	})

	status, _ := p.GetStatus(ctx)

	if status.Metrics["nodes"] != 1 {
		t.Errorf("expected nodes metric to be 1, got %v", status.Metrics["nodes"])
	}
}