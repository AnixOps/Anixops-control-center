// Package sdk provides a commercial-grade SDK for building AnixOps plugins and extensions.
//
// The SDK provides:
//   - Plugin system with lifecycle management
//   - Dependency injection container
//   - Structured logging
//   - Metrics collection
//   - Error handling with codes
//   - Configuration validation
//   - Graceful shutdown
//   - Event bus for pub/sub
//
// Example usage:
//
//	func main() {
//	    // Create SDK instance
//	    sdk := anixops.New(&anixops.Options{
//	        ConfigPath: "config.yaml",
//	    })
//
//	    // Register plugins
//	    sdk.RegisterPlugin("my-plugin", myPlugin)
//
//	    // Start the application
//	    if err := sdk.Start(); err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Wait for shutdown
//	    sdk.Wait()
//	}
package anixops

import (
	"context"
	"sync"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/config"
	"github.com/anixops/anixops-control-center/internal/core/container"
	"github.com/anixops/anixops-control-center/internal/core/errors"
	"github.com/anixops/anixops-control-center/internal/core/eventbus"
	"github.com/anixops/anixops-control-center/internal/core/lifecycle"
	"github.com/anixops/anixops-control-center/internal/core/logger"
	"github.com/anixops/anixops-control-center/internal/core/metrics"
	"github.com/anixops/anixops-control-center/internal/core/plugin"
	"github.com/anixops/anixops-control-center/internal/core/shutdown"
	"github.com/anixops/anixops-control-center/internal/core/validator"
)

// Version is the SDK version
const Version = "1.0.0"

// Options holds SDK initialization options
type Options struct {
	// ConfigPath is the path to the configuration file
	ConfigPath string

	// Config is a pre-loaded configuration
	Config *config.Config

	// Logger is a custom logger
	Logger *logger.Logger

	// MetricsProvider is a custom metrics provider
	MetricsProvider *metrics.MetricsProvider

	// ShutdownTimeout is the timeout for graceful shutdown
	ShutdownTimeout int // seconds

	// JSONLogging enables JSON formatted logging
	JSONLogging bool

	// LogLevel sets the logging level
	LogLevel string
}

// SDK is the main SDK instance
type SDK struct {
	mu           sync.RWMutex
	config       *config.Config
	container    *container.Container
	pluginMgr    *plugin.Manager
	eventBus     *eventbus.EventBus
	lifecycle    *lifecycle.Manager
	metrics      *metrics.MetricsProvider
	logger       *logger.Logger
	shutdown     *shutdown.Handler
	validator    *validator.ConfigValidator
	depManager   *plugin.DependencyManager
	hookRegistry *plugin.HookRegistry
	started      bool
	done         chan struct{}
}

// New creates a new SDK instance
func New(opts *Options) (*SDK, error) {
	s := &SDK{
		container:    container.New(),
		pluginMgr:    plugin.NewManager(),
		eventBus:     eventbus.New(),
		lifecycle:    lifecycle.NewManager(),
		metrics:      metrics.NewMetricsProvider(),
		shutdown:     shutdown.NewHandler(),
		validator:    validator.NewConfigValidator(),
		depManager:   plugin.NewDependencyManager(),
		hookRegistry: plugin.NewHookRegistry(),
		done:         make(chan struct{}),
	}

	// Initialize logger
	if opts != nil && opts.Logger != nil {
		s.logger = opts.Logger
	} else {
		logLevel := logger.InfoLevel
		if opts != nil && opts.LogLevel != "" {
			switch opts.LogLevel {
			case "debug":
				logLevel = logger.DebugLevel
			case "warn":
				logLevel = logger.WarnLevel
			case "error":
				logLevel = logger.ErrorLevel
			}
		}
		s.logger = logger.New(logger.Options{
			Level: logLevel,
			JSON:  opts != nil && opts.JSONLogging,
		})
	}

	// Load configuration
	if opts != nil {
		if opts.Config != nil {
			s.config = opts.Config
		} else if opts.ConfigPath != "" {
			cfg, err := config.Load(opts.ConfigPath)
			if err != nil {
				return nil, errors.NewError(errors.CodeConfigNotFound,
					"Failed to load configuration", errors.LevelError, 500).
					WithCause(err)
			}
			s.config = cfg
		} else {
			s.config = config.DefaultConfig()
		}

		// Set shutdown timeout
		if opts.ShutdownTimeout > 0 {
			s.shutdown.SetTimeout(time.Duration(opts.ShutdownTimeout) * time.Second)
		}
	} else {
		s.config = config.DefaultConfig()
	}

	// Register core services in container
	s.registerServices()

	return s, nil
}

// registerServices registers core services in the container
func (s *SDK) registerServices() {
	s.container.Register(container.ServiceConfig, s.config)
	s.container.Register(container.ServicePluginMgr, s.pluginMgr)
	s.container.Register(container.ServiceEventBus, s.eventBus)
	s.container.Register(container.ServiceLogger, s.logger)
	s.container.Register(container.ServiceMetrics, s.metrics)
	s.container.Register(container.ServiceScheduler, s.lifecycle)
}

// Config returns the configuration
func (s *SDK) Config() *config.Config {
	return s.config
}

// Container returns the dependency injection container
func (s *SDK) Container() *container.Container {
	return s.container
}

// PluginManager returns the plugin manager
func (s *SDK) PluginManager() *plugin.Manager {
	return s.pluginMgr
}

// EventBus returns the event bus
func (s *SDK) EventBus() *eventbus.EventBus {
	return s.eventBus
}

// Logger returns the logger
func (s *SDK) Logger() *logger.Logger {
	return s.logger
}

// Metrics returns the metrics provider
func (s *SDK) Metrics() *metrics.MetricsProvider {
	return s.metrics
}

// RegisterPlugin registers a plugin
func (s *SDK) RegisterPlugin(name string, p plugin.Plugin) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.pluginMgr.Register(name, p); err != nil {
		return errors.NewError(errors.CodePluginNotFound,
			"Failed to register plugin", errors.LevelError, 500).
			WithCause(err).
			WithContext("plugin", name)
	}

	// Register dependencies if plugin implements DependentPlugin
	if dp, ok := p.(plugin.DependentPlugin); ok {
		s.depManager.Register(name, dp.Dependencies())
	}

	// Register hooks if plugin implements HookProvider
	if hp, ok := p.(plugin.HookProvider); ok {
		for _, hook := range hp.Hooks() {
			s.hookRegistry.Register(hook)
		}
	}

	// Register capabilities as providers
	for _, cap := range p.Capabilities() {
		s.depManager.RegisterProvider(cap, name)
	}

	s.logger.Info("Plugin registered", logger.F("plugin", name))
	return nil
}

// UnregisterPlugin unregisters a plugin
func (s *SDK) UnregisterPlugin(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.pluginMgr.Unregister(name); err != nil {
		return errors.NewError(errors.CodePluginNotFound,
			"Failed to unregister plugin", errors.LevelError, 500).
			WithCause(err).
			WithContext("plugin", name)
	}

	s.logger.Info("Plugin unregistered", logger.F("plugin", name))
	return nil
}

// Start starts the SDK
func (s *SDK) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return errors.NewError(errors.CodeInternal,
			"SDK already started", errors.LevelError, 500)
	}
	s.mu.Unlock()

	s.logger.Info("Starting AnixOps SDK")

	// Run pre-start hooks
	if err := s.hookRegistry.Execute(ctx, plugin.HookPreStart, nil); err != nil {
		return err
	}

	// Resolve plugin dependencies and get start order
	plugins := s.pluginMgr.List()
	startOrder, err := s.depManager.Resolve(plugins)
	if err != nil {
		return errors.NewError(errors.CodePluginDependency,
			"Failed to resolve plugin dependencies", errors.LevelError, 500).
			WithCause(err)
	}

	// Initialize and start plugins in dependency order
	for _, name := range startOrder {
		p, ok := s.pluginMgr.Get(name)
		if !ok {
			continue
		}

		// Run pre-init hook
		if err := s.hookRegistry.Execute(ctx, plugin.HookPreInit, name); err != nil {
			return err
		}

		// Initialize
		cfg := s.pluginMgr.GetConfig(name)
		if cfg == nil {
			// Try to get config from SDK config
			if pluginConfigs, ok := s.config.Plugins.Ansible[name].(map[string]interface{}); ok {
				cfg = pluginConfigs
			}
		}
		if err := s.pluginMgr.InitPlugin(ctx, name); err != nil {
			return errors.NewError(errors.CodePluginInitFailed,
				"Failed to initialize plugin", errors.LevelError, 500).
				WithCause(err).
				WithContext("plugin", name)
		}

		// Run post-init hook
		if err := s.hookRegistry.Execute(ctx, plugin.HookPostInit, name); err != nil {
			return err
		}

		// Run pre-start hook
		if err := s.hookRegistry.Execute(ctx, plugin.HookPreStart, name); err != nil {
			return err
		}

		// Start
		if err := s.pluginMgr.StartPlugin(ctx, name); err != nil {
			return errors.NewError(errors.CodePluginStartFailed,
				"Failed to start plugin", errors.LevelError, 500).
				WithCause(err).
				WithContext("plugin", name)
		}

		// Run post-start hook
		if err := s.hookRegistry.Execute(ctx, plugin.HookPostStart, name); err != nil {
			return err
		}

		_ = p // use p to avoid unused variable
		s.logger.Info("Plugin started", logger.F("plugin", name))
	}

	// Register shutdown handler
	s.shutdown.Register(func(ctx context.Context) error {
		return s.Stop(ctx)
	})

	s.mu.Lock()
	s.started = true
	s.mu.Unlock()

	s.logger.Info("AnixOps SDK started successfully")
	return nil
}

// Stop stops the SDK
func (s *SDK) Stop(ctx context.Context) error {
	s.mu.RLock()
	if !s.started {
		s.mu.RUnlock()
		return nil
	}
	s.mu.RUnlock()

	s.logger.Info("Stopping AnixOps SDK")

	// Run pre-stop hooks
	if err := s.hookRegistry.Execute(ctx, plugin.HookPreStop, nil); err != nil {
		s.logger.Error("Pre-stop hook failed", logger.F("error", err))
	}

	// Stop plugins in reverse order
	plugins := s.pluginMgr.List()
	stopOrder, _ := s.depManager.Resolve(plugins)
	// Reverse the order
	for i, j := 0, len(stopOrder)-1; i < j; i, j = i+1, j-1 {
		stopOrder[i], stopOrder[j] = stopOrder[j], stopOrder[i]
	}

	for _, name := range stopOrder {
		if s.pluginMgr.GetState(name) == plugin.StateRunning {
			if err := s.pluginMgr.StopPlugin(ctx, name); err != nil {
				s.logger.Error("Failed to stop plugin",
					logger.F("plugin", name),
					logger.F("error", err))
			} else {
				s.logger.Info("Plugin stopped", logger.F("plugin", name))
			}
		}
	}

	// Run post-stop hooks
	if err := s.hookRegistry.Execute(ctx, plugin.HookPostStop, nil); err != nil {
		s.logger.Error("Post-stop hook failed", logger.F("error", err))
	}

	s.mu.Lock()
	s.started = false
	close(s.done) // Signal that SDK has stopped
	s.mu.Unlock()

	s.logger.Info("AnixOps SDK stopped")
	return nil
}

// Wait waits for shutdown signals
func (s *SDK) Wait() {
	s.shutdown.Wait()
}

// WaitForShutdown blocks until shutdown is complete
func (s *SDK) WaitForShutdown() <-chan struct{} {
	return s.done
}

// IsRunning returns true if the SDK is running
func (s *SDK) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.started
}

// ExecutePlugin executes a plugin action
func (s *SDK) ExecutePlugin(ctx context.Context, name, action string, params map[string]interface{}) (plugin.Result, error) {
	p, ok := s.pluginMgr.Get(name)
	if !ok {
		return plugin.Result{}, errors.NewError(errors.CodePluginNotFound,
			"Plugin not found", errors.LevelError, 404).
			WithContext("plugin", name)
	}

	execPlugin, ok := p.(plugin.ExecutablePlugin)
	if !ok {
		return plugin.Result{}, errors.NewError(errors.CodePluginNotExecutable,
			"Plugin is not executable", errors.LevelError, 400).
			WithContext("plugin", name)
	}

	// Run pre-execute hook
	if err := s.hookRegistry.Execute(ctx, plugin.HookPreExecute, map[string]interface{}{
		"plugin": name,
		"action": action,
		"params": params,
	}); err != nil {
		return plugin.Result{}, err
	}

	// Execute
	startTime := ctx.Value("start_time") // Track execution time
	_ = startTime
	result, err := execPlugin.Execute(ctx, action, params)

	// Record metrics
	s.metrics.PluginMetrics().IncExecutions()
	if duration, ok := ctx.Value("duration").(time.Duration); ok {
		s.metrics.PluginMetrics().RecordExecution(duration)
	}

	if err != nil {
		s.metrics.PluginMetrics().IncExecutionErrors()
		return result, err
	}

	// Run post-execute hook
	if err := s.hookRegistry.Execute(ctx, plugin.HookPostExecute, result); err != nil {
		s.logger.Error("Post-execute hook failed", logger.F("error", err))
	}

	return result, nil
}

// HealthCheck performs a health check
func (s *SDK) HealthCheck(ctx context.Context) map[string]error {
	return s.pluginMgr.HealthCheck(ctx)
}

// GetPluginInfo returns info for all plugins
func (s *SDK) GetPluginInfo() map[string]plugin.PluginInfo {
	return s.pluginMgr.GetInfo()
}

// DefaultSDK is the default SDK instance
var defaultSDK *SDK
var defaultOnce sync.Once

// Default returns the default SDK instance
func Default() *SDK {
	defaultOnce.Do(func() {
		sdk, err := New(nil)
		if err != nil {
			panic(err)
		}
		defaultSDK = sdk
	})
	return defaultSDK
}

// QuickStart creates and starts a default SDK instance
func QuickStart() (*SDK, error) {
	sdk, err := New(nil)
	if err != nil {
		return nil, err
	}
	if err := sdk.Start(context.Background()); err != nil {
		return nil, err
	}
	return sdk, nil
}