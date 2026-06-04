package testkitx

import (
	"context"
	"errors"
	"fmt"
)

type ErrorKind string

const (
	ErrorKindConfig      ErrorKind = "config"
	ErrorKindValidation  ErrorKind = "validation"
	ErrorKindConnection  ErrorKind = "connection"
	ErrorKindUnavailable ErrorKind = "unavailable"
	ErrorKindTimeout     ErrorKind = "timeout"
	ErrorKindAuth        ErrorKind = "auth"
	ErrorKindConflict    ErrorKind = "conflict"
	ErrorKindRateLimit   ErrorKind = "rate_limit"
	ErrorKindInternal    ErrorKind = "internal"
)

type Error struct {
	Kind      ErrorKind
	Op        string
	Message   string
	Cause     error
	Retryable bool
}

func NewError(kind ErrorKind, op string, message string, retryable bool) *Error {
	return newError(kind, op, message, retryable, nil)
}

func WrapError(kind ErrorKind, op string, message string, retryable bool, cause error) *Error {
	return newError(kind, op, message, retryable, cause)
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	message := string(e.Kind)
	if e.Op != "" {
		message += ": " + e.Op
	}
	if e.Message != "" {
		message += ": " + e.Message
	}
	if e.Message == "" && e.Cause != nil {
		message += ": " + e.Cause.Error()
	}
	return message
}

// Format implements fmt.Formatter.
// %v, %s: same as Error() output (kind: op: message)
// %+v: full detail with cause chain and retryable flag
// %#v: Go syntax representation
func (e *Error) Format(f fmt.State, c rune) {
	if e == nil {
		_, _ = fmt.Fprint(f, "<nil>")
		return
	}
	switch {
	case c == 'v' && f.Flag('#'):
		_, _ = fmt.Fprintf(f, "&testkitx.Error{Kind:%q, Op:%q, Message:%q, Cause:%#v, Retryable:%t}",
			e.Kind, e.Op, e.Message, e.Cause, e.Retryable)
	case c == 'v' && f.Flag('+'):
		_, _ = fmt.Fprintf(f, "%s", e.Error())
		_, _ = fmt.Fprintf(f, "\n    retryable: %t", e.Retryable)
		if e.Cause != nil {
			_, _ = fmt.Fprintf(f, "\n    cause: %s", e.Cause.Error())
		}
	default:
		_, _ = fmt.Fprint(f, e.Error())
	}
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func IsKind(err error, kind ErrorKind) bool {
	var target *Error
	if errors.As(err, &target) {
		return target.Kind == kind
	}
	return false
}

func newError(kind ErrorKind, op string, message string, retryable bool, cause error) *Error {
	if message == "" && cause != nil {
		message = cause.Error()
	}
	return &Error{
		Kind:      kind,
		Op:        op,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
	}
}

func validationError(op string, message string, cause error) *Error {
	return newError(ErrorKindValidation, op, message, false, cause)
}

func contextError(op string, cause error) *Error {
	kind := ErrorKindUnavailable
	retryable := false
	if errors.Is(cause, context.DeadlineExceeded) {
		kind = ErrorKindTimeout
		retryable = true
	}
	return newError(kind, op, "", retryable, cause)
}

func errorKind(err error) ErrorKind {
	var target *Error
	if errors.As(err, &target) {
		return target.Kind
	}
	return ErrorKindInternal
}
