package main

import (
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
)

func TestResolveServiceAddr(t *testing.T) {
	services := map[string]serviceaddr.Service{}

	authn := serviceaddr.Service{Name: "authn"}
	authn.Ports.Grpc = 50060
	services["authn"] = authn

	t.Run("explicit env wins", func(t *testing.T) {
		if got := resolveServiceAddr(" authn.example.internal:7000 ", services, "authn"); got != "authn.example.internal:7000" {
			t.Fatalf("expected explicit address, got %q", got)
		}
	})

	t.Run("host dev defaults to localhost", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "development")
		t.Setenv("B2B_SERVICE_DISCOVERY_HOST", "")

		if got := resolveServiceAddr("", services, "authn"); got != "localhost:50060" {
			t.Fatalf("expected localhost fallback, got %q", got)
		}
	})

	t.Run("production defaults to service name", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "production")
		t.Setenv("B2B_SERVICE_DISCOVERY_HOST", "")

		if got := resolveServiceAddr("", services, "authn"); got != "authn:50060" {
			t.Fatalf("expected production service-name fallback, got %q", got)
		}
	})

	t.Run("discovery host override works", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "development")
		t.Setenv("B2B_SERVICE_DISCOVERY_HOST", "host.docker.internal")

		if got := resolveServiceAddr("", services, "authn"); got != "host.docker.internal:50060" {
			t.Fatalf("expected discovery host override, got %q", got)
		}
	})

	t.Run("unknown environment falls back to runtime detection", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "staging")
		t.Setenv("B2B_SERVICE_DISCOVERY_HOST", "")

		got := resolveServiceAddr("", services, "authn")
		if got != "localhost:50060" && got != "authn:50060" {
			t.Fatalf("expected localhost or service-name fallback, got %q", got)
		}
	})

	t.Run("missing service returns empty", func(t *testing.T) {
		if got := resolveServiceAddr("", services, "missing"); got != "" {
			t.Fatalf("expected empty address for missing service, got %q", got)
		}
	})
}
