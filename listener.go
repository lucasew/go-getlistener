package getlistener

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
)

// GetAvailablePort get the number of an available port
func GetAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func getConfig() (string, int, error) {
	host := "127.0.0.1"
	port := 0
	envPort := os.Getenv("PORT")
	if envPort != "" {
		selectedPort, err := strconv.Atoi(envPort)
		if err != nil {
			return "", 0, fmt.Errorf("the environment variable PORT was provided to setup a port but has an invalid value: '%s'", envPort)
		}
		port = selectedPort
	}
	envHost := os.Getenv("HOST")
	if envHost != "" {
		host = envHost
		if host != "127.0.0.1" && host != "localhost" {
			slog.Warn(
				"SECURITY WARNING: The HOST environment variable is set to a non-local address, which may expose the service to the network. Please ensure this is intentional.",
				"host", host,
			)
		}
	}
	return host, port, nil
}
