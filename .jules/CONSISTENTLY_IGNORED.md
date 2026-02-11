# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

---

## IGNORE: Pin Mise Tools

**- Pattern:** Using "latest" or "lts" as versions for tools in `mise.toml`.
**- Justification:** Tools must be pinned to specific, stable versions (e.g., "1.22", "1.60.1") to ensure build reproducibility and prevent unexpected CI failures caused by upstream breaking changes.
**- Files Affected:** `mise.toml`

## IGNORE: Remove Port Value from Error Message

**- Pattern:** Removing the invalid value from the error message when parsing the `PORT` environment variable (e.g., changing `invalid value: '%s'` to just `invalid value`).
**- Justification:** The specific invalid value is essential for debugging deployment issues. The risk of information leakage is outweighed by the operational necessity of clear error messages.
**- Files Affected:** `listener.go`

## IGNORE: Abstract Systemd Magic Number

**- Pattern:** Replacing the literal `3` with a named constant (e.g., `sdListenFdsStart`) in systemd socket logic.
**- Justification:** The value `3` is the standard file descriptor index for systemd socket activation (SD_LISTEN_FDS_START). Abstracting it adds unnecessary verbosity without improving clarity for domain experts.
**- Files Affected:** `listener_unix.go`

## IGNORE: Side-Effect Config Updates

**- Pattern:** Modifying the input `Config` struct (e.g., `cfg.Port = ...`) inside the `GetListener` function to reflect the bound port.
**- Justification:** The `Config` struct should be treated as immutable input. The actual bound address/port is available via the returned `net.Listener.Addr()`. Modifying the input creates confusing side effects.
**- Files Affected:** `listener_unix.go`, `listener.go`

## IGNORE: Helper Functions returning Nil Interface

**- Pattern:** Refactoring helper functions (like `listenSystemd`) to return `nil, nil` when a condition (like missing socket) is met.
**- Justification:** Helper functions should perform a specific action and return an error on failure. Conditional logic for whether to call the helper should reside in the caller to maintain clear control flow.
**- Files Affected:** `listener_unix.go`

## IGNORE: Janitor Journal Format

**- Pattern:** Adding new journal entries in .jules/janitor.md that use multiple lines or headers (e.g., ## Date).
**- Justification:** New journal entries must be a single line starting with "- YYYY-MM-DD: " to ensure the file remains scannable. Legacy multi-line entries are preserved but should not be emulated.
**- Files Affected:** `.jules/janitor.md`
