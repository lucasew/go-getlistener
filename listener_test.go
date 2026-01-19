package getlistener

import (
	"net"
	"os"
	"testing"
)

func TestGetListener(t *testing.T) {
	// Ensure PORT is not set so we test the random port logic
	os.Unsetenv("PORT")
	os.Unsetenv("HOST")
	os.Unsetenv("LISTEN_PID")
	os.Unsetenv("LISTEN_FDS")

	ln, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener failed: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr()
	t.Logf("Listening on %s", addr)

	// Verify we can connect to it
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		t.Fatalf("Failed to connect to listener: %v", err)
	}
	conn.Close()
}
