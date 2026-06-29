# Project Conventions

This file documents the conventions and rules for the `go-getlistener` project. Agents must adhere to these guidelines.

## Codebase Map
- `cmd/` -> Executable applications and demos (e.g., `cmd/go-listener-demo/`).
- `/` (root) -> Core library logic (`getlistener` package).
- `dist/` -> Build artifacts (ignored in version control).
- `.github/` -> CI workflows.

## Error Handling
- **Centralized Reporting:** All code paths that handle unexpected errors MUST funnel through `getlistener.ReportError(err error, context string)`.
- **No Direct Panics/Logging:** Never call `panic`, `console.error`, or direct log functions (like `slog.Error`) at the call site for unrecoverable or unexpected errors. Always use the centralized `ReportError` function. In executables, this should be followed by a controlled exit (e.g., `os.Exit(1)`).
- **Error Messages:** Error messages should use formal language and correct grammar. Use "set up" as a verb phrase, "an issue" (not "a issue"), and "cannot handle" (instead of "can't deal with").

## Tooling
- Tasks are executed via `mise` (e.g., `mise run test`).
- Go version is 1.22.
- `golangci-lint` version is 1.60.1.
- Dependencies in `mise.toml` must be pinned to specific minor versions. Do not use "latest" or "lts".

## Platform Specifics
- Platform-specific files (e.g., `listener_unix.go`, `listener_windows.go`) must use explicit build tags (e.g., `//go:build unix`, `//go:build windows`). The `_unix.go` suffix alone is insufficient.
