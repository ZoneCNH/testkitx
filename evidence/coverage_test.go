package evidence

import (
	"testing"
	"time"
)

func TestFailedConstructor(t *testing.T) {
	t.Parallel()
	c := Failed("check1", "something broke")
	if c.Name != "check1" || c.Status != StatusFailed || c.Message != "something broke" {
		t.Fatalf("unexpected check: %+v", c)
	}
}

func TestSkippedConstructor(t *testing.T) {
	t.Parallel()
	c := Skipped("check2", "not applicable")
	if c.Name != "check2" || c.Status != StatusSkipped || c.Message != "not applicable" {
		t.Fatalf("unexpected check: %+v", c)
	}
}

func TestRunValidateEdgeCases(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	// missing suite
	r := Run{StartedAt: now, EndedAt: now, Cases: []Case{{ID: "x", Name: "y", Status: StatusPass}}}
	if err := r.Validate(); err == nil || err.Error() != "suite is required" {
		t.Fatalf("expected suite error, got %v", err)
	}

	// missing started_at
	r = Run{Suite: "s", EndedAt: now, Cases: []Case{{ID: "x", Name: "y", Status: StatusPass}}}
	if err := r.Validate(); err == nil || err.Error() != "started_at is required" {
		t.Fatalf("expected started_at error, got %v", err)
	}

	// missing ended_at
	r = Run{Suite: "s", StartedAt: now, Cases: []Case{{ID: "x", Name: "y", Status: StatusPass}}}
	if err := r.Validate(); err == nil || err.Error() != "ended_at is required" {
		t.Fatalf("expected ended_at error, got %v", err)
	}

	// ended before started
	r = Run{Suite: "s", StartedAt: now, EndedAt: now.Add(-time.Second), Cases: []Case{{ID: "x", Name: "y", Status: StatusPass}}}
	if err := r.Validate(); err == nil || err.Error() != "ended_at must not be before started_at" {
		t.Fatalf("expected time order error, got %v", err)
	}

	// empty cases
	r = Run{Suite: "s", StartedAt: now, EndedAt: now}
	if err := r.Validate(); err == nil || err.Error() != "cases must not be empty" {
		t.Fatalf("expected empty cases error, got %v", err)
	}

	// missing case id
	r = Run{Suite: "s", StartedAt: now, EndedAt: now, Cases: []Case{{Name: "y", Status: StatusPass}}}
	if err := r.Validate(); err == nil || err.Error() != "cases[0].id is required" {
		t.Fatalf("expected id error, got %v", err)
	}

	// missing case name
	r = Run{Suite: "s", StartedAt: now, EndedAt: now, Cases: []Case{{ID: "x", Status: StatusPass}}}
	if err := r.Validate(); err == nil || err.Error() != "cases[0].name is required" {
		t.Fatalf("expected name error, got %v", err)
	}

	// invalid case status
	r = Run{Suite: "s", StartedAt: now, EndedAt: now, Cases: []Case{{ID: "x", Name: "y", Status: "bad"}}}
	if err := r.Validate(); err == nil || err.Error() != "cases[0].status is invalid" {
		t.Fatalf("expected status error, got %v", err)
	}

	// negative duration
	r = Run{Suite: "s", StartedAt: now, EndedAt: now, Cases: []Case{{ID: "x", Name: "y", Status: StatusPass, Duration: -1}}}
	if err := r.Validate(); err == nil || err.Error() != "cases[0].duration must not be negative" {
		t.Fatalf("expected duration error, got %v", err)
	}
}

func TestReportValidateEdgeCases(t *testing.T) {
	t.Parallel()

	// empty schema_version
	r := Report{ID: "id", Status: StatusPassed, GeneratedAt: time.Now().UTC(), Checks: []Check{{Name: "c", Status: StatusPassed}}}
	if err := r.Validate(); err == nil || err.Error() != "schema_version is required" {
		t.Fatalf("expected schema_version error, got %v", err)
	}

	// invalid schema_version
	r = Report{SchemaVersion: "v9", ID: "id", Status: StatusPassed, GeneratedAt: time.Now().UTC(), Checks: []Check{{Name: "c", Status: StatusPassed}}}
	if err := r.Validate(); err == nil || err.Error() != "schema_version is invalid" {
		t.Fatalf("expected schema_version error, got %v", err)
	}

	// empty id
	r = NewReport("", Passed("c"))
	if err := r.Validate(); err == nil || err.Error() != "id is required" {
		t.Fatalf("expected id error, got %v", err)
	}

	// empty status
	r = Report{SchemaVersion: SchemaVersion, ID: "id", GeneratedAt: time.Now().UTC(), Checks: []Check{{Name: "c", Status: StatusPassed}}}
	if err := r.Validate(); err == nil || err.Error() != "status is required" {
		t.Fatalf("expected status error, got %v", err)
	}

	// invalid status
	r = Report{SchemaVersion: SchemaVersion, ID: "id", Status: "bad", GeneratedAt: time.Now().UTC(), Checks: []Check{{Name: "c", Status: StatusPassed}}}
	if err := r.Validate(); err == nil || err.Error() != "status is invalid" {
		t.Fatalf("expected status error, got %v", err)
	}

	// zero generated_at
	r = NewReport("id", Passed("c"))
	r.GeneratedAt = time.Time{}
	if err := r.Validate(); err == nil || err.Error() != "generated_at is required" {
		t.Fatalf("expected generated_at error, got %v", err)
	}

	// empty checks
	r = NewReport("id")
	if err := r.Validate(); err == nil || err.Error() != "checks must not be empty" {
		t.Fatalf("expected empty checks error, got %v", err)
	}

	// check missing name
	r = NewReport("id", Check{Status: StatusPassed})
	if err := r.Validate(); err == nil || err.Error() != "checks[0].name is required" {
		t.Fatalf("expected name error, got %v", err)
	}

	// check invalid status
	r = NewReport("id", Check{Name: "c", Status: "bad"})
	if err := r.Validate(); err == nil || err.Error() != "checks[0].status is invalid" {
		t.Fatalf("expected status error, got %v", err)
	}

	// aggregate mismatch
	r = NewReport("id", Passed("c"))
	r.Status = StatusFailed
	if err := r.Validate(); err == nil || err.Error() != "status must match aggregate check status" {
		t.Fatalf("expected aggregate mismatch error, got %v", err)
	}
}

func TestMarshalRejectsInvalidReport(t *testing.T) {
	t.Parallel()
	_, err := Marshal(Report{})
	if err == nil {
		t.Fatal("expected error for invalid report")
	}
}

func TestAggregateStatusAllSkipped(t *testing.T) {
	t.Parallel()
	status := aggregateStatus([]Check{
		{Name: "a", Status: StatusSkipped},
		{Name: "b", Status: StatusSkipped},
	})
	if status != StatusSkipped {
		t.Fatalf("expected skipped, got %s", status)
	}
}

func TestAggregateStatusEmpty(t *testing.T) {
	t.Parallel()
	status := aggregateStatus(nil)
	if status != StatusSkipped {
		t.Fatalf("expected skipped for empty, got %s", status)
	}
}

func TestNewReportWithMultipleChecks(t *testing.T) {
	t.Parallel()
	r := NewReport("id", Passed("a"), Skipped("b", "n/a"))
	if r.Status != StatusSkipped {
		t.Fatalf("expected skipped, got %s", r.Status)
	}
}

func TestWriteFileInvalidPath(t *testing.T) {
	t.Parallel()
	err := WriteFile("/nonexistent/dir/file.json", validRun())
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
