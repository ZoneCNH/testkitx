package fake

import (
	"fmt"
	"strings"
	"sync"
)

// LogLevel represents the severity of a log entry.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return fmt.Sprintf("LogLevel(%d)", l)
	}
}

// LogEntry records a single log event.
type LogEntry struct {
	Level   LogLevel
	Message string
	Fields  map[string]any
}

// Logger mirrors the observex.Logger interface consumed by FoundationX modules.
type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
}

// FakeLoggerImpl is a deterministic fake logger that records entries and
// supports post-hoc assertions. It implements Logger.
type FakeLoggerImpl struct {
	mu      sync.Mutex
	entries []LogEntry
}

// Compile-time contract: *FakeLoggerImpl implements Logger.
var _ Logger = (*FakeLoggerImpl)(nil)

// FakeLogger creates a new deterministic fake logger.
func FakeLogger() *FakeLoggerImpl {
	return &FakeLoggerImpl{}
}

// Debug records a debug-level log entry.
func (l *FakeLoggerImpl) Debug(msg string, fields ...any) {
	l.append(LevelDebug, msg, fieldsToMap(fields...))
}

// Info records an info-level log entry.
func (l *FakeLoggerImpl) Info(msg string, fields ...any) {
	l.append(LevelInfo, msg, fieldsToMap(fields...))
}

// Warn records a warn-level log entry.
func (l *FakeLoggerImpl) Warn(msg string, fields ...any) {
	l.append(LevelWarn, msg, fieldsToMap(fields...))
}

// Error records an error-level log entry.
func (l *FakeLoggerImpl) Error(msg string, fields ...any) {
	l.append(LevelError, msg, fieldsToMap(fields...))
}

// AssertLogged fails t if no log entry at the given level contains substr.
// substr is matched case-insensitively against the message.
func (l *FakeLoggerImpl) AssertLogged(t T, level LogLevel, contains string) {
	t.Helper()
	l.mu.Lock()
	defer l.mu.Unlock()
	contains = strings.ToLower(contains)
	for _, e := range l.entries {
		if e.Level == level && strings.Contains(strings.ToLower(e.Message), contains) {
			return
		}
	}
	t.Errorf("expected log entry [%s] containing %q, but none found (total entries: %d)",
		level, contains, len(l.entries))
}

// AssertNoErrors fails t if any Error-level entries were recorded.
func (l *FakeLoggerImpl) AssertNoErrors(t T) {
	t.Helper()
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, e := range l.entries {
		if e.Level == LevelError {
			t.Errorf("unexpected error log: %s", e.Message)
		}
	}
}

// Entries returns a snapshot of all recorded log entries.
func (l *FakeLoggerImpl) Entries() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]LogEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Reset clears all recorded entries.
func (l *FakeLoggerImpl) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = nil
}

func (l *FakeLoggerImpl) append(level LogLevel, msg string, fields map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, LogEntry{
		Level:   level,
		Message: msg,
		Fields:  fields,
	})
}

func fieldsToMap(kv ...any) map[string]any {
	if len(kv) == 0 {
		return nil
	}
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			key = fmt.Sprint(kv[i])
		}
		m[key] = kv[i+1]
	}
	return m
}

// T is the subset of testing.TB used by assertions.
type T interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}
