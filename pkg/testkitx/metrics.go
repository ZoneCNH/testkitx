package testkitx

const (
	MetricClientCreatedTotal           = "client_created_total"
	MetricClientClosedTotal            = "client_closed_total"
	MetricClientErrorsTotal            = "client_errors_total"
	MetricClientHealthStatus           = "client_health_status"
	MetricClientHealthLatencyMS        = "client_health_latency_ms"
	MetricClientRequestsTotal          = "client_requests_total"
	MetricClientRequestDurationSeconds = "client_request_duration_seconds"
	MetricClientRetriesTotal           = "client_retries_total"
	MetricClientInflight               = "client_inflight"
)

type Metrics interface {
	IncCounter(name string, labels map[string]string)
	ObserveHistogram(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
}

type NoopMetrics struct{}

//go:noinline
func (NoopMetrics) IncCounter(name string, labels map[string]string) {
	_ = name
}

//go:noinline
func (NoopMetrics) ObserveHistogram(name string, value float64, labels map[string]string) {
	_ = name
	_ = value
}

//go:noinline
func (NoopMetrics) SetGauge(name string, value float64, labels map[string]string) {
	_ = name
	_ = value
}
