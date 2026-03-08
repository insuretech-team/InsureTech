package service

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

// ── resolveCallerID ───────────────────────────────────────────────────────────

func TestResolveCallerID_FromMetadata(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs("x-user-id", "user-abc"))
	got := resolveCallerID(ctx, "fallback")
	if got != "user-abc" {
		t.Fatalf("expected user-abc, got %s", got)
	}
}

func TestResolveCallerID_FallsBackToDefault(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs("x-other-header", "something"))
	got := resolveCallerID(ctx, "default-user")
	if got != "default-user" {
		t.Fatalf("expected default-user, got %s", got)
	}
}

func TestResolveCallerID_FallsBackToSystem(t *testing.T) {
	got := resolveCallerID(context.Background(), "")
	if got != "system" {
		t.Fatalf("expected system, got %s", got)
	}
}

func TestResolveCallerID_TrimsWhitespace(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs("x-user-id", "  user-xyz  "))
	got := resolveCallerID(ctx, "fallback")
	if got != "user-xyz" {
		t.Fatalf("expected user-xyz (trimmed), got %q", got)
	}
}

func TestResolveCallerID_EmptyMetadataValue_UsesFallback(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs("x-user-id", "   "))
	got := resolveCallerID(ctx, "my-fallback")
	if got != "my-fallback" {
		t.Fatalf("expected my-fallback, got %s", got)
	}
}

func TestResolveCallerID_NoIncomingContext_UsesFallback(t *testing.T) {
	// Plain background context has no incoming metadata.
	got := resolveCallerID(context.Background(), "svc-caller")
	if got != "svc-caller" {
		t.Fatalf("expected svc-caller, got %s", got)
	}
}

// ── resolveTenantID ───────────────────────────────────────────────────────────

func TestResolveTenantID_ExplicitTakesPrecedence(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs("x-tenant-id", "from-meta"))
	got := resolveTenantID(ctx, "explicit-tenant")
	if got != "explicit-tenant" {
		t.Fatalf("expected explicit-tenant, got %s", got)
	}
}

func TestResolveTenantID_FallsBackToMetadata(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs("x-tenant-id", "meta-tenant"))
	got := resolveTenantID(ctx, "")
	if got != "meta-tenant" {
		t.Fatalf("expected meta-tenant, got %s", got)
	}
}

func TestResolveTenantID_FallsBackToDefault(t *testing.T) {
	got := resolveTenantID(context.Background(), "")
	// Should return the hard-coded default UUID when no env and no metadata.
	if got == "" {
		t.Fatal("expected non-empty tenant ID")
	}
}
