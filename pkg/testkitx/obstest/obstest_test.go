package obstest_test

import (
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/obstest"
)

func TestRecorderCapturesCountersAndLogs(t *testing.T) {
	recorder := obstest.NewRecorder()
	recorder.Inc("requests")
	recorder.Inc("requests")
	recorder.Log("started")
	if got := recorder.Count("requests"); got != 2 {
		t.Fatalf("Count() = %d, want 2", got)
	}
	if len(recorder.Logs) != 1 || recorder.Logs[0] != "started" {
		t.Fatalf("unexpected logs: %+v", recorder.Logs)
	}
}
