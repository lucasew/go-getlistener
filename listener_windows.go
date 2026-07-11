//go:build windows

package getlistener

import (
	"net"
)

// getListenerPlatform creates a standard TCP listener based on the configuration.
//
// On Windows, systemd socket activation is not supported, so this function
// always uses the HOST and PORT environment variables (or defaults).
func getListenerPlatform(cfg *Config) (net.Listener, error) {
	return listenTCP(cfg)
}
