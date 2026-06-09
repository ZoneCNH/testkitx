package contract_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/contract"
)

// mockTB implements testing.TB for Go 1.26 without calling runtime.Goexit on Fatalf.
type mockTB struct {
	testing.TB
	failed bool
}

func (m *mockTB) Helper()                              {}
func (m *mockTB) Fatalf(format string, args ...any)    { m.failed = true }
func (m *mockTB) Errorf(format string, args ...any)    { m.failed = true }
func (m *mockTB) FailNow()                             { m.failed = true }
func (m *mockTB) Failed() bool                         { return m.failed }
func (m *mockTB) Name() string                         { return "mock" }
func (m *mockTB) Log(args ...any)                      {}
func (m *mockTB) Logf(format string, args ...any)      {}
func (m *mockTB) Skip(args ...any)                     {}
func (m *mockTB) Skipf(format string, args ...any)     {}
func (m *mockTB) SkipNow()                             {}
func (m *mockTB) Skipped() bool                        { return false }
func (m *mockTB) TempDir() string                      { return os.TempDir() }
func (m *mockTB) Setenv(key, value string)             {}
func (m *mockTB) Cleanup(func())                       {}
func (m *mockTB) Error(args ...any)                    { m.failed = true }
func (m *mockTB) Fatal(args ...any)                    { m.failed = true }
func (m *mockTB) Fail()                                { m.failed = true }
func (m *mockTB) ArtifactDir() string                  { return os.TempDir() }
func (m *mockTB) Attr(key, value string)               {}
func (m *mockTB) Chdir(dir string)                     {}
func (m *mockTB) Context() context.Context             { return context.Background() }
func (m *mockTB) Output() io.Writer                    { return io.Discard }

func TestEvidenceValidateRejectsEmptyKind(t *testing.T) {
	t.Parallel()
	e := contract.Evidence{ContractID: "x", ContractPath: "p", SHA256: strings.Repeat("a", 64), Matched: true}
	if err := e.Validate(); err == nil || !strings.Contains(err.Error(), "kind is required") {
		t.Fatalf("expected kind error, got %v", err)
	}
}

func TestEvidenceValidateRejectsWrongKind(t *testing.T) {
	t.Parallel()
	e := contract.Evidence{Kind: "wrong", ContractID: "x", ContractPath: "p", SHA256: strings.Repeat("a", 64), Matched: true}
	if err := e.Validate(); err == nil || !strings.Contains(err.Error(), "kind must be contract_check") {
		t.Fatalf("expected kind error, got %v", err)
	}
}

func TestEvidenceValidateRejectsEmptyContractID(t *testing.T) {
	t.Parallel()
	e := contract.Evidence{Kind: "contract_check", ContractPath: "p", SHA256: strings.Repeat("a", 64), Matched: true}
	if err := e.Validate(); err == nil || !strings.Contains(err.Error(), "contract_id is required") {
		t.Fatalf("expected contract_id error, got %v", err)
	}
}

func TestEvidenceValidateRejectsEmptyPath(t *testing.T) {
	t.Parallel()
	e := contract.Evidence{Kind: "contract_check", ContractID: "x", SHA256: strings.Repeat("a", 64), Matched: true}
	if err := e.Validate(); err == nil || !strings.Contains(err.Error(), "contract_path is required") {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestEvidenceValidateRejectsEmptySHA(t *testing.T) {
	t.Parallel()
	e := contract.Evidence{Kind: "contract_check", ContractID: "x", ContractPath: "p", Matched: true}
	if err := e.Validate(); err == nil || !strings.Contains(err.Error(), "sha256 is required") {
		t.Fatalf("expected sha256 error, got %v", err)
	}
}

func TestEvidenceValidateRejectsShortSHA(t *testing.T) {
	t.Parallel()
	e := contract.Evidence{Kind: "contract_check", ContractID: "x", ContractPath: "p", SHA256: "abc", Matched: true}
	if err := e.Validate(); err == nil || !strings.Contains(err.Error(), "sha256 is invalid") {
		t.Fatalf("expected sha256 length error, got %v", err)
	}
}

func TestWriteEvidenceRejectsEmptyPath(t *testing.T) {
	t.Parallel()
	err := contract.WriteEvidence("", contract.Evidence{Kind: "contract_check", ContractID: "x", ContractPath: "p", SHA256: strings.Repeat("a", 64), Matched: true})
	if err == nil || !strings.Contains(err.Error(), "evidence path is required") {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestWriteEvidenceRejectsInvalidEvidence(t *testing.T) {
	t.Parallel()
	err := contract.WriteEvidence("/tmp/test.json", contract.Evidence{})
	if err == nil {
		t.Fatal("expected error for invalid evidence")
	}
}

func TestAssertHashRejectsEmptyContractID(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(path, []byte("data"), 0o644)
	hash, _ := contract.FileSHA256(path)
	m := &mockTB{}
	contract.AssertHash(m, "", path, hash, nil)
	if !m.failed {
		t.Fatal("expected failure on empty contract ID")
	}
}

func TestAssertHashMismatch(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(path, []byte("data"), 0o644)
	m := &mockTB{}
	contract.AssertHash(m, "id", path, strings.Repeat("0", 64), nil)
	if !m.failed {
		t.Fatal("expected failure on hash mismatch")
	}
}

func TestFileSHA256NonExistent(t *testing.T) {
	t.Parallel()
	_, err := contract.FileSHA256("/nonexistent/file.txt")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestHashDirNonExistent(t *testing.T) {
	t.Parallel()
	_, err := contract.HashDir("/nonexistent/dir")
	if err == nil {
		t.Fatal("expected error for nonexistent dir")
	}
}

func TestCopyMetadataNil(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "c.json")
	os.WriteFile(path, []byte(`{}`), 0o644)
	hash, _ := contract.FileSHA256(path)
	evidence := contract.AssertHash(t, "id", path, hash, nil)
	if evidence.Metadata != nil {
		t.Fatalf("expected nil metadata, got %v", evidence.Metadata)
	}
}

func TestCopyMetadataNonNil(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "c.json")
	os.WriteFile(path, []byte(`{}`), 0o644)
	hash, _ := contract.FileSHA256(path)
	evidence := contract.AssertHash(t, "id", path, hash, map[string]string{"k": "v"})
	if evidence.Metadata["k"] != "v" {
		t.Fatalf("expected metadata k=v, got %v", evidence.Metadata)
	}
}
