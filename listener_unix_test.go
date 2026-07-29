//go:build unix

package getlistener

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestGetSystemdSocketFD_NotPassed(t *testing.T) {
	t.Setenv("LISTEN_PID", "")
	t.Setenv("LISTEN_FDS", "")

	fd, err := GetSystemdSocketFD()
	if !errors.Is(err, ErrNotPassed) {
		t.Fatalf("err = %v, want ErrNotPassed", err)
	}
	if fd != 0 {
		t.Errorf("fd = %d, want 0", fd)
	}
}

func TestGetSystemdSocketFD_WrongPid(t *testing.T) {
	t.Setenv("LISTEN_PID", "1")
	t.Setenv("LISTEN_FDS", "1")

	fd, err := GetSystemdSocketFD()
	if !errors.Is(err, ErrWrongPid) {
		t.Fatalf("err = %v, want ErrWrongPid", err)
	}
	if fd != 0 {
		t.Errorf("fd = %d, want 0", fd)
	}
}

func TestGetSystemdSocketFD_MissingListenFds(t *testing.T) {
	t.Setenv("LISTEN_PID", fmt.Sprintf("%d", os.Getpid()))
	t.Setenv("LISTEN_FDS", "")

	fd, err := GetSystemdSocketFD()
	if !errors.Is(err, ErrUnsupportedCase) {
		t.Fatalf("err = %v, want ErrUnsupportedCase", err)
	}
	if fd != 0 {
		t.Errorf("fd = %d, want 0", fd)
	}
}

func TestGetSystemdSocketFD_MultipleSockets(t *testing.T) {
	t.Setenv("LISTEN_PID", fmt.Sprintf("%d", os.Getpid()))
	t.Setenv("LISTEN_FDS", "2")

	fd, err := GetSystemdSocketFD()
	if !errors.Is(err, ErrUnsupportedCase) {
		t.Fatalf("err = %v, want ErrUnsupportedCase", err)
	}
	if fd != 0 {
		t.Errorf("fd = %d, want 0", fd)
	}
}

func TestGetSystemdSocketFD_Ok(t *testing.T) {
	t.Setenv("LISTEN_PID", fmt.Sprintf("%d", os.Getpid()))
	t.Setenv("LISTEN_FDS", "1")

	fd, err := GetSystemdSocketFD()
	if err != nil {
		t.Fatalf("GetSystemdSocketFD: %v", err)
	}
	if fd != 3 {
		t.Errorf("fd = %d, want 3 (SD_LISTEN_FDS_START)", fd)
	}
}

func TestGetListener_SystemdWrongPidErrors(t *testing.T) {
	// Wrong LISTEN_PID must not silently fall back to TCP.
	t.Setenv("PORT", "")
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("LISTEN_PID", "1")
	t.Setenv("LISTEN_FDS", "1")

	ln, err := GetListener()
	if ln != nil {
		ln.Close()
		t.Fatal("expected nil listener when systemd PID mismatches")
	}
	if !errors.Is(err, ErrWrongPid) {
		t.Fatalf("err = %v, want ErrWrongPid", err)
	}
}
