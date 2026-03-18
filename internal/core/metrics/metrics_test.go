package metrics

import (
	"testing"
	"time"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("Registry is nil")
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()
	c := NewPluginMetricsCollector()
	r.Register(c)

	// Should not panic
}

func TestRegistryUnregister(t *testing.T) {
	r := NewRegistry()
	c := NewPluginMetricsCollector()
	r.Register(c)
	r.Unregister("plugins")

	// Should not panic
}

func TestRegistryCollect(t *testing.T) {
	r := NewRegistry()
	c := NewPluginMetricsCollector()
	r.Register(c)

	metrics := r.Collect()
	if len(metrics) == 0 {
		t.Error("Expected some metrics")
	}
}

func TestCounter(t *testing.T) {
	c := NewCounter("test_counter", "Test counter")

	c.Inc()
	if c.Value() != 1 {
		t.Errorf("Expected 1, got %f", c.Value())
	}

	c.Add(5)
	if c.Value() != 6 {
		t.Errorf("Expected 6, got %f", c.Value())
	}
}

func TestCounterMetric(t *testing.T) {
	c := NewCounter("test_counter", "Test counter")
	c.Inc()

	m := c.Metric()
	if m.Name != "test_counter" {
		t.Errorf("Expected name 'test_counter', got '%s'", m.Name)
	}
	if m.Type != MetricTypeCounter {
		t.Errorf("Expected type Counter, got %v", m.Type)
	}
	if m.Value != 1 {
		t.Errorf("Expected value 1, got %f", m.Value)
	}
}

func TestGauge(t *testing.T) {
	g := NewGauge("test_gauge", "Test gauge")

	g.Set(10)
	if g.Value() != 10 {
		t.Errorf("Expected 10, got %f", g.Value())
	}

	g.Inc()
	if g.Value() != 11 {
		t.Errorf("Expected 11, got %f", g.Value())
	}

	g.Dec()
	if g.Value() != 10 {
		t.Errorf("Expected 10, got %f", g.Value())
	}

	g.Add(5)
	if g.Value() != 15 {
		t.Errorf("Expected 15, got %f", g.Value())
	}

	g.Sub(5)
	if g.Value() != 10 {
		t.Errorf("Expected 10, got %f", g.Value())
	}
}

func TestGaugeMetric(t *testing.T) {
	g := NewGauge("test_gauge", "Test gauge")
	g.Set(42)

	m := g.Metric()
	if m.Name != "test_gauge" {
		t.Errorf("Expected name 'test_gauge', got '%s'", m.Name)
	}
	if m.Type != MetricTypeGauge {
		t.Errorf("Expected type Gauge, got %v", m.Type)
	}
	if m.Value != 42 {
		t.Errorf("Expected value 42, got %f", m.Value)
	}
}

func TestHistogram(t *testing.T) {
	h := NewHistogram("test_histogram", "Test histogram", []float64{0.1, 0.5, 1.0})

	h.Observe(0.05)
	h.Observe(0.3)
	h.Observe(0.7)
	h.Observe(2.0)

	metrics := h.Metrics()
	if len(metrics) == 0 {
		t.Error("Expected some metrics")
	}
}

func TestHistogramMetrics(t *testing.T) {
	h := NewHistogram("test_histogram", "Test histogram", []float64{0.1, 0.5, 1.0})

	h.Observe(0.5)
	h.Observe(1.5)

	metrics := h.Metrics()

	// Should have sum, count, and buckets
	foundSum := false
	foundCount := false
	for _, m := range metrics {
		if m.Name == "test_histogram_sum" {
			foundSum = true
		}
		if m.Name == "test_histogram_count" {
			foundCount = true
		}
	}

	if !foundSum {
		t.Error("Expected sum metric")
	}
	if !foundCount {
		t.Error("Expected count metric")
	}
}

func TestTimer(t *testing.T) {
	h := NewHistogram("test_timer", "Test timer", []float64{0.001, 0.01, 0.1})
	timer := NewTimer(h)

	time.Sleep(10 * time.Millisecond)
	timer.ObserveDuration()

	// Should have recorded something
	metrics := h.Metrics()
	foundSum := false
	for _, m := range metrics {
		if m.Name == "test_timer_sum" && m.Value > 0 {
			foundSum = true
		}
	}

	if !foundSum {
		t.Error("Timer should have recorded duration")
	}
}

func TestTimerDuration(t *testing.T) {
	h := NewHistogram("test_timer", "Test timer", []float64{0.1})
	timer := NewTimer(h)

	time.Sleep(20 * time.Millisecond)
	d := timer.Duration()

	if d < 20*time.Millisecond {
		t.Errorf("Duration should be at least 20ms, got %v", d)
	}
}

func TestPluginMetricsCollector(t *testing.T) {
	c := NewPluginMetricsCollector()

	// Test all methods
	c.IncPlugins()
	c.SetRunningPlugins(5)
	c.IncErrors()
	c.IncExecutions()
	c.IncExecutionErrors()
	c.RecordExecution(100 * time.Millisecond)

	metrics := c.Collect()
	if len(metrics) == 0 {
		t.Error("Expected some metrics")
	}

	// Check specific metrics exist
	names := make(map[string]bool)
	for _, m := range metrics {
		names[m.Name] = true
	}

	if !names["plugins_total"] {
		t.Error("Expected plugins_total metric")
	}
	if !names["plugins_running"] {
		t.Error("Expected plugins_running metric")
	}
	if !names["plugins_errors_total"] {
		t.Error("Expected plugins_errors_total metric")
	}
}

func TestPluginMetricsCollectorName(t *testing.T) {
	c := NewPluginMetricsCollector()

	if c.Name() != "plugins" {
		t.Errorf("Expected name 'plugins', got '%s'", c.Name())
	}
}

func TestNewMetricsProvider(t *testing.T) {
	p := NewMetricsProvider()
	if p == nil {
		t.Fatal("MetricsProvider is nil")
	}

	if p.Registry() == nil {
		t.Error("Registry should not be nil")
	}

	if p.PluginMetrics() == nil {
		t.Error("PluginMetrics should not be nil")
	}
}

func TestMetricsProviderCollect(t *testing.T) {
	p := NewMetricsProvider()
	p.PluginMetrics().IncPlugins()
	p.PluginMetrics().IncExecutions()

	metrics := p.Collect()
	if len(metrics) == 0 {
		t.Error("Expected some metrics")
	}
}

func TestDefaultMetricsProvider(t *testing.T) {
	if DefaultMetricsProvider == nil {
		t.Fatal("DefaultMetricsProvider is nil")
	}
}

func TestPluginMetricsFunc(t *testing.T) {
	c := PluginMetrics()
	if c == nil {
		t.Fatal("PluginMetrics() returned nil")
	}
}

func TestMetricTypeString(t *testing.T) {
	// MetricType values should be correct
	if MetricTypeCounter != 0 {
		t.Error("Counter should be 0")
	}
	if MetricTypeGauge != 1 {
		t.Error("Gauge should be 1")
	}
	if MetricTypeHistogram != 2 {
		t.Error("Histogram should be 2")
	}
	if MetricTypeSummary != 3 {
		t.Error("Summary should be 3")
	}
}