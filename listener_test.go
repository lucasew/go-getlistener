package getlistener

import (
	"testing"
)

func TestGetListener(t *testing.T) {
	ln, err := GetListener()
	if err != nil {
		t.Fatalf("GetListener failed: %v", err)
	}
	defer ln.Close()

	if ln.Addr() == nil {
		t.Errorf("Listener has no address")
	}
}
