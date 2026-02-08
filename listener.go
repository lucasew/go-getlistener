package getlistener

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
)

// Config holds the configuration for the TCP listener.
type Config struct {
	// Host is the IP address or hostname to bind to.
	// Defaults to "127.0.0.1" if not specified.
	Host string
	// Port is the TCP port number to listen on.
	// If 0, a random available port will be chosen.
	Port int
}

// GetAvailablePort returns a random available TCP port number.
//
// Warning: This function is vulnerable to Time-of-Check Time-of-Use (TOCTOU) race conditions.
// The port returned may be claimed by another process before it can be used.
// Prefer using GetListener with Port set to 0 to atomically bind to a random port.
func GetAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// loadConfig reads configuration from environment variables.
//
// It checks:
//   - PORT: The port number to listen on. Defaults to 0 (random port) if not set.
//   - HOST: The hostname or IP to bind to. Defaults to "127.0.0.1" if not set.
//
// It returns an error if PORT contains an invalid integer.
// It logs a security warning if HOST is set to a non-local address.
func loadConfig() (*Config, error) {
	cfg := &Config{
		Host: "127.0.0.1",
		Port: 0,
	}
	envPort := os.Getenv("PORT")
	if envPort != "" {
		selectedPort, err := strconv.Atoi(envPort)
		if err != nil {
			return nil, fmt.Errorf("the environment variable PORT was provided to setup a port but has an invalid value: '%s'", envPort)
		}
		cfg.Port = selectedPort
	}
	envHost := os.Getenv("HOST")
	if envHost != "" {
		cfg.Host = envHost
		if cfg.Host != "127.0.0.1" && cfg.Host != "localhost" {
			slog.Warn(
				"SECURITY WARNING: The HOST environment variable is set to a non-local address, which may expose the service to the network. Please ensure this is intentional.",
				"host", cfg.Host,
			)
		}
	}
	return cfg, nil
}
