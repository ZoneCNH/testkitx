// Package fake provides deterministic fake implementations of FoundationX
// L1 interfaces for use in tests. Every fake implements its corresponding
// interface and includes a compile-time assertion.
//
// All fakes are goroutine-safe and suitable for use with -race.
package fake
