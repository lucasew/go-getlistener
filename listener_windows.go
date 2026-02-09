package getlistener

import (
	"fmt"
	"net"
)

// GetListener creates a standard TCP listener based on the configuration.
//
// On Windows, systemd socket activation is not supported, so this function
// always uses the HOST and PORT environment variables (or defaults).
func GetListener() (net.Listener, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}
