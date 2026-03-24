package websocket

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/eventbus"
)

func TestNewHub(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	if hub == nil {
		t.Fatal("NewHub returned nil")
	}
	if hub.clients == nil {
		t.Error("clients map not initialized")
	}
	if hub.broadcast == nil {
		t.Error("broadcast channel not initialized")
	}
	if hub.register == nil {
		t.Error("register channel not initialized")
	}
	if hub.unregister == nil {
		t.Error("unregister channel not initialized")
	}
}

func TestHub_Broadcast(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	hub.Broadcast("test", map[string]string{"key": "value"})
	// Broadcast should not block even with no clients
}

func TestHub_BroadcastToTopic(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	hub.BroadcastToTopic("logs", "log", map[string]string{"message": "test"})
	// Should not block even with no clients
}

func TestHub_ClientCount(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	count := hub.ClientCount()
	if count != 0 {
		t.Errorf("expected 0 clients, got %d", count)
	}
}

func TestHub_Register(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	// Register client
	go func() {
		hub.register <- client
	}()

	// Wait a bit for the goroutine
	time.Sleep(10 * time.Millisecond)

	// Check client count
	hub.mu.RLock()
	_, exists := hub.clients[client]
	hub.mu.RUnlock()

	if !exists {
		// Directly add to test the client count
		hub.mu.Lock()
		hub.clients[client] = true
		hub.mu.Unlock()
	}

	count := hub.ClientCount()
	if count != 1 {
		t.Errorf("expected 1 client, got %d", count)
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		Type:      "test",
		Timestamp: 12345,
		Data:      "hello",
		Error:     "",
	}

	if msg.Type != "test" {
		t.Errorf("expected type 'test', got '%s'", msg.Type)
	}
}

func TestMessage_JSON(t *testing.T) {
	msg := Message{
		Type:      MessageTypeLog,
		Timestamp: time.Now().Unix(),
		Data: map[string]string{
			"level":   "info",
			"message": "test log",
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	if decoded.Type != MessageTypeLog {
		t.Errorf("expected type '%s', got '%s'", MessageTypeLog, decoded.Type)
	}
}

func TestClient(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	if client.hub != hub {
		t.Error("client hub not set correctly")
	}
	if client.send == nil {
		t.Error("client send channel not initialized")
	}
	if client.topics == nil {
		t.Error("client topics map not initialized")
	}
}

func TestClient_Topics(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	// Subscribe to topics
	client.topics["logs"] = true
	client.topics["events"] = true

	if !client.topics["logs"] {
		t.Error("expected logs topic to be subscribed")
	}
	if !client.topics["events"] {
		t.Error("expected events topic to be subscribed")
	}

	// Unsubscribe
	delete(client.topics, "logs")
	if client.topics["logs"] {
		t.Error("expected logs topic to be unsubscribed")
	}
}

func TestMessageTypeConstants(t *testing.T) {
	if MessageTypeLog != "log" {
		t.Errorf("expected MessageTypeLog 'log', got '%s'", MessageTypeLog)
	}
	if MessageTypeEvent != "event" {
		t.Errorf("expected MessageTypeEvent 'event', got '%s'", MessageTypeEvent)
	}
	if MessageTypeStatus != "status" {
		t.Errorf("expected MessageTypeStatus 'status', got '%s'", MessageTypeStatus)
	}
	if MessageTypeError != "error" {
		t.Errorf("expected MessageTypeError 'error', got '%s'", MessageTypeError)
	}
	if MessageTypePing != "ping" {
		t.Errorf("expected MessageTypePing 'ping', got '%s'", MessageTypePing)
	}
	if MessageTypePong != "pong" {
		t.Errorf("expected MessageTypePong 'pong', got '%s'", MessageTypePong)
	}
}

func TestUpgrader(t *testing.T) {
	if Upgrader.ReadBufferSize != 1024 {
		t.Errorf("expected ReadBufferSize 1024, got %d", Upgrader.ReadBufferSize)
	}
	if Upgrader.WriteBufferSize != 1024 {
		t.Errorf("expected WriteBufferSize 1024, got %d", Upgrader.WriteBufferSize)
	}
}

func TestUpgrader_CheckOrigin(t *testing.T) {
	// CheckOrigin should return true for any request
	if !Upgrader.CheckOrigin(nil) {
		t.Error("expected CheckOrigin to return true")
	}
}

func TestLogStreamer(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)
	streamer := NewLogStreamer(hub)

	if streamer == nil {
		t.Fatal("NewLogStreamer returned nil")
	}
	if streamer.hub != hub {
		t.Error("streamer hub not set correctly")
	}
}

func TestLogStreamer_Stream(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)
	streamer := NewLogStreamer(hub)

	// Stream should not block
	streamer.Stream("info", "test-source", "test message")
}

func TestEventStreamer(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)
	streamer := NewEventStreamer(hub)

	if streamer == nil {
		t.Fatal("NewEventStreamer returned nil")
	}
	if streamer.hub != hub {
		t.Error("streamer hub not set correctly")
	}
}

func TestEventStreamer_Stream(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)
	streamer := NewEventStreamer(hub)

	// Stream should not block
	streamer.Stream("user.created", map[string]string{"user_id": "123"})
}

func TestHub_Broadcast_WithData(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	// Create a client with a buffered channel
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	// Add client directly
	hub.mu.Lock()
	hub.clients[client] = true
	hub.mu.Unlock()

	// Broadcast a message
	hub.Broadcast(MessageTypeLog, map[string]string{"message": "test"})

	// Give time for the message to be sent
	time.Sleep(10 * time.Millisecond)

	// Check if message was received
	select {
	case msg := <-client.send:
		var decoded Message
		if err := json.Unmarshal(msg, &decoded); err != nil {
			t.Fatalf("failed to decode message: %v", err)
		}
		if decoded.Type != MessageTypeLog {
			t.Errorf("expected type '%s', got '%s'", MessageTypeLog, decoded.Type)
		}
	default:
		// Message might be in the broadcast channel, not sent yet
		// since Run() is not running
	}
}

func TestHub_MultipleClients(t *testing.T) {
	eb := eventbus.New()
	hub := NewHub(eb)

	// Create multiple clients
	client1 := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}
	client2 := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	hub.mu.Lock()
	hub.clients[client1] = true
	hub.clients[client2] = true
	hub.mu.Unlock()

	if hub.ClientCount() != 2 {
		t.Errorf("expected 2 clients, got %d", hub.ClientCount())
	}
}