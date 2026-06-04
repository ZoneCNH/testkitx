package main

import (
	"testing"
)

func TestToolVersionReportsMissingBinary(t *testing.T) {
	t.Parallel()
	got := toolVersion("definitely-missing-releasemanifest-test-binary")
	if got != "missing" {
		t.Fatalf("toolVersion missing binary = %q, want missing", got)
	}
}
