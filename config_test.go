package getlistener

import (
	"strings"
	"testing"
)

func TestLoadConfig_PortOutOfRange(t *testing.T) {
	cases := []string{"-1", "65536", "99999"}
	for _, port := range cases {
		t.Run(port, func(t *testing.T) {
			t.Setenv("PORT", port)
			t.Setenv("HOST", "127.0.0.1")

			cfg, err := loadConfig()
			if cfg != nil {
				t.Fatalf("expected nil config for PORT=%s", port)
			}
			if err == nil {
				t.Fatalf("expected error for PORT=%s", port)
			}
			if !strings.Contains(err.Error(), port) {
				t.Errorf("error should include invalid PORT value %q, got: %v", port, err)
			}
		})
	}
}

func TestLoadConfig_PortBoundaryOK(t *testing.T) {
	for _, port := range []string{"0", "1", "65535"} {
		t.Run(port, func(t *testing.T) {
			t.Setenv("PORT", port)
			t.Setenv("HOST", "127.0.0.1")

			cfg, err := loadConfig()
			if err != nil {
				t.Fatalf("loadConfig(PORT=%s): %v", port, err)
			}
			if cfg.Port < 0 || cfg.Port > 65535 {
				t.Errorf("Port = %d, want in 0..65535", cfg.Port)
			}
		})
	}
}
