package testkitx

import (
	"context"
	"time"
)

type HealthStatusValue string

const (
	HealthHealthy   HealthStatusValue = "healthy"
	HealthDegraded  HealthStatusValue = "degraded"
	HealthUnhealthy HealthStatusValue = "unhealthy"
)

type HealthStatus struct {
	Name      string            `json:"name"`
	Status    HealthStatusValue `json:"status"`
	Message   string            `json:"message,omitempty"`
	CheckedAt time.Time         `json:"checked_at"`
	LatencyMs int64             `json:"latency_ms"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

func newHealthStatus(name string, status HealthStatusValue, message string, start time.Time, metadata map[string]string) HealthStatus {
	return HealthStatus{
		Name:      name,
		Status:    status,
		Message:   message,
		CheckedAt: time.Now(),
		LatencyMs: time.Since(start).Milliseconds(),
		Metadata:  metadata,
	}
}

func (c *Client) HealthCheck(ctx context.Context) HealthStatus {
	start := time.Now()
	name := "testkitx"
	var metrics Metrics
	initialized := false
	closed := true
	var timeout time.Duration

	if c != nil {
		c.mu.Lock()
		name = c.cfg.Name
		metrics = c.metrics
		initialized = c.initialized
		closed = c.closed
		timeout = c.cfg.Timeout
		c.mu.Unlock()
		if name == "" {
			name = "testkitx"
		}
	}

	if ctx == nil {
		status := newHealthStatus(name, HealthUnhealthy, "context is required", start, nil)
		recordHealthMetric(metrics, status)
		return status
	}

	if err := ctx.Err(); err != nil {
		status := newHealthStatus(name, HealthUnhealthy, err.Error(), start, nil)
		recordHealthMetric(metrics, status)
		return status
	}

	if !initialized {
		status := newHealthStatus(name, HealthUnhealthy, "client is not initialized", start, nil)
		recordHealthMetric(metrics, status)
		return status
	}

	if closed {
		status := newHealthStatus(name, HealthUnhealthy, "client is closed", start, nil)
		recordHealthMetric(metrics, status)
		return status
	}

	if timeout > 0 {
		if deadline, ok := ctx.Deadline(); ok {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				message := context.DeadlineExceeded.Error()
				if err := ctx.Err(); err != nil {
					message = err.Error()
				}
				status := newHealthStatus(name, HealthUnhealthy, message, start, nil)
				recordHealthMetric(metrics, status)
				return status
			}
			if remaining < timeout {
				status := newHealthStatus(name, HealthDegraded, "context deadline is shorter than client timeout", start, map[string]string{
					"reason":  "deadline_below_timeout",
					"timeout": timeout.String(),
				})
				recordHealthMetric(metrics, status)
				return status
			}
		}
	}

	status := newHealthStatus(name, HealthHealthy, "ok", start, nil)
	recordHealthMetric(metrics, status)
	return status
}

func recordHealthMetric(metrics Metrics, status HealthStatus) {
	if metrics == nil {
		return
	}
	labels := map[string]string{
		"name":   status.Name,
		"status": string(status.Status),
	}
	metrics.SetGauge(MetricClientHealthStatus, healthGaugeValue(status.Status), labels)
	metrics.ObserveHistogram(MetricClientHealthLatencyMS, float64(status.LatencyMs), labels)
}

func healthGaugeValue(status HealthStatusValue) float64 {
	if status == HealthHealthy {
		return 1
	}
	return 0
}
