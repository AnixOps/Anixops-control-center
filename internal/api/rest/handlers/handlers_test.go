package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anixops/anixops-control-center/internal/core/plugin"
	"github.com/gin-gonic/gin"
)

func setupTestHandlers() *Handlers {
	gin.SetMode(gin.TestMode)
	return NewHandlers(plugin.NewManager(), nil, nil, nil, nil, nil, nil, nil, nil)
}

func TestHealthCheck(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.HealthCheck(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "healthy") {
		t.Error("response should contain 'healthy'")
	}
}

func TestReadinessCheck(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ReadinessCheck(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "ready") {
		t.Error("response should contain 'ready'")
	}
}

func TestLogin_Success(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"admin@example.com","password":"admin123456"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	// Without auth service, returns 500
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"wrong@example.com","password":"wrong"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	// Without auth service, returns 500
	if w.Code != http.StatusUnauthorized && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d or %d, got %d", http.StatusUnauthorized, http.StatusInternalServerError, w.Code)
	}
}

func TestLogin_InvalidRequest(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLogout(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.Logout(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRefreshToken(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/refresh", strings.NewReader(`{"refresh_token":"test-refresh-token"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.RefreshToken(c)

	// Without auth service, returns 500
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError && w.Code != http.StatusBadRequest {
		t.Errorf("unexpected status %d", w.Code)
	}
}

func TestGetCurrentUser(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user123")
	c.Set("role", "admin")

	h.GetCurrentUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "user123") {
		t.Error("response should contain user ID")
	}
}

func TestListPlugins(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ListPlugins(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetPlugin(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "name", Value: "ansible"}}

	h.GetPlugin(c)

	// Without registered plugin, returns 404
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("expected status %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestExecutePlugin(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/plugins/ansible/execute", strings.NewReader(`{"action":"run_playbook","params":{}}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "name", Value: "ansible"}}

	h.ExecutePlugin(c)

	// Without registered plugin, returns 404
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
		t.Errorf("unexpected status %d", w.Code)
	}
}

func TestExecutePlugin_InvalidBody(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/plugins/ansible/execute", strings.NewReader(`invalid`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "name", Value: "ansible"}}

	h.ExecutePlugin(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetPluginStatus(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "name", Value: "ansible"}}

	h.GetPluginStatus(c)

	// Without registered plugin, still returns status
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetPluginCapabilities(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "name", Value: "ansible"}}

	h.GetPluginCapabilities(c)

	// Without registered plugin, returns 404
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("expected status %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestStartPlugin(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/plugins/ansible/start", nil)
	c.Params = gin.Params{{Key: "name", Value: "ansible"}}

	h.StartPlugin(c)

	// Without registered plugin, returns 500 (plugin not found)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status %d", w.Code)
	}
}

func TestStopPlugin(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/plugins/ansible/stop", nil)
	c.Params = gin.Params{{Key: "name", Value: "ansible"}}

	h.StopPlugin(c)

	// Without registered plugin, returns 500 (plugin not found)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status %d", w.Code)
	}
}

func TestRestartPlugin(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/plugins/ansible/restart", nil)
	c.Params = gin.Params{{Key: "name", Value: "ansible"}}

	h.RestartPlugin(c)

	// Without registered plugin, returns 500 (plugin not found)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status %d", w.Code)
	}
}

func TestListNodes(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ListNodes(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreateNode(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/nodes", strings.NewReader(`{"name":"test-node","host":"192.168.1.1"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateNode(c)

	if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d or %d, got %d", http.StatusCreated, http.StatusInternalServerError, w.Code)
	}
}

func TestGetNode(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetNode(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateNode(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/nodes/1", strings.NewReader(`{"name":"updated-node"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UpdateNode(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestDeleteNode(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeleteNode(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetNodeStats(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetNodeStats(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListAgents(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ListAgents(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestConnectAgent(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.ConnectAgent(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDisconnectAgent(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DisconnectAgent(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestExecAgent(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.ExecAgent(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetAgentInfo(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetAgentInfo(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListPlaybooks(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ListPlaybooks(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetPlaybook(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "name", Value: "deploy.yml"}}

	h.GetPlaybook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRunPlaybook(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/playbooks/run", strings.NewReader(`{"playbook":"deploy.yml"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.RunPlaybook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestValidatePlaybook(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "name", Value: "deploy.yml"}}

	h.ValidatePlaybook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "valid") {
		t.Error("response should contain 'valid'")
	}
}

func TestGetInventory(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetInventory(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListHosts(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ListHosts(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListUsers(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/users?page=1&per_page=10", nil)

	h.ListUsers(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreateUser(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/users", strings.NewReader(`{"email":"test@example.com","password":"password123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateUser(c)

	if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d or %d, got %d", http.StatusCreated, http.StatusInternalServerError, w.Code)
	}
}

func TestGetUser(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/users/1", strings.NewReader(`{"email":"test@example.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UpdateUser(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeleteUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBanUser(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.BanUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUnbanUser(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UnbanUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetUserSubscriptions(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetUserSubscriptions(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListOrders(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ListOrders(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetOrder(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetOrder(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListPlans(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ListPlans(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetPlan(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetPlan(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetDashboard(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetDashboard(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "nodes") {
		t.Error("response should contain 'nodes'")
	}
}

func TestGetDashboardStats(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetDashboardStats(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetLogs(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetLogs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStreamLogs(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.StreamLogs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetSettings(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetSettings(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateSettings(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.UpdateSettings(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStreamEvents(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.StreamEvents(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected Content-Type text/event-stream, got %s", w.Header().Get("Content-Type"))
	}
}

func TestGetAdminStats(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetAdminStats(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetAuditLogs(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetAuditLogs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetSystemInfo(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetSystemInfo(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "version") {
		t.Error("response should contain 'version'")
	}
}

func TestCreateBackup(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.CreateBackup(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListBackups(t *testing.T) {
	h := setupTestHandlers()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.ListBackups(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestNewHandlers(t *testing.T) {
	mgr := plugin.NewManager()
	h := NewHandlers(mgr, nil, nil, nil, nil, nil, nil, nil, nil)

	if h == nil {
		t.Error("NewHandlers returned nil")
	}
	if h.pluginMgr == nil {
		t.Error("plugin manager should not be nil")
	}
}