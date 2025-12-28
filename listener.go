package getlistener

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

var (
	HOST    = "127.0.0.1"
	PORT    = 0
	initErr error
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

func init() {
	envPort := os.Getenv("PORT")
	if envPort != "" {
		selectedPort, err := strconv.Atoi(envPort)
		if err != nil {
			initErr = fmt.Errorf("the environment variable PORT was provided to setup a port but has an invalid value: '%s'", envPort)
			return
		}
		PORT = selectedPort
	}
	envHost := os.Getenv("HOST")
	if envHost != "" {
		HOST = envHost
	}
}
