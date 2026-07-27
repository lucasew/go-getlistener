package getlistener

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestLoadConfig_InvalidPortWrapsParseCause(t *testing.T) {
	t.Setenv("PORT", "not-a-port")
	t.Setenv("HOST", "127.0.0.1")

	cfg, err := loadConfig()
	if cfg != nil {
		t.Fatal("expected nil config for non-numeric PORT")
	}
	if err == nil {
		t.Fatal("expected error for non-numeric PORT")
	}
	// CONSISTENTLY_IGNORED: invalid PORT value must remain in the message.
	if !strings.Contains(err.Error(), "not-a-port") {
		t.Errorf("error should include invalid PORT value, got: %v", err)
	}
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fatalf("expected *strconv.NumError in chain, got: %v", err)
	}
	if numErr.Num != "not-a-port" {
		t.Errorf("NumError.Num = %q, want not-a-port", numErr.Num)
	}
}
