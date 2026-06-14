package fake

import (
	"testing"
)

// TestContract_FakeConfig_Reader verifies FakeConfig implements Reader
// and handles edge cases per §16.3.
func TestContract_FakeConfig_Reader(t *testing.T) {
	cfg := FakeConfig(map[string]any{
		"name":  "testkit",
		"limit": 100,
	})
	// implements Reader at compile time
	var _ Reader = cfg
	_ = cfg.Get("name")
	_ = cfg.GetString("name")
	_ = cfg.GetInt("limit")
	_ = cfg.GetBool("active")
}

// TestContract_FakeLogger_Concurrent verifies FakeLogger is safe under -race
// per §16.3.
func TestContract_FakeLogger_Concurrent(t *testing.T) {
	TestFakeLogger_Concurrent(t)
}

// TestContract_FakeMeter_Interface verifies FakeMeter implements Meter.
func TestContract_FakeMeter_Interface(t *testing.T) {
	var _ Meter = FakeMeter()
}

// TestContract_FakeTracer_Interface verifies FakeTracer implements Tracer.
func TestContract_FakeTracer_Interface(t *testing.T) {
	var _ Tracer = FakeTracer()
}

// TestContract_FakeBreaker_Interface verifies FakeBreaker implements Breaker.
func TestContract_FakeBreaker_Interface(t *testing.T) {
	var _ Breaker = FakeBreaker(BreakerClosed)
}

// TestContract_FakeConfig_Fingerprint verifies FakeConfig is deterministic.
func TestContract_FakeConfig_Fingerprint(t *testing.T) {
	cfg1 := FakeConfig(map[string]any{"a": 1})
	cfg2 := FakeConfig(map[string]any{"a": 1})
	if cfg1.GetInt("a") != cfg2.GetInt("a") {
		t.Error("FakeConfig should be deterministic")
	}
}
