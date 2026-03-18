//go:build cgo
// +build cgo

package sqlite

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// Use a temporary database file
	tmpFile := "/tmp/test_anixops.db"
	defer os.Remove(tmpFile)

	cfg := &Config{
		Driver:   "sqlite",
		Database: tmpFile,
	}

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if DB == nil {
		t.Fatal("DB is nil after Init")
	}
}

func TestGetDB(t *testing.T) {
	// Initialize first
	tmpFile := "/tmp/test_anixops_get.db"
	defer os.Remove(tmpFile)

	cfg := &Config{
		Driver:   "sqlite",
		Database: tmpFile,
	}
	_ = Init(cfg)

	db := GetDB()
	if db == nil {
		t.Error("GetDB returned nil")
	}
}

func TestClose(t *testing.T) {
	// Initialize first
	tmpFile := "/tmp/test_anixops_close.db"
	defer os.Remove(tmpFile)

	cfg := &Config{
		Driver:   "sqlite",
		Database: tmpFile,
	}
	_ = Init(cfg)

	err := Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Reset for other tests
	DB = nil
}

func TestInit_DefaultDriver(t *testing.T) {
	tmpFile := "/tmp/test_anixops_default.db"
	defer os.Remove(tmpFile)

	cfg := &Config{
		Driver:   "unknown",
		Database: tmpFile,
	}

	// Should default to sqlite
	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init with unknown driver failed: %v", err)
	}

	_ = Close()
}

func TestInit_Postgres(t *testing.T) {
	cfg := &Config{
		Driver:   "postgres",
		Database: "test",
	}

	err := Init(cfg)
	if err == nil {
		t.Error("expected error for postgres (not supported)")
	}
}

func TestConfig(t *testing.T) {
	cfg := &Config{
		Driver:   "sqlite",
		Database: "test.db",
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test123",
	}

	if cfg.Driver != "sqlite" {
		t.Errorf("expected driver 'sqlite', got %s", cfg.Driver)
	}
	if cfg.Database != "test.db" {
		t.Errorf("expected database 'test.db', got %s", cfg.Database)
	}
}