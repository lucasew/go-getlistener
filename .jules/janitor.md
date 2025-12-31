# Janitor's Journal

This journal is a record of CRITICAL refactoring learnings for this codebase. Its purpose is to prevent repeating past mistakes and to document architectural decisions.

## 2024-07-25 - `init()` functions should not terminate the application
**Issue:** The `init()` function in `listener.go` called `log.Fatalf` when it failed to parse the `PORT` environment variable. This is a critical anti-pattern in a library because it can unexpectedly terminate the entire host application.
**Root Cause:** The original code prioritized immediate feedback for a configuration error over library safety, not considering that libraries should not make termination decisions for the consumer.
**Solution:** The `init()` function was refactored to store any parsing errors in a package-level variable. The public `GetListener()` function now checks this variable upon being called and returns the error, allowing the application to handle it gracefully.
**Pattern:** Libraries must not call `os.Exit()` or `log.Fatalf()`. Initialization or configuration errors should be captured and exposed as return values from the library's public functions.
