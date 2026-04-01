package main

import (
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
)

func TestResolveServiceAddr(t *testing.T) {
	services := map[string]serviceaddr.Service{}

	authn := serviceaddr.Service{Name: "authn"}
	authn.Ports.Grpc = 50051
	services["authn"] = authn

	if got := resolveServiceAddr("explicit:6000", services, "authn"); got != "explicit:6000" {
		t.Fatalf("expected explicit address, got %q", got)
	}
	t.Setenv("ENVIRONMENT", "development")
	if got := resolveServiceAddr("", services, "authn"); got != "localhost:50051" {
		t.Fatalf("expected service address from map, got %q", got)
	}
	t.Setenv("ENVIRONMENT", "production")
	if got := resolveServiceAddr("", services, "authn"); got != "authn:50051" {
		t.Fatalf("expected production service address from map, got %q", got)
	}
	if got := resolveServiceAddr("", services, "missing"); got != "" {
		t.Fatalf("expected empty address for missing service, got %q", got)
	}
}
