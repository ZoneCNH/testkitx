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

func TestErrorNilReceiverError(t *testing.T) {
	var err *Error
	if got := err.Error(); got != "" {
		t.Fatalf("expected empty string for nil Error, got %q", got)
	}
}

func TestErrorUnwrapNilCause(t *testing.T) {
	t.Parallel()
	err := NewError(ErrorKindValidation, "op", "msg", false)
	if err.Unwrap() != nil {
		t.Fatalf("expected nil Unwrap for error without cause")
	}
}

func TestErrorUnwrapNilReceiver(t *testing.T) {
	var err *Error
	if err.Unwrap() != nil {
		t.Fatalf("expected nil Unwrap for nil Error")
	}
}

func TestIsKindWithNonError(t *testing.T) {
	t.Parallel()
	if IsKind(fmt.Errorf("plain error"), ErrorKindValidation) {
		t.Fatal("expected false for non-Error type")
	}
}

func TestErrorKindWithNonError(t *testing.T) {
	t.Parallel()
	kind := errorKind(fmt.Errorf("plain"))
	if kind != ErrorKindInternal {
		t.Fatalf("expected internal kind for non-Error, got %q", kind)
	}
}

func TestErrorWithMessageAndCause(t *testing.T) {
	t.Parallel()
	cause := fmt.Errorf("root cause")
	err := WrapError(ErrorKindConnection, "svc.Connect", "connection failed", true, cause)
	if got := err.Error(); !strings.Contains(got, "connection failed") {
		t.Fatalf("expected message in error string, got %q", got)
	}
	if !errors.Is(err, cause) {
		t.Fatal("expected cause to be unwrappable")
	}
}

func TestErrorWithOnlyCause(t *testing.T) {
	t.Parallel()
	cause := fmt.Errorf("root cause")
	err := WrapError(ErrorKindInternal, "op", "", false, cause)
	if got := err.Error(); !strings.Contains(got, "root cause") {
		t.Fatalf("expected cause message in error string, got %q", got)
	}
}

func TestErrorEmptyMessageWithCause(t *testing.T) {
	t.Parallel()
	cause := fmt.Errorf("underlying problem")
	err := WrapError(ErrorKindConnection, "svc.Call", "", true, cause)
	if got := err.Error(); !strings.Contains(got, "underlying problem") {
		t.Fatalf("expected cause text in error, got %q", got)
	}
	if err.Message != "underlying problem" {
		t.Fatalf("expected message to be set from cause, got %q", err.Message)
	}
}

func TestErrorNoOpNoMessage(t *testing.T) {
	t.Parallel()
	err := NewError(ErrorKindInternal, "", "", false)
	got := err.Error()
	if got != string(ErrorKindInternal) {
		t.Fatalf("expected just kind, got %q", got)
	}
}

func TestErrorWithOpOnly(t *testing.T) {
	t.Parallel()
	err := NewError(ErrorKindAuth, "login", "", false)
	got := err.Error()
	if got != "auth: login" {
		t.Fatalf("expected kind:op, got %q", got)
	}
}

func TestValidationErrorHelper(t *testing.T) {
	t.Parallel()
	err := validationError("validate.input", "bad value", nil)
	if err.Kind != ErrorKindValidation {
		t.Fatalf("expected validation kind, got %q", err.Kind)
	}
	if err.Retryable {
		t.Fatal("expected non-retryable")
	}
}

func TestContextErrorNonDeadline(t *testing.T) {
	t.Parallel()
	err := contextError("op", context.Canceled)
	if err.Kind != ErrorKindUnavailable {
		t.Fatalf("expected unavailable kind, got %q", err.Kind)
	}
	if err.Retryable {
		t.Fatal("expected non-retryable for canceled")
	}
}

func TestErrorNilReceiverReturnsEmpty(t *testing.T) {
	t.Parallel()
	var e *Error
	if got := e.Error(); got != "" {
		t.Fatalf("expected empty string for nil receiver, got %q", got)
	}
}

func TestErrorMessageEmptyWithCause(t *testing.T) {
	t.Parallel()
	e := WrapError(ErrorKindConnection, "op", "", false, fmt.Errorf("underlying"))
	if e.Error() != "connection: op: underlying" {
		t.Fatalf("unexpected error string: %q", e.Error())
	}
}

func TestErrorMessageEmptyWithCauseDirect(t *testing.T) {
	t.Parallel()
	// Construct Error directly to bypass newError which copies cause to message.
	e := &Error{Kind: ErrorKindConnection, Op: "op", Cause: fmt.Errorf("underlying")}
	got := e.Error()
	if !strings.Contains(got, "underlying") {
		t.Fatalf("expected cause in error string, got %q", got)
	}
}
