//go:build windows

package getlistener

import (
	"strings"
	"testing"
)

func TestGetListener_RejectsListenPIDOnWindows(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("LISTEN_PID", "1")
	t.Setenv("LISTEN_FDS", "1")

	ln, err := GetListener()
	if ln != nil {
		ln.Close()
		t.Fatal("expected nil listener when LISTEN_PID is set on windows")
	}
	if err == nil {
		t.Fatal("expected error when LISTEN_PID is set on windows")
	}
	if !strings.Contains(err.Error(), "not supported on windows") {
		t.Errorf("error should mention windows unsupported activation, got: %v", err)
	}
	if !strings.Contains(err.Error(), "LISTEN_PID") {
		t.Errorf("error should mention LISTEN_PID, got: %v", err)
	}
}
