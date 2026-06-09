package assertx

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// NotEqual asserts that want and got are NOT deeply equal.
func NotEqual[T any](t testing.TB, want, got T) {
	t.Helper()
	if reflect.DeepEqual(want, got) {
		t.Fatalf("expected values to differ, both are: %#v", want)
	}
}

// HasError asserts that err is non-nil.
func HasError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

// Contains asserts that haystack contains needle as a substring.
func Contains(t testing.TB, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected %q to contain %q", haystack, needle)
	}
}

// Len asserts that collection has the expected length.
func Len(t testing.TB, collection any, expected int) {
	t.Helper()
	v := reflect.ValueOf(collection)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		if v.Len() != expected {
			t.Fatalf("expected length %d, got %d", expected, v.Len())
		}
	default:
		t.Fatalf("Len: unsupported type %T", collection)
	}
}

// Sprintf is a helper for building assertion messages.
func Sprintf(format string, args ...any) string { return fmt.Sprintf(format, args...) }
