package eventbus

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	eb := New()
	if eb == nil {
		t.Fatal("New() returned nil")
	}
	if eb.handlers == nil {
		t.Error("handlers map not initialized")
	}
	if eb.buffer != 100 {
		t.Errorf("expected buffer 100, got %d", eb.buffer)
	}
}

func TestSubscribe(t *testing.T) {
	eb := New()
	called := false
	eb.Subscribe("test.event", func(ctx context.Context, event Event) error {
		called = true
		return nil
	})

	err := eb.Publish(context.Background(), Event{Type: "test.event"})
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestSubscribeAll(t *testing.T) {
	eb := New()
	var count int32

	eb.SubscribeAll(func(ctx context.Context, event Event) error {
		atomic.AddInt32(&count, 1)
		return nil
	})

	eb.Publish(context.Background(), Event{Type: "event1"})
	eb.Publish(context.Background(), Event{Type: "event2"})
	eb.Publish(context.Background(), Event{Type: "event3"})

	if atomic.LoadInt32(&count) != 3 {
		t.Errorf("expected 3 calls, got %d", count)
	}
}

func TestUnsubscribe(t *testing.T) {
	eb := New()
	called := false

	eb.Subscribe("test.event", func(ctx context.Context, event Event) error {
		called = true
		return nil
	})

	eb.Unsubscribe("test.event")
	eb.Publish(context.Background(), Event{Type: "test.event"})

	if called {
		t.Error("handler called after unsubscribe")
	}
}

func TestPublish_MultipleHandlers(t *testing.T) {
	eb := New()
	var count int32

	for i := 0; i < 3; i++ {
		eb.Subscribe("test.event", func(ctx context.Context, event Event) error {
			atomic.AddInt32(&count, 1)
			return nil
		})
	}

	eb.Publish(context.Background(), Event{Type: "test.event"})

	if atomic.LoadInt32(&count) != 3 {
		t.Errorf("expected 3 calls, got %d", count)
	}
}

func TestPublish_HandlerError(t *testing.T) {
	eb := New()

	eb.Subscribe("test.event", func(ctx context.Context, event Event) error {
		return nil
	})
	eb.Subscribe("test.event", func(ctx context.Context, event Event) error {
		return context.DeadlineExceeded
	})

	err := eb.Publish(context.Background(), Event{Type: "test.event"})
	if err == nil {
		t.Error("expected error from failing handler")
	}
}

func TestPublishAsync(t *testing.T) {
	eb := New()
	var called int32

	eb.Subscribe("test.event", func(ctx context.Context, event Event) error {
		atomic.StoreInt32(&called, 1)
		return nil
	})

	eb.PublishAsync(context.Background(), Event{Type: "test.event"})

	// Wait for async publish
	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&called) != 1 {
		t.Error("async handler not called")
	}
}

func TestChannel(t *testing.T) {
	eb := New()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := eb.Channel(ctx, "test.event")
	if ch == nil {
		t.Fatal("Channel returned nil")
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		eb.Publish(ctx, Event{Type: "test.event", Data: "test-data"})
	}()

	select {
	case event := <-ch:
		if event.Type != "test.event" {
			t.Errorf("expected event type 'test.event', got %s", event.Type)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for event")
	}
}

func TestSetBuffer(t *testing.T) {
	eb := New()
	eb.SetBuffer(50)
	if eb.buffer != 50 {
		t.Errorf("expected buffer 50, got %d", eb.buffer)
	}
}

func TestEventTypes(t *testing.T) {
	types := []string{
		EventSystemStart,
		EventSystemShutdown,
		EventPluginRegistered,
		EventPluginStarted,
		EventNodeAdded,
		EventTaskCreated,
		EventUserLogin,
		EventPlaybookStart,
	}

	for _, eventType := range types {
		if eventType == "" {
			t.Errorf("empty event type")
		}
	}
}