package models

import (
	"time"
)

// Payment represents a payment
type Payment struct {
	CustomerId string `json:"customer_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	CallbackReceivedAt time.Time `json:"callback_received_at,omitempty"`
	LedgerTransactionId string `json:"ledger_transaction_id,omitempty"`
	BankTranId string `json:"bank_tran_id,omitempty"`
	CardIssuerCountry string `json:"card_issuer_country,omitempty"`
	RiskLevel string `json:"risk_level,omitempty"`
	RiskTitle string `json:"risk_title,omitempty"`
	ManualProofFileId string `json:"manual_proof_file_id,omitempty"`
	TigerbeetleTransferId string `json:"tigerbeetle_transfer_id,omitempty"`
	Status interface{} `json:"status"`
	GatewayResponse string `json:"gateway_response,omitempty"`
	RetryCount int `json:"retry_count"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	ProviderReference string `json:"provider_reference,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
	CardType string `json:"card_type,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	InitiatedAt time.Time `json:"initiated_at"`
	ReceiptUrl string `json:"receipt_url,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	ReceiptDocumentId string `json:"receipt_document_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	TranId string `json:"tran_id,omitempty"`
	IpnReceivedAt time.Time `json:"ipn_received_at,omitempty"`
	Method *PaymentPaymentMethod `json:"method"`
	Amount *Money `json:"amount"`
	PayerId string `json:"payer_id"`
	FailureReason string `json:"failure_reason,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	CardIssuer string `json:"card_issuer,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	PayeeId string `json:"payee_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	CardBrand string `json:"card_brand,omitempty"`
	ManualReviewStatus interface{} `json:"manual_review_status"`
	VerifiedBy string `json:"verified_by,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	ReceiptFileId string `json:"receipt_file_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Type *PaymentType `json:"type"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Gateway string `json:"gateway,omitempty"`
	Provider string `json:"provider,omitempty"`
	ValId string `json:"val_id,omitempty"`
	ValidatedAt time.Time `json:"validated_at,omitempty"`
	PaymentId string `json:"payment_id"`
	Currency string `json:"currency"`
}
