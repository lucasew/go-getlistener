//go:build unix

package getlistener

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
)

var (
	// ErrNotPassed is returned when no socket is passed via systemd socket activation.
	ErrNotPassed       = errors.New("no socket passed")
	// ErrWrongPid is returned when the socket is passed to a different PID than the current process.
	ErrWrongPid        = errors.New("passed the socket to a different PID")
	// ErrUnsupportedCase is returned when the socket activation configuration is unsupported (e.g., multiple sockets).
	ErrUnsupportedCase = errors.New("this case is unsupported")
)

// GetSystemdSocketFD retrieves the file descriptor for the systemd socket.
//
// It validates the environment variables LISTEN_PID and LISTEN_FDS to ensure
// that the socket was intended for this process and that the configuration is supported.
//
// Returns 0 if no socket was passed (ErrNotPassed).
// Returns an error if the configuration is invalid or unsupported.
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
		return 0, fmt.Errorf("%w: LISTEN_PID specified but LISTEN_FDS not, this is an issue in your socket activation mechanism", ErrUnsupportedCase)
	}
	if envListenFds != "1" {
		return 0, fmt.Errorf("%w: this library cannot handle more than one socket being passed", ErrUnsupportedCase)
	}
	return 3, nil
}

// listenSystemd creates a listener from a systemd socket file descriptor.
func listenSystemd(fd int) (net.Listener, error) {
	f := os.NewFile(uintptr(fd), "sd_socket")
	defer func() {
		_ = f.Close()
	}()
	return net.FileListener(f)
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
// It prioritizes systemd socket activation if available.
// If not, it falls back to creating a standard TCP listener based on the configuration (HOST/PORT).
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
