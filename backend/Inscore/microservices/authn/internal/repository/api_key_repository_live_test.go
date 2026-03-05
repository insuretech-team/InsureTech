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

func TestApiKeyRepository_LiveDB_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("live DB")
	}
	ctx := context.Background()
	db := testAuthnDB(t)
	repo := NewApiKeyRepository(db)

	id := uuid.New().String()
	ownerID := uuid.New().String()
	keyHash := "hash_" + uuid.New().String()

	// cleanup
	_ = repo.Delete(ctx, id)
	t.Cleanup(func() { _ = repo.Delete(ctx, id) })

	k := &apikeyv1.ApiKey{
		Id:                 id,
		KeyHash:            keyHash,
		Name:               "test",
		OwnerType:          apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL,
		OwnerId:            ownerID,
		Scopes:             []string{"policy:read"},
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: 10,
		ExpiresAt:          timestamppb.New(time.Now().Add(1 * time.Hour)),
		IpWhitelist:        []string{"127.0.0.1"},
		AuditInfo:          nil,
	}
	require.NoError(t, repo.Create(ctx, k))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, got.Id)

	got2, err := repo.GetByKeyHash(ctx, keyHash)
	require.NoError(t, err)
	require.Equal(t, keyHash, got2.KeyHash)

	// list
	st := apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE
	keys, err := repo.ListByOwner(ctx, apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL, ownerID, &st, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, keys)

	// touch
	require.NoError(t, repo.TouchLastUsed(ctx, id, time.Now()))

	// revoke
	require.NoError(t, repo.Revoke(ctx, id))
}
