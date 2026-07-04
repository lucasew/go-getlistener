package errorhandler

import (
	"log/slog"
	"os"
)

// ReportError logs an unexpected error in a standardized way.
func ReportError(err error, msg string) {
	slog.Error(msg, "error", err)
}

// ReportErrorAndExit logs an unexpected error and exits the program.
func ReportErrorAndExit(err error, msg string) {
	ReportError(err, msg)
	os.Exit(1)
}
