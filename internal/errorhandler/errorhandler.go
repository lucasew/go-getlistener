package errorhandler

import (
	"log/slog"
	"os"
)

// ReportError provides centralized error reporting.
// All code paths handling unexpected errors MUST funnel through this function.
func ReportError(err error, context string) {
	slog.Error("Unexpected error occurred", "context", context, "error", err)
}

// ReportErrorAndExit provides centralized error reporting and exits the application.
// Useful for executables that cannot recover from an error.
func ReportErrorAndExit(err error, context string, code int) {
	ReportError(err, context)
	os.Exit(code)
}
