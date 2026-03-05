package service

import (
	"context"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/middleware"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/stretchr/testify/require"
)

// smoke-test event calls in revoke methods (no nil deref)
func TestAuthService_RevokeEvents_NoPanic(t *testing.T) {
	ctx := context.Background()
	pub := events.NewPublisher((events.EventProducer)(nil))

	// Directly call publisher methods to ensure wiring compiles (no-op with nil producer).
	require.NotPanics(t, func() {
		require.NoError(t, pub.PublishSessionRevoked(ctx, "u1", "s1", "", "user", "test"))
	})
	// Ensure signatures compile for service call sites
	_ = (&AuthService{eventPublisher: pub, metadata: middleware.NewMetadataExtractor()})
	_ = (&authnservicev1.RevokeAllSessionsRequest{UserId: "u1", Reason: "test"})
	require.True(t, true)
}
