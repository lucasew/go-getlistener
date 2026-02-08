package getlistener

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	// ErrNotPassed indicates that no systemd socket was passed to the process (LISTEN_PID not set).
	ErrNotPassed = errors.New("no socket passed")
	// ErrWrongPid indicates that the socket was passed to a different process ID.
	ErrWrongPid = errors.New("passed the socket to a different PID")
	// ErrUnsupportedCase indicates that the systemd socket configuration is not supported (e.g. multiple sockets).
	ErrUnsupportedCase = errors.New("this case is unsupported")
)

// GetSystemdSocketFD retrieves the file descriptor for the systemd-activated socket.
//
// It checks environment variables set by systemd:
//   - LISTEN_PID: Must match the current process ID.
//   - LISTEN_FDS: Must be "1" (only one socket is supported).
//
// It returns:
//   - 3 (SD_LISTEN_FDS_START) if a valid socket is passed.
//   - 0 and ErrNotPassed if LISTEN_PID is missing.
//   - 0 and an error wrapping ErrWrongPid or ErrUnsupportedCase if validation fails.
func GetSystemdSocketFD() (int, error) {
	envListenPid := os.Getenv("LISTEN_PID")
	if envListenPid == "" {
		return 0, ErrNotPassed
	}
	if envListenPid != fmt.Sprintf("%d", os.Getpid()) {
		return 0, fmt.Errorf("%w: %s instead of %d", ErrWrongPid, envListenPid, os.Getpid())
	}
	envListenFds := os.Getenv("LISTEN_FDS")
	if envListenFds == "" {
		return 0, fmt.Errorf("%w: LISTEN_PID specified but LISTEN_FDS not, this is a issue in your socket activation mechanism", ErrUnsupportedCase)
	}
	if envListenFds != "1" {
		return 0, fmt.Errorf("%w: this library can't deal with more than one socket being passed", ErrUnsupportedCase)
	}
	return 3, nil
}

// listenSystemd creates a listener from a systemd socket file descriptor.
// It wraps the file descriptor in a os.File and then converts it to a net.Listener.
func listenSystemd(fd int) (net.Listener, error) {
	f := os.NewFile(uintptr(fd), "sd_socket")
	defer f.Close()
	return net.FileListener(f)
}

// listenTCP creates a standard TCP listener based on the configuration.
// It uses net.Listen with the host and port specified in the config.
func listenTCP(cfg *Config) (net.Listener, error) {
	if cfg.Port == 0 {
		log.Printf("getlistener: PORT wasn't specified, using random one")
	}
	listenAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	log.Printf("getlistener: listening on %s", ln.Addr())
	return ln, nil
}

// GetListener creates a net.Listener, preferring systemd socket activation.
//
// It follows this logic:
//  1. Loads configuration from environment variables (PORT, HOST).
//     Validation of PORT occurs here, even if systemd socket is used.
//  2. Checks if a systemd socket is available via GetSystemdSocketFD.
//  3. If a systemd socket is provided (and valid), it returns a listener for that socket.
//  4. If no systemd socket is provided, it falls back to creating a standard TCP listener
//     using the loaded configuration.
//
// It returns an error if configuration loading fails or if systemd socket validation fails
// (other than ErrNotPassed).
func GetListener() (net.Listener, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	sdSocket, err := GetSystemdSocketFD()
	if err != nil && !errors.Is(err, ErrNotPassed) {
		return nil, err
	}
	if sdSocket != 0 {
		return listenSystemd(sdSocket)
	}
	return listenTCP(cfg)
}
