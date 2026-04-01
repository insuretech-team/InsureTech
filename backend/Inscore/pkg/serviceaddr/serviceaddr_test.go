package serviceaddr

import "testing"

func TestResolveGRPCAddr(t *testing.T) {
	t.Run("explicit wins", func(t *testing.T) {
		got := ResolveGRPCAddr(" authn.internal:7000 ", "", "authn", 50060)
		if got != "authn.internal:7000" {
			t.Fatalf("ResolveGRPCAddr explicit = %q", got)
		}
	})

	t.Run("development uses localhost", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "development")
		t.Setenv("SERVICE_DISCOVERY_HOST", "")

		got := ResolveGRPCAddr("", "", "authn", 50060)
		if got != "localhost:50060" {
			t.Fatalf("ResolveGRPCAddr development = %q", got)
		}
	})

	t.Run("production uses service key", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "production")
		t.Setenv("SERVICE_DISCOVERY_HOST", "")

		got := ResolveGRPCAddr("", "", "authn", 50060)
		if got != "authn:50060" {
			t.Fatalf("ResolveGRPCAddr production = %q", got)
		}
	})

	t.Run("override host wins over environment", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "development")

		got := ResolveGRPCAddr("", "host.docker.internal", "authn", 50060)
		if got != "host.docker.internal:50060" {
			t.Fatalf("ResolveGRPCAddr override = %q", got)
		}
	})

	t.Run("global override host works", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "development")
		t.Setenv("SERVICE_DISCOVERY_HOST", "10.0.0.15")

		got := ResolveGRPCAddr("", "", "authn", 50060)
		if got != "10.0.0.15:50060" {
			t.Fatalf("ResolveGRPCAddr global override = %q", got)
		}
	})

	t.Run("unknown environment falls back to runtime detection", func(t *testing.T) {
		previous := isRunningInDocker
		isRunningInDocker = func() bool { return false }
		t.Cleanup(func() {
			isRunningInDocker = previous
		})

		t.Setenv("ENVIRONMENT", "staging")
		t.Setenv("SERVICE_DISCOVERY_HOST", "")

		got := ResolveGRPCAddr("", "", "authn", 50060)
		if got != "localhost:50060" {
			t.Fatalf("ResolveGRPCAddr fallback = %q", got)
		}
	})

	t.Run("missing service data returns empty", func(t *testing.T) {
		got := ResolveGRPCAddr("", "", "authn", 0)
		if got != "" {
			t.Fatalf("ResolveGRPCAddr missing = %q", got)
		}
	})
}
