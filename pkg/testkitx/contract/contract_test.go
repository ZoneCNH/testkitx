package contract_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/contract"
)

func TestAssertHashAndWriteEvidence(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "contract.json")
	if err := os.WriteFile(path, []byte(`{"version":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := contract.FileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}

	metadata := map[string]string{"suite": "unit"}
	evidence := contract.AssertHash(t, "api-contract", path, hash, metadata)
	metadata["suite"] = "mutated"
	if !evidence.Matched || evidence.SHA256 != hash || evidence.Kind != "contract_check" {
		t.Fatalf("unexpected evidence: %+v", evidence)
	}
	if evidence.Metadata["suite"] != "unit" {
		t.Fatalf("expected evidence metadata to be immutable after AssertHash, got %+v", evidence.Metadata)
	}

	out := filepath.Join(t.TempDir(), "evidence.json")
	if err := contract.WriteEvidence(out, evidence); err != nil {
		t.Fatal(err)
	}
	encoded, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var decoded contract.Evidence
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ContractID != "api-contract" || decoded.Metadata["suite"] != "unit" {
		t.Fatalf("unexpected decoded evidence: %+v", decoded)
	}
}

func TestWriteEvidenceValidatesBeforeCreatingFile(t *testing.T) {
	t.Parallel()
	out := filepath.Join(t.TempDir(), "nested", "evidence.json")
	err := contract.WriteEvidence(out, contract.Evidence{Kind: "contract_check"})
	if err == nil || !strings.Contains(err.Error(), "contract_id is required") {
		t.Fatalf("expected contract_id validation failure, got %v", err)
	}
	if _, err := os.Stat(filepath.Dir(out)); !os.IsNotExist(err) {
		t.Fatalf("expected invalid evidence directory not to exist, got %v", err)
	}
}

func TestWriteEvidenceRejectsUnmatchedEvidence(t *testing.T) {
	t.Parallel()
	out := filepath.Join(t.TempDir(), "nested", "evidence.json")
	err := contract.WriteEvidence(out, contract.Evidence{
		Kind:         "contract_check",
		ContractID:   "api-contract",
		ContractPath: "contract.json",
		SHA256:       strings.Repeat("0", 64),
		Matched:      false,
	})
	if err == nil || !strings.Contains(err.Error(), "matched must be true") {
		t.Fatalf("expected matched validation failure, got %v", err)
	}
	if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
		t.Fatalf("expected invalid evidence file not to be created, got %v", statErr)
	}
}

func TestEvidenceValidateRejectsMalformedHash(t *testing.T) {
	t.Parallel()
	evidence := contract.Evidence{
		Kind:         "contract_check",
		ContractID:   "api-contract",
		ContractPath: "contract.json",
		SHA256:       "not-a-sha",
		Matched:      true,
	}
	if err := evidence.Validate(); err == nil || !strings.Contains(err.Error(), "sha256 is invalid") {
		t.Fatalf("expected sha validation failure, got %v", err)
	}
}

func TestWriteEvidenceCreatesDirectory(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "deep", "nested", "evidence.json")
	valid := contract.Evidence{
		Kind:         "contract_check",
		ContractID:   "test",
		ContractPath: "contract.json",
		SHA256:       strings.Repeat("a", 64),
		Matched:      true,
	}
	if err := contract.WriteEvidence(path, valid); err != nil {
		t.Fatalf("WriteEvidence: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty evidence file")
	}
}

func TestEvidenceValidateRejectsNonHexSHA(t *testing.T) {
	t.Parallel()
	e := contract.Evidence{
		Kind:         "contract_check",
		ContractID:   "x",
		ContractPath: "p",
		SHA256:       strings.Repeat("z", 64), // 'z' is not hex
		Matched:      true,
	}
	if err := e.Validate(); err == nil || !strings.Contains(err.Error(), "sha256 is invalid") {
		t.Fatalf("expected non-hex sha error, got %v", err)
	}
}
