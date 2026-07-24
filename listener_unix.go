//go:build unix

package getlistener

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
)

// clearSystemdActivationEnv drops LISTEN_* so children cannot re-inherit activation state.
func clearSystemdActivationEnv() {
	os.Unsetenv("LISTEN_PID")
	os.Unsetenv("LISTEN_FDS")
	os.Unsetenv("LISTEN_FDNAMES")
}

// parseSystemdListenFD validates LISTEN_PID/LISTEN_FDS and returns the socket FD.
// It does not modify the environment; callers clear after a successful claim.
func parseSystemdListenFD() (int, error) {
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
		return 0, fmt.Errorf("%w: LISTEN_FDS=%q (only exactly one socket is supported)", ErrUnsupportedCase, envListenFds)
	}
	return 3, nil
}

// GetSystemdSocketFD retrieves the file descriptor for the systemd socket.
//
// It validates the environment variables LISTEN_PID and LISTEN_FDS to ensure
// that the socket was intended for this process and that the configuration is supported.
//
// On success it unsets LISTEN_PID, LISTEN_FDS, and LISTEN_FDNAMES so child
// processes do not re-inherit socket activation state (systemd convention).
// Prefer GetListener when building the net.Listener: that path clears activation
// env only after FileListener succeeds, so a failed convert leaves LISTEN_* intact.
//
// Returns 0 if no socket was passed (ErrNotPassed).
// Returns an error if the configuration is invalid or unsupported.
func GetSystemdSocketFD() (int, error) {
	fd, err := parseSystemdListenFD()
	if err != nil {
		return 0, err
	}
	// Direct callers that take the raw FD are treated as having claimed it.
	clearSystemdActivationEnv()
	return fd, nil
}

// listenSystemd creates a listener from a systemd socket file descriptor.
func listenSystemd(fd int) (net.Listener, error) {
	// Systemd passes FDs without CLOEXEC; set it so they are not leaked across exec.
	syscall.CloseOnExec(fd)
	name := os.Getenv("LISTEN_FDNAMES")
	if name == "" {
		name = "sd_socket"
	}
	f := os.NewFile(uintptr(fd), name)
	defer f.Close()
	ln, err := net.FileListener(f)
	if err != nil {
		return nil, fmt.Errorf("systemd socket FileListener: %w", err)
	}
	return ln, nil
}

// getListenerPlatform creates a network listener based on the platform-specific logic.
func getListenerPlatform(cfg *Config) (net.Listener, error) {
	sdSocket, err := parseSystemdListenFD()
	if err != nil && !errors.Is(err, ErrNotPassed) {
		return nil, err
	}
	if sdSocket != 0 {
		ln, err := listenSystemd(sdSocket)
		if err != nil {
			// Leave LISTEN_* set so the failure is diagnosable; FD may already be closed.
			return nil, err
		}
		clearSystemdActivationEnv()
		return ln, nil
	}
	return listenTCP(cfg)
}
