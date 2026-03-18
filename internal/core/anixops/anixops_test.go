package anixops

import (
	"context"
	"testing"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/config"
	"github.com/anixops/anixops-control-center/internal/core/logger"
	"github.com/anixops/anixops-control-center/internal/core/plugin"
)

// MockPlugin is a mock plugin for testing
type MockPlugin struct {
	name    string
	version string
	started bool
}

func (m *MockPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name:    m.name,
		Version: m.version,
	}
}

func (m *MockPlugin) Init(ctx context.Context, config map[string]interface{}) error {
	return nil
}

func (m *MockPlugin) Start(ctx context.Context) error {
	m.started = true
	return nil
}

func (m *MockPlugin) Stop(ctx context.Context) error {
	m.started = false
	return nil
}

func (m *MockPlugin) HealthCheck(ctx context.Context) error {
	return nil
}

func (m *MockPlugin) Capabilities() []string {
	return []string{"test"}
}

// ExecutableMockPlugin is a mock plugin that can execute
type ExecutableMockPlugin struct {
	MockPlugin
}

func (m *ExecutableMockPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (plugin.Result, error) {
	return plugin.Result{Success: true, Data: action}, nil
}

func TestNewSDK(t *testing.T) {
	sdk, err := New(nil)
	if err != nil {
		t.Fatalf("Failed to create SDK: %v", err)
	}

	if sdk == nil {
		t.Fatal("SDK is nil")
	}

	if sdk.Config() == nil {
		t.Fatal("Config is nil")
	}

	if sdk.Container() == nil {
		t.Fatal("Container is nil")
	}

	if sdk.PluginManager() == nil {
		t.Fatal("PluginManager is nil")
	}

	if sdk.EventBus() == nil {
		t.Fatal("EventBus is nil")
	}

	if sdk.Logger() == nil {
		t.Fatal("Logger is nil")
	}

	if sdk.Metrics() == nil {
		t.Fatal("Metrics is nil")
	}
}

func TestSDKWithConfig(t *testing.T) {
	opts := &Options{
		LogLevel:        "debug",
		JSONLogging:     true,
		ShutdownTimeout: 60,
	}

	sdk, err := New(opts)
	if err != nil {
		t.Fatalf("Failed to create SDK: %v", err)
	}

	if sdk.Config() == nil {
		t.Fatal("Config is nil")
	}
}

func TestSDKWithPreloadedConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Server.Port = 9999

	opts := &Options{
		Config: cfg,
	}

	sdk, err := New(opts)
	if err != nil {
		t.Fatalf("Failed to create SDK: %v", err)
	}

	if sdk.Config().Server.Port != 9999 {
		t.Errorf("Expected port 9999, got %d", sdk.Config().Server.Port)
	}
}

func TestSDKWithCustomLogger(t *testing.T) {
	customLogger := logger.New(logger.Options{Level: logger.DebugLevel})

	opts := &Options{
		Logger: customLogger,
	}

	sdk, err := New(opts)
	if err != nil {
		t.Fatalf("Failed to create SDK: %v", err)
	}

	if sdk.Logger() != customLogger {
		t.Error("Should use custom logger")
	}
}

func TestRegisterPlugin(t *testing.T) {
	sdk, _ := New(nil)

	mockPlugin := &MockPlugin{name: "test-plugin", version: "1.0.0"}

	err := sdk.RegisterPlugin("test", mockPlugin)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}

	// Verify plugin is registered
	_, ok := sdk.PluginManager().Get("test")
	if !ok {
		t.Fatal("Plugin not found after registration")
	}
}

func TestRegisterDuplicatePlugin(t *testing.T) {
	sdk, _ := New(nil)

	mockPlugin := &MockPlugin{name: "test-plugin", version: "1.0.0"}

	sdk.RegisterPlugin("test", mockPlugin)
	err := sdk.RegisterPlugin("test", mockPlugin)

	if err == nil {
		t.Error("Expected error for duplicate plugin")
	}
}

func TestUnregisterPlugin(t *testing.T) {
	sdk, _ := New(nil)

	mockPlugin := &MockPlugin{name: "test-plugin", version: "1.0.0"}
	sdk.RegisterPlugin("test", mockPlugin)

	err := sdk.UnregisterPlugin("test")
	if err != nil {
		t.Fatalf("Failed to unregister plugin: %v", err)
	}

	_, ok := sdk.PluginManager().Get("test")
	if ok {
		t.Error("Plugin should be unregistered")
	}
}

func TestUnregisterNonexistentPlugin(t *testing.T) {
	sdk, _ := New(nil)

	err := sdk.UnregisterPlugin("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent plugin")
	}
}

func TestStartStop(t *testing.T) {
	sdk, _ := New(nil)

	mockPlugin := &MockPlugin{name: "test-plugin", version: "1.0.0"}
	sdk.RegisterPlugin("test", mockPlugin)

	ctx := context.Background()

	// Start
	err := sdk.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start SDK: %v", err)
	}

	if !sdk.IsRunning() {
		t.Fatal("SDK should be running")
	}

	if !mockPlugin.started {
		t.Fatal("Plugin should be started")
	}

	// Stop
	err = sdk.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop SDK: %v", err)
	}

	if sdk.IsRunning() {
		t.Fatal("SDK should not be running")
	}

	if mockPlugin.started {
		t.Fatal("Plugin should be stopped")
	}
}

func TestDoubleStart(t *testing.T) {
	sdk, _ := New(nil)

	mockPlugin := &MockPlugin{name: "test-plugin", version: "1.0.0"}
	sdk.RegisterPlugin("test", mockPlugin)

	ctx := context.Background()

	// First start
	err := sdk.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start SDK: %v", err)
	}

	// Second start should fail
	err = sdk.Start(ctx)
	if err == nil {
		t.Fatal("Expected error on double start")
	}
}

func TestStopWhenNotRunning(t *testing.T) {
	sdk, _ := New(nil)

	// Stop when not running should return nil
	err := sdk.Stop(context.Background())
	if err != nil {
		t.Fatalf("Stop should not fail when not running: %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	sdk, _ := New(nil)

	mockPlugin := &MockPlugin{name: "test-plugin", version: "1.0.0"}
	sdk.RegisterPlugin("test", mockPlugin)

	ctx := context.Background()
	sdk.Start(ctx)

	results := sdk.HealthCheck(ctx)
	if len(results) == 0 {
		t.Fatal("Expected health check results")
	}

	for name, err := range results {
		if err != nil {
			t.Errorf("Plugin %s health check failed: %v", name, err)
		}
	}
}

func TestGetPluginInfo(t *testing.T) {
	sdk, _ := New(nil)

	mockPlugin := &MockPlugin{name: "test-plugin", version: "1.0.0"}
	sdk.RegisterPlugin("test", mockPlugin)

	infos := sdk.GetPluginInfo()
	if len(infos) != 1 {
		t.Fatalf("Expected 1 plugin info, got %d", len(infos))
	}

	info, ok := infos["test"]
	if !ok {
		t.Fatal("Plugin info not found")
	}

	if info.Name != "test-plugin" {
		t.Errorf("Expected name 'test-plugin', got '%s'", info.Name)
	}

	if info.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", info.Version)
	}
}

func TestDefaultSDK(t *testing.T) {
	sdk := Default()
	if sdk == nil {
		t.Fatal("Default SDK is nil")
	}

	// Second call should return same instance
	sdk2 := Default()
	if sdk != sdk2 {
		t.Fatal("Default SDK should return same instance")
	}
}

func TestQuickStart(t *testing.T) {
	sdk, err := QuickStart()
	if err != nil {
		t.Fatalf("QuickStart failed: %v", err)
	}

	if !sdk.IsRunning() {
		t.Fatal("SDK should be running after QuickStart")
	}

	// Clean up
	sdk.Stop(context.Background())
}

func TestExecutePluginNotFound(t *testing.T) {
	sdk, _ := New(nil)
	sdk.Start(context.Background())

	_, err := sdk.ExecutePlugin(context.Background(), "nonexistent", "test", nil)
	if err == nil {
		t.Fatal("Expected error for nonexistent plugin")
	}
}

func TestExecutePluginNotExecutable(t *testing.T) {
	sdk, _ := New(nil)

	mockPlugin := &MockPlugin{name: "test-plugin", version: "1.0.0"}
	sdk.RegisterPlugin("test", mockPlugin)
	sdk.Start(context.Background())

	_, err := sdk.ExecutePlugin(context.Background(), "test", "action", nil)
	if err == nil {
		t.Fatal("Expected error for non-executable plugin")
	}
}

func TestExecutePlugin(t *testing.T) {
	sdk, _ := New(nil)

	execPlugin := &ExecutableMockPlugin{MockPlugin{name: "test-plugin", version: "1.0.0"}}
	sdk.RegisterPlugin("test", execPlugin)
	sdk.Start(context.Background())

	result, err := sdk.ExecutePlugin(context.Background(), "test", "test-action", map[string]interface{}{"key": "value"})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success")
	}
}

func TestVersion(t *testing.T) {
	if Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", Version)
	}
}

func TestShutdownTimeout(t *testing.T) {
	opts := &Options{
		ShutdownTimeout: 10,
	}

	sdk, err := New(opts)
	if err != nil {
		t.Fatalf("Failed to create SDK: %v", err)
	}

	// Verify shutdown handler is set
	if sdk.shutdown == nil {
		t.Fatal("Shutdown handler is nil")
	}
}

func TestWaitForShutdown(t *testing.T) {
	sdk, _ := New(nil)
	sdk.Start(context.Background())

	// Stop immediately
	go func() {
		time.Sleep(50 * time.Millisecond)
		sdk.Stop(context.Background())
	}()

	// This should return after shutdown
	select {
	case <-sdk.WaitForShutdown():
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for shutdown")
	}
}

func TestSDKAccessors(t *testing.T) {
	sdk, _ := New(nil)

	// Test all accessors return non-nil
	if sdk.Config() == nil {
		t.Error("Config() returned nil")
	}
	if sdk.Container() == nil {
		t.Error("Container() returned nil")
	}
	if sdk.PluginManager() == nil {
		t.Error("PluginManager() returned nil")
	}
	if sdk.EventBus() == nil {
		t.Error("EventBus() returned nil")
	}
	if sdk.Logger() == nil {
		t.Error("Logger() returned nil")
	}
	if sdk.Metrics() == nil {
		t.Error("Metrics() returned nil")
	}
}

func TestLogLevelOptions(t *testing.T) {
	tests := []struct {
		level    string
		expected int
	}{
		{"debug", 0},
		{"info", 1},
		{"warn", 2},
		{"error", 3},
		{"unknown", 1}, // defaults to info
	}

	for _, test := range tests {
		opts := &Options{
			LogLevel: test.level,
		}

		sdk, err := New(opts)
		if err != nil {
			t.Fatalf("Failed to create SDK with level %s: %v", test.level, err)
		}

		if sdk.Logger() == nil {
			t.Errorf("Logger should not be nil for level %s", test.level)
		}

		_ = sdk
	}
}