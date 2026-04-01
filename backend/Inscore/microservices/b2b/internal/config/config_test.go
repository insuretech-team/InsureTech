package config

import "testing"

func TestLoadReadsAuthServiceAddresses(t *testing.T) {
	t.Setenv("AUTHN_GRPC_ADDR", "localhost:50060")
	t.Setenv("AUTHZ_GRPC_ADDR", "localhost:50070")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AuthNServiceURL != "localhost:50060" {
		t.Fatalf("expected AuthNServiceURL to be loaded, got %q", cfg.AuthNServiceURL)
	}
	if cfg.AuthZServiceURL != "localhost:50070" {
		t.Fatalf("expected AuthZServiceURL to be loaded, got %q", cfg.AuthZServiceURL)
	}
}
