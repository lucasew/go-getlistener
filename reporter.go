package getlistener

import (
	"log/slog"
)

// ReportError is a centralized error-reporting function.
// All code paths that handle unexpected errors MUST funnel through this function.
func ReportError(err error, context string) {
	slog.Error("Unexpected error", "error", err, "context", context)
}
