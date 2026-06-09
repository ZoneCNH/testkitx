package testkitx

import "testing"

func TestNoopMetricsIncCounter(t *testing.T) {
	t.Parallel()
	var m NoopMetrics
	m.IncCounter("test_counter", map[string]string{"k": "v"})
}

func TestNoopMetricsObserveHistogram(t *testing.T) {
	t.Parallel()
	var m NoopMetrics
	m.ObserveHistogram("test_hist", 1.5, nil)
}

func TestNoopMetricsSetGauge(t *testing.T) {
	t.Parallel()
	var m NoopMetrics
	m.SetGauge("test_gauge", 42.0, map[string]string{})
}
