package requirex

func Equal[T comparable](t TestingT, want, got T) {
	t.Helper()
	if want != got {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
