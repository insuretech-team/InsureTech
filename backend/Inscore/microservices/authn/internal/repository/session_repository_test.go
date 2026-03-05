package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
)

func TestSessionRepository_RevokeAllByUserIDWithCount_ExcludesSessionID(t *testing.T) {
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	// Live DB enforces a strict chk_users_mobile_number + NOT NULL mobile_number.
	n := time.Now().UnixNano() % 1_000_000_000 // 9 digits
	mobile := fmt.Sprintf("+8801%09d", n)

	userID := uuid.New().String()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "test@example.com", "hashed_pw", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	sid1 := uuid.New().String()
	sid2 := uuid.New().String()
	sid3 := uuid.New().String()

	// Insert sessions using raw SQL to avoid schema drift issues.
	insertSessionMinimal(t, dbConn, sid1, userID, "SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	insertSessionMinimal(t, dbConn, sid2, userID, "SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	insertSessionMinimal(t, dbConn, sid3, userID, "SERVER_SIDE", true, time.Now().Add(24*time.Hour))

	revoked, err := sessionRepo.RevokeAllByUserIDWithCount(ctx, userID, sid2)
	require.NoError(t, err)
	require.Equal(t, int64(2), revoked)

	// sid2 should still be active
	s2, err := sessionRepo.GetByID(ctx, sid2)
	require.NoError(t, err)
	require.True(t, s2.IsActive)

	// others should be inactive, GetByID filters is_active=true so should fail
	_, err = sessionRepo.GetByID(ctx, sid1)
	require.Error(t, err)
	_, err = sessionRepo.GetByID(ctx, sid3)
	require.Error(t, err)
}
