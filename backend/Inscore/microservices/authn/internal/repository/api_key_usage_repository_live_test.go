package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	apikeyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/apikey/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestApiKeyUsageRepository_LiveDB_CreateListDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("live DB")
	}
	ctx := context.Background()
	db := testAuthnDB(t)
	keyRepo := NewApiKeyRepository(db)
	uRepo := NewApiKeyUsageRepository(db)

	keyID := uuid.New().String()
	ownerID := uuid.New().String()
	keyHash := "hash_" + uuid.New().String()
	_ = keyRepo.Delete(ctx, keyID)
	t.Cleanup(func() { _ = keyRepo.Delete(ctx, keyID) })

	require.NoError(t, keyRepo.Create(ctx, &apikeyv1.ApiKey{
		Id:                 keyID,
		KeyHash:            keyHash,
		Name:               "test",
		OwnerType:          apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL,
		OwnerId:            ownerID,
		Scopes:             []string{"policy:read"},
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: 10,
		AuditInfo:          nil,
	}))

	usageID := uuid.New().String()
	_, _ = uRepo.DeleteByApiKey(ctx, keyID)
	t.Cleanup(func() { _, _ = uRepo.DeleteByApiKey(ctx, keyID) })

	u := &apikeyv1.ApiKeyUsage{
		Id:             usageID,
		ApiKeyId:       keyID,
		Endpoint:       "/v1/policy",
		HttpMethod:     "GET",
		StatusCode:     200,
		ResponseTimeMs: 12,
		RequestIp:      "127.0.0.1",
		TraceId:        "trace",
		Timestamp:      timestamppb.New(time.Now()),
	}
	require.NoError(t, uRepo.Create(ctx, u))

	got, err := uRepo.GetByID(ctx, usageID)
	require.NoError(t, err)
	require.Equal(t, usageID, got.Id)

	from := time.Now().Add(-1 * time.Hour)
	to := time.Now().Add(1 * time.Hour)
	list, err := uRepo.ListByApiKey(ctx, keyID, &from, &to, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, list)
}
