package tracing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Tracer creates and manages spans
type Tracer struct {
	mu        sync.RWMutex
	name      string
	version   string
	sampler   Sampler
	processor Processor
	resource  []Attribute
}

// TracerOption configures a tracer
type TracerOption func(*Tracer)

// WithSampler sets the sampler
func WithSampler(sampler Sampler) TracerOption {
	return func(t *Tracer) { t.sampler = sampler }
}

// WithProcessor sets the processor
func WithProcessor(processor Processor) TracerOption {
	return func(t *Tracer) { t.processor = processor }
}

// WithVersion sets the tracer version
func WithVersion(version string) TracerOption {
	return func(t *Tracer) { t.version = version }
}

// WithResource sets resource attributes
func WithResource(attrs ...Attribute) TracerOption {
	return func(t *Tracer) { t.resource = attrs }
}

// NewTracer creates a new tracer
func NewTracer(name string, opts ...TracerOption) *Tracer {
	t := &Tracer{
		name:    name,
		sampler: &AlwaysSample{},
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Name returns the tracer name
func (t *Tracer) Name() string {
	return t.name
}

// Version returns the tracer version
func (t *Tracer) Version() string {
	return t.version
}

// Start starts a new span
func (t *Tracer) Start(ctx context.Context, name string, opts ...SpanStartOption) (context.Context, *Span) {
	// Generate IDs
	traceID := generateTraceID()
	spanID := generateSpanID()

	// Check for parent context
	var parentID SpanID
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		traceID = parentSpan.traceID
		parentID = parentSpan.spanID
	} else if sc := SpanContextFromContext(ctx); sc.IsValid() {
		traceID = sc.traceID
		parentID = sc.spanID
	}

	// Create span
	span := &Span{
		traceID:   traceID,
		spanID:    spanID,
		parentID:  parentID,
		name:      name,
		startTime: time.Now(),
		recorded:  true,
		tracer:    t,
	}

	// Apply options
	for _, opt := range opts {
		opt(span)
	}

	// Sample
	samplingResult := t.sampler.ShouldSample(SamplingParameters{
		TraceID:    traceID,
		Name:       name,
		Kind:       span.kind,
		Attributes: span.attributes,
	})

	span.sampled = samplingResult.Decision == SamplingDecisionRecordAndSample
	if len(samplingResult.Attributes) > 0 {
		span.attributes = append(span.attributes, samplingResult.Attributes...)
	}

	// Add resource attributes
	span.attributes = append(span.attributes, t.resource...)

	// Notify processor
	if t.processor != nil {
		t.processor.OnStart(span)
	}

	return ContextWithSpan(ctx, span), span
}

// export exports a completed span
func (t *Tracer) export(span *Span) {
	if t.processor != nil {
		t.processor.OnEnd(span)
	}
}

// SpanStartOption configures span start behavior
type SpanStartOption func(*Span)

// WithSpanKind sets the span kind
func WithSpanKind(kind SpanKind) SpanStartOption {
	return func(s *Span) { s.kind = kind }
}

// WithAttributes sets initial attributes
func WithAttributes(attrs ...Attribute) SpanStartOption {
	return func(s *Span) { s.attributes = attrs }
}

// WithStartTime sets the start time
func WithStartTime(t time.Time) SpanStartOption {
	return func(s *Span) { s.startTime = t }
}

// WithLinks adds links to other spans
func WithLinks(links ...SpanContext) SpanStartOption {
	return func(s *Span) {
		// Links are stored as attributes for simplicity
		for _, link := range links {
			s.attributes = append(s.attributes,
				String("link.trace_id", link.traceID.String()),
				String("link.span_id", link.spanID.String()),
			)
		}
	}
}

// WithNewRoot starts a new root span
func WithNewRoot() SpanStartOption {
	return func(s *Span) {
		s.traceID = generateTraceID()
		s.parentID = SpanID{}
	}
}

// TracerProvider creates and manages tracers
type TracerProvider struct {
	mu        sync.RWMutex
	tracers   map[string]*Tracer
	sampler   Sampler
	processor Processor
	resource  []Attribute
	version   string
}

// TracerProviderOption configures a tracer provider
type TracerProviderOption func(*TracerProvider)

// NewTracerProvider creates a new tracer provider
func NewTracerProvider(opts ...TracerProviderOption) *TracerProvider {
	p := &TracerProvider{
		tracers: make(map[string]*Tracer),
		sampler: &AlwaysSample{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Tracer creates or retrieves a tracer
func (p *TracerProvider) Tracer(name string, opts ...TracerOption) *Tracer {
	p.mu.RLock()
	if t, ok := p.tracers[name]; ok {
		p.mu.RUnlock()
		return t
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double check
	if t, ok := p.tracers[name]; ok {
		return t
	}

	// Create new tracer with provider defaults
	tracerOpts := []TracerOption{
		WithSampler(p.sampler),
	}
	if p.processor != nil {
		tracerOpts = append(tracerOpts, WithProcessor(p.processor))
	}
	if p.version != "" {
		tracerOpts = append(tracerOpts, WithVersion(p.version))
	}
	if len(p.resource) > 0 {
		tracerOpts = append(tracerOpts, WithResource(p.resource...))
	}
	tracerOpts = append(tracerOpts, opts...)

	t := NewTracer(name, tracerOpts...)
	p.tracers[name] = t
	return t
}

// SetSampler sets the default sampler
func (p *TracerProvider) SetSampler(sampler Sampler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sampler = sampler
}

// SetProcessor sets the processor
func (p *TracerProvider) SetProcessor(processor Processor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processor = processor
}

// Shutdown shuts down all tracers
func (p *TracerProvider) Shutdown(ctx context.Context) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.processor != nil {
		return p.processor.Shutdown(ctx)
	}
	return nil
}

// Global tracer provider
var globalProvider = NewTracerProvider()
var globalMu sync.RWMutex

// SetGlobalTracerProvider sets the global tracer provider
func SetGlobalTracerProvider(p *TracerProvider) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalProvider = p
}

// GetGlobalTracerProvider returns the global tracer provider
func GetGlobalTracerProvider() *TracerProvider {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalProvider
}

// GetTracer returns a tracer from the global provider
func GetTracer(name string) *Tracer {
	return GetGlobalTracerProvider().Tracer(name)
}

// Start starts a span using the global tracer provider
func Start(ctx context.Context, name string, opts ...SpanStartOption) (context.Context, *Span) {
	return GetTracer("default").Start(ctx, name, opts...)
}

// ID generation helpers
func generateTraceID() TraceID {
	var id TraceID
	_, _ = rand.Read(id[:])
	return id
}

func generateSpanID() SpanID {
	var id SpanID
	_, _ = rand.Read(id[:])
	return id
}

// Propagation formats
type Propagator interface {
	Inject(ctx context.Context, carrier TextMapCarrier)
	Extract(ctx context.Context, carrier TextMapCarrier) context.Context
}

// TextMapCarrier stores trace context as text
type TextMapCarrier interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}

// MapCarrier is a simple map-based carrier
type MapCarrier map[string]string

func (c MapCarrier) Get(key string) string    { return c[key] }
func (c MapCarrier) Set(key, value string)    { c[key] = value }
func (c MapCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// W3CTraceContextPropagator implements W3C trace context propagation
type W3CTraceContextPropagator struct{}

const (
	traceparentHeader = "traceparent"
	tracestateHeader  = "tracestate"
)

func (p *W3CTraceContextPropagator) Inject(ctx context.Context, carrier TextMapCarrier) {
	span := SpanFromContext(ctx)
	if span == nil {
		return
	}

	// Format: version-traceid-parentid-flags
	// version: 00
	// traceid: 32 hex chars
	// parentid: 16 hex chars
	// flags: 01 if sampled
	flags := "00"
	if span.sampled {
		flags = "01"
	}

	traceparent := fmt.Sprintf("00-%s-%s-%s",
		span.traceID.String(),
		span.spanID.String(),
		flags,
	)
	carrier.Set(traceparentHeader, traceparent)
}

func (p *W3CTraceContextPropagator) Extract(ctx context.Context, carrier TextMapCarrier) context.Context {
	traceparent := carrier.Get(traceparentHeader)
	if traceparent == "" {
		return ctx
	}

	// Parse traceparent
	// Format: version-traceid-parentid-flags
	var version, traceIDHex, spanIDHex, flags string
	_, err := fmt.Sscanf(traceparent, "%2s-%32s-%16s-%2s", &version, &traceIDHex, &spanIDHex, &flags)
	if err != nil {
		return ctx
	}

	if version != "00" {
		return ctx
	}

	traceID, err := parseTraceID(traceIDHex)
	if err != nil {
		return ctx
	}

	spanID, err := parseSpanID(spanIDHex)
	if err != nil {
		return ctx
	}

	sampled := flags == "01"

	sc := SpanContext{
		traceID: traceID,
		spanID:  spanID,
		sampled: sampled,
	}

	// Store in context
	return context.WithValue(ctx, spanKey, &Span{
		traceID: sc.traceID,
		spanID:  sc.spanID,
		sampled: sc.sampled,
	})
}

func parseTraceID(hexStr string) (TraceID, error) {
	var id TraceID
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return id, err
	}
	if len(decoded) != 16 {
		return id, fmt.Errorf("invalid trace ID length: %d", len(decoded))
	}
	copy(id[:], decoded)
	return id, nil
}

func parseSpanID(hexStr string) (SpanID, error) {
	var id SpanID
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return id, err
	}
	if len(decoded) != 8 {
		return id, fmt.Errorf("invalid span ID length: %d", len(decoded))
	}
	copy(id[:], decoded)
	return id, nil
}