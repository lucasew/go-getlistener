package getlistener

import (
	"net"
	"testing"
)

func TestGetListener_TCP(t *testing.T) {
	// Ensure no systemd variables are affecting the test
	t.Setenv("LISTEN_PID", "")
	t.Setenv("LISTEN_FDS", "")
	t.Setenv("PORT", "0") // Use random port

	// Call GetListener
	l, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener() error = %v", err)
	}
	defer l.Close()

	// Verify it's a TCP listener
	addr := l.Addr()
	if addr.Network() != "tcp" {
		t.Errorf("expected tcp network, got %s", addr.Network())
	}

	// Verify we can connect to it
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		t.Fatalf("failed to dial listener: %v", err)
	}
	conn.Close()
}
