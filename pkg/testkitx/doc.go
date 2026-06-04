// Package testkitx provides L1 test-only helper packages for Go library validation.
//
// It is intended for tests, fixtures, release evidence, boundary checks, and
// deterministic harnesses. Production packages must not import it, and it must
// not depend on real service providers, L2/x.go packages, or production secrets.
package testkitx
