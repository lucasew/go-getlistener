package getlistener

import (
	"net"
	"strconv"
	"strings"
	"testing"
)

func clearListenerEnv(t *testing.T) {
	t.Helper()
	t.Setenv("PORT", "")
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("LISTEN_PID", "")
	t.Setenv("LISTEN_FDS", "")
}

func TestGetListener(t *testing.T) {
	clearListenerEnv(t)

	ln, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener failed: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	if addr.Port == 0 {
		t.Errorf("Expected bound port to be non-zero, got 0")
	}
	t.Logf("Listening on %s", addr.String())
}

func TestGetListener_InvalidPORT(t *testing.T) {
	clearListenerEnv(t)
	t.Setenv("PORT", "not-a-port")

	ln, err := GetListener()
	if ln != nil {
		ln.Close()
		t.Fatal("expected nil listener for invalid PORT")
	}
	if err == nil {
		t.Fatal("expected error for invalid PORT")
	}
	if !strings.Contains(err.Error(), "not-a-port") {
		t.Errorf("error should include the invalid PORT value, got: %v", err)
	}
}

func TestGetListener_ExplicitPORT(t *testing.T) {
	clearListenerEnv(t)

	// Reserve an ephemeral port, release it, then ask GetListener to bind it.
	tmp, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	port := tmp.Addr().(*net.TCPAddr).Port
	tmp.Close()

	t.Setenv("PORT", strconv.Itoa(port))

	ln, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener with PORT=%d failed: %v", port, err)
	}
	defer ln.Close()

	got := ln.Addr().(*net.TCPAddr).Port
	if got != port {
		t.Errorf("bound port = %d, want %d", got, port)
	}
}

func TestGetListener_DefaultHost(t *testing.T) {
	clearListenerEnv(t)
	t.Setenv("HOST", "")
	t.Setenv("PORT", "")

	ln, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener failed: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	if !addr.IP.IsLoopback() {
		t.Errorf("default host should bind loopback, got %v", addr.IP)
	}
}

func TestLoadConfig_EmptyPORTUsesZero(t *testing.T) {
	clearListenerEnv(t)
	t.Setenv("PORT", "")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.Port != 0 {
		t.Errorf("Port = %d, want 0 (random)", cfg.Port)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("Host = %q, want 127.0.0.1", cfg.Host)
	}
}

func TestLoadConfig_CustomHost(t *testing.T) {
	clearListenerEnv(t)
	t.Setenv("HOST", "localhost")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.Host != "localhost" {
		t.Errorf("Host = %q, want localhost", cfg.Host)
	}
}
