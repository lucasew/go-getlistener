# Janitor's Journal

This journal is a record of CRITICAL refactoring learnings for this codebase. Its purpose is to prevent repeating past mistakes and to document architectural decisions.

## 2024-07-25 - `init()` functions should not terminate the application
**Issue:** The `init()` function in `listener.go` called `log.Fatalf` when it failed to parse the `PORT` environment variable. This is a critical anti-pattern in a library because it can unexpectedly terminate the entire host application.
**Root Cause:** The original code prioritized immediate feedback for a configuration error over library safety, not considering that libraries should not make termination decisions for the consumer.
**Solution:** The `init()` function was refactored to store any parsing errors in a package-level variable. The public `GetListener()` function now checks this variable upon being called and returns the error, allowing the application to handle it gracefully.
**Pattern:** Libraries must not call `os.Exit()` or `log.Fatalf()`. Initialization or configuration errors should be captured and exposed as return values from the library's public functions.

## 2026-01-19 - Avoid TOCTOU race conditions in port selection
**Issue:** The library used a check-then-act pattern to find an available port (`GetAvailablePort`), which involved listening on a random port, closing it, and then re-binding to it. This created a race condition where another process could steal the port in between.
**Root Cause:** A desire to "know" the port number before confirming the final listener was created, likely for logging purposes.
**Solution:** Refactored `GetListener` to bind directly to port 0 (`:0`), allowing the OS to atomically assign a free port. The port number is then retrieved from the active listener.
**Pattern:** Always prefer atomic operations over check-then-act patterns for resource acquisition. Use `net.Listen("tcp", ":0")` to safely bind to a random port.
