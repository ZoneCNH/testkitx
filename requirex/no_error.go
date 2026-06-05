package requirex

type TestingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

func NoError(t TestingT, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
