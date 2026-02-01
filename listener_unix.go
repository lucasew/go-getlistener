//go:build unix

package getlistener

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	ErrNotPassed       = errors.New("no socket passed")
	ErrWrongPid        = errors.New("passed the socket to a different PID")
	ErrUnsupportedCase = errors.New("this case is unsupported")
)

// GetSystemdSocketFD gets the systemd socket fd, gives 0 if not passed, error if passed wrong
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

func listenSystemd() (net.Listener, error) {
	sdSocket, err := GetSystemdSocketFD()
	if err != nil {
		if errors.Is(err, ErrNotPassed) {
			return nil, nil
		}
		return nil, err
	}
	if sdSocket != 0 {
		f := os.NewFile(uintptr(sdSocket), "sd_socket")
		log.Printf("getlistener: using socket activation on fd %d", sdSocket)
		return net.FileListener(f)
	}
	return nil, nil
}

func listenTCP(cfg *Config) (net.Listener, error) {
	if cfg.Port == 0 {
		log.Printf("getlistener: PORT wasn't specified, using random one")
		selectedPort, err := GetAvailablePort()
		if err != nil {
			return nil, err
		}
		cfg.Port = selectedPort
	}
	listenAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("getlistener: listening on %s", listenAddr)
	return net.Listen("tcp", listenAddr)
}

// GetListener returns a listener that is either a systemd socket or a TCP listener
func GetListener() (net.Listener, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	l, err := listenSystemd()
	if err != nil {
		return nil, err
	}
	if l != nil {
		return l, nil
	}

	return listenTCP(cfg)
}
