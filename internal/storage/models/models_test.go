package models

import (
	"testing"
	"time"
)

func TestUserTableName(t *testing.T) {
	user := User{}
	if user.TableName() != "users" {
		t.Errorf("expected table name 'users', got %s", user.TableName())
	}
}

func TestAPITokenTableName(t *testing.T) {
	token := APIToken{}
	if token.TableName() != "api_tokens" {
		t.Errorf("expected table name 'api_tokens', got %s", token.TableName())
	}
}

func TestAuditLogTableName(t *testing.T) {
	log := AuditLog{}
	if log.TableName() != "audit_logs" {
		t.Errorf("expected table name 'audit_logs', got %s", log.TableName())
	}
}

func TestUserFields(t *testing.T) {
	now := time.Now()
	user := User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         "admin",
		AuthProvider: "local",
		Enabled:      true,
		LastLoginAt:  &now,
	}

	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}
	if user.Role != "admin" {
		t.Errorf("expected role 'admin', got %s", user.Role)
	}
	if !user.Enabled {
		t.Error("expected enabled to be true")
	}
}

func TestAPITokenFields(t *testing.T) {
	now := time.Now()
	token := APIToken{
		ID:        1,
		UserID:    1,
		Name:      "test-token",
		Token:     "abc123",
		ExpiresAt: &now,
	}

	if token.ID != 1 {
		t.Errorf("expected ID 1, got %d", token.ID)
	}
	if token.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", token.UserID)
	}
	if token.Name != "test-token" {
		t.Errorf("expected name 'test-token', got %s", token.Name)
	}
}

func TestAuditLogFields(t *testing.T) {
	log := AuditLog{
		ID:        1,
		UserID:    1,
		Action:    "login",
		Resource:  "auth",
		IP:        "192.168.1.1",
		UserAgent: "Mozilla/5.0",
		Details:   "User logged in",
	}

	if log.ID != 1 {
		t.Errorf("expected ID 1, got %d", log.ID)
	}
	if log.Action != "login" {
		t.Errorf("expected action 'login', got %s", log.Action)
	}
	if log.IP != "192.168.1.1" {
		t.Errorf("expected IP '192.168.1.1', got %s", log.IP)
	}
}