package main

import (
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
)

func TestHelperFunctions(t *testing.T) {
	t.Setenv("ORDERS_TEST_ENV", "value")
	if got := envOrDefault("ORDERS_TEST_ENV", "fallback"); got != "value" {
		t.Fatalf("envOrDefault() = %q", got)
	}
	if got := envOrDefault("ORDERS_MISSING_ENV", "fallback"); got != "fallback" {
		t.Fatalf("envOrDefault missing = %q", got)
	}

	services := map[string]serviceaddr.Service{}

	payment := serviceaddr.Service{Name: "payment"}
	payment.Ports.Grpc = 50190
	services["payment"] = payment

	if got := resolveServiceAddr(" payment:6000 ", services, "payment"); got != "payment:6000" {
		t.Fatalf("resolveServiceAddr explicit = %q", got)
	}
	t.Setenv("ENVIRONMENT", "development")
	if got := resolveServiceAddr("", services, "payment"); got != "localhost:50190" {
		t.Fatalf("resolveServiceAddr inferred = %q", got)
	}
	t.Setenv("ENVIRONMENT", "production")
	if got := resolveServiceAddr("", services, "payment"); got != "payment:50190" {
		t.Fatalf("resolveServiceAddr production = %q", got)
	}
	if got := resolveServiceAddr("", services, "missing"); got != "" {
		t.Fatalf("resolveServiceAddr missing = %q", got)
	}
}
