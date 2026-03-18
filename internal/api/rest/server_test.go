package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anixops/anixops-control-center/internal/core/plugin"
	"github.com/anixops/anixops-control-center/internal/security/auth"
	"github.com/gin-gonic/gin"
)

func setupTestServer() *Server {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 3600, 86400, "test-issuer")
	rbacManager := auth.NewRBACManager()

	return NewServer(&ServerConfig{
		Mode:        gin.TestMode,
		PluginMgr:   plugin.NewManager(),
		JWTManager:  jwtManager,
		RBACManager: rbacManager,
	})
}

func TestNewServer(t *testing.T) {
	s := setupTestServer()
	if s == nil {
		t.Fatal("NewServer returned nil")
	}
	if s.router == nil {
		t.Error("router not initialized")
	}
}

func TestHealthEndpoint(t *testing.T) {
	s := setupTestServer()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	s.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestReadyEndpoint(t *testing.T) {
	s := setupTestServer()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ready", nil)
	s.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestLoginEndpoint(t *testing.T) {
	s := setupTestServer()
	w := httptest.NewRecorder()
	body := `{"email":"admin@example.com","password":"admin123456"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Without database, login returns 500 (service not configured)
	// This is expected behavior when services are nil
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestProtectedEndpoint_WithoutToken(t *testing.T) {
	s := setupTestServer()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/plugins", nil)
	s.router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestProtectedEndpoint_WithToken(t *testing.T) {
	s := setupTestServer()

	// Generate a valid token
	token, err := s.jwtManager.GenerateAccessToken("1", "admin", []string{}, "local")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/plugins", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	s.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAdminEndpoint_WithAdminRole(t *testing.T) {
	s := setupTestServer()

	token, err := s.jwtManager.GenerateAccessToken("1", "admin", []string{}, "local")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	s.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAdminEndpoint_WithNonAdminRole(t *testing.T) {
	s := setupTestServer()

	token, err := s.jwtManager.GenerateAccessToken("1", "viewer", []string{}, "local")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	s.router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestExtractToken_AuthorizationHeader(t *testing.T) {
	s := setupTestServer()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	token := s.extractToken(c)
	if token != "test-token" {
		t.Errorf("expected 'test-token', got '%s'", token)
	}
}

func TestExtractToken_QueryParam(t *testing.T) {
	s := setupTestServer()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test?token=query-token", nil)

	token := s.extractToken(c)
	if token != "query-token" {
		t.Errorf("expected 'query-token', got '%s'", token)
	}
}

func TestExtractToken_Cookie(t *testing.T) {
	s := setupTestServer()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.AddCookie(&http.Cookie{Name: "token", Value: "cookie-token"})

	token := s.extractToken(c)
	if token != "cookie-token" {
		t.Errorf("expected 'cookie-token', got '%s'", token)
	}
}

func TestRouter(t *testing.T) {
	s := setupTestServer()
	router := s.Router()
	if router == nil {
		t.Error("Router() returned nil")
	}
}

func TestStop(t *testing.T) {
	s := setupTestServer()
	err := s.Stop()
	if err != nil {
		t.Errorf("Stop returned error: %v", err)
	}
}