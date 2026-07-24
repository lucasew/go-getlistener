package getlistener

import (
	"errors"
	"testing"
)

func TestSentinelErrorsDistinct(t *testing.T) {
	// Sentinel errors must be defined for all GOOS so multiplatform callers can
	// use errors.Is without build tags.
	sentinels := []error{ErrNotPassed, ErrWrongPid, ErrUnsupportedCase}
	for i, a := range sentinels {
		if a == nil {
			t.Fatalf("sentinel %d is nil", i)
		}
		if a.Error() == "" {
			t.Errorf("sentinel %d has empty message", i)
		}
		for j, b := range sentinels {
			if i == j {
				continue
			}
			if errors.Is(a, b) {
				t.Errorf("sentinel %d unexpectedly matches %d", i, j)
			}
		}
	}
}
