package getlistener

import (
	"net"
	"testing"
)

func TestGetListener(t *testing.T) {
	// Ensure no environment variables interfere
	t.Setenv("PORT", "")
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("LISTEN_PID", "")
	t.Setenv("LISTEN_FDS", "")

	ln, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener failed: %v", err)
	}
	defer func() {
		_ = ln.Close()
	}()

	addr := ln.Addr().(*net.TCPAddr)
	if addr.Port == 0 {
		t.Errorf("Expected bound port to be non-zero, got 0")
	}
	t.Logf("Listening on %s", addr.String())
}
