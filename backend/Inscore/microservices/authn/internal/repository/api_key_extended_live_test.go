package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	apikeyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/apikey/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// ApiKeyRepository — additional coverage
// ---------------------------------------------------------------------------

func TestApiKeyRepository_ListByOwner_FiltersStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("live DB")
	}
	ctx := context.Background()
	db := testAuthnDB(t)
	repo := NewApiKeyRepository(db)

	ownerID := uuid.New().String()

	activeID := uuid.New().String()
	revokedID := uuid.New().String()
	activeHash := "hash_" + uuid.New().String()
	revokedHash := "hash_" + uuid.New().String()

	t.Cleanup(func() {
		_ = repo.Delete(ctx, activeID)
		_ = repo.Delete(ctx, revokedID)
	})

	// Create ACTIVE key
	requireNoError(t, repo.Create(ctx, &apikeyv1.ApiKey{
		Id:                 activeID,
		KeyHash:            activeHash,
		Name:               "active-key",
		OwnerType:          apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL,
		OwnerId:            ownerID,
		Scopes:             []string{"policy:read"},
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: 10,
		ExpiresAt:          timestamppb.New(time.Now().Add(1 * time.Hour)),
		IpWhitelist:        []string{"127.0.0.1"},
	}))

	// Create key then revoke it
	requireNoError(t, repo.Create(ctx, &apikeyv1.ApiKey{
		Id:                 revokedID,
		KeyHash:            revokedHash,
		Name:               "revoked-key",
		OwnerType:          apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL,
		OwnerId:            ownerID,
		Scopes:             []string{"policy:read"},
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: 10,
		ExpiresAt:          timestamppb.New(time.Now().Add(1 * time.Hour)),
		IpWhitelist:        []string{"127.0.0.1"},
	}))
	requireNoError(t, repo.Revoke(ctx, revokedID))

	// ListByOwner with status=ACTIVE filter
	st := apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE
	keys, err := repo.ListByOwner(ctx, apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL, ownerID, &st, 100, 0)
	require.NoError(t, err)

	// Collect returned IDs
	ids := make(map[string]bool, len(keys))
	for _, k := range keys {
		ids[k.Id] = true
	}

	require.True(t, ids[activeID], "ACTIVE key should be in results")
	require.False(t, ids[revokedID], "REVOKED key should not be in results")
}

func TestApiKeyRepository_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("live DB")
	}
	ctx := context.Background()
	db := testAuthnDB(t)
	repo := NewApiKeyRepository(db)

	randomID := uuid.New().String()
	_, err := repo.GetByID(ctx, randomID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestApiKeyRepository_UpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("live DB")
	}
	ctx := context.Background()
	db := testAuthnDB(t)
	repo := NewApiKeyRepository(db)

	id := uuid.New().String()
	keyHash := "hash_" + uuid.New().String()
	ownerID := uuid.New().String()

	t.Cleanup(func() { _ = repo.Delete(ctx, id) })

	requireNoError(t, repo.Create(ctx, &apikeyv1.ApiKey{
		Id:                 id,
		KeyHash:            keyHash,
		Name:               "status-test",
		OwnerType:          apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL,
		OwnerId:            ownerID,
		Scopes:             []string{"policy:read"},
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: 5,
		ExpiresAt:          timestamppb.New(time.Now().Add(1 * time.Hour)),
		IpWhitelist:        []string{},
	}))

	requireNoError(t, repo.UpdateStatus(ctx, id, apikeyv1.ApiKeyStatus_API_KEY_STATUS_EXPIRED))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, apikeyv1.ApiKeyStatus_API_KEY_STATUS_EXPIRED, got.Status)
}

// ---------------------------------------------------------------------------
// ApiKeyUsageRepository — additional coverage
// ---------------------------------------------------------------------------

func TestApiKeyUsageRepository_ListByApiKey_TimeFilter(t *testing.T) {
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

	t.Cleanup(func() {
		_, _ = uRepo.DeleteByApiKey(ctx, keyID)
		_ = keyRepo.Delete(ctx, keyID)
	})

	requireNoError(t, keyRepo.Create(ctx, &apikeyv1.ApiKey{
		Id:                 keyID,
		KeyHash:            keyHash,
		Name:               "time-filter-test",
		OwnerType:          apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL,
		OwnerId:            ownerID,
		Scopes:             []string{"policy:read"},
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: 10,
	}))

	// Old usage — 2 hours ago
	oldID := uuid.New().String()
	requireNoError(t, uRepo.Create(ctx, &apikeyv1.ApiKeyUsage{
		Id:         oldID,
		ApiKeyId:   keyID,
		Endpoint:   "/v1/policy",
		HttpMethod: "GET",
		StatusCode: 200,
		RequestIp:  "127.0.0.1",
		Timestamp:  timestamppb.New(time.Now().Add(-2 * time.Hour)),
	}))

	// Recent usage — now
	recentID := uuid.New().String()
	requireNoError(t, uRepo.Create(ctx, &apikeyv1.ApiKeyUsage{
		Id:         recentID,
		ApiKeyId:   keyID,
		Endpoint:   "/v1/policy",
		HttpMethod: "GET",
		StatusCode: 200,
		RequestIp:  "127.0.0.1",
		Timestamp:  timestamppb.New(time.Now()),
	}))

	// Filter: from = 30 minutes ago → only recent record should appear
	from := time.Now().Add(-30 * time.Minute)
	list, err := uRepo.ListByApiKey(ctx, keyID, &from, nil, 100, 0)
	require.NoError(t, err)

	ids := make(map[string]bool, len(list))
	for _, u := range list {
		ids[u.Id] = true
	}

	require.True(t, ids[recentID], "recent usage should be returned")
	require.False(t, ids[oldID], "old usage should be filtered out")
}

func TestApiKeyUsageRepository_DeleteByApiKey_Count(t *testing.T) {
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

	t.Cleanup(func() {
		_, _ = uRepo.DeleteByApiKey(ctx, keyID)
		_ = keyRepo.Delete(ctx, keyID)
	})

	requireNoError(t, keyRepo.Create(ctx, &apikeyv1.ApiKey{
		Id:                 keyID,
		KeyHash:            keyHash,
		Name:               "delete-count-test",
		OwnerType:          apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_INTERNAL,
		OwnerId:            ownerID,
		Scopes:             []string{"policy:read"},
		Status:             apikeyv1.ApiKeyStatus_API_KEY_STATUS_ACTIVE,
		RateLimitPerMinute: 10,
	}))

	// Pre-clean any leftover records
	_, _ = uRepo.DeleteByApiKey(ctx, keyID)

	// Insert 3 usage records
	for i := range 3 {
		requireNoError(t, uRepo.Create(ctx, &apikeyv1.ApiKeyUsage{
			Id:         uuid.New().String(),
			ApiKeyId:   keyID,
			Endpoint:   "/v1/policy",
			HttpMethod: "GET",
			StatusCode: int32(200 + i),
			RequestIp:  "127.0.0.1",
			Timestamp:  timestamppb.New(time.Now()),
		}))
	}

	rowsAffected, err := uRepo.DeleteByApiKey(ctx, keyID)
	require.NoError(t, err)
	require.EqualValues(t, 3, rowsAffected)
}
