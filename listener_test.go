package getlistener

import (
	"net"
	"testing"
)

func TestGetListener(t *testing.T) {
	ln, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener failed: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	if addr.Port == 0 {
		t.Errorf("Expected non-zero port, got 0")
	}
	t.Logf("Listening on port %d", addr.Port)
}
