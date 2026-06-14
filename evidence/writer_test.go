package evidence

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriterWritesStructuredJSON(t *testing.T) {
	t.Parallel()
	run := validRun()
	var buf bytes.Buffer
	if err := NewWriter(&buf).Write(run); err != nil {
		t.Fatalf("write evidence: %v", err)
	}
	var got Run
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("decode evidence: %v", err)
	}
	if got.Suite != "l2" || len(got.Cases) != 1 || got.Cases[0].ID != "common.lifecycle.start" {
		t.Fatalf("unexpected decoded run: %+v", got)
	}
}

func TestWriterRejectsIncompleteRun(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	if err := NewWriter(&buf).Write(Run{}); err == nil || !strings.Contains(err.Error(), "suite is required") {
		t.Fatalf("expected suite validation failure, got %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected invalid evidence not to be written, got %q", buf.String())
	}
}

func TestWriterRejectsEmptyCaseRun(t *testing.T) {
	t.Parallel()
	run := validRun()
	run.Cases = nil
	var buf bytes.Buffer
	if err := NewWriter(&buf).Write(run); err == nil || !strings.Contains(err.Error(), "cases must not be empty") {
		t.Fatalf("expected case validation failure, got %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected invalid evidence not to be written, got %q", buf.String())
	}
}

func TestWriteFileValidatesBeforeCreatingFile(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "nested", "evidence.json")
	if err := WriteFile(path, Run{}); err == nil || !strings.Contains(err.Error(), "suite is required") {
		t.Fatalf("expected suite validation failure, got %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected invalid evidence file not to exist, got %v", err)
	}
}

func TestWriteFileWritesNestedJSON(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "nested", "evidence.json")
	if err := WriteFile(path, validRun()); err != nil {
		t.Fatalf("write evidence file: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read evidence file: %v", err)
	}
	if !bytes.HasSuffix(data, []byte("\n")) {
		t.Fatalf("expected trailing newline in evidence JSON")
	}
	var decoded Run
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode evidence file: %v", err)
	}
	if decoded.Suite != "l2" || len(decoded.Cases) != 1 {
		t.Fatalf("unexpected decoded run: %+v", decoded)
	}
}

func TestMarshalWritesReportJSON(t *testing.T) {
	t.Parallel()
	report := NewReport("unit", Passed("go test"))
	data, err := Marshal(report)
	if err != nil {
		t.Fatalf("marshal evidence: %v", err)
	}
	if !bytes.HasSuffix(data, []byte("\n")) {
		t.Fatalf("expected trailing newline in evidence JSON")
	}
	if !json.Valid(data) {
		t.Fatalf("expected valid JSON: %s", data)
	}
	var decoded Report
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode report: %v", err)
	}
	if decoded.SchemaVersion != SchemaVersion || decoded.ID != "unit" || decoded.Status != StatusPassed {
		t.Fatalf("unexpected decoded report: %+v", decoded)
	}
}

func TestValidateRejectsIncompleteEvidence(t *testing.T) {
	t.Parallel()
	report := NewReport("", Passed("go test"))
	if err := report.Validate(); err == nil || !strings.Contains(err.Error(), "id is required") {
		t.Fatalf("expected id validation failure, got %v", err)
	}
	report = NewReport("suite", Check{Status: StatusPassed})
	if err := report.Validate(); err == nil || !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("expected check validation failure, got %v", err)
	}
	report = NewReport("suite", Check{Name: "go test", Status: Status("unknown")})
	if err := report.Validate(); err == nil || !strings.Contains(err.Error(), "status is invalid") {
		t.Fatalf("expected check status validation failure, got %v", err)
	}
	report = NewReport("suite", Passed("go test"))
	report.SchemaVersion = "testkitx.evidence.v0"
	if err := report.Validate(); err == nil || !strings.Contains(err.Error(), "schema_version is invalid") {
		t.Fatalf("expected schema validation failure, got %v", err)
	}
	report = NewReport("suite", Passed("go test"))
	report.GeneratedAt = time.Time{}
	if err := report.Validate(); err == nil || !strings.Contains(err.Error(), "generated_at is required") {
		t.Fatalf("expected generated_at validation failure, got %v", err)
	}
	report = NewReport("suite", Passed("go test"))
	report.Status = StatusFailed
	if err := report.Validate(); err == nil || !strings.Contains(err.Error(), "status must match aggregate") {
		t.Fatalf("expected aggregate status validation failure, got %v", err)
	}
	report = NewReport("suite")
	if err := report.Validate(); err == nil || !strings.Contains(err.Error(), "checks must not be empty") {
		t.Fatalf("expected empty checks validation failure, got %v", err)
	}
}

func TestDigestReturnsSHA256Hex(t *testing.T) {
	t.Parallel()
	got := Digest([]byte("evidence-digest-fixture"))
	if got != "46e6a840bf79bf3423bdda895e9a34df27a14927950b16f626b155265637a071" {
		t.Fatalf("unexpected digest: %s", got)
	}
}

func validRun() Run {
	return Run{
		Suite:     "l2",
		StartedAt: time.Unix(1, 0).UTC(),
		EndedAt:   time.Unix(2, 0).UTC(),
		Cases:     []Case{{ID: "common.lifecycle.start", Name: "start", Status: StatusPass}},
	}
}

func TestWriteFileError(t *testing.T) {
	t.Parallel()
	err := WriteFile("/nonexistent\x00dir/evidence.json", Run{Suite: "test"})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestWriteFileCreatesDir(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "deep", "nested", "evidence.json")
	run := Run{
		Suite:     "test",
		StartedAt: time.Now().Add(-time.Second),
		EndedAt:   time.Now(),
		Cases:     []Case{{ID: "1", Name: "test", Status: StatusPass}},
	}
	if err := WriteFile(path, run); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty evidence file")
	}
}

func TestWriteFileMkdirAllError(t *testing.T) {
	t.Parallel()
	run := validRun()
	err := WriteFile("/dev/null/impossible/evidence.json", run)
	if err == nil {
		t.Fatal("expected error for impossible directory")
	}
}

func TestWriteFileSimpleFilename(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "evidence.json")
	run := validRun()
	if err := WriteFile(path, run); err != nil {
		t.Fatalf("WriteFile with simple filename: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty evidence file")
	}
}

// failingWriter always returns an error on Write.
type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestWriteFileCreateError(t *testing.T) {
	t.Parallel()
	// /dev/null is not a directory, so os.Create("/dev/null/sub/file") fails.
	run := validRun()
	err := WriteFile("/dev/null/sub/evidence.json", run)
	if err == nil {
		t.Fatal("expected error creating file under /dev/null")
	}
}

func TestWriteFileEncoderError(t *testing.T) {
	t.Parallel()
	// Write to a file then close the read end of a pipe to cause encoder error.
	// Use os.CreateTemp to get a real file, then overwrite the writer.
	dir := t.TempDir()
	path := filepath.Join(dir, "evidence.json")
	run := validRun()
	// Create the file so MkdirAll is not needed.
	os.WriteFile(path, []byte{}, 0o644)

	// WriteFile creates a new file with os.Create, so we can't easily make
	// the encoder fail after Create succeeds. Instead test the
	// NewWriter + Write path directly with a failing writer.
	w := NewWriter(failingWriter{})
	err := w.Write(run)
	if err == nil {
		t.Fatal("expected encoder error with failing writer")
	}
}
