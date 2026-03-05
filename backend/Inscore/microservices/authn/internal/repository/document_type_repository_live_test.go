package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDocumentTypeRepository_LiveDB_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testAuthnDB(t)
	repo := NewDocumentTypeRepository(dbConn)

	dtID := uuid.New().String()
	code := "TESTCODE_" + uuid.New().String()[:8]

	// Pre-cleanup and deferred cleanup
	_ = repo.Delete(ctx, dtID)
	t.Cleanup(func() { _ = repo.Delete(ctx, dtID) })

	// --- Create ---
	desc := "Integration test document type"
	dt := &authnentityv1.DocumentType{
		DocumentTypeId: dtID,
		Code:           code,
		Name:           "Test Document Type",
		Description:    desc,
		IsActive:       true,
	}
	require.NoError(t, repo.Create(ctx, dt))

	// --- GetByID ---
	got, err := repo.GetByID(ctx, dtID)
	require.NoError(t, err)
	require.Equal(t, dtID, got.DocumentTypeId)
	require.Equal(t, code, got.Code)
	require.Equal(t, "Test Document Type", got.Name)
	require.True(t, got.IsActive)

	// --- GetByCode ---
	byCode, err := repo.GetByCode(ctx, code)
	require.NoError(t, err)
	require.Equal(t, dtID, byCode.DocumentTypeId)
	require.Equal(t, code, byCode.Code)

	// --- ListActive (verify created doc appears) ---
	activeList, err := repo.ListActive(ctx)
	require.NoError(t, err)
	found := false
	for _, item := range activeList {
		if item.DocumentTypeId == dtID {
			found = true
			break
		}
	}
	require.True(t, found, "created document type should appear in ListActive")

	// --- SetActive(false) ---
	require.NoError(t, repo.SetActive(ctx, dtID, false))

	deactivated, err := repo.GetByID(ctx, dtID)
	require.NoError(t, err)
	require.False(t, deactivated.IsActive)

	// Confirm it no longer appears in ListActive
	activeList2, err := repo.ListActive(ctx)
	require.NoError(t, err)
	for _, item := range activeList2 {
		require.NotEqual(t, dtID, item.DocumentTypeId, "deactivated doc type should not appear in ListActive")
	}

	// --- Delete ---
	require.NoError(t, repo.Delete(ctx, dtID))

	_, err = repo.GetByID(ctx, dtID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
