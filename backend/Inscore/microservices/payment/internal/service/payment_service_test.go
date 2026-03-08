package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/domain"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	paymententityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/entity/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakePaymentRepo struct {
	paymentsByID          map[string]*paymententityv1.Payment
	paymentsByIdempotency map[string]*paymententityv1.Payment
	refundsByID           map[string]*paymententityv1.PaymentRefund
	createdPayments       []*paymententityv1.Payment
	updatedPayments       []*paymententityv1.Payment
	createdRefunds        []*paymententityv1.PaymentRefund
}

func newFakePaymentRepo() *fakePaymentRepo {
	return &fakePaymentRepo{
		paymentsByID:          map[string]*paymententityv1.Payment{},
		paymentsByIdempotency: map[string]*paymententityv1.Payment{},
		refundsByID:           map[string]*paymententityv1.PaymentRefund{},
	}
}

func (r *fakePaymentRepo) CreatePayment(_ context.Context, payment *paymententityv1.Payment) error {
	cp := clonePayment(payment)
	r.paymentsByID[payment.GetPaymentId()] = cp
	if payment.GetIdempotencyKey() != "" {
		r.paymentsByIdempotency[payment.GetIdempotencyKey()] = cp
	}
	r.createdPayments = append(r.createdPayments, cp)
	return nil
}

func (r *fakePaymentRepo) GetPayment(_ context.Context, paymentID string) (*paymententityv1.Payment, error) {
	p, ok := r.paymentsByID[paymentID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return clonePayment(p), nil
}

func (r *fakePaymentRepo) GetPaymentByIdempotencyKey(_ context.Context, idempotencyKey string) (*paymententityv1.Payment, error) {
	p, ok := r.paymentsByIdempotency[idempotencyKey]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return clonePayment(p), nil
}

func (r *fakePaymentRepo) GetPaymentByOrderID(_ context.Context, orderID string) (*paymententityv1.Payment, error) {
	for _, p := range r.paymentsByID {
		if p.GetOrderId() == orderID {
			return clonePayment(p), nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *fakePaymentRepo) GetPaymentByProviderReference(_ context.Context, provider, providerReference string) (*paymententityv1.Payment, error) {
	for _, p := range r.paymentsByID {
		if !strings.EqualFold(p.GetProvider(), provider) {
			continue
		}
		for _, candidate := range []string{p.GetProviderReference(), p.GetTranId(), p.GetValId(), p.GetSessionKey()} {
			if candidate == providerReference {
				return clonePayment(p), nil
			}
		}
	}
	return nil, domain.ErrNotFound
}

func (r *fakePaymentRepo) GetPaymentByTranID(_ context.Context, tranID string) (*paymententityv1.Payment, error) {
	for _, p := range r.paymentsByID {
		if p.GetTranId() == tranID || p.GetTransactionId() == tranID {
			return clonePayment(p), nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *fakePaymentRepo) ListPayments(_ context.Context, _ domain.PaymentFilters) ([]*paymententityv1.Payment, int64, error) {
	out := make([]*paymententityv1.Payment, 0, len(r.paymentsByID))
	for _, p := range r.paymentsByID {
		out = append(out, clonePayment(p))
	}
	return out, int64(len(out)), nil
}

func (r *fakePaymentRepo) UpdatePayment(_ context.Context, payment *paymententityv1.Payment) error {
	cp := clonePayment(payment)
	r.paymentsByID[payment.GetPaymentId()] = cp
	if payment.GetIdempotencyKey() != "" {
		r.paymentsByIdempotency[payment.GetIdempotencyKey()] = cp
	}
	r.updatedPayments = append(r.updatedPayments, cp)
	return nil
}

func (r *fakePaymentRepo) CreateRefund(_ context.Context, refund *paymententityv1.PaymentRefund) error {
	cr := cloneRefund(refund)
	r.refundsByID[refund.GetRefundId()] = cr
	r.createdRefunds = append(r.createdRefunds, cr)
	return nil
}

func (r *fakePaymentRepo) GetRefund(_ context.Context, refundID string) (*paymententityv1.PaymentRefund, error) {
	refund, ok := r.refundsByID[refundID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return cloneRefund(refund), nil
}

type fakeGateway struct {
	initFn     func(context.Context, *domain.GatewaySessionRequest) (*domain.GatewaySessionResponse, error)
	validateFn func(context.Context, *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error)
	queryFn    func(context.Context, *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error)
	refundFn   func(context.Context, *domain.GatewayRefundRequest) (*domain.GatewayRefundResponse, error)
}

func (g *fakeGateway) InitSession(ctx context.Context, req *domain.GatewaySessionRequest) (*domain.GatewaySessionResponse, error) {
	if g.initFn != nil {
		return g.initFn(ctx, req)
	}
	return nil, errors.New("unexpected InitSession call")
}

func (g *fakeGateway) ValidatePayment(ctx context.Context, req *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error) {
	if g.validateFn != nil {
		return g.validateFn(ctx, req)
	}
	return nil, errors.New("unexpected ValidatePayment call")
}

func (g *fakeGateway) QueryPayment(ctx context.Context, req *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error) {
	if g.queryFn != nil {
		return g.queryFn(ctx, req)
	}
	return nil, errors.New("unexpected QueryPayment call")
}

func (g *fakeGateway) InitiateRefund(ctx context.Context, req *domain.GatewayRefundRequest) (*domain.GatewayRefundResponse, error) {
	if g.refundFn != nil {
		return g.refundFn(ctx, req)
	}
	return nil, errors.New("unexpected InitiateRefund call")
}

func TestPaymentService_InitiatePayment_SSLCommerz_UsesGatewayClient(t *testing.T) {
	repo := newFakePaymentRepo()
	gateway := &fakeGateway{
		initFn: func(_ context.Context, req *domain.GatewaySessionRequest) (*domain.GatewaySessionResponse, error) {
			require.Equal(t, "order-123", req.OrderID)
			require.Equal(t, "tenant-123", req.TenantID)
			require.Equal(t, "Alice Example", req.CustomerName)
			require.Equal(t, "https://example.com/v1/payments/sslcommerz/success?payment_id="+req.PaymentID, req.SuccessURL)
			require.Equal(t, "https://example.com/v1/payments/webhook/sslcommerz?payment_id="+req.PaymentID, req.IPNURL)
			return &domain.GatewaySessionResponse{
				Provider:       "SSLCOMMERZ",
				Status:         "INITIATED",
				GatewayPageURL: "https://sandbox.sslcommerz.com/EasyCheckOut/test",
				SessionKey:     "session-123",
				TranID:         req.TransactionID,
				RawFields: map[string]string{
					"GatewayPageURL": "https://sandbox.sslcommerz.com/EasyCheckOut/test",
					"sessionkey":     "session-123",
					"status":         "INITIATED",
				},
			}, nil
		},
	}
	svc := NewPaymentService(repo, nil, &config.Config{PublicBaseURL: "https://example.com"}, gateway)

	resp, err := svc.InitiatePayment(context.Background(), &paymentservicev1.InitiatePaymentRequest{
		UserId:         uuid.NewString(),
		Amount:         &commonv1.Money{Amount: 12500, Currency: "BDT"},
		Currency:       "BDT",
		PaymentMethod:  "CARD",
		IdempotencyKey: "idem-1",
		Metadata: map[string]string{
			"order_id":          "order-123",
			"tenant_id":         "tenant-123",
			"customer_name":     "Alice Example",
			"customer_email":    "alice@example.com",
			"customer_phone":    "+8801712345678",
			"customer_address":  "Dhaka",
			"customer_city":     "Dhaka",
			"customer_postcode": "1207",
			"customer_country":  "Bangladesh",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetPaymentId())
	require.NotEmpty(t, resp.GetTransactionId())
	require.Equal(t, "https://sandbox.sslcommerz.com/EasyCheckOut/test", resp.GetPaymentUrl())
	require.Equal(t, "SSLCOMMERZ", resp.GetProvider())
	require.Equal(t, resp.GetTransactionId(), resp.GetTranId())
	require.Equal(t, "session-123", resp.GetSessionKey())
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_INITIATED.String(), resp.GetStatus())
	require.Len(t, repo.createdPayments, 1)
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_INITIATED, repo.createdPayments[0].GetStatus())
	require.Equal(t, "SSLCOMMERZ", repo.createdPayments[0].GetGateway())
	require.Len(t, repo.updatedPayments, 1)
	require.Contains(t, repo.updatedPayments[0].GetGatewayResponse(), "session-123")
	require.Equal(t, "session-123", repo.updatedPayments[0].GetSessionKey())
	require.Equal(t, resp.GetTransactionId(), repo.updatedPayments[0].GetTranId())
}

func TestPaymentService_InitiatePayment_SSLCommerz_GatewayFailure_ReturnsError(t *testing.T) {
	repo := newFakePaymentRepo()
	gateway := &fakeGateway{
		initFn: func(context.Context, *domain.GatewaySessionRequest) (*domain.GatewaySessionResponse, error) {
			return nil, errors.New("gateway down")
		},
	}
	svc := NewPaymentService(repo, nil, &config.Config{PublicBaseURL: "https://example.com"}, gateway)

	resp, err := svc.InitiatePayment(context.Background(), &paymentservicev1.InitiatePaymentRequest{
		UserId:         uuid.NewString(),
		Amount:         &commonv1.Money{Amount: 12500, Currency: "BDT"},
		Currency:       "BDT",
		PaymentMethod:  "CARD",
		IdempotencyKey: "idem-2",
		Metadata: map[string]string{
			"order_id":          "order-123",
			"customer_name":     "Alice Example",
			"customer_email":    "alice@example.com",
			"customer_phone":    "+8801712345678",
			"customer_address":  "Dhaka",
			"customer_city":     "Dhaka",
			"customer_postcode": "1207",
			"customer_country":  "Bangladesh",
		},
	})
	require.Nil(t, resp)
	require.ErrorIs(t, err, ErrPaymentFailed)
	require.Len(t, repo.createdPayments, 1)
	require.Len(t, repo.updatedPayments, 1)
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_FAILED, repo.updatedPayments[0].GetStatus())
}

func TestPaymentService_VerifyPayment_SSLCommerz_ValidatesWithGateway(t *testing.T) {
	repo := newFakePaymentRepo()
	payment := &paymententityv1.Payment{
		PaymentId:       uuid.NewString(),
		TransactionId:   "tran-123",
		Type:            paymententityv1.PaymentType_PAYMENT_TYPE_PREMIUM,
		Method:          paymententityv1.PaymentMethod_PAYMENT_METHOD_CARD,
		Status:          paymententityv1.PaymentStatus_PAYMENT_STATUS_INITIATED,
		Amount:          &commonv1.Money{Amount: 12500, Currency: "BDT"},
		Currency:        "BDT",
		PayerId:         uuid.NewString(),
		InitiatedAt:     timestamppb.Now(),
		CreatedAt:       timestamppb.Now(),
		UpdatedAt:       timestamppb.Now(),
		Gateway:         "SSLCOMMERZ",
		GatewayResponse: `{"session_key":"session-123"}`,
		IdempotencyKey:  "idem-verify",
	}
	repo.paymentsByID[payment.GetPaymentId()] = clonePayment(payment)
	repo.paymentsByIdempotency[payment.GetIdempotencyKey()] = clonePayment(payment)

	gateway := &fakeGateway{
		validateFn: func(_ context.Context, req *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error) {
			require.Equal(t, "tran-123", req.TransactionID)
			return &domain.GatewayValidationResponse{
				Provider:          "SSLCOMMERZ",
				Status:            "VALID",
				TransactionID:     "tran-123",
				ValidationID:      "val-123",
				BankTransactionID: "bank-123",
				Amount:            &commonv1.Money{Amount: 12500, Currency: "BDT"},
				RiskLevel:         "0",
				RiskTitle:         "Safe",
				ValidatedAt:       time.Now().UTC(),
				RawFields:         map[string]string{"val_id": "val-123", "bank_tran_id": "bank-123"},
			}, nil
		},
	}
	svc := NewPaymentService(repo, nil, &config.Config{}, gateway)

	resp, err := svc.VerifyPayment(context.Background(), &paymentservicev1.VerifyPaymentRequest{
		PaymentId:      payment.GetPaymentId(),
		TransactionId:  "tran-123",
		PaymentMethod:  "CARD",
		IdempotencyKey: "idem-verify",
	})
	require.NoError(t, err)
	require.True(t, resp.GetVerified())
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS.String(), resp.GetStatus())
	require.Len(t, repo.updatedPayments, 1)
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS, repo.updatedPayments[0].GetStatus())
	require.Contains(t, repo.updatedPayments[0].GetGatewayResponse(), "val-123")
	require.Contains(t, repo.updatedPayments[0].GetGatewayResponse(), "bank-123")
}

func TestPaymentService_HandleGatewayWebhook_Success_VerifiesAndTracksCallback(t *testing.T) {
	repo := newFakePaymentRepo()
	payment := &paymententityv1.Payment{
		PaymentId:          uuid.NewString(),
		TransactionId:      "tran-123",
		TranId:             "tran-123",
		SessionKey:         "session-123",
		Type:               paymententityv1.PaymentType_PAYMENT_TYPE_PREMIUM,
		Method:             paymententityv1.PaymentMethod_PAYMENT_METHOD_CARD,
		Status:             paymententityv1.PaymentStatus_PAYMENT_STATUS_INITIATED,
		Amount:             &commonv1.Money{Amount: 12500, Currency: "BDT"},
		Currency:           "BDT",
		PayerId:            uuid.NewString(),
		Provider:           "SSLCOMMERZ",
		ManualReviewStatus: paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_NOT_REQUIRED,
		InitiatedAt:        timestamppb.Now(),
		CreatedAt:          timestamppb.Now(),
		UpdatedAt:          timestamppb.Now(),
	}
	repo.paymentsByID[payment.GetPaymentId()] = clonePayment(payment)

	gateway := &fakeGateway{
		validateFn: func(_ context.Context, req *domain.GatewayValidationRequest) (*domain.GatewayValidationResponse, error) {
			require.Equal(t, "val-123", req.TransactionID)
			return &domain.GatewayValidationResponse{
				Provider:          "SSLCOMMERZ",
				Status:            "VALID",
				TransactionID:     "tran-123",
				ValidationID:      "val-123",
				BankTransactionID: "bank-123",
				Amount:            &commonv1.Money{Amount: 12500, Currency: "BDT"},
				ValidatedAt:       time.Now().UTC(),
			}, nil
		},
	}
	svc := NewPaymentService(repo, nil, &config.Config{PublicBaseURL: "https://example.com"}, gateway)

	resp, err := svc.HandleGatewayWebhook(context.Background(), &paymentservicev1.HandleGatewayWebhookRequest{
		Provider:   "sslcommerz",
		Headers:    map[string]string{"x-payment-callback-type": "success"},
		RawPayload: []byte("value_a=" + payment.GetPaymentId() + "&val_id=val-123&tran_id=tran-123&sessionkey=session-123"),
	})
	require.NoError(t, err)
	require.True(t, resp.GetAccepted())
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS.String(), resp.GetStatus())

	updated, err := repo.GetPayment(context.Background(), payment.GetPaymentId())
	require.NoError(t, err)
	require.Equal(t, "val-123", updated.GetValId())
	require.Equal(t, "bank-123", updated.GetBankTranId())
	require.NotNil(t, updated.GetCallbackReceivedAt())
}

func TestPaymentService_HandleGatewayWebhook_Cancel_MarksCancelled(t *testing.T) {
	repo := newFakePaymentRepo()
	payment := &paymententityv1.Payment{
		PaymentId:          uuid.NewString(),
		TransactionId:      "tran-123",
		TranId:             "tran-123",
		SessionKey:         "session-123",
		Type:               paymententityv1.PaymentType_PAYMENT_TYPE_PREMIUM,
		Method:             paymententityv1.PaymentMethod_PAYMENT_METHOD_CARD,
		Status:             paymententityv1.PaymentStatus_PAYMENT_STATUS_INITIATED,
		Amount:             &commonv1.Money{Amount: 12500, Currency: "BDT"},
		Currency:           "BDT",
		PayerId:            uuid.NewString(),
		Provider:           "SSLCOMMERZ",
		ManualReviewStatus: paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_NOT_REQUIRED,
	}
	repo.paymentsByID[payment.GetPaymentId()] = clonePayment(payment)

	svc := NewPaymentService(repo, nil, &config.Config{}, nil)
	resp, err := svc.HandleGatewayWebhook(context.Background(), &paymentservicev1.HandleGatewayWebhookRequest{
		Provider:   "sslcommerz",
		Headers:    map[string]string{"x-payment-callback-type": "cancel"},
		RawPayload: []byte("value_a=" + payment.GetPaymentId() + "&tran_id=tran-123"),
	})
	require.NoError(t, err)
	require.True(t, resp.GetAccepted())
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_CANCELLED.String(), resp.GetStatus())
}

func TestPaymentService_GetPaymentByProviderReference_ReturnsPayment(t *testing.T) {
	repo := newFakePaymentRepo()
	payment := &paymententityv1.Payment{
		PaymentId:          uuid.NewString(),
		Provider:           "SSLCOMMERZ",
		ProviderReference:  "session-123",
		TranId:             "tran-123",
		Status:             paymententityv1.PaymentStatus_PAYMENT_STATUS_INITIATED,
		ManualReviewStatus: paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_NOT_REQUIRED,
	}
	repo.paymentsByID[payment.GetPaymentId()] = clonePayment(payment)

	svc := NewPaymentService(repo, nil, &config.Config{}, nil)
	resp, err := svc.GetPaymentByProviderReference(context.Background(), &paymentservicev1.GetPaymentByProviderReferenceRequest{
		Provider:          "sslcommerz",
		ProviderReference: "tran-123",
	})
	require.NoError(t, err)
	require.Equal(t, payment.GetPaymentId(), resp.GetPayment().GetPaymentId())
}

func TestPaymentService_ManualReview_Approve_GeneratesReceipt(t *testing.T) {
	repo := newFakePaymentRepo()
	payment := &paymententityv1.Payment{
		PaymentId:          uuid.NewString(),
		Type:               paymententityv1.PaymentType_PAYMENT_TYPE_PREMIUM,
		Method:             paymententityv1.PaymentMethod_PAYMENT_METHOD_BANK_TRANSFER,
		Status:             paymententityv1.PaymentStatus_PAYMENT_STATUS_PENDING,
		Amount:             &commonv1.Money{Amount: 20000, Currency: "BDT"},
		Currency:           "BDT",
		PayerId:            uuid.NewString(),
		ManualReviewStatus: paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_NOT_REQUIRED,
	}
	repo.paymentsByID[payment.GetPaymentId()] = clonePayment(payment)

	svc := NewPaymentService(repo, nil, &config.Config{PublicBaseURL: "https://example.com"}, nil)
	proofResp, err := svc.SubmitManualPaymentProof(context.Background(), &paymentservicev1.SubmitManualPaymentProofRequest{
		PaymentId:         payment.GetPaymentId(),
		ManualProofFileId: uuid.NewString(),
		SubmittedBy:       uuid.NewString(),
		Notes:             "bank deposit slip",
	})
	require.NoError(t, err)
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_MANUAL_REVIEW_REQUIRED.String(), proofResp.GetStatus())

	reviewerID := uuid.NewString()
	reviewResp, err := svc.ReviewManualPayment(context.Background(), &paymentservicev1.ReviewManualPaymentRequest{
		PaymentId:   payment.GetPaymentId(),
		Approved:    true,
		ReviewedBy:  reviewerID,
		ReviewNotes: "verified with bank statement",
	})
	require.NoError(t, err)
	require.Equal(t, paymententityv1.PaymentStatus_PAYMENT_STATUS_SUCCESS.String(), reviewResp.GetStatus())

	receiptResp, err := svc.GetPaymentReceipt(context.Background(), &paymentservicev1.GetPaymentReceiptRequest{PaymentId: payment.GetPaymentId()})
	require.NoError(t, err)
	require.NotEmpty(t, receiptResp.GetReceiptNumber())
	require.Equal(t, "https://example.com/v1/payments/"+payment.GetPaymentId()+"/receipt", receiptResp.GetReceiptUrl())
}

func clonePayment(in *paymententityv1.Payment) *paymententityv1.Payment {
	if in == nil {
		return nil
	}
	out := *in
	if in.Amount != nil {
		amt := *in.Amount
		out.Amount = &amt
	}
	return &out
}

func cloneRefund(in *paymententityv1.PaymentRefund) *paymententityv1.PaymentRefund {
	if in == nil {
		return nil
	}
	out := *in
	if in.RefundAmount != nil {
		amt := *in.RefundAmount
		out.RefundAmount = &amt
	}
	return &out
}
