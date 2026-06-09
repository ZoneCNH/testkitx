package contract_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/contract"
)

func writeTempFileE(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestVerifyFileHashMissingFile(t *testing.T) {
	err := contract.VerifyFileHash("/tmp/nonexistent-12345.bin", "abc123")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestVerifyFileHashMismatch(t *testing.T) {
	path := writeTempFileE(t, "test.txt", "hello world")
	err := contract.VerifyFileHash(path, "0000000000000000000000000000000000000000000000000000000000000000")
	if err == nil {
		t.Fatal("expected hash mismatch error")
	}
}

func TestCheckHashEmptyID(t *testing.T) {
	_, err := contract.CheckHash("", "/tmp/file", "hash", nil)
	if err == nil {
		t.Fatal("expected error for empty contract id")
	}
}

func TestCheckHashMissingFile(t *testing.T) {
	_, err := contract.CheckHash("c1", "/tmp/nonexistent-12345.bin", "hash", nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheckHashMismatch(t *testing.T) {
	path := writeTempFileE(t, "contract.txt", "data")
	_, err := contract.CheckHash("c1", path, "0000000000000000000000000000000000000000000000000000000000000000", nil)
	if err == nil {
		t.Fatal("expected hash mismatch error")
	}
}

func TestCheckHashSuccess(t *testing.T) {
	path := writeTempFileE(t, "contract.txt", "data")
	hash, _ := contract.HashFile(path)
	ev, err := contract.CheckHash("c1", path, hash, map[string]string{"k": "v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ev.Matched {
		t.Fatal("expected match")
	}
	if ev.ContractID != "c1" {
		t.Fatalf("got %q", ev.ContractID)
	}
}
