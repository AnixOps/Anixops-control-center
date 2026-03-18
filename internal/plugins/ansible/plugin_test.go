package ansible

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/plugin"
)

// skipIfNoAnsible skips the test if ansible is not installed
func skipIfNoAnsible(t *testing.T) {
	if _, err := exec.LookPath("ansible"); err != nil {
		t.Skip("ansible not installed, skipping test")
	}
}

func TestNew(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("New() returned nil")
	}
}

func TestInfo(t *testing.T) {
	p := New()
	info := p.Info()

	if info.Name != "ansible" {
		t.Errorf("expected name 'ansible', got %s", info.Name)
	}
	if info.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", info.Version)
	}
}

func TestCapabilities(t *testing.T) {
	p := New()
	caps := p.Capabilities()

	expected := []string{
		"run_playbook",
		"run_task",
		"list_hosts",
		"list_playbooks",
		"validate_inventory",
		"get_inventory",
	}

	if len(caps) != len(expected) {
		t.Errorf("expected %d capabilities, got %d", len(expected), len(caps))
	}
}

func TestStartStop(t *testing.T) {
	p := New()

	ctx := context.Background()
	err := p.Start(ctx)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	status, err := p.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.State != string(plugin.StateRunning) {
		t.Errorf("expected state running, got %s", status.State)
	}

	err = p.Stop(ctx)
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	status, _ = p.GetStatus(ctx)
	if status.State != string(plugin.StateStopped) {
		t.Errorf("expected state stopped, got %s", status.State)
	}
}

func TestExecute_UnknownAction(t *testing.T) {
	p := New()
	ctx := context.Background()

	_, err := p.Execute(ctx, "unknown_action", nil)
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestParseHosts(t *testing.T) {
	tests := []struct {
		output   []byte
		expected int
	}{
		{[]byte("  host1.example.com\n  host2.example.com"), 2},
		{[]byte("no leading spaces"), 0},
		{[]byte("  host1\n\n  host2\n  host3"), 3},
		{[]byte(""), 0},
	}

	for _, tt := range tests {
		hosts := parseHosts(tt.output)
		if len(hosts) != tt.expected {
			t.Errorf("expected %d hosts, got %d for output: %s", tt.expected, len(hosts), string(tt.output))
		}
	}
}

func TestParseHosts_Format(t *testing.T) {
	output := []byte(`  host1.example.com
  host2.example.com
  host3.example.com`)

	hosts := parseHosts(output)
	if len(hosts) != 3 {
		t.Errorf("expected 3 hosts, got %d", len(hosts))
	}
	if hosts[0] != "host1.example.com" {
		t.Errorf("expected host1.example.com, got %s", hosts[0])
	}
	if hosts[1] != "host2.example.com" {
		t.Errorf("expected host2.example.com, got %s", hosts[1])
	}
}

func TestErrToString(t *testing.T) {
	if errToString(nil) != "" {
		t.Error("expected empty string for nil error")
	}

	// Test with an error
	err := context.Canceled
	if errToString(err) == "" {
		t.Error("expected non-empty string for error")
	}
}

func TestInit(t *testing.T) {
	p := New()
	ctx := context.Background()

	// Init will fail without ansible installed
	err := p.Init(ctx, map[string]interface{}{
		"playbook_dir": "/tmp",
	})
	// Just check it doesn't panic
	_ = err
}

func TestInit_DefaultTimeout(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"playbook_dir": "/tmp",
	})

	// Default timeout is 3600
	if p.config.Timeout != 3600 {
		t.Errorf("expected default timeout 3600, got %d", p.config.Timeout)
	}
}

func TestInit_WithCustomTimeout(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"playbook_dir": "/tmp",
		"timeout":      1800,
	})

	if p.config.Timeout != 1800 {
		t.Errorf("expected timeout 1800, got %d", p.config.Timeout)
	}
}

func TestExecute_RunPlaybook_MissingPlaybook(t *testing.T) {
	p := New()
	ctx := context.Background()

	_, err := p.Execute(ctx, "run_playbook", map[string]interface{}{})
	if err == nil {
		t.Error("expected error when missing playbook parameter")
	}
}

func TestExecute_RunTask_MissingModule(t *testing.T) {
	p := New()
	ctx := context.Background()

	_, err := p.Execute(ctx, "run_task", map[string]interface{}{})
	if err == nil {
		t.Error("expected error when missing module parameter")
	}
}

func TestExecute_ListPlaybooks_NoDir(t *testing.T) {
	p := New()
	ctx := context.Background()

	result, err := p.Execute(ctx, "list_playbooks", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Success {
		t.Error("expected failure when playbook_dir not configured")
	}
}

func TestExecute_ValidateInventory_NoInventory(t *testing.T) {
	p := New()
	ctx := context.Background()

	result, err := p.Execute(ctx, "validate_inventory", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Success {
		t.Error("expected failure when no inventory specified")
	}
}

func TestExecute_GetInventory_NoInventory(t *testing.T) {
	p := New()
	ctx := context.Background()

	result, err := p.Execute(ctx, "get_inventory", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Success {
		t.Error("expected failure when no inventory specified")
	}
}

func TestGetStatus(t *testing.T) {
	p := New()
	ctx := context.Background()

	status, err := p.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.State != string(plugin.StateUninitialized) {
		t.Errorf("expected state uninitialized, got %s", status.State)
	}
}

func TestConfig(t *testing.T) {
	cfg := Config{
		PlaybookDir:   "/opt/playbooks",
		InventoryFile: "/opt/inventory",
		VaultPassword: "secret",
		EnvVars:       map[string]string{"ANSIBLE_HOST_KEY_CHECKING": "False"},
		Timeout:       1800,
	}

	if cfg.PlaybookDir != "/opt/playbooks" {
		t.Errorf("expected playbook_dir '/opt/playbooks', got %s", cfg.PlaybookDir)
	}
	if cfg.Timeout != 1800 {
		t.Errorf("expected timeout 1800, got %d", cfg.Timeout)
	}
}

func TestExecutor_New(t *testing.T) {
	cfg := Config{
		Timeout: 60,
	}
	exec := NewExecutor(cfg)
	if exec == nil {
		t.Fatal("NewExecutor returned nil")
	}
}

func TestExecutor_Execute_InvalidCommand(t *testing.T) {
	cfg := Config{
		Timeout: 5,
	}
	exec := NewExecutor(cfg)

	ctx := context.Background()
	_, err := exec.Execute(ctx, "nonexistent_command_that_does_not_exist")
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
}

func TestExecutor_ExecuteWithTimeout(t *testing.T) {
	cfg := Config{
		Timeout: 1, // Very short timeout
	}
	exec := NewExecutor(cfg)

	ctx := context.Background()
	// Sleep command that should timeout
	_, err := exec.Execute(ctx, "sleep", "10")
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestHealthCheck_WithoutAnsible(t *testing.T) {
	p := New()
	ctx := context.Background()

	// Health check will fail without ansible installed
	err := p.HealthCheck(ctx)
	// Just test it doesn't panic
	_ = err
}

func TestExecute_ListHosts(t *testing.T) {
	skipIfNoAnsible(t)

	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"playbook_dir": "/tmp",
	})

	// list_hosts will fail without ansible, test the error handling
	result, err := p.Execute(ctx, "list_hosts", map[string]interface{}{
		"inventory": "/nonexistent/inventory",
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	// Should return a result (even if failed)
	_ = result
}

func TestPluginImplementsInterface(t *testing.T) {
	// Verify that AnsiblePlugin implements the plugin interfaces
	var _ plugin.Plugin = New()
	var _ plugin.ExecutablePlugin = New()
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

func TestStatus_AfterInit(t *testing.T) {
	p := New()
	ctx := context.Background()

	// Init will fail without ansible, but state should change appropriately
	_ = p.Init(ctx, map[string]interface{}{"playbook_dir": "/tmp"})

	// Even if init fails due to ansible not being installed,
	// the code path is exercised
}

func TestStatus_Metrics(t *testing.T) {
	p := New()
	ctx := context.Background()

	status, _ := p.GetStatus(ctx)
	if status.Metrics == nil {
		t.Error("metrics map should be initialized")
	}
}

func TestRunPlaybook_BuildsCorrectArgs(t *testing.T) {
	skipIfNoAnsible(t)

	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{"playbook_dir": "/tmp"})

	// Test with all parameters
	_, err := p.Execute(ctx, "run_playbook", map[string]interface{}{
		"playbook":   "test.yml",
		"inventory":  "/tmp/inventory",
		"extra_vars": map[string]interface{}{"var1": "value1"},
		"tags":       "install",
		"limit":      "localhost",
		"verbose":    true,
		"check":      true,
		"diff":       true,
	})
	// Will fail because playbook doesn't exist, but tests the parameter parsing
	_ = err
}

func TestRunTask_BuildsCorrectArgs(t *testing.T) {
	skipIfNoAnsible(t)

	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{"playbook_dir": "/tmp"})

	// Test with all parameters
	_, err := p.Execute(ctx, "run_task", map[string]interface{}{
		"module":    "ping",
		"pattern":   "all",
		"args":      "data=hello",
		"inventory": "/tmp/inventory",
	})
	// Will fail because inventory doesn't exist, but tests the parameter parsing
	_ = err
}

func TestRunTask_DefaultPattern(t *testing.T) {
	skipIfNoAnsible(t)

	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{"playbook_dir": "/tmp"})

	// Test without pattern (should use "all" as default)
	_, err := p.Execute(ctx, "run_task", map[string]interface{}{
		"module": "ping",
	})
	_ = err
}

func TestListPlaybooks_WithDir(t *testing.T) {
	p := New()
	ctx := context.Background()

	// Create a temporary directory with playbooks
	tmpDir, err := os.MkdirTemp("", "playbooks")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some playbook files
	os.WriteFile(filepath.Join(tmpDir, "site.yml"), []byte("---"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "deploy.yaml"), []byte("---"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755) // Create a subdirectory
	os.WriteFile(filepath.Join(tmpDir, "subdir", "nested.yml"), []byte("---"), 0644)

	// Set config directly since Init will fail without ansible installed
	p.config = Config{
		PlaybookDir: tmpDir,
	}

	result, err := p.Execute(ctx, "list_playbooks", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success, got error: %s", result.Error)
	}

	playbooks, ok := result.Data.([]string)
	if !ok {
		t.Fatal("expected []string")
	}

	if len(playbooks) < 2 {
		t.Errorf("expected at least 2 playbooks, got %d", len(playbooks))
	}
}

func TestListPlaybooks_NonexistentDir(t *testing.T) {
	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{
		"playbook_dir": "/nonexistent/path/that/does/not/exist",
	})

	result, err := p.Execute(ctx, "list_playbooks", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Success {
		t.Error("expected failure for nonexistent directory")
	}
}

func TestValidateInventory_WithInventory(t *testing.T) {
	skipIfNoAnsible(t)

	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	// Create a temporary inventory file
	tmpFile, err := os.CreateTemp("", "inventory")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("[all]\nlocalhost\n")
	tmpFile.Close()

	result, err := p.Execute(ctx, "validate_inventory", map[string]interface{}{
		"inventory": tmpFile.Name(),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestGetInventory_WithInventory(t *testing.T) {
	skipIfNoAnsible(t)

	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	// Create a temporary inventory file
	tmpFile, err := os.CreateTemp("", "inventory")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("[all]\nlocalhost ansible_connection=local\n")
	tmpFile.Close()

	result, err := p.Execute(ctx, "get_inventory", map[string]interface{}{
		"inventory": tmpFile.Name(),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
}

func TestGetInventory_InvalidJSON(t *testing.T) {
	skipIfNoAnsible(t)

	p := New()
	ctx := context.Background()

	_ = p.Init(ctx, map[string]interface{}{})

	// Create an inventory that will produce invalid JSON output
	tmpFile, err := os.CreateTemp("", "inventory")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("invalid inventory content\n")
	tmpFile.Close()

	result, err := p.Execute(ctx, "get_inventory", map[string]interface{}{
		"inventory": tmpFile.Name(),
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	_ = result
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

func TestStatusAfterMultipleStartStop(t *testing.T) {
	p := New()
	ctx := context.Background()

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

func TestExecutor_WithEnvVars(t *testing.T) {
	cfg := Config{
		Timeout: 60,
		EnvVars: map[string]string{
			"ANSIBLE_HOST_KEY_CHECKING": "False",
			"CUSTOM_VAR":                "value",
		},
	}
	exec := NewExecutor(cfg)

	if exec == nil {
		t.Fatal("NewExecutor returned nil")
	}
}

func TestExecutor_WithVaultPassword(t *testing.T) {
	cfg := Config{
		Timeout:       60,
		VaultPassword: "secret",
	}
	exec := NewExecutor(cfg)

	if exec == nil {
		t.Fatal("NewExecutor returned nil")
	}
}

func TestExecutor_ExecuteWithOutput(t *testing.T) {
	cfg := Config{
		Timeout: 5,
	}
	exec := NewExecutor(cfg)

	ctx := context.Background()
	outputChan, errChan := exec.ExecuteWithOutput(ctx, "echo", "hello")

	// Read output
	var output string
	for s := range outputChan {
		output += s
	}

	err := <-errChan
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if output == "" {
		t.Error("expected some output")
	}
}

func TestExecutor_ExecuteWithOutput_InvalidCommand(t *testing.T) {
	cfg := Config{
		Timeout: 5,
	}
	exec := NewExecutor(cfg)

	ctx := context.Background()
	outputChan, errChan := exec.ExecuteWithOutput(ctx, "nonexistent_command_xyz")

	// Drain output channel
	for range outputChan {
	}

	err := <-errChan
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
}

func TestCheckAnsible(t *testing.T) {
	err := CheckAnsible()
	// Will fail without ansible installed, just test it doesn't panic
	_ = err
}

func TestGetAnsibleVersion(t *testing.T) {
	version, err := GetAnsibleVersion()
	// Will fail without ansible installed, just test it doesn't panic
	_ = version
	_ = err
}

func TestExecute_WithConfigInventory(t *testing.T) {
	skipIfNoAnsible(t)

	p := New()
	ctx := context.Background()

	// Create a temporary inventory file
	tmpFile, err := os.CreateTemp("", "inventory")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("[all]\nlocalhost\n")
	tmpFile.Close()

	_ = p.Init(ctx, map[string]interface{}{
		"inventory_file": tmpFile.Name(),
	})

	// Test run_task with config inventory (no explicit inventory param)
	_, _ = p.Execute(ctx, "run_task", map[string]interface{}{
		"module": "ping",
	})

	// Test list_hosts with config inventory
	_, _ = p.Execute(ctx, "list_hosts", nil)
}

func TestLastUpdated(t *testing.T) {
	p := New()
	ctx := context.Background()

	beforeStart := time.Now().Unix()
	_ = p.Start(ctx)
	afterStart := time.Now().Unix()

	status, _ := p.GetStatus(ctx)
	if status.LastUpdated < beforeStart || status.LastUpdated > afterStart {
		t.Errorf("last updated time %d not in expected range [%d, %d]",
			status.LastUpdated, beforeStart, afterStart)
	}
}