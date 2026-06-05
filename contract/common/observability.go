package common

import "github.com/ZoneCNH/testkitx/requirex"

const ObservabilityMetricsID = "common.observability.metrics"

type MetricsSource interface {
	Metric(name string) (float64, bool)
}

func RunObservabilityMetrics(t requirex.TestingT, source MetricsSource, name string) {
	t.Helper()
	value, ok := source.Metric(name)
	if !ok {
		t.Fatalf("expected metric %q to be recorded", name)
	}
	if value < 0 {
		t.Fatalf("expected metric %q to be non-negative, got %f", name, value)
	}
}
