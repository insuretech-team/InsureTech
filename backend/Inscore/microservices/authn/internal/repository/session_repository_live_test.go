package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ─── Create + GetByID ────────────────────────────────────────────────────────

func TestSessionRepository_LiveDB_CreateAndGetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	// Set up a user the sessions will belong to.
	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_create@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	sessionID := uuid.New().String()
	lookup := "lookup-" + uuid.New().String()

	session := &authnentityv1.Session{
		SessionId:          sessionID,
		UserId:             userID,
		SessionType:        authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE,
		SessionTokenHash:   "bcrypt-hash-placeholder",
		SessionTokenLookup: lookup,
		DeviceId:           "device-create-test",
		DeviceName:         "Test Browser",
		DeviceType:         authnentityv1.DeviceType_DEVICE_TYPE_WEB,
		IpAddress:          "127.0.0.1",
		UserAgent:          "Go-Test/1.0",
		IsActive:           true,
		ExpiresAt:          timestamppb.New(time.Now().Add(24 * time.Hour)),
	}

	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", sessionID).Delete(map[string]any{}).Error
	})

	got, err := sessionRepo.GetByID(ctx, sessionID)
	require.NoError(t, err)
	require.Equal(t, sessionID, got.SessionId)
	require.Equal(t, userID, got.UserId)
	require.True(t, got.IsActive)
}

// ─── GetByTokenLookup ────────────────────────────────────────────────────────

func TestSessionRepository_LiveDB_GetByTokenLookup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_lookup@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	sessionID := uuid.New().String()
	// Use a sha256-like deterministic string as the lookup value.
	lookup := "a3f1b2c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2"

	// Insert via helper to avoid schema drift on optional columns.
	insertSessionMinimal(t, dbConn, sessionID, userID, "SESSION_TYPE_SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", sessionID).Delete(map[string]any{}).Error
	})

	// Patch the lookup column directly so GetByTokenLookup can find it.
	err := dbConn.Table("authn_schema.sessions").
		Where("session_id = ?", sessionID).
		Update("session_token_lookup", lookup).Error
	require.NoError(t, err)

	got, err := sessionRepo.GetByTokenLookup(ctx, lookup)
	require.NoError(t, err)
	require.Equal(t, sessionID, got.SessionId)
}

// ─── GetByRefreshToken ───────────────────────────────────────────────────────

func TestSessionRepository_LiveDB_GetByRefreshToken(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_refresh@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	sessionID := uuid.New().String()
	refreshJTI := uuid.New().String()

	session := &authnentityv1.Session{
		SessionId:             sessionID,
		UserId:                userID,
		SessionType:           authnentityv1.SessionType_SESSION_TYPE_JWT,
		AccessTokenJti:        uuid.New().String(),
		RefreshTokenJti:       refreshJTI,
		AccessTokenExpiresAt:  timestamppb.New(time.Now().Add(15 * time.Minute)),
		RefreshTokenExpiresAt: timestamppb.New(time.Now().Add(7 * 24 * time.Hour)),
		DeviceId:              "device-jwt-test",
		DeviceName:            "Android Phone",
		DeviceType:            authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_ANDROID,
		IpAddress:             "10.0.0.1",
		UserAgent:             "Insuretech-Android/1.0",
		IsActive:              true,
		ExpiresAt:             timestamppb.New(time.Now().Add(7 * 24 * time.Hour)),
	}

	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", sessionID).Delete(map[string]any{}).Error
	})

	got, err := sessionRepo.GetByRefreshToken(ctx, refreshJTI)
	require.NoError(t, err)
	require.Equal(t, sessionID, got.SessionId)
	require.Equal(t, refreshJTI, got.RefreshTokenJti)
}

// ─── UpdateLastActivity ──────────────────────────────────────────────────────

func TestSessionRepository_LiveDB_UpdateLastActivity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_activity@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	sessionID := uuid.New().String()
	insertSessionMinimal(t, dbConn, sessionID, userID, "SESSION_TYPE_SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", sessionID).Delete(map[string]any{}).Error
	})

	// Record time before update.
	before := time.Now()
	time.Sleep(10 * time.Millisecond) // ensure measurable delta

	err := sessionRepo.UpdateLastActivity(ctx, sessionID)
	require.NoError(t, err)

	// Fetch raw last_activity_at to verify it was updated.
	var lastActivity time.Time
	err = dbConn.Raw(
		"SELECT last_activity_at FROM authn_schema.sessions WHERE session_id = ?",
		sessionID,
	).Scan(&lastActivity).Error
	require.NoError(t, err)
	require.True(t, lastActivity.After(before) || lastActivity.Equal(before),
		"expected last_activity_at to be updated, got %v (before=%v)", lastActivity, before)
}

// ─── UpdateTokens ────────────────────────────────────────────────────────────

func TestSessionRepository_LiveDB_UpdateTokens(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_tokens@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	sessionID := uuid.New().String()
	oldAccessJTI := uuid.New().String()
	oldRefreshJTI := uuid.New().String()

	session := &authnentityv1.Session{
		SessionId:             sessionID,
		UserId:                userID,
		SessionType:           authnentityv1.SessionType_SESSION_TYPE_JWT,
		AccessTokenJti:        oldAccessJTI,
		RefreshTokenJti:       oldRefreshJTI,
		AccessTokenExpiresAt:  timestamppb.New(time.Now().Add(15 * time.Minute)),
		RefreshTokenExpiresAt: timestamppb.New(time.Now().Add(7 * 24 * time.Hour)),
		DeviceId:              "device-token-update",
		DeviceName:            "iOS Device",
		DeviceType:            authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_IOS,
		IpAddress:             "192.168.1.1",
		UserAgent:             "Insuretech-iOS/1.0",
		IsActive:              true,
		ExpiresAt:             timestamppb.New(time.Now().Add(7 * 24 * time.Hour)),
	}

	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", sessionID).Delete(map[string]any{}).Error
	})

	newAccessJTI := uuid.New().String()
	newRefreshJTI := uuid.New().String()

	err = sessionRepo.UpdateTokens(ctx, sessionID, newAccessJTI, newRefreshJTI)
	require.NoError(t, err)

	// Verify the updated JTIs are persisted.
	var gotAccessJTI, gotRefreshJTI string
	err = dbConn.Raw(
		"SELECT access_token_jti, refresh_token_jti FROM authn_schema.sessions WHERE session_id = ?",
		sessionID,
	).Row().Scan(&gotAccessJTI, &gotRefreshJTI)
	require.NoError(t, err)
	require.Equal(t, newAccessJTI, gotAccessJTI)
	require.Equal(t, newRefreshJTI, gotRefreshJTI)
}

// ─── Revoke ──────────────────────────────────────────────────────────────────

func TestSessionRepository_LiveDB_Revoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_revoke@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	sessionID := uuid.New().String()
	insertSessionMinimal(t, dbConn, sessionID, userID, "SESSION_TYPE_SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", sessionID).Delete(map[string]any{}).Error
	})

	// Confirm session is active before revoke.
	got, err := sessionRepo.GetByID(ctx, sessionID)
	require.NoError(t, err)
	require.True(t, got.IsActive)

	// Revoke the session.
	err = sessionRepo.Revoke(ctx, sessionID)
	require.NoError(t, err)

	// GetByID filters is_active=true — should now return an error.
	_, err = sessionRepo.GetByID(ctx, sessionID)
	require.Error(t, err, "GetByID should fail for a revoked (inactive) session")
}

// ─── ListByUserID ─────────────────────────────────────────────────────────────

func TestSessionRepository_LiveDB_ListByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_list@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	// Insert two active SERVER_SIDE sessions and one inactive one.
	sid1 := uuid.New().String()
	sid2 := uuid.New().String()
	sid3 := uuid.New().String() // will be revoked

	insertSessionMinimal(t, dbConn, sid1, userID, "SESSION_TYPE_SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	insertSessionMinimal(t, dbConn, sid2, userID, "SESSION_TYPE_SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	insertSessionMinimal(t, dbConn, sid3, userID, "SESSION_TYPE_SERVER_SIDE", false, time.Now().Add(24*time.Hour))

	t.Cleanup(func() {
		for _, sid := range []string{sid1, sid2, sid3} {
			_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", sid).Delete(map[string]any{}).Error
		}
	})

	// activeOnly=true — should return only the 2 active ones.
	active, err := sessionRepo.ListByUserID(ctx, userID, true, nil)
	require.NoError(t, err)
	require.Len(t, active, 2, "expected exactly 2 active sessions")

	// activeOnly=false — should return all 3.
	all, err := sessionRepo.ListByUserID(ctx, userID, false, nil)
	require.NoError(t, err)
	require.Len(t, all, 3, "expected all 3 sessions regardless of active status")

	// Filter by SessionType SERVER_SIDE active only.
	stype := authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE
	filtered, err := sessionRepo.ListByUserID(ctx, userID, true, &stype)
	require.NoError(t, err)
	require.Len(t, filtered, 2, "expected 2 active SERVER_SIDE sessions")
	for _, s := range filtered {
		require.Equal(t, authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE, s.SessionType)
	}
}

// ─── ListByUserID with JWT sessionType filter ─────────────────────────────────

func TestSessionRepository_LiveDB_ListByUserID_JWTFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_list_jwt@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	// Create one JWT session via Create (populates all NOT NULL fields).
	jwtSessionID := uuid.New().String()
	jwtSession := &authnentityv1.Session{
		SessionId:             jwtSessionID,
		UserId:                userID,
		SessionType:           authnentityv1.SessionType_SESSION_TYPE_JWT,
		AccessTokenJti:        uuid.New().String(),
		RefreshTokenJti:       uuid.New().String(),
		AccessTokenExpiresAt:  timestamppb.New(time.Now().Add(15 * time.Minute)),
		RefreshTokenExpiresAt: timestamppb.New(time.Now().Add(7 * 24 * time.Hour)),
		DeviceId:              "device-jwt-list",
		DeviceName:            "Test Android",
		DeviceType:            authnentityv1.DeviceType_DEVICE_TYPE_MOBILE_ANDROID,
		IpAddress:             "10.1.2.3",
		UserAgent:             "Android/Go-Test",
		IsActive:              true,
		ExpiresAt:             timestamppb.New(time.Now().Add(7 * 24 * time.Hour)),
	}
	err := sessionRepo.Create(ctx, jwtSession)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", jwtSessionID).Delete(map[string]any{}).Error
	})

	// Also insert a SERVER_SIDE session.
	webSessionID := uuid.New().String()
	insertSessionMinimal(t, dbConn, webSessionID, userID, "SESSION_TYPE_SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	t.Cleanup(func() {
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", webSessionID).Delete(map[string]any{}).Error
	})

	jwtType := authnentityv1.SessionType_SESSION_TYPE_JWT
	jwtSessions, err := sessionRepo.ListByUserID(ctx, userID, true, &jwtType)
	require.NoError(t, err)
	require.Len(t, jwtSessions, 1, "expected exactly 1 JWT session")
	require.Equal(t, jwtSessionID, jwtSessions[0].SessionId)
}

// ─── CleanupExpiredSessions ───────────────────────────────────────────────────

func TestSessionRepository_LiveDB_CleanupExpiredSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testAuthnDB(t)

	userID := uuid.New().String()
	mobile := genValidMobile()
	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, mobile, "") })
	insertUserMinimal(t, dbConn, userID, mobile, "sess_cleanup@example.com", "hash", int32(authnentityv1.UserStatus_USER_STATUS_ACTIVE))
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionRepo := NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	// Insert one already-expired session and one valid session.
	expiredID := uuid.New().String()
	validID := uuid.New().String()

	insertSessionMinimal(t, dbConn, expiredID, userID, "SESSION_TYPE_SERVER_SIDE", true, time.Now().Add(-1*time.Hour))
	insertSessionMinimal(t, dbConn, validID, userID, "SESSION_TYPE_SERVER_SIDE", true, time.Now().Add(24*time.Hour))
	t.Cleanup(func() {
		// The expired one may already have been deleted by cleanup; best-effort.
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", expiredID).Delete(map[string]any{}).Error
		_ = dbConn.Table("authn_schema.sessions").Where("session_id = ?", validID).Delete(map[string]any{}).Error
	})

	deleted, err := sessionRepo.CleanupExpiredSessions(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, deleted, int64(1), "at least the one expired session should be deleted")

	// The valid session must still exist.
	got, err := sessionRepo.GetByID(ctx, validID)
	require.NoError(t, err)
	require.Equal(t, validID, got.SessionId)

	// The expired session row should be gone.
	var count int64
	err = dbConn.Raw(
		"SELECT COUNT(1) FROM authn_schema.sessions WHERE session_id = ?",
		expiredID,
	).Scan(&count).Error
	require.NoError(t, err)
	require.Equal(t, int64(0), count, "expired session should have been hard-deleted")
}
