package getlistener

import (
	"log"
	"net"
	"os"
	"strconv"
)

var (
	HOST = "127.0.0.1"
	PORT = 0
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
			log.Fatalf("the environment variable PORT was provided to setup a port but has an invalid value: '%s'", envPort)
			return
		}
		PORT = selectedPort
	}
	envHost := os.Getenv("HOST")
	if envHost != "" {
		HOST = envHost
		if HOST != "127.0.0.1" && HOST != "localhost" {
			log.Printf("getlistener: 🛡️ SECURITY WARNING: The HOST environment variable is set to '%s', which may expose the service to the network. Please ensure this is intentional.", HOST)
		}
	}
}
