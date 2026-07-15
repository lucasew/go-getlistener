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
			return nil, fmt.Errorf("the environment variable PORT was provided to set up a port but has an invalid value: '%s'", envPort)
		}
		if selectedPort < 0 || selectedPort > 65535 {
			return nil, fmt.Errorf("the environment variable PORT was provided to set up a port but has an invalid value: '%s'", envPort)
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

// listenTCP creates a standard TCP listener based on the configuration.
func listenTCP(cfg *Config) (net.Listener, error) {
	if cfg.Port == 0 {
		slog.Info("getlistener: PORT wasn't specified, using random one")
	}
	listenAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	slog.Info("getlistener: listening on", "addr", ln.Addr())
	return ln, nil
}

// GetListener creates a network listener.
//
// It loads the configuration and prioritizes systemd socket activation if available.
// If not, it falls back to creating a standard TCP listener based on the configuration (HOST/PORT).
func GetListener() (net.Listener, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return getListenerPlatform(cfg)
}
