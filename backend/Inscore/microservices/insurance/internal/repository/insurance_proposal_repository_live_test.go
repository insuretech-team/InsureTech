package repository_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/insurance/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	insurerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/insurer/entity/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	policyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/policy/entity/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	insuranceProposalLiveDBOnce sync.Once
	insuranceProposalLiveDB     *gorm.DB
	insuranceProposalLiveDBErr  error
)

func testInsuranceProposalDB(t *testing.T) *gorm.DB {
	t.Helper()
	insuranceProposalLiveDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		_ = env.Load()
		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../configs/database.yaml"
		}
		insuranceProposalLiveDBErr = db.InitializeManagerForService(configPath)
		if insuranceProposalLiveDBErr != nil {
			return
		}
		insuranceProposalLiveDB = db.GetDB()
	})

	if insuranceProposalLiveDBErr != nil {
		t.Skipf("skipping live DB test: %v", insuranceProposalLiveDBErr)
	}
	if insuranceProposalLiveDB == nil {
		t.Skip("skipping live DB test: db is nil")
	}
	return insuranceProposalLiveDB
}

type proposalFixtures struct {
	TenantID    string
	UserID      string
	ProductID   string
	PlanID      string
	QuotationID string
	OrderID     string
	InsurerID   string
}

func TestInsuranceProposalRepository_LiveDB_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	rootDB := testInsuranceProposalDB(t)
	tx := rootDB.Begin()
	require.NoError(t, tx.Error)
	defer tx.Rollback()

	ctx := context.Background()
	fx := createProposalFixtures(t, tx)
	repo := repository.NewInsuranceProposalRepository(tx)

	proposalID := uuid.NewString()
	proposalNumber := fmt.Sprintf("PRP-%s", strings.ToUpper(proposalID[:8]))
	now := timestamppb.New(time.Now().UTC())

	created, err := repo.Create(ctx, &policyv1.InsuranceProposal{
		ProposalId:         proposalID,
		ProposalNumber:     proposalNumber,
		TenantId:           fx.TenantID,
		OrderId:            fx.OrderID,
		QuotationId:        fx.QuotationID,
		CustomerId:         fx.UserID,
		InsurerId:          fx.InsurerID,
		ProductId:          fx.ProductID,
		PlanId:             fx.PlanID,
		ProposedPremium:    &commonv1.Money{Amount: 50000, Currency: "BDT"},
		ProposedSumInsured: &commonv1.Money{Amount: 1000000, Currency: "BDT"},
		Status:             policyv1.ProposalStatus_PROPOSAL_STATUS_SUBMITTED,
		SubmissionPayload:  `{"source":"live-test","step":"create"}`,
		DecisionReason:     "",
		SubmittedAt:        now,
		CorrelationId:      "live-test-correlation",
	})
	require.NoError(t, err)
	require.Equal(t, proposalID, created.ProposalId)
	require.Equal(t, fx.OrderID, created.OrderId)
	require.Equal(t, policyv1.ProposalStatus_PROPOSAL_STATUS_SUBMITTED, created.Status)

	fetched, err := repo.GetByID(ctx, proposalID)
	require.NoError(t, err)
	require.Equal(t, proposalNumber, fetched.ProposalNumber)
	require.Equal(t, fx.InsurerID, fetched.InsurerId)

	fetched.Status = policyv1.ProposalStatus_PROPOSAL_STATUS_REJECTED
	fetched.DecisionReason = "live db reject test"
	fetched.ReviewedByUserId = fx.UserID
	fetched.ReviewedAt = timestamppb.New(time.Now().UTC())
	fetched.InsurerResponsePayload = `{"decision":"rejected"}`

	updated, err := repo.Update(ctx, fetched)
	require.NoError(t, err)
	require.Equal(t, policyv1.ProposalStatus_PROPOSAL_STATUS_REJECTED, updated.Status)
	require.Equal(t, "live db reject test", updated.DecisionReason)

	listed, total, err := repo.List(
		ctx,
		fx.OrderID,
		fx.InsurerID,
		fx.UserID,
		policyv1.ProposalStatus_PROPOSAL_STATUS_REJECTED,
		1,
		10,
	)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int64(1))
	require.NotEmpty(t, listed)

	found := false
	for _, item := range listed {
		if item.ProposalId == proposalID {
			found = true
			break
		}
	}
	require.True(t, found, "expected created proposal in list results")

	require.NoError(t, repo.Delete(ctx, proposalID))
	_, err = repo.GetByID(ctx, proposalID)
	require.Error(t, err)
}

func createProposalFixtures(t *testing.T, tx *gorm.DB) *proposalFixtures {
	t.Helper()
	ctx := context.Background()

	fx := &proposalFixtures{
		TenantID:  uuid.NewString(),
		UserID:    uuid.NewString(),
		ProductID: uuid.NewString(),
		PlanID:    uuid.NewString(),
		OrderID:   uuid.NewString(),
		InsurerID: uuid.NewString(),
	}

	err := tx.WithContext(ctx).Exec(`
		INSERT INTO authn_schema.users
			(user_id, mobile_number, email, password_hash, status, user_type, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'ACTIVE', 'B2C_CUSTOMER', false, NOW(), NOW())`,
		fx.UserID,
		testProposalMobileNumber(fx.UserID),
		fmt.Sprintf("proposal-live-%s@test.local", fx.UserID[:8]),
		"$2a$10$testhashXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	).Error
	require.NoError(t, err)

	err = tx.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.products
			(product_id, product_code, product_name, category, base_premium,
			 min_sum_insured, max_sum_insured, min_tenure_months, max_tenure_months,
			 status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, 'LIFE', 50000, 100000, 10000000, 12, 120, 'ACTIVE', $4, NOW(), NOW())`,
		fx.ProductID,
		testProposalProductCode(fx.ProductID),
		fmt.Sprintf("Proposal Product %s", fx.ProductID[:8]),
		fx.UserID,
	).Error
	require.NoError(t, err)

	err = tx.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.product_plans
			(plan_id, product_id, plan_name, premium_amount, min_sum_insured, max_sum_insured, created_at, updated_at)
		VALUES ($1, $2, $3, 50000, 100000, 10000000, NOW(), NOW())`,
		fx.PlanID,
		fx.ProductID,
		fmt.Sprintf("Proposal Plan %s", fx.PlanID[:8]),
	).Error
	require.NoError(t, err)

	insurerRepo := repository.NewInsurerRepository(tx)
	_, err = insurerRepo.Create(ctx, &insurerv1.Insurer{
		Id:                 fx.InsurerID,
		Code:               fmt.Sprintf("INS-%s", strings.ToUpper(fx.InsurerID[:6])),
		Name:               fmt.Sprintf("Proposal Insurer %s", fx.InsurerID[:8]),
		Type:               insurerv1.InsurerType_INSURER_TYPE_LIFE,
		Status:             insurerv1.InsurerStatus_INSURER_STATUS_ACTIVE,
		TradeLicenseNumber: fmt.Sprintf("TL-%s", fx.InsurerID[:8]),
		TinNumber:          fmt.Sprintf("TIN-%s", fx.InsurerID[:8]),
		IdraLicenseNumber:  fmt.Sprintf("IDRA-%s", fx.InsurerID[:8]),
		ContactInfo: &commonv1.ContactInfo{
			MobileNumber: testProposalMobileNumber(fx.InsurerID),
			Email:        fmt.Sprintf("insurer-%s@test.local", fx.InsurerID[:8]),
		},
		PaidUpCapital: &commonv1.Money{Amount: 100000000, Currency: "BDT"},
	})
	require.NoError(t, err)

	fx.QuotationID = uuid.NewString()
	err = tx.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.quotations
			(quotation_id, business_id, insurer_name, plan_id, insurance_category, status, quotation_number, plan_name, created_by_user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'LIFE', 'APPROVED', $5, $6, $7, NOW(), NOW())`,
		fx.QuotationID,
		fx.UserID,
		fmt.Sprintf("Proposal Insurer %s", fx.InsurerID[:8]),
		fx.PlanID,
		fmt.Sprintf("QUO-%s", strings.ToUpper(fx.QuotationID[:8])),
		fmt.Sprintf("Proposal Plan %s", fx.PlanID[:8]),
		fx.UserID,
	).Error
	require.NoError(t, err)

	err = tx.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.orders
			(order_id, order_number, tenant_id, quotation_id, customer_id, product_id, plan_id,
			 status, total_payable, currency, payment_status, billing_status, fulfillment_status,
			 created_at, updated_at, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
		        'ORDER_STATUS_PAID', 50000, 'BDT', 'ORDER_PAYMENT_STATUS_PAID',
		        'ORDER_BILLING_STATUS_NOT_INVOICED', 'ORDER_FULFILLMENT_STATUS_NOT_STARTED',
		        NOW(), NOW(), NOW())`,
		fx.OrderID,
		fmt.Sprintf("ORD-%s", strings.ToUpper(fx.OrderID[:8])),
		fx.TenantID,
		fx.QuotationID,
		fx.UserID,
		fx.ProductID,
		fx.PlanID,
	).Error
	require.NoError(t, err)

	row := tx.WithContext(ctx).Raw(`
		SELECT order_id
		FROM insurance_schema.orders
		WHERE order_id = $1 AND status = $2`,
		fx.OrderID,
		ordersv1.OrderStatus_ORDER_STATUS_PAID.String(),
	).Row()
	var verifiedOrderID string
	require.NoError(t, row.Scan(&verifiedOrderID))
	require.Equal(t, fx.OrderID, verifiedOrderID)

	return fx
}

func testProposalMobileNumber(id string) string {
	digits := ""
	for _, c := range id {
		if c >= '0' && c <= '9' {
			digits += string(c)
			if len(digits) == 9 {
				break
			}
		}
	}
	for len(digits) < 9 {
		digits += "0"
	}
	return "+88017" + digits[:8]
}

func testProposalProductCode(id string) string {
	letters := ""
	for _, c := range strings.ToUpper(id) {
		if c >= 'A' && c <= 'F' {
			letters += string(c)
			if len(letters) == 3 {
				break
			}
		}
	}
	for len(letters) < 3 {
		letters += "A"
	}

	digits := ""
	for _, c := range id {
		if c >= '0' && c <= '9' {
			digits += string(c)
			if len(digits) == 3 {
				break
			}
		}
	}
	for len(digits) < 3 {
		digits += "0"
	}

	return letters + "-" + digits
}
