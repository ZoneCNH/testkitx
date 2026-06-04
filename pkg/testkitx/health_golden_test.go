package testkitx_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx"
	"github.com/ZoneCNH/testkitx/testkit"
)

func TestHealthStatusJSONGolden(t *testing.T) {
	t.Parallel()
	payload, err := json.Marshal(testkitx.HealthStatus{
		Name:      "testkitx",
		Status:    testkitx.HealthHealthy,
		Message:   "ok",
		CheckedAt: time.Unix(0, 0).UTC(),
		LatencyMs: 7,
		Metadata: map[string]string{
			"kind": "template",
		},
	})
	if err != nil {
		t.Fatalf("marshal health status: %v", err)
	}

	payload = append(payload, '\n')
	testkit.RequireGolden(t, "testdata/golden/health_status.json", payload)
}
