package internalrpc

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestOutgoingAndValidateIncoming(t *testing.T) {
	t.Setenv(secretEnvKey, "test-secret")
	ctx := OutgoingContext(context.Background(), "gateway")
	outgoing, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("expected outgoing metadata")
	}
	incoming := metadata.NewIncomingContext(context.Background(), outgoing.Copy())
	serviceName, err := ValidateIncoming(incoming, map[string]struct{}{"gateway": {}})
	if err != nil {
		t.Fatalf("ValidateIncoming() error = %v", err)
	}
	if serviceName != "gateway" {
		t.Fatalf("ValidateIncoming() service = %q", serviceName)
	}
}

func TestValidateIncomingRejectsInvalidSignature(t *testing.T) {
	t.Setenv(secretEnvKey, "test-secret")
	signed := OutgoingContext(context.Background(), "gateway")
	outgoing, ok := metadata.FromOutgoingContext(signed)
	if !ok {
		t.Fatal("expected outgoing metadata")
	}
	md := outgoing.Copy()
	md.Set(HeaderSignature, "bad")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	if _, err := ValidateIncoming(ctx, map[string]struct{}{"gateway": {}}); err == nil {
		t.Fatal("expected invalid signature error")
	}
}
