package middleware

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestExtractRequestContextAndHelpers(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-user-id", "user-1",
		"x-tenant-id", "tenant-1",
		"x-portal", "PORTAL_B2B",
		"x-session-id", "sess-1",
		"x-session-type", "JWT",
		"x-token-id", "token-1",
		"x-device-id", "device-1",
		"x-user-type", "B2C_CUSTOMER",
		"x-business-id", "org-1",
		"x-request-id", "trace-1",
	))

	rctx := ExtractRequestContext(ctx)
	if rctx.UserID != "user-1" || rctx.TenantID != "tenant-1" || rctx.Portal != "b2b" {
		t.Fatalf("unexpected request context: %+v", rctx)
	}
	if rctx.ActorUserID() != "user-1" {
		t.Fatalf("expected actor user id")
	}
	if !rctx.IsB2B() || rctx.IsSystemPortal() {
		t.Fatalf("unexpected portal helpers")
	}
	if !rctx.IsMobileOrAPI() || rctx.IsWebPortal() {
		t.Fatalf("unexpected session helpers")
	}

	if got := normPortal(" PORTAL_SYSTEM "); got != "system" {
		t.Fatalf("normPortal() = %q", got)
	}
	if got := firstMD(metadata.Pairs("k", " ", "k", "v"), "k"); got != "v" {
		t.Fatalf("firstMD() = %q", got)
	}
}
