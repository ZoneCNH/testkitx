package contract_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/contract"
)

func TestHashFileAndAssertFileHash(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "doc.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := contract.HashFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(hash) != 64 {
		t.Fatalf("expected 64-char hex hash, got %d chars: %s", len(hash), hash)
	}
	contract.AssertFileHash(t, path, hash)
}

func TestHashDirIsDeterministic(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	for _, f := range []struct{ name, body string }{
		{"a.txt", "alpha"},
		{"b.txt", "bravo"},
		{"sub/c.txt", "charlie"},
	} {
		p := filepath.Join(dir, f.name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(f.body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	h1, err := contract.HashDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := contract.HashDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Fatalf("HashDir not deterministic: %s != %s", h1, h2)
	}
}

func TestHashDirDetectsContentChange(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(path, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	h1, err := contract.HashDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("modified"), 0o644); err != nil {
		t.Fatal(err)
	}
	h2, err := contract.HashDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h2 {
		t.Fatal("HashDir should detect content change")
	}
}
