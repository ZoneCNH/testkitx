package evidence

import (
	"testing"
	"time"
)

func TestMarshalInvalidReportMissingSchema(t *testing.T) {
	t.Parallel()
	r := Report{ID: "test", Status: StatusPassed, GeneratedAt: time.Now().UTC(), Checks: []Check{{Name: "c", Status: StatusPassed}}}
	_, err := Marshal(r)
	if err == nil {
		t.Fatal("expected error for missing schema_version")
	}
}

func TestMarshalInvalidReportBadStatus(t *testing.T) {
	t.Parallel()
	r := Report{SchemaVersion: SchemaVersion, ID: "test", Status: "bogus", GeneratedAt: time.Now().UTC(), Checks: []Check{{Name: "c", Status: StatusPassed}}}
	_, err := Marshal(r)
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestValidateReportRejectsStatusMismatch(t *testing.T) {
	t.Parallel()
	r := Report{
		SchemaVersion: SchemaVersion,
		ID:            "test",
		Status:        StatusPassed,
		GeneratedAt:   time.Now().UTC(),
		Checks:        []Check{{Name: "c", Status: StatusFailed}},
	}
	if err := r.Validate(); err == nil || err.Error() != "status must match aggregate check status" {
		t.Fatalf("expected status mismatch error, got %v", err)
	}
}

func TestValidateReportRejectsEmptyCheckName(t *testing.T) {
	t.Parallel()
	r := Report{
		SchemaVersion: SchemaVersion,
		ID:            "test",
		Status:        StatusPassed,
		GeneratedAt:   time.Now().UTC(),
		Checks:        []Check{{Name: "", Status: StatusPassed}},
	}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for empty check name")
	}
}

func TestValidateReportRejectsInvalidCheckStatus(t *testing.T) {
	t.Parallel()
	r := Report{
		SchemaVersion: SchemaVersion,
		ID:            "test",
		Status:        StatusPassed,
		GeneratedAt:   time.Now().UTC(),
		Checks:        []Check{{Name: "c", Status: "invalid"}},
	}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for invalid check status")
	}
}

