package cliutil

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestPrintCLIError(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	code := PrintCLIError(&buf, errors.New("something broke"))
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(buf.String(), "ERROR: something broke") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestPrintCLIStatus(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	code := PrintCLIStatus(&buf, "all %s", "good")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if buf.String() != "all good" {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestPrintCLIMessage(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	code := PrintCLIMessage(&buf, 42, "msg %d", 7)
	if code != 42 {
		t.Fatalf("expected exit code 42, got %d", code)
	}
	if buf.String() != "msg 7" {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestPrintCLIMessageReturns1OnWriteError(t *testing.T) {
	t.Parallel()
	code := PrintCLIMessage(failWriter{}, 0, "hello")
	if code != 1 {
		t.Fatalf("expected exit code 1 on write error, got %d", code)
	}
}
