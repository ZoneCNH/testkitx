package contract_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/contract"
)

func TestAssertHashFileNotFound(t *testing.T) {
	t.Parallel()
	m := &mockTB{}
	contract.AssertHash(m, "test", filepath.Join(t.TempDir(), "nonexistent.json"), "abc", nil)
	if !m.failed {
		t.Fatal("expected failure on missing file")
	}
}

func TestWriteEvidenceWhitespacePath(t *testing.T) {
	t.Parallel()
	valid := contract.Evidence{
		Kind:         "contract_check",
		ContractID:   "test",
		ContractPath: "contract.json",
		SHA256:       strings.Repeat("a", 64),
		Matched:      true,
	}
	err := contract.WriteEvidence("   ", valid)
	if err == nil || !strings.Contains(err.Error(), "evidence path is required") {
		t.Fatalf("expected path required error, got %v", err)
	}
}
