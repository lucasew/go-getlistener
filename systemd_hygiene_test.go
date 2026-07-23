//go:build unix

package getlistener

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
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

func TestParseSystemdListenFD_DoesNotClearEnv(t *testing.T) {
	// GetListener clears only after FileListener succeeds; parse must be side-effect free.
	t.Setenv("LISTEN_PID", fmt.Sprintf("%d", os.Getpid()))
	t.Setenv("LISTEN_FDS", "1")
	t.Setenv("LISTEN_FDNAMES", "app.socket")

	fd, err := parseSystemdListenFD()
	if err != nil {
		t.Fatalf("parseSystemdListenFD: %v", err)
	}
	if fd != 3 {
		t.Errorf("fd = %d, want 3", fd)
	}
	for _, key := range []string{"LISTEN_PID", "LISTEN_FDS", "LISTEN_FDNAMES"} {
		if got := os.Getenv(key); got == "" {
			t.Errorf("%s cleared by parseSystemdListenFD; want preserved until listen succeeds", key)
		}
	}
}

func TestListenSystemd_NonSocket(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	_, err = listenSystemd(int(r.Fd()))
	if err == nil {
		t.Fatal("listenSystemd: expected error for pipe FD")
	}
	if !strings.Contains(err.Error(), "FileListener") {
		t.Errorf("error should mention FileListener, got: %v", err)
	}
}

func TestListenSystemd_FromTCP(t *testing.T) {
	base, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	f, err := base.(*net.TCPListener).File()
	if err != nil {
		base.Close()
		t.Fatal(err)
	}
	// Close the original listener so only the duplicated FD remains.
	if err := base.Close(); err != nil {
		f.Close()
		t.Fatal(err)
	}
	fd, err := syscall.Dup(int(f.Fd()))
	f.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("LISTEN_FDNAMES", "test.socket")
	ln, err := listenSystemd(fd)
	if err != nil {
		t.Fatalf("listenSystemd: %v", err)
	}
	defer ln.Close()

	addr, ok := ln.Addr().(*net.TCPAddr)
	if !ok || addr.Port == 0 {
		t.Errorf("unexpected addr %v", ln.Addr())
	}
}
