package getlistener

import (
	"net"
	"testing"
)

func TestListenTCP_IPv6Loopback(t *testing.T) {
	ln, err := listenTCP(&Config{Host: "::1", Port: 0})
	if err != nil {
		t.Fatalf("listenTCP(::1): %v", err)
	}
	defer ln.Close()

	addr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("Addr type %T, want *net.TCPAddr", ln.Addr())
	}
	if !addr.IP.IsLoopback() {
		t.Errorf("bound IP %v is not loopback", addr.IP)
	}
	if addr.Port == 0 {
		t.Error("expected non-zero ephemeral port")
	}
}

func TestListenTCP_IPv4Loopback(t *testing.T) {
	ln, err := listenTCP(&Config{Host: "127.0.0.1", Port: 0})
	if err != nil {
		t.Fatalf("listenTCP(127.0.0.1): %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	if addr.Port == 0 {
		t.Error("expected non-zero ephemeral port")
	}
}
