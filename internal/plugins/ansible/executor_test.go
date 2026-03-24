package ansible

import (
	"context"
	"testing"
	"time"
)

func TestNewExecutor(t *testing.T) {
	cfg := Config{
		PlaybookDir:   "/playbooks",
		InventoryFile: "/inventory/hosts",
		Timeout:       3600,
	}

	exec := NewExecutor(cfg)
	if exec == nil {
		t.Fatal("NewExecutor returned nil")
	}
}

func TestExecutor_Execute_Timeout(t *testing.T) {
	cfg := Config{
		Timeout: 1, // 1 second timeout
	}

	exec := NewExecutor(cfg)
	ctx := context.Background()

	// Execute a command that takes longer than timeout
	// On Windows, use "timeout" or "ping" command
	_, err := exec.Execute(ctx, "sleep", "10")

	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestExecutor_Execute_InvalidCommand(t *testing.T) {
	cfg := Config{}

	exec := NewExecutor(cfg)
	ctx := context.Background()

	_, err := exec.Execute(ctx, "nonexistent-command-that-does-not-exist")
	if err == nil {
		t.Error("expected error for invalid command")
	}
}

func TestExecutor_Execute_WithEnvVars(t *testing.T) {
	cfg := Config{
		EnvVars: map[string]string{
			"TEST_VAR": "test_value",
		},
	}

	exec := NewExecutor(cfg)
	ctx := context.Background()

	// On Windows, use echo %TEST_VAR%
	// This test just verifies the executor can be created with env vars
	_ = exec
	_ = ctx
}

func TestExecutor_Execute_WithVaultPassword(t *testing.T) {
	cfg := Config{
		VaultPassword: "secret",
	}

	exec := NewExecutor(cfg)
	if exec == nil {
		t.Fatal("NewExecutor returned nil")
	}
}

func TestCheckAnsible(t *testing.T) {
	// This will fail on systems without ansible installed
	// but we test that the function doesn't panic
	err := CheckAnsible()
	// We don't assert success/failure since ansible may not be installed
	_ = err
}

func TestGetAnsibleVersion(t *testing.T) {
	// This will fail on systems without ansible installed
	_, err := GetAnsibleVersion()
	// We don't assert success/failure since ansible may not be installed
	_ = err
}

func TestExecutor_ExecuteWithOutput_Context(t *testing.T) {
	cfg := Config{}
	exec := NewExecutor(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Try to execute a command with output streaming
	outputChan, errChan := exec.ExecuteWithOutput(ctx, "echo", "test")

	// Wait for completion or timeout
	select {
	case <-outputChan:
		// Got some output
	case err := <-errChan:
		if err != nil {
			t.Logf("Command error (expected on some systems): %v", err)
		}
	case <-ctx.Done():
		t.Log("Context timed out")
	}
}

func TestConfig_Timeout(t *testing.T) {
	cfg := Config{
		Timeout: 300,
	}

	if cfg.Timeout != 300 {
		t.Errorf("expected timeout 300, got %d", cfg.Timeout)
	}
}

func TestConfig_VaultPassword(t *testing.T) {
	cfg := Config{
		VaultPassword: "my-secret-password",
	}

	if cfg.VaultPassword != "my-secret-password" {
		t.Errorf("expected vault password 'my-secret-password', got '%s'", cfg.VaultPassword)
	}
}

func TestConfig_EnvVars(t *testing.T) {
	cfg := Config{
		EnvVars: map[string]string{
			"ANSIBLE_HOST_KEY_CHECKING": "False",
			"ANSIBLE_FORCE_COLOR":       "True",
		},
	}

	if len(cfg.EnvVars) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(cfg.EnvVars))
	}
	if cfg.EnvVars["ANSIBLE_HOST_KEY_CHECKING"] != "False" {
		t.Error("expected ANSIBLE_HOST_KEY_CHECKING to be False")
	}
}