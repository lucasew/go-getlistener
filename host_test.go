package getlistener

import "testing"

func TestIsLocalHost(t *testing.T) {
	local := []string{"localhost", "127.0.0.1", "::1", "127.0.0.2"}
	for _, h := range local {
		if !isLocalHost(h) {
			t.Errorf("isLocalHost(%q) = false, want true", h)
		}
	}
	// Unspecified and non-loopback must warn (non-local).
	nonLocal := []string{"0.0.0.0", "::", "192.168.1.1", "example.com", ""}
	for _, h := range nonLocal {
		if isLocalHost(h) {
			t.Errorf("isLocalHost(%q) = true, want false", h)
		}
	}
}
