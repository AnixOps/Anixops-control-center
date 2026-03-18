package mobile

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNewMobileClient(t *testing.T) {
	client := NewMobileClient("http://localhost:8080")
	if client == nil {
		t.Fatal("NewMobileClient returned nil")
	}
}

func TestSetToken(t *testing.T) {
	client := NewMobileClient("http://localhost:8080")
	client.SetToken("test-token")

	token := client.getToken()
	if token != "test-token" {
		t.Errorf("expected token 'test-token', got %s", token)
	}
}

func TestSetTimeout(t *testing.T) {
	client := NewMobileClient("http://localhost:8080")
	client.SetTimeout(60)

	if client.httpClient.Timeout != 60 {
		t.Errorf("expected timeout 60, got %v", client.httpClient.Timeout)
	}
}

func TestLoginResponseJSON(t *testing.T) {
	loginResp := LoginResponse{
		AccessToken:  "test-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		User: User{
			ID:    "1",
			Email: "test@example.com",
			Role:  "admin",
		},
	}

	data, err := json.Marshal(loginResp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded LoginResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.AccessToken != "test-token" {
		t.Errorf("expected access token 'test-token', got %s", decoded.AccessToken)
	}
}

func TestNodeJSON(t *testing.T) {
	node := Node{
		ID:          "1",
		Name:        "test-node",
		Host:        "192.168.1.1",
		Port:        443,
		Type:        "v2ray",
		Status:      "online",
		Region:      "tokyo",
		TrafficUp:   1024,
		TrafficDown: 2048,
		UserCount:   10,
		Enabled:     true,
	}

	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Node
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Name != "test-node" {
		t.Errorf("expected name 'test-node', got %s", decoded.Name)
	}
	if decoded.Status != "online" {
		t.Errorf("expected status 'online', got %s", decoded.Status)
	}
}

func TestDashboardJSON(t *testing.T) {
	dashboard := Dashboard{
		Nodes:       8,
		OnlineNodes: 6,
		Users:       357,
		ActiveSubs:  120,
		Traffic:     "1.2TB",
	}

	data, err := json.Marshal(dashboard)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Dashboard
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Nodes != 8 {
		t.Errorf("expected 8 nodes, got %d", decoded.Nodes)
	}
}

func TestPlanJSON(t *testing.T) {
	plan := Plan{
		ID:           "1",
		Name:         "Basic",
		Description:  "Basic plan",
		Price:        9.99,
		Currency:     "USD",
		Duration:     30,
		TrafficLimit: 100 * 1024 * 1024 * 1024, // 100GB
		Enabled:      true,
	}

	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Plan
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Name != "Basic" {
		t.Errorf("expected name 'Basic', got %s", decoded.Name)
	}
}

func TestPlaybookResultJSON(t *testing.T) {
	result := PlaybookResult{
		Success: true,
		Output:  "Playbook completed successfully",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded PlaybookResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !decoded.Success {
		t.Error("expected success to be true")
	}
}

func TestLogEntryJSON(t *testing.T) {
	log := LogEntry{
		ID:       "1",
		Time:     "12:34:56",
		Level:    "ERROR",
		Action:   "login",
		Resource: "user",
		Message:  "Something went wrong",
	}

	data, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded LogEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Level != "ERROR" {
		t.Errorf("expected level 'ERROR', got %s", decoded.Level)
	}
}

func TestNewWebSocketClient(t *testing.T) {
	client := NewWebSocketClient("ws://localhost:8080/ws")
	if client == nil {
		t.Fatal("NewWebSocketClient returned nil")
	}
}

func TestSetEventHandler(t *testing.T) {
	client := NewWebSocketClient("ws://localhost:8080/ws")

	called := false
	client.SetEventHandler(func(event string, data []byte) {
		called = true
	})

	if client.onEvent == nil {
		t.Error("event handler should not be nil")
	}

	// Call the handler
	client.onEvent("test", []byte{})
	if !called {
		t.Error("event handler was not called")
	}
}

func TestWebSocketClose(t *testing.T) {
	client := NewWebSocketClient("ws://localhost:8080/ws")

	err := client.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestMobileClientWithToken(t *testing.T) {
	client := NewMobileClient("http://localhost:8080")
	client.SetToken("my-auth-token")

	token := client.getToken()
	if token != "my-auth-token" {
		t.Errorf("expected token 'my-auth-token', got %s", token)
	}
}

func TestWebSocketConnectNil(t *testing.T) {
	client := NewWebSocketClient("ws://localhost:8080/ws")

	ctx := context.Background()
	// This will fail because there's no server, but we test it doesn't panic
	_ = client.Connect(ctx)
}

func TestUserJSON(t *testing.T) {
	user := User{
		ID:        "1",
		Email:     "test@example.com",
		Role:      "admin",
		Enabled:   true,
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", decoded.Email)
	}
}

func TestSubscriptionJSON(t *testing.T) {
	sub := Subscription{
		ID:           "1",
		UserID:       "1",
		PlanID:       "1",
		Status:       "active",
		StartDate:    "2024-01-01",
		EndDate:      "2024-02-01",
		TrafficUsed:  1024 * 1024 * 1024,
		TrafficLimit: 10 * 1024 * 1024 * 1024,
	}

	data, err := json.Marshal(sub)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Subscription
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Status != "active" {
		t.Errorf("expected status 'active', got %s", decoded.Status)
	}
}

func TestErrors(t *testing.T) {
	errors := []error{
		ErrNotConnected,
		ErrAuthentication,
		ErrInvalidResponse,
		ErrRequestFailed,
		ErrTokenExpired,
	}

	for _, err := range errors {
		if err == nil {
			t.Error("error should not be nil")
		}
	}
}