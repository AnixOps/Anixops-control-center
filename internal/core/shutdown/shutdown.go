package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/logger"
)

// Handler handles graceful shutdown
type Handler struct {
	mu           sync.Mutex
	shutdownFns  []func(ctx context.Context) error
	timeout      time.Duration
	shuttingDown bool
	done         chan struct{}
}

// NewHandler creates a new shutdown handler
func NewHandler() *Handler {
	return &Handler{
		shutdownFns: make([]func(ctx context.Context) error, 0),
		timeout:     30 * time.Second,
		done:        make(chan struct{}),
	}
}

// SetTimeout sets the shutdown timeout
func (h *Handler) SetTimeout(d time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.timeout = d
}

// Register registers a shutdown function
func (h *Handler) Register(fn func(ctx context.Context) error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.shutdownFns = append(h.shutdownFns, fn)
}

// IsShuttingDown returns true if shutdown is in progress
func (h *Handler) IsShuttingDown() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.shuttingDown
}

// Done returns a channel that is closed when shutdown completes
func (h *Handler) Done() <-chan struct{} {
	return h.done
}

// Shutdown initiates shutdown
func (h *Handler) Shutdown(ctx context.Context) error {
	h.mu.Lock()
	if h.shuttingDown {
		h.mu.Unlock()
		return nil
	}
	h.shuttingDown = true
	fns := make([]func(ctx context.Context) error, len(h.shutdownFns))
	copy(fns, h.shutdownFns)
	timeout := h.timeout
	h.mu.Unlock()

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	logger.Info("Starting graceful shutdown")

	// Run shutdown functions in reverse order
	for i := len(fns) - 1; i >= 0; i-- {
		if err := fns[i](ctx); err != nil {
			logger.Error("Shutdown function failed", logger.F("error", err))
		}
	}

	logger.Info("Graceful shutdown complete")
	close(h.done)
	return nil
}

// Wait waits for shutdown signals and handles them
func (h *Handler) Wait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sig := <-sigChan
	logger.Info("Received shutdown signal", logger.F("signal", sig.String()))

	_ = h.Shutdown(context.Background())
}

// WaitWithContext waits for shutdown signals or context cancellation
func (h *Handler) WaitWithContext(ctx context.Context) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", logger.F("signal", sig.String()))
		return h.Shutdown(ctx)
	case <-ctx.Done():
		logger.Info("Context cancelled, initiating shutdown")
		return h.Shutdown(ctx)
	}
}

// DefaultHandler is the default shutdown handler
var DefaultHandler = NewHandler()

// Register registers a shutdown function with the default handler
func Register(fn func(ctx context.Context) error) {
	DefaultHandler.Register(fn)
}

// Wait waits for shutdown signals with the default handler
func Wait() {
	DefaultHandler.Wait()
}

// WaitWithContext waits with context using the default handler
func WaitWithContext(ctx context.Context) error {
	return DefaultHandler.WaitWithContext(ctx)
}

// IsShuttingDown returns true if the default handler is shutting down
func IsShuttingDown() bool {
	return DefaultHandler.IsShuttingDown()
}

// Done returns the done channel of the default handler
func Done() <-chan struct{} {
	return DefaultHandler.Done()
}

// Shutdown initiates shutdown with the default handler
func Shutdown(ctx context.Context) error {
	return DefaultHandler.Shutdown(ctx)
}

// SetTimeout sets the timeout for the default handler
func SetTimeout(d time.Duration) {
	DefaultHandler.SetTimeout(d)
}