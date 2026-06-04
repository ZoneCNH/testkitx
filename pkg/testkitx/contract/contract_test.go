package contract_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/contract"
)

func TestAssertHashAndWriteEvidence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "contract.json")
	if err := os.WriteFile(path, []byte(`{"version":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := contract.FileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}

	evidence := contract.AssertHash(t, "api-contract", path, hash, map[string]string{"suite": "unit"})
	if !evidence.Matched || evidence.SHA256 != hash || evidence.Kind != "contract_check" {
		t.Fatalf("unexpected evidence: %+v", evidence)
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
