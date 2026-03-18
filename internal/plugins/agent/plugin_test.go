package agent

import (
	"context"
	"testing"
	"time"

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

	if info.Name != "agent" {
		t.Errorf("expected name 'agent', got %s", info.Name)
	}
	if info.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", info.Version)
	}
	if info.Description == "" {
		t.Error("description should not be empty")
	}
}

func TestCapabilities(t *testing.T) {
	p := New()
	caps := p.Capabilities()

	expected := []string{
		"connect",
		"disconnect",
		"exec",
		"upload",
		"download",
		"service_start",
		"service_stop",
		"service_status",
		"system_info",
		"list_agents",
	}

	if len(caps) != len(expected) {
		t.Errorf("expected %d capabilities, got %d", len(expected), len(caps))
	}
}

func TestInit(t *testing.T) {
	p := New()
	ctx := context.Background()

	err := p.Init(ctx, map[string]interface{}{
		"host": "192.168.1.1",
		"port": 8080,
		"timeout": 60,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if p.config.Timeout != 60 {
		t.Errorf("expected timeout 60, got %d", p.config.Timeout)
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
		"host": "192.168.1.1",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if p.config.Timeout != 30 {
		t.Errorf("expected default timeout 30, got %d", p.config.Timeout)
	}
}

func TestStartStop(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "192.168.1.1",
	})

	_ = p.Start(ctx)
	status, _ := p.GetStatus(ctx)
	if status.State != string(plugin.StateRunning) {
		t.Errorf("expected state running, got %s", status.State)
	}

	_ = p.Stop(ctx)
	status, _ = p.GetStatus(ctx)
	if status.State != string(plugin.StateStopped) {
		t.Errorf("expected state stopped, got %s", status.State)
	}
}

func TestHealthCheck_NotConnected(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})
	_ = p.Start(ctx)

	err := p.HealthCheck(ctx)
	if err == nil {
		t.Error("expected error when not connected")
	}
}

func TestExecute_UnknownAction(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	_, err := p.Execute(ctx, "unknown_action", nil)
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestExecute_Connect_MissingHost(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, err := p.Execute(ctx, "connect", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when missing host")
	}
}

func TestExecute_Connect_WithPort(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{"timeout": 1})

	result, err := p.Execute(ctx, "connect", map[string]interface{}{
		"host": "192.168.1.1",
		"port": 9090,
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Will fail due to connection, but test the path
	_ = result
}

func TestExecute_Exec_NotConnected(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, err := p.Execute(ctx, "exec", map[string]interface{}{
		"command": "ls",
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when not connected")
	}
}

func TestExecute_Exec_MissingCommand(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, err := p.Execute(ctx, "exec", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure when missing command")
	}
}

func TestExecute_Disconnect(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, err := p.Execute(ctx, "disconnect", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Error("expected success for disconnect")
	}
}

func TestExecute_Upload_NotImplemented(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, _ := p.Execute(ctx, "upload", nil)
	if result.Success {
		t.Error("expected failure - upload not implemented")
	}
}

func TestExecute_Download_NotImplemented(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, _ := p.Execute(ctx, "download", nil)
	if result.Success {
		t.Error("expected failure - download not implemented")
	}
}

func TestExecute_ServiceStart_NotConnected(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, _ := p.Execute(ctx, "service_start", map[string]interface{}{
		"service": "nginx",
	})
	if result.Success {
		t.Error("expected failure - not connected")
	}
}

func TestExecute_ServiceStart_MissingService(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, _ := p.Execute(ctx, "service_start", map[string]interface{}{})
	if result.Success {
		t.Error("expected failure - missing service")
	}
}

func TestExecute_ServiceStop_NotConnected(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, _ := p.Execute(ctx, "service_stop", map[string]interface{}{
		"service": "nginx",
	})
	if result.Success {
		t.Error("expected failure - not connected")
	}
}

func TestExecute_ServiceStatus_NotConnected(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, _ := p.Execute(ctx, "service_status", map[string]interface{}{
		"service": "nginx",
	})
	if result.Success {
		t.Error("expected failure - not connected")
	}
}

func TestExecute_SystemInfo_NotConnected(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	result, _ := p.Execute(ctx, "system_info", nil)
	if result.Success {
		t.Error("expected failure - not connected")
	}
}

func TestGetStatus(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	status, err := p.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if status.State != string(plugin.StateInitialized) {
		t.Errorf("expected state initialized, got %s", status.State)
	}
}

func TestGetStatus_AfterStart(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})
	_ = p.Start(ctx)

	status, err := p.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if status.State != string(plugin.StateRunning) {
		t.Errorf("expected state running, got %s", status.State)
	}
}

func TestGetStatus_AfterStop(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})
	_ = p.Start(ctx)
	_ = p.Stop(ctx)

	status, err := p.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if status.State != string(plugin.StateStopped) {
		t.Errorf("expected state stopped, got %s", status.State)
	}
}

func TestConfig(t *testing.T) {
	cfg := Config{
		Host:     "192.168.1.100",
		Port:     8080,
		APIToken: "test-token",
		Timeout:  60,
	}

	if cfg.Host != "192.168.1.100" {
		t.Errorf("expected host '192.168.1.100', got %s", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.APIToken != "test-token" {
		t.Errorf("expected token 'test-token', got %s", cfg.APIToken)
	}
	if cfg.Timeout != 60 {
		t.Errorf("expected timeout 60, got %d", cfg.Timeout)
	}
}

func TestStatus_InitialState(t *testing.T) {
	p := New()

	if p.status.State != string(plugin.StateUninitialized) {
		t.Errorf("expected initial state 'uninitialized', got %s", p.status.State)
	}

	if p.status.Health != string(plugin.HealthHealthy) {
		t.Errorf("expected initial health 'healthy', got %s", p.status.Health)
	}
}

func TestStatus_Metrics(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"host": "192.168.1.1",
	})

	status, _ := p.GetStatus(ctx)
	if status.Metrics == nil {
		t.Error("metrics map should be initialized")
	}
}

func TestLastUpdated(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})
	_ = p.Start(ctx)

	status, _ := p.GetStatus(ctx)
	if status.LastUpdated == 0 {
		t.Error("last updated should be set after start")
	}
}

func TestConcurrentAccess(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	// Test concurrent access
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			_, _ = p.GetStatus(ctx)
			_ = p.HealthCheck(ctx)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestExecute_Concurrent(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	// Test concurrent execute calls
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func() {
			_, _ = p.Execute(ctx, "disconnect", nil)
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}

func TestPluginImplementsInterface(t *testing.T) {
	// Verify that AgentPlugin implements the plugin interfaces
	var _ plugin.Plugin = New()
	var _ plugin.ExecutablePlugin = New()
	// Note: AgentPlugin doesn't implement ObservablePlugin (missing Watch method)
}

func TestStatusAfterMultipleStartStop(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	// Multiple start/stop cycles
	for i := 0; i < 3; i++ {
		_ = p.Start(ctx)
		status, _ := p.GetStatus(ctx)
		if status.State != string(plugin.StateRunning) {
			t.Errorf("cycle %d: expected state running, got %s", i, status.State)
		}

		_ = p.Stop(ctx)
		status, _ = p.GetStatus(ctx)
		if status.State != string(plugin.StateStopped) {
			t.Errorf("cycle %d: expected state stopped, got %s", i, status.State)
		}
	}
}

func TestTimeoutConfig(t *testing.T) {
	p := New()
	ctx := context.Background()

	// Test with very short timeout
	err := p.Init(ctx, map[string]interface{}{
		"host":    "192.168.1.1",
		"timeout": 1,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Timeout should be 1
	if p.config.Timeout != 1 {
		t.Errorf("expected timeout 1, got %d", p.config.Timeout)
	}
}

func TestTimeTracking(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	beforeStart := time.Now().Unix()
	_ = p.Start(ctx)
	afterStart := time.Now().Unix()

	status, _ := p.GetStatus(ctx)
	if status.LastUpdated < beforeStart || status.LastUpdated > afterStart {
		t.Errorf("last updated time %d not in expected range [%d, %d]",
			status.LastUpdated, beforeStart, afterStart)
	}
}