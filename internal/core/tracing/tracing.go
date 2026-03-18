package tracing

import (
	"context"
	"encoding/hex"
	"sync"
	"time"
)

// TraceID represents a trace identifier
type TraceID [16]byte

// String returns the hex representation of the trace ID
func (t TraceID) String() string {
	return hex.EncodeToString(t[:])
}

// IsValid returns true if the trace ID is valid
func (t TraceID) IsValid() bool {
	for _, b := range t {
		if b != 0 {
			return true
		}
	}
	return false
}

// SpanID represents a span identifier
type SpanID [8]byte

// String returns the hex representation of the span ID
func (s SpanID) String() string {
	return hex.EncodeToString(s[:])
}

// IsValid returns true if the span ID is valid
func (s SpanID) IsValid() bool {
	for _, b := range s {
		if b != 0 {
			return true
		}
	}
	return false
}

// SpanKind represents the kind of span
type SpanKind int

const (
	SpanKindUnspecified SpanKind = iota
	SpanKindInternal
	SpanKindServer
	SpanKindClient
	SpanKindProducer
	SpanKindConsumer
)

func (k SpanKind) String() string {
	switch k {
	case SpanKindInternal:
		return "internal"
	case SpanKindServer:
		return "server"
	case SpanKindClient:
		return "client"
	case SpanKindProducer:
		return "producer"
	case SpanKindConsumer:
		return "consumer"
	default:
		return "unspecified"
	}
}

// StatusCode represents the status of a span
type StatusCode int

const (
	StatusCodeUnset StatusCode = iota
	StatusCodeOK
	StatusCodeError
)

func (s StatusCode) String() string {
	switch s {
	case StatusCodeOK:
		return "ok"
	case StatusCodeError:
		return "error"
	default:
		return "unset"
	}
}

// Attribute represents a span attribute
type Attribute struct {
	Key   string
	Value interface{}
}

// String creates a string attribute
func String(key, value string) Attribute {
	return Attribute{Key: key, Value: value}
}

// Int creates an int attribute
func Int(key string, value int) Attribute {
	return Attribute{Key: key, Value: value}
}

// Int64 creates an int64 attribute
func Int64(key string, value int64) Attribute {
	return Attribute{Key: key, Value: value}
}

// Float64 creates a float64 attribute
func Float64(key string, value float64) Attribute {
	return Attribute{Key: key, Value: value}
}

// Bool creates a bool attribute
func Bool(key string, value bool) Attribute {
	return Attribute{Key: key, Value: value}
}

// Event represents a span event
type Event struct {
	Name       string
	Attributes []Attribute
	Timestamp  time.Time
}

// Span represents a unit of work in a trace
type Span struct {
	mu          sync.Mutex
	traceID     TraceID
	spanID      SpanID
	parentID    SpanID
	name        string
	kind        SpanKind
	startTime   time.Time
	endTime     time.Time
	status      StatusCode
	statusMsg   string
	attributes  []Attribute
	events      []Event
	recorded    bool
	sampled     bool
	tracer      *Tracer
}

// TraceID returns the trace ID
func (s *Span) TraceID() TraceID {
	return s.traceID
}

// SpanID returns the span ID
func (s *Span) SpanID() SpanID {
	return s.spanID
}

// ParentID returns the parent span ID
func (s *Span) ParentID() SpanID {
	return s.parentID
}

// Name returns the span name
func (s *Span) Name() string {
	return s.name
}

// Kind returns the span kind
func (s *Span) Kind() SpanKind {
	return s.kind
}

// StartTime returns the span start time
func (s *Span) StartTime() time.Time {
	return s.startTime
}

// EndTime returns the span end time
func (s *Span) EndTime() time.Time {
	return s.endTime
}

// Duration returns the span duration
func (s *Span) Duration() time.Duration {
	if s.endTime.IsZero() {
		return 0
	}
	return s.endTime.Sub(s.startTime)
}

// Status returns the span status
func (s *Span) Status() StatusCode {
	return s.status
}

// SetStatus sets the span status
func (s *Span) SetStatus(code StatusCode, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = code
	s.statusMsg = msg
}

// SetAttributes sets span attributes
func (s *Span) SetAttributes(attrs ...Attribute) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attributes = append(s.attributes, attrs...)
}

// AddEvent adds an event to the span
func (s *Span) AddEvent(name string, attrs ...Attribute) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, Event{
		Name:       name,
		Attributes: attrs,
		Timestamp:  time.Now(),
	})
}

// End ends the span
func (s *Span) End(options ...SpanEndOption) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.endTime.IsZero() {
		return // Already ended
	}

	for _, opt := range options {
		opt(s)
	}

	if s.endTime.IsZero() {
		s.endTime = time.Now()
	}

	// Export span if sampled
	if s.sampled && s.tracer != nil {
		s.tracer.export(s)
	}
}

// IsRecording returns true if the span is recording
func (s *Span) IsRecording() bool {
	return s.recorded
}

// SpanEndOption configures span end behavior
type SpanEndOption func(*Span)

// WithEndTime sets the end time
func WithEndTime(t time.Time) SpanEndOption {
	return func(s *Span) {
		s.endTime = t
	}
}

// WithStatus sets the status at end
func WithStatus(code StatusCode, msg string) SpanEndOption {
	return func(s *Span) {
		s.status = code
		s.statusMsg = msg
	}
}

// SpanContext contains trace context information
type SpanContext struct {
	traceID    TraceID
	spanID     SpanID
	sampled    bool
	traceFlags byte
}

// TraceID returns the trace ID
func (c SpanContext) TraceID() TraceID {
	return c.traceID
}

// SpanID returns the span ID
func (c SpanContext) SpanID() SpanID {
	return c.spanID
}

// IsSampled returns true if the span should be sampled
func (c SpanContext) IsSampled() bool {
	return c.sampled
}

// IsValid returns true if the context is valid
func (c SpanContext) IsValid() bool {
	return c.traceID.IsValid() && c.spanID.IsValid()
}

// SpanContextFromSpan extracts span context from a span
func SpanContextFromSpan(s *Span) SpanContext {
	return SpanContext{
		traceID: s.traceID,
		spanID:  s.spanID,
		sampled: s.sampled,
	}
}

// Context keys
type ctxKey struct{}

var spanKey = ctxKey{}

// ContextWithSpan returns a context with the span attached
func ContextWithSpan(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, spanKey, span)
}

// SpanFromContext extracts a span from context
func SpanFromContext(ctx context.Context) *Span {
	if span, ok := ctx.Value(spanKey).(*Span); ok {
		return span
	}
	return nil
}

// SpanContextFromContext extracts span context from context
func SpanContextFromContext(ctx context.Context) SpanContext {
	span := SpanFromContext(ctx)
	if span == nil {
		return SpanContext{}
	}
	return SpanContextFromSpan(span)
}

// Sampler decides whether a span should be sampled
type Sampler interface {
	ShouldSample(params SamplingParameters) SamplingResult
}

// SamplingParameters contains sampling parameters
type SamplingParameters struct {
	TraceID   TraceID
	Name      string
	Kind      SpanKind
	Attributes []Attribute
}

// SamplingDecision represents the sampling decision
type SamplingDecision int

const (
	SamplingDecisionDrop SamplingDecision = iota
	SamplingDecisionRecordAndSample
)

// SamplingResult represents the sampling result
type SamplingResult struct {
	Decision   SamplingDecision
	Attributes []Attribute
}

// AlwaysSample samples all spans
type AlwaysSample struct{}

func (s *AlwaysSample) ShouldSample(params SamplingParameters) SamplingResult {
	return SamplingResult{Decision: SamplingDecisionRecordAndSample}
}

// NeverSample samples no spans
type NeverSample struct{}

func (s *NeverSample) ShouldSample(params SamplingParameters) SamplingResult {
	return SamplingResult{Decision: SamplingDecisionDrop}
}

// TraceIDRatioBased samples based on a ratio
type TraceIDRatioBased struct {
	ratio float64
}

// NewTraceIDRatioBased creates a ratio-based sampler
func NewTraceIDRatioBased(ratio float64) *TraceIDRatioBased {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return &TraceIDRatioBased{ratio: ratio}
}

func (s *TraceIDRatioBased) ShouldSample(params SamplingParameters) SamplingResult {
	// Simple hash-based sampling
	h := uint64(0)
	for _, b := range params.TraceID[:] {
		h = h*31 + uint64(b)
	}
	if float64(h%10000)/10000 < s.ratio {
		return SamplingResult{Decision: SamplingDecisionRecordAndSample}
	}
	return SamplingResult{Decision: SamplingDecisionDrop}
}

// Exporter exports completed spans
type Exporter interface {
	ExportSpans(spans []*Span) error
	Shutdown(ctx context.Context) error
}

// Processor processes spans before export
type Processor interface {
	OnStart(span *Span)
	OnEnd(span *Span)
	Shutdown(ctx context.Context) error
}

// BatchProcessor batches spans for export
type BatchProcessor struct {
	mu        sync.Mutex
	spans     []*Span
	exporter  Exporter
	batchSize int
	interval  time.Duration
	timer     *time.Timer
	stopCh    chan struct{}
}

// BatchProcessorOption configures the batch processor
type BatchProcessorOption func(*BatchProcessor)

// WithBatchSize sets the batch size
func WithBatchSize(size int) BatchProcessorOption {
	return func(p *BatchProcessor) { p.batchSize = size }
}

// WithBatchInterval sets the export interval
func WithBatchInterval(d time.Duration) BatchProcessorOption {
	return func(p *BatchProcessor) { p.interval = d }
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(exporter Exporter, opts ...BatchProcessorOption) *BatchProcessor {
	p := &BatchProcessor{
		exporter:  exporter,
		batchSize: 512,
		interval:  5 * time.Second,
		spans:     make([]*Span, 0),
		stopCh:    make(chan struct{}),
	}
	for _, opt := range opts {
		opt(p)
	}
	go p.run()
	return p
}

func (p *BatchProcessor) run() {
	p.timer = time.NewTimer(p.interval)
	for {
		select {
		case <-p.stopCh:
			return
		case <-p.timer.C:
			p.flush()
			p.timer.Reset(p.interval)
		}
	}
}

func (p *BatchProcessor) flush() {
	p.mu.Lock()
	spans := p.spans
	p.spans = make([]*Span, 0)
	p.mu.Unlock()

	if len(spans) > 0 && p.exporter != nil {
		p.exporter.ExportSpans(spans)
	}
}

func (p *BatchProcessor) OnStart(span *Span) {}

func (p *BatchProcessor) OnEnd(span *Span) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.spans = append(p.spans, span)
	if len(p.spans) >= p.batchSize {
		go p.flush()
	}
}

func (p *BatchProcessor) Shutdown(ctx context.Context) error {
	close(p.stopCh)
	p.flush()
	if p.exporter != nil {
		return p.exporter.Shutdown(ctx)
	}
	return nil
}