package config_test

import (
	"testing"

	"github.com/anixops/anixops-control-center/internal/core/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	// Check defaults are set
	if cfg.Env != "development" {
		t.Errorf("expected env 'development', got '%s'", cfg.Env)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected server host '0.0.0.0', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected server port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Database.Driver != "sqlite" {
		t.Errorf("expected database driver 'sqlite', got '%s'", cfg.Database.Driver)
	}

	if cfg.Auth.Provider != "local" {
		t.Errorf("expected auth provider 'local', got '%s'", cfg.Auth.Provider)
	}
}

func TestLoadFromString(t *testing.T) {
	yaml := `
env: "production"
server:
  host: "127.0.0.1"
  port: 9000
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
`

	cfg, err := config.LoadFromString(yaml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Env != "production" {
		t.Errorf("expected env 'production', got '%s'", cfg.Env)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected server host '127.0.0.1', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != 9000 {
		t.Errorf("expected server port 9000, got %d", cfg.Server.Port)
	}

	if cfg.Database.Driver != "postgres" {
		t.Errorf("expected database driver 'postgres', got '%s'", cfg.Database.Driver)
	}
}

func TestSetDefaults(t *testing.T) {
	// Call internal defaults function through LoadFromString with empty config

	emptyYaml := ``
	loaded, err := config.LoadFromString(emptyYaml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that defaults are applied
	if loaded.Server.Port == 0 {
		t.Error("expected default server port to be set")
	}

	if loaded.Auth.JWT.Expire == 0 {
		t.Error("expected default JWT expire to be set")
	}
}

func TestJWTConfigDefaults(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Auth.JWT.Expire != 86400 {
		t.Errorf("expected JWT expire 86400, got %d", cfg.Auth.JWT.Expire)
	}

	if cfg.Auth.JWT.RefreshExpire != 604800 {
		t.Errorf("expected JWT refresh expire 604800, got %d", cfg.Auth.JWT.RefreshExpire)
	}
}

func TestGRPCConfigDefaults(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.GRPC.Port != 50052 {
		t.Errorf("expected gRPC port 50052, got %d", cfg.GRPC.Port)
	}

	if cfg.GRPC.KeepaliveTime != 30 {
		t.Errorf("expected keepalive time 30, got %d", cfg.GRPC.KeepaliveTime)
	}
}

func TestLoggingConfigDefaults(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Logging.Level != "info" {
		t.Errorf("expected logging level 'info', got '%s'", cfg.Logging.Level)
	}

	if cfg.Logging.Format != "text" {
		t.Errorf("expected logging format 'text', got '%s'", cfg.Logging.Format)
	}
}
