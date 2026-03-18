package metrics

import (
	"context"
	"sync"
	"time"
)

// MetricType represents the type of metric
type MetricType int

const (
	MetricTypeCounter   MetricType = iota
	MetricTypeGauge
	MetricTypeHistogram
	MetricTypeSummary
)

// Metric represents a single metric
type Metric struct {
	Name        string
	Type        MetricType
	Value       float64
	Labels      map[string]string
	Description string
	Timestamp   time.Time
}

// Collector collects metrics
type Collector interface {
	Collect() []Metric
	Name() string
}

// Registry holds all metrics collectors
type Registry struct {
	mu         sync.RWMutex
	collectors map[string]Collector
	metrics    map[string]*Metric
}

// NewRegistry creates a new metrics registry
func NewRegistry() *Registry {
	return &Registry{
		collectors: make(map[string]Collector),
		metrics:    make(map[string]*Metric),
	}
}

// Register registers a collector
func (r *Registry) Register(c Collector) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.collectors[c.Name()] = c
}

// Unregister unregisters a collector
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.collectors, name)
}

// Collect collects all metrics
func (r *Registry) Collect() []Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var metrics []Metric
	for _, c := range r.collectors {
		metrics = append(metrics, c.Collect()...)
	}
	return metrics
}

// Counter is a monotonically increasing counter
type Counter struct {
	mu    sync.RWMutex
	name  string
	value float64
	desc  string
}

// NewCounter creates a new counter
func NewCounter(name, desc string) *Counter {
	return &Counter{
		name: name,
		desc: desc,
	}
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Add adds a value to the counter
func (c *Counter) Add(v float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += v
}

// Value returns the current value
func (c *Counter) Value() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Metric returns the counter as a metric
func (c *Counter) Metric() Metric {
	return Metric{
		Name:        c.name,
		Type:        MetricTypeCounter,
		Value:       c.Value(),
		Description: c.desc,
		Timestamp:   time.Now(),
	}
}

// Gauge is a value that can go up and down
type Gauge struct {
	mu    sync.RWMutex
	name  string
	value float64
	desc  string
}

// NewGauge creates a new gauge
func NewGauge(name, desc string) *Gauge {
	return &Gauge{
		name: name,
		desc: desc,
	}
}

// Set sets the gauge value
func (g *Gauge) Set(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = v
}

// Inc increments the gauge by 1
func (g *Gauge) Inc() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value++
}

// Dec decrements the gauge by 1
func (g *Gauge) Dec() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value--
}

// Add adds a value to the gauge
func (g *Gauge) Add(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value += v
}

// Sub subtracts a value from the gauge
func (g *Gauge) Sub(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value -= v
}

// Value returns the current value
func (g *Gauge) Value() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

// Metric returns the gauge as a metric
func (g *Gauge) Metric() Metric {
	return Metric{
		Name:        g.name,
		Type:        MetricTypeGauge,
		Value:       g.Value(),
		Description: g.desc,
		Timestamp:   time.Now(),
	}
}

// Histogram tracks distribution of values
type Histogram struct {
	mu     sync.RWMutex
	name   string
	desc   string
	bounds []float64
	counts []uint64
	sum    float64
	count  uint64
}

// NewHistogram creates a new histogram
func NewHistogram(name, desc string, bounds []float64) *Histogram {
	return &Histogram{
		name:   name,
		desc:   desc,
		bounds: bounds,
		counts: make([]uint64, len(bounds)+1),
	}
}

// Observe records an observation
func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sum += v
	h.count++

	for i, bound := range h.bounds {
		if v <= bound {
			h.counts[i]++
			return
		}
	}
	h.counts[len(h.bounds)]++
}

// Metrics returns histogram metrics
func (h *Histogram) Metrics() []Metric {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var metrics []Metric

	// Sum
	metrics = append(metrics, Metric{
		Name:        h.name + "_sum",
		Type:        MetricTypeGauge,
		Value:       h.sum,
		Description: h.desc + " (sum)",
		Timestamp:   time.Now(),
	})

	// Count
	metrics = append(metrics, Metric{
		Name:        h.name + "_count",
		Type:        MetricTypeGauge,
		Value:       float64(h.count),
		Description: h.desc + " (count)",
		Timestamp:   time.Now(),
	})

	// Buckets
	for i, bound := range h.bounds {
		metrics = append(metrics, Metric{
			Name:        h.name + "_bucket",
			Type:        MetricTypeCounter,
			Value:       float64(h.counts[i]),
			Description: h.desc + " (bucket)",
			Labels:      map[string]string{"le": formatFloat(bound)},
			Timestamp:   time.Now(),
		})
	}

	return metrics
}

// Timer measures duration
type Timer struct {
	histogram *Histogram
	start     time.Time
}

// NewTimer creates a new timer
func NewTimer(h *Histogram) *Timer {
	return &Timer{
		histogram: h,
		start:     time.Now(),
	}
}

// ObserveDuration records the duration in seconds
func (t *Timer) ObserveDuration() {
	d := time.Since(t.start).Seconds()
	t.histogram.Observe(d)
}

// Duration returns the elapsed duration
func (t *Timer) Duration() time.Duration {
	return time.Since(t.start)
}

// PluginMetricsCollector collects plugin metrics
type PluginMetricsCollector struct {
	mu              sync.RWMutex
	pluginsTotal    *Counter
	pluginsRunning  *Gauge
	pluginsError    *Counter
	executionsTotal *Counter
	executionsError *Counter
	executionTime   *Histogram
}

// NewPluginMetricsCollector creates a new plugin metrics collector
func NewPluginMetricsCollector() *PluginMetricsCollector {
	return &PluginMetricsCollector{
		pluginsTotal:    NewCounter("plugins_total", "Total number of registered plugins"),
		pluginsRunning:  NewGauge("plugins_running", "Number of currently running plugins"),
		pluginsError:    NewCounter("plugins_errors_total", "Total number of plugin errors"),
		executionsTotal: NewCounter("plugin_executions_total", "Total number of plugin executions"),
		executionsError: NewCounter("plugin_executions_errors_total", "Total number of failed executions"),
		executionTime:   NewHistogram("plugin_execution_duration_seconds", "Plugin execution duration", []float64{0.001, 0.01, 0.1, 1, 10, 60}),
	}
}

// Name returns the collector name
func (c *PluginMetricsCollector) Name() string {
	return "plugins"
}

// Collect collects plugin metrics
func (c *PluginMetricsCollector) Collect() []Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var metrics []Metric
	metrics = append(metrics, c.pluginsTotal.Metric())
	metrics = append(metrics, c.pluginsRunning.Metric())
	metrics = append(metrics, c.pluginsError.Metric())
	metrics = append(metrics, c.executionsTotal.Metric())
	metrics = append(metrics, c.executionsError.Metric())
	metrics = append(metrics, c.executionTime.Metrics()...)
	return metrics
}

// IncPlugins increments total plugins
func (c *PluginMetricsCollector) IncPlugins() {
	c.pluginsTotal.Inc()
}

// SetRunningPlugins sets the running plugins count
func (c *PluginMetricsCollector) SetRunningPlugins(count int) {
	c.pluginsRunning.Set(float64(count))
}

// IncErrors increments error count
func (c *PluginMetricsCollector) IncErrors() {
	c.pluginsError.Inc()
}

// IncExecutions increments execution count
func (c *PluginMetricsCollector) IncExecutions() {
	c.executionsTotal.Inc()
}

// IncExecutionErrors increments execution error count
func (c *PluginMetricsCollector) IncExecutionErrors() {
	c.executionsError.Inc()
}

// RecordExecution records an execution duration
func (c *PluginMetricsCollector) RecordExecution(d time.Duration) {
	c.executionTime.Observe(d.Seconds())
}

// MetricsProvider provides metrics for the application
type MetricsProvider struct {
	registry *Registry
	plugins  *PluginMetricsCollector
}

// NewMetricsProvider creates a new metrics provider
func NewMetricsProvider() *MetricsProvider {
	registry := NewRegistry()
	plugins := NewPluginMetricsCollector()
	registry.Register(plugins)
	return &MetricsProvider{
		registry: registry,
		plugins:  plugins,
	}
}

// PluginMetrics returns the plugin metrics collector
func (p *MetricsProvider) PluginMetrics() *PluginMetricsCollector {
	return p.plugins
}

// Registry returns the metrics registry
func (p *MetricsProvider) Registry() *Registry {
	return p.registry
}

// Collect returns all metrics
func (p *MetricsProvider) Collect() []Metric {
	return p.registry.Collect()
}

// DefaultMetricsProvider is the default metrics provider
var DefaultMetricsProvider = NewMetricsProvider()

// PluginMetrics returns the default plugin metrics collector
func PluginMetrics() *PluginMetricsCollector {
	return DefaultMetricsProvider.PluginMetrics()
}

func formatFloat(v float64) string {
	if v == float64(int64(v)) {
		return string(rune(int64(v)))
	}
	return string(rune(int64(v * 100)))
}

// Metrics context key
type metricsCtxKey struct{}

// WithMetrics returns a context with metrics provider
func WithMetrics(ctx context.Context, provider *MetricsProvider) context.Context {
	return context.WithValue(ctx, metricsCtxKey{}, provider)
}

// FromContext retrieves the metrics provider from context
func FromContext(ctx context.Context) *MetricsProvider {
	if provider, ok := ctx.Value(metricsCtxKey{}).(*MetricsProvider); ok {
		return provider
	}
	return DefaultMetricsProvider
}