package repository

import (
	"context"
	"testing"

	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	"github.com/stretchr/testify/require"
)

func TestPartnerRepository_LiveDB_CreateGetUpdateList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testPartnerDB(t)
	// Use a fixed test encryption key (32 bytes).
	encKey := "0123456789abcdef0123456789abcdef"
	repo := NewPartnerRepository(dbConn, encKey)

	partnerID := newPartnerLiveID("partner")
	partner := &partnerv1.Partner{
		PartnerId:                 partnerID,
		OrganizationName:          "Live Test Org " + partnerID,
		Type:                      partnerv1.PartnerType_PARTNER_TYPE_CORPORATE,
		TradeLicense:              "TL-" + partnerID,
		TinNumber:                 "TIN-" + partnerID,
		BankAccount:               "1234567890",
		BankName:                  "Test Bank",
		BankBranch:                "Main Branch",
		ContactEmail:              partnerID + "@test.com",
		ContactPhone:              "+8801700000001",
		AcquisitionCommissionRate: 5.0,
		RenewalCommissionRate:     3.0,
		ClaimsAssistanceRate:      1.5,
	}
	t.Cleanup(func() { cleanupPartner(t, dbConn, partnerID) })

	// Create
	err := repo.Create(ctx, partner)
	require.NoError(t, err)
	require.NotEmpty(t, partner.PartnerId)

	// GetByID (verifies PII decryption round-trip)
	fetched, err := repo.GetByID(ctx, partnerID)
	require.NoError(t, err)
	require.Equal(t, partnerID, fetched.PartnerId)
	require.Equal(t, "Live Test Org "+partnerID, fetched.OrganizationName)
	require.Equal(t, "1234567890", fetched.BankAccount)      // decrypted
	require.Equal(t, "+8801700000001", fetched.ContactPhone) // decrypted

	// Update
	updatedPartner := &partnerv1.Partner{
		OrganizationName: "Updated Org Name",
	}
	err = repo.Update(ctx, partnerID, updatedPartner, []string{"organization_name"})
	require.NoError(t, err)

	reFetched, err := repo.GetByID(ctx, partnerID)
	require.NoError(t, err)
	require.Equal(t, "Updated Org Name", reFetched.OrganizationName)

	// ListWithFilters
	partners, total, err := repo.ListWithFilters(ctx, 50, 0, "", "")
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int32(1))
	require.NotEmpty(t, partners)

	// UpdateStatus
	err = repo.UpdateStatus(ctx, partnerID, partnerv1.PartnerStatus_PARTNER_STATUS_ACTIVE)
	require.NoError(t, err)

	active, err := repo.GetByID(ctx, partnerID)
	require.NoError(t, err)
	require.Equal(t, "PARTNER_STATUS_ACTIVE", active.Status.String())

	// SoftDelete
	err = repo.SoftDelete(ctx, partnerID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, partnerID)
	require.Error(t, err) // should not find deleted partner
}

func TestPartnerRepository_LiveDB_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	dbConn := testPartnerDB(t)
	repo := NewPartnerRepository(dbConn, "0123456789abcdef0123456789abcdef")

	_, err := repo.GetByID(ctx, "nonexistent-partner")
	require.Error(t, err)
}
