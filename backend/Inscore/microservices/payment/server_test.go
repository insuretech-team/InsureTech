package payment

import (
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
)

func TestResolveServiceAddr(t *testing.T) {
	services := map[string]serviceaddr.Service{}

	authz := serviceaddr.Service{Name: "authz"}
	authz.Ports.Grpc = 50070
	services["authz"] = authz

	if got := resolveServiceAddr(" authz.internal:7000 ", services, "authz"); got != "authz.internal:7000" {
		t.Fatalf("resolveServiceAddr explicit = %q", got)
	}

	t.Setenv("ENVIRONMENT", "development")
	if got := resolveServiceAddr("", services, "authz"); got != "localhost:50070" {
		t.Fatalf("resolveServiceAddr development = %q", got)
	}

	t.Setenv("ENVIRONMENT", "production")
	if got := resolveServiceAddr("", services, "authz"); got != "authz:50070" {
		t.Fatalf("resolveServiceAddr production = %q", got)
	}

	if got := resolveServiceAddr("", services, "missing"); got != "" {
		t.Fatalf("resolveServiceAddr missing = %q", got)
	}
}
