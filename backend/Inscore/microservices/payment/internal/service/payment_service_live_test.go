package service

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	_ "github.com/newage-saint/insuretech/backend/inscore/db"
	paymentcfg "github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	paymententityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/entity/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	paymentSvcTestDBOnce sync.Once
	paymentSvcTestDB     *gorm.DB
	paymentSvcTestDBErr  error
)

func testPaymentServiceLiveDB(t *testing.T) *gorm.DB {
	t.Helper()

	paymentSvcTestDBOnce.Do(func() {
		_ = logger.Initialize(logger.NoFileConfig())
		if err := env.Load(); err != nil {
			logger.Warnf("Warning: couldn't load .env: %v", err)
		}

		configPath := os.Getenv("INSCORE_DB_CONFIG")
		if configPath == "" {
			configPath = "../../../../database.yaml"
		}

		paymentSvcTestDBErr = db.InitializeManagerForService(configPath)
		if paymentSvcTestDBErr != nil {
			return
		}
		paymentSvcTestDB = db.GetDB()
		if paymentSvcTestDB != nil {
			paymentSvcTestDB = paymentSvcTestDB.Debug()
		}
	})

	if paymentSvcTestDBErr != nil {
		t.Skipf("skipping live DB test: failed to init DB: %v", paymentSvcTestDBErr)
	}
	if paymentSvcTestDB == nil {
		t.Skip("skipping live DB test: DB is nil")
	}
	return paymentSvcTestDB
}

func insertLivePaymentUser(t *testing.T, dbConn *gorm.DB) string {
	t.Helper()
	userID := uuid.NewString()
	mobileNumber := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	err := dbConn.Exec(
		`INSERT INTO authn_schema.users
		   (user_id, mobile_number, password_hash, status, user_type, created_at, updated_at)
		 VALUES (?, ?, 'test-hash', 'USER_STATUS_ACTIVE', 'USER_TYPE_B2C_CUSTOMER', NOW(), NOW())`,
		userID, mobileNumber,
	).Error
	require.NoError(t, err)
	return userID
}

func cleanupLivePaymentRows(t *testing.T, dbConn *gorm.DB, userID string) {
	t.Helper()
	_ = dbConn.Exec(`DELETE FROM payment_schema.payment_refunds WHERE payment_id IN (SELECT payment_id FROM payment_schema.payments WHERE payer_id = ?)`, userID).Error
	_ = dbConn.Exec(`DELETE FROM payment_schema.payments WHERE payer_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.users WHERE user_id = ?`, userID).Error
}

func TestPaymentService_LiveDB_SSLCommerz_InitiateAndVerify(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	dbConn := testPaymentServiceLiveDB(t)
	userID := insertLivePaymentUser(t, dbConn)
	t.Cleanup(func() { cleanupLivePaymentRows(t, dbConn, userID) })

	repo := repository.NewPaymentRepository(dbConn)
	gateway := &fakeGateway{
		initFn: func(_ context.Context, req *domain.GatewaySessionRequest) (*domain.GatewaySessionResponse, error) {
			return &domain.GatewaySessionResponse{
				Provider:       "SSLCOMMERZ",
				Status:         "INITIATED",
				GatewayPageURL: "https://sandbox.sslcommerz.com/EasyCheckOut/live-test",
				SessionKey:     "session-live-123",
				TranID:         req.TransactionID,
				RawFields: map[string]string{
					"GatewayPageURL": "https://sandbox.sslcommerz.com/EasyCheckOut/live-test",
					"sessionkey":     "session-live-123",
					"status":         "INITIATED",
				},
			}, nil
		},
		validateFn: func(_ context.Context, req *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error) {
			return &domain.GatewayValidationResponse{
				Provider:          "SSLCOMMERZ",
				Status:            "VALID",
				TransactionID:     req.TransactionID,
				ValidationID:      "val-live-123",
				BankTransactionID: "bank-live-123",
				Amount:            &commonv1.Money{Amount: 25000, Currency: "BDT"},
				RiskLevel:         "0",
				RiskTitle:         "Safe",
				ValidatedAt:       time.Now().UTC(),
				RawFields:         map[string]string{"val_id": "val-live-123", "bank_tran_id": "bank-live-123"},
			}, nil
		},
	}
	svc := NewPaymentService(repo, nil, &paymentcfg.Config{PublicBaseURL: "https://example.com"}, gateway)

	ctx := context.Background()
	initResp, err := svc.InitiatePayment(ctx, &paymentservicev1.InitiatePaymentRequest{
		UserId:         userID,
		Amount:         &commonv1.Money{Amount: 25000, Currency: "BDT"},
		Currency:       "BDT",
		PaymentMethod:  "CARD",
		IdempotencyKey: "live-" + uuid.NewString(),
		Metadata: map[string]string{
			"order_id":          "order-live-123",
			"tenant_id":         "00000000-0000-0000-0000-000000000001",
			"customer_name":     "Live Payment User",
			"customer_email":    "live_payment_" + uuid.NewString()[:8] + "@example.com",
			"customer_phone":    "+8801712345678",
			"customer_address":  "Dhaka",
			"customer_city":     "Dhaka",
			"customer_postcode": "1207",
			"customer_country":  "Bangladesh",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, initResp.GetPaymentId())
	require.Equal(t, "https://sandbox.sslcommerz.com/EasyCheckOut/live-test", initResp.GetPaymentUrl())

	payment, err := repo.GetPayment(ctx, initResp.GetPaymentId())
	require.NoError(t, err)
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_INITIATED, payment.GetStatus())
	require.Equal(t, "SSLCOMMERZ", payment.GetGateway())

	verifyResp, err := svc.VerifyPayment(ctx, &paymentservicev1.VerifyPaymentRequest{
		PaymentId:      initResp.GetPaymentId(),
		TransactionId:  initResp.GetTransactionId(),
		PaymentMethod:  "CARD",
		IdempotencyKey: "verify-" + initResp.GetPaymentId(),
	})
	require.NoError(t, err)
	require.True(t, verifyResp.GetVerified())

	verifiedPayment, err := repo.GetPayment(ctx, initResp.GetPaymentId())
	require.NoError(t, err)
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS, verifiedPayment.GetStatus())
	require.Contains(t, verifiedPayment.GetGatewayResponse(), "val-live-123")
	require.Contains(t, verifiedPayment.GetGatewayResponse(), "bank-live-123")
	require.Equal(t, "SSLCOMMERZ", verifiedPayment.GetProvider())
	require.Equal(t, "session-live-123", verifiedPayment.GetSessionKey())
}

func TestPaymentService_LiveDB_ManualReviewAndReceipt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}

	dbConn := testPaymentServiceLiveDB(t)
	payerID := insertLivePaymentUser(t, dbConn)
	reviewerID := insertLivePaymentUser(t, dbConn)
	t.Cleanup(func() {
		cleanupLivePaymentRows(t, dbConn, payerID)
		cleanupLivePaymentRows(t, dbConn, reviewerID)
	})

	repo := repository.NewPaymentRepository(dbConn)
	svc := NewPaymentService(repo, nil, &paymentcfg.Config{PublicBaseURL: "https://example.com"}, nil)

	ctx := context.Background()
	initResp, err := svc.InitiatePayment(ctx, &paymentservicev1.InitiatePaymentRequest{
		UserId:         payerID,
		Amount:         &commonv1.Money{Amount: 18000, Currency: "BDT"},
		Currency:       "BDT",
		PaymentMethod:  "BANK_TRANSFER",
		IdempotencyKey: "manual-" + uuid.NewString(),
	})
	require.NoError(t, err)

	_, err = svc.SubmitManualPaymentProof(ctx, &paymentservicev1.SubmitManualPaymentProofRequest{
		PaymentId:         initResp.GetPaymentId(),
		ManualProofFileId: uuid.NewString(),
		SubmittedBy:       payerID,
		Notes:             "uploaded transfer slip",
	})
	require.NoError(t, err)

	_, err = svc.ReviewManualPayment(ctx, &paymentservicev1.ReviewManualPaymentRequest{
		PaymentId:       initResp.GetPaymentId(),
		Approved:        true,
		ReviewedBy:      reviewerID,
		ReviewNotes:     "matched against bank statement",
		RejectionReason: "",
	})
	require.NoError(t, err)

	receiptResp, err := svc.GetPaymentReceipt(ctx, &paymentservicev1.GetPaymentReceiptRequest{PaymentId: initResp.GetPaymentId()})
	require.NoError(t, err)
	require.NotEmpty(t, receiptResp.GetReceiptNumber())

	payment, err := repo.GetPayment(ctx, initResp.GetPaymentId())
	require.NoError(t, err)
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS, payment.GetStatus())
	require.Equal(t, paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_APPROVED, payment.GetManualReviewStatus())
	require.Equal(t, reviewerID, payment.GetVerifiedBy())
	require.Equal(t, receiptResp.GetReceiptNumber(), payment.GetReceiptNumber())
}
