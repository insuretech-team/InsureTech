package grpcmeta

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestExtractRequestContext(t *testing.T) {
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
	if !rctx.IsB2B() || !rctx.IsMobileOrAPI() || rctx.IsSystemPortal() {
		t.Fatalf("unexpected helper behaviour: %+v", rctx)
	}
}

func TestExtractAllAndCookieHelpers(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "10.0.0.1, 10.0.0.2",
		"user-agent", "ua-test",
		"x-device-id", "dev-1",
		"cookie", "session_token=s123; theme=light",
		"x-csrf-token", "csrf-1",
		"authorization", "Bearer tok-1",
	))

	out := ExtractAll(ctx)
	if out.IPAddress != "10.0.0.1" || out.UserAgent != "ua-test" || out.Authorization != "tok-1" {
		t.Fatalf("unexpected request metadata: %+v", out)
	}
	if got := ParseCookie("k1=v1; k2=v2", "missing"); got != "" {
		t.Fatalf("expected missing cookie, got %q", got)
	}
}

func TestActorTenantAndCorrelationHelpers(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-actor-id", " actor-1 ",
		"tenant-id", "tenant-1",
		"x-request-id", "req-1",
	))

	if got := ActorID(ctx, "fallback"); got != "actor-1" {
		t.Fatalf("ActorID() = %q", got)
	}
	if got := TenantID(ctx, "fallback-tenant"); got != "tenant-1" {
		t.Fatalf("TenantID() = %q", got)
	}
	if got := CorrelationID(ctx); got != "req-1" {
		t.Fatalf("CorrelationID() = %q", got)
	}
}

func TestExtractIPAddressFallback(t *testing.T) {
	ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.IPAddr{IP: net.ParseIP("127.0.0.1")}})
	if got := ExtractIPAddress(ctx); got != "127.0.0.1" {
		t.Fatalf("ExtractIPAddress() = %q", got)
	}
}

func TestExtractIPAddressPeerTCPAddr(t *testing.T) {
	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 63222},
	})
	if got := ExtractIPAddress(ctx); got != "127.0.0.1" {
		t.Fatalf("ExtractIPAddress() = %q", got)
	}
}

func TestExtractIPAddressPeerTCPAddrIPv6(t *testing.T) {
	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{IP: net.ParseIP("::1"), Port: 63222},
	})
	if got := ExtractIPAddress(ctx); got != "::1" {
		t.Fatalf("ExtractIPAddress() = %q", got)
	}
}
