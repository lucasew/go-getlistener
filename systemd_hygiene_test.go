//go:build unix

package getlistener

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetSystemdSocketFD_ClearsActivationEnv(t *testing.T) {
	t.Setenv("LISTEN_PID", fmt.Sprintf("%d", os.Getpid()))
	t.Setenv("LISTEN_FDS", "1")
	t.Setenv("LISTEN_FDNAMES", "app.socket")

	fd, err := GetSystemdSocketFD()
	if err != nil {
		t.Fatalf("GetSystemdSocketFD: %v", err)
	}
	if fd != 3 {
		t.Errorf("fd = %d, want 3", fd)
	}
	for _, key := range []string{"LISTEN_PID", "LISTEN_FDS", "LISTEN_FDNAMES"} {
		if got := os.Getenv(key); got != "" {
			t.Errorf("%s still set to %q after successful claim", key, got)
		}
	}
}

func TestGetSystemdSocketFD_ZeroFdsMessage(t *testing.T) {
	t.Setenv("LISTEN_PID", fmt.Sprintf("%d", os.Getpid()))
	t.Setenv("LISTEN_FDS", "0")

	fd, err := GetSystemdSocketFD()
	if fd != 0 {
		t.Errorf("fd = %d, want 0", fd)
	}
	if !errors.Is(err, ErrUnsupportedCase) {
		t.Fatalf("err = %v, want ErrUnsupportedCase", err)
	}
	// Zero sockets must not be described as "more than one".
	if strings.Contains(err.Error(), "more than one") {
		t.Errorf("misleading multi-socket wording for LISTEN_FDS=0: %v", err)
	}
	if !strings.Contains(err.Error(), "LISTEN_FDS") || !strings.Contains(err.Error(), "0") {
		t.Errorf("error should mention LISTEN_FDS=0, got: %v", err)
	}
}

func TestGetSystemdSocketFD_KeepsEnvOnError(t *testing.T) {
	// Failed claim must not clear env (caller may inspect / retry policy).
	t.Setenv("LISTEN_PID", "1")
	t.Setenv("LISTEN_FDS", "1")

	_, err := GetSystemdSocketFD()
	if !errors.Is(err, ErrWrongPid) {
		t.Fatalf("err = %v, want ErrWrongPid", err)
	}
	if got := os.Getenv("LISTEN_PID"); got != "1" {
		t.Errorf("LISTEN_PID cleared on error, got %q", got)
	}
	if got := os.Getenv("LISTEN_FDS"); got != "1" {
		t.Errorf("LISTEN_FDS cleared on error, got %q", got)
	}
}
