package getlistener

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestListenTCP_WrapsListenError(t *testing.T) {
	// Bind first so the second listen fails with a real OS error (EADDRINUSE).
	first, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	port := first.Addr().(*net.TCPAddr).Port
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))

	_, err = listenTCP(&Config{Host: "127.0.0.1", Port: port})
	if err == nil {
		t.Fatal("expected error when port is already bound")
	}
	if !strings.Contains(err.Error(), addr) && !strings.Contains(err.Error(), fmt.Sprintf(":%d", port)) {
		t.Errorf("wrapped error should include listen address %q, got: %v", addr, err)
	}
	var opErr *net.OpError
	if !errors.As(err, &opErr) {
		t.Fatalf("errors.As(*net.OpError) failed; wrap should preserve cause: %v", err)
	}
}
