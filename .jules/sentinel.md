# Sentinel's Journal

This journal is a record of CRITICAL security learnings for this codebase. Its purpose is to prevent repeating past mistakes and to document architectural decisions.

## 2026-01-13 - Harden PORT parsing error message
**Vulnerability:** Information Leakage
**Learning:** The `init()` function in `listener.go` would expose the value of the `PORT` environment variable in an error message if it failed to parse. This could leak potentially sensitive configuration details to an attacker who could trigger the error.
**Prevention:** Error messages should be generic and not include sensitive data. Instead of returning the invalid value, the error now simply states that the value was invalid. This prevents information leakage while still providing enough information for a developer to debug the issue.
