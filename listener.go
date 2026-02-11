// Package getlistener provides a way to get a listener respecting systemd socket activation.
package getlistener

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
)

// Config holds the network configuration for the listener.
type Config struct {
	// Host is the hostname or IP address to listen on.
	// Defaults to "127.0.0.1" if not specified.
	Host string
	// Port is the port number to listen on.
	// If 0, a random available port will be chosen.
	Port int
}

// GetAvailablePort returns the number of an available TCP port.
//
// WARNING: This function is vulnerable to Time-of-Check Time-of-Use (TOCTOU) race conditions.
// The port returned may be claimed by another process between the time it is released by this function
// and the time it is used by the caller. It is recommended to let `net.Listen` choose a port by specifying port 0,
// rather than using this function.
func GetAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = listener.Close()
	}()
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// loadConfig loads the configuration from environment variables.
//
// It checks for the following environment variables:
// - PORT: The port to listen on. If invalid, an error is returned.
// - HOST: The host to listen on. Defaults to "127.0.0.1".
//
// If HOST is set to a non-local address, a security warning is logged.
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
