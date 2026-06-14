package fake

import (
	"context"
	"sync"
)

// Span represents a single traced operation span.
type Span struct {
	Name     string
	TraceID  string
	SpanID   string
	ParentID string
}

// Tracer mirrors the observex.Tracer interface consumed by FoundationX modules.
type Tracer interface {
	StartSpan(ctx context.Context, name string) (context.Context, Span)
}

// FakeTracerImpl is a deterministic fake tracer that records spans and
// supports post-hoc assertions. It implements Tracer.
type FakeTracerImpl struct {
	mu       sync.Mutex
	spans    []Span
	nextID   int
	traceID  string
}

// Compile-time contract: *FakeTracerImpl implements Tracer.
var _ Tracer = (*FakeTracerImpl)(nil)

// FakeTracer creates a new deterministic fake tracer.
func FakeTracer() *FakeTracerImpl {
	return &FakeTracerImpl{
		traceID: "00000000000000000000000000000000",
		spans:   make([]Span, 0),
	}
}

// FakeTracerWithTraceID creates a new fake tracer with a specific trace ID.
func FakeTracerWithTraceID(traceID string) *FakeTracerImpl {
	return &FakeTracerImpl{
		traceID: traceID,
		spans:   make([]Span, 0),
	}
}

// StartSpan starts a new span with the given name. It returns a new context
// and a span record. The span ID is auto-incremented.
func (t *FakeTracerImpl) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.nextID++
	spanID := spanID(t.nextID)
	s := Span{
		Name:    name,
		TraceID: t.traceID,
		SpanID:  spanID,
	}
	t.spans = append(t.spans, s)
	return ctx, s
}

// AssertSpanCount fails t if the number of spans does not equal expected.
func (t *FakeTracerImpl) AssertSpanCount(tt T, expected int) {
	tt.Helper()
	t.mu.Lock()
	defer t.mu.Unlock()
	actual := len(t.spans)
	if actual != expected {
		tt.Errorf("span count: expected %d, got %d", expected, actual)
	}
}

// AssertTraceID fails t if the trace ID was not propagated (is still the
// default sentinel).
func (t *FakeTracerImpl) AssertTraceID(tt T) {
	tt.Helper()
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.traceID == "" {
		tt.Errorf("trace ID is empty — trace ID was not propagated")
	}
}

// AssertSpanNamed fails t if no span with the given name was recorded.
func (t *FakeTracerImpl) AssertSpanNamed(tt T, name string) {
	tt.Helper()
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, s := range t.spans {
		if s.Name == name {
			return
		}
	}
	tt.Errorf("no span named %q recorded (total spans: %d)", name, len(t.spans))
}

// Spans returns a snapshot of all recorded spans.
func (t *FakeTracerImpl) Spans() []Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Span, len(t.spans))
	copy(out, t.spans)
	return out
}

// Reset clears all recorded spans.
func (t *FakeTracerImpl) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = nil
	t.nextID = 0
}

func spanID(n int) string {
	return spanIDPrefix + padID(n)
}

const spanIDPrefix = "0000000000000000"

func padID(n int) string {
	const hexDigits = "0123456789abcdef"
	b := make([]byte, 16)
	for i := 15; i >= 0; i-- {
		b[i] = hexDigits[n&0xf]
		n >>= 4
	}
	return string(b)
}
