package tracing

import (
	"context"
	"encoding/hex"
	"sync"
	"testing"
	"time"
)

func TestTraceIDString(t *testing.T) {
	id := TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	expected := "0102030405060708090a0b0c0d0e0f10"
	if id.String() != expected {
		t.Errorf("Expected %s, got %s", expected, id.String())
	}
}

func TestTraceIDIsValid(t *testing.T) {
	// Zero trace ID
	zeroID := TraceID{}
	if zeroID.IsValid() {
		t.Error("Zero TraceID should not be valid")
	}

	// Non-zero trace ID
	id := TraceID{1}
	if !id.IsValid() {
		t.Error("Non-zero TraceID should be valid")
	}
}

func TestSpanIDString(t *testing.T) {
	id := SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	expected := "0102030405060708"
	if id.String() != expected {
		t.Errorf("Expected %s, got %s", expected, id.String())
	}
}

func TestSpanIDIsValid(t *testing.T) {
	// Zero span ID
	zeroID := SpanID{}
	if zeroID.IsValid() {
		t.Error("Zero SpanID should not be valid")
	}

	// Non-zero span ID
	id := SpanID{1}
	if !id.IsValid() {
		t.Error("Non-zero SpanID should be valid")
	}
}

func TestSpanKindString(t *testing.T) {
	tests := []struct {
		kind     SpanKind
		expected string
	}{
		{SpanKindUnspecified, "unspecified"},
		{SpanKindInternal, "internal"},
		{SpanKindServer, "server"},
		{SpanKindClient, "client"},
		{SpanKindProducer, "producer"},
		{SpanKindConsumer, "consumer"},
		{SpanKind(99), "unspecified"},
	}

	for _, test := range tests {
		if test.kind.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.kind.String())
		}
	}
}

func TestStatusCodeString(t *testing.T) {
	tests := []struct {
		code     StatusCode
		expected string
	}{
		{StatusCodeUnset, "unset"},
		{StatusCodeOK, "ok"},
		{StatusCodeError, "error"},
		{StatusCode(99), "unset"},
	}

	for _, test := range tests {
		if test.code.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.code.String())
		}
	}
}

func TestAttributeCreators(t *testing.T) {
	tests := []struct {
		attr     Attribute
		key      string
		expected interface{}
	}{
		{String("key", "value"), "key", "value"},
		{Int("key", 42), "key", 42},
		{Int64("key", int64(42)), "key", int64(42)},
		{Float64("key", 3.14), "key", 3.14},
		{Bool("key", true), "key", true},
	}

	for _, test := range tests {
		if test.attr.Key != test.key {
			t.Errorf("Expected key %s, got %s", test.key, test.attr.Key)
		}
		if test.attr.Value != test.expected {
			t.Errorf("Expected value %v, got %v", test.expected, test.attr.Value)
		}
	}
}

func TestSpanBasicOperations(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	// Start span
	ctx, span := tracer.Start(ctx, "test-span")

	if span.Name() != "test-span" {
		t.Errorf("Expected name 'test-span', got %s", span.Name())
	}

	if !span.TraceID().IsValid() {
		t.Error("TraceID should be valid")
	}

	if !span.SpanID().IsValid() {
		t.Error("SpanID should be valid")
	}

	if span.ParentID().IsValid() {
		t.Error("Root span should not have parent")
	}

	// Add a small delay to ensure duration is measurable
	time.Sleep(time.Millisecond)

	// End span
	span.End()

	// Duration should be positive after End()
	duration := span.Duration()
	if duration <= 0 {
		t.Errorf("Duration should be positive after End(), got %v", duration)
	}
}

func TestSpanWithParent(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	// Create parent span
	ctx, parent := tracer.Start(ctx, "parent")
	parentID := parent.SpanID()
	traceID := parent.TraceID()

	// Create child span
	ctx, child := tracer.Start(ctx, "child")

	if child.TraceID() != traceID {
		t.Error("Child should have same TraceID as parent")
	}

	if child.ParentID() != parentID {
		t.Error("Child should have parent's SpanID as ParentID")
	}
}

func TestSpanAttributes(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test",
		WithAttributes(String("initial", "value")),
	)

	span.SetAttributes(
		String("key1", "value1"),
		Int("key2", 42),
	)

	// End to trigger export
	span.End()
}

func TestSpanEvents(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	span.AddEvent("event1", String("attr", "value"))
	span.End()
}

func TestSpanStatus(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")

	if span.Status() != StatusCodeUnset {
		t.Error("Initial status should be Unset")
	}

	span.SetStatus(StatusCodeOK, "all good")
	if span.Status() != StatusCodeOK {
		t.Error("Status should be OK")
	}

	span.End()
}

func TestSpanEndOptions(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")

	endTime := time.Now().Add(-time.Hour)
	span.End(WithEndTime(endTime))

	if !span.EndTime().Equal(endTime) {
		t.Error("EndTime should match option")
	}
}

func TestSpanEndWithStatus(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	span.End(WithStatus(StatusCodeError, "something went wrong"))

	if span.Status() != StatusCodeError {
		t.Error("Status should be Error")
	}
}

func TestSpanDoubleEnd(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	span.End()

	// Second end should be a no-op
	endTime := span.EndTime()
	span.End()

	if !span.EndTime().Equal(endTime) {
		t.Error("Second End() should not change EndTime")
	}
}

func TestSpanIsRecording(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")

	if !span.IsRecording() {
		t.Error("Span should be recording")
	}

	span.End()
}

func TestSpanKindOption(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test", WithSpanKind(SpanKindServer))

	if span.Kind() != SpanKindServer {
		t.Errorf("Expected SpanKindServer, got %v", span.Kind())
	}
}

func TestSpanStartTime(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	startTime := time.Now().Add(-time.Hour)
	_, span := tracer.Start(ctx, "test", WithStartTime(startTime))

	if !span.StartTime().Equal(startTime) {
		t.Error("StartTime should match option")
	}
}

func TestSpanWithLinks(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	// Create a span to link to
	_, other := tracer.Start(ctx, "other")
	otherSc := SpanContextFromSpan(other)

	// Create span with link
	_, span := tracer.Start(ctx, "test", WithLinks(otherSc))
	span.End()
}

func TestSpanWithNewRoot(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	// Create parent span
	_, parent := tracer.Start(ctx, "parent")
	parentTraceID := parent.TraceID()

	// Create new root span (ignores parent)
	_, root := tracer.Start(ctx, "root", WithNewRoot())

	if root.TraceID() == parentTraceID {
		t.Error("NewRoot span should have different TraceID")
	}

	if root.ParentID().IsValid() {
		t.Error("NewRoot span should not have parent")
	}
}

func TestSpanContext(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	span.sampled = true

	sc := SpanContextFromSpan(span)

	if sc.TraceID() != span.TraceID() {
		t.Error("TraceID should match")
	}

	if sc.SpanID() != span.SpanID() {
		t.Error("SpanID should match")
	}

	if !sc.IsSampled() {
		t.Error("Should be sampled")
	}

	if !sc.IsValid() {
		t.Error("Should be valid")
	}
}

func TestSpanContextInvalid(t *testing.T) {
	sc := SpanContext{}

	if sc.IsValid() {
		t.Error("Empty SpanContext should not be valid")
	}
}

func TestContextWithSpan(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	ctx = ContextWithSpan(ctx, span)

	extracted := SpanFromContext(ctx)
	if extracted != span {
		t.Error("Should extract same span")
	}
}

func TestSpanFromContextNil(t *testing.T) {
	ctx := context.Background()

	span := SpanFromContext(ctx)
	if span != nil {
		t.Error("Should be nil for empty context")
	}
}

func TestSpanContextFromContext(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	ctx = ContextWithSpan(ctx, span)

	sc := SpanContextFromContext(ctx)
	if !sc.IsValid() {
		t.Error("Should be valid")
	}
}

func TestSpanContextFromContextEmpty(t *testing.T) {
	ctx := context.Background()

	sc := SpanContextFromContext(ctx)
	if sc.IsValid() {
		t.Error("Should not be valid for empty context")
	}
}

func TestAlwaysSample(t *testing.T) {
	sampler := &AlwaysSample{}
	result := sampler.ShouldSample(SamplingParameters{})

	if result.Decision != SamplingDecisionRecordAndSample {
		t.Error("AlwaysSample should always sample")
	}
}

func TestNeverSample(t *testing.T) {
	sampler := &NeverSample{}
	result := sampler.ShouldSample(SamplingParameters{})

	if result.Decision != SamplingDecisionDrop {
		t.Error("NeverSample should never sample")
	}
}

func TestTraceIDRatioBased(t *testing.T) {
	// Always sample
	sampler := NewTraceIDRatioBased(1.0)
	result := sampler.ShouldSample(SamplingParameters{TraceID: TraceID{1}})
	if result.Decision != SamplingDecisionRecordAndSample {
		t.Error("Ratio 1.0 should always sample")
	}

	// Never sample
	sampler = NewTraceIDRatioBased(0.0)
	result = sampler.ShouldSample(SamplingParameters{TraceID: TraceID{1}})
	if result.Decision != SamplingDecisionDrop {
		t.Error("Ratio 0.0 should never sample")
	}

	// Test bounds
	sampler = NewTraceIDRatioBased(-0.5)
	if sampler.ratio != 0 {
		t.Error("Negative ratio should be clamped to 0")
	}

	sampler = NewTraceIDRatioBased(1.5)
	if sampler.ratio != 1 {
		t.Error("Ratio > 1 should be clamped to 1")
	}
}

func TestNewTracer(t *testing.T) {
	tracer := NewTracer("test")

	if tracer.Name() != "test" {
		t.Errorf("Expected name 'test', got %s", tracer.Name())
	}

	if tracer.Version() != "" {
		t.Error("Default version should be empty")
	}
}

func TestTracerWithOptions(t *testing.T) {
	sampler := &NeverSample{}
	tracer := NewTracer("test",
		WithSampler(sampler),
		WithVersion("1.0.0"),
		WithResource(String("service", "test")),
	)

	if tracer.Name() != "test" {
		t.Errorf("Expected name 'test', got %s", tracer.Name())
	}

	if tracer.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", tracer.Version())
	}
}

func TestNewTracerProvider(t *testing.T) {
	provider := NewTracerProvider()

	if provider == nil {
		t.Fatal("Provider is nil")
	}
}

func TestTracerProviderTracer(t *testing.T) {
	provider := NewTracerProvider()

	tracer1 := provider.Tracer("tracer1")
	if tracer1 == nil {
		t.Fatal("Tracer is nil")
	}

	// Same tracer should be returned
	tracer2 := provider.Tracer("tracer1")
	if tracer1 != tracer2 {
		t.Error("Same tracer should be returned for same name")
	}

	// Different tracer for different name
	tracer3 := provider.Tracer("tracer2")
	if tracer1 == tracer3 {
		t.Error("Different tracers for different names")
	}
}

func TestTracerProviderConcurrent(t *testing.T) {
	provider := NewTracerProvider()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = provider.Tracer("test")
		}()
	}
	wg.Wait()
}

func TestTracerProviderSetSampler(t *testing.T) {
	provider := NewTracerProvider()
	provider.SetSampler(&NeverSample{})

	// Get tracer after setting sampler
	tracer := provider.Tracer("test")
	if tracer.sampler == nil {
		t.Error("Tracer should have sampler")
	}
}

func TestTracerProviderShutdown(t *testing.T) {
	provider := NewTracerProvider()
	err := provider.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Shutdown should succeed: %v", err)
	}
}

func TestGlobalTracerProvider(t *testing.T) {
	provider := NewTracerProvider()
	SetGlobalTracerProvider(provider)

	retrieved := GetGlobalTracerProvider()
	if retrieved != provider {
		t.Error("Should retrieve same provider")
	}

	tracer := GetTracer("test")
	if tracer == nil {
		t.Error("Should get tracer")
	}
}

func TestStartGlobal(t *testing.T) {
	SetGlobalTracerProvider(NewTracerProvider())

	ctx := context.Background()
	ctx, span := Start(ctx, "test-operation")

	if span == nil {
		t.Fatal("Span is nil")
	}

	if span.Name() != "test-operation" {
		t.Errorf("Expected name 'test-operation', got %s", span.Name())
	}

	span.End()
}

func TestMapCarrier(t *testing.T) {
	carrier := MapCarrier{}

	carrier.Set("key1", "value1")
	carrier.Set("key2", "value2")

	if carrier.Get("key1") != "value1" {
		t.Error("Get should return set value")
	}

	keys := carrier.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestW3CTraceContextPropagatorInject(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	// Start span and use returned context
	ctx, span := tracer.Start(ctx, "test")
	span.sampled = true

	propagator := &W3CTraceContextPropagator{}
	carrier := MapCarrier{}

	// Use context that has the span
	propagator.Inject(ctx, carrier)

	traceparent := carrier.Get("traceparent")
	if traceparent == "" {
		t.Fatal("traceparent header should be set")
	}

	// Format: version-traceid-parentid-flags
	// Should start with 00-
	if traceparent[:3] != "00-" {
		t.Errorf("Invalid traceparent format: %s", traceparent)
	}

	// Should end with -01 for sampled
	if traceparent[len(traceparent)-2:] != "01" {
		t.Errorf("Sampled span should have flags 01: %s", traceparent)
	}
}

func TestW3CTraceContextPropagatorInjectNotSampled(t *testing.T) {
	tracer := NewTracer("test", WithSampler(&NeverSample{}))
	ctx := context.Background()

	// Use returned context
	ctx, span := tracer.Start(ctx, "test")
	defer span.End()

	propagator := &W3CTraceContextPropagator{}
	carrier := MapCarrier{}

	propagator.Inject(ctx, carrier)

	traceparent := carrier.Get("traceparent")
	if traceparent == "" {
		t.Fatal("traceparent header should be set")
	}

	// Should end with -00 for not sampled
	if traceparent[len(traceparent)-2:] != "00" {
		t.Errorf("Not sampled span should have flags 00: %s", traceparent)
	}
}

func TestW3CTraceContextPropagatorInjectNoSpan(t *testing.T) {
	propagator := &W3CTraceContextPropagator{}
	carrier := MapCarrier{}

	propagator.Inject(context.Background(), carrier)

	traceparent := carrier.Get("traceparent")
	if traceparent != "" {
		t.Error("traceparent should not be set without span")
	}
}

func TestW3CTraceContextPropagatorExtract(t *testing.T) {
	propagator := &W3CTraceContextPropagator{}

	traceID := "01234567890123456789012345678901"
	spanID := "0123456789012345"
	traceparent := "00-" + traceID + "-" + spanID + "-01"

	carrier := MapCarrier{}
	carrier.Set("traceparent", traceparent)

	ctx := propagator.Extract(context.Background(), carrier)

	span := SpanFromContext(ctx)
	if span == nil {
		t.Fatal("Span should be extracted")
	}

	expectedTraceID, _ := hex.DecodeString(traceID)
	var tid TraceID
	copy(tid[:], expectedTraceID)

	if span.TraceID() != tid {
		t.Errorf("TraceID mismatch: expected %s, got %s", tid.String(), span.TraceID().String())
	}
}

func TestW3CTraceContextPropagatorExtractInvalid(t *testing.T) {
	tests := []string{
		"",                          // empty
		"invalid",                   // wrong format
		"01-1234-5678-90",           // wrong version
		"00-invalid-trace-id-01",    // invalid trace ID (too short)
		"00-01234567890123456789012345678901-invalid-01", // invalid span ID (too short)
	}

	propagator := &W3CTraceContextPropagator{}

	for _, traceparent := range tests {
		carrier := MapCarrier{}
		carrier.Set("traceparent", traceparent)

		ctx := propagator.Extract(context.Background(), carrier)

		// Should not panic, just return empty context
		span := SpanFromContext(ctx)
		if traceparent == "" && span != nil {
			t.Error("Empty traceparent should not create span")
		}
	}
}

func TestBatchProcessor(t *testing.T) {
	exporter := &mockExporter{spans: make([]*Span, 0)}
	processor := NewBatchProcessor(exporter,
		WithBatchSize(2),
		WithBatchInterval(100*time.Millisecond),
	)

	tracer := NewTracer("test", WithProcessor(processor))
	ctx := context.Background()

	// Create spans
	_, span1 := tracer.Start(ctx, "span1")
	span1.End()

	_, span2 := tracer.Start(ctx, "span2")
	span2.End()

	// Wait for batch
	time.Sleep(200 * time.Millisecond)

	processor.Shutdown(context.Background())

	if len(exporter.spans) < 2 {
		t.Errorf("Expected at least 2 spans exported, got %d", len(exporter.spans))
	}
}

func TestBatchProcessorOnStart(t *testing.T) {
	exporter := &mockExporter{spans: make([]*Span, 0)}
	processor := NewBatchProcessor(exporter)

	// OnStart should be a no-op
	processor.OnStart(&Span{})
}

func TestBatchProcessorShutdown(t *testing.T) {
	exporter := &mockExporter{spans: make([]*Span, 0)}
	processor := NewBatchProcessor(exporter)

	err := processor.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Shutdown should succeed: %v", err)
	}
}

func TestBatchProcessorNilExporter(t *testing.T) {
	processor := NewBatchProcessor(nil)

	// Should not panic
	processor.OnEnd(&Span{})
	processor.flush()
	processor.Shutdown(context.Background())
}

type mockExporter struct {
	spans []*Span
	mu    sync.Mutex
}

func (e *mockExporter) ExportSpans(spans []*Span) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.spans = append(e.spans, spans...)
	return nil
}

func (e *mockExporter) Shutdown(ctx context.Context) error {
	return nil
}

func TestTracerProviderWithProcessor(t *testing.T) {
	exporter := &mockExporter{spans: make([]*Span, 0)}
	processor := NewBatchProcessor(exporter)

	provider := NewTracerProvider()
	provider.SetProcessor(processor)

	tracer := provider.Tracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	span.End()

	// Give time for export
	time.Sleep(100 * time.Millisecond)
	processor.Shutdown(context.Background())
}

func TestSpanDuration(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")

	// Duration is 0 before End()
	if span.Duration() != 0 {
		t.Error("Duration should be 0 before End()")
	}

	time.Sleep(10 * time.Millisecond)
	span.End()

	if span.Duration() < 10*time.Millisecond {
		t.Error("Duration should be at least 10ms")
	}
}

func TestSpanStartTimeEndTime(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	before := time.Now()
	_, span := tracer.Start(ctx, "test")
	after := time.Now()

	if span.StartTime().Before(before) || span.StartTime().After(after) {
		t.Error("StartTime should be around Now()")
	}

	if !span.EndTime().IsZero() {
		t.Error("EndTime should be zero before End()")
	}

	span.End()

	if span.EndTime().IsZero() {
		t.Error("EndTime should be set after End()")
	}
}

func TestParseTraceID(t *testing.T) {
	id, err := parseTraceID("0102030405060708090a0b0c0d0e0f10")
	if err != nil {
		t.Fatalf("parseTraceID failed: %v", err)
	}

	expected := TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	if id != expected {
		t.Errorf("Expected %v, got %v", expected, id)
	}
}

func TestParseSpanID(t *testing.T) {
	id, err := parseSpanID("0102030405060708")
	if err != nil {
		t.Fatalf("parseSpanID failed: %v", err)
	}

	expected := SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	if id != expected {
		t.Errorf("Expected %v, got %v", expected, id)
	}
}

func TestTracerExport(t *testing.T) {
	exporter := &mockExporter{spans: make([]*Span, 0)}
	processor := NewBatchProcessor(exporter)

	tracer := NewTracer("test", WithProcessor(processor))
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	span.End()

	// Give time for processing
	time.Sleep(100 * time.Millisecond)
	processor.Shutdown(context.Background())
}

func TestConcurrentSpans(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, span := tracer.Start(ctx, "concurrent-span")
			time.Sleep(time.Millisecond)
			span.End()
		}()
	}
	wg.Wait()
}

func TestSpanContextWithSampled(t *testing.T) {
	tracer := NewTracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")

	// Test SpanContextFromContext
	ctx = ContextWithSpan(ctx, span)
	sc := SpanContextFromContext(ctx)

	if !sc.IsValid() {
		t.Error("SpanContext should be valid")
	}
}

func TestTracerProviderWithResource(t *testing.T) {
	provider := NewTracerProvider()
	provider.resource = []Attribute{String("service.name", "test-service")}

	tracer := provider.Tracer("test")
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	span.End()

	// Resource attributes should be added to span
}

func TestTracerProviderWithVersion(t *testing.T) {
	provider := NewTracerProvider()
	provider.version = "1.0.0"

	tracer := provider.Tracer("test")

	if tracer.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", tracer.Version())
	}
}

func TestTracerWithProcessorOption(t *testing.T) {
	exporter := &mockExporter{spans: make([]*Span, 0)}
	processor := NewBatchProcessor(exporter)

	tracer := NewTracer("test", WithProcessor(processor))
	ctx := context.Background()

	_, span := tracer.Start(ctx, "test")
	span.End()

	time.Sleep(100 * time.Millisecond)
	processor.Shutdown(context.Background())
}

func TestSamplingWithAttributes(t *testing.T) {
	sampler := &testSampler{}
	tracer := NewTracer("test", WithSampler(sampler))

	ctx := context.Background()
	_, span := tracer.Start(ctx, "test",
		WithAttributes(String("key", "value")),
	)
	span.End()

	if !sampler.called {
		t.Error("Sampler should have been called")
	}
}

type testSampler struct {
	called bool
}

func (s *testSampler) ShouldSample(params SamplingParameters) SamplingResult {
	s.called = true
	return SamplingResult{Decision: SamplingDecisionRecordAndSample}
}

func TestSamplingResultAttributes(t *testing.T) {
	sampler := &attrSampler{attrs: []Attribute{String("sampled", "true")}}
	tracer := NewTracer("test", WithSampler(sampler))

	ctx := context.Background()
	_, span := tracer.Start(ctx, "test")
	span.End()
}

type attrSampler struct {
	attrs []Attribute
}

func (s *attrSampler) ShouldSample(params SamplingParameters) SamplingResult {
	return SamplingResult{
		Decision:   SamplingDecisionRecordAndSample,
		Attributes: s.attrs,
	}
}

func TestPropagatorRoundTrip(t *testing.T) {
	tracer := NewTracer("test")
	propagator := &W3CTraceContextPropagator{}

	// Create span
	ctx, span := tracer.Start(context.Background(), "test")
	span.sampled = true

	// Inject
	carrier := MapCarrier{}
	propagator.Inject(ctx, carrier)

	// Extract
	ctx2 := propagator.Extract(context.Background(), carrier)
	extractedSpan := SpanFromContext(ctx2)

	if extractedSpan == nil {
		t.Fatal("Should extract span")
	}

	if extractedSpan.TraceID() != span.TraceID() {
		t.Error("TraceID should match")
	}

	if extractedSpan.SpanID() != span.SpanID() {
		t.Error("SpanID should match")
	}
}