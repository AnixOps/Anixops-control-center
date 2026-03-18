package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	if tb == nil {
		t.Fatal("TokenBucket is nil")
	}
	if tb.Rate() != 10 {
		t.Errorf("Expected rate 10, got %f", tb.Rate())
	}
	if tb.Burst() != 5 {
		t.Errorf("Expected burst 5, got %d", tb.Burst())
	}
}

func TestTokenBucketAllow(t *testing.T) {
	tb := NewTokenBucket(10, 3)

	// Should allow 3 times (burst)
	for i := 0; i < 3; i++ {
		if !tb.Allow() {
			t.Errorf("Should allow operation %d", i+1)
		}
	}

	// Should not allow after burst exhausted
	if tb.Allow() {
		t.Error("Should not allow after burst exhausted")
	}
}

func TestTokenBucketWait(t *testing.T) {
	tb := NewTokenBucket(100, 1)

	// Consume the token
	if !tb.Allow() {
		t.Fatal("Should allow first operation")
	}

	// Wait should succeed after token refill
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := tb.Wait(ctx)
	if err != nil {
		t.Errorf("Wait should succeed: %v", err)
	}
}

func TestTokenBucketWaitCancelled(t *testing.T) {
	tb := NewTokenBucket(1, 1)

	// Consume the token
	if !tb.Allow() {
		t.Fatal("Should allow first operation")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := tb.Wait(ctx)
	if err == nil {
		t.Error("Wait should fail with cancelled context")
	}
}

func TestTokenBucketReserve(t *testing.T) {
	tb := NewTokenBucket(10, 5)

	res := tb.Reserve()
	if !res.OK() {
		t.Error("Reservation should be valid")
	}
}

func TestTokenBucketReserveDelay(t *testing.T) {
	tb := NewTokenBucket(100, 1)

	// Consume the token
	if !tb.Allow() {
		t.Fatal("Should allow first operation")
	}

	res := tb.Reserve()
	if !res.OK() {
		t.Error("Reservation should be valid")
	}

	delay := res.Delay()
	if delay <= 0 {
		t.Error("Delay should be positive when bucket empty")
	}
}

func TestTokenBucketReservationCancel(t *testing.T) {
	tb := NewTokenBucket(10, 1)

	res := tb.Reserve()
	res.Cancel()

	// Token should be returned
	if !tb.Allow() {
		t.Error("Should allow after cancel")
	}
}

func TestTokenBucketAllowN(t *testing.T) {
	tb := NewTokenBucket(10, 5)

	if !tb.AllowN(3) {
		t.Error("Should allow 3 tokens")
	}

	if tb.AllowN(3) {
		t.Error("Should not allow 3 more tokens (only 2 left)")
	}
}

func TestTokenBucketWaitN(t *testing.T) {
	tb := NewTokenBucket(100, 5)

	err := tb.WaitN(context.Background(), 3)
	if err != nil {
		t.Errorf("WaitN should succeed: %v", err)
	}
}

func TestTokenBucketWaitNExceedsBurst(t *testing.T) {
	tb := NewTokenBucket(100, 5)

	err := tb.WaitN(context.Background(), 10)
	if err == nil {
		t.Error("Should fail when n exceeds burst")
	}
}

func TestNewSlidingWindow(t *testing.T) {
	sw := NewSlidingWindow(5, time.Second)
	if sw == nil {
		t.Fatal("SlidingWindow is nil")
	}
}

func TestSlidingWindowAllow(t *testing.T) {
	sw := NewSlidingWindow(3, 100*time.Millisecond)

	// Should allow 3 times
	for i := 0; i < 3; i++ {
		if !sw.Allow() {
			t.Errorf("Should allow operation %d", i+1)
		}
	}

	// Should not allow
	if sw.Allow() {
		t.Error("Should not allow after limit reached")
	}

	// Wait for window to slide
	time.Sleep(150 * time.Millisecond)

	// Should allow again
	if !sw.Allow() {
		t.Error("Should allow after window slides")
	}
}

func TestSlidingWindowWait(t *testing.T) {
	sw := NewSlidingWindow(1, 50*time.Millisecond)

	// Consume the token
	if !sw.Allow() {
		t.Fatal("Should allow first operation")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := sw.Wait(ctx)
	if err != nil {
		t.Errorf("Wait should succeed: %v", err)
	}
}

func TestSlidingWindowWaitCancelled(t *testing.T) {
	sw := NewSlidingWindow(1, time.Second)

	// Consume the token
	if !sw.Allow() {
		t.Fatal("Should allow first operation")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sw.Wait(ctx)
	if err == nil {
		t.Error("Wait should fail with cancelled context")
	}
}

func TestNewLeakyBucket(t *testing.T) {
	lb := NewLeakyBucket(10*time.Millisecond, 5)
	if lb == nil {
		t.Fatal("LeakyBucket is nil")
	}
}

func TestLeakyBucketAllow(t *testing.T) {
	lb := NewLeakyBucket(time.Second, 3)

	// Should allow 3 times
	for i := 0; i < 3; i++ {
		if !lb.Allow() {
			t.Errorf("Should allow operation %d", i+1)
		}
	}

	// Should not allow
	if lb.Allow() {
		t.Error("Should not allow after capacity reached")
	}
}

func TestLeakyBucketWait(t *testing.T) {
	lb := NewLeakyBucket(10*time.Millisecond, 1)

	// Consume the token
	if !lb.Allow() {
		t.Fatal("Should allow first operation")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := lb.Wait(ctx)
	if err != nil {
		t.Errorf("Wait should succeed: %v", err)
	}
}

func TestLeakyBucketWaitCancelled(t *testing.T) {
	lb := NewLeakyBucket(time.Second, 1)

	// Consume the token
	if !lb.Allow() {
		t.Fatal("Should allow first operation")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := lb.Wait(ctx)
	if err == nil {
		t.Error("Wait should fail with cancelled context")
	}
}

func TestLimitExceededError(t *testing.T) {
	err := &LimitExceededError{Requested: 10, Limit: 5}
	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}

func TestIsLimitExceeded(t *testing.T) {
	if !IsLimitExceeded(&LimitExceededError{}) {
		t.Error("Should be limit exceeded error")
	}

	if IsLimitExceeded(nil) {
		t.Error("nil should not be limit exceeded error")
	}
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("Manager is nil")
	}
}

func TestManagerRegister(t *testing.T) {
	m := NewManager()
	tb := NewTokenBucket(10, 5)
	m.Register("test", tb)

	limiter, ok := m.Get("test")
	if !ok {
		t.Error("Should find registered limiter")
	}
	if limiter != tb {
		t.Error("Should return same limiter")
	}
}

func TestManagerGetNotFound(t *testing.T) {
	m := NewManager()

	_, ok := m.Get("nonexistent")
	if ok {
		t.Error("Should not find unregistered limiter")
	}
}

func TestManagerAllow(t *testing.T) {
	m := NewManager()
	m.Register("test", NewTokenBucket(10, 1))

	if !m.Allow("test") {
		t.Error("Should allow")
	}
	if m.Allow("test") {
		t.Error("Should not allow second time")
	}
}

func TestManagerAllowNoLimit(t *testing.T) {
	m := NewManager()

	// No limit registered, should always allow
	if !m.Allow("nonexistent") {
		t.Error("Should allow when no limit configured")
	}
}

func TestManagerWait(t *testing.T) {
	m := NewManager()
	m.Register("test", NewTokenBucket(100, 1))

	if !m.Allow("test") {
		t.Fatal("Should allow first operation")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := m.Wait(ctx, "test")
	if err != nil {
		t.Errorf("Wait should succeed: %v", err)
	}
}

func TestManagerWaitNoLimit(t *testing.T) {
	m := NewManager()

	err := m.Wait(context.Background(), "nonexistent")
	if err != nil {
		t.Errorf("Wait should succeed with no limit: %v", err)
	}
}

func TestManagerKeys(t *testing.T) {
	m := NewManager()
	m.Register("a", NewTokenBucket(1, 1))
	m.Register("b", NewTokenBucket(1, 1))

	keys := m.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestConcurrentTokenBucket(t *testing.T) {
	tb := NewTokenBucket(1000, 100)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				tb.Allow()
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestConcurrentManager(t *testing.T) {
	m := NewManager()
	m.Register("test", NewTokenBucket(1000, 100))

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				m.Allow("test")
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}