package config_test

import (
	"os"
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

func TestLoadFromString_InvalidYAML(t *testing.T) {
	yaml := `invalid: yaml: [unclosed`

	_, err := config.LoadFromString(yaml)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	// Create a temporary config file
	content := `
env: "test"
server:
  host: "0.0.0.0"
  port: 8080
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := config.Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Env != "test" {
		t.Errorf("expected env 'test', got '%s'", cfg.Env)
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

	if cfg.GRPC.KeepaliveTimeout != 10 {
		t.Errorf("expected keepalive timeout 10, got %d", cfg.GRPC.KeepaliveTimeout)
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

	if cfg.Logging.Output != "stdout" {
		t.Errorf("expected logging output 'stdout', got '%s'", cfg.Logging.Output)
	}
}

func TestDatabaseConfigDefaults(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Database.Driver != "sqlite" {
		t.Errorf("expected database driver 'sqlite', got '%s'", cfg.Database.Driver)
	}

	if cfg.Database.Database != "data/anixops.db" {
		t.Errorf("expected database 'data/anixops.db', got '%s'", cfg.Database.Database)
	}
}

func TestServerConfigDefaults(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Server.Mode != "debug" {
		t.Errorf("expected server mode 'debug', got '%s'", cfg.Server.Mode)
	}
}

func TestConfigSave(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Env = "test-save"

	// Create temp file
	tmpFile, err := os.CreateTemp("", "config-save-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Save config
	if err := cfg.Save(tmpFile.Name()); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Load and verify
	loaded, err := config.Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if loaded.Env != "test-save" {
		t.Errorf("expected env 'test-save', got '%s'", loaded.Env)
	}
}

func TestConfigSave_InvalidPath(t *testing.T) {
	cfg := config.DefaultConfig()

	// Try to save to an invalid path
	err := cfg.Save("/nonexistent/directory/config.yaml")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestOAuthConfig(t *testing.T) {
	yaml := `
auth:
  provider: "oauth"
  oauth:
    providers:
      - name: "google"
        client_id: "test-client-id"
        client_secret: "test-secret"
        redirect_url: "http://localhost/callback"
        scopes:
          - "email"
          - "profile"
`

	cfg, err := config.LoadFromString(yaml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Auth.Provider != "oauth" {
		t.Errorf("expected auth provider 'oauth', got '%s'", cfg.Auth.Provider)
	}

	if cfg.Auth.OAuth == nil {
		t.Fatal("expected OAuth config to be set")
	}

	if len(cfg.Auth.OAuth.Providers) != 1 {
		t.Fatalf("expected 1 OAuth provider, got %d", len(cfg.Auth.OAuth.Providers))
	}

	if cfg.Auth.OAuth.Providers[0].Name != "google" {
		t.Errorf("expected provider name 'google', got '%s'", cfg.Auth.OAuth.Providers[0].Name)
	}
}

func TestLDAPConfig(t *testing.T) {
	yaml := `
auth:
  provider: "ldap"
  ldap:
    url: "ldap://localhost:389"
    bind_dn: "cn=admin,dc=example,dc=com"
    bind_password: "secret"
    base_dn: "ou=users,dc=example,dc=com"
    filter: "(uid=%s)"
    tls: true
`

	cfg, err := config.LoadFromString(yaml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Auth.Provider != "ldap" {
		t.Errorf("expected auth provider 'ldap', got '%s'", cfg.Auth.Provider)
	}

	if cfg.Auth.LDAP == nil {
		t.Fatal("expected LDAP config to be set")
	}

	if cfg.Auth.LDAP.URL != "ldap://localhost:389" {
		t.Errorf("expected LDAP URL 'ldap://localhost:389', got '%s'", cfg.Auth.LDAP.URL)
	}
}

func TestSAMLConfig(t *testing.T) {
	yaml := `
auth:
  provider: "saml"
  saml:
    entity_id: "anixops-sp"
    idp_metadata_url: "https://idp.example.com/metadata"
`

	cfg, err := config.LoadFromString(yaml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Auth.Provider != "saml" {
		t.Errorf("expected auth provider 'saml', got '%s'", cfg.Auth.Provider)
	}

	if cfg.Auth.SAML == nil {
		t.Fatal("expected SAML config to be set")
	}

	if cfg.Auth.SAML.EntityID != "anixops-sp" {
		t.Errorf("expected entity_id 'anixops-sp', got '%s'", cfg.Auth.SAML.EntityID)
	}
}

func TestPluginsConfig(t *testing.T) {
	yaml := `
plugins:
  ansible:
    playbook_dir: "/opt/playbooks"
  v2board:
    host: "http://localhost:8080"
`

	cfg, err := config.LoadFromString(yaml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Plugins.Ansible == nil {
		t.Error("expected ansible plugin config")
	}

	if cfg.Plugins.V2board == nil {
		t.Error("expected v2board plugin config")
	}
}

func TestGRPCConfig(t *testing.T) {
	yaml := `
grpc:
  enable: true
  host: "0.0.0.0"
  port: 50051
  keepalive_time: 60
  keepalive_timeout: 20
`

	cfg, err := config.LoadFromString(yaml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !cfg.GRPC.Enable {
		t.Error("expected gRPC to be enabled")
	}

	if cfg.GRPC.Port != 50051 {
		t.Errorf("expected gRPC port 50051, got %d", cfg.GRPC.Port)
	}

	if cfg.GRPC.KeepaliveTime != 60 {
		t.Errorf("expected keepalive time 60, got %d", cfg.GRPC.KeepaliveTime)
	}
}

func TestLocalAuthConfig(t *testing.T) {
	yaml := `
auth:
  provider: "local"
  local:
    admin_email: "admin@example.com"
    admin_password: "admin123"
`

	cfg, err := config.LoadFromString(yaml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Auth.Local == nil {
		t.Fatal("expected local auth config to be set")
	}

	if cfg.Auth.Local.AdminEmail != "admin@example.com" {
		t.Errorf("expected admin email 'admin@example.com', got '%s'", cfg.Auth.Local.AdminEmail)
	}
}