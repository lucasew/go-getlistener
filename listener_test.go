package getlistener

import (
	"net"
	"os"
	"testing"
)

func TestGetListener_RandomPort(t *testing.T) {
	// Ensure no env vars affect this test
	os.Unsetenv("PORT")
	os.Unsetenv("LISTEN_PID")
	os.Unsetenv("LISTEN_FDS")

	ln, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener() error = %v", err)
	}
	defer ln.Close()

	if ln.Addr().Network() != "tcp" {
		t.Errorf("GetListener() network = %v, want tcp", ln.Addr().Network())
	}

	_, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Errorf("Failed to split host/port: %v", err)
	}
	if portStr == "0" {
		t.Errorf("Port should not be 0 after binding")
	}
}
