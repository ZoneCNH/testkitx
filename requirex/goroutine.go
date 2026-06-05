package requirex

func NoGoroutineLeak(t TestingT, before, after int) {
	t.Helper()
	if after > before {
		t.Fatalf("goroutine leak detected: before=%d after=%d", before, after)
	}
}
