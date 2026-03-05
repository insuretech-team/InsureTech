package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	voicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/voice/entity/v1"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestVoiceSessionRepository_LiveDB_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthnDB(t)
	repo := NewVoiceSessionRepository(dbConn)

	// Create a user to satisfy any FK constraint on user_id
	userID := uuid.New().String()
	mobile := genValidMobile()

	cleanupAuthnUser(ctx, t, dbConn, mobile, "")
	insertUserMinimal(t, dbConn, userID, mobile, "voice_test@example.com", "hash_voice", 0)
	t.Cleanup(func() { cleanupAuthnUser(ctx, t, dbConn, "", userID) })

	sessionID := uuid.New().String()
	externalSessionID := "ext-" + uuid.New().String()

	// Pre-cleanup and deferred cleanup
	_ = repo.Delete(ctx, sessionID)
	t.Cleanup(func() { _ = repo.Delete(ctx, sessionID) })

	// --- Create ---
	session := &voicev1.VoiceSession{
		Id:          sessionID,
		SessionId:   externalSessionID,
		UserId:      userID,
		PhoneNumber: mobile,
		Language:    "bn",
		Status:      voicev1.SessionStatus_SESSION_STATUS_ACTIVE,
		Intent:      "inquiry",
	}
	require.NoError(t, repo.Create(ctx, session))

	// --- GetByID ---
	got, err := repo.GetByID(ctx, sessionID)
	require.NoError(t, err)
	require.Equal(t, sessionID, got.Id)
	require.Equal(t, externalSessionID, got.SessionId)
	require.Equal(t, userID, got.UserId)
	require.Equal(t, "bn", got.Language)
	require.Equal(t, voicev1.SessionStatus_SESSION_STATUS_ACTIVE, got.Status)

	// --- GetByExternalSessionID ---
	byExt, err := repo.GetByExternalSessionID(ctx, externalSessionID)
	require.NoError(t, err)
	require.Equal(t, sessionID, byExt.Id)
	require.Equal(t, externalSessionID, byExt.SessionId)

	// --- ListByUser ---
	list, err := repo.ListByUser(ctx, userID, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, list)
	found := false
	for _, item := range list {
		if item.Id == sessionID {
			found = true
			break
		}
	}
	require.True(t, found, "created voice session should appear in ListByUser")

	// --- Complete (COMPLETED) ---
	endedAt := time.Now().UTC()
	duration := int32(120)
	require.NoError(t, repo.Complete(ctx, sessionID, voicev1.SessionStatus_SESSION_STATUS_COMPLETED, endedAt, &duration))

	// --- GetByID - verify status is COMPLETED ---
	completed, err := repo.GetByID(ctx, sessionID)
	require.NoError(t, err)
	require.Equal(t, voicev1.SessionStatus_SESSION_STATUS_COMPLETED, completed.Status)

	// --- Delete ---
	require.NoError(t, repo.Delete(ctx, sessionID))

	_, err = repo.GetByID(ctx, sessionID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
