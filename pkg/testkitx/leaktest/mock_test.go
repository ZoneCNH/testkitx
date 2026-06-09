package leaktest_test

import (
	"context"
	"io"
	"os"
	"testing"
)

type mockTB struct {
	testing.TB
	failed bool
}

func (m *mockTB) Helper()                           {}
func (m *mockTB) Fatalf(format string, args ...any)  { m.failed = true }
func (m *mockTB) Errorf(format string, args ...any)  { m.failed = true }
func (m *mockTB) FailNow()                           { m.failed = true }
func (m *mockTB) Failed() bool                       { return m.failed }
func (m *mockTB) Name() string                       { return "mock" }
func (m *mockTB) Log(args ...any)                    {}
func (m *mockTB) Logf(format string, args ...any)    {}
func (m *mockTB) Skip(args ...any)                   {}
func (m *mockTB) Skipf(format string, args ...any)   {}
func (m *mockTB) SkipNow()                           {}
func (m *mockTB) Skipped() bool                      { return false }
func (m *mockTB) TempDir() string                    { return os.TempDir() }
func (m *mockTB) Setenv(key, value string)           {}
func (m *mockTB) Cleanup(func())                     {}
func (m *mockTB) Error(args ...any)                  { m.failed = true }
func (m *mockTB) Fatal(args ...any)                  { m.failed = true }
func (m *mockTB) Fail()                              { m.failed = true }
func (m *mockTB) ArtifactDir() string                { return os.TempDir() }
func (m *mockTB) Attr(key, value string)             {}
func (m *mockTB) Chdir(dir string)                   {}
func (m *mockTB) Context() context.Context           { return context.Background() }
func (m *mockTB) Output() io.Writer                  { return io.Discard }
