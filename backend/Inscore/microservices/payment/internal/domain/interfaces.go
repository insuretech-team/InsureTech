package domain

import (
	"context"
	"errors"
	"time"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	paymententityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/entity/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
)

var ErrNotFound = errors.New("payment record not found")

type PaymentFilters struct {
	UserID        string
	PolicyID      string
	Status        string
	PaymentMethod string
	StartDate     *time.Time
	EndDate       *time.Time
	Limit         int32
	Offset        int
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *paymententityv1.Payment) error
	GetPayment(ctx context.Context, paymentID string) (*paymententityv1.Payment, error)
	GetPaymentByIdempotencyKey(ctx context.Context, idempotencyKey string) (*paymententityv1.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*paymententityv1.Payment, error)
	GetPaymentByProviderReference(ctx context.Context, provider, providerReference string) (*paymententityv1.Payment, error)
	GetPaymentByTranID(ctx context.Context, tranID string) (*paymententityv1.Payment, error)
	ListPayments(ctx context.Context, filters PaymentFilters) ([]*paymententityv1.Payment, int64, error)
	UpdatePayment(ctx context.Context, payment *paymententityv1.Payment) error
	CreateRefund(ctx context.Context, refund *paymententityv1.PaymentRefund) error
	GetRefund(ctx context.Context, refundID string) (*paymententityv1.PaymentRefund, error)
}

type GatewaySessionRequest struct {
	PaymentID        string
	TransactionID    string
	Amount           *commonv1.Money
	Currency         string
	SuccessURL       string
	FailURL          string
	CancelURL        string
	IPNURL           string
	OrderID          string
	TenantID         string
	CustomerName     string
	CustomerEmail    string
	CustomerPhone    string
	CustomerAddr1    string
	CustomerCity     string
	CustomerPostcode string
	CustomerCountry  string
	Metadata         map[string]string
}

type GatewaySessionResponse struct {
	Provider       string
	Status         string
	GatewayPageURL string
	SessionKey     string
	TranID         string
	RawFields      map[string]string
}

type GatewayValidationRequest struct {
	PaymentID     string
	TransactionID string
	SessionKey    string
}

type GatewayValidationResponse struct {
	Provider          string
	Status            string
	TransactionID     string
	ValidationID      string
	BankTransactionID string
	Amount            *commonv1.Money
	CardType          string
	CardBrand         string
	CardIssuer        string
	CardIssuerCountry string
	RiskLevel         string
	RiskTitle         string
	ValidatedAt       time.Time
	RawFields         map[string]string
}

type GatewayRefundRequest struct {
	PaymentID         string
	BankTransactionID string
	Amount            *commonv1.Money
	Reason            string
}

type GatewayRefundResponse struct {
	Provider    string
	Status      string
	RefundRefID string
	RawFields   map[string]string
}

type PaymentGateway interface {
	InitSession(ctx context.Context, req *GatewaySessionRequest) (*GatewaySessionResponse, error)
	ValidatePayment(ctx context.Context, req *GatewayValidationRequest) (*GatewayValidationResponse, error)
	QueryPayment(ctx context.Context, req *GatewayValidationRequest) (*GatewayValidationResponse, error)
	InitiateRefund(ctx context.Context, req *GatewayRefundRequest) (*GatewayRefundResponse, error)
}

type PaymentService interface {
	InitiatePayment(ctx context.Context, req *paymentservicev1.InitiatePaymentRequest) (*paymentservicev1.InitiatePaymentResponse, error)
	VerifyPayment(ctx context.Context, req *paymentservicev1.VerifyPaymentRequest) (*paymentservicev1.VerifyPaymentResponse, error)
	GetPayment(ctx context.Context, req *paymentservicev1.GetPaymentRequest) (*paymentservicev1.GetPaymentResponse, error)
	ListPayments(ctx context.Context, req *paymentservicev1.ListPaymentsRequest) (*paymentservicev1.ListPaymentsResponse, error)
	InitiateRefund(ctx context.Context, req *paymentservicev1.InitiateRefundRequest) (*paymentservicev1.InitiateRefundResponse, error)
	GetRefundStatus(ctx context.Context, req *paymentservicev1.GetRefundStatusRequest) (*paymentservicev1.GetRefundStatusResponse, error)
	ListPaymentMethods(ctx context.Context, req *paymentservicev1.ListPaymentMethodsRequest) (*paymentservicev1.ListPaymentMethodsResponse, error)
	AddPaymentMethod(ctx context.Context, req *paymentservicev1.AddPaymentMethodRequest) (*paymentservicev1.AddPaymentMethodResponse, error)
	ReconcilePayments(ctx context.Context, req *paymentservicev1.ReconcilePaymentsRequest) (*paymentservicev1.ReconcilePaymentsResponse, error)
	// Phase 2 — gateway webhook, manual proof, receipt
	HandleGatewayWebhook(ctx context.Context, req *paymentservicev1.HandleGatewayWebhookRequest) (*paymentservicev1.HandleGatewayWebhookResponse, error)
	GetPaymentByProviderReference(ctx context.Context, req *paymentservicev1.GetPaymentByProviderReferenceRequest) (*paymentservicev1.GetPaymentByProviderReferenceResponse, error)
	SubmitManualPaymentProof(ctx context.Context, req *paymentservicev1.SubmitManualPaymentProofRequest) (*paymentservicev1.SubmitManualPaymentProofResponse, error)
	ReviewManualPayment(ctx context.Context, req *paymentservicev1.ReviewManualPaymentRequest) (*paymentservicev1.ReviewManualPaymentResponse, error)
	GenerateReceipt(ctx context.Context, req *paymentservicev1.GenerateReceiptRequest) (*paymentservicev1.GenerateReceiptResponse, error)
	GetPaymentReceipt(ctx context.Context, req *paymentservicev1.GetPaymentReceiptRequest) (*paymentservicev1.GetPaymentReceiptResponse, error)
}
