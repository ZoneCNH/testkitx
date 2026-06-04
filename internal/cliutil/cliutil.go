package cliutil

import (
	"fmt"
	"io"
)

// PrintCLIError writes an error message to w and returns exit code 1.
func PrintCLIError(w io.Writer, err error) int {
	return PrintCLIMessage(w, 1, "ERROR: %v\n", err)
}

// PrintCLIStatus writes a status message to w and returns exit code 0.
func PrintCLIStatus(w io.Writer, format string, args ...any) int {
	return PrintCLIMessage(w, 0, format, args...)
}

// PrintCLIMessage writes a formatted message to w and returns the given exit code.
func PrintCLIMessage(w io.Writer, exitCode int, format string, args ...any) int {
	_, err := fmt.Fprintf(w, format, args...)
	if err != nil {
		return 1
	}
	return exitCode
}
