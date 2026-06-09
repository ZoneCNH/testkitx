package evidence

import (
	"strings"
	"testing"
	"time"
)


func TestMarshalValidReport(t *testing.T) {
	t.Parallel()
	report := NewReport("marshal-test", Passed("check1"))
	data, err := Marshal(report)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty marshal output")
	}
}

func TestAggregateStatusEmptyChecks(t *testing.T) {
	t.Parallel()
	got := aggregateStatus(nil)
	if got != StatusSkipped {
		t.Fatalf("empty checks: expected skip, got %q", got)
	}
	got = aggregateStatus([]Check{})
	if got != StatusSkipped {
		t.Fatalf("zero checks: expected skip, got %q", got)
	}
}

func TestAggregateStatusAllPassed(t *testing.T) {
	t.Parallel()
	checks := []Check{
		{Name: "a", Status: StatusPassed},
		{Name: "b", Status: StatusPassed},
	}
	if got := aggregateStatus(checks); got != StatusPassed {
		t.Fatalf("all passed: expected pass, got %q", got)
	}
}


func TestAggregateStatusMixedSkipAndPass(t *testing.T) {
	t.Parallel()
	checks := []Check{
		{Name: "a", Status: StatusPassed},
		{Name: "b", Status: StatusSkipped},
	}
	if got := aggregateStatus(checks); got != StatusSkipped {
		t.Fatalf("mixed pass+skip: expected skip, got %q", got)
	}
}

func TestAggregateStatusFailOverridesAll(t *testing.T) {
	t.Parallel()
	checks := []Check{
		{Name: "a", Status: StatusPassed},
		{Name: "b", Status: StatusFailed},
		{Name: "c", Status: StatusSkipped},
	}
	if got := aggregateStatus(checks); got != StatusFailed {
		t.Fatalf("fail present: expected fail, got %q", got)
	}
}

func TestAggregateStatusFailOnly(t *testing.T) {
	t.Parallel()
	checks := []Check{
		{Name: "a", Status: StatusFailed},
	}
	if got := aggregateStatus(checks); got != StatusFailed {
		t.Fatalf("single fail: expected fail, got %q", got)
	}
}

func TestValidateRunRejectsAllFields(t *testing.T) {
	t.Parallel()
	if err := (Run{}).Validate(); err == nil || !strings.Contains(err.Error(), "suite is required") {
		t.Fatalf("expected suite error, got %v", err)
	}
	if err := (Run{Suite: "s"}).Validate(); err == nil || !strings.Contains(err.Error(), "started_at") {
		t.Fatalf("expected started_at error, got %v", err)
	}
	if err := (Run{Suite: "s", StartedAt: time.Now()}).Validate(); err == nil || !strings.Contains(err.Error(), "ended_at") {
		t.Fatalf("expected ended_at error, got %v", err)
	}
	now := time.Now()
	if err := (Run{Suite: "s", StartedAt: now, EndedAt: now.Add(-time.Hour)}).Validate(); err == nil || !strings.Contains(err.Error(), "ended_at must not be before") {
		t.Fatalf("expected ended_at ordering error, got %v", err)
	}
	if err := (Run{Suite: "s", StartedAt: now, EndedAt: now}).Validate(); err == nil || !strings.Contains(err.Error(), "cases must not be empty") {
		t.Fatalf("expected cases error, got %v", err)
	}
	r := Run{Suite: "s", StartedAt: now, EndedAt: now, Cases: []Case{{Name: "n", Status: StatusPass}}}
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "id is required") {
		t.Fatalf("expected case id error, got %v", err)
	}
	r = Run{Suite: "s", StartedAt: now, EndedAt: now, Cases: []Case{{ID: "1", Status: StatusPass}}}
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("expected case name error, got %v", err)
	}
	r = Run{Suite: "s", StartedAt: now, EndedAt: now, Cases: []Case{{ID: "1", Name: "n", Status: Status("bad")}}}
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "status is invalid") {
		t.Fatalf("expected case status error, got %v", err)
	}
	r = Run{Suite: "s", StartedAt: now, EndedAt: now, Cases: []Case{{ID: "1", Name: "n", Status: StatusPass, Duration: -1}}}
	if err := r.Validate(); err == nil || !strings.Contains(err.Error(), "duration must not be negative") {
		t.Fatalf("expected duration error, got %v", err)
	}
}

func TestValidateReportRejectsInvalidStatus(t *testing.T) {
	t.Parallel()
	report := NewReport("test", Passed("c"))
	report.Status = Status("unknown")
	if err := report.Validate(); err == nil || !strings.Contains(err.Error(), "status is invalid") {
		t.Fatalf("expected status invalid error, got %v", err)
	}
}
