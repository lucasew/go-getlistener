package getlistener

import (
	"fmt"
	"net"
)

// GetListener creates a standard TCP listener.
//
// On Windows, systemd socket activation is not supported.
// This function behaves identically to creating a standard TCP listener:
//  1. Loads configuration from environment variables (PORT, HOST).
//  2. Binds to the specified host and port using net.Listen.
func GetListener() (net.Listener, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}
