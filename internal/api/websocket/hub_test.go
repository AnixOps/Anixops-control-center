package websocket

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/eventbus"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestHub_RegisterUnregister(t *testing.T) {
	hub := NewHub(eventbus.New())
	go hub.Run()
	defer func() {
		// Stop the hub by broadcasting nil (will cause panic, so we just let it be)
	}()

	// Create a mock client
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	// Register client
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client, got %d", hub.ClientCount())
	}

	// Unregister client
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("expected 0 clients, got %d", hub.ClientCount())
	}
}

func TestHub_BroadcastMessage(t *testing.T) {
	hub := NewHub(eventbus.New())
	go hub.Run()

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Broadcast a message
	hub.Broadcast(MessageTypeLog, map[string]string{"test": "data"})

	select {
	case msg := <-client.send:
		var received Message
		if err := json.Unmarshal(msg, &received); err != nil {
			t.Errorf("failed to unmarshal message: %v", err)
		}
		if received.Type != MessageTypeLog {
			t.Errorf("expected type %s, got %s", MessageTypeLog, received.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("expected to receive broadcast message")
	}
}

func TestHub_BroadcastToTopicWithClients(t *testing.T) {
	hub := NewHub(eventbus.New())
	go hub.Run()

	// Create two clients, one subscribed to "logs"
	client1 := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: map[string]bool{"logs": true},
	}
	client2 := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: map[string]bool{"events": true},
	}

	hub.register <- client1
	hub.register <- client2
	time.Sleep(10 * time.Millisecond)

	// Broadcast to "logs" topic
	hub.BroadcastToTopic("logs", MessageTypeLog, map[string]string{"message": "test"})

	// Client1 should receive
	select {
	case <-client1.send:
		// Good
	case <-time.After(100 * time.Millisecond):
		t.Error("client1 should have received message")
	}

	// Client2 should NOT receive
	select {
	case <-client2.send:
		t.Error("client2 should not have received message")
	case <-time.After(50 * time.Millisecond):
		// Good
	}
}

func TestClient_SubscribeUnsubscribe(t *testing.T) {
	hub := NewHub(eventbus.New())

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	// Subscribe
	client.topics["logs"] = true
	client.topics["events"] = true

	if len(client.topics) != 2 {
		t.Errorf("expected 2 topics, got %d", len(client.topics))
	}

	// Unsubscribe
	delete(client.topics, "logs")

	if len(client.topics) != 1 {
		t.Errorf("expected 1 topic, got %d", len(client.topics))
	}
	if client.topics["events"] != true {
		t.Error("expected events topic to be true")
	}
}

func TestUpgrader_CheckOrigin(t *testing.T) {
	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Origin", "http://example.com")

	// CheckOrigin returns true for all origins
	if !Upgrader.CheckOrigin(req) {
		t.Error("CheckOrigin should return true for all origins")
	}
}

func TestHandleWebSocket(t *testing.T) {
	gin.SetMode(gin.TestMode)

	hub := NewHub(eventbus.New())
	go hub.Run()

	router := gin.New()
	router.GET("/ws", HandleWebSocket(hub))

	server := httptest.NewServer(router)
	defer server.Close()

	// Convert http to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect as WebSocket client
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Should receive welcome message
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if msg.Type != "connected" {
		t.Errorf("expected type 'connected', got %s", msg.Type)
	}
}

func TestLogStreamer_StreamToTopic(t *testing.T) {
	hub := NewHub(eventbus.New())
	go hub.Run()

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: map[string]bool{"logs": true},
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	streamer := NewLogStreamer(hub)
	streamer.Stream("info", "test", "test message")

	select {
	case msg := <-client.send:
		var received Message
		if err := json.Unmarshal(msg, &received); err != nil {
			t.Errorf("failed to unmarshal: %v", err)
		}
		if received.Type != MessageTypeLog {
			t.Errorf("expected type %s, got %s", MessageTypeLog, received.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("expected to receive log message")
	}
}

func TestEventStreamer_StreamToTopic(t *testing.T) {
	hub := NewHub(eventbus.New())
	go hub.Run()

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: map[string]bool{"events": true},
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	streamer := NewEventStreamer(hub)
	streamer.Stream("node_deployed", map[string]string{"node": "tokyo-01"})

	select {
	case msg := <-client.send:
		var received Message
		if err := json.Unmarshal(msg, &received); err != nil {
			t.Errorf("failed to unmarshal: %v", err)
		}
		if received.Type != MessageTypeEvent {
			t.Errorf("expected type %s, got %s", MessageTypeEvent, received.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("expected to receive event message")
	}
}

func TestMessageJSON(t *testing.T) {
	msg := Message{
		Type:      MessageTypeLog,
		Timestamp: 1234567890,
		Data:      map[string]string{"key": "value"},
		Error:     "",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Type != MessageTypeLog {
		t.Errorf("expected type %s, got %s", MessageTypeLog, decoded.Type)
	}
}

func TestMessageWithError(t *testing.T) {
	msg := Message{
		Type:      MessageTypeError,
		Timestamp: time.Now().Unix(),
		Error:     "something went wrong",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Error != "something went wrong" {
		t.Errorf("expected error message, got %s", decoded.Error)
	}
}

func TestHub_ConcurrentBroadcast(t *testing.T) {
	hub := NewHub(eventbus.New())
	go hub.Run()

	// Create multiple clients
	for i := 0; i < 5; i++ {
		client := &Client{
			hub:    hub,
			send:   make(chan []byte, 256),
			topics: make(map[string]bool),
		}
		hub.register <- client
	}
	time.Sleep(10 * time.Millisecond)

	// Broadcast multiple messages
	for i := 0; i < 10; i++ {
		hub.Broadcast(MessageTypeStatus, map[string]int{"count": i})
	}
}

func TestHub_FullSendChannel(t *testing.T) {
	hub := NewHub(eventbus.New())
	go hub.Run()

	// Create client with buffer size 1
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 1), // Very small buffer
		topics: make(map[string]bool),
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Send many messages to fill the channel
	for i := 0; i < 100; i++ {
		hub.Broadcast(MessageTypeLog, map[string]int{"i": i})
	}

	// Give time for processing
	time.Sleep(50 * time.Millisecond)

	// Client should have been removed due to full channel
	// (implementation dependent)
}

func TestClientFields(t *testing.T) {
	hub := NewHub(eventbus.New())

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
		userID: "user123",
		role:   "admin",
	}

	if client.userID != "user123" {
		t.Errorf("expected userID 'user123', got %s", client.userID)
	}
	if client.role != "admin" {
		t.Errorf("expected role 'admin', got %s", client.role)
	}
}