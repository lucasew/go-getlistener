package errorhandler

import (
	"log/slog"
	"os"
)

// ReportError centrally logs or reports an unexpected error.
func ReportError(err error) {
	slog.Error("unexpected error", "error", err)
}

// ReportErrorAndExit centrally logs or reports an unexpected error and then cleanly exits.
// This should be used in executables for a controlled exit instead of panicking.
func ReportErrorAndExit(err error) {
	ReportError(err)
	os.Exit(1)
}
