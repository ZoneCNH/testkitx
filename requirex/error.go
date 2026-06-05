package requirex

import (
	"errors"
	"reflect"
)

type kindError interface {
	ErrorKind() string
}

type kindStringer interface {
	Kind() string
}

func ErrorKind(t TestingT, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error kind %q, got nil", want)
	}
	got, ok := errorKind(err)
	if !ok {
		t.Fatalf("expected error kind %q, got untyped error %T %[2]v", want, err)
	}
	if got != want {
		t.Fatalf("expected error kind %q, got %q", want, got)
	}
}

func ErrorKindOneOf(t TestingT, err error, wants ...string) {
	t.Helper()
	if len(wants) == 0 {
		t.Fatalf("expected at least one wanted error kind")
	}
	if err == nil {
		t.Fatalf("expected one of error kinds %v, got nil", wants)
	}
	got, ok := errorKind(err)
	if !ok {
		t.Fatalf("expected one of error kinds %v, got untyped error %T %[2]v", wants, err)
	}
	for _, want := range wants {
		if got == want {
			return
		}
	}
	t.Fatalf("expected one of error kinds %v, got %q", wants, got)
}

func errorKind(err error) (string, bool) {
	var withErrorKind kindError
	if errors.As(err, &withErrorKind) {
		return withErrorKind.ErrorKind(), true
	}
	var withKind kindStringer
	if errors.As(err, &withKind) {
		return withKind.Kind(), true
	}
	value := reflect.ValueOf(err)
	if !value.IsValid() || value.Kind() != reflect.Pointer || value.IsNil() {
		return "", false
	}
	elem := value.Elem()
	if elem.Kind() != reflect.Struct {
		return "", false
	}
	field := elem.FieldByName("Kind")
	if field.IsValid() && field.Kind() == reflect.String {
		return field.String(), true
	}
	return "", false
}
