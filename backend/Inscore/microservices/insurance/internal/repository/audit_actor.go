package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
)

func resolveAuditActor(ctx context.Context) (string, error) {
	actorID := grpcmeta.ActorID(ctx, "")
	if actorID == "" {
		return "", fmt.Errorf("missing actor context for audit trail")
	}
	return actorID, nil
}

func newAuditInfoJSON(ctx context.Context) (string, error) {
	actorID, err := resolveAuditActor(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"created_by":"%s","created_at":"%s"}`, actorID, time.Now().UTC().Format(time.RFC3339)), nil
}
