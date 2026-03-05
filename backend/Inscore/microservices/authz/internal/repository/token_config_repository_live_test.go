package repository

import (
	"context"
	"testing"

	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTokenConfigRepo_LiveDB_CreateListAndGetActive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthzDB(t)
	repo := NewTokenConfigRepo(dbConn)

	kid := newLiveID("kid")
	t.Cleanup(func() { cleanupTokenConfigByKID(t, dbConn, kid) })

	created, err := repo.Create(ctx, &authzentityv1.TokenConfig{
		Kid:           kid,
		Algorithm:     "RS256",
		PublicKeyPem:  "-----BEGIN PUBLIC KEY-----\nLIVE_TEST_KEY\n-----END PUBLIC KEY-----",
		PrivateKeyRef: "secret/authz/" + kid,
		IsActive:      false,
		CreatedAt:     timestamppb.Now(),
	})
	require.NoError(t, err)
	require.Equal(t, kid, created.Kid)

	cfgs, err := repo.List(ctx)
	require.NoError(t, err)
	found := false
	for _, cfg := range cfgs {
		if cfg.Kid == kid {
			found = true
			require.Equal(t, "RS256", cfg.Algorithm)
			require.Equal(t, "-----BEGIN PUBLIC KEY-----\nLIVE_TEST_KEY\n-----END PUBLIC KEY-----", cfg.PublicKeyPem)
			require.Equal(t, "secret/authz/"+kid, cfg.PrivateKeyRef)
			require.False(t, cfg.IsActive)
		}
	}
	require.True(t, found, "expected created token config in list")

	active, err := repo.GetActive(ctx)
	if err == nil {
		require.NotNil(t, active)
		require.NotEmpty(t, active.Kid)
		return
	}

	// Some environments may not have an active key seeded yet.
	fallbackKid := newLiveID("kid_active")
	t.Cleanup(func() { cleanupTokenConfigByKID(t, dbConn, fallbackKid) })
	_, err = repo.Create(ctx, &authzentityv1.TokenConfig{
		Kid:           fallbackKid,
		Algorithm:     "RS256",
		PublicKeyPem:  "-----BEGIN PUBLIC KEY-----\nLIVE_TEST_ACTIVE\n-----END PUBLIC KEY-----",
		PrivateKeyRef: "secret/authz/" + fallbackKid,
		IsActive:      true,
		CreatedAt:     timestamppb.Now(),
	})
	require.NoError(t, err)

	active, err = repo.GetActive(ctx)
	require.NoError(t, err)
	require.NotNil(t, active)
	require.NotEmpty(t, active.Kid)
}
