package fake

import "sync"

// Meter mirrors the observex.Meter interface consumed by FoundationX modules.
type Meter interface {
	IncCounter(name string, labels map[string]string)
	ObserveHistogram(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
}

// FakeMeterImpl is a deterministic fake meter that records operations and
// supports post-hoc assertions. It implements Meter.
type FakeMeterImpl struct {
	mu         sync.Mutex
	counters   map[string]float64
	histograms map[string][]float64
	gauges     map[string]float64
}

// Compile-time contract: *FakeMeterImpl implements Meter.
var _ Meter = (*FakeMeterImpl)(nil)

// FakeMeter creates a new deterministic fake meter.
func FakeMeter() *FakeMeterImpl {
	return &FakeMeterImpl{
		counters:   make(map[string]float64),
		histograms: make(map[string][]float64),
		gauges:     make(map[string]float64),
	}
}

// IncCounter increments the named counter by 1.
func (m *FakeMeterImpl) IncCounter(name string, labels map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name]++
}

// ObserveHistogram records a value for the named histogram.
func (m *FakeMeterImpl) ObserveHistogram(name string, value float64, labels map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.histograms[name] = append(m.histograms[name], value)
}

// SetGauge sets the named gauge to the given value.
func (m *FakeMeterImpl) SetGauge(name string, value float64, labels map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = value
}

// AssertCounterValue fails t if the named counter does not equal expected.
func (m *FakeMeterImpl) AssertCounterValue(t T, name string, expected float64) {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()
	actual := m.counters[name]
	if actual != expected {
		t.Errorf("counter %q: expected %v, got %v", name, expected, actual)
	}
}

// AssertHistogramRecorded fails t if the named histogram has no recorded values.
func (m *FakeMeterImpl) AssertHistogramRecorded(t T, name string) {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.histograms[name]) == 0 {
		t.Errorf("histogram %q: expected at least one recorded value, got none", name)
	}
}

// CounterValue returns the current value of the named counter.
func (m *FakeMeterImpl) CounterValue(name string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.counters[name]
}

// HistogramValues returns all recorded values for the named histogram.
func (m *FakeMeterImpl) HistogramValues(name string) []float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]float64, len(m.histograms[name]))
	copy(out, m.histograms[name])
	return out
}

// GaugeValue returns the current value of the named gauge.
func (m *FakeMeterImpl) GaugeValue(name string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.gauges[name]
}

// Reset clears all recorded metrics.
func (m *FakeMeterImpl) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters = make(map[string]float64)
	m.histograms = make(map[string][]float64)
	m.gauges = make(map[string]float64)
}
