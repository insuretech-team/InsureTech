package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/payment/internal/domain"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	paymententityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	paymentsTable = "payment_schema.payments"
	refundsTable  = "payment_schema.payment_refunds"
)

type PaymentRepository struct {
	db *gorm.DB
}

type paymentRow struct {
	PaymentID             string     `gorm:"column:payment_id"`
	TransactionID         string     `gorm:"column:transaction_id"`
	TigerbeetleTransferID string     `gorm:"column:tigerbeetle_transfer_id"`
	PolicyID              string     `gorm:"column:policy_id"`
	ClaimID               string     `gorm:"column:claim_id"`
	Type                  string     `gorm:"column:type"`
	Method                string     `gorm:"column:method"`
	Status                string     `gorm:"column:status"`
	Amount                int64      `gorm:"column:amount"`
	Currency              string     `gorm:"column:currency"`
	DecimalAmount         float64    `gorm:"column:decimal_amount"`
	PayerID               string     `gorm:"column:payer_id"`
	PayeeID               string     `gorm:"column:payee_id"`
	InitiatedAt           *time.Time `gorm:"column:initiated_at"`
	CompletedAt           *time.Time `gorm:"column:completed_at"`
	CreatedAt             *time.Time `gorm:"column:created_at"`
	UpdatedAt             *time.Time `gorm:"column:updated_at"`
	Gateway               string     `gorm:"column:gateway"`
	GatewayResponse       string     `gorm:"column:gateway_response"`
	ReceiptURL            string     `gorm:"column:receipt_url"`
	RetryCount            int32      `gorm:"column:retry_count"`
	FailureReason         string     `gorm:"column:failure_reason"`
	IdempotencyKey        string     `gorm:"column:idempotency_key"`
	OrderID               string     `gorm:"column:order_id"`
	InvoiceID             string     `gorm:"column:invoice_id"`
	TenantID              string     `gorm:"column:tenant_id"`
	CustomerID            string     `gorm:"column:customer_id"`
	OrganisationID        string     `gorm:"column:organisation_id"`
	PurchaseOrderID       string     `gorm:"column:purchase_order_id"`
	Provider              string     `gorm:"column:provider"`
	ProviderReference     string     `gorm:"column:provider_reference"`
	TranID                string     `gorm:"column:tran_id"`
	ValID                 string     `gorm:"column:val_id"`
	SessionKey            string     `gorm:"column:session_key"`
	BankTranID            string     `gorm:"column:bank_tran_id"`
	CardType              string     `gorm:"column:card_type"`
	CardBrand             string     `gorm:"column:card_brand"`
	CardIssuer            string     `gorm:"column:card_issuer"`
	CardIssuerCountry     string     `gorm:"column:card_issuer_country"`
	ValidatedAt           *time.Time `gorm:"column:validated_at"`
	ValidationStatus      string     `gorm:"column:validation_status"`
	RiskLevel             string     `gorm:"column:risk_level"`
	RiskTitle             string     `gorm:"column:risk_title"`
	CallbackReceivedAt    *time.Time `gorm:"column:callback_received_at"`
	IPNReceivedAt         *time.Time `gorm:"column:ipn_received_at"`
	ManualReviewStatus    string     `gorm:"column:manual_review_status"`
	ManualProofFileID     string     `gorm:"column:manual_proof_file_id"`
	VerifiedBy            string     `gorm:"column:verified_by"`
	VerifiedAt            *time.Time `gorm:"column:verified_at"`
	RejectionReason       string     `gorm:"column:rejection_reason"`
	ReceiptNumber         string     `gorm:"column:receipt_number"`
	ReceiptDocumentID     string     `gorm:"column:receipt_document_id"`
	ReceiptFileID         string     `gorm:"column:receipt_file_id"`
	LedgerTransactionID   string     `gorm:"column:ledger_transaction_id"`
}

type refundRow struct {
	RefundID        string     `gorm:"column:refund_id"`
	PaymentID       string     `gorm:"column:payment_id"`
	RefundPaymentID string     `gorm:"column:refund_payment_id"`
	RefundAmount    int64      `gorm:"column:refund_amount"`
	Currency        string     `gorm:"column:currency"`
	DecimalAmount   float64    `gorm:"column:decimal_amount"`
	Reason          string     `gorm:"column:reason"`
	Status          string     `gorm:"column:status"`
	ApprovedBy      string     `gorm:"column:approved_by"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	ProcessedAt     *time.Time `gorm:"column:processed_at"`
	CreatedAt       *time.Time `gorm:"column:created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *paymententityv1.Payment) error {
	payload := map[string]any{
		"payment_id":              payment.GetPaymentId(),
		"transaction_id":          payment.GetTransactionId(),
		"tigerbeetle_transfer_id": nullableValue(payment.GetTigerbeetleTransferId()),
		"policy_id":               nullableValue(payment.GetPolicyId()),
		"claim_id":                nullableValue(payment.GetClaimId()),
		"type":                    payment.GetType().String(),
		"method":                  payment.GetMethod().String(),
		"status":                  payment.GetStatus().String(),
		"amount":                  paymentAmount(payment),
		"currency":                paymentCurrency(payment),
		"payer_id":                payment.GetPayerId(),
		"payee_id":                nullableValue(payment.GetPayeeId()),
		"initiated_at":            tsToTime(payment.GetInitiatedAt()),
		"completed_at":            tsToTime(payment.GetCompletedAt()),
		"created_at":              tsToTime(payment.GetCreatedAt()),
		"updated_at":              tsToTime(payment.GetUpdatedAt()),
		"gateway":                 payment.GetGateway(),
		"gateway_response":        payment.GetGatewayResponse(),
		"receipt_url":             payment.GetReceiptUrl(),
		"retry_count":             payment.GetRetryCount(),
		"failure_reason":          payment.GetFailureReason(),
		"idempotency_key":         payment.GetIdempotencyKey(),
		"order_id":                nullableValue(payment.GetOrderId()),
		"invoice_id":              nullableValue(payment.GetInvoiceId()),
		"tenant_id":               nullableValue(payment.GetTenantId()),
		"customer_id":             nullableValue(payment.GetCustomerId()),
		"organisation_id":         nullableValue(payment.GetOrganisationId()),
		"purchase_order_id":       nullableValue(payment.GetPurchaseOrderId()),
		"provider":                nullableValue(payment.GetProvider()),
		"provider_reference":      nullableValue(payment.GetProviderReference()),
		"tran_id":                 nullableValue(payment.GetTranId()),
		"val_id":                  nullableValue(payment.GetValId()),
		"session_key":             nullableValue(payment.GetSessionKey()),
		"bank_tran_id":            nullableValue(payment.GetBankTranId()),
		"card_type":               nullableValue(payment.GetCardType()),
		"card_brand":              nullableValue(payment.GetCardBrand()),
		"card_issuer":             nullableValue(payment.GetCardIssuer()),
		"card_issuer_country":     nullableValue(payment.GetCardIssuerCountry()),
		"validated_at":            tsToTime(payment.GetValidatedAt()),
		"validation_status":       nullableValue(payment.GetValidationStatus()),
		"risk_level":              nullableValue(payment.GetRiskLevel()),
		"risk_title":              nullableValue(payment.GetRiskTitle()),
		"callback_received_at":    tsToTime(payment.GetCallbackReceivedAt()),
		"ipn_received_at":         tsToTime(payment.GetIpnReceivedAt()),
		"manual_review_status":    payment.GetManualReviewStatus().String(),
		"manual_proof_file_id":    nullableValue(payment.GetManualProofFileId()),
		"verified_by":             nullableValue(payment.GetVerifiedBy()),
		"verified_at":             tsToTime(payment.GetVerifiedAt()),
		"rejection_reason":        nullableValue(payment.GetRejectionReason()),
		"receipt_number":          nullableValue(payment.GetReceiptNumber()),
		"receipt_document_id":     nullableValue(payment.GetReceiptDocumentId()),
		"receipt_file_id":         nullableValue(payment.GetReceiptFileId()),
		"ledger_transaction_id":   nullableValue(payment.GetLedgerTransactionId()),
	}

	return r.db.WithContext(ctx).Table(paymentsTable).Create(payload).Error
}

func (r *PaymentRepository) GetPayment(ctx context.Context, paymentID string) (*paymententityv1.Payment, error) {
	var row paymentRow
	err := r.db.WithContext(ctx).Table(paymentsTable).Where("payment_id = ?", paymentID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToPayment(&row), nil
}

func (r *PaymentRepository) GetPaymentByIdempotencyKey(ctx context.Context, idempotencyKey string) (*paymententityv1.Payment, error) {
	var row paymentRow
	err := r.db.WithContext(ctx).Table(paymentsTable).Where("idempotency_key = ?", idempotencyKey).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToPayment(&row), nil
}

func (r *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID string) (*paymententityv1.Payment, error) {
	var row paymentRow
	err := r.db.WithContext(ctx).Table(paymentsTable).Where("order_id = ?", orderID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToPayment(&row), nil
}

func (r *PaymentRepository) GetPaymentByProviderReference(ctx context.Context, provider, providerReference string) (*paymententityv1.Payment, error) {
	var row paymentRow
	err := r.db.WithContext(ctx).
		Table(paymentsTable).
		Where("provider = ?", strings.TrimSpace(provider)).
		Where("(provider_reference = ? OR tran_id = ? OR val_id = ? OR session_key = ?)", providerReference, providerReference, providerReference, providerReference).
		Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToPayment(&row), nil
}

func (r *PaymentRepository) GetPaymentByTranID(ctx context.Context, tranID string) (*paymententityv1.Payment, error) {
	var row paymentRow
	err := r.db.WithContext(ctx).Table(paymentsTable).Where("tran_id = ? OR transaction_id = ?", tranID, tranID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToPayment(&row), nil
}

func (r *PaymentRepository) ListPayments(ctx context.Context, filters domain.PaymentFilters) ([]*paymententityv1.Payment, int64, error) {
	query := r.db.WithContext(ctx).Table(paymentsTable)
	if filters.UserID != "" {
		query = query.Where("payer_id = ?", filters.UserID)
	}
	if filters.PolicyID != "" {
		query = query.Where("policy_id = ?", filters.PolicyID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", strings.ToUpper(filters.Status))
	}
	if filters.PaymentMethod != "" {
		query = query.Where("method = ?", strings.ToUpper(filters.PaymentMethod))
	}
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", *filters.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filters.Limit
	if limit <= 0 {
		limit = 20
	}

	var rows []paymentRow
	err := query.Order("created_at DESC").Limit(int(limit)).Offset(filters.Offset).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	payments := make([]*paymententityv1.Payment, 0, len(rows))
	for i := range rows {
		payments = append(payments, rowToPayment(&rows[i]))
	}
	return payments, total, nil
}

func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *paymententityv1.Payment) error {
	updates := map[string]any{
		"transaction_id":          payment.GetTransactionId(),
		"tigerbeetle_transfer_id": nullableValue(payment.GetTigerbeetleTransferId()),
		"policy_id":               nullableValue(payment.GetPolicyId()),
		"claim_id":                nullableValue(payment.GetClaimId()),
		"type":                    payment.GetType().String(),
		"method":                  payment.GetMethod().String(),
		"status":                  payment.GetStatus().String(),
		"amount":                  paymentAmount(payment),
		"currency":                paymentCurrency(payment),
		"payer_id":                payment.GetPayerId(),
		"payee_id":                nullableValue(payment.GetPayeeId()),
		"initiated_at":            tsToTime(payment.GetInitiatedAt()),
		"completed_at":            tsToTime(payment.GetCompletedAt()),
		"updated_at":              tsToTime(payment.GetUpdatedAt()),
		"gateway":                 payment.GetGateway(),
		"gateway_response":        payment.GetGatewayResponse(),
		"receipt_url":             payment.GetReceiptUrl(),
		"retry_count":             payment.GetRetryCount(),
		"failure_reason":          payment.GetFailureReason(),
		"idempotency_key":         payment.GetIdempotencyKey(),
		"order_id":                nullableValue(payment.GetOrderId()),
		"invoice_id":              nullableValue(payment.GetInvoiceId()),
		"tenant_id":               nullableValue(payment.GetTenantId()),
		"customer_id":             nullableValue(payment.GetCustomerId()),
		"organisation_id":         nullableValue(payment.GetOrganisationId()),
		"purchase_order_id":       nullableValue(payment.GetPurchaseOrderId()),
		"provider":                nullableValue(payment.GetProvider()),
		"provider_reference":      nullableValue(payment.GetProviderReference()),
		"tran_id":                 nullableValue(payment.GetTranId()),
		"val_id":                  nullableValue(payment.GetValId()),
		"session_key":             nullableValue(payment.GetSessionKey()),
		"bank_tran_id":            nullableValue(payment.GetBankTranId()),
		"card_type":               nullableValue(payment.GetCardType()),
		"card_brand":              nullableValue(payment.GetCardBrand()),
		"card_issuer":             nullableValue(payment.GetCardIssuer()),
		"card_issuer_country":     nullableValue(payment.GetCardIssuerCountry()),
		"validated_at":            tsToTime(payment.GetValidatedAt()),
		"validation_status":       nullableValue(payment.GetValidationStatus()),
		"risk_level":              nullableValue(payment.GetRiskLevel()),
		"risk_title":              nullableValue(payment.GetRiskTitle()),
		"callback_received_at":    tsToTime(payment.GetCallbackReceivedAt()),
		"ipn_received_at":         tsToTime(payment.GetIpnReceivedAt()),
		"manual_review_status":    payment.GetManualReviewStatus().String(),
		"manual_proof_file_id":    nullableValue(payment.GetManualProofFileId()),
		"verified_by":             nullableValue(payment.GetVerifiedBy()),
		"verified_at":             tsToTime(payment.GetVerifiedAt()),
		"rejection_reason":        nullableValue(payment.GetRejectionReason()),
		"receipt_number":          nullableValue(payment.GetReceiptNumber()),
		"receipt_document_id":     nullableValue(payment.GetReceiptDocumentId()),
		"receipt_file_id":         nullableValue(payment.GetReceiptFileId()),
		"ledger_transaction_id":   nullableValue(payment.GetLedgerTransactionId()),
	}

	return r.db.WithContext(ctx).Table(paymentsTable).Where("payment_id = ?", payment.GetPaymentId()).Updates(updates).Error
}

func (r *PaymentRepository) CreateRefund(ctx context.Context, refund *paymententityv1.PaymentRefund) error {
	payload := map[string]any{
		"refund_id":         refund.GetRefundId(),
		"payment_id":        refund.GetPaymentId(),
		"refund_payment_id": nullableValue(refund.GetRefundPaymentId()),
		"refund_amount":     refundAmount(refund),
		"currency":          refundCurrency(refund),
		"reason":            refund.GetReason(),
		"status":            refund.GetStatus().String(),
		"approved_by":       nullableValue(refund.GetApprovedBy()),
		"approved_at":       tsToTime(refund.GetApprovedAt()),
		"processed_at":      tsToTime(refund.GetProcessedAt()),
		"created_at":        tsToTime(refund.GetCreatedAt()),
		"updated_at":        tsToTime(refund.GetUpdatedAt()),
	}

	return r.db.WithContext(ctx).Table(refundsTable).Create(payload).Error
}

func (r *PaymentRepository) GetRefund(ctx context.Context, refundID string) (*paymententityv1.PaymentRefund, error) {
	var row refundRow
	err := r.db.WithContext(ctx).Table(refundsTable).Where("refund_id = ?", refundID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToRefund(&row), nil
}

func rowToPayment(row *paymentRow) *paymententityv1.Payment {
	return &paymententityv1.Payment{
		PaymentId:             row.PaymentID,
		TransactionId:         row.TransactionID,
		TigerbeetleTransferId: row.TigerbeetleTransferID,
		PolicyId:              row.PolicyID,
		ClaimId:               row.ClaimID,
		Type:                  parsePaymentType(row.Type),
		Method:                parsePaymentMethod(row.Method),
		Status:                parsePaymentStatus(row.Status),
		Amount: &commonv1.Money{
			Amount:        row.Amount,
			Currency:      row.Currency,
			DecimalAmount: coalesceDecimalAmount(row.DecimalAmount, row.Amount),
		},
		Currency:            row.Currency,
		PayerId:             row.PayerID,
		PayeeId:             row.PayeeID,
		InitiatedAt:         timeToTS(row.InitiatedAt),
		CompletedAt:         timeToTS(row.CompletedAt),
		CreatedAt:           timeToTS(row.CreatedAt),
		UpdatedAt:           timeToTS(row.UpdatedAt),
		Gateway:             row.Gateway,
		GatewayResponse:     row.GatewayResponse,
		ReceiptUrl:          row.ReceiptURL,
		RetryCount:          row.RetryCount,
		FailureReason:       row.FailureReason,
		IdempotencyKey:      row.IdempotencyKey,
		OrderId:             row.OrderID,
		InvoiceId:           row.InvoiceID,
		TenantId:            row.TenantID,
		CustomerId:          row.CustomerID,
		OrganisationId:      row.OrganisationID,
		PurchaseOrderId:     row.PurchaseOrderID,
		Provider:            row.Provider,
		ProviderReference:   row.ProviderReference,
		TranId:              row.TranID,
		ValId:               row.ValID,
		SessionKey:          row.SessionKey,
		BankTranId:          row.BankTranID,
		CardType:            row.CardType,
		CardBrand:           row.CardBrand,
		CardIssuer:          row.CardIssuer,
		CardIssuerCountry:   row.CardIssuerCountry,
		ValidatedAt:         timeToTS(row.ValidatedAt),
		ValidationStatus:    row.ValidationStatus,
		RiskLevel:           row.RiskLevel,
		RiskTitle:           row.RiskTitle,
		CallbackReceivedAt:  timeToTS(row.CallbackReceivedAt),
		IpnReceivedAt:       timeToTS(row.IPNReceivedAt),
		ManualReviewStatus:  parseManualReviewStatus(row.ManualReviewStatus),
		ManualProofFileId:   row.ManualProofFileID,
		VerifiedBy:          row.VerifiedBy,
		VerifiedAt:          timeToTS(row.VerifiedAt),
		RejectionReason:     row.RejectionReason,
		ReceiptNumber:       row.ReceiptNumber,
		ReceiptDocumentId:   row.ReceiptDocumentID,
		ReceiptFileId:       row.ReceiptFileID,
		LedgerTransactionId: row.LedgerTransactionID,
	}
}

func rowToRefund(row *refundRow) *paymententityv1.PaymentRefund {
	return &paymententityv1.PaymentRefund{
		RefundId:        row.RefundID,
		PaymentId:       row.PaymentID,
		RefundPaymentId: row.RefundPaymentID,
		RefundAmount: &commonv1.Money{
			Amount:        row.RefundAmount,
			Currency:      row.Currency,
			DecimalAmount: coalesceDecimalAmount(row.DecimalAmount, row.RefundAmount),
		},
		Reason:      row.Reason,
		Status:      parseRefundStatus(row.Status),
		ApprovedBy:  row.ApprovedBy,
		ApprovedAt:  timeToTS(row.ApprovedAt),
		ProcessedAt: timeToTS(row.ProcessedAt),
		CreatedAt:   timeToTS(row.CreatedAt),
		UpdatedAt:   timeToTS(row.UpdatedAt),
	}
}

func paymentAmount(payment *paymententityv1.Payment) int64 {
	if payment.GetAmount() == nil {
		return 0
	}
	return payment.GetAmount().GetAmount()
}

func paymentCurrency(payment *paymententityv1.Payment) string {
	if payment.GetAmount() != nil && payment.GetAmount().GetCurrency() != "" {
		return payment.GetAmount().GetCurrency()
	}
	return payment.GetCurrency()
}

func refundAmount(refund *paymententityv1.PaymentRefund) int64 {
	if refund.GetRefundAmount() == nil {
		return 0
	}
	return refund.GetRefundAmount().GetAmount()
}

func refundCurrency(refund *paymententityv1.PaymentRefund) string {
	if refund.GetRefundAmount() == nil {
		return ""
	}
	return refund.GetRefundAmount().GetCurrency()
}

func tsToTime(value *timestamppb.Timestamp) *time.Time {
	if value == nil {
		return nil
	}
	t := value.AsTime().UTC()
	return &t
}

func timeToTS(value *time.Time) *timestamppb.Timestamp {
	if value == nil || value.IsZero() {
		return nil
	}
	return timestamppb.New(value.UTC())
}

func parsePaymentType(value string) paymententityv1.PaymentType {
	if parsed, ok := paymententityv1.PaymentType_value[strings.ToUpper(value)]; ok {
		return paymententityv1.PaymentType(parsed)
	}
	return paymententityv1.PaymentType_PAYMENT_TYPE_UNSPECIFIED
}

func parsePaymentMethod(value string) paymententityv1.PaymentMethod {
	if parsed, ok := paymententityv1.PaymentMethod_value[strings.ToUpper(value)]; ok {
		return paymententityv1.PaymentMethod(parsed)
	}
	return paymententityv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
}

func parsePaymentStatus(value string) paymententityv1.PaymentStatus {
	if parsed, ok := paymententityv1.PaymentStatus_value[strings.ToUpper(value)]; ok {
		return paymententityv1.PaymentStatus(parsed)
	}
	return paymententityv1.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
}

func parseRefundStatus(value string) paymententityv1.PaymentRefundStatus {
	if parsed, ok := paymententityv1.PaymentRefundStatus_value[strings.ToUpper(value)]; ok {
		return paymententityv1.PaymentRefundStatus(parsed)
	}
	return paymententityv1.PaymentRefundStatus_PAYMENT_REFUND_STATUS_UNSPECIFIED
}

func parseManualReviewStatus(value string) paymententityv1.ManualReviewStatus {
	if parsed, ok := paymententityv1.ManualReviewStatus_value[strings.ToUpper(value)]; ok {
		return paymententityv1.ManualReviewStatus(parsed)
	}
	return paymententityv1.ManualReviewStatus_MANUAL_REVIEW_STATUS_UNSPECIFIED
}

func MarshalGatewayResponse(value map[string]string) string {
	if len(value) == 0 {
		return ""
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(payload)
}

func UnmarshalGatewayResponse(value string) map[string]string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	result := map[string]string{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return map[string]string{"raw": value}
	}
	return result
}

func coalesceDecimalAmount(decimal float64, amount int64) float64 {
	if decimal != 0 {
		return decimal
	}
	if amount == 0 {
		return 0
	}
	return float64(amount) / 100
}

func nullableValue(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
