package service

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	partnerSvcDBOnce sync.Once
	partnerSvcDB     *gorm.DB
	partnerSvcDBErr  error
)

const testEncKey = "0123456789abcdef0123456789abcdef"

func testPartnerServiceDB(t *testing.T) *gorm.DB {
	t.Helper()

	partnerSvcDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}
		partnerSvcDBErr = db.InitializeManagerForService(configPath)
		if partnerSvcDBErr != nil {
			return
		}
		partnerSvcDB = db.GetDB()
	})

	if partnerSvcDBErr != nil {
		t.Skipf("skipping live DB test: %v", partnerSvcDBErr)
	}
	if partnerSvcDB == nil {
		t.Skip("skipping live DB test: db is nil")
	}
	return partnerSvcDB
}

func newLivePartnerService(t *testing.T) (*PartnerService, *gorm.DB) {
	t.Helper()
	dbConn := testPartnerServiceDB(t)
	pRepo := repository.NewPartnerRepository(dbConn, testEncKey)
	aRepo := repository.NewAgentRepository(dbConn, testEncKey)
	cRepo := repository.NewCommissionRepository(dbConn)
	pub := events.NewPublisher(nil, "partner-events") // nil producer: events dropped silently
	svc := NewPartnerService(pRepo, aRepo, cRepo, pub, nil)
	return svc, dbConn
}

func cleanupLivePartner(t *testing.T, db *gorm.DB, partnerID string) {
	t.Helper()
	_ = db.Exec(`DELETE FROM partner_schema.commissions WHERE partner_id = ?`, partnerID).Error
	_ = db.Exec(`DELETE FROM partner_schema.agents WHERE partner_id = ?`, partnerID).Error
	_ = db.Exec(`DELETE FROM partner_schema.partners WHERE partner_id = ?`, partnerID).Error
}

// TestPartnerService_Live_CRUDLifecycle tests Create → Get → Update → List → Delete
func TestPartnerService_Live_CRUDLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLivePartnerService(t)

	partnerID := "svc_partner_" + uuid.New().String()[:8]
	partner := &partnerv1.Partner{
		PartnerId:                 partnerID,
		OrganizationName:          "Service Live Org " + partnerID,
		Type:                      partnerv1.PartnerType_PARTNER_TYPE_AGENT_NETWORK,
		TradeLicense:              "TL-SVC-" + partnerID,
		TinNumber:                 "TIN-SVC-" + partnerID,
		ContactEmail:              partnerID + "@test.com",
		ContactPhone:              "+8801700000099",
		AcquisitionCommissionRate: 7.5,
		RenewalCommissionRate:     4.0,
	}
	t.Cleanup(func() { cleanupLivePartner(t, dbConn, partnerID) })

	// Create
	createResp, err := svc.CreatePartner(ctx, &partnerservicev1.CreatePartnerRequest{
		Partner: partner,
	})
	require.NoError(t, err)
	require.Equal(t, partnerID, createResp.PartnerId)

	// Get
	getResp, err := svc.GetPartner(ctx, &partnerservicev1.GetPartnerRequest{
		PartnerId: partnerID,
	})
	require.NoError(t, err)
	require.Equal(t, "Service Live Org "+partnerID, getResp.Partner.OrganizationName)

	// Update
	updateResp, err := svc.UpdatePartner(ctx, &partnerservicev1.UpdatePartnerRequest{
		PartnerId: partnerID,
		Partner: &partnerv1.Partner{
			OrganizationName: "Updated Service Org",
		},
		UpdateMask: []string{"organization_name"},
	})
	require.NoError(t, err)
	require.Equal(t, "Updated Service Org", updateResp.Partner.OrganizationName)

	// List
	listResp, err := svc.ListPartners(ctx, &partnerservicev1.ListPartnersRequest{
		PageSize: 50,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, listResp.TotalCount, int32(1))

	// Delete
	delResp, err := svc.DeletePartner(ctx, &partnerservicev1.DeletePartnerRequest{
		PartnerId: partnerID,
	})
	require.NoError(t, err)
	require.Contains(t, delResp.Message, "deleted")

	// Verify metrics
	snap := svc.MetricsSnapshot()
	require.GreaterOrEqual(t, snap["partner_created"], int64(1))
	require.GreaterOrEqual(t, snap["partner_fetched"], int64(1))
	require.GreaterOrEqual(t, snap["partner_updated"], int64(1))
}

// TestPartnerService_Live_VerifyAndStatusUpdate tests VerifyPartner and UpdatePartnerStatus
func TestPartnerService_Live_VerifyAndStatusUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, dbConn := newLivePartnerService(t)

	partnerID := "verify_" + uuid.New().String()[:8]
	t.Cleanup(func() { cleanupLivePartner(t, dbConn, partnerID) })

	_, err := svc.CreatePartner(ctx, &partnerservicev1.CreatePartnerRequest{
		Partner: &partnerv1.Partner{
			PartnerId:        partnerID,
			OrganizationName: "Verify Org " + partnerID,
			Type:             partnerv1.PartnerType_PARTNER_TYPE_CORPORATE,
			TradeLicense:     "TL-V-" + partnerID,
			ContactEmail:     partnerID + "@verify.com",
		},
	})
	require.NoError(t, err)

	// Verify partner
	verifyResp, err := svc.VerifyPartner(ctx, &partnerservicev1.VerifyPartnerRequest{
		PartnerId:        partnerID,
		VerificationType: "MANUAL",
		VerificationData: map[string]string{"verified_by": "admin-test"},
	})
	require.NoError(t, err)
	require.True(t, verifyResp.Verified)
	require.Equal(t, "APPROVED", verifyResp.VerificationStatus)
	require.Equal(t, "admin-test", verifyResp.VerifiedBy)

	// Update status to suspended
	statusResp, err := svc.UpdatePartnerStatus(ctx, &partnerservicev1.UpdatePartnerStatusRequest{
		PartnerId: partnerID,
		Status:    "PARTNER_STATUS_SUSPENDED",
	})
	require.NoError(t, err)
	require.Equal(t, "PARTNER_STATUS_SUSPENDED", statusResp.Partner.Status.String())
}

// TestPartnerService_Live_ValidationErrors checks that service rejects invalid input
func TestPartnerService_Live_ValidationErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	ctx := context.Background()
	svc, _ := newLivePartnerService(t)

	// CreatePartner with nil
	_, err := svc.CreatePartner(ctx, &partnerservicev1.CreatePartnerRequest{Partner: nil})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// GetPartner with empty ID
	_, err = svc.GetPartner(ctx, &partnerservicev1.GetPartnerRequest{PartnerId: ""})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// UpdatePartner with empty ID
	_, err = svc.UpdatePartner(ctx, &partnerservicev1.UpdatePartnerRequest{PartnerId: "", Partner: nil})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// DeletePartner with empty ID
	_, err = svc.DeletePartner(ctx, &partnerservicev1.DeletePartnerRequest{PartnerId: ""})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// VerifyPartner with empty fields
	_, err = svc.VerifyPartner(ctx, &partnerservicev1.VerifyPartnerRequest{PartnerId: "", VerificationType: ""})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)

	// UpdatePartnerStatus with invalid status
	_, err = svc.UpdatePartnerStatus(ctx, &partnerservicev1.UpdatePartnerStatusRequest{
		PartnerId: "x",
		Status:    "TOTALLY_INVALID",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidArgument)
}
