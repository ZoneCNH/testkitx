package testkitx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNewErrorFormatsKindOpAndMessage(t *testing.T) {
	t.Parallel()
	err := NewError(ErrorKindValidation, "testkitx.Test", "bad input", false)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Kind != ErrorKindValidation {
		t.Fatalf("expected validation kind, got %q", err.Kind)
	}
	if err.Retryable {
		t.Fatal("expected non-retryable error")
	}
	if got := err.Error(); !strings.Contains(got, "validation: testkitx.Test: bad input") {
		t.Fatalf("unexpected error string: %q", got)
	}
}

func TestWrapErrorPreservesCauseAndKind(t *testing.T) {
	t.Parallel()
	cause := context.DeadlineExceeded
	err := WrapError(ErrorKindTimeout, "testkitx.Test", "", true, cause)

	if !IsKind(err, ErrorKindTimeout) {
		t.Fatalf("expected timeout kind, got %v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause, got %v", err)
	}
	if !err.Retryable {
		t.Fatal("expected retryable error")
	}
}

func TestContextErrorClassifiesDeadlineAsRetryableTimeout(t *testing.T) {
	t.Parallel()
	err := contextError("testkitx.Test", context.DeadlineExceeded)
	if !IsKind(err, ErrorKindTimeout) {
		t.Fatalf("expected timeout kind, got %v", err)
	}
	if !err.Retryable {
		t.Fatal("expected deadline errors to be retryable")
	}
}

func TestErrorFormatDefault(t *testing.T) {
	err := NewError(ErrorKindConfig, "Config.Validate", "name is required", false)
	got := fmt.Sprintf("%v", err)
	want := "config: Config.Validate: name is required"
	if got != want {
		t.Fatalf("%%v: got %q, want %q", got, want)
	}
	// %s should produce the same output
	gotS := fmt.Sprintf("%s", err)
	if gotS != want {
		t.Fatalf("%%s: got %q, want %q", gotS, want)
	}
}

func TestErrorFormatVerbose(t *testing.T) {
	cause := errors.New("name must not be empty")
	err := WrapError(ErrorKindConfig, "Config.Validate", "name is required", false, cause)
	got := fmt.Sprintf("%+v", err)

	if !strings.Contains(got, "config: Config.Validate: name is required") {
		t.Fatalf("missing error line: %q", got)
	}
	if !strings.Contains(got, "retryable: false") {
		t.Fatalf("missing retryable line: %q", got)
	}
	if !strings.Contains(got, "cause: name must not be empty") {
		t.Fatalf("missing cause line: %q", got)
	}
}

func TestErrorFormatVerboseRetryable(t *testing.T) {
	err := NewError(ErrorKindTimeout, "svc.Call", "timed out", true)
	got := fmt.Sprintf("%+v", err)

	if !strings.Contains(got, "retryable: true") {
		t.Fatalf("expected retryable true in output: %q", got)
	}
	// No cause line when cause is nil
	if strings.Contains(got, "cause:") {
		t.Fatalf("unexpected cause line in output: %q", got)
	}
}

func TestErrorFormatGoSyntax(t *testing.T) {
	cause := errors.New("inner")
	err := WrapError(ErrorKindValidation, "Op", "msg", true, cause)
	got := fmt.Sprintf("%#v", err)

	if !strings.Contains(got, "&testkitx.Error{") {
		t.Fatalf("missing struct prefix: %q", got)
	}
	if !strings.Contains(got, `Kind:"validation"`) {
		t.Fatalf("missing kind: %q", got)
	}
	if !strings.Contains(got, `Retryable:true`) {
		t.Fatalf("missing retryable: %q", got)
	}
}

func TestErrorFormatNil(t *testing.T) {
	var err *Error
	got := fmt.Sprintf("%v", err)
	if got != "<nil>" {
		t.Fatalf("expected <nil>, got %q", got)
	}
}
