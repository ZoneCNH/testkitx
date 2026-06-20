package fake

import (
	"context"
	"testing"
	"time"
)

// ===== FakeConfig Tests (FR-001) =====

func TestFakeConfig_GetString(t *testing.T) {
	cfg := FakeConfig(map[string]any{
		"symbol": "BTCUSDT",
		"count":  42,
		"active": true,
	})

	if got := cfg.GetString("symbol"); got != "BTCUSDT" {
		t.Errorf("GetString(symbol) = %q, want %q", got, "BTCUSDT")
	}
	// Non-string key returns zero value
	if got := cfg.GetString("count"); got != "" {
		t.Errorf("GetString(count) = %q, want empty", got)
	}
}

func TestFakeConfig_GetInt(t *testing.T) {
	cfg := FakeConfig(map[string]any{
		"count":  42,
		"pi":     3.14,
		"symbol": "BTCUSDT",
	})
	if got := cfg.GetInt("count"); got != 42 {
		t.Errorf("GetInt(count) = %d, want 42", got)
	}
	if got := cfg.GetInt("missing"); got != 0 {
		t.Errorf("GetInt(missing) = %d, want 0", got)
	}
}

func TestFakeConfig_GetBool(t *testing.T) {
	cfg := FakeConfig(map[string]any{"active": true, "debug": false})
	if got := cfg.GetBool("active"); !got {
		t.Error("GetBool(active) = false, want true")
	}
	if got := cfg.GetBool("debug"); got {
		t.Error("GetBool(debug) = true, want false")
	}
	if got := cfg.GetBool("missing"); got {
		t.Error("GetBool(missing) = true, want false")
	}
}

func TestFakeConfig_NilValues(t *testing.T) {
	cfg := FakeConfig(nil)
	if got := cfg.Get("anything"); got != nil {
		t.Errorf("Get on nil map: got %v, want nil", got)
	}
	if got := cfg.GetString("x"); got != "" {
		t.Errorf("GetString on nil map: got %q, want empty", got)
	}
}

// ===== FakeLogger Tests (FR-002) =====

func TestFakeLogger_LogLevels(t *testing.T) {
	log := FakeLogger()
	log.Debug("debug msg")
	log.Info("info msg")
	log.Warn("warn msg")
	log.Error("error msg")

	entries := log.Entries()
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
}

func TestFakeLogger_AssertLogged(t *testing.T) {
	log := FakeLogger()
	log.Info("user login successful")

	// should not fail — use a helper TB that captures failures
	log.AssertLogged(t, LevelInfo, "login")
}

func TestFakeLogger_AssertNoErrors_Passes(t *testing.T) {
	log := FakeLogger()
	log.Info("all good")
	log.AssertNoErrors(t)
}

func TestFakeLogger_Entries_ReturnsCopy(t *testing.T) {
	log := FakeLogger()
	log.Info("first")
	entries1 := log.Entries()
	log.Info("second")
	entries2 := log.Entries()

	if len(entries1) != 1 {
		t.Errorf("first snapshot: want 1, got %d", len(entries1))
	}
	if len(entries2) != 2 {
		t.Errorf("second snapshot: want 2, got %d", len(entries2))
	}
}

func TestFakeLogger_Reset(t *testing.T) {
	log := FakeLogger()
	log.Info("msg")
	log.Reset()
	if len(log.Entries()) != 0 {
		t.Error("expected empty after reset")
	}
}

func TestFakeLogger_Fields(t *testing.T) {
	log := FakeLogger()
	log.Info("request", "method", "GET", "status", 200)

	entries := log.Entries()
	if len(entries) != 1 {
		t.Fatal("expected 1 entry")
	}
	if entries[0].Fields["method"] != "GET" {
		t.Errorf("field method: got %v, want GET", entries[0].Fields["method"])
	}
	if entries[0].Fields["status"] != 200 {
		t.Errorf("field status: got %v, want 200", entries[0].Fields["status"])
	}
	// odd number of fields — last key is dropped
}

func TestFakeLogger_Concurrent(t *testing.T) {
	log := FakeLogger()
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(n int) {
			for j := 0; j < 100; j++ {
				log.Info("msg")
			}
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if len(log.Entries()) != 1000 {
		t.Errorf("concurrent entries: want 1000, got %d", len(log.Entries()))
	}
}

// ===== FakeMeter Tests (FR-003) =====

func TestFakeMeter_AssertCounterValue(t *testing.T) {
	m := FakeMeter()
	m.IncCounter("requests", map[string]string{})
	m.IncCounter("requests", map[string]string{})
	m.IncCounter("requests", map[string]string{})

	m.AssertCounterValue(t, "requests", 3)
}

func TestFakeMeter_AssertHistogramRecorded(t *testing.T) {
	m := FakeMeter()
	m.ObserveHistogram("latency", 0.5, map[string]string{})
	m.ObserveHistogram("latency", 0.3, map[string]string{})

	m.AssertHistogramRecorded(t, "latency")
}

func TestFakeMeter_CounterValue(t *testing.T) {
	m := FakeMeter()
	m.IncCounter("errors", map[string]string{})
	m.IncCounter("errors", map[string]string{})

	if got := m.CounterValue("errors"); got != 2 {
		t.Errorf("CounterValue = %v, want 2", got)
	}
	if got := m.CounterValue("missing"); got != 0 {
		t.Errorf("CounterValue(missing) = %v, want 0", got)
	}
}

func TestFakeMeter_Reset(t *testing.T) {
	m := FakeMeter()
	m.IncCounter("x", nil)
	m.SetGauge("y", 1, map[string]string{})
	m.Reset()

	if got := m.CounterValue("x"); got != 0 {
		t.Error("counter not reset")
	}
	if got := m.GaugeValue("y"); got != 0 {
		t.Error("gauge not reset")
	}
}

// ===== FakeTracer Tests (FR-004) =====

func TestFakeTracer_StartSpan(t *testing.T) {
	tr := FakeTracer()
	_, s1 := tr.StartSpan(context.Background(), "operation-1")
	_, s2 := tr.StartSpan(context.Background(), "operation-2")

	if s1.Name != "operation-1" {
		t.Errorf("span 1 name = %q", s1.Name)
	}
	if s2.Name != "operation-2" {
		t.Errorf("span 2 name = %q", s2.Name)
	}
	if s1.TraceID == "" {
		t.Error("trace ID is empty")
	}
	if s1.SpanID == "" {
		t.Error("span ID is empty")
	}
	if s1.SpanID == s2.SpanID {
		t.Error("span IDs should be unique")
	}
}

func TestFakeTracer_AssertSpanCount(t *testing.T) {
	tr := FakeTracer()
	tr.StartSpan(context.Background(), "a")
	tr.StartSpan(context.Background(), "b")
	tr.StartSpan(context.Background(), "c")

	tr.AssertSpanCount(t, 3)
}

func TestFakeTracer_AssertTraceID(t *testing.T) {
	tr := FakeTracer()
	tr.StartSpan(context.Background(), "x")
	tr.AssertTraceID(t) // should pass — trace ID was propagated
}

func TestFakeTracer_AssertSpanNamed(t *testing.T) {
	tr := FakeTracer()
	tr.StartSpan(context.Background(), "checkout")
	tr.AssertSpanNamed(t, "checkout")
}

func TestFakeTracer_Reset(t *testing.T) {
	tr := FakeTracer()
	tr.StartSpan(context.Background(), "a")
	tr.Reset()
	tr.AssertSpanCount(t, 0)
}

// ===== FakeClock Tests (FR-005) =====

func TestFakeClock_Now(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := Clock(base)

	if got := c.Now(); !got.Equal(base) {
		t.Errorf("Now() = %v, want %v", got, base)
	}
}

func TestFakeClock_Advance(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := Clock(base)

	c.Advance(1 * time.Hour)
	if got := c.Now(); !got.Equal(base.Add(1 * time.Hour)) {
		t.Errorf("after Advance(1h): %v, want %v", got, base.Add(1*time.Hour))
	}

	c.Advance(30 * time.Minute)
	if got := c.Now(); !got.Equal(base.Add(90 * time.Minute)) {
		t.Errorf("after Advance(30m): got %v, want %v", got, base.Add(90*time.Minute))
	}
}

func TestFakeClock_Set(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := Clock(base)

	newTime := time.Date(2027, 6, 15, 12, 0, 0, 0, time.UTC)
	c.Set(newTime)

	if got := c.Now(); !got.Equal(newTime) {
		t.Errorf("after Set: %v, want %v", got, newTime)
	}
}

func TestFakeClock_NoAdvance_StaysSame(t *testing.T) {
	c := Clock(time.Unix(0, 0))
	first := c.Now()
	second := c.Now()
	if !first.Equal(second) {
		t.Error("Now() should return same value without Advance/Set")
	}
}

func TestFakeClock_Concurrent(t *testing.T) {
	c := Clock(time.Unix(0, 0))
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				c.Now()
				c.Advance(time.Millisecond)
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

// ===== FakeBreaker Tests (FR-006) =====

func TestFakeBreaker_InitialClosed(t *testing.T) {
	b := FakeBreaker(BreakerClosed)
	if !b.Allow() {
		t.Error("closed breaker should allow")
	}
	if b.State() != BreakerClosed {
		t.Error("state should be closed")
	}
}

func TestFakeBreaker_OpenDenies(t *testing.T) {
	b := FakeBreaker(BreakerOpen)
	if b.Allow() {
		t.Error("open breaker should deny")
	}
}

func TestFakeBreaker_Allow_DeniesOpen(t *testing.T) {
	b := FakeBreaker(BreakerClosed)
	b.RecordFailure() // opens the breaker
	if b.Allow() {
		t.Error("breaker should deny after RecordFailure")
	}
}

func TestFakeBreaker_RecordSuccess_Closes(t *testing.T) {
	b := FakeBreaker(BreakerOpen)
	b.RecordSuccess()
	if b.State() != BreakerClosed {
		t.Error("breaker should be closed after RecordSuccess")
	}
	if !b.Allow() {
		t.Error("should allow after closing")
	}
}

func TestFakeBreaker_HalfOpenAllows(t *testing.T) {
	b := FakeBreaker(BreakerHalfOpen)
	if !b.Allow() {
		t.Error("half-open breaker should allow")
	}
}

func TestFakeBreaker_SetState(t *testing.T) {
	b := FakeBreaker(BreakerHalfOpen)
	impl := b.(*FakeBreakerImpl)
	impl.SetState(BreakerClosed)
	if b.State() != BreakerClosed {
		t.Error("SetState did not take effect")
	}
}

// ===== Compile-time interface checks =====

func TestFakeImplementsContracts(t *testing.T) {
	// These calls keep the concrete fakes exercised by the test build.
	_ = FakeConfig(nil)
	var _ Logger = FakeLogger()
	var _ Meter = FakeMeter()
	var _ Tracer = FakeTracer()
	_ = FakeBreaker(BreakerClosed)
	_ = Clock(time.Now())
}
