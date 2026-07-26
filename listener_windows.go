//go:build windows

package getlistener

import (
	"fmt"
	"net"
	"os"
)

// getListenerPlatform creates a standard TCP listener based on the configuration.
//
// On Windows, systemd socket activation is not supported. If LISTEN_PID is set,
// return an error instead of silently falling back to HOST/PORT (matching unix
// fail-closed behavior when activation was requested but cannot be used).
// The error wraps ErrUnsupportedCase so multiplatform callers can use errors.Is.
func getListenerPlatform(cfg *Config) (net.Listener, error) {
	if os.Getenv("LISTEN_PID") != "" {
		return nil, fmt.Errorf("%w: systemd socket activation is not supported on windows (LISTEN_PID is set)", ErrUnsupportedCase)
	}
	return listenTCP(cfg)
}
