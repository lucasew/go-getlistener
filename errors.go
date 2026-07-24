package getlistener

import "errors"

var (
	// ErrNotPassed is returned when no socket is passed via systemd socket activation.
	ErrNotPassed = errors.New("no socket passed")
	// ErrWrongPid is returned when the socket is passed to a different PID than the current process.
	ErrWrongPid = errors.New("passed the socket to a different PID")
	// ErrUnsupportedCase is returned when the socket activation configuration is unsupported (e.g., multiple sockets).
	ErrUnsupportedCase = errors.New("this case is unsupported")
)
