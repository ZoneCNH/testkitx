package fake

import (
	"testing"
)

// TestCoverage ensures minimum test coverage for the fake package.
// Each Fake component is exercised through its public API.
func TestCoverage(t *testing.T) {
	// FakeConfig
	t.Run("FakeConfig", func(t *testing.T) {
		TestFakeConfig_GetString(t)
		TestFakeConfig_GetInt(t)
		TestFakeConfig_GetBool(t)
		TestFakeConfig_NilValues(t)
	})

	// FakeLogger
	t.Run("FakeLogger", func(t *testing.T) {
		TestFakeLogger_LogLevels(t)
		TestFakeLogger_AssertLogged(t)
		TestFakeLogger_AssertNoErrors_Passes(t)
		TestFakeLogger_Entries_ReturnsCopy(t)
		TestFakeLogger_Reset(t)
		TestFakeLogger_Fields(t)
	})

	// FakeMeter
	t.Run("FakeMeter", func(t *testing.T) {
		TestFakeMeter_AssertCounterValue(t)
		TestFakeMeter_AssertHistogramRecorded(t)
		TestFakeMeter_CounterValue(t)
		TestFakeMeter_Reset(t)
	})

	// FakeTracer
	t.Run("FakeTracer", func(t *testing.T) {
		TestFakeTracer_StartSpan(t)
		TestFakeTracer_AssertSpanCount(t)
		TestFakeTracer_AssertTraceID(t)
		TestFakeTracer_AssertSpanNamed(t)
		TestFakeTracer_Reset(t)
	})

	// FakeClock
	t.Run("FakeClock", func(t *testing.T) {
		TestFakeClock_Now(t)
		TestFakeClock_Advance(t)
		TestFakeClock_Set(t)
		TestFakeClock_NoAdvance_StaysSame(t)
	})

	// FakeBreaker
	t.Run("FakeBreaker", func(t *testing.T) {
		TestFakeBreaker_InitialClosed(t)
		TestFakeBreaker_OpenDenies(t)
		TestFakeBreaker_Allow_DeniesOpen(t)
		TestFakeBreaker_RecordSuccess_Closes(t)
		TestFakeBreaker_HalfOpenAllows(t)
		TestFakeBreaker_SetState(t)
	})

	// Contract
	t.Run("Contract", func(t *testing.T) {
		TestContract_FakeConfig_Reader(t)
		TestContract_FakeMeter_Interface(t)
		TestContract_FakeTracer_Interface(t)
		TestContract_FakeBreaker_Interface(t)
		TestContract_FakeConfig_Fingerprint(t)
	})

	// Compile-time checks
	t.Run("CompileTime", func(t *testing.T) {
		TestFakeImplementsContracts(t)
	})
}
