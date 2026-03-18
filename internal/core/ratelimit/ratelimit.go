package ratelimit

import (
	"context"
	"sync"
	"time"
)

// RateLimiter limits the rate of operations
type RateLimiter interface {
	// Wait waits until the rate limiter allows the operation
	Wait(ctx context.Context) error
	// Allow returns true if the operation is allowed
	Allow() bool
	// Reserve reserves a permit for future use
	Reserve() Reservation
}

// Reservation represents a reserved permit
type Reservation struct {
	ok        bool
	timeToAct time.Time
	limiter   *TokenBucket
	tokens    int
}

// OK returns true if the reservation is valid
func (r *Reservation) OK() bool {
	return r.ok
}

// Delay returns the duration to wait before the reservation can be used
func (r *Reservation) Delay() time.Duration {
	if !r.ok {
		return 0
	}
	return time.Until(r.timeToAct)
}

// Cancel cancels the reservation
func (r *Reservation) Cancel() {
	if !r.ok {
		return
	}
	r.limiter.cancel(r.tokens, r.timeToAct)
}

// TokenBucket implements a token bucket rate limiter
type TokenBucket struct {
	mu         sync.Mutex
	rate       float64       // tokens per second
	burst      int           // maximum bucket size
	tokens     float64       // current tokens
	lastUpdate time.Time     // last update time
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(rate float64, burst int) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastUpdate: time.Now(),
	}
}

// Wait waits until a token is available
func (tb *TokenBucket) Wait(ctx context.Context) error {
	return tb.waitN(ctx, 1)
}

// WaitN waits until n tokens are available
func (tb *TokenBucket) WaitN(ctx context.Context, n int) error {
	return tb.waitN(ctx, n)
}

func (tb *TokenBucket) waitN(ctx context.Context, n int) error {
	if n > tb.burst {
		return &LimitExceededError{Requested: n, Limit: tb.burst}
	}

	res := tb.reserveN(n)
	if !res.ok {
		return &LimitExceededError{Requested: n, Limit: 0}
	}

	delay := res.Delay()
	if delay <= 0 {
		return nil
	}

	t := time.NewTimer(delay)
	defer t.Stop()

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		res.Cancel()
		return ctx.Err()
	}
}

// Allow returns true if a token can be consumed immediately
func (tb *TokenBucket) Allow() bool {
	return tb.allowN(1)
}

// AllowN returns true if n tokens can be consumed immediately
func (tb *TokenBucket) AllowN(n int) bool {
	return tb.allowN(n)
}

func (tb *TokenBucket) allowN(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.updateTokens()
	if float64(n) <= tb.tokens {
		tb.tokens -= float64(n)
		return true
	}
	return false
}

// Reserve reserves a token
func (tb *TokenBucket) Reserve() Reservation {
	return tb.reserveN(1)
}

// ReserveN reserves n tokens
func (tb *TokenBucket) ReserveN(n int) Reservation {
	return tb.reserveN(n)
}

func (tb *TokenBucket) reserveN(n int) Reservation {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.updateTokens()

	// Calculate time to wait
	var waitTime time.Duration
	if tb.tokens < float64(n) {
		needed := float64(n) - tb.tokens
		waitTime = time.Duration(needed / tb.rate * float64(time.Second))
	}

	ok := n <= tb.burst
	tb.tokens -= float64(n)

	return Reservation{
		ok:        ok,
		timeToAct: time.Now().Add(waitTime),
		limiter:   tb,
		tokens:    n,
	}
}

func (tb *TokenBucket) cancel(tokens int, timeToAct time.Time) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.updateTokens()
	tb.tokens += float64(tokens)
}

func (tb *TokenBucket) updateTokens() {
	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate)
	tb.lastUpdate = now

	tb.tokens += float64(elapsed) / float64(time.Second) * tb.rate
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}
}

// Rate returns the rate limit
func (tb *TokenBucket) Rate() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.rate
}

// Burst returns the burst size
func (tb *TokenBucket) Burst() int {
	return tb.burst
}

// SlidingWindow implements a sliding window rate limiter
type SlidingWindow struct {
	mu        sync.Mutex
	rate      int
	window    time.Duration
	timestamps []time.Time
}

// NewSlidingWindow creates a new sliding window rate limiter
func NewSlidingWindow(rate int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		rate:       rate,
		window:     window,
		timestamps: make([]time.Time, 0),
	}
}

// Allow returns true if the operation is allowed
func (sw *SlidingWindow) Allow() bool {
	return sw.allowN(1)
}

// AllowN returns true if n operations are allowed
func (sw *SlidingWindow) AllowN(n int) bool {
	return sw.allowN(n)
}

func (sw *SlidingWindow) allowN(n int) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	// Remove old timestamps
	valid := make([]time.Time, 0, len(sw.timestamps))
	for _, ts := range sw.timestamps {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	sw.timestamps = valid

	// Check if we can allow
	if len(sw.timestamps)+n <= sw.rate {
		for i := 0; i < n; i++ {
			sw.timestamps = append(sw.timestamps, now)
		}
		return true
	}
	return false
}

// Wait waits until the operation is allowed
func (sw *SlidingWindow) Wait(ctx context.Context) error {
	for {
		if sw.Allow() {
			return nil
		}

		// Calculate wait time
		sw.mu.Lock()
		if len(sw.timestamps) == 0 {
			sw.mu.Unlock()
			continue
		}

		oldest := sw.timestamps[0]
		waitTime := oldest.Add(sw.window).Sub(time.Now())
		sw.mu.Unlock()

		if waitTime <= 0 {
			continue
		}

		select {
		case <-time.After(waitTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Reserve returns a reservation (not supported for sliding window)
func (sw *SlidingWindow) Reserve() Reservation {
	return Reservation{ok: false}
}

// LeakyBucket implements a leaky bucket rate limiter
type LeakyBucket struct {
	mu        sync.Mutex
	rate      time.Duration // time between drips
	capacity  int
	water     int
	lastDrip  time.Time
}

// NewLeakyBucket creates a new leaky bucket rate limiter
func NewLeakyBucket(rate time.Duration, capacity int) *LeakyBucket {
	return &LeakyBucket{
		rate:     rate,
		capacity: capacity,
		lastDrip: time.Now(),
	}
}

// Allow returns true if the operation is allowed
func (lb *LeakyBucket) Allow() bool {
	return lb.allowN(1)
}

func (lb *LeakyBucket) allowN(n int) bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.drip()

	if lb.water+n <= lb.capacity {
		lb.water += n
		return true
	}
	return false
}

func (lb *LeakyBucket) drip() {
	now := time.Now()
	elapsed := now.Sub(lb.lastDrip)
	drips := int(elapsed / lb.rate)
	lb.water -= drips
	if lb.water < 0 {
		lb.water = 0
	}
	lb.lastDrip = now
}

// Wait waits until the operation is allowed
func (lb *LeakyBucket) Wait(ctx context.Context) error {
	for {
		if lb.Allow() {
			return nil
		}

		select {
		case <-time.After(lb.rate):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Reserve returns a reservation (not supported for leaky bucket)
func (lb *LeakyBucket) Reserve() Reservation {
	return Reservation{ok: false}
}

// LimitExceededError is returned when rate limit is exceeded
type LimitExceededError struct {
	Requested int
	Limit     int
}

func (e *LimitExceededError) Error() string {
	return "rate limit exceeded"
}

// IsLimitExceeded checks if an error is a limit exceeded error
func IsLimitExceeded(err error) bool {
	_, ok := err.(*LimitExceededError)
	return ok
}

// Manager manages multiple rate limiters
type Manager struct {
	mu        sync.RWMutex
	limiters  map[string]RateLimiter
}

// NewManager creates a new rate limiter manager
func NewManager() *Manager {
	return &Manager{
		limiters: make(map[string]RateLimiter),
	}
}

// Register registers a rate limiter for a key
func (m *Manager) Register(key string, limiter RateLimiter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[key] = limiter
}

// Get gets a rate limiter for a key
func (m *Manager) Get(key string) (RateLimiter, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	limiter, ok := m.limiters[key]
	return limiter, ok
}

// Allow checks if an operation is allowed for a key
func (m *Manager) Allow(key string) bool {
	limiter, ok := m.Get(key)
	if !ok {
		return true // No limit configured
	}
	return limiter.Allow()
}

// Wait waits for permission to proceed for a key
func (m *Manager) Wait(ctx context.Context, key string) error {
	limiter, ok := m.Get(key)
	if !ok {
		return nil // No limit configured
	}
	return limiter.Wait(ctx)
}

// Keys returns all registered keys
func (m *Manager) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.limiters))
	for k := range m.limiters {
		keys = append(keys, k)
	}
	return keys
}