// Package mobile provides mobile-friendly APIs for AnixOps Control Center
// This package is designed to be compiled with gomobile for Android/iOS support
//
// Build requirements:
//   - Android: gomobile bind -target=android -o=dist/mobile/anixops-sdk.aar ./internal/mobile
//   - iOS: gomobile bind -target=ios -o=dist/mobile/AnixOpsSDK.xcframework ./internal/mobile
//
// See: https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile

package mobile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// SDK version information
const (
	SDKVersion = "1.5.0"
	SDKName    = "AnixOpsSDK"
)

// GetSDKVersion returns the SDK version
func GetSDKVersion() string {
	return SDKVersion
}

// Errors
var (
	ErrNotConnected     = errors.New("not connected to server")
	ErrAuthentication   = errors.New("authentication failed")
	ErrInvalidResponse  = errors.New("invalid response from server")
	ErrRequestFailed    = errors.New("request failed")
	ErrTokenExpired     = errors.New("token expired")
)

// Config holds client configuration
type Config struct {
	BaseURL     string
	Timeout     time.Duration
	AccessToken string
}

// MobileClient provides a simplified API for mobile applications
type MobileClient struct {
	config     Config
	httpClient *http.Client
	wsClient   *WebSocketClient
	mu         sync.RWMutex
}

// NewMobileClient creates a new mobile client
func NewMobileClient(baseURL string) *MobileClient {
	return &MobileClient{
		config: Config{
			BaseURL: baseURL,
			Timeout: 30 * time.Second,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetToken sets the authentication token
func (c *MobileClient) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.AccessToken = token
}

// SetTimeout sets the request timeout
func (c *MobileClient) SetTimeout(timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.Timeout = timeout
	c.httpClient.Timeout = timeout
}

// getToken returns the current access token
func (c *MobileClient) getToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.AccessToken
}

// Login authenticates and stores the token
func (c *MobileClient) Login(email, password string) (*LoginResponse, error) {
	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.doRequest("POST", "/api/v1/auth/login", req, false)
	if err != nil {
		return nil, err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(resp, &loginResp); err != nil {
		return nil, ErrInvalidResponse
	}

	c.SetToken(loginResp.AccessToken)
	return &loginResp, nil
}

// Register creates a new user account
func (c *MobileClient) Register(email, password, role string) (*User, error) {
	req := RegisterRequest{
		Email:    email,
		Password: password,
		Role:     role,
	}

	resp, err := c.doRequest("POST", "/api/v1/auth/register", req, false)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, ErrInvalidResponse
	}

	return &user, nil
}

// Logout logs out the current user
func (c *MobileClient) Logout() error {
	_, err := c.doRequest("POST", "/api/v1/auth/logout", nil, true)
	if err != nil {
		return err
	}
	c.SetToken("")
	return nil
}

// RefreshToken refreshes the access token
func (c *MobileClient) RefreshToken(refreshToken string) (*TokenResponse, error) {
	req := map[string]string{"refresh_token": refreshToken}
	resp, err := c.doRequest("POST", "/api/v1/auth/refresh", req, false)
	if err != nil {
		return nil, err
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(resp, &tokenResp); err != nil {
		return nil, ErrInvalidResponse
	}

	c.SetToken(tokenResp.AccessToken)
	return &tokenResp, nil
}

// GetCurrentUser gets the current authenticated user
func (c *MobileClient) GetCurrentUser() (*User, error) {
	resp, err := c.doRequest("GET", "/api/v1/auth/me", nil, true)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, ErrInvalidResponse
	}

	return &user, nil
}

// GetDashboard retrieves dashboard data
func (c *MobileClient) GetDashboard() (*Dashboard, error) {
	resp, err := c.doRequest("GET", "/api/v1/dashboard", nil, true)
	if err != nil {
		return nil, err
	}

	var dashboard Dashboard
	if err := json.Unmarshal(resp, &dashboard); err != nil {
		return nil, ErrInvalidResponse
	}

	return &dashboard, nil
}

// GetNodes retrieves the list of nodes
func (c *MobileClient) GetNodes() ([]Node, error) {
	resp, err := c.doRequest("GET", "/api/v1/nodes", nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Nodes []Node `json:"nodes"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, ErrInvalidResponse
	}

	return result.Nodes, nil
}

// GetNode retrieves a single node by ID
func (c *MobileClient) GetNode(id string) (*Node, error) {
	resp, err := c.doRequest("GET", "/api/v1/nodes/"+id, nil, true)
	if err != nil {
		return nil, err
	}

	var node Node
	if err := json.Unmarshal(resp, &node); err != nil {
		return nil, ErrInvalidResponse
	}

	return &node, nil
}

// CreateNode creates a new node
func (c *MobileClient) CreateNode(node *Node) (*Node, error) {
	resp, err := c.doRequest("POST", "/api/v1/nodes", node, true)
	if err != nil {
		return nil, err
	}

	var created Node
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, ErrInvalidResponse
	}

	return &created, nil
}

// UpdateNode updates an existing node
func (c *MobileClient) UpdateNode(id string, updates map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/v1/nodes/"+id, updates, true)
	return err
}

// DeleteNode deletes a node
func (c *MobileClient) DeleteNode(id string) error {
	_, err := c.doRequest("DELETE", "/api/v1/nodes/"+id, nil, true)
	return err
}

// GetUsers retrieves the list of users
func (c *MobileClient) GetUsers(page, perPage int) (*UserListResponse, error) {
	path := fmt.Sprintf("/api/v1/users?page=%d&per_page=%d", page, perPage)
	resp, err := c.doRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var result UserListResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, ErrInvalidResponse
	}

	return &result, nil
}

// GetUser retrieves a single user by ID
func (c *MobileClient) GetUser(id string) (*User, error) {
	resp, err := c.doRequest("GET", "/api/v1/users/"+id, nil, true)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, ErrInvalidResponse
	}

	return &user, nil
}

// RunPlaybook executes an Ansible playbook
func (c *MobileClient) RunPlaybook(playbook string, params map[string]interface{}) (*PlaybookResult, error) {
	req := map[string]interface{}{
		"playbook": playbook,
	}
	for k, v := range params {
		req[k] = v
	}

	resp, err := c.doRequest("POST", "/api/v1/playbooks/run", req, true)
	if err != nil {
		return nil, err
	}

	var result PlaybookResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, ErrInvalidResponse
	}

	return &result, nil
}

// ListPlaybooks lists available playbooks
func (c *MobileClient) ListPlaybooks() ([]string, error) {
	resp, err := c.doRequest("GET", "/api/v1/playbooks", nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Playbooks []string `json:"playbooks"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, ErrInvalidResponse
	}

	return result.Playbooks, nil
}

// GetLogs retrieves recent logs
func (c *MobileClient) GetLogs(page, perPage int) (*LogListResponse, error) {
	path := fmt.Sprintf("/api/v1/logs?page=%d&per_page=%d", page, perPage)
	resp, err := c.doRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var result LogListResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, ErrInvalidResponse
	}

	return &result, nil
}

// GetPlans retrieves available subscription plans
func (c *MobileClient) GetPlans() ([]Plan, error) {
	resp, err := c.doRequest("GET", "/api/v1/plans", nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Plans []Plan `json:"plans"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, ErrInvalidResponse
	}

	return result.Plans, nil
}

// GetSubscriptions retrieves the current user's subscriptions
func (c *MobileClient) GetSubscriptions() ([]Subscription, error) {
	resp, err := c.doRequest("GET", "/api/v1/auth/me", nil, true)
	if err != nil {
		return nil, err
	}

	// Parse user ID from response
	var user struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, ErrInvalidResponse
	}

	path := fmt.Sprintf("/api/v1/users/%s/subscriptions", user.ID)
	resp, err = c.doRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Subscriptions []Subscription `json:"subscriptions"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, ErrInvalidResponse
	}

	return result.Subscriptions, nil
}

// doRequest performs an HTTP request
func (c *MobileClient) doRequest(method, path string, body interface{}, auth bool) ([]byte, error) {
	c.mu.RLock()
	baseURL := c.config.BaseURL
	c.mu.RUnlock()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	url := baseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if auth {
		token := c.getToken()
		if token == "" {
			return nil, ErrNotConnected
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("%s: %s", ErrRequestFailed, errResp.Error)
		}
		return nil, fmt.Errorf("%w: status %d", ErrRequestFailed, resp.StatusCode)
	}

	return respBody, nil
}

// WebSocket methods

// ConnectWebSocket establishes a WebSocket connection
func (c *MobileClient) ConnectWebSocket(ctx context.Context) error {
	c.mu.RLock()
	baseURL := c.config.BaseURL
	token := c.config.AccessToken
	c.mu.RUnlock()

	// Convert HTTP URL to WebSocket URL
	wsURL := baseURL
	if len(wsURL) > 4 && wsURL[:4] == "http" {
		wsURL = "ws" + wsURL[4:]
	}
	wsURL = wsURL + "/api/v1/events?token=" + token

	c.wsClient = NewWebSocketClient(wsURL)
	return c.wsClient.Connect(ctx)
}

// DisconnectWebSocket closes the WebSocket connection
func (c *MobileClient) DisconnectWebSocket() error {
	if c.wsClient == nil {
		return nil
	}
	return c.wsClient.Close()
}

// SubscribeToEvents subscribes to real-time events
func (c *MobileClient) SubscribeToEvents(handler func(event string, data []byte)) error {
	if c.wsClient == nil {
		return ErrNotConnected
	}
	c.wsClient.SetEventHandler(handler)
	return nil
}

// Request and Response Types

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

// TokenResponse represents a token refresh response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// User represents a user
type User struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	Enabled   bool    `json:"enabled"`
	CreatedAt string  `json:"created_at"`
}

// UserListResponse represents a list of users
type UserListResponse struct {
	Users    []User `json:"users"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	Total    int64  `json:"total"`
}

// Node represents a proxy node
type Node struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Host         string  `json:"host"`
	Port         int     `json:"port"`
	Type         string  `json:"type"`
	Status       string  `json:"status"`
	Region       string  `json:"region"`
	TrafficUp    int64   `json:"traffic_up"`
	TrafficDown  int64   `json:"traffic_down"`
	UserCount    int     `json:"user_count"`
	Enabled      bool    `json:"enabled"`
}

// Dashboard represents dashboard statistics
type Dashboard struct {
	Nodes       int    `json:"nodes"`
	OnlineNodes int    `json:"online_nodes"`
	Users       int    `json:"users"`
	ActiveSubs  int    `json:"active_subs"`
	Traffic     string `json:"traffic"`
}

// Plan represents a subscription plan
type Plan struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	Duration     int     `json:"duration"`
	TrafficLimit int64   `json:"traffic_limit"`
	Enabled      bool    `json:"enabled"`
}

// Subscription represents a user subscription
type Subscription struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	PlanID       string `json:"plan_id"`
	Status       string `json:"status"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	TrafficUsed  int64  `json:"traffic_used"`
	TrafficLimit int64  `json:"traffic_limit"`
}

// PlaybookResult represents a playbook execution result
type PlaybookResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// LogEntry represents a log entry
type LogEntry struct {
	ID        string `json:"id"`
	Time      string `json:"created_at"`
	Level     string `json:"level"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Message   string `json:"details"`
}

// LogListResponse represents a list of logs
type LogListResponse struct {
	Logs    []LogEntry `json:"logs"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
	Total   int64      `json:"total"`
}

// WebSocketClient provides real-time updates for mobile
type WebSocketClient struct {
	url     string
	conn    *websocket.Conn
	onEvent func(event string, data []byte)
	mu      sync.RWMutex
}

// NewWebSocketClient creates a WebSocket client
func NewWebSocketClient(url string) *WebSocketClient {
	return &WebSocketClient{url: url}
}

// SetEventHandler sets the event handler
func (c *WebSocketClient) SetEventHandler(handler func(event string, data []byte)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onEvent = handler
}

// Connect establishes WebSocket connection
func (c *WebSocketClient) Connect(ctx context.Context) error {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	// Start reading messages
	go c.readMessages()

	return nil
}

// readMessages reads messages from the WebSocket connection
func (c *WebSocketClient) readMessages() {
	for {
		c.mu.RLock()
		conn := c.conn
		onEvent := c.onEvent
		c.mu.RUnlock()

		if conn == nil {
			return
		}

		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Connection closed or error
			return
		}

		if messageType == websocket.TextMessage && onEvent != nil {
			// Parse SSE event format
			// Format: event: <event_name>\ndata: <json_data>\n\n
			var event string
			var data []byte
			lines := bytes.Split(message, []byte("\n"))
			for _, line := range lines {
				if bytes.HasPrefix(line, []byte("event: ")) {
					event = string(bytes.TrimPrefix(line, []byte("event: ")))
				} else if bytes.HasPrefix(line, []byte("data: ")) {
					data = bytes.TrimPrefix(line, []byte("data: "))
				}
			}

			if event != "" && data != nil {
				onEvent(event, data)
			}
		}
	}
}

// Close closes the WebSocket connection
func (c *WebSocketClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}